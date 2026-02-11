#!/usr/bin/env bash

# Main procedure.
set -e

on_exit () {
    echo "generate-openapi.sh has exited in error"
}

trap on_exit ERR

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
cd "$DIR"

JAVA_VER=$(javap -verbose java.lang.String | grep "major version" | cut -d " " -f5)

if [[ $JAVA_VER -lt 55 ]]
   then
       echo "In order to use the Java openapi-generator-cli, you need Java 11 runtime at minimum (class file version 55.0). Alternatives are available here: https://openapi-generator.tech/"
       exit
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

mkdir -p .codegen
rm -rf .codegen

java -jar "$GENERATOR_JAR" generate -g java -i "$SPEC_FILE" -o .codegen \
    --skip-validate-spec  \
    --additional-properties=disallowAdditionalPropertiesIfNotPresent=false \
    --invoker-package com.taurushq.sdk.protect.openapi \
    --api-package com.taurushq.sdk.protect.openapi.api \
    --model-package com.taurushq.sdk.protect.openapi.model

cp -R .codegen/src/main/java openapi/src/main
rm -rf .codegen

# patch ApiClient to add TPV1 required functionnality
patch openapi/src/main/java/com/taurushq/sdk/protect/openapi/ApiClient.java < scripts/openapi-tpv1.patch