"""Tests for ECDSA signing utilities."""

import base64

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.crypto.signing import (
    get_public_key_from_private,
    sign_data,
    verify_signature,
)


class TestSignData:
    """Tests for sign_data function."""

    def test_sign_returns_base64_string(
        self, ecdsa_private_key: ec.EllipticCurvePrivateKey
    ) -> None:
        """Test that sign_data returns base64-encoded string."""
        data = b"test data to sign"
        signature = sign_data(ecdsa_private_key, data)

        assert isinstance(signature, str)
        # Should be valid base64
        decoded = base64.b64decode(signature)
        assert len(decoded) == 64  # P-256 produces 64-byte r||s signature

    def test_sign_produces_raw_rs_format(
        self, ecdsa_private_key: ec.EllipticCurvePrivateKey
    ) -> None:
        """Test that signature is in raw r||s format (not DER)."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        decoded = base64.b64decode(signature)
        # Raw r||s for P-256 is exactly 64 bytes (32 for r, 32 for s)
        assert len(decoded) == 64

        # Extract r and s
        r = int.from_bytes(decoded[:32], byteorder="big")
        s = int.from_bytes(decoded[32:], byteorder="big")

        # Both r and s should be positive and within curve order
        assert r > 0
        assert s > 0

    def test_sign_different_data_different_signatures(
        self, ecdsa_private_key: ec.EllipticCurvePrivateKey
    ) -> None:
        """Test that different data produces different signatures."""
        sig1 = sign_data(ecdsa_private_key, b"data 1")
        sig2 = sign_data(ecdsa_private_key, b"data 2")

        assert sig1 != sig2

    def test_sign_empty_data(self, ecdsa_private_key: ec.EllipticCurvePrivateKey) -> None:
        """Test signing empty data."""
        signature = sign_data(ecdsa_private_key, b"")
        assert isinstance(signature, str)
        decoded = base64.b64decode(signature)
        assert len(decoded) == 64

    def test_sign_large_data(self, ecdsa_private_key: ec.EllipticCurvePrivateKey) -> None:
        """Test signing large data."""
        large_data = b"x" * 10000
        signature = sign_data(ecdsa_private_key, large_data)
        assert isinstance(signature, str)
        # Signature size is fixed regardless of data size
        decoded = base64.b64decode(signature)
        assert len(decoded) == 64


class TestVerifySignature:
    """Tests for verify_signature function."""

    def test_verify_valid_signature(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test verifying a valid signature."""
        data = b"test data to sign"
        signature = sign_data(ecdsa_private_key, data)

        assert verify_signature(ecdsa_public_key, data, signature) is True

    def test_verify_wrong_data_fails(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that verification fails with wrong data."""
        signature = sign_data(ecdsa_private_key, b"original data")

        assert verify_signature(ecdsa_public_key, b"modified data", signature) is False

    def test_verify_wrong_key_fails(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        second_ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that verification fails with wrong public key."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        # Using different key should fail
        assert verify_signature(second_ecdsa_public_key, data, signature) is False

    def test_verify_corrupted_signature_fails(
        self, ecdsa_public_key: ec.EllipticCurvePublicKey
    ) -> None:
        """Test that corrupted signature fails verification."""
        data = b"test data"
        # Create an invalid signature
        invalid_sig = base64.b64encode(b"\x00" * 64).decode("utf-8")

        assert verify_signature(ecdsa_public_key, data, invalid_sig) is False

    def test_verify_invalid_base64_fails(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test that invalid base64 fails gracefully."""
        assert verify_signature(ecdsa_public_key, b"data", "not-valid-base64!!!") is False

    def test_verify_wrong_signature_length_fails(
        self, ecdsa_public_key: ec.EllipticCurvePublicKey
    ) -> None:
        """Test that wrong signature length fails."""
        data = b"test data"
        # Too short (32 bytes instead of 64)
        short_sig = base64.b64encode(b"\x01" * 32).decode("utf-8")
        assert verify_signature(ecdsa_public_key, data, short_sig) is False

        # Too long (96 bytes)
        long_sig = base64.b64encode(b"\x01" * 96).decode("utf-8")
        assert verify_signature(ecdsa_public_key, data, long_sig) is False

    def test_verify_empty_data(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test verifying signature of empty data."""
        data = b""
        signature = sign_data(ecdsa_private_key, data)
        assert verify_signature(ecdsa_public_key, data, signature) is True


class TestGetPublicKeyFromPrivate:
    """Tests for get_public_key_from_private function."""

    def test_extracts_public_key(self, ecdsa_private_key: ec.EllipticCurvePrivateKey) -> None:
        """Test extracting public key from private key."""
        public_key = get_public_key_from_private(ecdsa_private_key)

        assert isinstance(public_key, ec.EllipticCurvePublicKey)

    def test_public_key_can_verify(self, ecdsa_private_key: ec.EllipticCurvePrivateKey) -> None:
        """Test that extracted public key can verify signatures."""
        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        public_key = get_public_key_from_private(ecdsa_private_key)
        assert verify_signature(public_key, data, signature) is True

    def test_same_public_key_each_time(self, ecdsa_private_key: ec.EllipticCurvePrivateKey) -> None:
        """Test that same public key is extracted each time."""
        pk1 = get_public_key_from_private(ecdsa_private_key)
        pk2 = get_public_key_from_private(ecdsa_private_key)

        # Compare by signing and verifying
        data = b"test"
        sig = sign_data(ecdsa_private_key, data)
        assert verify_signature(pk1, data, sig) is True
        assert verify_signature(pk2, data, sig) is True


class TestSignAndVerifyRoundTrip:
    """Integration tests for sign and verify round-trip."""

    def test_sign_verify_json_payload(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test signing and verifying JSON payload (common use case)."""
        import json

        payload = json.dumps(
            {
                "hashes": [
                    "abc123",
                    "def456",
                    "ghi789",
                ]
            }
        ).encode("utf-8")

        signature = sign_data(ecdsa_private_key, payload)
        assert verify_signature(ecdsa_public_key, payload, signature) is True

    def test_sign_verify_hash_array(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test signing array of hashes (request approval use case)."""
        import json

        hashes = ["hash1", "hash2", "hash3"]
        hashes.sort()  # Sorted as per SDK pattern
        payload = json.dumps(hashes).encode("utf-8")

        signature = sign_data(ecdsa_private_key, payload)
        assert verify_signature(ecdsa_public_key, payload, signature) is True

    def test_multiple_signatures_all_valid(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test that multiple signatures of same data are all valid."""
        data = b"test data"

        # ECDSA produces different signatures each time due to random k
        signatures = [sign_data(ecdsa_private_key, data) for _ in range(5)]

        # All should be different but all valid
        assert len(set(signatures)) == 5  # All different
        for sig in signatures:
            assert verify_signature(ecdsa_public_key, data, sig) is True
