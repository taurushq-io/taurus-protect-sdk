# CLAUDE.md -- Python SDK

## Naming Conventions

**Taurus Product Names**: Always use hyphenated format: `Taurus-PROTECT`, `Taurus-CAPITAL`, `Taurus-EXPLORER`, `Taurus-PRIME`. Never use space-separated format like "Taurus PROTECT".

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: install + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Build package
./build.sh lint      # black, isort, flake8, mypy
./build.sh format    # Format code
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
./build.sh e2e       # Run E2E tests (requires API access)
./build.sh e2e-one <pattern>  # Run a single E2E test
```

**Single test:** `./build.sh unit-one <pattern>` (e.g., `test_approve`, `TestRequestService`)

**IMPORTANT:** Always use `./build.sh unit` (activates venv with protobuf), NOT bare `python -m pytest` (uses system Python). pytest needs `-o "addopts="` to override coverage flags without pytest-cov.

## Architecture

This SDK provides a Python client for the Taurus-PROTECT API, mirroring the Java and Go SDK architecture.

### Package Structure

- **taurus_protect/_internal/openapi**: Auto-generated OpenAPI client (DO NOT MODIFY)
- **taurus_protect/_internal/proto**: Auto-generated protobuf classes (DO NOT MODIFY)
- **taurus_protect/**: Public SDK package (main development target)
  - **services/**: Service layer wrapping OpenAPI calls
  - **models/**: Domain models (Pydantic) exposed to SDK users
  - **mappers/**: DTO to model conversion functions
  - **helpers/**: Signature verification, validation utilities
  - **cache/**: Thread-safe caching (rules container)
  - **crypto/**: TPV1 authentication and cryptographic utilities
- **tests/**: Unit and integration tests

### Key Patterns

- Services use context manager pattern (`with` statement)
- Properties for lazy service initialization with thread-safe locking
- Pydantic models for validation and serialization
- Type hints throughout (Python 3.9+ -- never use PEP 604 `str | int` syntax, use `Optional[Union[str, int]]`)
- Custom exceptions with `is_retryable()` methods
- HTTP client is urllib3 (via OpenAPI-generated rest.py)
- TPV1-HMAC-SHA256 authentication handled automatically by `AuthenticatedRESTClient` transport

### Available Services (38 + TaurusNetwork namespace)

The ProtectClient provides lazy-initialized properties for all services:

**Core Services**: `wallets`, `addresses`, `requests`, `transactions`, `governance_rules`, `balances`, `currencies`, `whitelisted_addresses`, `whitelisted_assets`

**Transaction/Request Management**: `audits`, `changes`, `fees`, `prices`

**Advanced Features**: `air_gap`, `staking`, `whitelisted_contracts`, `business_rules`, `reservations`, `multi_factor_signature`

**Administrative**: `users`, `groups`, `visibility_groups`, `config`, `webhooks`, `webhook_calls`, `tags`

**Specialized**: `assets`, `actions`, `blockchains`, `exchanges`, `fiat`, `fee_payers`, `health`, `jobs`, `scores`, `statistics`, `token_metadata`, `user_devices`

**Taurus Network** (namespace pattern):
```python
client.taurus_network.participants   # Participant management (5 methods)
client.taurus_network.pledges        # Pledge lifecycle (14 methods, ECDSA approval)
client.taurus_network.lending        # Offers + Agreements (14 methods)
client.taurus_network.settlements    # Settlement operations (6 methods)
client.taurus_network.sharing        # Address/Asset sharing (6 methods)
```

### TaurusNetwork Model Structure

TaurusNetwork models are in `models/taurus_network/` with **71 total models** across 5 files:
- `participant.py` - 7 models (Participant, MyParticipant, ParticipantSettings, etc.)
- `pledge.py` - 25 models (Pledge, PledgeAction, PledgeWithdrawal, enums, requests)
- `lending.py` - 13 models (LendingOffer, LendingAgreement, collaterals, requests)
- `settlement.py` - 11 models (Settlement, SettlementClip, transfers, requests)
- `sharing.py` - 15 models (SharedAddress, SharedAsset, proofs, requests)

Import pattern:
```python
from taurus_protect.models.taurus_network import (
    Participant, Pledge, LendingAgreement, Settlement, SharedAddress
)
```

## SDK Alignment with Java/Go (Source of Truth)

See `docs/SDK_ALIGNMENT_REPORT.md` (repository root) for the full comparison.

## Code Generation

### OpenAPI Generator

Uses `openapi-generator-cli` JAR (7.9.0) with `-g python`. Generated types are prefixed with `Tgvalidatord` (e.g., `TgvalidatordWallet`, `TgvalidatordAddress`).

**Java Requirement**: OpenAPI generation requires **Java 11+**. The script auto-detects Java 22 at `/Users/admin/Library/Java/JavaVirtualMachines/openjdk-22.0.2/Contents/Home` if `JAVA_HOME` is not set and system Java is older.

**Pydantic v2 Compatibility Fix**: The generated code contains patterns like:
```python
Optional[Union[Annotated[bytes, Field(strict=True)], Annotated[str, Field(strict=True)]]]
```
This causes `RuntimeError: Unable to apply constraint 'strict' to schema of type 'none'` with Pydantic v2. The `generate-openapi.sh` script includes a post-processing step that fixes this by replacing with `Optional[Union[bytes, str]]`.

### OpenAPI Type Naming Conventions

- Response types use `result` field (not `wallet`, `wallets`, etc.)
- Create operations often return only an ID, not the full object
- Pagination uses `total_items` and `offset` strings
- API request builders use `body=` parameter for POST/PUT requests
- API method names follow pattern: `{service}_service_{operation}` (e.g., `wallet_service_get_wallet_v2`, `request_service_approve_requests`)

### Protobuf Generator

Uses `protoc` directly with `--python_out` plugin (same approach as Java SDK). Generated files are flattened to `taurus_protect/_internal/proto/`.

**Known Issue - Import Paths**: After flattening, the generated Python files still have incorrect import paths (e.g., `from tp_messages import commitments_pb2`). These imports fail at runtime. For now, security features use Pydantic models instead of protobuf parsing. Full protobuf support requires fixing the import paths in the generated files.

## Service Implementation Pattern

Each service follows this pattern:
1. Wrap OpenAPI API service
2. Use error mapping for converting OpenAPI errors to domain errors
3. Use mapper functions to convert DTOs to domain models
4. Return pagination info when available

### TaurusNetwork Service Pattern

TaurusNetwork services have additional patterns:
- Located in `services/taurus_network/` subdirectory
- Use **cursor-based pagination** (not offset-based) via `CursorPagination` dataclass
- Services receive both `api_client` and specific API instance in `__init__`
- Some services (lending, settlement, sharing) define their own dataclass models inline for simplicity

**TaurusNetworkClient** (`services/taurus_network/_client.py`):
- Namespace client providing lazy-initialized access to 5 sub-services
- Uses `threading.RLock()` for thread-safe initialization
- Accessed via `client.taurus_network.{service_name}`

**Pledge Action Approval** (`services/taurus_network/pledge_service.py`):
- `approve_pledge_actions(actions, private_key, comment)` follows same ECDSA signing pattern as RequestService
- Sorts actions by ID, builds JSON hash array, signs with `crypto.sign_data()`

## Common Implementation Notes

### OpenAPI API Class Names

TaurusNetwork APIs follow this naming pattern:
- `TaurusNetworkParticipantApi` - Participant operations
- `TaurusNetworkPledgeApi` - Pledge operations
- `TaurusNetworkLendingApi` - Lending operations
- `TaurusNetworkSettlementApi` - Settlement operations
- `TaurusNetworkSharedAddressAssetApi` - Shared address/asset operations (note: not `TaurusNetworkSharedApi`)

### Model Export Pattern

Models are exported at multiple levels:
1. Individual model files (`models/taurus_network/pledge.py`)
2. Subpackage init (`models/taurus_network/__init__.py`)
3. Main models init (`models/__init__.py`)

This allows both specific and broad imports:
```python
# Specific import
from taurus_protect.models.taurus_network.pledge import Pledge, CreatePledgeRequest

# Broad import
from taurus_protect.models import Pledge
```

### Service Export Pattern

Services are exported in `services/__init__.py`. The TaurusNetworkClient is NOT exported directly - it's accessed via the ProtectClient's `taurus_network` property.

## Testing

### Running Tests

Use `./build.sh unit` for all unit tests. Use `./build.sh unit-one <pattern>` for specific tests (e.g., `test_approve`, `TestRequestService`).

### Integration Tests

Integration tests are in `tests/integration/`, disabled by default. Enable via environment variables:

```bash
export PROTECT_INTEGRATION_TEST=true
export PROTECT_API_HOST="https://your-api-host.com"
export PROTECT_API_KEY="your-api-key"
export PROTECT_API_SECRET="your-hex-encoded-secret"

./build.sh integration
```

### Shared Test Utilities (`tests/testutil/`)

All test config is centralized in `tests/testutil/`:
- `properties.py` — Key=value `.properties` file parser with `\n` escape for PEM keys
- `config.py` — `TestConfig` with multi-identity support (6 identities), env var overrides
- `helpers.py` — `skip_if_not_enabled()`, `get_test_client(index)`, `skip_if_insufficient_identities()`
- `test.properties.sample` — Sample config matching Java format

Integration and E2E `conftest.py` files delegate to testutil.

### E2E Tests

E2E tests are in `tests/e2e/`, with `conftest.py` delegating to testutil (uses `@pytest.mark.e2e` marker). Same env vars as integration tests.

```bash
./build.sh e2e                           # All E2E tests
./build.sh e2e-one test_multi_currency   # Single E2E test
```

## Build Troubleshooting

### Python SDK with older pip

The `build.sh` script automatically detects pip < 21.3 and falls back to non-editable install mode. You'll see a warning:

```
[WARN] pip X.X is too old for editable installs with pyproject.toml
[WARN] Upgrade pip with: pip3 install --upgrade pip
[WARN] Falling back to non-editable install
```

For best development experience (live code changes without reinstall), upgrade pip:

```bash
pip3 install --upgrade pip
```

## Lessons Learned

### Pydantic Model Immutability

When setting `frozen=True` on Pydantic models, models with mutable caching state cannot be frozen:
- `DecodedRulesContainer` - has `_hsm_public_key` cache
- `GovernanceRules` - has `_decoded_container` cache

Simple data models without internal state can and should use `frozen=True`.

### Use Monotonic Time for Cache Expiry

**Problem:** `time.time()` is vulnerable to system clock changes (NTP adjustments, manual changes).

**Solution:** Use `time.monotonic()` for cache TTL checks:
```python
# WRONG - affected by clock changes
self._cache_timestamp = time.time()
if time.time() - self._cache_timestamp > self._ttl:

# CORRECT - monotonic clock
self._cache_timestamp = time.monotonic()
if time.monotonic() - self._cache_timestamp > self._ttl:
```

### Protobuf to Model Conversion Pitfall

**Problem:** When converting protobuf messages to Python models, it's easy to mix up attribute naming conventions. Protobuf uses camelCase (`groupId`, `minimumSignatures`), while Python models use snake_case (`group_id`, `minimum_signatures`).

**Solution:** After calling a conversion function like `_sequential_thresholds_from_proto(pb)`, the result is a **Python model object**, not a protobuf. Access attributes using snake_case:

```python
# WRONG - accessing converted model with protobuf-style names
parallel_thresholds = [_sequential_thresholds_from_proto(pt) for pt in r.parallelThresholds]
for t in parallel_thresholds:
    t.thresholds[0].groupId  # ERROR: 'GroupThreshold' has no attribute 'groupId'

# CORRECT - use Python model attribute names after conversion
for t in parallel_thresholds:
    t.thresholds[0].group_id  # GroupThreshold model uses snake_case
    t.thresholds[0].minimum_signatures
```

**Protobuf Role Enum Conversion:** Use `Role.Name(role_int)` to convert enum integers to string names:
```python
from taurus_protect._internal.proto import request_reply_pb2
roles = [request_reply_pb2.Role.Name(role) for role in user.roles]  # [5] -> ["HSMSLOT"]
```

### Protobuf Silent Fallback Masking Test Failures

**Problem:** `user_signatures_from_base64()` silently falls back to JSON parsing when `google-protobuf` isn't installed (e.g., running with system Python instead of venv). Protobuf binary data isn't valid JSON, so the fallback returns `[]` -- tests appear to pass but with wrong results.

**Solution:** Always run tests via `./build.sh unit` which activates the venv with `google-protobuf` installed. Never use bare `python -m pytest` which uses system Python.

### OpenAPI Method Name Prefixes

The auto-generated OpenAPI method names use specific prefixes that don't always match the service wrapper names:

| Service Wrapper | OpenAPI Method Prefix | Example |
|---|---|---|
| `BusinessRuleService` | `rule_service_*` | `rule_service_get_business_rules_v2` |
| `ChangeService` | `change_service_*` | `change_service_create_change` |
| `JobService` | `job_service_*` | `job_service_get_jobs` |
| `ContractWhitelistingService` | `whitelist_service_*` | `whitelist_service_get_whitelisted_contracts` |

**Note:** BusinessRuleService uses the v2 endpoint (`rule_service_get_business_rules_v2`) for cursor-based pagination.

Always grep the actual OpenAPI API file to confirm method names before writing service wrappers.

### Python 3.9 Compatibility

If `pyproject.toml` says `requires-python = ">=3.9"`, never use PEP 604 syntax (`str | int`). Use `Optional[Union[str, int]]` instead. PEP 604 requires Python 3.10+.

### Python __init__.py Exports

All exception classes users may catch must be in `__all__`: `APIError`, `AuthenticationError`, `AuthorizationError`, `ConfigurationError`, `IntegrityError`, `NotFoundError`, `RateLimitError`, `RequestMetadataError`, `ValidationError`, `WhitelistError`.
