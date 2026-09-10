"""
Cross-SDK signed fixtures for the two verification thresholds.

Step 2 is the tenant-wide SuperAdmin threshold; step 5 is the per-group one. Different
scopes, same counting rule, so one file gates both.

The behaviour vectors cover everything expressible without key material; this file
covers what needs real signatures. Until it existed, the repo's own CLAUDE.md
recorded that NO automated cross-SDK gate covered the distinct-key rule, and the
five cases below lived as five hand-maintained copies in four suites.

The fixture carries PUBLIC keys and signatures only -- the gate verifies, it never
signs -- so there is no private key material in the repo.
"""

import base64
import json
from pathlib import Path

import pytest

from taurus_protect.crypto.keys import decode_public_key_pem
from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.signature_verifier import verify_governance_rules_signatures
from taurus_protect.helpers.whitelisted_address_verifier import WhitelistedAddressVerifier
from taurus_protect.models.governance_rules import (
    DecodedRulesContainer,
    GroupThreshold,
    RuleGroup,
    RuleUser,
    RuleUserSignature,
)
from taurus_protect.models.whitelisted_address import (
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)

# tests/unit/helpers -> tests/unit -> tests -> <sdk> -> <repo root>
_FIXTURES_PATH = (
    Path(__file__).resolve().parents[4] / "scripts" / "resources"
    / "verification-signed-fixtures.json"
)


def _load():
    if not _FIXTURES_PATH.exists():
        raise AssertionError(f"cannot read shared signed fixtures {_FIXTURES_PATH}")
    data = json.loads(_FIXTURES_PATH.read_text())

    counts = data["counts"]
    for key, expected in counts.items():
        assert len(data[key]) == expected, (
            f"{key}: got {len(data[key])} vectors, file declares {expected}"
        )
    assert data["rules_container_base64"], "fixture carries no rules container"
    return data


FIXTURES = _load()
CONTAINER = base64.b64decode(FIXTURES["rules_container_base64"])


@pytest.mark.parametrize(
    "case", FIXTURES["superadmin_threshold"], ids=lambda c: c["description"]
)
def test_superadmin_threshold(case):
    keys = [decode_public_key_pem(pem) for pem in case["super_admin_keys_pem"]]
    signatures = [
        RuleUserSignature(user_id=s["user_id"], signature=s["signature"])
        for s in case["signatures"]
    ]

    if case["expect"] == "error":
        with pytest.raises((IntegrityError, ValueError)):
            verify_governance_rules_signatures(
                CONTAINER, signatures, keys, case["min_valid_signatures"]
            )
        return

    verify_governance_rules_signatures(
        CONTAINER, signatures, keys, case["min_valid_signatures"]
    )


def test_fixture_covers_the_distinct_key_rule():
    """
    Without a case where the entry COUNT meets the threshold but the distinct-key
    count does not, the whole file would pass against an implementation that counts
    entries -- which is the regression these fixtures exist to catch.
    """
    assert any(
        c["expect"] == "error"
        and len(c["signatures"]) >= c["min_valid_signatures"] > 0
        for c in FIXTURES["superadmin_threshold"]
    ), "no vector distinguishes distinct-key counting from entry counting"


@pytest.mark.parametrize(
    "case", FIXTURES["group_threshold"], ids=lambda c: c["description"]
)
def test_group_threshold(case):
    """
    Step 5 -- the per-group threshold. A different threshold from the SuperAdmin one
    (per group, not tenant-wide) but the same counting rule, so it is gated from the
    same file: the four SDKs cannot drift on one without drifting on the other.
    """
    container = DecodedRulesContainer(
        users=[
            RuleUser(id=u["user_id"], public_key_pem=u["public_key_pem"], roles=["USER"])
            for u in case["users"]
        ],
        groups=[RuleGroup(id=case["group_id"], user_ids=case["group_user_ids"])],
    )
    signatures = [
        WhitelistSignatureEntry(
            hashes=s["hashes"],
            user_signature=WhitelistUserSignature(
                user_id=s["user_id"], signature=s["signature"]
            ),
        )
        for s in case["signatures"]
    ]
    precomputed = {
        idx: json.dumps(sig.hashes, separators=(",", ":")).encode("utf-8")
        for idx, sig in enumerate(signatures)
    }

    verifier = WhitelistedAddressVerifier([], 0)
    threshold = GroupThreshold(
        group_id=case["group_id"], minimum_signatures=case["minimum_signatures"]
    )
    failure = verifier._verify_group_threshold(
        threshold, container, signatures, case["metadata_hash"], precomputed
    )

    if case["expect"] == "error":
        assert failure is not None, "expected the threshold to fail, it passed"
    else:
        assert failure is None, f"expected the threshold to pass, got: {failure}"


def test_group_section_covers_the_distinct_signer_rule():
    """Same non-vacuity guard as the SuperAdmin section, applied per group."""
    found = False
    for c in FIXTURES["group_threshold"]:
        if c["expect"] != "error" or c["minimum_signatures"] <= 0:
            continue
        members = set(c["group_user_ids"])
        in_group = sum(1 for s in c["signatures"] if s["user_id"] in members)
        if in_group >= c["minimum_signatures"]:
            found = True
            break
    assert found, (
        "no group vector where in-group entry count meets the threshold but "
        "distinct signers do not; an implementation counting entries would pass"
    )
