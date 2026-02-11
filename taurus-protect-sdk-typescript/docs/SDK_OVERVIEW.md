# Taurus-PROTECT TypeScript SDK Overview

## Introduction

The Taurus-PROTECT TypeScript SDK provides a high-level TypeScript/JavaScript client for interacting with the Taurus-PROTECT API. It handles authentication, request signing, cryptographic verification of responses, and provides clean domain objects for all API operations.

## Architecture

### Package Structure

The SDK is organized as a single npm package with modular internal structure:

```
taurus-protect-sdk-typescript/
├── src/
│   ├── client.ts              # ProtectClient entry point
│   ├── errors.ts              # Exception hierarchy
│   ├── transport/             # HTTP middleware (TPV1)
│   ├── crypto/                # Cryptographic utilities
│   ├── services/              # High-level service layer
│   ├── models/                # Domain models
│   ├── mappers/               # DTO to model conversion
│   ├── helpers/               # Verification utilities
│   ├── cache/                 # RulesContainerCache
│   └── internal/
│       ├── openapi/           # Auto-generated OpenAPI client (DO NOT MODIFY)
│       └── proto/             # Auto-generated protobuf classes (DO NOT MODIFY)
├── tests/
│   ├── unit/                  # Unit tests
│   └── integration/           # Integration tests
└── dist/                      # Compiled output
```

| Directory | Purpose | Auto-Generated | Modifiable |
|-----------|---------|----------------|------------|
| **src/** | Public SDK source code | No | Yes |
| **src/internal/openapi/** | Low-level HTTP API client from OpenAPI specs | Yes | No |
| **src/internal/proto/** | Protocol Buffer classes (currently excluded) | Yes | No |

### Module Dependencies

```
┌─────────────────────────────────────────────────────────────┐
│                         client.ts                            │
│           (Entry point - Service Factory)                    │
├─────────────────────────────────────────────────────────────┤
│                             │                                │
│              ┌──────────────┼──────────────┐                │
│              ▼              ▼              ▼                │
│        ┌──────────┐  ┌──────────┐  ┌──────────────┐        │
│        │ services │  │  models  │  │   crypto     │        │
│        │  layer   │  │  layer   │  │   layer      │        │
│        └────┬─────┘  └──────────┘  └──────────────┘        │
│             │                                                │
│             ▼                                                │
│        ┌──────────────────────────┐                         │
│        │   internal/openapi       │                         │
│        │   (Generated HTTP API)   │                         │
│        └──────────────────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

## Module Details

### Client Entry Point

**Location:** `src/client.ts`

The `ProtectClient` class is the main entry point that developers interact with. It provides:

- **Factory Pattern** - Create clients via `ProtectClient.create(config)`
- **Lazy Initialization** - APIs and services are instantiated on first access
- **Resource Management** - `close()` method for cleanup
- **Configuration Validation** - Validates host URL, API credentials at creation time

**Key Properties:**
- 56 low-level API accessors (e.g., `walletsApi`, `addressesApi`)
- 26 high-level service getters (e.g., `wallets`, `addresses`)
- TaurusNetwork namespace for low-level Taurus Network API access
- Additional service classes in `src/services/` available for direct instantiation

### Services Layer

**Location:** `src/services/`

The services layer wraps OpenAPI calls with:
- Input validation
- Exception translation via `BaseService.execute()`
- Cryptographic verification (where applicable)
- Response mapping to domain objects

**Package Structure:**
```
src/services/
├── base.ts                    # BaseService abstract class
├── wallet-service.ts          # WalletService
├── address-service.ts         # AddressService
├── request-service.ts         # RequestService (with hash verification)
├── governance-rule-service.ts # GovernanceRuleService (with signature verification)
├── taurus-network/            # TaurusNetwork services
│   ├── participant-service.ts
│   ├── pledge-service.ts
│   ├── lending-service.ts
│   ├── settlement-service.ts
│   ├── sharing-service.ts
│   └── index.ts
└── index.ts                   # Service exports
```

**Service Count:** 26 high-level services accessible via ProtectClient getters. Additional service classes exist in `src/services/` (including TaurusNetwork services) for direct instantiation.

### Models Layer

**Location:** `src/models/`

Clean domain objects exposed to SDK users. All models use:
- TypeScript interfaces with `readonly` properties
- Enums for status and type fields
- Optional fields marked with `| undefined`

**Package Structure:**
```
src/models/
├── pagination.ts              # Pagination, PaginatedResult
├── wallet.ts                  # Wallet, WalletStatus, WalletBalance
├── address.ts                 # Address, AddressStatus
├── request.ts                 # Request, RequestStatus, RequestMetadata
├── transaction.ts             # Transaction, TransactionStatus
├── governance-rules.ts        # GovernanceRules, DecodedRulesContainer
├── taurus-network/            # TaurusNetwork models (72 models)
│   ├── participant.ts
│   ├── pledge.ts
│   ├── lending.ts
│   ├── settlement.ts
│   ├── sharing.ts
│   └── index.ts
└── index.ts                   # Model exports
```

### Crypto Layer

**Location:** `src/crypto/`

Cryptographic utilities for authentication and verification:

| File | Purpose |
|------|---------|
| `tpv1.ts` | TPV1-HMAC-SHA256 authentication |
| `signing.ts` | ECDSA P-256 signing/verification |
| `hashing.ts` | SHA-256 and HMAC functions |
| `keys.ts` | PEM key decoding (`decodePrivateKeyPem`, `decodePublicKeyPem`, `decodePublicKeysPem`, `encodePublicKeyPem`, `getPublicKeyFromPrivate`) |

### Helpers Layer

**Location:** `src/helpers/`

Verification utilities:

| File | Purpose |
|------|---------|
| `constant-time.ts` | Timing-safe comparison |
| `signature-verifier.ts` | SuperAdmin signature verification |
| `address-signature-verifier.ts` | HSM signature verification |

### Cache Layer

**Location:** `src/cache/`

| File | Purpose |
|------|---------|
| `rules-container-cache.ts` | Thread-safe governance rules caching with TTL |

### OpenAPI Module

**Location:** `src/internal/openapi/`

Auto-generated from OpenAPI specifications using `openapi-generator-cli`. Contains:

- HTTP client implementation (fetch-based)
- API endpoint classes (56 APIs)
- DTO classes (Data Transfer Objects)
- Runtime configuration

**Regeneration:**
```bash
./build.sh generate
```

## Key Design Patterns

### 1. Factory Pattern (ProtectClient)

`ProtectClient.create()` serves as the factory for creating clients:

```typescript
const client = ProtectClient.create({
  host: 'https://your-protect-instance.example.com',
  apiKey: 'your-api-key',
  apiSecret: 'your-hex-encoded-secret',
  superAdminKeysPem: ['-----BEGIN PUBLIC KEY-----...'],
  minValidSignatures: 2,
});

// Access services through the client
const wallets = await client.wallets.list();
const address = await client.addresses.get('address-id');
```

### 2. Lazy Initialization Pattern

All APIs and services are lazily initialized on first access:

```typescript
// Internal implementation
get wallets(): WalletService {
  this.ensureOpen();
  if (!this._walletService) {
    this._walletService = new WalletService(this.walletsApi);
  }
  return this._walletService;
}
```

### 3. BaseService Pattern

All services extend `BaseService` which provides:
- `execute()` method for wrapping API calls with error handling
- `handleError()` for converting OpenAPI errors to SDK errors

```typescript
// Service layer flow
Request → Service.execute() → OpenAPI Client → API
                                    ↓
Domain Object ← Mapper ← DTO ← Response
```

### 4. Exception Handling Pattern

Services convert OpenAPI errors to typed SDK exceptions:

```typescript
try {
  return await this.execute(async () => {
    const dto = await this.api.getWallet({ walletId });
    return mapWallet(dto);
  });
} catch (error) {
  if (error instanceof NotFoundError) {
    // Handle 404
  } else if (error instanceof APIError && error.isRetryable()) {
    // Handle retryable errors
  }
}
```

### 5. TaurusNetwork Namespace Pattern

TaurusNetwork APIs are grouped under a namespace:

```typescript
// Access via namespace (low-level API access)
const participants = await client.taurusNetwork.participantApi.getAllParticipants();
const pledges = await client.taurusNetwork.pledgeApi.getAllPledges();
const offers = await client.taurusNetwork.lendingApi.getAllLendingOffers();
```

## Build and Development

### Build Commands

Use `./build.sh` for standard operations:

```bash
# Build and test (USE THIS TO VERIFY CHANGES)
./build.sh

# Build only
./build.sh build

# Run tests only
./build.sh test

# Run linter
./build.sh lint

# Code generation
./build.sh generate

# Clean build artifacts
./build.sh clean
```

### Advanced: Direct npm Commands

For most workflows, use `./build.sh` above.

```bash
# Install dependencies
npm install

# Build TypeScript
npm run build

# Run tests
npm test

# Run linter
npm run lint
```

### Code Generation

OpenAPI client generation requires Java 11+:

```bash
./build.sh generate
```

This regenerates:
- `src/internal/openapi/` - TypeScript fetch client from OpenAPI specs

## Key Classes Reference

| Class | Purpose |
|-------|---------|
| `ProtectClient` | Entry point, service factory |
| `BaseService` | Base class for all services |
| `APIError` | Base exception with HTTP status |
| `IntegrityError` | Hash/signature verification failure |
| `RulesContainerCache` | Governance rules caching |

## Available Services

### Client Service Getters (26 high-level services)

These services are accessible directly as getters on `ProtectClient`:

**Core Services (8):**
- `wallets` - Wallet management
- `addresses` - Address management with HSM signature verification
- `requests` - Request management with hash verification and ECDSA approval
- `transactions` - Transaction retrieval
- `balances` - Balance queries
- `currencies` - Currency information
- `health` - Health check
- `jobs` - Background job monitoring

**Administrative Services (6):**
- `users` - User management
- `groups` - Group management
- `visibilityGroups` - Visibility group management
- `tags` - Tag management
- `webhooks` - Webhook configuration
- `audits` - Audit trail

**Security Services (3):**
- `governanceRules` - Governance rules with signature verification
- `whitelistedAddresses` - Whitelisted address management
- `whitelistedAssets` - Whitelisted asset management

**Advanced Services (4):**
- `airGap` - Air-gapped HSM operations
- `configService` - Tenant configuration
- `assets` - Asset queries
- `statistics` - Portfolio statistics

**Blockchain and Pricing Services (4):**
- `exchanges` - Exchange account management
- `prices` - Price information
- `fees` - Network fee information
- `feePayers` - Fee payer management

**Specialized Services (1):**
- `tokenMetadata` - Token metadata

### Service Classes Without Client Getters

These service classes exist in `src/services/` and can be instantiated directly, but do not have getters on `ProtectClient`:

- `WebhookCallService` - Webhook call history
- `StakingService` - Staking operations
- `ContractWhitelistingService` - Contract whitelisting
- `ReservationService` - Balance reservations
- `MultiFactorSignatureService` - MFA signatures
- `BusinessRuleService` - Business rule management
- `ChangeService` - Change tracking
- `BlockchainService` - Blockchain information
- `FiatService` - Fiat currency operations
- `ScoreService` - Risk scores
- `UserDeviceService` - User device management
- `ActionService` - Action management

These features are also accessible via the low-level OpenAPI-generated APIs (e.g., `client.businessRulesApi`, `client.stakingApi`).

### TaurusNetwork APIs (5 low-level APIs)

TaurusNetwork provides low-level API access (not high-level service wrappers):
- `taurusNetwork.participantApi` - Participant management
- `taurusNetwork.pledgeApi` - Pledge lifecycle
- `taurusNetwork.lendingApi` - Loan offers and agreements
- `taurusNetwork.settlementApi` - Settlement operations
- `taurusNetwork.sharedAddressAssetApi` - Address/asset sharing

## Related Documentation

- [Concepts](CONCEPTS.md) - TypeScript-specific model classes and exceptions
- [Services Reference](SERVICES.md) - Complete service API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns

### Common Documentation (applies to all SDKs)
- [Key Concepts](../../docs/CONCEPTS.md) - Full domain model
- [Authentication](../../docs/AUTHENTICATION.md) - TPV1 protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Verification flows
