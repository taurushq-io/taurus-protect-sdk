"""Every paged operation, and the Python method that wraps it.

One table serves the cross-SDK list-request vectors and the walk tests: an operationId
maps to how the service is built, how the method is called with Python keywords, which
canonical option names map to which keyword, and how the method returns its rows and
page. Every ``call`` returns ``(rows, page)``.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Callable, Dict, Optional, Tuple
from unittest.mock import MagicMock

from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi import (
    ActionsApi,
    AddressesApi,
    AddressWhitelistingApi,
    AssetsApi,
    AssetV2Api,
    AuditApi,
    BalancesApi,
    BusinessRulesApi,
    ChangesApi,
    ContractWhitelistingApi,
    CurrenciesApi,
    EarnApi,
    ExchangeApi,
    FeePayersApi,
    FiatApi,
    GovernanceRulesApi,
    GroupsApi,
    PricesApi,
    RequestsApi,
    ReservationsApi,
    TaurusNetworkLendingApi,
    TaurusNetworkPledgeApi,
    TaurusNetworkSettlementApi,
    TaurusNetworkSharedAddressAssetApi,
    TransactionsApi,
    UsersApi,
    WalletsApi,
    WebhookCallsApi,
    WebhooksApi,
)
from taurus_protect.models.address import ListAddressesOptions
from taurus_protect.models.audit import ListChangesOptions
from taurus_protect.models.pagination import CursorPage, Pagination
from taurus_protect.models.taurus_network.lending import (
    ListLendingAgreementsOptions,
    ListLendingOffersOptions,
)
from taurus_protect.models.taurus_network.pledge import (
    ListPledgeActionsOptions,
    ListPledgesOptions,
    ListPledgeWithdrawalsOptions,
)
from taurus_protect.models.taurus_network.settlement import (
    ListSettlementsForApprovalOptions,
    ListSettlementsOptions,
)
from taurus_protect.models.taurus_network.sharing import (
    ListSharedAddressesOptions,
    ListSharedAssetsOptions,
)
from taurus_protect.models.transaction import TransactionExport
from taurus_protect.services.action_service import ActionService
from taurus_protect.services.address_service import AddressService
from taurus_protect.services.asset_service import AssetService
from taurus_protect.services.audit_service import AuditService
from taurus_protect.services.balance_service import BalanceService
from taurus_protect.services.business_rule_service import BusinessRuleService
from taurus_protect.services.change_service import ChangeService
from taurus_protect.services.earn_service import EarnService
from taurus_protect.services.exchange_service import ExchangeService
from taurus_protect.services.fee_payer_service import FeePayerService
from taurus_protect.services.fiat_service import FiatService
from taurus_protect.services.governance_rule_service import GovernanceRuleService
from taurus_protect.services.group_service import GroupService
from taurus_protect.services.price_service import PriceService
from taurus_protect.services.request_service import RequestService
from taurus_protect.services.reservation_service import ReservationService
from taurus_protect.services.taurus_network.lending_service import LendingService
from taurus_protect.services.taurus_network.pledge_service import PledgeService
from taurus_protect.services.taurus_network.settlement_service import SettlementService
from taurus_protect.services.taurus_network.sharing_service import SharingService
from taurus_protect.services.transaction_service import TransactionService
from taurus_protect.services.user_service import UserService
from taurus_protect.services.wallet_service import WalletService
from taurus_protect.services.webhook_call_service import WebhookCallService
from taurus_protect.services.webhook_service import WebhookService
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService
from taurus_protect.services.whitelisted_asset_service import WhitelistedAssetService

Rows = Any
Page = Any


@dataclass
class Adapter:
    """How one paged operation is reached from Python."""

    build: Callable[[Any], Any]
    call: Callable[[Any, Dict[str, str], Dict[str, Any]], Tuple[Rows, Page]]
    # canonical option name -> Python keyword
    options: Dict[str, str] = field(default_factory=dict)
    # what the page comes back as: Pagination, CursorPage, TransactionExport or None
    page: Optional[type] = None


def super_admin_keys() -> list:
    return [ec.generate_private_key(ec.SECP256R1()).public_key()]


def _same(*names: str) -> Dict[str, str]:
    return {n: n for n in names}


PAGED = ("limit", "offset")
CURSOR = ("page_size", "cursor")


def _changes(svc: ChangeService, _p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
    result = svc.list(ListChangesOptions(**o))
    return result.changes, result.page


def _changes_for_approval(
    svc: ChangeService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    result = svc.list_for_approval(ListChangesOptions(**o))
    return result.changes, result.page


def _business_rules(
    svc: BusinessRuleService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    result = svc.list(**o)
    return result.rules, result.page


def _rules_history(
    svc: GovernanceRuleService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    result = svc.get_rules_history(**o)
    return result.rules, result.page


def _wla(
    svc: WhitelistedAddressService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    result = svc.list(**o)
    return result.addresses, result.pagination


def _wla_for_approval(
    svc: WhitelistedAddressService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    result = svc.list_for_approval(**o)
    return result.addresses, result.pagination


def _export(svc: TransactionService, _p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
    result = svc.export(**o)
    return result.content, result


def _price_history(svc: PriceService, p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
    return svc.get_historical(p["base"], p["quote"], **o), None


def _address_options(
    svc: AddressService, _p: Dict[str, str], o: Dict[str, Any]
) -> Tuple[Rows, Page]:
    return svc.list_with_options(ListAddressesOptions(**o))


def _options(method: str, options_type: type) -> Callable[..., Tuple[Rows, Page]]:
    def call(svc: Any, _p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
        return getattr(svc, method)(options_type(**o))

    return call


def _kwargs(method: str, *path: str) -> Callable[..., Tuple[Rows, Page]]:
    def call(svc: Any, p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
        return getattr(svc, method)(*(p[name] for name in path), **o)

    return call


def _wallet_tokens(svc: WalletService, p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
    return svc.get_tokens(int(p["id"]), **o)


def _assets(ac: Any) -> AssetService:
    rules_cache = MagicMock()
    return AssetService(
        ac,
        AssetsApi(ac),
        rules_cache,
        AssetV2Api(ac),
        address_service=AddressService(ac, AddressesApi(ac), rules_cache),
        whitelisted_address_service=WhitelistedAddressService(
            ac, AddressWhitelistingApi(ac), super_admin_keys(), 1
        ),
    )


def _asset_holders(svc: AssetService, p: Dict[str, str], o: Dict[str, Any]) -> Tuple[Rows, Page]:
    result = svc.query_asset_addresses(p["assetID"], **o)
    return result.addresses, result.page


def _lending(ac: Any) -> LendingService:
    return LendingService(ac, TaurusNetworkLendingApi(ac))


def _pledges(ac: Any) -> PledgeService:
    return PledgeService(ac, TaurusNetworkPledgeApi(ac))


def _settlements(ac: Any) -> SettlementService:
    return SettlementService(ac, TaurusNetworkSettlementApi(ac))


def _sharing(ac: Any) -> SharingService:
    return SharingService(ac, TaurusNetworkSharedAddressAssetApi(ac))


ADAPTERS: Dict[str, Adapter] = {
    # ---- offset lists
    "ActionService_GetActions": Adapter(
        lambda ac: ActionService(ac, ActionsApi(ac)), _kwargs("list"), _same(*PAGED), Pagination
    ),
    "FeePayerService_GetFeePayers": Adapter(
        lambda ac: FeePayerService(ac, FeePayersApi(ac)), _kwargs("list"), _same(*PAGED), Pagination
    ),
    "TransactionService_GetTransactions": Adapter(
        lambda ac: TransactionService(ac, TransactionsApi(ac)),
        _kwargs("list"),
        _same(*PAGED),
        Pagination,
    ),
    "UserService_GetGroups": Adapter(
        lambda ac: GroupService(ac, GroupsApi(ac)), _kwargs("list"), _same(*PAGED), Pagination
    ),
    "UserService_GetUsers": Adapter(
        lambda ac: UserService(ac, UsersApi(ac)), _kwargs("list"), _same(*PAGED), Pagination
    ),
    "WalletService_GetWalletsV2": Adapter(
        lambda ac: WalletService(ac, WalletsApi(ac)),
        _kwargs("list"),
        _same(*PAGED, "exclude_disabled"),
        Pagination,
    ),
    "WalletService_GetAddresses": Adapter(
        lambda ac: AddressService(ac, AddressesApi(ac), rules_cache=MagicMock()),
        _address_options,
        _same(*PAGED, "exclude_disabled"),
        Pagination,
    ),
    "WhitelistService_GetWhitelistedAddresses": Adapter(
        lambda ac: WhitelistedAddressService(ac, AddressWhitelistingApi(ac), super_admin_keys(), 1),
        _wla,
        _same(*PAGED),
        Pagination,
    ),
    "WhitelistService_GetWhitelistedAddressesForApproval": Adapter(
        lambda ac: WhitelistedAddressService(ac, AddressWhitelistingApi(ac), super_admin_keys(), 1),
        _wla_for_approval,
        _same(*PAGED),
        Pagination,
    ),
    "WhitelistService_GetWhitelistedContracts": Adapter(
        lambda ac: WhitelistedAssetService(ac, ContractWhitelistingApi(ac), super_admin_keys(), 1),
        _kwargs("list"),
        _same(*PAGED),
        Pagination,
    ),
    "WhitelistService_GetWhitelistedContractsForApproval": Adapter(
        lambda ac: WhitelistedAssetService(ac, ContractWhitelistingApi(ac), super_admin_keys(), 1),
        _kwargs("list_for_approval"),
        _same(*PAGED),
        Pagination,
    ),
    # ---- cannot page: one reply bounded by a limit
    "TransactionService_ExportTransactions": Adapter(
        lambda ac: TransactionService(ac, TransactionsApi(ac)),
        _export,
        _same("limit"),
        TransactionExport,
    ),
    "PriceService_GetPricesHistory": Adapter(
        lambda ac: PriceService(ac, PricesApi(ac), rules_cache=MagicMock()),
        _price_history,
        _same("limit"),
    ),
    # ---- cursor lists
    "AssetServiceV2_ListAssetOperationsV2": Adapter(
        _assets, _kwargs("list_asset_operations", "assetID"), _same(*CURSOR), CursorPage
    ),
    "AssetServiceV2_QueryAssetAddressesV2": Adapter(
        _assets, _asset_holders, _same(*CURSOR), CursorPage
    ),
    "AssetServiceV2_QueryAssetsV2": Adapter(
        _assets, _kwargs("query_assets"), _same(*CURSOR), CursorPage
    ),
    "WalletService_GetAssetAddresses": Adapter(
        _assets, _kwargs("get_addresses"), _same(*CURSOR, "currency"), CursorPage
    ),
    "WalletService_GetAssetWallets": Adapter(
        _assets, _kwargs("get_wallets"), _same(*CURSOR, "currency"), CursorPage
    ),
    "AuditService_GetAuditTrails": Adapter(
        lambda ac: AuditService(ac, AuditApi(ac)), _kwargs("list"), _same(*CURSOR), CursorPage
    ),
    "ChangeService_GetChanges": Adapter(
        lambda ac: ChangeService(ac, ChangesApi(ac)), _changes, _same(*CURSOR), CursorPage
    ),
    "ChangeService_GetChangesForApproval": Adapter(
        lambda ac: ChangeService(ac, ChangesApi(ac)),
        _changes_for_approval,
        _same(*CURSOR),
        CursorPage,
    ),
    "EarnService_GetRewards": Adapter(
        lambda ac: EarnService(ac, EarnApi(ac)), _kwargs("list_rewards"), _same(*CURSOR), CursorPage
    ),
    "ExchangeService_GetExchanges": Adapter(
        lambda ac: ExchangeService(ac, ExchangeApi(ac)), _kwargs("list"), _same(*CURSOR), CursorPage
    ),
    "FiatProviderService_GetFiatProviderAccounts": Adapter(
        lambda ac: FiatService(ac, FiatApi(ac), CurrenciesApi(ac)),
        _kwargs("list_fiat_provider_accounts"),
        _same(*CURSOR, "provider", "label"),
        CursorPage,
    ),
    "FiatProviderService_GetFiatProviderEntities": Adapter(
        lambda ac: FiatService(ac, FiatApi(ac), CurrenciesApi(ac)),
        _kwargs("list_fiat_provider_entities"),
        _same(*CURSOR),
        CursorPage,
    ),
    "PriceService_QueryPricesV2": Adapter(
        lambda ac: PriceService(ac, PricesApi(ac), rules_cache=MagicMock()),
        _kwargs("get_current"),
        _same(*CURSOR, "from_currency_id", "to_currency_ids"),
        CursorPage,
    ),
    "RequestService_GetRequestsV2": Adapter(
        lambda ac: RequestService(ac, RequestsApi(ac)),
        _kwargs("list"),
        _same(*CURSOR, "statuses"),
        CursorPage,
    ),
    # The approval queue has no status filter: ``statuses`` is not a parameter of
    # get_for_approval, so the call is refused by name before any request.
    "RequestService_GetRequestsForApprovalV2": Adapter(
        lambda ac: RequestService(ac, RequestsApi(ac)),
        _kwargs("get_for_approval"),
        _same(*CURSOR, "statuses", "external_request_ids"),
        CursorPage,
    ),
    "RuleService_GetBusinessRulesV2": Adapter(
        lambda ac: BusinessRuleService(ac, BusinessRulesApi(ac)),
        _business_rules,
        _same(*CURSOR),
        CursorPage,
    ),
    "WalletService_GetBalances": Adapter(
        lambda ac: BalanceService(ac, BalancesApi(ac)), _kwargs("list"), _same(*CURSOR), CursorPage
    ),
    "WalletService_GetNFTCollectionBalances": Adapter(
        lambda ac: BalanceService(ac, BalancesApi(ac)),
        _kwargs("list_nft_collections"),
        _same(*CURSOR),
        CursorPage,
    ),
    "WalletService_GetReservations": Adapter(
        lambda ac: ReservationService(ac, ReservationsApi(ac)),
        _kwargs("list"),
        _same(*CURSOR),
        CursorPage,
    ),
    "WebhookService_GetWebhookCalls": Adapter(
        lambda ac: WebhookCallService(ac, WebhookCallsApi(ac)),
        _kwargs("list"),
        _same(*CURSOR),
        CursorPage,
    ),
    "WebhookService_GetWebhooks": Adapter(
        lambda ac: WebhookService(ac, WebhooksApi(ac)), _kwargs("list"), _same(*CURSOR), CursorPage
    ),
    "TaurusNetworkService_GetLendingAgreements": Adapter(
        _lending,
        _options("list_lending_agreements", ListLendingAgreementsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetLendingAgreementsForApproval": Adapter(
        _lending,
        _options("list_lending_agreements_for_approval", ListLendingAgreementsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetLendingOffers": Adapter(
        _lending,
        _options("list_lending_offers", ListLendingOffersOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetPledgeActions": Adapter(
        _pledges,
        _options("list_pledge_actions", ListPledgeActionsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetPledgeActionsForApproval": Adapter(
        _pledges,
        _options("list_pledge_actions_for_approval", ListPledgeActionsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetPledges": Adapter(
        _pledges, _options("list_pledges", ListPledgesOptions), _same(*CURSOR), CursorPage
    ),
    "TaurusNetworkService_GetPledgesWithdrawals": Adapter(
        _pledges,
        _options("list_pledge_withdrawals", ListPledgeWithdrawalsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetSettlements": Adapter(
        _settlements,
        _options("list_settlements", ListSettlementsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetSettlementsForApproval": Adapter(
        _settlements,
        _options("list_settlements_for_approval", ListSettlementsForApprovalOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetSharedAddresses": Adapter(
        _sharing,
        _options("list_shared_addresses", ListSharedAddressesOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    "TaurusNetworkService_GetSharedAssets": Adapter(
        _sharing,
        _options("list_shared_assets", ListSharedAssetsOptions),
        _same(*CURSOR),
        CursorPage,
    ),
    # ---- token-only lists
    "RuleService_GetRulesHistory": Adapter(
        lambda ac: GovernanceRuleService(ac, GovernanceRulesApi(ac), super_admin_keys(), 1),
        _rules_history,
        _same(*CURSOR),
        CursorPage,
    ),
    "WalletService_GetWalletTokens": Adapter(
        lambda ac: WalletService(ac, WalletsApi(ac)), _wallet_tokens, _same(*CURSOR), CursorPage
    ),
}

NOT_WRAPPED: Dict[str, str] = {
    "FiatProviderService_GetFiatProviderCounterpartyAccounts": "no Python method wraps it",
    "FiatProviderService_GetFiatProviderOperations": "no Python method wraps it",
    "PriceService_ExportPricesHistory": "no Python method wraps it",
    "StatisticsService_GetAggregatedTagStats": "no Python method wraps it",
    "StatisticsService_GetPortfolioStatisticsHistory": "no Python method wraps it",
    "StakingService_GetStakeAccounts": (
        "no Python stake-accounts list; StakingService.get_staking_info reads the first "
        "row and asks for a page of 1"
    ),
}
