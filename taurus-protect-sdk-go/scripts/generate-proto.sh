#!/usr/bin/env bash

# Generate Go code from Protocol Buffer definitions.
# Uses common API definitions from the root scripts/resources folder.
# Requires: protoc, protoc-gen-go (go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)

set -e

on_exit() {
    echo "generate-proto.sh has exited in error"
}

trap on_exit ERR

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed. Install it from:"
    echo "  https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
# Note: Unlike Java, Go protobuf generation requires a separate plugin.
# This is the official Google-maintained plugin.
if ! command -v protoc-gen-go &> /dev/null; then
    echo "protoc-gen-go is not installed. Install it with:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Common API definitions in root folder
PROTO_BASE="$DIR/../scripts/resources/proto/schema"

if [[ ! -d "$PROTO_BASE" ]]; then
    echo "Proto directory not found at: $PROTO_BASE"
    exit 1
fi

echo
echo "==> protoc ($(which protoc)):"
echo "    $(protoc --version)"
echo

# Clean and create output directory
rm -rf internal/proto
mkdir -p internal/proto

# Go package for generated code
GO_PACKAGE="github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"

# Build M mappings for all proto files (required since proto files don't have go_package option)
M_OPTS=""

# Add mappings for v1 proto files
for file in $(find "$PROTO_BASE/v1" -maxdepth 1 -name '*.proto' -type f 2>/dev/null); do
    basename=$(basename "$file")
    M_OPTS="$M_OPTS --go_opt=M${basename}=${GO_PACKAGE}"
done

# Add mappings for common proto files (flatten to same package)
for file in $(find "$PROTO_BASE/common" -name '*.proto' -type f 2>/dev/null); do
    # Get relative path from common directory
    relpath=$(echo "$file" | sed "s|$PROTO_BASE/common/||")
    M_OPTS="$M_OPTS --go_opt=M${relpath}=${GO_PACKAGE}"
done

# Add mappings for Google well-known types (from google.golang.org/protobuf)
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/duration.proto=google.golang.org/protobuf/types/known/durationpb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/struct.proto=google.golang.org/protobuf/types/known/structpb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb"
M_OPTS="$M_OPTS --go_opt=Mgoogle/protobuf/descriptor.proto=google.golang.org/protobuf/types/descriptorpb"

# Add mappings for Google API protos (from google.golang.org/genproto)
M_OPTS="$M_OPTS --go_opt=Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations"
M_OPTS="$M_OPTS --go_opt=Mgoogle/api/http.proto=google.golang.org/genproto/googleapis/api/annotations"
M_OPTS="$M_OPTS --go_opt=Mgoogle/api/httpbody.proto=google.golang.org/genproto/googleapis/api/httpbody"
M_OPTS="$M_OPTS --go_opt=Mgoogle/api/visibility.proto=google.golang.org/genproto/googleapis/api/visibility"
M_OPTS="$M_OPTS --go_opt=Mgoogle/rpc/status.proto=google.golang.org/genproto/googleapis/rpc/status"
M_OPTS="$M_OPTS --go_opt=Mgoogle/rpc/code.proto=google.golang.org/genproto/googleapis/rpc/code"
M_OPTS="$M_OPTS --go_opt=Mgoogle/rpc/error_details.proto=google.golang.org/genproto/googleapis/rpc/errdetails"

# Add mappings for third-party protos
M_OPTS="$M_OPTS --go_opt=Mgithub.com/mwitkow/go-proto-validators/validator.proto=github.com/mwitkow/go-proto-validators"
M_OPTS="$M_OPTS --go_opt=Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
M_OPTS="$M_OPTS --go_opt=Mprotoc-gen-openapiv2/options/openapiv2.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"

echo "Generating Go code with package mappings..."
echo

# Generate all proto files in a single protoc invocation
ALL_PROTOS=""
for file in $(find "$PROTO_BASE" -type f -name '*.proto' | grep -v 'third_party/'); do
    ALL_PROTOS="$ALL_PROTOS $file"
done

echo "---> generating go files from all proto files"
protoc \
    --proto_path="$PROTO_BASE/third_party" \
    --proto_path="$PROTO_BASE/common" \
    --proto_path="$PROTO_BASE/v1" \
    --go_out=internal/proto \
    --go_opt=paths=source_relative \
    $M_OPTS \
    $ALL_PROTOS

echo

# Flatten subdirectories - move all .pb.go files to the root proto directory
echo "Flattening directory structure..."
find internal/proto -mindepth 2 -name "*.pb.go" -exec mv {} internal/proto/ \;
find internal/proto -mindepth 1 -type d -empty -delete

echo
echo "Protobuf code generated successfully."
echo "Files are in: internal/proto/"
