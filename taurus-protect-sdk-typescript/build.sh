#!/bin/bash
# Build script for Taurus-PROTECT TypeScript SDK
# Usage: ./build.sh [command]
#   (default) - Build and test
#   build     - Compile TypeScript only
#   test      - Run tests only
#   lint      - Run linter
#   generate  - Generate OpenAPI and Protobuf code
#   clean     - Clean build artifacts

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check Node.js version
check_node() {
    if ! command -v node &> /dev/null; then
        error "Node.js is not installed"
        exit 1
    fi

    NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_VERSION" -lt 18 ]; then
        error "Node.js 18+ is required (found: $(node -v))"
        exit 1
    fi

    info "Using Node.js $(node -v)"
}

# Check npm
check_npm() {
    if ! command -v npm &> /dev/null; then
        error "npm is not installed"
        exit 1
    fi
    info "Using npm $(npm -v)"
}

# Install dependencies if needed
install_deps() {
    if [ ! -d "node_modules" ]; then
        info "Installing dependencies..."
        npm install
    fi
}

# Build TypeScript
do_build() {
    info "Building TypeScript..."
    npm run build
    info "Build successful"
}

# Run unit tests
do_test() {
    info "Running unit tests..."
    npm test -- tests/unit/
    info "Unit tests completed"
}

# Run integration tests
do_integration_test() {
    info "Running integration tests..."
    npm test -- tests/integration/
    info "Integration tests completed"
}

# Run a single unit test by pattern
do_single_unit_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 unit-one <pattern>"
        cat << EOF

Examples:
  $0 unit-one "should approve request"    # Match test description
  $0 unit-one "RequestService"            # Match describe block
  $0 unit-one "hash verification"         # Match any test name
EOF
        exit 1
    fi
    info "Running single unit test: $pattern"
    npm test -- tests/unit/ --testNamePattern="$pattern" --verbose
    info "Test completed"
}

# Run a single integration test by pattern
do_single_integration_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 integration-one <pattern>"
        cat << EOF

Examples:
  $0 integration-one "should list wallets"    # Match test description
  $0 integration-one "WalletService"          # Match describe block
  $0 integration-one "pagination"             # Match any test name
EOF
        exit 1
    fi
    info "Running single integration test: $pattern"
    npm test -- tests/integration/ --testNamePattern="$pattern" --verbose
    info "Test completed"
}

# Run E2E tests
do_e2e_test() {
    info "Running E2E tests..."
    npm test -- tests/e2e/
    info "E2E tests completed"
}

# Run a single E2E test by pattern
do_single_e2e_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 e2e-one <pattern>"
        cat << EOF

Examples:
  $0 e2e-one "should complete transfer"     # Match test description
  $0 e2e-one "Multi-Currency"               # Match describe block
EOF
        exit 1
    fi
    info "Running single E2E test: $pattern"
    npm test -- tests/e2e/ --testNamePattern="$pattern" --verbose
    info "Test completed"
}

# Run linter
do_lint() {
    info "Running linter..."
    npm run lint
    info "Lint completed"
}

# Generate code
do_generate() {
    info "Generating OpenAPI client..."
    ./scripts/generate-openapi.sh

    info "Generating Protobuf types..."
    ./scripts/generate-proto.sh

    info "Code generation completed"
}

# Clean build artifacts
do_clean() {
    info "Cleaning build artifacts..."
    rm -rf dist coverage node_modules
    info "Clean completed"
}

# Main
check_node
check_npm

COMMAND="${1:-}"

usage() {
    cat << EOF
Usage: $0 [command] [args]

Commands:
    build              Compile TypeScript only
    test               Run unit tests (alias for unit)
    unit               Run unit tests only
    unit-one <pattern> Run a single unit test matching pattern
    integration        Run integration tests (requires API access)
    integration-one <pattern>  Run a single integration test matching pattern
    e2e                Run E2E tests (requires API access)
    e2e-one <pattern>  Run a single E2E test matching pattern
    lint               Run linter
    generate           Generate OpenAPI and Protobuf code
    clean              Clean build artifacts
    (default)          Build and run unit tests
    help               Show this help message

Environment Variables (for integration tests):
    PROTECT_INTEGRATION_TEST  Set to "true" to enable integration tests
    PROTECT_API_HOST          API host URL
    PROTECT_API_KEY           API key
    PROTECT_API_SECRET        API secret (hex-encoded)

Examples:
    $0                                  # Run build and unit tests
    $0 build                            # Compile TypeScript only
    $0 test                             # Run unit tests
    $0 unit-one "should approve"        # Run single unit test by description
    $0 unit-one "RequestService"        # Run tests in describe block
    $0 integration                      # Run integration tests
    $0 integration-one "should list"    # Run single integration test
    $0 e2e                                  # Run E2E tests
    $0 e2e-one "Multi-Currency"             # Run single E2E test
EOF
}

PATTERN="${2:-}"

case "$COMMAND" in
    "build")
        install_deps
        do_build
        ;;
    "test"|"unit")
        install_deps
        do_test
        ;;
    "unit-one")
        install_deps
        do_single_unit_test "$PATTERN"
        ;;
    "integration")
        install_deps
        do_integration_test
        ;;
    "integration-one")
        install_deps
        do_single_integration_test "$PATTERN"
        ;;
    "e2e")
        install_deps
        do_e2e_test
        ;;
    "e2e-one")
        install_deps
        do_single_e2e_test "$PATTERN"
        ;;
    "lint")
        install_deps
        do_lint
        ;;
    "generate")
        install_deps
        do_generate
        ;;
    "clean")
        do_clean
        ;;
    "help"|"--help"|"-h")
        usage
        ;;
    "")
        # Default: build and test
        install_deps
        do_build
        do_test
        ;;
    *)
        error "Unknown command: $COMMAND"
        usage
        exit 1
        ;;
esac
