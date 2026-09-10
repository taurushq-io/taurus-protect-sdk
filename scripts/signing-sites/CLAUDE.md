# CLAUDE.md — the signing-site inventory gate

`python3 scripts/signing-sites/check.py` enumerates every ECDSA signing call site across
the four SDKs and requires each to be declared in `manifest.json`.

## Why it exists

The same defect was made four times and found by **grep, not by a test**: `approveRequests`
checked that a metadata hash was non-empty and then signed it, so the approver attested to
a hash nothing had checked. A fifth site — Python's `approve_pledge_actions` — had no
verification at all and drifted alone for the same reason: nothing enumerated the signing
surface, so a new signer inherited no rule.

So the manifest is the point, not the count. Adding a signing site fails the gate until it
declares which of two shapes it is:

| Shape | Meaning |
|---|---|
| `verifies` | the bytes signed come from data this SDK verified in the same call |
| `signs-own-bytes` | it signs a document the caller is authoring, so there is nothing to verify against yet |

`signs-own-bytes` is currently only the four `approveRulesProposal` methods: a governance
proposal legitimately carries 0..N signatures and signing IS the approval step, so a
threshold is an enforcement-time invariant rather than a precondition. Do not "fix" those
to verify — that is the same mistake as routing mcpd's `view=proposal` through
`GetDecodedRulesContainer`.

Current state: **17 sites — 13 `verifies`, 4 `signs-own-bytes`.**

## Design decisions that matter

- **Keyed on the enclosing SYMBOL, not the line number.** A line number moves whenever
  anything above it is edited, which would force a re-classification on every unrelated
  change and train people to rubber-stamp the diff. Proven both ways: adding a comment
  above a signing call keeps the gate green; adding a new `crypto.SignData` call fails it.
  Renaming a signing method *does* require an `--update` and a re-review, which is correct.
- **`enclosing_symbol` walks upward and skips a keyword blocklist.** One regex serves all
  four languages, so `if (...) {` and `str(...)` matched as declarations until `NOT_A_SYMBOL`
  filtered them. There is a second TS alternative for a signature whose paren opens at
  end-of-line (`async approve(` with params on following lines) — without it three TS sites
  resolved to `list`, `getEnvelope` and a bare line number.
- **The crypto helpers are excluded, not classified.** `HELPER_PATHS` drops
  `crypto/`, `tpv1`, `CryptoTPV1.java`, `signing.ts|py` — those *define* the signing
  primitive rather than calling it.
- **A missing manifest hard-fails** rather than passing vacuously, and `--update` exits
  non-zero while anything is `UNCLASSIFIED`. Both verified.
- **It walks the filesystem, not `git grep`.** This repo carries a large uncommitted tree,
  so a tracked-files-only scan would miss exactly the new code under review.

## Adding a signing path

1. Write the code, verifying before you sign.
2. `python3 scripts/signing-sites/check.py --update`
3. Classify the new entry, and **read the diff** — the tool never carries a classification
   across a changed key, so a moved or renamed site comes back `UNCLASSIFIED` on purpose.
