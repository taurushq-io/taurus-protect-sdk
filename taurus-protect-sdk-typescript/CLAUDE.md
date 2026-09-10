# CLAUDE.md — TypeScript SDK

## Naming Conventions

**Taurus Product Names**: Always use hyphenated format: `Taurus-PROTECT`, `Taurus-CAPITAL`, `Taurus-EXPLORER`, `Taurus-PRIME`.

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: build + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Build only
./build.sh lint      # ESLint; config is .eslintrc.json (added — it was missing, so lint used to abort)
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
```

**Single test:** `./build.sh unit-one <pattern>` (e.g., `"should approve request"`, `"RequestService"`)

### Two gaps in the type/lint coverage — know both

**`tsc --noEmit` does NOT typecheck `tests/`.** The root `tsconfig.json` excludes them, so a change that
breaks a test file's types passes `tsc` and only fails under `ts-jest` — which reports it as a *suite that
failed to run*, so the test **count silently drops** instead of showing a failure. Watch the
`Test Suites: N failed` line, not just `Tests:`. **It recurred**: renaming the governance reads passed
`tsc --noEmit` but left `tests/unit/client/protect-client.test.ts` asserting `typeof svc.get`, which
took the whole 258-test suite out — the total silently went 1722 -> 1462 with an exit code of 1 and no
failing test named. Compare the `Tests:` total against the last known-good number after any rename.
(Making `superAdminKeysPem` required passed `tsc` and
broke `tests/unit/credentials.test.ts` exactly this way.)

**ESLint config used to be absent.** `package.json` declares `eslint`, `@typescript-eslint/*`,
`eslint-config-prettier`, `eslint-plugin-prettier` and a `"lint": "eslint src --ext .ts"` script, but no
config file existed, so lint aborted with "couldn't find a configuration file" — i.e. the gate was red,
not passing. `.eslintrc.json` now holds `eslint:recommended` + `plugin:@typescript-eslint/recommended`
with `src/internal/**` ignored; the whole SDK reports **0** findings. The declared prettier integration is
deliberately **not** enabled (it would surface formatting findings across the codebase) — enabling it is a
team decision. Note ESLint resolves `parser`/`plugins` relative to the **config file's** directory, so a
config outside the package cannot load them even with `--resolve-plugins-relative-to`.

## Architecture

### Package Structure

- **src/internal/openapi**: Auto-generated OpenAPI client (DO NOT MODIFY)
- **src/internal/proto**: Auto-generated protobuf (excluded from build — import path issues)
- **src/**: Public SDK package
  - **client.ts**: ProtectClient main entry point
  - **errors.ts**: Exception hierarchy (APIError, IntegrityError, etc.)
  - **crypto/**: TPV1 authentication and cryptographic utilities
  - **services/**: BaseService and implementations
  - **models/**: Domain models
  - **mappers/**: DTO to model conversion
  - **helpers/**: Signature verification, constant-time comparison
  - **cache/**: RulesContainerCache
  - **transport/**: TPV1 middleware

### Key Patterns

- Factory pattern: `ProtectClient.create(config)`
- Lazy initialization of API instances and services via getters
- Error hierarchy: APIError (base), ValidationError, AuthenticationError, IntegrityError
- BaseService pattern with `execute()` for error handling
- Internal errors throw `ServerError` (not generic `Error`) for `instanceof APIError` catching

### Available APIs (61 total)

**Core**: `walletsApi`, `addressesApi`, `requestsApi`, `transactionsApi`, `balancesApi`, `currenciesApi`
(59 are public accessors — `governanceRulesApi` and `addressWhitelistingApi` are deliberately private; see below)

**Blockchain Requests**: `requestsADAApi`, `requestsALGOApi`, `requestsContractsApi`, `requestsCosmosApi`, `requestsDOTApi`, `requestsFTMApi`, `requestsHederaApi`, `requestsICPApi`, `requestsMinaApi`, `requestsNEARApi`, `requestsSOLApi`, `requestsXLMApi`, `requestsXTZApi`

**Advanced**: `airGapApi`, `stakingApi`, `contractWhitelistingApi`, `businessRulesApi`, `reservationsApi`, `multiFactorSignatureApi`

**Administrative**: `usersApi`, `groupsApi`, `restrictedVisibilityGroupsApi`, `configApi`, `webhooksApi`, `webhookCallsApi`, `tagsApi`

**Specialized**: `assetsApi`, `actionsApi`, `blockchainApi`, `exchangeApi`, `fiatApi`, `feePayersApi`, `healthApi`, `jobsApi`, `scoresApi`, `statisticsApi`, `tokenMetadataApi`, `userDeviceApi`

**Taurus Network**: `client.taurusNetwork.{lendingApi, participantApi, pledgeApi, settlementApi, sharedAddressAssetApi}`

### High-Level Service Accessors (43 services: 38 core + 5 TaurusNetwork)

**Core** (8): `wallets`, `addresses`, `requests`, `transactions`, `balances`, `currencies`, `health`, `jobs`
**Administrative** (7): `users`, `groups`, `visibilityGroups`, `tags`, `webhooks`, `webhookCalls`, `audits`
**Security** (4): `governanceRules`, `whitelistedAddresses`, `whitelistedAssets`, `contractWhitelisting` (**writes only** — reads live on `whitelistedAssets`, the verified reader of the same endpoint)
**Advanced** (9): `staking`, `reservations`, `multiFactorSignature`, `businessRules`, `airGap`, `configService`, `assets`, `changes`, `statistics`
**Blockchain & Pricing** (6): `blockchains`, `exchanges`, `prices`, `fees`, `feePayers`, `fiatAccounts`
**Specialized** (4): `scores`, `tokenMetadata`, `userDevices`, `actions`
**TaurusNetwork** (5): `client.taurusNetwork.{participants, pledges, lending, settlements, sharing}`

### TaurusNetwork Models Structure

In `src/models/taurus-network/`: participant.ts (7), pledge.ts (25), lending.ts (14), settlement.ts (11), sharing.ts (15), index.ts re-exports all 72 models.

## Code Generation

### OpenAPI Generator

Uses `openapi-generator-cli` JAR (7.9.0) with `-g typescript-fetch`. Types prefixed with `Tgvalidatord`. **Java 11+ required.**

### Protobuf Generator

Uses `protoc` with `ts-proto` plugin. Excluded from build due to import path issues.

**Imports stay un-flattened after `generate-proto.sh`.** ts-proto generates relative imports based on the original proto directory layout (e.g., `import { CommitmentKind } from "./tp_messages/commitments";` for `request_reply.proto`). The script's flatten step moves the file but leaves the import path intact, breaking `tsc` even though `tsconfig.json` excludes `src/internal/proto/request*.ts` from compilation — TypeScript still reads referenced files during dep resolution. Run after the regen:

```bash
sed -i 's|from "\./tp_messages/|from "./|g' src/internal/proto/*.ts
```

Currently only `tp_messages/commitments.proto` ships as a non-Google subdir import; Google well-known imports (`./google/protobuf/timestamp` etc.) stay nested correctly because ts-proto resolves those via `node_modules`.

**`--ts_proto_opt=unknownFields=true` is required** in `generate-proto.sh` — the governance typed mapper (see repo-root `CLAUDE.md` → "Governance Rules Typed API") relies on ts-proto emitting a `_unknownFields` bag so unknown protobuf fields survive a decode→encode round-trip. Gotcha: `_unknownFields` is keyed by the full wire **tag** `(fieldNumber << 3) | wireType`, NOT the field number (e.g. field 500, varint → key `4000`). Encode/decode a cell with `X.encode(X.fromPartial({...})).finish()` / `X.decode(bytes)`. Cross-SDK cell byte-parity is enforced by `tests/unit/mappers/governance-cell-vectors.test.ts` against `scripts/resources/governance-cell-vectors.json`.

**Nested `_unknownFields` footgun (encoder).** ts-proto's `X.fromPartial({ children: [...] })` re-runs each child's `fromPartial`, which does **not** copy `_unknownFields` — so building a parent that way silently strips every nested node's preserved unknown fields (only the top-level container survives). `protobuf-rules-container-encode.ts` avoids this by building nested messages, assigning them **after** the parent `fromPartial`, and re-attaching via `attach()` (see the comment at the top of that file). Keep that pattern for any new nested node, and assert *nested* survival in the round-trip test — a top-level-only assertion passes while the invariant is broken.

This footgun bit again for the four nested rule-detail sub-messages: passing `evmCallContract` &co.
into `PbDetails.fromPartial({...})` dropped their `_unknownFields` and the round trip was not lossless.
They are now assigned onto `pbDetails` **after** `fromPartial` and wrapped in `attach()`.

**Map ordering is forced inside the generated encoders — not by the callers.** ts-proto emits map
entries via `Object.entries(message.properties).forEach(...)`, and **JavaScript enumerates
integer-like keys first, in ascending numeric order, no matter how the object was built.** So for
keys `{"10","2","a"}` stock ts-proto writes `2,10,a` while Go
(`Deterministic:true`), Java (`useDeterministicSerialization`) and Python (`deterministic=True`) all
sort lexicographically and write `10,2,a`. Those are the bytes a SuperAdmin signs on a rules
proposal, so it is a cross-SDK signature mismatch.

**Sorting in the caller cannot fix this** — that was the old `sortedMap()` / `sortPropertiesMaps()`
approach, and it silently did nothing for numeric keys, because re-inserting keys into a fresh plain
object does not change enumeration order. Both helpers are **deleted**; do not reintroduce them.
`scripts/generate-proto.sh` now patches all 7 `properties`-map encode sites in the generated
`request_reply.ts` to `.sort(...)` before `.forEach(...)`. If you regenerate, the script re-applies
it and logs the site count; if the count is ever 0 it warns, because that means the guarantee is
gone. `tests/unit/mappers/rules-container-determinism.test.ts` fails loudly if the patch is lost
(verified: reverting the generated file to stock makes it fail with `2,10,a`).

**`ruleSourceToBytes` is exported on purpose.** It was module-private, so the RuleSource lossless test
could only assert the *decode*; Go, Python and Java all re-encode and byte-compare. A regression that
dropped `RuleSource.raw` on encode passed here while breaking the other three SDKs' byte parity.

**Known TS-only wire divergence: empty packed repeated fields.** The generated `User.encode` writes
`writer.uint32(26).fork() … ldelim()` for `roles` *unconditionally*, so a user with **zero roles** makes
TS emit two extra bytes (`1a 00`) that Go/Java/Python omit. With roles present all four are byte-identical
(verified). Fixing it means patching `src/internal/proto` (DO NOT MODIFY) or adding a post-regen `sed` to
`generate-proto.sh` — **unresolved**; be aware if a container with role-less users must hash identically
across SDKs.

**Public API note:** `superAdminKeysPem` is **required** on `ProtectClientConfig` (not optional), so the
mandatory-keys invariant is enforced at compile time as in the other three SDKs. Don't add `?` or `?? []`
defaults back — an empty key set must not be representable downstream.

**The rules cache fetches through `governanceRules`, NOT `governanceApi()`.**
`ProtectClient.getRulesCache()` used to call the generated `ruleServiceGetRules` and
`rulesContainerFromBase64` directly, so the container was decoded **unverified** and
`rulesSignatures` was discarded — `AddressService` then verified addresses against an
HSM public key nothing had authenticated, and anyone able to influence that response
could have attacker-chosen addresses verify clean. Go, Java and Python all fetch through
their governance service; this SDK was the only one that did not, and no SDK tested the
path. The provider is now `() => this.governanceRules.getDecodedRulesContainer()`, which
verifies before decoding. Gated by `tests/unit/cache/rules-cache-verification.test.ts`
(all three cases fail if the provider is reverted). Note the `?? 1` default on
`minValidSignatures` at each service construction site is what keeps
`GovernanceRuleService`'s `minValidSignatures > 0` verification gate always true via the
public client — don't default it to 0.

**`addressWhitelistingApi` is NOT public either — do not re-add it.** Same hole, same shape:
the raw API's `whitelistServiceGetWhitelistedAddress` and its list siblings return the
unverified envelope straight off the wire, so a caller reaching it got attacker-controllable
`address`/`label`/`memo` with no signature check — a complete bypass of the one verified reader,
in this SDK alone (Go hides its generated client under `internal/`, Python under `_internal`,
Java never exposes it). It is now the private `rawAddressWhitelistingApi()`. **The rename is
load-bearing**: TypeScript `private` is erased at runtime, so keeping the old name would leave
`"addressWhitelistingApi" in client` true and make the reachability test vacuous — the same
reason `governanceRulesApi` became `governanceApi`. Guarded by "whitelisted-address low-level
API is not publicly reachable" in `tests/unit/client/protect-client.test.ts`; the api-getter
count there is now **54**, and `apiGetters` omits both.

**`governanceRulesApi` is NOT public — do not re-add it.** It is now the private
`governanceApi()`. The raw API's `ruleServiceGetRules` returns an unverified DTO and
`ruleServiceUpdateRulesProposal` accepts an arbitrary base64 blob, so exposing it made client-side
verification opt-out in this SDK alone (Go's generated client is unreachable under `internal/`).
Guarded by "governance low-level API is not publicly reachable" in
`tests/unit/client/protect-client.test.ts`. That file also asserts the api-getter **count** (now 54,
with `apiGetters` omitting governance AND address-whitelisting) — update both if the surface changes.

**`close()` must empty `authMiddleware`.** The api secret is captured inside the auth-middleware
closure and JavaScript cannot zero a captured string, so the only thing that helps is making it
unreachable: `close()` sets `this.authMiddleware.length = 0` on the array `Configuration` shares.
The old code overwrote `this.config.apiSecret`, which is `undefined` under the supported
`credentials:` path — so `close()` released nothing at all. Keep both.

**Renamed governance reads:** `getRules` / `getRulesById` / `getRulesProposal` / `getRulesHistory`
(were the bare `get`/`getById`/`getProposal`/`getHistory`), plus `getPublicKeys` and a
`SuperAdminPublicKey` model that did not exist here. `RawCell` carries `payload`, not `bytes`.
`verifySignatureWithKey` is the single-key verify the other three already had.

**Governance test facts.** `cellFamily` / `ruleCellToBytes` / `ruleCellFromBytes` are exported (test directly); `bigintToMagnitude` / `magnitudeToBigint` / `blockchain{To,From}Str` are module-private → exercise them through the public codec API. `src/models/governance-rules.ts` helpers (`findUserById`, `findGroupById`, `find*WhitelistingRules`, `getHsmPublicKey`, `hasUnknownFields`, `createEmptyRulesContainer`) are **free functions taking the container as arg 1**, not methods; `getHsmPublicKey` returns the raw PEM string (no parsing, caching, or validation — so only found/not-found map, not the Go/Python cached/invalid-PEM cases). An unknown cell-level blockchain round-trips as a numeric-passthrough typed cell (the raw number kept as a decimal string), so it stays typed rather than falling back to `RawCell`. Unlike `tg-protect-gui`, `jest.config.js` uses `ts-jest` **without** `isolatedModules`, so type errors fail the whole suite (they are NOT masked).

Signature tests: ECDSA is non-deterministic — verify a signature (and assert a tampered payload fails), don't byte-compare two signings. `crypto` `verifySignature(publicKey, data, sigB64)` takes the public key first.

## Common Implementation Notes

### OpenAPI Type Naming Conventions

- Types prefixed with `Tgvalidatord` (e.g., `TgvalidatordWallet`)
- Response types use `result` field
- Method names: `{service}Service{Operation}` (e.g., `walletServiceGetWalletsV2`)
- Pagination: `totalItems` and `offset` fields

### Model Export Pattern

```typescript
import { Wallet, WalletStatus } from '@taurushq/protect-sdk';
import { Participant, Pledge, Loan } from '@taurushq/protect-sdk';  // TaurusNetwork
```

### Service Export Pattern (Avoiding Duplicate Exports)

**Rule:** `src/services/index.ts` exports service classes ONLY. Types belong in `src/models/`.

```typescript
// CORRECT
export { StakingService } from "./staking-service";
// WRONG — causes TS2308 duplicate export
export { StakingService, type ADAStakePoolInfo } from "./staking-service";
```

If types are in a service file, move them: create `src/models/{domain}.ts`, update imports, export from `src/models/index.ts`.

Barrel pattern: `services/index.ts` and `services/taurus-network/index.ts` export classes only; `models/index.ts` and `models/taurus-network/index.ts` export types.

### OpenAPI Method Naming Notes

- Cancel request: `requestServiceCreateOutgoingCancelRequest` (not `...CreateCancelRequest`)
- `TgvalidatordInternalUser`: `firstName`, `lastName`, `email` (no `name`)
- `TgvalidatordAuditTrail`: `user?.email` via nested object (no `userEmail`)
- `TgvalidatordGetConfigTenantReply`: `config` property (not `result`)

### Model Property Names

- **DecodedRulesContainer**: `users`, `groups`, `minimumDistinctUserSignatures`, `minimumDistinctGroupSignatures`, `addressWhitelistingRules`, `contractAddressWhitelistingRules`, `enforcedRulesHash`, `timestamp` (NO `requestRules`)
- **ProtectClient**: `superAdminKeysPem` (not `superAdminKeys`)

### Verification Service Patterns

**WhitelistedAddressService/WhitelistedAssetService** — `withVerification(api, config)`, two
arguments (there is no separate governanceApi parameter):
```typescript
const service = WhitelistedAddressService.withVerification(api, {
  superAdminKeysPem: ['-----BEGIN PUBLIC KEY-----...'],
  minValidSignatures: 2,
  rulesContainerDecoder,
  userSignaturesDecoder,
});
```
Verification is **not optional**: `get`, `list` and `getEnvelope` all run the full pipeline on
both services. `WhitelistedAssetService` used to build a verifier and then never call it on
those paths — do not reintroduce a read path that skips `verify()`.

**GovernanceRuleService** — accepts `KeyObject[]` (not PEM strings):
```typescript
const govService = new GovernanceRuleService(governanceRulesApi, {
  superAdminKeys: [keyObject1, keyObject2],
  minValidSignatures: 2,
});
```

## Testing

`./build.sh unit` or `npm test -- --testPathPattern="tests/unit"`.

### There is NO network-stub library — prototype-spy the generated APIs

No nock, no msw, no fetch mock anywhere in `tests/unit`. Two patterns exist instead:

- **Service tests** inject hand-rolled mocks (`{...} as unknown as jest.Mocked<XApi>`),
  as in `tests/unit/services/address-service.test.ts`.
- **Client-level tests** construct a real client with `ProtectClient.create(...)` (see
  `tests/unit/client/protect-client.test.ts`) and, to get past the transport,
  `jest.spyOn(SomeApi.prototype, "methodName").mockResolvedValue({...} as never)`. That
  is how `tests/unit/cache/rules-cache-verification.test.ts` drives the real client
  wiring offline — necessary there because the defect *was* the wiring, so every
  injected-mock service test passed throughout. `as never` satisfies
  `mockResolvedValue` without building a full generated reply type; eslint only runs on
  `src`, but ts-jest still typechecks tests, so the cast is required.

Use a **wire-valid** container (`rulesContainerToBase64(createEmptyRulesContainer())`)
when testing that verification runs — a malformed one throws from the decoder instead,
and the test then passes whether verification ran or not.

**Unit Tests** (`tests/unit/`): Jest + ts-jest. Key files: `services/request-service.test.ts`, `services/governance-rule-service.test.ts`.

### Shared Test Utilities (`tests/testutil/`)

All test config is centralized in `tests/testutil/`:
- `properties.ts` — Key=value `.properties` file parser with `\n` escape for PEM keys
- `config.ts` — `TestConfig` with multi-identity support (6 identities), env var overrides
- `helpers.ts` — `skipIfNotEnabled()`, `getTestClient(index)`, `skipIfInsufficientIdentities()`
- `test.properties.sample` — Sample config matching Java format
- `index.ts` — Barrel exports

Integration config (`tests/integration/config.ts`, `helpers.ts`) remains unchanged for backward compatibility.

**Integration Tests** (`tests/integration/`): Require live API. May fail for non-code reasons (no base currency, protobuf format, API version).

**To TYPECHECK an integration test without an API, run jest on it and read the suite line.**
`tsc --noEmit` excludes `tests/`, so nothing in the normal gate compiles these files — a public
API removal leaves a break that only shows when someone runs them with live credentials.
`npx jest tests/integration/<file>` loads and typechecks the file, then skips the cases:
expect `Test Suites: 1 skipped`. **`1 failed` with "failed to run" is the type error** — and
note it drops the test count rather than showing a red test, the same trap as a broken unit
file. Caught a real one 2026-09-07 (`extended-services.test.ts` still calling a deleted
`contractWhitelisting.list`).

### Integration Test Guidelines

- Use high-level services (`client.addresses`) not raw APIs (`client.addressesApi`)
- Add timeouts: 30s for address list, 60s for pagination
- `Address.id` is `number` (not `string` like `TgvalidatordAddress.id`)

**Common field name mistakes:**

| Model | Wrong | Correct |
|---|---|---|
| `ListExchangesResult` | `.exchanges` | `.items` |
| `FiatProvider` | `.name` | `.provider` |
| `ParticipantService.list()` | Wrapped result | `Participant[]` directly |
| `StakeAccountType` | `"all"` | `"StakeAccountTypeSolana"` |

## Error Handling Pattern

```typescript
try {
  const wallets = await client.walletsApi.walletServiceGetWalletsV2();
} catch (error) {
  if (error instanceof IntegrityError) { /* security issue */ }
  else if (error instanceof APIError && error.isRetryable()) { /* retry */ }
  else if (error instanceof AuthorizationError) { /* error.requiredRoles */ }
}
```

**`BaseService.handleError` and `mapResponseError` are `async`** (`Promise<never>` / `Promise<APIError>`)
because the server's message and `error_code` live in the response body, and reading it (`response.clone().text()`
+ guarded `JSON.parse`) is async. Without that read the caller only ever sees the fetch layer's
`Response returned an error code`.

Two consequences: `execute` must `return await this.handleError(error)` — a bare call leaves a floating promise
and returns `undefined`, which type-checks under `Promise<T>` and fails silently. And because `handleError` is
`protected` on an exported class, the signature is public API: note it in release notes if it changes again.

## Verification surface added in the 2026-09-04 pass

- **`Verified<T>` (`src/helpers/verified.ts`)** — a non-exported `unique symbol` brand.
  `attestVerified` is the only producer and is not re-exported from the package barrel, so a
  caller cannot mint one without an explicit `as`. Erases at runtime. `getEnvelope` on the asset
  service returns `Verified<SignedWhitelistedAssetEnvelope>`.
  `tests/unit/helpers/verified.test.ts` guards it with `@ts-expect-error`: weakening the brand
  produces `TS2578: Unused '@ts-expect-error' directive` and the suite reports `Tests: 0`, so
  **watch the total, not just the pass/fail line**.
- **`isCryptoVerificationError` (`src/crypto/errors.ts`)** replaced five copies of
  `message.includes('key') || …`. That predicate swallowed ANY error containing the word —
  including a `TypeError` from a real bug, silently reported as a failed signature. It now keys
  on Node's stable `error.code` (`ERR_OSSL*`, `ERR_CRYPTO*`) and only falls back to message
  matching for this SDK's OWN messages, anchored with `startsWith`, never `includes`.
- **`containsHash` moved to `whitelist-hash-helper.ts`** and is exported. Each verifier had its
  own private copy; there are now exactly two hash-comparison functions in the SDK.
- **`verifyAddressSignature` throws, and no longer returns a boolean.** It also checks for an
  empty address, which it did not. `verifyAddressSignatures` fails on the first bad address
  instead of returning a `boolean[]` a caller could ignore.
- **`WhitelistedAssetService` verifies on every read path.** It built a verifier and used it at
  exactly one site; `get`/`list`/`getEnvelope` returned attacker-controllable payload data
  labelled "verified". `mapEnvelopeToAsset` is deleted — **do not reintroduce a mapper that
  takes a bare envelope.**
- **List results carry `excludedUnverified`.** A short page and a filtered page must not look
  the same; rows-returned-but-none-surviving throws `IntegrityError`.
- `resolveRuleKey` and `MAX_PAYLOAD_BYTES` live in `whitelist-hash-helper.ts`.

## Verification surface added in the 2026-09-07 pass

Cross-SDK rules are in the repo-root `CLAUDE.md`. TypeScript-specific:

- **`src/helpers/price-verifier.ts`** — `priceSignedBytes` / `verifyPrice` / `verifyPrices`.
  The canonical form is built by hand in validatord's field order, and
  `tests/unit/helpers/price-verifier.test.ts` pins it **by value** against the literal JSON:
  a reordered or extended object silently changes every signature this SDK accepts.
- **`keyFingerprint` is exported** from `signature-verifier.ts` — both whitelist verifiers'
  `verifyGroupThreshold` key their signer `Set` on it.
- **`mappers/contract-whitelist.ts` is attribute mappers ONLY.** `extractContractDetailsFromPayload`
  is deleted: it parsed the signed payload and returned its contents as fact, which made this
  SDK's copy of the bypass worse than the other three. `models/contract-whitelist.ts` is down
  to `WhitelistedContractAttribute` plus the two request types — the 12 read-side interfaces
  are gone. **Do not reintroduce a mapper or model that takes a bare contract envelope.**
- `whitelistedAssets` gained `listForApproval` and `approve`; `list`/`listForApproval` share a
  private `verifyPage`.
- `AssetService` and `PriceService` now take `rulesCache` as a **mandatory** constructor
  param, so `client.ts` must build it before them — the same shape `AddressService` already
  had. Service tests for those two need a `createMockRulesCache()` stub or the suite reports
  "failed to run" (and the count drops, per the note above).

## `List*Options` for whitelisting live in the SERVICE file, not `models/`

`ListWhitelistedAddressesOptions` (and the new
`ListWhitelistedAddressesForApprovalOptions`) are declared in
`src/services/whitelisted-address-service.ts` and re-exported as `type` from
`src/services/index.ts`. That contradicts the "services export classes ONLY, types live in
models" rule above — it is a pre-existing local exception, so follow it when adding a
sibling options type rather than relocating one and breaking the barrel. Grepping
`src/models/` for these names finds nothing, which is the confusing part.

## Governance verification is mandatory at construction

`GovernanceRuleServiceConfig` now REQUIRES `superAdminKeys` and a positive
`minValidSignatures`, and the constructor throws `ConfigurationError` on either. `config`
itself is no longer optional. The previous `config?.minValidSignatures ?? 0` default plus
the `minValidSignatures > 0` call-site gates made a zero-keys service constructible and
silently unverified — and the class is exported from the package barrel, so any consumer
could reach it. Do not re-add the defaults.

Consequence for tests: `new GovernanceRuleService(api)` no longer compiles, which surfaces
as **"Test suite failed to run"** and a DROPPED test count rather than a red test — the
trap this file already documents for `tsc` not covering `tests/`. `governance-rule-proposal.test.ts`
needed a `verifyingConfig()` helper for exactly this reason (proposal paths never verify,
so any P-256 public key is a valid fixture there).

## Lessons Learned (Non-Security)

### OpenAPI Type Compatibility

Use nullish coalescing for optional fields: `body: { comment: request.comment ?? '' }`.

### Service Implementation Checklist

1. **Model** (`src/models/{entity}.ts`) — interfaces with `readonly` properties
2. **Mapper** (`src/mappers/{entity}.ts`) — use `safeString`, `safeInt`, `safeBool`, `safeDate`, `safeMap`
3. **Service** (`src/services/{entity}-service.ts`) — extend `BaseService`, use `this.execute()`
4. **Barrel exports** — models in `models/index.ts`, service class in `services/index.ts`
5. **Client** (`src/client.ts`) — import, private field, lazy getter, cleanup in `close()`

### OpenAPI Trail Types

`TgvalidatordTrail.date` is `Date` (not `string`). Match generated types exactly.

### Model/Service Export Separation

TS2308/TS2724 errors when types exported from both models and services. Fix: services export only classes, types come from models only.
