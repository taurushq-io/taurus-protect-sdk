# Authentication & Security

## Overview

The Taurus-PROTECT TypeScript SDK uses a multi-layer security model:

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

The signature is computed over a message with the following format:

```
TPV1 <ApiKey> <Nonce> <Timestamp> <Method> <Host> <Path> [Query] [ContentType] [Body]
```

Using HMAC-SHA256 with the API secret as the key.

### Automatic Handling

The SDK handles TPV1 signing automatically through middleware. You only need to provide credentials during client initialization.

## Client Initialization

### Basic Setup

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key-uuid',
  apiSecret: 'your-api-secret-hex',
});

try {
  const wallets = await client.wallets.list();
  console.log(`Found ${wallets.items.length} wallets`);
} finally {
  client.close();
}
```

### With SuperAdmin Key Verification

For enhanced security with governance rules and whitelisted address verification:

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

// SuperAdmin public keys in PEM format
const superAdmin1 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`;

const superAdmin2 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`;

const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key-uuid',
  apiSecret: 'your-api-secret-hex',
  superAdminKeysPem: [superAdmin1, superAdmin2],
  minValidSignatures: 2,
});
```

### With Custom Cache TTL

```typescript
const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key-uuid',
  apiSecret: 'your-api-secret-hex',
  superAdminKeysPem: [superAdmin1, superAdmin2],
  minValidSignatures: 2,
  rulesCacheTtlMs: 600000, // 10 minutes
});
```

### Configuration Options

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `host` | string | Yes | API base URL |
| `apiKey` | string | Yes | UUID-format API key |
| `apiSecret` | string | Yes | Hex-encoded HMAC secret |
| `superAdminKeysPem` | string[] | No | SuperAdmin public keys for verification |
| `minValidSignatures` | number | No | Minimum SuperAdmin signatures required (default: 1) |
| `rulesCacheTtlMs` | number | No | Rules cache TTL in milliseconds (default: 5 minutes) |
| `timeout` | number | No | Request timeout in milliseconds (default: 30 seconds) |
| `middleware` | Middleware[] | No | Additional middleware to apply to requests |

## Environment-Based Configuration

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';
import * as fs from 'fs';

function createClientFromEnv(): ProtectClient {
  const host = process.env.TAURUS_API_HOST!;
  const apiKey = process.env.TAURUS_API_KEY!;
  const apiSecret = process.env.TAURUS_API_SECRET!;

  // Optional: Load SuperAdmin keys from files
  const keysPath = process.env.TAURUS_SUPERADMIN_KEYS_PATH;
  let superAdminKeysPem: string[] | undefined;
  if (keysPath) {
    superAdminKeysPem = keysPath
      .split(',')
      .map((p) => fs.readFileSync(p.trim(), 'utf-8'));
  }

  const minValidSignatures = parseInt(
    process.env.TAURUS_MIN_SIGNATURES || '2',
    10,
  );

  return ProtectClient.create({
    host,
    apiKey,
    apiSecret,
    superAdminKeysPem,
    minValidSignatures,
  });
}
```

## Cryptographic Operations

### Crypto Module Exports

The SDK exposes cryptographic utilities through the `crypto` module:

```typescript
import {
  // Hashing
  calculateHexHash,
  calculateSha256Bytes,
  calculateBase64Hmac,
  verifyBase64Hmac,
  constantTimeCompare,
  constantTimeCompareBytes,

  // Key handling
  decodePrivateKeyPem,
  decodePublicKeyPem,
  decodePublicKeysPem,
  encodePublicKeyPem,
  getPublicKeyFromPrivate,

  // ECDSA signing
  signData,
  verifySignature,

  // TPV1 authentication
  TPV1Auth,
  calculateSignedHeader,
} from '@taurushq/protect-sdk';
```

### Key Decoding

```typescript
import { decodePrivateKeyPem, decodePublicKeyPem, decodePublicKeysPem } from '@taurushq/protect-sdk';

// Decode private key from PEM
const pemPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MEECAQAwEwYHKoZIzj0CAQYIKoZIzj0DAQcEJzAlAgEBBCB...
-----END EC PRIVATE KEY-----`;
const privateKey = decodePrivateKeyPem(pemPrivateKey);

// Decode public key from PEM
const pemPublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`;
const publicKey = decodePublicKeyPem(pemPublicKey);

// Decode multiple public keys
const publicKeys = decodePublicKeysPem([pem1, pem2, pem3]);
```

### ECDSA Signing (SHA256 with Raw r||s Format)

The SDK uses ECDSA P-256 with raw r||s signature format, compatible with Java's SHA256withPLAIN-ECDSA:

```typescript
import { signData, verifySignature } from '@taurushq/protect-sdk';

// Sign data
const data = Buffer.from('message to sign', 'utf-8');
const base64Signature = signData(privateKey, data);

// Verify signature
const isValid = verifySignature(publicKey, data, base64Signature);
```

### Hash Computation (SHA-256)

```typescript
import { calculateHexHash } from '@taurushq/protect-sdk';

// Compute hex-encoded SHA-256 hash
const hexHash = calculateHexHash('data to hash');
// Returns: hex-encoded hash string
```

### HMAC Operations

```typescript
import { calculateBase64Hmac, verifyBase64Hmac, constantTimeCompare } from '@taurushq/protect-sdk';

// Compute HMAC-SHA256 (base64 encoded)
const secret = Buffer.from(apiSecretHex, 'hex');
const base64Hmac = calculateBase64Hmac(secret, 'data to sign');

// Verify HMAC (constant-time comparison)
const isValid = verifyBase64Hmac(secret, 'data to sign', base64Hmac);

// Constant-time string comparison (prevents timing attacks)
const areEqual = constantTimeCompare(hash1, hash2);
```

## SuperAdmin Public Keys

### Purpose

SuperAdmin public keys are used to verify:

1. **Governance Rules** - Rules defining approval thresholds and user groups
2. **Whitelisted Addresses** - Cryptographically signed address whitelists
3. **Whitelisted Assets** - Cryptographically signed asset/contract whitelists

### Configuration

SuperAdmin keys are typically provided by your Taurus-PROTECT administrator. You need:

- Multiple public keys (for multi-signature verification)
- Minimum signatures threshold (how many must be valid)

### Verification Flow

```
+-----------------------------------------------------------+
|                 Governance Rules                           |
|  +-------------------+                                    |
|  | rulesContainer    |  (Base64-encoded JSON)             |
|  | rulesSignatures   |  (List of ECDSA signatures)        |
|  +-------------------+                                    |
+-----------------------------------------------------------+
                          |
                          v
+-----------------------------------------------------------+
|            Signature Verification                          |
|  For each signature:                                       |
|    1. Decode base64 signature                             |
|    2. Verify against SuperAdmin public keys               |
|    3. Count valid signatures (track distinct user IDs)    |
|  Require: validCount >= minValidSignatures                |
+-----------------------------------------------------------+
```

### Verification with GovernanceRuleService

```typescript
import { ProtectClient, decodePublicKeysPem } from '@taurushq/protect-sdk';

// Create client
const client = ProtectClient.create({ host, apiKey, apiSecret });

// Configure verification with KeyObject instances
const superAdminKeys = decodePublicKeysPem([superAdmin1Pem, superAdmin2Pem]);

// Create service with verification
const govService = client.governanceRules;

// Get rules with automatic verification
const rules = await govService.get();
```

## Request Approval Signing

When approving transaction requests, the SDK signs the request hashes using ECDSA:

```typescript
import { ProtectClient, decodePrivateKeyPem } from '@taurushq/protect-sdk';

// Load user's private key
const privateKey = decodePrivateKeyPem(userPrivateKeyPem);

// Approve request (signs internally)
const count = await client.requests.approveRequest(request, privateKey);
console.log(`Approved: ${count} request(s)`);
```

### Signing Process

The approval signing process follows these steps:

1. **Sort requests by ID** (numeric ascending)
2. **Build JSON array of hashes**: `JSON.stringify(requests.map(r => r.metadata.hash))`
3. **Sign with ECDSA** using raw r||s format (64 bytes, base64 encoded)
4. **Submit signed approval** to API

Example implementation:

```typescript
// Internal flow (handled by RequestService)
const sortedRequests = [...requests].sort((a, b) => a.id - b.id);
const hashes = sortedRequests.map(r => r.metadata!.hash);
const hashesJson = JSON.stringify(hashes);
const signature = signData(privateKey, Buffer.from(hashesJson, 'utf-8'));

// Submit to API
await requestsApi.requestServiceApproveRequests({
  body: {
    ids: sortedRequests.map(r => String(r.id)),
    signature,
    comment: 'Approved via SDK',
  },
});
```

## Security Best Practices

### API Credentials

- Store API secrets securely (environment variables, secret managers)
- Never commit credentials to version control
- Rotate API keys periodically
- Use separate API keys for different environments

### Private Keys

- Use hardware security modules (HSM) for production private keys
- Limit access to signing operations
- Audit all signing activities
- Consider using AWS KMS, Azure Key Vault, or HashiCorp Vault

### SuperAdmin Keys

- Verify SuperAdmin keys through secure out-of-band channels
- Update keys when administrators change
- Use appropriate signature thresholds (minimum 2 recommended)
- Store public keys in version control or configuration management

### Client Lifecycle

Always close the client when done to clean up resources:

```typescript
const client = ProtectClient.create(config);
try {
  // Use the client
} finally {
  client.close();
}
```

## Error Handling

### IntegrityError

Thrown when cryptographic verification fails:

```typescript
import { ProtectClient, IntegrityError } from '@taurushq/protect-sdk';

try {
  const request = await client.requests.get(requestId);
} catch (error) {
  if (error instanceof IntegrityError) {
    // Hash mismatch or signature verification failed
    console.error('Security error:', error.message);
    // This indicates potential tampering - investigate!
  }
}
```

### Common Causes

| Exception | Cause |
|-----------|-------|
| `IntegrityError: Request hash mismatch` | Payload tampering detected |
| `IntegrityError: Insufficient valid signatures` | SuperAdmin key mismatch or not enough signatures |
| `IntegrityError: Invalid signature for address` | HSM signature verification failed |
| `WhitelistError: verification failed` | Whitelist signature verification failed |

## TPV1Auth Class

For advanced use cases, you can use the `TPV1Auth` class directly:

```typescript
import { TPV1Auth } from '@taurushq/protect-sdk';

const auth = new TPV1Auth('api-key', 'hex-encoded-secret');

try {
  // Generate authorization header for a request
  const header = auth.signRequest(
    'POST',
    'api.protect.taurushq.com',
    '/v1/wallets',
    undefined,           // query string
    'application/json',  // content-type
    '{"name":"my-wallet"}' // body
  );
  // header = "TPV1-HMAC-SHA256 ApiKey=... Nonce=... Timestamp=... Signature=..."
} finally {
  // Securely wipe the secret from memory
  auth.close();
}
```

### Parsing URLs

```typescript
const { host, path, query } = TPV1Auth.parseUrl(
  'https://api.example.com/v1/wallets?limit=10'
);
// host = "api.example.com"
// path = "/v1/wallets"
// query = "limit=10"
```

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples for common operations
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Detailed verification flow
