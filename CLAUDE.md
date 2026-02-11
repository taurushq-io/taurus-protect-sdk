# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# General guidelines

## Workflow orchestration

### 1. Plan Mode Default

* Enter plan mode for ANY non-trivial task (3+ steps or architectural decisions)
* If something goes sideways, STOP and re-plan immediately – don't keep pushing
* Use plan mode for verification steps, not just building
* Write detailed specs upfront to reduce ambiguity
* Write down the plan (to be able to restart later if needed)
* At the end of each plan, give me a list of unresolved questions to answer, if any.

### 2. Subagent Strategy

* Use subagents or teams liberally to keep main context window clean
* Offload research, exploration, and parallel analysis to subagents or teams
* For complex problems, throw more compute at it via subagents or teams
* One task per subagent for focused execution

### 3. Self-Improvement Loop

* After ANY correction from the user: update the "Lessons Learned" section below
* Write rules for yourself that prevent the same mistake
* Ruthlessly iterate on these lessons until mistake rate drops
* Review lessons at session start for relevant project

### 4. Verification Before Done

* Never mark a task complete without proving it works
* Diff behavior between main and your changes when relevant
* Ask yourself: "Would a staff engineer approve this?"
* Run tests, check logs, demonstrate correctness

### 5. Demand Elegance (Balanced)

* For non-trivial changes: pause and ask "is there a more elegant way?"
* If a fix feels hacky: "Knowing everything I know now, implement the elegant solution"
* Skip this for simple, obvious fixes – don't over-engineer
* Challenge your own work before presenting it

### 6. Autonomous Bug Fixing

* When given a bug report: just fix it. Don't ask for hand-holding
* Point at logs, errors, failing tests – then resolve them
* Zero context switching required from the user
* Go fix failing CI tests without being told how

## Task Management

* Plan First: Write the plan to .claude/tasks/ with checkable items
* Verify Plan: Check in before starting implementation
* Track Progress: Mark items complete as you go
* Explain Changes: High-level summary at each step
* Document Results: Add review section to .claude/tasks/
* Capture Lessons: learn (command: `learn.md`) and update `CLAUDE.md` for shared or global project
  context, and `CLAUDE.local.md` for private, developer-specific notes after corrections

## Core Principles

* Simplicity First: Make every change as simple as possible. Impact minimal code.
* No Laziness: Find root causes. No temporary fixes. Senior developer standards.
* Minimal Impact: Changes should only touch what's necessary. Avoid introducing bugs.

## Repository Overview

This is a monorepo containing SDKs for the Taurus-PROTECT API, a cryptocurrency custody and transaction management platform.

### SDK Directories

| Directory | Language | Status | Services |
|-----------|----------|--------|----------|
| `taurus-protect-sdk-java/` | Java | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-go/` | Go | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-python/` | Python | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-typescript/` | TypeScript | Active development | 43 (38 + 5 TaurusNetwork) |

**Service Parity (All SDKs):** All four SDKs have 43 high-level service wrappers (38 core services + 5 TaurusNetwork services). TypeScript also provides access to all 56 OpenAPI-generated APIs. All SDKs split Taurus Network into 5 services (Participants, Pledges, Lending, Settlements, Sharing).

**TypeScript note:** ProtectClient exposes 26 high-level service getters (not 43). 43 service implementations exist in `src/services/` but 17 are not wired as client getters.

## SDK Alignment Reference

When aligning SDKs, the **Java SDK is the source of truth**. A comprehensive alignment report exists at [`docs/SDK_ALIGNMENT_REPORT.md`](docs/SDK_ALIGNMENT_REPORT.md).

## Security Invariants (Cross-SDK)

### Implemented Security Features (All SDKs)

1. **Request approval with private key signing** - `RequestService.approveRequest(request, privateKey)` signs hashes with ECDSA
2. **Address signature verification** - `AddressService.getAddress()` verifies signatures using `RulesContainerCache`
3. **Request hash verification** - `RequestService.getRequest()` verifies `computedHash == providedHash` using constant-time comparison
4. **Whitelisted address 6-step verification** - Full verification flow with legacy hash computation for backward compatibility
5. **Whitelisted asset 5-step verification** - `WhitelistedAssetService` verification with SuperAdmin keys
6. **Governance rules signature verification** - `GovernanceRuleService.verifyGovernanceRules()` with cryptographic ECDSA verification against SuperAdmin keys

### Anti-Patterns

- Never use bare `except Exception` / `catch {}` — always catch specific exceptions
- Always add explicit null/nil checks before constant-time comparison
- Log security cleanup failures — don't silently ignore secret wiping failures
- Pre-fetch shared resources before loops to avoid N+1 patterns

### Cross-SDK Security Rules

- **Constant-time comparison**: Perform dummy comparison on length mismatch; never break/return early in multi-signature loops
- **Address verification mandatory**: `RulesContainerCache` must be provided at construction, never optional
- **Legacy hash**: Address and Asset verifiers use DIFFERENT functions (`ComputeLegacyHashes` vs `ComputeAssetLegacyHashes`)
- **Request hash errors**: INCLUDE computed/provided hash values (not secrets — derived from payload client already has). Whitelisted address/asset errors should NOT include hash values.
- **Field sourcing**: Security-critical fields MUST come from verified payload only, never from unverified DTO
- **Hash exists + no payload**: MUST fail explicitly (never silently return)
- **P-256 curve validation**: All SDKs MUST validate ECDSA keys use P-256 (secp256r1) before use
- **Verification functions must not mutate input**: Return state in result struct instead
- **RulesContainerCache requires SuperAdmin verification**: The rules container must be verified by SuperAdmin key signatures before trusting HSM public keys

## Documentation Structure

Documentation is organized hierarchically to avoid duplication:

- **`docs/`** - Common documentation shared across all SDKs:
  - `README.md` - Documentation index with links to all SDK docs
  - `CONCEPTS.md` - Domain model, entities (Wallet, Address, Request, Transaction, etc.)
  - `AUTHENTICATION.md` - TPV1-HMAC-SHA256 protocol, API credentials, SuperAdmin keys
  - `INTEGRITY_VERIFICATION.md` - Cryptographic verification flows (6-step whitelisted address verification)
  - `BUSINESS_RULES.md` - Business rules, change approval system, entity scopes, and dual-admin workflow

- **`taurus-protect-sdk-java/docs/`** - Java SDK-specific documentation
- **`taurus-protect-sdk-go/docs/`** - Go SDK-specific documentation
- **`taurus-protect-sdk-python/docs/`** - Python SDK-specific documentation
- **`taurus-protect-sdk-typescript/docs/`** - TypeScript SDK-specific documentation

Each SDK directory has:
- `README.md` - Entry point with quick start, services overview, and build commands
- `docs/SDK_OVERVIEW.md` - Architecture, package structure, design patterns
- `docs/SERVICES.md` - Complete API reference for all 43 services
- `docs/CONCEPTS.md` - SDK-specific model classes and exceptions
- `docs/AUTHENTICATION.md` - SDK-specific authentication implementation
- `docs/USAGE_EXAMPLES.md` - Code examples and patterns
- `docs/WHITELISTED_ADDRESS_VERIFICATION.md` - Verification flow implementation

SDK-specific docs reference common docs for shared concepts. When adding features that apply to all SDKs, update the common docs. When adding SDK-specific implementation details, update the SDK-specific docs.

### Documentation Maintenance

When adding new services or features:
1. Update all four SDK `docs/SERVICES.md` files to maintain service parity documentation
2. Update the service count in `docs/SDK_OVERVIEW.md` if adding new services
3. Update the services table in each SDK's `README.md`
4. Ensure cross-references in `docs/CONCEPTS.md` files include all SDKs

## Cross-SDK Lessons Learned

### Service Naming Consistency

When naming services, ensure consistency:

| Service | Java | Go | Python | TypeScript |
|---------|------|-----|--------|------------|
| Business Rules | BusinessRuleService | BusinessRuleService | BusinessRuleService | BusinessRuleService |
| Sharing | TaurusNetworkSharingService | TaurusNetworkSharingService | TaurusNetworkSharingService | SharingService |

Note: Python uses `snake_case` for file names (`business_rule_service.py`) but `PascalCase` for class names.

### Error Handling Alignment

**isRetryable() Pattern:** All SDKs should return `true` for retryable errors:
- HTTP 429 (rate limit)
- HTTP 5xx (server errors)

### RequestMetadataAmount — String Types

The API returns `valueFrom`, `valueTo`, `rate` as JSON strings for arbitrary-precision. All SDKs use string types with a `jsonValueToString()` helper for backward compatibility. Keep `decimals` as integer.

### Whitelisted Address/Asset Field Sourcing

**Key Rules (All SDKs):**
1. `WhitelistedAddressService.list()` MUST verify each envelope (strict mode - fail on first error)
2. `WhitelistedAssetService` must source `name`, `symbol`, `contract_address`, `blockchain`, `network` only from payload
3. If payload is missing a field, the result must be `None` (not DTO value)
4. Non-security fields (`status`, `action`, `rule`, `created_at`) can come from DTO
5. Remove any mapping methods that bypass the verified envelope path

**DTO field extraction:**
- `createdAt` — Extracted from trails array (find "created" action)
- `attributes` — Extracted from DTO attributes array (structured array in Go, key-value map in others)

### Test Infrastructure (Cross-SDK)

**Test hierarchy:** All SDKs use 3 tiers: `testutil/` (shared config + helpers), `integration/` (API tests), `e2e/` (multi-identity workflow tests).

**Config pattern:** All SDKs use `test.properties` file (key=value format) with env var overrides, supporting 6 identities (3 API users + 3 SuperAdmin public keys). Java's `testutil/TestConfig.java` is the reference implementation.

**testutil modules:**
| SDK | Location | Key files |
|-----|----------|-----------|
| Java | `client/src/test/java/.../testutil/` | `TestConfig.java`, `TestHelper.java` |
| Go | `test/testutil/` | `config.go`, `helpers.go`, `properties.go` |
| Python | `tests/testutil/` | `config.py`, `helpers.py`, `properties.py` |
| TypeScript | `tests/testutil/` | `config.ts`, `helpers.ts`, `properties.ts` |

**Integration config delegation:** Integration and E2E directories delegate to testutil (thin wrappers for backward compatibility). Don't duplicate config.

### Integration Test Patterns

- Use SDK-specific helpers (`getTestClientWithVerification()`) for SuperAdmin key config
- Test both positive cases (all fields match) and negative cases (one field differs)
- Fields to test for WhitelistedAddress: blockchain, network, address, customerId, label, memo, addressType
- Fields to test for WhitelistedAsset: blockchain, network, address, name, symbol, decimals, customerId, label

### ChangeService Alignment (comment→changeComment)

When creating changes via the API, the SDK `comment` field must be mapped to OpenAPI `changeComment` field. This applies across all SDKs:
- Java: MapStruct `@Mapping(source = "comment", target = "changeComment")`
- Go: Manual mapping in `CreateChange()`
- Python: `create_change_request_to_dto()` maps `comment` → `change_comment`
- TypeScript: `createChangeRequestToDto()` maps `comment` → `changeComment`

### BusinessRule Pagination (v2 API)

All SDKs must use the v2 API endpoint for listing business rules (`ruleServiceGetBusinessRulesV2`), which returns cursor-based pagination. The v1 endpoint uses offset-based pagination and is deprecated.

### Documentation Verification Quick-Checks

- `RequestMetadata.payload` — intentionally omitted from all SDKs (security)
- `RequestStatus` enum — uses CANCELED (not CANCELLED), PENDING (not PENDING_APPROVAL), BROADCASTED (not BROADCAST)
- Service naming — always singular: BusinessRuleService, WebhookCallService
- All SDKs use 6-step address verification (Step 6: Parse WhitelistedAddress from verified payload)
- Cross-references in common docs must include ALL 4 SDKs
