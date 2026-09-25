#!/usr/bin/env python3
"""Generate scripts/resources/pagination-vectors.json from the reference rules below.

The rules ARE the cross-SDK contract (repo-root CLAUDE.md § "Pagination (cross-SDK)"); every SDK's unit
suite loads the generated file and asserts its own helpers produce the same values. Edit the rules or the
cases here, never the JSON by hand, then bump the count constant in all four loaders.
"""
import json
import re
from pathlib import Path

OUT = Path(__file__).resolve().parents[1] / "resources" / "pagination-vectors.json"
DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE, PRICE_HISTORY_MAX = 20, 100, 365
MAX_COUNT = 2**53 - 1  # the largest integer TypeScript represents exactly; parity across all four SDKs
COUNT = re.compile(r"0|[1-9][0-9]*")


class Invalid(Exception):
    pass


def parse_count(value):
    """A server count: absent means 0 (zero values are omitted on the wire)."""
    if value is None:
        return 0
    if not isinstance(value, str) or not COUNT.fullmatch(value) or int(value) > MAX_COUNT:
        raise Invalid(value)
    return int(value)


def offset_page(rule, limit, offset, rows, excluded, total, reply_offset):
    server_total = parse_count(total)
    if rule == "reply_offset":
        next_offset = parse_count(reply_offset) if reply_offset is not None else offset + rows
    elif rule in ("plus_rows", "plus_server_rows"):
        next_offset = offset + rows
    elif rule == "plus_min_rows_limit":
        next_offset = offset + min(rows, limit)
    elif rule == "plus_limit":
        next_offset = offset + limit
    else:
        raise ValueError(rule)
    return {
        "limit": limit,
        "offset": offset,
        "total_items": max(0, server_total - excluded),
        "next_offset": next_offset,
        "has_more": offset < next_offset < server_total,
    }


def cursor_page(page_size, cursor, total, has_total):
    """cursor{currentPage, hasNext} families (v1, v2, requestCursor, body cursor)."""
    has_more = bool(cursor and cursor.get("hasNext"))
    current = (cursor or {}).get("currentPage") or ""
    if has_more and not current:
        raise Invalid("hasNext without currentPage")
    return {
        "page_size": page_size,
        "next_cursor": current if has_more else "",
        "has_more": has_more,
        "total_items": parse_count(total) if has_total else None,
    }


def token_page(page_size, token, total, has_total):
    """Token-only operations: the reply token is the next cursor."""
    return {
        "page_size": page_size,
        "next_cursor": token or "",
        "has_more": bool(token),
        "total_items": parse_count(total) if has_total else None,
    }


def resolve_size(kind, value):
    maximum = {"page": MAX_PAGE_SIZE, "price_history": PRICE_HISTORY_MAX, "export": None}[kind]
    if value is None or value == 0:
        return DEFAULT_PAGE_SIZE
    if value < 0 or (maximum is not None and value > maximum):
        raise Invalid(value)
    return value


def resolve_offset(value):
    if value is None:
        return 0
    if value < 0:
        raise Invalid(value)
    return value


def expect(fn, *args):
    try:
        return {"expect": fn(*args)}
    except Invalid:
        return {"expect_error": True}


OFFSET_RULE_OPS = {
    "WalletService_GetWalletsV2": "reply_offset",
    "WalletService_GetAddresses": "reply_offset",
    "TransactionService_GetTransactions": "plus_rows",
    "FeePayerService_GetFeePayers": "plus_rows",
    "ActionService_GetActions": "plus_rows",
    "UserService_GetUsers": "plus_min_rows_limit",
    "UserService_GetGroups": "plus_min_rows_limit",
    "WhitelistService_GetWhitelistedAddresses": "plus_server_rows",
    "WhitelistService_GetWhitelistedAddressesForApproval": "plus_server_rows",
    "WhitelistService_GetWhitelistedContracts": "plus_limit",
    "WhitelistService_GetWhitelistedContractsForApproval": "plus_limit",
}
CURSOR_OPS = [
    "RequestService_GetRequestsV2", "RequestService_GetRequestsForApprovalV2",
    "ChangeService_GetChanges", "ChangeService_GetChangesForApproval",
    "AuditService_GetAuditTrails", "RuleService_GetBusinessRulesV2", "StakingService_GetStakeAccounts",
    "TaurusNetworkService_GetSharedAddresses", "TaurusNetworkService_GetSharedAssets",
    "TaurusNetworkService_GetLendingAgreements", "TaurusNetworkService_GetLendingAgreementsForApproval",
    "TaurusNetworkService_GetLendingOffers", "TaurusNetworkService_GetPledges",
    "TaurusNetworkService_GetPledgeActions", "TaurusNetworkService_GetPledgeActionsForApproval",
    "TaurusNetworkService_GetPledgesWithdrawals", "TaurusNetworkService_GetSettlements",
    "TaurusNetworkService_GetSettlementsForApproval",
    "WalletService_GetBalances", "WalletService_GetAssetAddresses", "WalletService_GetAssetWallets",
    "WalletService_GetNFTCollectionBalances", "WalletService_GetReservations", "ExchangeService_GetExchanges",
    "FiatProviderService_GetFiatProviderAccounts", "FiatProviderService_GetFiatProviderCounterpartyAccounts",
    "FiatProviderService_GetFiatProviderOperations", "FiatProviderService_GetFiatProviderEntities",
    "StatisticsService_GetAggregatedTagStats", "StatisticsService_GetPortfolioStatisticsHistory",
    "WebhookService_GetWebhooks", "WebhookService_GetWebhookCalls", "PriceService_QueryPricesV2",
    "AssetServiceV2_QueryAssetsV2", "AssetServiceV2_QueryAssetAddressesV2", "AssetServiceV2_ListAssetOperationsV2",
    "EarnService_GetRewards",
]
TOKEN_OPS = ["WalletService_GetWalletTokens", "RuleService_GetRulesHistory"]
# Operations that cannot page: (size kind, reply carries a count). The transaction export takes an `offset`
# on the wire, but validatord ignores it and always exports from the first matching row.
LIMIT_ONLY_OPS = {
    "PriceService_GetPricesHistory": ("price_history", False),
    "PriceService_ExportPricesHistory": ("export", False),
    "TransactionService_ExportTransactions": ("export", True),
}
# Cursor operations whose reply carries a count (absent on the wire means 0).
TOTAL_OPS = {"WalletService_GetBalances", "WalletService_GetWalletTokens", "RuleService_GetRulesHistory",
             "WalletService_GetAssetAddresses", "WalletService_GetAssetWallets"}


def operations():
    ops = {op: {"rule": rule, "total": True} for op, rule in OFFSET_RULE_OPS.items()}
    ops.update({op: {"rule": "cursor", "total": op in TOTAL_OPS} for op in CURSOR_OPS})
    ops.update({op: {"rule": "token", "total": op in TOTAL_OPS} for op in TOKEN_OPS})
    ops.update({op: {"rule": "limit_only", "size": kind, "total": total} for op, (kind, total) in LIMIT_ONLY_OPS.items()})
    return dict(sorted(ops.items()))

# (description, rule, limit, offset, server rows, sdk-excluded rows, reply totalItems, reply offset)
OFFSET_CASES = [
    ("empty page {} (zero values omitted)", "reply_offset", 20, 0, 0, 0, None, None),
    ("reply offset IS the next offset (the old current-offset reading dropped the last page)", "reply_offset", 20, 0, 20, 0, "25", "20"),
    ("reply offset: last page", "reply_offset", 20, 20, 5, 0, "25", "25"),
    ("reply offset: exact multiple, last full page", "reply_offset", 20, 20, 20, 0, "40", "40"),
    ("reply offset absent although rows came back: fall back to offset + rows", "reply_offset", 20, 0, 3, 0, "3", None),
    ("reply offset: offset past the total", "reply_offset", 20, 60, 0, 0, "40", "60"),
    ("plus_rows: first of three pages", "plus_rows", 20, 0, 20, 0, "45", None),
    ("plus_rows: middle page", "plus_rows", 20, 20, 20, 0, "45", None),
    ("plus_rows: last page", "plus_rows", 20, 40, 5, 0, "45", None),
    ("plus_rows: upper-bound total, trailing empty page makes no progress", "plus_rows", 20, 40, 0, 0, "45", None),
    ("plus_rows: empty page {}", "plus_rows", 20, 0, 0, 0, None, None),
    ("plus_min_rows_limit: synthetic row appended beyond limit", "plus_min_rows_limit", 20, 0, 21, 0, "30", None),
    ("plus_min_rows_limit: short last page with synthetic row", "plus_min_rows_limit", 20, 20, 11, 0, "30", None),
    ("plus_server_rows: SDK excluded two rows, next offset stays in server space", "plus_server_rows", 20, 0, 20, 2, "30", None),
    ("plus_server_rows: server dropped a bad-signature row mid-page", "plus_server_rows", 20, 0, 19, 0, "30", None),
    ("plus_server_rows: exclusions exceed the total, total clamps at 0", "plus_server_rows", 20, 0, 3, 5, "3", None),
    ("plus_limit: skipped rows keep their slot, short page is not the end", "plus_limit", 20, 0, 17, 0, "45", None),
    ("plus_limit: last page", "plus_limit", 20, 40, 5, 0, "45", None),
    ("plus_limit: empty page {}", "plus_limit", 20, 0, 0, 0, None, None),
    ("count at the largest exact integer", "plus_rows", 20, 0, 20, 0, "9007199254740991", None),
    ("count above 2^53 - 1 is rejected", "plus_rows", 20, 0, 20, 0, "9007199254740992", None),
    ("non-numeric count is rejected", "plus_rows", 20, 0, 20, 0, "abc", None),
    ("negative count is rejected", "plus_rows", 20, 0, 20, 0, "-1", None),
    ("empty-string count is rejected", "plus_rows", 20, 0, 20, 0, "", None),
    ("non-canonical count is rejected", "plus_rows", 20, 0, 20, 0, "020", None),
    ("non-numeric reply offset is rejected", "reply_offset", 20, 0, 20, 0, "25", "x"),
]

# (description, family, page size, reply cursor object or token, reply total, has_total)
CURSOR_CASES = [
    ("v2 middle page", "cursor", 20, {"currentPage": "g2dg-Wq_bQ", "hasNext": True}, None, False),
    ("v2 last page: currentPage present, hasNext omitted", "cursor", 20, {"currentPage": "g2dg-Wq_bQ"}, None, False),
    ("empty page {}: no cursor object", "cursor", 20, None, None, False),
    ("cursor object present but empty", "cursor", 20, {}, None, False),
    ("v1 standard-base64 token with + / = survives verbatim", "cursor", 20, {"currentPage": "eyJhIjoiKy8ifQ+/cQ==", "hasNext": True}, None, False),
    ("hasNext without currentPage is rejected", "cursor", 20, {"hasNext": True}, None, False),
    ("request-cursor family with a total", "cursor", 20, {"currentPage": "abc", "hasNext": True}, "57", True),
    ("request-cursor family: zero total omitted", "cursor", 20, None, None, True),
    ("token: next page", "token", 20, "eyJpZCI6NDF9", "120", True),
    ("token: last page (token omitted)", "token", 20, None, "120", True),
    ("token with + / = survives verbatim", "token", 20, "ab+/cQ==", None, True),
    ("token: non-numeric total is rejected", "token", 20, "YWJj", "x", True),
]

# (description, kind, input)
SIZE_CASES = [
    ("unset → default", "page", None), ("zero → default", "page", 0), ("minimum", "page", 1),
    ("default", "page", 20), ("maximum", "page", 100), ("above maximum", "page", 101), ("negative", "page", -1),
    ("price history unset → default", "price_history", None), ("price history maximum", "price_history", 365),
    ("price history above maximum", "price_history", 366), ("price history negative", "price_history", -1),
    ("export unset → default", "export", None), ("export has no SDK maximum", "export", 5000), ("export negative", "export", -1),
]
OFFSET_INPUT_CASES = [("unset → 0", None), ("zero", 0), ("positive", 40), ("negative", -1)]


def main():
    offset = [
        dict(description=d, rule=r, request={"limit": l, "offset": o}, served_rows=n, sdk_excluded=x,
             reply={k: v for k, v in (("totalItems", t), ("offset", ro)) if v is not None},
             **expect(offset_page, r, l, o, n, x, t, ro))
        for d, r, l, o, n, x, t, ro in OFFSET_CASES
    ]
    cursor = []
    for d, fam, size, reply, total, has_total in CURSOR_CASES:
        entry = dict(description=d, family=fam, page_size=size, has_total=has_total)
        if fam == "cursor":
            entry["reply_cursor"] = reply
            entry.update(expect(cursor_page, size, reply, total, has_total))
        else:
            entry["reply_token"] = reply
            entry.update(expect(token_page, size, reply, total, has_total))
        entry["reply_total"] = total
        cursor.append(entry)
    sizes = [dict(description=d, kind=k, input=v, **expect(resolve_size, k, v)) for d, k, v in SIZE_CASES]
    offsets = [dict(description=d, input=v, **expect(resolve_offset, v)) for d, v in OFFSET_INPUT_CASES]
    doc = {
        "_comment": "Generated by scripts/pagination-vectors/generate.py from the rules in repo-root CLAUDE.md "
                    "§ Pagination (cross-SDK). Do not edit by hand; regenerate and bump every loader's counts.",
        "constants": {"default_page_size": DEFAULT_PAGE_SIZE, "max_page_size": MAX_PAGE_SIZE,
                      "price_history_max": PRICE_HISTORY_MAX, "max_count": MAX_COUNT},
        "counts": {"operations": len(operations()), "offset": len(offset), "cursor": len(cursor),
                   "page_size": len(sizes), "offset_input": len(offsets)},
        "operations": operations(),
        "offset": offset,
        "cursor": cursor,
        "page_size": sizes,
        "offset_input": offsets,
    }
    OUT.write_text(json.dumps(doc, indent=2) + "\n")
    print(f"wrote {OUT} {doc['counts']}")


if __name__ == "__main__":
    main()
