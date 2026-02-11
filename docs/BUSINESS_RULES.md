# Business Rules & Change Approval System

This document describes the business rules configuration system and the dual-admin change approval workflow used by Taurus-PROTECT. These concepts are shared across all SDK implementations (Java, Go, Python, TypeScript).

## Overview

Business rules are configurable policies that govern transaction limits, currency enabling, trading schedules, and counterparty exposure. Rules are scoped to entity levels (global, currency, wallet, address) where more specific rules override broader ones.

All modifications to business rules — and to most other administrative entities — go through a **dual-admin approval workflow** via the Changes API. One admin proposes a change; a different admin approves it.

---

## Business Rule Structure

Each business rule has the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | uint64 | Unique rule identifier |
| `tenantId` | uint64 | Tenant this rule belongs to |
| `currency` | string | Currency code (e.g., `ETH`, `XLM`) |
| `ruleKey` | string | Rule name (e.g., `max_outgoing_transaction_per_day`) |
| `ruleValue` | string | Rule value (numeric limit, boolean, schedule, etc.) |
| `ruleGroup` | string | Category grouping |
| `ruleDescription` | string | Human-readable description |
| `ruleValidation` | string | Validation constraint for the value |
| `entityType` | string | Scope level (see Entity Types) |
| `entityID` | string | Target entity identifier |
| `currencyInfo` | object | Detailed currency metadata |
| `walletId` | uint64 | **Deprecated** — use `entityType`/`entityID` instead |
| `addressId` | uint64 | **Deprecated** — use `entityType`/`entityID` instead |

---

## Entity Types (Scope Hierarchy)

Rules apply at different levels of specificity. More specific rules override broader ones.

```
┌─────────────────────────────────────────────────┐
│              RULE SCOPE HIERARCHY                │
└─────────────────────────────────────────────────┘

    ┌──────────┐
    │  global  │  Tenant-wide defaults
    └────┬─────┘
         │
    ┌────▼─────┐
    │ currency │  Per-currency overrides (e.g., XLM, ETH)
    └────┬─────┘
         │
    ┌────▼─────┐
    │  wallet  │  Per-wallet overrides
    └────┬─────┘
         │
    ┌────▼─────┐
    │ address  │  Per-address overrides
    └──────────┘

    Additional entity types:
    ┌──────────┐  ┌──────────────────┐  ┌────────────────┐
    │ exchange │  │ exchange_account │  │ tn_participant │
    └──────────┘  └──────────────────┘  └────────────────┘
```

| Entity Type | entityID | Description |
|-------------|----------|-------------|
| `global` | _(blank)_ | Tenant-wide default |
| `currency` | Currency code (e.g., `XLM`) | Per-currency override |
| `wallet` | Wallet ID | Per-wallet override |
| `address` | Address ID | Per-address override |
| `exchange` | Exchange label | Per-exchange setting |
| `exchange_account` | Account ID | Per-exchange-account setting |
| `tn_participant` | Participant ID | Per-Taurus-Network-participant |

---

## Rule Keys & Groups

### Transaction & Amount Limits

| Rule Key | Rule Group | Typical Values | Description |
|----------|-----------|----------------|-------------|
| `transactions_enabled` | workflow | `0` / `1` | Master switch for outgoing transactions |
| `max_outgoing_transaction_per_day` | maximum outgoing transactions per day | Integer (e.g., `666`) | Max outgoing txns per day |
| `max_amount_smallest_unit` | maximum amount per transaction | Integer (e.g., `1000000000`) | Max amount per txn in smallest unit |
| `max_outgoing_fiat_amount` | maximum fiat amount per outgoing transaction | Integer (e.g., `1500`) | Max fiat equivalent per outgoing txn |
| `max_transaction_per_day` | _(varies)_ | Integer (e.g., `14`) | Max total txns per day |
| `coin_enabled` | coins enabling | `0` / `1` | Whether a coin is enabled for trading |
| `default_coin_enabled` | coins enabling | `0` / `1` | Default coin enabled state for currency |
| `default_max_outgoing_transaction_per_day` | _(varies)_ | Integer | Default max outgoing txns per day |
| `default_max_outgoing_fiat_amount` | _(varies)_ | Integer | Default max outgoing fiat amount |

### Workflow & Schedule

| Rule Key | Rule Group | Typical Values | Description |
|----------|-----------|----------------|-------------|
| `seconds_before_hsm_signature` | workflow | Integer (e.g., `30`) | Delay before HSM signing |
| `minutes_between_approvals` | workflow | Integer (e.g., `2880`) | Cooldown between approvals |
| `business_timezone` | trading schedule | IANA timezone (e.g., `Europe/Vienna`) | Business timezone |
| `business_weekdays` | trading schedule | JSON array (e.g., `[Sunday,...,Saturday]`) | Active business days |
| `business_days_from` | trading schedule | Day name (e.g., `Sunday`) | Business week start |
| `business_days_to` | trading schedule | Day name (e.g., `Friday`) | Business week end |
| `business_hours_min` | trading schedule | Time (e.g., `00:00`) | Business hours start |
| `business_hours_max` | trading schedule | Time (e.g., `23:59`) | Business hours end |
| `transactions_time_window` | business hours and validation | Schedule string | Transfer window description |

### Counterparty Exposure

| Rule Key | Rule Group | Typical Values | Description |
|----------|-----------|----------------|-------------|
| `counterparty_exposure_limit` | counterparty exposure limit | Integer (e.g., `3000000`) | Max exposure to a counterparty |
| `counterparty_exposure_health_threshold` | _(varies)_ | Integer | Health warning threshold |

---

## Change Approval System

### Overview

All modifications to business rules (and other entities) go through a dual-admin approval workflow via the Changes API. One admin proposes a change; a different admin approves it. Changes are applied immediately upon approval.

### Change Lifecycle

```
┌─────────────────────────────────────────────────────────┐
│                 CHANGE LIFECYCLE                         │
└─────────────────────────────────────────────────────────┘

  Admin 1                   System                 Admin 2
  ───────                   ──────                 ───────
     │                         │                      │
     │  POST /changes          │                      │
     │  (propose change)       │                      │
     │────────────────────────►│                      │
     │                         │                      │
     │  changeId               │                      │
     │◄────────────────────────│                      │
     │                         │   Status: Created    │
     │                         │                      │
     │                         │  GET /changes/       │
     │                         │  for-approval        │
     │                         │◄─────────────────────│
     │                         │                      │
     │                         │  POST /changes/      │
     │                         │  {id}/approve        │
     │                         │◄─────────────────────│
     │                         │                      │
     │                         │  Status: Approved    │
     │                         │  (applied immediately)│
     │                         │                      │
```

### Change Status

| Status | Description |
|--------|-------------|
| `Created` | Change proposed, awaiting approval |
| `Approved` | Change approved and applied |
| `Rejected` | Change rejected by approver |
| `Canceled` | Change canceled by proposer |

### CreateChangeRequest Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | **Yes** | See Actions table |
| `entity` | string | **Yes** | See Entities table |
| `entityId` | uint64 | No | ID of entity to change |
| `entityUUID` | UUID | No | UUID of entity (for UUID-based entities) |
| `changes` | map | Conditional | Field-to-value map (required for `create`/`update`) |
| `changeComment` | string | No | Human-readable description |

### Actions

| Action | Description |
|--------|-------------|
| `create` | Create a new entity |
| `update` | Modify an existing entity |
| `delete` | Delete an entity |
| `resetpassword` | Reset a user's password |
| `resettotp` | Reset a user's TOTP |
| `resetkeycontainer` | Reset a user's key container |
| `assign` | Assign a user to a group |
| `unassign` | Remove a user from a group |

---

## Changeable Entities & Valid Fields

The `changes` map key names are **all lowercase**. The following table lists the supported entities, actions, and valid change fields.

| Entity | Supported Actions | Valid Change Fields |
|--------|-------------------|---------------------|
| `user` | create, update, delete, resetpassword, resettotp, resetkeycontainer | `firstname`, `lastname`, `status`, `roles`, `externaluserid`, `username`, `publickey`, `email`, `userid`, `keycontainer` |
| `group` | create, update, delete | `name`, `externalgroupid`, `description`, `groupemail` |
| `usergroup` | create, update, delete | `name`, `externalgroupid`, `description` |
| `businessrule` | create, update, delete | `rulekey`, `rulevalue`, `rulewalletid` |
| `exchange` | create, update, delete | `name`, `symbol`, `country`, `website` |
| `price` | create, update, delete | `blockchain`, `currencyfrom`, `currencyto`, `decimals`, `rate`, `source`, `currencyfromid`, `currencytoid` |
| `action` | create, update, delete | `label`, `autoApprove`, `trigger`, `tasks`, `state` |
| `feepayer` | create, update, delete | `name`, `network`, `blockchain`, `address` |
| `userapikey` | create, update, delete | `key`, `description`, `permissions` |
| `securitydomain` | create, update, delete | `name`, `description`, `mode`, `openid_configuration_url` |
| `visibilitygroup` | create, update, delete | `name`, `description`, `members` |
| `uservisibilitygroup` | create, update, delete | `userid`, `visibilitygroupid` |
| `wallet` | create, update, delete | `address`, `network`, `type`, `balance` |
| `whitelistedaddress` | create, update, delete | `address`, `network`, `type`, `description` |
| `manualaccountbalancefreeze` | create, update, delete | `account`, `reason`, `duration` |
| `manualutxofreeze` | create, update, delete | `utxo`, `reason`, `duration` |
| `taurusnetworkparticipant` | create, update, delete | _(see API docs)_ |
| `autotransfereventhandler` | create, update, delete | `monitored_wallet_id`, `payer_address_id`, `trigger_type`, `minimum_fiat_value_token`, `transfer_amount_factor_percentage`, `status` |

---

## Entity-Specific Processing Rules

### All Entities

- You cannot create a duplicate change. The API will return an error.
- Once created, changes cannot be edited — reject and create a new one.
- Once approved, changes are applied immediately.

### User

- Super admin users cannot be deleted.
- A user's Super Admin role cannot be removed.
- Only `firstname`, `lastname`, `email`, `username`, `status`, and `roles` can be changed for Super Admin users.
- All users must have an `externalUserID` (typically an email address).

### Group

- If SSO is enabled, groups changed via SCIM cannot be changed manually.
- Groups are automatically reset for SSO users on the next login, creating a change request.

### Account Freezes

Two entities exist for account freezes, depending on the blockchain model:

- `manualaccountbalancefreeze` — for account-based blockchains (e.g., Ethereum)
- `manualutxofreeze` — for UTXO-based blockchains (e.g., Bitcoin)

---

## API Endpoints

### Business Rules

| Endpoint | Method | Description | Required Role |
|----------|--------|-------------|---------------|
| `/api/rest/v2/businessrules` | GET | List business rules (with filters) | Admin / AdminReadOnly |
| `/api/rest/v1/businessrules/transactions_enabled` | PUT | Toggle requests enabling | Admin |

> **Note:** The v1 `GET /api/rest/v1/businessrules` endpoint is deprecated. Use the v2 endpoint instead.

### Changes

| Endpoint | Method | Description | Required Role |
|----------|--------|-------------|---------------|
| `/api/rest/v1/changes` | GET | List changes (with filters) | Admin |
| `/api/rest/v1/changes` | POST | Create a change | Admin |
| `/api/rest/v1/changes/{id}` | GET | Get a specific change | Admin |
| `/api/rest/v1/changes/{id}/approve` | POST | Approve a single change | Admin |
| `/api/rest/v1/changes/{id}/reject` | POST | Reject a single change | Admin |
| `/api/rest/v1/changes/approve` | POST | Approve multiple changes | Admin |
| `/api/rest/v1/changes/reject` | POST | Reject multiple changes | Admin |
| `/api/rest/v1/changes/for-approval` | GET | List changes pending approval | Admin |

---

## Query Filters

### Business Rules (v2)

| Parameter | Description |
|-----------|-------------|
| `ids` | Filter by specific rule IDs |
| `ruleKeys` | Filter by rule key names |
| `ruleGroups` | Filter by rule group |
| `currencyIds` | Filter by currency |
| `entityType` | Filter by scope level (`global`, `currency`, `wallet`, `address`, `exchange`, `exchange_account`, `tn_participant`) |
| `entityIDs` | Filter by entity identifiers |
| `cursor.*` | Cursor pagination (`pageRequest`, `pageSize`, `currentPage`) |

### Changes

| Parameter | Description |
|-----------|-------------|
| `entity` | Filter by entity type |
| `status` | Filter by change status (`Created` / `Approved` / `Rejected` / `Canceled`) |
| `creatorId` | Filter by proposer |
| `entityIDs` | Filter by entity IDs |
| `entityUUIDs` | Filter by entity UUIDs |
| `sortOrder` | `ASC` or `DESC` (default `DESC`) |
| `cursor.*` | Cursor pagination (`pageRequest`, `pageSize`, `currentPage`) |

---

## SDK Services Reference

| SDK | BusinessRuleService | ChangeService |
|-----|---------------------|---------------|
| Java | `client.getBusinessRuleService()` | `client.getChangeService()` |
| Go | `client.BusinessRules()` | `client.Changes()` |
| Python | `client.business_rules` | `client.changes` |
| TypeScript | `client.businessRules` | `client.changes` |

---

## Related Documentation

### Common Documentation
- [Concepts](CONCEPTS.md) — Core domain entities
- [Authentication](AUTHENTICATION.md) — TPV1 authentication protocol
- [Integrity Verification](INTEGRITY_VERIFICATION.md) — Governance rules verification

### SDK-Specific Documentation
- [Java Services](../taurus-protect-sdk-java/docs/SERVICES.md) — Java API reference
- [Go Services](../taurus-protect-sdk-go/docs/SERVICES.md) — Go API reference
- [Python Services](../taurus-protect-sdk-python/docs/SERVICES.md) — Python API reference
- [TypeScript Services](../taurus-protect-sdk-typescript/docs/SERVICES.md) — TypeScript API reference
