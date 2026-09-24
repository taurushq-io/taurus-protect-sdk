# Services Reference

This document provides complete API documentation for all 44 services in the Taurus-PROTECT Python SDK.

## Service Overview

The SDK provides services organized into two categories: core services (39) and TaurusNetwork services (5).

### Core Services

| Service | Access | Purpose |
|---------|--------|---------|
| `WalletService` | `client.wallets` | Wallet management |
| `AddressService` | `client.addresses` | Address management with signature verification |
| `RequestService` | `client.requests` | Transaction requests with ECDSA approval |
| `TransactionService` | `client.transactions` | Transaction queries |
| `BalanceService` | `client.balances` | Balance queries |
| `CurrencyService` | `client.currencies` | Currency information |
| `GovernanceRuleService` | `client.governance_rules` | Governance rules with SuperAdmin verification |
| `WhitelistedAddressService` | `client.whitelisted_addresses` | Address whitelisting |
| `WhitelistedAssetService` | `client.whitelisted_assets` | Asset/contract whitelisting |
| `AuditService` | `client.audits` | Audit log queries |
| `ChangeService` | `client.changes` | Configuration change tracking |
| `FeeService` | `client.fees` | Fee information |
| `PriceService` | `client.prices` | Price data |
| `AirGapService` | `client.air_gap` | Air-gap signing |
| `StakingService` | `client.staking` | Staking operations |
| `ContractWhitelistingService` | `client.contract_whitelisting` | Smart contract whitelisting |
| `BusinessRuleService` | `client.business_rules` | Business rules |
| `ReservationService` | `client.reservations` | Balance reservations |
| `MultiFactorSignatureService` | `client.multi_factor_signature` | Multi-factor signatures |
| `UserService` | `client.users` | User management |
| `GroupService` | `client.groups` | User groups |
| `VisibilityGroupService` | `client.visibility_groups` | Visibility groups |
| `ConfigService` | `client.config` | System configuration |
| `WebhookService` | `client.webhooks` | Webhook management |
| `WebhookCallService` | `client.webhook_calls` | Webhook call history |
| `TagService` | `client.tags` | Tag management |
| `AssetService` | `client.assets` | Asset information |
| `ActionService` | `client.actions` | Action management |
| `BlockchainService` | `client.blockchains` | Blockchain information |
| `EarnService` | `client.earn` | Earn rewards |
| `ExchangeService` | `client.exchanges` | Exchange integration |
| `FiatService` | `client.fiat` | Fiat operations |
| `FeePayerService` | `client.fee_payers` | Fee payer management |
| `HealthService` | `client.health` | API health checks |
| `JobService` | `client.jobs` | Background jobs |
| `ScoreService` | `client.scores` | Risk scoring |
| `StatisticsService` | `client.statistics` | Platform statistics |
| `TokenMetadataService` | `client.token_metadata` | Token metadata |
| `UserDeviceService` | `client.user_devices` | User device management |

### TaurusNetwork Services

| Service | Access | Purpose |
|---------|--------|---------|
| `ParticipantService` | `client.taurus_network.participants` | Participant management |
| `PledgeService` | `client.taurus_network.pledges` | Pledge lifecycle |
| `LendingService` | `client.taurus_network.lending` | Lending offers and agreements |
| `SettlementService` | `client.taurus_network.settlements` | Settlement operations |
| `SharingService` | `client.taurus_network.sharing` | Address/asset sharing |

---

## Core Services

### WalletService

Provides wallet management operations including creation, retrieval, and balance history.

**Access:** `client.wallets`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(wallet_id)` | `wallet_id: int` | `Wallet` | Get wallet by ID |
| `list(limit, offset, exclude_disabled)` | `limit: Optional[int]` (20, max 100), `offset: Optional[int]`, `exclude_disabled: Optional[bool]` | `Tuple[List[Wallet], Pagination]` | List wallets |
| `list_with_options(options)` | `options: ListWalletsOptions` | `Tuple[List[Wallet], Pagination]` | List with every filter (name, sort, tags, balance, chain, ids) |
| `get_by_name(name, limit, offset, exclude_disabled)` | `name: str`, `limit`, `offset`, `exclude_disabled` | `Tuple[List[Wallet], Pagination]` | Find wallets by name |
| `create(request)` | `request: CreateWalletRequest` | `Wallet` | Create wallet |
| `create_wallet(...)` | See below | `Wallet` | Create with explicit params |
| `create_attribute(wallet_id, key, value)` | `wallet_id: int`, `key: str`, `value: str` | `None` | Add attribute |
| `get_balance_history(wallet_id, interval_hours)` | `wallet_id: int`, `interval_hours: int` | `List[BalanceHistoryPoint]` | Get balance history |
| `get_tokens(wallet_id, page_size, cursor)` | `wallet_id: int`, `page_size: Optional[int]`, `cursor: Optional[str]` | `Tuple[List[AssetBalance], CursorPage]` | Get token balances, one page at a time |

#### Example

```python
from taurus_protect.models import CreateWalletRequest

# List wallets (continue with offset=pagination.next_offset while pagination.has_more)
wallets, pagination = client.wallets.list(limit=50)
print(f"Total: {pagination.total_items}")

# Get single wallet
wallet = client.wallets.get(123)
print(f"{wallet.name}: {wallet.balance.total_confirmed}")

# Create wallet
request = CreateWalletRequest(
    blockchain="ETH",
    network="mainnet",
    name="Trading Wallet",
)
wallet = client.wallets.create(request)
```

---

### AddressService

Provides address management with mandatory signature verification.

**Access:** `client.addresses`

**Security:** All addresses retrieved are automatically verified using the rules container public keys when configured.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(address_id)` | `address_id: int` | `Address` | Get address (with verification) |
| `list(wallet_id, limit, offset, exclude_disabled)` | `wallet_id: int`, `limit: Optional[int]` (20, max 100), `offset: Optional[int]`, `exclude_disabled: Optional[bool]` | `Tuple[List[Address], Pagination]` | List addresses |
| `list_with_options(options)` | `options: ListAddressesOptions` | `Tuple[List[Address], Pagination]` | List with filtering; `exclude_disabled` sends `includeDisabledAddresses` |
| `create(request)` | `request: CreateAddressRequest` | `Address` | Create address |
| `create_address(...)` | See below | `Address` | Create with explicit params |
| `create_attribute(address_id, key, value)` | `address_id: int`, `key: str`, `value: str` | `None` | Add attribute |
| `delete_attribute(address_id, attribute_id)` | `address_id: int`, `attribute_id: int` | `None` | Delete attribute |
| `get_proof_of_reserve(address_id, challenge)` | `address_id: int`, `challenge: str` | `Any` | Get proof of reserve |

#### Example

```python
from taurus_protect.models import CreateAddressRequest

# List addresses for wallet
addresses, pagination = client.addresses.list(wallet_id=123, limit=50)

# Create address
request = CreateAddressRequest(
    wallet_id="123",
    label="Customer Deposit",
    comment="Primary deposit address",
)
address = client.addresses.create(request)
print(f"Address: {address.address}")
```

---

### RequestService

Provides transaction request management with ECDSA approval signing.

**Access:** `client.requests`

**Security Features:**
- Hash verification using constant-time comparison
- ECDSA signature-based request approval

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(request_id)` | `request_id: int` | `Request` | Get request (with hash verification) |
| `list(page_size, cursor, ...)` | Filters (dates, currency, statuses, types, ids, sort) | `Tuple[List[Request], CursorPage]` | List requests (v2) |
| `get_for_approval(page_size, cursor, ...)` | Filters (currency, types, exclude_types, ids, sort) | `Tuple[List[Request], CursorPage]` | Get pending approvals (v2); no status filter |
| `approve_requests(requests, private_key, comment)` | `requests: List[Request]`, `private_key: EllipticCurvePrivateKey` | `int` | Approve with signature; refuses metadata whose hash was not verified |
| `approve_request(request, private_key, comment)` | Single request | `int` | Approve single request |
| `reject_requests(request_ids, comment)` | `request_ids: List[int]`, `comment: str` | `None` | Reject requests |
| `reject_request(request_id, comment)` | `request_id: int`, `comment: str` | `None` | Reject single request |
| `create_internal_transfer(...)` | From/to address IDs, amount | `Request` | Create internal transfer |
| `create_internal_transfer_from_wallet(...)` | From wallet ID | `Request` | Create from omnibus wallet |
| `create_external_transfer(...)` | To whitelisted address | `Request` | Create external transfer |
| `create_external_transfer_from_wallet(...)` | From wallet to whitelisted | `Request` | Create external from wallet |
| `create_cancel_request(address_id, nonce)` | `address_id: int`, `nonce: int` | `Request` | Create cancel request |
| `create_incoming_request(...)` | From exchange | `Request` | Create incoming request |

#### Example

```python
from cryptography.hazmat.primitives.serialization import load_pem_private_key

# Load private key
with open("key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=None)

# Get requests pending approval
requests, _ = client.requests.get_for_approval(page_size=10)

# Approve with ECDSA signature
if requests:
    signed_count = client.requests.approve_requests(requests, private_key)
    print(f"Approved {signed_count} request(s)")

# Create internal transfer
request = client.requests.create_internal_transfer(
    from_address_id=123,
    to_address_id=456,
    amount="1000000000000000000",  # 1 ETH in wei
)
```

---

### TransactionService

Provides transaction query operations.

**Access:** `client.transactions`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(transaction_id)` | `transaction_id: int` | `Transaction` | Get transaction |
| `list(...)` | Multiple filters | `Tuple[List[Transaction], Pagination]` | List transactions |
| `list_by_address(address, limit, offset)` | `address: str` | `Tuple[List[Transaction], Pagination]` | List an address's transactions |
| `export(..., limit, format)` | Export filters | `TransactionExport` | Export in one reply (no offset: the server always starts at row 0) |
| `export_csv(..., limit)` | Export filters | `TransactionExport` | `export` with `format="csv"` |

#### Example

```python
# List recent transactions
transactions, _ = client.transactions.list(limit=100)
for tx in transactions:
    print(f"{tx.hash}: {tx.amount} {tx.currency}")
```

---

### BalanceService

Provides balance query operations.

**Access:** `client.balances`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(currency, page_size, cursor, token_id)` | Filters | `Tuple[List[AssetBalance], CursorPage]` | List balances (request cursor; the page carries the total) |
| `list_nft_collections(blockchain, network, page_size, cursor, ...)` | Filters | `Tuple[List[NFTCollectionBalance], CursorPage]` | List NFT collection balances |

---

### GovernanceRuleService

Provides governance rules with SuperAdmin signature verification.

**Access:** `client.governance_rules`

**Security:** Automatically verifies SuperAdmin signatures when keys are configured.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_rules()` | None | `Optional[GovernanceRules]` | Get current governance rules |
| `get_rules_by_id(rules_id)` | `rules_id: str` | `Optional[GovernanceRules]` | Get governance rules by ID |
| `get_rules_proposal()` | None | `Optional[GovernanceRules]` | Get proposed governance rules |
| `get_rules_history(page_size, cursor)` | `page_size: Optional[int]` (20, max 100), `cursor: Optional[str]` | `GovernanceRulesHistoryResult` | Get governance rules history; `result.page` continues the walk |
| `get_public_keys()` | None | `List[SuperAdminPublicKey]` | Get SuperAdmin public keys |
| `get_decoded_rules_container(rules)` | `rules: GovernanceRules` | `DecodedRulesContainer` | Decode rules container from governance rules |
| `verify_governance_rules(rules)` | `rules: GovernanceRules` | `GovernanceRules` | Verify SuperAdmin signatures on rules |
| `update_rules_proposal(container)` | `container: DecodedRulesContainer` | `None` | Submit a typed container as a governance proposal (SuperAdmin only) |
| `approve_rules_proposal(private_key, comment)` | `private_key: EllipticCurvePrivateKey`, `comment: str = ""` | `None` | Sign the pending proposal and submit the approval |
| `reject_rules_proposal(comment)` | `comment: str = ""` | `None` | Reject the pending rules proposal |

`update_rules_proposal` encodes the container internally and strips the server-controlled
`enforced_rules_hash` and `timestamp`. The endpoint returns no body — call
`get_rules_proposal()` to read the persisted proposal (rules reads are cached server-side,
so an immediate read-back may be stale).

`approve_rules_proposal` signs the pending proposal's rules container (SHA-256 + P-256
ECDSA, base64 raw r||s). `expected_container_hash` PINS the content: pass
`proposal_container_hash()` of the proposal you reviewed, and the call aborts without
signing if the re-fetched container differs. Review with `get_rules_proposal()` +
`decode_proposal_for_review()` — **not** `get_decoded_rules_container()`, which verifies
unconditionally and so always fails on a pending proposal (it legitimately carries 0..N
signatures; signing IS the approval step).

**Typed governance models:** `DecodedRulesContainer` is a lossless typed container that
round-trips through `rules_container_to_bytes` / `rules_container_from_base64`. Cells are
decoded into the typed `RuleCell` union (`FiatAmountAny`, `FiatAmountRange`,
`SourceInternalWallet`, `StringEqualValue`, ...), with `RawCell` preserving cells from
newer schemas verbatim.

Cross-SDK cell wire-format parity is pinned by the shared golden vectors at
`scripts/resources/governance-cell-vectors.json` (monorepo root), consumed by every SDK's
test suite.

#### Example

```python
# Get current governance rules
rules = client.governance_rules.get_rules()
if rules:
    print(f"Rules locked: {rules.locked}")

    # Decode the rules container
    decoded = client.governance_rules.get_decoded_rules_container(rules)
```

---

### WhitelistedAddressService

Provides whitelisted address management with 6-step verification.

**Access:** `client.whitelisted_addresses`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(address_id)` | `address_id: int` | `WhitelistedAddress` | Get with verification |
| `list(...)` | Filters | `WhitelistedAddressListResult` | List addresses; carries `excluded_unverified` |
| `list_for_approval(...)` | Filters | `WhitelistedAddressListResult` | Rows awaiting approval, verified as in `list` |
| `approve(selection, private_key, comment)` | `selection: WhitelistedAddressApproval`, `private_key`, `comment: str` | `None` | Sign an approval, all-or-nothing, pinned to the reviewed rows |

**`approve` takes the ROWS a verified read returned, not bare ids.** Mint the selection with
`result.select(*ids)` or `result.select_all()` off a `list_for_approval` result; it carries
the metadata hash each row had **at review time**, and the approval aborts if the re-read
hash differs. See the note under `WhitelistedAssetService.approve` for the substitution
attack the pin defeats — it is the same one, and verification alone does not catch it.

```python
result = client.whitelisted_addresses.list_for_approval()
selection = result.select_all()
client.whitelisted_addresses.approve(selection, private_key, "reviewed")
```

---

### WhitelistedAssetService

Provides whitelisted asset/contract management with 5-step verification plus the parse from
the verified payload. **This is the only verified reader of `/whitelists/contracts`** — a
whitelisted asset and a whitelisted contract are one server entity, and
`ContractWhitelistingService` is write-only.

**Access:** `client.whitelisted_assets`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(asset_id)` | `asset_id: int` | `WhitelistedAsset` | Get with verification |
| `list(...)` | Filters | `Tuple[List[WhitelistedAsset], Pagination]` | List assets; a short page is not the end, follow `has_more` |
| `list_for_approval(ids, limit, offset)` | `ids: Optional[List[str]]`, `limit: Optional[int]` (20, max 100), `offset: Optional[int]` | `Tuple[List[WhitelistedAsset], Pagination]` | List assets awaiting approval, verified as in `list` |
| `approve(selection, private_key, comment)` | `selection: WhitelistedAssetApproval`, `private_key`, `comment: str` | `None` | Sign an approval, all-or-nothing, pinned to the reviewed rows |

`list_for_approval` exists here because the for-approval read used to live only on the
unverified contract service, so the rows an approver inspects were never checked against
governance.

`approve` re-reads and verifies each asset and signs the hashes *those* rows carry. Any row
that is missing or fails verification aborts the whole call and nothing is signed: the API
takes one signature covering the whole batch, so a partial approval would mean the caller
believes they approved more than they did.

**`approve` takes the ROWS a verified read returned, not bare ids.** Mint the selection with
`WhitelistedAssetApproval.select(assets, *ids)` or `.select_all(assets)`; it carries the
metadata hash each row had **at review time**, and the approval aborts if the re-read hash
differs. Without that pin a response-controlling server could answer the id-filtered
re-read with a *different* row — one whose existing signatures already satisfy the container
it presents — and harvest a genuine approver signature over content the approver never saw.
Verification alone does not catch it: the substituted row is a real, validly-signed entry,
just not the one that was reviewed. An empty selection raises rather than meaning "approve
nothing".

```python
assets, _ = client.whitelisted_assets.list_for_approval()
selection = WhitelistedAssetApproval.select(assets, *[a.id for a in assets])
client.whitelisted_assets.approve(selection, private_key, "reviewed")
```

---

### CurrencyService

Provides currency information including native currencies and tokens.

**Access:** `client.currencies`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(show_disabled, include_logo)` | `show_disabled: bool = False`, `include_logo: bool = False` | `List[Currency]` | List all currencies |
| `get(currency_id)` | `currency_id: str` | `Currency` | Get currency by ID |
| `get_by_blockchain(blockchain, network, contract_address, token_id)` | `blockchain: str`, `network: str`, `contract_address: Optional[str]`, `token_id: Optional[str]` | `Currency` | Get currency by blockchain/network |
| `get_base_currency()` | None | `Currency` | Get tenant's base currency |

#### Example

```python
# List all enabled currencies
currencies = client.currencies.list()
for currency in currencies:
    print(f"{currency.symbol}: {currency.name}")

# Get currency by blockchain
eth = client.currencies.get_by_blockchain("ETH", "mainnet")
print(f"{eth.name} has {eth.decimals} decimals")

# Get base currency (e.g., USD, EUR, CHF)
base = client.currencies.get_base_currency()
print(f"Base currency: {base.symbol}")
```

---

### AuditService

Provides audit event query operations.

**Access:** `client.audits`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(page_size, cursor, ...)` | Filters (user, entities, actions, dates, sort) | `Tuple[List[Audit], CursorPage]` | List audit events |
| `get(audit_id)` | `audit_id: str` | `Audit` | Get audit event by ID (walks the pages) |

#### Example

```python
# List audit events
audits, page = client.audits.list(page_size=50)
for audit in audits:
    print(f"{audit.id}: {audit.description}")
```

---

### ChangeService

Provides configuration change tracking and approval operations.

**Access:** `client.changes`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(options)` | `options: ListChangesOptions` | `ChangeResult` | List changes; `result.page` continues the walk |
| `list_for_approval(options)` | `options: ListChangesOptions` | `ChangeResult` | List changes pending approval |
| `get(change_id)` | `change_id: str` | `Change` | Get change by ID |
| `approve_change(change_id)` | `change_id: str` | `None` | Approve a change |
| `approve_changes(change_ids)` | `change_ids: List[str]` | `None` | Approve multiple changes |
| `reject_change(change_id)` | `change_id: str` | `None` | Reject a change |
| `reject_changes(change_ids)` | `change_ids: List[str]` | `None` | Reject multiple changes |

#### Example

```python
# List changes
result = client.changes.list(ListChangesOptions(page_size=50))
for change in result.changes:
    print(f"{change.id}: {change.description}")

# Approve a change
client.changes.approve_change("change-123")
```

---

### FeeService

Provides transaction fee estimation operations.

**Access:** `client.fees`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `estimate(currency, amount, destination)` | `currency: str`, `amount: Optional[str]`, `destination: Optional[str]` | `FeeEstimate` | Estimate fee for a transaction |
| `list()` | None | `List[FeeEstimate]` | List fee estimates for all currencies |

#### Example

```python
# Estimate fee for an ETH transfer
estimate = client.fees.estimate(currency="ETH", amount="1.5")
print(f"Low: {estimate.fee_low}, Medium: {estimate.fee_medium}, High: {estimate.fee_high}")

# List all fee estimates
fees = client.fees.list()
for fee in fees:
    print(f"{fee.currency}: {fee.fee_medium}")
```

---

### PriceService

Provides cryptocurrency price data.

**Access:** `client.prices`

`rate` and `decimals` feed amount conversion, so an unverified price is a wrong number a
caller acts on. `get_current` lists through `QueryPricesV2` and verifies each price against the `PRICEUPDATER` keys in the
SuperAdmin-verified rules container, which is why the service takes the cache as a mandatory
parameter. Whether prices must be signed is the **container's** call: no `PRICEUPDATER`
configured means this tenant does not sign prices and the price passes through; a
`PRICEUPDATER` configured plus a price carrying no signatures raises `IntegrityError`.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_current(*, from_currency_id, to_currency_ids, only_primary, sort_order, page_size, cursor)` | Keyword-only filters | `Tuple[List[Price], CursorPage]` | List current prices (from / fromTo / to filter) |
| `get_historical(base_currency, quote_currency, limit)` | `base_currency: str`, `quote_currency: str`, `limit: Optional[int]` (20, max 365) | `List[PriceHistoryPoint]` | Get historical prices (one reply) |

#### Example

```python
# Current prices quoted from one currency
prices, page = client.prices.get_current(from_currency_id="<currency-id>", page_size=100)
for price in prices:
    print(f"{price.currency_from}/{price.currency_to}: {price.rate}")

# Get BTC/USD price history
history = client.prices.get_historical("BTC", "USD", limit=100)
for point in history:
    print(f"{point.timestamp}: {point.rate}")
```

---

### AirGapService

Provides air-gap (cold storage) signing operations for high-security environments.

**Access:** `client.air_gap`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_unsigned_payload(request_id)` | `request_id: int` | `bytes` | Get unsigned payload for offline signing |
| `submit_signed_payload(request_id, signed_payload)` | `request_id: int`, `signed_payload: bytes` | `None` | Submit signed payload |

#### Example

```python
# Get unsigned payload for offline signing
payload = client.air_gap.get_unsigned_payload(request_id=123)
print(f"Payload size: {len(payload)} bytes")

# After signing offline, submit the signed payload
client.air_gap.submit_signed_payload(request_id=123, signed_payload=signed_data)
```

---

### StakingService

Provides multi-chain staking information including validators and staking positions.

**Access:** `client.staking`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list_validators(blockchain, network, *, ids)` | `blockchain: str`, `network: str = "mainnet"`, `ids: Optional[List[str]]` | `List[Validator]` | Every validator (ETH only; the endpoint does not page) |
| `get_staking_info(address_id)` | `address_id: int` | `StakingInfo` | Get staking info for address |

#### Example

```python
# List ETH validators
validators = client.staking.list_validators(blockchain="ETH")
for v in validators:
    print(f"{v.name}: {v.commission}% commission")

# Get staking info for an address
info = client.staking.get_staking_info(address_id=123)
print(f"Staked: {info.staked_amount}, Rewards: {info.rewards}")
```

---

### ContractWhitelistingService

Provides smart contract whitelisting **WRITE** operations.

**Access:** `client.contract_whitelisting`

> **Reads live on `client.whitelisted_assets`.** A whitelisted contract and a whitelisted
> asset are one server entity (`/whitelists/contracts`); the `get`/`list` that used to sit here
> returned the DTO with no verification, which made the verified reader avoidable.
> There is no `delete`: the endpoint is deprecated with no replacement and cannot succeed.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `create(address, name, blockchain, network, abi)` | `address: str`, `name: str`, `blockchain: str`, `network: Optional[str]`, `abi: Optional[str]` | `int` | Create whitelisted contract |
| `approve_whitelisted_contracts(contract_ids, signature, comment)` | `contract_ids: List[str]`, `signature: str`, `comment: Optional[str]` | `None` | **Deprecated** — signature is an opaque blob over unverified hashes; use `whitelisted_assets.approve` |
| `create_attribute(contract_id, key, value)` | `contract_id: str`, `key: str`, `value: str` | `None` | Add attribute |
| `get_attribute(contract_id, key)` | `contract_id: str`, `key: str` | `Optional[str]` | Get attribute value |

#### Example

```python
# Create a whitelisted contract
contract_id = client.contract_whitelisting.create(
    address="0x1234...",
    name="USDC Token",
    blockchain="ETH",
)

# Read it back through the verified reader
assets, pagination = client.whitelisted_assets.list(blockchain="ETH")
for asset in assets:
    print(f"{asset.name}: {asset.contract_address}")
```

---

### BusinessRuleService

Provides business rule management for custom transaction validation logic.

**Access:** `client.business_rules`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(page_size, current_page, page_request, ..., cursor)` | Filters | `BusinessRuleResult` | List business rules (v2); `result.page` continues the walk |
| `list_by_wallet(wallet_id, ...)` | `wallet_id: int` | `BusinessRuleResult` | Rules of one wallet |
| `list_by_currency(currency_id, ...)` | `currency_id: str` | `BusinessRuleResult` | Rules of one currency |

#### Example

```python
# Walk every business rule
cursor = None
while True:
    result = client.business_rules.list(page_size=100, cursor=cursor)
    for rule in result.rules:
        print(f"{rule.rule_key}: {rule.rule_value}")
    if not result.page.has_more:
        break
    cursor = result.page.next_cursor
```

---

### ReservationService

Provides balance reservation management for pending transactions.

**Access:** `client.reservations`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(reservation_id)` | `reservation_id: int` | `Reservation` | Get reservation by ID |
| `list(page_size, cursor, *, kind, kinds, address, address_id)` | Filters | `Tuple[List[Reservation], CursorPage]` | List reservations |

#### Example

```python
# List the reservations on an address
reservations, page = client.reservations.list(address_id="123")
for r in reservations:
    print(f"{r.id}: {r.amount} {r.currency} ({r.kind})")
```

---

### MultiFactorSignatureService

Provides multi-factor signature operations: a **second approval channel** over the same
entities the SDK otherwise protects (a request, a whitelisted address, a whitelisted
contract). The caller is the second-factor signing device, and the signature it submits is
precisely the artefact a compromised server cannot forge on its own.

> **Security — `payload_to_sign` is UNVERIFIED server data.**
> `get_multi_factor_signature_info` returns the bytes the server asks you to sign, and the
> SDK cannot check them: the reply carries only `{id, payload_to_sign[], entity_type}` with
> **no entity id**, so nothing in it can be joined back to the entities the request was
> created for. A compromised server can therefore answer with the metadata hash of an
> entity of its choosing under the expected kind.
>
> **Bind it yourself.** `create_multi_factor_signatures` takes the entity IDs — keep them,
> re-read those entities through the verifying reader for that kind (`client.requests`,
> `client.whitelisted_addresses`, `client.whitelisted_assets`) and require each
> `payload_to_sign` element to equal the locally recomputed, verified metadata hash.
> Tracked in `TODOS.md`.

**Access:** `client.multi_factor_signature`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_multi_factor_signature_info(id)` | `id: str` | `MultiFactorSignatureInfo` | Get the signature request info. **`payload_to_sign` is unverified** |
| `create_multi_factor_signatures(entity_ids, entity_type)` | `entity_ids: List[str]`, `entity_type: Union[MultiFactorSignatureEntityType, str]` | `MultiFactorSignatureResult` | Create a batch of signature requests |
| `approve_multi_factor_signature(id, signature, comment)` | `id: str`, `signature: str`, `comment: str = ""` | `MultiFactorSignatureApprovalResult` | Submit a caller-produced signature. **Opaque to the SDK** |
| `reject_multi_factor_signature(id, comment)` | `id: str`, `comment: str = ""` | `None` | Reject a signature request |

#### Example

```python
from taurus_protect.models import MultiFactorSignatureEntityType

# 1. Create the batch, and KEEP the entity ids -- they are the only thing that can bind
#    the payload the server later asks you to sign to entities you can verify.
entity_ids = ["123", "124"]
batch = client.multi_factor_signature.create_multi_factor_signatures(
    entity_ids, MultiFactorSignatureEntityType.REQUEST
)

# 2. Read the payload -- UNVERIFIED.
info = client.multi_factor_signature.get_multi_factor_signature_info(batch.id)

# 3. Bind it before signing: re-read each entity through the VERIFYING reader and require
#    the payload to be a hash it produced. Skip this and you sign bytes the server chose.
verified_hashes = set()
for entity_id in entity_ids:
    request = client.requests.get(int(entity_id))   # verifies the metadata hash
    if request.metadata and request.metadata.hash:
        verified_hashes.add(request.metadata.hash)
for payload in info.payload_to_sign:
    if payload not in verified_hashes:
        raise RuntimeError(f"refusing to sign an unrecognised payload: {payload}")

# 4. Sign with the second-factor key and submit.
result = client.multi_factor_signature.approve_multi_factor_signature(
    id=info.id, signature=my_second_factor_signature, comment="reviewed"
)
print(f"signatures now: {result.signature_count}")
```

---

### UserService

Provides user management operations.

**Access:** `client.users`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(user_id)` | `user_id: str` | `User` | Get user by ID |
| `get_current()` | None | `User` | Get current authenticated user |
| `list(limit, offset)` | `limit: Optional[int]` (20, max 100), `offset: Optional[int]` | `Tuple[List[User], Pagination]` | List users |
| `get_users_by_email(emails)` | `emails: List[str]` | `List[User]` | Get users by email addresses (reads every page) |
| `create_user_attribute(user_id, key, value)` | `user_id: str`, `key: str`, `value: str` | `None` | Create user attribute |

#### Example

```python
# Get current user
me = client.users.get_current()
print(f"Logged in as: {me.email}")

# List users
users, pagination = client.users.list(limit=50)
for user in users:
    print(f"{user.email}: {user.status}")
```

---

### GroupService

Provides user group management operations.

**Access:** `client.groups`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(group_id)` | `group_id: str` | `Group` | Get group by ID |
| `list(limit, offset)` | `limit: Optional[int]` (20, max 100), `offset: Optional[int]` | `Tuple[List[Group], Pagination]` | List groups |

#### Example

```python
# List groups
groups, pagination = client.groups.list()
for group in groups:
    print(f"{group.name}: {len(group.users)} users")
```

---

### VisibilityGroupService

Provides visibility group management for controlling resource access.

**Access:** `client.visibility_groups`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(group_id)` | `group_id: str` | `VisibilityGroup` | Get visibility group by ID |
| `list()` | None | `List[VisibilityGroup]` | Every visibility group (the endpoint does not page) |
| `get_users(group_id)` | `group_id: str` | `List[Any]` | Get users in a visibility group |

#### Example

```python
# List visibility groups
groups = client.visibility_groups.list()
for group in groups:
    print(f"{group.name}: {group.user_count} users")
```

---

### ConfigService

Provides tenant configuration and feature flag operations.

**Access:** `client.config`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get()` | None | `TenantConfig` | Get tenant configuration |
| `get_features()` | None | `List[Feature]` | Get enabled features |

#### Example

```python
# Get tenant configuration
config = client.config.get()
print(f"Tenant: {config.tenant_id}")

# Get enabled features
features = client.config.get_features()
for feature in features:
    print(f"{feature.name}: {'enabled' if feature.enabled else 'disabled'}")
```

---

### WebhookService

Provides webhook management for event notifications.

**Access:** `client.webhooks`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(page_size, cursor, *, type, url, sort_order)` | Filters | `Tuple[List[Webhook], CursorPage]` | List webhooks |
| `get(webhook_id)` | `webhook_id: str` | `Webhook` | Get webhook by ID (walks the pages) |
| `create(url, events)` | `url: str`, `events: List[str]` | `Webhook` | Create webhook |
| `delete(webhook_id)` | `webhook_id: str` | `None` | Delete webhook |

#### Example

```python
# Create a webhook
webhook = client.webhooks.create(
    url="https://example.com/webhook",
    events=["REQUEST_CREATED", "REQUEST_APPROVED"],
)
print(f"Created webhook: {webhook.id}")

# List webhooks
webhooks, page = client.webhooks.list()
for wh in webhooks:
    print(f"{wh.id}: {wh.url} ({wh.status})")

# Delete a webhook
client.webhooks.delete("webhook-123")
```

---

### WebhookCallService

Provides webhook call history and delivery status tracking.

**Access:** `client.webhook_calls`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_webhook_calls(event_id, webhook_id, status, sort_order, cursor, page_size)` | Optional filters; `cursor` is a `next_cursor` or an `ApiRequestCursor` | `WebhookCallResult` | Get webhook calls; `result.page` continues the walk |
| `list(webhook_id, event_id, status, sort_order, page_size, cursor)` | Optional filters | `Tuple[List[WebhookCall], CursorPage]` | List webhook calls |
| `get(call_id)` | `call_id: str` | `WebhookCall` | Get webhook call by ID (walks the pages) |

#### Example

```python
# Get webhook call history
result = client.webhook_calls.get_webhook_calls(webhook_id="webhook-123")
for call in result.calls:
    print(f"{call.id}: {call.status}")

# List failed calls
calls, page = client.webhook_calls.list(status="FAILED", page_size=50)
```

---

### TagService

Provides tag management for organizing entities.

**Access:** `client.tags`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(tag_id)` | `tag_id: str` | `Tag` | Get tag by ID |
| `list(*, query, ids)` | `query: Optional[str]`, `ids: Optional[List[str]]` | `List[Tag]` | Every matching tag (the endpoint does not page) |
| `create(name, color)` | `name: str`, `color: str` | `Tag` | Create tag |
| `delete(tag_id)` | `tag_id: str` | `None` | Delete tag |

#### Example

```python
# Create a tag
tag = client.tags.create("Important", "#FF0000")
print(f"Created tag: {tag.id}")

# List tags
tags = client.tags.list()
for tag in tags:
    print(f"{tag.name}: {tag.color}")
```

---

### AssetService

Provides asset information and balance queries across wallets and addresses.

**Access:** `client.assets`

`get_addresses` verifies every address's HSM signature — the same check `AddressService` runs
— and **fails fast** on the first that does not verify. It returns the same entity, so
returning it unverified made `AddressService`'s mandatory verification avoidable. It also
returns domain `Address` objects rather than the raw generated DTOs it used to.

The rows of `query_asset_addresses` (v2) carry no signature, so each page is confirmed through
the verified readers: an `ADDRESS_TYPE_V2_INTERNAL` row by re-reading its `address_id` through
the HSM-verified managed-address list (at most 50 ids per request), an
`ADDRESS_TYPE_V2_WHITELISTED` row by re-reading its `whitelisted_address_id` through the
6-step whitelist list (at most 100 ids per request). A row is kept only when the verified
address with that id has the same address string; it then carries `verified=True` and the
verified address. Any other row (EXTERNAL, untyped) is returned with `verified=False`: on-chain
data, never a Taurus-PROTECT address. A row that cannot be confirmed is excluded and named in
`excluded_unverified` (`{id, reason}`); rows came back but none survived raises
`IntegrityError`, and a verified reader's own failure aborts the call. Exclusions never move
the cursor, and a page with no INTERNAL/WHITELISTED rows makes no extra request.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(currency, page_size, cursor)` | `currency: str` | `Tuple[List[Wallet], CursorPage]` | The wallets holding an asset (same as `get_wallets`) |
| `get(asset_id)` | `asset_id: str` | `Asset` | Get asset by currency; `NotFoundError` when no address holds it |
| `get_wallets(currency, page_size, cursor, *, wallet_id, wallet_name)` | `currency: str` | `Tuple[List[Wallet], CursorPage]` | Wallets holding an asset (request cursor, with total) |
| `get_addresses(currency, page_size, cursor, *, wallet_id, address_id, addresses)` | `currency: str` | `Tuple[List[Address], CursorPage]` | Addresses holding an asset, HSM-verified |
| `query_assets(*, blockchain, network, symbol, contract_address, label, currency_name, page_size, cursor)` | Filters | `Tuple[List[AssetV2], CursorPage]` | Asset definitions (v2) |
| `query_asset_addresses(asset_id, *, address_type, kyc_status, page_size, cursor)` | `asset_id: str` | `QueryAssetAddressesResult` | Addresses holding a v2 asset; INTERNAL/WHITELISTED rows verified, the rest `verified=False` |
| `list_asset_operations(asset_id, *, type, status, page_size, cursor)` | `asset_id: str` | `Tuple[List[AssetOperationV2], CursorPage]` | Operations of a v2 asset |

#### Example

```python
# Wallets holding ETH
wallets, page = client.assets.get_wallets("ETH", page_size=50)
for wallet in wallets:
    print(f"{wallet.name}: {wallet.balance}")

# v2 asset definitions
assets, page = client.assets.query_assets(blockchain="CANTON", page_size=50)

# Holders of a v2 asset: only verified rows are Taurus-PROTECT addresses
holders = client.assets.query_asset_addresses(assets[0].id, page_size=50)
for holder in holders.addresses:
    print(holder.address, holder.address_type, holder.verified)
for excluded in holders.excluded_unverified:
    print(f"withheld {excluded.id}: {excluded.reason}")

# Get asset by ID
asset = client.assets.get("BTC")
print(f"Decimals: {asset.decimals}")
```

---

### ActionService

Provides action management operations for pending or completed operations.

**Access:** `client.actions`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(action_id)` | `action_id: str` | `Action` | Get action by ID |
| `list(limit, offset)` | `limit: Optional[int]` (20, max 100), `offset: Optional[int]` | `Tuple[List[Action], Pagination]` | List actions |

#### Example

```python
# List actions
actions, pagination = client.actions.list(limit=50)
for action in actions:
    print(f"{action.label}: {action.status}")
```

---

### BlockchainService

Provides blockchain information and supported network queries.

**Access:** `client.blockchains`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(blockchain, network, include_block_height)` | `blockchain: Optional[str]`, `network: Optional[str]`, `include_block_height: bool = False` | `List[Blockchain]` | List supported blockchains |
| `get(blockchain, network, include_block_height)` | `blockchain: str`, `network: str = "mainnet"`, `include_block_height: bool = False` | `Blockchain` | Get blockchain info |
| `get_by_id(blockchain_id)` | `blockchain_id: str` | `Blockchain` | Get blockchain by composite ID (e.g., "ETH_mainnet") |

#### Example

```python
# List all blockchains
blockchains = client.blockchains.list()
for bc in blockchains:
    print(f"{bc.name} ({bc.network}): {bc.native_currency}")

# Get specific blockchain with block height
blockchain = client.blockchains.get("ETH", "mainnet", include_block_height=True)
print(f"Block height: {blockchain.block_height}")
```

---

### ExchangeService

Provides exchange account integration and management.

**Access:** `client.exchanges`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(page_size, cursor, currency_id, exchange_label, status, only_positive_balance, *, sort_order, include_base_currency_valuation)` | Multiple filters | `Tuple[List[Exchange], CursorPage]` | List exchange accounts |
| `get(exchange_id)` | `exchange_id: str` | `Exchange` | Get exchange account by ID |
| `list_counterparties()` | None | `List[Any]` | List exchange counterparties |
| `get_withdrawal_fee(exchange_id, to_address_id, amount)` | `exchange_id: str`, `to_address_id: Optional[str]`, `amount: Optional[str]` | `Any` | Get withdrawal fees |

#### Example

```python
# List exchange accounts
exchanges, page = client.exchanges.list(page_size=50)
for exchange in exchanges:
    print(f"{exchange.name} ({exchange.exchange_label}): {exchange.balance}")

# Get withdrawal fee
fee = client.exchanges.get_withdrawal_fee(exchange_id="123")
```

---

### EarnService

Provides earn rewards.

**Access:** `client.earn`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list_rewards(recipient_address_id, page_size, cursor)` | `recipient_address_id: Optional[str]` | `Tuple[List[EarnReward], CursorPage]` | List earn rewards |

#### Example

```python
rewards, page = client.earn.list_rewards(page_size=100)
for reward in rewards:
    print(f"{reward.recipient_address}: {reward.amount} {reward.token_symbol}")
```

---

### FiatService

Provides fiat currency operations including provider accounts and exchange rates.

**Access:** `client.fiat`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list_fiat_provider_accounts(provider, label, page_size, cursor, *, account_type, sort_order)` | `provider: str`, `label: str` | `Tuple[List[FiatProviderAccount], CursorPage]` | List a provider's accounts |
| `list_fiat_provider_entities(provider, label, page_size, cursor, *, sort_order)` | Optional filters | `Tuple[List[FiatProviderEntity], CursorPage]` | List fiat provider entities |
| `get_account(account_id)` | `account_id: str` | `FiatProviderAccount` | Get fiat provider account |
| `get_base_currency()` | None | `FiatCurrency` | Get configured base currency |
| `get_rate(from_currency, to_currency)` | `from_currency: str`, `to_currency: str` | `ExchangeRate` | Get exchange rate |
| `list_providers()` | None | `List[Any]` | List fiat providers |

#### Example

```python
# List a provider's fiat accounts
accounts, page = client.fiat.list_fiat_provider_accounts("provider", "label")
for account in accounts:
    print(f"{account.name}: {account.balance} {account.currency_code}")

# Get exchange rate
rate = client.fiat.get_rate("USD", "EUR")
print(f"1 USD = {rate.rate} EUR")
```

---

### FeePayerService

Provides fee payer management for sponsored transactions.

**Access:** `client.fee_payers`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(limit, offset, blockchain, network)` | `limit: Optional[int]` (20, max 100), `offset: Optional[int]`, `blockchain: Optional[str]`, `network: Optional[str]` | `Tuple[List[FeePayer], Pagination]` | List fee payers |
| `get(fee_payer_id)` | `fee_payer_id: str` | `FeePayer` | Get fee payer by ID |

#### Example

```python
# List fee payers
fee_payers, pagination = client.fee_payers.list(blockchain="SOL")
for fp in fee_payers:
    print(f"{fp.id}: {fp.address} ({fp.balance})")
```

---

### HealthService

Provides API health check operations.

**Access:** `client.health`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `check()` | None | `HealthStatus` | Check API health |
| `get_all_health_checks(tenant_id, fail_if_unhealthy)` | `tenant_id: Optional[str]`, `fail_if_unhealthy: bool = False` | `GetAllHealthChecksResult` | Get all health checks with optional filtering |

#### Example

```python
# Check API health
health = client.health.check()
print(f"API Status: {health.status}")

# Get all health checks
all_checks = client.health.get_all_health_checks()
for component_name, component in (all_checks.components or {}).items():
    print(f"Component: {component_name}")
```

---

### JobService

Provides background job management and monitoring.

**Access:** `client.jobs`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list()` | None | `List[Job]` | Every job (the endpoint does not page) |
| `get(job_id)` | `job_id: str` | `Job` | Get job by ID |

#### Example

```python
# List jobs
jobs = client.jobs.list()
for job in jobs:
    print(job.name)
```

---

### ScoreService

Provides address and transaction risk scoring from compliance/AML providers.

**Access:** `client.scores`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_address_score(address_id, provider)` | `address_id: str`, `provider: Optional[str]` | `List[Score]` | Get risk scores for an address |
| `get_transaction_score(tx_hash, provider)` | `tx_hash: str`, `provider: Optional[str]` | `List[Score]` | Get risk scores for a transaction |
| `refresh_whitelisted_address_score(address_id, provider)` | `address_id: str`, `provider: Optional[str]` | `List[Score]` | Refresh scores for whitelisted address |

#### Example

```python
# Get risk score for an address
scores = client.scores.get_address_score(address_id="123")
for score in scores:
    print(f"{score.provider}: {score.score}")

# Refresh scores for a whitelisted address
scores = client.scores.refresh_whitelisted_address_score(
    address_id="456", provider="chainalysis"
)
```

---

### StatisticsService

Provides portfolio and transaction statistics.

**Access:** `client.statistics`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_summary()` | None | `Optional[PortfolioStatistics]` | Get portfolio statistics |
| `get_transaction_stats(from_date, to_date)` | `from_date: Optional[datetime]`, `to_date: Optional[datetime]` | `TransactionStatistics` | Get transaction statistics |

#### Example

```python
# Get portfolio summary
summary = client.statistics.get_summary()
print(f"Total wallets: {summary.wallets_count}")
print(f"Total addresses: {summary.addresses_count}")
```

---

### TokenMetadataService

Provides token metadata for ERC tokens, FA tokens, and CryptoPunks.

**Access:** `client.token_metadata`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(blockchain, contract_address, token_id, network, with_data)` | `blockchain: str`, `contract_address: str`, `token_id: str = "0"`, `network: str = "mainnet"`, `with_data: bool = False` | `Optional[TokenMetadata]` | Get ERC token metadata |
| `get_erc(network, contract_address, token_id, blockchain, with_data)` | `network: str`, `contract_address: str`, `token_id: str`, `blockchain: Optional[str]`, `with_data: bool` | `Optional[TokenMetadata]` | Get ERC721/ERC1155 metadata |
| `get_fa(network, contract_address, token_id, with_data)` | `network: str`, `contract_address: str`, `token_id: str = "0"`, `with_data: bool = False` | `Optional[FATokenMetadata]` | Get Tezos FA token metadata |
| `get_crypto_punk(network, contract_address, punk_id, blockchain)` | `network: str`, `contract_address: str`, `punk_id: str`, `blockchain: Optional[str]` | `Optional[CryptoPunkMetadata]` | Get CryptoPunk metadata |

#### Example

```python
# Get ERC token metadata
metadata = client.token_metadata.get("ETH", "0x1234...", network="mainnet")
print(f"Token: {metadata.name}")

# Get NFT metadata with token ID
metadata = client.token_metadata.get_erc(
    network="mainnet", contract_address="0x1234...", token_id="42"
)
```

---

### UserDeviceService

Provides user device pairing management for multi-factor authentication.

**Access:** `client.user_devices`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `create_pairing()` | None | `UserDevicePairing` | Create pairing request (Step 1) |
| `start_pairing(pairing_id, nonce, encryption_key)` | `pairing_id: str`, `nonce: str`, `encryption_key: str` | `None` | Start pairing (Step 2) |
| `approve_pairing(pairing_id, nonce)` | `pairing_id: str`, `nonce: str` | `None` | Approve pairing (Step 3) |
| `get_pairing_status(pairing_id, nonce)` | `pairing_id: str`, `nonce: str` | `UserDevicePairingInfo` | Get pairing status |

#### Example

```python
# Create a new device pairing
pairing = client.user_devices.create_pairing()
print(f"Pairing ID: {pairing.pairing_id}")

# Start the pairing process (device provides nonce and key)
client.user_devices.start_pairing(
    pairing_id=pairing.pairing_id,
    nonce="123456",
    encryption_key="...",
)

# Approve the pairing
client.user_devices.approve_pairing(
    pairing_id=pairing.pairing_id,
    nonce="123456",
)

# Check pairing status
info = client.user_devices.get_pairing_status(pairing.pairing_id, "123456")
print(f"Status: {info.status}")
```

---

## TaurusNetwork Services

TaurusNetwork services are accessed through the `client.taurus_network` namespace.

### ParticipantService

Provides Taurus Network participant management.

**Access:** `client.taurus_network.participants`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_my_participant()` | None | `MyParticipant` | Get current participant info |
| `get_participant(participant_id)` | `participant_id: str` | `Participant` | Get participant by ID |
| `list_participants(...)` | Filters | `Tuple[List[Participant], Optional[Pagination]]` | List participants |
| `update_settings(settings)` | `settings: ParticipantSettings` | `Participant` | Update settings |

#### Example

```python
# Get my participant info
me = client.taurus_network.participants.get_my_participant()
print(f"Participant: {me.name} (ID: {me.id})")

# List all participants
participants, _ = client.taurus_network.participants.list_participants()
for p in participants:
    print(f"{p.name}: {p.status}")
```

---

### PledgeService

Provides pledge lifecycle management with ECDSA approval signing.

**Access:** `client.taurus_network.pledges`

**Security:** Pledge action approval uses ECDSA signatures.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_pledge(pledge_id)` | `pledge_id: str` | `Pledge` | Get pledge |
| `list_pledges(opts)` | `opts: ListPledgesOptions` | `Tuple[List[Pledge], CursorPage]` | List pledges |
| `create_pledge(req)` | `req: CreatePledgeRequest` | `Tuple[Pledge, PledgeAction]` | Create pledge |
| `update_pledge(pledge_id, req)` | `pledge_id: str`, `req: UpdatePledgeRequest` | `Pledge` | Update pledge |
| `add_pledge_collateral(pledge_id, req)` | `pledge_id: str`, `req: AddPledgeCollateralRequest` | `Tuple[Pledge, PledgeAction]` | Add collateral |
| `withdraw_pledge(pledge_id, req)` | `pledge_id: str`, `req: WithdrawPledgeRequest` | `Tuple[PledgeWithdrawal, PledgeAction]` | Withdraw (pledgee) |
| `initiate_withdraw_pledge(pledge_id, req)` | `pledge_id: str`, `req: InitiateWithdrawPledgeRequest` | `Tuple[PledgeWithdrawal, PledgeAction]` | Initiate withdrawal (pledgor) |
| `unpledge(pledge_id)` | `pledge_id: str` | `Tuple[Pledge, PledgeAction]` | Unpledge all funds |
| `reject_pledge(pledge_id, req)` | `pledge_id: str`, `req: RejectPledgeRequest` | `Pledge` | Reject pledge |
| `list_pledge_actions(opts)` | `opts: ListPledgeActionsOptions` | `Tuple[List[PledgeAction], CursorPage]` | List actions (hash verified) |
| `list_pledge_actions_for_approval(opts)` | `opts: ListPledgeActionsOptions` | `Tuple[List[PledgeAction], CursorPage]` | Get pending actions (hash verified) |
| `approve_pledge_actions(actions, private_key, comment)` | `actions: List[PledgeAction]`, `private_key` | `int` | Approve with signature |
| `reject_pledge_actions(req)` | `req: RejectPledgeActionsRequest` | `int` | Reject actions |
| `list_pledge_withdrawals(opts)` | `opts: ListPledgeWithdrawalsOptions` | `Tuple[List[PledgeWithdrawal], CursorPage]` | List withdrawals |

#### Example

```python
from taurus_protect.models.taurus_network import CreatePledgeRequest

# Create a pledge
request = CreatePledgeRequest(
    shared_address_id="addr-123",
    currency_id="ETH",
    amount="1000000000000000000",
    pledge_type="PLEDGEE_WITHDRAWALS_RIGHTS",
)
pledge, action = client.taurus_network.pledges.create_pledge(request)
print(f"Created pledge {pledge.id}, action {action.id} pending approval")

# Approve pledge actions
actions, _ = client.taurus_network.pledges.list_pledge_actions_for_approval()
if actions:
    count = client.taurus_network.pledges.approve_pledge_actions(
        actions, private_key
    )
    print(f"Approved {count} action(s)")
```

---

### LendingService

Provides lending offer and agreement management.

**Access:** `client.taurus_network.lending`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_lending_offer(offer_id)` | `offer_id: str` | `LendingOffer` | Get offer |
| `list_lending_offers(opts)` | `opts: ListLendingOffersOptions` | `Tuple[List[LendingOffer], CursorPage]` | List offers |
| `create_lending_offer(req)` | `req: CreateLendingOfferRequest` | `LendingOffer` | Create offer |
| `cancel_lending_offer(offer_id)` | `offer_id: str` | `LendingOffer` | Cancel offer |
| `get_lending_agreement(agreement_id)` | `agreement_id: str` | `LendingAgreement` | Get agreement |
| `list_lending_agreements(opts)` | `opts: ListLendingAgreementsOptions` | `Tuple[List[LendingAgreement], CursorPage]` | List agreements |
| `list_lending_agreements_for_approval(opts)` | `opts: ListLendingAgreementsOptions` | `Tuple[List[LendingAgreement], CursorPage]` | List agreements pending approval (`ids` filter) |
| `accept_lending_offer(offer_id, req)` | `offer_id: str`, `req: AcceptLendingOfferRequest` | `LendingAgreement` | Accept offer |
| `repay_lending_agreement(agreement_id, req)` | `agreement_id: str`, `req: RepayLendingAgreementRequest` | `LendingAgreement` | Repay loan |

#### Example

```python
# List available lending offers
offers, _ = client.taurus_network.lending.list_lending_offers()
for offer in offers:
    print(f"Offer {offer.id}: {offer.amount} {offer.currency_id} at {offer.interest_rate}%")

# Accept an offer
agreement = client.taurus_network.lending.accept_lending_offer(
    offer_id="offer-123",
    req=AcceptLendingOfferRequest(
        collateral_shared_address_id="addr-456",
        collateral_amount="2000000000000000000",
    ),
)
```

---

### SettlementService

Provides settlement operations for Taurus Network.

**Access:** `client.taurus_network.settlements`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_settlement(settlement_id)` | `settlement_id: str` | `Settlement` | Get settlement |
| `list_settlements(opts)` | `opts: ListSettlementsOptions` | `Tuple[List[Settlement], CursorPage]` | List settlements |
| `list_settlements_for_approval(opts)` | `opts: ListSettlementsForApprovalOptions` | `Tuple[List[Settlement], CursorPage]` | List settlements pending approval |
| `create_settlement(req)` | `req: CreateSettlementRequest` | `Settlement` | Create settlement |
| `approve_settlement(settlement_id)` | `settlement_id: str` | `Settlement` | Approve settlement |
| `reject_settlement(settlement_id, comment)` | `settlement_id: str`, `comment: str` | `Settlement` | Reject settlement |

#### Example

```python
# List pending settlements
settlements, _ = client.taurus_network.settlements.list_settlements()
for s in settlements:
    print(f"Settlement {s.id}: {s.status}")
```

---

### SharingService

Provides address and asset sharing operations.

**Access:** `client.taurus_network.sharing`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list_shared_addresses(opts)` | `opts: ListSharedAddressesOptions` | `Tuple[List[SharedAddress], CursorPage]` | List shared addresses |
| `create_shared_address(req)` | `req: CreateSharedAddressRequest` | `SharedAddress` | Share an address |
| `revoke_shared_address(address_id)` | `address_id: str` | `None` | Revoke sharing |
| `list_shared_assets(opts)` | `opts: ListSharedAssetsOptions` | `Tuple[List[SharedAsset], CursorPage]` | List shared assets |
| `create_shared_asset(req)` | `req: CreateSharedAssetRequest` | `SharedAsset` | Share an asset |
| `revoke_shared_asset(asset_id)` | `asset_id: str` | `None` | Revoke sharing |

#### Example

```python
# List shared addresses
addresses, _ = client.taurus_network.sharing.list_shared_addresses()
for addr in addresses:
    print(f"Shared: {addr.address} with {addr.target_participant_name}")

# Share an address
shared = client.taurus_network.sharing.create_shared_address(
    CreateSharedAddressRequest(
        internal_address_id="123",
        target_participant_id="participant-456",
        permissions=["VIEW", "RECEIVE"],
    )
)
```

---

## Exception Handling

All services raise consistent exceptions:

| Exception | HTTP Code | Description |
|-----------|-----------|-------------|
| `ValidationError` | 400 | Input validation failed |
| `AuthenticationError` | 401 | Invalid credentials |
| `AuthorizationError` | 403 | Insufficient permissions |
| `NotFoundError` | 404 | Resource not found |
| `RateLimitError` | 429 | Rate limit exceeded |
| `ServerError` | 5xx | Server error |
| `IntegrityError` | - | Cryptographic verification failed |
| `WhitelistError` | - | Whitelist signature verification failed |

### Example

```python
from taurus_protect.errors import (
    APIError,
    NotFoundError,
    RateLimitError,
    IntegrityError,
)

try:
    wallet = client.wallets.get(999999)
except NotFoundError:
    print("Wallet not found")
except RateLimitError as e:
    delay = e.suggested_retry_delay()
    print(f"Rate limited, retry after {delay.total_seconds()}s")
except IntegrityError as e:
    # Security error - do not retry
    print(f"Verification failed: {e.message}")
except APIError as e:
    if e.is_retryable():
        print(f"Retryable: {e.message}")
```

---

## Pagination Patterns

Page sizes default to 20 and may not exceed 100 (`DEFAULT_PAGE_SIZE` / `MAX_PAGE_SIZE`);
above 100 or negative raises `ValueError` before any request. See
[Key Concepts](CONCEPTS.md#pagination-model) for the contract.

### Offset lists

Continue with `offset=pagination.next_offset` while `pagination.has_more`; never compute the
next offset yourself.

```python
all_wallets = []
offset = 0

while True:
    wallets, pagination = client.wallets.list(limit=100, offset=offset)
    all_wallets.extend(wallets)
    if not pagination.has_more:
        break
    offset = pagination.next_offset
```

### Cursor lists

Continue with `cursor=page.next_cursor` while `page.has_more`.

```python
cursor = None
while True:
    balances, page = client.balances.list(page_size=100, cursor=cursor)
    for balance in balances:
        print(balance)
    if not page.has_more:
        break
    cursor = page.next_cursor
```

### Options-Based Pagination

For advanced filtering, the page window rides on the options:

```python
from taurus_protect.models import ListWalletsOptions

options = ListWalletsOptions(currency="ETH", exclude_disabled=True, limit=50)
wallets, pagination = client.wallets.list_with_options(options)
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and design
- [Authentication](AUTHENTICATION.md) - Security details
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
- [Key Concepts](CONCEPTS.md) - Domain model

<!-- BEGIN GENERATED METHOD INDEX -->

## Complete Method Index

Generated from the python source by `scripts/api-surface/docs.py`; regenerate with
`./build.sh docs`. Every method below exists in the SDK, and `./build.sh docs --check`
fails if this list drifts or if the prose above documents a method that does not.

44 services, 210 public methods.

### ActionService

- `get(action_id: 'str') -> 'Action'` — Get an action by ID.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[Action], Pagination]'` — List actions, one page at a time.

### AddressService

- `create(request: 'CreateAddressRequest') -> 'Address'` — Create a new address.
- `create_address(wallet_id: 'int', label: 'str', comment: 'str' = '', customer_id: 'str' = '') -> 'Address'` — Create a new address with explicit parameters.
- `create_attribute(address_id: 'int', key: 'str', value: 'str') -> 'None'` — Create an attribute for an address.
- `delete_attribute(address_id: 'int', attribute_id: 'int') -> 'None'` — Delete an attribute from an address.
- `get(address_id: 'int') -> 'Address'` — Get an address by ID with mandatory signature verification.
- `get_proof_of_reserve(address_id: 'int', challenge: 'Optional[str]' = None) -> 'Any'` — Get the proof of reserve for an address.
- `list(wallet_id: 'int', limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, exclude_disabled: 'Optional[bool]' = None) -> 'Tuple[List[Address], Pagination]'` — List a wallet's addresses, one page at a time, with mandatory signature verification.
- `list_with_options(options: 'Optional[ListAddressesOptions]' = None) -> 'Tuple[List[Address], Pagination]'` — List addresses with filters, one page at a time, with mandatory signature verification.

### AirGapService

- `get_unsigned_payload(request_id: 'int') -> 'bytes'` — Get unsigned transaction payload for offline signing.
- `submit_signed_payload(request_id: 'int', signed_payload: 'bytes') -> 'None'` — Submit a signed payload back to the system.

### AssetService

- `get(asset_id: 'str') -> 'Asset'` — Get an asset by its currency ID or symbol.
- `get_addresses(currency: 'str', page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, wallet_id: 'Optional[str]' = None, address_id: 'Optional[str]' = None, addresses: 'Optional[List[str]]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Address], CursorPage]'` — List the addresses holding an asset, one page at a time, each signature verified.
- `get_wallets(currency: 'str', page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, wallet_id: 'Optional[str]' = None, wallet_name: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Wallet], CursorPage]'` — List the wallets holding an asset, one page at a time.
- `list(currency: 'str', page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None) -> 'Tuple[List[Wallet], CursorPage]'` — List the wallets holding an asset, one page at a time; same as :meth:`get_wallets`.
- `list_asset_operations(asset_id: 'str', *, type: 'Optional[str]' = None, status: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[AssetOperationV2], CursorPage]'` — List the operations of a v2 asset, one page at a time.
- `query_asset_addresses(asset_id: 'str', *, address_type: 'Optional[str]' = None, kyc_status: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'QueryAssetAddressesResult'` — List the addresses holding a v2 asset, one page at a time.
- `query_assets(*, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, symbol: 'Optional[str]' = None, contract_address: 'Optional[str]' = None, label: 'Optional[str]' = None, currency_name: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[AssetV2], CursorPage]'` — List asset definitions (v2 API), one page at a time.

### AuditService

- `export_audit_trails(external_user_id: 'Optional[str]' = None, entities: 'Optional[List[str]]' = None, actions: 'Optional[List[str]]' = None, from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, format: 'Optional[str]' = None) -> 'str'` — Export audit trails in the specified format.
- `get(audit_id: 'str') -> 'Audit'` — Get an audit event by ID.
- `list(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, external_user_id: 'Optional[str]' = None, entities: 'Optional[List[str]]' = None, actions: 'Optional[List[str]]' = None, from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, sort_by: 'Optional[List[str]]' = None, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Audit], CursorPage]'` — List audit trails, one page at a time.

### BalanceService

- `list(currency: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, token_id: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[AssetBalance], CursorPage]'` — List the tenant's balances, one page at a time, optionally by currency.
- `list_nft_collections(blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, query: 'Optional[str]' = None, only_positive_balance: 'Optional[bool]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[NFTCollectionBalance], CursorPage]'` — List the tenant's NFT collection balances, one page at a time.

### BlockchainService

- `get(blockchain: 'str', network: 'str' = 'mainnet', include_block_height: 'bool' = False) -> 'Blockchain'` — Get blockchain information.
- `get_by_id(blockchain_id: 'str') -> 'Blockchain'` — Get blockchain by composite ID.
- `list(blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, include_block_height: 'bool' = False) -> 'List[Blockchain]'` — List supported blockchains.

### BusinessRuleService

- `list(page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None, rule_keys: 'Optional[List[str]]' = None, wallet_ids: 'Optional[List[str]]' = None, currency_ids: 'Optional[List[str]]' = None, entity_type: 'Optional[str]' = None, entity_ids: 'Optional[List[str]]' = None, *, cursor: 'Optional[str]' = None, ids: 'Optional[List[str]]' = None, rule_groups: 'Optional[List[str]]' = None, address_ids: 'Optional[List[str]]' = None, level: 'Optional[str]' = None) -> 'BusinessRuleResult'` — List business rules, one page at a time (v2 API).
- `list_by_currency(currency_id: 'str', page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None, *, cursor: 'Optional[str]' = None) -> 'BusinessRuleResult'` — List business rules for a specific currency, one page at a time.
- `list_by_wallet(wallet_id: 'int', page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None, *, cursor: 'Optional[str]' = None) -> 'BusinessRuleResult'` — List business rules for a specific wallet, one page at a time.
- `update_transactions_enabled(enabled: 'bool') -> 'None'` — Enable or disable transaction processing for the tenant.

### ChangeService

- `approve_change(change_id: 'str') -> 'None'` — Approve a change.
- `approve_changes(change_ids: 'List[str]') -> 'None'` — Approve multiple changes.
- `create_change(request: 'CreateChangeRequest') -> 'str'` — Create a change request.
- `get(change_id: 'str') -> 'Change'` — Get a change by ID.
- `list(options: 'Optional[ListChangesOptions]' = None) -> 'ChangeResult'` — List changes, one page at a time.
- `list_for_approval(options: 'Optional[ListChangesOptions]' = None) -> 'ChangeResult'` — List changes pending approval, one page at a time.
- `reject_change(change_id: 'str') -> 'None'` — Reject a change.
- `reject_changes(change_ids: 'List[str]') -> 'None'` — Reject multiple changes.

### ConfigService

- `get() -> 'TenantConfig'` — Get tenant configuration.
- `get_features() -> 'List[Feature]'` — Get enabled features for the tenant.

### ContractWhitelistingService

- `approve_whitelisted_contracts(contract_ids: 'List[str]', signature: 'str', comment: 'Optional[str]' = None) -> 'None'` — Approve whitelisted contracts with a signature.
- `create(address: 'str', name: 'str', blockchain: 'str', network: 'Optional[str]' = None, abi: 'Optional[str]' = None) -> 'int'` — Create a whitelisted contract request.
- `create_attribute(contract_id: 'str', key: 'str', value: 'str') -> 'None'` — Create an attribute on a whitelisted contract.
- `get_attribute(contract_id: 'str', key: 'str') -> 'Optional[str]'` — Get an attribute value from a whitelisted contract.

### CurrencyService

- `get(currency_id: 'str') -> 'Currency'` — Get a currency by ID.
- `get_base_currency() -> 'Currency'` — Get the base currency configured for the tenant.
- `get_by_blockchain(blockchain: 'str', network: 'str', contract_address: 'Optional[str]' = None, token_id: 'Optional[str]' = None) -> 'Currency'` — Get a currency by blockchain and network.
- `list(show_disabled: 'bool' = False, include_logo: 'bool' = False) -> 'List[Currency]'` — Get all currencies.

### EarnService

- `list_rewards(recipient_address_id: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[EarnReward], CursorPage]'` — List earn rewards, one page at a time.

### ExchangeService

- `get(exchange_id: 'str') -> 'Exchange'` — Get an exchange account by ID.
- `get_withdrawal_fee(exchange_id: 'str', to_address_id: 'Optional[str]' = None, amount: 'Optional[str]' = None) -> 'Any'` — Get withdrawal fees for an exchange account.
- `list(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, currency_id: 'Optional[str]' = None, exchange_label: 'Optional[str]' = None, status: 'Optional[str]' = None, only_positive_balance: 'Optional[bool]' = None, *, sort_order: 'Optional[str]' = None, include_base_currency_valuation: 'Optional[bool]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Exchange], CursorPage]'` — List exchange accounts, one page at a time.
- `list_counterparties() -> 'List[Any]'` — List exchange counterparties with their exposure limits.

### FeePayerService

- `get(fee_payer_id: 'str') -> 'FeePayer'` — Get a fee payer by ID.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'Tuple[List[FeePayer], Pagination]'` — List fee payers, one page at a time.

### FeeService

- `estimate(currency: 'str', amount: 'Optional[str]' = None, destination: 'Optional[str]' = None) -> 'FeeEstimate'` — Estimate transaction fee.
- `list() -> 'List[FeeEstimate]'` — List fee estimates for all supported currencies.

### FiatService

- `get_account(account_id: 'str') -> 'FiatProviderAccount'` — Get a fiat provider account by ID.
- `get_base_currency() -> 'FiatCurrency'` — Get the configured base currency.
- `get_rate(from_currency: 'str', to_currency: 'str') -> 'ExchangeRate'` — Get the exchange rate between two currencies.
- `list_fiat_provider_accounts(provider: 'str', label: 'str', page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, account_type: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[FiatProviderAccount], CursorPage]'` — List a fiat provider's accounts, one page at a time.
- `list_fiat_provider_entities(provider: 'Optional[str]' = None, label: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[FiatProviderEntity], CursorPage]'` — List the entities registered with fiat providers, one page at a time.
- `list_providers() -> 'List[Any]'` — List available fiat providers.

### GovernanceRuleService

- `approve_rules_proposal(private_key: 'EllipticCurvePrivateKey', comment: 'str', expected_container_hash: 'str') -> 'None'` — Sign the pending proposal's rules container and submit the approval.
- `decode_proposal_for_review(rules: 'GovernanceRules') -> 'DecodedRulesContainer'` — Decode a PENDING proposal so a SuperAdmin can inspect it before approving.
- `get_decoded_rules_container(rules: 'GovernanceRules') -> 'DecodedRulesContainer'` — Get the decoded rules container from governance rules.
- `get_public_keys() -> 'List[SuperAdminPublicKey]'` — Get the list of SuperAdmin public keys.
- `get_rules() -> 'Optional[GovernanceRules]'` — Get the currently enforced governance rules.
- `get_rules_by_id(rules_id: 'str') -> 'Optional[GovernanceRules]'` — Get a governance ruleset by its ID.
- `get_rules_history(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None) -> 'GovernanceRulesHistoryResult'` — Get the history of governance rules, one page at a time.
- `get_rules_proposal() -> 'Optional[GovernanceRules]'` — Get the proposed governance rules.
- `proposal_container_hash(rules: 'GovernanceRules') -> 'str'` — Return the canonical SHA-256 hex digest of a ruleset's decoded container.
- `reject_rules_proposal(comment: 'str') -> 'None'` — Reject the pending rules proposal with a comment (SuperAdmin only).
- `update_rules_proposal(container: 'DecodedRulesContainer') -> 'None'` — Submit a rules container as a governance proposal (SuperAdmin only).
- `verify_governance_rules(rules: 'GovernanceRules') -> 'GovernanceRules'` — Verify that governance rules have enough valid SuperAdmin signatures.

### GroupService

- `get(group_id: 'str') -> 'Group'` — Get a group by ID.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[Group], Pagination]'` — List groups, one page at a time.

### HealthService

- `check() -> 'HealthStatus'` — Check the API health status.
- `get_all_health_checks(tenant_id: 'Optional[str]' = None, fail_if_unhealthy: 'bool' = False) -> 'GetAllHealthChecksResult'` — Get all health checks with optional filtering.

### JobService

- `get(job_id: 'str') -> 'Job'` — Get a job by ID.
- `list() -> 'List[Job]'` — List jobs: every job the endpoint returns, which does not page.

### LendingService

- `cancel_lending_agreement(lending_agreement_id: 'str') -> 'None'` — Cancel a lending agreement.
- `create_attachment(agreement_id: 'str', request: 'CreateLendingAgreementAttachmentRequest') -> 'str'` — Add an attachment to a lending agreement.
- `create_lending_agreement(request: 'CreateLendingAgreementRequest') -> 'str'` — Create a new lending agreement.
- `create_lending_agreement_attachment(lending_agreement_id: 'str', request: 'CreateLendingAgreementAttachmentRequest') -> 'str'` — Add an attachment to a lending agreement.
- `create_lending_offer(request: 'CreateLendingOfferRequest') -> 'str'` — Create a new lending offer.
- `delete_lending_offer(offer_id: 'str') -> 'None'` — Delete a specific lending offer.
- `delete_lending_offers() -> 'None'` — Delete all lending offers for the current participant.
- `get_lending_agreement(lending_agreement_id: 'str') -> 'LendingAgreement'` — Get a lending agreement by ID.
- `get_lending_offer(offer_id: 'str') -> 'LendingOffer'` — Get a lending offer by ID.
- `list_attachments(agreement_id: 'str') -> 'List[LendingAgreementAttachment]'` — List attachments for a lending agreement.
- `list_lending_agreement_attachments(lending_agreement_id: 'str') -> 'List[LendingAgreementAttachment]'` — List attachments for a lending agreement.
- `list_lending_agreements(options: 'Optional[ListLendingAgreementsOptions]' = None) -> 'Tuple[List[LendingAgreement], CursorPage]'` — List lending agreements, one page at a time.
- `list_lending_agreements_for_approval(options: 'Optional[ListLendingAgreementsOptions]' = None) -> 'Tuple[List[LendingAgreement], CursorPage]'` — List lending agreements pending approval, one page at a time.
- `list_lending_offers(options: 'Optional[ListLendingOffersOptions]' = None) -> 'Tuple[List[LendingOffer], CursorPage]'` — List lending offers, one page at a time.
- `repay_lending_agreement(lending_agreement_id: 'str', request: 'RepayLendingAgreementRequest') -> 'None'` — Record repayment for a lending agreement.
- `update_lending_agreement(lending_agreement_id: 'str', request: 'UpdateLendingAgreementRequest') -> 'None'` — Update a lending agreement.

### MultiFactorSignatureService

- `approve_multi_factor_signature(id: 'str', signature: 'str', comment: 'str' = '') -> 'MultiFactorSignatureApprovalResult'` — Approve a multi-factor signature request.
- `create_multi_factor_signatures(entity_ids: 'List[str]', entity_type: 'Union[MultiFactorSignatureEntityType, str]') -> 'MultiFactorSignatureResult'` — Create a batch of multi-factor signature requests.
- `get_multi_factor_signature_info(id: 'str') -> 'MultiFactorSignatureInfo'` — Retrieve information about a multi-factor signature request.
- `reject_multi_factor_signature(id: 'str', comment: 'str' = '') -> 'None'` — Reject a multi-factor signature request.

### ParticipantService

- `create_participant_attribute(participant_id: 'str', request: 'CreateParticipantAttributeRequest') -> 'None'` — Create an attribute for a participant.
- `delete_participant_attribute(participant_id: 'str', attribute_id: 'str') -> 'None'` — Delete an attribute for a participant.
- `get(participant_id: 'str', options: 'Optional[GetParticipantOptions]' = None) -> 'Participant'` — Get a participant by ID.
- `get_my_participant() -> 'MyParticipant'` — Get the current participant with settings.
- `list(options: 'Optional[ListParticipantsOptions]' = None) -> 'List[Participant]'` — List visible Taurus Network participants.

### PledgeService

- `add_pledge_collateral(pledge_id: 'str', req: 'AddPledgeCollateralRequest') -> 'Tuple[Pledge, PledgeAction]'` — Add collateral to an existing pledge.
- `approve_pledge_actions(actions: 'List[PledgeAction]', private_key: 'EllipticCurvePrivateKey', comment: 'str' = 'approving via taurus-protect-sdk-python') -> 'int'` — Approve multiple pledge actions with ECDSA signature.
- `create_pledge(req: 'CreatePledgeRequest') -> 'Tuple[Pledge, PledgeAction]'` — Create a new pledge.
- `get_pledge(pledge_id: 'str') -> 'Pledge'` — Get a pledge by ID.
- `initiate_withdraw_pledge(pledge_id: 'str', req: 'InitiateWithdrawPledgeRequest') -> 'Tuple[PledgeWithdrawal, PledgeAction]'` — Initiate withdrawal from a pledge (pledgor operation).
- `list_pledge_actions(opts: 'Optional[ListPledgeActionsOptions]' = None) -> 'Tuple[List[PledgeAction], CursorPage]'` — List pledge actions, one page at a time, each hash verified against its payload.
- `list_pledge_actions_for_approval(opts: 'Optional[ListPledgeActionsOptions]' = None) -> 'Tuple[List[PledgeAction], CursorPage]'` — List pledge actions pending approval, one page at a time.
- `list_pledge_withdrawals(opts: 'Optional[ListPledgeWithdrawalsOptions]' = None) -> 'Tuple[List[PledgeWithdrawal], CursorPage]'` — List pledge withdrawals, one page at a time.
- `list_pledges(opts: 'Optional[ListPledgesOptions]' = None) -> 'Tuple[List[Pledge], CursorPage]'` — List pledges, one page at a time.
- `reject_pledge(pledge_id: 'str', req: 'RejectPledgeRequest') -> 'Pledge'` — Reject a pledge.
- `reject_pledge_actions(req: 'RejectPledgeActionsRequest') -> 'int'` — Reject multiple pledge actions.
- `unpledge(pledge_id: 'str') -> 'Tuple[Pledge, PledgeAction]'` — Unpledge all funds from a pledge.
- `update_pledge(pledge_id: 'str', req: 'UpdatePledgeRequest') -> 'Pledge'` — Update a pledge's default destination.
- `withdraw_pledge(pledge_id: 'str', req: 'WithdrawPledgeRequest') -> 'Tuple[PledgeWithdrawal, PledgeAction]'` — Withdraw from a pledge (pledgee operation).

### PriceService

- `get_current(*, from_currency_id: 'Optional[str]' = None, to_currency_ids: 'Optional[List[str]]' = None, only_primary: 'Optional[bool]' = None, sort_order: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Price], CursorPage]'` — List current prices, one page at a time, each signature verified.
- `get_historical(base_currency: 'str', quote_currency: 'str', limit: 'Optional[int]' = None) -> 'List[PriceHistoryPoint]'` — Get historical prices for a currency pair.

### RequestService

- `approve_request(request: 'Request', private_key: 'EllipticCurvePrivateKey', comment: 'str' = 'approving via taurus-protect-sdk-python') -> 'int'` — Approve a single request with ECDSA signature.
- `approve_requests(requests: 'List[Request]', private_key: 'EllipticCurvePrivateKey', comment: 'str' = 'approving via taurus-protect-sdk-python') -> 'int'` — Approve multiple requests with ECDSA signature.
- `create_cancel_request(address_id: 'int', nonce: 'int') -> 'Request'` — Create a cancel request for a pending transaction.
- `create_external_transfer(from_address_id: 'int', to_whitelisted_address_id: 'int', amount: 'str') -> 'Request'` — Create an external transfer to a whitelisted address.
- `create_external_transfer_from_wallet(from_wallet_id: 'int', to_whitelisted_address_id: 'int', amount: 'str') -> 'Request'` — Create an external transfer from an omnibus wallet.
- `create_incoming_request(from_exchange_id: 'int', to_address_id: 'int', amount: 'str') -> 'Request'` — Create an incoming request from an exchange.
- `create_internal_transfer(from_address_id: 'int', to_address_id: 'int', amount: 'str') -> 'Request'` — Create an internal transfer request between addresses.
- `create_internal_transfer_from_wallet(from_wallet_id: 'int', to_address_id: 'int', amount: 'str') -> 'Request'` — Create an internal transfer from an omnibus wallet.
- `get(request_id: 'int') -> 'Request'` — Get a request by ID with mandatory hash verification.
- `get_for_approval(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, currency_id: 'Optional[str]' = None, types: 'Optional[List[str]]' = None, exclude_types: 'Optional[List[str]]' = None, ids: 'Optional[List[str]]' = None, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Request], CursorPage]'` — List requests pending the caller's approval, one page at a time.
- `list(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency_id: 'Optional[str]' = None, statuses: 'Optional[List[RequestStatus]]' = None, *, types: 'Optional[List[str]]' = None, ids: 'Optional[List[str]]' = None, external_request_ids: 'Optional[List[str]]' = None, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Request], CursorPage]'` — List requests with filtering, one page at a time.
- `reject_request(request_id: 'int', comment: 'str') -> 'None'` — Reject a single request.
- `reject_requests(request_ids: 'List[int]', comment: 'str') -> 'None'` — Reject multiple requests.

### ReservationService

- `get(reservation_id: 'int') -> 'Reservation'` — Get a reservation by ID.
- `list(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, kind: 'Optional[str]' = None, kinds: 'Optional[List[str]]' = None, address: 'Optional[str]' = None, address_id: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Reservation], CursorPage]'` — List reservations, one page at a time.

### ScoreService

- `get_address_score(address_id: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Get risk score for an address.
- `get_transaction_score(tx_hash: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Get risk score for a transaction.
- `refresh_whitelisted_address_score(address_id: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Refresh risk score for a whitelisted address.

### SettlementService

- `cancel_settlement(settlement_id: 'str') -> 'None'` — Cancel a settlement.
- `create_settlement(request: 'CreateSettlementRequest') -> 'str'` — Create a new settlement.
- `get_settlement(settlement_id: 'str') -> 'Settlement'` — Get a settlement by ID.
- `list_settlements(options: 'Optional[ListSettlementsOptions]' = None) -> 'Tuple[List[Settlement], CursorPage]'` — List settlements, one page at a time.
- `list_settlements_for_approval(options: 'Optional[ListSettlementsForApprovalOptions]' = None) -> 'Tuple[List[Settlement], CursorPage]'` — List settlements pending approval, one page at a time.
- `replace_settlement(settlement_id: 'str', request: 'CreateSettlementRequest') -> 'None'` — Replace a settlement with new attributes.

### SharingService

- `list_shared_addresses(options: 'Optional[ListSharedAddressesOptions]' = None) -> 'Tuple[List[SharedAddress], CursorPage]'` — List shared addresses, one page at a time.
- `list_shared_assets(options: 'Optional[ListSharedAssetsOptions]' = None) -> 'Tuple[List[SharedAsset], CursorPage]'` — List shared whitelisted assets, one page at a time.
- `share_address(request: 'ShareAddressRequest') -> 'None'` — Share an address with a Taurus Network participant.
- `share_whitelisted_asset(request: 'ShareWhitelistedAssetRequest') -> 'None'` — Share a whitelisted asset with a Taurus Network participant.
- `unshare_address(shared_address_id: 'str') -> 'None'` — Unshare an address with a Taurus Network participant.
- `unshare_whitelisted_asset(shared_asset_id: 'str') -> 'None'` — Unshare a whitelisted asset with a Taurus Network participant.

### StakingService

- `get_staking_info(address_id: 'int') -> 'StakingInfo'` — Get staking information for an address.
- `list_validators(blockchain: 'str', network: 'str' = 'mainnet', *, ids: 'Optional[List[str]]' = None) -> 'List[Validator]'` — List validators for a blockchain: every validator the endpoint returns.

### StatisticsService

- `get_summary() -> 'Optional[PortfolioStatistics]'` — Get summary statistics for the portfolio.
- `get_transaction_stats(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None) -> 'TransactionStatistics'` — Get transaction statistics for a date range.

### TagService

- `create(name: 'str', color: 'str') -> 'Tag'` — Create a new tag.
- `delete(tag_id: 'str') -> 'None'` — Delete a tag.
- `get(tag_id: 'str') -> 'Tag'` — Get a tag by ID.
- `list(*, query: 'Optional[str]' = None, ids: 'Optional[List[str]]' = None) -> 'List[Tag]'` — List tags: every tag the endpoint returns, which does not page.

### TokenMetadataService

- `get(blockchain: 'str', contract_address: 'str', token_id: 'str' = '0', network: 'str' = 'mainnet', with_data: 'bool' = False) -> 'Optional[TokenMetadata]'` — Get token metadata for an ERC token.
- `get_crypto_punk(network: 'str', contract_address: 'str', punk_id: 'str', blockchain: 'Optional[str]' = None) -> 'Optional[CryptoPunkMetadata]'` — Get CryptoPunk token metadata.
- `get_erc(network: 'str', contract_address: 'str', token_id: 'str', blockchain: 'Optional[str]' = None, with_data: 'bool' = False) -> 'Optional[TokenMetadata]'` — Get ERC token metadata (ERC721 or ERC1155).
- `get_fa(network: 'str', contract_address: 'str', token_id: 'str' = '0', with_data: 'bool' = False) -> 'Optional[FATokenMetadata]'` — Get FA token metadata (Tezos FA1.2 or FA2).

### TransactionService

- `export(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency: 'Optional[str]' = None, direction: 'Optional[str]' = None, limit: 'Optional[int]' = None, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, format: 'Optional[str]' = None) -> 'TransactionExport'` — Export transactions in one reply.
- `export_csv(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency: 'Optional[str]' = None, direction: 'Optional[str]' = None, limit: 'Optional[int]' = None, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'TransactionExport'` — Export transactions as CSV: :meth:`export` with ``format="csv"``.
- `get(transaction_id: 'int') -> 'Transaction'` — Get a single transaction by ID.
- `get_by_hash(tx_hash: 'str') -> 'Transaction'` — Get a transaction by its blockchain hash.
- `list(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency: 'Optional[str]' = None, direction: 'Optional[str]' = None, limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'Tuple[List[Transaction], Pagination]'` — List transactions with filtering, one page at a time.
- `list_by_address(address: 'str', limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[Transaction], Pagination]'` — List transactions for a specific blockchain address, one page at a time.

### UserDeviceService

- `approve_pairing(pairing_id: 'str', nonce: 'str') -> 'None'` — Approve a user device pairing request (Step 3).
- `create_pairing() -> 'UserDevicePairing'` — Create a new user device pairing request (Step 1).
- `get(device_id: 'str') -> 'UserDevicePairingInfo'` — Get user device pairing by ID.
- `get_pairing_status(pairing_id: 'str', nonce: 'str') -> 'UserDevicePairingInfo'` — Get the status of a user device pairing request.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[UserDevicePairing], Pagination]'` — List user device pairings.
- `start_pairing(pairing_id: 'str', nonce: 'str', encryption_key: 'str') -> 'None'` — Start a user device pairing request (Step 2).

### UserService

- `create_user_attribute(user_id: 'str', key: 'str', value: 'str') -> 'None'` — Create an attribute for a user.
- `get(user_id: 'str') -> 'User'` — Get a user by ID.
- `get_current() -> 'User'` — Get the current authenticated user.
- `get_users_by_email(emails: 'List[str]') -> 'List[User]'` — Get users by their email addresses, reading every page.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[User], Pagination]'` — List users, one page at a time.

### VisibilityGroupService

- `get(group_id: 'str') -> 'VisibilityGroup'` — Get a visibility group by ID.
- `get_users(group_id: 'str') -> 'List[Any]'` — Get users in a visibility group.
- `list() -> 'List[VisibilityGroup]'` — List visibility groups: every group the endpoint returns, which does not page.

### WalletService

- `create(request: 'CreateWalletRequest') -> 'Wallet'` — Create a new wallet.
- `create_attribute(wallet_id: 'int', key: 'str', value: 'str') -> 'None'` — Create an attribute for a wallet.
- `create_wallet(blockchain: 'str', network: 'str', name: 'str', is_omnibus: 'bool' = False, comment: 'str' = '', customer_id: 'str' = '') -> 'Wallet'` — Create a new wallet with explicit parameters.
- `get(wallet_id: 'int') -> 'Wallet'` — Get a wallet by ID.
- `get_balance_history(wallet_id: 'int', interval_hours: 'int') -> 'List[BalanceHistoryPoint]'` — Get wallet balance history.
- `get_by_name(name: 'str', limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, exclude_disabled: 'Optional[bool]' = None) -> 'Tuple[List[Wallet], Pagination]'` — List wallets whose name matches (case-insensitive, partial), one page at a time.
- `get_tokens(wallet_id: 'int', page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None) -> 'Tuple[List[AssetBalance], CursorPage]'` — List a wallet's token balances, one page at a time.
- `list(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, exclude_disabled: 'Optional[bool]' = None) -> 'Tuple[List[Wallet], Pagination]'` — List wallets, one page at a time.
- `list_with_options(options: 'Optional[ListWalletsOptions]' = None) -> 'Tuple[List[Wallet], Pagination]'` — List wallets with every filter the endpoint supports.

### WebhookCallService

- `get(call_id: 'str') -> 'WebhookCall'` — Get a webhook call by ID.
- `get_webhook_calls(event_id: 'Optional[str]' = None, webhook_id: 'Optional[str]' = None, status: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, cursor: 'Optional[Union[str, ApiRequestCursor]]' = None, *, page_size: 'Optional[int]' = None) -> 'WebhookCallResult'` — Retrieve webhook call history, one page at a time.
- `list(webhook_id: 'Optional[str]' = None, event_id: 'Optional[str]' = None, status: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None) -> 'Tuple[List[WebhookCall], CursorPage]'` — List webhook calls, one page at a time.

### WebhookService

- `create(url: 'str', events: 'List[str]') -> 'Webhook'` — Create a new webhook.
- `delete(webhook_id: 'str') -> 'None'` — Delete a webhook.
- `get(webhook_id: 'str') -> 'Webhook'` — Get a webhook by ID.
- `list(page_size: 'Optional[int]' = None, cursor: 'Optional[str]' = None, *, type: 'Optional[str]' = None, url: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'Tuple[List[Webhook], CursorPage]'` — List webhooks, one page at a time.

### WhitelistedAddressService

- `approve(selection: 'WhitelistedAddressApproval', private_key: 'Any', comment: 'str') -> 'None'` — Sign and submit an approval for the reviewed whitelisted addresses,
- `get(whitelisted_address_id: 'int') -> 'WhitelistedAddress'` — Get a whitelisted address by ID with verification.
- `get_envelope(whitelisted_address_id: 'int') -> 'SignedWhitelistedAddressEnvelope'` — Get the signed envelope for a whitelisted address.
- `list(currency: 'Optional[str]' = None, limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, *, ids: 'Optional[List[str]]' = None, include_for_approval: 'bool' = False) -> 'WhitelistedAddressListResult'` — List whitelisted addresses with cryptographic verification.
- `list_for_approval(limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, *, ids: 'Optional[List[str]]' = None, include_already_signed_by_user: 'bool' = False) -> 'WhitelistedAddressListResult'` — List whitelisted addresses awaiting approval, verified exactly as ``list`` is.

### WhitelistedAssetService

- `approve(selection: 'WhitelistedAssetApproval', private_key: 'Any', comment: 'str') -> 'None'` — Sign and submit an approval for the reviewed whitelisted assets, all-or-nothing.
- `get(asset_id: 'int') -> 'WhitelistedAsset'` — Get a whitelisted asset by ID.
- `list(blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, limit: 'Optional[int]' = None, offset: 'Optional[int]' = None, *, ids: 'Optional[List[str]]' = None, include_for_approval: 'bool' = False) -> 'Tuple[List[WhitelistedAsset], Pagination]'` — List whitelisted assets, one page at a time.
- `list_for_approval(ids: 'Optional[List[str]]' = None, limit: 'Optional[int]' = None, offset: 'Optional[int]' = None) -> 'Tuple[List[WhitelistedAsset], Pagination]'` — List whitelisted assets awaiting approval, verified as in list().

<!-- END GENERATED METHOD INDEX -->
