# Taurus-PROTECT SDK Documentation

This monorepo contains SDKs for the **Taurus-PROTECT API**, a cryptocurrency custody and transaction management platform.

## Documentation

### Common Documentation

For concepts shared across all SDKs, see [docs/](docs/):

| Document | Description |
|----------|-------------|
| [Key Concepts](docs/CONCEPTS.md) | Domain model, entities, relationships, request lifecycle |
| [Authentication](docs/AUTHENTICATION.md) | TPV1 authentication protocol, security best practices |
| [Integrity Verification](docs/INTEGRITY_VERIFICATION.md) | Cryptographic verification flows |

### SDK-Specific Documentation

- **Java SDK**: [taurus-protect-sdk-java/](taurus-protect-sdk-java/) | [docs](taurus-protect-sdk-java/docs/)
- **Go SDK**: [taurus-protect-sdk-go/](taurus-protect-sdk-go/) | [docs](taurus-protect-sdk-go/docs/)
- **Python SDK**: [taurus-protect-sdk-python/](taurus-protect-sdk-python/) | [docs](taurus-protect-sdk-python/docs/)
- **TypeScript SDK**: [taurus-protect-sdk-typescript/](taurus-protect-sdk-typescript/) | [docs](taurus-protect-sdk-typescript/docs/)

---

## Repository Structure

```
taurus-protect-sdk/
├── docs/                              # Common documentation
│   ├── CONCEPTS.md                    # Domain model (shared)
│   ├── AUTHENTICATION.md              # TPV1 protocol (shared)
│   └── INTEGRITY_VERIFICATION.md      # Verification flows (shared)
├── scripts/resources/                 # Shared resources
│   ├── jars/
│   │   └── openapi-generator-cli-7.9.0.jar
│   ├── proto/schema/                  # Protobuf definitions (98 files)
│   │   ├── v1/                        # v1 API schemas
│   │   ├── common/                    # Shared types
│   │   └── third_party/               # Google/gRPC dependencies
│   └── swagger/
│       └── apis.swagger.json          # OpenAPI 2.0 specification
├── taurus-protect-sdk-java/           # Java SDK
├── taurus-protect-sdk-go/             # Go SDK
├── taurus-protect-sdk-python/         # Python SDK
└── taurus-protect-sdk-typescript/     # TypeScript SDK
```

---

## Java SDK

**Requirements:** Java 8+ (runtime), Maven 3.6+ (build), Java 11+ (code generation)

```bash
cd taurus-protect-sdk-java
./build.sh           # Full build (compile + test)
./build.sh unit      # Unit tests only
```

**Documentation:** [taurus-protect-sdk-java/docs/](taurus-protect-sdk-java/docs/)
- [SDK Overview](taurus-protect-sdk-java/docs/SDK_OVERVIEW.md) - Architecture, modules, code generation
- [Services Reference](taurus-protect-sdk-java/docs/SERVICES.md) - Complete API reference
- [Usage Examples](taurus-protect-sdk-java/docs/USAGE_EXAMPLES.md) - Code examples

---

## Go SDK

**Requirements:** Go 1.21+ (runtime), Java 11+ (code generation), protoc + protoc-gen-go (protobuf)

```bash
cd taurus-protect-sdk-go
./build.sh build     # Compile only
./build.sh unit      # Unit tests only
```

**Documentation:** [taurus-protect-sdk-go/docs/](taurus-protect-sdk-go/docs/)
- [SDK Overview](taurus-protect-sdk-go/docs/SDK_OVERVIEW.md) - Architecture, packages, code generation
- [Services Reference](taurus-protect-sdk-go/docs/SERVICES.md) - Complete API reference
- [Usage Examples](taurus-protect-sdk-go/docs/USAGE_EXAMPLES.md) - Code examples

---

## Python SDK

**Requirements:** Python 3.9+ (runtime), Java 11+ (code generation), protoc (protobuf)

```bash
cd taurus-protect-sdk-python
./build.sh           # Full build (install + test)
./build.sh unit      # Unit tests only
```

**Documentation:** [taurus-protect-sdk-python/docs/](taurus-protect-sdk-python/docs/)
- [SDK Overview](taurus-protect-sdk-python/docs/SDK_OVERVIEW.md) - Architecture, packages, code generation
- [Services Reference](taurus-protect-sdk-python/docs/SERVICES.md) - Complete API reference
- [Usage Examples](taurus-protect-sdk-python/docs/USAGE_EXAMPLES.md) - Code examples

---

## TypeScript SDK

**Requirements:** Node.js 18+ (runtime), npm 9+ (package manager), Java 11+ (code generation)

```bash
cd taurus-protect-sdk-typescript
./build.sh build     # Compile only
./build.sh unit      # Unit tests only
```

**Documentation:** [taurus-protect-sdk-typescript/docs/](taurus-protect-sdk-typescript/docs/)
- [SDK Overview](taurus-protect-sdk-typescript/docs/SDK_OVERVIEW.md) - Architecture, packages, code generation
- [Services Reference](taurus-protect-sdk-typescript/docs/SERVICES.md) - Complete API reference
- [Usage Examples](taurus-protect-sdk-typescript/docs/USAGE_EXAMPLES.md) - Code examples

---

## SDK Comparison

| Aspect | Java | Go | Python | TypeScript |
|--------|------|-----|--------|------------|
| **Language Version** | Java 8+ | Go 1.21+ | Python 3.9+ | Node.js 18+ |
| **Build Tool** | Maven | `go` tool | pip/setuptools | npm |
| **HTTP Client** | OkHttp | Standard `net/http` | urllib3 | fetch API |
| **DTO Mapping** | MapStruct (compile-time) | Manual functions | Manual functions + Pydantic models | Manual functions |
| **Services** | 43 | 43 | 43 | 43 |
| **Testing** | JUnit 5 | Go testing | pytest | Jest |
| **Static Analysis** | SpotBugs, PMD, Checkstyle | golangci-lint | black, flake8, mypy | ESLint, TypeScript |

---

## Shared Resources

All SDKs generate code from the same specifications:

### OpenAPI Specification
- **Location**: `scripts/resources/swagger/apis.swagger.json`
- **Format**: Swagger 2.0 / OpenAPI 2.0
- **Content**: 56 REST API services

### Protocol Buffer Schemas
- **Location**: `scripts/resources/proto/schema/`
- **Content**: 98 .proto files defining:
  - v1 API message types
  - Common shared types
  - Third-party dependencies (Google protobuf, gRPC gateway)

---

## Disclaimer

This software is provided **for educational, experimental, and development purposes only**. It is not intended for production use. See [DISCLAIMER.md](DISCLAIMER.md) for full details.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
