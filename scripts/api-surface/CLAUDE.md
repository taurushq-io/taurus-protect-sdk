# API-surface pipeline

Two of the five cross-SDK alignment gates live here (see the repo-root `CLAUDE.md` for the full
list). Four extractors emit a common `scripts/resources/api-surface.<lang>.json`; one `docs.py`
and one `diff.py` consume all four.

```
extract-go.go      ─┐
extract-java.py     ├─► scripts/resources/api-surface.<lang>.json ─┬─► docs.py  → <sdk>/docs/SERVICES.md
extract-python.py   │                                              └─► diff.py  → parity report
extract-ts.mjs     ─┘
generate.sh <sdk|all> [render|check]      # the entry point every build.sh docs target calls
selftest.sh                               # proves diff.py detects a gap
```

## Why each extractor uses the language's own tooling

Regex parsing would silently miss methods, and a missing method makes the differ report parity that
isn't there — worse than no gate.

| SDK | Tool | Constraint that drove the choice |
|---|---|---|
| Go | `go/ast` (stdlib) | **Not** `go/packages`: it needs a go1.26 toolchain in this container. `parser.ParseDir` needs no module loading. Service is taken from the method receiver, so the walk needs no section context. |
| Python | `inspect` on the imported package | Records what a caller can actually reach, not what the source appears to define. Needs the SDK's `.venv`. |
| TypeScript | TypeScript compiler API | Classifies private/protected the way `tsc` does. The script sits outside the package, so `typescript` is resolved with `createRequire` off the SDK root. |
| Java | `javap -public` on `client/target/classes` | Avoids writing a Java parser; overloads and visibility come out right. **Requires `mvn compile -o -pl client` first** — a stale `target/classes` yields a stale surface. |

Exclude abstract bases: TS's `BaseService` made that SDK report 44 services where the others report
43.

## diff.py — what is a gate and what is advisory

- **Service parity is a hard gate** (exit 1). A service missing from an SDK either gets implemented
  or gets recorded in `aliases.json`.
- **Method-count deltas are advisory.** Method names legitimately differ by idiom — Go `ListWallets`,
  Python `list`, Java `getWallets`. Normalising that away either hides real gaps or invents false
  ones, so the differ prints spreads for a human to read instead of failing on them. It was the
  spread that revealed Java's TaurusNetwork services are systematically thin (lending 4 vs 14/16/14,
  pledge 3 vs 14) — see `TODOS.md`.
- `aliases.json` groups per-SDK idiomatic names for the same service. **Every entry must correspond
  to a row in `docs/SDK_ALIGNMENT_REPORT.md` → "Remaining known differences"**; that pairing is what
  keeps an alias a recorded decision rather than a silent exemption.
- Keep `selftest.sh` green. It runs the differ against `testdata/known-gap/` (four surfaces, one
  service deliberately absent from TypeScript) and fails if the differ reports parity or does not
  name the offending SDK.

## docs.py — render and check

`render` injects a generated method index between `<!-- BEGIN/END GENERATED METHOD INDEX -->` in
`<sdk>/docs/SERVICES.md`. `check` fails on a stale index **or a documented method that does not
exist**.

**Do not regenerate those files wholesale.** They hold ~7k lines of hand-written prose —
descriptions, parameter tables, Key Models — that no extractor can reproduce. Only the marker region
is rewritten.

The phantom-method scan reads two doc shapes, because both are in use:

1. `#### MethodName` headings under a `## SomeService` heading.
2. Bare signature lines inside a fenced block — **but only inside a `### Methods` section.** Other
   fences are usage examples, where `for (...)` and `client.x.y()` parse as declarations. Scoping to
   the Methods section is what took the false-positive count to zero; without it Python reported 83
   phantoms, all noise.

A `##` heading that is not a service clears the service scope, so signatures in a TaurusNetwork
low-level section are not blamed on whichever service preceded them. Those sections document
generated `*Api` classes, which this surface deliberately does not cover.

First run found ~100 documented-but-nonexistent methods (Go 21, Java 26, TypeScript 4, plus a Go
TaurusNetwork section naming the wrong receiver type on every signature).

## Adding a service or method

1. Implement it in all four SDKs (or record the difference in `aliases.json` + the report).
2. `./scripts/api-surface/generate.sh all` — regenerates surfaces, re-renders the four indexes, runs
   the differ.
3. `./scripts/api-surface/generate.sh all check` — must print "docs match the code" for all four.
4. `./scripts/api-surface/selftest.sh`.

`JAVA_HOME` must be exported and the Java client compiled before step 2; see the repo-root
`CLAUDE.local.md`.
