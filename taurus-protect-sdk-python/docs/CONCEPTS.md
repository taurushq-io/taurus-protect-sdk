# Key Concepts

This document covers Python-specific implementation details for the Taurus-PROTECT SDK.

For shared domain concepts (Wallets, Addresses, Requests, Transactions, etc.), see the [Common Concepts](../../docs/CONCEPTS.md) documentation.

## Python Model Classes

### Pydantic Models

All domain models are implemented as frozen Pydantic v2 models, providing:

- **Immutability** - Models cannot be modified after creation
- **Validation** - Automatic input validation
- **Serialization** - JSON serialization/deserialization
- **Type Safety** - Full type annotations

### Model Pattern

```python
from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime

class Wallet(BaseModel):
    """Represents a cryptocurrency wallet."""

    model_config = {"frozen": True}  # Makes model immutable

    id: str = Field(..., description="Unique wallet identifier")
    name: str = Field(..., description="Human-readable name")
    blockchain: str = Field(..., description="Blockchain type (ETH, BTC, etc.)")
    network: str = Field(..., description="Network (mainnet, testnet)")
    currency: str = Field(..., description="Primary currency identifier")
    is_omnibus: bool = Field(default=False, description="Whether this is an omnibus wallet")
    customer_id: Optional[str] = Field(default=None, description="External customer reference")
    balance: Optional[Balance] = Field(default=None, description="Current balance")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
```

### Model Benefits

```python
from taurus_protect.models import Wallet

# Models are immutable (frozen)
wallet = Wallet(id="1", name="Test", blockchain="ETH", network="mainnet", currency="ETH")
# wallet.name = "Changed"  # This would raise an error

# Models are hashable (can be used in sets/dict keys)
wallet_set = {wallet}

# Automatic validation
try:
    invalid = Wallet(id=None, name="Test", ...)  # Raises ValidationError
except ValueError:
    pass

# JSON serialization
json_data = wallet.model_dump_json()
wallet_copy = Wallet.model_validate_json(json_data)
```

### Request Models

For create/update operations, the SDK provides dedicated request models:

```python
from taurus_protect.models import CreateWalletRequest, CreateAddressRequest

# Create request models are also Pydantic models
request = CreateWalletRequest(
    blockchain="ETH",
    network="mainnet",
    name="My Wallet",
    is_omnibus=False,
    comment="Optional description",
)

# They provide validation before sending to API
```

---

## Exception Hierarchy

The SDK defines a comprehensive exception hierarchy for handling errors.

### Class Diagram

```
Exception
├── APIError (base class for HTTP errors)
│   ├── ValidationError (400)
│   ├── AuthenticationError (401)
│   ├── AuthorizationError (403)
│   ├── NotFoundError (404)
│   ├── RateLimitError (429)
│   └── ServerError (5xx)
├── IntegrityError (cryptographic verification failure)
├── WhitelistError (whitelist signature verification failure)
├── ConfigurationError (SDK configuration error)
└── RequestMetadataError (request metadata parsing failure)
```

### APIError (Base Class)

Base exception for all HTTP API errors.

```python
from taurus_protect.errors import APIError

class APIError(Exception):
    """Base exception for all Taurus-PROTECT API errors."""

    message: str          # Human-readable error message
    code: int             # HTTP status code
    description: str      # Short error description
    error_code: str       # Application-specific error code
    retry_after: timedelta  # Suggested retry delay (for rate limits)
    original_error: Exception  # Underlying exception
```

**Helper Methods:**

| Method | Returns | Description |
|--------|---------|-------------|
| `is_retryable()` | `bool` | True for 429 and 5xx errors |
| `is_client_error()` | `bool` | True for 4xx errors |
| `is_server_error()` | `bool` | True for 5xx errors |
| `suggested_retry_delay()` | `timedelta` | Recommended wait before retry |

### ValidationError

Raised when input validation fails (HTTP 400).

```python
from taurus_protect.errors import ValidationError

try:
    wallet = client.wallets.get(-1)  # Invalid wallet_id
except ValidationError as e:
    print(f"Validation failed: {e.message}")
    # code = 400
```

### AuthenticationError

Raised when authentication fails (HTTP 401).

```python
from taurus_protect.errors import AuthenticationError

try:
    # With invalid credentials
    wallet = client.wallets.get(123)
except AuthenticationError as e:
    print(f"Auth failed: {e.message}")
    # Check API key and secret
```

### AuthorizationError

Raised when the authenticated user lacks permission (HTTP 403).

```python
from taurus_protect.errors import AuthorizationError

try:
    # Action requires higher privileges
    wallet = client.wallets.create(...)
except AuthorizationError as e:
    print(f"Permission denied: {e.message}")
```

### NotFoundError

Raised when a resource doesn't exist (HTTP 404).

```python
from taurus_protect.errors import NotFoundError

try:
    wallet = client.wallets.get(999999)
except NotFoundError as e:
    print(f"Not found: {e.message}")
```

### RateLimitError

Raised when rate limit is exceeded (HTTP 429).

```python
from taurus_protect.errors import RateLimitError
import time

try:
    # Too many requests
    wallet = client.wallets.get(123)
except RateLimitError as e:
    delay = e.suggested_retry_delay()
    print(f"Rate limited. Retry after: {delay.total_seconds()}s")
    time.sleep(delay.total_seconds())
    # Retry the request
```

### ServerError

Raised for server-side errors (HTTP 5xx).

```python
from taurus_protect.errors import ServerError

try:
    wallet = client.wallets.get(123)
except ServerError as e:
    if e.is_retryable():
        delay = e.suggested_retry_delay()
        print(f"Server error. Retry in {delay.total_seconds()}s")
```

### IntegrityError

Raised when cryptographic verification fails. **Security-critical - never retry.**

```python
from taurus_protect.errors import IntegrityError

try:
    # Hash verification might fail
    request = client.requests.get(123)
except IntegrityError as e:
    print(f"SECURITY: {e.message}")
    # Do NOT retry
    # Alert security team
```

**Causes:**
- Request hash mismatch (computed vs provided)
- Invalid address signature
- Insufficient SuperAdmin signatures
- Invalid SuperAdmin signature

### WhitelistError

Raised when whitelist verification fails.

```python
from taurus_protect.errors import WhitelistError

try:
    # Whitelist signature verification
    address = client.whitelisted_addresses.get(123)
except WhitelistError as e:
    print(f"Whitelist verification failed: {e.message}")
```

**Causes:**
- Invalid hash in metadata
- Insufficient governance rule signatures
- Missing required fields in payload

### ConfigurationError

Raised when SDK configuration is invalid.

```python
from taurus_protect.errors import ConfigurationError

try:
    client = ProtectClient.create(
        host="",  # Empty host
        api_key="key",
        api_secret="secret",
    )
except ConfigurationError as e:
    print(f"Config error: {e.message}")
```

**Causes:**
- Empty host, api_key, or api_secret
- Invalid SuperAdmin key format
- min_valid_signatures exceeds number of keys

### RequestMetadataError

Raised when request metadata cannot be parsed or extracted.

```python
from taurus_protect.errors import RequestMetadataError

try:
    # When parsing request metadata fails
    metadata = parse_request_metadata(data)
except RequestMetadataError as e:
    print(f"Metadata error: {e.message}")
    if e.cause:
        print(f"Underlying cause: {e.cause}")
```

**Causes:**
- Missing required fields in the metadata payload
- Malformed JSON structure
- Type mismatches when parsing values

---

## Helper Methods

### is_retryable()

Determines if an error should be retried.

```python
from taurus_protect.errors import APIError

try:
    result = client.wallets.get(123)
except APIError as e:
    if e.is_retryable():
        # Safe to retry: 429 (rate limit) or 5xx (server error)
        pass
    else:
        # Don't retry: 4xx client errors
        raise
```

**Returns True for:**
- HTTP 429 (Rate Limit)
- HTTP 5xx (Server Errors)

**Returns False for:**
- HTTP 4xx (except 429)
- All other status codes

### suggested_retry_delay()

Returns the recommended wait time before retrying.

```python
from taurus_protect.errors import APIError
import time

try:
    result = client.wallets.get(123)
except APIError as e:
    if e.is_retryable():
        delay = e.suggested_retry_delay()
        print(f"Waiting {delay.total_seconds()} seconds before retry")
        time.sleep(delay.total_seconds())
        # Retry...
```

**Returns:**
- For 429: `retry_after` header value or 1 second default
- For 5xx: 5 seconds
- For other errors: 0 seconds (timedelta(0))

---

## Pagination Model

The SDK uses an offset-based pagination model.

### Pagination Class

```python
from pydantic import BaseModel, Field
from typing import Optional

class Pagination(BaseModel):
    """Pagination information for list responses."""

    model_config = {"frozen": True}

    total_items: int = Field(..., description="Total number of items available")
    offset: int = Field(..., description="Current offset")
    limit: int = Field(..., description="Items per page")
    has_more: bool = Field(..., description="Whether more items exist")
```

### Usage Pattern

```python
# Services return (items, pagination) tuples
wallets, pagination = client.wallets.list(limit=50, offset=0)

if pagination:
    print(f"Total: {pagination.total_items}")
    print(f"Showing: {pagination.offset} to {pagination.offset + len(wallets)}")
    print(f"Has more: {pagination.has_more}")
```

---

## TaurusNetwork Models

TaurusNetwork services use specialized models located in `taurus_protect.models.taurus_network`.

### Model Organization (71 total models)

| File | Models | Purpose |
|------|--------|---------|
| `participant.py` | 7 | Participant, MyParticipant, ParticipantSettings |
| `pledge.py` | 26 | Pledge, PledgeAction, PledgeWithdrawal, requests |
| `lending.py` | 13 | LendingOffer, LendingAgreement, collaterals |
| `settlement.py` | 11 | Settlement, SettlementClip, transfers |
| `sharing.py` | 14 | SharedAddress, SharedAsset, proofs |

### Import Patterns

```python
# Specific imports
from taurus_protect.models.taurus_network.pledge import Pledge, CreatePledgeRequest

# Broad imports
from taurus_protect.models.taurus_network import (
    Participant,
    Pledge,
    LendingAgreement,
    Settlement,
    SharedAddress,
)

# All TaurusNetwork models
from taurus_protect.models import Pledge  # Re-exported at top level
```

---

## Type Annotations

The SDK uses comprehensive type annotations for IDE support and static analysis.

### Common Patterns

```python
from typing import List, Optional, Tuple

# Service methods return typed tuples
def list(
    self,
    limit: int = 50,
    offset: int = 0,
) -> Tuple[List[Wallet], Optional[Pagination]]:
    ...

# Optional parameters
def get(
    self,
    wallet_id: int,
    include_balance: Optional[bool] = None,
) -> Wallet:
    ...

# Generic error handling
from taurus_protect.errors import APIError

def safe_get(client, wallet_id: int) -> Optional[Wallet]:
    try:
        return client.wallets.get(wallet_id)
    except APIError:
        return None
```

### Type Checking

The SDK supports strict mypy type checking:

```bash
mypy taurus_protect --exclude='taurus_protect/_internal'
```

---

## Related Documentation

- [Common Concepts](../../docs/CONCEPTS.md) - Shared domain model
- [SDK Overview](SDK_OVERVIEW.md) - Architecture details
- [Services Reference](SERVICES.md) - API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples
