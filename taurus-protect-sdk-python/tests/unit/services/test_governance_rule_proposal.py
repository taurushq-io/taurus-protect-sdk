"""Unit tests for the governance proposal lifecycle:
update_rules_proposal / approve_rules_proposal / reject_rules_proposal.
"""

import base64
from unittest.mock import MagicMock

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import APIError, IntegrityError
from taurus_protect.mappers.rules_container_encode import rules_container_to_base64
from taurus_protect.models.governance_rules import (
    DecodedRulesContainer,
    GovernanceRules,
    RuleColumn,
    TransactionRules,
)
from taurus_protect.services.governance_rule_service import GovernanceRuleService


def _service():
    return GovernanceRuleService(
        api_client=MagicMock(),
        governance_rules_api=MagicMock(),
        super_admin_keys=[],
        min_valid_signatures=0,
    )


def _pin_for(container_b64: str) -> str:
    """The pin an approver passes back: sha256 of the DECODED container bytes."""
    import hashlib

    return hashlib.sha256(base64.b64decode(container_b64)).hexdigest()


def _sample_container() -> DecodedRulesContainer:
    return DecodedRulesContainer(
        minimum_distinct_user_signatures=1,
        transaction_rules=[TransactionRules(key="ETH/transfer", columns=[RuleColumn(type="RuleFiatAmount")])],
    )


def test_update_rules_proposal_submits_encoded_container():
    svc = _service()
    container = _sample_container()

    svc.update_rules_proposal(container)

    svc._api.rule_service_update_rules_proposal.assert_called_once()
    body = svc._api.rule_service_update_rules_proposal.call_args.kwargs["body"]
    assert body.rules_container == rules_container_to_base64(container)


def test_update_rules_proposal_none_raises():
    svc = _service()
    with pytest.raises(ValueError):
        svc.update_rules_proposal(None)
    svc._api.rule_service_update_rules_proposal.assert_not_called()


def test_approve_rules_proposal_signs_pending_container():
    private_key = ec.generate_private_key(ec.SECP256R1())
    public_key = private_key.public_key()
    pending_b64 = rules_container_to_base64(_sample_container())

    svc = _service()
    svc.get_rules_proposal = lambda: GovernanceRules(rules_container=pending_b64)  # type: ignore[method-assign]

    svc.approve_rules_proposal(private_key, "looks good", _pin_for(pending_b64))

    svc._api.rule_service_approve_rules_proposal.assert_called_once()
    body = svc._api.rule_service_approve_rules_proposal.call_args.kwargs["body"]
    assert body.comment == "looks good"
    # signature verifies against the decoded pending bytes; ECDSA is non-deterministic
    # so we verify rather than compare.
    data = base64.b64decode(pending_b64)
    assert verify_signature(public_key, data, body.signature) is True
    assert verify_signature(public_key, data + b"\x01", body.signature) is False


def test_approve_rules_proposal_no_pending_raises():
    private_key = ec.generate_private_key(ec.SECP256R1())
    svc = _service()
    svc.get_rules_proposal = lambda: None  # type: ignore[method-assign]
    with pytest.raises(ValueError):
        svc.approve_rules_proposal(private_key, "x", "a" * 64)
    svc._api.rule_service_approve_rules_proposal.assert_not_called()


def test_approve_rules_proposal_none_key_raises():
    svc = _service()
    with pytest.raises(ValueError):
        svc.approve_rules_proposal(None, "x", "a" * 64)


def test_reject_rules_proposal_submits_comment():
    svc = _service()
    svc.reject_rules_proposal("not compliant")
    svc._api.rule_service_reject_rules_proposal.assert_called_once()
    body = svc._api.rule_service_reject_rules_proposal.call_args.kwargs["body"]
    assert body.comment == "not compliant"


def test_approve_rules_proposal_requires_an_explicit_comment():
    """comment is positional in Go, Java and TypeScript. It used to default to "" here,
    so a Python caller could approve a governance change with no rationale recorded
    while the same call in any other SDK would not compile."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    svc = _service()
    svc.get_rules_proposal = lambda: GovernanceRules(rules_container=rules_container_to_base64(_sample_container()))  # type: ignore[method-assign]

    with pytest.raises(TypeError):
        svc.approve_rules_proposal(private_key)  # type: ignore[call-arg]

    with pytest.raises(TypeError):
        svc.reject_rules_proposal()  # type: ignore[call-arg]

    # An empty comment is still accepted when passed deliberately.
    pending = rules_container_to_base64(_sample_container())
    svc.approve_rules_proposal(private_key, "", _pin_for(pending))
    body = svc._api.rule_service_approve_rules_proposal.call_args.kwargs["body"]
    assert body.comment == ""


def test_update_rules_proposal_maps_api_error():
    svc = _service()
    svc._api.rule_service_update_rules_proposal.side_effect = ApiException(status=500, reason="boom")
    with pytest.raises(APIError):
        svc.update_rules_proposal(_sample_container())


def test_approve_rules_proposal_maps_api_error():
    private_key = ec.generate_private_key(ec.SECP256R1())
    svc = _service()
    svc.get_rules_proposal = lambda: GovernanceRules(rules_container=rules_container_to_base64(_sample_container()))  # type: ignore[method-assign]
    svc._api.rule_service_approve_rules_proposal.side_effect = ApiException(status=500, reason="boom")
    with pytest.raises(APIError):
        svc.approve_rules_proposal(
            private_key, "x", _pin_for(rules_container_to_base64(_sample_container()))
        )


def test_reject_rules_proposal_maps_api_error():
    svc = _service()
    svc._api.rule_service_reject_rules_proposal.side_effect = ApiException(status=500, reason="boom")
    with pytest.raises(APIError):
        svc.reject_rules_proposal("x")


def test_approve_rules_proposal_refuses_when_pending_container_changed():
    """The content-binding gate.

    A server able to shape responses can serve the benign proposal to the review call and
    a DIFFERENT container to the re-fetch inside approve_rules_proposal. Without the pin
    that yields a genuine SuperAdmin signature over attacker-chosen bytes, which then
    verifies clean everywhere -- including in verifiers that never trusted that server.
    The approval must abort and submit nothing.
    """
    private_key = ec.generate_private_key(ec.SECP256R1())
    reviewed_b64 = rules_container_to_base64(_sample_container())
    reviewed = base64.b64decode(reviewed_b64)
    # The re-fetch serves the reviewed bytes with a field appended -- the shape a
    # merge-on-concatenation protobuf attack produces.
    substituted = base64.b64encode(reviewed + b"\x40\x01").decode("ascii")

    svc = _service()
    svc.get_rules_proposal = lambda: GovernanceRules(rules_container=substituted)  # type: ignore[method-assign]

    with pytest.raises(IntegrityError, match="changed since it was reviewed"):
        svc.approve_rules_proposal(private_key, "lgtm", _pin_for(reviewed_b64))
    svc._api.rule_service_approve_rules_proposal.assert_not_called()


def test_approve_rules_proposal_requires_a_pin():
    """The pin is mandatory: an empty one would restore the unpinned behaviour silently."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    svc = _service()
    called = []
    svc.get_rules_proposal = lambda: called.append(1)  # type: ignore[method-assign]

    with pytest.raises(ValueError, match="expected_container_hash"):
        svc.approve_rules_proposal(private_key, "lgtm", "")
    assert not called, "the pin must be validated before any fetch"
    svc._api.rule_service_approve_rules_proposal.assert_not_called()
