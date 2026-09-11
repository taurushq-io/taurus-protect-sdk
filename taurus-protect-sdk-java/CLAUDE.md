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

**Build/test without `mvn clean` (this environment):** `./build.sh` runs `mvn clean compile` first, and `clean` fails deleting the large generated `proto/target/classes/.../proto/v1` dir (owned by the dev user — not a permission issue, the delete just chokes). `build.sh` also needs `java` on **PATH**, not just `JAVA_HOME`. Iterate against the already-installed `proto`/`openapi` `1.0-SNAPSHOT` artifacts in `~/.m2` instead:
```bash
export JAVA_HOME=/workspace/java/jdk-17.0.19+10 && export PATH="$JAVA_HOME/bin:$PATH"
mvn test -o -pl client -Dspotbugs.skip=true -Dpmd.skip=true -Dcheckstyle.skip=true            # whole client suite
mvn test -o -pl client -Dtest='FooTest,BarTest' -Dspotbugs.skip=true -Dpmd.skip=true -Dcheckstyle.skip=true
mvn compile -o -pl client -D...                                                                 # main only (fast MapStruct check)
```
`-o` (offline) + `-pl client` (no `-am`) skips the proto/openapi rebuild. Piping `mvn` to `tail`/`head` buffers output until exit — redirect to a file and poll.

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

**Advanced Features**: `getAirGapService()`, `getStakingService()`, `getContractWhitelistingService()` (**writes only** — reads live on `getWhitelistedAssetService()`, the verified reader of the same endpoint), `getBusinessRuleService()`, `getReservationService()`

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
- `scripts/generate-openapi.sh` calls `patch` (BSD/GNU `patch` binary) to apply `scripts/openapi-tpv1.patch` to the regenerated `ApiClient.java`. The patch (1) collapses 4 `auth.*` imports into a wildcard, (2) replaces `new ApiKeyAuth("header","Authorization")` with `new ApiKeyTPV1Auth()` in both constructors, and (3) inserts `setApiKeyTPV1`/`setApiSecretTPV1` helper methods. If `patch` is missing, apply by hand — the diff is tiny.
- **`auth/ApiKeyTPV1Auth.java` and `auth/ApiKeyTPV1Exception.java` are committed files, not regenerated.** They are referenced by `ApiClientTPV1.java` and by the patched `ApiClient.java` (the `scripts/openapi-tpv1.patch` only renames references — it never recreates the helpers). Treat them as part of the SDK source. They were once deleted in `e4500b2 cleanup` and had to be restored from `873e982`; do not re-delete them under a "looks generated" assumption.

### Protobuf
- Uses `protoc` directly
- Generated classes flattened to proto module
- **Runtime is `protobuf-java 4.29.3` (`pom.xml:protobuf-version`).** This requires `protoc >= 21` to generate compatible Java code. Older `protoc` (≤3.20.x) emits removed APIs (`makeExtensionsImmutable()`, `Address.newLongList()`) and `mvn compile` fails on the proto module with "cannot find symbol" errors. The validatord toolchain at `/workspace/tg-validatord/scripts/tools/latest/bin/protoc` is too old; download a newer protoc release from `https://github.com/protocolbuffers/protobuf/releases` and put it earlier on PATH before regen.

## Static Analysis

Code must pass SpotBugs, PMD, and Checkstyle. The openapi module is excluded from these checks. Checkstyle config is in
`checkstyle.xml`. PMD excludes generated mappers (`*Impl.java`).

```bash
mvn checkstyle:check pmd:check spotbugs:check -pl client     # NO -o: spotbugs needs to fetch findsecbugs
```

**`spotbugs:check` fails in offline mode** — `findsecbugs-plugin` is not in the local `~/.m2` cache, and
`-o` aborts before any analysis. Maven Central is reachable; see the repo-root `CLAUDE.local.md`.

Baselines on committed master: **Checkstyle 0** (a real gate — a violation is a regression), PMD 10
(already red), SpotBugs 0. Run these on any non-trivial change: `RuleCell.java` alone carried 37 PMD
violations that nobody saw because the file was untracked.

**Common PMD/SpotBugs patterns to handle:**

- **Empty catch blocks**: a comment inside the body is **NOT** enough for this ruleset — PMD flags
  commented empty catches too. Either give the body a statement (`continue;` in a retry loop) or use
  `@SuppressWarnings("PMD.EmptyCatchBlock")` on the enclosing method
- **GuardLogStatement**: every `LOGGER.log/warning/fine` needs an `if (LOGGER.isLoggable(Level.X))` guard
- **IdenticalCatchBranches**: merge them with multi-catch (`catch (A | B e)`) — but keep a branch separate
  when it behaves differently, e.g. `NoSuchAlgorithmException` in `SignatureVerifier` must keep failing fast
- **PreserveStackTrace**: when translating an exception, chain the cause (`ex.initCause(e)`)
- **MethodLength (Checkstyle, max 150)**: `RuleCellCodec.encode` / `decodeTyped` are exhaustive switches
  over the 36 cell types and hit the limit. They are kept under it by extracting the integer families into
  `encodeIntegerCell` / `decodeIntegerCell` — extract another family rather than raising the limit, and
  re-run the golden-vector suite afterwards to prove the wire bytes did not move
- **ConstantsInInterface**: MapStruct mapper interfaces use `INSTANCE` constant - suppress with
  `@SuppressWarnings("PMD.ConstantsInInterface")`
- **Redundant null checks**: SpotBugs flags null checks on fields marked `@Nonnull` in OpenAPI models - respect the API
  contract and remove unnecessary null checks
- **Classes with only private constructors**: Must be marked `final`
- **CPD (Copy-Paste Detector)**: `minimumTokens` threshold in pom.xml controls duplication sensitivity (currently 700)

## Testing

### Unit test deps are JUnit ONLY — no Mockito, no HTTP stub library

`client/pom.xml` declares `junit-jupiter-engine` and `junit-jupiter` and nothing else for
test scope. There is no Mockito, no WireMock, no MockWebServer, and no `HttpServer`-based
stub anywhere in `client/src/test/java`. Consequences when writing a test:

- **You cannot mock a service or stub the transport.** A test that needs a service to
  return canned data has to construct the real object and drive the method that takes the
  data as an argument. Concretely: the rules-container verification gate asserts on
  `governanceRuleService.getDecodedRulesContainer(rules)` with a hand-built
  `GovernanceRules`, because `RulesContainerCache` can only be driven through a live
  `ApiClient`. Say so in the test comment rather than implying full-path coverage.
- A `GovernanceRuleService` for tests is `new GovernanceRuleService(new ApiClient(), new
  ApiExceptionMapper(), keys, minValidSignatures)` — see `RulesContainerCacheTest`'s
  `@BeforeAll`, which generates a real P-256 key via BouncyCastle.

### Building a rules container in a test

`RulesContainerMapper.INSTANCE.toBase64String(container)` /
`.fromBase64String(base64)`. `fromBase64String` throws the checked
`com.google.protobuf.InvalidProtocolBufferException`, so a test calling it needs
`throws Exception` — a compile error, but one that reads as unrelated to the test.

**Use a wire-valid container when testing that verification runs.** The model's
`getDecodedRulesContainer` verifies and *then* `parseFrom`s, so a malformed blob throws
`IntegrityException` from the parse whether verification ran or not — the test passes
either way. `RulesContainerCacheTest.wireValidContainer_decodesCleanly` exists to hold
that guarantee; keep it beside the verification tests.

### Shared cross-SDK fixtures: the loader lives in `testutil`

`testutil/SignedFixtures.load()` is the single reader of
`scripts/resources/verification-signed-fixtures.json` — public, and in `testutil` rather than
beside either consumer, because the file's two sections are exercised from **different
packages**: the SuperAdmin threshold against `helper.SignatureVerifier`
(`helper/SignedFixturesTest`), and the per-group one against a **package-private service
method** (`service/SignedFixturesGroupThresholdTest`), since whitelist verification lives in
the services here and not in `helper/`.

That split is why there are two test classes for one file. Do not duplicate the loader to
avoid it: `load()` carries the per-section count assertion that is the file's own
"was this section actually consumed?" guard, and two copies is how the two would drift on it.

### Test Configuration

Credentials are loaded from `client/src/test/resources/test.properties` (git-ignored), with environment variable overrides. Copy `test.properties.sample` to get started. The `TestConfig` class (in `testutil` package) loads identities with multi-identity support (API creds, private keys, SuperAdmin public keys).

### Integration Tests

Integration tests are located in `client/src/test/java/.../integration/` and are excluded from default test runs via surefire (`**/*IntegrationTest.java`).

**They are still COMPILED by `mvn test`** — surefire excludes them from *running*, not from
compilation — so a public-API removal that breaks one surfaces as a normal build failure here.
That is not true of the sibling SDKs (Go's `go build` skips `_test.go` entirely, TS's `tsc`
excludes `tests/`), so a cross-SDK deletion needs the per-language sweep in the repo-root
`CLAUDE.md` → "Deleting from a public surface".

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

### Gson cannot deserialize the domain models under JDK 17 — use reflection in a test

`new Gson().fromJson(json, Request.class)` fails with `JsonIO Failed making field
'java.time.OffsetDateTime#dateTime' accessible`. The surefire add-opens profile only opens
`java.base/java.lang` (for `Throwable.detailMessage`), and `java.time` internals stay
closed — so any model carrying an `OffsetDateTime` is un-deserializable by Gson here, and
most of them do.

This matters when a test needs to forge a private field, e.g. proving `approveRequests`
re-verifies rather than trusting `RequestMetadata.hashVerified` (private, no setter). Gson
would be the realistic attack path but cannot build the object, so set the field directly:

```java
Field f = RequestMetadata.class.getDeclaredField("hashVerified");
f.setAccessible(true);
f.setBoolean(md, true);
```

`setAccessible` works because the class is the SDK's own, not `java.base`. Catch
`ReflectiveOperationException` and `fail(...)` with a message naming the field, so a rename
reports as a stale test rather than a mysterious skip.

### A DTO fixture needs a status label or the mapper rejects it

`RequestMapper.INSTANCE.fromDTO(dto)` throws `IllegalArgumentException: Request status label
must not be null or empty`. So a `TgvalidatordRequest` fixture built for a metadata test
still needs `dto.setStatus("APPROVING")` — the failure names status, not the thing you were
testing, which reads as an unrelated break.

### JDK 9+ surefire add-opens (JDK-conditional)

`ApiExceptionMapperTest` deserializes JSON into `com.taurushq.sdk.protect.client.model.ApiException` via Gson. `ApiException extends Exception`, so Gson reflects into `java.lang.Throwable.detailMessage` — fine on JDK 8, but JDK 9+ strong encapsulation refuses it without an explicit add-opens flag.

**The argLine is gated by a JDK profile**, not added to the main surefire config, because the JDK 8 launcher rejects the flag (`Unrecognized option: --add-opens` → "Could not create the Java Virtual Machine") and would crash every test fork on systems where `java` is JDK 8. The build script's `MIN_JAVA_VERSION=8` plus the user's listed Corretto-8 setup mean JDK 8 is a supported target.

Lives in `client/pom.xml` as a profile activated by `<jdk>[9,)</jdk>`:

```xml
<profiles>
  <profile>
    <id>jdk9-plus-add-opens</id>
    <activation>
      <jdk>[9,)</jdk>
    </activation>
    <build>
      <plugins>
        <plugin>
          <artifactId>maven-surefire-plugin</artifactId>
          <configuration>
            <argLine>--add-opens java.base/java.lang=ALL-UNNAMED</argLine>
          </configuration>
        </plugin>
      </plugins>
    </build>
  </profile>
</profiles>
```

Verify activation with `mvn help:active-profiles -pl client` → expect `jdk9-plus-add-opens` listed on JDK 9+ and absent on JDK 8.

Symptoms:
- Argline not gated, run on JDK 8 → `Unrecognized option: --add-opens` and `The forked VM terminated without properly saying goodbye`.
- Argline missing on JDK 17+ → `JsonIO Failed making field 'java.lang.Throwable#detailMessage' accessible` on 5 `ApiExceptionMapperTest` cases.

**XML comment gotcha:** the literal string `--` is forbidden inside an XML comment, so a `pom.xml` comment like `<!-- explains --add-opens -->` is a parse error. Maven reports `Non-parseable POM ... in comment after two dashes (--) next character must be > not a`. Spell the flag out in prose ("the add-opens flag") rather than embedding the literal CLI form.

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

## Verification surface added in the 2026-09-04 pass

- **`UnverifiedMetadataException extends RequestMetadataException`.** `RequestMetadata` no
  longer parses in `setPayloadAsString` — that ran at MapStruct mapping time, before anything
  could reject the payload, and threw `NullPointerException` on a null one. Parsing is now lazy
  inside `verifiedPayload()`, which throws unless `hashVerified` is set. Every extraction method
  routes through it.
- **`SignatureVerifier.containsHash`** is the per-signature half of the hash-comparison pair.
  Both whitelist services carried their own `List.contains` copies — `String.equals`, not
  constant-time, early-returning — while `verifyHashCoverage` sat in `helper/` with **zero
  callers**. All four sites now route through the helper.
- **`WhitelistHashHelper.resolveRuleKey`** returns the `(blockchain, network)` pair from the
  signed payload as a two-element array.
- **`approveRequests`/`approveRequest` gained 3-argument overloads** taking a comment. Overloads
  rather than a changed signature, matching how `TransactionService` handles its extra filters.
- **The approve sort works on a copy.** It sorted the caller's list in place, which mutated
  their argument and threw `UnsupportedOperationException` on an immutable one.
- **Whitelist verification still lives in the services, not `helper/`.** That is exactly where
  the non-constant-time comparisons hid. Only the crypto helpers were relocated; the full
  extraction is written up in `TODOS.md`, and the alignment report records the layering as an
  accepted difference *with that TODO as the plan to stop accepting it*.

## Verification surface added in the 2026-09-10 security-scan pass

Go was the reference SDK for this pass; cross-SDK reasoning is in the repo-root `CLAUDE.md`.
Java-specific:

- **CHECKED vs UNCHECKED is a CLASSIFICATION rule here, and getting it wrong re-opens a fixed
  finding.** Throw the checked `WhitelistException` for anything that is one ROW's problem;
  reserve the unchecked `IntegrityException` / `ContainerIntegrityException` for what
  invalidates the whole call. `WhitelistHashHelper.rejectDuplicateObjectKeys` first shipped
  throwing `IntegrityException`, which silently escapes every `catch (WhitelistException)` — so
  one unparseable row would have aborted a whole listing, exactly the failure
  `WhitelistedAddressListResult`'s javadoc claims was fixed. Both public parse entry points
  declare `throws WhitelistException`, so two pre-existing tests caught it
  (`WhitelistHashHelperTest.testParseWhitelistedAddressFromJson_InvalidJson`,
  `VerificationBehaviourVectorsTest`). Keep that alignment when adding a check to either parser.
- **`RulesContainerCache` released its single-flight flag only on the `catch (ApiException)`
  path.** `doFetch()` runs governance verification, whose `IntegrityException` is unchecked, and
  the Gson error path can raise `StackOverflowError`; either escaped before `fetching` was
  cleared, so every later caller parked on the untimed `lock.wait()` forever — **one crafted
  `/rules` response wedged address, asset and price verification process-wide.** Both
  `getDecodedRulesContainer` and `invalidate` now release through a `finally` +
  `releaseTheFetchSlot()`. `fetchException` is deliberately NOT set there: it is typed
  `ApiException`, and the escaping throwable reaches its own caller anyway. Gated by
  `RulesContainerCacheTest.{getDecodedRulesContainer,invalidate}_failedFetchReleasesTheSlot`,
  which are 5-second timeouts on another thread — they were the tests that proved the fix was
  still missing after everything else had landed.
- **The GENERATED `openapi.ApiException.getMessage()` interpolates the ENTIRE response body**
  (`"…HTTP response code: %s%nHTTP response body: %s…"`, and it is an override, so `super.getMessage()`
  is only the first `%s`). That defeated `ApiExceptionMapper`'s own `MAX_ERROR_BODY_BYTES` ceiling in
  the one case the ceiling exists for: an over-ceiling body was correctly not PARSED, and then the
  fallback copied `e.getMessage()` into the SDK exception anyway — so a hostile body still reached an
  integrator's log in full. `truncateForMessage` / `MAX_MESSAGE_CHARS` (2048) bound it, with a
  `[truncated]` marker so a cut message is distinguishable from a short one. A body *excerpt* is
  diagnostically useful, so the fix is a bound, not removal. Found by the new
  `ApiExceptionMapperHostileBodyTest`; the red run reports `131248 chars for a body of 65566`.
- **MapStruct silently produced an ALL-NULL bean from an enum source.**
  `MultiFactorSignatureMapper.fromEntityTypeDTO` maps a bare enum to a bean with `id`/`kind`;
  with no explicit mapping, MapStruct's default bean mapping sets **no** target properties, so
  `getKind()` was null and the mapper test's `assertNotNull` passed. It is a hand-written
  `default` method now. `id` stays null on purpose — the reply carries no entity id, and that
  absence is the MFA blocker itself. This is the same class as the `is*()`-setter trap already
  documented below: MapStruct's failures here are silent, so assert VALUES, never non-null.
- **`markVerified` takes the VERIFIED PAYLOAD as a parameter**, on both envelopes, and refuses
  an absent one rather than falling back to `metadata.getPayloadAsString()`. The delivered text
  is not always what a signature covered (the legacy strips are not injective, and Gson keeps
  the LAST of two duplicate keys), so parsing the delivered payload is precisely the injection.
  `metadata.payloadAsString` is left untouched — a caller needs it to reproduce `metadata.hash`.
- **The approval pin types have no public constructor.** `WhitelistedAddressApproval` /
  `WhitelistedAssetApproval` are minted only by `WhitelistedAddressListResult.select(ids)` /
  `selectAll()` and `WhitelistedAssetResult.select(ids)` / `selectAll()`, so a test in the
  `service` package must build a result first — which is the point of the type, not friction.
- **Java's address approve now re-reads through the NORMALIZED list path**, where containers are
  response-level and label-verified, rather than the per-row in-band containers it used before.
  That converged it onto Go/Python/TypeScript; the asset side still uses in-band, as all four do.
- **Three test files re-implemented the legacy-hash regexes inline** because the production
  methods were private. All three now call `WhitelistHashHelper.computeLegacyPayloadVariants` /
  `AssetHashHelper.computeAssetLegacyHashes`: `crypto/CrossSdkCryptoVectorTest`,
  `service/WhitelistedAddressServiceLegacyHashTest`, `helper/WhitelistVerificationFlowTest`.
  A copy here is worse than no test — the cross-SDK oracle would keep asserting the pre-fix
  semantics while the SDK moved.
- **PMD is now 0, down from 10 on master.** Almost all of it was unused imports and unnecessary
  fully-qualified names left behind when logic moved, plus two `PreserveStackTrace` in
  `RequestService` (chain the cause: `new IntegrityException(msg, e)`). Watch out for a
  find-and-replace that strips a package qualifier from an **import** line — `java.util.Map<` →
  `Map<` is safe, but `java.security.NoSuchAlgorithmException` → `NoSuchAlgorithmException`
  rewrote the import statement itself into `import NoSuchAlgorithmException;`.

## Verification surface added in the 2026-09-07 pass

Cross-SDK rules are in the repo-root `CLAUDE.md`. Java-specific:

- **`RequestMetadata.setHashVerified` is GONE.** It was `public`, and the payload gate read
  only that flag, so any caller could unlock real payload data from metadata nothing verified.
  `public boolean verifyAndMaterialise()` replaces it: it performs the hash check and sets the
  flag in one operation, so the two cannot be separated. A package-private setter was not an
  option — `RequestMetadata` (`…client.model`) and `RequestService` (`…client.service`) are in
  different packages. `RequestService.verifyMetadataHash` is now a one-line delegate.
- **`helper/PriceVerifier.java`** + `model/PriceSignature.java` + `signatures` on `Price`. Note
  `PriceMapper` carries `@Mapping(target = "signatures", ignore = true)`: the generated
  `TgvalidatordCurrencyPrice` has no such field even though `apis.swagger.json` declares it —
  a stale-snapshot symptom, so populating it needs codegen, not a mapper change.
- **`SignatureVerifier.keyFingerprint` is `public static`**, and both services'
  `verifyGroupThreshold` is package-private so its tests can reach it.
- **`model/WhitelistedAssetResult`** exists because the asset list had no page total while the
  contract list did. Its `hasMore(currentOffset, pageSize)` is overflow-safe
  (`totalItems > currentOffset && totalItems - currentOffset > pageSize`), deliberately unlike
  the deleted `WhitelistedContractAddressResult`'s `(currentOffset + pageSize) < totalItems`.
- **Ordering matters in `approveRequests`.** Java checks metadata *before* `privateKey`, the
  reverse of Go, so the new hash-verified refusal sits **after** `checkNotNull(privateKey)` —
  otherwise `approveRequests_throwsOnNullPrivateKey` starts failing on the wrong error.
- **Test-fixture trap:** `requestWith` builds RAW mapped rows, some deliberately tampered.
  Blanket-verifying them at construction breaks the test whose subject is the service dropping
  them. Leave that fixture unverified.
- Service happy-paths are not unit-testable here (JUnit only, no Mockito or HTTP stub), so this
  pass's Java coverage is on the models and on argument validation — say so in the test rather
  than implying full-path coverage.

## Lessons Learned (Non-Security)

### AuthorizationException derives requiredRoles in its constructors

`AuthorizationException(message, ...)` calls the static `parseRequiredRoles(message)` itself, so
`ApiExceptionMapper.createTypedException` needed no change — it already passes `parsed.getMessage()`. Keep the
derivation in the exception: putting it in the mapper would leave directly-constructed exceptions with empty
roles. `getRequiredRoles()` returns an unmodifiable list. See the cross-SDK contract in the repo-root CLAUDE.md.

### Client Authentication (Credentials)

`Credentials` (`client/.../client/Credentials.java`) is an abstract sum-type with static factories `apiKey(k,s)` / `bearerToken(t)` / `bearerTokenProvider(sup)` and a package-private `applyTo(ApiClient)`. `ProtectClient.create(host, Credentials, keys, minSig[, ttl])` and `builder().credentials(Credentials)` are the clean path; the flat `create(host, apiKey, apiSecret, …)` and builder `credentials(k,s)`/`apiKey`/`apiSecret` are `@Deprecated` (delegate to `Credentials.apiKey`). `createFromPem` is NOT deprecated (a PEM-decoding convenience). SuperAdmin keys are mandatory for every mechanism — enforced because the Governance/WhitelistedAddress/WhitelistedAsset service constructors already `checkArgument(!superAdminPublicKeys.isEmpty())` (no service change was needed for the keys-mandatory decision). `@Deprecated` is safe: `pom.xml` sets `failOnWarnings=false` + `failOnError=false`, so deprecation warnings across the test suite don't fail the build.

### Building a whitelist signature entry in a test — three Java-only shapes

`hashes` is a `private final List<String>` initialised inline, exposed only through `getHashes()`.
Populate it with `getHashes().add(...)` / `addAll(...)`; `setHashes(...)` does not exist and reaching
for it is a compile error, not a silent no-op.

The other two bite when porting a fixture from another SDK, because **this SDK is the outlier
on both** (the four-way comparison is in the repo-root `CLAUDE.md` → "Signed fixtures"; keep
the two in step):

- **`WhitelistUserSignature.signature` is `byte[]`, not a base64 `String`.** Go, Python and
  TypeScript all hold base64, so a shared JSON fixture must be `Base64.getDecoder().decode(...)`d
  here and consumed verbatim there.
- **The nested setter is `WhitelistSignature.setSignature(WhitelistUserSignature)`** — not
  `setUserSignature`. The getter is `getSignature()` and returns the nested object, so the name
  reads like the raw signature and is not.

`RuleUser` accepts either form: `setPublicKeyPem(String)` or
`setPublicKey(CryptoTPV1.decodePublicKey(pem))`. The verifiers read `getPublicKey()`, so a
PEM-only fixture user silently contributes nothing to a threshold — set the decoded key.

### Cross-SDK surface added here (keep it)

- `SignatureVerifier.verifyHashCoverage(hash, signatures)` — Java was the only SDK without it, so
  callers were comparing hashes by hand, which is exactly where a non constant-time compare creeps
  in. The loop deliberately does **not** break on a match: returning early leaks which signature
  matched through timing.
- `BusinessRuleService.updateTransactionsEnabled(boolean)` — the transactions kill switch, previously
  Go-only even though the generated op exists in all four.
- `TransactionService.getTransactions(...)` / `exportTransactions(...)` gained 8-argument overloads
  carrying `blockchain` + `network`; the 6-argument forms delegate with nulls, so no caller breaks.
  **8 params is exactly `ParameterNumber max` in `checkstyle.xml`** — a ninth filter needs an options
  object, not another parameter.
- `GovernanceRuleService.verifyGovernanceRules(rules)` — single-argument overload using the configured
  threshold, which is the cross-SDK shape.

### JUnit fails fast — one `@Test` per invariant

`SignatureVerifierTest` packed all five `minValidSignatures` distinct-key cases into a single `@Test`
with sequential asserts. JUnit stops at the first failure, so a break in case 2 meant cases 3-5 never
ran and the report showed one failure instead of four. They are five separate tests now (Go uses
`t.Run` subtests and Python/TS separate tests for the same reason). Shared signatures live in a small
private `Fixture` class rather than being recomputed per test.

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

**Java 9+ *APIs* are the sharper edge, and `source`/`target` do NOT catch them.** `var` is a
syntax error under `-source 8`, so it fails everywhere. A Java 9+ *library method* does not:
`source`/`target` only pick the language level and bytecode version, while javac still links
against the **running** JDK's class library. So on a JDK 17 dev machine `List.of(...)` (Java 9)
and `"a".repeat(n)` (Java 11) compile clean, every local gate stays green, and the build breaks
only on a real JDK 8 — which is what shipped in `GovernanceRuleServiceTest` and
`AuthorizationErrorVectorsTest`.

Java 8 equivalents: `Arrays.asList(...)` for `List.of` (the house style — 158 uses in
`client/src/test`), `String.join("", Collections.nCopies(n, s))` for `String.repeat`. Also absent
in 8: `Map.of`/`Set.of`/`Map.entry`, `String.isBlank`/`strip`/`lines`, `Optional.isEmpty`,
`Stream.toList`, `Files.readString`/`writeString`, `InputStream.readAllBytes`.

**Now enforced locally** by a `maven.compiler.release` property in the `jdk9-plus-release`
profile in `client/pom.xml`, which pins the class library to 8 as well and turns the above into
a compile error on a JDK 9+ machine. Verified by reintroducing `"a".repeat(64)` and getting the
same `cannot find symbol: method repeat(int)` on JDK 17 that a JDK 8 build reports.

Two things not to change about it:

- **It is JDK-conditional, and must stay so** — the JDK 8 javac has no such flag. Same
  `<jdk>[9,)</jdk>` pattern as `jdk9-plus-add-opens` beside it, so on a JDK 8 machine both are
  inactive and nothing changes. `maven.compiler.release` stays **unset** in the parent
  `pom.xml` properties for that reason.
- **It is scoped to the `client` module on purpose.** Setting it in the parent, where it also
  reaches the generated `openapi`/`proto` modules, makes MapStruct fail in a *full reactor*
  build with `No implementation was created for WhitelistedAddressMapper ... erroneous element
  ...Whitelist.WhitelistedAddressOrBuilder` — while `mvn -pl client` against those modules'
  installed jars passes, so the breakage only appears in `./build.sh`. `client` is the
  hand-written module, so that is also where the rule is worth enforcing.

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

### Governance Rules Typed Mapper (`RulesContainerMapper` + `RuleCellCodec`)

Implements the cross-SDK typed governance API (see repo-root `CLAUDE.md` → "Governance Rules Typed API"). Java-specific choices:

- **Cells stay `List<ByteString>` on `RuleLine`**; the typed `RuleCell` union is decoded/encoded via a standalone `RuleCellCodec`, NOT auto-decoded into the model. MapStruct maps each `Line` independently and can't see the sibling `columns` needed to type a cell, so a typed `List<RuleCell>` field is impractical. Users call `RuleCellCodec.decode(columnType, bytes)`.
- `RuleCell` is one file: an abstract base + nested `public static final` subclasses (Java 8 — no sealed/records), with `equals`/`hashCode` via commons-lang3 `EqualsBuilder`/`HashCodeBuilder` reflection.
- Lossless plumbing: all 15 container nodes extend `RulesNode` (unknownFields) / `RulesNodeWithProperties` (+`properties` map). MapStruct populates them via `@Mapping(target="unknownFields", expression="java(unknownBytes(proto))")` on every `fromProto` — this MUST be explicit because `unknownFields` name-clashes with protobuf's own `Message.getUnknownFields()` (returns `UnknownFieldSet`, not `ByteString`) and auto-mapping is a hard compile error. Encode is hand-written `toProto`/`toBytes`/`toBase64String` default methods that re-attach unknown fields via `builder.setUnknownFields(UnknownFieldSet.parseFrom(bytes))` and strip `enforcedRulesHash`/`timestamp`.
- `RuleSource` carries a `raw` `ByteString` fallback: whitelisting source cells that don't round-trip byte-identically keep their verbatim bytes.
- **`CosmosDetails` is a node class, not a flattened list.** It was `List<String> cosmosMethodSignatures`
  on `TransactionRuleDetails`, which made unknown-field preservation structurally impossible for that
  sub-message. It is now `CosmosDetails extends RulesNode` with `getMethodSignatures()`, matching the
  other three SDKs. MapStruct needs an **explicit** `fromProto(...CosmosDetails)` with
  `@Mapping(source = "methodSignaturesList", target = "methodSignatures")` plus the usual
  `unknownFields` expression — without its own method, auto-mapping hits the `UnknownFieldSet` vs
  `ByteString` name clash and fails to compile.
- **Encode is deterministic**: `toBytes` goes through `deterministicBytes(Message)`, which writes via a
  `CodedOutputStream` with `useDeterministicSerialization()`. Plain `toByteArray()` emits
  `map<string, bytes> properties` in unspecified order, so the same container would encode to different
  bytes across runs and differently from the other SDKs. Don't reintroduce `toByteArray()` on the
  container path — **and that includes the JSON bridge**: `RulesContainerJsonMapper` used plain
  `toByteArray()` in both `rulesContainerBase64FromJson` and `ruleMessageBase64FromJson`, which made
  the bytes a proposal is signed over order-dependent. It now delegates to
  `RulesContainerMapper.INSTANCE.deterministicBytes(...)` via a private static helper — `deterministicBytes`
  is an interface `default` method, so it is only reachable through `INSTANCE`, not statically. Keep the
  single implementation; do not copy the serializer.
- **Nested rule-detail unknown fields are re-attached on encode.** `toProtoDetails` builds each nested
  builder and calls `reattach(builder, node.getUnknownFields())` before `setX(...)`. Decode captured them
  all along; the encode side used to drop them.
- **`AddressWhitelistingRules.includeNetworkInPayload` is a model-only field with NO proto backing** — there is no `setIncludeNetworkInPayload` on the proto builder; don't try to encode it.

**Governance test-authoring notes:**
- `RuleCellCodec.cellFamily(RawCell)` returns the cell's `columnType` (not `""`); `decode(unknownColumn, emptyBytes)` returns a `RawCell` (not null).
- Container-level enums are **numeric-passthrough, matching the other three SDKs**: decode reads the `*Value` int accessors (`typeValue`, `domainValue`, `subDomainValue`, `blockchainValue`, `rolesValueList`) and renders a value this SDK does not know as its decimal string — `UNRECOGNIZED` does not carry the number, so sourcing the enum itself would lose it. `enumValue(EnumClass, name)` resolves a known name, passes a decimal string through, returns `0` for null/empty, and **throws `IllegalArgumentException`** on a name that is neither; `rolesFromStrings` passes unknown role numbers through rather than dropping the role. Port Go's enum passthrough and authored-error tests — they apply here now.
- The cell-level `blockchainToInt` is **private** and throws `NumberFormatException` on a non-numeric, non-enum name; it is numeric-passthrough like the container-level enums.
- `RuleCellCodec.magnitude` / `fromMagnitude` are **package-private statics** — call them directly from a same-package `RuleCellCodecTest`.
- `RuleCell` has reflection-based `equals`, so `assertEquals(expectedCell, decoded)` works; `RuleSource` has **no** `equals` — compare field-by-field (getType + variant getter) or by re-encoded bytes.
- Inject per-node unknown fields with `RequestReply.<Node>.newBuilder()...setUnknownFields(UnknownFieldSet.newBuilder().addField(500, Field.newBuilder().addVarint(42).build()).build())`, then assert each decoded model node's `hasUnknownFields()`.
- Service tests stay validation-only (project forbids Mockito — no network stubs); the encode/decode paths are covered via the mapper tests.
