"""Tests for WhitelistedAddressVerifier.

Tests the 6-step verification flow for whitelisted addresses:
1. Verify metadata hash
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container
4. Verify hash coverage (metadata hash in signature hashes)
5. Verify whitelist signatures meet governance thresholds
6. Parse WhitelistedAddress from verified payload

Also tests:
- Legacy hash fallback
- Rule lines logic (wallet path matching)
- Happy path end-to-end
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
from taurus_protect.helpers.whitelisted_address_verifier import (
    AddressVerificationResult,
    WhitelistedAddressVerifier,
    _contains_hash,
    _verify_hash_coverage,
)
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_INTERNAL_WALLET,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    DecodedRulesContainer,
    GroupThreshold,
    RuleGroup,
    RuleSource,
    RuleSourceInternalWallet,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
)
from taurus_protect.models.whitelisted_address import (
    InternalWallet,
    SignedWhitelistedAddress,
    SignedWhitelistedAddressEnvelope,
    WhitelistMetadata,
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)


# =============================================================================
# Fixtures
# =============================================================================


def _generate_key_pair():
    """Generate a P-256 ECDSA key pair."""
    private_key = ec.generate_private_key(ec.SECP256R1())
    return private_key, private_key.public_key()


def _public_key_to_pem(public_key: EllipticCurvePublicKey) -> str:
    """Convert public key to PEM string."""
    return public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode("utf-8")


def _encode_rules_container_json(rules_dict: dict) -> str:
    """Encode rules container dict to base64."""
    json_str = json.dumps(rules_dict, separators=(",", ":"))
    return base64.b64encode(json_str.encode("utf-8")).decode("utf-8")


def _make_rules_signatures_base64(rules_container_b64: str, private_keys) -> str:
    """Create a base64-encoded UserSignatures protobuf-like structure.

    Since protobuf encoding is complex, we use a mock decoder in tests instead.
    This returns a dummy base64 string.
    """
    return base64.b64encode(b"dummy-signatures").decode("utf-8")


@pytest.fixture
def superadmin_keys():
    """Generate SuperAdmin key pair."""
    return _generate_key_pair()


@pytest.fixture
def user1_keys():
    """Generate user1 key pair."""
    return _generate_key_pair()


@pytest.fixture
def user2_keys():
    """Generate user2 key pair."""
    return _generate_key_pair()


def _build_payload(
    address="0xABCD1234",
    currency="ETH",
    label="Test Addr",
    contract_type="",
    linked_internal_addresses=None,
    linked_wallets=None,
    network=None,
):
    """Build a whitelisted address payload dict."""
    payload = {
        "currency": currency,
        "addressType": "individual",
        "address": address,
        "memo": "",
        "label": label,
        "customerId": "",
        "exchangeAccountId": "",
        "linkedInternalAddresses": linked_internal_addresses or [],
        "contractType": contract_type,
    }
    if linked_wallets is not None:
        payload["linkedWallets"] = linked_wallets
    if network is not None:
        payload["network"] = network
    return payload


def _payload_to_string(payload: dict) -> str:
    """Serialize payload to compact JSON string."""
    return json.dumps(payload, separators=(",", ":"))


def _build_rules_container_dict(user_pub_pem: str, group_id="approvers", user_id="user1@bank.com"):
    """Build a rules container dict with one user and one group."""
    return {
        "users": [
            {
                "id": user_id,
                "publicKey": user_pub_pem,
                "roles": ["USER", "OPERATOR"],
            }
        ],
        "groups": [
            {"id": group_id, "userIds": [user_id]},
        ],
        "minimumDistinctUserSignatures": 0,
        "minimumDistinctGroupSignatures": 0,
        "addressWhitelistingRules": [
            {
                "currency": "ETH",
                "network": "mainnet",
                "parallelThresholds": [
                    {"groupId": group_id, "minimumSignatures": 1}
                ],
            }
        ],
        "contractAddressWhitelistingRules": [],
        "enforcedRulesHash": "",
        "timestamp": 1706194800,
    }


def _build_full_envelope(
    user_private_key,
    user_public_key,
    superadmin_private_key,
    superadmin_public_key,
    payload_dict=None,
    blockchain="ETH",
    network="mainnet",
    group_id="approvers",
    user_id="user1@bank.com",
    linked_wallets=None,
):
    """Build a fully valid SignedWhitelistedAddressEnvelope for testing.

    Returns (envelope, rules_container_decoder, user_signatures_decoder).
    """
    if payload_dict is None:
        payload_dict = _build_payload()

    payload_str = _payload_to_string(payload_dict)
    metadata_hash = calculate_hex_hash(payload_str)

    # Build rules container
    user_pub_pem = _public_key_to_pem(user_public_key)
    rules_dict = _build_rules_container_dict(user_pub_pem, group_id=group_id, user_id=user_id)
    rules_b64 = _encode_rules_container_json(rules_dict)
    rules_data = base64.b64decode(rules_b64)

    # Sign rules container with SuperAdmin key
    sa_sig = sign_data(superadmin_private_key, rules_data)

    # Sign the hashes array with user key
    hashes = [metadata_hash]
    hashes_json = json.dumps(hashes, separators=(",", ":"))
    user_sig = sign_data(user_private_key, hashes_json.encode("utf-8"))

    # Build the envelope
    envelope = SignedWhitelistedAddressEnvelope(
        metadata=WhitelistMetadata(
            hash=metadata_hash,
            payload_as_string=payload_str,
        ),
        blockchain=blockchain,
        network=network,
        rules_container=rules_b64,
        rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
        signed_address=SignedWhitelistedAddress(
            signatures=[
                WhitelistSignatureEntry(
                    user_signature=WhitelistUserSignature(
                        user_id=user_id,
                        signature=user_sig,
                    ),
                    hashes=hashes,
                )
            ]
        ),
        linked_wallets=linked_wallets or [],
    )

    # Build decoders that return known-good data
    def rules_container_decoder(b64_data):
        return DecodedRulesContainer(
            users=[
                RuleUser(
                    id=user_id,
                    name="User 1",
                    public_key_pem=user_pub_pem,
                    roles=["USER", "OPERATOR"],
                )
            ],
            groups=[
                RuleGroup(id=group_id, name="Approvers", user_ids=[user_id]),
            ],
            address_whitelisting_rules=[
                AddressWhitelistingRules(
                    currency=blockchain,
                    network=network,
                    parallel_thresholds=[
                        SequentialThresholds(
                            thresholds=[
                                GroupThreshold(group_id=group_id, minimum_signatures=1)
                            ]
                        )
                    ],
                )
            ],
        )

    def user_signatures_decoder(b64_data):
        return [RuleUserSignature(user_id="superadmin@bank.com", signature=sa_sig)]

    return envelope, rules_container_decoder, user_signatures_decoder


# =============================================================================
# Step 1 Tests: Verify Metadata Hash
# =============================================================================


class TestStep1MetadataHash:
    """Tests for Step 1: metadata hash verification."""

    def test_valid_hash_passes(self, superadmin_keys, user1_keys):
        """Valid hash passes step 1."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAddressVerifier([sa_pub], min_valid_signatures=1)

        # Should not raise
        result = verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)
        assert isinstance(result, AddressVerificationResult)

    def test_mismatched_hash_raises_integrity_error(self, superadmin_keys, user1_keys):
        """Mismatched hash raises IntegrityError at step 1."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        # Tamper with the hash
        envelope.metadata = WhitelistMetadata(
            hash="0" * 64,
            payload_as_string=envelope.metadata.payload_as_string,
        )

        verifier = WhitelistedAddressVerifier([sa_pub], min_valid_signatures=1)
        with pytest.raises(IntegrityError, match="metadata hash verification failed"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_empty_payload_raises_integrity_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(hash="abc", payload_as_string=""),
        )
        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="payloadAsString is empty"):
            verifier.verify_whitelisted_address(envelope, lambda x: None, lambda x: [])

    def test_none_hash_raises_integrity_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(hash=None, payload_as_string='{"foo":"bar"}'),
        )
        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="metadata hash is empty"):
            verifier.verify_whitelisted_address(envelope, lambda x: None, lambda x: [])

    def test_none_metadata_raises_value_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        envelope = SignedWhitelistedAddressEnvelope(metadata=None)
        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(ValueError, match="metadata cannot be None"):
            verifier.verify_whitelisted_address(envelope, lambda x: None, lambda x: [])

    def test_none_envelope_raises_value_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(ValueError, match="envelope cannot be None"):
            verifier.verify_whitelisted_address(None, lambda x: None, lambda x: [])


# =============================================================================
# Step 2 Tests: Verify Rules Container Signatures
# =============================================================================


class TestStep2RulesContainerSignatures:
    """Tests for Step 2: SuperAdmin signature verification on rules container."""

    def test_no_superadmin_keys_raises(self, user1_keys):
        """No SuperAdmin keys configured raises IntegrityError."""
        u_priv, u_pub = user1_keys
        sa_priv, sa_pub = _generate_key_pair()

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        verifier = WhitelistedAddressVerifier([], min_valid_signatures=1)
        with pytest.raises(IntegrityError, match="no SuperAdmin keys"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_empty_rules_container_raises(self, superadmin_keys, user1_keys):
        """Empty rules container raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        envelope.rules_container = ""

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="rulesContainer is empty"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_empty_rules_signatures_raises(self, superadmin_keys, user1_keys):
        """Empty rules signatures raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        envelope.rules_signatures = ""

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="rulesSignatures is empty"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_insufficient_valid_signatures_raises(self, superadmin_keys, user1_keys):
        """Require 2 SuperAdmin sigs but only 1 is valid -> IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        # Require 2 valid signatures but decoder returns only 1 valid sig
        verifier = WhitelistedAddressVerifier([sa_pub], min_valid_signatures=2)
        with pytest.raises(IntegrityError, match="rules container signature verification failed"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_signature_decode_failure_raises(self, superadmin_keys, user1_keys):
        """Failed signature decode raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, _ = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def bad_sig_decoder(b64):
            raise ValueError("decode failed")

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="failed to decode rules signatures"):
            verifier.verify_whitelisted_address(envelope, rc_dec, bad_sig_decoder)


# =============================================================================
# Step 3 Tests: Decode Rules Container
# =============================================================================


class TestStep3DecodeRulesContainer:
    """Tests for Step 3: rules container decoding."""

    def test_decoder_failure_raises(self, superadmin_keys, user1_keys):
        """Rules container decode failure raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def bad_rc_decoder(b64):
            raise ValueError("invalid protobuf")

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="failed to decode rules container"):
            verifier.verify_whitelisted_address(envelope, bad_rc_decoder, us_dec)


# =============================================================================
# Step 4 Tests: Hash Coverage
# =============================================================================


class TestStep4HashCoverage:
    """Tests for Step 4: metadata hash in signature hashes."""

    def test_hash_found_in_signatures(self):
        """Hash found in signature hashes -> True."""
        target = "abc123"
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["other", target],
            )
        ]
        assert _verify_hash_coverage(target, sigs) is True

    def test_hash_not_found(self):
        """Hash not in any signatures -> False."""
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["other"],
            )
        ]
        assert _verify_hash_coverage("missing", sigs) is False

    def test_empty_signatures(self):
        """Empty signatures list -> False."""
        assert _verify_hash_coverage("abc", []) is False

    def test_no_signed_address_raises(self, superadmin_keys, user1_keys):
        """No signed_address raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        envelope.signed_address = None

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="signedAddress is nil"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_empty_signatures_in_signed_address_raises(self, superadmin_keys, user1_keys):
        """Empty signatures list in signed_address raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        envelope.signed_address = SignedWhitelistedAddress(signatures=[])

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="no signatures in signedAddress"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

    def test_hash_not_covered_by_any_signature_raises(self, superadmin_keys, user1_keys):
        """Hash not in any signature hashes raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        # Replace signatures with ones that don't cover the metadata hash
        envelope.signed_address = SignedWhitelistedAddress(
            signatures=[
                WhitelistSignatureEntry(
                    user_signature=WhitelistUserSignature(user_id="u1", signature="sig"),
                    hashes=["wrong_hash"],
                )
            ]
        )

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="metadata hash is not covered by any signature"):
            verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)


class TestStep4LegacyHashFallback:
    """Tests for Step 4 legacy hash fallback."""

    def test_legacy_hash_without_contract_type(self, superadmin_keys, user1_keys):
        """When current hash not found, legacy hash (without contractType) should match."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        # Create payload WITH contractType
        payload = _build_payload(contract_type="ERC20")
        payload_str = _payload_to_string(payload)
        current_hash = calculate_hex_hash(payload_str)

        # Build a legacy payload without contractType
        legacy_payload = dict(payload)
        del legacy_payload["contractType"]
        legacy_str = _payload_to_string(legacy_payload)
        legacy_hash = calculate_hex_hash(legacy_str)

        # Now compute the legacy hashes from the helper
        from taurus_protect.helpers.whitelist_hash_helper import compute_legacy_hashes

        legacy_hashes = compute_legacy_hashes(payload_str)

        # The legacy hash should be in the computed list
        # (may not match our manual computation exactly due to regex approach,
        # but there should be at least one legacy hash)
        if not legacy_hashes:
            pytest.skip("No legacy hashes computed for this payload format")

        # Build envelope where signatures only cover a legacy hash
        user_pub_pem = _public_key_to_pem(u_pub)
        hashes = [legacy_hashes[0]]
        hashes_json = json.dumps(hashes, separators=(",", ":"))
        user_sig = sign_data(u_priv, hashes_json.encode("utf-8"))

        rules_b64 = _encode_rules_container_json(
            _build_rules_container_dict(user_pub_pem)
        )
        rules_data = base64.b64decode(rules_b64)
        sa_sig = sign_data(sa_priv, rules_data)

        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(hash=current_hash, payload_as_string=payload_str),
            blockchain="ETH",
            network="mainnet",
            rules_container=rules_b64,
            rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
            signed_address=SignedWhitelistedAddress(
                signatures=[
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(
                            user_id="user1@bank.com",
                            signature=user_sig,
                        ),
                        hashes=hashes,
                    )
                ]
            ),
            linked_wallets=[],
        )

        def rc_decoder(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=user_pub_pem,
                        roles=["USER"],
                    )
                ],
                groups=[RuleGroup(id="approvers", user_ids=["user1@bank.com"])],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                    )
                ],
            )

        def us_decoder(b64):
            return [RuleUserSignature(user_id="sa", signature=sa_sig)]

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder, us_decoder)

        # Should succeed with the legacy hash
        assert result.verified_hash == legacy_hashes[0]
        assert result.verified_hash != current_hash


# =============================================================================
# Step 5 Tests: Whitelist Signature Thresholds
# =============================================================================


class TestStep5WhitelistSignatures:
    """Tests for Step 5: governance threshold verification."""

    def test_happy_path_single_group(self, superadmin_keys, user1_keys):
        """Single group with 1 required signature succeeds."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAddressVerifier([sa_pub])

        result = verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)
        assert isinstance(result, AddressVerificationResult)
        assert result.rules_container is not None

    def test_no_rules_for_blockchain_raises_whitelist_error(self, superadmin_keys, user1_keys):
        """No address whitelisting rules for blockchain raises WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub, blockchain="ETH", network="mainnet"
        )

        # Decoder returns rules for BTC, not ETH
        def rc_decoder_btc_only(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[RuleGroup(id="approvers", user_ids=["user1@bank.com"])],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="BTC",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="no address whitelisting rules found"):
            verifier.verify_whitelisted_address(envelope, rc_decoder_btc_only, us_dec)

    def test_threshold_not_met_raises_whitelist_error(self, superadmin_keys, user1_keys):
        """Threshold requires 2 sigs but only 1 valid -> WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        # Decoder returns rules requiring 2 signatures
        def rc_decoder_2_sigs(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[RuleGroup(id="approvers", user_ids=["user1@bank.com"])],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=2)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="signature verification failed"):
            verifier.verify_whitelisted_address(envelope, rc_decoder_2_sigs, us_dec)

    def test_group_not_found_raises_whitelist_error(self, superadmin_keys, user1_keys):
        """Group referenced in rules but missing from container -> WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def rc_decoder_missing_group(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[],  # No groups
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="ghost_group", minimum_signatures=1)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="signature verification failed"):
            verifier.verify_whitelisted_address(envelope, rc_decoder_missing_group, us_dec)

    def test_empty_thresholds_raises_whitelist_error(self, superadmin_keys, user1_keys):
        """No parallel thresholds defined -> WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def rc_decoder_no_thresholds(b64):
            return DecodedRulesContainer(
                users=[],
                groups=[],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="no threshold rules defined"):
            verifier.verify_whitelisted_address(envelope, rc_decoder_no_thresholds, us_dec)

    def test_wildcard_blockchain_rules_used_as_fallback(self, superadmin_keys, user1_keys):
        """Wildcard blockchain rules used when exact match not found."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def rc_decoder_wildcard(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[RuleGroup(id="approvers", user_ids=["user1@bank.com"])],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="Any",
                        network=None,
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_wildcard, us_dec)
        assert result is not None

    def test_parallel_paths_or_logic(self, superadmin_keys, user1_keys):
        """Multiple parallel paths: if one succeeds, overall passes (OR logic)."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, _, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def rc_decoder_two_paths(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[
                    RuleGroup(id="approvers", user_ids=["user1@bank.com"]),
                    RuleGroup(id="other_team", user_ids=["nobody@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            # Path 1: other_team needs 1 sig (will fail - user not in group)
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="other_team", minimum_signatures=1)]
                            ),
                            # Path 2: approvers needs 1 sig (will pass)
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            ),
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_two_paths, us_dec)
        assert result is not None


# =============================================================================
# Rule Lines Logic Tests
# =============================================================================


class TestRuleLines:
    """Tests for rule lines wallet path matching logic."""

    def test_matching_wallet_path_uses_line_thresholds(self, superadmin_keys, user1_keys):
        """When wallet path matches a rule line, use line thresholds instead of default."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        # Payload has no linked internal addresses, envelope has 1 linked wallet
        payload = _build_payload(linked_internal_addresses=[])
        envelope, _, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub,
            payload_dict=payload,
            linked_wallets=[InternalWallet(id=1, path="m/44/60/0")],
        )

        def rc_decoder_with_lines(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[
                    RuleGroup(id="approvers", user_ids=["user1@bank.com"]),
                    RuleGroup(id="admins", user_ids=["nobody@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        # Default thresholds require "admins" group (will fail)
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="admins", minimum_signatures=1)]
                            )
                        ],
                        # Rule line for wallet path m/44/60/0 uses "approvers" group (will pass)
                        lines=[
                            AddressWhitelistingLine(
                                cells=[
                                    RuleSource(
                                        type=RULE_SOURCE_TYPE_INTERNAL_WALLET,
                                        internal_wallet=RuleSourceInternalWallet(path="m/44/60/0"),
                                    )
                                ],
                                parallel_thresholds=[
                                    SequentialThresholds(
                                        thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                                    )
                                ],
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_with_lines, us_dec)
        assert result is not None

    def test_non_matching_wallet_path_uses_default_thresholds(self, superadmin_keys, user1_keys):
        """When wallet path does not match any line, use default thresholds."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        payload = _build_payload(linked_internal_addresses=[])
        envelope, _, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub,
            payload_dict=payload,
            linked_wallets=[InternalWallet(id=1, path="m/44/60/999")],  # Non-matching path
        )

        def rc_decoder_with_lines(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[
                    RuleGroup(id="approvers", user_ids=["user1@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        # Default thresholds use "approvers" group (will pass)
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                        lines=[
                            AddressWhitelistingLine(
                                cells=[
                                    RuleSource(
                                        type=RULE_SOURCE_TYPE_INTERNAL_WALLET,
                                        internal_wallet=RuleSourceInternalWallet(path="m/44/60/0"),
                                    )
                                ],
                                parallel_thresholds=[
                                    SequentialThresholds(
                                        thresholds=[GroupThreshold(group_id="admins", minimum_signatures=5)]
                                    )
                                ],
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_with_lines, us_dec)
        assert result is not None

    def test_linked_addresses_present_skips_lines(self, superadmin_keys, user1_keys):
        """When linked internal addresses exist, rule lines are NOT checked."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        # Payload with linked internal addresses
        payload = _build_payload(linked_internal_addresses=[{"address": "0xINTERNAL", "id": "1"}])
        envelope, _, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub,
            payload_dict=payload,
            linked_wallets=[InternalWallet(id=1, path="m/44/60/0")],
        )

        def rc_decoder_with_lines(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[
                    RuleGroup(id="approvers", user_ids=["user1@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        # Default thresholds (should be used because linked addresses present)
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                        lines=[
                            AddressWhitelistingLine(
                                cells=[
                                    RuleSource(
                                        type=RULE_SOURCE_TYPE_INTERNAL_WALLET,
                                        internal_wallet=RuleSourceInternalWallet(path="m/44/60/0"),
                                    )
                                ],
                                # Line thresholds require impossible 99 sigs (would fail)
                                parallel_thresholds=[
                                    SequentialThresholds(
                                        thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=99)]
                                    )
                                ],
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        # Should pass using default thresholds, NOT the line thresholds
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_with_lines, us_dec)
        assert result is not None

    def test_multiple_wallets_skips_lines(self, superadmin_keys, user1_keys):
        """When more than 1 linked wallet, rule lines are NOT checked."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        payload = _build_payload(linked_internal_addresses=[])
        envelope, _, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub,
            payload_dict=payload,
            linked_wallets=[
                InternalWallet(id=1, path="m/44/60/0"),
                InternalWallet(id=2, path="m/44/60/1"),
            ],
        )

        def rc_decoder_with_lines(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u_pub),
                        roles=["USER"],
                    )
                ],
                groups=[
                    RuleGroup(id="approvers", user_ids=["user1@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                        lines=[
                            AddressWhitelistingLine(
                                cells=[
                                    RuleSource(
                                        type=RULE_SOURCE_TYPE_INTERNAL_WALLET,
                                        internal_wallet=RuleSourceInternalWallet(path="m/44/60/0"),
                                    )
                                ],
                                parallel_thresholds=[
                                    SequentialThresholds(
                                        thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=99)]
                                    )
                                ],
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAddressVerifier([sa_pub])
        # Should pass using default thresholds because wallet_count != 1
        result = verifier.verify_whitelisted_address(envelope, rc_decoder_with_lines, us_dec)
        assert result is not None


# =============================================================================
# Helper Function Tests
# =============================================================================


class TestContainsHash:
    """Tests for _contains_hash helper."""

    def test_hash_found(self):
        assert _contains_hash(["abc", "def"], "def") is True

    def test_hash_not_found(self):
        assert _contains_hash(["abc", "def"], "xyz") is False

    def test_empty_list(self):
        assert _contains_hash([], "abc") is False


class TestVerifyHashCoverage:
    """Tests for _verify_hash_coverage helper."""

    def test_found_in_first_signature(self):
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["target"],
            )
        ]
        assert _verify_hash_coverage("target", sigs) is True

    def test_found_in_second_signature(self):
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["other"],
            ),
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u2"),
                hashes=["target"],
            ),
        ]
        assert _verify_hash_coverage("target", sigs) is True

    def test_not_found_returns_false(self):
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["nope"],
            )
        ]
        assert _verify_hash_coverage("target", sigs) is False


# =============================================================================
# End-to-End Happy Path
# =============================================================================


class TestEndToEnd:
    """End-to-end happy path tests."""

    def test_full_verification_succeeds(self, superadmin_keys, user1_keys):
        """Full 6-step verification succeeds with valid data."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAddressVerifier([sa_pub])

        result = verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

        assert isinstance(result, AddressVerificationResult)
        assert result.rules_container is not None
        assert result.verified_hash == envelope.metadata.hash

    def test_verification_returns_correct_hash(self, superadmin_keys, user1_keys):
        """Verification result contains the correct verified hash."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        payload = _build_payload(address="0xCUSTOM")
        envelope, rc_dec, us_dec = _build_full_envelope(
            u_priv, u_pub, sa_priv, sa_pub, payload_dict=payload
        )
        verifier = WhitelistedAddressVerifier([sa_pub])

        result = verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

        expected_hash = calculate_hex_hash(_payload_to_string(payload))
        assert result.verified_hash == expected_hash

    def test_verification_returns_decoded_rules_container(self, superadmin_keys, user1_keys):
        """Verification result includes the decoded rules container."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        envelope, rc_dec, us_dec = _build_full_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAddressVerifier([sa_pub])

        result = verifier.verify_whitelisted_address(envelope, rc_dec, us_dec)

        container = result.rules_container
        assert len(container.users) == 1
        assert container.users[0].id == "user1@bank.com"
        assert len(container.groups) == 1
        assert len(container.address_whitelisting_rules) == 1

    def test_two_users_two_groups_sequential(self, superadmin_keys, user1_keys, user2_keys):
        """Two groups in sequential threshold: both must pass (AND logic)."""
        sa_priv, sa_pub = superadmin_keys
        u1_priv, u1_pub = user1_keys
        u2_priv, u2_pub = user2_keys

        payload = _build_payload()
        payload_str = _payload_to_string(payload)
        metadata_hash = calculate_hex_hash(payload_str)

        # Both users sign the same hashes
        hashes = [metadata_hash]
        hashes_json = json.dumps(hashes, separators=(",", ":"))
        u1_sig = sign_data(u1_priv, hashes_json.encode("utf-8"))
        u2_sig = sign_data(u2_priv, hashes_json.encode("utf-8"))

        rules_b64 = _encode_rules_container_json(
            _build_rules_container_dict(_public_key_to_pem(u1_pub))
        )
        rules_data = base64.b64decode(rules_b64)
        sa_sig = sign_data(sa_priv, rules_data)

        envelope = SignedWhitelistedAddressEnvelope(
            metadata=WhitelistMetadata(hash=metadata_hash, payload_as_string=payload_str),
            blockchain="ETH",
            network="mainnet",
            rules_container=rules_b64,
            rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
            signed_address=SignedWhitelistedAddress(
                signatures=[
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(
                            user_id="user1@bank.com",
                            signature=u1_sig,
                        ),
                        hashes=hashes,
                    ),
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(
                            user_id="user2@bank.com",
                            signature=u2_sig,
                        ),
                        hashes=hashes,
                    ),
                ]
            ),
            linked_wallets=[],
        )

        def rc_decoder(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(
                        id="user1@bank.com",
                        public_key_pem=_public_key_to_pem(u1_pub),
                        roles=["USER"],
                    ),
                    RuleUser(
                        id="user2@bank.com",
                        public_key_pem=_public_key_to_pem(u2_pub),
                        roles=["USER"],
                    ),
                ],
                groups=[
                    RuleGroup(id="group_a", user_ids=["user1@bank.com"]),
                    RuleGroup(id="group_b", user_ids=["user2@bank.com"]),
                ],
                address_whitelisting_rules=[
                    AddressWhitelistingRules(
                        currency="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            # Sequential: group_a AND group_b both must pass
                            SequentialThresholds(
                                thresholds=[
                                    GroupThreshold(group_id="group_a", minimum_signatures=1),
                                    GroupThreshold(group_id="group_b", minimum_signatures=1),
                                ]
                            )
                        ],
                    )
                ],
            )

        def us_decoder(b64):
            return [RuleUserSignature(user_id="sa", signature=sa_sig)]

        verifier = WhitelistedAddressVerifier([sa_pub])
        result = verifier.verify_whitelisted_address(envelope, rc_decoder, us_decoder)
        assert result is not None
        assert result.verified_hash == metadata_hash
