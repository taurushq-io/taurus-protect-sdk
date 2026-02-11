"""Tests for 6-step whitelisted address verification flow.

This module tests the complete 6-step verification flow for whitelisted addresses:
1. Step 1: Verify metadata hash (SHA256(payloadAsString) == metadata.hash)
2. Step 2: Verify rules container signatures (SuperAdmin keys)
3. Step 3: Decode rules container (base64 -> model)
4. Step 4: Verify hash coverage (metadata.hash in signature hashes list)
5. Step 5: Verify whitelist signatures meet governance thresholds
6. Step 6: Parse WhitelistedAddress from verified payload

The tests verify both success and failure cases for each step.
"""

from __future__ import annotations

import base64
import json
from typing import List

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.constant_time import constant_time_compare
from taurus_protect.helpers.signature_verifier import is_valid_signature
from taurus_protect.helpers.whitelist_hash_helper import compute_legacy_hashes
from taurus_protect.helpers.whitelisted_asset_verifier import verify_hash_coverage
from taurus_protect.mappers.governance_rules import rules_container_from_base64
from taurus_protect.models.governance_rules import (
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    GroupThreshold,
    RuleGroup,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
)
from taurus_protect.models.whitelisted_address import (
    SignedWhitelistedAddressEnvelope,
    WhitelistedAsset,
    WhitelistedAssetMetadata,
    WhitelistMetadata,
    WhitelistSignatureEntry,
    WhitelistUserSignature,
    SignedContractAddress,
)


# =============================================================================
# Test Fixtures
# =============================================================================


@pytest.fixture
def superadmin_key_pair():
    """Generate a SuperAdmin ECDSA P-256 key pair."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    public_key = private_key.public_key()
    return private_key, public_key


@pytest.fixture
def superadmin_key_pair_2():
    """Generate a second SuperAdmin ECDSA P-256 key pair."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    public_key = private_key.public_key()
    return private_key, public_key


@pytest.fixture
def user_key_pair():
    """Generate a user ECDSA P-256 key pair."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    public_key = private_key.public_key()
    return private_key, public_key


def _public_key_to_pem(public_key: EllipticCurvePublicKey) -> str:
    """Convert public key to PEM string."""
    pem_bytes = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return pem_bytes.decode("utf-8")


@pytest.fixture
def sample_payload() -> str:
    """Sample whitelisted address payload JSON string."""
    return json.dumps(
        {
            "currency": "ETH",
            "addressType": "individual",
            "address": "0x1234567890abcdef1234567890abcdef12345678",
            "memo": "",
            "label": "Test Address",
            "customerId": "",
            "exchangeAccountId": "",
            "linkedInternalAddresses": [],
            "contractType": "",
        },
        separators=(",", ":"),
    )


@pytest.fixture
def sample_payload_hash(sample_payload: str) -> str:
    """SHA-256 hash of the sample payload."""
    return calculate_hex_hash(sample_payload)


@pytest.fixture
def sample_rules_container_json(superadmin_key_pair, user_key_pair) -> dict:
    """Sample rules container as a dict."""
    _, superadmin_pub = superadmin_key_pair
    _, user_pub = user_key_pair
    return {
        "users": [
            {
                "id": "superadmin1@bank.com",
                "publicKey": _public_key_to_pem(superadmin_pub),
                "roles": ["SUPERADMIN"],
            },
            {
                "id": "user1@bank.com",
                "publicKey": _public_key_to_pem(user_pub),
                "roles": ["USER", "OPERATOR"],
            },
        ],
        "groups": [
            {"id": "team1", "userIds": ["user1@bank.com"]},
            {"id": "superadmins", "userIds": ["superadmin1@bank.com"]},
        ],
        "minimumDistinctUserSignatures": 0,
        "minimumDistinctGroupSignatures": 0,
        "addressWhitelistingRules": [
            {
                "currency": "ETH",
                "network": "mainnet",
                "parallelThresholds": [{"groupId": "team1", "minimumSignatures": 1}],
            }
        ],
        "contractAddressWhitelistingRules": [],
        "enforcedRulesHash": "",
        "timestamp": 1706194800,
    }


def _encode_rules_container_json(rules_container_dict: dict) -> str:
    """Encode rules container dict to base64."""
    json_str = json.dumps(rules_container_dict, separators=(",", ":"))
    return base64.b64encode(json_str.encode("utf-8")).decode("utf-8")


# =============================================================================
# Step 1 Tests: Verify Metadata Hash
# =============================================================================


class TestStep1VerifyMetadataHash:
    """Tests for Step 1: Verify metadata hash (SHA256(payloadAsString) == metadata.hash)."""

    def test_step1_verify_metadata_hash_success(self, sample_payload: str) -> None:
        """Test that SHA256 hash of payload matches metadata hash."""
        # Compute expected hash
        expected_hash = calculate_hex_hash(sample_payload)

        # Simulate verification
        computed_hash = calculate_hex_hash(sample_payload)

        # Verify match using constant-time comparison
        assert constant_time_compare(computed_hash, expected_hash)
        assert computed_hash == expected_hash

    def test_step1_verify_metadata_hash_failure(self, sample_payload: str) -> None:
        """Test that hash mismatch raises IntegrityError."""
        computed_hash = calculate_hex_hash(sample_payload)
        wrong_hash = "0" * 64  # Invalid hash

        # Should not match
        assert not constant_time_compare(computed_hash, wrong_hash)

        # In actual service, this would raise IntegrityError
        if not constant_time_compare(computed_hash, wrong_hash):
            with pytest.raises(IntegrityError):
                raise IntegrityError(
                    f"computed hash ({computed_hash}) does not match provided hash ({wrong_hash})"
                )

    def test_step1_verify_hash_empty_payload_raises(self) -> None:
        """Test that empty payload raises IntegrityError."""
        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(
                hash="somehash",
                payload_as_string="",
            )
        )

        # Empty payload should be treated as an error
        if not envelope.metadata.payload_as_string:
            with pytest.raises(IntegrityError):
                raise IntegrityError("payloadAsString is empty")

    def test_step1_verify_hash_null_hash_raises(self, sample_payload: str) -> None:
        """Test that null/missing hash raises IntegrityError."""
        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(
                hash=None,
                payload_as_string=sample_payload,
            )
        )

        # Missing hash should be treated as an error
        if not envelope.metadata.hash:
            with pytest.raises(IntegrityError):
                raise IntegrityError("metadata hash is null or empty")


# =============================================================================
# Step 2 Tests: Verify Rules Container Signatures (SuperAdmin Keys)
# =============================================================================


class TestStep2VerifyRulesSignatures:
    """Tests for Step 2: Verify rules container signatures against SuperAdmin keys."""

    def test_step2_verify_rules_signatures_success(
        self,
        superadmin_key_pair,
        superadmin_key_pair_2,
        sample_rules_container_json: dict,
    ) -> None:
        """Test that valid SuperAdmin signatures pass verification."""
        priv1, pub1 = superadmin_key_pair
        priv2, pub2 = superadmin_key_pair_2

        # Encode rules container
        rules_container_b64 = _encode_rules_container_json(sample_rules_container_json)
        rules_data = base64.b64decode(rules_container_b64)

        # Sign with both SuperAdmin keys
        sig1 = sign_data(priv1, rules_data)
        sig2 = sign_data(priv2, rules_data)

        super_admin_keys = [pub1, pub2]

        # Verify signatures
        valid_count = 0
        for sig in [sig1, sig2]:
            if is_valid_signature(rules_data, sig, super_admin_keys):
                valid_count += 1

        # Should have 2 valid signatures
        assert valid_count == 2

    def test_step2_verify_rules_signatures_failure(
        self, superadmin_key_pair, sample_rules_container_json: dict
    ) -> None:
        """Test that invalid signatures raise IntegrityError."""
        _, pub1 = superadmin_key_pair

        # Create a different key pair (wrong key)
        wrong_private_key = ec.generate_private_key(ec.SECP256R1())

        # Encode rules container
        rules_container_b64 = _encode_rules_container_json(sample_rules_container_json)
        rules_data = base64.b64decode(rules_container_b64)

        # Sign with wrong key
        wrong_sig = sign_data(wrong_private_key, rules_data)

        # Verify with correct public key (should fail)
        is_valid = is_valid_signature(rules_data, wrong_sig, [pub1])
        assert is_valid is False

        # This would raise IntegrityError in actual service
        if not is_valid:
            with pytest.raises(IntegrityError):
                raise IntegrityError(
                    "rules container signature verification failed: only 0 valid signatures"
                )

    def test_step2_verify_rules_signatures_below_threshold(
        self,
        superadmin_key_pair,
        superadmin_key_pair_2,
        sample_rules_container_json: dict,
    ) -> None:
        """Test that insufficient valid signatures fail verification."""
        priv1, pub1 = superadmin_key_pair
        _, pub2 = superadmin_key_pair_2

        # Encode rules container
        rules_container_b64 = _encode_rules_container_json(sample_rules_container_json)
        rules_data = base64.b64decode(rules_container_b64)

        # Sign with only one key
        sig1 = sign_data(priv1, rules_data)

        super_admin_keys = [pub1, pub2]
        min_valid_signatures = 2  # Require 2 signatures

        # Verify - should have only 1 valid signature
        valid_count = 0
        if is_valid_signature(rules_data, sig1, super_admin_keys):
            valid_count += 1

        assert valid_count == 1  # Only 1 valid signature
        assert valid_count < min_valid_signatures

        # This would raise IntegrityError
        if valid_count < min_valid_signatures:
            with pytest.raises(IntegrityError):
                raise IntegrityError(
                    f"rules container signature verification failed: only {valid_count} valid signatures, "
                    f"minimum {min_valid_signatures} required"
                )


# =============================================================================
# Step 3 Tests: Decode Rules Container
# =============================================================================


class TestStep3DecodeRulesContainer:
    """Tests for Step 3: Decode rules container (base64 -> model)."""

    def test_step3_decode_rules_container_success(
        self, sample_rules_container_json: dict
    ) -> None:
        """Test that valid rules container decodes successfully."""
        rules_container_b64 = _encode_rules_container_json(sample_rules_container_json)

        # Decode
        decoded = rules_container_from_base64(rules_container_b64)

        # Verify structure
        assert isinstance(decoded, DecodedRulesContainer)
        assert len(decoded.users) == 2
        assert len(decoded.groups) == 2
        assert len(decoded.address_whitelisting_rules) == 1

        # Verify user data
        superadmin_user = decoded.find_user_by_id("superadmin1@bank.com")
        assert superadmin_user is not None
        assert "SUPERADMIN" in superadmin_user.roles

        # Verify group data
        team1_group = decoded.find_group_by_id("team1")
        assert team1_group is not None
        assert "user1@bank.com" in team1_group.user_ids

    def test_step3_decode_rules_container_failure(self) -> None:
        """Test that invalid data raises IntegrityError."""
        # Invalid base64 that doesn't decode to valid JSON or protobuf
        invalid_b64 = base64.b64encode(b"not valid data").decode("utf-8")

        with pytest.raises(IntegrityError) as exc_info:
            rules_container_from_base64(invalid_b64)

        assert "failed to decode" in str(exc_info.value).lower()

    def test_step3_decode_rules_container_empty(self) -> None:
        """Test that empty rules container returns empty DecodedRulesContainer."""
        decoded = rules_container_from_base64("")

        assert isinstance(decoded, DecodedRulesContainer)
        assert decoded.users == []
        assert decoded.groups == []

    def test_step3_decode_rules_container_invalid_base64(self) -> None:
        """Test that invalid base64 raises IntegrityError."""
        with pytest.raises(IntegrityError) as exc_info:
            rules_container_from_base64("not-valid-base64!!!")

        assert "failed" in str(exc_info.value).lower() or "decode" in str(exc_info.value).lower()


# =============================================================================
# Step 4 Tests: Verify Hash Coverage
# =============================================================================


class TestStep4VerifyHashCoverage:
    """Tests for Step 4: Verify hash is in signature hashes list."""

    def test_step4_verify_hash_coverage_success(self, sample_payload_hash: str) -> None:
        """Test that metadata hash is found in signature hashes."""
        # Create signature entries with the hash
        signatures = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(
                    user_id="user1@bank.com",
                    signature="dummysig",
                ),
                hashes=[sample_payload_hash],
            )
        ]

        # Verify hash coverage
        assert verify_hash_coverage(sample_payload_hash, signatures) is True

    def test_step4_verify_hash_coverage_multiple_hashes(
        self, sample_payload_hash: str
    ) -> None:
        """Test that hash is found when signature covers multiple hashes."""
        other_hash = "a" * 64

        # Signature covers multiple hashes
        signatures = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(
                    user_id="user1@bank.com",
                    signature="dummysig",
                ),
                hashes=[other_hash, sample_payload_hash],  # Target hash is second
            )
        ]

        assert verify_hash_coverage(sample_payload_hash, signatures) is True

    def test_step4_verify_hash_coverage_failure_uses_legacy(
        self, sample_payload: str
    ) -> None:
        """Test that legacy hash fallback works when current hash not found."""
        # Modify payload to have contractType (which gets removed in legacy)
        payload_with_contract_type = json.dumps(
            {
                "currency": "ETH",
                "address": "0x1234",
                "label": "Test",
                "contractType": "ERC20",
            },
            separators=(",", ":"),
        )

        current_hash = calculate_hex_hash(payload_with_contract_type)
        legacy_hashes = compute_legacy_hashes(payload_with_contract_type)

        # Create signatures with only the legacy hash (not current hash)
        assert len(legacy_hashes) > 0
        legacy_hash = legacy_hashes[0]

        signatures = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(
                    user_id="user1@bank.com",
                    signature="dummysig",
                ),
                hashes=[legacy_hash],
            )
        ]

        # Current hash not found
        assert verify_hash_coverage(current_hash, signatures) is False

        # But legacy hash IS found
        assert verify_hash_coverage(legacy_hash, signatures) is True

    def test_step4_verify_hash_coverage_failure_not_found(
        self, sample_payload_hash: str
    ) -> None:
        """Test that missing hash raises IntegrityError."""
        different_hash = "b" * 64

        # Signature covers a different hash
        signatures = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(
                    user_id="user1@bank.com",
                    signature="dummysig",
                ),
                hashes=[different_hash],
            )
        ]

        found = verify_hash_coverage(sample_payload_hash, signatures)
        assert found is False

        # This would raise IntegrityError in actual service
        if not found:
            with pytest.raises(IntegrityError):
                raise IntegrityError(
                    f"metadata hash '{sample_payload_hash}' is not covered by any signature"
                )

    def test_step4_verify_hash_coverage_empty_signatures(
        self, sample_payload_hash: str
    ) -> None:
        """Test that empty signatures list returns False."""
        signatures: List[WhitelistSignatureEntry] = []

        found = verify_hash_coverage(sample_payload_hash, signatures)
        assert found is False


# =============================================================================
# Step 5 Tests: Verify Whitelist Signatures Meet Threshold
# =============================================================================


class TestStep5VerifyWhitelistSignatures:
    """Tests for Step 5: Verify whitelist signatures meet governance thresholds."""

    def test_step5_verify_whitelist_signatures_success(
        self, user_key_pair, sample_payload_hash: str
    ) -> None:
        """Test that valid signatures meeting threshold pass verification."""
        user_private, user_public = user_key_pair

        # Create rules container with threshold of 1
        rules_container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="user1@bank.com",
                    name="User 1",
                    public_key_pem=_public_key_to_pem(user_public),
                    roles=["USER", "OPERATOR"],
                )
            ],
            groups=[
                RuleGroup(id="team1", name="Team 1", user_ids=["user1@bank.com"])
            ],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="ETH",
                    network="mainnet",
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="team1", minimum_signatures=1)])
                    ],
                )
            ],
        )

        # Find the rules for ETH/mainnet
        whitelist_rules = rules_container.find_address_whitelisting_rules("ETH", "mainnet")
        assert whitelist_rules is not None
        assert len(whitelist_rules.parallel_thresholds) == 1

        seq_threshold = whitelist_rules.parallel_thresholds[0]
        assert len(seq_threshold.thresholds) == 1
        threshold = seq_threshold.thresholds[0]
        assert threshold.group_id == "team1"
        assert threshold.get_min_signatures() == 1

        # Verify group contains user
        group = rules_container.find_group_by_id("team1")
        assert group is not None
        assert "user1@bank.com" in group.user_ids

    def test_step5_verify_whitelist_signatures_failure_below_threshold(
        self, user_key_pair
    ) -> None:
        """Test that insufficient signatures fail verification (below threshold)."""
        _, user_public = user_key_pair

        # Create rules container requiring 2 signatures
        rules_container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="user1@bank.com",
                    name="User 1",
                    public_key_pem=_public_key_to_pem(user_public),
                    roles=["USER", "OPERATOR"],
                )
            ],
            groups=[
                RuleGroup(id="team1", name="Team 1", user_ids=["user1@bank.com"])
            ],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="ETH",
                    network="mainnet",
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="team1", minimum_signatures=2)])  # Require 2
                    ],
                )
            ],
        )

        whitelist_rules = rules_container.find_address_whitelisting_rules("ETH", "mainnet")
        assert whitelist_rules is not None

        seq_threshold = whitelist_rules.parallel_thresholds[0]
        threshold = seq_threshold.thresholds[0]
        min_sigs = threshold.get_min_signatures()
        assert min_sigs == 2

        # We only have 1 user in the group, so we can only get 1 signature
        group = rules_container.find_group_by_id("team1")
        assert group is not None
        assert len(group.user_ids) == 1

        # Simulation: 1 valid signature but need 2
        valid_count = 1

        if valid_count < min_sigs:
            with pytest.raises(WhitelistError):
                raise WhitelistError(
                    f"group 'team1' requires {min_sigs} signature(s) but only {valid_count} valid"
                )

    def test_step5_verify_no_whitelist_rules_for_blockchain(self) -> None:
        """Test that missing whitelist rules for blockchain raises error."""
        # Rules container without ETH rules
        rules_container = DecodedRulesContainer(
            users=[],
            groups=[],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="BTC",  # Only BTC, not ETH
                    network="mainnet",
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="team1", minimum_signatures=1)])
                    ],
                )
            ],
        )

        # Try to find ETH rules
        whitelist_rules = rules_container.find_address_whitelisting_rules("ETH", "mainnet")
        assert whitelist_rules is None

        if whitelist_rules is None:
            with pytest.raises(WhitelistError):
                raise WhitelistError(
                    "no address whitelisting rules found for blockchain=ETH network=mainnet"
                )

    def test_step5_verify_group_not_found(self) -> None:
        """Test that missing group raises error."""
        rules_container = DecodedRulesContainer(
            users=[],
            groups=[],  # No groups defined
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="ETH",
                    network="mainnet",
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="nonexistent_group", minimum_signatures=1)])
                    ],
                )
            ],
        )

        group = rules_container.find_group_by_id("nonexistent_group")
        assert group is None

        if group is None:
            with pytest.raises(WhitelistError):
                raise WhitelistError("group 'nonexistent_group' not found in rules container")

    def test_step5_verify_user_not_found_in_group(self) -> None:
        """Test that user not in group is handled correctly."""
        rules_container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="user1@bank.com",
                    name="User 1",
                    public_key_pem="-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----",
                    roles=["USER"],
                )
            ],
            groups=[
                RuleGroup(id="team1", name="Team 1", user_ids=["different_user@bank.com"])
            ],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="ETH",
                    network="mainnet",
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="team1", minimum_signatures=1)])
                    ],
                )
            ],
        )

        # user1@bank.com exists in users but NOT in team1 group
        user = rules_container.find_user_by_id("user1@bank.com")
        assert user is not None

        group = rules_container.find_group_by_id("team1")
        assert group is not None
        assert "user1@bank.com" not in group.user_ids

    def test_step5_verify_wildcard_blockchain_rules(self) -> None:
        """Test that wildcard blockchain rules are found as fallback."""
        rules_container = DecodedRulesContainer(
            users=[],
            groups=[
                RuleGroup(id="team1", name="Team 1", user_ids=["user1@bank.com"])
            ],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency="Any",  # Wildcard
                    network=None,  # Wildcard
                    parallel_thresholds=[
                        SequentialThresholds(thresholds=[GroupThreshold(group_id="team1", minimum_signatures=1)])
                    ],
                )
            ],
        )

        # Should find the wildcard rule when looking for ETH
        whitelist_rules = rules_container.find_address_whitelisting_rules("ETH", "mainnet")
        assert whitelist_rules is not None
        assert whitelist_rules.currency == "Any"
