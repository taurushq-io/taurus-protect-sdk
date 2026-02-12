#!/usr/bin/env bash

# Generate Python code from Protocol Buffer definitions.
# Uses common API definitions from the root scripts/resources folder.
# Requires: protoc

set -e

on_exit() {
    echo "generate-proto.sh has exited in error"
}

trap on_exit ERR

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Common API definitions in root folder
PROTO_BASE="$DIR/../scripts/resources/proto/schema"

if [[ ! -d "$PROTO_BASE" ]]; then
    error "Proto directory not found at: $PROTO_BASE"
fi

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    error "protoc is not installed. Install it from: https://github.com/protocolbuffers/protobuf/releases"
fi

echo
echo "==> protoc ($(which protoc)):"
echo "    $(protoc --version)"
echo

# Clean and create output directory
rm -rf taurus_protect/_internal/proto
mkdir -p taurus_protect/_internal/proto

info "Generating Python protobuf files..."

# Generate Python files from all proto files (excluding third_party)
for file in $(find "$PROTO_BASE" -type f -name '*.proto' | grep -v 'third_party/google/protobuf/'); do
    if [[ -f "$file" ]]; then
        echo "---> generating python files from '$file'"
        protoc \
            --proto_path="$PROTO_BASE/third_party" \
            --proto_path="$PROTO_BASE/common" \
            --proto_path="$PROTO_BASE/v1" \
            --python_out=taurus_protect/_internal/proto \
            "$file" 2>/dev/null || true
    fi
done

# Create __init__.py
echo '"""Auto-generated protobuf classes. DO NOT MODIFY."""' > taurus_protect/_internal/proto/__init__.py

# Flatten subdirectories - move all _pb2.py files to the root proto directory
echo
info "Flattening directory structure..."
find taurus_protect/_internal/proto -mindepth 2 -name "*_pb2.py" -exec mv {} taurus_protect/_internal/proto/ \; 2>/dev/null || true
find taurus_protect/_internal/proto -mindepth 1 -type d -empty -delete 2>/dev/null || true

# Fix import paths: change 'import X_pb2' to 'from . import X_pb2'
# This is necessary because all pb2 files are flattened into the same directory
echo
info "Fixing import paths for package-relative imports..."
for file in taurus_protect/_internal/proto/*_pb2.py; do
    if [[ -f "$file" ]]; then
        # macOS sed requires '' after -i, Linux doesn't
        if [[ "$(uname)" == "Darwin" ]]; then
            # Replace 'import X_pb2 as' with 'from . import X_pb2 as'
            sed -i '' 's/^import \([a-zA-Z0-9_]*_pb2\) as/from . import \1 as/g' "$file"
            # Replace 'import X_pb2$' (at end of line) with 'from . import X_pb2'
            sed -i '' 's/^import \([a-zA-Z0-9_]*_pb2\)$/from . import \1/g' "$file"
            # Replace 'from X import Y_pb2 as Z' with 'from . import Y_pb2 as Z'
            sed -i '' 's/^from [a-zA-Z0-9_]* import \([a-zA-Z0-9_]*_pb2\) as/from . import \1 as/g' "$file"
        else
            sed -i 's/^import \([a-zA-Z0-9_]*_pb2\) as/from . import \1 as/g' "$file"
            sed -i 's/^import \([a-zA-Z0-9_]*_pb2\)$/from . import \1/g' "$file"
            # Replace 'from X import Y_pb2 as Z' with 'from . import Y_pb2 as Z'
            sed -i 's/^from [a-zA-Z0-9_]* import \([a-zA-Z0-9_]*_pb2\) as/from . import \1 as/g' "$file"
        fi
    fi
done

# Count generated files
count=$(find taurus_protect/_internal/proto -name "*_pb2.py" | wc -l | tr -d ' ')

echo
info "Generated $count protobuf Python files"
info "Protobuf code generated successfully."
info "Files are in: taurus_protect/_internal/proto/"
