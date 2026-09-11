# Services Reference

This document provides comprehensive documentation for the Taurus-PROTECT TypeScript SDK services.

## Service Overview

The SDK provides 43 high-level services (38 on `ProtectClient` plus 5 on the `taurusNetwork` namespace), each with domain models and validation. Low-level OpenAPI access remains available for every feature alongside the high-level services.

### High-Level Services (26)

| Service | Access | Purpose |
|---------|--------|---------|
| [WalletService](#walletservice) | `client.wallets` | Create and manage blockchain wallets |
| [AddressService](#addressservice) | `client.addresses` | Create and manage addresses with signature verification |
| [RequestService](#requestservice) | `client.requests` | Transaction requests with approval workflow |
| [TransactionService](#transactionservice) | `client.transactions` | Query blockchain transactions |
| [BalanceService](#balanceservice) | `client.balances` | Asset and NFT balances |
| [CurrencyService](#currencyservice) | `client.currencies` | Currency metadata |
| [GovernanceRuleService](#governanceruleservice) | `client.governanceRules` | Governance rules with signature verification |
| [WhitelistedAddressService](#whitelistedaddressservice) | `client.whitelistedAddresses` | Whitelisted addresses with cryptographic verification |
| [WhitelistedAssetService](#whitelistedassetservice) | `client.whitelistedAssets` | Asset/contract whitelisting with verification |
| [AuditService](#auditservice) | `client.audits` | Audit log queries |
| [FeeService](#feeservice) | `client.fees` | Transaction fee information |
| [PriceService](#priceservice) | `client.prices` | Price data and conversion |
| [AirGapService](#airgapservice) | `client.airGap` | Air-gap signing operations |
| [UserService](#userservice) | `client.users` | User management |
| [GroupService](#groupservice) | `client.groups` | User group management |
| [VisibilityGroupService](#visibilitygroupservice) | `client.visibilityGroups` | Visibility group management |
| [ConfigService](#configservice) | `client.configService` | System configuration |
| [WebhookService](#webhookservice) | `client.webhooks` | Webhook management |
| [TagService](#tagservice) | `client.tags` | Tag management |
| [AssetService](#assetservice) | `client.assets` | Asset information |
| [ExchangeService](#exchangeservice) | `client.exchanges` | Exchange integration |
| [FeePayerService](#feepayerservice) | `client.feePayers` | Fee payer management |
| [HealthService](#healthservice) | `client.health` | API health checks |
| [JobService](#jobservice) | `client.jobs` | Background job management |
| [StatisticsService](#statisticsservice) | `client.statistics` | Platform statistics |
| [TokenMetadataService](#tokenmetadataservice) | `client.tokenMetadata` | Token metadata information |

### Low-Level API Access

Every feature also has a low-level OpenAPI-generated API on the client (`client.<name>Api`),
useful for endpoints or parameters the high-level service does not surface. There is no
longer any feature reachable *only* through a low-level API — all 43 services have
high-level getters.

### TaurusNetwork APIs

TaurusNetwork provides low-level API access for Taurus Network operations. Use the OpenAPI-generated APIs directly.

| API | Access | Purpose |
|-----|--------|---------|
| [TaurusNetworkParticipantApi](#taurusnetworkparticipantapi) | `client.taurusNetwork.participantApi` | Participant management |
| [TaurusNetworkPledgeApi](#taurusnetworkpledgeapi) | `client.taurusNetwork.pledgeApi` | Pledge lifecycle operations |
| [TaurusNetworkLendingApi](#taurusnetworklendingapi) | `client.taurusNetwork.lendingApi` | Lending offers and agreements |
| [TaurusNetworkSettlementApi](#taurusnetworksettlementapi) | `client.taurusNetwork.settlementApi` | Settlement operations |
| [TaurusNetworkSharedAddressAssetApi](#taurusnetworksharedaddressassetapi) | `client.taurusNetwork.sharedAddressAssetApi` | Address and asset sharing |

---

## WalletService

**Purpose:** Creates and manages blockchain wallets with balance tracking.

**Location:** `src/services/wallet-service.ts`

### Methods

#### get

Retrieves a wallet by ID.

```typescript
get(id: number): Promise<Wallet>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | number | The wallet ID |

**Returns:** `Wallet` - The wallet with balance information

**Throws:** `ValidationError` if id is invalid, `NotFoundError` if wallet not found

**Example:**
```typescript
const wallet = await client.wallets.get(123);
console.log(`Wallet: ${wallet.name}, Balance: ${wallet.balance}`);
```

#### list

Lists wallets with pagination.

```typescript
list(options?: ListWalletsOptions): Promise<ListWalletsResult>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| options.limit | number | Maximum results per page (optional) |
| options.offset | number | Pagination offset (optional) |
| options.name | string | Filter by name (optional) |
| options.query | string | Search query (optional) |

**Returns:** `ListWalletsResult` - Wallets list with pagination

**Example:**
```typescript
const result = await client.wallets.list({ limit: 50, offset: 0 });
console.log(`Total wallets: ${result.totalItems}`);
for (const wallet of result.wallets) {
  console.log(`${wallet.name}: ${wallet.blockchain}/${wallet.network}`);
}
```

#### create

Creates a new blockchain wallet.

```typescript
create(request: CreateWalletRequest): Promise<Wallet>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.blockchain | string | Blockchain identifier (e.g., "ETH", "BTC") |
| request.network | string | Network identifier (e.g., "mainnet", "testnet") |
| request.name | string | Human-readable wallet name |
| request.isOmnibus | boolean | Whether this is an omnibus wallet (optional) |
| request.comment | string | Optional comment |
| request.customerId | string | External customer ID (optional) |

**Returns:** `Wallet` - The created wallet

**Example:**
```typescript
const wallet = await client.wallets.create({
  blockchain: 'ETH',
  network: 'mainnet',
  name: 'My Treasury Wallet',
  isOmnibus: false,
  comment: 'Production wallet',
  customerId: 'CUST-001',
});
console.log(`Created wallet ID: ${wallet.id}`);
```

#### createAttribute

Adds a custom attribute to a wallet.

```typescript
createAttribute(walletId: number, key: string, value: string): Promise<void>
```

#### deleteAttribute

Removes an attribute from a wallet.

```typescript
deleteAttribute(walletId: number, attributeId: number): Promise<void>
```

#### getBalanceHistory

Gets historical balance data for a wallet.

```typescript
getBalanceHistory(walletId: number, intervalHours: number): Promise<BalanceHistoryPoint[]>
```

### Key Models

- `Wallet` - id, name, blockchain, network, balance, isOmnibus, customerId, attributes
- `ListWalletsResult` - wallets list with totalItems and pagination
- `BalanceHistoryPoint` - timestamp, balance values

---

## AddressService

**Purpose:** Manages blockchain addresses within wallets with signature verification.

**Location:** `src/services/address-service.ts`

### Methods

#### get

Retrieves an address with **mandatory signature verification**.

```typescript
get(id: number): Promise<Address>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | number | The address ID |

**Returns:** `Address` - The verified address

**Note:** This method performs cryptographic verification using the rules container cache. Throws `IntegrityError` if verification fails.

**Example:**
```typescript
const address = await client.addresses.get(456);
console.log(`Address: ${address.address}`);
console.log(`Balance: ${address.balance}`);
```

#### list

Lists addresses for a wallet with **signature verification**.

```typescript
list(walletId: number, limit?: number, offset?: number): Promise<Address[]>
```

#### listWithOptions

Lists addresses with advanced filtering options.

```typescript
listWithOptions(options?: ListAddressesOptions): Promise<ListAddressesResult>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| options.walletId | number | Parent wallet ID (optional) |
| options.limit | number | Maximum results (optional) |
| options.offset | number | Pagination offset (optional) |
| options.address | string | Filter by address (optional) |
| options.label | string | Filter by label (optional) |

#### create / createAddress

Creates a new address in a wallet, verifying the HSM signature on the reply.

```typescript
create(request: CreateAddressRequest): Promise<Address>
createAddress(walletId: number, label: string, comment?: string, customerId?: string): Promise<Address>
```

Both routes go through the same verification seam as `get` and `list`: **a non-empty
`address` string is never returned unverified.** This is the highest-value moment for
substitution, because the caller is about to publish or fund a fresh deposit address.

Address creation can be asynchronous, so the branch is on the address string rather than
on the server-controlled `status`:

| Reply | Result |
|---|---|
| `address` empty (`status` `created` / `creating`) | returned as-is — no destination yet, so nothing to verify and nothing to misuse. Re-read with `get(id)` once the status advances |
| `address` present, HSM `signature` present | signature verified against the HSMSLOT key from the rules container |
| `address` present, `signature` absent | `IntegrityError` — the server's address string is withheld rather than handed back unchecked |

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.walletId | number | Parent wallet ID |
| request.label | string | Address label |
| request.comment | string | Optional comment |
| request.customerId | string | External customer ID (optional) |

**Example:**
```typescript
const address = await client.addresses.create({
  walletId: 123,
  label: 'Customer Deposit Address',
  comment: 'For customer CUST-001',
  customerId: 'CUST-001',
});
console.log(`Created address: ${address.address}`);
```

#### createAttribute / deleteAttribute

Manages custom attributes on addresses.

```typescript
createAttribute(addressId: number, key: string, value: string): Promise<void>
deleteAttribute(addressId: number, attributeId: number): Promise<void>
```

#### getProofOfReserve

Gets proof of reserve for an address.

```typescript
getProofOfReserve(addressId: number, challenge: string): Promise<ProofOfReserve>
```

### Key Models

- `Address` - id, address, walletId, label, customerId, balance, attributes
- `ListAddressesResult` - addresses list with pagination
- `ProofOfReserve` - cryptographic proof data

---

## RequestService

**Purpose:** Creates, approves, and manages transaction requests with cryptographic signing.

**Location:** `src/services/request-service.ts`

### Methods

#### get

Retrieves a request with **hash verification**.

```typescript
get(id: number): Promise<Request>
```

**Verification:** Computes SHA-256 hash of metadata and compares with provided hash using constant-time comparison.

**Throws:** `IntegrityError` if hash verification fails

#### list

Lists requests with filtering.

```typescript
list(options?: ListRequestsOptions): Promise<RequestResult>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| options.from | Date | Start date (optional) |
| options.to | Date | End date (optional) |
| options.currencyId | string | Currency filter (optional) |
| options.statuses | RequestStatus[] | Status filter (optional) |
| options.cursor | RequestCursor | Pagination cursor (optional) |

#### listForApproval

Gets requests pending approval for the current user.

```typescript
listForApproval(cursor?: RequestCursor): Promise<RequestResult>
```

#### approveRequest / approveRequests

Signs and approves requests using a private key.

> **A request whose metadata hash has not been verified is refused.** The signature attests
> to those hashes, so each must be one verification cleared against its payload — the flag was
> set on every read path and read by nobody. The check runs after the hash-present check, so
> absent metadata still reports as absent.

```typescript
approveRequest(request: Request, privateKey: KeyObject): Promise<number>
approveRequests(requests: Request[], privateKey: KeyObject): Promise<number>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request(s) | Request / Request[] | Request(s) to approve |
| privateKey | KeyObject | User's ECDSA signing key |

**Returns:** Number of signatures performed

**Example:**
```typescript
import { createPrivateKey } from 'crypto';

const privateKey = createPrivateKey({ key: myPrivateKeyPem, format: 'pem' });
const request = await client.requests.get(requestId);

// Review metadata before signing
console.log(`Amount: ${request.metadata?.amount}`);
console.log(`Destination: ${request.metadata?.destinationAddress}`);

// Approve
const sigCount = await client.requests.approveRequest(request, privateKey);
console.log(`Signatures: ${sigCount}`);
```

#### rejectRequest / rejectRequests

Rejects requests with a comment.

```typescript
rejectRequest(requestId: number, comment: string): Promise<void>
rejectRequests(requestIds: number[], comment: string): Promise<void>
```

#### Create Transfer Requests

```typescript
// Internal transfer between addresses
createInternalTransferRequest(
  fromAddressId: number,
  toAddressId: number,
  amount: string,
  currencyId: string,
  comment?: string
): Promise<Request>

// Internal transfer from wallet (auto-selects source address)
createInternalTransferFromWalletRequest(
  fromWalletId: number,
  toAddressId: number,
  amount: string,
  currencyId: string,
  comment?: string
): Promise<Request>

// External transfer to whitelisted address
createExternalTransferRequest(
  fromAddressId: number,
  toWhitelistedAddressId: number,
  amount: string,
  currencyId: string,
  comment?: string
): Promise<Request>

// External transfer from wallet
createExternalTransferFromWalletRequest(
  fromWalletId: number,
  toWhitelistedAddressId: number,
  amount: string,
  currencyId: string,
  comment?: string
): Promise<Request>

// Incoming transfer from exchange
createIncomingRequest(
  fromExchangeId: number,
  toAddressId: number,
  amount: string,
  currencyId: string,
  comment?: string
): Promise<Request>

// Cancel pending transaction
createCancelRequest(
  addressId: number,
  nonce: string,
  comment?: string
): Promise<Request>
```

### Key Models

- `Request` - id, status, type, currency, currencyInfo, metadata, approvers, needsApprovalFrom, tags, memo, rule, createdAt, updatedAt
- `RequestMetadata` - hash (SHA-256 hex-encoded), payloadAsString (raw payload for hash computation). Note: `payload` field is intentionally omitted for security; use `JSON.parse(payloadAsString)` instead.
- `RequestResult` - requests list with pagination cursor
- `RequestStatus` - 'CREATED', 'PENDING', 'APPROVING', 'APPROVED', 'HSM_READY', 'HSM_SIGNED', 'BROADCASTING', 'BROADCASTED', 'MINED', 'CONFIRMED', 'REJECTED', 'CANCELED', 'PERMANENT_FAILURE', 'EXPIRED', and others (see `src/models/request.ts` for full list of 30+ statuses)

---

## TransactionService

**Purpose:** Retrieves and analyzes blockchain transactions.

**Location:** `src/services/transaction-service.ts`

### Methods

#### get

Retrieves a transaction by ID.

```typescript
get(id: number): Promise<Transaction>
```

#### getByHash

Retrieves a transaction by blockchain hash.

```typescript
getByHash(hash: string): Promise<Transaction>
```

#### list

Lists transactions with filtering.

```typescript
list(options?: ListTransactionsOptions): Promise<TransactionResult>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| options.from | Date | Start date (optional) |
| options.to | Date | End date (optional) |
| options.currencyId | string | Currency filter (optional) |
| options.direction | string | "incoming" or "outgoing" (optional) |
| options.limit | number | Maximum results (optional) |
| options.offset | number | Pagination offset (optional) |

#### listByRequest

Lists transactions for a specific request.

```typescript
listByRequest(requestId: number): Promise<Transaction[]>
```

#### listByAddress

Lists transactions for a specific address.

```typescript
listByAddress(address: string, limit?: number, offset?: number): Promise<Transaction[]>
```

**Example:**
```typescript
const transactions = await client.transactions.list({
  from: new Date('2024-01-01'),
  to: new Date('2024-12-31'),
  currencyId: 'ETH',
  direction: 'outgoing',
  limit: 100,
});

for (const tx of transactions.transactions) {
  console.log(`${tx.hash}: ${tx.value} ${tx.currency}`);
}
```

### Key Models

- `Transaction` - id, hash, status, currency, blockchain, network, sources, destinations, value, fee, blockNumber, direction

---

## BalanceService

**Purpose:** Retrieves asset and NFT collection balances with pagination.

**Location:** `src/services/balance-service.ts`

### Methods

#### list

Lists balances with pagination.

```typescript
list(options?: ListBalancesOptions): Promise<BalanceResult>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| options.currencyId | string | Filter by currency (optional) |
| options.pageSize | number | Page size (optional) |
| options.currentPage | string | Current page cursor (optional) |

**Example:**
```typescript
let result = await client.balances.list({ pageSize: 100 });
do {
  for (const balance of result.balances) {
    console.log(`${balance.asset}: ${balance.balance}`);
  }
  if (result.pagination?.hasNext) {
    result = await client.balances.list({
      pageSize: 100,
      currentPage: result.pagination.currentPage,
    });
  }
} while (result.pagination?.hasNext);
```

#### listNFTCollections

Lists NFT collection balances.

```typescript
listNFTCollections(options?: ListNFTCollectionBalancesOptions): Promise<NFTCollectionBalanceResult>
```

### Key Models

- `BalanceResult` - balances list with cursor pagination
- `AssetBalance` - asset info with available/pending balances
- `NFTCollectionBalance` - NFT collection balance information

---

## CurrencyService

**Purpose:** Manages and retrieves currency metadata.

**Location:** `src/services/currency-service.ts`

### Methods

```typescript
list(options?: ListCurrenciesOptions): Promise<Currency[]>
get(currencyId: string): Promise<Currency>
getByBlockchain(blockchain: string, network: string): Promise<Currency>
getBaseCurrency(): Promise<Currency>
```

**Example:**
```typescript
// Get all currencies
const currencies = await client.currencies.list({ showDisabled: false });

// Get specific currency
const eth = await client.currencies.get('ETH');
console.log(`${eth.name}: ${eth.decimals} decimals`);

// Get tenant's base currency
const baseCurrency = await client.currencies.getBaseCurrency();
```

### Key Models

- `Currency` - id, name, symbol, blockchain, network, decimals, logo

---

## GovernanceRuleService

**Purpose:** Manages governance rules with SuperAdmin signature verification.

**Location:** `src/services/governance-rule-service.ts`

### Methods

#### getRules

Gets current governance rules with **signature verification**.

```typescript
getRules(): Promise<GovernanceRules>
```

#### getRulesById

Gets governance rules by ID.

```typescript
getRulesById(id: number): Promise<GovernanceRules>
```

#### getRulesProposal

Gets pending rules proposal (SuperAdmin only).

```typescript
getRulesProposal(): Promise<GovernanceRules | null>
```

#### getRulesHistory

Gets historical governance rules.

```typescript
getRulesHistory(options?: GovernanceRulesHistoryOptions): Promise<GovernanceRulesHistoryResult>
```

#### getDecodedRulesContainer

Decodes the rules container.

```typescript
getDecodedRulesContainer(): Promise<DecodedRulesContainer>
```

#### verifyGovernanceRules

Manually verifies rules against SuperAdmin keys.

```typescript
verifyGovernanceRules(rules: GovernanceRules, minValidSignatures: number): Promise<GovernanceRules>
```

#### updateRulesProposal

Submits a typed rules container as a governance proposal (SuperAdmin only). The container
is encoded to the wire format internally; the server-controlled `enforcedRulesHash` and
`timestamp` fields are stripped. The endpoint returns no body — callers needing the
persisted proposal should call `getRulesProposal` (note: rules reads are cached
server-side, so an immediate read-back may be stale).

```typescript
updateRulesProposal(container: DecodedRulesContainer): Promise<void>
```

#### approveRulesProposal

Signs the pending proposal's rules container with a SuperAdmin private key (SHA-256 +
P-256 ECDSA, base64 raw r||s) and submits the approval. The signature binds the exact
pending content — review it first via `getRulesProposal` + `getDecodedRulesContainer`.

```typescript
approveRulesProposal(privateKey: KeyObject, comment: string): Promise<void>
```

#### rejectRulesProposal

Rejects the pending rules proposal with a comment (SuperAdmin only).

```typescript
rejectRulesProposal(comment: string): Promise<void>
```

**Example:**
```typescript
const rules = await client.governanceRules.getRules();
console.log(`Rules locked: ${rules.locked}`);

// Decode rules container
const decoded = await client.governanceRules.getDecodedRulesContainer();
console.log(`Groups: ${decoded.groups?.length}`);
```

### Key Models

- `GovernanceRules` - rulesContainer, rulesSignatures, locked, trails
- `DecodedRulesContainer` - lossless typed rules container (users, groups, transaction and
  whitelisting rules); round-trips through `rulesContainerToBase64` /
  `rulesContainerFromBase64`
- `RuleCell` - typed transaction-rule cell union covering every cell type
  (`FiatAmountAny`, `FiatAmountRange`, `SourceInternalWallet`, `StringEqualValue`, ...);
  `RawCell` preserves cells from newer schemas verbatim. Decode and encode a cell with
  `ruleCellFromBytes(columnType, bytes)` / `ruleCellToBytes(columnType, cell)`

Cross-SDK cell wire-format parity is pinned by the shared golden vectors at
`scripts/resources/governance-cell-vectors.json` (monorepo root), consumed by every SDK's
test suite.

---

## WhitelistedAddressService

**Purpose:** Manages whitelisted addresses with comprehensive cryptographic verification.

**Location:** `src/services/whitelisted-address-service.ts`

### Methods

#### get

Gets a whitelisted address by ID.

```typescript
get(id: number): Promise<WhitelistedAddress>
```

#### getWithVerification

Gets a whitelisted address with **full verification**.

```typescript
getWithVerification(id: number): Promise<WhitelistedAddress>
```

**Verification Steps:**
1. Metadata hash verification (SHA-256)
2. Rules container signature verification (SuperAdmin)
3. Hash coverage verification
4. Whitelist signature verification (governance thresholds)

#### getEnvelope

Gets the signed envelope with all verification details.

```typescript
getEnvelope(id: number): Promise<SignedWhitelistedAddressEnvelope>
```

#### list

Lists whitelisted addresses with filtering.

```typescript
list(options?: ListWhitelistedAddressesOptions): Promise<ListWhitelistedAddressesResult>
```

**Example:**
```typescript
// Get verified whitelisted address
const wlAddress = await client.whitelistedAddresses.getWithVerification(123);
console.log(`Address: ${wlAddress.address}`);
console.log(`Blockchain: ${wlAddress.blockchain}/${wlAddress.network}`);

// List all whitelisted addresses
const result = await client.whitelistedAddresses.list({
  blockchain: 'ETH',
  network: 'mainnet',
  limit: 50,
});
```

### Key Models

- `WhitelistedAddress` - blockchain, network, address, addressType, memo, label, linkedInternalAddresses, linkedWallets
- `SignedWhitelistedAddressEnvelope` - signedAddress, metadata, rulesContainer, rulesSignatures, approvers, trails

---

## WhitelistedAssetService

**Purpose:** Manages whitelisted assets/contracts with cryptographic verification.

**Location:** `src/services/whitelisted-asset-service.ts`

### Methods

#### get

Gets a whitelisted asset by ID.

```typescript
get(id: string): Promise<WhitelistedAsset>
```

#### getWithVerification

Gets a whitelisted asset with **full verification**.

```typescript
getWithVerification(id: string): Promise<WhitelistedAsset>
```

#### getEnvelope

Gets the signed envelope with verification details.

```typescript
getEnvelope(id: string): Promise<SignedWhitelistedAssetEnvelope>
```

#### list

Lists whitelisted assets with filtering. Verification is **lenient** here: an unverifiable
row is excluded and named on `excludedUnverified` rather than failing the call, but a page
where rows came back and none survived throws `IntegrityError`.

```typescript
list(options?: ListWhitelistedAssetsOptions): Promise<ListWhitelistedAssetsResult>
```

#### listForApproval

Lists whitelisted assets awaiting approval, verified as in `list`. Without this the only
reader of the for-approval endpoint was the unverified contract service, so the rows an
approver inspects were never checked against governance.

```typescript
listForApproval(options?: ListWhitelistedAssetsForApprovalOptions): Promise<ListWhitelistedAssetsResult>
```

#### approve

Signs and submits an approval, **all-or-nothing**. Each asset is re-read and verified, and
the hashes those rows carry are what gets signed; a row that is missing or fails verification
aborts the whole call and nothing is signed. The API takes one signature covering the whole
batch, so a partial approval would mean the caller believes they approved more than they did.

```typescript
approve(ids: number[], privateKey: KeyObject, comment: string): Promise<void>
```

### Key Models

- `WhitelistedAsset` - id, blockchain, network, contractAddress, symbol, name, decimals, kind
- `ListWhitelistedAssetsResult` - items, pagination, excludedUnverified

---

## HealthService

**Purpose:** Checks system health status.

**Location:** `src/services/health-service.ts`

### Methods

#### check

Performs a basic health check.

```typescript
check(): Promise<HealthStatus>
```

#### getGlobalStatus

Gets detailed health status with all components.

```typescript
getGlobalStatus(): Promise<GlobalHealthStatus>
```

**Example:**
```typescript
const health = await client.health.check();
console.log(`Status: ${health.status}`);

const globalHealth = await client.health.getGlobalStatus();
for (const [name, group] of Object.entries(globalHealth.groups || {})) {
  console.log(`${name}: ${group.status}`);
}
```

### Key Models

- `HealthStatus` - status, version
- `GlobalHealthStatus` - groups with component statuses

---

## UserService

**Purpose:** Retrieves and manages user information.

**Location:** `src/services/user-service.ts`

### Methods

```typescript
get(userId: string): Promise<User>
getCurrentUser(): Promise<User>
list(options?: ListUsersOptions): Promise<ListUsersResult>
```

**Example:**
```typescript
// Get current user
const me = await client.users.getCurrentUser();
console.log(`Logged in as: ${me.email}`);

// List all users
const result = await client.users.list({ limit: 100 });
for (const user of result.users) {
  console.log(`${user.firstName} ${user.lastName}: ${user.email}`);
}
```

### Key Models

- `User` - id, email, firstName, lastName, roles, attributes

---

## GroupService

**Purpose:** Manages user groups for approval workflows.

**Location:** `src/services/group-service.ts`

### Methods

```typescript
get(groupId: string): Promise<Group>
list(options?: ListGroupsOptions): Promise<ListGroupsResult>
```

### Key Models

- `Group` - id, name, members, threshold

---

## VisibilityGroupService

**Purpose:** Manages visibility groups for resource access control.

**Location:** `src/services/visibility-group-service.ts`

### Methods

```typescript
list(): Promise<VisibilityGroup[]>
getUsersByVisibilityGroup(visibilityGroupId: string): Promise<User[]>
```

### Key Models

- `VisibilityGroup` - id, name

---

## WebhookService

**Purpose:** Manages webhooks for receiving real-time event notifications.

**Location:** `src/services/webhook-service.ts`

### Methods

#### create

Creates a new webhook configuration.

```typescript
create(request: CreateWebhookRequest): Promise<Webhook>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.url | string | URL to receive webhook notifications (HTTPS) |
| request.type | string | Event type (e.g., "TRANSACTION", "REQUEST") |
| request.secret | string | Secret for signing webhook payloads |

**Example:**
```typescript
const webhook = await client.webhooks.create({
  url: 'https://example.com/webhook',
  type: 'TRANSACTION',
  secret: 'my-secret-key',
});
console.log(`Created webhook: ${webhook.id}`);
```

#### list

Lists webhooks with filtering.

```typescript
list(options?: ListWebhooksOptions): Promise<WebhookResult>
```

#### get

Gets a webhook by ID.

```typescript
get(webhookId: string): Promise<Webhook>
```

#### delete

Deletes a webhook.

```typescript
delete(webhookId: string): Promise<void>
```

### Key Models

- `Webhook` - id, url, type, status, createdAt

---

## WebhookCallService

**Purpose:** Queries webhook call history.

**Location:** `src/services/webhook-call-service.ts`

> **Note:** This service does not have a client getter on `ProtectClient`. Use the low-level API via `client.webhookCallsApi` or instantiate `WebhookCallService` directly.

### Methods

```typescript
list(options?: ListWebhookCallsOptions): Promise<WebhookCallResult>
get(webhookCallId: string): Promise<WebhookCall>
```

### Key Models

- `WebhookCall` - id, webhookId, status, requestBody, responseCode, timestamp

---

## AuditService

**Purpose:** Queries audit trail events.

**Location:** `src/services/audit-service.ts`

### Methods

```typescript
list(options?: ListAuditTrailsOptions): Promise<AuditTrailResult>
```

**Example:**
```typescript
const result = await client.audits.list({
  entity: 'REQUEST',
  action: 'APPROVE',
  from: new Date('2024-01-01'),
  to: new Date(),
  pageSize: 100,
});

for (const audit of result.audits) {
  console.log(`${audit.entity} ${audit.action} by ${audit.user?.email}`);
}
```

### Key Models

- `AuditTrail` - id, entity, action, user, creationDate, details

---

## TagService

**Purpose:** Manages tags for organizing resources.

**Location:** `src/services/tag-service.ts`

### Methods

```typescript
list(): Promise<Tag[]>
get(tagId: string): Promise<Tag>
create(request: CreateTagRequest): Promise<Tag>
delete(tagId: string): Promise<void>
```

**Example:**
```typescript
// Create a tag
const tag = await client.tags.create({ value: 'high-priority' });

// List all tags
const tags = await client.tags.list();
for (const t of tags) {
  console.log(`Tag: ${t.value}`);
}
```

### Key Models

- `Tag` - id, value

---

## StakingService

**Purpose:** Retrieves staking information across multiple proof-of-stake blockchains.

**Location:** `src/services/staking-service.ts`

### Methods

#### getADAStakePoolInfo

Retrieves Cardano stake pool information.

```typescript
getADAStakePoolInfo(network: string, stakePoolId: string): Promise<ADAStakePoolInfo>
```

**Example:**
```typescript
const poolInfo = await client.staking.getADAStakePoolInfo('mainnet', 'pool1abc123...');
console.log(`Pool pledge: ${poolInfo.pledge}`);
```

#### getETHValidatorsInfo

Retrieves Ethereum validator information.

```typescript
getETHValidatorsInfo(network: string, validatorIds: string[]): Promise<ETHValidatorInfo[]>
```

#### getFTMValidatorInfo

Retrieves Fantom validator information.

```typescript
getFTMValidatorInfo(network: string, validatorAddress: string): Promise<FTMValidatorInfo>
```

#### getICPNeuronInfo

Retrieves Internet Computer neuron information.

```typescript
getICPNeuronInfo(network: string, neuronId: string): Promise<ICPNeuronInfo>
```

#### getNEARValidatorInfo

Retrieves NEAR Protocol validator information.

```typescript
getNEARValidatorInfo(network: string, validatorAddress: string): Promise<NEARValidatorInfo>
```

#### getStakeAccounts

Lists stake accounts with pagination.

```typescript
getStakeAccounts(options: GetStakeAccountsOptions): Promise<StakeAccountResult>
```

#### getXTZStakingRewards

Retrieves Tezos staking rewards for an address.

```typescript
getXTZStakingRewards(network: string, addressId: string, from?: Date, to?: Date): Promise<XTZStakingRewards>
```

### Key Models

- `ADAStakePoolInfo` - pledge, margin, fixedCost, activeStake
- `ETHValidatorInfo` - publicKey, balance, status
- `FTMValidatorInfo` - stakedAmount, status
- `ICPNeuronInfo` - stake, votingPower, dissolveDelay
- `NEARValidatorInfo` - stake, fee
- `XTZStakingRewards` - totalRewards, cycles

---

## ContractWhitelistingService

**Purpose:** WRITE operations on whitelisted smart contract addresses (ERC20 tokens, NFTs, FA2 tokens).

**Location:** `src/services/contract-whitelisting-service.ts`

**Accessor:** `client.contractWhitelisting`

> **Reads live on `WhitelistedAssetService`.** A whitelisted contract and a whitelisted asset
> are one server entity (`/whitelists/contracts`); the `get`/`list`/`listForApproval` that used
> to sit here returned the envelope with no verification, and this SDK's mapper additionally
> parsed the signed payload and returned its contents as fact. Use
> `client.whitelistedAssets.get` / `.list` / `.listForApproval`, which run the six-step chain.

### Methods

#### create

Creates a new whitelisted contract.

```typescript
create(request: CreateWhitelistedContractRequest): Promise<string>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.blockchain | string | Blockchain identifier |
| request.network | string | Network (e.g., "mainnet") |
| request.contractAddress | string | Smart contract address |
| request.symbol | string | Token symbol (e.g., "USDC") |
| request.name | string | Human-readable name |
| request.decimals | number | Token decimals (0 for NFTs) |
| request.kind | string | Contract kind (e.g., "erc20", "erc721") |
| request.tokenId | string | Token ID for NFTs (optional) |

**Example:**
```typescript
const id = await client.contractWhitelisting.create({
  blockchain: 'ETH',
  network: 'mainnet',
  contractAddress: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
  symbol: 'USDC',
  name: 'USD Coin',
  decimals: 6,
  kind: 'erc20',
});
console.log(`Created whitelist entry: ${id}`);
```

#### approve / reject

```typescript
approve(ids: string[], signature: string, comment: string): Promise<void>
reject(ids: string[], comment: string): Promise<void>
```

> `approve` is **deprecated**: the signature is an opaque blob over hashes nothing verified,
> so the caller cannot know what they signed. Use `client.whitelistedAssets.approve`, which
> re-reads and verifies the rows first.

#### update

```typescript
update(id: string, request: UpdateWhitelistedContractRequest): Promise<void>
```

#### createAttribute / getAttribute / deleteAttribute

```typescript
createAttribute(contractId: string, request: CreateAttributeRequest): Promise<Attribute[]>
getAttribute(contractId: string, attributeId: string): Promise<Attribute>
deleteAttribute(contractId: string, attributeId: string): Promise<void>
```

### Key Models

- `WhitelistedContractAttribute` - id, key, value, contentType, type, subType, isFile

---

## BusinessRuleService

**Purpose:** Manages transaction approval rules.

**Location:** `src/services/business-rule-service.ts`

### Methods

```typescript
list(options?: ListBusinessRulesOptions): Promise<BusinessRuleResult>
get(id: string): Promise<BusinessRule>
```

### Key Models

- `BusinessRule` - rule definition with conditions and actions

---

## ChangeService

**Purpose:** Manages configuration changes and approval workflows.

**Location:** `src/services/change-service.ts`

### Methods

```typescript
get(id: number): Promise<Change>
list(options?: ListChangesOptions): Promise<ChangeResult>
listForApproval(options?: ListChangesOptions): Promise<ChangeResult>
approve(id: number, comment?: string): Promise<void>
approveMany(ids: number[], comment?: string): Promise<void>
reject(id: number, comment: string): Promise<void>
rejectMany(ids: number[], comment: string): Promise<void>
```

**Example:**
```typescript
// List changes pending approval
const result = await client.changes.listForApproval();
for (const change of result.changes) {
  console.log(`${change.entity} ${change.operation}: ${change.status}`);
}

// Approve a change
await client.changes.approve(changeId, 'Approved via SDK');
```

### Key Models

- `Change` - id, entity, operation, status, payload, trails

---

## PriceService

**Purpose:** Provides price data and currency conversion.

**Location:** `src/services/price-service.ts`

> **Prices are signature-verified.** `rate` and `decimals` feed amount conversion, so an
> unverified price is a wrong number a caller acts on. `list` verifies each price against the
> `PRICEUPDATER` keys in the SuperAdmin-verified rules container. Whether prices must be
> signed is the **container's** call: no `PRICEUPDATER` configured means this tenant does not
> sign prices and the price passes through; a `PRICEUPDATER` configured plus a price with no
> signatures is an `IntegrityError`.

### Methods

```typescript
list(): Promise<Price[]>
getHistory(base: string, quote: string, limit?: number): Promise<PriceHistoryPoint[]>
convert(currency: string, amount: string, targetCurrencyIds: string[]): Promise<ConversionResult[]>
```

**Example:**
```typescript
// Get all prices
const prices = await client.prices.list();

// Convert 1 ETH to USD
const conversions = await client.prices.convert('ETH', '1000000000000000000', ['USD']);
for (const result of conversions) {
  console.log(`${result.targetCurrency}: ${result.convertedAmount}`);
}
```

### Key Models

- `Price` - currency pair, rate, decimals, signatures
- `PriceSignature` - userId, signature
- `PriceHistoryPoint` - timestamp, price
- `ConversionResult` - target currency, converted amount

---

## FeeService

**Purpose:** Retrieves transaction fee information.

**Location:** `src/services/fee-service.ts`

### Methods

```typescript
getFees(currency: string): Promise<Fee[]>
getFeesV2(currency: string): Promise<FeeV2[]>
```

### Key Models

- `Fee` - currency, feeType, amount, unit

---

## ScoreService

**Purpose:** Manages and refreshes address risk/compliance scores.

**Location:** `src/services/score-service.ts`

### Methods

```typescript
refreshAddressScore(addressId: number, scoreProvider?: string): Promise<Score[]>
refreshWhitelistedAddressScore(whitelistedAddressId: number, scoreProvider?: string): Promise<Score[]>
```

**Example:**
```typescript
const scores = await client.scores.refreshAddressScore(addressId, 'chainalysis');
for (const score of scores) {
  console.log(`${score.provider}: ${score.score}`);
}
```

### Key Models

- `Score` - provider, score value, timestamp, details

---

## AirGapService

**Purpose:** Provides air-gap signing operations for offline transaction signing.

**Location:** `src/services/air-gap-service.ts`

### Methods

```typescript
getOutgoingAirGap(requestId: number): Promise<AirGapRequestData>
getOutgoingAirGapAddresses(requestId: number): Promise<AirGapAddress[]>
submitIncomingAirGap(payload: string): Promise<void>
```

---

## ReservationService

**Purpose:** Manages balance reservations for addresses.

**Location:** `src/services/reservation-service.ts`

### Methods

```typescript
list(options?: ListReservationsOptions): Promise<ReservationResult>
get(reservationId: string): Promise<Reservation>
getUtxo(reservationId: string): Promise<ReservationUtxo[]>
```

### Key Models

- `Reservation` - id, addressId, amount, status, expiresAt

---

## MultiFactorSignatureService

**Purpose:** Manages multi-factor signature operations for enhanced security.

**Location:** `src/services/multi-factor-signature-service.ts`

### Methods

#### get

Gets multi-factor signature info by ID.

```typescript
get(id: string): Promise<MultiFactorSignatureInfo>
```

#### create

Creates a multi-factor signature batch.

```typescript
create(request: CreateMultiFactorSignatureRequest): Promise<string>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.entityType | MultiFactorSignatureEntityType | 'REQUEST', 'WHITELISTED_ADDRESS', or 'WHITELISTED_CONTRACT' |
| request.entityIds | string[] | IDs of entities to sign |

#### approve

Approves a multi-factor signature.

```typescript
approve(request: ApproveMultiFactorSignatureRequest): Promise<void>
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request.id | string | Multi-factor signature ID |
| request.signature | string | Base64-encoded ECDSA signature |
| request.comment | string | Optional approval comment |

#### reject

Rejects a multi-factor signature.

```typescript
reject(request: RejectMultiFactorSignatureRequest): Promise<void>
```

**Example:**
```typescript
// Create a multi-factor signature for requests
const mfsId = await client.multiFactorSignature.create({
  entityType: 'REQUEST',
  entityIds: ['123', '456'],
});

// Get the signature info
const info = await client.multiFactorSignature.get(mfsId);

// Approve with signature
await client.multiFactorSignature.approve({
  id: mfsId,
  signature: 'base64EncodedSignature',
  comment: 'Approved via SDK',
});
```

---

## ConfigService

**Purpose:** Retrieves tenant configuration settings.

**Location:** `src/services/config-service.ts`

### Methods

```typescript
getTenantConfig(): Promise<TenantConfig>
```

**Example:**
```typescript
const config = await client.configService.getTenantConfig();
console.log(`Tenant ID: ${config.tenantId}`);
console.log(`Base currency: ${config.baseCurrency}`);
console.log(`MFA mandatory: ${config.mfaMandatory}`);
```

### Key Models

- `TenantConfig` - tenantId, baseCurrency, mfaMandatory, protectEngineVersion

---

## AssetService

**Purpose:** Retrieves asset information for addresses and wallets.

**Location:** `src/services/asset-service.ts`

> **`getAssetAddresses` verifies every address's HSM signature**, the same check
> `AddressService` runs, and **fails fast** on the first address that does not verify. It
> returns the same entity, so returning it unverified made `AddressService`'s mandatory
> verification avoidable.

### Methods

```typescript
getAssetAddresses(currencyId: string, options?: GetAssetAddressesOptions): Promise<AssetAddressResult>
getAssetWallets(currencyId: string, options?: GetAssetWalletsOptions): Promise<AssetWalletResult>
```

---

## ActionService

**Purpose:** Lists and retrieves actions in the system.

**Location:** `src/services/action-service.ts`

### Methods

```typescript
list(options?: ListActionsOptions): Promise<ActionResult>
get(actionId: string): Promise<Action>
```

---

## BlockchainService

**Purpose:** Retrieves blockchain information.

**Location:** `src/services/blockchain-service.ts`

### Methods

```typescript
list(): Promise<BlockchainInfo[]>
get(blockchain: string, network: string): Promise<BlockchainInfo>
```

**Example:**
```typescript
const blockchains = await client.blockchains.list();
for (const bc of blockchains) {
  console.log(`${bc.name} (${bc.symbol}): ${bc.network}`);
}
```

### Key Models

- `BlockchainInfo` - symbol, name, network, features

---

## ExchangeService

**Purpose:** Manages exchange integrations.

**Location:** `src/services/exchange-service.ts`

### Methods

```typescript
list(): Promise<Exchange[]>
get(exchangeId: string): Promise<Exchange>
getCounterparties(exchangeId: string): Promise<Counterparty[]>
getWithdrawalFee(exchangeId: string, options: GetWithdrawalFeeOptions): Promise<WithdrawalFee>
export(exchangeId: string, options: ExportExchangeOptions): Promise<string>
```

### Key Models

- `Exchange` - id, name, type, status

---

## FiatService

**Purpose:** Manages fiat currency operations with fiat providers.

**Location:** `src/services/fiat-service.ts`

### Methods

```typescript
getFiatProviders(): Promise<FiatProvider[]>
getFiatProviderAccount(providerId: string, accountId: string): Promise<FiatProviderAccount>
getFiatProviderAccounts(providerId: string): Promise<FiatProviderAccount[]>
getFiatProviderCounterpartyAccount(providerId: string, accountId: string): Promise<FiatProviderCounterpartyAccount>
getFiatProviderCounterpartyAccounts(providerId: string): Promise<FiatProviderCounterpartyAccount[]>
getFiatProviderOperation(providerId: string, operationId: string): Promise<FiatProviderOperation>
getFiatProviderOperations(providerId: string, options?: GetFiatProviderOperationsOptions): Promise<FiatProviderOperation[]>
```

### Key Models

- `FiatProvider` - id, name, status
- `FiatProviderAccount` - account details
- `FiatProviderOperation` - operation details

---

## FeePayerService

**Purpose:** Manages fee payer configurations.

**Location:** `src/services/fee-payer-service.ts`

### Methods

```typescript
list(blockchain?: string, network?: string): Promise<FeePayer[]>
get(feePayerId: string): Promise<FeePayer>
```

### Key Models

- `FeePayer` - id, blockchain, network, address, balance

---

## JobService

**Purpose:** Manages background jobs.

**Location:** `src/services/job-service.ts`

### Methods

```typescript
list(options?: ListJobsOptions): Promise<JobResult>
get(jobId: string): Promise<Job>
getStatus(jobId: string): Promise<JobStatus>
```

### Key Models

- `Job` - id, type, status, progress, createdAt

---

## StatisticsService

**Purpose:** Retrieves platform statistics.

**Location:** `src/services/statistics-service.ts`

### Methods

```typescript
getPortfolioStatistics(): Promise<PortfolioStatistics>
```

**Example:**
```typescript
const stats = await client.statistics.getPortfolioStatistics();
console.log(`Total balance: ${stats.totalBalance}`);
console.log(`Wallets: ${stats.walletsCount}`);
console.log(`Addresses: ${stats.addressesCount}`);
```

### Key Models

- `PortfolioStatistics` - totalBalance, totalBalanceBaseCurrency, walletsCount, addressesCount

---

## TokenMetadataService

**Purpose:** Retrieves token metadata for various token standards.

**Location:** `src/services/token-metadata-service.ts`

### Methods

#### getERCTokenMetadata (deprecated)

```typescript
getERCTokenMetadata(options: GetERCTokenMetadataOptions): Promise<TokenMetadata>
```

#### getEVMERCTokenMetadata

Preferred method for ERC token metadata on EVM chains.

```typescript
getEVMERCTokenMetadata(options: GetEVMERCTokenMetadataOptions): Promise<TokenMetadata>
```

**Example:**
```typescript
// Get ERC-20 token metadata
const usdc = await client.tokenMetadata.getEVMERCTokenMetadata({
  network: 'mainnet',
  contract: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
  blockchain: 'ETH',
});
console.log(`Token: ${usdc.name}, Decimals: ${usdc.decimals}`);

// Get ERC-721 NFT metadata with image data
const nft = await client.tokenMetadata.getEVMERCTokenMetadata({
  network: 'mainnet',
  contract: '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D',
  tokenId: '1234',
  withData: true,
  blockchain: 'ETH',
});
console.log(`NFT: ${nft.name}, URI: ${nft.uri}`);
```

#### getFATokenMetadata

Retrieves FA token metadata (Tezos FA1.2/FA2).

```typescript
getFATokenMetadata(options: GetFATokenMetadataOptions): Promise<TokenMetadata>
```

#### getCryptoPunkMetadata

Retrieves CryptoPunk metadata.

```typescript
getCryptoPunkMetadata(options: GetCryptoPunkMetadataOptions): Promise<CryptoPunkMetadata>
```

### Key Models

- `TokenMetadata` - name, symbol, decimals, description, uri, dataType, base64Data
- `CryptoPunkMetadata` - punkId, punkAttributes, image

---

## UserDeviceService

**Purpose:** Manages user device pairing for mobile app integration.

**Location:** `src/services/user-device-service.ts`

### Methods

#### createPairing

Creates a new device pairing request.

```typescript
createPairing(): Promise<UserDevicePairing>
```

#### startPairing

Starts the pairing process.

```typescript
startPairing(pairingId: string, options: StartPairingOptions): Promise<void>
```

#### getPairingStatus

Gets the status of a pairing request.

```typescript
getPairingStatus(pairingId: string, nonce: string): Promise<UserDevicePairingInfo>
```

#### approvePairing

Approves a device pairing request.

```typescript
approvePairing(pairingId: string, options: ApprovePairingOptions): Promise<void>
```

**Example:**
```typescript
// Step 1: Create pairing
const pairing = await client.userDevices.createPairing();
console.log(`Pairing ID: ${pairing.pairingId}`);

// Step 2: Start pairing
await client.userDevices.startPairing(pairing.pairingId, {
  nonce: '123456',
  publicKey: 'base64-encoded-public-key',
});

// Step 3: Approve pairing
await client.userDevices.approvePairing(pairing.pairingId, { nonce: '123456' });

// Get API key
const info = await client.userDevices.getPairingStatus(pairing.pairingId, '123456');
if (info.apiKey) {
  console.log(`Device paired! API Key: ${info.apiKey}`);
}
```

### Key Models

- `UserDevicePairing` - pairingId
- `UserDevicePairingInfo` - status, apiKey

---

## TaurusNetwork APIs

TaurusNetwork provides low-level API access through a namespace pattern. These are OpenAPI-generated APIs for direct control over API calls.

```typescript
// Access TaurusNetwork APIs
client.taurusNetwork.participantApi
client.taurusNetwork.pledgeApi
client.taurusNetwork.lendingApi
client.taurusNetwork.settlementApi
client.taurusNetwork.sharedAddressAssetApi
```

**Note:** TaurusNetwork currently provides only low-level API access. The examples below show how to use these APIs directly.

---

## TaurusNetworkParticipantApi

**Purpose:** Low-level API for Taurus Network participant management.

**Access:** `client.taurusNetwork.participantApi`

### Methods

#### getMyParticipant

Retrieves the current participant with settings.

```typescript
getMyParticipant(): Promise<TgvalidatordGetMyParticipantReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.participantApi.getMyParticipant();
const myParticipant = response.result;
console.log(`My participant: ${myParticipant?.participant?.name}`);
console.log(`Settings: ${JSON.stringify(myParticipant?.settings)}`);
```

#### getParticipant

Retrieves a participant by ID.

```typescript
getParticipant(params: { participantId: string; includeTotalPledgesValuation?: boolean }): Promise<TgvalidatordGetParticipantReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.participantApi.getParticipant({
  participantId: 'participant-id',
  includeTotalPledgesValuation: true,
});
const participant = response.result?.participant;
console.log(`Outgoing pledges: ${participant?.outgoingTotalPledgesValuationBaseCurrency}`);
```

#### getAllParticipants

Lists visible participants.

```typescript
getAllParticipants(params?: { includeTotalPledgesValuation?: boolean }): Promise<TgvalidatordGetAllParticipantsReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.participantApi.getAllParticipants();
const participants = response.result?.participants ?? [];
for (const p of participants) {
  console.log(`${p.name}: ${p.country}`);
}
```

### Key Response Types

- `TgvalidatordGetMyParticipantReply` - Contains result with participant and settings
- `TgvalidatordGetParticipantReply` - Contains result with participant details
- `TgvalidatordGetAllParticipantsReply` - Contains result with participants array

---

## TaurusNetworkPledgeApi

**Purpose:** Low-level API for Taurus Network pledge lifecycle operations.

**Access:** `client.taurusNetwork.pledgeApi`

### Methods

#### getPledge

Retrieves a pledge by ID.

```typescript
getPledge(params: { pledgeId: string }): Promise<TgvalidatordGetPledgeReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.pledgeApi.getPledge({ pledgeId: 'pledge-id' });
const pledge = response.result?.pledge;
console.log(`Pledge amount: ${pledge?.pledgedAmount}`);
```

#### getAllPledges

Lists pledges with filtering.

```typescript
getAllPledges(params?: {
  ownerParticipantId?: string;
  targetParticipantId?: string;
  sharedAddressIds?: string[];
  currencyId?: string;
  statuses?: string[];
  pageSize?: number;
}): Promise<TgvalidatordGetAllPledgesReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.pledgeApi.getAllPledges({
  statuses: ['ACTIVE'],
  pageSize: 50,
});
const pledges = response.result?.pledges ?? [];
for (const pledge of pledges) {
  console.log(`${pledge.id}: ${pledge.status}`);
}
```

#### createPledge

Creates a new pledge.

```typescript
createPledge(params: { body: TgvalidatordCreatePledgeRequest }): Promise<TgvalidatordCreatePledgeReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.pledgeApi.createPledge({
  body: {
    sharedAddressId: '123',
    currencyId: 'ETH',
    amount: '1000000000000000000',
    pledgeType: 'PLEDGEE_WITHDRAWALS_RIGHTS',
  },
});
console.log(`Created pledge: ${response.result?.pledgeAction?.pledgeId}`);
```

#### getAllPledgeActions

Lists pledge actions.

```typescript
getAllPledgeActions(params?: { pageSize?: number }): Promise<TgvalidatordGetAllPledgeActionsReply>
```

#### getAllPledgeActionsForApproval

Lists pledge actions pending approval.

```typescript
getAllPledgeActionsForApproval(params?: { pageSize?: number }): Promise<TgvalidatordGetAllPledgeActionsForApprovalReply>
```

### Key Response Types

- `TgvalidatordGetPledgeReply` - Contains result with pledge details
- `TgvalidatordGetAllPledgesReply` - Contains result with pledges array and pagination
- `TgvalidatordCreatePledgeReply` - Contains result with created pledge action

---

## TaurusNetworkLendingApi

**Purpose:** Low-level API for lending offers and agreements in the Taurus Network.

**Access:** `client.taurusNetwork.lendingApi`

### Methods

#### getAllLendingOffers

Lists lending offers.

```typescript
getAllLendingOffers(params?: { pageSize?: number }): Promise<TgvalidatordGetAllLendingOffersReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.lendingApi.getAllLendingOffers({ pageSize: 50 });
const offers = response.result?.lendingOffers ?? [];
for (const offer of offers) {
  console.log(`${offer.id}: ${offer.annualPercentageYieldMainUnit} APY`);
}
```

#### getLendingOffer

Gets a specific lending offer.

```typescript
getLendingOffer(params: { lendingOfferId: string }): Promise<TgvalidatordGetLendingOfferReply>
```

#### createLendingOffer

Creates a new lending offer.

```typescript
createLendingOffer(params: { body: TgvalidatordCreateLendingOfferRequest }): Promise<TgvalidatordCreateLendingOfferReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.lendingApi.createLendingOffer({
  body: {
    currencyId: 'ETH',
    amount: '10000000000000000000',
    annualPercentageYield: '500',
    duration: 'P30D',
  },
});
console.log(`Created offer: ${response.result?.lendingOffer?.id}`);
```

#### getAllLendingAgreements

Lists lending agreements.

```typescript
getAllLendingAgreements(params?: { pageSize?: number }): Promise<TgvalidatordGetAllLendingAgreementsReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.lendingApi.getAllLendingAgreements({ pageSize: 50 });
const agreements = response.result?.lendingAgreements ?? [];
console.log(`Found ${agreements.length} agreements`);
```

#### getAllLendingAgreementsForApproval

Lists lending agreements pending approval.

```typescript
getAllLendingAgreementsForApproval(params?: { pageSize?: number }): Promise<TgvalidatordGetAllLendingAgreementsForApprovalReply>
```

### Key Response Types

- `TgvalidatordGetAllLendingOffersReply` - Contains result with lendingOffers array
- `TgvalidatordGetAllLendingAgreementsReply` - Contains result with lendingAgreements array

---

## TaurusNetworkSettlementApi

**Purpose:** Low-level API for settlements in the Taurus Network.

**Access:** `client.taurusNetwork.settlementApi`

### Methods

#### getSettlement

Retrieves a settlement by ID.

```typescript
getSettlement(params: { settlementId: string }): Promise<TgvalidatordGetSettlementReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.settlementApi.getSettlement({ settlementId: 'settlement-id' });
const settlement = response.result?.settlement;
console.log(`Settlement status: ${settlement?.status}`);
```

#### getAllSettlements

Lists settlements with filtering.

```typescript
getAllSettlements(params?: {
  counterParticipantId?: string;
  statuses?: string[];
  pageSize?: number;
}): Promise<TgvalidatordGetAllSettlementsReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.settlementApi.getAllSettlements({
  statuses: ['PENDING'],
  pageSize: 50,
});
const settlements = response.result?.settlements ?? [];
for (const settlement of settlements) {
  console.log(`${settlement.id}: ${settlement.status}`);
}
```

#### getAllSettlementsForApproval

Lists settlements pending approval.

```typescript
getAllSettlementsForApproval(params?: { pageSize?: number }): Promise<TgvalidatordGetAllSettlementsForApprovalReply>
```

#### createSettlement

Creates a new settlement.

```typescript
createSettlement(params: { body: TgvalidatordCreateSettlementRequest }): Promise<TgvalidatordCreateSettlementReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.settlementApi.createSettlement({
  body: {
    targetParticipantId: 'participant-456',
    firstLegParticipantId: 'participant-123',
    firstLegAssets: [{
      sourceSharedAddressId: 'addr-1',
      destinationSharedAddressId: 'addr-2',
      currencyId: 'ETH',
      amount: '1000000000000000000',
    }],
    secondLegAssets: [{
      sourceSharedAddressId: 'addr-3',
      destinationSharedAddressId: 'addr-4',
      currencyId: 'USDC',
      amount: '2000000000',
    }],
  },
});
console.log(`Created settlement: ${response.result?.settlement?.id}`);
```

### Key Response Types

- `TgvalidatordGetSettlementReply` - Contains result with settlement details
- `TgvalidatordGetAllSettlementsReply` - Contains result with settlements array and pagination

---

## TaurusNetworkSharedAddressAssetApi

**Purpose:** Low-level API for address and asset sharing in the Taurus Network.

**Access:** `client.taurusNetwork.sharedAddressAssetApi`

### Methods

#### getAllSharedAddressAssets

Lists shared address assets with filtering.

```typescript
getAllSharedAddressAssets(params?: {
  participantId?: string;
  ownerParticipantId?: string;
  targetParticipantId?: string;
  blockchain?: string;
  network?: string;
  statuses?: string[];
  pageSize?: number;
}): Promise<TgvalidatordGetAllSharedAddressAssetsReply>
```

**Example:**
```typescript
const response = await client.taurusNetwork.sharedAddressAssetApi.getAllSharedAddressAssets({
  ownerParticipantId: myParticipantId,
  pageSize: 50,
});
const sharedAssets = response.result?.sharedAddressAssets ?? [];
for (const asset of sharedAssets) {
  console.log(`${asset.id}: ${asset.address} (${asset.blockchain})`);
}
```

#### getSharedAddressAsset

Gets a specific shared address asset.

```typescript
getSharedAddressAsset(params: { sharedAddressAssetId: string }): Promise<TgvalidatordGetSharedAddressAssetReply>
```

#### shareAddress

Shares an address with another participant.

```typescript
shareAddress(params: { body: TgvalidatordShareAddressRequest }): Promise<TgvalidatordShareAddressReply>
```

**Example:**
```typescript
await client.taurusNetwork.sharedAddressAssetApi.shareAddress({
  body: {
    addressId: 'addr-123',
    toParticipantId: 'participant-456',
    keyValueAttributes: [{ key: 'purpose', value: 'settlement' }],
  },
});
```

#### unshareAddress

Unshares an address.

```typescript
unshareAddress(params: { sharedAddressId: string }): Promise<TgvalidatordUnshareAddressReply>
```

#### shareWhitelistedAsset / unshareWhitelistedAsset

Shares or unshares a whitelisted asset.

```typescript
shareWhitelistedAsset(params: { body: TgvalidatordShareWhitelistedAssetRequest }): Promise<TgvalidatordShareWhitelistedAssetReply>
unshareWhitelistedAsset(params: { sharedAssetId: string }): Promise<TgvalidatordUnshareWhitelistedAssetReply>
```

### Key Response Types

- `TgvalidatordGetAllSharedAddressAssetsReply` - Contains result with sharedAddressAssets array and pagination
- `TgvalidatordGetSharedAddressAssetReply` - Contains result with shared address asset details

---

## Exception Handling

All services follow a consistent exception pattern:

```typescript
import { ProtectClient, APIError, ValidationError, NotFoundError, IntegrityError } from '@taurushq/protect-sdk';

try {
  const wallet = await client.wallets.get(walletId);
} catch (error) {
  if (error instanceof IntegrityError) {
    // Hash/signature verification failed - security issue
    console.error('Security error:', error.message);
  } else if (error instanceof NotFoundError) {
    // Resource not found
    console.error('Not found:', error.message);
  } else if (error instanceof ValidationError) {
    // Invalid input
    console.error('Validation error:', error.message);
  } else if (error instanceof APIError) {
    console.error('API Error:', error.message);
    console.error('HTTP Status:', error.statusCode);

    if (error.isRetryable()) {
      // Rate limit or transient error - can retry
      await sleep(error.suggestedRetryDelayMs());
    }
  }
}
```

### Exception Types

| Exception | When Thrown |
|-----------|-------------|
| `APIError` | General API errors (network, auth, server errors) |
| `ValidationError` | Invalid input parameters |
| `NotFoundError` | Resource not found |
| `IntegrityError` | Hash verification or signature verification failed |
| `AuthenticationError` | Authentication failed |

---

## Pagination Patterns

### Cursor-Based Pagination

Used by: BalanceService, BusinessRuleService, RequestService, TaurusNetwork services

```typescript
let result = await service.list({ pageSize: 100 });
const allItems = [...result.items];

while (result.pagination?.hasNext) {
  result = await service.list({
    pageSize: 100,
    currentPage: result.pagination.currentPage,
    pageRequest: 'NEXT',
  });
  allItems.push(...result.items);
}
```

### Offset-Based Pagination

Used by: WalletService, AddressService, TransactionService, UserService

```typescript
const limit = 100;
let offset = 0;
const allItems: Item[] = [];

let result;
do {
  result = await service.list({ limit, offset });
  allItems.push(...result.items);
  offset += limit;
} while (result.items.length === limit);
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
- [Concepts](CONCEPTS.md) - Domain models and entities

<!-- BEGIN GENERATED METHOD INDEX -->

## Complete Method Index

Generated from the typescript source by `scripts/api-surface/docs.py`; regenerate with
`./build.sh docs`. Every method below exists in the SDK, and `./build.sh docs --check`
fails if this list drifts or if the prose above documents a method that does not.

43 services, 209 public methods.

### ActionService

- `get(actionId: string): Promise<ActionEnvelope>` — Retrieves a specific action by its ID.
- `list(): Promise<ActionEnvelope[]>` — Lists all actions.
- `list(options: ListActionsOptions): Promise<ActionEnvelope[]>` — Lists actions with optional filters.
- `list(options?: ListActionsOptions): Promise<ActionEnvelope[]>`

### AddressService

- `create(request: CreateAddressRequest): Promise<Address>` — Creates a new address, with mandatory signature verification of the reply.
- `createAddress(walletId: number, label: string, comment?: string, customerId?: string): Promise<Address>` — Creates an address with explicit parameters.
- `createAttribute(addressId: number, key: string, value: string): Promise<void>` — Creates an attribute for an address.
- `deleteAttribute(addressId: number, attributeId: number): Promise<void>` — Deletes an attribute from an address.
- `get(addressId: number): Promise<Address>` — Gets an address by ID with mandatory signature verification.
- `getProofOfReserve(addressId: number, challenge?: string): Promise<TgvalidatordGetAddressProofOfReserveReply["result"]>` — Gets the proof of reserve for an address.
- `list(walletId: number, options?: Omit<ListAddressesOptions, "walletId">): Promise<{ items: Address[]; pagination: Pagination | undefined; }>` — Lists addresses for a wallet with mandatory signature verification.
- `listWithOptions(options?: ListAddressesOptions): Promise<{ items: Address[]; pagination: Pagination | undefined; }>` — Lists addresses with full filtering options.

### AirGapService

- `getOutgoingAirGap(options: GetOutgoingAirGapOptions): Promise<Blob>` — Exports HSM-ready requests for cold HSM signing.
- `getOutgoingAirGapAddresses(options: GetOutgoingAirGapAddressOptions): Promise<Blob>` — Exports addresses for cold HSM signing.
- `submitIncomingAirGap(options: SubmitIncomingAirGapOptions): Promise<void>` — Imports signed requests from the cold HSM.

### AssetService

- `getAssetAddresses(options: GetAssetAddressesOptions): Promise<Address[]>` — Retrieves addresses that hold a specific asset.
- `getAssetWallets(options: GetAssetWalletsOptions): Promise<Wallet[]>` — Retrieves wallets that hold a specific asset.

### AuditService

- `exportAuditTrails(options?: { externalUserId?: string; entities?: string[]; actions?: string[]; creationDateFrom?: Date; creationDateTo?: Date; format?: string; }): Promise<string>` — Export audit trails to a formatted string (CSV or JSON).
- `list(options?: ListAuditTrailsOptions): Promise<AuditTrail[]>` — Lists audit trails with optional filtering.

### BalanceService

- `list(options?: ListBalancesOptions): Promise<AssetBalance[]>` — Lists asset balances for the tenant.
- `listNFTCollections(options: ListNFTCollectionBalancesOptions): Promise<NFTCollectionBalance[]>` — Lists NFT collection balances for the tenant.

### BlockchainService

- `get(blockchain: string, network: string, includeBlockHeight?: boolean): Promise<Blockchain>` — Gets a blockchain by symbol and network.
- `list(options?: ListBlockchainsOptions): Promise<Blockchain[]>` — Lists all supported blockchains.

### BusinessRuleService

- `get(ruleId: string): Promise<BusinessRule>` — Gets a business rule by ID.
- `list(options?: ListBusinessRulesOptions): Promise<ListBusinessRulesResult>` — Lists business rules with optional filtering and pagination.
- `updateTransactionsEnabled(enabled: boolean): Promise<void>` — Enables or disables transaction processing for the tenant (the

### ChangeService

- `approve(id: string): Promise<void>` — Approves a change.
- `approveMany(ids: string[]): Promise<void>` — Approves multiple changes.
- `create(request: CreateChangeRequest): Promise<string>` — Creates a change request.
- `get(id: string): Promise<Change>` — Gets a change by ID.
- `list(options?: ListChangesOptions): Promise<ListChangesResult>` — Lists changes with optional filtering.
- `listForApproval(options?: ListChangesForApprovalOptions): Promise<ListChangesResult>` — Lists changes pending approval.
- `reject(id: string): Promise<void>` — Rejects a change.
- `rejectMany(ids: string[]): Promise<void>` — Rejects multiple changes.

### ConfigService

- `getTenantConfig(): Promise<TenantConfig>` — Retrieves the tenant configuration.

### ContractWhitelistingService

- `approve(ids: string[], signature: string, comment: string): Promise<void>` — Approves one or more whitelisted contract addresses.
- `create(request: CreateWhitelistedContractRequest): Promise<string>` — Creates a new whitelisted contract address.
- `createAttribute(contractId: string, key: string, value: string, options?: { contentType?: string; type?: string; subType?: string; }): Promise<void>` — Creates an attribute on a whitelisted contract.
- `deleteAttribute(contractId: string, attributeId: string): Promise<void>` — Deletes an attribute from a whitelisted contract.
- `getAttribute(contractId: string, attributeId: string): Promise<WhitelistedContractAttribute>` — Gets an attribute from a whitelisted contract.
- `reject(ids: string[], comment: string): Promise<void>` — Rejects a whitelisted contract.
- `update(id: string, request: UpdateWhitelistedContractRequest): Promise<void>` — Updates an existing whitelisted contract.

### CurrencyService

- `get(currencyId: string): Promise<Currency>` — Gets a currency by ID.
- `getBaseCurrency(): Promise<Currency>` — Gets the base currency configured for the tenant.
- `getByBlockchain(options: GetCurrencyByBlockchainOptions): Promise<Currency>` — Gets a currency by blockchain and network.
- `list(options?: ListCurrenciesOptions): Promise<Currency[]>` — Lists all currencies.

### ExchangeService

- `export(format?: string): Promise<string>` — Exports all exchange accounts to a specified format.
- `get(id: string): Promise<Exchange>` — Gets an exchange account by ID.
- `getCounterparties(): Promise<ExchangeCounterparty[]>` — Gets all exchange counterparties.
- `getWithdrawalFee(exchangeId: string, options?: GetWithdrawalFeeOptions): Promise<ExchangeWithdrawalFee | undefined>` — Gets the withdrawal fee for a transfer from an exchange.
- `list(options?: ListExchangesOptions): Promise<ListExchangesResult>` — Lists exchange accounts.

### FeePayerService

- `get(id: string): Promise<FeePayer>` — Gets a fee payer by ID.
- `list(options?: ListFeePayersOptions): Promise<FeePayer[]>` — Lists fee payers with optional filtering.

### FeeService

- `getFees(): Promise<Fee[]>` — Retrieves current network fees for all supported blockchains (v1 API).
- `getFeesV2(): Promise<FeeV2[]>` — Retrieves current native currency fees for all supported blockchains (v2 API).

### FiatService

- `getFiatProviderAccount(id: string): Promise<FiatProviderAccount>` — Retrieves a fiat provider account by ID.
- `getFiatProviderAccounts(options: ListFiatProviderAccountsOptions): Promise<FiatProviderAccountResult>` — Retrieves fiat provider accounts with optional filtering.
- `getFiatProviderCounterpartyAccount(id: string): Promise<FiatProviderCounterpartyAccount>` — Retrieves a fiat provider counterparty account by ID.
- `getFiatProviderCounterpartyAccounts(options: ListFiatProviderCounterpartyAccountsOptions): Promise<FiatProviderCounterpartyAccountResult>` — Retrieves fiat provider counterparty accounts with optional filtering.
- `getFiatProviderOperation(id: string): Promise<FiatProviderOperation>` — Retrieves a fiat provider operation by ID.
- `getFiatProviderOperations(options?: ListFiatProviderOperationsOptions): Promise<FiatProviderOperationResult>` — Retrieves fiat provider operations with optional filtering.
- `getFiatProviders(): Promise<FiatProvider[]>` — Retrieves all configured fiat providers.

### GovernanceRuleService

- `approveRulesProposal(privateKey: KeyObject, comment: string, expectedContainerHash: string): Promise<void>`
- `decodeProposalForReview(rules: GovernanceRules): DecodedRulesContainer`
- `getDecodedRulesContainer(): Promise<DecodedRulesContainer>` — Gets the decoded rules container from the current governance rules.
- `getPublicKeys(): Promise<SuperAdminPublicKey[]>` — Gets the SuperAdmin public keys the server has configured.
- `getRules(): Promise<GovernanceRules | undefined>` — Gets the currently enforced governance rules.
- `getRulesById(rulesId: string): Promise<GovernanceRules | undefined>` — Gets a governance ruleset by its ID.
- `getRulesHistory(options?: ListGovernanceRulesHistoryOptions): Promise<GovernanceRulesHistoryResult>` — Gets the history of governance rules.
- `getRulesProposal(): Promise<GovernanceRules | undefined>` — Gets the proposed governance rules.
- `proposalContainerHash(rules: GovernanceRules): string`
- `rejectRulesProposal(comment: string): Promise<void>` — Rejects the pending rules proposal with a comment (SuperAdmin only).
- `updateRulesProposal(container: DecodedRulesContainer): Promise<void>`
- `verifyGovernanceRules(rules: GovernanceRules): GovernanceRules` — Verifies that governance rules have enough valid SuperAdmin signatures.

### GroupService

- `get(groupId: string): Promise<Group>` — Gets a group by ID.
- `list(options?: ListGroupsOptions): Promise<PaginatedResult<Group>>` — Lists groups with pagination.

### HealthService

- `check(): Promise<HealthStatus>` — Checks the API health status.
- `getGlobalStatus(): Promise<HealthStatus>` — Checks the global component status.

### JobService

- `get(name: string): Promise<Job>` — Gets a job by name.
- `getStatus(name: string, id: string): Promise<JobStatus>` — Gets the status of a specific job execution.
- `list(): Promise<Job[]>` — Lists all jobs.

### LendingService

- `cancelLendingAgreement(lendingAgreementId: string): Promise<void>` — Cancels a lending agreement.
- `createLendingAgreement(request: CreateLendingAgreementRequest): Promise<string>` — Creates a new lending agreement.
- `createLendingAgreementAttachment(lendingAgreementId: string, request: CreateLendingAgreementAttachmentRequest): Promise<string>` — Adds an attachment to a lending agreement.
- `createLendingOffer(request: CreateLendingOfferRequest): Promise<string>` — Creates a new lending offer.
- `deleteLendingOffer(offerId: string): Promise<void>` — Deletes a specific lending offer.
- `deleteLendingOffers(): Promise<void>` — Deletes all lending offers for the current participant.
- `getLendingAgreement(lendingAgreementId: string): Promise<LendingAgreement>` — Gets a lending agreement by ID.
- `getLendingOffer(offerId: string): Promise<LendingOffer>` — Gets a lending offer by ID.
- `listLendingAgreementAttachments(lendingAgreementId: string): Promise<LendingAgreementAttachment[]>` — Lists attachments for a lending agreement.
- `listLendingAgreements(options?: ListLendingAgreementsOptions): Promise<{ agreements: LendingAgreement[]; pagination?: CursorPagination; }>` — Lists lending agreements.
- `listLendingAgreementsForApproval(options?: ListLendingAgreementsOptions): Promise<{ agreements: LendingAgreement[]; pagination?: CursorPagination; }>` — Lists lending agreements pending approval.
- `listLendingOffers(options?: ListLendingOffersOptions): Promise<{ offers: LendingOffer[]; pagination?: CursorPagination; }>` — Lists lending offers.
- `repayLendingAgreement(lendingAgreementId: string, request: RepayLendingAgreementRequest): Promise<void>` — Records repayment for a lending agreement.
- `updateLendingAgreement(lendingAgreementId: string, request: UpdateLendingAgreementRequest): Promise<void>` — Updates a lending agreement.

### MultiFactorSignatureService

- `approve(request: ApproveMultiFactorSignatureRequest): Promise<void>`
- `create(request: CreateMultiFactorSignatureRequest): Promise<string>` — Create a multi-factor signature batch.
- `get(id: string): Promise<MultiFactorSignatureInfo>`
- `reject(request: RejectMultiFactorSignatureRequest): Promise<void>` — Reject a multi-factor signature.

### ParticipantService

- `createParticipantAttribute(participantId: string, request: CreateParticipantAttributeRequest): Promise<void>` — Creates an attribute for a participant.
- `deleteParticipantAttribute(participantId: string, attributeId: string): Promise<void>` — Deletes an attribute for a participant.
- `get(participantId: string, options?: GetParticipantOptions): Promise<Participant>` — Gets a participant by ID.
- `getMyParticipant(): Promise<MyParticipant>` — Gets the current participant with settings.
- `list(options?: ListParticipantsOptions): Promise<Participant[]>` — Lists visible Taurus Network participants.

### PledgeService

- `addCollateral(pledgeId: string, request: AddPledgeCollateralRequest): Promise<AddCollateralResult>` — Adds collateral to an existing pledge.
- `approvePledgeActions(actions: PledgeAction[], privateKey: KeyObject, comment?: string): Promise<number>` — Approves one or more pledge actions, verifying every metadata hash it is about to
- `createPledge(request: CreatePledgeRequest): Promise<CreatePledgeResult>` — Creates a new pledge.
- `get(pledgeId: string): Promise<Pledge>` — Gets a pledge by ID.
- `initiateWithdrawPledge(pledgeId: string, request: InitiateWithdrawPledgeRequest): Promise<WithdrawPledgeResult>` — Initiates withdrawal from a pledge (pledgor operation).
- `list(options?: ListPledgesOptions): Promise<{ pledges: Pledge[]; pagination?: CursorPagination; }>` — Lists pledges with optional filtering.
- `listPledgeActions(options?: ListPledgeActionsOptions): Promise<{ actions: PledgeAction[]; pagination?: CursorPagination; }>` — Lists pledge actions with optional filtering.
- `listPledgeActionsForApproval(options?: ListPledgeActionsOptions): Promise<{ actions: PledgeAction[]; pagination?: CursorPagination; }>` — Lists pledge actions pending approval.
- `listPledgeWithdrawals(options?: ListPledgeWithdrawalsOptions): Promise<{ withdrawals: PledgeWithdrawal[]; pagination?: CursorPagination; }>` — Lists pledge withdrawals with optional filtering.
- `rejectPledge(pledgeId: string, request: RejectPledgeRequest): Promise<void>` — Rejects a pledge.
- `rejectPledgeActions(request: RejectPledgeActionsRequest): Promise<void>` — Rejects multiple pledge actions.
- `unpledge(pledgeId: string): Promise<UnpledgeResult>` — Unpledges all funds from a pledge.
- `updatePledge(pledgeId: string, request: UpdatePledgeRequest): Promise<void>` — Updates a pledge's default destination.
- `withdrawPledge(pledgeId: string, request: WithdrawPledgeRequest): Promise<WithdrawPledgeResult>` — Withdraws from a pledge (pledgee operation).

### PriceService

- `convert(options: ConvertOptions): Promise<ConversionResult[]>` — Converts an amount from one currency to target currencies.
- `getHistory(options: GetPriceHistoryOptions): Promise<PriceHistoryPoint[]>` — Gets price history for a currency pair.
- `list(): Promise<Price[]>` — Lists all current prices.

### RequestService

- `approveRequest(request: Request, privateKey: KeyObject, comment?: string): Promise<number>` — Approve a single request with ECDSA signature.
- `approveRequests(requests: Request[], privateKey: KeyObject, comment?: string): Promise<number>` — Approve multiple requests with ECDSA signature.
- `createCancelRequest(addressId: number, nonce: bigint | number): Promise<Request>` — Create a cancel request for a pending transaction.
- `createExternalTransferFromWalletRequest(options: CreateExternalTransferFromWalletOptions): Promise<Request>` — Create an external transfer request from a wallet to a whitelisted address.
- `createExternalTransferRequest(options: CreateExternalTransferOptions): Promise<Request>` — Create an external transfer request to a whitelisted address.
- `createIncomingRequest(options: CreateIncomingRequestOptions): Promise<Request>` — Create an incoming request from an exchange.
- `createInternalTransferFromWalletRequest(options: CreateInternalTransferFromWalletOptions): Promise<Request>` — Create an internal transfer request from a wallet.
- `createInternalTransferRequest(options: CreateInternalTransferOptions): Promise<Request>` — Create an internal transfer request between addresses.
- `get(requestId: number): Promise<Request>` — Get a request by ID with mandatory hash verification.
- `list(options?: ListRequestsOptions): Promise<ListRequestsResult>` — List requests with filtering and pagination.
- `listForApproval(options?: ListRequestsForApprovalOptions): Promise<ListRequestsResult>` — List requests pending approval.
- `rejectRequest(requestId: number, comment: string): Promise<void>` — Reject a single request.
- `rejectRequests(requestIds: number[], comment: string): Promise<void>` — Reject multiple requests.

### ReservationService

- `get(id: string): Promise<Reservation>` — Gets a reservation by ID.
- `getUtxo(id: string): Promise<ReservationUtxo>` — Gets the UTXO details for a reservation.
- `list(options?: ListReservationsOptions): Promise<Reservation[]>` — Lists all reservations with optional filtering.

### ScoreService

- `refreshAddressScore(addressId: number, scoreProvider: string): Promise<Score[]>` — Refreshes the compliance scores for an internal address.
- `refreshWhitelistedAddressScore(whitelistedAddressId: number, scoreProvider: string): Promise<Score[]>` — Refreshes the compliance scores for a whitelisted external address.

### SettlementService

- `cancel(settlementId: string): Promise<void>` — Cancels a settlement.
- `create(request: CreateSettlementRequest): Promise<string>` — Creates a new settlement.
- `get(settlementId: string): Promise<Settlement>` — Gets a settlement by ID.
- `list(options?: ListSettlementsOptions): Promise<{ settlements: Settlement[]; pagination?: CursorPagination; }>` — Lists settlements with optional filtering.
- `listForApproval(options?: ListSettlementsForApprovalOptions): Promise<{ settlements: Settlement[]; pagination?: CursorPagination; }>` — Lists settlements pending approval.
- `replace(settlementId: string, request: ReplaceSettlementRequest): Promise<void>` — Replaces (updates) a settlement.

### SharingService

- `listSharedAddresses(options?: ListSharedAddressesOptions): Promise<{ sharedAddresses: SharedAddress[]; pagination?: CursorPagination; }>` — Lists shared addresses with optional filtering.
- `listSharedAssets(options?: ListSharedAssetsOptions): Promise<{ sharedAssets: SharedAsset[]; pagination?: CursorPagination; }>` — Lists shared assets with optional filtering.
- `shareAddress(request: ShareAddressRequest): Promise<void>` — Shares an internal address with a Taurus Network participant.
- `shareWhitelistedAsset(request: ShareWhitelistedAssetRequest): Promise<void>` — Shares a whitelisted asset with a Taurus Network participant.
- `unshareAddress(sharedAddressId: string): Promise<void>` — Unshares an address from a Taurus Network participant.
- `unshareWhitelistedAsset(sharedAssetId: string): Promise<void>` — Unshares an asset from a Taurus Network participant.

### StakingService

- `getADAStakePoolInfo(network: string, stakePoolId: string): Promise<ADAStakePoolInfo>` — Retrieves information about a Cardano stake pool.
- `getETHValidatorsInfo(network: string, ids: string[]): Promise<ETHValidatorInfo[]>` — Retrieves information about Ethereum validators.
- `getFTMValidatorInfo(network: string, validatorAddress: string): Promise<FTMValidatorInfo>` — Retrieves information about a Fantom validator.
- `getICPNeuronInfo(network: string, neuronId: string): Promise<ICPNeuronInfo>` — Retrieves information about an Internet Computer Protocol neuron.
- `getNEARValidatorInfo(network: string, validatorAddress: string): Promise<NEARValidatorInfo>` — Retrieves information about a NEAR Protocol validator.
- `getStakeAccounts(options?: ListStakeAccountsOptions): Promise<StakeAccountResult>` — Retrieves stake accounts with optional filtering.
- `getXTZStakingRewards(options: GetXTZStakingRewardsOptions): Promise<XTZStakingRewards>` — Retrieves Tezos staking rewards for an address over a time period.

### StatisticsService

- `getPortfolioStatistics(): Promise<PortfolioStatistics>` — Retrieves aggregated portfolio statistics.

### TagService

- `create(request: CreateTagRequest): Promise<Tag>` — Creates a new tag.
- `delete(tagId: string): Promise<void>` — Deletes a tag.
- `get(tagId: string): Promise<Tag>` — Gets a tag by ID.
- `list(options?: ListTagsOptions): Promise<Tag[]>` — Lists tags.

### TokenMetadataService

- `getCryptoPunkMetadata(options: GetCryptoPunkMetadataOptions): Promise<CryptoPunkMetadata>` — Retrieves CryptoPunk metadata.
- `getERCTokenMetadata(options: GetERCTokenMetadataOptions): Promise<TokenMetadata>` — Retrieves ERC token metadata (ERC-20, ERC-721, ERC-1155).
- `getEVMERCTokenMetadata(options: GetEVMERCTokenMetadataOptions): Promise<TokenMetadata>` — Retrieves ERC token metadata for EVM-compatible chains.
- `getFATokenMetadata(options: GetFATokenMetadataOptions): Promise<TokenMetadata>` — Retrieves FA token metadata (Tezos FA1.2/FA2 standards).

### TransactionService

- `exportTransactions(options?: { fromDate?: Date; toDate?: Date; currency?: string; direction?: string; limit?: number; offset?: number; format?: string; blockchain?: string; network?: string; }): Promise<string>` — Export transactions to a formatted string (CSV or JSON).
- `get(transactionId: string): Promise<Transaction>` — Gets a transaction by ID.
- `getByHash(txHash: string): Promise<Transaction>` — Gets a transaction by its blockchain hash.
- `list(options?: ListTransactionsOptions): Promise<PaginatedResult<Transaction>>` — Lists transactions with pagination and optional filtering.
- `listByAddress(address: string, options?: { limit?: number; offset?: number; }): Promise<PaginatedResult<Transaction>>` — Lists transactions for a specific blockchain address.
- `listByRequest(requestId: string, options?: { limit?: number; offset?: number; }): Promise<PaginatedResult<Transaction>>` — Lists transactions associated with a specific request ID.

### UserDeviceService

- `approvePairing(pairingId: string, options: ApprovePairingOptions): Promise<void>` — Approves a device pairing request.
- `createPairing(): Promise<UserDevicePairing>` — Creates a new device pairing request.
- `getPairingStatus(pairingId: string, nonce: string): Promise<UserDevicePairingInfo>` — Gets the status of a device pairing request.
- `startPairing(pairingId: string, options: StartPairingOptions): Promise<void>` — Starts the device pairing process.

### UserService

- `get(userId: string): Promise<User>` — Gets a user by ID.
- `getCurrentUser(): Promise<User>` — Gets the current authenticated user.
- `list(options?: ListUsersOptions): Promise<PaginatedResult<User>>` — Lists users with pagination.

### VisibilityGroupService

- `getUsersByVisibilityGroup(visibilityGroupId: string): Promise<User[]>` — Gets users assigned to a specific visibility group.
- `list(): Promise<VisibilityGroup[]>` — Lists all visibility groups.

### WalletService

- `create(request: CreateWalletRequest): Promise<Wallet>` — Creates a new wallet.
- `createAttribute(walletId: number, key: string, value: string): Promise<void>` — Creates an attribute for a wallet.
- `deleteAttribute(walletId: number, attributeId: string): Promise<void>` — Deletes an attribute from a wallet.
- `get(walletId: number): Promise<Wallet>` — Gets a wallet by ID.
- `getBalanceHistory(walletId: number, intervalHours: number): Promise<BalanceHistoryPoint[]>` — Gets the balance history for a wallet.
- `getWalletTokens(walletId: number, limit?: number): Promise<AssetBalance[]>` — Gets the list of tokens (asset balances) for a wallet.
- `list(options?: ListWalletsOptions): Promise<PaginatedResult<Wallet>>` — Lists wallets with pagination and optional filtering.

### WebhookCallService

- `get(callId: string): Promise<WebhookCall>` — Gets a webhook call by ID.
- `list(options?: ListWebhookCallsOptions): Promise<WebhookCallResult>` — Lists webhook calls with optional filtering.

### WebhookService

- `create(request: CreateWebhookRequest): Promise<Webhook>` — Creates a new webhook.
- `delete(webhookId: string): Promise<void>` — Deletes a webhook.
- `get(webhookId: string): Promise<Webhook>` — Gets a webhook by ID.
- `list(options?: ListWebhooksOptions): Promise<Webhook[]>` — Lists webhooks.

### WhitelistedAddressService

- `approve(selection: WhitelistedAddressApproval, privateKey: KeyObject, comment: string): Promise<void>`
- `get(addressId: string): Promise<WhitelistedAddress>` — Gets a whitelisted address by ID with mandatory verification.
- `getEnvelope(addressId: string): Promise<Verified<SignedWhitelistedAddressEnvelope>>` — Gets the signed envelope for a whitelisted address, after verifying it.
- `getWithVerification(addressId: string): Promise<WhitelistedAddressVerificationResult>` — Gets a whitelisted address by ID with full verification.
- `list(options?: ListWhitelistedAddressesOptions): Promise<ListWhitelistedAddressesResult>` — Lists whitelisted addresses with mandatory verification.
- `listForApproval(options?: ListWhitelistedAddressesForApprovalOptions): Promise<ListWhitelistedAddressesResult>`
- `withVerification(api: AddressWhitelistingApi, config: WhitelistedAddressServiceConfig): WhitelistedAddressService` — Creates a WhitelistedAddressService with verification enabled.

### WhitelistedAssetService

- `approve(selection: WhitelistedAssetApproval, privateKey: KeyObject, comment: string): Promise<void>`
- `get(assetId: number): Promise<WhitelistedAsset>`
- `getEnvelope(assetId: number): Promise<Verified<SignedWhitelistedAssetEnvelope>>` — Gets the signed envelope for a whitelisted asset, after verifying it.
- `getWithVerification(assetId: number): Promise<WhitelistedAssetVerificationResult>` — Gets a whitelisted asset by ID with full verification.
- `list(options?: ListWhitelistedAssetsOptions): Promise<ListWhitelistedAssetsResult>` — Lists whitelisted assets.
- `listForApproval(options?: ListWhitelistedAssetsForApprovalOptions): Promise<ListWhitelistedAssetsResult>`
- `withVerification(api: ContractWhitelistingApi, config: WhitelistedAssetServiceConfig): WhitelistedAssetService` — Creates a WhitelistedAssetService with verification enabled.

<!-- END GENERATED METHOD INDEX -->
