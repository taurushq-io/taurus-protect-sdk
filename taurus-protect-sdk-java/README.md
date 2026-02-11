# Taurus-PROTECT Java SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Java SDK for interacting with the Taurus-PROTECT API, providing secure cryptocurrency custody and transaction management capabilities.

## Documentation

| Document | Description |
|----------|-------------|
| [Key Concepts](docs/CONCEPTS.md) | Domain model, entities, and workflows |
| [SDK Overview](docs/SDK_OVERVIEW.md) | Architecture, modules, and design patterns |
| [Authentication](docs/AUTHENTICATION.md) | TPV1 authentication and cryptographic operations |
| [Services Reference](docs/SERVICES.md) | Complete API documentation for all services |
| [Usage Examples](docs/USAGE_EXAMPLES.md) | Code examples and common patterns |
| [Whitelisted Address Verification](docs/WHITELISTED_ADDRESS_VERIFICATION.md) | Detailed verification flow |

## Quick Start

### Prerequisites

- Java 8 or higher (tested with Corretto 8)
- Maven 3.6+

### Installation

Build and install the SDK locally:

```bash
mvn clean install
```

Add the dependency to your project:

```xml
<dependency>
    <groupId>com.taurushq.sdk</groupId>
    <artifactId>taurus-protect-sdk-client</artifactId>
    <version>1.0-SNAPSHOT</version>
</dependency>
```

### Client Initialization

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import java.util.Arrays;

String host = "https://api.protect.taurushq.com";
String apiKey = "your-api-key";
String apiSecret = "your-api-secret";

List<String> superAdminPublicKeysPem = Arrays.asList(
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----"
);

ProtectClient client = ProtectClient.createFromPem(
    host,
    apiKey,
    apiSecret,
    superAdminPublicKeysPem,
    2  // minimum valid signatures
);
```

See [Authentication](docs/AUTHENTICATION.md) for more initialization options.

## Services

The SDK provides 43 services organized into core services and the TaurusNetwork namespace.

### Core Services (38 services)

| Service | Purpose |
|---------|---------|
| `WalletService` | Wallet creation, retrieval, balance history |
| `AddressService` | Address management, proof of reserve |
| `RequestService` | Transaction requests and approvals |
| `TransactionService` | Transaction queries and export |
| `BalanceService` | Balance queries across assets |
| `CurrencyService` | Currency and blockchain information |
| `GovernanceRuleService` | Governance rules with signature verification |
| `WhitelistedAddressService` | Address whitelisting with verification |
| `WhitelistedAssetService` | Asset/contract whitelisting |
| `AuditService` | Audit log queries |
| `ChangeService` | Configuration change tracking |
| `FeeService` | Transaction fee information |
| `PriceService` | Price data and history |
| `AirGapService` | Air-gap signing operations |
| `StakingService` | Multi-chain staking information |
| `ContractWhitelistingService` | Smart contract whitelisting |
| `BusinessRuleService` | Business rule management |
| `ReservationService` | Balance reservations |
| `MultiFactorSignatureService` | Multi-factor signature operations |
| `UserService` | User management |
| `GroupService` | User group management |
| `VisibilityGroupService` | Visibility group management |
| `ConfigService` | System configuration |
| `WebhookService` | Webhook management |
| `WebhookCallsService` | Webhook call history |
| `TagService` | Tag management |
| `AssetService` | Asset information |
| `ActionService` | Action management |
| `BlockchainService` | Blockchain information |
| `ExchangeService` | Exchange integration |
| `FiatService` | Fiat currency operations |
| `FeePayerService` | Fee payer management |
| `HealthService` | API health checks |
| `JobService` | Background job management |
| `ScoreService` | Risk scoring operations |
| `StatisticsService` | Platform statistics |
| `TokenMetadataService` | Token metadata information |
| `UserDeviceService` | User device management |

### TaurusNetwork Services (5 services)

| Service | Access | Purpose |
|---------|--------|---------|
| `TaurusNetworkParticipantService` | `client.taurusNetwork().participants()` | Participant management |
| `TaurusNetworkPledgeService` | `client.taurusNetwork().pledges()` | Pledge lifecycle operations |
| `TaurusNetworkLendingService` | `client.taurusNetwork().lending()` | Lending offers and agreements |
| `TaurusNetworkSettlementService` | `client.taurusNetwork().settlements()` | Settlement operations |
| `TaurusNetworkSharingService` | `client.taurusNetwork().sharing()` | Address and asset sharing |

See [Services Reference](docs/SERVICES.md) for complete API documentation.

## Basic Usage

```java
// Create a wallet
Wallet wallet = client.getWalletService().createWallet("ETH", "mainnet", "My Wallet", false);

// Create an address
Address address = client.getAddressService().createAddress(wallet.getId(), "Payment Address", "", "");

// Create and approve a transfer request
Request request = client.getRequestService().createInternalTransferRequest(
    fromAddressId, toAddressId, BigInteger.valueOf(1000000)
);

PrivateKey privateKey = CryptoTPV1.decodePrivateKey(pemPrivateKey);
client.getRequestService().approveRequest(request, privateKey);
```

See [Usage Examples](docs/USAGE_EXAMPLES.md) for comprehensive examples.

## Build Commands

```bash
./build.sh                                    # Full build (compile + test)
./build.sh build                              # Compile only
./build.sh unit                               # Unit tests only
./build.sh integration                        # Integration tests (requires API access)
./build.sh unit-one "TestClassName#methodName" # Single test
```

## Code Generation

```bash
./scripts/generate-openapi.sh   # Regenerate OpenAPI client (requires Java 11+)
./scripts/generate-proto.sh     # Regenerate protobuf classes
```

## License

This code is copyright (c) 2025 Taurus SA. It is released under the [MIT license](./LICENSE).
