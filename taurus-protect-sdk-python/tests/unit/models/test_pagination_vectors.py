"""Cross-SDK pagination vectors: the same reply must give the same page in all four SDKs.

``scripts/resources/pagination-vectors.json`` is generated from the contract in the
repo-root ``CLAUDE.md`` ("Pagination (cross-SDK)"). Every entry runs against the ONE
builder of its family, with the raw reply JSON decoded by the generated reply model
first, so a field-name drift in the generated client fails here too.

Tests are parametrized by the pinned counts, not by the file: a missing or truncated
file fails every test instead of collecting none.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any, Dict

import pytest

from taurus_protect._internal.openapi.models import (
    TgvalidatordGetSignedWhitelistedAddressEnvelopesReply,
    TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply,
    TgvalidatordGetTransactionsReply,
    TgvalidatordGetUsersReply,
    TgvalidatordGetWalletsInfoReply,
    TgvalidatordResponseCursor,
)
from taurus_protect.errors import APIError
from taurus_protect.models.pagination import (
    DEFAULT_PAGE_SIZE,
    LIMIT_ONLY,
    MAX_COUNT,
    MAX_PAGE_SIZE,
    OFFSET_RULES,
    PLUS_LIMIT,
    PLUS_MIN_ROWS_LIMIT,
    PLUS_ROWS,
    PLUS_SERVER_ROWS,
    PRICE_HISTORY_MAX_LIMIT,
    REPLY_OFFSET,
    cursor_page,
    offset_pagination,
    resolve_offset,
    resolve_page_size,
)

VECTORS_PATH = (
    Path(__file__).resolve().parents[4] / "scripts" / "resources" / "pagination-vectors.json"
)

# Bump these together with the file; a mismatch is the signal that the loader is stale.
EXPECTED_COUNTS = {
    "operations": 53,
    "offset": 26,
    "cursor": 12,
    "page_size": 14,
    "offset_input": 4,
}

# One generated reply model per offset rule: the vector's reply JSON is decoded by it.
REPLY_MODEL = {
    REPLY_OFFSET: TgvalidatordGetWalletsInfoReply,
    PLUS_ROWS: TgvalidatordGetTransactionsReply,
    PLUS_MIN_ROWS_LIMIT: TgvalidatordGetUsersReply,
    PLUS_SERVER_ROWS: TgvalidatordGetSignedWhitelistedAddressEnvelopesReply,
    PLUS_LIMIT: TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply,
}

# The pagination family each operation uses in THIS SDK, including the operations it
# does not wrap: the vectors name every paged operation, so a new one fails here until
# it is classified.
SDK_RULES = {
    "ActionService_GetActions": PLUS_ROWS,
    "AssetServiceV2_ListAssetOperationsV2": "cursor",
    "AssetServiceV2_QueryAssetAddressesV2": "cursor",
    "AssetServiceV2_QueryAssetsV2": "cursor",
    "AuditService_GetAuditTrails": "cursor",
    "ChangeService_GetChanges": "cursor",
    "ChangeService_GetChangesForApproval": "cursor",
    "EarnService_GetRewards": "cursor",
    "ExchangeService_GetExchanges": "cursor",
    "FeePayerService_GetFeePayers": PLUS_ROWS,
    "FiatProviderService_GetFiatProviderAccounts": "cursor",
    "FiatProviderService_GetFiatProviderCounterpartyAccounts": "cursor",
    "FiatProviderService_GetFiatProviderEntities": "cursor",
    "FiatProviderService_GetFiatProviderOperations": "cursor",
    "PriceService_ExportPricesHistory": LIMIT_ONLY,
    "PriceService_GetPricesHistory": LIMIT_ONLY,
    "PriceService_QueryPricesV2": "cursor",
    "RequestService_GetRequestsForApprovalV2": "cursor",
    "RequestService_GetRequestsV2": "cursor",
    "RuleService_GetBusinessRulesV2": "cursor",
    "RuleService_GetRulesHistory": "token",
    "StakingService_GetStakeAccounts": "cursor",
    "StatisticsService_GetAggregatedTagStats": "cursor",
    "StatisticsService_GetPortfolioStatisticsHistory": "cursor",
    "TaurusNetworkService_GetLendingAgreements": "cursor",
    "TaurusNetworkService_GetLendingAgreementsForApproval": "cursor",
    "TaurusNetworkService_GetLendingOffers": "cursor",
    "TaurusNetworkService_GetPledgeActions": "cursor",
    "TaurusNetworkService_GetPledgeActionsForApproval": "cursor",
    "TaurusNetworkService_GetPledges": "cursor",
    "TaurusNetworkService_GetPledgesWithdrawals": "cursor",
    "TaurusNetworkService_GetSettlements": "cursor",
    "TaurusNetworkService_GetSettlementsForApproval": "cursor",
    "TaurusNetworkService_GetSharedAddresses": "cursor",
    "TaurusNetworkService_GetSharedAssets": "cursor",
    "TransactionService_ExportTransactions": LIMIT_ONLY,
    "TransactionService_GetTransactions": PLUS_ROWS,
    "UserService_GetGroups": PLUS_MIN_ROWS_LIMIT,
    "UserService_GetUsers": PLUS_MIN_ROWS_LIMIT,
    "WalletService_GetAddresses": REPLY_OFFSET,
    "WalletService_GetAssetAddresses": "cursor",
    "WalletService_GetAssetWallets": "cursor",
    "WalletService_GetBalances": "cursor",
    "WalletService_GetNFTCollectionBalances": "cursor",
    "WalletService_GetReservations": "cursor",
    "WalletService_GetWalletTokens": "token",
    "WalletService_GetWalletsV2": REPLY_OFFSET,
    "WebhookService_GetWebhookCalls": "cursor",
    "WebhookService_GetWebhooks": "cursor",
    "WhitelistService_GetWhitelistedAddresses": PLUS_SERVER_ROWS,
    "WhitelistService_GetWhitelistedAddressesForApproval": PLUS_SERVER_ROWS,
    "WhitelistService_GetWhitelistedContracts": PLUS_LIMIT,
    "WhitelistService_GetWhitelistedContractsForApproval": PLUS_LIMIT,
}


def _load() -> Dict[str, Any]:
    if not VECTORS_PATH.is_file():
        raise AssertionError(f"pagination vectors missing: {VECTORS_PATH}")
    with VECTORS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


try:
    _VECTORS: Dict[str, Any] = _load()
except AssertionError:
    _VECTORS = {}


def _vectors() -> Dict[str, Any]:
    # Re-raises a missing file inside every test rather than at collection time.
    return _VECTORS or _load()


def _ids(section: str) -> list:
    entries = _VECTORS.get(section) or []
    ids = []
    for i in range(EXPECTED_COUNTS[section]):
        label = entries[i]["description"] if i < len(entries) else "missing"
        ids.append(f"{i:02d}-{label}")
    return ids


def test_counts_match_the_file() -> None:
    data = _vectors()
    assert data["counts"] == EXPECTED_COUNTS
    for section, count in EXPECTED_COUNTS.items():
        assert len(data[section]) == count, section


def test_constants_match_the_file() -> None:
    constants = _vectors()["constants"]
    assert constants == {
        "default_page_size": DEFAULT_PAGE_SIZE,
        "max_page_size": MAX_PAGE_SIZE,
        "price_history_max": PRICE_HISTORY_MAX_LIMIT,
        "max_count": MAX_COUNT,
    }


def test_every_operation_is_classified_the_same_way() -> None:
    operations = _vectors()["operations"]
    assert set(operations) == set(SDK_RULES)
    for operation, entry in operations.items():
        assert entry["rule"] == SDK_RULES[operation], operation
        if entry["rule"] in OFFSET_RULES:
            assert entry["total"] is True, operation


@pytest.mark.parametrize("index", range(EXPECTED_COUNTS["offset"]), ids=_ids("offset"))
def test_offset_vector(index: int) -> None:
    vector = _vectors()["offset"][index]
    reply = REPLY_MODEL[vector["rule"]].from_dict(vector["reply"])
    reply_offset = reply.offset if vector["rule"] == REPLY_OFFSET else None

    def build():
        return offset_pagination(
            vector["rule"],
            limit=vector["request"]["limit"],
            offset=vector["request"]["offset"],
            served_rows=vector["served_rows"],
            total_items=reply.total_items,
            reply_offset=reply_offset,
            excluded=vector["sdk_excluded"],
        )

    if vector.get("expect_error"):
        with pytest.raises(APIError):
            build()
        return
    assert build().model_dump() == vector["expect"], vector["description"]


@pytest.mark.parametrize("index", range(EXPECTED_COUNTS["cursor"]), ids=_ids("cursor"))
def test_cursor_vector(index: int) -> None:
    vector = _vectors()["cursor"][index]

    def build():
        if vector["family"] == "token":
            return cursor_page(
                vector["page_size"],
                token=vector["reply_token"],
                total=vector["reply_total"],
                has_total=vector["has_total"],
            )
        return cursor_page(
            vector["page_size"],
            TgvalidatordResponseCursor.from_dict(vector["reply_cursor"]),
            total=vector["reply_total"],
            has_total=vector["has_total"],
        )

    if vector.get("expect_error"):
        with pytest.raises(APIError):
            build()
        return
    assert build().model_dump() == vector["expect"], vector["description"]


_PAGE_SIZE_KINDS = {
    "page": ("page_size", MAX_PAGE_SIZE),
    "price_history": ("limit", PRICE_HISTORY_MAX_LIMIT),
    "export": ("limit", None),
}


@pytest.mark.parametrize("index", range(EXPECTED_COUNTS["page_size"]), ids=_ids("page_size"))
def test_page_size_vector(index: int) -> None:
    vector = _vectors()["page_size"][index]
    name, maximum = _PAGE_SIZE_KINDS[vector["kind"]]

    if vector.get("expect_error"):
        with pytest.raises(ValueError, match=name):
            resolve_page_size(vector["input"], name, maximum)
        return
    assert resolve_page_size(vector["input"], name, maximum) == vector["expect"]


@pytest.mark.parametrize("index", range(EXPECTED_COUNTS["offset_input"]), ids=_ids("offset_input"))
def test_offset_input_vector(index: int) -> None:
    vector = _vectors()["offset_input"][index]

    if vector.get("expect_error"):
        with pytest.raises(ValueError, match="offset"):
            resolve_offset(vector["input"])
        return
    assert resolve_offset(vector["input"]) == vector["expect"]
