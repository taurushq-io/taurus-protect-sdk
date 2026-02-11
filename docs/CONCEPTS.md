# Taurus-PROTECT - Key Concepts

This document explains the core domain concepts of the Taurus-PROTECT platform. These concepts are shared across all SDK implementations (Java, Go, Python, TypeScript).

## Platform Overview

Taurus-PROTECT is a cryptocurrency custody platform that provides secure management of digital assets. The platform manages:

- **Wallets and Addresses** - Hierarchical structure for organizing blockchain accounts
- **Requests** - Transaction intents that go through an approval workflow
- **Transactions** - Actual blockchain movements (incoming deposits, outgoing transfers)
- **Whitelisted Addresses** - Pre-approved external destinations for transfers
- **Governance Rules** - Multi-signature approval policies signed by SuperAdmins

---

## Core Entities

### Wallet

A **Wallet** is a logical container that groups related blockchain addresses.

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `name` | Human-readable identifier |
| `blockchain` | The blockchain type (ETH, BTC, SOL, etc.) |
| `network` | The network (mainnet, testnet, etc.) |
| `currency` | Primary currency identifier |
| `isOmnibus` | If true, this wallet pools funds from multiple customers |
| `customerId` | Optional external customer reference |
| `balance` | Aggregated balance across all addresses |

**Key points:**
- Each wallet is tied to a specific blockchain and network
- A wallet can contain multiple addresses
- Wallet balance is the sum of all its address balances

### Address

An **Address** is an on-chain identifier within a wallet where funds can be received and sent.

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `walletId` | Parent wallet reference |
| `address` | The blockchain address string |
| `label` | Human-readable name |
| `customerId` | Optional customer reference |
| `status` | Status: `created`, `creating`, `signed`, `observed`, `confirmed` |
| `balance` | Current holdings (available + pending) |
| `addressPath` | HD derivation path |

**Key points:**
- Addresses are derived using HD (Hierarchical Deterministic) key derivation
- Each address has a cryptographic signature that proves it was created by Taurus-PROTECT
- SDKs can verify address signatures using governance rules

### Request

A **Request** represents a transaction intent that requires approval before execution.

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `status` | Current state in the approval workflow |
| `currency` | The asset being transferred |
| `type` | Request type (internal, external, incoming) |
| `metadata` | Transaction details with cryptographic hash |
| `approvers` | Who needs to approve this request |

**Request types:**
- **Internal Transfer** - Between addresses within Taurus-PROTECT
- **External Transfer** - To a whitelisted external address
- **Incoming Transfer** - From an exchange to an internal address
- **Cancel** - Cancel a pending transaction (for replaceable transactions)

### Request Metadata

Contains the cryptographically signed details of a request.

| Property | Description |
|----------|-------------|
| `hash` | SHA-256 hash of the payload |
| `payload` | Parsed transaction details (source, destination, amount) |
| `payloadAsString` | Raw JSON string for verification |

### Transaction

A **Transaction** represents an actual movement of funds on the blockchain.

| Property | Description |
|----------|-------------|
| `id` | Internal identifier |
| `hash` | Blockchain transaction hash |
| `direction` | `incoming` or `outgoing` |
| `status` | Blockchain confirmation status |
| `currency` | Currency identifier |
| `amount` | Amount transferred |
| `fee` | Network fee paid |
| `block` | Block number |

**Key points:**
- Incoming transactions are detected automatically by the platform
- Outgoing transactions are created when approved requests are broadcast
- A transaction may or may not be linked to a request (e.g., direct deposits have no request)

### Whitelisted Address

A **WhitelistedAddress** is a pre-approved external destination for outgoing transfers.

| Property | Description |
|----------|-------------|
| `id` | Unique identifier |
| `blockchain` | Target blockchain |
| `network` | Target network |
| `address` | The external address |
| `name` | Human-readable label |
| `status` | Approval status |
| `rulesContainer` | Base64-encoded governance rules |
| `signedAddress` | Payload with user signatures |
| `approvers` | Required approval structure |

**Key points:**
- External transfers can only be sent to whitelisted addresses
- Whitelisting requires governance approval
- Each whitelisted address has cryptographic verification (metadata hash + signatures)

### Balance

Represents the balance state for a wallet or address.

| Property | Description |
|----------|-------------|
| `totalConfirmed` | Total confirmed balance |
| `totalUnconfirmed` | Total including unconfirmed |
| `availableConfirmed` | Available confirmed balance |
| `availableUnconfirmed` | Available including unconfirmed |
| `reservedConfirmed` | Reserved confirmed balance |
| `reservedUnconfirmed` | Reserved including unconfirmed |

---

## Entity Relationships

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           ENTITY RELATIONSHIPS                           │
└─────────────────────────────────────────────────────────────────────────┘

    ┌──────────┐         1:N         ┌──────────┐
    │  Wallet  │ ──────────────────► │ Address  │
    └──────────┘                     └──────────┘
         │                                │
         │                                │
         ▼                                ▼
    ┌──────────┐                    ┌──────────┐
    │ Balance  │                    │ Balance  │
    └──────────┘                    └──────────┘


    ┌──────────┐                    ┌─────────────────┐
    │ Request  │ ──────────────────►│ RequestMetadata │
    └──────────┘                    └─────────────────┘
         │
         │ 0..1:1
         ▼
    ┌─────────────┐
    │ Transaction │  (may exist independently)
    └─────────────┘


    ┌───────────────────┐         ┌──────────────────────────┐
    │ GovernanceRuleset │ ───────►│  RuleUserSignature[]     │
    └───────────────────┘         └──────────────────────────┘


    ┌──────────────────┐          ┌───────────────────────────┐
    │WhitelistedAddress│ ────────►│ SignedWhitelistedAddress  │
    └──────────────────┘          └───────────────────────────┘
                                          │
                                          ▼
                                  ┌──────────────────────┐
                                  │ WhitelistSignature[] │
                                  └──────────────────────┘
```

---

## Request Lifecycle

Requests move through a series of states from creation to completion:

```
                              ┌──────────────────────────────────────────┐
                              │           REQUEST LIFECYCLE              │
                              └──────────────────────────────────────────┘

    ┌─────────┐      ┌──────────┐      ┌────────────┐      ┌────────────┐
    │ CREATED │ ───► │ APPROVING│ ───► │ HSM_SIGNED │ ───► │BROADCASTING│
    └─────────┘      └──────────┘      └────────────┘      └────────────┘
                           │                                      │
                           │                                      ▼
                           ▼                               ┌─────────────┐
                    ┌──────────┐                           │ BROADCASTED │
                    │ REJECTED │                           └─────────────┘
                    └──────────┘                                  │
                                                   ┌──────────────┴──────────────┐
                                                   ▼                             ▼
                                            ┌───────────┐                 ┌────────┐
                                            │ CONFIRMED │                 │ FAILED │
                                            └───────────┘                 └────────┘
```

| Status | Description |
|--------|-------------|
| `CREATED` | Request has been created and is being validated |
| `APPROVING` | Awaiting required approver signatures |
| `REJECTED` | One or more approvers rejected the request |
| `HSM_SIGNED` | Transaction has been cryptographically signed by HSM |
| `BROADCASTING` | Transaction is being submitted to the blockchain |
| `BROADCASTED` | Transaction has been submitted to the blockchain |
| `CONFIRMED` | Transaction confirmed on-chain |
| `FAILED` | Transaction failed (reverted, out of gas, etc.) |

**Approval flow:**
1. User creates a transfer request
2. Request enters `APPROVING` state
3. Required approvers review and sign the request
4. Once threshold is met, request moves to `HSM_SIGNED`
5. System broadcasts the transaction
6. Status updates as blockchain confirms

---

## Governance Model

The governance model defines who can approve what, using a multi-signature approach.

### Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                     GOVERNANCE STRUCTURE                         │
└─────────────────────────────────────────────────────────────────┘

    ┌───────────────────┐
    │ GovernanceRuleset │  (signed by SuperAdmins)
    └───────────────────┘
              │
              ▼
    ┌───────────────────┐
    │  RulesContainer   │  (protobuf-encoded rules)
    └───────────────────┘
              │
              ├─────────────────────────────────┐
              ▼                                 ▼
    ┌─────────────────────┐          ┌──────────────────────────┐
    │   ApproversGroups   │          │ AddressWhitelistingRules │
    └─────────────────────┘          └──────────────────────────┘
              │
              ▼
    ┌─────────────────────┐
    │ Individual Users    │  (with their public keys)
    └─────────────────────┘
```

### Key Concepts

**SuperAdmin:**
- Highest privilege level in the system
- Signs governance rules
- Multiple SuperAdmin signatures required (configurable threshold)

**Approvers Group:**
- A named group of users who can approve requests
- Has a threshold (e.g., "2 of 3 must approve")

**Parallel Approvers Groups:**
- Multiple independent groups that must all approve (AND logic)
- Example: "Operations team AND Compliance team must both approve"

**Sequential Approvers Groups:**
- Groups within a parallel path that must approve in sequence
- Enables complex approval workflows

### Approval Thresholds

```
ParallelThresholds (OR paths - any path can succeed)
  └── SequentialThresholds (AND groups within a path)
        └── GroupThreshold
              ├── groupId
              └── minimumSignatures
```

---

## Related Documentation

### Common Documentation (This Directory)
- [Authentication & TPV1](AUTHENTICATION.md) - API authentication protocol
- [Integrity Verification](INTEGRITY_VERIFICATION.md) - Cryptographic verification flows

### SDK-Specific Documentation

**Java SDK** (`taurus-protect-sdk-java/docs/`)
- [SDK Overview](../taurus-protect-sdk-java/docs/SDK_OVERVIEW.md) - Java architecture and modules
- [Services Reference](../taurus-protect-sdk-java/docs/SERVICES.md) - Java API documentation
- [Usage Examples](../taurus-protect-sdk-java/docs/USAGE_EXAMPLES.md) - Java code examples

**Go SDK** (`taurus-protect-sdk-go/docs/`)
- [SDK Overview](../taurus-protect-sdk-go/docs/SDK_OVERVIEW.md) - Go architecture and packages
- [Services Reference](../taurus-protect-sdk-go/docs/SERVICES.md) - Go API documentation
- [Usage Examples](../taurus-protect-sdk-go/docs/USAGE_EXAMPLES.md) - Go code examples

**Python SDK** (`taurus-protect-sdk-python/docs/`)
- [SDK Overview](../taurus-protect-sdk-python/docs/SDK_OVERVIEW.md) - Python architecture and packages
- [Services Reference](../taurus-protect-sdk-python/docs/SERVICES.md) - Python API documentation
- [Usage Examples](../taurus-protect-sdk-python/docs/USAGE_EXAMPLES.md) - Python code examples

**TypeScript SDK** (`taurus-protect-sdk-typescript/docs/`)
- [SDK Overview](../taurus-protect-sdk-typescript/docs/SDK_OVERVIEW.md) - TypeScript architecture and packages
- [Services Reference](../taurus-protect-sdk-typescript/docs/SERVICES.md) - TypeScript API documentation
- [Usage Examples](../taurus-protect-sdk-typescript/docs/USAGE_EXAMPLES.md) - TypeScript code examples