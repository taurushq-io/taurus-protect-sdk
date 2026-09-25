#!/usr/bin/env python3
"""Adopt named operations, parameters and definitions from tg-validatord's swagger into the shared
snapshot (scripts/resources/swagger/apis.swagger.json), additively.

The snapshot feeds all four generators and is deliberately NOT refreshed wholesale (see the repo-root
CLAUDE.md "Shared swagger ... is a STALE SNAPSHOT"). This copies exactly what is named, plus every
definition the copied parts reference that the snapshot lacks, and refuses anything that would remove
or retype an existing member — so a regeneration can only ADD API surface.

    scripts/swagger/adopt-from-validatord.py \\
        --validatord-swagger ../tg-validatord/api/swagger/v1/apis.swagger.json \\
        --op PriceService_QueryPricesV2 \\
        --param WalletService_GetAddresses:includeDisabledAddresses \\
        --def tgvalidatordCurrencyPrice
"""
import argparse
import json
import sys
from pathlib import Path

DEFAULT_SPEC = Path(__file__).resolve().parents[1] / "resources" / "swagger" / "apis.swagger.json"
TYPE_KEYS = ("type", "format", "$ref", "items", "enum")


def go_style_dump(obj):
    """The snapshot is written by Go's encoding/json: 2-space indent, HTML characters escaped."""
    text = json.dumps(obj, indent=2, ensure_ascii=False)
    for char, escaped in (("<", "\\u003c"), (">", "\\u003e"), ("&", "\\u0026"),
                          ("\u2028", "\\u2028"), ("\u2029", "\\u2029")):
        text = text.replace(char, escaped)
    return text + "\n"


def operations(spec):
    found = {}
    for path, item in spec.get("paths", {}).items():
        for method, op in item.items():
            if isinstance(op, dict) and "operationId" in op:
                found[op["operationId"]] = (path, method, op)
    return found


def refs_in(obj, into):
    if isinstance(obj, dict):
        for key, value in obj.items():
            if key == "$ref" and isinstance(value, str) and value.startswith("#/definitions/"):
                into.add(value.split("/")[-1])
            else:
                refs_in(value, into)
    elif isinstance(obj, list):
        for value in obj:
            refs_in(value, into)
    return into


def closure(names, definitions):
    seen, todo = set(), list(names)
    while todo:
        name = todo.pop()
        if name in seen:
            continue
        seen.add(name)
        if name not in definitions:
            raise SystemExit(f"validatord swagger has no definition {name!r}")
        todo.extend(refs_in(definitions[name], set()) - seen)
    return seen


def shape(prop):
    return {k: prop.get(k) for k in TYPE_KEYS}


def assert_additive(name, old, new):
    """Every member of the snapshot's definition must survive with the same type."""
    old_props, new_props = old.get("properties", {}), new.get("properties", {})
    problems = [f"{name}.{prop} removed" for prop in old_props if prop not in new_props]
    problems += [f"{name}.{prop} retyped: {shape(old_props[prop])} -> {shape(new_props[prop])}"
                 for prop in old_props if prop in new_props and shape(old_props[prop]) != shape(new_props[prop])]
    for key in ("type", "enum"):
        if old.get(key) != new.get(key):
            problems.append(f"{name}.{key} changed: {old.get(key)} -> {new.get(key)}")
    if problems:
        raise SystemExit("refusing a non-additive definition change:\n  " + "\n  ".join(problems))


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--validatord-swagger", required=True, type=Path)
    parser.add_argument("--spec", type=Path, default=DEFAULT_SPEC)
    parser.add_argument("--op", action="append", default=[], help="operationId to copy")
    parser.add_argument("--param", action="append", default=[], help="operationId:parameterName to add")
    parser.add_argument("--def", dest="defs", action="append", default=[], help="definition to replace additively")
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    source = json.loads(args.validatord_swagger.read_text())
    target = json.loads(args.spec.read_text())
    src_ops, dst_ops = operations(source), operations(target)
    src_defs, dst_defs = source["definitions"], target["definitions"]
    wanted_defs, log = set(), []

    for op_id in args.op:
        if op_id not in src_ops:
            raise SystemExit(f"validatord swagger has no operation {op_id!r}")
        path, method, op = src_ops[op_id]
        if op_id in dst_ops:
            if dst_ops[op_id][2] != op:
                raise SystemExit(f"{op_id} already exists in the snapshot with different content; adopt it by hand")
            continue
        target["paths"].setdefault(path, {})[method] = op
        wanted_defs |= refs_in(op, set())
        log.append(f"operation {op_id} ({method.upper()} {path})")

    for spec_ in args.param:
        op_id, _, name = spec_.partition(":")
        src_params = src_ops[op_id][2].get("parameters", [])
        dst_params = dst_ops[op_id][2].setdefault("parameters", [])
        if any(p.get("name") == name for p in dst_params):
            raise SystemExit(f"{op_id} already has parameter {name!r}")
        index = next(i for i, p in enumerate(src_params) if p.get("name") == name)
        new_param = src_params[index]
        # Keep validatord's order: insert after the nearest preceding parameter the snapshot has.
        position = 0
        for prev in reversed(src_params[:index]):
            hit = next((i for i, p in enumerate(dst_params) if p.get("name") == prev.get("name")), None)
            if hit is not None:
                position = hit + 1
                break
        dst_params.insert(position, new_param)
        wanted_defs |= refs_in(new_param, set())
        log.append(f"parameter {op_id}:{name}")

    for name in args.defs:
        if name in dst_defs:
            assert_additive(name, dst_defs[name], src_defs[name])
        dst_defs[name] = src_defs[name]
        wanted_defs |= refs_in(src_defs[name], set())
        log.append(f"definition {name} (replaced additively)")

    for name in sorted(closure(wanted_defs, src_defs)):
        if name in dst_defs:
            if dst_defs[name] != src_defs[name] and name not in args.defs:
                raise SystemExit(f"{name} is referenced but differs from validatord; pass --def {name} to adopt it")
            continue
        dst_defs[name] = src_defs[name]
        log.append(f"definition {name} (new)")

    # The snapshot keeps both maps sorted; a stable order keeps regeneration diffs minimal.
    target["paths"] = dict(sorted(target["paths"].items()))
    target["definitions"] = dict(sorted(dst_defs.items()))
    missing = sorted(refs_in(target, set()) - set(target["definitions"]))
    if missing:
        raise SystemExit(f"unresolved $refs after the patch: {missing}")

    print("\n".join(log))
    print(f"{len(log)} change(s)")
    if not args.dry_run:
        args.spec.write_text(go_style_dump(target))


if __name__ == "__main__":
    sys.exit(main())
