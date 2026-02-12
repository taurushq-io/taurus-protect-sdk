# Taurus-PROTECT TypeScript SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A TypeScript SDK for interacting with the Taurus-PROTECT API, providing secure cryptocurrency custody and transaction management capabilities.

## Documentation

| Document | Description |
|----------|-------------|
| [Key Concepts](docs/CONCEPTS.md) | TypeScript models, exceptions, and domain concepts |
| [SDK Overview](docs/SDK_OVERVIEW.md) | Architecture, packages, and design patterns |
| [Authentication](docs/AUTHENTICATION.md) | TPV1 authentication and cryptographic operations |
| [Services Reference](docs/SERVICES.md) | Complete API documentation for all 43 services |
| [Usage Examples](docs/USAGE_EXAMPLES.md) | TypeScript code examples and common patterns |
| [Whitelisted Address Verification](docs/WHITELISTED_ADDRESS_VERIFICATION.md) | 6-step verification flow |

## Quick Start

### Prerequisites

- Node.js 18.0.0 or higher
- npm or yarn package manager

### Installation

Install the SDK using npm:

```bash
npm install @taurushq/protect-sdk
```

Or using yarn:

```bash
yarn add @taurushq/protect-sdk
```

### Client Initialization

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

// Basic initialization
const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key',
  apiSecret: 'your-api-secret-hex',
});

// Use the client
const wallets = await client.wallets.list();
for (const wallet of wallets.items) {
  console.log(`${wallet.name}: ${wallet.currency}`);
}

// Clean up when done
client.close();
```

For production use with SuperAdmin key verification:

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

const superAdminKeys = [
  '-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----',
  '-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----',
];

const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key',
  apiSecret: 'your-api-secret-hex',
  superAdminKeysPem: superAdminKeys,
  minValidSignatures: 2,
});

// All governance rule verifications will require 2 valid signatures
```

See [Authentication](docs/AUTHENTICATION.md) for more initialization options.

## Services

The SDK provides 43 services (38 core + 5 TaurusNetwork), with 26 available as high-level wrappers via ProtectClient getters and the remainder accessible through low-level OpenAPI APIs.

### High-Level Services (26 via ProtectClient getters)

These services provide domain models, validation, and simplified interfaces.

| Service | Access | Purpose |
|---------|--------|---------|
| `WalletService` | `client.wallets` | Wallet creation, retrieval, balance history |
| `AddressService` | `client.addresses` | Address management with signature verification |
| `RequestService` | `client.requests` | Transaction requests and approvals with ECDSA signing |
| `TransactionService` | `client.transactions` | Transaction queries and export |
| `BalanceService` | `client.balances` | Balance queries across assets |
| `CurrencyService` | `client.currencies` | Currency and blockchain information |
| `GovernanceRuleService` | `client.governanceRules` | Governance rules with signature verification |
| `WhitelistedAddressService` | `client.whitelistedAddresses` | Address whitelisting with 6-step verification |
| `WhitelistedAssetService` | `client.whitelistedAssets` | Asset/contract whitelisting |
| `AuditService` | `client.audits` | Audit log queries |
| `FeeService` | `client.fees` | Transaction fee information |
| `PriceService` | `client.prices` | Price data and history |
| `AirGapService` | `client.airGap` | Air-gap signing operations |
| `UserService` | `client.users` | User management |
| `GroupService` | `client.groups` | User group management |
| `VisibilityGroupService` | `client.visibilityGroups` | Visibility group management |
| `ConfigService` | `client.configService` | System configuration |
| `WebhookService` | `client.webhooks` | Webhook management |
| `TagService` | `client.tags` | Tag management |
| `AssetService` | `client.assets` | Asset information |
| `ExchangeService` | `client.exchanges` | Exchange integration |
| `FeePayerService` | `client.feePayers` | Fee payer management |
| `HealthService` | `client.health` | API health checks |
| `JobService` | `client.jobs` | Background job management |
| `StatisticsService` | `client.statistics` | Platform statistics |
| `TokenMetadataService` | `client.tokenMetadata` | Token metadata information |

### Low-Level API Access (for services without high-level wrappers)

Some features are available only through low-level OpenAPI-generated APIs:

| API | Access | Purpose |
|-----|--------|---------|
| `ChangesApi` | `client.changesApi` | Configuration change tracking |
| `BusinessRulesApi` | `client.businessRulesApi` | Business rule management |
| `ReservationsApi` | `client.reservationsApi` | Balance reservations |
| `MultiFactorSignatureApi` | `client.multiFactorSignatureApi` | Multi-factor signature operations |
| `ContractWhitelistingApi` | `client.contractWhitelistingApi` | Smart contract whitelisting |
| `StakingApi` | `client.stakingApi` | Multi-chain staking information |
| `ActionsApi` | `client.actionsApi` | Action management |
| `BlockchainApi` | `client.blockchainApi` | Blockchain information |
| `FiatApi` | `client.fiatApi` | Fiat currency operations |
| `ScoresApi` | `client.scoresApi` | Risk scoring operations |
| `UserDeviceApi` | `client.userDeviceApi` | User device management |
| `WebhookCallsApi` | `client.webhookCallsApi` | Webhook call history |

### TaurusNetwork APIs (5 APIs)

TaurusNetwork provides low-level API access for Taurus Network operations. Use these APIs directly for full control over the API calls.

| API | Access | Purpose |
|-----|--------|---------|
| `TaurusNetworkParticipantApi` | `client.taurusNetwork.participantApi` | Participant management |
| `TaurusNetworkPledgeApi` | `client.taurusNetwork.pledgeApi` | Pledge lifecycle operations |
| `TaurusNetworkLendingApi` | `client.taurusNetwork.lendingApi` | Lending offers and agreements |
| `TaurusNetworkSettlementApi` | `client.taurusNetwork.settlementApi` | Settlement operations |
| `TaurusNetworkSharedAddressAssetApi` | `client.taurusNetwork.sharedAddressAssetApi` | Address and asset sharing |

See [Services Reference](docs/SERVICES.md) for complete API documentation.

## Basic Usage

### List Wallets

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key',
  apiSecret: 'your-api-secret-hex',
});

// List wallets with pagination
const { items: wallets, pagination } = await client.wallets.list({ limit: 50 });

for (const wallet of wallets) {
  console.log(`${wallet.name}: ${wallet.currency} (${wallet.blockchain}/${wallet.network})`);
}

if (pagination) {
  console.log(`Total wallets: ${pagination.totalItems}`);
}

client.close();
```

### Create a Wallet

```typescript
const wallet = await client.wallets.create({
  blockchain: 'ETH',
  network: 'mainnet',
  name: 'Trading Wallet',
});
console.log(`Created wallet ID: ${wallet.id}`);
```

### Approve Transaction Requests

```typescript
import { ProtectClient, signData } from '@taurushq/protect-sdk';
import { createPrivateKey } from 'crypto';

// Load your private key
const privateKey = createPrivateKey({
  key: pemPrivateKey,
  format: 'pem',
});

// Get requests pending approval
const { items: requests } = await client.requests.listForApproval({ limit: 10 });

if (requests.length > 0) {
  // Approve with ECDSA signature
  const signedCount = await client.requests.approveRequests(requests, privateKey);
  console.log(`Approved ${signedCount} request(s)`);
}
```

### TaurusNetwork Operations

```typescript
// Get all participants
const participantsResponse = await client.taurusNetwork.participantApi.getAllParticipants();
const participants = participantsResponse.result?.participants ?? [];
for (const p of participants) {
  console.log(`Participant: ${p.name}`);
}

// Get my participant info
const myParticipantResponse = await client.taurusNetwork.participantApi.getMyParticipant();
const me = myParticipantResponse.result;
console.log(`My participant: ${me?.participant?.name}`);

// List pledges
const pledgesResponse = await client.taurusNetwork.pledgeApi.getAllPledges();
const pledges = pledgesResponse.result?.pledges ?? [];
for (const pledge of pledges) {
  console.log(`Pledge ${pledge.id}: ${pledge.status}`);
}

// List shared address assets
const sharedResponse = await client.taurusNetwork.sharedAddressAssetApi.getAllSharedAddressAssets();
const sharedAssets = sharedResponse.result?.sharedAddressAssets ?? [];
```

### Error Handling

```typescript
import {
  ProtectClient,
  APIError,
  IntegrityError,
  NotFoundError,
  RateLimitError,
} from '@taurushq/protect-sdk';

try {
  const wallet = await client.wallets.get(999999);
} catch (error) {
  if (error instanceof NotFoundError) {
    console.log('Wallet not found');
  } else if (error instanceof RateLimitError) {
    const delay = error.suggestedRetryDelayMs();
    console.log(`Rate limited, retry after ${delay}ms`);
    await new Promise(resolve => setTimeout(resolve, delay));
  } else if (error instanceof IntegrityError) {
    // Security error - DO NOT retry
    console.error(`Integrity verification failed: ${error.message}`);
  } else if (error instanceof APIError) {
    if (error.isRetryable()) {
      console.log(`Retryable error: ${error.message}`);
    } else {
      console.log(`Non-retryable error: ${error.message}`);
    }
  }
}
```

See [Usage Examples](docs/USAGE_EXAMPLES.md) for comprehensive examples.

## Build Commands

```bash
# Build and test (default) - USE THIS TO VERIFY CHANGES
./build.sh

# Build only
./build.sh build

# Run tests only
./build.sh test

# Run linter
./build.sh lint

# Code generation (OpenAPI + protobuf)
./build.sh generate

# Clean build artifacts
./build.sh clean
```

## Development

### Running Tests

```bash
# Run all unit tests (includes coverage)
./build.sh unit

# Run a specific test by name pattern
./build.sh unit-one "test-tpv1"

# Run a specific test case
./build.sh unit-one "should sign request"
```

### Integration Tests

Integration tests require environment configuration:

```bash
export PROTECT_INTEGRATION_TEST=true
export PROTECT_API_HOST="https://your-api-host.com"
export PROTECT_API_KEY="your-api-key"
export PROTECT_API_SECRET="your-hex-encoded-secret"

./build.sh integration
```

## Low-Level API Access

In addition to the high-level services, the SDK provides direct access to 56 OpenAPI-generated APIs:

```typescript
// Direct API access for advanced use cases
const response = await client.walletsApi.walletServiceGetWalletsV2();
const addresses = await client.addressesApi.addressServiceGetAddresses();

// TaurusNetwork APIs
const participants = await client.taurusNetwork.participantApi.getAllParticipants();
const pledges = await client.taurusNetwork.pledgeApi.getAllPledges();
```

## License

This code is copyright (c) 2025 Taurus SA. It is released under the [MIT license](./LICENSE).
