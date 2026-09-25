# Vendored openapi-generator templates (Python)

Only the files in this directory override the generator's built-ins. Everything else still comes
from `../../jars/openapi-generator-cli-7.9.0.jar`, so this directory is deliberately tiny — it is
not a fork of the template set. `taurus-protect-sdk-python/scripts/generate-openapi.sh` passes it
with `-t`.

## Why it exists

`model_enum.mustache` is overridden for one reason: an enum value the client does not know must
not fail the decode.

The built-in template generates a plain `class X(str, Enum)`. `X("NEW_VALUE")` raises
`ValueError`, and a pydantic model field typed `X` raises `ValidationError`, so one value the
server adds — a new token type, stake account state, entity kind — fails the WHOLE reply before
any SDK code runs. The override adds a `_missing_` classmethod that returns a pseudo-member
keeping the raw string:

- `X("NEW_VALUE").value == "NEW_VALUE"`, and it compares equal to the raw string;
- a known value still returns the canonical member (`X("A") is X.A`);
- pydantic honours it for model fields, lists and `validate_assignment`, and `to_dict()` /
  `to_json()` write the raw value back.

This is the same contract as the Go, Java and TypeScript SDKs — an unknown enum value keeps its
raw string, a value the caller supplies is sent verbatim — and it is gated, not inspected: every
SDK consumes `../../decode-tolerance-vectors.json`.

Pseudo-members are not cached (`X("NEW") is X("NEW")` is False) and `str()` of one reads
`X.UNKNOWN`; use `.value` for the raw string. Proven on Python 3.12.3 with pydantic 2.13.4;
not run on the 3.14 dev venv or the 3.9 floor.

## Two places, kept in step

The committed generated enum modules under `taurus_protect/_internal/openapi/models/` carry the
`_missing_` block because `generate-openapi.sh` rendered them from this template. The script
`rm -rf`s that directory before regenerating, so an edit to a generated module alone is lost on
the next run. Change this template, then regenerate.

## On upgrading the generator

This file was extracted from openapi-generator-cli **7.9.0**. If that JAR is upgraded, re-extract
`python/model_enum.mustache` from the new one and re-apply the `_missing_` block, rather than
keeping this copy — otherwise every unrelated upstream template improvement is silently pinned to
the 7.9.0 version.

```bash
cd scripts/resources/templates/python
unzip -o -j ../../jars/openapi-generator-cli-<version>.jar "python/model_enum.mustache"
# then re-apply the _missing_ block after from_json, and regenerate
```
