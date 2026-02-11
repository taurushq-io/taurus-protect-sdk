# CLAUDE.md — TypeScript SDK

## Naming Conventions

**Taurus Product Names**: Always use hyphenated format: `Taurus-PROTECT`, `Taurus-CAPITAL`, `Taurus-EXPLORER`, `Taurus-PRIME`.

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: build + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Build only
./build.sh lint      # Linter
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
```

**Single test:** `./build.sh unit-one <pattern>` (e.g., `"should approve request"`, `"RequestService"`)

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

### Available APIs (56 total)

**Core**: `walletsApi`, `addressesApi`, `requestsApi`, `transactionsApi`, `governanceRulesApi`, `balancesApi`, `currenciesApi`, `addressWhitelistingApi`

**Blockchain Requests**: `requestsADAApi`, `requestsALGOApi`, `requestsContractsApi`, `requestsCosmosApi`, `requestsDOTApi`, `requestsFTMApi`, `requestsHederaApi`, `requestsICPApi`, `requestsMinaApi`, `requestsNEARApi`, `requestsSOLApi`, `requestsXLMApi`, `requestsXTZApi`

**Advanced**: `airGapApi`, `stakingApi`, `contractWhitelistingApi`, `businessRulesApi`, `reservationsApi`, `multiFactorSignatureApi`

**Administrative**: `usersApi`, `groupsApi`, `restrictedVisibilityGroupsApi`, `configApi`, `webhooksApi`, `webhookCallsApi`, `tagsApi`

**Specialized**: `assetsApi`, `actionsApi`, `blockchainApi`, `exchangeApi`, `fiatApi`, `feePayersApi`, `healthApi`, `jobsApi`, `scoresApi`, `statisticsApi`, `tokenMetadataApi`, `userDeviceApi`

**Taurus Network**: `client.taurusNetwork.{lendingApi, participantApi, pledgeApi, settlementApi, sharedAddressAssetApi}`

### High-Level Service Accessors (43 services: 38 core + 5 TaurusNetwork)

**Core** (8): `wallets`, `addresses`, `requests`, `transactions`, `balances`, `currencies`, `health`, `jobs`
**Administrative** (7): `users`, `groups`, `visibilityGroups`, `tags`, `webhooks`, `webhookCalls`, `audits`
**Security** (4): `governanceRules`, `whitelistedAddresses`, `whitelistedAssets`, `contractWhitelisting`
**Advanced** (9): `staking`, `reservations`, `multiFactorSignature`, `businessRules`, `airGap`, `configService`, `assets`, `changes`, `statistics`
**Blockchain & Pricing** (6): `blockchains`, `exchanges`, `prices`, `fees`, `feePayers`, `fiat`
**Specialized** (4): `scores`, `tokenMetadata`, `userDevices`, `actions`
**TaurusNetwork** (5): `client.taurusNetwork.{participants, pledges, lending, settlements, sharing}`

### TaurusNetwork Models Structure

In `src/models/taurus-network/`: participant.ts (7), pledge.ts (25), lending.ts (14), settlement.ts (11), sharing.ts (15), index.ts re-exports all 72 models.

## Code Generation

### OpenAPI Generator

Uses `openapi-generator-cli` JAR (7.9.0) with `-g typescript-fetch`. Types prefixed with `Tgvalidatord`. **Java 11+ required.**

### Protobuf Generator

Uses `protoc` with `ts-proto` plugin. Excluded from build due to import path issues.

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

**WhitelistedAddressService/WhitelistedAssetService** — use static factory for verification:
```typescript
const service = WhitelistedAddressService.withVerification(api, governanceApi, {
  superAdminKeysPem: ['-----BEGIN PUBLIC KEY-----...'],
  minValidSignatures: 2,
});
```

**GovernanceRuleService** — accepts `KeyObject[]` (not PEM strings):
```typescript
const govService = new GovernanceRuleService(governanceRulesApi, {
  superAdminKeys: [keyObject1, keyObject2],
  minValidSignatures: 2,
});
```

## Testing

`./build.sh unit` or `npm test -- --testPathPattern="tests/unit"`.

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
}
```

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
