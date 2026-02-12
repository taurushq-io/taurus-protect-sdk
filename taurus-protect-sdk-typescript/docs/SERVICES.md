# Services Reference

This document provides comprehensive documentation for the Taurus-PROTECT TypeScript SDK services.

## Service Overview

The SDK provides 26 high-level services accessible as getters on `ProtectClient`, with domain models and validation. Additional service classes exist in `src/services/` for direct instantiation, and low-level API access is available for all features.

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

The following features are available through low-level OpenAPI-generated APIs:

| API | Access | Purpose |
|-----|--------|---------|
| ChangesApi | `client.changesApi` | Configuration change approvals |
| BusinessRulesApi | `client.businessRulesApi` | Transaction approval rules |
| ReservationsApi | `client.reservationsApi` | Balance reservations |
| MultiFactorSignatureApi | `client.multiFactorSignatureApi` | Multi-factor signature operations |
| ContractWhitelistingApi | `client.contractWhitelistingApi` | Smart contract address whitelisting |
| StakingApi | `client.stakingApi` | Multi-chain staking information and validators |
| ActionsApi | `client.actionsApi` | Action management |
| BlockchainApi | `client.blockchainApi` | Blockchain information |
| FiatApi | `client.fiatApi` | Fiat currency operations |
| ScoresApi | `client.scoresApi` | Risk/compliance scoring |
| UserDeviceApi | `client.userDeviceApi` | User device management |
| WebhookCallsApi | `client.webhookCallsApi` | Webhook call history |

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

Creates a new address in a wallet.

```typescript
create(request: CreateAddressRequest): Promise<Address>
createAddress(walletId: number, label: string, comment?: string, customerId?: string): Promise<Address>
```

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

#### get

Gets current governance rules with **signature verification**.

```typescript
get(): Promise<GovernanceRules>
```

#### getById

Gets governance rules by ID.

```typescript
getById(id: number): Promise<GovernanceRules>
```

#### getProposal

Gets pending rules proposal (SuperAdmin only).

```typescript
getProposal(): Promise<GovernanceRules | null>
```

#### getHistory

Gets historical governance rules.

```typescript
getHistory(options?: GovernanceRulesHistoryOptions): Promise<GovernanceRulesHistoryResult>
```

#### getDecodedRulesContainer

Decodes the rules container.

```typescript
getDecodedRulesContainer(rules: GovernanceRules): Promise<DecodedRulesContainer>
```

#### verifyGovernanceRules

Manually verifies rules against SuperAdmin keys.

```typescript
verifyGovernanceRules(rules: GovernanceRules, minValidSignatures: number): Promise<GovernanceRules>
```

**Example:**
```typescript
const rules = await client.governanceRules.get();
console.log(`Rules locked: ${rules.locked}`);

// Decode rules container
const decoded = await client.governanceRules.getDecodedRulesContainer(rules);
console.log(`Groups: ${decoded.groups?.length}`);
```

### Key Models

- `GovernanceRules` - rulesContainer, rulesSignatures, locked, trails
- `DecodedRulesContainer` - groups, users, thresholds, addressWhitelistingRules

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

Lists whitelisted assets with filtering.

```typescript
list(options?: ListWhitelistedAssetsOptions): Promise<ListWhitelistedAssetsResult>
```

### Key Models

- `WhitelistedAsset` - id, blockchain, network, contractAddress, symbol, name, decimals, kind

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

**Purpose:** Manages whitelisted smart contract addresses (ERC20 tokens, NFTs, FA2 tokens).

**Location:** `src/services/contract-whitelisting-service.ts`

> **Note:** This service does not have a client getter on `ProtectClient`. Use the low-level API via `client.contractWhitelistingApi` or instantiate `ContractWhitelistingService` directly.

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

#### get / list / listForApproval

```typescript
get(id: string): Promise<SignedWhitelistedContractEnvelope>
list(options?: ListWhitelistedContractsOptions): Promise<WhitelistedContractResult>
listForApproval(options?: ListWhitelistedContractsForApprovalOptions): Promise<WhitelistedContractResult>
```

#### approve / reject

```typescript
approve(ids: string[], signature: string, comment?: string): Promise<void>
reject(ids: string[], comment: string): Promise<void>
```

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

- `SignedWhitelistedContractEnvelope` - id, blockchain, network, status, metadata, approvers, trails
- `WhitelistedContractResult` - contracts list with pagination

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

- `Price` - currency pair, price, timestamp
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
