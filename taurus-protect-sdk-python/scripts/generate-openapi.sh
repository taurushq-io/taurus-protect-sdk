#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(dirname "$PROJECT_ROOT")"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Configuration
GENERATOR_JAR="$REPO_ROOT/scripts/resources/jars/openapi-generator-cli-7.9.0.jar"
OPENAPI_SPEC="$REPO_ROOT/scripts/resources/swagger/apis.swagger.json"
OUTPUT_DIR="$PROJECT_ROOT/taurus_protect/_internal/openapi"
TEMP_DIR="$PROJECT_ROOT/.codegen"

# Java configuration - use Java 11+ if available
# Check for Java 22 in user library first (common location for homebrew/manual installs)
if [[ -x "/Users/admin/Library/Java/JavaVirtualMachines/openjdk-22.0.2/Contents/Home/bin/java" ]]; then
    JAVA_CMD="/Users/admin/Library/Java/JavaVirtualMachines/openjdk-22.0.2/Contents/Home/bin/java"
elif [[ -n "${JAVA_HOME:-}" && -x "${JAVA_HOME}/bin/java" ]]; then
    JAVA_CMD="${JAVA_HOME}/bin/java"
else
    JAVA_CMD="java"  # Fall back to system java
fi

# Check prerequisites
check_java() {
    if ! command -v "$JAVA_CMD" &> /dev/null && [[ ! -x "$JAVA_CMD" ]]; then
        error "Java is required but not installed. Please install Java 11+."
    fi

    local java_version
    java_version=$("$JAVA_CMD" -version 2>&1 | head -1 | cut -d'"' -f2)
    local major_version
    # Handle both "1.8.x" and "11.x" formats
    if [[ "$java_version" == 1.* ]]; then
        major_version=$(echo "$java_version" | cut -d'.' -f2)
    else
        major_version=$(echo "$java_version" | cut -d'.' -f1)
    fi
    if [[ "$major_version" -lt 11 ]]; then
        error "Java 11+ is required. Found version: $java_version"
    fi
    info "Using Java version: $("$JAVA_CMD" -version 2>&1 | head -1)"
}

check_files() {
    if [[ ! -f "$GENERATOR_JAR" ]]; then
        error "OpenAPI Generator JAR not found at: $GENERATOR_JAR"
    fi

    if [[ ! -f "$OPENAPI_SPEC" ]]; then
        error "OpenAPI spec not found at: $OPENAPI_SPEC"
    fi
}

generate() {
    info "Generating Python OpenAPI client..."

    # Clean previous generation
    rm -rf "$TEMP_DIR"
    mkdir -p "$TEMP_DIR"

    # Generate Python client
    "$JAVA_CMD" -jar "$GENERATOR_JAR" generate \
        -g python \
        -i "$OPENAPI_SPEC" \
        -o "$TEMP_DIR" \
        --skip-validate-spec \
        --additional-properties=packageName=taurus_protect._internal.openapi \
        --additional-properties=projectName=taurus-protect-openapi \
        --additional-properties=pythonAttrNoneIfUnset=true \
        --additional-properties=disallowAdditionalPropertiesIfNotPresent=false \
        --global-property=apiTests=false \
        --global-property=modelTests=false \
        --global-property=apiDocs=false \
        --global-property=modelDocs=false

    info "Copying generated files to $OUTPUT_DIR"

    # Clean and copy
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"

    # Copy the generated package
    if [[ -d "$TEMP_DIR/taurus_protect/_internal/openapi" ]]; then
        cp -r "$TEMP_DIR/taurus_protect/_internal/openapi/"* "$OUTPUT_DIR/"
    elif [[ -d "$TEMP_DIR/openapi" ]]; then
        cp -r "$TEMP_DIR/openapi/"* "$OUTPUT_DIR/"
    else
        # Find the generated package directory
        local pkg_dir
        pkg_dir=$(find "$TEMP_DIR" -type d -name "openapi" | head -1)
        if [[ -n "$pkg_dir" ]]; then
            cp -r "$pkg_dir/"* "$OUTPUT_DIR/"
        else
            error "Could not find generated openapi package in $TEMP_DIR"
        fi
    fi

    # Create __init__.py if not exists
    if [[ ! -f "$OUTPUT_DIR/__init__.py" ]]; then
        echo '"""Auto-generated OpenAPI client. DO NOT MODIFY."""' > "$OUTPUT_DIR/__init__.py"
    fi

    # Fix Pydantic v2 compatibility issues
    # The generated code uses strict=True on Union types with None which doesn't work
    info "Fixing Pydantic v2 compatibility issues..."
    find "$OUTPUT_DIR/models" -name "*.py" -exec sed -i '' \
        's/Optional\[Union\[Annotated\[bytes, Field(strict=True)\], Annotated\[str, Field(strict=True)\]\]\]/Optional[Union[bytes, str]]/g' {} \;

    # Clean up
    rm -rf "$TEMP_DIR"

    info "OpenAPI generation completed successfully"
}

main() {
    info "Starting OpenAPI code generation for Python SDK"
    check_java
    check_files
    generate
    info "Done!"
}

main "$@"
