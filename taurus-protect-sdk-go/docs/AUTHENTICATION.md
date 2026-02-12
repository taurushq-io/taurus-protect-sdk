# Authentication & Security

## Overview

The Taurus-PROTECT SDK uses a multi-layer security model:

1. **TPV1 Authentication** - HMAC-based API request signing
2. **SuperAdmin Verification** - ECDSA signature verification for governance rules
3. **Data Integrity** - SHA-256 hash verification for payloads

## TPV1 Authentication Scheme

All API requests are signed using the TPV1 (Taurus Protocol Version 1) scheme.

### Header Format

```
Authorization: TPV1-HMAC-SHA256 ApiKey=<api_key> Nonce=<uuid> Timestamp=<unix_ms> Signature=<base64_sig>
```

### Signature Computation

The signature is computed over:
```
TPV1 <ApiKey> <Nonce> <Timestamp> <Method> <Host> <Path> <Query> <ContentType> <Body>
```

Using HMAC-SHA256 with the API secret as the key.

### Automatic Handling

The SDK handles TPV1 signing automatically via the custom HTTP transport. You only need to provide credentials during client initialization.

## Client Initialization

### Basic Initialization

```go
import (
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

client, err := protect.NewClient(
    "https://api.taurus-protect.com",
    protect.WithCredentials(apiKey, apiSecret),
)
if err != nil {
    return err
}
defer client.Close()
```

### With SuperAdmin Keys (Recommended)

```go
import (
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

// SuperAdmin public keys in PEM format
superAdminKeys := []string{
    `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`,
    `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`,
}

client, err := protect.NewClient(
    "https://api.taurus-protect.com",
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeysPEM(superAdminKeys),
    protect.WithMinValidSignatures(2),
)
if err != nil {
    return err
}
defer client.Close()
```

### With Pre-Parsed Keys

```go
import (
    "crypto/ecdsa"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Parse keys manually
key1, err := crypto.DecodePublicKeyPEM(pemKey1)
if err != nil {
    return err
}
key2, err := crypto.DecodePublicKeyPEM(pemKey2)
if err != nil {
    return err
}

client, err := protect.NewClient(
    host,
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeys([]*ecdsa.PublicKey{key1, key2}),
    protect.WithMinValidSignatures(2),
)
```

### With Custom Cache TTL

```go
client, err := protect.NewClient(
    host,
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeysPEM(superAdminKeys),
    protect.WithMinValidSignatures(2),
    protect.WithRulesCacheTTL(10 * time.Minute),
)
```

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `WithCredentials` | (string, string) | API key and hex-encoded secret (required) |
| `WithSuperAdminKeysPEM` | []string | SuperAdmin public keys in PEM format |
| `WithSuperAdminKeys` | []*ecdsa.PublicKey | Pre-parsed SuperAdmin public keys |
| `WithMinValidSignatures` | int | Minimum SuperAdmin signatures required |
| `WithRulesCacheTTL` | time.Duration | Rules cache TTL (default: 5 minutes) |
| `WithHTTPClient` | *http.Client | Custom HTTP client |
| `WithHTTPTimeout` | time.Duration | HTTP request timeout (default: 30s) |

## Cryptographic Operations

### crypto Package

**Location:** `pkg/protect/crypto/tpv1.go`

#### Key Decoding

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"

// Decode private key from PEM
pemPrivateKey := `-----BEGIN EC PRIVATE KEY-----
...
-----END EC PRIVATE KEY-----`
privateKey, err := crypto.DecodePrivateKeyPEM(pemPrivateKey)

// Decode public key from PEM
pemPublicKey := `-----BEGIN PUBLIC KEY-----
...
-----END PUBLIC KEY-----`
publicKey, err := crypto.DecodePublicKeyPEM(pemPublicKey)
```

#### ECDSA Signing (P-256 / secp256r1)

```go
// Sign data
data := []byte("message to sign")
base64Signature, err := crypto.SignData(privateKey, data)

// Verify signature
isValid, err := crypto.VerifySignature(publicKey, data, base64Signature)
```

#### Hash Computation (SHA-256)

```go
// Compute hex-encoded SHA-256 hash
hexHash := crypto.CalculateHexHash(data)
```

#### HMAC Operations

```go
// Compute HMAC-SHA256 (base64 encoded)
base64Hmac := crypto.CalculateBase64HMAC(secret, data)

// Verify HMAC (constant-time comparison)
valid := crypto.CheckBase64HMAC(secret, data, base64Hmac)
```

#### Secure Memory Wiping

```go
// Securely zero sensitive data
sensitiveData := []byte("secret")
defer crypto.Wipe(sensitiveData)
```

## SuperAdmin Public Keys

### Purpose

SuperAdmin public keys are used to verify:

1. **Governance Rules** - Rules defining approval thresholds and user groups
2. **Whitelisted Addresses** - Cryptographically signed address whitelists

### Configuration

SuperAdmin keys are typically provided by your Taurus-PROTECT administrator. You need:

- Multiple public keys (for multi-signature verification)
- Minimum signatures threshold (how many must be valid)

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

## Request Approval Signing

When approving transaction requests, you sign on behalf of the user:

```go
import (
    "context"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Load user's private key
privateKey, err := crypto.DecodePrivateKeyPEM(userPrivateKeyPem)
if err != nil {
    return err
}

// Get request to approve
request, err := client.Requests().GetRequest(ctx, requestID)
if err != nil {
    return err
}

// Approve the request - SDK handles signing internally
signedCount, err := client.Requests().ApproveRequest(ctx, request, privateKey)
if err != nil {
    return err
}
fmt.Printf("Approved %d request(s)\n", signedCount)
```

### Signing Process

1. SDK extracts request metadata hashes and sorts requests by ID
2. SDK builds JSON array of hashes and signs with user's private key (ECDSA P-256)
3. SDK submits signed approval to API and returns count of signed requests

## Security Best Practices

### API Credentials

- Store API secrets securely (environment variables, secret managers)
- Never commit credentials to version control
- Rotate API keys periodically
- Use the minimum required permissions

### Private Keys

- Use hardware security modules (HSM) for production private keys
- Limit access to signing operations
- Audit all signing activities
- Use `crypto.Wipe()` to clear sensitive data from memory

### SuperAdmin Keys

- Verify SuperAdmin keys through secure channels
- Update keys when administrators change
- Use appropriate signature thresholds (minimum 2 recommended)

### HTTP Client

```go
// Use custom HTTP client with timeouts
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    },
}

client, err := protect.NewClient(
    host,
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithHTTPClient(httpClient),
)
```

## Resource Cleanup

Always close the client when done to securely wipe credentials from memory:

```go
// Recommended: use defer for cleanup
client, err := protect.NewClient(host, protect.WithCredentials(apiKey, apiSecret))
if err != nil {
    return err
}
defer client.Close()

// The Close() method:
// 1. Wipes the API secret from memory using crypto.Wipe()
// 2. Invalidates the HTTP transport
// 3. Clears cached governance rules
```

For sensitive intermediate data, use `crypto.Wipe()` directly:

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"

// Securely zero sensitive byte slices
secretData := []byte("sensitive")
defer crypto.Wipe(secretData)

// Wipe private key material after use
keyBytes := loadKeyFromVault()
defer crypto.Wipe(keyBytes)
privateKey, err := crypto.DecodePrivateKeyPEM(string(keyBytes))
```

## Error Handling

### IntegrityError

Returned when cryptographic verification fails:

```go
import (
    "errors"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

addr, err := client.WhitelistedAddresses().GetWhitelistedAddress(ctx, id)
if err != nil {
    if protect.IsIntegrityError(err) {
        // Hash mismatch or insufficient valid signatures
        log.Error("Verification failed:", err)
        return err
    }
    return err
}
```

### Common Causes

| Error | Cause |
|-------|-------|
| `IntegrityError: hash mismatch` | Payload tampering detected |
| `IntegrityError: insufficient signatures` | SuperAdmin key mismatch or not enough valid signatures |
| `WhitelistError: verification failed` | User signature verification failed |

### Error Type Checking

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

if apiErr, ok := protect.IsAPIError(err); ok {
    fmt.Printf("API Error: %d - %s\n", apiErr.Code, apiErr.Message)
    if apiErr.IsRetryable() {
        // Implement retry logic
    }
}

if protect.IsIntegrityError(err) {
    // Cryptographic verification failed - investigate!
}

if protect.IsWhitelistError(err) {
    // Whitelist-specific error
}
```

## Environment-Based Configuration

```go
import (
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func createClientFromEnv() (*protect.Client, error) {
    host := os.Getenv("TAURUS_API_HOST")
    apiKey := os.Getenv("TAURUS_API_KEY")
    apiSecret := os.Getenv("TAURUS_API_SECRET")

    // SuperAdmin keys as newline-separated PEM strings
    keysEnv := os.Getenv("TAURUS_SUPERADMIN_KEYS")
    var superAdminKeys []string
    if keysEnv != "" {
        superAdminKeys = strings.Split(keysEnv, "\\n\\n")
    }

    minSigs := 2
    if minSigsEnv := os.Getenv("TAURUS_MIN_SIGNATURES"); minSigsEnv != "" {
        minSigs, _ = strconv.Atoi(minSigsEnv)
    }

    opts := []protect.ClientOption{
        protect.WithCredentials(apiKey, apiSecret),
    }

    if len(superAdminKeys) > 0 {
        opts = append(opts,
            protect.WithSuperAdminKeysPEM(superAdminKeys),
            protect.WithMinValidSignatures(minSigs),
        )
    }

    return protect.NewClient(host, opts...)
}
```

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Detailed verification flow