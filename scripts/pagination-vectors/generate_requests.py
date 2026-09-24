#!/usr/bin/env python3
"""Generate scripts/resources/list-request-vectors.json: what every paged operation must put on the wire.

Keyed by operationId so no SDK's method naming wins; each SDK's loader maps an operationId to its own
method and the canonical options below to its option fields. Paging parameter names are read from the
shared swagger at generation time, so a vector can never name a parameter the spec does not have.
"""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SPEC = ROOT / "resources" / "swagger" / "apis.swagger.json"
RULES = ROOT / "resources" / "pagination-vectors.json"
OUT = ROOT / "resources" / "list-request-vectors.json"
TOKEN = "ab+/cQ=="  # carries every base64 character that needs URL-escaping; canonical, so an SDK whose
# generated cursor type is bytes (Java) decodes and re-encodes it to the same text


def spec_operations():
    spec = json.loads(SPEC.read_text())
    ops = {}
    for path, item in spec["paths"].items():
        for method, op in item.items():
            if isinstance(op, dict) and "operationId" in op:
                ops[op["operationId"]] = (method.upper(), path, op, spec["definitions"])
    return ops


def paging_shape(op_id, rule, op, definitions):
    """Where page size and continuation go for this operation."""
    query = {p["name"] for p in op.get("parameters", []) if p.get("in") == "query"}
    if rule == "limit_only":
        # Decided by the rule, not the spec: the transaction export declares an offset the server ignores.
        if "limit" not in query:
            raise SystemExit(f"{op_id}: limit_only without a limit parameter")
        return {"style": "limit_only"}
    body = next((p for p in op.get("parameters", []) if p.get("in") == "body"), None)
    body_props = {}
    if body:
        body_props = definitions.get(body["schema"].get("$ref", "").split("/")[-1], {}).get("properties", {})
    if "requestCursor" in body_props:
        return {"style": "body_cursor", "field": "requestCursor"}
    if "cursor" in body_props and "$ref" in body_props["cursor"]:
        return {"style": "body_cursor", "field": "cursor"}
    if "requestCursor.pageSize" in query:
        return {"style": "query_cursor", "prefix": "requestCursor"}
    if "cursor.pageSize" in query:
        return {"style": "query_cursor", "prefix": "cursor"}
    if {"limit", "cursor"} <= query:
        return {"style": "token", "size": "limit", "token": "cursor"}
    if {"limit", "offset"} <= query:
        return {"style": "offset"}
    if "limit" in query:
        return {"style": "limit_only"}
    raise SystemExit(f"{op_id}: no recognisable paging parameters")


PATH_SENTINELS = {"id": "7", "assetID": "a1", "base": "BTC", "quote": "USD"}
# Options a call cannot omit, with where they land on the wire.
REQUIRED = {
    "FiatProviderService_GetFiatProviderAccounts": ({"provider": "r-provider", "label": "r-label"},
                                                    {"query": [["label", "r-label"], ["provider", "r-provider"]]}),
    "FiatProviderService_GetFiatProviderCounterpartyAccounts": ({"provider": "r-provider", "label": "r-label"},
                                                                {"query": [["label", "r-label"], ["provider", "r-provider"]]}),
    "WalletService_GetAssetAddresses": ({"currency": "r-currency"}, {"body": {"asset": {"currency": "r-currency"}}}),
    "WalletService_GetAssetWallets": ({"currency": "r-currency"}, {"body": {"asset": {"currency": "r-currency"}}}),
}
# Parameters every SDK sends on every call of the operation.
ALWAYS = {"WhitelistService_GetWhitelistedAddresses": [["rulesContainerNormalized", "true"]]}


def path_params(op):
    return {p["name"]: PATH_SENTINELS.get(p["name"], f"p-{p['name']}") for p in op.get("parameters", []) if p.get("in") == "path"}


def with_fixed(op_id, vector):
    """Fold the operation's required and always-sent parameters into a vector."""
    required, required_wire = REQUIRED.get(op_id, ({}, {}))
    vector["options"] = {**required, **vector["options"]}
    if "expect" in vector:
        wire_ = vector["expect"]
        if "query" in wire_ or "query" in required_wire or op_id in ALWAYS:
            wire_["query"] = sorted(wire_.get("query", []) + required_wire.get("query", []) + ALWAYS.get(op_id, []))
        if "body" in required_wire:
            wire_["body"] = {**required_wire["body"], **wire_.get("body", {})}
    return vector


def wire(shape, size, continuation=None):
    """The expected request for a page size and an optional continuation (offset int or cursor str)."""
    style = shape["style"]
    if style == "offset":
        query = [["limit", str(size)]] + ([["offset", str(continuation)]] if continuation else [])
        return {"query": sorted(query)}
    if style in ("limit_only", "token"):
        query = [["limit", str(size)]]
        if continuation:
            query.append([shape["token"], continuation])
        return {"query": sorted(query)}
    cursor = {"pageSize": str(size)}
    if continuation:
        cursor.update({"currentPage": continuation, "pageRequest": "NEXT"})
    if style == "query_cursor":
        return {"query": [[f"{shape['prefix']}.{k}", v] for k, v in sorted(cursor.items())]}
    return {"body": {shape["field"]: cursor}}


def vectors_for(op_id, entry, shape, params):
    base = {"operation": op_id, "path_params": params}
    offset_style = shape["style"] == "offset"
    size_opt = "limit" if shape["style"] in ("offset", "limit_only") else "page_size"
    out = [dict(base, description="defaults: page size 20 is always sent", options={}, expect=wire(shape, 20))]
    if shape["style"] == "limit_only":
        kind = entry.get("size", "page")
        top = 365 if kind == "price_history" else 5000
        out.append(dict(base, description="largest accepted limit", options={size_opt: top}, expect=wire(shape, top)))
        if kind == "price_history":
            out.append(dict(base, description="above the price-history maximum", options={size_opt: 366}, expect_error=True))
        out.append(dict(base, description="negative limit", options={size_opt: -1}, expect_error=True))
        return out
    cont = ({"limit": 7, "offset": 14}, 14) if offset_style else ({"page_size": 7, "cursor": TOKEN}, TOKEN)
    out.append(dict(base, description="continuation page", options=cont[0], expect=wire(shape, 7, cont[1])))
    out.append(dict(base, description="maximum page size", options={size_opt: 100}, expect=wire(shape, 100)))
    out.append(dict(base, description="page size above the maximum", options={size_opt: 101}, expect_error=True))
    out.append(dict(base, description="negative page size", options={size_opt: -1}, expect_error=True))
    if offset_style:
        out.append(dict(base, description="negative offset", options={"offset": -1}, expect_error=True))
    return out


# Filter vectors that pin a decision of the plan, beyond paging.
FILTER_VECTORS = [
    ("WalletService_GetWalletsV2", "exclude_disabled reaches the wire", {"exclude_disabled": True},
     {"query": [["excludeDisabled", "true"], ["limit", "20"]]}),
    ("WalletService_GetAddresses", "exclude_disabled maps to includeDisabledAddresses=exclude", {"exclude_disabled": True},
     {"query": [["includeDisabledAddresses", "exclude"], ["limit", "20"]]}),
    ("PriceService_QueryPricesV2", "from currency alone", {"from_currency_id": "c1"},
     {"body": {"cursor": {"pageSize": "20"}, "from": {"currencyFromId": "c1"}}}),
    ("PriceService_QueryPricesV2", "from currency with target currencies", {"from_currency_id": "c1", "to_currency_ids": ["c2", "c3"]},
     {"body": {"cursor": {"pageSize": "20"}, "fromTo": {"currencyFromId": "c1", "currencyToIds": ["c2", "c3"]}}}),
    ("PriceService_QueryPricesV2", "target currencies alone", {"to_currency_ids": ["c2", "c3"]},
     {"body": {"cursor": {"pageSize": "20"}, "to": {"currencyToIds": ["c2", "c3"]}}}),
    ("RequestService_GetRequestsForApprovalV2", "statuses cannot be applied to the approval queue", {"statuses": ["CREATED"]}, None),
    # The swagger and the controller take externalRequestIDs, but the approval paginator drops it before the query.
    ("RequestService_GetRequestsForApprovalV2", "external_request_ids cannot be applied to the approval queue",
     {"external_request_ids": ["e1"]}, None),
]


def main():
    ops = spec_operations()
    rules = json.loads(RULES.read_text())["operations"]
    methods, vectors = {}, []
    for op_id, entry in rules.items():
        if op_id not in ops:
            raise SystemExit(f"{op_id} is in pagination-vectors.json but not in the shared swagger")
        method, path, op, definitions = ops[op_id]
        shape = paging_shape(op_id, entry["rule"], op, definitions)
        params = path_params(op)
        methods[op_id] = {"method": method, "path": path, "paging": shape, "rule": entry["rule"]}
        vectors.extend(with_fixed(op_id, v) for v in vectors_for(op_id, entry, shape, params))
    for op_id, description, options, expected in FILTER_VECTORS:
        v = {"operation": op_id, "path_params": path_params(ops[op_id][2]), "description": description, "options": options}
        v.update({"expect": json.loads(json.dumps(expected))} if expected else {"expect_error": True})
        vectors.append(with_fixed(op_id, v))
    doc = {
        "_comment": "Generated by scripts/pagination-vectors/generate_requests.py. Query pairs are compared as an "
                    "exact multiset after URL-decoding; bodies structurally. Do not edit by hand.",
        "counts": {"methods": len(methods), "vectors": len(vectors)},
        "methods": methods,
        "vectors": vectors,
    }
    OUT.write_text(json.dumps(doc, indent=2) + "\n")
    print(f"wrote {OUT} {doc['counts']}")


if __name__ == "__main__":
    main()
