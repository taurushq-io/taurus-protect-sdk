# Taurus-PROTECT Go SDK Overview

## Introduction

The Taurus-PROTECT Go SDK provides a high-level Go client for interacting with the Taurus-PROTECT API. It handles authentication, request signing, cryptographic verification of responses, and provides clean domain objects for all API operations.

## Architecture

### Package Structure

The SDK is organized into the following packages:

```
taurus-protect-sdk-go/
├── pkg/protect/                    # Main SDK package
│   ├── client.go                   # Client entry point
│   ├── options.go                  # Functional options
│   ├── transport.go                # HTTP transport with TPV1
│   ├── errors.go                   # Error types
│   ├── service/                    # Service implementations (43 services: 38 core + 5 TaurusNetwork)
│   ├── model/                      # Domain models (46 models)
│   ├── mapper/                     # DTO converters (83 mappers)
│   ├── cache/                      # Rules container caching
│   └── crypto/                     # Cryptographic utilities
├── internal/
│   ├── openapi/                    # Auto-generated OpenAPI client (DO NOT MODIFY)
│   └── proto/                      # Auto-generated Protobuf classes (DO NOT MODIFY)
└── scripts/
    ├── generate-openapi.sh         # OpenAPI code generation
    └── generate-proto.sh           # Protobuf code generation
```

| Package | Purpose | Auto-Generated | Modifiable |
|---------|---------|----------------|------------|
| **pkg/protect** | High-level Go SDK with business logic | No | Yes |
| **internal/openapi** | Low-level HTTP API client from OpenAPI specs | Yes | No |
| **internal/proto** | Protocol Buffer classes for binary messages | Yes | No |

### Package Dependencies

```
┌─────────────────────────────────────────────────────────┐
│                    pkg/protect                           │
│  (High-level SDK - Business Logic & Verification)       │
├─────────────────────────────────────────────────────────┤
│                         │                                │
│              ┌──────────┴──────────┐                    │
│              ▼                     ▼                    │
│     ┌──────────────┐       ┌──────────────┐            │
│     │internal/     │       │internal/     │            │
│     │openapi       │       │proto         │            │
│     │(HTTP API)    │       │(Protobuf)    │            │
│     └──────────────┘       └──────────────┘            │
└─────────────────────────────────────────────────────────┘
```

## Package Details

### Client Package

**Location:** `pkg/protect/`

The client package is the main SDK that developers interact with. It provides:

- **Client** - Entry point and service factory
- **Services** - Business logic layer (43 services: 38 core + 5 TaurusNetwork)
- **Models** - Clean domain objects
- **Mappers** - Functions for DTO conversion
- **Crypto** - Cryptographic verification utilities

**Package Structure:**
```
pkg/protect/
├── client.go              # Entry point, lazy service initialization
├── options.go             # Functional options for configuration
├── transport.go           # HTTP transport with TPV1 signing
├── errors.go              # Error types and helpers
├── service/               # Service layer
│   ├── wallet.go
│   ├── address.go
│   ├── request.go
│   └── ... (43 services total: 38 core + 5 TaurusNetwork)
├── model/                 # Domain models
│   ├── wallet.go
│   ├── address.go
│   ├── request.go
│   └── ... (46 models total)
├── mapper/                # DTO converters
│   ├── wallet.go
│   ├── address.go
│   └── ... (83 mappers)
├── cache/                 # Caching
│   └── rules_container.go
└── crypto/                # Cryptographic utilities
    └── tpv1.go
```

### Internal OpenAPI Package

**Location:** `internal/openapi/`

Auto-generated from OpenAPI specifications. Contains:

- HTTP client implementation
- API endpoint methods
- DTO types (Data Transfer Objects)
- Request/response handling

**Regeneration:**
```bash
./scripts/generate-openapi.sh  # Requires Java 11+
```

### Internal Proto Package

**Location:** `internal/proto/`

Auto-generated Protocol Buffer classes for binary message handling.

**Regeneration:**
```bash
./scripts/generate-proto.sh  # Requires protoc + protoc-gen-go
```

## Key Design Patterns

### 1. Functional Options Pattern

The SDK uses functional options for client configuration:

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

client, err := protect.NewClient(
    "https://api.taurus-protect.com",
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeysPEM(pemKeys),
    protect.WithMinValidSignatures(2),
    protect.WithRulesCacheTTL(10 * time.Minute),
)
```

### 2. Lazy Service Initialization

Services are initialized on first access using double-checked locking:

```go
// Services are accessed through client methods
wallets := client.Wallets()       // Initialized on first call
addresses := client.Addresses()   // Initialized on first call
requests := client.Requests()     // Initialized on first call
```

### 3. Service Layer Pattern

Each service wraps OpenAPI calls with:
- Input validation
- Error mapping
- Response mapping to domain objects

```go
// Service layer flow
Request → Service → OpenAPI Client → API
                         ↓
Domain Object ← Mapper ← DTO ← Response
```

### 4. Mapper Functions

Mappers convert between OpenAPI DTOs and domain models:

```go
// Mapper function signature
func WalletFromDTO(dto *openapi.TgvalidatordWalletInfo) *model.Wallet

// Usage in service
wallet := mapper.WalletFromDTO(apiResponse.Result)
```

### 5. Error Handling Pattern

Services map OpenAPI errors to domain errors:

```go
resp, httpResp, err := s.api.WalletServiceGetWalletV2(ctx, walletID).Execute()
if err != nil {
    return nil, s.errMapper.MapError(err, httpResp)
}
```

**Error Hierarchy:**
- `APIError` - General API errors (HTTP status, error code, message)
- `IntegrityError` - Hash/signature verification failures
- `WhitelistError` - Whitelist-specific verification failures
- `RequestMetadataError` - Request metadata validation failures

## Build & Development

### Build Commands

```bash
# Build and test (default) - USE THIS TO VERIFY CHANGES
./build.sh

# Build only
./build.sh build

# Run unit tests (includes verbose output and coverage)
./build.sh unit

# Run a single unit test by name pattern
./build.sh unit-one TestMapWallet

# Run integration tests (requires API access)
./build.sh integration

# Run a single integration test by name pattern
./build.sh integration-one TestListWallets

# Run linter (requires golangci-lint)
./build.sh lint

# Code generation (OpenAPI + protobuf)
./build.sh generate

# Clean build artifacts
./build.sh clean
```

### Advanced: Direct Go Commands

```bash
# Build the project
go build ./...

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/protect/...

# Run linting
golangci-lint run
```

### Code Generation

```bash
# Regenerate OpenAPI client (requires Java 11+ runtime)
./scripts/generate-openapi.sh

# Regenerate protobuf classes (requires protoc + protoc-gen-go)
./scripts/generate-proto.sh
```

### Installing Protobuf Tools

```bash
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Key Types Reference

| Type | Purpose |
|------|---------|
| `Client` | Entry point, service factory |
| `TPV1Auth` | Cryptographic utilities (signing, hashing) |
| `ErrorMapper` | Error translation |
| `RulesContainerCache` | Governance rules caching |
| `APIError` | Structured API error |

## Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `WithCredentials` | Required | API key and secret |
| `WithSuperAdminKeysPEM` | Optional | SuperAdmin public keys in PEM format |
| `WithSuperAdminKeys` | Optional | SuperAdmin public keys as `*ecdsa.PublicKey` |
| `WithMinValidSignatures` | 1 | Minimum SuperAdmin signatures required |
| `WithRulesCacheTTL` | 5 min | Rules container cache TTL |
| `WithHTTPClient` | Default | Custom HTTP client |
| `WithHTTPTimeout` | 30s | HTTP request timeout |

## Thread Safety

The SDK is designed for concurrent use:

- `Client` uses `sync.RWMutex` for lazy service initialization
- `RulesContainerCache` uses `sync.RWMutex` for thread-safe caching
- Services are stateless after initialization
- HTTP client is shared and thread-safe

## Resource Cleanup

The client implements `io.Closer`:

```go
client, err := protect.NewClient(host, opts...)
if err != nil {
    return err
}
defer client.Close()

// Use client...
```

## Related Documentation

- [Key Concepts](CONCEPTS.md) - Domain model and entities
- [Authentication](AUTHENTICATION.md) - TPV1 authentication and cryptographic operations
- [Services Reference](SERVICES.md) - Complete service API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Detailed verification logic