# SDK Overview

This document provides an architectural overview of the Taurus-PROTECT Python SDK.

## Introduction

The Python SDK provides a type-safe, Pythonic interface to the Taurus-PROTECT API. It is designed to be:

- **Simple** - Clean API with sensible defaults
- **Type-safe** - Full type annotations for IDE support and static analysis
- **Secure** - Built-in cryptographic verification and secure credential handling
- **Thread-safe** - Safe for use in multi-threaded applications

## Package Structure

```
taurus_protect/
├── __init__.py                 # Public exports (ProtectClient, models, errors)
├── client.py                   # ProtectClient - main entry point
├── errors.py                   # Exception hierarchy
│
├── services/                   # Service layer (38 core services)
│   ├── __init__.py
│   ├── _base.py               # BaseService with common functionality
│   ├── wallet_service.py
│   ├── address_service.py
│   ├── request_service.py
│   ├── ...
│   └── taurus_network/        # TaurusNetwork namespace (5 services)
│       ├── __init__.py
│       ├── _client.py         # TaurusNetworkClient namespace
│       ├── participant_service.py
│       ├── pledge_service.py
│       ├── lending_service.py
│       ├── settlement_service.py
│       └── sharing_service.py
│
├── models/                     # Pydantic domain models
│   ├── __init__.py
│   ├── wallet.py
│   ├── address.py
│   ├── request.py
│   ├── transaction.py
│   ├── balance.py
│   ├── governance_rules.py
│   ├── whitelisted_address.py
│   ├── pagination.py
│   ├── ...
│   └── taurus_network/        # TaurusNetwork models (71 models)
│       ├── __init__.py
│       ├── participant.py     # 7 models
│       ├── pledge.py          # 26 models
│       ├── lending.py         # 13 models
│       ├── settlement.py      # 11 models
│       └── sharing.py         # 14 models
│
├── mappers/                    # DTO to model conversion
│   ├── __init__.py
│   ├── _base.py               # Safe conversion helpers
│   ├── wallet.py
│   ├── address.py
│   ├── request.py
│   └── ...
│
├── helpers/                    # Verification utilities
│   ├── __init__.py
│   ├── constant_time.py       # Timing-attack-safe comparison
│   ├── signature_verifier.py  # SuperAdmin signature verification
│   ├── address_signature_verifier.py
│   ├── whitelisted_address_verifier.py  # 6-step address verification
│   ├── whitelisted_asset_verifier.py
│   ├── whitelist_integrity_helper.py    # Whitelist integrity checks
│   └── whitelist_hash_helper.py         # Legacy hash computation
│
├── cache/                      # Caching infrastructure
│   ├── __init__.py
│   └── rules_container_cache.py  # Thread-safe rules caching
│
├── crypto/                     # Cryptographic utilities
│   ├── __init__.py
│   ├── tpv1.py               # TPV1-HMAC-SHA256 authentication
│   ├── signing.py            # ECDSA signing and verification
│   ├── hashing.py            # SHA-256 hash computation
│   ├── keys.py               # Key encoding/decoding
│   └── authenticated_rest.py # TPV1-signed REST client
│
└── _internal/                  # Auto-generated code (DO NOT MODIFY)
    ├── openapi/               # OpenAPI-generated client
    └── proto/                 # Protobuf-generated classes
```

## Module Descriptions

### services/

The service layer provides the primary API for interacting with Taurus-PROTECT. Each service wraps the auto-generated OpenAPI client with:

- Type-safe method signatures with Pydantic models
- Error mapping to domain-specific exceptions
- Automatic pagination handling
- Security verification (hash checks, signature validation)

**Key services:**

| Service | Purpose |
|---------|---------|
| `WalletService` | Wallet CRUD and balance history |
| `AddressService` | Address management with signature verification |
| `RequestService` | Transaction requests with ECDSA approval signing |
| `GovernanceRuleService` | Governance rules with SuperAdmin verification |
| `WhitelistedAddressService` | Whitelisted addresses with 6-step verification |
| `WhitelistedAssetService` | Whitelisted assets with 5-step verification |

### models/

Domain models are implemented as frozen Pydantic models, providing:

- Immutability (models are read-only after creation)
- Automatic validation
- JSON serialization/deserialization
- Rich type annotations with descriptions

**Example model:**

```python
from pydantic import BaseModel, Field

class Wallet(BaseModel):
    """Represents a cryptocurrency wallet."""

    model_config = {"frozen": True}

    id: str = Field(..., description="Unique wallet identifier")
    name: str = Field(..., description="Human-readable name")
    blockchain: str = Field(..., description="Blockchain type (ETH, BTC, etc.)")
    network: str = Field(..., description="Network (mainnet, testnet)")
    currency: str = Field(..., description="Primary currency identifier")
    is_omnibus: bool = Field(default=False, description="Whether this is an omnibus wallet")
```

### mappers/

Mapper functions convert between OpenAPI DTOs and domain models. They use safe conversion helpers that handle None values gracefully:

```python
from taurus_protect.mappers._base import safe_string, safe_int, safe_datetime

def wallet_from_dto(dto) -> Optional[Wallet]:
    if dto is None:
        return None
    return Wallet(
        id=safe_int(getattr(dto, "id", None), 0),
        name=safe_string(getattr(dto, "name", None), ""),
        blockchain=safe_string(getattr(dto, "blockchain", None), ""),
        # ...
    )
```

### helpers/

Verification utilities implement security-critical operations:

- **constant_time.py** - Timing-attack-safe string/bytes comparison using `hmac.compare_digest()`
- **signature_verifier.py** - SuperAdmin ECDSA signature verification
- **address_signature_verifier.py** - Address signature verification using HSM slot public key
- **whitelisted_asset_verifier.py** - 5-step verification for whitelisted assets

### cache/

The `RulesContainerCache` provides thread-safe caching of decoded governance rules:

- Configurable TTL (default 5 minutes)
- Thread-safe with `threading.RLock()`
- Auto-refreshes from GovernanceRuleService when expired
- Used by AddressService for HSM key lookup during signature verification

### crypto/

Cryptographic utilities for authentication and signing:

| Module | Purpose |
|--------|---------|
| `tpv1.py` | TPV1-HMAC-SHA256 request signing |
| `signing.py` | ECDSA P-256 signing and verification |
| `hashing.py` | SHA-256 hash computation |
| `keys.py` | PEM key encoding/decoding |

## Key Design Patterns

### Facade Pattern (ProtectClient)

The `ProtectClient` class serves as the main entry point, providing a unified interface to all services:

```python
with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
) as client:
    # Access services through properties
    wallets, _ = client.wallets.list()
    addresses, _ = client.addresses.list(wallet_id=123)
    requests, _ = client.requests.get_for_approval()
```

### Lazy Initialization with RLock

Services are lazily initialized on first access, with thread-safe locking:

```python
class ProtectClient:
    def __init__(self, ...):
        self._lock = threading.RLock()
        self._wallet_service: Optional[WalletService] = None

    @property
    def wallets(self) -> WalletService:
        self._check_not_closed()
        with self._lock:
            if self._wallet_service is None:
                self._wallet_service = WalletService(...)
            return self._wallet_service
```

Benefits:
- Services only created when needed
- Thread-safe initialization
- Single instance per client

### Context Manager for Resource Cleanup

The client implements the context manager protocol for proper resource cleanup:

```python
def __enter__(self) -> "ProtectClient":
    return self

def __exit__(self, exc_type, exc_val, exc_tb) -> None:
    self.close()

def close(self) -> None:
    """Close the client and securely wipe credentials."""
    with self._lock:
        if not self._closed:
            if self._auth:
                self._auth.close()  # Wipes API secret from memory
            self._closed = True
```

### Pydantic Frozen Models

All domain models are immutable (frozen) Pydantic models:

```python
class Wallet(BaseModel):
    model_config = {"frozen": True}  # Makes model immutable

    id: int
    name: str
    # ...
```

Benefits:
- Thread safety (no shared mutable state)
- Hashable (can be used in sets/dict keys)
- Clear data contracts

### Error Hierarchy with Retry Helpers

The exception hierarchy provides rich error information with retry guidance:

```python
class APIError(Exception):
    def is_retryable(self) -> bool:
        """True for 429 and 5xx errors."""
        return self.code == 429 or self.code >= 500

    def suggested_retry_delay(self) -> timedelta:
        """Recommended wait time before retry."""
        if self.code == 429:
            return self.retry_after or timedelta(seconds=1)
        if self.code >= 500:
            return timedelta(seconds=5)
        return timedelta(0)
```

## Key Classes Reference

| Class | Module | Purpose |
|-------|--------|---------|
| `ProtectClient` | `client.py` | Main SDK entry point |
| `TPV1Auth` | `crypto/tpv1.py` | Request authentication |
| `Wallet`, `Address`, `Request` | `models/` | Core domain models |
| `APIError`, `IntegrityError` | `errors.py` | Exception types |
| `RulesContainerCache` | `cache/` | Governance rules caching |
| `WhitelistedAssetVerifier` | `helpers/` | 5-step asset verification |

## Build and Development

### Build Commands

```bash
# Install and test (default)
./build.sh

# Build package
./build.sh build

# Run tests
./build.sh test

# Run linters
./build.sh lint

# Format code
./build.sh format

# Code generation
./build.sh generate
```

### Code Generation

The SDK uses code generation for the OpenAPI client and protobuf classes:

```bash
# Generate OpenAPI client (requires Java 11+)
./build.sh generate
```

**Important:** Never modify files in `_internal/openapi/` or `_internal/proto/` - they are auto-generated.

### Type Checking

The SDK uses strict mypy type checking:

```bash
mypy taurus_protect --exclude='taurus_protect/_internal'
```

## Related Documentation

- [Authentication](AUTHENTICATION.md) - TPV1 authentication details
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
- [Key Concepts](CONCEPTS.md) - Domain model and exceptions
- [Common Concepts](../../docs/CONCEPTS.md) - Shared domain concepts
