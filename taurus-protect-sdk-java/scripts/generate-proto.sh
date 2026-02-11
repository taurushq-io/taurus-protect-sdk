#!/usr/bin/env bash

# Main procedure.
set -e

on_exit () {
    echo "generate-proto.sh has exited in error"
}

trap on_exit ERR

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$DIR"

# Common API definitions in root folder
PROTO_BASE="$DIR/../scripts/resources/proto/schema"

if [[ ! -d "$PROTO_BASE" ]]; then
    echo "Proto directory not found at: $PROTO_BASE"
    exit 1
fi

# patch proto source as a workaround for protoc generating code that does not compile
if grep -q -F "VersionKind version = 1;" "$PROTO_BASE/v1/request_reply.proto"; then
  patch "$PROTO_BASE/v1/request_reply.proto" < scripts/proto-tpv1.patch
fi

echo
echo "==> protoc (`which protoc`):"
echo "    `protoc --version`"
echo


for file in  `find "$PROTO_BASE" -type f | grep '\.proto' | grep -v 'third_party/google/protobuf/'`
do
    if [[ -f $file ]]; then
        echo "---> generating java files from '$file'"
        protoc \
            --proto_path="$PROTO_BASE/third_party" \
            --proto_path="$PROTO_BASE/common" \
            --proto_path="$PROTO_BASE/v1" \
            --java_out=proto/src/main/java \
            $file

        echo
    fi
done






