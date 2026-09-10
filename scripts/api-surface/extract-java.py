#!/usr/bin/env python3
"""Extracts the Java SDK's public service surface into the shared api-surface JSON.

Reads the COMPILED classes with javap rather than parsing .java: what javap prints is
exactly the public API, so overloads, inherited members and visibility come out right
without a Java parser. Requires `mvn compile -o -pl client` to have run.
"""

from __future__ import annotations

import json
import re
import subprocess
import sys
from pathlib import Path

# javap prints e.g.
#   public com.taurushq...BusinessRuleResult getBusinessRules(...ApiRequestCursor) throws ...;
SIGNATURE = re.compile(
    r"^\s*public\s+(?:final\s+|static\s+|synchronized\s+)*"
    r"(?P<ret>[\w.$<>\[\],?\s]+?)\s+"
    r"(?P<name>\w+)\((?P<args>[^)]*)\)"
)


def simple(type_name: str) -> str:
    """Strips packages so signatures read like the other SDKs' output."""
    return re.sub(r"\b[\w$]+\.", "", type_name.strip())


def methods_of(class_file: Path, classes_root: Path) -> tuple[str, list[dict]]:
    fqcn = ".".join(class_file.relative_to(classes_root).with_suffix("").parts)
    proc = subprocess.run(
        ["javap", "-public", "-classpath", str(classes_root), fqcn],
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        raise SystemExit(f"javap failed for {fqcn}: {proc.stderr.strip()}")

    methods = []
    class_name = fqcn.rsplit(".", 1)[-1]
    for line in proc.stdout.splitlines():
        m = SIGNATURE.match(line)
        if not m:
            continue
        name = m.group("name")
        # constructors print as the class name; skip them and any accessor of the class
        if name == class_name:
            continue
        args = ", ".join(simple(a) for a in m.group("args").split(",") if a.strip())
        methods.append(
            {
                "name": name,
                "signature": f"{name}({args}): {simple(m.group('ret'))}",
                "doc": "",
            }
        )

    methods.sort(key=lambda x: (x["name"], x["signature"]))
    return class_name, methods


def main() -> int:
    if len(sys.argv) < 3:
        print("usage: extract-java.py <sdk-root> <out.json>", file=sys.stderr)
        return 2

    sdk_root, out = Path(sys.argv[1]).resolve(), Path(sys.argv[2])
    classes_root = sdk_root / "client" / "target" / "classes"
    service_dir = classes_root / "com/taurushq/sdk/protect/client/service"

    if not service_dir.is_dir():
        print(
            f"compiled classes not found at {service_dir}\n"
            "run: mvn compile -o -pl client",
            file=sys.stderr,
        )
        return 1

    services = []
    for class_file in sorted(service_dir.glob("*Service.class")):
        if "$" in class_file.name:
            continue
        name, methods = methods_of(class_file, classes_root)
        services.append({"name": name, "methods": methods})

    services.sort(key=lambda s: s["name"])
    surface = {
        "sdk": "java",
        "source": "client/src/main/java/.../client/service",
        "services": services,
    }
    out.write_text(json.dumps(surface, indent=2) + "\n")
    print(f"java: {len(services)} services")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
