# Usage Examples

This document provides comprehensive Python code examples for the Taurus-PROTECT SDK.

## Table of Contents

1. [Client Setup](#client-setup)
2. [Wallet Management](#wallet-management)
3. [Address Management](#address-management)
4. [Request Approval Workflow](#request-approval-workflow)
5. [Transaction Queries](#transaction-queries)
6. [Balance Queries](#balance-queries)
7. [Whitelisted Addresses](#whitelisted-addresses)
8. [TaurusNetwork Operations](#taurusnetwork-operations)
9. [Pagination Patterns](#pagination-patterns)
10. [Error Handling](#error-handling)

---

## Client Setup

### Basic Initialization

```python
from taurus_protect import ProtectClient

# Using context manager (recommended)
with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
) as client:
    # Client is automatically closed when exiting the block
    wallets, _ = client.wallets.list()
    print(f"Found {len(wallets)} wallets")
```

### With SuperAdmin Key Verification

```python
from taurus_protect import ProtectClient

# Load SuperAdmin public keys (ECDSA P-256 in PEM format)
super_admin_keys = [
    """-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----""",
    """-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----""",
]

with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
    super_admin_keys_pem=super_admin_keys,
    min_valid_signatures=2,  # Require 2-of-N signatures
    rules_cache_ttl=300.0,   # Cache rules for 5 minutes
    timeout=30.0,            # 30 second timeout
) as client:
    # Governance rules will be verified with SuperAdmin signatures
    rules = client.governance_rules.get_rules()
```

### Environment-Based Configuration

```python
import os
from taurus_protect import ProtectClient

def create_client_from_env() -> ProtectClient:
    """Create a client using environment variables."""
    return ProtectClient.create(
        host=os.environ["PROTECT_API_HOST"],
        api_key=os.environ["PROTECT_API_KEY"],
        api_secret=os.environ["PROTECT_API_SECRET"],
    )

# Usage
with create_client_from_env() as client:
    health = client.health.check()
    print(f"API Status: {health.status}")
```

### Manual Resource Management

```python
from taurus_protect import ProtectClient

# Without context manager (remember to close!)
client = ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
)

try:
    wallets, _ = client.wallets.list()
finally:
    client.close()  # Securely wipes credentials from memory
```

---

## Wallet Management

### List All Wallets

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Simple listing
    wallets, pagination = client.wallets.list(limit=50, offset=0)

    for wallet in wallets:
        print(f"Wallet {wallet.id}: {wallet.name}")
        print(f"  Blockchain: {wallet.blockchain}/{wallet.network}")
        print(f"  Currency: {wallet.currency}")
        if wallet.balance:
            print(f"  Balance: {wallet.balance.total_confirmed}")
```

### Find Wallets by Name

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    wallets, _ = client.wallets.get_by_name("Trading", limit=10)
    for wallet in wallets:
        print(f"{wallet.name}: {wallet.id}")
```

### Advanced Wallet Filtering

```python
from taurus_protect.models import ListWalletsOptions

with ProtectClient.create(host, api_key, api_secret) as client:
    options = ListWalletsOptions(
        currency="ETH",
        exclude_disabled=True,
        query="customer",  # Search term
        limit=100,
        offset=0,
    )
    wallets, pagination = client.wallets.list_with_options(options)
```

### Create a Wallet

```python
from taurus_protect.models import CreateWalletRequest

with ProtectClient.create(host, api_key, api_secret) as client:
    # Using request object
    request = CreateWalletRequest(
        blockchain="ETH",
        network="mainnet",
        name="Customer Deposits",
        is_omnibus=True,
        comment="Pool wallet for customer deposits",
        customer_id="CUST-001",
    )
    wallet = client.wallets.create(request)
    print(f"Created wallet: {wallet.id}")

    # Or using explicit parameters
    wallet = client.wallets.create_wallet(
        blockchain="BTC",
        network="mainnet",
        name="Cold Storage",
        is_omnibus=False,
    )
```

### Wallet Attributes

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Add custom attribute
    client.wallets.create_attribute(
        wallet_id=123,
        key="department",
        value="treasury",
    )
```

### Balance History

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Get balance history with 1-hour intervals
    history = client.wallets.get_balance_history(
        wallet_id=123,
        interval_hours=1,
    )
    for point in history:
        print(f"{point.timestamp}: {point.balance}")
```

### Token Balances

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Get all token balances for a wallet
    tokens = client.wallets.get_tokens(wallet_id=123, limit=100)
    for token in tokens:
        print(f"{token.currency}: {token.total_confirmed}")
```

---

## Address Management

### List Addresses

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # List addresses for a specific wallet
    addresses, pagination = client.addresses.list(
        wallet_id=123,
        limit=50,
        offset=0,
    )

    for addr in addresses:
        print(f"Address {addr.id}: {addr.address}")
        print(f"  Label: {addr.label}")
        print(f"  Status: {addr.status}")
        if addr.balance:
            print(f"  Balance: {addr.balance.total_confirmed}")
```

### Get Single Address

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Address signature is automatically verified
    address = client.addresses.get(address_id=456)
    print(f"Address: {address.address}")
    print(f"Path: {address.address_path}")
```

### Create Address

```python
from taurus_protect.models import CreateAddressRequest

with ProtectClient.create(host, api_key, api_secret) as client:
    request = CreateAddressRequest(
        wallet_id="123",
        label="Customer Deposit",
        comment="Primary deposit address for customer",
        customer_id="CUST-001",
    )
    address = client.addresses.create(request)
    print(f"New address: {address.address}")

    # Or using explicit parameters
    address = client.addresses.create_address(
        wallet_id=123,
        label="Trading Address",
        comment="For trading operations",
    )
```

### Address Attributes

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Add attribute
    client.addresses.create_attribute(
        address_id=456,
        key="purpose",
        value="customer_deposit",
    )

    # Delete attribute
    client.addresses.delete_attribute(
        address_id=456,
        attribute_id=789,
    )
```

### Proof of Reserve

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Get cryptographic proof of reserve
    proof = client.addresses.get_proof_of_reserve(
        address_id=456,
        challenge="random-challenge-string",
    )
    print(f"Proof: {proof}")
```

---

## Request Approval Workflow

### Complete Approval Flow

```python
from cryptography.hazmat.primitives.serialization import load_pem_private_key
from taurus_protect import ProtectClient

# Load your approval private key
with open("approval_key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=None)

with ProtectClient.create(host, api_key, api_secret) as client:
    # Step 1: Create a transfer request
    request = client.requests.create_internal_transfer(
        from_address_id=123,
        to_address_id=456,
        amount="1000000000000000000",  # 1 ETH in wei
    )
    print(f"Created request {request.id}, status: {request.status}")

    # Step 2: Get requests pending approval
    pending_requests, _ = client.requests.get_for_approval(limit=10)
    print(f"Found {len(pending_requests)} request(s) pending approval")

    # Step 3: Approve with ECDSA signature
    if pending_requests:
        signed_count = client.requests.approve_requests(
            pending_requests,
            private_key,
            comment="Batch approval via SDK",
        )
        print(f"Approved {signed_count} request(s)")
```

### Create Different Transfer Types

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Internal transfer (address to address)
    internal = client.requests.create_internal_transfer(
        from_address_id=123,
        to_address_id=456,
        amount="1000000",
    )

    # Internal transfer from omnibus wallet
    from_wallet = client.requests.create_internal_transfer_from_wallet(
        from_wallet_id=100,
        to_address_id=456,
        amount="1000000",
    )

    # External transfer (to whitelisted address)
    external = client.requests.create_external_transfer(
        from_address_id=123,
        to_whitelisted_address_id=789,
        amount="1000000",
    )

    # External transfer from omnibus wallet
    external_wallet = client.requests.create_external_transfer_from_wallet(
        from_wallet_id=100,
        to_whitelisted_address_id=789,
        amount="1000000",
    )

    # Cancel pending transaction
    cancel = client.requests.create_cancel_request(
        address_id=123,
        nonce=42,
    )
```

### Reject Requests

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Reject single request
    client.requests.reject_request(
        request_id=12345,
        comment="Amount exceeds daily limit",
    )

    # Reject multiple requests
    client.requests.reject_requests(
        request_ids=[12345, 12346, 12347],
        comment="Batch rejection - policy violation",
    )
```

### Filter Requests by Status

```python
from taurus_protect.models import RequestStatus

with ProtectClient.create(host, api_key, api_secret) as client:
    # Get only approved requests (ready for broadcast)
    approved, _ = client.requests.list(
        statuses=[RequestStatus.APPROVED, RequestStatus.BROADCAST],
        limit=50,
    )

    # Get failed or rejected requests
    failed, _ = client.requests.list(
        statuses=[RequestStatus.FAILED, RequestStatus.REJECTED],
        limit=50,
    )

    # Get pending requests awaiting approval
    pending, _ = client.requests.list(
        statuses=[RequestStatus.PENDING],
        limit=50,
    )
```

---

## Transaction Queries

### List Transactions

```python
from datetime import datetime, timedelta

with ProtectClient.create(host, api_key, api_secret) as client:
    # Get recent transactions
    transactions, pagination = client.transactions.list(
        limit=100,
        offset=0,
    )

    for tx in transactions:
        print(f"TX {tx.hash}")
        print(f"  Direction: {tx.direction}")
        print(f"  Amount: {tx.amount} {tx.currency}")
        print(f"  Status: {tx.status}")
        print(f"  Block: {tx.block}")
```

---

## Balance Queries

### Get Balances

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    balances, _ = client.balances.list()
    for balance in balances:
        print(f"{balance.currency}: {balance.total_confirmed}")
```

---

## Whitelisted Addresses

### List Whitelisted Addresses

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    addresses, pagination = client.whitelisted_addresses.list(limit=50)
    for addr in addresses:
        print(f"Whitelisted: {addr.name}")
        print(f"  Address: {addr.address}")
        print(f"  Blockchain: {addr.blockchain}/{addr.network}")
        print(f"  Status: {addr.status}")
```

### Get Verified Whitelisted Address

```python
with ProtectClient.create(host, api_key, api_secret,
                          super_admin_keys_pem=super_admin_keys,
                          min_valid_signatures=2) as client:
    # Address is automatically verified (6-step verification)
    address = client.whitelisted_addresses.get(address_id=123)
    print(f"Verified: {address.name} - {address.address}")
```

---

## TaurusNetwork Operations

### Participant Management

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Get my participant info
    me = client.taurus_network.participants.get_my_participant()
    print(f"I am: {me.name} (ID: {me.id})")

    # List all participants
    participants, _ = client.taurus_network.participants.list_participants()
    for p in participants:
        print(f"Participant: {p.name}")
```

### Pledge Operations

```python
from taurus_protect.models.taurus_network import (
    CreatePledgeRequest,
    ListPledgesOptions,
)

with ProtectClient.create(host, api_key, api_secret) as client:
    # Create a pledge
    request = CreatePledgeRequest(
        shared_address_id="shared-addr-123",
        currency_id="ETH",
        amount="1000000000000000000",
        pledge_type="PLEDGEE_WITHDRAWALS_RIGHTS",
    )
    pledge, action = client.taurus_network.pledges.create_pledge(request)
    print(f"Created pledge {pledge.id}")
    print(f"Action {action.id} pending approval")

    # List pledges
    options = ListPledgesOptions(
        statuses=["ACTIVE"],
        direction="OUTGOING",
    )
    pledges, _ = client.taurus_network.pledges.list_pledges(opts=options)

    # Approve pledge actions
    actions, _ = client.taurus_network.pledges.list_pledge_actions_for_approval()
    if actions:
        count = client.taurus_network.pledges.approve_pledge_actions(
            actions, private_key, "Batch approval"
        )
        print(f"Approved {count} action(s)")
```

### Shared Addresses

```python
from taurus_protect.models.taurus_network import CreateSharedAddressRequest

with ProtectClient.create(host, api_key, api_secret) as client:
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
    print(f"Shared address: {shared.id}")
```

---

## Pagination Patterns

### Iterate Through All Pages

```python
from typing import List
from taurus_protect.models import Wallet

def get_all_wallets(client) -> List[Wallet]:
    """Fetch all wallets using pagination."""
    all_wallets = []
    offset = 0
    limit = 100

    while True:
        wallets, pagination = client.wallets.list(limit=limit, offset=offset)
        all_wallets.extend(wallets)

        # Check if we've fetched all items
        if pagination is None:
            break
        if offset + limit >= pagination.total_items:
            break

        offset += limit

    return all_wallets

# Usage
with ProtectClient.create(host, api_key, api_secret) as client:
    wallets = get_all_wallets(client)
    print(f"Total wallets: {len(wallets)}")
```

### Generator Pattern for Memory Efficiency

```python
from typing import Iterator
from taurus_protect.models import Wallet

def iter_wallets(client, page_size: int = 100) -> Iterator[Wallet]:
    """Iterate through wallets without loading all into memory."""
    offset = 0

    while True:
        wallets, pagination = client.wallets.list(limit=page_size, offset=offset)

        for wallet in wallets:
            yield wallet

        if pagination is None or offset + page_size >= pagination.total_items:
            break

        offset += page_size

# Usage
with ProtectClient.create(host, api_key, api_secret) as client:
    for wallet in iter_wallets(client):
        print(f"Processing wallet: {wallet.name}")
```

---

## Error Handling

### Comprehensive Error Handling

```python
from taurus_protect.errors import (
    APIError,
    ValidationError,
    AuthenticationError,
    AuthorizationError,
    NotFoundError,
    RateLimitError,
    ServerError,
    IntegrityError,
    WhitelistError,
    ConfigurationError,
)
import time

def fetch_wallet_safely(client, wallet_id: int):
    """Fetch a wallet with comprehensive error handling."""
    max_retries = 3
    retry_count = 0

    while retry_count < max_retries:
        try:
            return client.wallets.get(wallet_id)

        except NotFoundError:
            print(f"Wallet {wallet_id} not found")
            return None

        except ValidationError as e:
            print(f"Invalid input: {e.message}")
            return None

        except AuthenticationError as e:
            print(f"Authentication failed: {e.message}")
            raise  # Don't retry auth errors

        except AuthorizationError as e:
            print(f"Permission denied: {e.message}")
            raise  # Don't retry permission errors

        except RateLimitError as e:
            delay = e.suggested_retry_delay()
            print(f"Rate limited, waiting {delay.total_seconds()}s...")
            time.sleep(delay.total_seconds())
            retry_count += 1

        except ServerError as e:
            if e.is_retryable():
                delay = e.suggested_retry_delay()
                print(f"Server error, retrying in {delay.total_seconds()}s...")
                time.sleep(delay.total_seconds())
                retry_count += 1
            else:
                raise

        except IntegrityError as e:
            # SECURITY: Never retry integrity errors
            print(f"SECURITY: Integrity verification failed: {e.message}")
            raise

        except APIError as e:
            print(f"API error: {e.message} (code: {e.code})")
            if not e.is_retryable():
                raise
            retry_count += 1

    raise Exception("Max retries exceeded")

# Usage
with ProtectClient.create(host, api_key, api_secret) as client:
    wallet = fetch_wallet_safely(client, 123)
    if wallet:
        print(f"Wallet: {wallet.name}")
```

### Using is_retryable()

```python
from taurus_protect.errors import APIError

with ProtectClient.create(host, api_key, api_secret) as client:
    try:
        wallet = client.wallets.get(123)
    except APIError as e:
        if e.is_retryable():
            delay = e.suggested_retry_delay()
            print(f"Retryable error. Suggested delay: {delay}")
        else:
            print(f"Non-retryable error: {e.message}")
```

### Handling Security Errors

```python
from taurus_protect.errors import IntegrityError, WhitelistError

with ProtectClient.create(host, api_key, api_secret,
                          super_admin_keys_pem=super_admin_keys) as client:
    try:
        address = client.whitelisted_addresses.get(123)
    except IntegrityError as e:
        # Hash mismatch or insufficient SuperAdmin signatures
        # This is a SECURITY issue - do not retry
        print(f"SECURITY ALERT: {e.message}")
        # Log and alert security team
    except WhitelistError as e:
        # Whitelist signature verification failed
        print(f"Whitelist verification failed: {e.message}")
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture details
- [Services Reference](SERVICES.md) - Complete API documentation
- [Authentication](AUTHENTICATION.md) - Security implementation
- [Key Concepts](CONCEPTS.md) - Domain model and exceptions
