#!/usr/bin/env python3
"""Generates and gates the per-SDK service reference from the extracted API surface.

Two subcommands:

  render  Injects a generated method index into <sdk>/docs/SERVICES.md between markers.
          The hand-written prose (descriptions, parameter tables, Key Models) is left
          alone — only the region between the markers is rewritten.

  check   Verifies that every method the prose documents actually exists in the SDK, and
          that the generated region is up to date. This is the gate: the four SERVICES.md
          files had accumulated ~96 documented methods that do not exist, which is worse
          than missing documentation because a reader has no way to tell.

Usage:
  docs.py render <sdk> <surface.json> <services.md>
  docs.py check  <sdk> <surface.json> <services.md>
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

BEGIN = "<!-- BEGIN GENERATED METHOD INDEX -->"
END = "<!-- END GENERATED METHOD INDEX -->"

# All four docs list methods two ways: as `#### MethodName` headings, and as bare
# signature lines inside a fenced block under a `## SomeService` section. Both shapes
# have to be checked — the fenced signatures are where most of the phantom entries were.
SERVICE_HEADING = re.compile(r"^#{2,3}\s+`?(\w*Service\w*)`?\s*$")
METHOD_HEADING = re.compile(r"^####\s+`?([A-Za-z_][\w]*)`?\s*(?:\(.*\))?\s*$")
ANY_HEADING = re.compile(r"^#{1,6}\s+(.*?)\s*$")
# Only fences inside a "Methods" section are signature lists; elsewhere a fence is a
# usage example, where `for (...)` and `client.x.y()` would read as declarations.
METHODS_SECTION = re.compile(r"^#{3,4}\s+Methods?\b", re.I)
# Headings that are structure, not method names.
NON_METHOD_HEADINGS = {"methods", "method", "example", "examples", "key models", "usage"}

FENCE = re.compile(r"^```")

# Signature lines inside fences, per language. Deliberately narrow: a false positive
# here would turn the gate into noise, so each pattern anchors on syntax that only a
# declaration has.
FENCED_SIGNATURE = {
    # func (s *WalletService) GetWallet(ctx context.Context, ...) (...)
    "go": re.compile(r"^func\s+\(\s*\w+\s+\*?(?P<svc>\w+)\s*\)\s+(?P<name>\w+)\s*\("),
    # def get_wallet(self, wallet_id: int) -> Wallet:
    "python": re.compile(r"^\s*def\s+(?P<name>\w+)\s*\("),
    # WalletResult getWallets(ApiRequestCursor cursor)   /  public void close()
    "java": re.compile(
        r"^\s*(?:public\s+|protected\s+)?(?:static\s+|final\s+)*"
        r"[A-Za-z_][\w<>\[\],.?\s]*\s+(?P<name>[a-z]\w*)\s*\("
    ),
    # async getWallet(walletId: number): Promise<Wallet>   — a return type is required,
    # which is what separates a declaration from a call in an example.
    "typescript": re.compile(
        r"^\s*(?:async\s+)?(?P<name>[a-z]\w*)\s*(?:<[^>]*>)?\([^)]*\)\s*:\s*\S"
    ),
}


def load_surface(path: Path) -> dict[str, set[str]]:
    data = json.loads(Path(path).read_text())
    return {s["name"]: {m["name"] for m in s["methods"]} for s in data["services"]}


def documented(md_path: Path, sdk: str) -> list[tuple[str, str, int]]:
    """(service, method, line) for every method the prose claims exists."""
    out: list[tuple[str, str, int]] = []
    service = None
    in_fence = False
    in_generated = False
    in_methods_section = False
    fenced = FENCED_SIGNATURE.get(sdk)

    for lineno, line in enumerate(md_path.read_text().splitlines(), start=1):
        if line.strip() == BEGIN:
            in_generated = True
            continue
        if line.strip() == END:
            in_generated = False
            continue
        if in_generated:
            continue
        if FENCE.match(line):
            in_fence = not in_fence
            continue

        if in_fence:
            if not fenced or not in_methods_section:
                continue
            m = fenced.match(line)
            if not m:
                continue
            # Go carries the service in the receiver, so it does not need the section.
            svc = m.groupdict().get("svc") or service
            if svc:
                out.append((svc, m.group("name"), lineno))
            continue

        m = SERVICE_HEADING.match(line)
        if m:
            service = m.group(1)
            in_methods_section = False
            continue

        if METHODS_SECTION.match(line):
            in_methods_section = True
            continue

        m = METHOD_HEADING.match(line)
        if m and service:
            if m.group(1).strip().lower() not in NON_METHOD_HEADINGS:
                out.append((service, m.group(1), lineno))
            continue

        # Any other heading ends the Methods section (Key Models, Example, next service).
        if ANY_HEADING.match(line) and not METHOD_HEADING.match(line):
            heading = ANY_HEADING.match(line).group(1).strip().lower()
            if heading in NON_METHOD_HEADINGS - {"methods", "method"} or line.startswith("## "):
                in_methods_section = False
            # A top-level heading that is not a service (e.g. the TaurusNetwork
            # low-level API section) ends the service scope, so its signatures are not
            # blamed on whichever service happened to come before it. Those sections
            # document generated Api classes, which this surface does not cover.
            if line.startswith("## ") and not SERVICE_HEADING.match(line):
                service = None
    return out


def render_index(sdk: str, surface: dict[str, set[str]], surface_json: Path) -> str:
    data = json.loads(surface_json.read_text())
    lines = [
        BEGIN,
        "",
        "## Complete Method Index",
        "",
        f"Generated from the {sdk} source by `scripts/api-surface/docs.py`; regenerate with",
        "`./build.sh docs`. Every method below exists in the SDK, and `./build.sh docs --check`",
        "fails if this list drifts or if the prose above documents a method that does not.",
        "",
        f"{len(data['services'])} services, "
        f"{sum(len(s['methods']) for s in data['services'])} public methods.",
        "",
    ]
    for svc in data["services"]:
        lines.append(f"### {svc['name']}")
        lines.append("")
        if not svc["methods"]:
            lines.append("_No public methods._")
            lines.append("")
            continue
        for meth in svc["methods"]:
            doc = meth.get("doc") or ""
            doc = re.sub(r"\s+", " ", doc).strip()
            suffix = f" — {doc}" if doc else ""
            lines.append(f"- `{meth['signature']}`{suffix}")
        lines.append("")
    lines.append(END)
    return "\n".join(lines)


def replace_region(text: str, region: str) -> str:
    if BEGIN in text and END in text:
        head = text[: text.index(BEGIN)]
        tail = text[text.index(END) + len(END) :]
        return head + region + tail
    sep = "" if text.endswith("\n") else "\n"
    return text + sep + "\n" + region + "\n"


def main() -> int:
    if len(sys.argv) != 5:
        print(__doc__, file=sys.stderr)
        return 2

    cmd, sdk, surface_json, md = sys.argv[1], sys.argv[2], Path(sys.argv[3]), Path(sys.argv[4])
    surface = load_surface(surface_json)
    region = render_index(sdk, surface, surface_json)

    if cmd == "render":
        md.write_text(replace_region(md.read_text(), region))
        print(f"{sdk}: method index written to {md}")
        return 0

    if cmd != "check":
        print(__doc__, file=sys.stderr)
        return 2

    failures = 0

    phantoms = [
        (svc, meth, line)
        for svc, meth, line in documented(md, sdk)
        if svc in surface and meth not in surface[svc]
    ]
    if phantoms:
        failures += 1
        print(f"{sdk}: {len(phantoms)} documented method(s) do not exist:")
        for svc, meth, line in phantoms:
            print(f"  {md}:{line}: {svc}.{meth}")

    unknown_services = sorted(
        {svc for svc, _, _ in documented(md, sdk)} - set(surface) - {"BaseService"}
    )
    if unknown_services:
        failures += 1
        print(f"{sdk}: documented service(s) do not exist: {', '.join(unknown_services)}")

    if region not in md.read_text():
        failures += 1
        print(f"{sdk}: generated method index is stale — run ./build.sh docs")

    if not failures:
        print(f"{sdk}: docs match the code.")
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
