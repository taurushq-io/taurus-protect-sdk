# Key Concepts - TypeScript SDK

This document covers TypeScript SDK-specific aspects of Taurus-PROTECT concepts. For the full domain model documentation, see the [Common Concepts](../../docs/CONCEPTS.md) document.

## Common Documentation

The following documentation applies to all Taurus-PROTECT SDKs:

- [Key Concepts](../../docs/CONCEPTS.md) - Domain model, entities, relationships
- [Authentication & TPV1](../../docs/AUTHENTICATION.md) - API authentication protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Cryptographic verification flows

---

## TypeScript Model Interfaces

The TypeScript SDK represents domain entities as immutable TypeScript interfaces in `src/models/`. All properties use the `readonly` modifier to ensure immutability.

### Wallet

```typescript
interface Wallet {
  readonly id: string;
  readonly name: string;
  readonly status: WalletStatus;
  readonly type: WalletType;
  readonly blockchain: string | undefined;
  readonly network: string | undefined;
  readonly currency: string | undefined;
  readonly isOmnibus: boolean;
  readonly balance: WalletBalance | undefined;
  readonly createdAt: Date | undefined;
  readonly updatedAt: Date | undefined;
  readonly comment: string | undefined;
  readonly customerId: string | undefined;
  readonly addressesCount: number | undefined;
  readonly attributes: WalletAttribute[];
  readonly visibilityGroupId: string | undefined;
  readonly externalWalletId: string | undefined;
  readonly tags: string[];
}
```

### Address

```typescript
interface Address {
  readonly id: string;
  readonly walletId: string;
  readonly address: string;
  readonly alternateAddress: string | undefined;
  readonly label: string | undefined;
  readonly comment: string | undefined;
  readonly currency: string;
  readonly customerId: string | undefined;
  readonly externalAddressId: string | undefined;
  readonly addressPath: string | undefined;
  readonly addressIndex: string | undefined;
  readonly nonce: string | undefined;
  readonly status: string | undefined;
  readonly signature: string | undefined;
  readonly disabled: boolean;
  readonly canUseAllFunds: boolean;
  readonly createdAt: Date | undefined;
  readonly updatedAt: Date | undefined;
  readonly attributes: AddressAttribute[];
  readonly linkedWhitelistedAddressIds: string[];
}
```

### Request

```typescript
interface Request {
  readonly id: number;
  readonly type: RequestType | string;
  readonly status: RequestStatus | string;
  readonly metadata: RequestMetadata | undefined;
  readonly hash: string;
  readonly payloadAsString: string | undefined;
  readonly tenantId: string | undefined;
  readonly currency: string | undefined;
  readonly currencyInfo: CurrencyInfo | undefined;
  readonly memo: string | undefined;
  readonly rule: string | undefined;
  readonly externalRequestId: string | undefined;
  readonly requestBundleId: string | undefined;
  readonly needsApprovalFrom: string[];
  readonly approvers: Approvers | undefined;
  readonly createdAt: Date | undefined;
  readonly updatedAt: Date | undefined;
  readonly tags: string[];
}
```

### RequestMetadata

The `RequestMetadata` interface contains cryptographic hash information for integrity verification:

```typescript
interface RequestMetadata {
  /** SHA-256 hash of the request payload (hex-encoded) */
  readonly hash: string;
  /** The raw payload used for hash computation */
  readonly payloadAsString: string | undefined;
  // SECURITY: payload field intentionally omitted.
  // Use JSON.parse(payloadAsString) for secure data extraction.
}
```

**Note:** The `payload` field is intentionally omitted for security. A parsed payload object could be tampered with by an attacker while `payloadAsString` remains unchanged (hash still verifies). Use `JSON.parse(metadata.payloadAsString)` to securely extract data after hash verification.

### Transaction

```typescript
interface Transaction {
  readonly id: string;
  readonly direction: TransactionDirection | string | undefined;
  readonly currency: string | undefined;
  readonly currencyInfo: TransactionCurrencyInfo | undefined;
  readonly blockchain: string | undefined;
  readonly network: string | undefined;
  readonly hash: string | undefined;
  readonly block: string | undefined;
  readonly confirmationBlock: string | undefined;
  readonly amount: string | undefined;
  readonly amountMainUnit: string | undefined;
  readonly fee: string | undefined;
  readonly feeMainUnit: string | undefined;
  readonly type: string | undefined;
  readonly status: TransactionStatus | string | undefined;
  readonly isConfirmed: boolean;
  readonly sources: TransactionAddressInfo[];
  readonly destinations: TransactionAddressInfo[];
  readonly transactionId: string | undefined;
  readonly uniqueId: string | undefined;
  readonly requestId: string | undefined;
  readonly requestVisible: boolean | undefined;
  readonly receptionDate: Date | undefined;
  readonly confirmationDate: Date | undefined;
  readonly attributes: TransactionAttribute[];
  readonly forkNumber: string | undefined;
}
```

### Balance

```typescript
interface AssetBalance {
  readonly currencyId?: string;
  readonly currency?: string;
  readonly blockchain?: string;
  readonly network?: string;
  readonly contractAddress?: string;
  readonly tokenId?: string;
  readonly balance?: string;
  readonly fiatValue?: string;
  readonly fiatCurrency?: string;
}
```

### GovernanceRules

```typescript
interface GovernanceRules {
  /** Base64-encoded rules container */
  readonly rulesContainer: string | undefined;
  /** List of SuperAdmin signatures on the rules */
  readonly rulesSignatures: RuleUserSignature[];
  /** Whether the rules are locked */
  readonly locked: boolean;
  /** Creation timestamp */
  readonly creationDate: Date | undefined;
  /** Last update timestamp */
  readonly updateDate: Date | undefined;
  /** Audit trail of rule changes */
  readonly trails: RulesTrail[];
}
```

### DecodedRulesContainer

The decoded rules container contains all governance rules including users, groups, and whitelisting thresholds:

```typescript
interface DecodedRulesContainer {
  readonly users: RuleUser[];
  readonly groups: RuleGroup[];
  readonly minimumDistinctUserSignatures: number;
  readonly minimumDistinctGroupSignatures: number;
  readonly addressWhitelistingRules: AddressWhitelistingRules[];
  readonly contractAddressWhitelistingRules: ContractAddressWhitelistingRules[];
  readonly enforcedRulesHash: string | undefined;
  readonly timestamp: number;
}
```

---

## TypeScript Enums

### WalletStatus

```typescript
enum WalletStatus {
  UNKNOWN = 'UNKNOWN',
  ACTIVE = 'ACTIVE',
  PENDING = 'PENDING',
  DISABLED = 'DISABLED',
}
```

### WalletType

```typescript
enum WalletType {
  UNKNOWN = 'UNKNOWN',
  STANDARD = 'STANDARD',
  AIRGAP = 'AIRGAP',
  STAKING = 'STAKING',
}
```

### AddressStatus

```typescript
enum AddressStatus {
  UNKNOWN = "UNKNOWN",
  CREATED = "CREATED",
  CREATING = "CREATING",
  SIGNED = "SIGNED",
  OBSERVED = "OBSERVED",
  CONFIRMED = "CONFIRMED",
  ACTIVE = "ACTIVE",
  PENDING = "PENDING",
  DISABLED = "DISABLED",
}
```

### RequestStatus

The RequestStatus enum has 30+ values. Common statuses include:

```typescript
enum RequestStatus {
  UNKNOWN = "UNKNOWN",
  CREATED = "CREATED",
  PENDING = "PENDING",
  APPROVING = "APPROVING",
  APPROVED = "APPROVED",
  HSM_READY = "HSM_READY",
  HSM_SIGNED = "HSM_SIGNED",
  BROADCASTING = "BROADCASTING",
  BROADCASTED = "BROADCASTED",
  MINED = "MINED",
  CONFIRMED = "CONFIRMED",
  REJECTED = "REJECTED",
  CANCELED = "CANCELED",
  EXPIRED = "EXPIRED",
  PERMANENT_FAILURE = "PERMANENT_FAILURE",
  // ... and others (READY, SENT, BUNDLE_*, HSM_FAILED, etc.)
  // See src/models/request.ts for the full list
}
```

### RequestType

```typescript
enum RequestType {
  UNKNOWN = "UNKNOWN",
  TRANSFER = "TRANSFER",
  SIGN = "SIGN",
  CONTRACT_CALL = "CONTRACT_CALL",
  PAYMENT = "payment",
  FUNCTION_CALL = "function_call",
}
```

### TransactionStatus

```typescript
enum TransactionStatus {
  NEW = "NEW",
  SUCCESS = "SUCCESS",
  PENDING = "PENDING",
  TEMPORARY_FAILURE = "TEMPORARY_FAILURE",
  INVALID = "INVALID",
  TIMEOUT = "TIMEOUT",
  MINED = "MINED",
  LOST = "LOST",
  EXPIRED = "EXPIRED",
}
```

### TransactionDirection

```typescript
enum TransactionDirection {
  INCOMING = "incoming",
  OUTGOING = "outgoing",
}
```

---

## TypeScript Exception Types

The SDK uses a comprehensive error hierarchy for different error scenarios:

### Error Hierarchy

```
Error
├── APIError (base for all HTTP errors)
│   ├── ValidationError (400)
│   ├── AuthenticationError (401)
│   ├── AuthorizationError (403)
│   ├── NotFoundError (404)
│   ├── RateLimitError (429)
│   └── ServerError (5xx)
├── IntegrityError (hash/signature verification)
├── WhitelistError (whitelist verification)
├── ConfigurationError (SDK configuration)
└── RequestMetadataError (metadata parsing)
```

### APIError

Base class for all API errors with HTTP status codes:

```typescript
class APIError extends Error {
  readonly statusCode: number;
  readonly errorCode: string | undefined;
  readonly body: unknown | undefined;
  readonly retryAfterMs: number | undefined;
  readonly cause: Error | undefined;

  isRetryable(): boolean;
  isClientError(): boolean;
  isServerError(): boolean;
  suggestedRetryDelayMs(): number | undefined;
}
```

### Specific Error Types

| Exception | HTTP Status | When Thrown |
|-----------|-------------|-------------|
| `ValidationError` | 400 | Input validation failed |
| `AuthenticationError` | 401 | Invalid or missing credentials |
| `AuthorizationError` | 403 | Insufficient permissions; `requiredRoles` names the roles that would satisfy the check |
| `NotFoundError` | 404 | Resource not found |
| `RateLimitError` | 429 | Rate limit exceeded |
| `ServerError` | 5xx | Server-side error |

### Security Exceptions

| Exception | When Thrown |
|-----------|-------------|
| `IntegrityError` | Hash/signature verification failures (NEVER retry) |
| `WhitelistError` | Whitelist verification failures |
| `ConfigurationError` | Invalid SDK configuration |
| `RequestMetadataError` | Metadata parsing failures |

### Example Usage

```typescript
import {
  ProtectClient,
  APIError,
  IntegrityError,
  NotFoundError,
  RateLimitError,
} from '@taurushq/protect-sdk';

try {
  const client = ProtectClient.create({
    host: 'https://protect.example.com',
    credentials: Credentials.apiKey('your-api-key', 'your-hex-secret'),
  });

  const wallet = await client.wallets.get('wallet-123');
} catch (error) {
  if (error instanceof IntegrityError) {
    // Hash/signature verification failed - security issue
    // NEVER retry - investigate the cause
    console.error('Security error:', error.message);
  } else if (error instanceof NotFoundError) {
    // Resource not found
    console.error('Wallet not found');
  } else if (error instanceof RateLimitError) {
    // Rate limited - wait and retry
    const delay = error.suggestedRetryDelayMs();
    await new Promise(resolve => setTimeout(resolve, delay));
    // Retry the request
  } else if (error instanceof APIError) {
    if (error.isRetryable()) {
      // Server error or rate limit - can retry
      const delay = error.suggestedRetryDelayMs() ?? 5000;
      await new Promise(resolve => setTimeout(resolve, delay));
    } else {
      // Client error - check inputs
      console.error(`API error: ${error.statusCode} - ${error.message}`);
    }
  }
}
```

---

## Pagination

Every list method follows the cross-SDK contract (repo-root `CLAUDE.md` § "Pagination
(cross-SDK)", pinned by `scripts/resources/pagination-vectors.json` and
`list-request-vectors.json`). Two families, one rule for page sizes:

- **Page size**: `DEFAULT_PAGE_SIZE` (20) when unset or 0, at most `MAX_PAGE_SIZE` (100), and
  always sent. Above 100, negative or fractional → `ValidationError` before any request; a
  negative `offset` likewise. Exempt because they cannot page: price history (`limit` ≤ 365,
  default 20) and the transaction export (`limit` default 20, no maximum — it always exports
  from the first matching row, so it takes no offset and reports the server's `totalItems`).
- **A reply field the SDK cannot read raises `PaginationError`**: a count that is not a
  canonical decimal in [0, 2^53 - 1], a non-numeric reply offset, or a cursor that reports
  `hasNext` without a `currentPage`. It is never read as 0 or "no more pages".
- The pagination value is **never undefined** on success, the empty page `{}` included.

### Offset-Based Pagination

Wallets, addresses, transactions, users, groups, fee payers, actions and the whitelists — the
lists validatord offers no cursor for:

```typescript
interface Pagination {
  readonly limit: number;       // the page size sent
  readonly offset: number;      // the request's offset, never the reply's
  readonly totalItems: number;  // server total, minus rows this SDK withheld (never below 0)
  readonly nextOffset: number;  // pass back as `offset` while hasMore
  readonly hasMore: boolean;
}

interface PaginatedResult<T> {
  readonly items: T[];
  readonly pagination: Pagination;
}
```

`nextOffset` follows the endpoint's rule — the reply's `offset` for wallets and addresses,
`offset + rows` for transactions/fee payers/actions, `offset + min(rows, limit)` for users and
groups (a synthetic row can be appended), `offset + rows the server returned` for whitelisted
addresses, `offset + limit` for whitelisted contracts — and `hasMore` is
`nextOffset > offset && nextOffset < server total`. Rows withheld as unverifiable reduce
`totalItems` only; they never shift `nextOffset`.

**Example — walk every page:**

```typescript
let offset = 0;
for (;;) {
  const page = await client.wallets.list({ limit: 100, offset });
  page.items.forEach((w) => console.log(w.name));
  if (!page.pagination.hasMore) break;
  offset = page.pagination.nextOffset;
}
```

### Cursor-Based Pagination

Every other list — requests, changes, audit trails, business rules, balances, exchanges, fiat,
prices, reservations, webhooks, the v2 asset registry, earn rewards and Taurus-NETWORK. Options
extend `CursorPageOptions` (`pageSize`, `cursor`); lists that had them keep the low-level
`currentPage` / `pageRequest`, which cannot be combined with `cursor`:

```typescript
interface CursorPage {
  readonly pageSize: number;     // the page size sent
  readonly nextCursor: string;   // "" when there is no next page
  readonly hasMore: boolean;     // the reply cursor's hasNext (or the token's presence)
  readonly totalItems?: number;  // only where the server returns a total
}
```

Each result carries it as `pagination` next to its rows (`requests`, `changes`, `rules`,
`items`, ...). Passing `cursor` sends `currentPage=<cursor>` + `pageRequest=NEXT` + the page
size. Cursors are opaque base64 strings, possibly containing `+ / =`; pass them back verbatim.
Wallet tokens and governance-rules history use a token instead of a cursor object; the
contract is the same.

**Example — walk every page:**

```typescript
let cursor: string | undefined;
do {
  const page = await client.requests.list({ pageSize: 100, cursor });
  page.requests.forEach((r) => console.log(r.id, r.status));
  cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
} while (cursor);
```

---

## Mappers

The TypeScript SDK uses mapper functions in `src/mappers/` to convert between OpenAPI DTOs and domain models. Mappers are pure functions that handle:

- Property name conversion (camelCase from snake_case)
- Type coercion with safe helpers
- Nested object mapping
- Default values for missing fields

**Safe Conversion Helpers:**

```typescript
// From src/mappers/helpers.ts
safeString(value: unknown): string | undefined;
safeInt(value: unknown): number | undefined;
safeBool(value: unknown): boolean;
safeDate(value: unknown): Date | undefined;
safeMap<T>(value: T[] | undefined, mapper: (item: T) => R): R[];
```

**Example Mapper:**

```typescript
// src/mappers/wallet.ts
export function mapWallet(dto: TgvalidatordWallet): Wallet {
  return {
    id: safeString(dto.id) ?? '',
    name: safeString(dto.name) ?? '',
    status: mapWalletStatus(dto.status),
    type: mapWalletType(dto.type),
    blockchain: safeString(dto.blockchain),
    network: safeString(dto.network),
    // ... additional fields
  };
}
```

---

## TaurusNetwork Models

TaurusNetwork models are organized in `src/models/taurus-network/` with 72 models across 5 domains:

| File | Models | Purpose |
|------|--------|---------|
| `participant.ts` | 7 | Participant, ParticipantType, etc. |
| `pledge.ts` | 25 | Pledge, PledgeStatus, CollateralLimits, etc. |
| `lending.ts` | 14 | Loan, LoanStatus, LoanOffer, etc. |
| `settlement.ts` | 11 | Settlement, SettlementInstruction, etc. |
| `sharing.ts` | 15 | SharedAddressAsset, SharingAgreement, etc. |

**Import Pattern:**

```typescript
// Import specific models
import { Participant, Pledge, Loan } from '@taurushq/protect-sdk';

// Or import from taurus-network submodule
import { Settlement, SharedAddressAsset } from '@taurushq/protect-sdk';
```

---

## Related Documentation

### TypeScript SDK Specific
- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns

### Common (applies to all SDKs)
- [Key Concepts](../../docs/CONCEPTS.md) - Full domain model
- [Authentication](../../docs/AUTHENTICATION.md) - TPV1 protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Verification flows

### Other SDKs
- [Java SDK Documentation](../../taurus-protect-sdk-java/docs/) - Java SDK reference
- [Go SDK Documentation](../../taurus-protect-sdk-go/docs/) - Go SDK reference
- [Python SDK Documentation](../../taurus-protect-sdk-python/docs/) - Python SDK reference

## AuthorizationError.requiredRoles

A 403 is raised as `AuthorizationError`, which exposes `requiredRoles: readonly string[]`.

```typescript
try {
  await client.wallets.create(/* ... */);
} catch (error) {
  if (error instanceof AuthorizationError) {
    console.log(
      error.requiredRoles.length > 0
        ? `Permission denied; requires one of: ${error.requiredRoles.join(', ')}`
        : 'Permission denied'
    );
  }
}
```

The error carries the **roles that would satisfy the failed check**, so a caller can say
which role to request instead of only "forbidden". The list is empty when the denial was
not role-based (a disabled endpoint, a visibility restriction). One entry means that role
is required; several mean any one of them suffices.

> Surface the role list, never the server's message. The role names are safe to display;
> the surrounding text is server-controlled. A consumer that forwards errors into another
> system should render roles through an allowlist (the roles are lowercase alphanumeric).
