# CLAUDE.md -- Python SDK

## Naming Conventions

**Taurus Product Names**: Always use hyphenated format: `Taurus-PROTECT`, `Taurus-CAPITAL`, `Taurus-EXPLORER`, `Taurus-PRIME`. Never use space-separated format like "Taurus PROTECT".

## Quick Reference

**Build & test:**
```bash
./build.sh           # Default: install + unit tests
./build.sh unit      # Unit tests only
./build.sh build     # Build package
./build.sh lint      # black, isort, flake8, mypy — ALREADY RED ON MASTER, see below
./build.sh format    # Format code
./build.sh generate  # OpenAPI + protobuf code generation
./build.sh clean     # Clean artifacts
./build.sh e2e       # Run E2E tests (requires API access)
./build.sh e2e-one <pattern>  # Run a single E2E test
```

**Single test:** `./build.sh unit-one <pattern>` (e.g., `test_approve`, `TestRequestService`)

**IMPORTANT:** Always use `./build.sh unit` (activates venv with protobuf), NOT bare `python -m pytest` (uses system Python). pytest needs `-o "addopts="` to override coverage flags without pytest-cov.

### `./build.sh lint` has never passed — do not chase it

Measured on committed master: **96 files fail flake8** (mostly E501 at the 100-char limit), **79 files
would be reformatted by black**, and **mypy reports 3008 errors in 124 files**. It is abandoned, not a
gate. Reformatting the package to chase it is a repo-wide project and would bury a review in churn.

What is still worth doing on a change: keep *new* files free of **real** findings (F401 dead imports,
F841 unused locals) and leave E501 alone, since the surrounding code already exceeds it. Check just
your files:

```bash
.venv/bin/python -m flake8 <paths> --max-line-length=100 | grep -v E501
```

The unit suite is the real signal: `.venv/bin/python -m pytest tests/unit -o "addopts=" -q`.

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

**Advanced Features**: `air_gap`, `staking`, `whitelisted_contracts` (**writes only** — reads live on `whitelisted_assets`, the verified reader of the same endpoint), `business_rules`, `reservations`, `multi_factor_signature`

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

### Governance test facts that cost time to find

- **The container encoder is `rules_container_to_base64`, in
  `taurus_protect/mappers/rules_container_encode.py`** — re-exported from
  `taurus_protect.mappers`, so import it from there. It is NOT in
  `mappers/governance_rules.py` (that file holds the *decode* side,
  `rules_container_from_base64`), and it is not named `container_to_base64` or
  `to_base64`; grepping those finds nothing.
- **`GovernanceRuleService.__init__` takes FOUR args**:
  `(api_client, governance_rules_api, super_admin_keys, min_valid_signatures)` — the
  api_client *and* the api. For a unit test, a `MagicMock()` for the first and a
  `MagicMock()` whose `rule_service_get_rules()` returns a reply works; the keys must be
  real `EllipticCurvePublicKey`s (`ec.generate_private_key(ec.SECP256R1()).public_key()`).
- **The reply DTO needs `rules_container`, `rules_signatures`, `locked`** —
  `_map_rules_from_dto` reads those; signature DTOs need `user_id` and `signature`.
- **Governance rules are verified TWICE**: once in `get_rules()` and again in
  `get_decoded_rules_container()`. So a test asserting that verification happens stays
  green if you remove only one site — take out both when checking it is non-vacuous.
- **Use a wire-valid container when testing that verification runs.** A malformed blob
  raises `IntegrityError` from the decoder, so the test passes whether verification ran
  or not. `TestRulesContainerCacheVerification` keeps a
  `test_the_container_decodes_so_only_signatures_can_fail` guard beside its two
  verification tests for exactly this reason — this trap bit once already.

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

### `build.sh` DELETES `.venv` when pip is unusable — and offline it cannot rebuild it

`build.sh` probes pip and, on failure, logs `Existing .venv is broken (pip not usable) -- recreating` and
**replaces the venv**. In an offline container the fresh venv then has no `pytest`, so `./build.sh unit` dies
with `pytest: command not found` — and any previously working interpreter is already gone. It destroys state to
diagnose it.

Recovery without network, from the `uv` cache (`uv` is on PATH; wheels for pytest, cryptography, pydantic,
protobuf, urllib3 are cached):
```bash
uv pip install --offline --python .venv/bin/python -e ".[dev]"
.venv/bin/python -m pytest tests/unit -o "addopts=" -q     # -o addopts= drops the coverage flags
```
Run tests via `.venv/bin/python -m pytest` afterwards rather than `build.sh unit`, which re-probes pip. Note
`conftest.py` imports `cryptography`, so pytest alone is not enough — install the `[dev]` extra, not just pytest.

### `generate-openapi.sh` Pydantic-fix sed is macOS-only

Line 116 of `scripts/generate-openapi.sh` uses `sed -i '' 's/.../.../' file` (BSD sed). On Linux GNU sed treats the `''` as the input file path and prints `sed: can't read s/...` for every model file. The script **continues anyway** and reports success, but the Pydantic v2 fix never lands and the SDK fails at runtime with `RuntimeError: Unable to apply constraint 'strict' to schema of type 'none'`.

Workaround on Linux until the script is fixed: re-run the Pydantic patch manually after the script:
```bash
find taurus_protect/_internal/openapi/models -name "*.py" -exec sed -i \
  's/Optional\[Union\[Annotated\[bytes, Field(strict=True)\], Annotated\[str, Field(strict=True)\]\]\]/Optional[Union[bytes, str]]/g' {} \;
```

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

## Verification surface added in the 2026-09-04 pass

- **`UnverifiedMetadataError(RequestMetadataError)`** (`taurus_protect/errors.py`, exported from
  the package root). `RequestMetadata._require_verified()` gates `_get_payload_value`, so
  `get_source_address` / `get_destination_address` / `get_amount` raise rather than returning
  `None` on metadata nothing verified. Metadata with no payload is NOT an error.
- **`RequestService._verify_and_mark`** is the single place verification and the
  `hash_verified` flag are set together. Keeping them apart is what let `get()` verify and then
  return a request whose flag was still `False`.
- **The hash-comparison pair lives in `helpers/whitelist_hash_helper.py`**: `verify_hash_coverage`
  and `contains_hash`. Both verifiers carried their own copies and they disagreed — the address
  one returned on the first match, the asset one did not. Do not re-add a local copy.
- **`resolve_rule_key`** (same module) picks the `(blockchain, network)` pair from the signed
  payload; `MAX_PAYLOAD_BYTES` bounds it first.
- **`WhitelistedAsset` gained `decimals` and `token_id`**, sourced from the verified payload
  only. Java and TypeScript already exposed them.
- The list path is **lenient**: an unverifiable envelope is excluded and logged rather than
  raising, and rows-returned-but-none-surviving raises `IntegrityError`.

### The address payload fixture deliberately omits `network`

`_build_payload(network=None)` in `tests/unit/helpers/test_whitelisted_address_verifier.py` is
correct, not an oversight: real signed payloads omit `network` unless the rule sets
`includeNetworkInPayload`. Don't "fix" the default — it is what keeps the fallback path covered.

## Verification surface added in the 2026-09-07 pass

Cross-SDK rules are in the repo-root `CLAUDE.md`. Python-specific:

- **`helpers/price_verifier.py`** — `price_signed_bytes` / `verify_price` / `verify_prices`.
  `models/statistics.py` gained `PriceSignature` and `Price.signatures` (it had no signature
  field at all); the forward ref needs `List` imported in that module.
- **`key_fingerprint` is no longer `_`-prefixed** in `helpers/signature_verifier.py` — both
  verifiers' `_verify_group_threshold` key their signer `set` on it.
- **`whitelist_integrity_helper`'s `verify` parameter is DELETED**, and its three functions are
  off `helpers/__init__.py`'s public surface. `verify=False` skipped everything and
  `verify=True` ran step 1 of six — exactly the skip-verification path the repo records as
  removed. Its own internal caller passed `verify=True` and had to be fixed with it.
- **The asset verifier no longer shadows the DTO chain.** `lookup_blockchain = asset.blockchain
  or dto_blockchain` fed the payload's own value back in as the "DTO side", so
  `resolve_rule_key`'s payload-vs-DTO mismatch check could never fire. A test that calls the
  inner method directly does **not** gate this — the shadowing lived in the caller, so drive it
  through `verify_whitelisted_asset`.
- `asset_service.get_addresses` maps to the domain `Address` (it returned raw generated DTOs)
  and verifies each HSM signature. Its three `except` funnels had to widen to
  `isinstance(e, (APIError, IntegrityError, ValueError))` — otherwise the new `IntegrityError`
  was re-wrapped as a **retryable** `ServerError`. The same funnel bug remains elsewhere; see
  `TODOS.md`.
- **A bare `MagicMock` metadata makes `hash_verified` truthy**, so an approve test can satisfy
  a hash-verified gate by accident and prove nothing. Set the attribute explicitly.

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

### Credentials — client construction (positional-arg footgun)

`ProtectClient.create` is `create(host, credentials=None, api_key=None, api_secret=None, super_admin_keys_pem=None, min_valid_signatures=1, …)`. The api-key params are **deprecated keyword** args, so a positional `create(host, api_key, api_secret)` now MISBINDS the key string to `credentials` and fails. Use `credentials=Credentials.api_key(k, s)` (preferred) or keyword `api_key=…, api_secret=…`. `Credentials` (`taurus_protect/credentials.py`) exposes `api_key` / `bearer_token` / `bearer_token_provider`; SuperAdmin keys are always required (no keys-optional bearer path). `create_from_pem` forwards api_key/api_secret as **keywords** to avoid the same misbind.

### Deleting a service-local model class leaves two orphans

Some domain models live **inside** their service module rather than under `models/` —
`WhitelistedContract` was defined in `services/contract_whitelisting_service.py`. Removing one
means two more edits the interpreter will not prompt for until import time:

- **`services/__init__.py`** re-exports it twice: in the `from … import (…)` block *and* in
  `__all__`. Missing either is an `ImportError` on the whole package.
- **A test module may exist solely for it.** `tests/unit/mappers/test_contract_whitelist.py`
  was 42 lines testing only that class, so it became a collection error rather than a failure —
  pytest reports `Interrupted: 1 error during collection` and runs **nothing**, which looks like
  a broken environment, not a stale test.

Also clear `__pycache__` after deleting a test file (`find tests -name __pycache__ -type d
-exec rm -rf {} +`); a stale `.pyc` for a removed module keeps showing up in greps.

### Extracting a method: re-indent explicitly, never by prefix replacement

Pulling a loop out of a method into a new one means dedenting one level. A blanket
`line.replace('            ', '        ')` corrupts every line nested DEEPER than the base,
because it rewrites the first match anywhere in the leading whitespace — the symptom is
`IndentationError: expected an indented block after 'if' statement`, several lines away
from where you were working. Re-indent by stripping a fixed count from the line start, or
paste the body with the target indentation.

### `WhitelistedAddressService.__init__` takes `whitelisting_api`

Not `address_whitelisting_api` — the keyword differs from the class name and from the
asset service's `pledge_api`/`api` conventions. A wrong kwarg is a `TypeError` at
construction inside the test, which reads as a broken fixture rather than a typo.

### Append new list-filter params, don't insert them

`TransactionService.list()` / `export_csv()` gained `blockchain` and `network` (the API accepts both
and Go/TS already exposed them; this SDK hardcoded `None`). They are appended **after** `offset`
rather than grouped with the other filters, because inserting a parameter mid-signature silently
rebinds any positional call — the same footgun as `ProtectClient.create` above. Prefer appending, or
make new params keyword-only.

### `WhitelistedAsset` is the envelope here (accepted difference)

Go, Java and TS split asset-from-envelope and expose a separate envelope getter. Python merged them:
`WhitelistedAsset` already carries `metadata`, `rules_container`, `rules_signatures` and
`signed_contract_address`, so a caller has raw access from `get()`. **Do not add a
`get_envelope()`** — it would return the same object under a second name. Recorded in the report's
"Remaining known differences".

### Governance Rules Typed Mapper — protobuf unknown-field capture (upb backend)

The governance-rules mappers (`mappers/governance_rules.py`, `mappers/rules_container_encode.py`, `mappers/rule_cell_codec.py`) DO parse protobuf via `taurus_protect._internal.proto.request_reply_pb2` — they implement the cross-SDK typed governance API (see repo-root `CLAUDE.md`).

For schema-evolution losslessness, unknown protobuf fields must be preserved per node. **The installed `protobuf` runtime uses the C/upb backend, where `msg.UnknownFields()` raises `NotImplementedError`.** Capture unknowns without that API: re-parse the message, `ClearField` every known field, then `SerializeToString` — what remains is the unknown-field bytes. Store them on the Pydantic node (`unknown_fields: bytes`) and on encode reattach with `node_pb.MergeFromString(unknown_bytes)` (merges only the unknown-to-us fields — no known-field duplication). Cross-SDK cell byte-parity is enforced by `tests/unit/mappers/test_governance_cell_vectors.py` against `scripts/resources/governance-cell-vectors.json`.

**Serialization must be deterministic.** `rules_container_to_bytes`, `_rule_source_to_bytes` and both
encode entry points in `mappers/rules_container_json.py` call `SerializeToString(deterministic=True)`.
Without it, `map<string, bytes> properties` (present at five levels) emits in unspecified order, so the
same container encodes differently across runs and differs from the other SDKs. The JSON bridge was
missing the flag even though this note already said "add it to any new encode entry point" — the
bytes it produces are what a governance proposal is signed over, so **check every new encode path
against this rule, not just the typed encoder.**

**Catch specific protobuf errors, never bare `Exception`.** The codec and mapper paths catch
`DecodeError` (from `google.protobuf.message`) plus `UnicodeDecodeError`/`ValueError` where a cell
payload is decoded as UTF-8 — a cell field is protobuf `bytes`, so non-UTF-8 content is legal on the
wire and must degrade that one cell to a `RawCell` rather than abort the container.

**`capture_unknown` is needed on nested nodes too**, including the four rule-detail sub-messages and
the whitelisting `RuleSource` (outer message, the `Any`-arm payload check, and the decoded inner
payload). See the repo-root `CLAUDE.md` for the shared vectors that pin all of this.
