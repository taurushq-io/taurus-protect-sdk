#!/usr/bin/env bash
set -euo pipefail

# Script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Minimum Go version
MIN_GO_VERSION="1.24"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions for colored output
info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Check Go version
check_go_version() {
    if ! command -v go &> /dev/null; then
        error "Go is not installed"
    fi

    local version
    version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    if [[ "$(printf '%s\n' "$MIN_GO_VERSION" "$version" | sort -V | head -n1)" != "$MIN_GO_VERSION" ]]; then
        error "Go $MIN_GO_VERSION or higher is required (found $version)"
    fi
    info "Using Go $version"
}

# Build all packages
build() {
    info "Building all packages..."
    go build ./...
    info "Build successful"
}

# Run unit tests
run_unit_tests() {
    info "Running unit tests..."
    go test -cover ./pkg/protect/...
    info "Unit tests completed"
}

# Run integration tests
run_integration_tests() {
    info "Running integration tests..."
    go test -v -timeout 15m ./test/integration/...
    info "Integration tests completed"
}

# Run a single unit test by pattern
run_single_unit_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 unit-one <pattern>\n\nExamples:\n  $0 unit-one TestMapWallet           # Run tests matching pattern\n  $0 unit-one 'Test.*Request'         # Regex pattern\n  $0 unit-one TestConstantTime        # Specific test function"
    fi
    info "Running single unit test: $pattern"
    go test -v -run "$pattern" ./pkg/protect/...
    info "Test completed"
}

# Run a single integration test by pattern
run_single_integration_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 integration-one <pattern>\n\nExamples:\n  $0 integration-one TestListWallets    # Run tests matching pattern\n  $0 integration-one 'Test.*Address'    # Regex pattern\n  $0 integration-one TestHealthCheck    # Specific test function"
    fi
    info "Running single integration test: $pattern"
    go test -v -timeout 15m -run "$pattern" ./test/integration/...
    info "Test completed"
}

# Run E2E tests
run_e2e_tests() {
    info "Running E2E tests..."
    go test -v -timeout 30m ./test/e2e/...
    info "E2E tests completed"
}

# Run a single E2E test by pattern
run_single_e2e_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 e2e-one <pattern>\n\nExamples:\n  $0 e2e-one TestIntegration_MultiCurrencyE2E  # Run E2E test\n  $0 e2e-one 'Test.*E2E'                       # Regex pattern"
    fi
    info "Running single E2E test: $pattern"
    go test -v -timeout 30m -run "$pattern" ./test/e2e/...
    info "Test completed"
}

# Run linter
lint() {
    if command -v golangci-lint &> /dev/null; then
        info "Running linter..."
        golangci-lint run
        info "Lint completed"
    else
        warn "golangci-lint not installed, skipping lint"
    fi
}

# Run code generation
generate() {
    info "Running code generation..."
    if [[ -x "./scripts/generate-openapi.sh" ]]; then
        ./scripts/generate-openapi.sh
    fi
    if [[ -x "./scripts/generate-proto.sh" ]]; then
        ./scripts/generate-proto.sh
    fi
    info "Code generation completed"
}

# Clean build artifacts
clean() {
    info "Cleaning..."
    go clean ./...
    info "Clean completed"
}

# Show usage
usage() {
    cat << EOF
Usage: $0 [command] [args]

Commands:
    build              Build all packages
    test               Run unit tests with coverage
    unit               Run unit tests only (alias for test)
    unit-one <pattern> Run a single unit test matching pattern (regex)
    integration        Run integration tests (requires API access)
    integration-one <pattern>  Run a single integration test matching pattern (regex)
    e2e                Run E2E tests (requires API access)
    e2e-one <pattern>  Run a single E2E test matching pattern (regex)
    lint               Run golangci-lint (if installed)
    generate           Run OpenAPI and protobuf code generation
    clean              Clean build artifacts
    all                Build and test (default)
    help               Show this help message

Environment Variables (for integration tests):
    PROTECT_INTEGRATION_TEST  Set to "true" to enable integration tests
    PROTECT_API_HOST          API host URL
    PROTECT_API_KEY           API key
    PROTECT_API_SECRET        API secret (hex-encoded)

Examples:
    $0                              # Run build and unit tests
    $0 build                        # Build only
    $0 test                         # Run unit tests
    $0 unit-one TestMapWallet       # Run single unit test
    $0 unit-one 'Test.*Request'     # Run tests matching regex
    $0 integration                  # Run integration tests
    $0 integration-one TestListWallets  # Run single integration test
    $0 e2e                              # Run E2E tests
    $0 e2e-one TestIntegration_MultiCurrencyE2E  # Run single E2E test
EOF
}

# Main
main() {
    check_go_version

    local command="${1:-all}"

    case "$command" in
        build)
            build
            ;;
        test|unit)
            build
            run_unit_tests
            ;;
        unit-one)
            build
            run_single_unit_test "${2:-}"
            ;;
        integration)
            build
            run_integration_tests
            ;;
        integration-one)
            build
            run_single_integration_test "${2:-}"
            ;;
        e2e)
            build
            run_e2e_tests
            ;;
        e2e-one)
            build
            run_single_e2e_test "${2:-}"
            ;;
        lint)
            lint
            ;;
        generate)
            generate
            ;;
        clean)
            clean
            ;;
        all)
            build
            run_unit_tests
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            error "Unknown command: $command. Use '$0 help' for usage."
            ;;
    esac
}

main "$@"
