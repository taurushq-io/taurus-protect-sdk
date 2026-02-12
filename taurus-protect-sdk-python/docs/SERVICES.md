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
| `approve_requests(requests, private_key, comment)` | `requests: List[Request]`, `private_key: EllipticCurvePrivateKey` | `int` | Approve with signature |
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

Provides whitelisted asset/contract management with 5-step verification.

**Access:** `client.whitelisted_assets`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(asset_id)` | `asset_id: int` | `WhitelistedAsset` | Get with verification |
| `list(...)` | Filters | `Tuple[List[WhitelistedAsset], Optional[Pagination]]` | List assets |

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

Provides smart contract whitelisting management.

**Access:** `client.contract_whitelisting`

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `get(contract_id)` | `contract_id: int` | `WhitelistedContract` | Get whitelisted contract |
| `list(blockchain, network, limit, offset)` | `blockchain: Optional[str]`, `network: Optional[str]`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[WhitelistedContract], Optional[Pagination]]` | List whitelisted contracts |
| `create(address, name, blockchain, network, abi)` | `address: str`, `name: str`, `blockchain: str`, `network: Optional[str]`, `abi: Optional[str]` | `int` | Create whitelisted contract |
| `delete(contract_id)` | `contract_id: int` | `None` | Delete whitelisted contract |
| `approve_whitelisted_contracts(contract_ids, signature, comment)` | `contract_ids: List[str]`, `signature: str`, `comment: Optional[str]` | `None` | Approve contracts with signature |
| `create_attribute(contract_id, key, value)` | `contract_id: str`, `key: str`, `value: str` | `None` | Add attribute |
| `get_attribute(contract_id, key)` | `contract_id: str`, `key: str` | `Optional[str]` | Get attribute value |

#### Example

```python
# List whitelisted contracts
contracts, pagination = client.contract_whitelisting.list(blockchain="ETH")
for contract in contracts:
    print(f"{contract.name}: {contract.address}")

# Create a whitelisted contract
contract_id = client.contract_whitelisting.create(
    address="0x1234...",
    name="USDC Token",
    blockchain="ETH",
)
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

#### Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `list(limit, offset)` | `limit: int = 50`, `offset: int = 0` | `Tuple[List[Asset], Optional[Pagination]]` | List assets |
| `get(asset_id)` | `asset_id: str` | `Asset` | Get asset by ID |
| `get_wallets(currency, limit, offset)` | `currency: str`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Any], Optional[Pagination]]` | Get wallet balances for an asset |
| `get_addresses(currency, limit, offset)` | `currency: str`, `limit: int = 50`, `offset: int = 0` | `Tuple[List[Any], Optional[Pagination]]` | Get address balances for an asset |

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
