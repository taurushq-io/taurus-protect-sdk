"""Two-page walks over every paged operation this SDK wraps.

Each walk runs through the transport stub, so the rows come out of the generated reply
model's real field. A list that reads a field the model does not have returns no rows
and fails here; so does a continuation that drops the cursor, the ``NEXT`` request or
the page size, and a walk that sends anything after ``has_more`` turned false.
"""

from __future__ import annotations

import base64
import json
from pathlib import Path
from typing import Any, Dict, List
from unittest.mock import patch

import pytest

from taurus_protect.models.pagination import CursorPage, Pagination
from taurus_protect.models.whitelisted_address import (
    SignedWhitelistedAddressEnvelope,
    WhitelistedAddress,
)
from tests.unit.services.list_adapters import ADAPTERS
from tests.unit.transport_stub import RecordedRequest, StubTransport, api_client

RESOURCES = Path(__file__).resolve().parents[4] / "scripts" / "resources"

# Standard base64 carrying all three characters that must be escaped on the wire.
TOKEN = base64.b64encode(b"\xfb\xef\xbe\xff").decode("ascii")
assert "+" in TOKEN and "/" in TOKEN and "=" in TOKEN

PAGE_SIZE = 2

# Per operation: the JSON key the generated reply model reads rows from, one minimal row
# its mapper accepts, and any options or path parameters the operation requires.
WALKS: Dict[str, Dict[str, Any]] = {
    # offset
    "ActionService_GetActions": {"rows": "result", "row": {"id": "1"}},
    "FeePayerService_GetFeePayers": {"rows": "result", "row": {"id": "1"}},
    "TransactionService_GetTransactions": {"rows": "result", "row": {"id": "1"}},
    "UserService_GetGroups": {"rows": "result", "row": {"id": "1"}},
    "UserService_GetUsers": {"rows": "result", "row": {"id": "1"}},
    "WalletService_GetWalletsV2": {"rows": "result", "row": {"id": "1"}},
    "WalletService_GetAddresses": {"rows": "result", "row": {"id": "1", "walletId": "7"}},
    "WhitelistService_GetWhitelistedAddresses": {"rows": "result", "row": {"id": "1"}},
    "WhitelistService_GetWhitelistedAddressesForApproval": {"rows": "result", "row": {"id": "1"}},
    "WhitelistService_GetWhitelistedContracts": {"rows": "result", "row": {"id": "1"}},
    "WhitelistService_GetWhitelistedContractsForApproval": {"rows": "result", "row": {"id": "1"}},
    # cursor
    "AssetServiceV2_ListAssetOperationsV2": {
        "rows": "result",
        "row": {"id": "op1"},
        "path": {"assetID": "a1"},
    },
    "AssetServiceV2_QueryAssetAddressesV2": {
        "rows": "result",
        "row": {"address": "addr", "addressID": "1"},
        "path": {"assetID": "a1"},
    },
    "AssetServiceV2_QueryAssetsV2": {"rows": "result", "row": {"id": "a1"}},
    "WalletService_GetAssetAddresses": {
        "rows": "addresses",
        "row": {"id": "1", "walletId": "7"},
        "options": {"currency": "ETH"},
    },
    "WalletService_GetAssetWallets": {
        "rows": "wallets",
        "row": {"id": "1"},
        "options": {"currency": "ETH"},
    },
    "AuditService_GetAuditTrails": {"rows": "result", "row": {"id": "1"}},
    "ChangeService_GetChanges": {"rows": "result", "row": {"id": "1"}},
    "ChangeService_GetChangesForApproval": {"rows": "result", "row": {"id": "1"}},
    "EarnService_GetRewards": {"rows": "rewards", "row": {"id": "1"}},
    "ExchangeService_GetExchanges": {"rows": "result", "row": {"id": "1"}},
    "FiatProviderService_GetFiatProviderAccounts": {
        "rows": "result",
        "row": {"id": "1"},
        "options": {"provider": "p", "label": "l"},
    },
    "FiatProviderService_GetFiatProviderEntities": {"rows": "result", "row": {"id": "1"}},
    "PriceService_QueryPricesV2": {
        "rows": "result",
        "row": {"currencyFrom": "BTC", "currencyTo": "USD", "rate": "1"},
    },
    "RequestService_GetRequestsV2": {"rows": "result", "row": {"id": "1"}},
    "RequestService_GetRequestsForApprovalV2": {"rows": "result", "row": {"id": "1"}},
    "RuleService_GetBusinessRulesV2": {"rows": "result", "row": {"id": "1"}},
    "WalletService_GetBalances": {"rows": "balances", "row": {"asset": {"currency": "ETH"}}},
    "WalletService_GetNFTCollectionBalances": {"rows": "balances", "row": {"name": "c"}},
    "WalletService_GetReservations": {"rows": "result", "row": {"id": "1"}},
    "WebhookService_GetWebhookCalls": {"rows": "calls", "row": {"id": "1"}},
    "WebhookService_GetWebhooks": {"rows": "webhooks", "row": {"id": "1"}},
    "TaurusNetworkService_GetLendingAgreements": {"rows": "lendingAgreements", "row": {"id": "1"}},
    "TaurusNetworkService_GetLendingAgreementsForApproval": {"rows": "result", "row": {"id": "1"}},
    "TaurusNetworkService_GetLendingOffers": {"rows": "lendingOffers", "row": {"id": "1"}},
    "TaurusNetworkService_GetPledgeActions": {"rows": "result", "row": {"id": "1"}},
    "TaurusNetworkService_GetPledgeActionsForApproval": {"rows": "result", "row": {"id": "1"}},
    "TaurusNetworkService_GetPledges": {"rows": "pledges", "row": {"id": "1"}},
    "TaurusNetworkService_GetPledgesWithdrawals": {"rows": "withdrawals", "row": {"id": "1"}},
    "TaurusNetworkService_GetSettlements": {"rows": "result", "row": {"id": "1"}},
    "TaurusNetworkService_GetSettlementsForApproval": {"rows": "result", "row": {"id": "1"}},
    "TaurusNetworkService_GetSharedAddresses": {"rows": "sharedAddresses", "row": {"id": "1"}},
    "TaurusNetworkService_GetSharedAssets": {"rows": "sharedAssets", "row": {"id": "1"}},
    # token
    "WalletService_GetWalletTokens": {
        "rows": "balances",
        "row": {"asset": {"currency": "ETH"}},
        "path": {"id": "7"},
    },
    # History rows are withheld unless SuperAdmin-signed; the walk counts rows+withheld.
    "RuleService_GetRulesHistory": {"rows": "result", "row": {"rulesContainer": "e30="}},
}


def _load(name: str) -> Dict[str, Any]:
    with (RESOURCES / name).open(encoding="utf-8") as handle:
        return json.load(handle)


METHODS = _load("list-request-vectors.json")["methods"]
OPERATIONS = _load("pagination-vectors.json")["operations"]

PAGED = sorted(op for op, a in ADAPTERS.items() if a.page in (Pagination, CursorPage))


def test_every_paged_adapter_has_a_walk() -> None:
    assert set(PAGED) == set(WALKS)


def _unverified_rows() -> Any:
    """Whitelist rows need a signed envelope chain; for the walk, accept them as mapped."""

    def verified_addresses(self: Any, rows: List[Any], _cache: Any) -> Any:
        envelopes = [
            SignedWhitelistedAddressEnvelope(
                id=str(dto.id),
                verified_whitelisted_address=WhitelistedAddress(id=str(dto.id)),
            )
            for dto in rows
        ]
        return envelopes, []

    return [
        patch(
            "taurus_protect.services.whitelisted_address_service."
            "WhitelistedAddressService._verified_addresses",
            verified_addresses,
        ),
        patch(
            "taurus_protect.services.whitelisted_asset_service."
            "WhitelistedAssetService._verified_asset",
            lambda self, asset, dto=None: asset,
        ),
    ]


def _cursor_reply(op: str, rows: List[Any], current_page: str, has_next: bool) -> Dict[str, Any]:
    walk = WALKS[op]
    style = METHODS[op]["paging"]["style"]
    reply: Dict[str, Any] = {walk["rows"]: rows}
    if style == "token":
        if has_next:
            reply["next" if op == "WalletService_GetWalletTokens" else "cursor"] = current_page
    else:
        cursor: Dict[str, Any] = {"currentPage": current_page}
        if has_next:
            cursor["hasNext"] = True
        reply["cursor"] = cursor
    if OPERATIONS[op]["total"]:
        key = (
            "total"
            if op in ("WalletService_GetBalances", "WalletService_GetWalletTokens")
            else "totalItems"
        )
        reply[key] = "3"
    return reply


def _assert_first_page_request(op: str, request: RecordedRequest) -> None:
    paging = METHODS[op]["paging"]
    style = paging["style"]
    if style == "query_cursor":
        prefix = paging["prefix"]
        assert request.param(f"{prefix}.pageSize") == str(PAGE_SIZE)
        assert request.param(f"{prefix}.currentPage") is None
        assert request.param(f"{prefix}.pageRequest") is None
    elif style == "body_cursor":
        assert request.body[paging["field"]] == {"pageSize": str(PAGE_SIZE)}
    elif style == "token":
        assert request.param(paging["size"]) == str(PAGE_SIZE)
        assert request.param(paging["token"]) is None
    else:
        assert request.param("limit") == str(PAGE_SIZE)
        assert request.param("offset") is None


def _assert_continuation_request(op: str, request: RecordedRequest) -> None:
    paging = METHODS[op]["paging"]
    style = paging["style"]
    if style == "query_cursor":
        prefix = paging["prefix"]
        assert request.param(f"{prefix}.currentPage") == TOKEN
        assert request.param(f"{prefix}.pageRequest") == "NEXT"
        assert request.param(f"{prefix}.pageSize") == str(PAGE_SIZE)
        assert "+" not in request.raw_query and "%2B" in request.raw_query
        assert "%3D" in request.raw_query
    elif style == "body_cursor":
        assert request.body[paging["field"]] == {
            "pageSize": str(PAGE_SIZE),
            "currentPage": TOKEN,
            "pageRequest": "NEXT",
        }
    elif style == "token":
        assert request.param(paging["token"]) == TOKEN
        assert request.param(paging["size"]) == str(PAGE_SIZE)
        assert "+" not in request.raw_query and "%2B" in request.raw_query
    else:
        assert request.param("limit") == str(PAGE_SIZE)
        assert request.param("offset") == str(PAGE_SIZE)


@pytest.mark.parametrize("op", PAGED)
def test_two_page_walk(op: str) -> None:
    adapter = ADAPTERS[op]
    walk = WALKS[op]
    row = walk["row"]
    options = dict(walk.get("options", {}))
    path = walk.get("path", {})
    is_offset = adapter.page is Pagination

    if is_offset:
        first_reply: Dict[str, Any] = {walk["rows"]: [row, row], "totalItems": "3"}
        last_reply: Dict[str, Any] = {walk["rows"]: [row], "totalItems": "3"}
        if OPERATIONS[op]["rule"] == "reply_offset":
            first_reply["offset"] = "2"
            last_reply["offset"] = "3"
    else:
        first_reply = _cursor_reply(op, [row, row], TOKEN, True)
        last_reply = _cursor_reply(op, [row], "bGFzdA==", False)

    ac = api_client()
    service = adapter.build(ac)
    patches = _unverified_rows()
    for p in patches:
        p.start()
    try:
        with StubTransport(first_reply, last_reply) as transport:
            size_key = "limit" if is_offset else "page_size"
            rows, page = _call(op, service, path, {**options, size_key: PAGE_SIZE})
            assert len(rows) == 2, "rows read from the reply field"
            assert page.has_more is True
            if is_offset:
                assert page.next_offset == PAGE_SIZE
                continuation = {**options, "limit": PAGE_SIZE, "offset": page.next_offset}
            else:
                assert page.next_cursor == TOKEN
                assert page.page_size == PAGE_SIZE
                continuation = {**options, "page_size": PAGE_SIZE, "cursor": page.next_cursor}

            rows, page = _call(op, service, path, continuation)
            assert page.has_more is False
            if not is_offset:
                assert page.next_cursor == ""
    finally:
        for p in patches:
            p.stop()

    assert len(transport.requests) == 2
    _assert_first_page_request(op, transport.requests[0])
    _assert_continuation_request(op, transport.requests[1])


def _call(op: str, service: Any, path: Dict[str, str], options: Dict[str, Any]) -> Any:
    if op == "RuleService_GetRulesHistory":
        result = service.get_rules_history(**options)
        # Unsigned history rows are withheld and reported, not dropped: count both.
        return [*result.rules, *result.excluded_unverified], result.page
    return ADAPTERS[op].call(service, path, options)
