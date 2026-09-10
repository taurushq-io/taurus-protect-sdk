"""
Governance reads must verify SuperAdmin signatures.

The governance container carries the HSM public key that address verification trusts,
so an unverified one is the worst object this SDK can hand back. ``get_rules_history``
returned it unverified and nothing tested that.
"""

from __future__ import annotations

from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.errors import IntegrityError
from taurus_protect.mappers import rules_container_to_base64
from taurus_protect.models.governance_rules import DecodedRulesContainer
from taurus_protect.services.governance_rule_service import (
    GovernanceRuleService,
    _ruleset_verification_key,
)


def wire_valid_unsigned_container() -> str:
    """
    A container that DECODES cleanly but carries no signatures. It has to decode: a
    malformed blob raises from the decoder instead, and the test then passes whether
    verification ran or not. This trap has bitten this repo twice.
    """
    return rules_container_to_base64(DecodedRulesContainer(minimum_distinct_user_signatures=1))


def test_wire_valid_container_decodes_cleanly() -> None:
    """The guard that keeps the tests below non-vacuous."""
    from taurus_protect.mappers import rules_container_from_base64

    assert rules_container_from_base64(wire_valid_unsigned_container()) is not None


def _service() -> tuple:
    api = MagicMock()
    service = GovernanceRuleService(
        api_client=MagicMock(),
        governance_rules_api=api,
        super_admin_keys=[ec.generate_private_key(ec.SECP256R1()).public_key()],
        min_valid_signatures=1,
    )
    return service, api


def _rules_dto(container: str) -> MagicMock:
    dto = MagicMock()
    dto.rules_container = container
    dto.rules_signatures = []
    dto.locked = False
    return dto


@pytest.mark.parametrize("read", ["get_rules", "get_rules_by_id"])
def test_single_ruleset_reads_refuse_unsigned_container(read: str) -> None:
    service, api = _service()
    reply = MagicMock()
    reply.result = _rules_dto(wire_valid_unsigned_container())
    api.rule_service_get_rules.return_value = reply
    api.rule_service_get_rules_by_id.return_value = reply

    with pytest.raises(IntegrityError):
        service.get_rules() if read == "get_rules" else service.get_rules_by_id("1")


def test_history_excludes_and_names_unverified_entries() -> None:
    """
    History is LENIENT where the single reads are strict, deliberately: a SuperAdmin key
    rotation makes every pre-rotation ruleset unverifiable, so a strict page would deny
    the whole audit trail from the rotation onwards. It must still NAME what it dropped.
    """
    service, api = _service()
    container = wire_valid_unsigned_container()

    dto = _rules_dto(container)
    dto.creation_date = datetime(2026, 1, 1, tzinfo=timezone.utc)
    reply = MagicMock()
    reply.result = [dto, dto]
    reply.cursor = None
    reply.total_items = "2"
    api.rule_service_get_rules_history.return_value = reply

    result = service.get_rules_history(page_size=10)

    assert result.rules == [], "no entry verified, so none may be returned"
    assert len(result.excluded_unverified) == 2, "both entries must be named as excluded"
    assert all(ex.reason for ex in result.excluded_unverified)
    # The server counted 2; the caller can read 0.
    assert result.total_items == "0"


def test_memo_keys_on_the_signature_set() -> None:
    """
    The memo must key on the signature set, not just the container bytes: otherwise
    re-signing the same container would be served from a stale hit.
    """
    from taurus_protect.models.governance_rules import GovernanceRules, RuleUserSignature

    container = wire_valid_unsigned_container()
    base = GovernanceRules(rules_container=container)
    resigned = GovernanceRules(
        rules_container=container,
        rules_signatures=[RuleUserSignature(user_id="u1", signature="sig")],
    )
    assert _ruleset_verification_key(base) != _ruleset_verification_key(resigned)

    # Signature order is server-controlled, so it must not change the key.
    a = GovernanceRules(
        rules_container=container,
        rules_signatures=[
            RuleUserSignature(user_id="u1", signature="s1"),
            RuleUserSignature(user_id="u2", signature="s2"),
        ],
    )
    b = GovernanceRules(
        rules_container=container,
        rules_signatures=[
            RuleUserSignature(user_id="u2", signature="s2"),
            RuleUserSignature(user_id="u1", signature="s1"),
        ],
    )
    assert _ruleset_verification_key(a) == _ruleset_verification_key(b)


def test_memo_ignores_user_id() -> None:
    """
    user_id is not read by verification, so it must not participate in the memo key --
    otherwise a field the server controls freely and verification ignores can influence a
    verification-SKIP decision.
    """
    from taurus_protect.models.governance_rules import GovernanceRules, RuleUserSignature

    container = wire_valid_unsigned_container()
    a = GovernanceRules(
        rules_container=container,
        rules_signatures=[RuleUserSignature(user_id="u1", signature="s1")],
    )
    b = GovernanceRules(
        rules_container=container,
        rules_signatures=[RuleUserSignature(user_id="attacker-chosen", signature="s1")],
    )
    assert _ruleset_verification_key(a) == _ruleset_verification_key(b)


def test_memo_key_is_injective() -> None:
    """
    The memo key must be INJECTIVE. A plain concatenation leaves the boundary between the
    container and the signature list uncommitted, so a response-controlling attacker can
    shift bytes across it and make a MODIFIED container collide with a genuine one's key
    -- inheriting its "already verified" status and skipping ECDSA entirely.

    Each pair below collides under an unprefixed key and must not collide here. Reverting
    _ruleset_verification_key to a concatenation takes this test red.
    """
    import base64

    from taurus_protect.models.governance_rules import GovernanceRules, RuleUserSignature

    container = wire_valid_unsigned_container()
    raw = base64.b64decode(container, validate=True)

    # Fold the signature list into the container: genuine (C, [sig]) vs (C||sig, []).
    # Three tail shapes, one per plausible separator-free encoding, so this test is red
    # against ALL of them rather than only the one that happens to be in place:
    #   b"s1"         -> plain concatenation
    #   b"\x00s1"     -> a 0x00 written before each entry
    #   b"\x00\x00s1" -> 0x00 separator plus an empty user_id field
    signed = GovernanceRules(
        rules_container=container,
        rules_signatures=[RuleUserSignature(signature="s1")],
    )
    for tail in (b"s1", b"\x00s1", b"\x00\x00s1"):
        folded = base64.b64encode(raw + tail).decode("ascii")
        assert _ruleset_verification_key(signed) != _ruleset_verification_key(
            GovernanceRules(rules_container=folded)
        ), f"collision with container tail {tail!r}"

    # One signature split into two must not collide with the concatenation.
    one = GovernanceRules(
        rules_container=container,
        rules_signatures=[RuleUserSignature(signature="abcd")],
    )
    two = GovernanceRules(
        rules_container=container,
        rules_signatures=[
            RuleUserSignature(signature="ab"),
            RuleUserSignature(signature="cd"),
        ],
    )
    assert _ruleset_verification_key(one) != _ruleset_verification_key(two)


def test_memo_skips_undecodable_container() -> None:
    """
    An undecodable container has no stable identity, so it must not be memoised -- and
    verification must still get the final say.
    """
    from taurus_protect.models.governance_rules import GovernanceRules

    assert _ruleset_verification_key(GovernanceRules(rules_container="!!!not base64!!!")) is None


def test_container_decode_rejects_out_of_alphabet_bytes() -> None:
    """
    Lenient base64 decoding silently discards out-of-alphabet characters and absorbs the
    rest, so a container carrying embedded separators would decode to the genuine bytes
    with ATTACKER bytes appended. Protobuf treats concatenation as merge, so that adds
    entries to repeated fields such as `users` -- an attacker-chosen HSM or PRICEUPDATER
    key in the trust root. Strict decoding is what refuses it.
    """
    import base64

    from taurus_protect.mappers import rules_container_from_base64

    container = wire_valid_unsigned_container()
    raw = base64.b64decode(container, validate=True)
    appended = base64.b64encode(b"appended").decode("ascii")
    smuggled = base64.b64encode(raw).decode("ascii") + "\x00" + appended

    with pytest.raises(IntegrityError):
        rules_container_from_base64(smuggled)


def test_memo_does_not_cache_failures() -> None:
    """A failure must resurface on every call rather than being remembered."""
    service, api = _service()
    reply = MagicMock()
    reply.result = _rules_dto(wire_valid_unsigned_container())
    api.rule_service_get_rules.return_value = reply

    with pytest.raises(IntegrityError):
        service.get_rules()
    with pytest.raises(IntegrityError):
        service.get_rules()
