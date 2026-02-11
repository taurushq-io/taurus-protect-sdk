# Taurus-PROTECT SDK Documentation

This directory contains documentation that is common across all Taurus-PROTECT SDK implementations.

## Common Documentation

| Document | Description |
|----------|-------------|
| [CONCEPTS.md](CONCEPTS.md) | Domain model, entities, relationships, request lifecycle, governance model |
| [AUTHENTICATION.md](AUTHENTICATION.md) | TPV1 authentication protocol, API credentials, security best practices |
| [INTEGRITY_VERIFICATION.md](INTEGRITY_VERIFICATION.md) | Cryptographic verification flows for governance rules and whitelisted addresses |
| [BUSINESS_RULES.md](BUSINESS_RULES.md) | Business rules, change approval system, entities, actions, and scopes |
| [SDK_ALIGNMENT_REPORT.md](SDK_ALIGNMENT_REPORT.md) | Cross-SDK alignment report (services, security, models, documentation) |

## SDK-Specific Documentation

For language-specific implementation details, see:

### Java SDK

Location: `taurus-protect-sdk-java/docs/`

| Document | Description |
|----------|-------------|
| [SDK_OVERVIEW.md](../taurus-protect-sdk-java/docs/SDK_OVERVIEW.md) | Java architecture, modules, build commands |
| [CONCEPTS.md](../taurus-protect-sdk-java/docs/CONCEPTS.md) | Java model classes, exceptions |
| [AUTHENTICATION.md](../taurus-protect-sdk-java/docs/AUTHENTICATION.md) | Java authentication implementation |
| [SERVICES.md](../taurus-protect-sdk-java/docs/SERVICES.md) | Complete Java API reference |
| [USAGE_EXAMPLES.md](../taurus-protect-sdk-java/docs/USAGE_EXAMPLES.md) | Java code examples |
| [WHITELISTED_ADDRESS_VERIFICATION.md](../taurus-protect-sdk-java/docs/WHITELISTED_ADDRESS_VERIFICATION.md) | Java verification implementation |

### Go SDK

Location: `taurus-protect-sdk-go/docs/`

| Document | Description |
|----------|-------------|
| [SDK_OVERVIEW.md](../taurus-protect-sdk-go/docs/SDK_OVERVIEW.md) | Go architecture, packages, build commands |
| [CONCEPTS.md](../taurus-protect-sdk-go/docs/CONCEPTS.md) | Go model types, error handling |
| [AUTHENTICATION.md](../taurus-protect-sdk-go/docs/AUTHENTICATION.md) | Go authentication implementation |
| [SERVICES.md](../taurus-protect-sdk-go/docs/SERVICES.md) | Complete Go API reference |
| [USAGE_EXAMPLES.md](../taurus-protect-sdk-go/docs/USAGE_EXAMPLES.md) | Go code examples |
| [WHITELISTED_ADDRESS_VERIFICATION.md](../taurus-protect-sdk-go/docs/WHITELISTED_ADDRESS_VERIFICATION.md) | Go verification implementation |

### Python SDK

Location: `taurus-protect-sdk-python/docs/`

| Document | Description |
|----------|-------------|
| [SDK_OVERVIEW.md](../taurus-protect-sdk-python/docs/SDK_OVERVIEW.md) | Python architecture, packages, build commands |
| [CONCEPTS.md](../taurus-protect-sdk-python/docs/CONCEPTS.md) | Python model classes, exceptions |
| [AUTHENTICATION.md](../taurus-protect-sdk-python/docs/AUTHENTICATION.md) | Python authentication implementation |
| [SERVICES.md](../taurus-protect-sdk-python/docs/SERVICES.md) | Complete Python API reference |
| [USAGE_EXAMPLES.md](../taurus-protect-sdk-python/docs/USAGE_EXAMPLES.md) | Python code examples |
| [WHITELISTED_ADDRESS_VERIFICATION.md](../taurus-protect-sdk-python/docs/WHITELISTED_ADDRESS_VERIFICATION.md) | Python verification implementation |

### TypeScript SDK

Location: `taurus-protect-sdk-typescript/docs/`

| Document | Description |
|----------|-------------|
| [SDK_OVERVIEW.md](../taurus-protect-sdk-typescript/docs/SDK_OVERVIEW.md) | TypeScript architecture, packages, build commands |
| [CONCEPTS.md](../taurus-protect-sdk-typescript/docs/CONCEPTS.md) | TypeScript model types, error handling |
| [AUTHENTICATION.md](../taurus-protect-sdk-typescript/docs/AUTHENTICATION.md) | TypeScript authentication implementation |
| [SERVICES.md](../taurus-protect-sdk-typescript/docs/SERVICES.md) | Complete TypeScript API reference |
| [USAGE_EXAMPLES.md](../taurus-protect-sdk-typescript/docs/USAGE_EXAMPLES.md) | TypeScript code examples |
| [WHITELISTED_ADDRESS_VERIFICATION.md](../taurus-protect-sdk-typescript/docs/WHITELISTED_ADDRESS_VERIFICATION.md) | TypeScript verification implementation |

## Quick Start

1. Read [CONCEPTS.md](CONCEPTS.md) to understand the domain model
2. Read [AUTHENTICATION.md](AUTHENTICATION.md) to understand API authentication
3. Choose your SDK and follow the SDK-specific documentation

## Documentation Structure

```
taurus-protect-sdk/
├── docs/                           # Common documentation (this directory)
│   ├── README.md                   # This file
│   ├── CONCEPTS.md                 # Domain model (shared)
│   ├── AUTHENTICATION.md           # TPV1 protocol (shared)
│   ├── INTEGRITY_VERIFICATION.md   # Verification flows (shared)
│   └── BUSINESS_RULES.md           # Business rules & change approval
│
├── taurus-protect-sdk-java/
│   └── docs/                       # Java-specific documentation
│       ├── SDK_OVERVIEW.md
│       ├── CONCEPTS.md             # Java models & exceptions
│       ├── AUTHENTICATION.md       # Java auth implementation
│       ├── SERVICES.md
│       ├── USAGE_EXAMPLES.md
│       └── WHITELISTED_ADDRESS_VERIFICATION.md
│
├── taurus-protect-sdk-go/
│   └── docs/                       # Go-specific documentation
│       ├── SDK_OVERVIEW.md
│       ├── CONCEPTS.md             # Go types & error handling
│       ├── AUTHENTICATION.md       # Go auth implementation
│       ├── SERVICES.md
│       ├── USAGE_EXAMPLES.md
│       └── WHITELISTED_ADDRESS_VERIFICATION.md
│
├── taurus-protect-sdk-python/
│   └── docs/                       # Python-specific documentation
│       ├── SDK_OVERVIEW.md
│       ├── CONCEPTS.md             # Python models & exceptions
│       ├── AUTHENTICATION.md       # Python auth implementation
│       ├── SERVICES.md
│       ├── USAGE_EXAMPLES.md
│       └── WHITELISTED_ADDRESS_VERIFICATION.md
│
└── taurus-protect-sdk-typescript/
    └── docs/                       # TypeScript-specific documentation
        ├── SDK_OVERVIEW.md
        ├── CONCEPTS.md             # TypeScript types & error handling
        ├── AUTHENTICATION.md       # TypeScript auth implementation
        ├── SERVICES.md
        ├── USAGE_EXAMPLES.md
        └── WHITELISTED_ADDRESS_VERIFICATION.md
```
