#!/usr/bin/env bash

# Generate TypeScript code from Protocol Buffer definitions.
# Uses common API definitions from the root scripts/resources folder.
# Uses ts-proto for TypeScript protobuf generation.
# Requires: protoc, ts-proto (installed via npm)

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

on_exit() {
    echo -e "${RED}[ERROR]${NC} generate-proto.sh has exited in error"
}

trap on_exit ERR

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(dirname "$PROJECT_ROOT")"

# Common API definitions in root folder
PROTO_BASE="$REPO_ROOT/scripts/resources/proto/schema"
OUTPUT_DIR="$PROJECT_ROOT/src/internal/proto"

check_protoc() {
    if ! command -v protoc &> /dev/null; then
        error "protoc is not installed. Install it from: https://github.com/protocolbuffers/protobuf/releases"
    fi
    info "protoc: $(which protoc)"
    info "$(protoc --version)"
}

check_ts_proto() {
    # Check if ts-proto plugin is available
    local ts_proto_plugin="$PROJECT_ROOT/node_modules/.bin/protoc-gen-ts_proto"

    if [[ ! -x "$ts_proto_plugin" ]]; then
        # Try global installation
        if command -v protoc-gen-ts_proto &> /dev/null; then
            ts_proto_plugin="$(which protoc-gen-ts_proto)"
            info "Using global ts-proto plugin: $ts_proto_plugin"
        else
            error "ts-proto is not installed. Run 'npm install' in $PROJECT_ROOT to install dependencies."
        fi
    else
        info "Using local ts-proto plugin: $ts_proto_plugin"
    fi

    # Export for use in generate function
    export TS_PROTO_PLUGIN="$ts_proto_plugin"
}

check_proto_dir() {
    if [[ ! -d "$PROTO_BASE" ]]; then
        error "Proto directory not found at: $PROTO_BASE"
    fi
    info "Found proto directory: $PROTO_BASE"
}

generate() {
    info "Generating TypeScript protobuf files..."

    # Clean and create output directory
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"

    # Find all proto files (excluding third_party/google/protobuf)
    local proto_files
    proto_files=$(find "$PROTO_BASE" -type f -name '*.proto' | grep -v 'third_party/google/protobuf/' || true)

    if [[ -z "$proto_files" ]]; then
        warn "No proto files found in $PROTO_BASE"
        return
    fi

    # Generate TypeScript files using ts-proto
    for file in $proto_files; do
        if [[ -f "$file" ]]; then
            info "Generating TypeScript from: $(basename "$file")"
            protoc \
                --proto_path="$PROTO_BASE/third_party" \
                --proto_path="$PROTO_BASE/common" \
                --proto_path="$PROTO_BASE/v1" \
                --plugin="protoc-gen-ts_proto=$TS_PROTO_PLUGIN" \
                --ts_proto_out="$OUTPUT_DIR" \
                --ts_proto_opt=esModuleInterop=true \
                --ts_proto_opt=outputJsonMethods=true \
                --ts_proto_opt=outputEncodeMethods=true \
                --ts_proto_opt=outputPartialMethods=true \
                --ts_proto_opt=useExactTypes=false \
                --ts_proto_opt=env=browser \
                "$file" 2>/dev/null || {
                    warn "Skipped: $(basename "$file") (may have unsupported imports)"
                }
        fi
    done

    # Flatten subdirectories - move all .ts files to the root proto directory
    info "Flattening directory structure..."
    find "$OUTPUT_DIR" -mindepth 2 -name "*.ts" -exec mv {} "$OUTPUT_DIR/" \; 2>/dev/null || true
    find "$OUTPUT_DIR" -mindepth 1 -type d -empty -delete 2>/dev/null || true

    # Create index.ts to export all generated types
    info "Creating index.ts..."
    local index_file="$OUTPUT_DIR/index.ts"
    echo '// Auto-generated protobuf types. DO NOT MODIFY.' > "$index_file"
    echo '' >> "$index_file"

    for ts_file in "$OUTPUT_DIR"/*.ts; do
        if [[ -f "$ts_file" && "$(basename "$ts_file")" != "index.ts" ]]; then
            local module_name
            module_name=$(basename "$ts_file" .ts)
            echo "export * from './${module_name}';" >> "$index_file"
        fi
    done

    # Count generated files
    local count
    count=$(find "$OUTPUT_DIR" -name "*.ts" -not -name "index.ts" | wc -l | tr -d ' ')

    info "Generated $count TypeScript protobuf files"
    info "Protobuf code generated successfully"
    info "Files are in: $OUTPUT_DIR"
}

main() {
    info "Starting protobuf code generation for TypeScript SDK"
    check_protoc
    check_ts_proto
    check_proto_dir
    generate
    info "Done!"
}

main "$@"
