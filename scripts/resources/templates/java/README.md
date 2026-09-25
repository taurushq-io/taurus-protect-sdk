# Vendored openapi-generator templates (Java)

Only the files in this directory override the generator's built-ins. Everything else still comes
from `../../jars/openapi-generator-cli-7.9.0.jar`, so this directory is deliberately tiny — it is
not a fork of the template set. `taurus-protect-sdk-java/scripts/generate-openapi.sh` passes it with
`-t`; library-specific overrides live under `libraries/<library>/`, exactly as in the JAR.

## Why it exists

`libraries/okhttp-gson/modelEnum.mustache` turns every generated enum into an **open class**: a
`final class` whose known values are `public static final` constants, and whose `fromValue` keeps a
value the client was not generated with instead of throwing.

The built-in template emits a Java `enum` whose `fromValue` throws
`IllegalArgumentException("Unexpected value ...")`, and both the Gson adapter and
`validateJsonElement` go through it. One enum value added on the server therefore failed the WHOLE
reply decode, before any SDK code ran — a new token type under `currency.tokenInfo` would fail every
balances page. A Java `enum` cannot carry a value it does not declare, and
`enumUnknownDefaultCase` only maps it to a sentinel, losing the value. The other three SDKs keep the
raw value (Go and TypeScript natively once decoding is lenient, Python through `_missing_`), and the
four must behave the same way, so the type itself has to be open.

What stays the same for callers: `X.CONSTANT` still names the known values (singletons, so `==`
works for them), `getValue()`, `toString()`, `values()` and `fromValue(...)` keep their signatures.
What changes: an unknown value is a distinct instance equal by value, and `switch`, `name()`,
`ordinal()` and `EnumSet`/`EnumMap` are no longer available on generated enum types.

The behaviour is gated, not inspected: `scripts/resources/decode-tolerance-vectors.json` is consumed
by `DecodeToleranceVectorsTest` here and by the Go, Python and TypeScript suites, and
`generate-openapi.sh` fails loudly if any generated model still throws `Unexpected value`.

## On upgrading the generator

`modelEnum.mustache` was extracted from openapi-generator-cli **7.9.0** and rewritten. If that JAR
is upgraded, re-extract the upstream file from the new one and re-apply the open-class shape, rather
than keeping this copy — otherwise every unrelated upstream template improvement is silently pinned
to the 7.9.0 version.

```bash
cd scripts/resources/templates/java/libraries/okhttp-gson
unzip -p ../../../../jars/openapi-generator-cli-<version>.jar \
    "Java/libraries/okhttp-gson/modelEnum.mustache" > /tmp/modelEnum.upstream.mustache
# diff against this file, re-apply the open-class changes, then regenerate and run
# DecodeToleranceVectorsTest
```
