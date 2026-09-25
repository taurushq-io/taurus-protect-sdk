# Services Reference

This document provides comprehensive documentation for all 44 services in the Taurus-PROTECT Java SDK.

## Service Overview

The SDK provides services organized into two categories: core services (39) and TaurusNetwork services (5).

### Core Services

| Service | Purpose |
|---------|---------|
| [WalletService](#walletservice) | Create and manage blockchain wallets |
| [AddressService](#addressservice) | Create and manage addresses with signature verification |
| [RequestService](#requestservice) | Transaction requests with approval workflow |
| [TransactionService](#transactionservice) | Query blockchain transactions |
| [BalanceService](#balanceservice) | Asset and NFT balances |
| [CurrencyService](#currencyservice) | Currency metadata |
| [GovernanceRuleService](#governanceruleservice) | Governance rules with signature verification |
| [WhitelistedAddressService](#whitelistedaddressservice) | Whitelisted addresses with cryptographic verification |
| [WhitelistedAssetService](#whitelistedassetservice) | Asset/contract whitelisting with verification |
| [AuditService](#auditservice) | Audit log queries |
| [ChangeService](#changeservice) | Configuration change approvals |
| [FeeService](#feeservice) | Transaction fee information |
| [PriceService](#priceservice) | Price data and conversion |
| [AirGapService](#airgapservice) | Air-gap signing operations |
| [StakingService](#stakingservice) | Multi-chain staking information and validators |
| [ContractWhitelistingService](#contractwhitelistingservice) | Smart contract address whitelisting |
| [BusinessRuleService](#businessruleservice) | Transaction approval rules |
| [ReservationService](#reservationservice) | Balance reservations |
| [MultiFactorSignatureService](#multifactorsignatureservice) | Multi-factor signature operations |
| [UserService](#userservice) | User management |
| [GroupService](#groupservice) | User group management |
| [VisibilityGroupService](#visibilitygroupservice) | Visibility group management |
| [ConfigService](#configservice) | System configuration |
| [WebhookService](#webhookservice) | Webhook management |
| [WebhookCallsService](#webhookcallsservice) | Webhook call history |
| [TagService](#tagservice) | Tag management |
| [AssetService](#assetservice) | Asset information |
| [ActionService](#actionservice) | Action management |
| [BlockchainService](#blockchainservice) | Blockchain information |
| [EarnService](#earnservice) | Earn rewards credited to addresses |
| [ExchangeService](#exchangeservice) | Exchange integration |
| [FiatService](#fiatservice) | Fiat currency operations |
| [FeePayerService](#feepayerservice) | Fee payer management |
| [HealthService](#healthservice) | API health checks |
| [JobService](#jobservice) | Background job management |
| [ScoreService](#scoreservice) | Risk/compliance scoring |
| [StatisticsService](#statisticsservice) | Platform statistics |
| [TokenMetadataService](#tokenmetadataservice) | Token metadata information |
| [UserDeviceService](#userdeviceservice) | User device management |

### TaurusNetwork Services

| Service | Access | Purpose |
|---------|--------|---------|
| [TaurusNetworkParticipantService](#taurusnetworkparticipantservice) | `client.taurusNetwork().participants()` | Participant management |
| [TaurusNetworkPledgeService](#taurusnetworkpledgeservice) | `client.taurusNetwork().pledges()` | Pledge lifecycle operations |
| [TaurusNetworkLendingService](#taurusnetworklendingservice) | `client.taurusNetwork().lending()` | Lending offers and agreements |
| [TaurusNetworkSettlementService](#taurusnetworksettlementservice) | `client.taurusNetwork().settlements()` | Settlement operations |
| [TaurusNetworkSharingService](#taurusnetworksharingservice) | `client.taurusNetwork().sharing()` | Address and asset sharing |

---

## Pagination

Every list method follows the cross-SDK pagination contract (`CONCEPTS.md` → "Pagination"):

- **Page size** — always sent. `0`/`null` is the default, `Pagination.DEFAULT_PAGE_SIZE` (20);
  a negative size, one above `Pagination.MAX_PAGE_SIZE` (100), or a negative offset throws
  `IllegalArgumentException` naming the option before any request.
- **Offset lists** (wallets, addresses, transactions, users, groups, fee payers, actions, the
  whitelists) take `int limit, long offset` and return a result with `getPagination()`
  (`OffsetPagination`): continue from `getNextOffset()` while `hasMore()`.
- **Cursor lists** take `Integer pageSize, String cursor` — `cursor` is a previous page's
  `getPage().getNextCursor()`, `null` for the first page — and return a result with `getPage()`
  (`CursorPage`). The `ApiRequestCursor` overloads are the low-level form (`PREVIOUS`/`LAST`);
  `null` means the first page with the default size.
- **Exempt:** `PriceService.getPriceHistory` (limit ≤ 365) and
  `TransactionService.exportTransactions` (limit only; the server ignores an export offset).

---

## WalletService

**Purpose:** Creates and manages blockchain wallets with balance tracking.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/WalletService.java`

### Methods

#### createWallet

Creates a new blockchain wallet.

```java
Wallet createWallet(String blockchain, String network, String walletName, boolean isOmnibus)
Wallet createWallet(String blockchain, String network, String walletName, boolean isOmnibus, String comment)
Wallet createWallet(String blockchain, String network, String walletName, boolean isOmnibus, String comment, String customerId)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| blockchain | String | Blockchain identifier (e.g., "ETH", "BTC") |
| network | String | Network identifier (e.g., "mainnet", "testnet") |
| walletName | String | Human-readable wallet name |
| isOmnibus | boolean | Whether this is an omnibus wallet |
| comment | String | Optional comment |
| customerId | String | Optional external customer ID |

**Returns:** `Wallet` - The created wallet

**Example:**
```java
Wallet wallet = client.getWalletService().createWallet(
    "ETH", "mainnet", "My Treasury Wallet", false, "Production wallet", "CUST-001"
);
System.out.println("Created wallet ID: " + wallet.getId());
```

#### getWallet

Retrieves a wallet by ID.

```java
Wallet getWallet(long walletId) throws ApiException
```

**Returns:** `Wallet` with balance information

#### getWallets

Lists a page of wallets (offset list; see [Pagination](#pagination)).

```java
WalletResult getWallets(int limit, long offset) throws ApiException
WalletResult getWallets(int limit, long offset, Boolean excludeDisabled) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| limit | int | Page size, 0 for the default (20), at most 100 |
| offset | long | 0 for the first page, then `getPagination().getNextOffset()` |
| excludeDisabled | Boolean | `true` hides every disabled wallet; by default only currency-disabled ones are hidden |

#### getWalletsByName

Searches wallets by name.

```java
WalletResult getWalletsByName(String name, int limit, long offset) throws ApiException
```

#### createWalletAttribute

Adds a custom attribute to a wallet.

```java
void createWalletAttribute(long walletId, String key, String value) throws ApiException
```

#### getWalletBalanceHistory

Gets historical balance data.

```java
List<BalanceHistoryPoint> getWalletBalanceHistory(long walletId, int intervalHours) throws ApiException
```

#### getWalletTokens

Lists a page of a wallet's token balances (token-paged; the page carries the server total).

```java
WalletTokensResult getWalletTokens(long walletId, Integer pageSize) throws ApiException
WalletTokensResult getWalletTokens(long walletId, Integer pageSize, String cursor) throws ApiException
```

### Key Models

- `WalletResult` - `getWallets()` + `getPagination()` (`OffsetPagination`)
- `WalletTokensResult` - `getBalances()` + `getPage()` (`CursorPage`)
- `Wallet` - id, name, blockchain, network, balance, isOmnibus, customerId, attributes
- `BalanceHistoryPoint` - timestamp, balance values
- `AssetBalance` - asset info with balance

---

## AddressService

**Purpose:** Manages blockchain addresses within wallets with signature verification.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/AddressService.java`

### Methods

#### createAddress

Creates a new address in a wallet, and **verifies the HSM signature on the address it returns**.

```java
Address createAddress(long walletId, String label, String comment, String customerId) throws ApiException
```

> **The rule is: never return a non-empty address string that has not been verified.** The
> create reply carries the same `signature` field the read paths verify, so a signed address is
> verified here exactly as `getAddress` does and a failure raises `IntegrityException`.
> Asynchronous creation is real, though — `status` is one of `created`, `creating`, `signed`,
> `observed`, `confirmed` — so when the reply carries an address with **no** signature the SDK
> refuses to hand back the server's address string and returns the id and `status` instead.
> Re-read through `getAddress` once the status advances. The branch is on the address STRING,
> not on `status`, because `status` is server-controlled.
>
> `getAddress`, `getAddresses`, `createAddress` and `AssetService.getAssetAddresses` all route
> through one private seam. This finding existed because `getAssetAddresses` was fixed and
> `createAddress` was missed, so a seam rather than a third copy is the point.

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| walletId | long | Parent wallet ID |
| label | String | Address label |
| comment | String | Optional comment |
| customerId | String | External customer ID |

**Returns:** `Address` - The created address

#### getAddress

Retrieves an address with **mandatory signature verification**.

```java
Address getAddress(long id) throws ApiException
```

**Note:** This method performs cryptographic verification using the rules container cache.

#### getAddresses

Lists a page of addresses with **signature verification** (offset list).

```java
AddressResult getAddresses(long walletId, int limit, long offset) throws ApiException
AddressResult getAddresses(Long walletId, int limit, long offset, Boolean excludeDisabled) throws ApiException
```

`walletId` null lists every wallet's addresses; `excludeDisabled` true hides disabled addresses
(sent as `includeDisabledAddresses=exclude`).

#### createAddressAttribute

Adds a custom attribute.

```java
void createAddressAttribute(long addressId, String key, String value) throws ApiException
```

#### deleteAddressAttribute

Removes an attribute.

```java
void deleteAddressAttribute(long addressId, long attributeId) throws ApiException
```

#### getAddressProofOfReserve

Gets proof of reserve for an address.

```java
TgvalidatordProofOfReserve getAddressProofOfReserve(long addressId, String challenge) throws ApiException
```

### Key Models

- `Address` - id, address, walletId, label, customerId, balance, attributes

---

## RequestService

**Purpose:** Creates, approves, and manages transaction requests with cryptographic signing.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/RequestService.java`

### Methods

#### Create Requests

```java
// Internal transfer between addresses
Request createInternalTransferRequest(long fromAddressId, long toAddressId, BigInteger amount)

// Internal transfer from wallet (auto-selects source address)
Request createInternalTransferFromWalletRequest(long fromWalletId, long toAddressId, BigInteger amount)

// External transfer to whitelisted address
Request createExternalTransferRequest(long fromAddressId, long toWhitelistedAddressId, BigInteger amount)

// External transfer from wallet
Request createExternalTransferFromWalletRequest(long fromWalletId, long toWhitelistedAddressId, BigInteger amount)

// Incoming transfer from exchange
Request createIncomingRequest(long fromExchangeId, long toAddressId, BigInteger amount)

// Cancel pending transaction
Request createCancelRequest(long addressId, long nonce)
```

#### getRequest

Retrieves a request with **hash verification**.

```java
Request getRequest(long id) throws ApiException
```

**Verification:** Computes SHA-256 hash of metadata and compares with provided hash.

#### getRequests

Lists a page of requests with filtering (cursor list; see [Pagination](#pagination)).

```java
RequestResult getRequests(OffsetDateTime from, OffsetDateTime to, String currencyId,
                          List<RequestStatus> statuses, Integer pageSize, String cursor) throws ApiException
RequestResult getRequests(OffsetDateTime from, OffsetDateTime to, String currencyId,
                          List<RequestStatus> statuses, ApiRequestCursor cursor) throws ApiException
```

#### getRequestsForApproval

Gets a page of requests pending approval for the current user.

```java
RequestResult getRequestsForApproval(Integer pageSize, String cursor) throws ApiException
RequestResult getRequestsForApproval(ApiRequestCursor cursor) throws ApiException
```

#### approveRequest / approveRequests

Signs and approves requests using a private key.

> **A request whose metadata hash has not been verified is refused** with
> `IntegrityException`. The signature attests to those hashes, so each must be one
> verification cleared against its payload. `RequestMetadata` has no `setHashVerified` — the
> flag can only be set by `verifyAndMaterialise()`, which performs the hash check in the same
> operation, so the payload accessors cannot be unlocked without the check having run. The
> refusal is ordered after `checkNotNull(privateKey)` so argument errors stay argument errors.

```java
int approveRequest(Request request, PrivateKey privateKey) throws ApiException
int approveRequests(List<Request> requests, PrivateKey privateKey) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request(s) | Request / List<Request> | Request(s) to approve (metadata must be hash-VERIFIED) |
| privateKey | PrivateKey | User's signing key |

**Returns:** Number of signatures performed

**Example:**
```java
PrivateKey myKey = CryptoTPV1.decodePrivateKey(myPrivateKeyPem);
Request request = client.getRequestService().getRequest(requestId);

// Review metadata before signing
System.out.println("Amount: " + request.getMetadata().getAmount());
System.out.println("Destination: " + request.getMetadata().getDestinationAddress());

// Approve
int sigs = client.getRequestService().approveRequest(request, myKey);
System.out.println("Signatures: " + sigs);
```

#### rejectRequest / rejectRequests

Rejects requests with a comment.

```java
void rejectRequest(long requestId, String comment) throws ApiException
void rejectRequests(List<Long> requestIds, String comment) throws ApiException
```

### Key Models

- `Request` - id, status, currency, metadata, signedRequests, approvers, trails
- `RequestMetadata` - hash, payloadAsString, payloadAsJson; plus extraction methods: getRequestId(), getSourceAddress(), getDestinationAddress(), getAmount(), getCurrency(), getRulesKey()
- `RequestResult` - requests list with pagination cursor
- `RequestStatus` - Enum: CREATED, PENDING, APPROVING, APPROVED, REJECTED, HSM_READY, HSM_SIGNED, BROADCASTED, MINED, CONFIRMED, PARTIALLY_CONFIRMED, PERMANENT_FAILURE, CANCELED, EXPIRED, etc.

---

## TransactionService

**Purpose:** Retrieves and analyzes blockchain transactions.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/TransactionService.java`

### Methods

#### getTransactionById

```java
Transaction getTransactionById(long id) throws ApiException
```

#### getTransactionByHash

```java
Transaction getTransactionByHash(String hash) throws ApiException
```

#### getTransactions

Lists a page of transactions with filtering (offset list).

```java
TransactionResult getTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                  String direction, int limit, long offset) throws ApiException
TransactionResult getTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                  String direction, String blockchain, String network,
                                  int limit, long offset) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| from | OffsetDateTime | Start date (optional) |
| to | OffsetDateTime | End date (optional) |
| currency | String | Currency filter (optional) |
| direction | String | "incoming" or "outgoing" (optional) |
| limit | int | Page size, 0 for the default (20), at most 100 |
| offset | long | 0 for the first page, then `getPagination().getNextOffset()` |

#### getTransactionsByAddress

```java
TransactionResult getTransactionsByAddress(String address, int limit, long offset) throws ApiException
```

#### exportTransactions

Exports transactions. The export cannot page — the server ignores any offset — so it takes a limit
only (default 20, no SDK maximum); the result carries the text and the server's total, so a
truncated export can be detected. No format is sent unless given (the server default is JSON).

```java
TransactionExportResult exportTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                           String direction, int limit) throws ApiException
TransactionExportResult exportTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                           String direction, String blockchain, String network,
                                           int limit) throws ApiException
TransactionExportResult exportTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                           String direction, String blockchain, String network,
                                           String format, int limit) throws ApiException
```

### Key Models

- `Transaction` - id, hash, currency, blockchain, network, sources, destinations (List<AddressInfo>), amount (BigInteger), fee, block, direction

---

## BalanceService

**Purpose:** Retrieves asset and NFT collection balances with pagination.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/BalanceService.java`

### Methods

#### getBalances

```java
BalanceResult getBalances(String currency, Integer pageSize, String cursor) throws ApiException
BalanceResult getBalances(ApiRequestCursor cursor) throws ApiException
BalanceResult getBalances(String currency, ApiRequestCursor cursor) throws ApiException
```

**Example:**
```java
String cursor = null;
BalanceResult result;
do {
    result = client.getBalanceService().getBalances(null, 100, cursor);
    for (AssetBalance balance : result.getBalances()) {
        System.out.println(balance.getAsset() + ": " + balance.getBalance());
    }
    cursor = result.getPage().getNextCursor();
} while (result.getPage().hasMore());
```

#### getNFTCollectionBalances

```java
NFTCollectionBalanceResult getNFTCollectionBalances(String blockchain, String network,
                                                     Integer pageSize, String cursor) throws ApiException
NFTCollectionBalanceResult getNFTCollectionBalances(String blockchain, String network,
                                                     ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `BalanceResult` - balances list + `getPage()` (`CursorPage`, with the server total)
- `AssetBalance` - asset info with available/pending balances
- `NFTCollectionBalance` - NFT collection balances

---

## CurrencyService

**Purpose:** Manages and retrieves currency metadata.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/CurrencyService.java`

### Methods

```java
List<Currency> getCurrencies() throws ApiException
List<Currency> getCurrencies(boolean showDisabled, boolean includeLogo) throws ApiException
Currency getCurrency(String currencyId) throws ApiException
Currency getCurrencyByBlockchain(String blockchain, String network) throws ApiException
Currency getBaseCurrency() throws ApiException
```

### Key Models

- `Currency` - id, name, symbol, blockchain, network, decimals, logo

---

## ScoreService

**Purpose:** Manages and refreshes address risk/compliance scores.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/ScoreService.java`

### Methods

```java
List<Score> refreshAddressScore(long addressId, String scoreProvider) throws ApiException
List<Score> refreshWhitelistedAddressScore(long addressId, String scoreProvider) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| addressId | long | Address ID |
| scoreProvider | String | Provider name (e.g., "chainalysis", "elliptic") |

### Key Models

- `Score` - provider, score value, timestamp, details

---

## PriceService

**Purpose:** Provides price data and currency conversion.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/PriceService.java`

`rate` and `decimals` feed amount conversion, so an unverified price is a wrong number a
caller acts on. `getPrices` (the cursor-paged `QueryPricesV2` endpoint) verifies each price against the `PRICEUPDATER` keys in the
SuperAdmin-verified rules container, which is why the service takes the
`RulesContainerCache` as a mandatory constructor argument. Whether prices must be signed is
the **container's** call: no `PRICEUPDATER` configured means this tenant does not sign prices
and the price passes through; a `PRICEUPDATER` configured plus a price carrying no signatures
throws `IntegrityException`.

### Methods

```java
PriceResult getPrices() throws ApiException
PriceResult getPrices(String fromCurrencyId, List<String> toCurrencyIds, Boolean onlyPrimary,
                      String sortOrder, Integer pageSize, String cursor) throws ApiException
List<PriceHistoryPoint> getPriceHistory(String base, String quote, int limit) throws ApiException
List<ConversionResult> convert(String currency, String amount, List<String> targetCurrencyIds) throws ApiException
```

The currency filter of `getPrices`: `fromCurrencyId` alone lists that currency's prices,
`toCurrencyIds` alone the prices into those currencies, both together the prices of one into the
others, neither every price. `getPriceHistory` cannot page: `limit` is 0 for the default (20) and
at most 365 daily points, newest first.

### Key Models

- `PriceResult` - prices, `getBaseCurrency()`, `getPage()` (`CursorPage`)
- `Price` - blockchain, currencyFrom, currencyTo, decimals, rate, signatures
- `PriceSignature` - userId, signature
- `PriceHistoryPoint` - timestamp, price
- `ConversionResult` - target currency, converted amount

---

## UserService

**Purpose:** Retrieves and manages user information.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/UserService.java`

### Methods

```java
User getMe() throws ApiException
User getUser(String userId) throws ApiException
UserResult getUsers(int limit, long offset) throws ApiException
List<User> getUsersByEmail(List<String> emails) throws ApiException      // walks every page
void createUserAttribute(String userId, String key, String value) throws ApiException
```

### Key Models

- `User` - id, email, name, roles, group memberships (`getGroups()`)
- `enforcedInRules` (user and each membership) is `true`/`false` from `getMe`, `getUsers` and
  `getUsersByEmail`, and `null` from `getUser`, whose endpoint does not compute it.
  `publicKeyEnforcedInRules` is computed by `getMe` only and is `null` everywhere else.

---

## ChangeService

**Purpose:** Manages configuration changes and approval workflows.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/ChangeService.java`

### Methods

```java
Change getChange(long id) throws ApiException
ChangeResult getChanges(String entity, String status, Integer pageSize, String cursor) throws ApiException
ChangeResult getChanges(String entity, String status, ApiRequestCursor cursor) throws ApiException
ChangeResult getChangesForApproval(Integer pageSize, String cursor) throws ApiException
ChangeResult getChangesForApproval(ApiRequestCursor cursor) throws ApiException
void approveChange(long id) throws ApiException
void approveChanges(List<Long> ids) throws ApiException
void rejectChange(long id) throws ApiException
void rejectChanges(List<Long> ids) throws ApiException
```

### Key Models

- `Change` - id, entity, operation, status, payload, trails
- `ChangeResult` - changes list + `getPage()` (`CursorPage`)

---

## BusinessRuleService

**Purpose:** Manages transaction approval rules.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/BusinessRuleService.java`

### Methods

```java
BusinessRuleResult getBusinessRules(Integer pageSize, String cursor) throws ApiException
BusinessRuleResult getBusinessRules(ApiRequestCursor cursor) throws ApiException
BusinessRuleResult getBusinessRulesByWallet(long walletId, Integer pageSize, String cursor) throws ApiException
BusinessRuleResult getBusinessRulesByWallet(long walletId, ApiRequestCursor cursor) throws ApiException
BusinessRuleResult getBusinessRulesByCurrency(String currencyId, Integer pageSize, String cursor) throws ApiException
BusinessRuleResult getBusinessRulesByCurrency(String currencyId, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `BusinessRule` - rule definition with conditions and actions
- `BusinessRuleResult` - rules list + `getPage()` (`CursorPage`)

---

## GovernanceRuleService

**Purpose:** Manages governance rules with SuperAdmin signature verification.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/GovernanceRuleService.java`

### Methods

#### getRules

Gets current governance rules with **signature verification**.

```java
GovernanceRules getRules() throws ApiException
GovernanceRules getRulesById(String id) throws ApiException
```

#### getRulesHistory

Token-paged: pass `getPage().getNextCursor()` back as `cursor`. Entries whose SuperAdmin
signatures do not verify are withheld (`getExcludedUnverified()`) and the page total is reduced
by them.

```java
GovernanceRulesHistoryResult getRulesHistory(Integer pageSize) throws ApiException
GovernanceRulesHistoryResult getRulesHistory(Integer pageSize, String cursor) throws ApiException
```

#### getRulesProposal

Gets pending rules proposal (SuperAdmin only).

```java
GovernanceRules getRulesProposal() throws ApiException
```

#### getPublicKeys

Lists SuperAdmin public keys.

```java
List<SuperAdminPublicKey> getPublicKeys() throws ApiException
```

#### verifyGovernanceRules

Manually verify rules against SuperAdmin keys.

```java
GovernanceRules verifyGovernanceRules(GovernanceRules rules, int minValidSignatures) throws ApiException
```

#### getDecodedRulesContainer

Decodes the rules container protobuf.

```java
DecodedRulesContainer getDecodedRulesContainer(GovernanceRules rules) throws ApiException
```

#### updateRulesProposal

Submits a typed rules container as a governance proposal (SuperAdmin only). The container
is encoded to the wire format internally; the server-controlled `enforcedRulesHash` and
`timestamp` fields are stripped. The endpoint returns no body — callers needing the
persisted proposal should call `getRulesProposal` (note: rules reads are cached
server-side, so an immediate read-back may be stale).

```java
void updateRulesProposal(DecodedRulesContainer container) throws ApiException
```

#### approveRulesProposal

Signs the pending proposal's rules container with a SuperAdmin private key (SHA-256 +
P-256 ECDSA, base64 raw r||s) and submits the approval. The signature binds the exact
pending content — review it first via `getRulesProposal` + `getDecodedRulesContainer`.

```java
void approveRulesProposal(PrivateKey privateKey, String comment) throws ApiException
```

#### rejectRulesProposal

Rejects the pending rules proposal with a comment (SuperAdmin only).

```java
void rejectRulesProposal(String comment) throws ApiException
```

### Key Models

- `GovernanceRules` - rulesContainer, rulesSignatures, locked, trails
- `DecodedRulesContainer` - lossless typed rules container (users, groups, transaction and
  whitelisting rules); round-trips through `RulesContainerMapper.toBase64String` /
  `fromBytes`
- `RuleCell` - typed transaction-rule cell union covering every cell type
  (`FiatAmountAny`, `FiatAmountRange`, `SourceInternalWallet`, `StringEqualValue`, ...);
  `RawCell` preserves cells from newer schemas verbatim. Decode and encode a cell with
  `RuleCellCodec.decode(columnType, bytes)` / `encode(columnType, cell)` — cells stay
  `List<ByteString>` on `RuleLine`
- `SuperAdminPublicKey` - id, publicKey, name

Cross-SDK cell wire-format parity is pinned by the shared golden vectors at
`scripts/resources/governance-cell-vectors.json` (monorepo root), consumed by every SDK's
test suite.

---

## WhitelistedAddressService

**Purpose:** Manages whitelisted addresses with comprehensive cryptographic verification.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/WhitelistedAddressService.java`

### Methods

#### getWhitelistedAddress

Gets a whitelisted address with **full verification**.

```java
WhitelistedAddress getWhitelistedAddress(long id) throws ApiException, WhitelistException
```

**Verification Steps:**
1. Metadata hash verification (SHA-256)
2. Rules container signature verification (SuperAdmin)
3. Hash coverage verification
4. Whitelist signature verification (governance thresholds)

#### getWhitelistedAddressEnvelope

Gets the signed envelope with all verification details.

```java
SignedWhitelistedAddressEnvelope getWhitelistedAddressEnvelope(long id) throws ApiException, WhitelistException
```

#### getWhitelistedAddresses

Lists a page of whitelisted addresses, verified (offset list). Rows that fail verification are
withheld and listed in `getExcludedUnverified()`; the page total is reduced by them while the
next offset stays in the server's row space.

```java
WhitelistedAddressListResult getWhitelistedAddresses(int limit, long offset)
WhitelistedAddressListResult getWhitelistedAddresses(int limit, long offset, String blockchain)
WhitelistedAddressListResult getWhitelistedAddresses(int limit, long offset, String blockchain, String network)
WhitelistedAddressListResult getWhitelistedAddresses(int limit, long offset, String blockchain, String network,
                                                     boolean rulesContainerNormalized)
WhitelistedAddressListResult getWhitelistedAddressesForApproval(int limit, long offset, List<String> ids,
                                                                Boolean includeAlreadySignedByUser)
```

### Key Models

- `WhitelistedAddress` - blockchain, network, address, addressType, memo, label, linkedInternalAddresses, linkedWallets
- `SignedWhitelistedAddressEnvelope` - signedAddress, metadata, rulesContainer, rulesSignatures, approvers, trails
- `WhitelistSignature` - hashes, signature (WhitelistUserSignature)

### Verification Details

See [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) for detailed verification flow.

---

## WebhookService

**Purpose:** Manages webhooks for receiving real-time event notifications.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/WebhookService.java`

### Methods

#### createWebhook

Creates a new webhook configuration.

```java
Webhook createWebhook(String url, String type, String secret) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| url | String | URL to receive webhook notifications (HTTPS) |
| type | String | Event type (e.g., "TRANSACTION", "REQUEST") |
| secret | String | Secret for signing webhook payloads |

**Returns:** `Webhook` - The created webhook

**Example:**
```java
Webhook webhook = client.getWebhookService().createWebhook(
    "https://example.com/webhook",
    "TRANSACTION",
    "my-secret-key"
);
System.out.println("Created webhook: " + webhook.getId());
```

#### getWebhooks

Lists webhooks with optional filtering.

```java
WebhookResult getWebhooks(String type, String url, Integer pageSize, String cursor) throws ApiException
WebhookResult getWebhooks(String type, String url, ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| type | String | Filter by webhook type (optional) |
| url | String | Filter by URL (optional) |
| pageSize | Integer | Page size, null or 0 for the default (20), at most 100 |
| cursor | String | A previous page's `getPage().getNextCursor()`, null for the first page |

#### deleteWebhook

Deletes a webhook configuration.

```java
void deleteWebhook(String webhookId) throws ApiException
```

#### updateWebhookStatus

Enables or disables a webhook.

```java
Webhook updateWebhookStatus(String webhookId, WebhookStatus status) throws ApiException
```

**Example:**
```java
// Disable a webhook
client.getWebhookService().updateWebhookStatus(webhookId, WebhookStatus.DISABLED);

// Re-enable it
client.getWebhookService().updateWebhookStatus(webhookId, WebhookStatus.ENABLED);
```

### Key Models

- `Webhook` - id, url, type, status, createdAt
- `WebhookStatus` - Enum: ENABLED, DISABLED, TIMEOUT
- `WebhookResult` - webhooks list with pagination

---

## StakingService

**Purpose:** Retrieves staking information across multiple proof-of-stake blockchains.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/StakingService.java`

### Methods

#### getADAStakePoolInfo

Retrieves Cardano stake pool information.

```java
ADAStakePoolInfo getADAStakePoolInfo(String network, String stakePoolId) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| network | String | Network (e.g., "mainnet", "preprod") |
| stakePoolId | String | Stake pool ID (Bech32 format) |

**Example:**
```java
ADAStakePoolInfo poolInfo = client.getStakingService()
    .getADAStakePoolInfo("mainnet", "pool1abc123...");
System.out.println("Pool pledge: " + poolInfo.getPledge());
```

#### getETHValidatorsInfo

Retrieves Ethereum validator information.

```java
List<ETHValidatorInfo> getETHValidatorsInfo(String network, List<String> ids) throws ApiException
```

**Example:**
```java
List<ETHValidatorInfo> validators = client.getStakingService()
    .getETHValidatorsInfo("mainnet", Arrays.asList("validator1", "validator2"));
for (ETHValidatorInfo v : validators) {
    System.out.println("Validator: " + v.getPublicKey() + ", Balance: " + v.getBalance());
}
```

#### getFTMValidatorInfo

Retrieves Fantom validator information.

```java
FTMValidatorInfo getFTMValidatorInfo(String network, String validatorAddress) throws ApiException
```

#### getICPNeuronInfo

Retrieves Internet Computer neuron information.

```java
ICPNeuronInfo getICPNeuronInfo(String network, String neuronId) throws ApiException
```

#### getNEARValidatorInfo

Retrieves NEAR Protocol validator information.

```java
NEARValidatorInfo getNEARValidatorInfo(String network, String validatorAddress) throws ApiException
```

#### getStakeAccounts

Lists stake accounts with pagination.

```java
StakeAccountResult getStakeAccounts(String addressId, String accountType,
                                     String accountAddress, Integer pageSize, String cursor) throws ApiException
StakeAccountResult getStakeAccounts(String addressId, String accountType,
                                     String accountAddress, ApiRequestCursor cursor) throws ApiException
```

#### getXTZStakingRewards

Retrieves Tezos staking rewards for an address.

```java
XTZStakingRewards getXTZStakingRewards(String network, String addressId,
                                        OffsetDateTime from, OffsetDateTime to) throws ApiException
```

**Example:**
```java
XTZStakingRewards rewards = client.getStakingService()
    .getXTZStakingRewards("mainnet", "address-123",
        OffsetDateTime.now().minusDays(30), OffsetDateTime.now());
System.out.println("Total rewards: " + rewards.getTotalRewards());
```

### Key Models

- `ADAStakePoolInfo` - pledge, margin, fixedCost, activeStake
- `ETHValidatorInfo` - publicKey, balance, status
- `FTMValidatorInfo` - stakedAmount, status
- `ICPNeuronInfo` - stake, votingPower, dissolveDelay
- `NEARValidatorInfo` - stake, fee
- `StakeAccountResult` - stake accounts with pagination
- `XTZStakingRewards` - totalRewards, cycles

---

## ContractWhitelistingService

**Purpose:** WRITE operations on whitelisted smart contract addresses (ERC20 tokens, NFTs, FA2 tokens).

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/ContractWhitelistingService.java`

> **Reads live on `WhitelistedAssetService`.** A whitelisted contract and a whitelisted asset
> are one server entity (`/whitelists/contracts`); the `getWhitelistedContract`,
> `getWhitelistedContracts`, `getWhitelistedContractsWithFilters` and
> `getWhitelistedContractsForApproval` that used to sit here returned the envelope with no
> verification, which made the verified reader avoidable. Use
> `client.getWhitelistedAssetService().getWhitelistedAssets(...)` /
> `.getWhitelistedAssetsForApproval(...)`, which run the six-step chain.

### Methods

#### createWhitelistedContract

Creates a new whitelisted contract address.

```java
String createWhitelistedContract(String blockchain, String network, String contractAddress,
                                  String symbol, String name, int decimals,
                                  String kind, String tokenId) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| blockchain | String | Blockchain identifier (e.g., "ETH", "MATIC") |
| network | String | Network (e.g., "mainnet", "goerli") |
| contractAddress | String | Smart contract address |
| symbol | String | Token symbol (e.g., "USDC") |
| name | String | Human-readable name |
| decimals | int | Token decimals (0 for NFTs) |
| kind | String | Contract kind (e.g., "erc20", "erc721") |
| tokenId | String | Token ID for NFTs (null for fungible tokens) |

**Returns:** `String` - The ID of the created whitelist entry

**Example:**
```java
String id = client.getContractWhitelistingService().createWhitelistedContract(
    "ETH", "mainnet", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
    "USDC", "USD Coin", 6, "erc20", null);
System.out.println("Created whitelist entry: " + id);
```

#### approveWhitelistedContracts

Approves one or more whitelisted contract addresses.

```java
@Deprecated
void approveWhitelistedContracts(List<String> ids, String signature, String comment) throws ApiException
```

> **Deprecated.** The signature is an opaque blob over hashes nothing verified, so the caller
> cannot know what they signed. Use
> `WhitelistedAssetService.approveWhitelistedAssets(ids, privateKey, comment)`, which re-reads
> and verifies the rows first.

**Example — write here, read through the verified reader:**
```java
String id = client.getContractWhitelistingService().createWhitelistedContract(
    "ETH", "mainnet", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
    "USDC", "USD Coin", 6, "erc20", null);

WhitelistedAssetResult result = client.getWhitelistedAssetService()
    .getWhitelistedAssets(50, 0, "ETH", "mainnet", null, null, null, null);
for (SignedWhitelistedAssetEnvelope asset : result.getAssets()) {
    System.out.println(asset.getWhitelistedAsset().getContractAddress());
}
```

#### updateWhitelistedContract

Updates an existing whitelisted contract.

```java
void updateWhitelistedContract(String id, String symbol, String name, int decimals) throws ApiException
```

#### createAttribute

Creates an attribute on a whitelisted contract.

```java
List<Attribute> createAttribute(String contractId, String key, String value,
                                 String contentType, String type, String subType) throws ApiException
```

#### getAttribute

Retrieves an attribute from a whitelisted contract.

```java
Attribute getAttribute(String contractId, String attributeId) throws ApiException
```

### Key Models

- `Attribute` - key, value, contentType, type, subType

---

## WhitelistedAssetService

**Purpose:** Manages whitelisted assets/contracts with cryptographic verification. **This is
the only verified reader of `/whitelists/contracts`** — a whitelisted asset and a whitelisted
contract are one server entity, and `ContractWhitelistingService` is write-only.

**Access:** `client.getWhitelistedAssetService()`

### Methods

#### getWhitelistedAsset

Gets a whitelisted asset by ID with verification.

```java
WhitelistedAsset getWhitelistedAsset(long id) throws ApiException, WhitelistException
```

#### getWhitelistedAssets

Lists a page of whitelisted assets, verified (offset list). Every overload returns a
`WhitelistedAssetResult` with `getPagination()`. Skipped rows keep their SQL slot on this
endpoint, so a short page is not the end: continue while `getPagination().hasMore()`.

```java
WhitelistedAssetResult getWhitelistedAssets(int limit, long offset)
WhitelistedAssetResult getWhitelistedAssets(int limit, long offset, String blockchain)
WhitelistedAssetResult getWhitelistedAssets(int limit, long offset, String blockchain, String network)
WhitelistedAssetResult getWhitelistedAssets(int limit, long offset,
                                            String blockchain, String network,
                                            String query, Boolean includeForApproval,
                                            List<String> kindTypes, List<String> ids)
        throws ApiException, WhitelistException
```

#### getWhitelistedAssetsForApproval

Lists whitelisted assets awaiting approval, verified as in `getWhitelistedAssets`. Without
this the only reader of the for-approval endpoint was the unverified contract service, so the
rows an approver inspects were never checked against governance.

```java
WhitelistedAssetResult getWhitelistedAssetsForApproval(int limit, long offset, List<String> ids)
        throws ApiException, WhitelistException
```

#### approveWhitelistedAssets

Signs and submits an approval, **all-or-nothing**. Each asset is re-read and verified, and the
hashes those rows carry are what gets signed; a row that is missing or fails verification
aborts the whole call and nothing is signed. The API takes one signature covering the whole
batch, so a partial approval would mean the caller believes they approved more than they did.

```java
void approveWhitelistedAssets(WhitelistedAssetApproval selection, PrivateKey privateKey,
        String comment) throws ApiException, WhitelistException
```

> **`selection` is the ROWS a verified read returned, not bare ids.** Mint it with
> `result.select(ids)` or `result.selectAll()` off a `getWhitelistedAssetsForApproval` result.
> It carries the metadata hash each row had **at review time**, and the approval aborts if the
> re-read hash differs.
>
> Without that pin a response-controlling server could answer the id-filtered re-read with a
> *different* row — one whose existing signatures already satisfy the container it presents —
> and harvest a genuine approver signature over content the approver never saw. Verification
> alone does not catch it: the substituted row is a real, validly-signed entry, just not the
> one that was reviewed. Same mitigation as `approveRulesProposal`'s mandatory
> `expectedContainerHash`. `WhitelistedAssetApproval` has no public constructor, so the pin
> cannot be forgotten; an empty selection raises rather than meaning "approve nothing".
>
> The value **signed** is still the row's current `metadata.hash`, never the legacy variant
> step 4 matched — validatord rebuilds the hash array from the current schema and verifies the
> submitted signature against those bytes.

```java
WhitelistedAssetResult reviewed = client.getWhitelistedAssetService()
        .getWhitelistedAssetsForApproval(50, 0, null);
client.getWhitelistedAssetService()
        .approveWhitelistedAssets(reviewed.selectAll(), approverKey, "reviewed");
```

`WhitelistedAddressService.approveWhitelistedAddresses(WhitelistedAddressApproval, PrivateKey,
String)` is the address peer, pinned the same way off a `WhitelistedAddressListResult`. Its
re-read goes through the **normalized** list path, where containers are response-level and
label-verified, rather than the per-row in-band containers it used before.

### Key Models

- `WhitelistedAsset` - id, blockchain, network, status, metadata, signedContractAddress
- `WhitelistedAssetResult` - assets list with totalItems and an overflow-safe `hasMore(currentOffset, pageSize)`

---

## AuditService

**Purpose:** Queries audit trail events.

**Access:** `client.getAuditService()`

### Methods

#### getAuditTrails

Lists audit events with filtering.

```java
AuditTrailResult getAuditTrails(String externalUserId, List<String> entities, List<String> actions,
                                OffsetDateTime from, OffsetDateTime to, Integer pageSize, String cursor)
        throws ApiException
AuditTrailResult getAuditTrails(String externalUserId, List<String> entities, List<String> actions,
                                OffsetDateTime from, OffsetDateTime to, ApiRequestCursor cursor)
        throws ApiException
```

### Key Models

- `AuditTrailResult` - audit trails + `getPage()` (`CursorPage`)
- `AuditTrail` - entity, action, details, creationDate

---

## FeeService

**Purpose:** Retrieves transaction fee information.

**Access:** `client.getFeeService()`

### Methods

#### getFees

Gets the current network fee of every currency (the V2 endpoint).

```java
List<Fee> getFees() throws ApiException
```

### Key Models

- `Fee` - currencyId, value, denom, currencyInfo (`Currency`), updateDate

---

## AirGapService

**Purpose:** Provides air-gap signing operations for offline transaction signing.

**Access:** `client.getAirGapService()`

### Methods

#### getOutgoingAirGap

Gets an air-gap request for offline signing.

```java
AirGapRequest getOutgoingAirGap(long requestId) throws ApiException
```

#### submitIncomingAirGap

Submits a signature for an air-gap request.

```java
void submitIncomingAirGap(long requestId, String signature) throws ApiException
```

---

## ReservationService

**Purpose:** Manages balance reservations for addresses.

**Access:** `client.getReservationService()`

### Methods

#### getReservations

Lists reservations.

```java
ReservationResult getReservations() throws ApiException
ReservationResult getReservations(String kind, String address, String addressId, List<String> kinds,
                                  Integer pageSize, String cursor) throws ApiException
ReservationResult getReservations(String kind, String address, String addressId, List<String> kinds,
                                  ApiRequestCursor cursor) throws ApiException
```

A continuation sends `pageRequest=NEXT` with the cursor and the page size (it used to send the
cursor alone, which the server cannot page on).

### Key Models

- `ReservationResult` - reservations + `getPage()` (`CursorPage`)
- `Reservation` - id, addressId, amount, status, expiresAt

---

## MultiFactorSignatureService

**Purpose:** Manages multi-factor signature operations for enhanced security.

**Access:** `client.getMultiFactorSignatureService()`

### Methods

Multi-factor signatures are a **second approval channel** over the same entities this SDK
otherwise protects (a request, a whitelisted address, a whitelisted contract). The caller is
the second-factor signing device, and the signature it submits is precisely the artefact a
compromised server cannot forge on its own.

> **Security — `payloadToSign` is UNVERIFIED server data, and `approve` is opaque to the SDK.**
> The reply carries only `{id, payloadToSign[], entityType}` with **no entity id**, so nothing
> in it can be joined back to the entities the request was created for and the SDK has nothing
> to check against. A compromised server can therefore answer with the metadata hash of an
> entity of its choosing under the expected kind; sign it and the server holds a valid
> MobileAppSigner approval over an entity nobody reviewed.
>
> **Bind it yourself.** `createMultiFactorSignatures` takes the entity IDs — keep them, re-read
> those entities through the verifying reader for that kind (`RequestService`,
> `WhitelistedAddressService`, `WhitelistedAssetService`) and require each `payloadToSign`
> element to equal the locally recomputed, verified metadata hash. Tracked in `TODOS.md`; this
> is the one read path in the SDK that returns bytes intended for a signing key without
> verifying them, and it is deliberate-but-unresolved rather than an oversight.

#### getMultiFactorSignatureInfo

Retrieves a multi-factor signature request. **`payloadToSign` is unverified** — see above.

```java
MultiFactorSignatureInfo getMultiFactorSignatureInfo(String id) throws ApiException
```

#### createMultiFactorSignatures

Creates a batch of multi-factor signature requests. Keep the entity IDs.

```java
MultiFactorSignatureResult createMultiFactorSignatures(List<String> entityIDs,
        MultiFactorSignatureEntityType entityType) throws ApiException
```

#### approveMultiFactorSignature

Submits a caller-produced signature. **The SDK never learns what it covers** — see above.

```java
MultiFactorSignatureApprovalResult approveMultiFactorSignature(String id, String signature,
        String comment) throws ApiException
```

#### rejectMultiFactorSignature

Rejects a multi-factor signature request.

```java
void rejectMultiFactorSignature(String id, String comment) throws ApiException
```

---

## GroupService

**Purpose:** Manages user groups for approval workflows.

**Access:** `client.getGroupService()`

### Methods

#### getGroups

Lists user groups.

```java
GroupResult getGroups() throws ApiException
GroupResult getGroups(int limit, long offset, List<String> ids, List<String> externalGroupIds,
                      String query) throws ApiException
```

### Key Models

- `GroupResult` - groups + `getPagination()` (`OffsetPagination`)
- `Group` - id, name, members, threshold; `enforcedInRules` on the group and each of its users
  is always `true`/`false` from `getGroups`

---

## VisibilityGroupService

**Purpose:** Manages visibility groups for resource access control.

**Access:** `client.getVisibilityGroupService()`

### Methods

#### getVisibilityGroups

Lists visibility groups.

```java
List<VisibilityGroup> getVisibilityGroups(int limit, int offset) throws ApiException
```

### Key Models

- `VisibilityGroup` - id, name, members

---

## ConfigService

**Purpose:** Manages system configuration settings.

**Access:** `client.getConfigService()`

### Methods

#### getTenantConfig

Gets tenant configuration.

```java
TenantConfig getTenantConfig() throws ApiException
```

### Key Models

- `TenantConfig` - tenantId, settings, features

---

## WebhookCallsService

**Purpose:** Queries webhook call history.

**Access:** `client.getWebhookCallsService()`

### Methods

#### getWebhookCalls

Lists webhook calls with filtering.

```java
WebhookCallResult getWebhookCalls(String eventId, String webhookId, String status, String sortOrder,
                                  Integer pageSize, String cursor) throws ApiException
WebhookCallResult getWebhookCalls(String eventId, String webhookId, String status, String sortOrder,
                                  ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `WebhookCall` - id, webhookId, status, requestBody, responseCode, timestamp

---

## TagService

**Purpose:** Manages tags for organizing resources.

**Access:** `client.getTagService()`

### Methods

#### getTags

Lists all tags.

```java
List<Tag> getTags() throws ApiException
```

#### createTag

Creates a new tag.

```java
Tag createTag(String name, String color) throws ApiException
```

### Key Models

- `Tag` - id, name, color

---

## AssetService

**Purpose:** Retrieves asset information.

**Access:** `client.getAssetService()`

`getAssetAddresses` verifies every address's HSM signature — the same check `AddressService`
runs — and **fails fast** on the first that does not verify. It returns the same entity, so
returning it unverified made `AddressService`'s mandatory verification avoidable. The service
therefore takes the `RulesContainerCache` as a mandatory constructor argument.

`queryAssetAddresses` rows carry no signature, so the service also takes the `AddressService`
and `WhitelistedAddressService` it confirms them through (`ProtectClient` passes its own):

| Row type | Confirmed by | Kept row |
|---|---|---|
| `ADDRESS_TYPE_V2_INTERNAL` | its `addressID`, re-read through the HSM-verified managed-address list (≤ 50 ids per request) | `isVerified() == true`, address from the verified address |
| `ADDRESS_TYPE_V2_WHITELISTED` | its `whitelistedAddressID`, re-read through the six-step verified whitelist (≤ 100 ids per request) | `isVerified() == true`, address from the verified envelope |
| anything else (EXTERNAL, no type) | nothing, no extra request | `isVerified() == false` |

An INTERNAL/WHITELISTED row with no id, one the verified read does not return or rejects, or
one whose address differs from the verified one is withheld and named in
`getExcludedUnverified()` (id = the addressID / whitelistedAddressID, else the address); rows
came back but none survived throws `IntegrityException`. Failures of the readers themselves
propagate: a request error, a rules container without an HSMSLOT key
(`ContainerIntegrityException`), a whitelist container that fails verification
(`WhitelistException`). Exclusions never move the cursor. The generated client rejects an
`addressType` value it does not know while parsing the reply, so an unknown type fails the call.

### Methods

```java
AssetAddressesResult getAssetAddresses(String currency) throws ApiException
AssetAddressesResult getAssetAddresses(String currency, String walletId, String addressId,
                                       Integer pageSize, String cursor) throws ApiException
AssetAddressesResult getAssetAddresses(String currency, String walletId, String addressId,
                                       ApiRequestCursor cursor) throws ApiException
AssetWalletsResult getAssetWallets(String currency) throws ApiException
AssetWalletsResult getAssetWallets(String currency, Integer pageSize, String cursor) throws ApiException
AssetWalletsResult getAssetWallets(String currency, ApiRequestCursor cursor) throws ApiException
AssetV2Result queryAssets(String blockchain, String network, String symbol, String contractAddress,
                          String label, String currencyName, Integer pageSize, String cursor) throws ApiException
AssetAddressV2Result queryAssetAddresses(String assetId, String addressType, String kycStatus,
                                         Integer pageSize, String cursor)
        throws ApiException, WhitelistException
AssetOperationV2Result listAssetOperations(String assetId, String type, String status,
                                           Integer pageSize, String cursor) throws ApiException
```

`getAssetAddresses`/`getAssetWallets` page through the body `requestCursor` only, and their
pages carry the server total. `queryAssets`, `queryAssetAddresses` and `listAssetOperations` are
the v2 asset service; filter values are the wire enum names (e.g. `ADDRESS_TYPE_V2_INTERNAL`,
`KYC_STATUS_V2_APPROVED`, `ASSET_OPERATION_TYPE_V2_MINT`), and an unknown one is rejected by name.

### Key Models

- `AssetAddressesResult` / `AssetWalletsResult` - verified addresses / wallets + `getPage()`
- `AssetV2` - id, label, assetType, status, blockchain, network, currencyId, name, symbol, decimals, contractAddress, attributes, cantonNativeToken
- `AssetAddressV2` - address, kycStatus, balance, addressType, addressId, whitelistedAddressId, `isVerified()`
- `AssetAddressV2Result` - addresses, `getExcludedUnverified()` + `getPage()`
- `AssetOperationV2` - id, assetId, type, status, dates, failure/blocking reason, and the type's details (flattened)
- `Asset` - id, symbol, name, blockchain, network, contractAddress, decimals

---

## ActionService

**Purpose:** Manages actions in the system.

**Access:** `client.getActionService()`

### Methods

#### getActions

Lists a page of actions (offset list).

```java
ActionResult getActions() throws ApiException
ActionResult getActions(int limit, long offset, List<String> ids) throws ApiException
```

---

## BlockchainService

**Purpose:** Retrieves blockchain information.

**Access:** `client.getBlockchainService()`

### Methods

#### getBlockchains

Lists supported blockchains.

```java
List<Blockchain> getBlockchains() throws ApiException
```

### Key Models

- `Blockchain` - id, name, networks, features

---

## EarnService

**Purpose:** Lists the earn rewards credited to addresses (for example Merkl token rewards).

**Access:** `client.getEarnService()`

### Methods

#### listRewards

Lists a page of rewards (cursor list).

```java
EarnRewardResult listRewards(String recipientAddressId, Integer pageSize, String cursor) throws ApiException
```

### Key Models

- `EarnRewardResult` - rewards + `getPage()` (`CursorPage`)
- `EarnReward` - id, recipientAddressId, recipientAddress, rewardType, amount, claimed, pending, token address/symbol/assetId

---

## ExchangeService

**Purpose:** Manages exchange integrations.

**Access:** `client.getExchangeService()`

### Methods

#### getExchange

Lists configured exchanges.

```java
List<Exchange> getExchange() throws ApiException
```

### Key Models

- `Exchange` - id, name, type, status

---

## FiatService

**Purpose:** Manages fiat currency operations.

**Access:** `client.getFiatService()`

### Methods

#### getFiatProviderAccounts / getFiatProviderCounterpartyAccounts / getFiatProviderOperations / listFiatProviderEntities

Cursor lists. `provider` and `label` are required for the two account lists.

```java
FiatProviderAccountResult getFiatProviderAccounts(String provider, String label, String accountType,
                                                  String sortOrder, Integer pageSize, String cursor) throws ApiException
FiatProviderCounterpartyAccountResult getFiatProviderCounterpartyAccounts(String provider, String label,
        String counterpartyId, String sortOrder, Integer pageSize, String cursor) throws ApiException
FiatProviderOperationResult getFiatProviderOperations(String provider, String label, String sortOrder,
                                                      Integer pageSize, String cursor) throws ApiException
FiatProviderEntityResult listFiatProviderEntities(String provider, String label, String sortOrder,
                                                  Integer pageSize, String cursor) throws ApiException
```

The first three also keep an `ApiRequestCursor` overload.

---

## FeePayerService

**Purpose:** Manages fee payer configurations.

**Access:** `client.getFeePayerService()`

### Methods

#### getFeePayers

Lists fee payers.

```java
FeePayerResult getFeePayers() throws ApiException
FeePayerResult getFeePayers(int limit, long offset, List<String> ids, String blockchain,
                            String network) throws ApiException
```

---

## HealthService

**Purpose:** Checks system health status.

**Access:** `client.getHealthService()`

### Methods

#### getAllHealthChecks

Checks API health.

```java
HealthStatus getAllHealthChecks() throws ApiException
```

### Key Models

- `HealthStatus` - status, version, timestamp

---

## JobService

**Purpose:** Manages background jobs.

**Access:** `client.getJobService()`

### Methods

#### getJobs

Lists background jobs.

```java
List<Job> getJobs(String status, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `Job` - id, type, status, progress, createdAt

---

## StatisticsService

**Purpose:** Retrieves platform statistics.

**Access:** `client.getStatisticsService()`

### Methods

#### getPortfolioStatistics

Gets portfolio statistics.

```java
PortfolioStatistics getPortfolioStatistics() throws ApiException
```

### Key Models

- `PortfolioStatistics` - totalValue, assetBreakdown, changePercent

---

## TokenMetadataService

**Purpose:** Retrieves token metadata information.

**Access:** `client.getTokenMetadataService()`

### Methods

#### getEVMERCTokenMetadata

Gets ERC token metadata (ERC-20/721/1155) on an EVM chain. The deprecated `GetERCTokenMetadata`
endpoint is not wrapped.

```java
TokenMetadata getEVMERCTokenMetadata(String network, String contract, String tokenId, Boolean withData,
                                     String blockchain) throws ApiException
```

### Key Models

- `TokenMetadata` - name, symbol, decimals, totalSupply, logoUrl

---

## UserDeviceService

**Purpose:** Manages user device registrations.

**Access:** `client.getUserDeviceService()`

### Methods

### Key Models

- `UserDevice` - id, deviceType, name, lastUsed, status

---

## TaurusNetwork Services

The TaurusNetwork services are accessed through a namespace pattern:

```java
// Access TaurusNetwork services
client.taurusNetwork().participants()
client.taurusNetwork().pledges()
client.taurusNetwork().lending()
client.taurusNetwork().settlements()
client.taurusNetwork().sharing()
```

---

## TaurusNetworkParticipantService

**Purpose:** Provides access to Taurus Network participant management.

**Access:** `client.taurusNetwork().participants()`

### Methods

#### getMyParticipant

Retrieves the current participant.

```java
Participant getMyParticipant() throws ApiException
```

**Example:**
```java
Participant me = client.taurusNetwork().participants().getMyParticipant();
System.out.println("My participant ID: " + me.getId());
```

#### get

Retrieves a participant by ID.

```java
Participant get(String participantId, Boolean includeTotalPledgesValuation) throws ApiException
```

#### list

Retrieves multiple participants by IDs.

```java
List<Participant> list(List<String> participantIds, Boolean includeTotalPledgesValuation) throws ApiException
```

### Key Models

- `Participant` - id, name, country, publicKey, totalPledgesValuation

---

## TaurusNetworkPledgeService

**Purpose:** Provides access to Taurus Network pledge lifecycle operations.

**Access:** `client.taurusNetwork().pledges()`

### Methods

#### get

Retrieves a pledge by ID.

```java
Pledge get(String pledgeId) throws ApiException
```

#### list

Retrieves pledges with optional filtering.

```java
PledgeResult list(String ownerParticipantId, String targetParticipantId,
                        List<String> sharedAddressIds, String currencyId,
                        String sortOrder, Integer pageSize, String cursor) throws ApiException
PledgeResult list(String ownerParticipantId, String targetParticipantId,
                        List<String> sharedAddressIds, String currencyId,
                        String sortOrder, ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| ownerParticipantId | String | Filter by owner participant ID (optional) |
| targetParticipantId | String | Filter by target participant ID (optional) |
| sharedAddressIds | List<String> | Filter by shared address IDs (optional) |
| currencyId | String | Filter by currency ID (optional) |
| sortOrder | String | Sort order: "ASC" or "DESC" (optional) |
| pageSize | Integer | Page size, null or 0 for the default (20), at most 100 |
| cursor | String | A previous page's `getPage().getNextCursor()`, null for the first page |

#### listWithdrawals

Retrieves pledge withdrawals, optionally for one pledge.

```java
PledgeWithdrawalResult listWithdrawals(String pledgeId, String withdrawalStatus,
                                             String sortOrder, Integer pageSize, String cursor) throws ApiException
PledgeWithdrawalResult listWithdrawals(String pledgeId, String withdrawalStatus,
                                             String sortOrder, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `Pledge` - id, ownerParticipantId, targetParticipantId, amount, status, currency
- `PledgeResult` - pledges list with pagination cursor
- `PledgeAction` - id, pledgeId, actionType, status, metadata
- `PledgeWithdrawal` - withdrawal details with status

---

## TaurusNetworkLendingService

**Purpose:** Provides access to lending offers and agreements in the Taurus Network.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/TaurusNetworkLendingService.java`

### Methods

#### getLendingOffer

Retrieves a lending offer by ID.

```java
LendingOffer getLendingOffer(String offerId) throws ApiException
```

**Example:**
```java
LendingOffer offer = client.taurusNetwork().lending()
    .getLendingOffer("offer-123");
System.out.println("Offer rate: " + offer.getRate());
```

#### getLendingOffers

Retrieves lending offers with optional filtering.

```java
LendingOfferResult getLendingOffers(List<String> currencyIds, String participantId,
                                     String duration, String sortOrder,
                                     Integer pageSize, String cursor) throws ApiException
LendingOfferResult getLendingOffers(List<String> currencyIds, String participantId,
                                     String duration, String sortOrder,
                                     ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| currencyIds | List<String> | Filter by currency IDs (optional) |
| participantId | String | Filter by participant ID (optional) |
| duration | String | Filter by duration (optional) |
| sortOrder | String | Sort order: "ASC" or "DESC" (optional) |
| pageSize | Integer | Page size, null or 0 for the default (20), at most 100 |
| cursor | String | A previous page's `getPage().getNextCursor()`, null for the first page |

**Example:**
```java
LendingOfferResult result = client.taurusNetwork().lending()
    .getLendingOffers(Arrays.asList("ETH"), null, null, "DESC", 20, null);
for (LendingOffer offer : result.getOffers()) {
    System.out.println("Offer: " + offer.getId());
}
```

#### getLendingAgreement

Retrieves a lending agreement by ID.

```java
LendingAgreement getLendingAgreement(String agreementId) throws ApiException
```

#### getLendingAgreements

Retrieves lending agreements with optional filtering.

```java
LendingAgreementResult getLendingAgreements(String sortOrder, Integer pageSize, String cursor) throws ApiException
LendingAgreementResult getLendingAgreements(String sortOrder, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `LendingOffer` - id, participantId, currencyId, amount, rate, duration, status
- `LendingOfferResult` - offers list with pagination cursor
- `LendingAgreement` - id, offerId, borrowerParticipantId, lenderParticipantId, amount, status
- `LendingAgreementResult` - agreements list with pagination cursor

---

## TaurusNetworkSettlementService

**Purpose:** Provides access to settlements in the Taurus Network.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/TaurusNetworkSettlementService.java`

### Methods

#### getSettlement

Retrieves a settlement by ID.

```java
Settlement getSettlement(String settlementId) throws ApiException
```

**Example:**
```java
Settlement settlement = client.taurusNetwork().settlements()
    .getSettlement("settlement-123");
System.out.println("Settlement status: " + settlement.getStatus());
```

#### getSettlements

Retrieves settlements with optional filtering.

```java
SettlementResult getSettlements(String counterParticipantId, List<String> statuses,
                                 String sortOrder, Integer pageSize, String cursor) throws ApiException
SettlementResult getSettlements(String counterParticipantId, List<String> statuses,
                                 String sortOrder, ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| counterParticipantId | String | Filter by counter participant ID (optional) |
| statuses | List<String> | Filter by statuses (optional) |
| sortOrder | String | Sort order: "ASC" or "DESC" (optional) |
| pageSize | Integer | Page size, null or 0 for the default (20), at most 100 |
| cursor | String | A previous page's `getPage().getNextCursor()`, null for the first page |

**Example:**
```java
SettlementResult result = client.taurusNetwork().settlements()
    .getSettlements(null, Arrays.asList("PENDING", "COMPLETED"), "DESC", 20, null);
for (Settlement settlement : result.getSettlements()) {
    System.out.println("Settlement: " + settlement.getId() + ", Status: " + settlement.getStatus());
}
```

### Key Models

- `Settlement` - id, counterParticipantId, amount, currency, status, createdAt
- `SettlementResult` - settlements list with pagination cursor

---

## TaurusNetworkSharingService

**Purpose:** Provides access to address and asset sharing in the Taurus Network.

**Access:** `client.taurusNetwork().sharing()`

### Methods

#### listSharedAddresses

Retrieves shared addresses with optional filtering.

```java
SharedAddressResult listSharedAddresses(String participantId, String ownerParticipantId,
                                        String targetParticipantId, String blockchain,
                                        String network, List<String> ids,
                                        String sortOrder, ApiRequestCursor cursor) throws ApiException
```

Seven filters leave no room for separate page arguments (Checkstyle caps a method at eight), so
the page is passed as `Pagination.page(pageSize, cursor)`.

### Key Models

- `SharedAddress` - id, blockchain, network, address, ownerParticipantId, targetParticipantId, permissions
- `SharedAddressResult` - addresses list with pagination cursor
- `SharedAsset` - id, assetId, ownerParticipantId, targetParticipantId, permissions
- `SharedAssetResult` - assets list with pagination cursor

---

## Exception Handling

All services follow a consistent exception pattern:

```java
try {
    Wallet wallet = client.getWalletService().getWallet(walletId);
} catch (ApiException e) {
    System.err.println("API Error: " + e.getMessage());
    System.err.println("Error Code: " + e.getErrorCode());
    System.err.println("HTTP Status: " + e.getCode());
}
```

### Exception Types

| Exception | Type | When Thrown |
|-----------|------|-------------|
| `ApiException` | Checked | General API errors (network, auth, validation) |
| `AuthenticationException` | Checked (extends `ApiException`) | Authentication failures (HTTP 401) |
| `AuthorizationException` | Checked (extends `ApiException`) | Authorization failures (HTTP 403) |
| `NotFoundException` | Checked (extends `ApiException`) | Resource not found (HTTP 404) |
| `RateLimitException` | Checked (extends `ApiException`) | Rate limit exceeded (HTTP 429) |
| `ServerException` | Checked (extends `ApiException`) | Server errors (HTTP 5xx) |
| `ValidationException` | Checked (extends `ApiException`) | Input validation errors (HTTP 400) |
| `IntegrityException` | Unchecked (extends `SecurityException`) | Hash/signature verification failed |
| `WhitelistException` | Checked | Whitelist-specific verification errors |
| `ConfigurationException` | Checked | Client configuration errors |
| `RequestMetadataException` | Checked | Metadata payload parsing errors |

---

## Pagination Patterns

### Cursor-Based (Recommended)

Used by: BalanceService, BusinessRuleService, ChangeService, RequestService

```java
ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);
do {
    Result result = service.getItems(cursor);
    // Process items
    cursor = result.nextCursor(100);
} while (result.hasNext());
```

### Offset-Based

Used by: WalletService, AddressService, TransactionService, UserService

```java
int limit = 100;
int offset = 0;
List<Item> allItems = new ArrayList<>();
List<Item> page;
do {
    page = service.getItems(limit, offset);
    allItems.addAll(page);
    offset += limit;
} while (!page.isEmpty());
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Verification details

<!-- BEGIN GENERATED METHOD INDEX -->

## Complete Method Index

Generated from the java source by `scripts/api-surface/docs.py`; regenerate with
`./build.sh docs`. Every method below exists in the SDK, and `./build.sh docs --check`
fails if this list drifts or if the prose above documents a method that does not.

44 services, 222 public methods.

### ActionService

- `getAction(String): ActionEnvelope`
- `getActions(): ActionResult`
- `getActions(int, long, List<String>): ActionResult`

### AddressService

- `createAddress(CreateAddressRequest): Address`
- `createAddress(long, String, String, String): Address`
- `createAddressAttribute(long, String, String): void`
- `deleteAddressAttribute(long, long): void`
- `getAddress(long): Address`
- `getAddressProofOfReserve(long, String): TgvalidatordProofOfReserve`
- `getAddresses(Long, int, long, Boolean): AddressResult`
- `getAddresses(long, int, long): AddressResult`

### AirGapService

- `getOutgoingAirGap(List<String>): File`
- `submitIncomingAirGap(String): void`

### AssetService

- `getAssetAddresses(String): AssetAddressesResult`
- `getAssetAddresses(String, String, String, ApiRequestCursor): AssetAddressesResult`
- `getAssetAddresses(String, String, String, Integer, String): AssetAddressesResult`
- `getAssetWallets(String): AssetWalletsResult`
- `getAssetWallets(String, ApiRequestCursor): AssetWalletsResult`
- `getAssetWallets(String, Integer, String): AssetWalletsResult`
- `listAssetOperations(String, String, String, Integer, String): AssetOperationV2Result`
- `queryAssetAddresses(String, String, String, Integer, String): AssetAddressV2Result`
- `queryAssets(String, String, String, String, String, String, Integer, String): AssetV2Result`

### AuditService

- `exportAuditTrails(String, List<String>, List<String>, OffsetDateTime, OffsetDateTime, String): String`
- `getAuditTrails(String, List<String>, List<String>, OffsetDateTime, OffsetDateTime, ApiRequestCursor): AuditTrailResult`
- `getAuditTrails(String, List<String>, List<String>, OffsetDateTime, OffsetDateTime, Integer, String): AuditTrailResult`

### BalanceService

- `getBalances(ApiRequestCursor): BalanceResult`
- `getBalances(String, ApiRequestCursor): BalanceResult`
- `getBalances(String, Integer, String): BalanceResult`
- `getNFTCollectionBalances(String, String, ApiRequestCursor): NFTCollectionBalanceResult`
- `getNFTCollectionBalances(String, String, Integer, String): NFTCollectionBalanceResult`

### BlockchainService

- `getBlockchains(): List<BlockchainInfo>`
- `getBlockchains(String, String, Boolean): List<BlockchainInfo>`

### BusinessRuleService

- `getBusinessRules(ApiRequestCursor): BusinessRuleResult`
- `getBusinessRules(Integer, String): BusinessRuleResult`
- `getBusinessRulesByCurrency(String, ApiRequestCursor): BusinessRuleResult`
- `getBusinessRulesByCurrency(String, Integer, String): BusinessRuleResult`
- `getBusinessRulesByWallet(long, ApiRequestCursor): BusinessRuleResult`
- `getBusinessRulesByWallet(long, Integer, String): BusinessRuleResult`
- `updateTransactionsEnabled(boolean): void`

### ChangeService

- `approveChange(String): void`
- `approveChanges(List<String>): void`
- `createChange(CreateChangeRequest): String`
- `getChange(String): Change`
- `getChanges(String, String, ApiRequestCursor): ChangeResult`
- `getChanges(String, String, Integer, String): ChangeResult`
- `getChangesForApproval(ApiRequestCursor): ChangeResult`
- `getChangesForApproval(Integer, String): ChangeResult`
- `rejectChange(String): void`
- `rejectChanges(List<String>): void`

### ConfigService

- `getTenantConfig(): TenantConfig`

### ContractWhitelistingService

- `approveWhitelistedContracts(List<String>, String, String): void`
- `createAttribute(String, String, String, String, String, String): List<Attribute>`
- `createWhitelistedContract(String, String, String, String, String, int, String, String): String`
- `getAttribute(String, String): Attribute`
- `updateWhitelistedContract(String, String, String, int): void`

### CurrencyService

- `getBaseCurrency(): Currency`
- `getCurrencies(): List<Currency>`
- `getCurrencies(boolean, boolean): List<Currency>`
- `getCurrency(String): Currency`
- `getCurrencyByBlockchain(String, String): Currency`

### EarnService

- `listRewards(String, Integer, String): EarnRewardResult`

### ExchangeService

- `exportExchanges(String): String`
- `getExchange(String): Exchange`
- `getExchangeCounterparties(): List<ExchangeCounterparty>`
- `getExchangeWithdrawalFee(String, String, String): ExchangeWithdrawalFee`

### FeePayerService

- `getFeePayer(String): FeePayer`
- `getFeePayers(): FeePayerResult`
- `getFeePayers(int, long, List<String>, String, String): FeePayerResult`

### FeeService

- `getFees(): List<Fee>`

### FiatService

- `getFiatProviderAccount(String): FiatProviderAccount`
- `getFiatProviderAccounts(String, String, String, String, ApiRequestCursor): FiatProviderAccountResult`
- `getFiatProviderAccounts(String, String, String, String, Integer, String): FiatProviderAccountResult`
- `getFiatProviderCounterpartyAccount(String): FiatProviderCounterpartyAccount`
- `getFiatProviderCounterpartyAccounts(String, String, String, String, ApiRequestCursor): FiatProviderCounterpartyAccountResult`
- `getFiatProviderCounterpartyAccounts(String, String, String, String, Integer, String): FiatProviderCounterpartyAccountResult`
- `getFiatProviderOperation(String): FiatProviderOperation`
- `getFiatProviderOperations(String, String, String, ApiRequestCursor): FiatProviderOperationResult`
- `getFiatProviderOperations(String, String, String, Integer, String): FiatProviderOperationResult`
- `getFiatProviders(): List<FiatProvider>`
- `listFiatProviderEntities(String, String, String, Integer, String): FiatProviderEntityResult`

### GovernanceRuleService

- `approveRulesProposal(PrivateKey, String, String): void`
- `decodeProposalForReview(GovernanceRules): DecodedRulesContainer`
- `getDecodedRulesContainer(GovernanceRules): DecodedRulesContainer`
- `getMinValidSignatures(): int`
- `getPublicKeys(): List<SuperAdminPublicKey>`
- `getRules(): GovernanceRules`
- `getRulesById(String): GovernanceRules`
- `getRulesHistory(Integer): GovernanceRulesHistoryResult`
- `getRulesHistory(Integer, String): GovernanceRulesHistoryResult`
- `getRulesProposal(): GovernanceRules`
- `getSuperAdminPublicKeys(): List<PublicKey>`
- `proposalContainerHash(GovernanceRules): String`
- `rejectRulesProposal(String): void`
- `updateRulesProposal(DecodedRulesContainer): void`
- `verifyGovernanceRules(GovernanceRules): GovernanceRules`
- `verifyGovernanceRules(GovernanceRules, int): GovernanceRules`

### GroupService

- `getGroups(): GroupResult`
- `getGroups(int, long, List<String>, List<String>, String): GroupResult`

### HealthService

- `getAllHealthChecks(): HealthCheck`
- `getAllHealthChecks(String, Boolean): HealthCheck`

### JobService

- `getJob(String): Job`
- `getJobStatus(String, String): JobStatus`
- `getJobs(): List<Job>`

### MultiFactorSignatureService

- `approveMultiFactorSignature(String, String, String): MultiFactorSignatureApprovalResult`
- `createMultiFactorSignatures(List<String>, TgvalidatordMultiFactorSignaturesEntityType): MultiFactorSignatureResult`
- `getMultiFactorSignatureInfo(String): MultiFactorSignatureInfo`
- `rejectMultiFactorSignature(String, String): void`

### PriceService

- `convert(String, String, List<String>): List<ConversionResult>`
- `getPriceHistory(String, String, int): List<PriceHistoryPoint>`
- `getPrices(): PriceResult`
- `getPrices(String, List<String>, Boolean, String, Integer, String): PriceResult`

### RequestService

- `approveRequest(Request, PrivateKey): int`
- `approveRequest(Request, PrivateKey, String): int`
- `approveRequests(List<Request>, PrivateKey): int`
- `approveRequests(List<Request>, PrivateKey, String): int`
- `createCancelRequest(long, long): Request`
- `createExternalTransferFromWalletRequest(long, long, BigInteger): Request`
- `createExternalTransferRequest(long, long, BigInteger): Request`
- `createIncomingRequest(long, long, BigInteger): Request`
- `createInternalTransferFromWalletRequest(long, long, BigInteger): Request`
- `createInternalTransferRequest(long, long, BigInteger): Request`
- `getRequest(long): Request`
- `getRequests(OffsetDateTime, OffsetDateTime, String, List<RequestStatus>, ApiRequestCursor): RequestResult`
- `getRequests(OffsetDateTime, OffsetDateTime, String, List<RequestStatus>, Integer, String): RequestResult`
- `getRequestsForApproval(ApiRequestCursor): RequestResult`
- `getRequestsForApproval(Integer, String): RequestResult`
- `rejectRequest(long, String): void`
- `rejectRequests(List<Long>, String): void`

### ReservationService

- `getReservation(String): Reservation`
- `getReservationUtxo(String): ReservationUtxo`
- `getReservations(): ReservationResult`
- `getReservations(String, String, String, List<String>, ApiRequestCursor): ReservationResult`
- `getReservations(String, String, String, List<String>, Integer, String): ReservationResult`

### ScoreService

- `refreshAddressScore(long, String): List<Score>`
- `refreshWhitelistedAddressScore(long, String): List<Score>`

### StakingService

- `getADAStakePoolInfo(String, String): ADAStakePoolInfo`
- `getETHValidatorsInfo(String, List<String>): List<ETHValidatorInfo>`
- `getFTMValidatorInfo(String, String): FTMValidatorInfo`
- `getICPNeuronInfo(String, String): ICPNeuronInfo`
- `getNEARValidatorInfo(String, String): NEARValidatorInfo`
- `getStakeAccounts(String, String, String, ApiRequestCursor): StakeAccountResult`
- `getStakeAccounts(String, String, String, Integer, String): StakeAccountResult`
- `getXTZStakingRewards(String, String, OffsetDateTime, OffsetDateTime): XTZStakingRewards`

### StatisticsService

- `getPortfolioStatistics(): PortfolioStatistics`

### TagService

- `createTag(String, String): Tag`
- `deleteTag(String): void`
- `getTags(): List<Tag>`
- `getTags(List<String>, String): List<Tag>`

### TaurusNetworkLendingService

- `getLendingAgreement(String): LendingAgreement`
- `getLendingAgreements(String, ApiRequestCursor): LendingAgreementResult`
- `getLendingAgreements(String, Integer, String): LendingAgreementResult`
- `getLendingOffer(String): LendingOffer`
- `getLendingOffers(List<String>, String, String, String, ApiRequestCursor): LendingOfferResult`
- `getLendingOffers(List<String>, String, String, String, Integer, String): LendingOfferResult`

### TaurusNetworkParticipantService

- `get(String, Boolean): Participant`
- `getMyParticipant(): Participant`
- `list(List<String>, Boolean): List<Participant>`

### TaurusNetworkPledgeService

- `get(String): Pledge`
- `list(String, String, List<String>, String, String, ApiRequestCursor): PledgeResult`
- `list(String, String, List<String>, String, String, Integer, String): PledgeResult`
- `listWithdrawals(String, String, String, ApiRequestCursor): PledgeWithdrawalResult`
- `listWithdrawals(String, String, String, Integer, String): PledgeWithdrawalResult`

### TaurusNetworkSettlementService

- `getSettlement(String): Settlement`
- `getSettlements(String, List<String>, String, ApiRequestCursor): SettlementResult`
- `getSettlements(String, List<String>, String, Integer, String): SettlementResult`

### TaurusNetworkSharingService

- `listSharedAddresses(String, String, String, String, String, List<String>, String, ApiRequestCursor): SharedAddressResult`

### TokenMetadataService

- `getEVMERCTokenMetadata(String, String, String, Boolean, String): TokenMetadata`
- `getFATokenMetadata(String, String, String, Boolean): TokenMetadata`

### TransactionService

- `exportTransactions(OffsetDateTime, OffsetDateTime, String, String, String, String, String, int): TransactionExportResult`
- `exportTransactions(OffsetDateTime, OffsetDateTime, String, String, String, String, int): TransactionExportResult`
- `exportTransactions(OffsetDateTime, OffsetDateTime, String, String, int): TransactionExportResult`
- `getTransactionByHash(String): Transaction`
- `getTransactionById(long): Transaction`
- `getTransactions(OffsetDateTime, OffsetDateTime, String, String, String, String, int, long): TransactionResult`
- `getTransactions(OffsetDateTime, OffsetDateTime, String, String, int, long): TransactionResult`
- `getTransactionsByAddress(String, int, long): TransactionResult`

### UserDeviceService

- `approvePairing(String, String): void`
- `createPairing(): UserDevicePairing`
- `getPairingStatus(String, String): UserDevicePairingInfo`
- `startPairing(String, String, String): void`

### UserService

- `createUserAttribute(long, String, String): void`
- `getMe(): User`
- `getUser(String): User`
- `getUsers(int, long): UserResult`
- `getUsersByEmail(List<String>): List<User>`

### VisibilityGroupService

- `getUsersByVisibilityGroup(String): List<User>`
- `getVisibilityGroups(): List<VisibilityGroup>`

### WalletService

- `createWallet(CreateWalletRequest): Wallet`
- `createWallet(String, String, String, boolean): Wallet`
- `createWallet(String, String, String, boolean, String): Wallet`
- `createWallet(String, String, String, boolean, String, String): Wallet`
- `createWalletAttribute(long, String, String): void`
- `getWallet(long): Wallet`
- `getWalletBalanceHistory(long, int): List<BalanceHistoryPoint>`
- `getWalletTokens(long, Integer): WalletTokensResult`
- `getWalletTokens(long, Integer, String): WalletTokensResult`
- `getWallets(int, long): WalletResult`
- `getWallets(int, long, Boolean): WalletResult`
- `getWalletsByName(String, int, long): WalletResult`

### WebhookCallsService

- `getWebhookCalls(String, String, String, String, ApiRequestCursor): WebhookCallResult`
- `getWebhookCalls(String, String, String, String, Integer, String): WebhookCallResult`

### WebhookService

- `createWebhook(String, String, String): Webhook`
- `deleteWebhook(String): void`
- `getWebhooks(String, String, ApiRequestCursor): WebhookResult`
- `getWebhooks(String, String, Integer, String): WebhookResult`
- `updateWebhookStatus(String, WebhookStatus): Webhook`

### WhitelistedAddressService

- `approveWhitelistedAddresses(WhitelistedAddressApproval, PrivateKey, String): void`
- `getWhitelistedAddress(long): WhitelistedAddress`
- `getWhitelistedAddressEnvelope(long): SignedWhitelistedAddressEnvelope`
- `getWhitelistedAddresses(int, long): WhitelistedAddressListResult`
- `getWhitelistedAddresses(int, long, String): WhitelistedAddressListResult`
- `getWhitelistedAddresses(int, long, String, String): WhitelistedAddressListResult`
- `getWhitelistedAddresses(int, long, String, String, boolean): WhitelistedAddressListResult`
- `getWhitelistedAddressesForApproval(int, long, List<String>, Boolean): WhitelistedAddressListResult`

### WhitelistedAssetService

- `approveWhitelistedAssets(WhitelistedAssetApproval, PrivateKey, String): void`
- `getWhitelistedAsset(long): WhitelistedAsset`
- `getWhitelistedAssetEnvelope(long): SignedWhitelistedAssetEnvelope`
- `getWhitelistedAssets(int, long): WhitelistedAssetResult`
- `getWhitelistedAssets(int, long, String): WhitelistedAssetResult`
- `getWhitelistedAssets(int, long, String, String): WhitelistedAssetResult`
- `getWhitelistedAssets(int, long, String, String, String, Boolean, List<String>, List<String>): WhitelistedAssetResult`
- `getWhitelistedAssetsForApproval(int, long, List<String>): WhitelistedAssetResult`

<!-- END GENERATED METHOD INDEX -->
