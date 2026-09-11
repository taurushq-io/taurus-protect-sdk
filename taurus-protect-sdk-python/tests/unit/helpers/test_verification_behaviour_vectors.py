"""
Cross-SDK behaviour vectors for the verification primitives.

These are the invariants that had drifted apart before: which (blockchain, network)
pair selects the governance rules, and whether a hash is covered by a signature.
Both are pure input -> outcome mappings, so all four SDKs assert them from one file
rather than from four hand-maintained copies -- the arrangement that let Java's
hash-coverage go non-constant-time while its peers did not.

To add a case: append to the shared file, bump the matching count there, and consume
it in all four suites.
"""

import json
from typing import Any
from pathlib import Path

import pytest

from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.whitelist_hash_helper import resolve_rule_key
from taurus_protect.helpers.whitelist_hash_helper import contains_hash
from taurus_protect.helpers.whitelist_hash_helper import verify_hash_coverage
from taurus_protect.models.governance_rules import GovernanceRules, RuleUserSignature
from taurus_protect.services.governance_rule_service import _ruleset_verification_key

# tests/unit/helpers/ -> tests/unit -> tests -> <sdk> -> <repo root>
_VECTORS_PATH = (
    Path(__file__).resolve().parents[4] / "scripts" / "resources"
    / "verification-behaviour-vectors.json"
)


def _load():
    if not _VECTORS_PATH.exists():
        raise AssertionError(
            f"cannot read shared verification behaviour vectors {_VECTORS_PATH}"
        )
    data = json.loads(_VECTORS_PATH.read_text())

    # Counts are asserted so a case added to the shared file without being consumed
    # here fails loudly rather than being silently ignored by this SDK.
    counts = data["counts"]
    for key, expected in counts.items():
        assert len(data[key]) == expected, (
            f"{key}: got {len(data[key])} vectors, file declares {expected}"
        )
    return data


VECTORS = _load()


class _Sig:
    """Minimal stand-in for a whitelist signature entry."""

    def __init__(self, hashes):
        self.hashes = hashes


@pytest.mark.parametrize("case", VECTORS["rule_key"], ids=lambda c: c["description"])
def test_rule_key(case):
    if case["expect"] == "error":
        with pytest.raises(IntegrityError):
            resolve_rule_key(
                case["payload_as_string"], case["dto_blockchain"], case["dto_network"]
            )
        return

    blockchain, network = resolve_rule_key(
        case["payload_as_string"], case["dto_blockchain"], case["dto_network"]
    )
    assert (blockchain, network) == (case["blockchain"], case["network"])


@pytest.mark.parametrize(
    "case", VECTORS["hash_coverage"], ids=lambda c: c["description"]
)
def test_hash_coverage(case):
    signatures = [_Sig(h) for h in case["signatures"]]
    assert verify_hash_coverage(case["hash"], signatures) is case["expect"]


@pytest.mark.parametrize(
    "case", VECTORS["contains_hash"], ids=lambda c: c["description"]
)
def test_contains_hash(case):
    assert contains_hash(case["hashes"], case["hash"]) is case["expect"]


def _ruleset(spec) -> GovernanceRules:
    """Build a GovernanceRules from a memo_key vector side."""
    return GovernanceRules(
        rules_container=spec["rules_container"],
        rules_signatures=[
            RuleUserSignature(user_id=sig["user_id"], signature=sig["signature"])
            for sig in spec["signatures"]
        ],
    )


def test_memo_key_section_can_catch_a_colliding_key():
    """Non-vacuity: without a 'distinct' pair the section would pass against a key
    function that returns a constant."""
    assert any(c["expect"] == "distinct" for c in VECTORS["memo_key"]), (
        "memo_key: no 'distinct' vector, so the section cannot catch a colliding key"
    )


@pytest.mark.parametrize("case", VECTORS["memo_key"], ids=lambda c: c["description"])
def test_memo_key(case):
    """The governance verification memo key decides whether ECDSA runs at all, so it
    MUST be injective. An unprefixed concatenation leaves the boundary between the
    container and the signature list uncommitted, and a response-controlling attacker
    can then shift bytes across it to make a MODIFIED container inherit a genuine one's
    "already verified" status -- skipping signature verification on the document that
    carries every HSM key the SDK trusts.
    """
    key_a = _ruleset_verification_key(_ruleset(case["a"]))
    key_b = _ruleset_verification_key(_ruleset(case["b"]))
    assert key_a is not None and key_b is not None, "both containers must be decodable"

    if case["expect"] == "distinct":
        assert key_a != key_b, (
            "distinct (container, signatures) pairs must not share a memo key"
        )
    elif case["expect"] == "same":
        assert key_a == key_b, (
            "these inputs describe the same verified document, so the key must match"
        )
    else:
        raise AssertionError(f"unknown expect {case['expect']!r}")


# ---------------------------------------------------------------------------------
# rule_tier_candidates -- which rule tiers apply when the network is UNSIGNED
# ---------------------------------------------------------------------------------
#
# Until 2026-09-10 this section and `legacy_hash` were consumed by the Go suite ALONE,
# even though the file's `counts` block declares both and every loader asserts the counts.
# Asserting a section's LENGTH proves the file is well formed; it does not prove the
# behaviour is checked. That asymmetry is how Python came to ship the single-tier lookup
# this section exists to forbid, while three SDKs had the candidate-set fix.


def _rules_container_with(rules: list, family: str) -> Any:
    """A container carrying only the rule tiers a vector describes."""
    from taurus_protect.models.governance_rules import (
        AddressWhitelistingRules,
        ContractAddressWhitelistingRules,
        DecodedRulesContainer,
    )

    if family == "address":
        return DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency=r["blockchain"], network=r["network"])
                for r in rules
            ]
        )
    return DecodedRulesContainer(
        contract_address_whitelisting_rules=[
            ContractAddressWhitelistingRules(
                blockchain=r["blockchain"], network=r["network"]
            )
            for r in rules
        ]
    )


@pytest.mark.parametrize("case", VECTORS["rule_tier_candidates"], ids=lambda c: c["description"])
def test_rule_tier_candidates_for_the_address_family(case: dict) -> None:
    """Every tier an unsigned DTO network could have selected must be returned.

    ``include_network_in_payload`` has no proto backing in any SDK, so when the signed
    payload omits ``network`` there is no authenticated way to learn whether that was
    legitimate. A single lookup on the DTO's value hands a response-controlling server the
    choice of WHICH quorum the row must meet; the answer is to enforce every reachable
    tier instead.
    """
    container = _rules_container_with(case["rules"], "address")

    got = container.find_address_whitelisting_rule_candidates(case["blockchain"])

    assert [(r.currency or "", r.network or "") for r in got] == [
        (e["blockchain"], e["network"]) for e in case["expect_candidates"]
    ], case["description"]


@pytest.mark.parametrize("case", VECTORS["rule_tier_candidates"], ids=lambda c: c["description"])
def test_rule_tier_candidates_for_the_contract_family(case: dict) -> None:
    """The asset peer, driven by the same vectors.

    Worth its own pass rather than trusting the address one: in THIS SDK the two families
    name the chain field differently (``currency`` vs ``blockchain``), so a walk reading
    one name treats every contract rule as a wildcard global default -- fail-OPEN, since
    the global default is the broadest tier.
    """
    container = _rules_container_with(case["rules"], "contract")

    got = container.find_contract_address_whitelisting_rule_candidates(case["blockchain"])

    assert [(r.blockchain or "", r.network or "") for r in got] == [
        (e["blockchain"], e["network"]) for e in case["expect_candidates"]
    ], case["description"]


# ---------------------------------------------------------------------------------
# legacy_hash -- WHICH payload step 6 parses
# ---------------------------------------------------------------------------------
#
# This section asserts the PARSED PAYLOAD, not a hash, which is why it lives here rather
# than in `docs/test-vectors/crypto-test-vectors.json`: the legacy-strip injection moves no
# hash at all, so a hash-value gate is structurally blind to it.


@pytest.mark.parametrize(
    "case", VECTORS["legacy_hash"], ids=lambda c: c["description"]
)
def test_legacy_hash_matches_the_payload_the_signature_covered(case: dict) -> None:
    """Step 4 must carry forward the payload it matched, not the one the server delivered.

    The three legacy strips are not injective, so a server can append a duplicate
    ``,"label":"X"`` (or a ``contractType`` a row never had) immediately before the closing
    brace: the strip recovers the genuinely signed bytes, every signature check passes, and
    a step 6 that parses the DELIVERED text returns the appended value as *verified*.
    """
    from taurus_protect.crypto.hashing import calculate_hex_hash
    from taurus_protect.helpers.whitelist_hash_helper import (
        compute_legacy_payload_variants,
    )

    covered_hash = calculate_hex_hash(case["signed_payload"])

    matched_payload = None
    if calculate_hex_hash(case["delivered_payload"]) == covered_hash:
        matched_payload = case["delivered_payload"]
    else:
        for variant in compute_legacy_payload_variants(case["delivered_payload"]):
            if variant.hash == covered_hash:
                matched_payload = variant.payload
                break

    if case["expect"] == "no_match":
        assert matched_payload is None, (
            f"expected no variant to be covered, but {matched_payload!r} matched"
        )
        return

    assert matched_payload is not None, "expected a covered variant, found none"
    assert matched_payload == case["expect_matched_payload"], case["description"]
