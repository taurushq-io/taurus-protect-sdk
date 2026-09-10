#!/usr/bin/env python3
"""Cross-SDK API-surface differ — the third alignment gate.

The other two gates cover wire bytes (governance-cell-vectors.json) and error parsing
(authorization-error-vectors.json). Nothing covered the SERVICE SURFACE, so a method
added to one SDK and forgotten in the other three was invisible until someone diffed by
hand — which is how getPublicKeys, UpdateTransactionsEnabled and the transaction chain
filters each ended up in only one or two SDKs.

Service-level parity is a hard gate: exit 1 when a service is missing from an SDK and
the pair is not recorded in aliases.json. Method-level output is ADVISORY and printed
for review, because method names legitimately differ by language idiom (Go ListWallets
vs Python list vs Java getWallets) and normalising that away would either hide real gaps
or invent false ones.
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

LANGS = ("go", "java", "python", "typescript")


def canonical(name: str) -> str:
    name = re.sub(r"Service$", "", name)
    name = re.sub(r"^TaurusNetwork", "", name)
    return re.sub(r"[^a-z0-9]", "", name.lower())


def load(surface_dir: Path) -> dict[str, dict]:
    surfaces = {}
    for lang in LANGS:
        path = surface_dir / f"api-surface.{lang}.json"
        if not path.exists():
            raise SystemExit(
                f"missing {path}\nrun each SDK's ./build.sh docs (or scripts/api-surface/generate.sh)"
            )
        surfaces[lang] = json.loads(path.read_text())
    return surfaces


def alias_map(path: Path) -> dict[str, str]:
    """Maps every alias onto its group's canonical name."""
    if not path.exists():
        return {}
    data = json.loads(path.read_text())
    mapping = {}
    for group in data.get("service_groups", []):
        for alias in group["aliases"]:
            mapping[alias] = group["canonical"]
    return mapping


def main() -> int:
    here = Path(__file__).resolve().parent
    surface_dir = here.parent / "resources"
    surfaces = load(surface_dir)
    aliases = alias_map(here / "aliases.json")

    def key(name: str) -> str:
        c = canonical(name)
        return aliases.get(c, c)

    by_lang = {
        lang: {key(s["name"]): s for s in surfaces[lang]["services"]} for lang in LANGS
    }
    every = sorted(set().union(*(set(v) for v in by_lang.values())))

    missing = []
    for svc in every:
        absent = [lang for lang in LANGS if svc not in by_lang[lang]]
        if absent:
            missing.append((svc, absent))

    print(f"Services: {len(every)} canonical")
    for lang in LANGS:
        total = sum(len(s["methods"]) for s in surfaces[lang]["services"])
        print(f"  {lang:11s} {len(by_lang[lang]):3d} services  {total:4d} public methods")

    if missing:
        print("\nSERVICE PARITY FAILURES (add the method, or record the pair in aliases.json):")
        for svc, absent in missing:
            present = {
                lang: by_lang[lang][svc]["name"] for lang in LANGS if svc in by_lang[lang]
            }
            print(f"  {svc}: missing from {', '.join(absent)}   (present as {present})")
    else:
        print("\nService parity: OK — every service exists in all four SDKs.")

    # Advisory: where one SDK exposes a very different number of methods for a service,
    # that is where a forgotten port usually hides.
    print("\nMethod-count deltas (advisory — method names differ by language idiom):")
    rows = []
    for svc in every:
        counts = {
            lang: len(by_lang[lang][svc]["methods"]) if svc in by_lang[lang] else 0
            for lang in LANGS
        }
        spread = max(counts.values()) - min(counts.values())
        if spread >= 2:
            rows.append((spread, svc, counts))
    for spread, svc, counts in sorted(rows, reverse=True):
        detail = "  ".join(f"{lang[:2]}={counts[lang]}" for lang in LANGS)
        print(f"  spread {spread:2d}  {svc:26s} {detail}")
    if not rows:
        print("  none with a spread of 2 or more.")

    return 1 if missing else 0


if __name__ == "__main__":
    raise SystemExit(main())
