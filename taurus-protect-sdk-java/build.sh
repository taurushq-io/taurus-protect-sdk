#!/usr/bin/env bash
set -euo pipefail

# Script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Minimum Java version
MIN_JAVA_VERSION="8"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions for colored output
info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Check Java version
check_java_version() {
    if ! command -v java &> /dev/null; then
        error "Java is not installed"
    fi

    local version
    version=$(java -version 2>&1 | head -1 | cut -d'"' -f2 | cut -d'.' -f1)
    # Handle versions like "1.8" -> 8
    if [[ "$version" == "1" ]]; then
        version=$(java -version 2>&1 | head -1 | cut -d'"' -f2 | cut -d'.' -f2)
    fi

    if [[ "$version" -lt "$MIN_JAVA_VERSION" ]]; then
        error "Java $MIN_JAVA_VERSION or higher is required (found $version)"
    fi
    info "Using Java $version"
}

# Check Maven installation
check_maven() {
    if ! command -v mvn &> /dev/null; then
        error "Maven is not installed"
    fi
    local mvn_version
    mvn_version=$(mvn -version 2>&1 | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
    info "Using Maven $mvn_version"
}

# Build all modules
build() {
    info "Building all modules..."
    mvn clean compile -q
    info "Build successful"
}

# Run tests
run_tests() {
    info "Running tests..."
    mvn test -q
    info "Tests completed"
}

# Run integration tests
run_integration_tests() {
    info "Running integration tests..."
    export PROTECT_INTEGRATION_TEST=true
    mvn test -Dtest="*IntegrationTest" -pl client -q
    info "Integration tests completed"
}

# Run a single unit test by pattern
run_single_unit_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 unit-one <pattern>\n\nExamples:\n  $0 unit-one RequestServiceTest              # Run all tests in class\n  $0 unit-one RequestServiceTest#testApprove  # Run specific test method\n  $0 unit-one '*Mapper*'                       # Wildcard patterns"
    fi
    info "Running single unit test: $pattern"
    mvn test -Dtest="$pattern" -pl client
    info "Test completed"
}

# Run a single integration test by pattern
run_single_integration_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 integration-one <pattern>\n\nExamples:\n  $0 integration-one WalletIntegrationTest              # Run all tests in class\n  $0 integration-one WalletIntegrationTest#testList     # Run specific test method\n  $0 integration-one '*Address*'                        # Wildcard patterns"
    fi
    info "Running single integration test: $pattern"
    export PROTECT_INTEGRATION_TEST=true
    mvn test -Dtest="$pattern" -pl client
    info "Test completed"
}

# Run E2E tests
run_e2e_tests() {
    info "Running E2E tests..."
    export PROTECT_INTEGRATION_TEST=true
    mvn test -Dtest="*E2ETest" -pl client -q
    info "E2E tests completed"
}

# Run a single E2E test by pattern
run_single_e2e_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 e2e-one <pattern>\n\nExamples:\n  $0 e2e-one MultiCurrencyE2ETest              # Run all tests in class\n  $0 e2e-one MultiCurrencyE2ETest#multiCurrency  # Run specific test method\n  $0 e2e-one '*BusinessRule*'                    # Wildcard patterns"
    fi
    info "Running single E2E test: $pattern"
    export PROTECT_INTEGRATION_TEST=true
    mvn test -Dtest="$pattern" -pl client
    info "Test completed"
}

# Run full verification (compile + test + static analysis)
verify() {
    info "Running full verification (compile + test + static analysis)..."
    mvn clean verify
    info "Verification completed"
}

# Run static analysis
lint() {
    info "Running static analysis..."
    mvn spotbugs:check pmd:check checkstyle:check -q
    info "Static analysis completed"
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
    mvn clean -q
    info "Clean completed"
}

# Fast local install (skip checks)
install_fast() {
    info "Installing locally (skipping checks)..."
    mvn install -DskipTests -Dspotbugs.skip=true -Dpmd.skip=true \
        -Dcheckstyle.skip=true -Dmaven.javadoc.skip=true -q
    info "Install completed"
}

# Show usage
usage() {
    cat << EOF
Usage: $0 [command] [args]

Commands:
    build              Compile all modules
    test               Build and run unit tests (alias for unit)
    unit               Run unit tests only
    unit-one <pattern> Run a single unit test matching pattern
    integration        Run integration tests (requires API access)
    integration-one <pattern>  Run a single integration test matching pattern
    e2e                Run E2E tests (requires API access)
    e2e-one <pattern>  Run a single E2E test matching pattern
    verify             Full verification (compile + test + static analysis)
    lint               Run static analysis (SpotBugs, PMD, Checkstyle)
    generate           Run OpenAPI and protobuf code generation
    clean              Clean build artifacts
    install            Fast local install (skip checks)
    all                Build and unit test (default)
    help               Show this help message

Environment Variables (for integration tests):
    PROTECT_INTEGRATION_TEST  Set to "true" to enable integration tests
    PROTECT_API_HOST          API host URL
    PROTECT_API_KEY           API key
    PROTECT_API_SECRET        API secret (hex-encoded)

Examples:
    $0                                          # Run build and test
    $0 build                                    # Compile only
    $0 unit-one RequestServiceTest              # Run single test class
    $0 unit-one RequestServiceTest#testApprove  # Run single test method
    $0 integration                              # Run integration tests
    $0 integration-one WalletIntegrationTest    # Run single integration test
    $0 e2e                                          # Run E2E tests
    $0 e2e-one MultiCurrencyE2ETest                 # Run single E2E test
    $0 verify                                   # Full verification with static analysis
    $0 install                                  # Fast local install
EOF
}

# Main
main() {
    check_java_version
    check_maven

    local command="${1:-all}"

    case "$command" in
        build)
            build
            ;;
        test|unit)
            build
            run_tests
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
        verify)
            verify
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
        install)
            install_fast
            ;;
        all)
            build
            run_tests
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
