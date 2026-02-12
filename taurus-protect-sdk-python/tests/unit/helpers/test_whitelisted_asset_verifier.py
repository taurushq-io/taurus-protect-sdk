"""Tests for WhitelistedAssetVerifier.

Tests the 5-step verification flow for whitelisted assets (contracts):
1. Verify metadata hash
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container
4. Verify hash coverage (metadata hash in signature hashes)
5. Verify whitelist signatures meet governance thresholds

Also tests:
- Legacy hash fallback for backward compatibility
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
from taurus_protect.helpers.whitelisted_asset_verifier import (
    AssetVerificationResult,
    WhitelistedAssetVerifier,
    verify_hash_coverage,
)
from taurus_protect.models.governance_rules import (
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    GroupThreshold,
    RuleGroup,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
)
from taurus_protect.models.whitelisted_address import (
    SignedContractAddress,
    WhitelistedAsset,
    WhitelistedAssetMetadata,
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


def _build_asset_payload(
    blockchain="ETH",
    network="mainnet",
    contract_address="0xUSDC",
    name="USDC",
    symbol="USDC",
    decimals=6,
    is_nft=None,
    kind_type=None,
):
    """Build a whitelisted asset payload dict."""
    payload = {
        "blockchain": blockchain,
        "network": network,
        "contractAddress": contract_address,
        "name": name,
        "symbol": symbol,
        "decimals": decimals,
    }
    if is_nft is not None:
        payload["isNFT"] = is_nft
    if kind_type is not None:
        payload["kindType"] = kind_type
    return payload


def _payload_to_string(payload: dict) -> str:
    """Serialize payload to compact JSON string."""
    return json.dumps(payload, separators=(",", ":"))


def _build_full_asset_envelope(
    user_private_key,
    user_public_key,
    superadmin_private_key,
    superadmin_public_key,
    payload_dict=None,
    blockchain="ETH",
    network="mainnet",
    group_id="approvers",
    user_id="user1@bank.com",
):
    """Build a fully valid WhitelistedAsset with decoders for testing.

    Returns (asset, rules_container_decoder, user_signatures_decoder).
    """
    if payload_dict is None:
        payload_dict = _build_asset_payload()

    payload_str = _payload_to_string(payload_dict)
    metadata_hash = calculate_hex_hash(payload_str)

    # Build rules container
    user_pub_pem = _public_key_to_pem(user_public_key)
    rules_b64 = _encode_rules_container_json({
        "users": [{"id": user_id, "publicKey": user_pub_pem, "roles": ["USER"]}],
        "groups": [{"id": group_id, "userIds": [user_id]}],
        "contractAddressWhitelistingRules": [
            {
                "blockchain": blockchain,
                "network": network,
                "parallelThresholds": [
                    {"groupId": group_id, "minimumSignatures": 1}
                ],
            }
        ],
    })
    rules_data = base64.b64decode(rules_b64)

    # Sign rules container with SuperAdmin key
    sa_sig = sign_data(superadmin_private_key, rules_data)

    # Sign the hashes array with user key
    hashes = [metadata_hash]
    hashes_json = json.dumps(hashes, separators=(",", ":"))
    user_sig = sign_data(user_private_key, hashes_json.encode("utf-8"))

    # Build the asset
    asset = WhitelistedAsset(
        id="asset-1",
        blockchain=blockchain,
        network=network,
        metadata=WhitelistedAssetMetadata(
            hash=metadata_hash,
            payload_as_string=payload_str,
        ),
        rules_container=rules_b64,
        rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
        signed_contract_address=SignedContractAddress(
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
    )

    # Build decoders
    def rules_container_decoder(b64_data):
        return DecodedRulesContainer(
            users=[
                RuleUser(
                    id=user_id,
                    name="User 1",
                    public_key_pem=user_pub_pem,
                    roles=["USER"],
                )
            ],
            groups=[
                RuleGroup(id=group_id, name="Approvers", user_ids=[user_id]),
            ],
            contract_address_whitelisting_rules=[
                ContractAddressWhitelistingRules(
                    blockchain=blockchain,
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

    return asset, rules_container_decoder, user_signatures_decoder


# =============================================================================
# Step 1 Tests: Verify Metadata Hash
# =============================================================================


class TestStep1MetadataHash:
    """Tests for Step 1: metadata hash verification."""

    def test_valid_hash_passes(self, superadmin_keys, user1_keys):
        """Valid hash passes step 1."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=1)

        result = verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)
        assert isinstance(result, AssetVerificationResult)

    def test_mismatched_hash_raises_integrity_error(self, superadmin_keys, user1_keys):
        """Mismatched hash raises IntegrityError at step 1."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        # Tamper with the hash - create new asset with bad hash
        tampered = WhitelistedAsset(
            id=asset.id,
            blockchain=asset.blockchain,
            network=asset.network,
            metadata=WhitelistedAssetMetadata(
                hash="0" * 64,
                payload_as_string=asset.metadata.payload_as_string,
            ),
            rules_container=asset.rules_container,
            rules_signatures=asset.rules_signatures,
            signed_contract_address=asset.signed_contract_address,
        )

        verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=1)
        with pytest.raises(IntegrityError, match="metadata hash verification failed"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)

    def test_empty_payload_raises_integrity_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys

        asset = WhitelistedAsset(
            id="1",
            metadata=WhitelistedAssetMetadata(hash="abc", payload_as_string=""),
        )
        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="payloadAsString is empty"):
            verifier.verify_whitelisted_asset(asset, lambda x: None, lambda x: [])

    def test_none_hash_raises_integrity_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys

        asset = WhitelistedAsset(
            id="1",
            metadata=WhitelistedAssetMetadata(hash=None, payload_as_string='{"foo":"bar"}'),
        )
        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="metadata hash is empty"):
            verifier.verify_whitelisted_asset(asset, lambda x: None, lambda x: [])

    def test_none_metadata_raises_value_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        asset = WhitelistedAsset(id="1", metadata=None)
        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(ValueError, match="metadata cannot be None"):
            verifier.verify_whitelisted_asset(asset, lambda x: None, lambda x: [])

    def test_none_asset_raises_value_error(self, superadmin_keys):
        _, sa_pub = superadmin_keys
        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(ValueError, match="asset cannot be None"):
            verifier.verify_whitelisted_asset(None, lambda x: None, lambda x: [])


# =============================================================================
# Step 2 Tests: Rules Container Signatures
# =============================================================================


class TestStep2RulesContainerSignatures:
    """Tests for Step 2: SuperAdmin signature verification."""

    def test_no_superadmin_keys_raises(self, user1_keys):
        """No SuperAdmin keys configured raises IntegrityError."""
        u_priv, u_pub = user1_keys
        sa_priv, sa_pub = _generate_key_pair()

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        verifier = WhitelistedAssetVerifier([], min_valid_signatures=1)
        with pytest.raises(IntegrityError, match="no SuperAdmin keys"):
            verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)

    def test_empty_rules_container_raises(self, superadmin_keys, user1_keys):
        """Empty rules container raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        tampered = WhitelistedAsset(
            id=asset.id,
            metadata=asset.metadata,
            rules_container="",
            rules_signatures=asset.rules_signatures,
            signed_contract_address=asset.signed_contract_address,
            blockchain=asset.blockchain,
            network=asset.network,
        )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="rulesContainer is empty"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)

    def test_empty_rules_signatures_raises(self, superadmin_keys, user1_keys):
        """Empty rules signatures raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        tampered = WhitelistedAsset(
            id=asset.id,
            metadata=asset.metadata,
            rules_container=asset.rules_container,
            rules_signatures="",
            signed_contract_address=asset.signed_contract_address,
            blockchain=asset.blockchain,
            network=asset.network,
        )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="rulesSignatures is empty"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)

    def test_insufficient_valid_signatures_raises(self, superadmin_keys, user1_keys):
        """Require 2 SuperAdmin sigs but only 1 is valid."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=2)
        with pytest.raises(IntegrityError, match="rules container signature verification failed"):
            verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)

    def test_signature_decode_failure_raises(self, superadmin_keys, user1_keys):
        """Failed signature decode raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, _ = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def bad_sig_decoder(b64):
            raise ValueError("decode failed")

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="failed to decode rules signatures"):
            verifier.verify_whitelisted_asset(asset, rc_dec, bad_sig_decoder)


# =============================================================================
# Step 3 Tests: Decode Rules Container
# =============================================================================


class TestStep3DecodeRulesContainer:
    """Tests for Step 3: rules container decoding."""

    def test_decoder_failure_raises(self, superadmin_keys, user1_keys):
        """Rules container decode failure raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, _, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def bad_rc_decoder(b64):
            raise ValueError("invalid protobuf")

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="failed to decode rules container"):
            verifier.verify_whitelisted_asset(asset, bad_rc_decoder, us_dec)


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
        assert verify_hash_coverage(target, sigs) is True

    def test_hash_not_found(self):
        """Hash not in any signatures -> False."""
        sigs = [
            WhitelistSignatureEntry(
                user_signature=WhitelistUserSignature(user_id="u1"),
                hashes=["other"],
            )
        ]
        assert verify_hash_coverage("missing", sigs) is False

    def test_empty_signatures(self):
        """Empty signatures list -> False."""
        assert verify_hash_coverage("abc", []) is False

    def test_no_signed_contract_address_raises(self, superadmin_keys, user1_keys):
        """No signed_contract_address raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        tampered = WhitelistedAsset(
            id=asset.id,
            metadata=asset.metadata,
            rules_container=asset.rules_container,
            rules_signatures=asset.rules_signatures,
            signed_contract_address=None,
            blockchain=asset.blockchain,
            network=asset.network,
        )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="signedContractAddress is nil"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)

    def test_empty_signatures_in_signed_contract_address_raises(self, superadmin_keys, user1_keys):
        """Empty signatures list raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        tampered = WhitelistedAsset(
            id=asset.id,
            metadata=asset.metadata,
            rules_container=asset.rules_container,
            rules_signatures=asset.rules_signatures,
            signed_contract_address=SignedContractAddress(signatures=[]),
            blockchain=asset.blockchain,
            network=asset.network,
        )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="no signatures in signedContractAddress"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)

    def test_hash_not_covered_raises(self, superadmin_keys, user1_keys):
        """Hash not in any signature hashes raises IntegrityError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        tampered = WhitelistedAsset(
            id=asset.id,
            metadata=asset.metadata,
            rules_container=asset.rules_container,
            rules_signatures=asset.rules_signatures,
            signed_contract_address=SignedContractAddress(
                signatures=[
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(user_id="u1", signature="sig"),
                        hashes=["wrong_hash"],
                    )
                ]
            ),
            blockchain=asset.blockchain,
            network=asset.network,
        )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(IntegrityError, match="metadata hash is not covered by any signature"):
            verifier.verify_whitelisted_asset(tampered, rc_dec, us_dec)


# =============================================================================
# Step 4 Legacy Hash Tests
# =============================================================================


class TestStep4LegacyHashFallback:
    """Tests for Step 4 legacy hash fallback for assets."""

    def test_legacy_hash_without_is_nft(self, superadmin_keys, user1_keys):
        """When current hash not found, legacy hash (without isNFT) should match in step 4.

        Step 4 uses legacy hash fallback to find hash coverage.
        Step 5 uses metadata.hash to check signatures, so the signature's hashes
        list must also include metadata.hash for the full flow to pass.
        In production, the signatures typically cover both the original and legacy hashes.
        """
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        # Create payload WITH isNFT
        payload = _build_asset_payload(is_nft=False)
        payload_str = _payload_to_string(payload)
        current_hash = calculate_hex_hash(payload_str)

        # Compute legacy hashes
        from taurus_protect.helpers.whitelist_hash_helper import compute_asset_legacy_hashes

        legacy_hashes = compute_asset_legacy_hashes(payload_str)

        if not legacy_hashes:
            pytest.skip("No legacy hashes computed for this payload format")

        # Build signatures that cover BOTH the legacy hash (for step 4 coverage)
        # and the current hash (for step 5 threshold verification).
        # This mirrors production where signers cover the metadata hash.
        user_pub_pem = _public_key_to_pem(u_pub)
        hashes = [legacy_hashes[0], current_hash]
        hashes_json = json.dumps(hashes, separators=(",", ":"))
        user_sig = sign_data(u_priv, hashes_json.encode("utf-8"))

        rules_b64 = _encode_rules_container_json({
            "users": [{"id": "user1@bank.com", "publicKey": user_pub_pem, "roles": ["USER"]}],
            "groups": [{"id": "approvers", "userIds": ["user1@bank.com"]}],
            "contractAddressWhitelistingRules": [
                {
                    "blockchain": "ETH",
                    "network": "mainnet",
                    "parallelThresholds": [
                        {"groupId": "approvers", "minimumSignatures": 1}
                    ],
                }
            ],
        })
        rules_data = base64.b64decode(rules_b64)
        sa_sig = sign_data(sa_priv, rules_data)

        asset = WhitelistedAsset(
            id="asset-legacy",
            blockchain="ETH",
            network="mainnet",
            metadata=WhitelistedAssetMetadata(
                hash=current_hash,
                payload_as_string=payload_str,
            ),
            rules_container=rules_b64,
            rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
            signed_contract_address=SignedContractAddress(
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
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="ETH",
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

        verifier = WhitelistedAssetVerifier([sa_pub])
        result = verifier.verify_whitelisted_asset(asset, rc_decoder, us_decoder)
        assert result is not None

    def test_legacy_hash_passes_full_verification(self, superadmin_keys, user1_keys):
        """When signatures cover a legacy hash, the verified hash from step 4
        is propagated to step 5, so the full verification succeeds."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        # Create payload WITH isNFT
        payload = _build_asset_payload(is_nft=False)
        payload_str = _payload_to_string(payload)
        current_hash = calculate_hex_hash(payload_str)

        from taurus_protect.helpers.whitelist_hash_helper import compute_asset_legacy_hashes

        legacy_hashes = compute_asset_legacy_hashes(payload_str)
        if not legacy_hashes:
            pytest.skip("No legacy hashes computed for this payload format")

        # Signatures only cover the legacy hash, NOT the current metadata hash
        user_pub_pem = _public_key_to_pem(u_pub)
        hashes = [legacy_hashes[0]]
        hashes_json = json.dumps(hashes, separators=(",", ":"))
        user_sig = sign_data(u_priv, hashes_json.encode("utf-8"))

        rules_b64 = _encode_rules_container_json({
            "users": [{"id": "user1@bank.com", "publicKey": user_pub_pem, "roles": ["USER"]}],
            "groups": [{"id": "approvers", "userIds": ["user1@bank.com"]}],
            "contractAddressWhitelistingRules": [
                {
                    "blockchain": "ETH",
                    "network": "mainnet",
                    "parallelThresholds": [
                        {"groupId": "approvers", "minimumSignatures": 1}
                    ],
                }
            ],
        })
        rules_data = base64.b64decode(rules_b64)
        sa_sig = sign_data(sa_priv, rules_data)

        asset = WhitelistedAsset(
            id="asset-legacy-only",
            blockchain="ETH",
            network="mainnet",
            metadata=WhitelistedAssetMetadata(
                hash=current_hash,
                payload_as_string=payload_str,
            ),
            rules_container=rules_b64,
            rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
            signed_contract_address=SignedContractAddress(
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
        )

        def rc_decoder(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(id="user1@bank.com", public_key_pem=user_pub_pem, roles=["USER"])
                ],
                groups=[RuleGroup(id="approvers", user_ids=["user1@bank.com"])],
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="ETH",
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

        verifier = WhitelistedAssetVerifier([sa_pub])
        # Step 4 passes (legacy hash found) and returns verified_hash.
        # Step 5 now uses verified_hash (the legacy hash) which IS in sig.hashes,
        # so the full verification succeeds.
        result = verifier.verify_whitelisted_asset(asset, rc_decoder, us_decoder)
        assert result is not None


# =============================================================================
# Step 5 Tests: Whitelist Signatures
# =============================================================================


class TestStep5WhitelistSignatures:
    """Tests for Step 5: governance threshold verification."""

    def test_happy_path(self, superadmin_keys, user1_keys):
        """Single group with 1 required signature succeeds."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAssetVerifier([sa_pub])

        result = verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)
        assert isinstance(result, AssetVerificationResult)
        assert result.rules_container is not None

    def test_no_rules_for_blockchain_raises(self, superadmin_keys, user1_keys):
        """No contract whitelisting rules for blockchain raises WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, _, us_dec = _build_full_asset_envelope(
            u_priv, u_pub, sa_priv, sa_pub, blockchain="ETH", network="mainnet"
        )

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
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="BTC",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=1)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="no contract address whitelisting rules found"):
            verifier.verify_whitelisted_asset(asset, rc_decoder_btc_only, us_dec)

    def test_threshold_not_met_raises(self, superadmin_keys, user1_keys):
        """Threshold requires 2 sigs but only 1 valid -> WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, _, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

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
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="ETH",
                        network="mainnet",
                        parallel_thresholds=[
                            SequentialThresholds(
                                thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=2)]
                            )
                        ],
                    )
                ],
            )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="signature verification failed"):
            verifier.verify_whitelisted_asset(asset, rc_decoder_2_sigs, us_dec)

    def test_empty_thresholds_raises(self, superadmin_keys, user1_keys):
        """No parallel thresholds defined -> WhitelistError."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, _, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)

        def rc_decoder_no_thresholds(b64):
            return DecodedRulesContainer(
                users=[],
                groups=[],
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="ETH",
                        network="mainnet",
                        parallel_thresholds=[],
                    )
                ],
            )

        verifier = WhitelistedAssetVerifier([sa_pub])
        with pytest.raises(WhitelistError, match="no threshold rules defined"):
            verifier.verify_whitelisted_asset(asset, rc_decoder_no_thresholds, us_dec)


# =============================================================================
# End-to-End Happy Path
# =============================================================================


class TestEndToEnd:
    """End-to-end happy path tests."""

    def test_full_verification_succeeds(self, superadmin_keys, user1_keys):
        """Full 5-step verification succeeds with valid data."""
        sa_priv, sa_pub = superadmin_keys
        u_priv, u_pub = user1_keys

        asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
        verifier = WhitelistedAssetVerifier([sa_pub])

        result = verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)

        assert isinstance(result, AssetVerificationResult)
        assert result.rules_container is not None

    def test_two_users_two_groups_sequential(self, superadmin_keys, user1_keys, user2_keys):
        """Two groups in sequential threshold: both must pass (AND logic)."""
        sa_priv, sa_pub = superadmin_keys
        u1_priv, u1_pub = user1_keys
        u2_priv, u2_pub = user2_keys

        payload = _build_asset_payload()
        payload_str = _payload_to_string(payload)
        metadata_hash = calculate_hex_hash(payload_str)

        hashes = [metadata_hash]
        hashes_json = json.dumps(hashes, separators=(",", ":"))
        u1_sig = sign_data(u1_priv, hashes_json.encode("utf-8"))
        u2_sig = sign_data(u2_priv, hashes_json.encode("utf-8"))

        rules_b64 = _encode_rules_container_json({
            "users": [], "groups": [],
        })
        rules_data = base64.b64decode(rules_b64)
        sa_sig = sign_data(sa_priv, rules_data)

        asset = WhitelistedAsset(
            id="asset-2",
            blockchain="ETH",
            network="mainnet",
            metadata=WhitelistedAssetMetadata(hash=metadata_hash, payload_as_string=payload_str),
            rules_container=rules_b64,
            rules_signatures=base64.b64encode(b"dummy").decode("utf-8"),
            signed_contract_address=SignedContractAddress(
                signatures=[
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(user_id="user1@bank.com", signature=u1_sig),
                        hashes=hashes,
                    ),
                    WhitelistSignatureEntry(
                        user_signature=WhitelistUserSignature(user_id="user2@bank.com", signature=u2_sig),
                        hashes=hashes,
                    ),
                ]
            ),
        )

        def rc_decoder(b64):
            return DecodedRulesContainer(
                users=[
                    RuleUser(id="user1@bank.com", public_key_pem=_public_key_to_pem(u1_pub), roles=["USER"]),
                    RuleUser(id="user2@bank.com", public_key_pem=_public_key_to_pem(u2_pub), roles=["USER"]),
                ],
                groups=[
                    RuleGroup(id="group_a", user_ids=["user1@bank.com"]),
                    RuleGroup(id="group_b", user_ids=["user2@bank.com"]),
                ],
                contract_address_whitelisting_rules=[
                    ContractAddressWhitelistingRules(
                        blockchain="ETH",
                        network="mainnet",
                        parallel_thresholds=[
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

        verifier = WhitelistedAssetVerifier([sa_pub])
        result = verifier.verify_whitelisted_asset(asset, rc_decoder, us_decoder)
        assert result is not None
