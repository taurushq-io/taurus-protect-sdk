# CLAUDE.md — Java SDK

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: compile + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Compile only
./build.sh verify    # Full verification (compile + test + static analysis)
./build.sh lint      # SpotBugs, PMD, Checkstyle
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
./build.sh install   # Fast local install (skip checks)
./build.sh e2e       # Run E2E tests (requires API access)
./build.sh e2e-one <pattern>  # Run a single E2E test
```

**Single test:** `./build.sh unit-one ClassName#methodName` (e.g., `RequestServiceTest#testApprove`)

## Architecture

Three modules:
- **openapi**: Auto-generated OpenAPI client (DO NOT MODIFY). Uses TPV1 auth (HMAC-based signing).
- **proto**: Auto-generated protobuf classes (DO NOT MODIFY).
- **client**: High-level SDK wrapping openapi. Main development target.

### Client Module Layers

1. **ProtectClient** (`client/.../ProtectClient.java`) — entry point, lazy service initialization
2. **Services** (`client/.../service/*Service.java`) — business logic wrapping OpenAPI calls
3. **Mappers** (`client/.../mapper/*Mapper.java`) — MapStruct interfaces converting DTOs to models
4. **Models** (`client/.../model/*.java`) — clean domain objects for SDK users

### Key Patterns

- Services catch `com.taurushq.sdk.protect.openapi.ApiException` and rethrow as `com.taurushq.sdk.protect.client.model.ApiException`
- MapStruct generates mapper implementations at compile time (implementations are in `target/generated-sources`)
- The `ApiExceptionMapper` extracts structured error info from raw API responses
- TPV1 authentication is handled via `ApiKeyTPV1Auth` which signs requests with HMAC
- Do not use deprecated methods or classes
- Java 8 target — no `var` keyword (Java 10+), use explicit type declarations

### Available Services (38 + TaurusNetwork namespace)

The ProtectClient provides lazy-initialized getters for all services:

**Core Services**: `getWalletService()`, `getAddressService()`, `getRequestService()`, `getTransactionService()`, `getGovernanceRuleService()`, `getBalanceService()`, `getCurrencyService()`, `getWhitelistedAddressService()`, `getWhitelistedAssetService()`

**Transaction/Request Management**: `getAuditService()`, `getChangeService()`, `getFeeService()`, `getPriceService()`

**Advanced Features**: `getAirGapService()`, `getStakingService()`, `getContractWhitelistingService()`, `getBusinessRuleService()`, `getReservationService()`

**Administrative**: `getUserService()`, `getGroupService()`, `getVisibilityGroupService()`, `getConfigService()`, `getWebhookService()`, `getWebhookCallsService()`, `getTagService()`

**Specialized**: `getAssetService()`, `getActionService()`, `getBlockchainService()`, `getExchangeService()`, `getFiatService()`, `getFeePayerService()`, `getHealthService()`, `getJobService()`, `getScoreService()`, `getStatisticsService()`, `getTokenMetadataService()`, `getUserDeviceService()`, `getMultiFactorSignatureService()`

**Taurus Network** (namespace pattern):
```java
client.taurusNetwork().participants()   // Participant management
client.taurusNetwork().pledges()        // Pledge lifecycle
client.taurusNetwork().lending()        // Offers + Agreements
client.taurusNetwork().settlements()    // Settlement operations
client.taurusNetwork().sharing()        // Address/Asset sharing
```

## Code Generation

### OpenAPI Generator
- Uses `openapi-generator-cli` JAR (7.9.0) with `-g java`
- Generated types prefixed with `Tgvalidatord`
- Requires Java 11+ runtime

### Protobuf
- Uses `protoc` directly
- Generated classes flattened to proto module

## Static Analysis

Code must pass SpotBugs, PMD, and Checkstyle. The openapi module is excluded from these checks. Checkstyle config is in
`checkstyle.xml`. PMD excludes generated mappers (`*Impl.java`). Only do static analysis when a user explicitly asks for
it.

**Common PMD/SpotBugs patterns to handle:**

- **Empty catch blocks**: PMD requires a comment inside the catch block body (not just a Javadoc), or use
  `@SuppressWarnings("PMD.EmptyCatchBlock")` for intentional empty catches
- **ConstantsInInterface**: MapStruct mapper interfaces use `INSTANCE` constant - suppress with
  `@SuppressWarnings("PMD.ConstantsInInterface")`
- **Redundant null checks**: SpotBugs flags null checks on fields marked `@Nonnull` in OpenAPI models - respect the API
  contract and remove unnecessary null checks
- **Classes with only private constructors**: Must be marked `final`
- **CPD (Copy-Paste Detector)**: `minimumTokens` threshold in pom.xml controls duplication sensitivity (currently 700)

## Testing

### Test Configuration

Credentials are loaded from `client/src/test/resources/test.properties` (git-ignored), with environment variable overrides. Copy `test.properties.sample` to get started. The `TestConfig` class (in `testutil` package) loads identities with multi-identity support (API creds, private keys, SuperAdmin public keys).

### Integration Tests

Integration tests are located in `client/src/test/java/.../integration/` and are excluded from default test runs via surefire (`**/*IntegrationTest.java`).

**Structure:**
- Shared test utilities live in `testutil/` package: `TestConfig.java` (config) and `TestHelper.java` (helpers like `skipIfNotEnabled()`, `getTestClient()`)
- 15 domain-specific test classes:
  - `WalletIntegrationTest`, `AddressIntegrationTest`, `RequestIntegrationTest`
  - `TransactionIntegrationTest`, `UserIntegrationTest`, `BalanceIntegrationTest`
  - `GovernanceIntegrationTest`, `WhitelistedAddressIntegrationTest`, `WhitelistedAssetIntegrationTest`
  - `HealthIntegrationTest`, `AdminIntegrationTest`, `BlockchainIntegrationTest`, `MiscIntegrationTest`

### E2E Tests

E2E tests are located in `client/src/test/java/.../e2e/` and are excluded from default test runs via surefire (`**/*E2ETest.java`). They reuse `TestConfig` and `TestHelper` from the `testutil` package.

- `MultiCurrencyE2ETest` - Multi-currency parallel transfer lifecycle
- `BusinessRuleChangeE2ETest` - Business rule change proposal/approval lifecycle

**Running:**
```bash
./build.sh integration                              # All integration tests
./build.sh e2e                                      # All E2E tests
./build.sh e2e-one MultiCurrencyE2ETest             # Single E2E test

# With custom credentials
export PROTECT_API_HOST="https://your-api.com"
export PROTECT_API_KEY="your-key"
export PROTECT_API_SECRET="your-secret"
./build.sh e2e
```

**Environment Variables:**
- `PROTECT_INTEGRATION_TEST` - Set to "true" to enable
- `PROTECT_API_HOST` - API host URL
- `PROTECT_API_KEY` - API key
- `PROTECT_API_SECRET` - API secret (hex-encoded)

### Model Field Reference

Key model classes and their actual field names (to avoid compilation errors):

- **User**: `getFirstName()`, `getLastName()`, `getEmail()`, `getId()` (no `getName()`)
- **AuditTrail**: `getEntity()`, `getAction()`, `getDetails()`, `getCreationDate()` (no `getUserName()`)
- **HealthComponent**: `getGroups()` returns Map<String, HealthGroup> (no `getStatus()`)
- **TenantConfig**: `getTenantId()`, `getBaseCurrency()` (no `getName()`)
- **PortfolioStatistics**: `getTotalBalance()`, `getTotalBalanceBaseCurrency()`, `getWalletsCount()`, `getAddressesCount()` (no `getTotalValue()` or `getCurrency()`)
- **Transaction**: `getSources()` and `getDestinations()` return List<AddressInfo> (no `getSourceAddress()`)
- **Tag**: `getValue()` (not `getName()`)
- **BlockchainInfo**: `getSymbol()`, `getNetwork()`, `getName()` (not `getCurrency()`)

**Service Method Signatures:**
- `RequestService.getRequests(OffsetDateTime, OffsetDateTime, String, List<RequestStatus>, ApiRequestCursor)` - uses cursor pagination
- `GroupService.getGroups(String limit, String offset, List<String> ids, List<String> externalGroupIds, String query)` - String params
- `VisibilityGroupService.getVisibilityGroups()` - no pagination parameters
- `AuditService.getAuditTrails(...)` - returns `AuditTrailResult`, not `List<AuditTrail>`
- `ProtectClient.create(...)` - always requires SuperAdmin keys; use `createFromPem()` for PEM-encoded keys

## Lessons Learned (Non-Security)

### Thread-Safe Lazy Initialization

**Problem:** Lazy initialization in `getHsmPublicKey()` and `getDecodedRulesContainer()` has race conditions when multiple threads access simultaneously.

**Solution:** Use synchronized blocks with dedicated lock objects:
```java
private final Object hsmKeyLock = new Object();

public PublicKey getHsmPublicKey() {
    synchronized (hsmKeyLock) {
        if (!hsmPublicKeyResolved) {
            hsmPublicKey = findHsmPublicKey();
            hsmPublicKeyResolved = true;
        }
        return hsmPublicKey;
    }
}
```

### MapStruct Silent Mapping Failures with Non-Standard Setters

**Problem:** MapStruct follows JavaBean conventions and requires `set*()` methods for property mapping. If a model class uses `is*()` as a setter name (e.g., `public void isDisabled(boolean disabled)`), MapStruct treats it as a getter and silently skips the field. The field retains its default value (`false` for booleans), which can cause critical bugs — for example, disabled addresses appearing enabled.

**Solution:**
1. **Always name setters with `set*` prefix** — even for boolean properties. The getter can use `is*()` per JavaBean convention, but the setter must be `set*()`.
2. **When DTO getter names don't match** (e.g., OpenAPI generates `getIsToken()` for a JSON field `isToken`, but the model setter is `setToken()`), add explicit `@Mapping` annotations to the MapStruct interface:
```java
@Mapping(source = "isToken", target = "token")
@Mapping(source = "isERC20", target = "ERC20")
Currency fromDTO(TgvalidatordCurrency currency);
```
3. **After adding new boolean fields to models**, always verify the generated `*MapperImpl.java` in `target/generated-sources/` includes the mapping call. Missing calls mean the field is silently dropped.

**Affected models (fixed):** `Address.disabled`, `Wallet.disabled`, `Wallet.omnibus`, `Currency.token/ERC20/UTXOBased/accountBased/fiat/FA12/FA20/NFT`

### Java 8 Compatibility — No `var` Keyword

**Problem:** The Java SDK targets Java 8, which doesn't support the `var` keyword (introduced in Java 10).

**Solution:** Always use explicit type declarations in Java SDK code:
```java
// WRONG - Java 8 doesn't support var
for (var attr : envelope.getAttributes()) {
    // ...
}

// CORRECT - explicit type
for (Attribute attr : envelope.getAttributes()) {
    // ...
}
```

### MapStruct Version Policy

Always use stable releases. Current stable: `1.6.3`. Never use `-Beta`, `-RC`, or `-SNAPSHOT` versions in production SDKs.

### RequestMetadataAmount — String Types

**Problem:** API returns `valueFrom`, `valueTo`, `rate` as JSON strings to support arbitrary-precision amounts exceeding 64-bit limits. Using `long`/`double` causes silent data loss (e.g., `getAsLong()` returns 0 on string input).

**Solution:** Use `String` for all amount fields. Add `jsonValueToString(Object)` helper in `RequestMetadata.java` that handles both `String` and `Number` JSON inputs via `BigDecimal.toPlainString()` for lossless conversion.

### ApiRequestCursor Cannot Be Null

**Problem:** `BusinessRuleService.getBusinessRules(cursor)` and similar methods require a non-null `ApiRequestCursor`. Passing `null` causes `NullPointerException`.

**Solution:** Always construct a proper cursor:
```java
ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);
BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(cursor);
```

### ProtectClient Secret Cleanup — Validate Reflection Targets

**Problem:** `close()` uses reflection to access private `apiSecret` field in `ApiKeyTPV1Auth`. If OpenAPI internals change, this breaks silently.

**Solution:** Added `ProtectClientTest.testSecretCleanupReflectionTarget()` that validates the field exists, is `byte[]` type, and is accessible. This test catches breakage from OpenAPI regeneration.
