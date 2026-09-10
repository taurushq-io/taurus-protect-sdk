# Services Reference

This document provides complete API documentation for all 43 services in the Taurus-PROTECT Python SDK.

## Service Overview

The SDK provides services organized into two categories: core services (38) and TaurusNetwork services (5).

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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Wallet], Optional[Pagination]]` | List wallets |
| `list_with_options(options)` | `options: ListWalletsOptions` | `Tuple[List[Wallet], Optional[Pagination]]` | List with full filtering |
| `get_by_name(name, limit, offset)` | `name: str`, `limit: int`, `offset: int` | `Tuple[List[Wallet], Optional[Pagination]]` | Find wallets by name |
| `create(request)` | `request: CreateWalletRequest` | `Wallet` | Create wallet |
| `create_wallet(...)` | See below | `Wallet` | Create with explicit params |
| `create_attribute(wallet_id, key, value)` | `wallet_id: int`, `key: str`, `value: str` | `None` | Add attribute |
| `get_balance_history(wallet_id, interval_hours)` | `wallet_id: int`, `interval_hours: int` | `List[BalanceHistoryPoint]` | Get balance history |
| `get_tokens(wallet_id, limit)` | `wallet_id: int`, `limit: int` | `List[AssetBalance]` | Get token balances |

#### Example

```python
from taurus_protect.models import CreateWalletRequest

# List wallets
wallets, pagination = client.wallets.list(limit=50)
print(f"Total: {pagination.total_items if pagination else len(wallets)}")

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
| `list(wallet_id, limit, offset)` | `wallet_id: int`, `limit: int`, `offset: int` | `Tuple[List[Address], Optional[Pagination]]` | List addresses |
| `list_with_options(options)` | `options: ListAddressesOptions` | `Tuple[List[Address], Optional[Pagination]]` | List with filtering |
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
| `list(limit, offset, ...)` | Multiple filters | `Tuple[List[Request], Optional[Pagination]]` | List requests |
| `get_for_approval(limit, offset)` | `limit: int`, `offset: int` | `Tuple[List[Request], Optional[Pagination]]` | Get pending approvals |
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
requests, _ = client.requests.get_for_approval(limit=10)

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
| `list(...)` | Multiple filters | `Tuple[List[Transaction], Optional[Pagination]]` | List transactions |
| `export(...)` | Export filters | Export response | Export transactions |

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
| `list(...)` | Filters | `Tuple[List[Balance], Optional[Pagination]]` | List balances |
| `get_totals(...)` | Filters | `BalanceTotals` | Get balance totals |

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
| `get_rules_history(page_size, cursor)` | `page_size: int = 50`, `cursor: Optional[bytes] = None` | `Tuple[List[GovernanceRules], Optional[bytes]]` | Get governance rules history |
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
| `list(...)` | Filters | `Tuple[List[WhitelistedAddress], Optional[Pagination]]` | List addresses |

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
| `list(...)` | Filters | `Tuple[List[WhitelistedAsset], Optional[Pagination]]` | List assets |
| `list_for_approval(ids, limit, offset)` | `ids: Optional[List[str]]`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[WhitelistedAsset], Optional[Pagination]]` | List assets awaiting approval, verified as in `list` |
| `approve(ids, private_key, comment)` | `ids: List[int]`, `private_key`, `comment: str` | `None` | Sign an approval, all-or-nothing |

`list_for_approval` exists here because the for-approval read used to live only on the
unverified contract service, so the rows an approver inspects were never checked against
governance.

`approve` re-reads and verifies each asset and signs the hashes *those* rows carry. Any row
that is missing or fails verification aborts the whole call and nothing is signed: the API
takes one signature covering the whole batch, so a partial approval would mean the caller
believes they approved more than they did.

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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Audit], Optional[Pagination]]` | List audit events |
| `get(audit_id)` | `audit_id: str` | `Audit` | Get audit event by ID |

#### Example

```python
# List audit events
audits, pagination = client.audits.list(limit=50)
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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Change], Optional[Pagination]]` | List changes |
| `get(change_id)` | `change_id: str` | `Change` | Get change by ID |
| `approve_change(change_id)` | `change_id: str` | `None` | Approve a change |
| `approve_changes(change_ids)` | `change_ids: List[str]` | `None` | Approve multiple changes |
| `reject_change(change_id)` | `change_id: str` | `None` | Reject a change |
| `reject_changes(change_ids)` | `change_ids: List[str]` | `None` | Reject multiple changes |

#### Example

```python
# List pending changes
changes, pagination = client.changes.list(limit=50)
for change in changes:
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
caller acts on. `get_current` verifies each price against the `PRICEUPDATER` keys in the
SuperAdmin-verified rules container, which is why the service takes the cache as a mandatory
parameter. Whether prices must be signed is the **container's** call: no `PRICEUPDATER`
configured means this tenant does not sign prices and the price passes through; a
`PRICEUPDATER` configured plus a price carrying no signatures raises `IntegrityError`.

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_current(currency)` | `currency: Optional[str] = None` | `List[Price]` | Get current prices |
| `get_historical(base_currency, quote_currency, limit)` | `base_currency: str`, `quote_currency: str`, `limit: Optional[int]` | `List[PriceHistoryPoint]` | Get historical prices |

#### Example

```python
# Get all current prices
prices = client.prices.get_current()
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
| `list_validators(blockchain, network, limit, offset)` | `blockchain: str`, `network: str = "mainnet"`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Validator], Optional[Pagination]]` | List validators |
| `get_staking_info(address_id)` | `address_id: int` | `StakingInfo` | Get staking info for address |

#### Example

```python
# List ETH validators
validators, pagination = client.staking.list_validators(blockchain="ETH", limit=50)
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

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `create(address, name, blockchain, network, abi)` | `address: str`, `name: str`, `blockchain: str`, `network: Optional[str]`, `abi: Optional[str]` | `int` | Create whitelisted contract |
| `delete(contract_id)` | `contract_id: int` | `None` | Delete whitelisted contract |
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
| `get(rule_id)` | `rule_id: int` | `BusinessRule` | Get business rule by ID |
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[BusinessRule], Optional[Pagination]]` | List business rules |

#### Example

```python
# List business rules
rules, pagination = client.business_rules.list()
for rule in rules:
    print(f"{rule.name}: {'enabled' if rule.enabled else 'disabled'}")
```

---

### ReservationService

Provides balance reservation management for pending transactions.

**Access:** `client.reservations`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(reservation_id)` | `reservation_id: int` | `Reservation` | Get reservation by ID |
| `list(wallet_id, limit, offset)` | `wallet_id: Optional[int]`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Reservation], Optional[Pagination]]` | List reservations |
| `cancel(reservation_id)` | `reservation_id: int` | `None` | Cancel a reservation |

#### Example

```python
# List reservations for a wallet
reservations, pagination = client.reservations.list(wallet_id=123)
for r in reservations:
    print(f"{r.id}: {r.amount} {r.currency} ({r.status})")

# Cancel a reservation
client.reservations.cancel(reservation_id=456)
```

---

### MultiFactorSignatureService

Provides multi-factor signature operations for high-value transactions.

**Access:** `client.multi_factor_signature`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get_challenge(challenge_id)` | `challenge_id: str` | `MultiFactorSignatureChallenge` | Get challenge by ID |
| `list_challenges(request_id, limit, offset)` | `request_id: Optional[int]`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[MultiFactorSignatureChallenge], Optional[Pagination]]` | List challenges |
| `create_challenge(request_id, challenge_type)` | `request_id: int`, `challenge_type: str` | `str` | Create a challenge |
| `verify_challenge(challenge_id, response)` | `challenge_id: str`, `response: str` | `bool` | Verify a challenge response |

#### Example

```python
# Create a challenge for a request
challenge_id = client.multi_factor_signature.create_challenge(
    request_id=123, challenge_type="TOTP"
)

# Verify the challenge
is_valid = client.multi_factor_signature.verify_challenge(
    challenge_id=challenge_id, response="123456"
)
print(f"Verification: {'passed' if is_valid else 'failed'}")
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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[User], Optional[Pagination]]` | List users |
| `get_users_by_email(emails)` | `emails: List[str]` | `List[User]` | Get users by email addresses |
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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Group], Optional[Pagination]]` | List groups |

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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[VisibilityGroup], Optional[Pagination]]` | List visibility groups |
| `get_users(group_id)` | `group_id: str` | `List[Any]` | Get users in a visibility group |

#### Example

```python
# List visibility groups
groups, pagination = client.visibility_groups.list()
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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Webhook], Optional[Pagination]]` | List webhooks |
| `get(webhook_id)` | `webhook_id: str` | `Webhook` | Get webhook by ID |
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
webhooks, pagination = client.webhooks.list()
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
| `get_webhook_calls(event_id, webhook_id, status, sort_order, cursor)` | All optional filters | `WebhookCallResult` | Get webhook calls with filtering |
| `list(webhook_id, event_id, status, sort_order, limit, cursor)` | Optional filters, `limit: int = 50` | `Tuple[List[WebhookCall], Optional[Pagination]]` | List webhook calls |
| `get(call_id)` | `call_id: str` | `WebhookCall` | Get webhook call by ID |

#### Example

```python
# Get webhook call history
result = client.webhook_calls.get_webhook_calls(webhook_id="webhook-123")
for call in result.calls:
    print(f"{call.id}: {call.status}")

# List failed calls
calls, pagination = client.webhook_calls.list(status="FAILED", limit=50)
```

---

### TagService

Provides tag management for organizing entities.

**Access:** `client.tags`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(tag_id)` | `tag_id: str` | `Tag` | Get tag by ID |
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `List[Tag]` | List tags |
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

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Asset], Optional[Pagination]]` | List assets |
| `get(asset_id)` | `asset_id: str` | `Asset` | Get asset by ID |
| `get_wallets(currency, limit, offset)` | `currency: str`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Any], Optional[Pagination]]` | Get wallet balances for an asset |
| `get_addresses(currency, limit, offset)` | `currency: str`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Address], Optional[Pagination]]` | Get address balances for an asset, HSM-verified |

#### Example

```python
# List assets
assets, pagination = client.assets.list(limit=50)
for asset in assets:
    print(f"{asset.symbol}: {asset.name}")

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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Action], Optional[Pagination]]` | List actions |

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
| `list(limit, offset, currency_id, exchange_label, status, only_positive_balance)` | Multiple filters | `Tuple[List[Exchange], Optional[Pagination]]` | List exchange accounts |
| `get(exchange_id)` | `exchange_id: str` | `Exchange` | Get exchange account by ID |
| `list_counterparties()` | None | `List[Any]` | List exchange counterparties |
| `get_withdrawal_fee(exchange_id, to_address_id, amount)` | `exchange_id: str`, `to_address_id: Optional[str]`, `amount: Optional[str]` | `Any` | Get withdrawal fees |

#### Example

```python
# List exchange accounts
exchanges, pagination = client.exchanges.list(limit=50)
for exchange in exchanges:
    print(f"{exchange.name} ({exchange.exchange_label}): {exchange.balance}")

# Get withdrawal fee
fee = client.exchanges.get_withdrawal_fee(exchange_id="123")
```

---

### FiatService

Provides fiat currency operations including provider accounts and exchange rates.

**Access:** `client.fiat`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[FiatProviderAccount], Optional[Pagination]]` | List fiat provider accounts |
| `get_account(account_id)` | `account_id: str` | `FiatProviderAccount` | Get fiat provider account |
| `get_base_currency()` | None | `FiatCurrency` | Get configured base currency |
| `get_rate(from_currency, to_currency)` | `from_currency: str`, `to_currency: str` | `ExchangeRate` | Get exchange rate |
| `list_providers()` | None | `List[Any]` | List fiat providers |

#### Example

```python
# List fiat provider accounts
accounts, pagination = client.fiat.list()
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
| `list(limit, offset, blockchain, network)` | `limit: int = 50`, `offset: int = 0`, `blockchain: Optional[str]`, `network: Optional[str]` | `Tuple[List[FeePayer], Optional[Pagination]]` | List fee payers |
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
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Job], Optional[Pagination]]` | List jobs |
| `get(job_id)` | `job_id: str` | `Job` | Get job by ID |

#### Example

```python
# List jobs
jobs, pagination = client.jobs.list(limit=50)
for job in jobs:
    print(f"{job.id}: {job.description}")
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
| `list_pledges(opts)` | `opts: ListPledgesOptions` | `Tuple[List[Pledge], Optional[Pagination]]` | List pledges |
| `create_pledge(req)` | `req: CreatePledgeRequest` | `Tuple[Pledge, PledgeAction]` | Create pledge |
| `update_pledge(pledge_id, req)` | `pledge_id: str`, `req: UpdatePledgeRequest` | `Pledge` | Update pledge |
| `add_pledge_collateral(pledge_id, req)` | `pledge_id: str`, `req: AddPledgeCollateralRequest` | `Tuple[Pledge, PledgeAction]` | Add collateral |
| `withdraw_pledge(pledge_id, req)` | `pledge_id: str`, `req: WithdrawPledgeRequest` | `Tuple[PledgeWithdrawal, PledgeAction]` | Withdraw (pledgee) |
| `initiate_withdraw_pledge(pledge_id, req)` | `pledge_id: str`, `req: InitiateWithdrawPledgeRequest` | `Tuple[PledgeWithdrawal, PledgeAction]` | Initiate withdrawal (pledgor) |
| `unpledge(pledge_id)` | `pledge_id: str` | `Tuple[Pledge, PledgeAction]` | Unpledge all funds |
| `reject_pledge(pledge_id, req)` | `pledge_id: str`, `req: RejectPledgeRequest` | `Pledge` | Reject pledge |
| `list_pledge_actions(opts)` | `opts: ListPledgeActionsOptions` | `Tuple[List[PledgeAction], Optional[Pagination]]` | List actions |
| `list_pledge_actions_for_approval(opts)` | `opts: ListPledgeActionsOptions` | `Tuple[List[PledgeAction], Optional[Pagination]]` | Get pending actions |
| `approve_pledge_actions(actions, private_key, comment)` | `actions: List[PledgeAction]`, `private_key` | `int` | Approve with signature |
| `reject_pledge_actions(req)` | `req: RejectPledgeActionsRequest` | `int` | Reject actions |
| `list_pledge_withdrawals(opts)` | `opts: ListPledgeWithdrawalsOptions` | `Tuple[List[PledgeWithdrawal], Optional[Pagination]]` | List withdrawals |

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
| `list_lending_offers(opts)` | `opts: ListLendingOffersOptions` | `Tuple[List[LendingOffer], Optional[Pagination]]` | List offers |
| `create_lending_offer(req)` | `req: CreateLendingOfferRequest` | `LendingOffer` | Create offer |
| `cancel_lending_offer(offer_id)` | `offer_id: str` | `LendingOffer` | Cancel offer |
| `get_lending_agreement(agreement_id)` | `agreement_id: str` | `LendingAgreement` | Get agreement |
| `list_lending_agreements(opts)` | `opts: ListLendingAgreementsOptions` | `Tuple[List[LendingAgreement], Optional[Pagination]]` | List agreements |
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
| `list_settlements(opts)` | `opts: ListSettlementsOptions` | `Tuple[List[Settlement], Optional[Pagination]]` | List settlements |
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
| `list_shared_addresses(opts)` | `opts: ListSharedAddressesOptions` | `Tuple[List[SharedAddress], Optional[Pagination]]` | List shared addresses |
| `create_shared_address(req)` | `req: CreateSharedAddressRequest` | `SharedAddress` | Share an address |
| `revoke_shared_address(address_id)` | `address_id: str` | `None` | Revoke sharing |
| `list_shared_assets(opts)` | `opts: ListSharedAssetsOptions` | `Tuple[List[SharedAsset], Optional[Pagination]]` | List shared assets |
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

### Offset-Based Pagination

Most services use offset-based pagination:

```python
all_wallets = []
offset = 0
limit = 50

while True:
    wallets, pagination = client.wallets.list(limit=limit, offset=offset)
    all_wallets.extend(wallets)

    if pagination is None or offset + limit >= pagination.total_items:
        break
    offset += limit
```

### Options-Based Pagination

For advanced filtering:

```python
from taurus_protect.models import ListWalletsOptions

options = ListWalletsOptions(
    currency="ETH",
    exclude_disabled=True,
    limit=50,
    offset=0,
)
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

43 services, 206 public methods.

### ActionService

- `get(action_id: 'str') -> 'Action'` — Get an action by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Action], Optional[Pagination]]'` — List actions with pagination.

### AddressService

- `create(request: 'CreateAddressRequest') -> 'Address'` — Create a new address.
- `create_address(wallet_id: 'int', label: 'str', comment: 'str' = '', customer_id: 'str' = '') -> 'Address'` — Create a new address with explicit parameters.
- `create_attribute(address_id: 'int', key: 'str', value: 'str') -> 'None'` — Create an attribute for an address.
- `delete_attribute(address_id: 'int', attribute_id: 'int') -> 'None'` — Delete an attribute from an address.
- `get(address_id: 'int') -> 'Address'` — Get an address by ID with mandatory signature verification.
- `get_proof_of_reserve(address_id: 'int', challenge: 'Optional[str]' = None) -> 'Any'` — Get the proof of reserve for an address.
- `list(wallet_id: 'int', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Address], Optional[Pagination]]'` — List addresses for a wallet with mandatory signature verification.
- `list_with_options(options: 'Optional[ListAddressesOptions]' = None) -> 'Tuple[List[Address], Optional[Pagination]]'` — List addresses with full filtering options.

### AirGapService

- `get_unsigned_payload(request_id: 'int') -> 'bytes'` — Get unsigned transaction payload for offline signing.
- `submit_signed_payload(request_id: 'int', signed_payload: 'bytes') -> 'None'` — Submit a signed payload back to the system.

### AssetService

- `get(asset_id: 'str') -> 'Asset'` — Get an asset by ID.
- `get_addresses(currency: 'str', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Address], Optional[Pagination]]'` — Get address balances for a specific asset.
- `get_wallets(currency: 'str', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Any], Optional[Pagination]]'` — Get wallet balances for a specific asset.
- `list(currency: 'str' = 'ETH', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Asset], Optional[Pagination]]'` — List asset wallets for a given currency.

### AuditService

- `export_audit_trails(external_user_id: 'Optional[str]' = None, entities: 'Optional[List[str]]' = None, actions: 'Optional[List[str]]' = None, from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, format: 'Optional[str]' = None) -> 'str'` — Export audit trails in the specified format.
- `get(audit_id: 'str') -> 'Audit'` — Get an audit event by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Audit], Optional[Pagination]]'` — List audit events with pagination.

### BalanceService

- `list(currency: 'Optional[str]' = None, limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[AssetBalance], Optional[Pagination]]'` — Get all balances for the tenant, optionally filtered by currency.
- `list_nft_collections(blockchain: 'str', network: 'str', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[NFTCollectionBalance], Optional[Pagination]]'` — Get NFT collection balances for the tenant.

### BlockchainService

- `get(blockchain: 'str', network: 'str' = 'mainnet', include_block_height: 'bool' = False) -> 'Blockchain'` — Get blockchain information.
- `get_by_id(blockchain_id: 'str') -> 'Blockchain'` — Get blockchain by composite ID.
- `list(blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, include_block_height: 'bool' = False) -> 'List[Blockchain]'` — List supported blockchains.

### BusinessRuleService

- `list(page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None, rule_keys: 'Optional[List[str]]' = None, wallet_ids: 'Optional[List[str]]' = None, currency_ids: 'Optional[List[str]]' = None, entity_type: 'Optional[str]' = None, entity_ids: 'Optional[List[str]]' = None) -> 'BusinessRuleResult'` — List business rules with cursor-based pagination (v2 API).
- `list_by_currency(currency_id: 'str', page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'BusinessRuleResult'` — List business rules for a specific currency.
- `list_by_wallet(wallet_id: 'int', page_size: 'Optional[int]' = None, current_page: 'Optional[str]' = None, page_request: 'Optional[str]' = None) -> 'BusinessRuleResult'` — List business rules for a specific wallet.
- `update_transactions_enabled(enabled: 'bool') -> 'None'` — Enable or disable transaction processing for the tenant.

### ChangeService

- `approve_change(change_id: 'str') -> 'None'` — Approve a change.
- `approve_changes(change_ids: 'List[str]') -> 'None'` — Approve multiple changes.
- `create_change(request: 'CreateChangeRequest') -> 'str'` — Create a change request.
- `get(change_id: 'str') -> 'Change'` — Get a change by ID.
- `list(options: 'Optional[ListChangesOptions]' = None) -> 'ChangeResult'` — List changes with cursor-based pagination.
- `list_for_approval(options: 'Optional[ListChangesOptions]' = None) -> 'ChangeResult'` — List changes pending approval with cursor-based pagination.
- `reject_change(change_id: 'str') -> 'None'` — Reject a change.
- `reject_changes(change_ids: 'List[str]') -> 'None'` — Reject multiple changes.

### ConfigService

- `get() -> 'TenantConfig'` — Get tenant configuration.
- `get_features() -> 'List[Feature]'` — Get enabled features for the tenant.

### ContractWhitelistingService

- `approve_whitelisted_contracts(contract_ids: 'List[str]', signature: 'str', comment: 'Optional[str]' = None) -> 'None'` — Approve whitelisted contracts with a signature.
- `create(address: 'str', name: 'str', blockchain: 'str', network: 'Optional[str]' = None, abi: 'Optional[str]' = None) -> 'int'` — Create a whitelisted contract request.
- `create_attribute(contract_id: 'str', key: 'str', value: 'str') -> 'None'` — Create an attribute on a whitelisted contract.
- `delete(contract_id: 'int') -> 'None'` — Delete a whitelisted contract.
- `get_attribute(contract_id: 'str', key: 'str') -> 'Optional[str]'` — Get an attribute value from a whitelisted contract.

### CurrencyService

- `get(currency_id: 'str') -> 'Currency'` — Get a currency by ID.
- `get_base_currency() -> 'Currency'` — Get the base currency configured for the tenant.
- `get_by_blockchain(blockchain: 'str', network: 'str', contract_address: 'Optional[str]' = None, token_id: 'Optional[str]' = None) -> 'Currency'` — Get a currency by blockchain and network.
- `list(show_disabled: 'bool' = False, include_logo: 'bool' = False) -> 'List[Currency]'` — Get all currencies.

### ExchangeService

- `get(exchange_id: 'str') -> 'Exchange'` — Get an exchange account by ID.
- `get_withdrawal_fee(exchange_id: 'str', to_address_id: 'Optional[str]' = None, amount: 'Optional[str]' = None) -> 'Any'` — Get withdrawal fees for an exchange account.
- `list(limit: 'int' = 50, offset: 'int' = 0, currency_id: 'Optional[str]' = None, exchange_label: 'Optional[str]' = None, status: 'Optional[str]' = None, only_positive_balance: 'bool' = False) -> 'Tuple[List[Exchange], Optional[Pagination]]'` — List exchange accounts with pagination.
- `list_counterparties() -> 'List[Any]'` — List exchange counterparties with their exposure limits.

### FeePayerService

- `get(fee_payer_id: 'str') -> 'FeePayer'` — Get a fee payer by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'Tuple[List[FeePayer], Optional[Pagination]]'` — List fee payers with pagination.

### FeeService

- `estimate(currency: 'str', amount: 'Optional[str]' = None, destination: 'Optional[str]' = None) -> 'FeeEstimate'` — Estimate transaction fee.
- `list() -> 'List[FeeEstimate]'` — List fee estimates for all supported currencies.

### FiatService

- `get_account(account_id: 'str') -> 'FiatProviderAccount'` — Get a fiat provider account by ID.
- `get_base_currency() -> 'FiatCurrency'` — Get the configured base currency.
- `get_rate(from_currency: 'str', to_currency: 'str') -> 'ExchangeRate'` — Get the exchange rate between two currencies.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[FiatProviderAccount], Optional[Pagination]]'` — List fiat provider accounts with pagination.
- `list_providers() -> 'List[Any]'` — List available fiat providers.

### GovernanceRuleService

- `approve_rules_proposal(private_key: 'EllipticCurvePrivateKey', comment: 'str', expected_container_hash: 'str') -> 'None'` — Sign the pending proposal's rules container and submit the approval.
- `decode_proposal_for_review(rules: 'GovernanceRules') -> 'DecodedRulesContainer'` — Decode a PENDING proposal so a SuperAdmin can inspect it before approving.
- `get_decoded_rules_container(rules: 'GovernanceRules') -> 'DecodedRulesContainer'` — Get the decoded rules container from governance rules.
- `get_public_keys() -> 'List[SuperAdminPublicKey]'` — Get the list of SuperAdmin public keys.
- `get_rules() -> 'Optional[GovernanceRules]'` — Get the currently enforced governance rules.
- `get_rules_by_id(rules_id: 'str') -> 'Optional[GovernanceRules]'` — Get a governance ruleset by its ID.
- `get_rules_history(page_size: 'int' = 50, cursor: 'Optional[str]' = None) -> 'GovernanceRulesHistoryResult'` — Get the history of governance rules with cursor-based pagination.
- `get_rules_proposal() -> 'Optional[GovernanceRules]'` — Get the proposed governance rules.
- `proposal_container_hash(rules: 'GovernanceRules') -> 'str'` — Return the canonical SHA-256 hex digest of a ruleset's decoded container.
- `reject_rules_proposal(comment: 'str') -> 'None'` — Reject the pending rules proposal with a comment (SuperAdmin only).
- `update_rules_proposal(container: 'DecodedRulesContainer') -> 'None'` — Submit a rules container as a governance proposal (SuperAdmin only).
- `verify_governance_rules(rules: 'GovernanceRules') -> 'GovernanceRules'` — Verify that governance rules have enough valid SuperAdmin signatures.

### GroupService

- `get(group_id: 'str') -> 'Group'` — Get a group by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Group], Optional[Pagination]]'` — List groups with pagination.

### HealthService

- `check() -> 'HealthStatus'` — Check the API health status.
- `get_all_health_checks(tenant_id: 'Optional[str]' = None, fail_if_unhealthy: 'bool' = False) -> 'GetAllHealthChecksResult'` — Get all health checks with optional filtering.

### JobService

- `get(job_id: 'str') -> 'Job'` — Get a job by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Job], Optional[Pagination]]'` — List jobs with pagination.

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
- `list_lending_agreements(options: 'Optional[ListLendingAgreementsOptions]' = None) -> 'Tuple[List[LendingAgreement], Optional[CursorPagination]]'` — List lending agreements.
- `list_lending_agreements_for_approval(options: 'Optional[ListLendingAgreementsOptions]' = None) -> 'Tuple[List[LendingAgreement], Optional[CursorPagination]]'` — List lending agreements pending approval.
- `list_lending_offers(options: 'Optional[ListLendingOffersOptions]' = None) -> 'Tuple[List[LendingOffer], Optional[CursorPagination]]'` — List lending offers.
- `repay_lending_agreement(lending_agreement_id: 'str', request: 'RepayLendingAgreementRequest') -> 'None'` — Record repayment for a lending agreement.
- `update_lending_agreement(lending_agreement_id: 'str', request: 'UpdateLendingAgreementRequest') -> 'None'` — Update a lending agreement.

### MultiFactorSignatureService

- `create_challenge(request_id: 'int', challenge_type: 'str') -> 'str'` — Create a new multi-factor signature challenge.
- `get_challenge(challenge_id: 'str') -> 'MultiFactorSignatureChallenge'` — Get a multi-factor signature challenge by ID.
- `list_challenges(request_id: 'Optional[int]' = None, limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[MultiFactorSignatureChallenge], Optional[Pagination]]'` — List multi-factor signature challenges.
- `verify_challenge(challenge_id: 'str', response: 'str') -> 'bool'` — Verify a multi-factor signature challenge response.

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
- `list_pledge_actions(opts: 'Optional[ListPledgeActionsOptions]' = None) -> 'Tuple[List[PledgeAction], Optional[Pagination]]'` — List all pledge actions with optional filtering.
- `list_pledge_actions_for_approval(opts: 'Optional[ListPledgeActionsOptions]' = None) -> 'Tuple[List[PledgeAction], Optional[Pagination]]'` — List pledge actions pending approval.
- `list_pledge_withdrawals(opts: 'Optional[ListPledgeWithdrawalsOptions]' = None) -> 'Tuple[List[PledgeWithdrawal], Optional[Pagination]]'` — List pledge withdrawals with optional filtering.
- `list_pledges(opts: 'Optional[ListPledgesOptions]' = None) -> 'Tuple[List[Pledge], Optional[Pagination]]'` — List pledges with optional filtering.
- `reject_pledge(pledge_id: 'str', req: 'RejectPledgeRequest') -> 'Pledge'` — Reject a pledge.
- `reject_pledge_actions(req: 'RejectPledgeActionsRequest') -> 'int'` — Reject multiple pledge actions.
- `unpledge(pledge_id: 'str') -> 'Tuple[Pledge, PledgeAction]'` — Unpledge all funds from a pledge.
- `update_pledge(pledge_id: 'str', req: 'UpdatePledgeRequest') -> 'Pledge'` — Update a pledge's default destination.
- `withdraw_pledge(pledge_id: 'str', req: 'WithdrawPledgeRequest') -> 'Tuple[PledgeWithdrawal, PledgeAction]'` — Withdraw from a pledge (pledgee operation).

### PriceService

- `get_current(currency: 'Optional[str]' = None) -> 'List[Price]'` — Get current prices for all currencies or a specific currency.
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
- `get_for_approval(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Request], Optional[Pagination]]'` — Get requests pending approval.
- `list(limit: 'int' = 50, offset: 'int' = 0, from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency_id: 'Optional[str]' = None, statuses: 'Optional[List[RequestStatus]]' = None) -> 'Tuple[List[Request], Optional[Pagination]]'` — List requests with filtering and pagination.
- `reject_request(request_id: 'int', comment: 'str') -> 'None'` — Reject a single request.
- `reject_requests(request_ids: 'List[int]', comment: 'str') -> 'None'` — Reject multiple requests.

### ReservationService

- `cancel(reservation_id: 'int') -> 'None'` — Cancel a reservation.
- `get(reservation_id: 'int') -> 'Reservation'` — Get a reservation by ID.
- `list(wallet_id: 'Optional[int]' = None, limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Reservation], Optional[Pagination]]'` — List reservations.

### ScoreService

- `get_address_score(address_id: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Get risk score for an address.
- `get_transaction_score(tx_hash: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Get risk score for a transaction.
- `refresh_whitelisted_address_score(address_id: 'str', provider: 'Optional[str]' = None) -> 'List[Score]'` — Refresh risk score for a whitelisted address.

### SettlementService

- `cancel_settlement(settlement_id: 'str') -> 'None'` — Cancel a settlement.
- `create_settlement(request: 'CreateSettlementRequest') -> 'str'` — Create a new settlement.
- `get_settlement(settlement_id: 'str') -> 'Settlement'` — Get a settlement by ID.
- `list_settlements(options: 'Optional[ListSettlementsOptions]' = None) -> 'Tuple[List[Settlement], Optional[CursorPagination]]'` — List settlements.
- `list_settlements_for_approval(options: 'Optional[ListSettlementsForApprovalOptions]' = None) -> 'Tuple[List[Settlement], Optional[CursorPagination]]'` — List settlements pending approval.
- `replace_settlement(settlement_id: 'str', request: 'CreateSettlementRequest') -> 'None'` — Replace a settlement with new attributes.

### SharingService

- `list_shared_addresses(options: 'Optional[ListSharedAddressesOptions]' = None) -> 'Tuple[List[SharedAddress], Optional[CursorPagination]]'` — List shared addresses.
- `list_shared_assets(options: 'Optional[ListSharedAssetsOptions]' = None) -> 'Tuple[List[SharedAsset], Optional[CursorPagination]]'` — List shared whitelisted assets.
- `share_address(request: 'ShareAddressRequest') -> 'None'` — Share an address with a Taurus Network participant.
- `share_whitelisted_asset(request: 'ShareWhitelistedAssetRequest') -> 'None'` — Share a whitelisted asset with a Taurus Network participant.
- `unshare_address(shared_address_id: 'str') -> 'None'` — Unshare an address with a Taurus Network participant.
- `unshare_whitelisted_asset(shared_asset_id: 'str') -> 'None'` — Unshare a whitelisted asset with a Taurus Network participant.

### StakingService

- `get_staking_info(address_id: 'int') -> 'StakingInfo'` — Get staking information for an address.
- `list_validators(blockchain: 'str', network: 'str' = 'mainnet', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Validator], Optional[Pagination]]'` — List validators for a blockchain.

### StatisticsService

- `get_summary() -> 'Optional[PortfolioStatistics]'` — Get summary statistics for the portfolio.
- `get_transaction_stats(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None) -> 'TransactionStatistics'` — Get transaction statistics for a date range.

### TagService

- `create(name: 'str', color: 'str') -> 'Tag'` — Create a new tag.
- `delete(tag_id: 'str') -> 'None'` — Delete a tag.
- `get(tag_id: 'str') -> 'Tag'` — Get a tag by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'List[Tag]'` — List tags.

### TokenMetadataService

- `get(blockchain: 'str', contract_address: 'str', token_id: 'str' = '0', network: 'str' = 'mainnet', with_data: 'bool' = False) -> 'Optional[TokenMetadata]'` — Get token metadata for an ERC token.
- `get_crypto_punk(network: 'str', contract_address: 'str', punk_id: 'str', blockchain: 'Optional[str]' = None) -> 'Optional[CryptoPunkMetadata]'` — Get CryptoPunk token metadata.
- `get_erc(network: 'str', contract_address: 'str', token_id: 'str', blockchain: 'Optional[str]' = None, with_data: 'bool' = False) -> 'Optional[TokenMetadata]'` — Get ERC token metadata (ERC721 or ERC1155).
- `get_fa(network: 'str', contract_address: 'str', token_id: 'str' = '0', with_data: 'bool' = False) -> 'Optional[FATokenMetadata]'` — Get FA token metadata (Tezos FA1.2 or FA2).

### TransactionService

- `export_csv(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency: 'Optional[str]' = None, direction: 'Optional[str]' = None, limit: 'int' = 1000, offset: 'int' = 0, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'str'` — Export transactions to CSV format.
- `get(transaction_id: 'int') -> 'Transaction'` — Get a single transaction by ID.
- `get_by_hash(tx_hash: 'str') -> 'Transaction'` — Get a transaction by its blockchain hash.
- `list(from_date: 'Optional[datetime]' = None, to_date: 'Optional[datetime]' = None, currency: 'Optional[str]' = None, direction: 'Optional[str]' = None, limit: 'int' = 50, offset: 'int' = 0, blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None) -> 'Tuple[List[Transaction], Optional[Pagination]]'` — List transactions with filtering.
- `list_by_address(address: 'str', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Transaction], Optional[Pagination]]'` — List transactions for a specific blockchain address.

### UserDeviceService

- `approve_pairing(pairing_id: 'str', nonce: 'str') -> 'None'` — Approve a user device pairing request (Step 3).
- `create_pairing() -> 'UserDevicePairing'` — Create a new user device pairing request (Step 1).
- `get(device_id: 'str') -> 'UserDevicePairingInfo'` — Get user device pairing by ID.
- `get_pairing_status(pairing_id: 'str', nonce: 'str') -> 'UserDevicePairingInfo'` — Get the status of a user device pairing request.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[UserDevicePairing], Optional[Pagination]]'` — List user device pairings.
- `start_pairing(pairing_id: 'str', nonce: 'str', encryption_key: 'str') -> 'None'` — Start a user device pairing request (Step 2).

### UserService

- `create_user_attribute(user_id: 'str', key: 'str', value: 'str') -> 'None'` — Create an attribute for a user.
- `get(user_id: 'str') -> 'User'` — Get a user by ID.
- `get_current() -> 'User'` — Get the current authenticated user.
- `get_users_by_email(emails: 'List[str]') -> 'List[User]'` — Get users by their email addresses.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[User], Optional[Pagination]]'` — List users with pagination.

### VisibilityGroupService

- `get(group_id: 'str') -> 'VisibilityGroup'` — Get a visibility group by ID.
- `get_users(group_id: 'str') -> 'List[Any]'` — Get users in a visibility group.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[VisibilityGroup], Optional[Pagination]]'` — List visibility groups with pagination.

### WalletService

- `create(request: 'CreateWalletRequest') -> 'Wallet'` — Create a new wallet.
- `create_attribute(wallet_id: 'int', key: 'str', value: 'str') -> 'None'` — Create an attribute for a wallet.
- `create_wallet(blockchain: 'str', network: 'str', name: 'str', is_omnibus: 'bool' = False, comment: 'str' = '', customer_id: 'str' = '') -> 'Wallet'` — Create a new wallet with explicit parameters.
- `get(wallet_id: 'int') -> 'Wallet'` — Get a wallet by ID.
- `get_balance_history(wallet_id: 'int', interval_hours: 'int') -> 'List[BalanceHistoryPoint]'` — Get wallet balance history.
- `get_by_name(name: 'str', limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Wallet], Optional[Pagination]]'` — Get wallets by name with pagination.
- `get_tokens(wallet_id: 'int', limit: 'int' = 50) -> 'List[AssetBalance]'` — Get wallet tokens (asset balances).
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Wallet], Optional[Pagination]]'` — List wallets with pagination.
- `list_with_options(options: 'Optional[ListWalletsOptions]' = None) -> 'Tuple[List[Wallet], Optional[Pagination]]'` — List wallets with full filtering options.

### WebhookCallService

- `get(call_id: 'str') -> 'WebhookCall'` — Get a webhook call by ID.
- `get_webhook_calls(event_id: 'Optional[str]' = None, webhook_id: 'Optional[str]' = None, status: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, cursor: 'Optional[ApiRequestCursor]' = None) -> 'WebhookCallResult'` — Retrieve webhook call history with optional filtering.
- `list(webhook_id: 'Optional[str]' = None, event_id: 'Optional[str]' = None, status: 'Optional[str]' = None, sort_order: 'Optional[str]' = None, limit: 'int' = 50, cursor: 'Optional[str]' = None) -> 'Tuple[List[WebhookCall], Optional[Pagination]]'` — List webhook calls with optional filtering.

### WebhookService

- `create(url: 'str', events: 'List[str]') -> 'Webhook'` — Create a new webhook.
- `delete(webhook_id: 'str') -> 'None'` — Delete a webhook.
- `get(webhook_id: 'str') -> 'Webhook'` — Get a webhook by ID.
- `list(limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[Webhook], Optional[Pagination]]'` — List webhooks with pagination.

### WhitelistedAddressService

- `approve(ids: 'List[int]', private_key: 'Any', comment: 'str') -> 'None'` — Sign and submit an approval for the given whitelisted addresses, all-or-nothing.
- `get(whitelisted_address_id: 'int') -> 'WhitelistedAddress'` — Get a whitelisted address by ID with verification.
- `get_envelope(whitelisted_address_id: 'int') -> 'SignedWhitelistedAddressEnvelope'` — Get the signed envelope for a whitelisted address.
- `list(currency: 'Optional[str]' = None, limit: 'int' = 50, offset: 'int' = 0, *, ids: 'Optional[List[str]]' = None, include_for_approval: 'bool' = False) -> 'WhitelistedAddressListResult'` — List whitelisted addresses with cryptographic verification.
- `list_for_approval(limit: 'int' = 50, offset: 'int' = 0, *, ids: 'Optional[List[str]]' = None, include_already_signed_by_user: 'bool' = False) -> 'WhitelistedAddressListResult'` — List whitelisted addresses awaiting approval, verified exactly as ``list`` is.

### WhitelistedAssetService

- `approve(ids: 'List[int]', private_key: 'Any', comment: 'str') -> 'None'` — Sign and submit an approval for the given whitelisted assets, all-or-nothing.
- `get(asset_id: 'int') -> 'WhitelistedAsset'` — Get a whitelisted asset by ID.
- `list(blockchain: 'Optional[str]' = None, network: 'Optional[str]' = None, limit: 'int' = 50, offset: 'int' = 0, *, ids: 'Optional[List[str]]' = None, include_for_approval: 'bool' = False) -> 'Tuple[List[WhitelistedAsset], Optional[Pagination]]'` — List whitelisted assets.
- `list_for_approval(ids: 'Optional[List[str]]' = None, limit: 'int' = 50, offset: 'int' = 0) -> 'Tuple[List[WhitelistedAsset], Optional[Pagination]]'` — List whitelisted assets awaiting approval, verified as in list().

<!-- END GENERATED METHOD INDEX -->
