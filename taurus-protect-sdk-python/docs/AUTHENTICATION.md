# Authentication

This document describes the authentication mechanisms in the Taurus-PROTECT Python SDK.

## Overview

The SDK implements a multi-layer security model:

1. **TPV1 Authentication** - HMAC-based API request signing for all API calls
2. **SuperAdmin Verification** - ECDSA signature verification for governance rules
3. **Request Approval** - ECDSA signing for transaction approval
4. **Data Integrity** - SHA-256 hash verification for payloads

## TPV1-HMAC-SHA256 Scheme

All API requests are automatically signed using the TPV1 (Taurus Protocol Version 1) scheme.

### Authorization Header Format

```
Authorization: TPV1-HMAC-SHA256 ApiKey=<api_key> Nonce=<uuid> Timestamp=<unix_ms> Signature=<base64_sig>
```

| Component | Description |
|-----------|-------------|
| `ApiKey` | UUID-format API key identifying the caller |
| `Nonce` | Random UUID for replay prevention |
| `Timestamp` | Unix timestamp in milliseconds |
| `Signature` | Base64-encoded HMAC-SHA256 signature |

### Signature Computation

The signature is computed over a canonical message string:

```
TPV1 <ApiKey> <Nonce> <Timestamp> <Method> <Host> <Path> [Query] [ContentType] [Body]
```

Optional components are only included if present.

### Implementation

The SDK handles TPV1 signing automatically through the `TPV1Auth` class:

```python
from taurus_protect.crypto.tpv1 import TPV1Auth

# Create auth handler
auth = TPV1Auth(api_key="your-api-key", api_secret_hex="your-secret-hex")

# Sign a request (called automatically by the SDK)
header = auth.sign_request(
    method="POST",
    host="api.protect.taurushq.com",
    path="/api/rest/v1/wallets",
    content_type="application/json",
    body='{"name": "My Wallet"}',
)
# Returns: "TPV1-HMAC-SHA256 ApiKey=... Nonce=... Timestamp=... Signature=..."

# Close to securely wipe credentials
auth.close()
```

## Client Initialization

### Basic Initialization

```python
from taurus_protect import ProtectClient

with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
) as client:
    # API requests are automatically signed
    wallets, _ = client.wallets.list()
```

### Initialization Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `host` | `str` | Required | API host URL |
| `api_key` | `str` | Required | API key for authentication |
| `api_secret` | `str` | Required | API secret as hex-encoded string |
| `super_admin_keys_pem` | `List[str]` | `None` | List of PEM-encoded SuperAdmin public keys |
| `min_valid_signatures` | `int` | `1` | Minimum valid signatures required |
| `rules_cache_ttl` | `float` | `300.0` | Rules container cache TTL in seconds |
| `timeout` | `float` | `30.0` | HTTP request timeout in seconds |

### With SuperAdmin Keys

For production use, provide SuperAdmin public keys to enable governance rule verification:

```python
from taurus_protect import ProtectClient

super_admin_keys = [
    """-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----""",
    """-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----""",
]

with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
    super_admin_keys_pem=super_admin_keys,
    min_valid_signatures=2,  # Require 2 valid signatures
) as client:
    # Governance rule verification will require 2 valid SuperAdmin signatures
    pass
```

### Environment-Based Configuration

```python
import os
from taurus_protect import ProtectClient

def create_client() -> ProtectClient:
    """Create client from environment variables."""
    host = os.environ["PROTECT_API_HOST"]
    api_key = os.environ["PROTECT_API_KEY"]
    api_secret = os.environ["PROTECT_API_SECRET"]

    # Optional: Load SuperAdmin keys from file
    super_admin_keys = None
    keys_path = os.environ.get("PROTECT_SUPER_ADMIN_KEYS_PATH")
    if keys_path:
        with open(keys_path) as f:
            super_admin_keys = [key.strip() for key in f.read().split("-----END PUBLIC KEY-----")
                               if key.strip()]
            super_admin_keys = [k + "-----END PUBLIC KEY-----" for k in super_admin_keys]

    return ProtectClient.create(
        host=host,
        api_key=api_key,
        api_secret=api_secret,
        super_admin_keys_pem=super_admin_keys,
        min_valid_signatures=int(os.environ.get("PROTECT_MIN_SIGNATURES", "1")),
    )
```

## Cryptographic Operations

### ECDSA Signing

The SDK uses ECDSA P-256 (secp256r1) for request approval and signature verification:

```python
from taurus_protect.crypto.signing import sign_data, verify_signature
from cryptography.hazmat.primitives.serialization import load_pem_private_key

# Load private key
with open("private_key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=None)

# Sign data
data = b'["hash1", "hash2"]'
signature = sign_data(private_key, data)
# Returns base64-encoded raw r||s signature

# Verify signature
public_key = private_key.public_key()
is_valid = verify_signature(public_key, data, signature)
```

### SHA-256 Hashing

```python
from taurus_protect.crypto.hashing import calculate_hex_hash

# Compute SHA-256 hash
payload = '{"from_address_id": "123", "to_address_id": "456", "amount": "1000"}'
hash_hex = calculate_hex_hash(payload)
# Returns: "a1b2c3d4..."
```

### Constant-Time Comparison

To prevent timing attacks, all hash comparisons use constant-time algorithms:

```python
from taurus_protect.helpers.constant_time import constant_time_compare

# Safe comparison (prevents timing attacks)
is_equal = constant_time_compare(computed_hash, provided_hash)
```

## Request Approval Signing

When approving transaction requests, the SDK signs the request hashes with ECDSA:

### Approval Flow

1. Retrieve requests to approve
2. Sort requests by ID (numeric order)
3. Build JSON array of hashes
4. Sign the JSON array with ECDSA P-256
5. Submit signature to API

### Implementation

```python
from cryptography.hazmat.primitives.serialization import load_pem_private_key
from taurus_protect import ProtectClient

# Load your approval private key
with open("approval_key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=None)

with ProtectClient.create(host, api_key, api_secret) as client:
    # Get requests pending approval
    requests, _ = client.requests.get_for_approval(limit=10)

    if requests:
        # Approve with ECDSA signature
        # Internally:
        # 1. Sorts requests by ID
        # 2. Builds: '["hash1", "hash2", ...]'
        # 3. Signs with ECDSA
        signed_count = client.requests.approve_requests(
            requests,
            private_key,
            comment="Batch approval via SDK",
        )
        print(f"Approved {signed_count} request(s)")
```

### What Gets Signed

The SDK signs a JSON array of metadata hashes:

```python
# Internal signing process:
sorted_requests = sorted(requests, key=lambda r: int(r.id))
hashes = [r.metadata.hash for r in sorted_requests]
to_sign = json.dumps(hashes)  # '["abc123...", "def456..."]'
signature = sign_data(private_key, to_sign.encode("utf-8"))
```

## SuperAdmin Key Configuration

### Key Format

SuperAdmin keys are ECDSA P-256 public keys in PEM format:

```
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----
```

### Verification Flow

When SuperAdmin keys are configured, the SDK verifies governance rules:

```
┌─────────────────────────────────────────────────────────────┐
│                 Governance Rules Verification                 │
├─────────────────────────────────────────────────────────────┤
│  1. Fetch rules container (base64-encoded)                  │
│  2. Decode base64 to raw bytes                              │
│  3. For each signature:                                     │
│     a. Decode base64 signature                              │
│     b. Try verification against each SuperAdmin key         │
│     c. Count if valid                                       │
│  4. Require: validCount >= minValidSignatures               │
└─────────────────────────────────────────────────────────────┘
```

### Configuration Example

```python
from taurus_protect import ProtectClient

# Production configuration with 2-of-3 SuperAdmin verification
super_admin_keys = [
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
]

client = ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key=api_key,
    api_secret=api_secret,
    super_admin_keys_pem=super_admin_keys,
    min_valid_signatures=2,  # 2 of 3 required
)
```

## Request Hash Verification

All requests fetched via `RequestService.get()` have their hash verified:

```python
def _verify_request_hash(self, request: Request) -> None:
    """Verify the hash using constant-time comparison."""
    if request.metadata is None:
        return

    payload = request.metadata.payload_as_string
    provided_hash = request.metadata.hash

    if not payload:
        return

    # Compute hash
    computed_hash = calculate_hex_hash(payload)

    # Constant-time comparison (prevents timing attacks)
    if not hmac.compare_digest(computed_hash, provided_hash or ""):
        raise IntegrityError("request hash verification failed")
```

## Security Best Practices

### Credential Storage

```python
# DO: Use environment variables
api_key = os.environ["PROTECT_API_KEY"]
api_secret = os.environ["PROTECT_API_SECRET"]

# DO: Use secret managers in production
from azure.keyvault.secrets import SecretClient
# or
from google.cloud import secretmanager

# DON'T: Hardcode credentials
api_key = "abc123"  # Never do this!
```

### Private Key Handling

```python
# DO: Load from secure storage
with open("private_key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=b"passphrase")

# DO: Use HSM for production
# Configure your HSM provider's PKCS#11 interface

# DON'T: Store keys in code
key_pem = "-----BEGIN PRIVATE KEY-----..."  # Never do this!
```

### Resource Cleanup

```python
# DO: Use context manager for automatic cleanup
with ProtectClient.create(...) as client:
    # Credentials are securely wiped when exiting the block
    pass

# OR: Explicitly close
client = ProtectClient.create(...)
try:
    # Use client
    pass
finally:
    client.close()  # Securely wipes credentials
```

## Error Handling

### Authentication Errors

| Error | HTTP Code | Cause | Resolution |
|-------|-----------|-------|------------|
| `AuthenticationError` | 401 | Invalid API key or signature | Verify credentials |
| `AuthorizationError` | 403 | Valid auth but insufficient permissions | Check API key permissions |

### Integrity Errors

| Error | Cause | Resolution |
|-------|-------|------------|
| `IntegrityError` | Hash mismatch | Data may be tampered - investigate |
| `IntegrityError` | Insufficient SuperAdmin signatures | Check key configuration |
| `WhitelistError` | User signature verification failed | Invalid approver signature |
| `ConfigurationError` | Invalid SDK configuration | Check initialization parameters |

### Error Type Checking

```python
from taurus_protect.errors import ApiException, IntegrityError, WhitelistError

# Check error type and retryability
try:
    wallet = client.wallets.get(123)
except ApiException as e:
    if e.is_retryable():
        # True for 429 (rate limit) and 5xx (server errors)
        # Implement retry with backoff
        pass
    if e.code == 401:
        # Authentication failure - check credentials
        pass
    if e.code == 403:
        # Authorization failure - check permissions
        pass
except IntegrityError as e:
    # Cryptographic verification failed - do NOT retry
    # Investigate potential data tampering
    pass
except WhitelistError as e:
    # Whitelist-specific verification failure
    pass
```

### Example Error Handling

```python
from taurus_protect.errors import (
    AuthenticationError,
    AuthorizationError,
    IntegrityError,
    ConfigurationError,
)

try:
    with ProtectClient.create(host, api_key, api_secret) as client:
        wallet = client.wallets.get(123)
except ConfigurationError as e:
    print(f"Invalid configuration: {e.message}")
except AuthenticationError as e:
    print(f"Authentication failed: {e.message}")
except AuthorizationError as e:
    print(f"Permission denied: {e.message}")
except IntegrityError as e:
    # SECURITY: Do not retry - investigate the cause
    print(f"Integrity verification failed: {e.message}")
```

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and design patterns
- [Services Reference](SERVICES.md) - Complete API documentation
- [Key Concepts](CONCEPTS.md) - Domain model and exceptions
- [Common Authentication](../../docs/AUTHENTICATION.md) - Shared authentication concepts
