"""Tests for address_signature_verifier.

Tests the HSM signature verification for blockchain addresses:
- Valid signatures pass verification
- Invalid/tampered signatures fail
- None inputs raise appropriate errors
- Missing HSM key raises IntegrityError
"""

from __future__ import annotations

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.address_signature_verifier import verify_address_signature
from taurus_protect.models.address import Address
from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser


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


def _build_rules_container_with_hsm(hsm_pub_pem: str) -> DecodedRulesContainer:
    """Build a rules container with an HSM slot user."""
    return DecodedRulesContainer(
        users=[
            RuleUser(
                id="hsm-slot-1",
                name="HSM Slot",
                public_key_pem=hsm_pub_pem,
                roles=["HSMSLOT"],
            )
        ],
    )


@pytest.fixture
def hsm_keys():
    """Generate HSM key pair."""
    return _generate_key_pair()


@pytest.fixture
def hsm_container(hsm_keys):
    """Build a rules container with HSM key."""
    _, hsm_pub = hsm_keys
    return _build_rules_container_with_hsm(_public_key_to_pem(hsm_pub))


# =============================================================================
# None/Missing Input Tests
# =============================================================================


class TestNilInputs:
    """Tests for None input handling."""

    def test_none_address_raises_value_error(self, hsm_container):
        with pytest.raises(ValueError, match="address cannot be None"):
            verify_address_signature(None, hsm_container)

    def test_none_rules_container_raises_value_error(self):
        addr = Address(id="1", wallet_id="w1", address="0xABC", signature="sig")
        with pytest.raises(ValueError, match="rules_container cannot be None"):
            verify_address_signature(addr, None)


# =============================================================================
# Missing Field Tests
# =============================================================================


class TestMissingFields:
    """Tests for missing address fields."""

    def test_empty_signature_raises_integrity_error(self, hsm_container):
        addr = Address(id="1", wallet_id="w1", address="0xABC", signature="")
        with pytest.raises(IntegrityError, match="has no signature"):
            verify_address_signature(addr, hsm_container)

    def test_none_signature_raises_integrity_error(self, hsm_container):
        addr = Address(id="1", wallet_id="w1", address="0xABC", signature=None)
        with pytest.raises(IntegrityError, match="has no signature"):
            verify_address_signature(addr, hsm_container)

    def test_empty_address_string_raises_integrity_error(self, hsm_container):
        addr = Address(id="1", wallet_id="w1", address="", signature="somesig")
        with pytest.raises(IntegrityError, match="has no blockchain address"):
            verify_address_signature(addr, hsm_container)


# =============================================================================
# Missing HSM Key Tests
# =============================================================================


class TestMissingHSMKey:
    """Tests for missing HSM key in rules container."""

    def test_no_hsm_user_raises_integrity_error(self):
        container = DecodedRulesContainer(
            users=[
                RuleUser(
                    id="regular-user",
                    name="Regular User",
                    public_key_pem="-----BEGIN PUBLIC KEY-----\nfake\n-----END PUBLIC KEY-----",
                    roles=["USER", "OPERATOR"],
                )
            ],
        )
        addr = Address(id="1", wallet_id="w1", address="0xABC", signature="somesig")
        with pytest.raises(IntegrityError, match="HSMSLOT"):
            verify_address_signature(addr, container)

    def test_empty_users_raises_integrity_error(self):
        container = DecodedRulesContainer(users=[])
        addr = Address(id="1", wallet_id="w1", address="0xABC", signature="somesig")
        with pytest.raises(IntegrityError, match="HSMSLOT"):
            verify_address_signature(addr, container)


# =============================================================================
# Valid Signature Tests
# =============================================================================


class TestValidSignature:
    """Tests for successful signature verification."""

    def test_valid_signature_passes(self, hsm_keys, hsm_container):
        hsm_priv, _ = hsm_keys
        address_str = "0x1234567890abcdef"
        sig = sign_data(hsm_priv, address_str.encode("utf-8"))

        addr = Address(id="1", wallet_id="w1", address=address_str, signature=sig)
        # Should not raise
        verify_address_signature(addr, hsm_container)

    def test_valid_long_address(self, hsm_keys, hsm_container):
        """Verify with a longer address string (e.g., Solana)."""
        hsm_priv, _ = hsm_keys
        address_str = "7xJfzrHn2H5vKJRZzHQEEMnDH7N8pZLaVbrcxBz3Y4wT"
        sig = sign_data(hsm_priv, address_str.encode("utf-8"))

        addr = Address(id="2", wallet_id="w1", address=address_str, signature=sig)
        verify_address_signature(addr, hsm_container)


# =============================================================================
# Invalid Signature Tests
# =============================================================================


class TestInvalidSignature:
    """Tests for signature verification failures."""

    def test_signature_for_different_address_fails(self, hsm_keys, hsm_container):
        hsm_priv, _ = hsm_keys
        # Sign a different address than what we verify
        sig = sign_data(hsm_priv, b"0xDIFFERENT_ADDRESS")

        addr = Address(id="1", wallet_id="w1", address="0x1234567890abcdef", signature=sig)
        with pytest.raises(IntegrityError, match="Address signature verification failed"):
            verify_address_signature(addr, hsm_container)

    def test_wrong_key_fails(self):
        """Sign with one key, verify with another."""
        hsm_priv, _ = _generate_key_pair()
        _, other_pub = _generate_key_pair()
        container = _build_rules_container_with_hsm(_public_key_to_pem(other_pub))

        address_str = "0x1234567890abcdef"
        sig = sign_data(hsm_priv, address_str.encode("utf-8"))

        addr = Address(id="1", wallet_id="w1", address=address_str, signature=sig)
        with pytest.raises(IntegrityError, match="Address signature verification failed"):
            verify_address_signature(addr, container)

    def test_malformed_signature_fails(self, hsm_container):
        addr = Address(id="1", wallet_id="w1", address="0x1234567890abcdef", signature="not-valid-base64!!!")
        with pytest.raises(IntegrityError, match="Address signature verification failed"):
            verify_address_signature(addr, hsm_container)

    def test_truncated_signature_fails(self, hsm_keys, hsm_container):
        """Truncated (partial) signature should fail."""
        hsm_priv, _ = hsm_keys
        address_str = "0x1234567890abcdef"
        sig = sign_data(hsm_priv, address_str.encode("utf-8"))

        # Truncate the signature
        import base64
        sig_bytes = base64.b64decode(sig)
        truncated = base64.b64encode(sig_bytes[:len(sig_bytes) // 2]).decode("utf-8")

        addr = Address(id="1", wallet_id="w1", address=address_str, signature=truncated)
        with pytest.raises(IntegrityError, match="Address signature verification failed"):
            verify_address_signature(addr, hsm_container)
