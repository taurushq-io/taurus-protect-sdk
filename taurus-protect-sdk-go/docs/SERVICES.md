# Services Reference

This document provides comprehensive documentation for all 43 services in the Taurus-PROTECT Go SDK.

## Service Overview

The SDK provides services organized into two categories: core services (38) and TaurusNetwork services (5).

### Core Services

| Service | Access | Purpose |
|---------|--------|---------|
| [WalletService](#walletservice) | `client.Wallets()` | Create and manage blockchain wallets |
| [AddressService](#addressservice) | `client.Addresses()` | Create and manage addresses within wallets |
| [RequestService](#requestservice) | `client.Requests()` | Transaction requests with approval workflow |
| [TransactionService](#transactionservice) | `client.Transactions()` | Query blockchain transactions |
| [BalanceService](#balanceservice) | `client.Balances()` | Query asset balances |
| [CurrencyService](#currencyservice) | `client.Currencies()` | Currency metadata |
| [GovernanceRuleService](#governanceruleservice) | `client.GovernanceRules()` | Governance rules with signature verification |
| [WhitelistedAddressService](#whitelistedaddressservice) | `client.WhitelistedAddresses()` | External address whitelisting |
| [WhitelistedAssetService](#whitelistedassetservice) | `client.WhitelistedAssets()` | Token/contract whitelisting |
| [WhitelistedContractService](#whitelistedcontractservice) | `client.WhitelistedContracts()` | Smart contract whitelisting |
| [AuditService](#auditservice) | `client.Audits()` | Audit trail querying |
| [ChangeService](#changeservice) | `client.Changes()` | Configuration changes |
| [FeeService](#feeservice) | `client.Fees()` | Transaction fee information |
| [PriceService](#priceservice) | `client.Prices()` | Price data and conversion |
| [AirGapService](#airgapservice) | `client.AirGap()` | Air-gap signing operations |
| [StakingService](#stakingservice) | `client.Staking()` | Staking operations |
| [BusinessRuleService](#businessruleservice) | `client.BusinessRules()` | Business rules |
| [ReservationService](#reservationservice) | `client.Reservations()` | Balance reservations |
| [MultiFactorSignatureService](#multifactorsignatureservice) | `client.MultiFactorSignature()` | Multi-factor signatures |
| [UserService](#userservice) | `client.Users()` | User management |
| [GroupService](#groupservice) | `client.Groups()` | User group management |
| [VisibilityGroupService](#visibilitygroupservice) | `client.VisibilityGroups()` | Visibility group management |
| [ConfigService](#configservice) | `client.Config()` | System configuration |
| [WebhookService](#webhookservice) | `client.Webhooks()` | Webhook configuration |
| [WebhookCallService](#webhookcallservice) | `client.WebhookCalls()` | Webhook call history |
| [TagService](#tagservice) | `client.Tags()` | Tag management |
| [AssetService](#assetservice) | `client.Assets()` | Asset information |
| [ActionService](#actionservice) | `client.Actions()` | Action management |
| [BlockchainService](#blockchainservice) | `client.Blockchains()` | Blockchain information |
| [ExchangeService](#exchangeservice) | `client.Exchanges()` | Exchange integration |
| [FiatService](#fiatservice) | `client.Fiat()` | Fiat currency operations |
| [FeePayerService](#feepayerservice) | `client.FeePayers()` | Fee payer management |
| [HealthService](#healthservice) | `client.Health()` | System health checks |
| [JobService](#jobservice) | `client.Jobs()` | Background job management |
| [ScoreService](#scoreservice) | `client.Scores()` | Address risk scores |
| [StatisticsService](#statisticsservice) | `client.Statistics()` | Platform statistics |
| [TokenMetadataService](#tokenmetadataservice) | `client.TokenMetadata()` | Token metadata |
| [UserDeviceService](#userdeviceservice) | `client.UserDevices()` | User device management |

### TaurusNetwork Services

| Service | Access | Purpose |
|---------|--------|---------|
| [TaurusNetworkParticipantService](#taurusnetworkparticipantservice) | `client.TaurusNetwork().Participants()` | Participant management |
| [TaurusNetworkPledgeService](#taurusnetworkpledgeservice) | `client.TaurusNetwork().Pledges()` | Pledge lifecycle operations |
| [TaurusNetworkLendingService](#taurusnetworklendingservice) | `client.TaurusNetwork().Lending()` | Lending offers and agreements |
| [TaurusNetworkSettlementService](#taurusnetworksettlementservice) | `client.TaurusNetwork().Settlements()` | Settlement operations |
| [TaurusNetworkSharingService](#taurusnetworksharingservice) | `client.TaurusNetwork().Sharing()` | Address and asset sharing |

---

## WalletService

**Purpose:** Creates and manages blockchain wallets with balance tracking.

**Access:** `client.Wallets()`

### Methods

#### GetWallet

Retrieves a wallet by ID.

```go
func (s *WalletService) GetWallet(ctx context.Context, walletID string) (*model.Wallet, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Request context |
| walletID | string | Wallet identifier |

**Returns:** `*model.Wallet`, `error`

**Example:**
```go
wallet, err := client.Wallets().GetWallet(ctx, "wallet-123")
if err != nil {
    return err
}
fmt.Printf("Wallet: %s, Balance: %s\n", wallet.Name, wallet.Balance.AvailableConfirmed)
```

#### ListWallets

Lists wallets with pagination and filtering.

```go
func (s *WalletService) ListWallets(ctx context.Context, opts *model.ListWalletsOptions) ([]*model.Wallet, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Limit | int64 | Maximum results per page |
| opts.Offset | int64 | Pagination offset |
| opts.Currency | string | Filter by currency symbol |
| opts.Query | string | Search wallet names |
| opts.ExcludeDisabled | bool | Exclude disabled wallets |

**Returns:** `[]*model.Wallet`, `*model.Pagination`, `error`

#### CreateWallet

Creates a new blockchain wallet.

```go
func (s *WalletService) CreateWallet(ctx context.Context, req *model.CreateWalletRequest) (*model.Wallet, error)
```

**Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| Name | string | Human-readable wallet name (required) |
| Currency | string | Currency symbol (required, e.g., "ETH", "BTC") |
| Comment | string | Optional description |
| CustomerID | string | Optional external customer ID |
| ExternalWalletID | string | Optional external identifier |
| VisibilityGroupID | string | Optional visibility group to assign |

**Returns:** `*model.Wallet`, `error`

**Example:**
```go
wallet, err := client.Wallets().CreateWallet(ctx, &model.CreateWalletRequest{
    Name:       "Treasury Wallet",
    Currency:   "ETH",
    CustomerID: "CUST-001",
})
```

### Key Models

- `model.Wallet` - ID, Name, Currency, Blockchain, Network, Balance, IsOmnibus, CustomerID, Attributes
- `model.ListWalletsOptions` - Limit, Offset, Currency, Query, ExcludeDisabled
- `model.CreateWalletRequest` - Name, Currency, Comment, CustomerID, ExternalWalletID, VisibilityGroupID

---

## AddressService

**Purpose:** Manages blockchain addresses within wallets.

**Access:** `client.Addresses()`

### Methods

#### GetAddress

Retrieves an address by ID.

```go
func (s *AddressService) GetAddress(ctx context.Context, addressID string) (*model.Address, error)
```

**Returns:** `*model.Address`, `error`

#### ListAddresses

Lists addresses for a wallet with pagination.

```go
func (s *AddressService) ListAddresses(ctx context.Context, walletID string, opts *model.ListAddressesOptions) ([]*model.Address, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| walletID | string | Parent wallet ID |
| opts.Limit | int64 | Maximum results per page |
| opts.Offset | int64 | Pagination offset |

#### CreateAddress

Creates a new address in a wallet.

```go
func (s *AddressService) CreateAddress(ctx context.Context, walletID string, req *model.CreateAddressRequest) (*model.Address, error)
```

**Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| Label | string | Address label |
| Comment | string | Optional comment |
| CustomerID | string | External customer ID |

**Example:**
```go
address, err := client.Addresses().CreateAddress(ctx, walletID, &model.CreateAddressRequest{
    Label:      "Customer Deposit",
    Comment:    "Auto-generated",
    CustomerID: "USER-789",
})
fmt.Printf("Address: %s\n", address.Address)
```

#### DeleteAddress

Deletes an address by ID.

```go
func (s *AddressService) DeleteAddress(ctx context.Context, addressID string) error
```

#### UpdateAddress

Updates address properties.

```go
func (s *AddressService) UpdateAddress(ctx context.Context, addressID string, req *model.UpdateAddressRequest) (*model.Address, error)
```

### Key Models

- `model.Address` - ID, WalletID, Address, Label, CustomerID, Balance, Status, Attributes

---

## RequestService

**Purpose:** Creates, approves, and manages transaction requests with cryptographic signing.

**Access:** `client.Requests()`

### Methods

#### GetRequest

Retrieves a request by ID.

```go
func (s *RequestService) GetRequest(ctx context.Context, requestID string) (*model.Request, error)
```

**Returns:** `*model.Request` with `Metadata` containing the transaction details and hash.

#### ListRequests

Lists requests with filtering and pagination.

```go
func (s *RequestService) ListRequests(ctx context.Context, opts *model.ListRequestsOptions) ([]*model.Request, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Limit | int64 | Maximum results per page |
| opts.Offset | int64 | Pagination offset |
| opts.Status | string | Filter by status |
| opts.Currency | string | Filter by currency |
| opts.From | time.Time | Start date filter |
| opts.To | time.Time | End date filter |

#### CreateOutgoingRequest

Creates an external transfer request.

```go
func (s *RequestService) CreateOutgoingRequest(ctx context.Context, req *model.CreateOutgoingRequest) (*model.Request, error)
```

**Parameters:**
| Field | Type | Description |
|-------|------|-------------|
| Amount | string | Amount to transfer in smallest unit (required) |
| FromAddressID | string | Source address ID (either this or FromWalletID required) |
| FromWalletID | string | Source wallet ID for omnibus wallets |
| ToAddressID | string | Destination address ID (either this or ToWhitelistedAddressID required) |
| ToWhitelistedAddressID | string | Destination whitelisted address ID |
| FeeLimit | string | Maximum fee amount |
| GasLimit | string | Maximum gas for the transaction |
| Comment | string | Reconciliation note |
| ExternalRequestID | string | Optional external identifier |

**Example:**
```go
request, err := client.Requests().CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
    FromAddressID:          "addr-123",
    ToWhitelistedAddressID: "wla-456",
    Amount:                 "1000000000000000000", // 1 ETH in wei
})
fmt.Printf("Request ID: %s, Status: %s\n", request.ID, request.Status)
```

#### ApproveRequest

Approves a single request with ECDSA signing. The SDK handles hash extraction and signing internally.

```go
func (s *RequestService) ApproveRequest(ctx context.Context, request *model.Request, privateKey *ecdsa.PrivateKey) (int, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Request context |
| request | *model.Request | Request to approve (must have Metadata.Hash) |
| privateKey | *ecdsa.PrivateKey | User's ECDSA private key for signing |

**Returns:** `int` (number of requests signed), `error`

**Example:**
```go
// Get request to approve
request, err := client.Requests().GetRequest(ctx, requestID)
if err != nil {
    return err
}

// Approve with private key - SDK handles signing
signedCount, err := client.Requests().ApproveRequest(ctx, request, privateKey)
if err != nil {
    return err
}
fmt.Printf("Approved %d request(s)\n", signedCount)
```

#### ApproveRequests

Approves multiple requests with ECDSA signing. Requests are sorted by ID before signing.

```go
func (s *RequestService) ApproveRequests(ctx context.Context, requests []*model.Request, privateKey *ecdsa.PrivateKey) (int, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Request context |
| requests | []*model.Request | Requests to approve (must have Metadata.Hash) |
| privateKey | *ecdsa.PrivateKey | User's ECDSA private key for signing |

**Returns:** `int` (number of requests signed), `error`

#### RejectRequest

Rejects a request with a comment.

```go
func (s *RequestService) RejectRequest(ctx context.Context, requestID string, comment string) error
```

### Key Models

- `model.Request` - ID, Status, Currency, Type, Metadata, NeedsApprovalFrom
- `model.RequestMetadata` - Hash, PayloadAsString (use `ParsePayloadEntries()` for structured access)
- `model.RequestStatus` - CREATED, APPROVING, HSM_SIGNED, BROADCASTING, BROADCASTED, CONFIRMED, REJECTED, FAILED

---

## TransactionService

**Purpose:** Retrieves and analyzes blockchain transactions.

**Access:** `client.Transactions()`

### Methods

#### ListTransactions

Lists transactions with filtering.

```go
func (s *TransactionService) ListTransactions(ctx context.Context, opts *model.ListTransactionsOptions) ([]*model.Transaction, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Limit | int64 | Maximum results per page |
| opts.Offset | int64 | Pagination offset |
| opts.Direction | string | "incoming" or "outgoing" |
| opts.Currency | string | Currency filter |
| opts.From | time.Time | Start date |
| opts.To | time.Time | End date |

**Example:**
```go
txs, pagination, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
    Limit:     100,
    Direction: "outgoing",
    Currency:  "ETH",
})
for _, tx := range txs {
    fmt.Printf("%s: %s -> %s\n", tx.Hash, tx.Amount, tx.Direction)
}
```

### Key Models

- `model.Transaction` - ID, Hash, Direction, Type, Currency, Amount, Fee, Block, Sources, Destinations

---

## BalanceService

**Purpose:** Retrieves total asset balances across the tenant.

**Access:** `client.Balances()`

### Methods

#### GetBalances

Gets total balances for the tenant, grouped by asset. An asset is identified by a full triplet of attributes: blockchain, contract number, and token ID.

```go
func (s *BalanceService) GetBalances(ctx context.Context, opts *model.GetBalancesOptions) (*model.GetBalancesResult, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Currency | string | Filter by currency ID or symbol |
| opts.TokenID | string | Filter by token ID |
| opts.Limit | int64 | Maximum results per page |
| opts.Cursor | string | Pagination cursor for next page |

**Returns:** `*model.GetBalancesResult`, `error`

**Example:**
```go
result, err := client.Balances().GetBalances(ctx, &model.GetBalancesOptions{
    Currency: "ETH",
    Limit:    100,
})
if err != nil {
    return err
}

for _, balance := range result.Balances {
    fmt.Printf("%s: %s available\n", balance.Asset.Currency, balance.Balance.AvailableConfirmed)
}

// Paginate with cursor
if result.NextCursor != "" {
    nextResult, err := client.Balances().GetBalances(ctx, &model.GetBalancesOptions{
        Cursor: result.NextCursor,
    })
}
```

### Key Models

- `model.GetBalancesOptions` - Currency, TokenID, Limit, Cursor
- `model.GetBalancesResult` - Balances, Total, NextCursor
- `model.AssetBalance` - Asset, Balance
- `model.Balance` - TotalConfirmed, TotalUnconfirmed, AvailableConfirmed, AvailableUnconfirmed, ReservedConfirmed, ReservedUnconfirmed

---

## CurrencyService

**Purpose:** Manages and retrieves currency metadata.

**Access:** `client.Currencies()`

### Methods

```go
func (s *CurrencyService) GetCurrencies(ctx context.Context) ([]*model.Currency, error)
```

### Key Models

- `model.Currency` - ID, Name, Symbol, Blockchain, Network, Decimals

---

## GovernanceRuleService

**Purpose:** Manages governance rules with SuperAdmin signature verification.

**Access:** `client.GovernanceRules()`

### Methods

#### GetRules

Gets current governance rules.

```go
func (s *GovernanceRuleService) GetRules(ctx context.Context) (*model.GovernanceRuleset, error)
```

**Returns:** `*model.GovernanceRuleset` with rulesContainer and signatures.

#### GetRulesByID

Gets governance rules by version ID.

```go
func (s *GovernanceRuleService) GetRulesByID(ctx context.Context, rulesID string) (*model.GovernanceRuleset, error)
```

#### GetRulesHistory

Gets historical governance rules with cursor-based pagination.

```go
func (s *GovernanceRuleService) GetRulesHistory(ctx context.Context, opts *model.ListRulesHistoryOptions) (*model.GovernanceRulesHistoryResult, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Limit | int64 | Maximum results per page |
| opts.Cursor | string | Pagination cursor from previous request |

**Returns:** `*model.GovernanceRulesHistoryResult` with `Rules`, `TotalItems`, `Cursor`

#### GetRulesProposal

Gets pending rules proposal (SuperAdmin only).

```go
func (s *GovernanceRuleService) GetRulesProposal(ctx context.Context) (*model.GovernanceRuleset, error)
```

#### GetPublicKeys

Lists SuperAdmin public keys.

```go
func (s *GovernanceRuleService) GetPublicKeys(ctx context.Context) ([]*model.SuperAdminPublicKey, error)
```

#### VerifyGovernanceRules

Verifies governance rules signatures against configured SuperAdmin keys. Returns the decoded rules container if verification passes.

```go
func (s *GovernanceRuleService) VerifyGovernanceRules(ctx context.Context, rules *model.GovernanceRuleset) (*model.DecodedRulesContainer, error)
```

### Key Models

- `model.GovernanceRuleset` - RulesContainer, Signatures, Locked, Trails
- `model.RuleUserSignature` - UserID, Signature
- `model.SuperAdminPublicKey` - ID, PublicKey, Name

---

## WhitelistedAddressService

**Purpose:** Manages whitelisted addresses with cryptographic verification.

**Access:** `client.WhitelistedAddresses()`

### Methods

#### GetWhitelistedAddress

Gets a whitelisted address by ID.

```go
func (s *WhitelistedAddressService) GetWhitelistedAddress(ctx context.Context, id string) (*model.WhitelistedAddress, error)
```

**Returns:** `*model.WhitelistedAddress` with verification data.

#### ListWhitelistedAddresses

Lists whitelisted addresses with filtering.

```go
func (s *WhitelistedAddressService) ListWhitelistedAddresses(ctx context.Context, opts *model.ListWhitelistedAddressesOptions) ([]*model.WhitelistedAddress, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Limit | int64 | Maximum results per page |
| opts.Offset | int64 | Pagination offset |
| opts.Blockchain | string | Filter by blockchain |
| opts.Network | string | Filter by network |

#### CreateWhitelistedAddress

Creates a new whitelisted address.

```go
func (s *WhitelistedAddressService) CreateWhitelistedAddress(ctx context.Context, req *model.CreateWhitelistedAddressRequest) (*model.WhitelistedAddress, error)
```

#### ApproveWhitelistedAddress

Approves a whitelisted address with signatures.

```go
func (s *WhitelistedAddressService) ApproveWhitelistedAddress(ctx context.Context, id string, signature string, comment string) error
```

### Key Models

- `model.WhitelistedAddress` - ID, Blockchain, Network, Address, Name, Status, RulesContainer, SignedAddress, Approvers
- `model.SignedWhitelistedAddress` - Payload, Signatures
- `model.WhitelistSignature` - UserSignature, Hashes

---

## WhitelistedAssetService

**Purpose:** Manages whitelisted tokens and contracts.

**Access:** `client.WhitelistedAssets()`

### Methods

```go
func (s *WhitelistedAssetService) GetWhitelistedAsset(ctx context.Context, id string) (*model.WhitelistedAsset, error)
func (s *WhitelistedAssetService) GetWhitelistedAssetEnvelope(ctx context.Context, id string) (*model.WhitelistedAssetEnvelope, error)
func (s *WhitelistedAssetService) ListWhitelistedAssets(ctx context.Context, opts *model.ListWhitelistedAssetsOptions) ([]*model.WhitelistedAsset, *model.Pagination, error)
func (s *WhitelistedAssetService) CreateWhitelistedAsset(ctx context.Context, req *model.CreateWhitelistedAssetRequest) (*model.WhitelistedAsset, error)
```

### Key Models

- `model.WhitelistedAsset` - ID, Blockchain, Network, Status, Metadata, SignedContractAddress, Approvers

---

## UserService

**Purpose:** Retrieves and manages user information.

**Access:** `client.Users()`

### Methods

```go
func (s *UserService) GetMe(ctx context.Context) (*model.User, error)
func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error)
func (s *UserService) ListUsers(ctx context.Context, opts *model.ListUsersOptions) ([]*model.User, *model.Pagination, error)
```

### Key Models

- `model.User` - ID, Email, Name, Roles, Attributes

---

## GroupService

**Purpose:** Manages user groups for approval workflows.

**Access:** `client.Groups()`

### Methods

```go
func (s *GroupService) GetGroup(ctx context.Context, groupID string) (*model.Group, error)
func (s *GroupService) ListGroups(ctx context.Context, opts *model.ListGroupsOptions) ([]*model.Group, *model.Pagination, error)
```

### Key Models

- `model.Group` - ID, Name, Members, Threshold

---

## AuditService

**Purpose:** Queries audit trail events.

**Access:** `client.Audits()`

### Methods

```go
func (s *AuditService) ListAudits(ctx context.Context, opts *model.ListAuditsOptions) ([]*model.Audit, *model.Pagination, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.Entity | string | Entity type filter |
| opts.Action | string | Action filter |
| opts.From | time.Time | Start date |
| opts.To | time.Time | End date |

---

## ChangeService

**Purpose:** Manages configuration changes and approval workflows.

**Access:** `client.Changes()`

### Methods

```go
func (s *ChangeService) GetChange(ctx context.Context, changeID string) (*model.Change, error)
func (s *ChangeService) ListChanges(ctx context.Context, opts *model.ListChangesOptions) ([]*model.Change, *model.Pagination, error)
func (s *ChangeService) ApproveChange(ctx context.Context, changeID string) error
func (s *ChangeService) RejectChange(ctx context.Context, changeID string, comment string) error
```

### Key Models

- `model.Change` - ID, Entity, Operation, Status, Payload, Trails

---

## PriceService

**Purpose:** Provides price data and currency conversion.

**Access:** `client.Prices()`

### Methods

```go
func (s *PriceService) GetConversions(ctx context.Context, opts *model.GetConversionsOptions) ([]*model.Conversion, error)
func (s *PriceService) ConvertPrices(ctx context.Context, currency string, amount string, targets []string) ([]*model.ConversionResult, error)
```

### Key Models

- `model.Conversion` - Currency pair, rate, timestamp
- `model.ConversionResult` - Target currency, converted amount

---

## WebhookService

**Purpose:** Manages webhooks for receiving real-time event notifications.

**Access:** `client.Webhooks()`

### Methods

```go
func (s *WebhookService) CreateWebhook(ctx context.Context, req *model.CreateWebhookRequest) (*model.Webhook, error)
func (s *WebhookService) GetWebhook(ctx context.Context, webhookID string) (*model.Webhook, error)
func (s *WebhookService) ListWebhooks(ctx context.Context, opts *model.ListWebhooksOptions) ([]*model.Webhook, *model.Pagination, error)
func (s *WebhookService) DeleteWebhook(ctx context.Context, webhookID string) error
func (s *WebhookService) UpdateWebhookStatus(ctx context.Context, webhookID string, status string) (*model.Webhook, error)
```

**Example:**
```go
webhook, err := client.Webhooks().CreateWebhook(ctx, &model.CreateWebhookRequest{
    URL:    "https://example.com/webhook",
    Type:   "TRANSACTION",
    Secret: "my-secret-key",
})
```

### Key Models

- `model.Webhook` - ID, URL, Type, Status, CreatedAt

---

## StakingService

**Purpose:** Retrieves staking information across multiple proof-of-stake blockchains.

**Access:** `client.Staking()`

### Methods

```go
func (s *StakingService) GetStakeAccounts(ctx context.Context, opts *model.GetStakeAccountsOptions) ([]*model.StakeAccount, *model.Pagination, error)
```

### Key Models

- `model.StakeAccount` - AccountAddress, AccountType, Balance, Status

---

## FeeService

**Purpose:** Retrieves transaction fee information.

**Access:** `client.Fees()`

### Methods

```go
func (s *FeeService) GetFees(ctx context.Context, currency string) ([]*model.Fee, error)
```

### Key Models

- `model.Fee` - Currency, FeeType, Amount, Unit

---

## AirGapService

**Purpose:** Provides air-gap signing operations for offline transaction signing.

**Access:** `client.AirGap()`

### Methods

```go
func (s *AirGapService) GetAirGapRequest(ctx context.Context, requestID string) (*model.AirGapRequest, error)
func (s *AirGapService) SubmitAirGapSignature(ctx context.Context, requestID string, signature string) error
```

---

## ReservationService

**Purpose:** Manages balance reservations for addresses.

**Access:** `client.Reservations()`

### Methods

```go
func (s *ReservationService) CreateReservation(ctx context.Context, addressID string, amount string, comment string) (*model.Reservation, error)
func (s *ReservationService) ListReservations(ctx context.Context, addressID string, opts *model.ListReservationsOptions) ([]*model.Reservation, *model.Pagination, error)
func (s *ReservationService) CancelReservation(ctx context.Context, reservationID string) error
```

### Key Models

- `model.Reservation` - ID, AddressID, Amount, Status, ExpiresAt

---

## MultiFactorSignatureService

**Purpose:** Manages multi-factor signature operations for enhanced security.

**Access:** `client.MultiFactorSignature()`

### Methods

```go
func (s *MultiFactorSignatureService) ListMultiFactorSignatures(ctx context.Context, opts *model.ListMultiFactorSignaturesOptions) ([]*model.MultiFactorSignature, *model.Pagination, error)
func (s *MultiFactorSignatureService) ApproveMultiFactorSignature(ctx context.Context, id string, signature string) error
```

---

## ConfigService

**Purpose:** Manages system configuration settings.

**Access:** `client.Config()`

### Methods

```go
func (s *ConfigService) GetTenantConfig(ctx context.Context) (*model.TenantConfig, error)
```

### Key Models

- `model.TenantConfig` - TenantID, Settings, Features

---

## WebhookCallService

**Purpose:** Queries webhook call history.

**Access:** `client.WebhookCalls()`

### Methods

```go
func (s *WebhookCallService) ListWebhookCalls(ctx context.Context, opts *model.ListWebhookCallsOptions) (*model.ListWebhookCallsResult, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| opts.EventID | string | Filter by event ID |
| opts.WebhookID | string | Filter by webhook ID |
| opts.Status | string | Filter by call status |
| opts.CurrentPage | string | Current page cursor |
| opts.PageRequest | string | Page request cursor |
| opts.PageSize | int64 | Page size |
| opts.SortOrder | string | Sort order |

**Returns:** `*model.ListWebhookCallsResult`, `error`

### Key Models

- `model.ListWebhookCallsResult` - WebhookCalls, CurrentPage, HasPrevious, HasNext

---

## AssetService

**Purpose:** Retrieves asset information.

**Access:** `client.Assets()`

### Methods

```go
func (s *AssetService) ListAssets(ctx context.Context, opts *model.ListAssetsOptions) ([]*model.Asset, *model.Pagination, error)
```

### Key Models

- `model.Asset` - ID, Symbol, Name, Blockchain, Network, ContractAddress, Decimals

---

## ActionService

**Purpose:** Manages actions in the system.

**Access:** `client.Actions()`

### Methods

```go
func (s *ActionService) ListActions(ctx context.Context, opts *model.ListActionsOptions) ([]*model.Action, *model.Pagination, error)
```

---

## BlockchainService

**Purpose:** Retrieves blockchain information.

**Access:** `client.Blockchains()`

### Methods

```go
func (s *BlockchainService) ListBlockchains(ctx context.Context) ([]*model.Blockchain, error)
```

### Key Models

- `model.Blockchain` - ID, Name, Networks, Features

---

## ExchangeService

**Purpose:** Manages exchange integrations.

**Access:** `client.Exchanges()`

### Methods

```go
func (s *ExchangeService) ListExchanges(ctx context.Context) ([]*model.Exchange, error)
```

### Key Models

- `model.Exchange` - ID, Name, Type, Status

---

## FiatService

**Purpose:** Manages fiat currency operations.

**Access:** `client.Fiat()`

### Methods

```go
func (s *FiatService) ListFiatCurrencies(ctx context.Context) ([]*model.FiatCurrency, error)
```

---

## FeePayerService

**Purpose:** Manages fee payer configurations.

**Access:** `client.FeePayers()`

### Methods

```go
func (s *FeePayerService) ListFeePayers(ctx context.Context, blockchain string, network string) ([]*model.FeePayer, error)
```

---

## HealthService

**Purpose:** Checks system health status.

**Access:** `client.Health()`

### Methods

```go
func (s *HealthService) Check(ctx context.Context) (*model.HealthStatus, error)
```

### Key Models

- `model.HealthStatus` - Status, Version, Timestamp

---

## JobService

**Purpose:** Manages background jobs.

**Access:** `client.Jobs()`

### Methods

```go
func (s *JobService) ListJobs(ctx context.Context, opts *model.ListJobsOptions) ([]*model.Job, *model.Pagination, error)
```

### Key Models

- `model.Job` - ID, Type, Status, Progress, CreatedAt

---

## StatisticsService

**Purpose:** Retrieves platform statistics.

**Access:** `client.Statistics()`

### Methods

```go
func (s *StatisticsService) GetPortfolioStatistics(ctx context.Context) (*model.PortfolioStatistics, error)
```

### Key Models

- `model.PortfolioStatistics` - TotalValue, AssetBreakdown, ChangePercent

---

## TokenMetadataService

**Purpose:** Retrieves token metadata information.

**Access:** `client.TokenMetadata()`

### Methods

```go
func (s *TokenMetadataService) GetTokenMetadata(ctx context.Context, blockchain string, network string, contractAddress string) (*model.TokenMetadata, error)
```

### Key Models

- `model.TokenMetadata` - Name, Symbol, Decimals, TotalSupply, LogoURL

---

## UserDeviceService

**Purpose:** Manages user device registrations.

**Access:** `client.UserDevices()`

### Methods

```go
func (s *UserDeviceService) ListUserDevices(ctx context.Context, userID string) ([]*model.UserDevice, error)
```

### Key Models

- `model.UserDevice` - ID, DeviceType, Name, LastUsed, Status

---

## TaurusNetwork Services

The TaurusNetwork services are accessed through a namespace pattern:

```go
// Access TaurusNetwork services
client.TaurusNetwork().Participants()
client.TaurusNetwork().Pledges()
client.TaurusNetwork().Lending()
client.TaurusNetwork().Settlements()
client.TaurusNetwork().Sharing()
```

---

## TaurusNetworkParticipantService

**Purpose:** Provides access to Taurus Network participant management.

**Access:** `client.TaurusNetwork().Participants()`

### Methods

```go
func (s *ParticipantService) GetMyParticipant(ctx context.Context) (*model.Participant, error)
func (s *ParticipantService) GetParticipant(ctx context.Context, participantID string, includeTotalPledgesValuation bool) (*model.Participant, error)
func (s *ParticipantService) ListParticipants(ctx context.Context, opts *model.ListParticipantsOptions) ([]*model.Participant, *model.Pagination, error)
```

**Example:**
```go
me, err := client.TaurusNetwork().Participants().GetMyParticipant(ctx)
if err != nil {
    return err
}
fmt.Printf("My participant ID: %s\n", me.ID)
```

### Key Models

- `model.Participant` - ID, Name, Country, PublicKey, TotalPledgesValuation

---

## TaurusNetworkPledgeService

**Purpose:** Provides access to Taurus Network pledge lifecycle operations.

**Access:** `client.TaurusNetwork().Pledges()`

### Methods

```go
func (s *PledgeService) GetPledge(ctx context.Context, pledgeID string) (*model.Pledge, error)
func (s *PledgeService) ListPledges(ctx context.Context, opts *model.ListPledgesOptions) ([]*model.Pledge, *model.CursorPagination, error)
func (s *PledgeService) CreatePledge(ctx context.Context, req *model.CreatePledgeRequest) (*model.Pledge, *model.PledgeAction, error)
func (s *PledgeService) ListPledgeWithdrawals(ctx context.Context, pledgeID string, opts *model.ListPledgeWithdrawalsOptions) ([]*model.PledgeWithdrawal, *model.CursorPagination, error)
func (s *PledgeService) ListPledgeActionsForApproval(ctx context.Context, opts *model.ListPledgeActionsOptions) ([]*model.PledgeAction, *model.CursorPagination, error)
func (s *PledgeService) ApprovePledgeActions(ctx context.Context, actions []*model.PledgeAction, privateKey *ecdsa.PrivateKey) (int, error)
```

### Key Models

- `model.Pledge` - ID, OwnerParticipantID, TargetParticipantID, Amount, Status, Currency
- `model.PledgeAction` - ID, PledgeID, ActionType, Status, Metadata
- `model.PledgeWithdrawal` - Withdrawal details with status

---

## TaurusNetworkLendingService

**Purpose:** Provides access to lending offers and agreements in the Taurus Network.

**Access:** `client.TaurusNetwork().Lending()`

### Methods

```go
func (s *LendingService) GetLendingOffer(ctx context.Context, offerID string) (*model.LendingOffer, error)
func (s *LendingService) ListLendingOffers(ctx context.Context, opts *model.ListLendingOffersOptions) ([]*model.LendingOffer, *model.CursorPagination, error)
func (s *LendingService) CreateLendingOffer(ctx context.Context, req *model.CreateLendingOfferRequest) (*model.LendingOffer, error)
func (s *LendingService) CancelLendingOffer(ctx context.Context, offerID string) (*model.LendingOffer, error)
func (s *LendingService) GetLendingAgreement(ctx context.Context, agreementID string) (*model.LendingAgreement, error)
func (s *LendingService) ListLendingAgreements(ctx context.Context, opts *model.ListLendingAgreementsOptions) ([]*model.LendingAgreement, *model.CursorPagination, error)
func (s *LendingService) AcceptLendingOffer(ctx context.Context, offerID string, req *model.AcceptLendingOfferRequest) (*model.LendingAgreement, error)
```

### Key Models

- `model.LendingOffer` - ID, ParticipantID, CurrencyID, Amount, Rate, Duration, Status
- `model.LendingAgreement` - ID, OfferID, BorrowerParticipantID, LenderParticipantID, Amount, Status

---

## TaurusNetworkSettlementService

**Purpose:** Provides access to settlements in the Taurus Network.

**Access:** `client.TaurusNetwork().Settlements()`

### Methods

```go
func (s *SettlementService) GetSettlement(ctx context.Context, settlementID string) (*model.Settlement, error)
func (s *SettlementService) ListSettlements(ctx context.Context, opts *model.ListSettlementsOptions) ([]*model.Settlement, *model.CursorPagination, error)
func (s *SettlementService) CreateSettlement(ctx context.Context, req *model.CreateSettlementRequest) (*model.Settlement, error)
func (s *SettlementService) ApproveSettlement(ctx context.Context, settlementID string) (*model.Settlement, error)
func (s *SettlementService) RejectSettlement(ctx context.Context, settlementID string, comment string) (*model.Settlement, error)
```

### Key Models

- `model.Settlement` - ID, CounterParticipantID, Amount, Currency, Status, CreatedAt

---

## TaurusNetworkSharingService

**Purpose:** Provides access to address and asset sharing in the Taurus Network.

**Access:** `client.TaurusNetwork().Sharing()`

### Methods

```go
func (s *SharingService) ListSharedAddresses(ctx context.Context, opts *model.ListSharedAddressesOptions) ([]*model.SharedAddress, *model.CursorPagination, error)
func (s *SharingService) CreateSharedAddress(ctx context.Context, req *model.CreateSharedAddressRequest) (*model.SharedAddress, error)
func (s *SharingService) RevokeSharedAddress(ctx context.Context, sharedAddressID string) error
func (s *SharingService) ListSharedAssets(ctx context.Context, opts *model.ListSharedAssetsOptions) ([]*model.SharedAsset, *model.CursorPagination, error)
func (s *SharingService) CreateSharedAsset(ctx context.Context, req *model.CreateSharedAssetRequest) (*model.SharedAsset, error)
func (s *SharingService) RevokeSharedAsset(ctx context.Context, sharedAssetID string) error
```

### Key Models

- `model.SharedAddress` - ID, Blockchain, Network, Address, OwnerParticipantID, TargetParticipantID, Permissions
- `model.SharedAsset` - ID, AssetID, OwnerParticipantID, TargetParticipantID, Permissions

---

## Error Handling

All services follow a consistent error pattern:

```go
wallet, err := client.Wallets().GetWallet(ctx, walletID)
if err != nil {
    if apiErr, ok := protect.IsAPIError(err); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.Code, apiErr.Message)
        if apiErr.Code == 404 {
            fmt.Println("Wallet not found")
        }
    }
    return err
}
```

### Error Types

| Error Type | When Returned |
|------------|---------------|
| `APIError` | General API errors (network, auth, validation) |
| `IntegrityError` | Hash verification or signature verification failed |
| `WhitelistError` | Whitelist-specific verification errors |
| `RequestMetadataError` | Request metadata validation errors |

### Sentinel Errors

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

errors.Is(err, protect.ErrValidation)      // 400
errors.Is(err, protect.ErrAuthentication)  // 401
errors.Is(err, protect.ErrAuthorization)   // 403
errors.Is(err, protect.ErrNotFound)        // 404
errors.Is(err, protect.ErrRateLimit)       // 429
errors.Is(err, protect.ErrServer)          // 500
```

---

## Pagination Patterns

### Standard Pattern

```go
var allWallets []*model.Wallet
opts := &model.ListWalletsOptions{Limit: 100, Offset: 0}

for {
    wallets, pagination, err := client.Wallets().ListWallets(ctx, opts)
    if err != nil {
        return err
    }

    allWallets = append(allWallets, wallets...)

    if !pagination.HasMore {
        break
    }
    opts.Offset = pagination.Offset + pagination.Limit
}
```

### Pagination Model

```go
type Pagination struct {
    Limit      int64  // Items per page
    Offset     int64  // Current offset
    TotalItems int64  // Total available items
    HasMore    bool   // More items available
}
```

---

## Context Usage

All service methods accept `context.Context` as the first parameter:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

wallet, err := client.Wallets().GetWallet(ctx, walletID)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // Cancel on signal
    <-sigChan
    cancel()
}()

wallets, _, err := client.Wallets().ListWallets(ctx, opts)
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Verification details