#!/usr/bin/env python3
"""
Cross-SDK gate: every ECDSA signing site must declare whether it verifies first.

Why this exists. The same defect was made four times and found by grep, not by a test:
`approveRequests` checked that a metadata hash was non-empty and then signed it, so the
approver attested to a hash nothing had checked. A fifth site (Python's
`approve_pledge_actions`) had no verification at all and drifted alone for the same
reason -- nothing enumerated the signing surface, so a new signer inherited no rule.

So the manifest is the point, not the count: adding a signing site forces an entry, and
the entry has to say which of two shapes it is.

  verifies         the bytes signed come from data this SDK verified in the same call
  signs-own-bytes  it signs a document the caller is authoring, so there is nothing to
                   verify against yet (a governance proposal carries 0..N signatures --
                   signing IS the approval step)

Keyed on the enclosing symbol, not the line number: a line number moves whenever anything
above it is edited, which would force a re-classification on every unrelated change and
train people to rubber-stamp the diff.

Run:  python3 scripts/signing-sites/check.py           # verify against the manifest
      python3 scripts/signing-sites/check.py --update  # rewrite the manifest
"""

from __future__ import annotations

import json
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
MANIFEST = Path(__file__).parent / "manifest.json"

# The signing entry points, per language. Anything that turns bytes into a signature.
SIGN_CALLS = [
    "sign_data(",
    "signData(",
    "calculateBase64Signature(",
    "crypto.SignData(",
]

# Where production source lives, and what is not production source.
TREES = {
    "go": ("taurus-protect-sdk-go/pkg", ("_test.go",)),
    "java": ("taurus-protect-sdk-java/client/src/main", ()),
    "python": ("taurus-protect-sdk-python/taurus_protect", ()),
    "typescript": ("taurus-protect-sdk-typescript/src", ()),
}

# The crypto helpers themselves DEFINE the signing primitive; they are not call sites.
HELPER_PATHS = re.compile(
    r"(/crypto/|crypto_|tpv1|TPV1\.java|/signing\.(ts|py)$|/auth/CryptoTPV1)"
)


# Declaration lines, per language. Keyed on the enclosing symbol rather than the line
# number: a line number moves whenever anything above it is edited, which would force a
# re-classification on every unrelated change and train people to rubber-stamp the diff.
DECL = re.compile(
    r"^\s*(?:"
    r"func\s+(?:\([^)]*\)\s*)?(?P<go>\w+)"                       # Go
    r"|(?:public|private|protected)\s+[\w<>\[\], .]+?\s(?P<java>\w+)\s*\("  # Java
    r"|def\s+(?P<py>\w+)"                                          # Python
    # TS: params on one line, or a multi-line signature whose paren opens at EOL.
    r"|(?:async\s+)?(?:private\s+|public\s+)?(?P<ts>\w+)\s*\([^)]*\)\s*[:{]"
    r"|(?:async\s+)?(?:private\s+|public\s+)?(?P<ts2>\w+)\s*\(\s*$"
    r")"
)


# Control flow and calls look like declarations to the patterns above. Skipping them
# rather than tightening the regexes per language keeps this a ~20-line scanner; a wrong
# key would be caught on review, and the fallback is still the line number.
NOT_A_SYMBOL = {
    "if", "for", "while", "switch", "catch", "try", "else", "do", "return", "throw",
    "new", "str", "String", "await", "async", "function", "get", "set", "constructor",
}


def enclosing_symbol(lines: list[str], lineno: int) -> str:
    """The nearest declaration at or above `lineno`, or the line number if none."""
    for i in range(lineno - 1, -1, -1):
        m = DECL.match(lines[i])
        if not m:
            continue
        for name in ("go", "java", "py", "ts", "ts2"):
            sym = m.group(name)
            if sym and sym not in NOT_A_SYMBOL:
                return sym
    return str(lineno)


def discover() -> dict[str, list[str]]:
    """Every signing call site, as `<path>:<line>` keyed by SDK."""
    found: dict[str, list[str]] = {}
    for sdk, (tree, skip_suffixes) in TREES.items():
        base = ROOT / tree
        hits: list[str] = []
        # git grep would miss untracked files, and this repo carries a large
        # uncommitted tree -- so walk the filesystem.
        for path in sorted(base.rglob("*")):
            if not path.is_file():
                continue
            if path.suffix not in {".go", ".java", ".py", ".ts"}:
                continue
            rel = path.relative_to(ROOT).as_posix()
            if any(rel.endswith(s) for s in skip_suffixes):
                continue
            if HELPER_PATHS.search("/" + rel):
                continue
            try:
                lines = path.read_text(encoding="utf-8").splitlines()
            except UnicodeDecodeError:
                continue
            for n, line in enumerate(lines, start=1):
                stripped = line.strip()
                # Comments are not call sites. This is deliberately crude: a false
                # positive costs one manifest entry, a false negative costs the gate.
                if stripped.startswith(("//", "*", "#", "/*")):
                    continue
                if any(call in line for call in SIGN_CALLS):
                    hits.append(f"{rel}:{enclosing_symbol(lines, n)}")
        found[sdk] = hits
    return found


def load_manifest() -> dict:
    if not MANIFEST.exists():
        sys.exit(
            f"{MANIFEST} is missing. The gate hard-fails rather than passing vacuously: "
            "run --update and review the result."
        )
    return json.loads(MANIFEST.read_text(encoding="utf-8"))


def main() -> int:
    found = discover()

    if "--update" in sys.argv:
        existing = json.loads(MANIFEST.read_text(encoding="utf-8")) if MANIFEST.exists() else {}
        sites = existing.get("sites", {})
        out = {}
        for sdk, hits in found.items():
            out[sdk] = {
                site: sites.get(sdk, {}).get(site, "UNCLASSIFIED")
                for site in hits
            }
        MANIFEST.write_text(
            json.dumps({"_comment": __doc__.strip().splitlines(), "sites": out}, indent=2)
            + "\n",
            encoding="utf-8",
        )
        print(f"wrote {MANIFEST}")
        unclassified = [
            f"{sdk} {site}"
            for sdk, s in out.items()
            for site, kind in s.items()
            if kind == "UNCLASSIFIED"
        ]
        if unclassified:
            print("\nCLASSIFY THESE before committing:")
            for u in unclassified:
                print(f"  {u}")
            return 1
        return 0

    manifest = load_manifest()
    declared = manifest.get("sites", {})
    failures: list[str] = []

    for sdk, hits in found.items():
        known = declared.get(sdk, {})
        for site in hits:
            if site not in known:
                failures.append(
                    f"{sdk}: UNDECLARED signing site {site} — add it to the manifest as "
                    "'verifies' or 'signs-own-bytes'"
                )
            elif known[site] == "UNCLASSIFIED":
                failures.append(f"{sdk}: {site} is still UNCLASSIFIED")
        for site in known:
            if site not in hits:
                failures.append(
                    f"{sdk}: manifest names {site}, which no longer signs anything — "
                    "line numbers move, so re-run --update and re-review"
                )

    total = sum(len(h) for h in found.values())
    if failures:
        print(f"signing-site gate FAILED ({total} sites found)\n")
        for f in failures:
            print(f"  {f}")
        return 1

    by_kind: dict[str, int] = {}
    for sdk, known in declared.items():
        for kind in known.values():
            by_kind[kind] = by_kind.get(kind, 0) + 1
    print(f"signing-site gate OK: {total} sites, all declared")
    for kind, n in sorted(by_kind.items()):
        print(f"  {kind}: {n}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
