#!/usr/bin/env bash

# Generate Go client from OpenAPI specification.
# Uses common API definitions from the root scripts/resources folder.

set -e

on_exit() {
    echo "generate-openapi.sh has exited in error"
}

trap on_exit ERR

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

# Check Java version (openapi-generator-cli requires Java 11+)
JAVA_VER=$(javap -verbose java.lang.String | grep "major version" | cut -d " " -f5)

if [[ $JAVA_VER -lt 55 ]]; then
    echo "In order to use the Java openapi-generator-cli, you need Java 11 runtime at minimum (class file version 55.0)."
    exit 1
fi

# Common API definitions in root folder
RESOURCES_DIR="$DIR/../scripts/resources"
GENERATOR_JAR="$RESOURCES_DIR/jars/openapi-generator-cli-7.9.0.jar"
SPEC_FILE="$RESOURCES_DIR/swagger/apis.swagger.json"

if [[ ! -f "$GENERATOR_JAR" ]]; then
    echo "OpenAPI generator JAR not found at: $GENERATOR_JAR"
    exit 1
fi

if [[ ! -f "$SPEC_FILE" ]]; then
    echo "OpenAPI spec not found at: $SPEC_FILE"
    exit 1
fi

echo "Generating Go client from: $SPEC_FILE"

# Clean previous generation
rm -rf .codegen internal/openapi
mkdir -p .codegen internal/openapi

# Vendored template overrides. Only the files present here override the generator's built-ins;
# everything else still comes from the JAR, so this directory stays deliberately tiny.
#
# It exists because `decode` must reject a 2xx whose typed reply could not be populated. An empty
# body, or the JSON literal `null`, otherwise leaves the reply pointer nil and returns no error —
# and every hand-written service dereferences the reply immediately after checking err (~110 sites
# across 38 files), so a server could turn any routine read into a process-killing nil-pointer
# panic. Patching internal/openapi/client.go by hand is not enough: this script `rm -rf`s that
# directory above, so the fix has to live in the template to survive regeneration.
TEMPLATE_DIR="$RESOURCES_DIR/templates/go"

# Generate Go client with enumClassPrefix to avoid const conflicts
java -jar "$GENERATOR_JAR" generate -g go -i "$SPEC_FILE" -o .codegen \
    -t "$TEMPLATE_DIR" \
    --skip-validate-spec \
    --additional-properties=packageName=openapi \
    --additional-properties=isGoSubmodule=true \
    --additional-properties=enumClassPrefix=true

# Copy generated files to internal/openapi
cp -R .codegen/*.go internal/openapi/ 2>/dev/null || true

# Clean up
rm -rf .codegen

# --- post-generation patch -----------------------------------------------------
#
# Drop TgvalidatordMetadata.payload.
#
# The proto declares it `google.protobuf.Value` (any JSON value) and the swagger
# correctly emits an untyped `{}`, but openapi-generator's Go target maps that to
# map[string]interface{} while the wire sends an ARRAY — so every requests read
# failed with "cannot unmarshal array into ... map[string]interface {}".
# `--type-mappings=AnyType=interface{}` does not change it.
#
# Removing the field also keeps the raw object out of the process entirely, which
# is what pkg/protect/mapper/request.go has always intended: the object can be
# altered while payloadAsString and its hash stay consistent, so only the hashed
# string is ever read.
#
# This MUST fail loudly. A silent no-op reintroduces both the decode failure and
# the tamperable field, and nothing downstream would say so.
patch_metadata_payload() {
    local file="internal/openapi/model_tgvalidatord_metadata.go"

    if [[ ! -f "$file" ]]; then
        echo "ERROR: $file not found after generation" >&2
        return 1
    fi
    if ! grep -q 'Payload map\[string\]interface{}' "$file"; then
        echo "ERROR: expected generated field 'Payload map[string]interface{}' in $file." >&2
        echo "       The generator's output changed. Re-check whether the payload field" >&2
        echo "       still needs removing before editing this guard." >&2
        return 1
    fi

    python3 - "$file" <<'PYEOF'
import re, sys
path = sys.argv[1]
src = open(path).read()

before = src
src = re.sub(r'\n\tPayload map\[string\]interface\{\} `json:"payload,omitempty"`', '', src, count=1)
for name in ("GetPayload", "GetPayloadOk", "HasPayload", "SetPayload"):
    src = re.sub(r'\n// ' + name + r' [^\n]*\n(?://[^\n]*\n)*func \(o \*?TgvalidatordMetadata\) '
                 + name + r'\([^\n]*\{.*?\n\}\n', '\n', src, count=1, flags=re.S)
src = re.sub(r'\tif !IsNil\(o\.Payload\) \{\n\t\ttoSerialize\["payload"\] = o\.Payload\n\t\}\n', '', src, count=1)

if src == before:
    sys.exit("ERROR: post-generation patch matched nothing in " + path)
if re.search(r'\bo\.Payload\b|\bPayload\s+map\[string\]|"payload"', src):
    sys.exit("ERROR: payload references remain in " + path + " after patching")
if "PayloadAsString" not in src:
    sys.exit("ERROR: patch removed PayloadAsString from " + path)

open(path, "w").write(src)
print("post-generation: removed TgvalidatordMetadata.payload")
PYEOF
}

if ! patch_metadata_payload; then
    echo "" >&2
    echo "Generation ABORTED: the payload patch did not apply." >&2
    echo "Leaving it unapplied would restore a decode failure on every requests read." >&2
    exit 1
fi

echo ""
echo "OpenAPI client generated successfully."
echo "Files are in: internal/openapi/"
echo ""
echo "Note: You may need to add TPV1 authentication to the generated client."
