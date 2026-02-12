# Taurus PROTECT Java SDK Overview

## Introduction

The Taurus PROTECT Java SDK provides a high-level Java client for interacting with the Taurus PROTECT API. It handles authentication, request signing, cryptographic verification of responses, and provides clean domain objects for all API operations.

## Architecture

### Module Structure

The SDK is organized into three Maven modules:

```
taurus-protect-sdk-java/
├── client/     # High-level SDK (main development target)
├── openapi/    # Auto-generated OpenAPI client (DO NOT MODIFY)
└── proto/      # Auto-generated Protobuf classes (DO NOT MODIFY)
```

| Module | Purpose | Auto-Generated | Modifiable |
|--------|---------|----------------|------------|
| **client** | High-level Java SDK with business logic | No | Yes |
| **openapi** | Low-level HTTP API client from OpenAPI specs | Yes | No |
| **proto** | Protocol Buffer classes for binary messages | Yes | No |

### Module Dependencies

```
┌─────────────────────────────────────────────────────────┐
│                      client                              │
│  (High-level SDK - Business Logic & Verification)       │
├─────────────────────────────────────────────────────────┤
│                         │                                │
│              ┌──────────┴──────────┐                    │
│              ▼                     ▼                    │
│        ┌──────────┐          ┌──────────┐              │
│        │ openapi  │          │  proto   │              │
│        │(HTTP API)│          │(Protobuf)│              │
│        └──────────┘          └──────────┘              │
└─────────────────────────────────────────────────────────┘
```

## Module Details

### Client Module

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/`

The client module is the main SDK that developers interact with. It provides:

- **ProtectClient** - Entry point and service factory
- **Services** - Business logic layer (43 services: 38 core + 5 TaurusNetwork)
- **Models** - Clean domain objects
- **Mappers** - MapStruct interfaces for DTO conversion
- **Helpers** - Cryptographic verification utilities

**Package Structure:**
```
com.taurushq.sdk.protect.client
├── ProtectClient.java          # Entry point
├── service/                    # Service layer
│   ├── WalletService.java
│   ├── AddressService.java
│   ├── RequestService.java
│   └── ... (43 services total: 38 core + 5 TaurusNetwork)
├── model/                      # Domain models
│   ├── Wallet.java
│   ├── Address.java
│   ├── Request.java
│   └── ... (55+ models)
├── mapper/                     # MapStruct interfaces
│   ├── WalletMapper.java
│   ├── AddressMapper.java
│   └── ...
├── helper/                     # Verification utilities
│   ├── SignatureVerifier.java
│   ├── WhitelistHashHelper.java
│   └── ...
└── cache/                      # Caching
    └── RulesContainerCache.java
```

### OpenAPI Module

**Location:** `openapi/src/main/java/com/taurushq/sdk/protect/openapi/`

Auto-generated from OpenAPI specifications. Contains:

- HTTP client implementation (OkHttp3)
- API endpoint classes
- DTO classes (Data Transfer Objects)
- Authentication handlers (TPV1)

**Regeneration:**
```bash
./scripts/generate-openapi.sh  # Requires Java 11+
```

### Proto Module

**Location:** `proto/src/main/java/`

Auto-generated Protocol Buffer classes for binary message handling.

**Regeneration:**
```bash
./scripts/generate-proto.sh  # Requires protoc
```

## Key Design Patterns

### 1. Facade Pattern (ProtectClient)

`ProtectClient` serves as the main entry point, instantiating and providing access to all services:

```java
ProtectClient client = ProtectClient.createFromPem(host, apiKey, apiSecret, superAdminKeys, 2);

// Access services through the facade
client.getWalletService().createWallet(...);
client.getAddressService().getAddress(...);
client.getRequestService().approveRequest(...);
```

### 2. Service Layer Pattern

Each service wraps OpenAPI calls with:
- Input validation
- Exception translation
- Cryptographic verification (where applicable)
- Response mapping to domain objects

```java
// Service layer flow
Request → Service → OpenAPI Client → API
                         ↓
Domain Object ← Mapper ← DTO ← Response
```

### 3. MapStruct Mappers

MapStruct generates mapper implementations at compile time:

```java
// Mapper interface
@Mapper
public interface WalletMapper {
    WalletMapper INSTANCE = Mappers.getMapper(WalletMapper.class);
    Wallet fromDTO(TgvalidatordWallet dto);
}

// Usage in service
Wallet wallet = WalletMapper.INSTANCE.fromDTO(apiResponse);
```

Generated implementations are located in `target/generated-sources/`.

### 4. Exception Handling Pattern

Services catch OpenAPI exceptions and rethrow as client exceptions:

```java
try {
    return walletApi.getWallet(walletId);
} catch (com.taurushq.sdk.protect.openapi.ApiException e) {
    throw apiExceptionMapper.toApiException(e);
}
```

**Exception Hierarchy:**
- `ApiException` (checked) - General API errors
  - `AuthenticationException` - Authentication failures (401)
  - `AuthorizationException` - Authorization failures (403)
  - `NotFoundException` - Resource not found (404)
  - `RateLimitException` - Rate limit exceeded (429)
  - `ServerException` - Server errors (5xx)
  - `ValidationException` - Input validation errors (400)
- `IntegrityException` (unchecked, extends `SecurityException`) - Hash/signature verification failures
- `WhitelistException` (checked) - Whitelist-specific verification failures
- `ConfigurationException` (checked) - Client configuration errors
- `RequestMetadataException` (checked) - Metadata payload parsing errors

## Build & Development

### Build Commands

```bash
# Full build (compile + test) - USE THIS TO VERIFY CHANGES
./build.sh

# Compile only
./build.sh build

# Run unit tests only
./build.sh unit

# Run integration tests (requires API access)
./build.sh integration

# Run a single test
./build.sh unit-one "RequestStatusTest#testFromStatus"

# Full verification (compile + test + static analysis)
./build.sh verify

# Run static analysis (SpotBugs, PMD, Checkstyle)
./build.sh lint

# Clean build artifacts
./build.sh clean
```

### Advanced: Direct Maven Commands

```bash
# Run specific test class
mvn test -Dtest=RequestStatusTest

# Run specific test method
mvn test -Dtest=RequestStatusTest#testFromStatus

# Install locally (skip checks)
mvn clean install -DskipTests

# Run static analysis directly
mvn verify
```

### Code Generation

```bash
# Regenerate OpenAPI client (requires Java 11+ runtime)
./scripts/generate-openapi.sh

# Regenerate protobuf classes (requires protoc)
./scripts/generate-proto.sh
```

### Static Analysis

The project enforces code quality through:

| Tool | Purpose | Config File |
|------|---------|-------------|
| SpotBugs | Bug detection | `spotbugs-exclude.xml` |
| PMD | Code complexity | Excludes openapi, generated mappers |
| Checkstyle | Style enforcement | `checkstyle.xml` |

The openapi module is excluded from static analysis checks.

## Key Classes Reference

| Class | Purpose |
|-------|---------|
| `ProtectClient` | Entry point, service factory |
| `CryptoTPV1` | Cryptographic utilities (signing, hashing) |
| `ApiExceptionMapper` | Exception translation |
| `RulesContainerCache` | Governance rules caching |
| `SignatureVerifier` | Signature verification |

## Related Documentation

- [Authentication](AUTHENTICATION.md) - TPV1 authentication and cryptographic operations
- [Services Reference](SERVICES.md) - Complete service API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Detailed verification logic
