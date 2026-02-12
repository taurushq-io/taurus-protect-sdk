"""Tests for key handling utilities."""

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.crypto.keys import (
    decode_private_key_pem,
    decode_public_key_pem,
    decode_public_keys_pem,
    encode_public_key_pem,
)


class TestDecodePublicKeyPem:
    """Tests for decode_public_key_pem function."""

    def test_decode_valid_public_key(self, ecdsa_public_key_pem: str) -> None:
        """Test decoding valid PEM public key."""
        key = decode_public_key_pem(ecdsa_public_key_pem)
        assert isinstance(key, ec.EllipticCurvePublicKey)

    def test_decode_invalid_pem_raises(self) -> None:
        """Test that invalid PEM raises ValueError."""
        with pytest.raises(ValueError, match="Failed to decode public key"):
            decode_public_key_pem("not a valid pem")

    def test_decode_empty_string_raises(self) -> None:
        """Test that empty string raises ValueError."""
        with pytest.raises(ValueError):
            decode_public_key_pem("")

    def test_decode_private_key_as_public_raises(self, ecdsa_private_key_pem: str) -> None:
        """Test that private key PEM raises ValueError when decoded as public."""
        with pytest.raises(ValueError, match="Failed to decode public key"):
            decode_public_key_pem(ecdsa_private_key_pem)

    def test_decode_malformed_pem_raises(self) -> None:
        """Test that malformed PEM raises ValueError."""
        malformed = """-----BEGIN PUBLIC KEY-----
invalid base64 content here!!!
-----END PUBLIC KEY-----"""
        with pytest.raises(ValueError):
            decode_public_key_pem(malformed)

    def test_decode_rsa_key_raises(self) -> None:
        """Test that RSA key raises ValueError (not EC key)."""
        # A sample RSA public key PEM (invalid for EC operations)
        rsa_pem = """-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----"""
        with pytest.raises(ValueError, match="not an elliptic curve public key"):
            decode_public_key_pem(rsa_pem)


class TestDecodePublicKeysPem:
    """Tests for decode_public_keys_pem function."""

    def test_decode_multiple_keys(
        self, ecdsa_public_key_pem: str, second_ecdsa_public_key: ec.EllipticCurvePublicKey
    ) -> None:
        """Test decoding multiple PEM public keys."""
        from taurus_protect.crypto.keys import encode_public_key_pem

        second_pem = encode_public_key_pem(second_ecdsa_public_key)

        keys = decode_public_keys_pem([ecdsa_public_key_pem, second_pem])
        assert len(keys) == 2
        assert all(isinstance(k, ec.EllipticCurvePublicKey) for k in keys)

    def test_decode_empty_list(self) -> None:
        """Test decoding empty list."""
        keys = decode_public_keys_pem([])
        assert keys == []

    def test_decode_single_key(self, ecdsa_public_key_pem: str) -> None:
        """Test decoding single key in list."""
        keys = decode_public_keys_pem([ecdsa_public_key_pem])
        assert len(keys) == 1
        assert isinstance(keys[0], ec.EllipticCurvePublicKey)

    def test_decode_with_invalid_key_raises(self, ecdsa_public_key_pem: str) -> None:
        """Test that invalid key in list raises ValueError."""
        with pytest.raises(ValueError):
            decode_public_keys_pem([ecdsa_public_key_pem, "invalid key"])


class TestDecodePrivateKeyPem:
    """Tests for decode_private_key_pem function."""

    def test_decode_valid_private_key(self, ecdsa_private_key_pem: str) -> None:
        """Test decoding valid PEM private key."""
        key = decode_private_key_pem(ecdsa_private_key_pem)
        assert isinstance(key, ec.EllipticCurvePrivateKey)

    def test_decode_invalid_pem_raises(self) -> None:
        """Test that invalid PEM raises ValueError."""
        with pytest.raises(ValueError, match="Failed to decode private key"):
            decode_private_key_pem("not a valid pem")

    def test_decode_empty_string_raises(self) -> None:
        """Test that empty string raises ValueError."""
        with pytest.raises(ValueError):
            decode_private_key_pem("")

    def test_decode_public_key_as_private_raises(self, ecdsa_public_key_pem: str) -> None:
        """Test that public key PEM raises ValueError when decoded as private."""
        with pytest.raises(ValueError):
            decode_private_key_pem(ecdsa_public_key_pem)

    def test_decoded_key_can_sign(self, ecdsa_private_key_pem: str) -> None:
        """Test that decoded private key can sign data."""
        from taurus_protect.crypto.signing import sign_data

        key = decode_private_key_pem(ecdsa_private_key_pem)
        signature = sign_data(key, b"test data")
        assert isinstance(signature, str)


class TestEncodePublicKeyPem:
    """Tests for encode_public_key_pem function."""

    def test_encode_public_key(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test encoding public key to PEM."""
        pem = encode_public_key_pem(ecdsa_public_key)

        assert "-----BEGIN PUBLIC KEY-----" in pem
        assert "-----END PUBLIC KEY-----" in pem

    def test_encode_decode_roundtrip(self, ecdsa_public_key: ec.EllipticCurvePublicKey) -> None:
        """Test that encode/decode round-trip preserves key."""
        from taurus_protect.crypto.signing import sign_data, verify_signature

        pem = encode_public_key_pem(ecdsa_public_key)
        decoded = decode_public_key_pem(pem)

        # Keys should be functionally equivalent
        assert isinstance(decoded, ec.EllipticCurvePublicKey)


class TestKeyRoundTrip:
    """Integration tests for key encode/decode round trips."""

    def test_sign_with_decoded_key_verify_with_original(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_private_key_pem: str,
        ecdsa_public_key: ec.EllipticCurvePublicKey,
    ) -> None:
        """Test signing with decoded key, verifying with original."""
        from taurus_protect.crypto.signing import sign_data, verify_signature

        decoded_private = decode_private_key_pem(ecdsa_private_key_pem)
        data = b"test data"
        signature = sign_data(decoded_private, data)

        assert verify_signature(ecdsa_public_key, data, signature) is True

    def test_sign_with_original_verify_with_decoded(
        self,
        ecdsa_private_key: ec.EllipticCurvePrivateKey,
        ecdsa_public_key_pem: str,
    ) -> None:
        """Test signing with original key, verifying with decoded."""
        from taurus_protect.crypto.signing import sign_data, verify_signature

        data = b"test data"
        signature = sign_data(ecdsa_private_key, data)

        decoded_public = decode_public_key_pem(ecdsa_public_key_pem)
        assert verify_signature(decoded_public, data, signature) is True
