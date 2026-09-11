# Vendored openapi-generator templates (Go)

Only the files in this directory override the generator's built-ins. Everything else still comes
from `../../jars/openapi-generator-cli-7.9.0.jar`, so this directory is deliberately tiny — it is
not a fork of the template set.

## Why it exists

`client.mustache` is overridden for one reason: the generated `decode` must reject a 2xx response
whose typed reply could not be populated.

An empty body, or a body consisting of the JSON literal `null`, otherwise leaves the caller's
reply pointer nil and returns **no error** — and every hand-written service in `pkg/protect`
dereferences the reply immediately after checking `err`. Measured before the fix: ~110 unchecked
`resp.<Field>` dereferences across 38 service files, and even the apparently-careful
`if resp.Result == nil` is itself a dereference of a nil `resp`.

That makes it a whole-process kill rather than a failed call, because nothing in the SDK recovers
panics. In `tg-protect-mcpd`, which runs one client per tenant in a single daemon, one empty 200
aimed at one tenant takes down every tenant, repeatably on every retry. A server that merely
errors or returns malformed JSON cannot do that.

The guard is `assertReplyDecoded`. Only the "asked to populate a typed reply and got nil" shape is
an error: `*string`, `*os.File`, `[]byte`, value types and `google.protobuf.Empty` endpoints all
legitimately accept an empty body.

## Two places, kept in step

The same change lives in **both**:

- `taurus-protect-sdk-go/internal/openapi/client.go` — the committed generated file, which is what
  actually ships today
- this template — which is what makes it survive `./build.sh generate`, since
  `scripts/generate-openapi.sh` does `rm -rf internal/openapi` before regenerating

Editing only the generated file means the next regeneration silently reverts the fix. Editing only
the template means the shipped code is unfixed until someone regenerates. Change both.

## On upgrading the generator

This file was extracted from openapi-generator-cli **7.9.0**. If that JAR is upgraded, re-extract
`go/client.mustache` from the new one and re-apply the `assertReplyDecoded` change, rather than
keeping this copy — otherwise every unrelated upstream template improvement is silently pinned to
the 7.9.0 version.

```bash
cd scripts/resources/templates/go
unzip -o -j ../../jars/openapi-generator-cli-<version>.jar "go/client.mustache"
# then re-apply the decode guard, and diff against internal/openapi/client.go
```
