# Services Reference

This document provides comprehensive documentation for all 43 services in the Taurus-PROTECT Java SDK.

## Service Overview

The SDK provides services organized into two categories: core services (38) and TaurusNetwork services (5).

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

Lists wallets with pagination.

```java
List<Wallet> getWallets(int limit, int offset) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| limit | int | Maximum results per page |
| offset | int | Pagination offset |

#### getWalletsByName

Searches wallets by name.

```java
List<Wallet> getWalletsByName(String name, int limit, int offset) throws ApiException
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

Lists token balances for a wallet.

```java
List<AssetBalance> getWalletTokens(long walletId, int limit) throws ApiException
```

### Key Models

- `Wallet` - id, name, blockchain, network, balance, isOmnibus, customerId, attributes
- `BalanceHistoryPoint` - timestamp, balance values
- `AssetBalance` - asset info with balance

---

## AddressService

**Purpose:** Manages blockchain addresses within wallets with signature verification.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/AddressService.java`

### Methods

#### createAddress

Creates a new address in a wallet.

```java
Address createAddress(long walletId, String label, String comment, String customerId) throws ApiException
```

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

Lists addresses for a wallet with **signature verification**.

```java
List<Address> getAddresses(long walletId, int limit, int offset) throws ApiException
```

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

Lists requests with filtering and pagination.

```java
RequestResult getRequests(OffsetDateTime from, OffsetDateTime to, String currencyId,
                          List<RequestStatus> statuses, ApiRequestCursor cursor) throws ApiException
```

#### getRequestsForApproval

Gets requests pending approval for the current user.

```java
RequestResult getRequestsForApproval(ApiRequestCursor cursor) throws ApiException
```

#### approveRequest / approveRequests

Signs and approves requests using a private key.

```java
int approveRequest(Request request, PrivateKey privateKey) throws ApiException
int approveRequests(List<Request> requests, PrivateKey privateKey) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| request(s) | Request / List<Request> | Request(s) to approve |
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

Lists transactions with filtering.

```java
List<Transaction> getTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                                   String direction, int limit, int offset) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| from | OffsetDateTime | Start date (optional) |
| to | OffsetDateTime | End date (optional) |
| currency | String | Currency filter (optional) |
| direction | String | "incoming" or "outgoing" (optional) |

#### getTransactionsByAddress

```java
List<Transaction> getTransactionsByAddress(String address, int limit, int offset) throws ApiException
```

#### exportTransactions

Exports transactions as CSV.

```java
String exportTransactions(OffsetDateTime from, OffsetDateTime to, String currency,
                          String direction, int limit, int offset) throws ApiException
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
BalanceResult getBalances(ApiRequestCursor cursor) throws ApiException
BalanceResult getBalances(String currency, ApiRequestCursor cursor) throws ApiException
```

**Example:**
```java
ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);
BalanceResult result;
do {
    result = client.getBalanceService().getBalances(cursor);
    for (AssetBalance balance : result.getBalances()) {
        System.out.println(balance.getAsset() + ": " + balance.getBalance());
    }
    cursor = result.nextCursor(100);
} while (result.hasNext());
```

#### getNFTCollectionBalances

```java
NFTCollectionBalanceResult getNFTCollectionBalances(String blockchain, String network,
                                                     ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `BalanceResult` - balances list with pagination
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

### Methods

```java
List<Price> getPrices() throws ApiException
List<PriceHistoryPoint> getPriceHistory(String base, String quote, int limit) throws ApiException
List<ConversionResult> convert(String currency, BigDecimal amount, List<String> targetCurrencyIds) throws ApiException
```

### Key Models

- `Price` - currency pair, price, timestamp
- `PriceHistoryPoint` - timestamp, price
- `ConversionResult` - target currency, converted amount

---

## UserService

**Purpose:** Retrieves and manages user information.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/UserService.java`

### Methods

```java
User getMe() throws ApiException
List<User> getUsers(int limit, int offset) throws ApiException
List<User> getUsersByEmail(List<String> emails) throws ApiException
void createUserAttribute(String userId, String key, String value) throws ApiException
```

### Key Models

- `User` - id, email, name, roles, attributes

---

## ChangeService

**Purpose:** Manages configuration changes and approval workflows.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/ChangeService.java`

### Methods

```java
Change getChange(long id) throws ApiException
ChangeResult getChanges(String entity, String status, ApiRequestCursor cursor) throws ApiException
ChangeResult getChangesForApproval(ApiRequestCursor cursor) throws ApiException
void approveChange(long id) throws ApiException
void approveChanges(List<Long> ids) throws ApiException
void rejectChange(long id) throws ApiException
void rejectChanges(List<Long> ids) throws ApiException
```

### Key Models

- `Change` - id, entity, operation, status, payload, trails
- `ChangeResult` - changes list with pagination

---

## BusinessRuleService

**Purpose:** Manages transaction approval rules.

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/BusinessRuleService.java`

### Methods

```java
BusinessRuleResult getBusinessRules(ApiRequestCursor cursor) throws ApiException
BusinessRuleResult getBusinessRulesByWallet(long walletId, ApiRequestCursor cursor) throws ApiException
BusinessRuleResult getBusinessRulesByCurrency(String currencyId, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `BusinessRule` - rule definition with conditions and actions
- `BusinessRuleResult` - rules list with pagination

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

```java
GovernanceRulesHistoryResult getRulesHistory(int pageSize) throws ApiException
GovernanceRulesHistoryResult getRulesHistory(int pageSize, byte[] cursor) throws ApiException
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

### Key Models

- `GovernanceRules` - rulesContainer, rulesSignatures, locked, trails
- `DecodedRulesContainer` - groups, users, thresholds, addressWhitelistingRules
- `SuperAdminPublicKey` - id, publicKey, name

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

Lists whitelisted addresses with filtering.

```java
List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(int limit, int offset)
List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(int limit, int offset, String blockchain)
List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(int limit, int offset, String blockchain, String network)
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
WebhookResult getWebhooks(String type, String url, ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| type | String | Filter by webhook type (optional) |
| url | String | Filter by URL (optional) |
| cursor | ApiRequestCursor | Pagination cursor (optional) |

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

**Purpose:** Manages whitelisted smart contract addresses (ERC20 tokens, NFTs, FA2 tokens).

**Location:** `client/src/main/java/com/taurushq/sdk/protect/client/service/ContractWhitelistingService.java`

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
void approveWhitelistedContracts(List<String> ids, String signature, String comment) throws ApiException
```

#### getWhitelistedContract

Retrieves a single whitelisted contract by ID.

```java
SignedWhitelistedContractAddressEnvelope getWhitelistedContract(String id) throws ApiException
```

#### getWhitelistedContracts

Lists whitelisted contracts with filtering and pagination.

```java
WhitelistedContractAddressResult getWhitelistedContracts(String blockchain, String network,
                                                          String query, Boolean isNFT,
                                                          Integer limit, Integer offset) throws ApiException
```

**Example:**
```java
WhitelistedContractAddressResult result = client.getContractWhitelistingService()
    .getWhitelistedContracts("ETH", "mainnet", null, false, 50, 0);
for (SignedWhitelistedContractAddressEnvelope contract : result.getContracts()) {
    System.out.println(contract.getBlockchain() + ": " + contract.getId());
}
```

#### getWhitelistedContractsForApproval

Lists whitelisted contracts pending approval.

```java
WhitelistedContractAddressResult getWhitelistedContractsForApproval(List<String> ids,
                                                                      Integer limit, Integer offset) throws ApiException
```

#### updateWhitelistedContract

Updates an existing whitelisted contract.

```java
void updateWhitelistedContract(String id, String symbol, String name, int decimals) throws ApiException
```

#### deleteWhitelistedContract

Deletes a whitelisted contract.

```java
String deleteWhitelistedContract(String id, String comment) throws ApiException
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

- `SignedWhitelistedContractAddressEnvelope` - id, blockchain, network, status, signedContractAddress, metadata, approvers, trails
- `WhitelistedContractAddressResult` - contracts list with totalItems and pagination helpers
- `Attribute` - key, value, contentType, type, subType

---

## WhitelistedAssetService

**Purpose:** Manages whitelisted assets/contracts with cryptographic verification.

**Access:** `client.getWhitelistedAssetService()`

### Methods

#### getWhitelistedAsset

Gets a whitelisted asset by ID with verification.

```java
WhitelistedAsset getWhitelistedAsset(String id) throws ApiException
```

#### getWhitelistedAssets

Lists whitelisted assets with filtering.

```java
WhitelistedAssetResult getWhitelistedAssets(String blockchain, String network, int limit, int offset) throws ApiException
```

### Key Models

- `WhitelistedAsset` - id, blockchain, network, status, metadata, signedContractAddress

---

## AuditService

**Purpose:** Queries audit trail events.

**Access:** `client.getAuditService()`

### Methods

#### getAudits

Lists audit events with filtering.

```java
AuditResult getAudits(String entity, String action, OffsetDateTime from, OffsetDateTime to, ApiRequestCursor cursor) throws ApiException
```

### Key Models

- `Audit` - id, entity, action, user, timestamp, details

---

## FeeService

**Purpose:** Retrieves transaction fee information.

**Access:** `client.getFeeService()`

### Methods

#### getFees

Gets fee information for a currency.

```java
List<Fee> getFees(String currency) throws ApiException
```

### Key Models

- `Fee` - currency, feeType, amount, unit

---

## AirGapService

**Purpose:** Provides air-gap signing operations for offline transaction signing.

**Access:** `client.getAirGapService()`

### Methods

#### getAirGapRequest

Gets an air-gap request for offline signing.

```java
AirGapRequest getAirGapRequest(long requestId) throws ApiException
```

#### submitAirGapSignature

Submits a signature for an air-gap request.

```java
void submitAirGapSignature(long requestId, String signature) throws ApiException
```

---

## ReservationService

**Purpose:** Manages balance reservations for addresses.

**Access:** `client.getReservationService()`

### Methods

#### createReservation

Creates a balance reservation.

```java
Reservation createReservation(long addressId, String amount, String comment) throws ApiException
```

#### getReservations

Lists reservations.

```java
ReservationResult getReservations(long addressId, ApiRequestCursor cursor) throws ApiException
```

#### cancelReservation

Cancels a reservation.

```java
void cancelReservation(long reservationId) throws ApiException
```

### Key Models

- `Reservation` - id, addressId, amount, status, expiresAt

---

## MultiFactorSignatureService

**Purpose:** Manages multi-factor signature operations for enhanced security.

**Access:** `client.getMultiFactorSignatureService()`

### Methods

#### getMultiFactorSignatures

Lists pending multi-factor signature requests.

```java
List<MultiFactorSignature> getMultiFactorSignatures(ApiRequestCursor cursor) throws ApiException
```

#### approveMultiFactorSignature

Approves a multi-factor signature request.

```java
void approveMultiFactorSignature(String id, String signature) throws ApiException
```

---

## GroupService

**Purpose:** Manages user groups for approval workflows.

**Access:** `client.getGroupService()`

### Methods

#### getGroups

Lists user groups.

```java
List<Group> getGroups(int limit, int offset) throws ApiException
```

#### getGroup

Gets a group by ID.

```java
Group getGroup(String groupId) throws ApiException
```

### Key Models

- `Group` - id, name, members, threshold

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

#### getVisibilityGroup

Gets a visibility group by ID.

```java
VisibilityGroup getVisibilityGroup(String id) throws ApiException
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
WebhookCallResult getWebhookCalls(String webhookId, String status, ApiRequestCursor cursor) throws ApiException
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

### Methods

#### getAssets

Lists assets.

```java
List<Asset> getAssets(String blockchain, String network, int limit, int offset) throws ApiException
```

### Key Models

- `Asset` - id, symbol, name, blockchain, network, contractAddress, decimals

---

## ActionService

**Purpose:** Manages actions in the system.

**Access:** `client.getActionService()`

### Methods

#### getActions

Lists actions.

```java
List<Action> getActions(ApiRequestCursor cursor) throws ApiException
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

## ExchangeService

**Purpose:** Manages exchange integrations.

**Access:** `client.getExchangeService()`

### Methods

#### getExchanges

Lists configured exchanges.

```java
List<Exchange> getExchanges() throws ApiException
```

### Key Models

- `Exchange` - id, name, type, status

---

## FiatService

**Purpose:** Manages fiat currency operations.

**Access:** `client.getFiatService()`

### Methods

#### getFiatCurrencies

Lists supported fiat currencies.

```java
List<FiatCurrency> getFiatCurrencies() throws ApiException
```

---

## FeePayerService

**Purpose:** Manages fee payer configurations.

**Access:** `client.getFeePayerService()`

### Methods

#### getFeePayers

Lists fee payers.

```java
List<FeePayer> getFeePayers(String blockchain, String network) throws ApiException
```

---

## HealthService

**Purpose:** Checks system health status.

**Access:** `client.getHealthService()`

### Methods

#### check

Checks API health.

```java
HealthStatus check() throws ApiException
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

#### getTokenMetadata

Gets metadata for a token.

```java
TokenMetadata getTokenMetadata(String blockchain, String network, String contractAddress) throws ApiException
```

### Key Models

- `TokenMetadata` - name, symbol, decimals, totalSupply, logoUrl

---

## UserDeviceService

**Purpose:** Manages user device registrations.

**Access:** `client.getUserDeviceService()`

### Methods

#### getUserDevices

Lists devices for a user.

```java
List<UserDevice> getUserDevices(String userId) throws ApiException
```

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

#### getParticipant

Retrieves a participant by ID.

```java
Participant getParticipant(String participantId, Boolean includeTotalPledgesValuation) throws ApiException
```

#### getParticipants

Retrieves multiple participants by IDs.

```java
List<Participant> getParticipants(List<String> participantIds, Boolean includeTotalPledgesValuation) throws ApiException
```

### Key Models

- `Participant` - id, name, country, publicKey, totalPledgesValuation

---

## TaurusNetworkPledgeService

**Purpose:** Provides access to Taurus Network pledge lifecycle operations.

**Access:** `client.taurusNetwork().pledges()`

### Methods

#### getPledge

Retrieves a pledge by ID.

```java
Pledge getPledge(String pledgeId) throws ApiException
```

#### getPledges

Retrieves pledges with optional filtering.

```java
PledgeResult getPledges(String ownerParticipantId, String targetParticipantId,
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
| cursor | ApiRequestCursor | Pagination cursor (optional) |

#### getPledgeWithdrawals

Retrieves pledge withdrawals for a specific pledge.

```java
PledgeWithdrawalResult getPledgeWithdrawals(String pledgeId, String withdrawalStatus,
                                             String sortOrder, ApiRequestCursor cursor) throws ApiException
```

#### approvePledgeActions

Approves pledge actions with ECDSA signature.

```java
int approvePledgeActions(List<PledgeAction> actions, PrivateKey privateKey) throws ApiException
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
                                     ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| currencyIds | List<String> | Filter by currency IDs (optional) |
| participantId | String | Filter by participant ID (optional) |
| duration | String | Filter by duration (optional) |
| sortOrder | String | Sort order: "ASC" or "DESC" (optional) |
| cursor | ApiRequestCursor | Pagination cursor (optional) |

**Example:**
```java
LendingOfferResult result = client.taurusNetwork().lending()
    .getLendingOffers(Arrays.asList("ETH"), null, null, "DESC", null);
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
                                 String sortOrder, ApiRequestCursor cursor) throws ApiException
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| counterParticipantId | String | Filter by counter participant ID (optional) |
| statuses | List<String> | Filter by statuses (optional) |
| sortOrder | String | Sort order: "ASC" or "DESC" (optional) |
| cursor | ApiRequestCursor | Pagination cursor (optional) |

**Example:**
```java
SettlementResult result = client.taurusNetwork().settlements()
    .getSettlements(null, Arrays.asList("PENDING", "COMPLETED"), "DESC", null);
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

#### getSharedAddresses

Retrieves shared addresses with optional filtering.

```java
SharedAddressResult getSharedAddresses(String participantId, String ownerParticipantId,
                                        String targetParticipantId, String blockchain,
                                        String network, List<String> ids,
                                        String sortOrder, ApiRequestCursor cursor) throws ApiException
```

#### createSharedAddress

Creates a shared address.

```java
SharedAddress createSharedAddress(String internalAddressId, String targetParticipantId,
                                   List<String> permissions) throws ApiException
```

#### revokeSharedAddress

Revokes a shared address.

```java
void revokeSharedAddress(String sharedAddressId) throws ApiException
```

#### getSharedAssets

Retrieves shared assets with optional filtering.

```java
SharedAssetResult getSharedAssets(String participantId, ApiRequestCursor cursor) throws ApiException
```

#### createSharedAsset

Creates a shared asset.

```java
SharedAsset createSharedAsset(String assetId, String targetParticipantId,
                               List<String> permissions) throws ApiException
```

#### revokeSharedAsset

Revokes a shared asset.

```java
void revokeSharedAsset(String sharedAssetId) throws ApiException
```

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
