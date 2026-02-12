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

# Generate Go client with enumClassPrefix to avoid const conflicts
java -jar "$GENERATOR_JAR" generate -g go -i "$SPEC_FILE" -o .codegen \
    --skip-validate-spec \
    --additional-properties=packageName=openapi \
    --additional-properties=isGoSubmodule=true \
    --additional-properties=enumClassPrefix=true

# Copy generated files to internal/openapi
cp -R .codegen/*.go internal/openapi/ 2>/dev/null || true

# Clean up
rm -rf .codegen

echo ""
echo "OpenAPI client generated successfully."
echo "Files are in: internal/openapi/"
echo ""
echo "Note: You may need to add TPV1 authentication to the generated client."
