"""Tests for governance rules domain models."""

from datetime import datetime, timezone
from unittest.mock import MagicMock, patch

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.models.governance_rules import (
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    GovernanceRules,
    GovernanceRulesTrail,
    GroupThreshold,
    RuleGroup,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
)


class TestRuleUserSignature:
    """Tests for RuleUserSignature model."""

    def test_create_signature(self) -> None:
        """Test creating a user signature."""
        sig = RuleUserSignature(user_id="user1", signature="base64sig==")

        assert sig.user_id == "user1"
        assert sig.signature == "base64sig=="

    def test_default_values(self) -> None:
        """Test default values."""
        sig = RuleUserSignature()

        assert sig.user_id is None
        assert sig.signature is None


class TestRuleUser:
    """Tests for RuleUser model."""

    def test_create_user(self) -> None:
        """Test creating a rule user."""
        user = RuleUser(
            id="user1",
            name="Admin User",
            public_key_pem="-----BEGIN PUBLIC KEY-----\nABC\n-----END PUBLIC KEY-----",
            roles=["ADMIN", "SUPERADMIN"],
        )

        assert user.id == "user1"
        assert user.name == "Admin User"
        assert user.public_key_pem is not None
        assert user.roles == ["ADMIN", "SUPERADMIN"]

    def test_default_values(self) -> None:
        """Test default values."""
        user = RuleUser()

        assert user.id is None
        assert user.name is None
        assert user.public_key_pem is None
        assert user.roles == []


class TestRuleGroup:
    """Tests for RuleGroup model."""

    def test_create_group(self) -> None:
        """Test creating a rule group."""
        group = RuleGroup(
            id="group1",
            name="Approvers",
            user_ids=["user1", "user2", "user3"],
        )

        assert group.id == "group1"
        assert group.name == "Approvers"
        assert group.user_ids == ["user1", "user2", "user3"]

    def test_default_values(self) -> None:
        """Test default values."""
        group = RuleGroup()

        assert group.id is None
        assert group.name is None
        assert group.user_ids == []


class TestGroupThreshold:
    """Tests for GroupThreshold model."""

    def test_create_threshold(self) -> None:
        """Test creating a group threshold."""
        threshold = GroupThreshold(group_id="group1", threshold=2)

        assert threshold.group_id == "group1"
        assert threshold.threshold == 2

    def test_default_values(self) -> None:
        """Test default values."""
        threshold = GroupThreshold()

        assert threshold.group_id is None
        assert threshold.threshold == 0


class TestSequentialThresholds:
    """Tests for SequentialThresholds model."""

    def test_create_sequential(self) -> None:
        """Test creating sequential thresholds."""
        seq = SequentialThresholds(
            thresholds=[
                GroupThreshold(group_id="g1", threshold=1),
                GroupThreshold(group_id="g2", threshold=2),
            ]
        )

        assert len(seq.thresholds) == 2
        assert seq.thresholds[0].group_id == "g1"
        assert seq.thresholds[1].threshold == 2

    def test_default_values(self) -> None:
        """Test default values."""
        seq = SequentialThresholds()

        assert seq.thresholds == []


class TestAddressWhitelistingRules:
    """Tests for AddressWhitelistingRules model."""

    def test_create_rules(self) -> None:
        """Test creating address whitelisting rules."""
        rules = AddressWhitelistingRules(
            currency="ETH",
            network="mainnet",
            parallel_thresholds=[
                SequentialThresholds(thresholds=[GroupThreshold(group_id="g1", threshold=2)])
            ],
        )

        assert rules.currency == "ETH"
        assert rules.network == "mainnet"
        assert len(rules.parallel_thresholds) == 1
        assert len(rules.lines) == 0

    def test_default_values(self) -> None:
        """Test default values."""
        rules = AddressWhitelistingRules()

        assert rules.currency is None
        assert rules.network is None
        assert rules.parallel_thresholds == []
        assert rules.lines == []


class TestContractAddressWhitelistingRules:
    """Tests for ContractAddressWhitelistingRules model."""

    def test_create_rules(self) -> None:
        """Test creating contract address whitelisting rules."""
        rules = ContractAddressWhitelistingRules(
            blockchain="ETH",
            network="mainnet",
            parallel_thresholds=[
                SequentialThresholds(thresholds=[GroupThreshold(group_id="g1", threshold=1)])
            ],
        )

        assert rules.blockchain == "ETH"
        assert rules.network == "mainnet"
        assert len(rules.parallel_thresholds) == 1
        assert len(rules.parallel_thresholds[0].thresholds) == 1
        assert rules.parallel_thresholds[0].thresholds[0].group_id == "g1"


class TestDecodedRulesContainer:
    """Tests for DecodedRulesContainer model."""

    def test_create_container(self) -> None:
        """Test creating a decoded rules container."""
        container = DecodedRulesContainer(
            users=[RuleUser(id="user1", name="User 1")],
            groups=[RuleGroup(id="group1", name="Group 1")],
            minimum_distinct_user_signatures=2,
            minimum_distinct_group_signatures=1,
            timestamp=1704067200,
        )

        assert len(container.users) == 1
        assert len(container.groups) == 1
        assert container.minimum_distinct_user_signatures == 2
        assert container.minimum_distinct_group_signatures == 1
        assert container.timestamp == 1704067200

    def test_default_values(self) -> None:
        """Test default values."""
        container = DecodedRulesContainer()

        assert container.users == []
        assert container.groups == []
        assert container.minimum_distinct_user_signatures == 0
        assert container.minimum_distinct_group_signatures == 0
        assert container.address_whitelisting_rules == []
        assert container.contract_address_whitelisting_rules == []
        assert container.enforced_rules_hash is None
        assert container.timestamp == 0
        assert container.minimum_commitment_signatures == 0
        assert container.engine_identities == []

    def test_find_user_by_id(self) -> None:
        """Test finding a user by ID."""
        container = DecodedRulesContainer(
            users=[
                RuleUser(id="user1", name="User 1"),
                RuleUser(id="user2", name="User 2"),
            ]
        )

        found = container.find_user_by_id("user1")
        assert found is not None
        assert found.name == "User 1"

        not_found = container.find_user_by_id("user3")
        assert not_found is None

    def test_find_group_by_id(self) -> None:
        """Test finding a group by ID."""
        container = DecodedRulesContainer(
            groups=[
                RuleGroup(id="group1", name="Group 1"),
                RuleGroup(id="group2", name="Group 2"),
            ]
        )

        found = container.find_group_by_id("group1")
        assert found is not None
        assert found.name == "Group 1"

        not_found = container.find_group_by_id("group3")
        assert not_found is None

    def test_find_address_whitelisting_rules_exact_match(self) -> None:
        """Test finding rules with exact blockchain and network match."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="ETH", network="mainnet"),
                AddressWhitelistingRules(currency="ETH", network="testnet"),
            ]
        )

        found = container.find_address_whitelisting_rules("ETH", "mainnet")
        assert found is not None
        assert found.network == "mainnet"

    def test_find_address_whitelisting_rules_blockchain_only(self) -> None:
        """Test finding rules with blockchain match and wildcard network."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="BTC", network=None),  # Wildcard
            ]
        )

        found = container.find_address_whitelisting_rules("BTC", "mainnet")
        assert found is not None
        assert found.currency == "BTC"

    def test_find_address_whitelisting_rules_global_default(self) -> None:
        """Test finding global default rules."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="Any", network=None),  # Global
            ]
        )

        found = container.find_address_whitelisting_rules("ETH", "mainnet")
        assert found is not None
        assert found.currency == "Any"

    def test_find_address_whitelisting_rules_priority(self) -> None:
        """Test that exact match takes priority over global default."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="Any", network=None),  # Global default
                AddressWhitelistingRules(currency="ETH", network="mainnet"),  # Exact
            ]
        )

        found = container.find_address_whitelisting_rules("ETH", "mainnet")
        assert found is not None
        assert found.currency == "ETH"
        assert found.network == "mainnet"

    def test_find_address_whitelisting_rules_empty_is_wildcard(self) -> None:
        """Test that empty string is treated as wildcard."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="", network=""),  # Empty = wildcard
            ]
        )

        found = container.find_address_whitelisting_rules("SOL", "mainnet")
        assert found is not None

    def test_find_address_whitelisting_rules_not_found(self) -> None:
        """Test when no matching rules are found."""
        container = DecodedRulesContainer(
            address_whitelisting_rules=[
                AddressWhitelistingRules(currency="BTC", network="mainnet"),
            ]
        )

        found = container.find_address_whitelisting_rules("ETH", "mainnet")
        assert found is None

    def test_get_hsm_public_key_found(self, ecdsa_public_key_pem: str) -> None:
        """Test getting HSM public key when user with HSMSLOT role exists."""
        container = DecodedRulesContainer(
            users=[
                RuleUser(id="user1", name="Normal User", roles=["USER"]),
                RuleUser(
                    id="hsm_user",
                    name="HSM Slot",
                    public_key_pem=ecdsa_public_key_pem,
                    roles=["HSMSLOT"],
                ),
            ]
        )

        key = container.get_hsm_public_key()
        assert key is not None
        assert isinstance(key, ec.EllipticCurvePublicKey)

    def test_get_hsm_public_key_not_found(self) -> None:
        """Test getting HSM public key when no HSMSLOT user exists."""
        container = DecodedRulesContainer(
            users=[
                RuleUser(id="user1", name="Normal User", roles=["USER"]),
            ]
        )

        key = container.get_hsm_public_key()
        assert key is None

    def test_get_hsm_public_key_cached(self, ecdsa_public_key_pem: str) -> None:
        """Test that HSM public key is cached after first lookup."""
        container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="hsm_user",
                    name="HSM Slot",
                    public_key_pem=ecdsa_public_key_pem,
                    roles=["HSMSLOT"],
                ),
            ]
        )

        # First call
        key1 = container.get_hsm_public_key()
        # Second call should be cached
        key2 = container.get_hsm_public_key()

        assert key1 is key2  # Same instance

    def test_get_hsm_public_key_invalid_pem(self) -> None:
        """Test that invalid PEM is handled gracefully."""
        container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="hsm_user",
                    name="HSM Slot",
                    public_key_pem="invalid-pem",
                    roles=["HSMSLOT"],
                ),
            ]
        )

        key = container.get_hsm_public_key()
        assert key is None

    def test_is_wildcard_static_method(self) -> None:
        """Test _is_wildcard static method."""
        assert DecodedRulesContainer._is_wildcard(None) is True
        assert DecodedRulesContainer._is_wildcard("") is True
        assert DecodedRulesContainer._is_wildcard("Any") is True
        assert DecodedRulesContainer._is_wildcard("any") is True
        assert DecodedRulesContainer._is_wildcard("ANY") is True
        assert DecodedRulesContainer._is_wildcard("ETH") is False
        assert DecodedRulesContainer._is_wildcard("mainnet") is False


class TestGovernanceRulesTrail:
    """Tests for GovernanceRulesTrail model."""

    def test_create_trail(self) -> None:
        """Test creating an audit trail entry."""
        now = datetime.now(timezone.utc)
        trail = GovernanceRulesTrail(
            user_id="user1",
            action="APPROVE",
            timestamp=now,
        )

        assert trail.user_id == "user1"
        assert trail.action == "APPROVE"
        assert trail.timestamp == now

    def test_default_values(self) -> None:
        """Test default values."""
        trail = GovernanceRulesTrail()

        assert trail.user_id is None
        assert trail.action is None
        assert trail.timestamp is None


class TestGovernanceRules:
    """Tests for GovernanceRules model."""

    def test_create_rules(self) -> None:
        """Test creating governance rules."""
        now = datetime.now(timezone.utc)
        rules = GovernanceRules(
            rules_container="base64encoded==",
            rules_signatures=[
                RuleUserSignature(user_id="admin1", signature="sig1"),
                RuleUserSignature(user_id="admin2", signature="sig2"),
            ],
            locked=True,
            creation_date=now,
            update_date=now,
        )

        assert rules.rules_container == "base64encoded=="
        assert len(rules.rules_signatures) == 2
        assert rules.locked is True
        assert rules.creation_date == now
        assert rules.update_date == now

    def test_default_values(self) -> None:
        """Test default values."""
        rules = GovernanceRules()

        assert rules.rules_container is None
        assert rules.rules_signatures == []
        assert rules.locked is False
        assert rules.creation_date is None
        assert rules.update_date is None
        assert rules.trails == []
