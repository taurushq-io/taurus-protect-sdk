#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

MIN_PYTHON_VERSION="3.9"
VENV_DIR=".venv"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

check_python_version() {
    if ! command -v python3 &> /dev/null; then
        error "Python 3 is not installed"
    fi

    local version
    version=$(python3 -c 'import sys; print(f"{sys.version_info.major}.{sys.version_info.minor}")')
    info "Using Python $version"
}

setup_venv() {
    if [[ ! -d "$VENV_DIR" ]]; then
        info "Creating virtual environment in $VENV_DIR..."
        python3 -m venv "$VENV_DIR"
    fi
    # shellcheck disable=SC1091
    source "$VENV_DIR/bin/activate"
    info "Using virtual environment: $VENV_DIR"
}

# Check if pip version supports PEP 660 editable installs (requires pip >= 21.3)
check_pip_version() {
    local pip_version
    pip_version=$(pip --version | awk '{print $2}')
    local major minor
    major=$(echo "$pip_version" | cut -d. -f1)
    minor=$(echo "$pip_version" | cut -d. -f2)

    # pip >= 21.3 required for PEP 660 (editable installs with pyproject.toml only)
    if [[ "$major" -lt 21 ]] || { [[ "$major" -eq 21 ]] && [[ "$minor" -lt 3 ]]; }; then
        warn "pip $pip_version is too old for editable installs with pyproject.toml"
        warn "Upgrade pip with: pip install --upgrade pip"
        return 1
    fi
    return 0
}

install_deps() {
    setup_venv
    info "Installing dependencies..."
    if check_pip_version; then
        pip install -e ".[dev]" -q
        info "Dependencies installed"
    else
        warn "Falling back to non-editable install"
        pip install ".[dev]" -q
        info "Dependencies installed (non-editable mode)"
    fi
}

build() {
    info "Building package..."
    python3 -m build
    info "Build successful"
}

run_unit_tests() {
    info "Running unit tests..."
    pytest tests/unit/
    info "Unit tests completed"
}

run_integration_tests() {
    info "Running integration tests..."
    pytest tests/integration/
    info "Integration tests completed"
}

# Run a single unit test by pattern
run_single_unit_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 unit-one <pattern>\n\nExamples:\n  $0 unit-one test_approve_request                           # Test function name\n  $0 unit-one tests/unit/services/test_request.py::test_x    # Full path\n  $0 unit-one 'test_*wallet*'                                # Wildcard pattern\n  $0 unit-one TestRequestService                             # Test class name"
    fi
    info "Running single unit test: $pattern"
    if [[ "$pattern" == *"::"* ]] || [[ "$pattern" == *"/"* ]]; then
        pytest -v "$pattern"
    else
        pytest -v -k "$pattern" tests/unit/
    fi
    info "Test completed"
}

# Run a single integration test by pattern
run_single_integration_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 integration-one <pattern>\n\nExamples:\n  $0 integration-one test_list_wallets                              # Test function name\n  $0 integration-one tests/integration/test_wallet.py::test_list    # Full path\n  $0 integration-one 'test_*address*'                               # Wildcard pattern"
    fi
    info "Running single integration test: $pattern"
    if [[ "$pattern" == *"::"* ]] || [[ "$pattern" == *"/"* ]]; then
        pytest -v "$pattern"
    else
        pytest -v -k "$pattern" tests/integration/
    fi
    info "Test completed"
}

# Run E2E tests
run_e2e_tests() {
    info "Running E2E tests..."
    pytest tests/e2e/
    info "E2E tests completed"
}

# Run a single E2E test by pattern
run_single_e2e_test() {
    local pattern="$1"
    if [[ -z "$pattern" ]]; then
        error "Usage: $0 e2e-one <pattern>\n\nExamples:\n  $0 e2e-one test_multi_currency_transfer_e2e                   # Test function name\n  $0 e2e-one tests/e2e/test_multi_currency_e2e.py::test_x       # Full path\n  $0 e2e-one 'test_*e2e*'                                       # Wildcard pattern"
    fi
    info "Running single E2E test: $pattern"
    if [[ "$pattern" == *"::"* ]] || [[ "$pattern" == *"/"* ]]; then
        pytest -v "$pattern"
    else
        pytest -v -k "$pattern" tests/e2e/
    fi
    info "Test completed"
}

lint() {
    info "Running linters..."
    black --check taurus_protect tests
    isort --check-only taurus_protect tests
    flake8 taurus_protect tests --max-line-length=100 --exclude=taurus_protect/_internal
    mypy taurus_protect --exclude='taurus_protect/_internal'
    info "Lint completed"
}

format_code() {
    info "Formatting code..."
    black taurus_protect tests
    isort taurus_protect tests
    info "Format completed"
}

generate() {
    info "Running code generation..."
    if [[ -x "./scripts/generate-openapi.sh" ]]; then
        ./scripts/generate-openapi.sh
    else
        warn "OpenAPI generation script not found or not executable"
    fi
    if [[ -x "./scripts/generate-proto.sh" ]]; then
        ./scripts/generate-proto.sh
    else
        warn "Protobuf generation script not found or not executable"
    fi
    info "Code generation completed"
}

clean() {
    info "Cleaning..."
    rm -rf build/ dist/ *.egg-info .pytest_cache .mypy_cache .coverage htmlcov/
    find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
    find . -type f -name "*.pyc" -delete 2>/dev/null || true
    info "Clean completed"
}

usage() {
    cat << EOF
Usage: $0 [command] [args]

Commands:
    build              Build the package
    test               Run unit tests (alias for unit)
    unit               Run unit tests only
    unit-one <pattern> Run a single unit test matching pattern
    integration        Run integration tests (requires API access)
    integration-one <pattern>  Run a single integration test matching pattern
    e2e                Run E2E tests (requires API access)
    e2e-one <pattern>  Run a single E2E test matching pattern
    lint               Run linters (black, isort, flake8, mypy)
    format             Format code with black and isort
    generate           Run OpenAPI and protobuf code generation
    clean              Clean build artifacts
    install            Install package in development mode
    all                Install dependencies and run unit tests (default)
    help               Show this help message

Environment Variables (for integration tests):
    PROTECT_INTEGRATION_TEST  Set to "true" to enable integration tests
    PROTECT_API_HOST          API host URL
    PROTECT_API_KEY           API key
    PROTECT_API_SECRET        API secret (hex-encoded)

Examples:
    $0                                # Run install and unit tests
    $0 build                          # Build only
    $0 test                           # Run unit tests
    $0 unit-one test_approve          # Run single unit test by name
    $0 unit-one TestRequestService    # Run tests in class
    $0 integration                    # Run integration tests
    $0 integration-one test_list_wallets  # Run single integration test
    $0 e2e                                # Run E2E tests
    $0 e2e-one test_multi_currency        # Run single E2E test
    $0 lint                           # Run linters
    $0 format                         # Format code
EOF
}

main() {
    check_python_version

    local command="${1:-all}"

    case "$command" in
        build)
            setup_venv
            build
            ;;
        test|unit)
            install_deps
            run_unit_tests
            ;;
        unit-one)
            setup_venv
            run_single_unit_test "${2:-}"
            ;;
        integration)
            setup_venv
            run_integration_tests
            ;;
        integration-one)
            setup_venv
            run_single_integration_test "${2:-}"
            ;;
        e2e)
            setup_venv
            run_e2e_tests
            ;;
        e2e-one)
            setup_venv
            run_single_e2e_test "${2:-}"
            ;;
        lint)
            setup_venv
            lint
            ;;
        format)
            setup_venv
            format_code
            ;;
        generate)
            generate
            ;;
        clean)
            clean
            ;;
        install)
            install_deps
            ;;
        all)
            install_deps
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
