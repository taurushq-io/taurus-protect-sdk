# Vendored openapi-generator templates (TypeScript, `typescript-fetch`)

Only the files in this directory override the generator's built-ins. Everything else still comes
from `../../jars/openapi-generator-cli-7.9.0.jar`, so this directory is deliberately tiny — it is
not a fork of the template set.

## Why it exists

`modelGeneric.mustache` and `modelGenericInterfaces.mustache` are overridden for one reason: a
field the server sends that this client does not know must be **kept**, not dropped.

All four SDKs follow one decode-tolerance contract, gated by the shared
`scripts/resources/decode-tolerance-vectors.json`:

- an unknown field never fails a decode; it lands in the generated model's additional-properties
  map, and serializing writes it back after the known fields;
- an unknown enum value stays the raw string;
- a caller-supplied enum value is sent verbatim.

Go, Java and Python keep unknown fields through their generators
(`disallowAdditionalPropertiesIfNotPresent=false`). `typescript-fetch` 7.9.0 is invoked with the same
flag but still drops them: its `...json` spread only fires for a schema with a typed
`additionalProperties`, and it could not be used anyway — some wire names differ from the TS
property names (`BatchSignatureID` is `batchSignatureID`), so spreading the whole reply would
duplicate them. The overrides therefore:

- add `additionalProperties?: { [key: string]: any }` to every model interface that has properties;
- make `FromJSONTyped` copy every key that is not one of the model's wire names (a per-model
  `…WireKeys` set) into it, attached only when non-empty so existing fixtures do not change;
- make `ToJSONTyped` spread it back after the known fields, the same order as Go's `ToMap`, Java's
  adapter and Python's `to_dict`.

A schema with its own typed `additionalProperties` (`ProtobufAny`) keeps the upstream index
signature and spread untouched. Enums need no override: `XFromJSONTyped` is a plain cast, so an
unknown value already passes through.

## Two places, kept in step

The generated files in `taurus-protect-sdk-typescript/src/internal/openapi/models/` are what
actually ships. `scripts/generate-openapi.sh` does `rm -rf` of that directory before regenerating
and passes this directory with `-t`, so change the template and regenerate — never edit the
generated models by hand.

## On upgrading the generator

These files were extracted from openapi-generator-cli **7.9.0**. If that JAR is upgraded,
re-extract both from the new one and re-apply the changes above (search for `WireKeys` and
`additionalProperties`), rather than keeping these copies — otherwise every unrelated upstream
template improvement is silently pinned to the 7.9.0 version.

```bash
cd scripts/resources/templates/typescript-fetch
unzip -o -j ../../jars/openapi-generator-cli-<version>.jar \
  "typescript-fetch/modelGeneric.mustache" "typescript-fetch/modelGenericInterfaces.mustache"
# then re-apply the additionalProperties change and run tests/unit/decode-tolerance-vectors.test.ts
```
