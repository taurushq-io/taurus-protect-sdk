# Taurus-PROTECT Go SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go SDK for interacting with the Taurus-PROTECT API, providing secure cryptocurrency custody and transaction management capabilities.

## Documentation

| Document | Description |
|----------|-------------|
| [Key Concepts](docs/CONCEPTS.md) | Go model types, exceptions, and domain concepts |
| [SDK Overview](docs/SDK_OVERVIEW.md) | Architecture, packages, and design patterns |
| [Authentication](docs/AUTHENTICATION.md) | TPV1 authentication and cryptographic operations |
| [Services Reference](docs/SERVICES.md) | Complete API documentation for all 43 services |
| [Usage Examples](docs/USAGE_EXAMPLES.md) | Go code examples and common patterns |
| [Whitelisted Address Verification](docs/WHITELISTED_ADDRESS_VERIFICATION.md) | 6-step verification flow |

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Go modules enabled

### Installation

Install the SDK using Go modules:

```bash
go get github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go
```

Or add to your `go.mod`:

```go
require github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go v1.0.0
```

### Dependencies

The SDK requires the following packages (installed automatically):

| Package | Purpose |
|---------|---------|
| `google.golang.org/protobuf` | Protocol buffer support |
| Standard library | HTTP client, crypto, encoding |

### Client Initialization

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

// Basic initialization
client, err := protect.NewClient(
    "https://api.protect.taurushq.com",
    protect.WithCredentials(apiKey, apiSecret),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Use the client
wallets, _, err := client.Wallets().ListWallets(ctx, nil)
```

For production use with SuperAdmin key verification:

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

superAdminKeys := []string{
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----",
}

client, err := protect.NewClient(
    "https://api.protect.taurushq.com",
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeysPEM(superAdminKeys),
    protect.WithMinValidSignatures(2),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// All governance rule verifications will require 2 valid signatures
```

See [Authentication](docs/AUTHENTICATION.md) for more initialization options.

## Services

The SDK provides 43 services organized into core services and the TaurusNetwork namespace.

### Core Services (38 services)

| Service | Access | Purpose |
|---------|--------|---------|
| `WalletService` | `client.Wallets()` | Wallet creation, retrieval, balance history |
| `AddressService` | `client.Addresses()` | Address management, proof of reserve |
| `RequestService` | `client.Requests()` | Transaction requests and approvals |
| `TransactionService` | `client.Transactions()` | Transaction queries and export |
| `BalanceService` | `client.Balances()` | Balance queries across assets |
| `CurrencyService` | `client.Currencies()` | Currency and blockchain information |
| `GovernanceRuleService` | `client.GovernanceRules()` | Governance rules with signature verification |
| `WhitelistedAddressService` | `client.WhitelistedAddresses()` | Address whitelisting with verification |
| `WhitelistedAssetService` | `client.WhitelistedAssets()` | Asset/contract whitelisting |
| `AuditService` | `client.Audits()` | Audit log queries |
| `ChangeService` | `client.Changes()` | Configuration change tracking |
| `FeeService` | `client.Fees()` | Transaction fee information |
| `PriceService` | `client.Prices()` | Price data and history |
| `AirGapService` | `client.AirGap()` | Air-gap signing operations |
| `StakingService` | `client.Staking()` | Multi-chain staking information |
| `WhitelistedContractService` | `client.WhitelistedContracts()` | Smart contract whitelisting |
| `BusinessRuleService` | `client.BusinessRules()` | Business rule management |
| `ReservationService` | `client.Reservations()` | Balance reservations |
| `MultiFactorSignatureService` | `client.MultiFactorSignature()` | Multi-factor signature operations |
| `UserService` | `client.Users()` | User management |
| `GroupService` | `client.Groups()` | User group management |
| `VisibilityGroupService` | `client.VisibilityGroups()` | Visibility group management |
| `ConfigService` | `client.Config()` | System configuration |
| `WebhookService` | `client.Webhooks()` | Webhook management |
| `WebhookCallService` | `client.WebhookCalls()` | Webhook call history |
| `TagService` | `client.Tags()` | Tag management |
| `AssetService` | `client.Assets()` | Asset information |
| `ActionService` | `client.Actions()` | Action management |
| `BlockchainService` | `client.Blockchains()` | Blockchain information |
| `ExchangeService` | `client.Exchanges()` | Exchange integration |
| `FiatService` | `client.Fiat()` | Fiat currency operations |
| `FeePayerService` | `client.FeePayers()` | Fee payer management |
| `HealthService` | `client.Health()` | API health checks |
| `JobService` | `client.Jobs()` | Background job management |
| `ScoreService` | `client.Scores()` | Risk scoring operations |
| `StatisticsService` | `client.Statistics()` | Platform statistics |
| `TokenMetadataService` | `client.TokenMetadata()` | Token metadata information |
| `UserDeviceService` | `client.UserDevices()` | User device management |

### TaurusNetwork Services (5 services)

| Service | Access | Purpose |
|---------|--------|---------|
| `ParticipantService` | `client.TaurusNetwork().Participants()` | Participant management |
| `PledgeService` | `client.TaurusNetwork().Pledges()` | Pledge lifecycle operations |
| `LendingService` | `client.TaurusNetwork().Lending()` | Lending offers and agreements |
| `SettlementService` | `client.TaurusNetwork().Settlements()` | Settlement operations |
| `SharingService` | `client.TaurusNetwork().Sharing()` | Address and asset sharing |

See [Services Reference](docs/SERVICES.md) for complete API documentation.

## Basic Usage

### List Wallets

```go
ctx := context.Background()

// List wallets with pagination
wallets, pagination, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
    Limit: 50,
})
if err != nil {
    return err
}

for _, wallet := range wallets {
    fmt.Printf("%s: %s (%s/%s)\n", wallet.Name, wallet.Currency, wallet.Blockchain, wallet.Network)
}

if pagination != nil {
    fmt.Printf("Total wallets: %d\n", pagination.TotalItems)
}
```

### Create a Wallet

```go
wallet, err := client.Wallets().CreateWallet(ctx, &model.CreateWalletRequest{
    Name:     "Trading Wallet",
    Currency: "ETH",
})
if err != nil {
    return err
}
fmt.Printf("Created wallet ID: %s\n", wallet.ID)
```

### Approve Transaction Requests

```go
import (
    "crypto/ecdsa"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// Load your private key (takes a PEM string, not bytes)
privateKey, err := crypto.DecodePrivateKeyPEM(pemString)
if err != nil {
    return err
}

// Get requests pending approval
requests, _, err := client.Requests().ListRequestsForApproval(ctx, &model.ListRequestsOptions{
    Limit: 10,
})
if err != nil {
    return err
}

if len(requests) > 0 {
    // Approve with ECDSA signature
    signedCount, err := client.Requests().ApproveRequests(ctx, requests, privateKey)
    if err != nil {
        return err
    }
    fmt.Printf("Approved %d request(s)\n", signedCount)
}
```

### TaurusNetwork Operations

```go
// Get my participant info
me, err := client.TaurusNetwork().Participants().GetMyParticipant(ctx)
if err != nil {
    return err
}
fmt.Printf("Participant: %s\n", me.Name)

// List pledges
pledges, cursor, err := client.TaurusNetwork().Pledges().ListPledges(ctx, &taurusnetwork.ListPledgesOptions{
    Limit: 10,
})
if err != nil {
    return err
}
for _, pledge := range pledges {
    fmt.Printf("Pledge %s: %s\n", pledge.ID, pledge.Status)
}

// List shared addresses
addresses, _, err := client.TaurusNetwork().Sharing().ListSharedAddresses(ctx, nil)
if err != nil {
    return err
}
```

### Error Handling

```go
import (
    "errors"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

wallet, err := client.Wallets().GetWallet(ctx, "999999")
if err != nil {
    if errors.Is(err, protect.ErrNotFound) {
        fmt.Println("Wallet not found")
    } else if errors.Is(err, protect.ErrRateLimit) {
        if apiErr, ok := protect.IsAPIError(err); ok && apiErr.RetryAfter > 0 {
            fmt.Printf("Rate limited, retry after %v\n", apiErr.RetryAfter)
            time.Sleep(apiErr.RetryAfter)
        }
    } else if protect.IsIntegrityError(err) {
        // Security error - DO NOT retry
        fmt.Printf("Integrity verification failed: %v\n", err)
    } else if apiErr, ok := protect.IsAPIError(err); ok {
        if apiErr.IsRetryable() {
            fmt.Printf("Retryable error: %s\n", apiErr.Message)
        } else {
            fmt.Printf("Non-retryable error: %s\n", apiErr.Message)
        }
    }
    return err
}
```

See [Usage Examples](docs/USAGE_EXAMPLES.md) for comprehensive examples.

## Build Commands

```bash
# Build and test (default) - USE THIS TO VERIFY CHANGES
./build.sh

# Build only
./build.sh build

# Run unit tests only
./build.sh unit

# Run linter (requires golangci-lint)
./build.sh lint

# Code generation (OpenAPI + protobuf)
./build.sh generate

# Clean build artifacts
./build.sh clean
```

## Development

### Running Tests

```bash
# Run all unit tests (includes verbose output and coverage)
./build.sh unit

# Run a single test by name pattern
./build.sh unit-one TestMapWallet

# Run all integration tests (requires API access, see below)
./build.sh integration

# Run a single integration test by name pattern
./build.sh integration-one TestListWallets
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

## License

This code is copyright (c) 2025 Taurus SA. It is released under the [MIT license](./LICENSE).
