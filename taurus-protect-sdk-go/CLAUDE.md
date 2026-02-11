# CLAUDE.md -- Go SDK

## Naming Conventions

**Taurus Product Names**: Always use hyphenated format: `Taurus-PROTECT`, `Taurus-CAPITAL`, `Taurus-EXPLORER`, `Taurus-PRIME`.

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: build + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Build only
./build.sh lint      # golangci-lint
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
./build.sh e2e       # Run E2E tests (requires API access)
./build.sh e2e-one <pattern>  # Run a single E2E test
```

**Single test:** `./build.sh unit-one <pattern>` (e.g., `TestMapWallet`, `Test.*Request`)

## Architecture

### Package Structure
- **internal/openapi**: Auto-generated OpenAPI client (DO NOT MODIFY)
- **internal/proto**: Auto-generated protobuf classes (DO NOT MODIFY)
- **pkg/protect**: Public SDK package
  - **service/**: Service layer wrapping OpenAPI calls
  - **model/**: Domain models exposed to users
  - **mapper/**: DTO to model conversion functions
  - **helper/**: Signature verification, validation utilities
  - **cache/**: Thread-safe caching (rules container)
  - **crypto/**: TPV1 authentication and cryptographic utilities
- **test/integration**: Integration tests

### Key Patterns
- Services accept `context.Context` as first parameter
- Errors use Go 1.13+ wrapping (`errors.Is`, `errors.As`)
- Functional options pattern for client configuration
- `internal/` packages not importable by external code
- Lazy initialization via double-checked locking with `sync.RWMutex`
- TPV1-HMAC-SHA256 authentication handled by HTTP transport middleware

### Available Services (38 + TaurusNetwork)

**Core**: `Wallets()`, `Addresses()`, `Requests()`, `Transactions()`, `GovernanceRules()`, `Balances()`, `Currencies()`, `WhitelistedAddresses()`, `WhitelistedAssets()`

**Transaction/Request**: `Audits()`, `Changes()`, `Fees()`, `Prices()`

**Advanced**: `AirGap()`, `Staking()`, `WhitelistedContracts()`, `BusinessRules()`, `Reservations()`, `MultiFactorSignature()`

**Administrative**: `Users()`, `Groups()`, `VisibilityGroups()`, `Config()`, `Webhooks()`, `WebhookCalls()`, `Tags()`

**Specialized**: `Assets()`, `Actions()`, `Blockchains()`, `Exchanges()`, `Fiat()`, `FeePayers()`, `Health()`, `Jobs()`, `Scores()`, `Statistics()`, `TokenMetadata()`, `UserDevices()`

**Taurus Network** (namespace): `client.TaurusNetwork().Participants()`, `.Pledges()`, `.Lending()`, `.Settlements()`, `.Sharing()`

SDK alignment: see `docs/SDK_ALIGNMENT_REPORT.md` (repository root). Java SDK is source of truth.

## Code Generation

### OpenAPI Generator
Uses `openapi-generator-cli` JAR with `-g go`:
- `enumClassPrefix=true` -- avoids enum constant redeclaration conflicts
- Types prefixed with `Tgvalidatord` (e.g., `TgvalidatordWallet`, `TgvalidatordAddress`)
- Response types use `.Result` field; create operations often return only an ID
- Pagination: `TotalItems`/`Offset` strings; cursor: `Cursor.CurrentPage`/`HasPrevious`/`HasNext`
- Request builders use `.Body(req)` for POST/PUT
- `TgvalidatordWallet` (create) vs `TgvalidatordWalletInfo` (get/list) -- different types
- **Body struct fields use `string`, not `*string`** -- use `Field: value` not `Field: &value`
- **Field names may differ between Get and List responses** -- check actual OpenAPI model

### WhitelistedAddress Metadata Payload Mapping
Fields from `dto.Metadata.Payload` (`map[string]interface{}`):
`"label"`, `"address"`, `"memo"`, `"customerId"`, `"contractType"`, `"addressType"`, `"exchangeAccountId"` (string->int64), `"linkedInternalAddresses"` (`[]{id, label}`), `"linkedWallets"` (`[]{id, path, label}`)

### Protobuf Generator
Uses `protoc` with `protoc-gen-go`. Requires M mappings for go_package. Proto files flattened to `internal/proto/`.

## Service Implementation Pattern

1. Wrap OpenAPI API service (e.g., `openapi.WalletsAPIService`)
2. Use `ErrorMapper` for converting OpenAPI errors to domain errors
3. Use mapper functions to convert DTOs to domain models
4. Return pagination info when available

### Helper Functions (Mapper Utilities)
All safe pointer helpers in `pkg/protect/mapper/helpers.go`: `safeString`, `safeBool`, `safeInt64`, `safeFloat32`, `safeTime`, `stringPtr`, `boolPtr`, `int64Ptr`.

**Important:**
- **service package**: `stringPtr()` defined separately (packages can't share unexported functions)
- In tests, define `testStringPtr()` to avoid conflicts with production helpers

## Testing

### Running Tests
```bash
./build.sh unit                      # Preferred
go test ./pkg/protect/...            # Direct
go test -v -cover ./pkg/protect/...  # Verbose + coverage
```

### Test Structure
- `pkg/protect/mapper/*_test.go` -- DTO to model mapping
- `pkg/protect/service/*_test.go` -- Input validation, error mapping
- `pkg/protect/cache/rules_container_test.go` -- Cache operations, concurrency
- `pkg/protect/crypto/tpv1_test.go` -- TPV1 auth, HMAC, signatures
- `pkg/protect/transport_test.go` -- HTTP transport with TPV1 signing

### Integration Tests
Located in `test/integration/`, disabled by default. Enable via env vars (`PROTECT_INTEGRATION_TEST=true`, `PROTECT_API_HOST`, `PROTECT_API_KEY`, `PROTECT_API_SECRET`) or hard-coded defaults in `test/integration/config.go`.

```bash
./build.sh integration                                      # All
go test -v -timeout 15m ./test/integration/... -run Wallet  # Specific
```

26 tests covering: wallets, addresses, requests, transactions, users, balances, currencies, governance, health, groups, config, blockchains, whitelisting, tags, statistics.

### Shared Test Utilities (`test/testutil/`)

All test config is centralized in `test/testutil/` (package `testutil`):
- `properties.go` — Key=value `.properties` file parser with `\n` escape for PEM keys
- `config.go` — `TestConfig` with multi-identity support (6 identities), env var overrides
- `helpers.go` — `SkipIfNotEnabled()`, `GetTestClient(t, index)`, `SkipIfInsufficientIdentities()`
- `test.properties.sample` — Sample config matching Java format

Integration and E2E directories delegate to testutil via thin wrappers.

### E2E Tests
Located in `test/e2e/`, delegates config to `test/testutil/`.

```bash
./build.sh e2e                                                 # All
./build.sh e2e-one TestIntegration_MultiCurrencyE2E            # Specific
```

**Timeout:** Default 600s is insufficient. Use `-timeout 15m` for integration, `-timeout 30m` for E2E (`build.sh` does this automatically).

## Build Troubleshooting

Integration test failures (e.g., "cannot unmarshal array") are typically API schema mismatches, not code bugs. **Unit tests are the reliable verification:** `go test ./pkg/protect/... -v`

## Lessons Learned (Non-Security)

### Pagination Overflow Prevention
Use overflow-safe comparison: `totalItems > offset && totalItems - offset > limit` (not `offset + limit < totalItems`).

### Helper Function Consolidation
All safe pointer helpers consolidated in `pkg/protect/mapper/helpers.go`.

### Go Module Version Policy
Target oldest supported Go version (1.21+). Remove `toolchain` directive when lowering.

### Request Body Preservation in Transport
Read and restore `req.Body` before cloning -- cloning consumes the body, breaking retry middleware. Pattern: `ReadAll` -> reset original with `NopCloser(NewReader)` -> clone with fresh `NopCloser(NewReader)`.
