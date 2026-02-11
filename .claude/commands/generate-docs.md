# Generate Project Documentation

Analyze the global structure, every project structure and source code to generate comprehensive documentation.

## Instructions

1. **Scan the project** for all public APIs, services, and models
2. **Generate documentation** in Markdown format covering:
    - Project overview and architecture
    - Getting started guide (installation, configuration)
    - Entities and concepts, especially around integrity checking and metadata verification
    - Entities lifecycle (requests, whitelisted addresses, whitelisted contracts)
    - API reference for all public classes and methods
    - Usage examples with code snippets
    - Configuration options
3. **Output structure**:
    - `docs/README.md` - Main overview
    - `docs/getting-started.md` - Setup instructions
    - `docs/api/` - API reference per module
    - `docs/examples/` - Usage examples

## Focus on:

- Public interfaces and their contracts
- Method parameters, return types, and exceptions
- Required vs optional parameters
- Thread safety considerations
- Any deprecations or breaking changes

## Format:

- Use clear headings and code blocks
- Include Javadoc-style descriptions or equivalent
- Add cross-references between related classes

Generate documentation for: $ARGUMENTS
