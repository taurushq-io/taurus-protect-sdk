"""Tests for rules container decoding from mappers/governance_rules.py."""

from __future__ import annotations

import base64
import json
from pathlib import Path
from typing import Any, Dict

import pytest

from taurus_protect.errors import IntegrityError
from taurus_protect.mappers.governance_rules import (
    rules_container_from_base64,
    user_signatures_from_base64,
)


def load_fixture() -> Dict[str, Any]:
    """Load the test fixture."""
    fixture_path = Path(__file__).parent.parent / "fixtures" / "whitelisted_address_raw_response.json"
    with open(fixture_path) as f:
        return json.load(f)


def get_rules_container_base64() -> str:
    """Get base64-encoded rulesContainerJson from fixture."""
    fixture = load_fixture()
    # The fixture has rulesContainerJson as a dict, encode it as JSON then base64
    rules_json = json.dumps(fixture["rulesContainerJson"])
    return base64.b64encode(rules_json.encode("utf-8")).decode("ascii")


def get_rules_signatures_base64() -> str:
    """Get rulesSignatures from fixture (already base64 protobuf)."""
    fixture = load_fixture()
    return fixture["rulesSignatures"]


# =============================================================================
# Category A: Rules Container Decoding Tests (9 tests)
# =============================================================================


class TestRulesContainerFromBase64:
    """Tests for rules_container_from_base64 function."""

    def test_decode_rules_container_from_base64_success(self) -> None:
        """Test successful decoding of base64 rules container."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        assert decoded is not None
        assert len(decoded.users) > 0
        assert len(decoded.groups) > 0

    def test_decoded_container_has_users(self) -> None:
        """Test decoded container has 4 users."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        assert len(decoded.users) == 4
        user_ids = [u.id for u in decoded.users]
        assert "superadmin1@bank.com" in user_ids
        assert "superadmin2@bank.com" in user_ids
        assert "team1@bank.com" in user_ids
        assert "hsmslot@bank.com" in user_ids

    def test_decoded_container_has_groups(self) -> None:
        """Test decoded container has 2 groups."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        assert len(decoded.groups) == 2
        group_ids = [g.id for g in decoded.groups]
        assert "team1" in group_ids
        assert "superadmins" in group_ids

        # Verify group membership
        team1_group = next(g for g in decoded.groups if g.id == "team1")
        assert "team1@bank.com" in team1_group.user_ids

        superadmins_group = next(g for g in decoded.groups if g.id == "superadmins")
        assert "superadmin1@bank.com" in superadmins_group.user_ids
        assert "superadmin2@bank.com" in superadmins_group.user_ids

    def test_decoded_users_have_pem_keys(self) -> None:
        """Test decoded users have public_key_pem populated."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        for user in decoded.users:
            assert user.public_key_pem is not None, f"User {user.id} missing public_key_pem"
            assert user.public_key_pem.startswith("-----BEGIN PUBLIC KEY-----")
            assert "-----END PUBLIC KEY-----" in user.public_key_pem

    def test_decoded_users_have_roles(self) -> None:
        """Test decoded users have roles array populated."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        for user in decoded.users:
            assert isinstance(user.roles, list), f"User {user.id} roles is not a list"
            assert len(user.roles) > 0, f"User {user.id} has empty roles"

    def test_find_superadmin_users(self) -> None:
        """Test finding 2 users with SUPERADMIN role."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        superadmin_users = [u for u in decoded.users if "SUPERADMIN" in u.roles]
        assert len(superadmin_users) == 2

        superadmin_ids = [u.id for u in superadmin_users]
        assert "superadmin1@bank.com" in superadmin_ids
        assert "superadmin2@bank.com" in superadmin_ids

    def test_find_hsmslot_user(self) -> None:
        """Test finding 1 user with HSMSLOT role."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        hsmslot_users = [u for u in decoded.users if "HSMSLOT" in u.roles]
        assert len(hsmslot_users) == 1
        assert hsmslot_users[0].id == "hsmslot@bank.com"

    def test_decoded_has_address_whitelisting_rules(self) -> None:
        """Test decoded container has 1 ALGO/mainnet address whitelisting rule."""
        base64_data = get_rules_container_base64()
        decoded = rules_container_from_base64(base64_data)

        assert len(decoded.address_whitelisting_rules) == 1
        rule = decoded.address_whitelisting_rules[0]
        assert rule.currency == "ALGO"
        assert rule.network == "mainnet"

        # Check parallel thresholds (now List[SequentialThresholds])
        assert len(rule.parallel_thresholds) == 1
        seq_threshold = rule.parallel_thresholds[0]
        assert len(seq_threshold.thresholds) == 1
        threshold = seq_threshold.thresholds[0]
        assert threshold.group_id == "team1"
        assert threshold.minimum_signatures == 1

    def test_invalid_base64_raises_error(self) -> None:
        """Test invalid base64 raises IntegrityError."""
        with pytest.raises(IntegrityError) as exc_info:
            rules_container_from_base64("not-valid-base64!!!")

        assert "Failed to decode rules container" in str(exc_info.value)


# =============================================================================
# Category B: User Signatures Decoding Tests (5 tests)
# =============================================================================


class TestUserSignaturesFromBase64:
    """Tests for user_signatures_from_base64 function."""

    def test_decode_user_signatures_from_base64_success(self) -> None:
        """Test successful decoding of base64 user signatures."""
        base64_data = get_rules_signatures_base64()
        signatures = user_signatures_from_base64(base64_data)

        assert signatures is not None
        assert len(signatures) > 0

    def test_signatures_contain_user_ids(self) -> None:
        """Test signatures contain expected user IDs."""
        base64_data = get_rules_signatures_base64()
        signatures = user_signatures_from_base64(base64_data)

        user_ids = [sig.user_id for sig in signatures]
        assert "superadmin1@bank.com" in user_ids
        assert "superadmin2@bank.com" in user_ids

    def test_signatures_contain_signature_bytes(self) -> None:
        """Test each signature has non-empty signature field."""
        base64_data = get_rules_signatures_base64()
        signatures = user_signatures_from_base64(base64_data)

        for sig in signatures:
            assert sig.signature is not None, f"Signature for {sig.user_id} is None"
            assert len(sig.signature) > 0, f"Signature for {sig.user_id} is empty"

    def test_signatures_count_matches_expected(self) -> None:
        """Test exactly 2 signatures are present."""
        base64_data = get_rules_signatures_base64()
        signatures = user_signatures_from_base64(base64_data)

        assert len(signatures) == 2

    def test_invalid_signatures_base64_raises_error(self) -> None:
        """Test invalid base64 returns empty list (graceful degradation)."""
        # Note: user_signatures_from_base64 returns empty list on error instead of raising
        # This is different from rules_container_from_base64 which raises IntegrityError
        signatures = user_signatures_from_base64("not-valid-base64!!!")
        assert signatures == []


class TestEmptyInputs:
    """Tests for empty/null input handling."""

    def test_empty_rules_container_returns_empty(self) -> None:
        """Test empty string returns empty container."""
        decoded = rules_container_from_base64("")
        assert decoded is not None
        assert len(decoded.users) == 0
        assert len(decoded.groups) == 0

    def test_empty_signatures_returns_empty_list(self) -> None:
        """Test empty string returns empty list."""
        signatures = user_signatures_from_base64("")
        assert signatures == []
