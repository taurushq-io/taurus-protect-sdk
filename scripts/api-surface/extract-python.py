#!/usr/bin/env python3
"""Extracts the Python SDK's public service surface into the shared api-surface JSON.

Uses `inspect` against the imported package rather than parsing source, so what is
recorded is exactly what a caller can reach.
"""

from __future__ import annotations

import inspect
import json
import pkgutil
import sys
from pathlib import Path


def collect(services_pkg) -> list[dict]:
    services: list[dict] = []

    for module_info in pkgutil.walk_packages(
        services_pkg.__path__, prefix=services_pkg.__name__ + "."
    ):
        if module_info.name.rsplit(".", 1)[-1].startswith("_"):
            continue
        module = __import__(module_info.name, fromlist=["*"])

        for class_name, cls in inspect.getmembers(module, inspect.isclass):
            if not class_name.endswith("Service") or cls.__module__ != module_info.name:
                continue

            methods = []
            for name, fn in inspect.getmembers(cls, inspect.isfunction):
                if name.startswith("_") or fn.__qualname__.split(".")[0] != class_name:
                    continue
                try:
                    sig = str(inspect.signature(fn)).replace("self, ", "").replace("self", "")
                except (TypeError, ValueError):
                    sig = "(...)"
                doc = (inspect.getdoc(fn) or "").strip().split("\n")[0]
                methods.append({"name": name, "signature": f"{name}{sig}", "doc": doc})

            services.append(
                {"name": class_name, "methods": sorted(methods, key=lambda m: m["name"])}
            )

    return sorted(services, key=lambda s: s["name"])


def main() -> int:
    if len(sys.argv) < 3:
        print("usage: extract-python.py <sdk-root> <out.json>", file=sys.stderr)
        return 2

    sdk_root, out = Path(sys.argv[1]).resolve(), Path(sys.argv[2])
    sys.path.insert(0, str(sdk_root))

    import taurus_protect.services as services_pkg

    surface = {
        "sdk": "python",
        "source": "taurus_protect/services",
        "services": collect(services_pkg),
    }
    out.write_text(json.dumps(surface, indent=2) + "\n")
    print(f"python: {len(surface['services'])} services")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
