# Authentication & Security

## Overview

The Taurus PROTECT SDK uses a multi-layer security model:

1. **TPV1 Authentication** - HMAC-based API request signing
2. **SuperAdmin Verification** - ECDSA signature verification for governance rules
3. **Data Integrity** - SHA-256 hash verification for payloads

## TPV1 Authentication Scheme

All API requests are signed using the TPV1 (Taurus Protocol Version 1) scheme.

### Header Format

```
Authorization: TPV1-HMAC-SHA256 ApiKey=<api_key> Nonce=<uuid> Timestamp=<unix_ts> Signature=<base64_sig>
```

### Signature Computation

The signature is computed over:
```
TPV1 <ApiKey> <Nonce> <Timestamp> <Method> <Host> <Path> <Query> <ContentType> <Body>
```

Using HMAC-SHA256 with the API secret as the key.

### Automatic Handling

The SDK handles TPV1 signing automatically. You only need to provide credentials during client initialization.

## Client Initialization

### Using PEM-Encoded Keys (Recommended)

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import java.util.Arrays;

// SuperAdmin public keys in PEM format
String superAdmin1 = "-----BEGIN PUBLIC KEY-----\nMFkw....\n-----END PUBLIC KEY-----";
String superAdmin2 = "-----BEGIN PUBLIC KEY-----\nMFkw....\n-----END PUBLIC KEY-----";

ProtectClient client = ProtectClient.createFromPem(
    "https://api.taurus-protect.com",           // API host
    "your-api-key-uuid",                         // API key
    "your-api-secret-hex",                       // API secret (hex-encoded)
    Arrays.asList(superAdmin1, superAdmin2),     // SuperAdmin public keys
    2                                            // Minimum valid signatures
);
```

### Using PublicKey Objects

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import java.security.PublicKey;
import java.util.List;

// Decode keys first
List<PublicKey> superAdminKeys = CryptoTPV1.decodePublicKeys(pemKeysList);

ProtectClient client = ProtectClient.create(
    "https://api.taurus-protect.com",
    "your-api-key-uuid",
    "your-api-secret-hex",
    superAdminKeys,
    2                                            // Minimum valid signatures
);
```

### With Custom Cache TTL

```java
ProtectClient client = ProtectClient.create(
    host,
    apiKey,
    apiSecret,
    superAdminKeys,
    minValidSignatures,
    600000L                                      // Rules cache TTL in milliseconds (10 min)
);
```

### Constructor Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `host` | String | API base URL |
| `apiKey` | String | UUID-format API key |
| `apiSecret` | String | Hex-encoded HMAC secret |
| `superAdminPublicKeys` | List<PublicKey> or List<String> | SuperAdmin public keys for verification |
| `minValidSignatures` | int | Minimum SuperAdmin signatures required |
| `rulesContainerCacheTtlMs` | long | Optional cache TTL (default: 5 minutes) |

## Environment-Based Configuration

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import java.util.Arrays;

// Load credentials from environment variables
String host = System.getenv("TAURUS_API_HOST");
String apiKey = System.getenv("TAURUS_API_KEY");
String apiSecret = System.getenv("TAURUS_API_SECRET");

// SuperAdmin keys (comma-separated PEM file paths)
String keysPath = System.getenv("TAURUS_SUPERADMIN_KEYS_PATH");
List<String> superAdminKeys = new ArrayList<>();
if (keysPath != null) {
    for (String path : keysPath.split(",")) {
        superAdminKeys.add(new String(Files.readAllBytes(Paths.get(path.trim()))));
    }
}

int minSigs = Integer.parseInt(
    System.getenv().getOrDefault("TAURUS_MIN_SIGNATURES", "2")
);

ProtectClient client = ProtectClient.createFromPem(
    host, apiKey, apiSecret, superAdminKeys, minSigs
);
```

## Cryptographic Operations

### CryptoTPV1 Utility Class

**Location:** `openapi/src/main/java/com/taurushq/sdk/protect/openapi/auth/CryptoTPV1.java`

#### Key Decoding

```java
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;

// Decode private key from PEM
String pemPrivateKey = "-----BEGIN EC PRIVATE KEY-----\n...\n-----END EC PRIVATE KEY-----";
PrivateKey privateKey = CryptoTPV1.decodePrivateKey(pemPrivateKey);

// Decode public key from PEM
String pemPublicKey = "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----";
PublicKey publicKey = CryptoTPV1.decodePublicKey(pemPublicKey);

// Decode multiple public keys
List<String> pemKeys = Arrays.asList(pem1, pem2, pem3);
List<PublicKey> publicKeys = CryptoTPV1.decodePublicKeys(pemKeys);

// Encode public key to PEM
String pem = CryptoTPV1.encodePublicKey(publicKey);
```

#### ECDSA Signing (SHA256withPLAIN-ECDSA)

```java
// Sign data
byte[] data = "message to sign".getBytes(StandardCharsets.UTF_8);
String base64Signature = CryptoTPV1.calculateBase64Signature(privateKey, data);

// Verify signature
boolean isValid = CryptoTPV1.verifyBase64Signature(publicKey, data, base64Signature);
```

#### Hash Computation (SHA-256)

```java
// Compute hex-encoded SHA-256 hash
String hexHash = CryptoTPV1.calculateHexHash("data to hash");
```

#### HMAC Operations

```java
// Compute HMAC-SHA256 (base64 encoded)
byte[] secret = hexStringToBytes(apiSecret);
String base64Hmac = CryptoTPV1.calculateBase64Hmac(secret, data);

// Verify HMAC (constant-time comparison)
boolean valid = CryptoTPV1.checkBase64Hmac(secret, data, base64Hmac);

// Hex-encoded HMAC
String hexHmac = CryptoTPV1.calculateHexHmac(secretString, data);
boolean valid = CryptoTPV1.checkHexHmac(secretString, data, hexHmac);
```

## SuperAdmin Public Keys

### Purpose

SuperAdmin public keys are used to verify:

1. **Governance Rules** - Rules defining approval thresholds and user groups
2. **Whitelisted Addresses** - Cryptographically signed address whitelists

### Configuration

SuperAdmin keys are typically provided by your Taurus PROTECT administrator. You need:

- Multiple public keys (for multi-signature verification)
- Minimum signatures threshold (how many must be valid)

### Verification Flow

```
┌─────────────────────────────────────────────────────────┐
│                 Governance Rules                         │
│  ┌─────────────────┐                                    │
│  │ rulesContainer  │  (Base64-encoded protobuf)        │
│  │ rulesSignatures │  (List of ECDSA signatures)       │
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

When approving transaction requests, the SDK signs on behalf of the user:

```java
// Load user's private key
PrivateKey userPrivateKey = CryptoTPV1.decodePrivateKey(userPrivateKeyPem);

// Approve request (signs internally)
int signaturesPerformed = client.getRequestService().approveRequest(request, userPrivateKey);
```

### Signing Process

1. Extract request metadata hash
2. Sign hash with user's private key (SHA256withPLAIN-ECDSA)
3. Submit signed approval to API

## Security Best Practices

### API Credentials

- Store API secrets securely (environment variables, secret managers)
- Never commit credentials to version control
- Rotate API keys periodically

### Private Keys

- Use hardware security modules (HSM) for production private keys
- Limit access to signing operations
- Audit all signing activities

### SuperAdmin Keys

- Verify SuperAdmin keys through secure channels
- Update keys when administrators change
- Use appropriate signature thresholds (minimum 2 recommended)

## Exception Handling

### IntegrityException

Thrown when cryptographic verification fails:

```java
try {
    WhitelistedAddress addr = client.getWhitelistedAddressService().getWhitelistedAddress(id);
} catch (IntegrityException e) {
    // Hash mismatch or insufficient valid signatures
    System.err.println("Verification failed: " + e.getMessage());
}
```

### Common Causes

| Exception | Cause |
|-----------|-------|
| `IntegrityException: computed hash must equal provided hash` | Payload tampering detected |
| `IntegrityException: only N valid signatures found, minimum M required` | SuperAdmin key mismatch or insufficient signatures |
| `WhitelistException: signature verification failed` | User signature verification failed |

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Services Reference](SERVICES.md) - Complete API documentation
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Detailed verification flow
