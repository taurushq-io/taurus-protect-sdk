"""Pytest fixtures for Taurus-PROTECT SDK tests."""

from __future__ import annotations

from datetime import datetime, timedelta, timezone

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

# =============================================================================
# API Credentials Fixtures
# =============================================================================


@pytest.fixture
def api_key() -> str:
    """Sample API key for testing."""
    return "test-api-key"


@pytest.fixture
def api_secret_hex() -> str:
    """Sample hex-encoded API secret for testing (32 bytes = 64 hex chars)."""
    return "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"


@pytest.fixture
def host() -> str:
    """Sample API host for testing."""
    return "https://api.test.taurushq.com"


# =============================================================================
# Crypto Key Fixtures
# =============================================================================


@pytest.fixture
def ecdsa_private_key() -> ec.EllipticCurvePrivateKey:
    """Generate a fresh ECDSA P-256 private key for testing."""
    return ec.generate_private_key(ec.SECP256R1())


@pytest.fixture
def ecdsa_public_key(ecdsa_private_key: ec.EllipticCurvePrivateKey) -> ec.EllipticCurvePublicKey:
    """Get public key from the test private key."""
    return ecdsa_private_key.public_key()


@pytest.fixture
def ecdsa_private_key_pem(ecdsa_private_key: ec.EllipticCurvePrivateKey) -> str:
    """PEM-encoded private key for testing."""
    pem_bytes = ecdsa_private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption(),
    )
    return pem_bytes.decode("utf-8")


@pytest.fixture
def ecdsa_public_key_pem(ecdsa_public_key: ec.EllipticCurvePublicKey) -> str:
    """PEM-encoded public key for testing."""
    pem_bytes = ecdsa_public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return pem_bytes.decode("utf-8")


@pytest.fixture
def super_admin_keys_pem(ecdsa_public_key_pem: str) -> list:
    """List containing a single SuperAdmin PEM key for tests requiring mandatory keys."""
    return [ecdsa_public_key_pem]


@pytest.fixture
def second_ecdsa_private_key() -> ec.EllipticCurvePrivateKey:
    """Second ECDSA P-256 private key for multi-signature tests."""
    return ec.generate_private_key(ec.SECP256R1())


@pytest.fixture
def second_ecdsa_public_key(
    second_ecdsa_private_key: ec.EllipticCurvePrivateKey,
) -> ec.EllipticCurvePublicKey:
    """Second public key for multi-signature tests."""
    return second_ecdsa_private_key.public_key()


# =============================================================================
# Sample Data Fixtures
# =============================================================================


@pytest.fixture
def sample_wallet_dto() -> dict:
    """Sample wallet DTO for mapper tests."""
    return {
        "id": "123",
        "tenant_id": "tenant-456",
        "name": "Test Wallet",
        "wallet_type": "STANDARD",
        "status": "ACTIVE",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T12:00:00Z",
    }


@pytest.fixture
def sample_address_dto() -> dict:
    """Sample address DTO for mapper tests."""
    return {
        "id": 789,
        "wallet_id": "123",
        "address": "0x1234567890abcdef1234567890abcdef12345678",
        "blockchain": "ETH",
        "network": "mainnet",
        "path": "m/44'/60'/0'/0/0",
        "status": "ACTIVE",
        "created_at": "2024-01-15T10:30:00Z",
    }


@pytest.fixture
def sample_request_dto() -> dict:
    """Sample request DTO for mapper tests."""
    return {
        "id": 456,
        "request_type": "EXTERNAL_TRANSFER",
        "status": "PENDING_APPROVAL",
        "hash": "abc123def456",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T12:00:00Z",
    }


@pytest.fixture
def sample_transaction_dto() -> dict:
    """Sample transaction DTO for mapper tests."""
    return {
        "id": "tx-123",
        "request_id": 456,
        "tx_hash": "0xabcdef123456",
        "blockchain": "ETH",
        "network": "mainnet",
        "status": "CONFIRMED",
        "amount": "1000000000000000000",
        "created_at": "2024-01-15T10:30:00Z",
    }


@pytest.fixture
def sample_user_dto() -> dict:
    """Sample user DTO for mapper tests."""
    return {
        "id": "user-123",
        "external_user_id": "ext-user-456",
        "tenant_id": "tenant-789",
        "username": "testuser",
        "email": "test@example.com",
        "first_name": "Test",
        "last_name": "User",
        "status": "ACTIVE",
        "roles": ["ADMIN", "VIEWER"],
        "created_at": "2024-01-15T10:30:00Z",
    }


# =============================================================================
# Datetime Fixtures
# =============================================================================


@pytest.fixture
def sample_datetime() -> datetime:
    """Sample datetime for testing."""
    return datetime(2024, 1, 15, 10, 30, 0, tzinfo=timezone.utc)


@pytest.fixture
def sample_datetime_str() -> str:
    """Sample ISO datetime string for testing."""
    return "2024-01-15T10:30:00Z"


@pytest.fixture
def sample_timedelta() -> timedelta:
    """Sample timedelta for testing."""
    return timedelta(seconds=30)


# =============================================================================
# Governance Rules Fixtures
# =============================================================================


@pytest.fixture
def sample_governance_rules_dto() -> dict:
    """Sample governance rules DTO for verification tests."""
    return {
        "id": "rules-123",
        "tenant_id": "tenant-456",
        "rules_container": None,  # Will be set by specific tests
        "rules_signatures": [],
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T12:00:00Z",
    }


# =============================================================================
# Pagination Fixtures
# =============================================================================


@pytest.fixture
def sample_pagination_dto() -> dict:
    """Sample pagination DTO for service tests."""
    return {
        "total_items": "100",
        "offset": "0",
        "limit": "50",
    }


@pytest.fixture
def sample_cursor_pagination_dto() -> dict:
    """Sample cursor-based pagination DTO for TaurusNetwork tests."""
    return {
        "current_page": "eyJwYWdlIjoxfQ==",  # Base64 encoded cursor
        "has_previous": False,
        "has_next": True,
    }
