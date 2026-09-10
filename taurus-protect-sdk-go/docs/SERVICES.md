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
| [WhitelistedContractService](#whitelistedcontractservice) | `client.WhitelistedContracts()` | Smart contract whitelisting — WRITES only; reads live on `WhitelistedAssets()` |
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

> **A request whose metadata hash has not been verified is refused.** The signature attests
> to those hashes, so each must be one verification cleared against its payload — the
> `HashVerified` flag was set on every read path and read by nobody. The check runs after the
> hash-present check, so absent metadata still reports as absent.

```go
func (s *RequestService) ApproveRequests(ctx context.Context, requests []*model.Request, privateKey *ecdsa.PrivateKey) (int, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| ctx | context.Context | Request context |
| requests | []*model.Request | Requests to approve (must carry a VERIFIED Metadata.Hash) |
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

#### UpdateRulesProposal

Submits a typed rules container as a governance proposal (SuperAdmin only). The container is encoded to the wire format internally; the server-controlled `enforcedRulesHash` and `timestamp` fields are stripped. The endpoint returns no body — callers needing the persisted proposal should call `GetRulesProposal` (note: rules reads are cached server-side, so an immediate read-back may be stale).

```go
func (s *GovernanceRuleService) UpdateRulesProposal(ctx context.Context, container *model.DecodedRulesContainer) error
```

#### ApproveRulesProposal

Signs the pending proposal's rules container with a SuperAdmin private key (SHA-256 + P-256 ECDSA, base64 raw r||s) and submits the approval. `expectedContainerHash` PINS the content: pass `ProposalContainerHash` of the proposal you reviewed, and the call aborts without signing if the re-fetched container differs. Review with `GetRulesProposal` + `DecodeProposalForReview` — **not** `GetDecodedRulesContainer`, which verifies unconditionally and so always fails on a pending proposal (it legitimately carries 0..N signatures; signing IS the approval step).

```go
func (s *GovernanceRuleService) ApproveRulesProposal(ctx context.Context, privateKey *ecdsa.PrivateKey, comment string) error
```

#### RejectRulesProposal

Rejects the pending rules proposal with a comment (SuperAdmin only).

```go
func (s *GovernanceRuleService) RejectRulesProposal(ctx context.Context, comment string) error
```

### Key Models

- `model.GovernanceRuleset` - RulesContainer, Signatures, Locked, Trails
- `model.DecodedRulesContainer` - lossless typed rules container (users, groups, transaction/whitelisting rules); round-trips through `mapper.RulesContainerToBase64` / `RulesContainerFromBase64`
- `model.RuleCell` - typed transaction-rule cell union covering every cell type (`FiatAmountAny`, `FiatAmountRange`, `SourceInternalWallet`, `DestinationInternalWallet`, `StringEqualValue`, ...); `RawCell` preserves cells from newer schemas verbatim. The `*Any` values are the protobuf zero values, so their serialized form is the empty cell; `nil` is accepted at encode as an empty cell.
- `model.RuleUserSignature` - UserID, Signature
- `model.SuperAdminPublicKey` - ID, PublicKey, Name

Cross-SDK cell wire-format parity is pinned by the shared golden vectors at `scripts/resources/governance-cell-vectors.json` (monorepo root), consumed by every SDK's test suite.

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

### Key Models

- `model.WhitelistedAddress` - ID, Blockchain, Network, Address, Name, Status, RulesContainer, SignedAddress, Approvers
- `model.SignedWhitelistedAddress` - Payload, Signatures
- `model.WhitelistSignature` - UserSignature, Hashes

---

## WhitelistedAssetService

**Purpose:** Manages whitelisted tokens and contracts. **This is the only verified reader of
`/whitelists/contracts`** — a whitelisted asset and a whitelisted contract are one server
entity, and `WhitelistedContractService` is write-only.

**Access:** `client.WhitelistedAssets()`

### Methods

```go
func (s *WhitelistedAssetService) GetWhitelistedAsset(ctx context.Context, id string) (*model.WhitelistedAsset, error)
func (s *WhitelistedAssetService) GetWhitelistedAssetEnvelope(ctx context.Context, id string) (*model.WhitelistedAssetEnvelope, error)
func (s *WhitelistedAssetService) ListWhitelistedAssets(ctx context.Context, opts *model.ListWhitelistedAssetsOptions) ([]*model.WhitelistedAsset, *model.Pagination, error)
func (s *WhitelistedAssetService) ListWhitelistedAssetsForApproval(ctx context.Context, opts *model.ListWhitelistedAssetsForApprovalOptions) ([]*model.WhitelistedAsset, *model.Pagination, error)
func (s *WhitelistedAssetService) ApproveWhitelistedAssets(ctx context.Context, ids []string, privateKey *ecdsa.PrivateKey, comment string) error
```

`Get` and `List` populate `ContractAddress`, `Name`, `Symbol`, `Decimals` and `TokenID` by
running step 6 on the verified payload — without them the verified reader could not answer
"is this contract whitelisted, and at what address?", which is what sent callers to the
unverified one.

`ApproveWhitelistedAssets` is **all-or-nothing**: each asset is re-read and verified, the
hashes those rows carry are what gets signed, and any row that is missing or fails
verification aborts the whole call. The API takes one signature covering the whole batch, so a
partial approval would mean the caller believes they approved more than they did.
`WhitelistedContractService.ApproveWhitelistedContract` — which takes an opaque signature over
hashes nothing verified — is deprecated in its favour.

### Key Models

- `model.WhitelistedAsset` - ID, Blockchain, Network, Status, Metadata, SignedContractAddress, Approvers, ContractAddress, Name, Symbol, Decimals, TokenID

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
func (s *AuditService) ListAuditTrails(ctx context.Context, opts *model.ListAuditsOptions) ([]*model.Audit, *model.Pagination, error)
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

`Rate` and `Decimals` feed amount conversion, so an unverified price is a wrong number a
caller acts on. Every price read is verified against the `PRICEUPDATER` keys in the
SuperAdmin-verified rules container, which is why the service takes the cache as a mandatory
parameter. Whether prices must be signed is the **container's** call: no `PRICEUPDATER`
configured means this tenant does not sign prices and the price passes through; a
`PRICEUPDATER` configured plus a price carrying no signatures is an `IntegrityError`.

### Methods

```go
func (s *PriceService) Convert(ctx context.Context, currency string, amount string, targets []string) ([]*model.ConversionResult, error)
```

### Key Models

- `model.Conversion` - Currency pair, rate, timestamp
- `model.ConversionResult` - Target currency, converted amount
- `model.Price` - Blockchain, CurrencyFrom, CurrencyTo, Decimals, Rate, Signatures

---

## WebhookService

**Purpose:** Manages webhooks for receiving real-time event notifications.

**Access:** `client.Webhooks()`

### Methods

```go
func (s *WebhookService) CreateWebhook(ctx context.Context, req *model.CreateWebhookRequest) (*model.Webhook, error)
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
func (s *StakingService) ListStakeAccounts(ctx context.Context, opts *model.GetStakeAccountsOptions) ([]*model.StakeAccount, *model.Pagination, error)
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
func (s *AirGapService) GetOutgoingAirGap(ctx context.Context, requestID string) (*model.AirGapRequest, error)
func (s *AirGapService) SubmitIncomingAirGap(ctx context.Context, requestID string, signature string) error
```

---

## ReservationService

**Purpose:** Manages balance reservations for addresses.

**Access:** `client.Reservations()`

### Methods

```go
func (s *ReservationService) ListReservations(ctx context.Context, addressID string, opts *model.ListReservationsOptions) ([]*model.Reservation, *model.Pagination, error)
```

### Key Models

- `model.Reservation` - ID, AddressID, Amount, Status, ExpiresAt

---

## MultiFactorSignatureService

**Purpose:** Manages multi-factor signature operations for enhanced security.

**Access:** `client.MultiFactorSignature()`

### Methods

```go
func (s *MultiFactorSignatureService) GetMultiFactorSignatureInfo(ctx context.Context, opts *model.ListMultiFactorSignaturesOptions) ([]*model.MultiFactorSignature, *model.Pagination, error)
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

`GetAssetAddresses` verifies every address's HSM signature — the same check `AddressService`
runs — and **fails fast** on the first that does not verify. It returns the same entity, so
returning it unverified made `AddressService`'s mandatory verification avoidable. The service
therefore takes the rules cache as a mandatory parameter and panics on a nil one.

### Methods

```go
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
func (s *FiatService) ListFiatProviders(ctx context.Context) ([]*model.FiatCurrency, error)
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
func (s *HealthService) GetAllHealthChecks(ctx context.Context) (*model.HealthStatus, error)
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
func (s *TokenMetadataService) GetERCTokenMetadata(ctx context.Context, blockchain string, network string, contractAddress string) (*model.TokenMetadata, error)
```

### Key Models

- `model.TokenMetadata` - Name, Symbol, Decimals, TotalSupply, LogoURL

---

## UserDeviceService

**Purpose:** Manages user device registrations.

**Access:** `client.UserDevices()`

### Methods

```go
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
func (s *TaurusNetworkParticipantService) GetMyParticipant(ctx context.Context) (*model.Participant, error)
func (s *TaurusNetworkParticipantService) GetParticipant(ctx context.Context, participantID string, includeTotalPledgesValuation bool) (*model.Participant, error)
func (s *TaurusNetworkParticipantService) ListParticipants(ctx context.Context, opts *model.ListParticipantsOptions) ([]*model.Participant, *model.Pagination, error)
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
func (s *TaurusNetworkPledgeService) GetPledge(ctx context.Context, pledgeID string) (*model.Pledge, error)
func (s *TaurusNetworkPledgeService) ListPledges(ctx context.Context, opts *model.ListPledgesOptions) ([]*model.Pledge, *model.CursorPagination, error)
func (s *TaurusNetworkPledgeService) CreatePledge(ctx context.Context, req *model.CreatePledgeRequest) (*model.Pledge, *model.PledgeAction, error)
func (s *TaurusNetworkPledgeService) ListPledgeWithdrawals(ctx context.Context, pledgeID string, opts *model.ListPledgeWithdrawalsOptions) ([]*model.PledgeWithdrawal, *model.CursorPagination, error)
func (s *TaurusNetworkPledgeService) ListPledgeActionsForApproval(ctx context.Context, opts *model.ListPledgeActionsOptions) ([]*model.PledgeAction, *model.CursorPagination, error)
func (s *TaurusNetworkPledgeService) ApprovePledgeActions(ctx context.Context, actions []*model.PledgeAction, privateKey *ecdsa.PrivateKey) (int, error)
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
func (s *TaurusNetworkLendingService) GetLendingOffer(ctx context.Context, offerID string) (*model.LendingOffer, error)
func (s *TaurusNetworkLendingService) ListLendingOffers(ctx context.Context, opts *model.ListLendingOffersOptions) ([]*model.LendingOffer, *model.CursorPagination, error)
func (s *TaurusNetworkLendingService) CreateLendingOffer(ctx context.Context, req *model.CreateLendingOfferRequest) (*model.LendingOffer, error)
func (s *TaurusNetworkLendingService) DeleteLendingOffer(ctx context.Context, offerID string) (*model.LendingOffer, error)
func (s *TaurusNetworkLendingService) GetLendingAgreement(ctx context.Context, agreementID string) (*model.LendingAgreement, error)
func (s *TaurusNetworkLendingService) ListLendingAgreements(ctx context.Context, opts *model.ListLendingAgreementsOptions) ([]*model.LendingAgreement, *model.CursorPagination, error)
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
func (s *TaurusNetworkSettlementService) GetSettlement(ctx context.Context, settlementID string) (*model.Settlement, error)
func (s *TaurusNetworkSettlementService) ListSettlements(ctx context.Context, opts *model.ListSettlementsOptions) ([]*model.Settlement, *model.CursorPagination, error)
func (s *TaurusNetworkSettlementService) CreateSettlement(ctx context.Context, req *model.CreateSettlementRequest) (*model.Settlement, error)
```

### Key Models

- `model.Settlement` - ID, CounterParticipantID, Amount, Currency, Status, CreatedAt

---

## TaurusNetworkSharingService

**Purpose:** Provides access to address and asset sharing in the Taurus Network.

**Access:** `client.TaurusNetwork().Sharing()`

### Methods

```go
func (s *TaurusNetworkSharingService) ListSharedAddresses(ctx context.Context, opts *model.ListSharedAddressesOptions) ([]*model.SharedAddress, *model.CursorPagination, error)
func (s *TaurusNetworkSharingService) ShareAddress(ctx context.Context, req *model.CreateSharedAddressRequest) (*model.SharedAddress, error)
func (s *TaurusNetworkSharingService) UnshareAddress(ctx context.Context, sharedAddressID string) error
func (s *TaurusNetworkSharingService) ListSharedAssets(ctx context.Context, opts *model.ListSharedAssetsOptions) ([]*model.SharedAsset, *model.CursorPagination, error)
func (s *TaurusNetworkSharingService) ShareWhitelistedAsset(ctx context.Context, req *model.CreateSharedAssetRequest) (*model.SharedAsset, error)
func (s *TaurusNetworkSharingService) UnshareWhitelistedAsset(ctx context.Context, sharedAssetID string) error
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

<!-- BEGIN GENERATED METHOD INDEX -->

## Complete Method Index

Generated from the go source by `scripts/api-surface/docs.py`; regenerate with
`./build.sh docs`. Every method below exists in the SDK, and `./build.sh docs --check`
fails if this list drifts or if the prose above documents a method that does not.

43 services, 194 public methods.

### ActionService

- `GetAction(ctx context.Context, id string) (*model.Action, error)` — GetAction retrieves a single action by ID.
- `ListActions(ctx context.Context, opts *model.ListActionsOptions) (*model.ListActionsResult, error)` — ListActions retrieves a list of actions with optional filtering and pagination.

### AddressService

- `CreateAddress(ctx context.Context, req *model.CreateAddressRequest) (*model.Address, error)` — CreateAddress creates a new address in a wallet.
- `CreateAddressAttribute(ctx context.Context, addressID string, key string, value string) ([]model.AddressAttribute, error)` — CreateAddressAttribute creates an attribute on an address.
- `DeleteAddressAttribute(ctx context.Context, addressID string, attributeID string) error` — DeleteAddressAttribute deletes an attribute from an address.
- `GetAddress(ctx context.Context, addressID string) (*model.Address, error)` — GetAddress retrieves an address by ID with mandatory signature verification.
- `GetAddressProofOfReserve(ctx context.Context, addressID string, challenge string) (*model.ProofOfReserve, error)` — GetAddressProofOfReserve retrieves the proof of reserve for an address.
- `ListAddresses(ctx context.Context, opts *model.ListAddressesOptions) ([]*model.Address, *model.Pagination, error)` — ListAddresses retrieves a list of addresses with mandatory signature verification.

### AirGapService

- `GetOutgoingAirGap(ctx context.Context, req *model.GetOutgoingAirGapRequest) (*model.GetOutgoingAirGapResult, error)` — GetOutgoingAirGap exports HSM ready requests for cold HSM.
- `SubmitIncomingAirGap(ctx context.Context, req *model.SubmitIncomingAirGapRequest) error` — SubmitIncomingAirGap imports signed requests from the cold HSM.

### AssetService

- `GetAssetAddresses(ctx context.Context, req *model.GetAssetAddressesRequest) (*model.GetAssetAddressesResult, error)` — GetAssetAddresses retrieves address-level balances for a specific asset.
- `GetAssetWallets(ctx context.Context, req *model.GetAssetWalletsRequest) (*model.GetAssetWalletsResult, error)` — GetAssetWallets retrieves wallet-level balances for a specific asset.

### AuditService

- `ExportAuditTrails(ctx context.Context, opts *model.ExportAuditTrailsOptions) (*model.ExportAuditTrailsResult, error)` — ExportAuditTrails exports audit trails in the specified format (CSV or JSON).
- `ListAuditTrails(ctx context.Context, opts *model.ListAuditTrailsOptions) (*model.ListAuditTrailsResult, error)` — ListAuditTrails retrieves a list of audit trails with optional filtering and pagination.

### BalanceService

- `GetBalances(ctx context.Context, opts *model.GetBalancesOptions) (*model.GetBalancesResult, error)` — GetBalances retrieves the total balances for the tenant, for each asset.

### BlockchainService

- `ListBlockchains(ctx context.Context, opts *model.ListBlockchainsOptions) ([]*model.Blockchain, error)` — ListBlockchains retrieves a list of all enabled blockchains.

### BusinessRuleService

- `ListBusinessRules(ctx context.Context, opts *model.ListBusinessRulesOptions) (*model.ListBusinessRulesResult, error)` — ListBusinessRules retrieves a list of business rules with optional filtering and pagination.
- `UpdateTransactionsEnabled(ctx context.Context, req *model.UpdateTransactionsEnabledRequest) error` — UpdateTransactionsEnabled toggles the transactions enabled business rule.

### ChangeService

- `ApproveChange(ctx context.Context, changeID string) error` — ApproveChange approves a single change by ID.
- `ApproveChanges(ctx context.Context, changeIDs []string) error` — ApproveChanges approves multiple changes by their IDs.
- `CreateChange(ctx context.Context, req *model.CreateChangeRequest) (*model.CreateChangeResult, error)` — CreateChange creates a new change request.
- `GetChange(ctx context.Context, changeID string) (*model.Change, error)` — GetChange retrieves a change by ID.
- `ListChanges(ctx context.Context, opts *model.ListChangesOptions) (*model.ListChangesResult, error)` — ListChanges retrieves a list of changes.
- `ListChangesForApproval(ctx context.Context, opts *model.ListChangesForApprovalOptions) (*model.ListChangesResult, error)` — ListChangesForApproval retrieves a list of changes pending approval for the current user.
- `RejectChange(ctx context.Context, changeID string) error` — RejectChange rejects a single change by ID.
- `RejectChanges(ctx context.Context, changeIDs []string) error` — RejectChanges rejects multiple changes by their IDs.

### ConfigService

- `GetTenantConfig(ctx context.Context) (*model.TenantConfig, error)` — GetTenantConfig retrieves the configuration of the tenant for the connected user.

### CurrencyService

- `GetCurrencies(ctx context.Context, opts *model.ListCurrenciesOptions) ([]*model.Currency, error)` — GetCurrencies retrieves a list of all currencies.
- `GetCurrency(ctx context.Context, opts *model.GetCurrencyOptions) (*model.Currency, error)` — GetCurrency retrieves a single currency by its filter criteria.

### ExchangeService

- `ExportExchanges(ctx context.Context, opts *model.ExportExchangesOptions) (*model.ExportExchangesResult, error)` — ExportExchanges exports exchange accounts in the specified format (CSV or JSON).
- `GetExchange(ctx context.Context, id string) (*model.Exchange, error)` — GetExchange retrieves a single exchange account by ID.
- `ListExchangeCounterparties(ctx context.Context) (*model.ListExchangeCounterpartiesResult, error)` — ListExchangeCounterparties retrieves a list of exchange counterparties with their exposure and limits.
- `ListExchanges(ctx context.Context, opts *model.ListExchangesOptions) (*model.ListExchangesResult, error)` — ListExchanges retrieves a list of exchange accounts with optional filtering and pagination.

### FeePayerService

- `GetChecksum(ctx context.Context, req *model.ChecksumRequest) (*model.ChecksumResult, error)` — GetChecksum computes a checksum for the provided data.
- `GetFeePayer(ctx context.Context, id string) (*model.FeePayer, error)` — GetFeePayer retrieves a single fee payer by ID.
- `ListFeePayers(ctx context.Context, opts *model.ListFeePayersOptions) (*model.ListFeePayersResult, error)` — ListFeePayers retrieves a list of fee payers with optional filtering.

### FeeService

- `GetFees(ctx context.Context) (*model.GetFeesResult, error)` — GetFees retrieves a list of fee estimates.
- `GetFeesV2(ctx context.Context) (*model.GetFeesV2Result, error)` — GetFeesV2 retrieves a list of native currency fee estimates.

### FiatService

- `GetFiatProviderAccount(ctx context.Context, id string) (*model.FiatProviderAccount, error)` — GetFiatProviderAccount retrieves a fiat provider account by ID.
- `ListFiatProviderAccounts(ctx context.Context, opts *model.ListFiatProviderAccountsOptions) (*model.ListFiatProviderAccountsResult, error)` — ListFiatProviderAccounts retrieves a list of fiat provider accounts with optional filtering and pagination.
- `ListFiatProviders(ctx context.Context) (*model.ListFiatProvidersResult, error)` — ListFiatProviders retrieves a list of all enabled fiat providers and their valuations.

### GovernanceRuleService

- `ApproveRulesProposal(ctx context.Context, privateKey *ecdsa.PrivateKey, comment string, expectedContainerHash string) error` — ApproveRulesProposal signs the pending rules proposal with a SuperAdmin
- `DecodeProposalForReview(rules *model.GovernanceRuleset) (*model.DecodedRulesContainer, error)` — DecodeProposalForReview decodes a PENDING rules proposal so a SuperAdmin can inspect it
- `GetDecodedRulesContainer(rules *model.GovernanceRuleset) (*model.DecodedRulesContainer, error)` — GetDecodedRulesContainer decodes and verifies a GovernanceRuleset's rules container.
- `GetPublicKeys(ctx context.Context) ([]*model.SuperAdminPublicKey, error)` — GetPublicKeys retrieves the list of SuperAdmin public keys.
- `GetRules(ctx context.Context) (*model.GovernanceRuleset, error)` — GetRules retrieves the currently enforced governance rules.
- `GetRulesByID(ctx context.Context, id string) (*model.GovernanceRuleset, error)` — GetRulesByID retrieves a governance ruleset by its ID.
- `GetRulesHistory(ctx context.Context, opts *model.ListRulesHistoryOptions) (*model.GovernanceRulesHistoryResult, error)` — GetRulesHistory retrieves the history of governance rules with pagination.
- `GetRulesProposal(ctx context.Context) (*model.GovernanceRuleset, error)` — GetRulesProposal retrieves the proposed governance rules.
- `MinValidSignatures() int` — MinValidSignatures returns the minimum number of valid signatures required.
- `ProposalContainerHash(rules *model.GovernanceRuleset) (string, error)` — ProposalContainerHash returns the canonical SHA-256 hex digest of a ruleset's decoded
- `RejectRulesProposal(ctx context.Context, comment string) error` — RejectRulesProposal rejects the pending rules proposal with a comment.
- `SuperAdminKeys() []*ecdsa.PublicKey` — SuperAdminKeys returns the configured SuperAdmin public keys.
- `UpdateRulesProposal(ctx context.Context, container *model.DecodedRulesContainer) error` — UpdateRulesProposal submits a rules container as a governance proposal.
- `VerifyGovernanceRules(rules *model.GovernanceRuleset) (*model.GovernanceRuleset, error)` — VerifyGovernanceRules verifies the SuperAdmin signatures on the governance rules and

### GroupService

- `ListGroups(ctx context.Context, opts *model.ListGroupsOptions) (*model.ListGroupsResult, error)` — ListGroups retrieves a list of groups with optional filtering and pagination.

### HealthService

- `GetAllHealthChecks(ctx context.Context, opts *model.GetAllHealthChecksOptions) (*model.GetAllHealthChecksResult, error)` — GetAllHealthChecks retrieves all health checks with optional filtering.

### JobService

- `GetJob(ctx context.Context, name string) (*model.Job, error)` — GetJob retrieves a single job by name.
- `GetJobStatus(ctx context.Context, name string, id string) (*model.JobStatus, error)` — GetJobStatus retrieves the status of a specific job execution.
- `ListJobs(ctx context.Context) ([]*model.Job, error)` — ListJobs retrieves all registered jobs.

### MultiFactorSignatureService

- `ApproveMultiFactorSignature(ctx context.Context, id string, signature string, comment string) (*model.MultiFactorSignatureApprovalResult, error)` — ApproveMultiFactorSignature approves a multi-factor signature request.
- `CreateMultiFactorSignatures(ctx context.Context, entityIDs []string, entityType model.MultiFactorSignatureEntityType) (*model.MultiFactorSignatureResult, error)` — CreateMultiFactorSignatures creates a batch of multi-factor signature requests.
- `GetMultiFactorSignatureInfo(ctx context.Context, id string) (*model.MultiFactorSignatureInfo, error)` — GetMultiFactorSignatureInfo retrieves information about a multi-factor signature request.
- `RejectMultiFactorSignature(ctx context.Context, id string, comment string) error` — RejectMultiFactorSignature rejects a multi-factor signature request.

### PriceService

- `Convert(ctx context.Context, opts *model.ConvertOptions) (*model.ConversionResult, error)` — Convert converts an amount from one currency to other currencies.
- `ExportPriceHistory(ctx context.Context, opts *model.ExportPriceHistoryOptions) (*model.ExportPriceHistoryResult, error)` — ExportPriceHistory exports the price history in a specified format.
- `GetPriceHistory(ctx context.Context, opts *model.GetPriceHistoryOptions) (*model.GetPriceHistoryResult, error)` — GetPriceHistory retrieves the price history for a currency pair.
- `GetPrices(ctx context.Context) (*model.GetPricesResult, error)` — GetPrices retrieves all available currency prices.

### RequestService

- `ApproveRequest(ctx context.Context, request *model.Request, privateKey *ecdsa.PrivateKey, comment ...string) (int, error)` — ApproveRequest approves a single request using a private key for signing.
- `ApproveRequests(ctx context.Context, requests []*model.Request, privateKey *ecdsa.PrivateKey, comment ...string) (int, error)` — ApproveRequests approves multiple requests using a private key for signing.
- `CreateCancelRequest(ctx context.Context, addressID string, nonce string) (*model.Request, error)` — CreateCancelRequest creates a cancel request for a pending transaction.
- `CreateExternalTransferFromWalletRequest(ctx context.Context, fromWalletID, toWhitelistedAddressID, amount string) (*model.Request, error)` — CreateExternalTransferFromWalletRequest creates an external transfer from an omnibus wallet.
- `CreateExternalTransferRequest(ctx context.Context, fromAddressID, toWhitelistedAddressID, amount string) (*model.Request, error)` — CreateExternalTransferRequest creates an external transfer to a whitelisted address.
- `CreateIncomingRequest(ctx context.Context, req *model.CreateIncomingRequest) (*model.Request, error)` — CreateIncomingRequest creates an incoming request to log an incoming transaction from an exchange.
- `CreateInternalTransferFromWalletRequest(ctx context.Context, fromWalletID, toAddressID, amount string) (*model.Request, error)` — CreateInternalTransferFromWalletRequest creates an internal transfer from an omnibus wallet.
- `CreateInternalTransferRequest(ctx context.Context, fromAddressID, toAddressID, amount string) (*model.Request, error)` — CreateInternalTransferRequest creates an internal transfer request from one address to another.
- `CreateOutgoingRequest(ctx context.Context, req *model.CreateOutgoingRequest) (*model.Request, error)` — CreateOutgoingRequest creates a new outgoing (withdrawal) request.
- `GetRequest(ctx context.Context, requestID string) (*model.Request, error)` — GetRequest retrieves a request by ID with hash verification.
- `ListRequests(ctx context.Context, opts *model.ListRequestsOptions) (*model.RequestResult, error)` — ListRequests retrieves a list of requests using cursor-based pagination.
- `ListRequestsForApproval(ctx context.Context, opts *model.ListRequestsOptions) (*model.RequestResult, error)` — ListRequestsForApproval retrieves requests pending approval for the current user.
- `RejectRequest(ctx context.Context, requestID string, comment string) error` — RejectRequest rejects a single request with a comment.
- `RejectRequests(ctx context.Context, requestIDs []string, comment string) error` — RejectRequests rejects multiple requests with a comment.

### ReservationService

- `GetReservation(ctx context.Context, id string) (*model.Reservation, error)` — GetReservation retrieves a single reservation by ID.
- `GetReservationUTXO(ctx context.Context, id string) (*model.ReservationUTXO, error)` — GetReservationUTXO retrieves the UTXO associated with a reservation.
- `ListReservations(ctx context.Context, opts *model.ListReservationsOptions) (*model.ListReservationsResult, error)` — ListReservations retrieves a list of reservations with optional filtering and pagination.

### ScoreService

- `RefreshAddressScore(ctx context.Context, addressID string, provider string) (*model.RefreshScoreResult, error)` — RefreshAddressScore refreshes the risk score for an address using the specified provider.
- `RefreshWLAScore(ctx context.Context, addressID string, provider string) (*model.RefreshScoreResult, error)` — RefreshWLAScore refreshes the risk score for a whitelisted address using the specified provider.

### StakingService

- `GetADAStakePoolInfo(ctx context.Context, network, stakePoolID string) (*model.ADAStakePoolInfo, error)` — GetADAStakePoolInfo retrieves information about an ADA stake pool.
- `GetETHValidatorsInfo(ctx context.Context, network string, ids []string) ([]*model.ETHValidatorInfo, error)` — GetETHValidatorsInfo retrieves information about ETH validators.
- `GetFTMValidatorInfo(ctx context.Context, network, validatorAddress string) (*model.FTMValidatorInfo, error)` — GetFTMValidatorInfo retrieves information about an FTM validator.
- `GetICPNeuronInfo(ctx context.Context, network, neuronID string) (*model.ICPNeuronInfo, error)` — GetICPNeuronInfo retrieves information about an ICP neuron.
- `GetNEARValidatorInfo(ctx context.Context, network, validatorAddress string) (*model.NEARValidatorInfo, error)` — GetNEARValidatorInfo retrieves information about a NEAR validator.
- `GetXTZAddressStakingRewards(ctx context.Context, network, addressID string, opts *model.GetXTZStakingRewardsOptions) (*model.XTZStakingReward, error)` — GetXTZAddressStakingRewards retrieves staking rewards for an XTZ address.
- `ListStakeAccounts(ctx context.Context, opts *model.ListStakeAccountsOptions) (*model.ListStakeAccountsResult, error)` — ListStakeAccounts retrieves a list of stake accounts with optional filtering and pagination.

### StatisticsService

- `GetPortfolioStatistics(ctx context.Context) (*model.PortfolioStatistics, error)` — GetPortfolioStatistics retrieves the global portfolio statistics.
- `GetPortfolioStatisticsHistory(ctx context.Context, opts *model.GetPortfolioStatisticsHistoryOptions) (*model.GetPortfolioStatisticsHistoryResult, error)` — GetPortfolioStatisticsHistory retrieves the portfolio statistics history with optional filtering and pagination.
- `ListTagStatistics(ctx context.Context, opts *model.ListTagStatisticsOptions) (*model.ListTagStatisticsResult, error)` — ListTagStatistics retrieves a list of tag statistics with optional filtering and pagination.

### TagService

- `CreateTag(ctx context.Context, req *model.CreateTagRequest) (string, error)` — CreateTag creates a new tag and returns its ID.
- `DeleteTag(ctx context.Context, id string) error` — DeleteTag deletes a tag by ID.
- `ListTags(ctx context.Context, opts *model.ListTagsOptions) ([]*model.Tag, error)` — ListTags retrieves a list of tags with optional filtering.

### TaurusNetworkLendingService

- `CancelLendingAgreement(ctx context.Context, lendingAgreementID string) error` — CancelLendingAgreement cancels a lending agreement that is not yet approved by the lender.
- `CreateLendingAgreement(ctx context.Context, req *taurusnetwork.CreateLendingAgreementRequest) (*taurusnetwork.LendingAgreement, error)` — CreateLendingAgreement creates a new lending agreement.
- `CreateLendingAgreementAttachment(ctx context.Context, lendingAgreementID string, req *taurusnetwork.CreateLendingAgreementAttachmentRequest) error` — CreateLendingAgreementAttachment creates an attachment for a lending agreement.
- `CreateLendingOffer(ctx context.Context, req *taurusnetwork.CreateLendingOfferRequest) (*taurusnetwork.LendingOffer, error)` — CreateLendingOffer creates a new lending offer.
- `DeleteLendingOffer(ctx context.Context, offerID string) error` — DeleteLendingOffer deletes a lending offer by ID.
- `DeleteLendingOffers(ctx context.Context) error` — DeleteLendingOffers deletes all lending offers created by the current participant.
- `GetLendingAgreement(ctx context.Context, lendingAgreementID string) (*taurusnetwork.LendingAgreement, error)` — GetLendingAgreement retrieves a lending agreement by ID.
- `GetLendingOffer(ctx context.Context, offerID string) (*taurusnetwork.LendingOffer, error)` — GetLendingOffer retrieves a lending offer by ID.
- `ListLendingAgreementAttachments(ctx context.Context, lendingAgreementID string) ([]*taurusnetwork.LendingAgreementAttachment, error)` — ListLendingAgreementAttachments retrieves attachments for a lending agreement.
- `ListLendingAgreements(ctx context.Context, opts *taurusnetwork.ListLendingAgreementsOptions) (*taurusnetwork.ListLendingAgreementsResult, error)` — ListLendingAgreements retrieves a list of lending agreements with optional filtering and pagination.
- `ListLendingAgreementsForApproval(ctx context.Context, opts *taurusnetwork.ListLendingAgreementsForApprovalOptions) (*taurusnetwork.ListLendingAgreementsForApprovalResult, error)` — ListLendingAgreementsForApproval retrieves lending agreements pending approval.
- `ListLendingOffers(ctx context.Context, opts *taurusnetwork.ListLendingOffersOptions) (*taurusnetwork.ListLendingOffersResult, error)` — ListLendingOffers retrieves a list of lending offers with optional filtering and pagination.
- `RepayLendingAgreement(ctx context.Context, lendingAgreementID string, req *taurusnetwork.RepayLendingAgreementRequest) error` — RepayLendingAgreement repays a lending agreement.
- `UpdateLendingAgreement(ctx context.Context, lendingAgreementID string, req *taurusnetwork.UpdateLendingAgreementRequest) error` — UpdateLendingAgreement updates a lending agreement.

### TaurusNetworkParticipantService

- `CreateParticipantAttribute(ctx context.Context, participantID string, req *taurusnetwork.CreateParticipantAttributeRequest) error` — CreateParticipantAttribute creates an attribute for a given participant.
- `DeleteParticipantAttribute(ctx context.Context, participantID string, attributeID string) error` — DeleteParticipantAttribute deletes an attribute for a given participant.
- `GetMyParticipant(ctx context.Context) (*taurusnetwork.GetMyParticipantResult, error)` — GetMyParticipant returns the current participant linked to the tenant with additional Taurus-NETWORK settings.
- `GetParticipant(ctx context.Context, participantID string, opts *taurusnetwork.GetParticipantOptions) (*taurusnetwork.TnParticipant, error)` — GetParticipant returns a Taurus-NETWORK participant based on the provided ID.
- `ListParticipants(ctx context.Context, opts *taurusnetwork.ListParticipantsOptions) (*taurusnetwork.ListParticipantsResult, error)` — ListParticipants returns a list of visible Taurus-NETWORK participants connected to the current participant.

### TaurusNetworkPledgeService

- `AddPledgeCollateral(ctx context.Context, pledgeID string, req *taurusnetwork.AddPledgeCollateralRequest) (*taurusnetwork.AddPledgeCollateralResponse, error)` — AddPledgeCollateral adds collateral to an existing pledge.
- `ApprovePledgeActions(ctx context.Context, req *taurusnetwork.ApprovePledgeActionsRequest) (*taurusnetwork.ApprovePledgeActionsResponse, error)` — ApprovePledgeActions approves one or more pledge actions.
- `CreatePledge(ctx context.Context, req *taurusnetwork.CreatePledgeRequest) (*taurusnetwork.CreatePledgeResponse, error)` — CreatePledge creates a new pledge.
- `GetPledge(ctx context.Context, pledgeID string) (*taurusnetwork.Pledge, error)` — GetPledge retrieves a pledge by ID.
- `InitiateWithdrawPledge(ctx context.Context, pledgeID string, req *taurusnetwork.InitiateWithdrawPledgeRequest) (*taurusnetwork.InitiateWithdrawPledgeResponse, error)` — InitiateWithdrawPledge initiates a withdrawal from a pledge (pledgor initiates).
- `ListPledgeActions(ctx context.Context, opts *taurusnetwork.ListPledgeActionsOptions) ([]taurusnetwork.PledgeAction, *model.CursorPagination, error)` — ListPledgeActions retrieves a list of pledge actions.
- `ListPledgeActionsForApproval(ctx context.Context, opts *taurusnetwork.ListPledgeActionsForApprovalOptions) ([]taurusnetwork.PledgeAction, *model.CursorPagination, error)` — ListPledgeActionsForApproval retrieves a list of pledge actions pending approval.
- `ListPledgeWithdrawals(ctx context.Context, opts *taurusnetwork.ListPledgeWithdrawalsOptions) ([]taurusnetwork.PledgeWithdrawal, *model.CursorPagination, error)` — ListPledgeWithdrawals retrieves a list of pledge withdrawals.
- `ListPledges(ctx context.Context, opts *taurusnetwork.ListPledgesOptions) ([]*taurusnetwork.Pledge, *model.CursorPagination, error)` — ListPledges retrieves a list of pledges with optional filters.
- `RejectPledge(ctx context.Context, pledgeID string, req *taurusnetwork.RejectPledgeRequest) error` — RejectPledge rejects a pledge.
- `RejectPledgeActions(ctx context.Context, req *taurusnetwork.RejectPledgeActionsRequest) error` — RejectPledgeActions rejects one or more pledge actions.
- `Unpledge(ctx context.Context, pledgeID string) (*taurusnetwork.UnpledgeResponse, error)` — Unpledge unpledges funds from a pledge.
- `UpdatePledge(ctx context.Context, pledgeID string, req *taurusnetwork.UpdatePledgeRequest) error` — UpdatePledge updates a pledge's modifiable fields.
- `WithdrawPledge(ctx context.Context, pledgeID string, req *taurusnetwork.WithdrawPledgeRequest) (*taurusnetwork.WithdrawPledgeResponse, error)` — WithdrawPledge creates a withdrawal request from a pledge (pledgee initiates).

### TaurusNetworkSettlementService

- `CancelSettlement(ctx context.Context, settlementID string) error` — CancelSettlement cancels a settlement.
- `CreateSettlement(ctx context.Context, req *taurusnetwork.CreateSettlementRequest) (*taurusnetwork.CreateSettlementResult, error)` — CreateSettlement creates a new settlement.
- `GetSettlement(ctx context.Context, settlementID string) (*taurusnetwork.Settlement, error)` — GetSettlement retrieves a settlement by ID.
- `ListSettlements(ctx context.Context, opts *taurusnetwork.ListSettlementsOptions) (*taurusnetwork.ListSettlementsResult, error)` — ListSettlements retrieves a list of settlements with optional filtering and pagination.
- `ListSettlementsForApproval(ctx context.Context, opts *taurusnetwork.ListSettlementsForApprovalOptions) (*taurusnetwork.ListSettlementsForApprovalResult, error)` — ListSettlementsForApproval retrieves a list of settlements pending approval.
- `ReplaceSettlement(ctx context.Context, settlementID string, req *taurusnetwork.ReplaceSettlementRequest) error` — ReplaceSettlement replaces a settlement with new attributes at the target side before approval.

### TaurusNetworkSharingService

- `ListSharedAddresses(ctx context.Context, opts *taurusnetwork.ListSharedAddressesOptions) (*taurusnetwork.ListSharedAddressesResult, error)` — ListSharedAddresses retrieves a list of shared addresses with optional filtering and pagination.
- `ListSharedAssets(ctx context.Context, opts *taurusnetwork.ListSharedAssetsOptions) (*taurusnetwork.ListSharedAssetsResult, error)` — ListSharedAssets retrieves a list of shared assets with optional filtering and pagination.
- `ShareAddress(ctx context.Context, request *taurusnetwork.ShareAddressRequest) error` — ShareAddress shares an internal address with a Taurus-NETWORK participant.
- `ShareWhitelistedAsset(ctx context.Context, request *taurusnetwork.ShareWhitelistedAssetRequest) error` — ShareWhitelistedAsset shares an asset with a Taurus-NETWORK participant.
- `UnshareAddress(ctx context.Context, sharedAddressID string) error` — UnshareAddress unshares an address with a Taurus-NETWORK participant.
- `UnshareWhitelistedAsset(ctx context.Context, sharedAssetID string) error` — UnshareWhitelistedAsset unshares an asset with a Taurus-NETWORK participant.

### TokenMetadataService

- `GetCryptoPunkMetadata(ctx context.Context, network, contract, tokenID string) (*model.CryptoPunkMetadata, error)` — GetCryptoPunkMetadata retrieves metadata for a CryptoPunk NFT.
- `GetERCTokenMetadata(ctx context.Context, network, contract, tokenID string, opts *model.GetERCTokenMetadataOptions) (*model.ERCTokenMetadata, error)` — GetERCTokenMetadata retrieves metadata for an ERC-721 or ERC-1155 token.

### TransactionService

- `ExportTransactions(ctx context.Context, opts *model.ExportTransactionsOptions) (string, error)` — ExportTransactions exports transactions in the specified format.
- `GetTransaction(ctx context.Context, txID string) (*model.Transaction, error)` — GetTransaction retrieves a transaction by ID.
- `GetTransactionByHash(ctx context.Context, hash string) (*model.Transaction, error)` — GetTransactionByHash retrieves a transaction by its blockchain hash.
- `ListTransactions(ctx context.Context, opts *model.ListTransactionsOptions) ([]*model.Transaction, *model.Pagination, error)` — ListTransactions retrieves a list of transactions.
- `ListTransactionsByAddress(ctx context.Context, address string, opts *model.ListTransactionsByAddressOptions) ([]*model.Transaction, *model.Pagination, error)` — ListTransactionsByAddress retrieves a list of transactions for a specific address.

### UserDeviceService

- `ApprovePairing(ctx context.Context, pairingID string, req *model.ApprovePairingRequest) error` — ApprovePairing approves a user device pairing request (Step 3).
- `CreatePairing(ctx context.Context) (*model.CreatePairingResult, error)` — CreatePairing creates a new user device pairing request (Step 1).
- `GetPairingStatus(ctx context.Context, pairingID string, nonce string) (*model.PairingStatusResult, error)` — GetPairingStatus retrieves the status of a user device pairing request.
- `StartPairing(ctx context.Context, pairingID string, req *model.StartPairingRequest) error` — StartPairing starts a user device pairing request (Step 2).

### UserService

- `CreateUserAttribute(ctx context.Context, userID, key, value string) error` — CreateUserAttribute creates an attribute for a user.
- `GetMe(ctx context.Context) (*model.User, error)` — GetMe retrieves the currently authenticated user.
- `GetUser(ctx context.Context, id string) (*model.User, error)` — GetUser retrieves a user by ID.
- `GetUsersByEmail(ctx context.Context, emails []string) ([]*model.User, error)` — GetUsersByEmail retrieves users by their email addresses.
- `ListUsers(ctx context.Context, opts *model.ListUsersOptions) (*model.ListUsersResult, error)` — ListUsers retrieves a list of users with optional filtering.

### VisibilityGroupService

- `GetUsersByVisibilityGroupID(ctx context.Context, visibilityGroupID string) ([]*model.User, error)` — GetUsersByVisibilityGroupID retrieves users in a visibility group by its ID.
- `ListVisibilityGroups(ctx context.Context) ([]*model.VisibilityGroup, error)` — ListVisibilityGroups retrieves a list of all restricted visibility groups.

### WalletService

- `CreateWallet(ctx context.Context, req *model.CreateWalletRequest) (*model.Wallet, error)` — CreateWallet creates a new wallet.
- `CreateWalletAttribute(ctx context.Context, walletID, key, value string) error` — CreateWalletAttribute creates a custom attribute on a wallet.
- `DeleteWalletAttribute(ctx context.Context, walletID, attributeID string) error` — DeleteWalletAttribute deletes a custom attribute from a wallet.
- `GetWallet(ctx context.Context, walletID string) (*model.Wallet, error)` — GetWallet retrieves a wallet by ID.
- `GetWalletBalanceHistory(ctx context.Context, walletID string, intervalHours int) ([]*model.BalanceHistoryPoint, error)` — GetWalletBalanceHistory retrieves balance history for a wallet.
- `GetWalletTokens(ctx context.Context, walletID string, opts *model.GetWalletTokensOptions) ([]*model.AssetBalance, error)` — GetWalletTokens retrieves token balances for a wallet.
- `ListWallets(ctx context.Context, opts *model.ListWalletsOptions) ([]*model.Wallet, *model.Pagination, error)` — ListWallets retrieves a list of wallets.

### WebhookCallService

- `ListWebhookCalls(ctx context.Context, opts *model.ListWebhookCallsOptions) (*model.ListWebhookCallsResult, error)` — ListWebhookCalls retrieves a list of webhook calls with optional filtering and pagination.

### WebhookService

- `CreateWebhook(ctx context.Context, req *model.CreateWebhookRequest) (*model.CreateWebhookResult, error)` — CreateWebhook creates a new webhook configuration.
- `DeleteWebhook(ctx context.Context, id string) error` — DeleteWebhook deletes a webhook configuration by ID.
- `ListWebhooks(ctx context.Context, opts *model.ListWebhooksOptions) (*model.ListWebhooksResult, error)` — ListWebhooks retrieves a list of webhooks with optional filtering and pagination.
- `UpdateWebhookStatus(ctx context.Context, id string, enabled bool) (*model.Webhook, error)` — UpdateWebhookStatus updates a webhook's status (enabled/disabled).

### WhitelistedAddressService

- `ApproveWhitelistedAddresses(ctx context.Context, ids []string, privateKey *ecdsa.PrivateKey, comment string) error` — ApproveWhitelistedAddresses signs and submits an approval for the given whitelisted
- `GetWhitelistedAddress(ctx context.Context, id string) (*model.WhitelistedAddress, error)` — GetWhitelistedAddress retrieves a whitelisted address by ID.
- `GetWhitelistedAddressEnvelope(ctx context.Context, id string) (*model.WhitelistedAddressEnvelope, error)` — GetWhitelistedAddressEnvelope retrieves a whitelisted address envelope by ID and performs
- `ListWhitelistedAddresses(ctx context.Context, opts *model.ListWhitelistedAddressesOptions) (*model.WhitelistedAddressResult, error)` — ListWhitelistedAddresses retrieves a list of whitelisted addresses, verifying every
- `ListWhitelistedAddressesForApproval(ctx context.Context, opts *model.ListWhitelistedAddressesForApprovalOptions) (*model.WhitelistedAddressResult, error)` — ListWhitelistedAddressesForApproval retrieves whitelisted addresses awaiting approval,

### WhitelistedAssetService

- `ApproveWhitelistedAssets(ctx context.Context, ids []string, privateKey *ecdsa.PrivateKey, comment string) error` — ApproveWhitelistedAssets signs and submits an approval for the given whitelisted
- `GetWhitelistedAsset(ctx context.Context, id string) (*model.WhitelistedAsset, error)` — GetWhitelistedAsset retrieves a whitelisted asset by ID.
- `GetWhitelistedAssetEnvelope(ctx context.Context, id string) (*model.WhitelistedAssetEnvelope, error)` — GetWhitelistedAssetEnvelope retrieves a whitelisted asset envelope by ID and performs
- `ListWhitelistedAssets(ctx context.Context, opts *model.ListWhitelistedAssetsOptions) ([]*model.WhitelistedAsset, *model.Pagination, error)` — ListWhitelistedAssets retrieves a list of whitelisted assets.
- `ListWhitelistedAssetsForApproval(ctx context.Context, opts *model.ListWhitelistedAssetsForApprovalOptions) ([]*model.WhitelistedAsset, *model.Pagination, error)` — ListWhitelistedAssetsForApproval retrieves assets awaiting approval, verified the same

### WhitelistedContractService

- `ApproveWhitelistedContract(ctx context.Context, ids []string, signature string, comment string) error` — ApproveWhitelistedContract approves whitelisted contracts.
- `CreateWhitelistedContract(ctx context.Context, req *model.CreateWhitelistedContractRequest) (string, error)` — CreateWhitelistedContract creates a new whitelisted contract.
- `CreateWhitelistedContractAttribute(ctx context.Context, contractID string, req *model.CreateWhitelistedContractAttributeRequest) ([]model.WhitelistedContractAttribute, error)` — CreateWhitelistedContractAttribute creates an attribute on a whitelisted contract.
- `CreateWhitelistedContractAttributes(ctx context.Context, contractID string, reqs []model.CreateWhitelistedContractAttributeRequest) ([]model.WhitelistedContractAttribute, error)` — CreateWhitelistedContractAttributes creates multiple attributes on a whitelisted contract.
- `DeleteWhitelistedContract(ctx context.Context, id string, comment string) (string, error)` — DeleteWhitelistedContract deletes a whitelisted contract.
- `DeleteWhitelistedContractAttribute(ctx context.Context, contractID string, attributeID string) error` — DeleteWhitelistedContractAttribute deletes an attribute from a whitelisted contract.
- `GetWhitelistedContractAttribute(ctx context.Context, contractID string, attributeID string) (*model.WhitelistedContractAttribute, error)` — GetWhitelistedContractAttribute retrieves an attribute from a whitelisted contract.
- `RejectWhitelistedContract(ctx context.Context, ids []string, comment string) error` — RejectWhitelistedContract rejects whitelisted contracts.
- `UpdateWhitelistedContract(ctx context.Context, id string, req *model.UpdateWhitelistedContractRequest) (string, error)` — UpdateWhitelistedContract updates an existing whitelisted contract.

<!-- END GENERATED METHOD INDEX -->
