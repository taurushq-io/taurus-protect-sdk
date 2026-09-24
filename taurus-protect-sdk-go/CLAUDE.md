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

### Available Services (39 + TaurusNetwork)

**Core**: `Wallets()`, `Addresses()`, `Requests()`, `Transactions()`, `GovernanceRules()`, `Balances()`, `Currencies()`, `WhitelistedAddresses()`, `WhitelistedAssets()`

**Transaction/Request**: `Audits()`, `Changes()`, `Fees()`, `Prices()`

**Advanced**: `AirGap()`, `Staking()`, `WhitelistedContracts()` (**writes only** — reads live on `WhitelistedAssets()`, the verified reader of the same endpoint), `BusinessRules()`, `Reservations()`, `MultiFactorSignature()`

**Administrative**: `Users()`, `Groups()`, `VisibilityGroups()`, `Config()`, `Webhooks()`, `WebhookCalls()`, `Tags()`

**Specialized**: `Assets()`, `Actions()`, `Blockchains()`, `Earn()`, `Exchanges()`, `Fiat()`, `FeePayers()`, `Health()`, `Jobs()`, `Scores()`, `Statistics()`, `TokenMetadata()`, `UserDevices()`

**Taurus Network** (namespace): `client.TaurusNetwork().Participants()`, `.Pledges()`, `.Lending()`, `.Settlements()`, `.Sharing()`

SDK alignment: see `docs/SDK_ALIGNMENT_REPORT.md` (repository root). Java SDK is source of truth.

## Code Generation

### OpenAPI Generator
Uses `openapi-generator-cli` JAR with `-g go`:
- `enumClassPrefix=true` -- avoids enum constant redeclaration conflicts
- `disallowAdditionalPropertiesIfNotPresent=false` -- unknown fields land in each model's `AdditionalProperties` instead of failing the reply; with `templates/go/model_enum.mustache`, unknown enum values keep their raw string. Same contract in all four SDKs, gated by `scripts/resources/decode-tolerance-vectors.json`
  - A decoded DTO always carries a non-nil `AdditionalProperties` (empty when nothing was unknown): compare
    domain models, not generated DTOs, in tests. Enum `NewXFromValue` never errors; `IsValid()` is the
    known-value check. Hand-written code casts caller values (`openapi.X(v)`) and sends them verbatim.
  - `generate-openapi.sh` aborts if any model still contains `DisallowUnknownFields()` or `is not a valid`,
    and `patch_metadata_payload` requires exactly one `delete(additionalProperties, "payload")` — that
    line keeps the raw payload out of `AdditionalProperties` now that the field is removed.
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
4. Page through the helpers in `service/pagination.go` — see "Pagination" below. Never parse a
   reply count or cursor inline.

### The service layer is where filters get LOST — check the generated client first

`List*Options` is hand-written and routinely exposes a fraction of what the endpoint takes. Measured
2026-09-03, before a widening pass: addresses **6 options against 33 client parameters**, whitelisted
addresses 10 against 30, transactions 8 against 21, wallets 5 against 12. A filter that exists but is
unreachable forces every consumer to over-fetch and narrow in memory.

So **do not conclude "the API cannot do that" from the options struct**. Grep the generated builder methods:

```bash
grep -n "func (r ApiWalletServiceGetAddressesRequest) " internal/openapi/api_wallets.go
```

That grep is also the only reliable check for whether a parameter exists at all — see the next section.

Some intent is already recorded in the source: `GetTransaction` is
`TransactionServiceGetTransactions(ctx).Ids([]string{txID}).Limit("1")` with the comment *"Use the list
endpoint with IDs filter"*, and `GetUsersByEmail` is `UserServiceGetUsers(ctx).Emails(emails)`. The single-get
and list methods are frequently ONE endpoint, so a consumer collapsing them loses nothing — but confirm per
pair, because `GetWallet`/`GetUser`/`GetRequest` really are distinct RPCs.

Two shapes to preserve when widening:

- **A tri-state filter needs a pointer.** `ListUsersOptions.TotpEnabled` is `*bool`: as a plain `bool` an
  explicit `false` ("users WITHOUT TOTP") is indistinguishable from "no filter".
- **A nested filter group needs a pointer too.** `ListAddressesOptions.Score *AddressScoreFilter` — nil omits
  every score parameter. A zero-valued nested struct would send empty provider params on every address call.

### The shared swagger is a STALE SNAPSHOT — see the repo-root `CLAUDE.md`

`scripts/resources/swagger/apis.swagger.json` sits at the SDK **repo root**, feeds all four generators, and
has drifted from `tg-validatord/api/swagger/`. The root `CLAUDE.md` carries the rule and the four-SDK regen
procedure; what it means *here* is that a `List*Options` gap may be un-fillable without codegen. Confirmed
absent from the generated Go client 2026-09-03:

| Missing | Live in validatord | Consequence |
|---|---|---|
| `GetPriceByID` | `api/swagger/v1/price-service.swagger.json` | no single-price read (`PriceService_QueryPricesV2` was patched in 2026-09-24 and backs `ListPrices`) |
| `Entities` (plural), `FieldKey`, `FieldValue` on `ChangeServiceGetChanges` | endpoint has 11 fields | `ListChangesOptions` can only reach 8 |

**Always grep `internal/openapi/` before writing a builder call.** A plausible-looking `.Actions(...)` on the
whitelist request does not exist; the failure is a compile error, but only after you have written the option
field, the mapping and the consumer.

Adopting one of these is a **codegen job, not plumbing**: patch just the new path plus its schemas into
`apis.swagger.json` from validatord's per-service spec, then re-run `generate-openapi.sh`. Do NOT do a
wholesale refresh — that pulls in every validatord change since the snapshot, an uncontrolled diff across all
services. And the spec is shared by four SDKs, so either regenerate all four or state plainly that the spec
now leads them.

### Deprecated endpoints and fields the SDK wires (or used to)

The generated client exposes both generations of a deprecated field, so "wire up the score filters" naively
adds six dead parameters. Known cases:

- **Endpoints.** `ListWallets` called `WalletServiceGetWalletsInfo` (deprecated, with a `//nolint:staticcheck`
  admitting it) → now `GetWalletsV2`, a drop-in: same request message, same reply. The v1 `Fees().GetFees`
  (deprecated `FeeServiceGetFees`) and the v1 unpaged `Prices().GetPrices` are removed — `GetFeesV2` and
  `ListPrices` (`QueryPricesV2`) replace them. `DeleteWhitelistedContract` was deprecated with **no
  replacement** and could no longer delete anything, so it is not wrapped.
- **Fields.** On `GetAddresses`, the flat score params (`scoreProvider`, `scoreInBelow`, `scoreOutBelow`,
  `scoreExclusive`, `coinfirmScoreGreater`, `chainalysisScoreGreater`) are all `deprecated = true` in
  `wallet-service.proto` — use the nested `scoreFilter` generation. `ListWhitelistedAddressesOptions.Currency`
  mapped to the deprecated `currency`; `exchangeAccountId`, `contractType` and `isNFT` are likewise superseded
  by their plural/nested forms.

Deprecation is only visible in the proto, not in the generated Go — check
`tg-validatord/api/proto/v1/` when adding a call.

### Helper Functions (Mapper Utilities)
All safe pointer helpers in `pkg/protect/mapper/helpers.go`: `safeString`, `safeBool`, `safeInt64`, `safeFloat32`, `safeTime`, `stringPtr`, `boolPtr`, `int64Ptr`.

**Important:**
- **service package**: `stringPtr()` defined separately (packages can't share unexported functions)
- In tests, define `testStringPtr()` to avoid conflicts with production helpers

## Pagination — `service/pagination.go` is the only place it happens

The cross-SDK contract is in the repo-root `CLAUDE.md` ("Pagination (cross-SDK)"). In this SDK:

```
options ─▶ resolveOffsetWindow / resolveCursorWindow / resolveSize ─▶ apply*(req, window) ─▶ request
reply ───▶ offsetPagination(rule, window, served, excluded, offsetReply) ─▶ *model.Pagination (never nil)
       └─▶ cursorPage(pageSize, cursorReply{Cursor | TokenOnly+Token, HasTotal+Total, Excluded}) ─▶ model.CursorPage
```

- **Every list always sends a page size.** 0 → `model.DefaultPageSize` (20); above
  `model.MaxPageSize` (100) or negative → an error naming the option, before any request. Limit-only
  endpoints use `resolveSize` with their own maximum: price history 365, the prices-history export
  and the transaction export none (`max 0`).
- **One offset rule per endpoint** (`offsetRule`): wallets/addresses read the reply `offset` —
  validatord returns the NEXT page's offset there, which the old code read as the current one and
  so dropped the last page; users/groups `offset + min(rows, limit)` (a synthetic row can be
  appended); whitelisted addresses `offset + rows the server returned` (SDK exclusions reduce only
  `TotalItems`); whitelisted contracts `offset + limit`; everything else `offset + rows`.
  `HasMore = next > offset && next < serverTotal`.
- **Cursor lists** take `PageSize` + `Cursor`; `Cursor` is sent as `currentPage` with
  `pageRequest=NEXT` and cannot be combined with the low-level `CurrentPage`/`PageRequest`.
  `NextCursor` is set only when the reply has `hasNext` — never derive a cursor any other way.
  The body-cursor endpoints (`QueryPricesV2`, the v2 asset queries, asset addresses/wallets) take
  `window.body()`.
- **The transaction export ignores `offset` server-side** (validatord always exports from row 0),
  so `ExportTransactions` is limit-only and reports `TotalItems`.
- **Counts** are canonical decimals in `[0, 2^53-1]` or an error wrapping
  `model.ErrMalformedPagination` (`parseCount`): a count silently read as 0 ends a walk early, and
  2^53-1 keeps TypeScript reading the same number.
- **Internal id-filtered reads batch by `chunkIDs(ids, size)`**: ≤ 100 ids for the whitelist
  re-reads and `GetUsersByEmail`, ≤ 50 (`maxAddressIDsPerRequest`, validatord's cap) for managed
  addresses.

Tests, all through the real generated client against `httptest`:

| File | Pins |
|---|---|
| `pagination_test.go` | the helpers: every rule, `{}`, overflow, malformed counts, windows, chunking |
| `pagination_vectors_test.go` | `scripts/resources/pagination-vectors.json`, plus each wrapped operation's rule probed through its real method |
| `list_request_vectors_test.go` | `scripts/resources/list-request-vectors.json` — the explicit `listAdapters` table maps operationId → Go method and canonical option → field; `notWrappedOperations` lists the three paged ops Go does not wrap, with reasons |
| `pagination_walk_test.go` | walks per family visit every row once and stop (reply offset, synthetic row, trailing empty page, v2 cursor with `+ / =` on the wire, body cursor, token, requestCursor) |
| `whitelisted_address_pagination_test.go` | a signed whitelist fixture, and exclusions never shifting the walk |
| `list_endpoints_test.go` | endpoint policy: QueryPricesV2 verification, the approval-queue rejections, chunked re-reads, the new v2 asset / Earn / fiat-entity lists |

Adding a list method: wire it through the helpers, add an adapter row (the loader fails on an
unmapped operation or option), and bump the loaders' count constants when the coordinator
regenerates the vector files.

Fixture traps that cost time:

- A `walkClient`/`httptest` handler must set `Content-Type: application/json`, or the generated
  client fails every call with "undefined response type".
- `mapper.RulesContainerToBase64` encodes `RuleUser.PublicKeyPEM`, not the parsed `PublicKey`, and
  rejects role names the proto does not know (`USER` is not one; `WHITELISTEDADDRESSAPPROVER` is).
  A container built with only `PublicKey` verifies at step 2 and then fails step 5 with "user has
  no public key".
- A test that panics stops the whole test binary, so a mutation run must re-run the tests that
  never reported.

## `QueryAssetAddresses` rows are confirmed through the verified readers

`AssetServiceV2_QueryAssetAddressesV2` rows carry no signature (repo-root `CLAUDE.md`, Cross-SDK
Security Rules). `AssetService` therefore takes the client's `AddressService` and
`WhitelistedAddressService` (`NewAssetService` panics on nil; `Client.Assets()` resolves both
before taking the client lock, which is not reentrant) and completes each page in
`verifiedAssetAddresses`:

- INTERNAL → `AddressService.verifiedAddressesByID`: GetAddresses + addressIds, every row through
  the existing `verifiedAddress` seam. An `IntegrityError` on a row is an exclusion; anything
  else, or a container without an HSMSLOT key, aborts — `helper.VerifyAddressSignature` reports a
  missing HSM key as a per-row `IntegrityError`, so the container is checked once up front.
- WHITELISTED → `WhitelistedAddressService.verifiedAddressesByID(ids, includeForApproval)`, the
  seam the approval re-read also uses (with `true`; the holders read passes `false`).
- Anything else is `Verified = false`. The generated enum rejects an unknown `addressType` at
  decode, so only EXTERNAL and untyped rows reach this branch in Go.

Tests: `asset_holders_verification_test.go` (HSM-signed managed addresses and 6-step-signed
whitelist rows through one `httptest` fixture).

## Testing

### Running Tests
```bash
./build.sh unit                      # Preferred
go test ./pkg/protect/...            # Direct
go test -v -cover ./pkg/protect/...  # Verbose + coverage
```

**Coverage caveat:** per-package `go test -cover` counts only statements executed by *that package's own* tests. Governance code is exercised cross-package (e.g. `mapper` round-trip tests call `model.HasUnknownFields`; `service` tests call `mapper.RulesContainerToBytes`), so per-package numbers understate real coverage — a function can read 0.0% while being fully exercised by another package's tests. Measure combined coverage with `-coverpkg`:
```bash
go test -coverpkg=./pkg/protect/... -coverprofile=cov.out ./pkg/protect/...
go tool cover -func=cov.out | grep -E 'rule_cell|rules_container'
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

### Regen scripts assume macOS / a JDK install

- `scripts/generate-openapi.sh` checks Java via `javap`, which lives in the JDK only — a JRE-only install (e.g., the bundled sonar-scanner JRE) makes the script fail before it gets to the JAR. Workaround: invoke `java -jar scripts/resources/jars/openapi-generator-cli-7.9.0.jar generate -g go -i ../scripts/resources/swagger/apis.swagger.json -o .codegen --skip-validate-spec --additional-properties=packageName=openapi --additional-properties=isGoSubmodule=true --additional-properties=enumClassPrefix=true` directly, then `cp -R .codegen/*.go internal/openapi/` (the openapi-generator-cli puts every model + api file at `.codegen/` top level for `-g go`, so the flat `cp` works).
- `scripts/generate-proto.sh` needs `protoc` and `protoc-gen-go` on PATH. The validatord toolchain at `/workspace/tg-validatord/scripts/tools/latest/bin/` has both.

## Lint and cross-compilation gates

`./build.sh lint` runs `golangci-lint run` against `.golangci.yml` (which relaxes `errcheck` and a few
`staticcheck` checks for `_test.go` only — production code stays strictly linted). It reports **0 issues**;
keep it there.

**Also build for 32-bit.** `GOARCH=386 go build ./pkg/...` caught a real break: `n > math.MaxUint32` with
`n int` does not compile where `int` is 32 bits. Compare in `uint64` (`uint64(n) > math.MaxUint32`).
`go vet ./pkg/...` and `go test ./pkg/...` (7 packages) complete the set.

**Vet the WHOLE module — `go vet ./...`, not `./pkg/...` — after removing anything public.**
`go build` never compiles `_test.go`, and `./pkg/...` excludes `test/integration` and
`test/e2e`, so a deleted method leaves a compile error nothing in the standard gate reports.
Measured 2026-09-07: deleting `WhitelistedContractService.ListWhitelistedContracts` left
`test/integration/extended_services_test.go:297` broken while build, `vet ./pkg/...`, lint,
`GOARCH=386` and all 7 unit packages were green. `go vet ./...` found it in one run. The
cross-SDK version of this rule is in the repo-root `CLAUDE.md`.

### Governance reads all verify now — and there is no keyless constructor

`GetRules`, `GetRulesByID` and `GetRulesHistory` each verify the SuperAdmin signatures
before returning, and `GetDecodedRulesContainer` verifies **unconditionally**. Two things
that used to be here and are gone: the `if len(s.superAdminKeys) > 0` skip (fail-open — it
returned an unverified container, the document the HSM key is read from, with no error) and
`NewGovernanceRuleService`, the "without verification" constructor that made the skip look
necessary. Only `NewGovernanceRuleServiceWithVerification` remains.

For a consumer: a `*model.GovernanceRuleset` from any read has been verified.
`VerifyGovernanceRules(rules)` is still exported for a container decoded some other way,
and `SuperAdminKeys()` / `MinValidSignatures()` remain exported so a caller can report
signing progress without re-implementing the threshold.

**`GetRulesHistory` is LENIENT where the two single reads are strict** — it excludes and
names unverifiable entries in `ExcludedUnverified` with `TotalItems` reduced, and does not
error when nothing survives. A SuperAdmin key rotation makes every pre-rotation ruleset
unverifiable, so a strict page would deny the whole audit trail from the rotation onwards.
`tg-protect-mcpd` has the test that establishes this; it surfaces the exclusions as
`unreadable_entries`.

**Verification is memoised** (`verifiedRuleset` / `rulesetVerificationKey`), successes
only. It is not a second container cache — `cache.RulesContainerCache` already holds the
decoded+verified container for the address/asset/price paths. Do not add a third.

**The memo key must stay INJECTIVE, and that is a security property, not tidiness.** The key
decides whether ECDSA runs at all, so it is length-prefixed: an 8-byte big-endian length
before the decoded container bytes, then the signature count, then each sorted signature
length-prefixed. It was `sha256(container || 0x00 || userId || 0x00 || sig || …)`, which
leaves the container/signature boundary uncommitted — a response-controlling attacker could
shift bytes across it (moving container content into a `userId`, which verification never
reads) and make a MODIFIED container inherit a genuine one's "already verified" status,
skipping verification on the document every HSM key is read from. Interleaving a `0x00`
separator does NOT fix it; only length prefixes do. `userId` is excluded on purpose, and the
key covers the DECODED bytes so it cannot identify something other than what was verified.
An undecodable container returns `ok=false` and is not memoised. Gated by the shared
`memo_key` vectors (`pkg/protect/service/memo_key_vectors_test.go`) plus
`TestVerificationMemoKeyIsInjective`.

### Test fixtures: `rulesSignatures` is an ARRAY, and a string makes the test vacuous

A governance reply fixture with `"rulesSignatures": ""` dies in the generated client with
`json: cannot unmarshal string into ... []openapi.TgvalidatordRuleUserSignature` — *before*
verification runs. A test asserting "this read refuses an unsigned container" then passes
for the wrong reason and keeps passing when verification is removed. Caught exactly that
way here; use `[]any{}`. Same family as the wire-valid-container trap: build the fixture so
the only thing that can fail is the thing under test.

## Verification surface added in the 2026-09-04 pass

- **Request accessors return `(value, error)`.** `GetSourceAddress`, `GetDestinationAddress`,
  `GetMetadataCurrency`, `GetMetadataRequestID` and `GetAmount` return `ErrMetadataUnverified`
  when the metadata carries a payload verification has not cleared. Empty metadata (no hash, no
  payload) is NOT an error — an early-status request has nothing to read. `tg-protect-mcpd` is
  the only external consumer and was updated with it.
- **`helper.VerifiedAsset`** is the witness type returned by `VerifyAsset`. Its fields are
  unexported, so a value built outside `helper` returns `nil` from every accessor. Go cannot
  make it unforgeable — `helper.VerifiedAsset{}` compiles anywhere — so the property is
  *useless if forged*, which is what `RequestMetadata.entries` already relies on. Don't
  document it as unforgeable.
- **`helper.ParseWhitelistedAssetFromJSON`** is step 6 for assets, and
  `model.WhitelistedAsset` gained `Name`/`Symbol`/`ContractAddress`/`Decimals`/`TokenID` —
  none of which existed, so a Go caller could not obtain the verified asset identity at all.
- **`helper.ResolveRuleKey`** picks the `(blockchain, network)` pair from the signed payload.
  `MaxPayloadBytes` (1 MiB) bounds it first.
- **`ApproveRequests`/`ApproveRequest` take a variadic `comment ...string`** — variadic rather
  than a new positional parameter so existing callers keep compiling.
- **`isWildcard` is case-insensitive** (`strings.EqualFold`). It compared against `"Any"`
  exactly, which selected a different rule tier here than in the other three SDKs.

### ECDSA test signatures must pad both halves to 32 bytes

`append(r.Bytes(), s.Bytes()...)` is wrong and fails roughly 1 run in 128: `r` and `s` are
uniform in `[1, n-1]`, so one of them periodically has a leading zero byte that `Bytes()` drops,
yielding a 63-byte signature. `signature_verifier_test.go` has a `padTo32Bytes` helper — use it.
One test in that file did not, and read as an intermittent verification failure.

## Verification surface added in the 2026-09-10 security-scan pass

Go was the reference SDK for this pass, so these are the shapes the other three were ported to.
Cross-SDK reasoning is in the repo-root `CLAUDE.md`; what matters *here*:

- **`helper.LegacyPayloadVariant` carries the payload, not just the hash.**
  `ComputeLegacyPayloadVariants` / `ComputeAssetLegacyPayloadVariants` are what verification uses;
  `ComputeLegacyHashes` / `ComputeAssetLegacyHashes` remain as thin projections **because they are
  the cross-SDK vector oracle's entry point** (`cross_sdk_vectors_test.go` →
  `docs/test-vectors/crypto-test-vectors.json`). Keeping the old names intact is what let that gate
  keep passing unchanged, which is the correct outcome: the fix moves no hash.
- **`verifyHashInSignedHashes` returns `(hash, payload, err)`** in both verifiers, and
  `VerificationResult` / `AssetVerificationResult` gained `VerifiedPayload`. Step 6 parses THAT.
  `Metadata.PayloadAsString` is deliberately left untouched — a caller needs it to reproduce
  `metadata.hash`.
- **`rejectDuplicateObjectKeys`** is a structural `json.Decoder.Token()` pre-pass with per-object
  key sets (siblings sharing a key are fine), depth-bounded at 32. It sits inside
  `ParseWhitelistedAddressFromJSON` / `ParseWhitelistedAssetFromJSON` next to the existing
  `MaxPayloadBytes` guard, so a new caller cannot forget either. `DisallowUnknownFields` is the
  wrong tool and would reject real payloads.
- **`mapper.WhitelistedAddressFromDTO` now returns an error, and that ORDERING MATTERED.**
  It used to swallow the parse failure (`if err == nil && parsed != nil`), leaving every identity
  field at its zero value with no signal. That had to be fixed *before* duplicate-key rejection
  landed: step 6 parses the clean variant, so a rejected duplicate in the mapper would otherwise
  have produced a "verified" address with an empty address string — and worse, step 5 reads
  `LinkedInternalAddresses`/`LinkedWallets` from the mapper, so both empty means
  `shouldCheckRuleLines == false` and a silent fall-through to the container default thresholds.
  A fail-open through a new door. No fixture has duplicate keys, so the suite would not have
  caught it.
- **The verifying list path maps ROW BY ROW.** `WhitelistedAddressesFromDTO` (all-or-nothing) is
  kept for other callers, but `verifiedAddresses` calls the single-row mapper so an unparseable
  payload becomes one `ExcludedWhitelistedAddress` rather than failing the page.
- **`applyVerifiedIdentity`** (service) is what stopped `_, err :=` discarding the verifier's
  step-6 result on every read path except `GetWhitelistedAddressEnvelope`. `Network` is the one
  field that keeps its DTO value when the payload omits it — safe only because rule selection no
  longer lets an unsigned network reach a weaker quorum.
- **`ResolveRuleKeyWithSource`** is the source-aware form; `ResolveRuleKey` keeps its two-value
  shape because the shared `rule_key` vectors assert exactly that. Same projection trick as the
  legacy hashes, same reason.
- **`model.WhitelistedAssetResult` now exists**, so `ListWhitelistedAssets` /
  `ListWhitelistedAssetsForApproval` return a result struct rather than a bare tuple. It has no
  `ExcludedUnverified` yet: that list is still STRICT, and making it lenient needs a `logger` on
  `WhitelistedAssetService`, whose constructor takes no `opts ...ServiceOption` — a separate change.
  **This broke 5 `tg-protect-mcpd` call sites**; expect to bump it deliberately.
- **`AddressService.verifiedAddress` is the ONE seam** for a `*model.Address`. It checks the
  address STRING, not the status, because `status` is server-controlled: empty address → return
  (nothing to misuse), address + signature → verify, address without signature → refuse.
- **`ApprovePledgeActions` takes actions + key and signs inside the SDK.** It is now
  `verifies` in the signing manifest — the gate found it unprompted, which is the gate working.
  `verifyPledgeActionMetadata` covers the two list paths. Note the list paths return
  `[]taurusnetwork.PledgeAction` (values), so the approval takes values too rather than pointers.
- **`TPV1Transport` has a `Host` field and `newHTTPClient` takes the host.** Both guards are
  required (redirect refusal + host pin) for the reason already documented for `bearerTransport`
  a few sections below — and TPV1's consequence is worse, because it mints a *fresh valid
  signature* for the attacker's host rather than replaying one. `apiKeyCredentials.apply` no
  longer discards the host (`_ string`) and fails closed on an unparseable one.
  An empty `Host` means unpinned; only the pre-existing transport tests pass `""`.
- **The nil-reply guard is in TWO places that must stay in step**: `internal/openapi/client.go`
  (the committed generated file) and `scripts/resources/templates/go/client.mustache`. The repo
  had **no vendored templates** before this, and `scripts/generate-openapi.sh` does
  `rm -rf internal/openapi` before regenerating, so a hand-edit to the generated file alone is
  silently reverted on the next `./build.sh generate`. If you add another generated-client fix,
  put it in both and pass `-t` to the generator.

## Verification surface added in the 2026-09-07 pass

Cross-SDK rules are in the repo-root `CLAUDE.md` (one reader per entity; all-or-nothing
approve; approval must not sign an unverified hash; prices verified with the container
deciding). Go-specific mechanics:

- **`helper/group_threshold.go` is the ONE step-5 site.** Go previously had *three* identical
  walk methods **per verifier** (`tryVerifyAllPaths`, `verifySequentialThresholds`,
  `verifyGroupThreshold`) — six duplicated methods, all removed. `precomputeHashesJSON` was
  already package-level. Selecting *which* thresholds apply stays per-entity, because the
  address side's rule-line tree has no contract-rule counterpart.
- **`helper.VerifyAndDecodeRulesContainer`** is the shared step-2+3; `buildRulesContainerCache`
  (hash-keyed, for the address list) and `newInlineContainerCache` (base64-keyed, for the asset
  list) both go through it. The two caches differ because the **contracts list endpoint has no
  `rulesContainerNormalized` parameter and no `RulesContainers` array** — each envelope carries
  its container inline, so there is no hash to key on.
- **`helper/price_verifier.go`** — `PriceSignedBytes` / `VerifyPrice` / `VerifyPrices`. The
  `canonicalPrice` struct's field ORDER is the signed byte sequence (`encoding/json` emits
  declaration order): do not reorder or add fields. `model.PriceUpdaterKeys()` is the
  `GetHsmPublicKey` analogue, `sync.Once`-cached.
- **`model.KeyFingerprint(*ecdsa.PublicKey)`** is the canonical fingerprint;
  `helper.keyFingerprint` delegates and `RuleUser.KeyFingerprint()` caches per user.
- `MaxPayloadBytes` is enforced **inside** `ParseWhitelistedAddressFromJSON` /
  `ParseWhitelistedAssetFromJSON`, not at the call sites — a new caller cannot forget it.
- `ApproveWhitelistedAssets` sorts IDs numerically and rejects a non-numeric one (it would have
  no defined position in the signed array).

## Lessons Learned (Non-Security)

### Cell codec: the lossless guard, and nil vs RawCell

`ruleCellFromBytes` is a thin wrapper that byte-compares: it calls `ruleCellFromBytesTyped`, re-encodes the
typed result and falls back to a verbatim `RawCell` unless the bytes match. An unknown-*field* check alone is
too weak — a non-canonical encoding carries no unknown field yet still re-encodes differently, and a
`RuleIntegerGreaterNegValue` cell with payload `0x00` normalizes to zero and flips to the `Value` arm, silently
changing a signed rule's comparison operator.

**Public entry points:** `RuleCellToBytes` / `RuleCellFromBytes` / `CellFamily` are thin exported
wrappers over the lowercase internals. They exist because only the internals were exported before,
so an external Go caller could not encode or decode a `model.RuleCell` at all while Java, Python and
TS all could. `rule_cell_codec_exported_test.go` lives in package `mapper_test` on purpose — an
in-package test would pass against the unexported functions and prove nothing about the public
surface.

Two behaviours to preserve when touching it:
- **`nil` and `RawCell` are different results.** An *empty* cell under a column type this SDK does not know
  returns `nil` (nothing to preserve — the inner switch has explicit `len(data) == 0` checks); a *non-empty*
  one returns a `RawCell`. Collapsing the two breaks `TestRuleCellCodec_UnknownColumnEmptyBytesIsNil` and
  `TestRuleCellCodec_MatchAnySemantics`. Note Java differs here and returns a `RawCell` for the empty case.
- **`cellFamily(RawCell)` returns the cell's `ColumnType`**, matching Java/Python/TS. It used to
  return `""`, throwing away the only grammar information a raw cell still carries.
- **`RawCell.Payload`** is the field name (was `Bytes`), and every variant has `Kind() string`
  returning its own type name — `TestRuleCell_KindMatchesTheTypeNameForEveryVariant` pins all 36, so
  a new variant without a `Kind()` fails there rather than silently lacking a discriminator.
- The inner `ruleCellFromBytesTyped` keeps its own `raw()` closure returning a `RawCell` **without** touching
  `decodeStats`; the wrapper owns the `stats.rawCells` increment, so a cell is never counted twice.

### Deterministic marshaling applies to the JSON bridge too

`deterministicMarshal` (`rules_container_encode.go`) is the only correct marshaler on any path whose output is
hashed, signed or compared. `rules_container_json.go` used plain `proto.Marshal` in both
`RulesContainerBase64FromJSON` and `RuleMessageBase64FromJSON`; both now use `deterministicMarshal`. Note
`protojson` has no representation for unknown fields, so the JSON pair is still lossy for schema-newer data
despite its doc comment — don't advertise it as lossless.

### Error Types: ONE family now — `pkg/protect` aliases the live types

There used to be three same-named `APIError` families and only `service.APIError` was live, so every
documented `errors.Is(err, protect.ErrNotFound)` / `protect.IsAPIError(...)` branch was dead code. They
are consolidated:

- `service.APIError` (`pkg/protect/service/errors.go`) is the one real type — all 43 services funnel
  through `NewErrorMapper()` -> `MapError` -> `buildAPIError` into it. It carries `RetryAfter` and the
  `Is` / `IsRetryable` / `IsClientError` / `IsServerError` / `SuggestedRetryDelay` methods (Go had
  none of those before, while the report claimed all four SDKs implemented `isRetryable()`).
- `pkg/protect/errors.go` is now **type aliases plus sentinels**: `protect.APIError = service.APIError`,
  `AuthorizationError = service.AuthorizationError`, `IntegrityError`/`WhitelistError`/
  `ConfigurationError`/`RequestMetadataError` = the `model` ones. Its own structs and the six
  `ValidationError()`-style constructors are gone (zero callers).
- `model.APIError` + its six typed wrappers + `NewAPIError` are **deleted** (zero references).
- `errors.Is` matching works because `(*service.APIError).Is` compares by status code, with the 500
  sentinel matching any 5xx. It cannot live in `protect` — `protect` imports `service`, so the reverse
  is an import cycle. Keep the method in `service`.
- `model.IntegrityError`/`WhitelistError` have `Is` methods so `ErrIntegrity`/`ErrWhitelist` match
  regardless of message. Verification helpers still return the `model` types.

Guarded by `TestSentinelsMatchTheErrorsServicesActuallyReturn` in `pkg/protect/errors_test.go`. If a
parallel error type reappears, that test is what fails.

`(*APIError).Error()` handles description-only and empty cases explicitly — never return a bare
`" (code=404)"`, which loses all context in a log.

Mapping lives in `pkg/protect/service/errors.go` (not `wallet.go`). `buildAPIError(code, statusLine, rawBody)`
is the pure core, split out because `openapi.GenericOpenAPIError` has unexported fields in another package and
cannot be constructed in a test — test the core directly, and use one `httptest` case for the wiring.

### Embedded error types need an explicit Unwrap

`AuthorizationError` embeds `*APIError`, which promotes `(*APIError).Unwrap` returning `APIError.Err` — nil on
the mapper path. So `errors.As(err, &APIError{})` walks straight past the embedded value and reports **false**
unless the outer type defines its own `Unwrap() error { return e.APIError }`. Any future typed subtype that
embeds `*APIError` needs the same, or it silently breaks every caller matching the base type.
`TestAuthorizationErrorAsBothTargets` guards this.

### Pagination Overflow Prevention
`offsetPagination` adds with `addSaturated` (stops at `math.MaxInt64`), so no rule wraps a next
offset into a negative on caller-supplied values; `TestOffsetPaginationNeverOverflows` pins it.

### Helper Function Consolidation
All safe pointer helpers consolidated in `pkg/protect/mapper/helpers.go`.

### Go Module Version Policy
Target oldest supported Go version (1.21+). Remove `toolchain` directive when lowering.

### Request Body Preservation in Transport
Read and restore `req.Body` before cloning -- cloning consumes the body, breaking retry middleware. Pattern: `ReadAll` -> reset original with `NopCloser(NewReader)` -> clone with fresh `NopCloser(NewReader)`.

### The bearer token must not follow a redirect — and `net/http` cannot help you

`bearerTransport` attaches `Authorization` inside `RoundTrip`, i.e. BELOW redirect handling.
`http.Client.do` strips `Authorization` on a cross-host redirect, but only for headers present
on the request handed to `Do` — a header set in the transport is invisible to it, so there is
nothing to strip and the transport happily mints a fresh token for whatever host the server
named in a `Location`. (Same caveat as `golang.org/x/oauth2`'s `Transport`.) That turned "can
tamper with responses" into "holds a replayable credential", and Go was the only one of the
four SDKs affected — TS/Python/Java attach above the redirect layer where the platform strips.

Two guards, both required: `newBearerHTTPClient` sets
`CheckRedirect: return http.ErrUseLastResponse` (an API client has no reason to follow one),
and `bearerTransport` refuses to send the token to any host but the configured one — which is
what holds if a caller supplies a redirect-following policy of their own. `Credentials.apply`
therefore takes the host; an unparseable one fails closed. `newHTTPClient` /
`newBearerHTTPClient` also propagate `CheckRedirect` and `Jar` from a caller-supplied client
instead of silently dropping them. Gated by `TestBearerTokenIsNotSentToARedirectTarget` and
`TestBearerTokenIsNotSentCrossHostEvenWhenRedirectsAreFollowed` — both go red, with the token
visible in the failure message, if either guard is removed.
