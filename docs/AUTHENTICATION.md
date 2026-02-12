# Authentication & Security

This document describes the authentication and security model used by the Taurus-PROTECT API. These concepts are shared across all SDK implementations.

## Overview

The Taurus-PROTECT platform uses a multi-layer security model:

1. **TPV1 Authentication** - HMAC-based API request signing
2. **SuperAdmin Verification** - ECDSA signature verification for governance rules
3. **Data Integrity** - SHA-256 hash verification for payloads

---

## TPV1 Authentication Scheme

All API requests are signed using the TPV1 (Taurus Protocol Version 1) scheme. This provides:

- **Authentication** - Proves the request came from a valid API key holder
- **Integrity** - Ensures the request hasn't been tampered with
- **Replay Prevention** - Timestamp and nonce prevent replay attacks

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
TPV1 <ApiKey> <Nonce> <Timestamp> <Method> <Host> <Path> <Query> <ContentType> <Body>
```

| Field | Description |
|-------|-------------|
| `Method` | HTTP method (GET, POST, PUT, DELETE) |
| `Host` | API host (without port for standard ports) |
| `Path` | Request path (e.g., `/api/rest/v1/wallets`) |
| `Query` | Query string (empty string if none) |
| `ContentType` | Content-Type header (normalized for JSON) |
| `Body` | Request body (empty string if none) |

### Algorithm

1. Construct the canonical message string
2. Compute HMAC-SHA256 using the API secret as the key
3. Base64-encode the result
4. Include in the Authorization header

### Automatic Handling

All SDKs handle TPV1 signing automatically. You only need to provide credentials during client initialization.

---

## API Credentials

### Required Credentials

| Credential | Format | Description |
|------------|--------|-------------|
| API Key | UUID | Identifies the API caller |
| API Secret | Hex string | Secret key for HMAC signing |

### Obtaining Credentials

API credentials are provisioned through the Taurus-PROTECT administration interface. Contact your administrator to obtain credentials for your application.

### Credential Security

- Store secrets securely (environment variables, secret managers)
- Never commit credentials to version control
- Rotate API keys periodically
- Use the minimum required permissions

---

## SuperAdmin Public Keys

### Purpose

SuperAdmin public keys are used to verify:

1. **Governance Rules** - Rules defining approval thresholds and user groups
2. **Whitelisted Addresses** - Cryptographically signed address whitelists
3. **Configuration Changes** - System-level configuration modifications

### Key Format

SuperAdmin keys are ECDSA P-256 (secp256r1) public keys in PEM format:

```
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----
```

### Configuration

When initializing an SDK client, you can provide SuperAdmin public keys to enable automatic verification:

- **Multiple keys** - Provide all known SuperAdmin public keys
- **Minimum signatures** - Specify the threshold for valid signatures
- **Cache TTL** - Configure how long verified rules are cached

### Verification Flow

```
┌─────────────────────────────────────────────────────────┐
│                 Governance Rules                         │
│  ┌─────────────────┐                                    │
│  │ rulesContainer  │  (Base64-encoded protobuf)        │
│  │ signatures[]    │  (List of ECDSA signatures)       │
│  └─────────────────┘                                    │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│            Signature Verification                        │
│  For each signature:                                     │
│    1. Decode base64 signature                           │
│    2. Verify against SuperAdmin public keys             │
│    3. Count valid signatures                            │
│  Require: validCount >= minValidSignatures              │
└─────────────────────────────────────────────────────────┘
```

---

## Cryptographic Primitives

### Algorithms Used

| Operation | Algorithm | Notes |
|-----------|-----------|-------|
| API Signing | HMAC-SHA256 | TPV1 request authentication |
| Hash | SHA-256 | Payload integrity verification |
| Signatures | ECDSA P-256 | SuperAdmin and user signatures |
| Key Encoding | PEM / Base64 | Key transport format |

### Signature Format

ECDSA signatures are encoded as Base64 raw `r||s` format (plain ECDSA, not DER-encoded).

### Constant-Time Comparison

All hash comparisons use constant-time algorithms to prevent timing attacks.

---

## Request Approval Signing

When approving transaction requests, users sign with their private key:

### Process

1. Retrieve the request and its metadata
2. Extract the metadata hash (SHA-256 of the payload)
3. Sign the hash with user's private key (ECDSA P-256)
4. Submit the signature to the API

### What Gets Signed

Users sign the `metadata.hash` field, which is a SHA-256 hash of the transaction details:
- Source address
- Destination address
- Amount
- Currency
- Additional transaction parameters

This ensures users are approving the exact transaction details.

---

## Security Best Practices

### API Credentials

- Use separate API keys for different environments (dev, staging, prod)
- Implement key rotation procedures
- Monitor API key usage for anomalies
- Revoke compromised keys immediately

### Private Keys

- Use hardware security modules (HSM) for production private keys
- Implement key ceremony procedures for critical keys
- Limit access to signing operations
- Audit all signing activities

### SuperAdmin Keys

- Verify SuperAdmin keys through secure out-of-band channels
- Update keys when administrators change
- Use appropriate signature thresholds (minimum 2 recommended)
- Maintain secure backup of key material

### Network Security

- Always use HTTPS (TLS 1.2+)
- Implement certificate pinning for high-security applications
- Use appropriate network segmentation
- Monitor for man-in-the-middle attacks

---

## Error Handling

### Authentication Errors

| Error | Cause | Resolution |
|-------|-------|------------|
| 401 Unauthorized | Invalid API key or signature | Verify credentials |
| 403 Forbidden | Valid auth but insufficient permissions | Check API key permissions |

### Verification Errors

| Error | Cause | Resolution |
|-------|-------|------------|
| IntegrityError | Hash mismatch | Data may be tampered - investigate |
| IntegrityError | Insufficient signatures | Check SuperAdmin key configuration |
| WhitelistError | Signature verification failed | User signature invalid |

---

## Related Documentation

### Common Documentation
- [Key Concepts](CONCEPTS.md) - Domain model and entities
- [Integrity Verification](INTEGRITY_VERIFICATION.md) - Detailed verification flows

### SDK-Specific Documentation

**Java SDK**
- [Authentication](../taurus-protect-sdk-java/docs/AUTHENTICATION.md) - Java-specific auth implementation
- [Usage Examples](../taurus-protect-sdk-java/docs/USAGE_EXAMPLES.md) - Java code examples

**Go SDK**
- [Authentication](../taurus-protect-sdk-go/docs/AUTHENTICATION.md) - Go-specific auth implementation
- [Usage Examples](../taurus-protect-sdk-go/docs/USAGE_EXAMPLES.md) - Go code examples

**Python SDK**
- [Authentication](../taurus-protect-sdk-python/docs/AUTHENTICATION.md) - Python-specific auth implementation
- [Usage Examples](../taurus-protect-sdk-python/docs/USAGE_EXAMPLES.md) - Python code examples

**TypeScript SDK**
- [Authentication](../taurus-protect-sdk-typescript/docs/AUTHENTICATION.md) - TypeScript-specific auth implementation
- [Usage Examples](../taurus-protect-sdk-typescript/docs/USAGE_EXAMPLES.md) - TypeScript code examples