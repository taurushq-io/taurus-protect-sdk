#!/usr/bin/env bash

# Generate TypeScript client from OpenAPI specification.
# Uses common API definitions from the root scripts/resources folder.
# Uses typescript-fetch generator for modern TypeScript.

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

on_exit() {
    echo -e "${RED}[ERROR]${NC} generate-openapi.sh has exited in error"
}

trap on_exit ERR

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(dirname "$PROJECT_ROOT")"

# Common API definitions in root folder
RESOURCES_DIR="$REPO_ROOT/scripts/resources"
GENERATOR_JAR="$RESOURCES_DIR/jars/openapi-generator-cli-7.9.0.jar"
SPEC_FILE="$RESOURCES_DIR/swagger/apis.swagger.json"
OUTPUT_DIR="$PROJECT_ROOT/src/internal/openapi"
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
    info "Found OpenAPI Generator JAR: $GENERATOR_JAR"

    if [[ ! -f "$SPEC_FILE" ]]; then
        error "OpenAPI spec not found at: $SPEC_FILE"
    fi
    info "Found OpenAPI spec: $SPEC_FILE"
}

generate() {
    info "Generating TypeScript OpenAPI client..."

    # Clean previous generation
    rm -rf "$TEMP_DIR"
    mkdir -p "$TEMP_DIR"

    # Generate TypeScript client using typescript-fetch generator
    "$JAVA_CMD" -jar "$GENERATOR_JAR" generate \
        -g typescript-fetch \
        -i "$SPEC_FILE" \
        -o "$TEMP_DIR" \
        --skip-validate-spec \
        --additional-properties=typescriptThreePlus=true \
        --additional-properties=supportsES6=true \
        --additional-properties=withInterfaces=true \
        --additional-properties=npmName=@taurushq/protect-openapi \
        --additional-properties=npmVersion=0.0.1 \
        --additional-properties=disallowAdditionalPropertiesIfNotPresent=false \
        --global-property=apiTests=false \
        --global-property=modelTests=false \
        --global-property=apiDocs=false \
        --global-property=modelDocs=false

    info "Copying generated files to $OUTPUT_DIR"

    # Clean and copy
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"

    # The generator outputs to $TEMP_DIR/src/ for typescript-fetch
    local SRC_DIR="$TEMP_DIR/src"

    # Copy the generated TypeScript files
    if [[ -d "$SRC_DIR" ]]; then
        # Copy all TypeScript files from src
        cp "$SRC_DIR/"*.ts "$OUTPUT_DIR/" 2>/dev/null || true

        # Copy models and apis directories
        if [[ -d "$SRC_DIR/models" ]]; then
            cp -r "$SRC_DIR/models" "$OUTPUT_DIR/"
        fi
        if [[ -d "$SRC_DIR/apis" ]]; then
            cp -r "$SRC_DIR/apis" "$OUTPUT_DIR/"
        fi
    elif [[ -d "$TEMP_DIR" ]]; then
        # Fallback: try direct copy from TEMP_DIR
        cp -r "$TEMP_DIR/"*.ts "$OUTPUT_DIR/" 2>/dev/null || true
        if [[ -d "$TEMP_DIR/models" ]]; then
            cp -r "$TEMP_DIR/models" "$OUTPUT_DIR/"
        fi
        if [[ -d "$TEMP_DIR/apis" ]]; then
            cp -r "$TEMP_DIR/apis" "$OUTPUT_DIR/"
        fi
    else
        error "Could not find generated files in $TEMP_DIR"
    fi

    # Clean up
    rm -rf "$TEMP_DIR"

    info "OpenAPI generation completed successfully"
    info "Files are in: $OUTPUT_DIR"
}

main() {
    info "Starting OpenAPI code generation for TypeScript SDK"
    check_java
    check_files
    generate
    info "Done!"
}

main "$@"
