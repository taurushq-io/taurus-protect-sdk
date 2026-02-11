# Taurus-PROTECT Python SDK

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Python SDK for interacting with the Taurus-PROTECT API, providing secure cryptocurrency custody and transaction management capabilities.

## Documentation

| Document | Description |
|----------|-------------|
| [Key Concepts](docs/CONCEPTS.md) | Python models, exceptions, and domain concepts |
| [SDK Overview](docs/SDK_OVERVIEW.md) | Architecture, packages, and design patterns |
| [Authentication](docs/AUTHENTICATION.md) | TPV1 authentication and cryptographic operations |
| [Services Reference](docs/SERVICES.md) | Complete API documentation for all 43 services |
| [Usage Examples](docs/USAGE_EXAMPLES.md) | Python code examples and common patterns |
| [Whitelisted Address Verification](docs/WHITELISTED_ADDRESS_VERIFICATION.md) | 6-step verification flow |

## Quick Start

### Prerequisites

- Python 3.9 or higher
- pip package manager

### Installation

Install the SDK using pip:

```bash
pip install taurus-protect-sdk
```

Or install from source:

```bash
git clone https://github.com/taurushq-io/taurus-protect-sdk.git
cd taurus-protect-sdk/taurus-protect-sdk-python
pip install -e .
```

### Dependencies

The SDK requires the following packages (installed automatically):

| Package | Version | Purpose |
|---------|---------|---------|
| `urllib3` | >=1.26.0 | HTTP client for API calls (via OpenAPI-generated client) |
| `pydantic` | >=2.0.0 | Data validation and models |
| `cryptography` | >=41.0.0 | ECDSA signing and verification |
| `protobuf` | >=4.21.0 | Protocol buffer support |

### Client Initialization

```python
from taurus_protect import ProtectClient

# Basic initialization
with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
) as client:
    # Use the client
    wallets, _ = client.wallets.list()
```

For production use with SuperAdmin key verification:

```python
from taurus_protect import ProtectClient

super_admin_keys = [
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----",
]

with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key="your-api-key",
    api_secret="your-api-secret-hex",
    super_admin_keys_pem=super_admin_keys,
    min_valid_signatures=2,
) as client:
    # All governance rule verifications will require 2 valid signatures
    pass
```

See [Authentication](docs/AUTHENTICATION.md) for more initialization options.

## Services

The SDK provides 43 services organized into core services and the TaurusNetwork namespace.

### Core Services (38 services)

| Service | Access | Purpose |
|---------|--------|---------|
| `WalletService` | `client.wallets` | Wallet creation, retrieval, balance history |
| `AddressService` | `client.addresses` | Address management, proof of reserve |
| `RequestService` | `client.requests` | Transaction requests and approvals |
| `TransactionService` | `client.transactions` | Transaction queries and export |
| `BalanceService` | `client.balances` | Balance queries across assets |
| `CurrencyService` | `client.currencies` | Currency and blockchain information |
| `GovernanceRuleService` | `client.governance_rules` | Governance rules with signature verification |
| `WhitelistedAddressService` | `client.whitelisted_addresses` | Address whitelisting with verification |
| `WhitelistedAssetService` | `client.whitelisted_assets` | Asset/contract whitelisting |
| `AuditService` | `client.audits` | Audit log queries |
| `ChangeService` | `client.changes` | Configuration change tracking |
| `FeeService` | `client.fees` | Transaction fee information |
| `PriceService` | `client.prices` | Price data and history |
| `AirGapService` | `client.air_gap` | Air-gap signing operations |
| `StakingService` | `client.staking` | Multi-chain staking information |
| `ContractWhitelistingService` | `client.contract_whitelisting` | Smart contract whitelisting |
| `BusinessRuleService` | `client.business_rules` | Business rule management |
| `ReservationService` | `client.reservations` | Balance reservations |
| `MultiFactorSignatureService` | `client.multi_factor_signature` | Multi-factor signature operations |
| `UserService` | `client.users` | User management |
| `GroupService` | `client.groups` | User group management |
| `VisibilityGroupService` | `client.visibility_groups` | Visibility group management |
| `ConfigService` | `client.config` | System configuration |
| `WebhookService` | `client.webhooks` | Webhook management |
| `WebhookCallService` | `client.webhook_calls` | Webhook call history |
| `TagService` | `client.tags` | Tag management |
| `AssetService` | `client.assets` | Asset information |
| `ActionService` | `client.actions` | Action management |
| `BlockchainService` | `client.blockchains` | Blockchain information |
| `ExchangeService` | `client.exchanges` | Exchange integration |
| `FiatService` | `client.fiat` | Fiat currency operations |
| `FeePayerService` | `client.fee_payers` | Fee payer management |
| `HealthService` | `client.health` | API health checks |
| `JobService` | `client.jobs` | Background job management |
| `ScoreService` | `client.scores` | Risk scoring operations |
| `StatisticsService` | `client.statistics` | Platform statistics |
| `TokenMetadataService` | `client.token_metadata` | Token metadata information |
| `UserDeviceService` | `client.user_devices` | User device management |

### TaurusNetwork Services (5 services)

| Service | Access | Purpose |
|---------|--------|---------|
| `ParticipantService` | `client.taurus_network.participants` | Participant management |
| `PledgeService` | `client.taurus_network.pledges` | Pledge lifecycle operations |
| `LendingService` | `client.taurus_network.lending` | Lending offers and agreements |
| `SettlementService` | `client.taurus_network.settlements` | Settlement operations |
| `SharingService` | `client.taurus_network.sharing` | Address and asset sharing |

See [Services Reference](docs/SERVICES.md) for complete API documentation.

## Basic Usage

### List Wallets

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # List wallets with pagination
    wallets, pagination = client.wallets.list(limit=50, offset=0)

    for wallet in wallets:
        print(f"{wallet.name}: {wallet.currency} ({wallet.blockchain}/{wallet.network})")

    if pagination:
        print(f"Total wallets: {pagination.total_items}")
```

### Create a Wallet

```python
from taurus_protect.models import CreateWalletRequest

with ProtectClient.create(host, api_key, api_secret) as client:
    request = CreateWalletRequest(
        blockchain="ETH",
        network="mainnet",
        name="Trading Wallet",
    )
    wallet = client.wallets.create(request)
    print(f"Created wallet ID: {wallet.id}")
```

### Approve Transaction Requests

```python
from cryptography.hazmat.primitives.serialization import load_pem_private_key

# Load your private key
with open("private_key.pem", "rb") as f:
    private_key = load_pem_private_key(f.read(), password=None)

with ProtectClient.create(host, api_key, api_secret) as client:
    # Get requests pending approval
    requests, _ = client.requests.get_for_approval(limit=10)

    if requests:
        # Approve with ECDSA signature
        signed_count = client.requests.approve_requests(requests, private_key)
        print(f"Approved {signed_count} request(s)")
```

### TaurusNetwork Operations

```python
with ProtectClient.create(host, api_key, api_secret) as client:
    # Get my participant info
    me = client.taurus_network.participants.get_my_participant()
    print(f"Participant: {me.name}")

    # List pledges
    pledges, cursor = client.taurus_network.pledges.list_pledges(limit=10)
    for pledge in pledges:
        print(f"Pledge {pledge.id}: {pledge.status}")

    # List shared addresses
    addresses, _ = client.taurus_network.sharing.list_shared_addresses()
```

### Error Handling

```python
from taurus_protect.errors import (
    APIError,
    NotFoundError,
    RateLimitError,
    IntegrityError,
)
import time

with ProtectClient.create(host, api_key, api_secret) as client:
    try:
        wallet = client.wallets.get(999999)
    except NotFoundError:
        print("Wallet not found")
    except RateLimitError as e:
        # Use built-in retry helper
        delay = e.suggested_retry_delay()
        print(f"Rate limited, retry after {delay.total_seconds()}s")
        time.sleep(delay.total_seconds())
    except IntegrityError as e:
        # Security error - DO NOT retry
        print(f"Integrity verification failed: {e.message}")
    except APIError as e:
        if e.is_retryable():
            print(f"Retryable error: {e.message}")
        else:
            print(f"Non-retryable error: {e.message}")
```

See [Usage Examples](docs/USAGE_EXAMPLES.md) for comprehensive examples.

## Build Commands

```bash
# Install and test (default) - USE THIS TO VERIFY CHANGES
./build.sh

# Build package only
./build.sh build

# Run unit tests only
./build.sh unit

# Run integration tests (requires API access)
./build.sh integration

# Run tests only
./build.sh test

# Run linters (black, isort, flake8, mypy)
./build.sh lint

# Format code
./build.sh format

# Code generation (OpenAPI + protobuf)
./build.sh generate

# Clean build artifacts
./build.sh clean
```

## Development

### Running Tests

```bash
# Run all unit tests (includes coverage)
./build.sh unit

# Run a specific test file
./build.sh unit-one test_tpv1

# Run a specific test by name
./build.sh unit-one test_sign_request
```

### Integration Tests

Integration tests require environment configuration:

```bash
export PROTECT_INTEGRATION_TEST=true
export PROTECT_API_HOST="https://your-api-host.com"
export PROTECT_API_KEY="your-api-key"
export PROTECT_API_SECRET="your-hex-encoded-secret"

./build.sh integration
```

## License

This code is copyright (c) 2025 Taurus SA. It is released under the [MIT license](./LICENSE).
