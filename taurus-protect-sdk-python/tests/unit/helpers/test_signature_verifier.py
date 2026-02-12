"""Tests for signature verification utilities."""

import base64

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.signature_verifier import (
    is_valid_signature,
    verify_raw_signature,
)


class TestIsValidSignature:
    """Tests for is_valid_signature function."""

    def test_valid_signature_with_matching_key(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that valid signature returns True."""
        data = b"test data to sign"
        signature = sign_data(ecdsa_private_key, data)

        result = is_valid_signature(data, signature, [ecdsa_public_key])
        assert result is True

    def test_valid_signature_with_multiple_keys(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
        second_ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that valid signature returns True when key is in list."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        # Correct key is second in the list
        result = is_valid_signature(data, signature, [second_ecdsa_public_key, ecdsa_public_key])
        assert result is True

    def test_invalid_signature_with_wrong_key(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        second_ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that invalid signature returns False."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        # Only wrong key in list
        result = is_valid_signature(data, signature, [second_ecdsa_public_key])
        assert result is False

    def test_invalid_signature_empty_key_list(
        self, ecdsa_private_key: ec.EllipticCurvePrivateKey
    ) -> None:
        """Test that empty key list returns False."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        result = is_valid_signature(data, signature, [])
        assert result is False

    def test_corrupted_signature(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test that corrupted signature returns False."""
        data = b"test data"
        corrupted_sig = base64.b64encode(b"\x00" * 64).decode("utf-8")

        result = is_valid_signature(data, corrupted_sig, [ecdsa_public_key])
        assert result is False

    def test_invalid_base64_signature(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test that invalid base64 returns False."""
        data = b"test data"

        result = is_valid_signature(data, "not-valid-base64!!!", [ecdsa_public_key])
        assert result is False

    def test_wrong_data(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that signature of different data returns False."""
        signature = sign_data(ecdsa_private_key, b"original data")

        result = is_valid_signature(b"different data", signature, [ecdsa_public_key])
        assert result is False


class TestVerifyRawSignature:
    """Tests for verify_raw_signature function."""

    def test_valid_raw_signature(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test verifying valid raw signature."""
        data = b"test data"
        signature_b64 = sign_data(ecdsa_private_key, data)
        signature_bytes = base64.b64decode(signature_b64)

        result = verify_raw_signature(data, signature_bytes, ecdsa_public_key)
        assert result is True

    def test_invalid_raw_signature(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test verifying invalid raw signature."""
        data = b"test data"
        invalid_sig = b"\x00" * 64

        result = verify_raw_signature(data, invalid_sig, ecdsa_public_key)
        assert result is False

    def test_none_data_returns_false(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test that None data returns False."""
        result = verify_raw_signature(None, b"sig", ecdsa_public_key)  # type: ignore
        assert result is False

    def test_none_signature_returns_false(
        self, ecdsa_public_key: ec.EllipticCurvePublicKey
    ) -> None:
        """Test that None signature returns False."""
        result = verify_raw_signature(b"data", None, ecdsa_public_key)  # type: ignore
        assert result is False

    def test_none_key_returns_false(self) -> None:
        """Test that None key returns False."""
        result = verify_raw_signature(b"data", b"sig", None)  # type: ignore
        assert result is False


class TestVerifyGovernanceRulesIntegration:
    """Integration tests for governance rules verification.

    Note: Full verify_governance_rules tests require GovernanceRules model,
    so we test the underlying signature validation logic here.
    """

    def test_multiple_valid_signatures(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
        second_ecdsa_private_key: ec.EllipticCurvePrivateKey,
        second_ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test validating multiple signatures (common in governance rules)."""
        data = b"rules container data"

        sig1 = sign_data(ecdsa_private_key, data)
        sig2 = sign_data(second_ecdsa_private_key, data)

        keys = [ecdsa_public_key, second_ecdsa_public_key]

        # Both signatures should be valid against the key list
        assert is_valid_signature(data, sig1, keys) is True
        assert is_valid_signature(data, sig2, keys) is True

    def test_signature_count_threshold(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
        second_ecdsa_private_key: ec.EllipticCurvePrivateKey,
        second_ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test counting valid signatures against threshold."""
        data = b"rules container data"

        # Create signatures
        signatures = [
            sign_data(ecdsa_private_key, data),
            sign_data(second_ecdsa_private_key, data),
        ]

        keys = [ecdsa_public_key, second_ecdsa_public_key]

        # Count valid signatures
        valid_count = sum(1 for sig in signatures if is_valid_signature(data, sig, keys))

        assert valid_count == 2
