"""Asset service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.errors import APIError, IntegrityError, NotFoundError, WhitelistError
from taurus_protect.helpers.address_signature_verifier import verified_address
from taurus_protect.mappers._base import safe_bool, safe_int, safe_string
from taurus_protect.mappers.address import address_from_dto
from taurus_protect.mappers.wallet import wallets_from_dto
from taurus_protect.models.address import Address
from taurus_protect.models.asset import (
    AssetAddressV2,
    AssetOperationV2,
    AssetV2,
    QueryAssetAddressesResult,
)
from taurus_protect.models.blockchain import Asset
from taurus_protect.models.pagination import CursorPage, cursor_page, cursor_request
from taurus_protect.models.wallet import Wallet
from taurus_protect.models.whitelisted_address import ExcludedWhitelistedAddress
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect.services.address_service import AddressService
    from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService

_INTERNAL = "ADDRESS_TYPE_V2_INTERNAL"
_WHITELISTED = "ADDRESS_TYPE_V2_WHITELISTED"


def asset_from_dto(dto: Any) -> Optional[Asset]:
    """
    Convert an OpenAPI asset DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        Asset model or None if dto is None.
    """
    if dto is None:
        return None

    return Asset(
        id=safe_string(getattr(dto, "id", None) or getattr(dto, "currency_id", None)),
        name=getattr(dto, "name", None),
        symbol=getattr(dto, "symbol", None) or getattr(dto, "currency", None),
        blockchain=getattr(dto, "blockchain", None),
        network=getattr(dto, "network", None),
        decimals=safe_int(getattr(dto, "decimals", 0)),
        logo_url=getattr(dto, "logo_url", None) or getattr(dto, "logoUrl", None),
        enabled=safe_bool(getattr(dto, "enabled", True)),
        is_token=safe_bool(getattr(dto, "is_token", False) or getattr(dto, "isToken", False)),
        contract_address=getattr(dto, "contract_address", None)
        or getattr(dto, "contractAddress", None),
    )


def assets_from_dto(dtos: Any) -> List[Asset]:
    """
    Convert a list of OpenAPI asset DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of Asset models.
    """
    if dtos is None:
        return []
    return [a for dto in dtos if (a := asset_from_dto(dto)) is not None]


def _enum_value(value: Any) -> Optional[str]:
    return value.value if value is not None and hasattr(value, "value") else value


def asset_v2_from_dto(dto: Any) -> Optional[AssetV2]:
    """Convert a generated ``TgvalidatordAssetResourceV2`` to the domain model."""
    if dto is None:
        return None

    return AssetV2(
        id=safe_string(dto.id),
        label=dto.label,
        asset_type=dto.asset_type,
        status=_enum_value(dto.status),
        blockchain=dto.blockchain,
        network=dto.network,
        currency_id=dto.currency_id,
        name=dto.name,
        symbol=dto.symbol,
        decimals=dto.decimals,
        contract_address=dto.contract_address,
        attributes={a.key or "": a.value or "" for a in dto.attributes or []},
        version=dto.version,
        created_at=dto.created_at,
        updated_at=dto.updated_at,
    )


def asset_address_v2_from_dto(dto: Any) -> Optional[AssetAddressV2]:
    """Convert a generated ``TgvalidatordAssetAddressV2`` to the domain model."""
    if dto is None:
        return None

    return AssetAddressV2(
        address=dto.address,
        address_id=dto.address_id,
        whitelisted_address_id=dto.whitelisted_address_id,
        address_type=_enum_value(dto.address_type),
        kyc_status=_enum_value(dto.kyc_status),
        balance=dto.balance,
    )


def asset_operation_v2_from_dto(dto: Any) -> Optional[AssetOperationV2]:
    """Convert a generated ``TgvalidatordAssetOperationV2`` to the domain model."""
    if dto is None:
        return None

    return AssetOperationV2(
        id=safe_string(dto.id),
        asset_id=dto.asset_id,
        type=_enum_value(dto.type),
        status=_enum_value(dto.status),
        initiated_by_address_id=dto.initiated_by_address_id,
        failure_reason=_enum_value(dto.failure_reason),
        blocking_reason=_enum_value(dto.blocking_reason),
        created_at=dto.created_at,
        updated_at=dto.updated_at,
    )


class AssetService(BaseService):
    """
    Service for asset operations.

    Provides the wallets and addresses holding an asset, and the v2 asset API
    (asset definitions, their addresses and operations).

    Example:
        >>> # Walk every wallet holding ETH
        >>> cursor = None
        >>> while True:
        ...     wallets, page = client.assets.get_wallets("ETH", page_size=100, cursor=cursor)
        ...     for wallet in wallets:
        ...         print(f"{wallet.name}: {wallet.balance}")
        ...     if not page.has_more:
        ...         break
        ...     cursor = page.next_cursor
    """

    def __init__(
        self,
        api_client: Any,
        assets_api: Any,
        rules_cache: Any,
        assets_v2_api: Any = None,
        *,
        address_service: Optional["AddressService"] = None,
        whitelisted_address_service: Optional["WhitelistedAddressService"] = None,
    ) -> None:
        """
        Initialize asset service.

        Address signature verification is mandatory: get_addresses returns the same
        Address entity AddressService verifies, signature and all, and used to hand back
        the raw generated DTOs unverified — so the mandatory verification there could be
        walked around by asking for the same rows here.

        The v2 asset-holder rows carry no signature, so :meth:`query_asset_addresses`
        confirms them through the two verified readers passed in, the same instances the
        client hands out, rather than through a second copy of either verification.

        Args:
            api_client: The OpenAPI client instance.
            assets_api: The AssetsAPI service from OpenAPI client.
            rules_cache: Rules container cache supplying the HSM key. Required.
            assets_v2_api: The AssetV2Api service; built from ``api_client`` when omitted.
            address_service: The verifying managed-address service. Required.
            whitelisted_address_service: The verifying whitelisted-address service.
                Required.

        Raises:
            ValueError: If rules_cache or either verifying service is None.
        """
        super().__init__(api_client)
        if rules_cache is None:
            raise ValueError(
                "rules_cache cannot be None - address signature verification is mandatory"
            )
        if address_service is None or whitelisted_address_service is None:
            raise ValueError(
                "address_service and whitelisted_address_service cannot be None - "
                "asset holders are verified through them"
            )
        self._address_service = address_service
        self._whitelisted_address_service = whitelisted_address_service
        if assets_v2_api is None:
            from taurus_protect._internal.openapi.api.asset_v2_api import AssetV2Api

            assets_v2_api = AssetV2Api(api_client)
        self._assets_api = assets_api
        self._assets_v2_api = assets_v2_api
        self._rules_cache = rules_cache

    def list(
        self,
        currency: str,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
    ) -> Tuple[List[Wallet], CursorPage]:
        """
        List the wallets holding an asset, one page at a time; same as :meth:`get_wallets`.

        Args:
            currency: The currency ID or symbol (e.g., "ETH").
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.

        Returns:
            Tuple of (wallets, page); the page carries the server's total.

        Raises:
            ValueError: If currency is empty or paging options are invalid.
            APIError: If API request fails.
        """
        return self.get_wallets(currency, page_size=page_size, cursor=cursor)

    def get(self, asset_id: str) -> Asset:
        """
        Get an asset by its currency ID or symbol.

        Read from the currency of the first address holding the asset.

        Args:
            asset_id: The asset ID to retrieve (e.g., "BTC", "ETH").

        Returns:
            The asset.

        Raises:
            ValueError: If asset_id is invalid.
            NotFoundError: If no address holds the asset.
            APIError: If API request fails.
        """
        self._validate_required(asset_id, "asset_id")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_addresses_request import (
                TgvalidatordGetAssetAddressesRequest,
            )

            body = TgvalidatordGetAssetAddressesRequest(
                asset=TgvalidatordAsset(currency=asset_id),
                request_cursor=self._request_cursor_body(cursor_request(1)),
            )
            resp = self._assets_api.wallet_service_get_asset_addresses(body=body)

            rows = resp.addresses or []
            if not rows:
                raise NotFoundError(f"Asset {asset_id} not found")

            first = rows[0]
            info = first.currency_info
            if info is None:
                return Asset(id=asset_id, symbol=first.currency or asset_id)
            return Asset(
                id=info.id or asset_id,
                name=info.name,
                symbol=info.symbol or first.currency or asset_id,
                blockchain=info.blockchain,
                network=info.network,
                decimals=safe_int(info.decimals),
                enabled=bool(info.enabled) if info.enabled is not None else True,
                is_token=bool(info.is_token),
                contract_address=info.contract_address,
            )
        except Exception as e:
            # Consistent with the siblings in this file: IntegrityError is not an
            # APIError, so omitting it here turns a security failure into a retryable
            # ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_wallets(
        self,
        currency: str,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        wallet_id: Optional[str] = None,
        wallet_name: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Wallet], CursorPage]:
        """
        List the wallets holding an asset, one page at a time.

        Args:
            currency: The currency ID or symbol.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            wallet_id: Filter by wallet ID.
            wallet_name: Filter by wallet name.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (wallets, page); the page carries the server's total.

        Raises:
            ValueError: If currency is empty or paging options are invalid.
            APIError: If API request fails.
        """
        self._validate_required(currency, "currency")
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_wallets_request import (
                TgvalidatordGetAssetWalletsRequest,
            )

            body = TgvalidatordGetAssetWalletsRequest(
                asset=TgvalidatordAsset(currency=currency),
                wallet_id=wallet_id,
                wallet_name=wallet_name,
                request_cursor=self._request_cursor_body(req),
            )
            resp = self._assets_api.wallet_service_get_asset_wallets(body=body)

            rows = resp.wallets or []
            wallets = wallets_from_dto(rows)
            page = cursor_page(
                req.page_size,
                resp.cursor,
                total=resp.total_items,
                has_total=True,
                excluded=len(rows) - len(wallets),
            )
            return wallets, page
        except Exception as e:
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_addresses(
        self,
        currency: str,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        wallet_id: Optional[str] = None,
        address_id: Optional[str] = None,
        addresses: Optional[List[str]] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Address], CursorPage]:
        """
        List the addresses holding an asset, one page at a time, each signature verified.

        Args:
            currency: The currency ID or symbol.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            wallet_id: Filter by wallet ID.
            address_id: Filter by address ID.
            addresses: Filter by blockchain addresses.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (verified addresses, page); the page carries the server's total.

        Raises:
            ValueError: If currency is empty or paging options are invalid.
            IntegrityError: If an address signature does not verify.
            APIError: If API request fails.
        """
        self._validate_required(currency, "currency")
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_addresses_request import (
                TgvalidatordGetAssetAddressesRequest,
            )

            body = TgvalidatordGetAssetAddressesRequest(
                asset=TgvalidatordAsset(currency=currency),
                wallet_id=wallet_id,
                address_id=address_id,
                addresses=addresses,
                request_cursor=self._request_cursor_body(req),
            )
            resp = self._assets_api.wallet_service_get_asset_addresses(body=body)

            rows = resp.addresses or []

            # Fail-fast, as AddressService does: one unverifiable address is not a row to
            # skip past when the caller is choosing where funds go. Through the SHARED
            # seam, so the "address string with no signature is withheld" rule cannot
            # differ between this reader and AddressService -- this path was fixed to
            # verify once already while create_address was missed.
            verified: List[Address] = []
            if rows:
                rules_container = self._rules_cache.get_decoded_rules_container()
                for dto in rows:
                    address = address_from_dto(dto)
                    if address is not None:
                        verified.append(verified_address(address, rules_container))

            page = cursor_page(
                req.page_size,
                resp.cursor,
                total=resp.total_items,
                has_total=True,
                excluded=len(rows) - len(verified),
            )
            return verified, page
        except Exception as e:
            # IntegrityError is NOT an APIError, so without naming it here a failed
            # address signature check would be remapped to a retryable ServerError and a
            # caller following isRetryable() would retry a suspected forgery.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def query_assets(
        self,
        *,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        symbol: Optional[str] = None,
        contract_address: Optional[str] = None,
        label: Optional[str] = None,
        currency_name: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[AssetV2], CursorPage]:
        """
        List asset definitions (v2 API), one page at a time.

        Args:
            blockchain: Filter by blockchain.
            network: Filter by network.
            symbol: Filter by symbol.
            contract_address: Filter by contract address.
            label: Filter by label.
            currency_name: Filter by currency name.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (assets, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_get_assets_request_v2 import (
                TgvalidatordGetAssetsRequestV2,
            )

            body = TgvalidatordGetAssetsRequestV2(
                cursor=self._request_cursor_body(req),
                blockchain=blockchain,
                network=network,
                symbol=symbol,
                contract_address=contract_address,
                label=label,
                currency_name=currency_name,
            )
            resp = self._assets_v2_api.asset_service_v2_query_assets_v2(body=body)

            assets = [a for dto in resp.result or [] if (a := asset_v2_from_dto(dto)) is not None]
            return assets, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def query_asset_addresses(
        self,
        asset_id: str,
        *,
        address_type: Optional[str] = None,
        kyc_status: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> QueryAssetAddressesResult:
        """
        List the addresses holding a v2 asset, one page at a time.

        The rows carry no signature, so each INTERNAL and WHITELISTED row is confirmed
        against its verified counterpart and marked ``verified``, with the address taken
        from that verified read; a row that cannot be confirmed is withheld and named in
        ``excluded_unverified``. Every other row (EXTERNAL, untyped) is returned with
        ``verified`` False.

        Args:
            asset_id: The v2 asset ID.
            address_type: Filter by address type (e.g. ``ADDRESS_TYPE_V2_INTERNAL``).
            kyc_status: Filter by KYC status (e.g. ``KYC_STATUS_V2_APPROVED``). Both filters
                are sent as given: a value this SDK does not know is the server's to judge.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            The kept rows, the page and the excluded rows. Exclusions never move the
            cursor.

        Raises:
            ValueError: If asset_id is empty or paging options are invalid.
            IntegrityError: If rows came back and none could be verified, or a verified
                reader failed as a whole.
            APIError: If API request fails.
        """
        self._validate_required(asset_id, "asset_id")
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        from taurus_protect._internal.openapi.models.asset_service_v2_query_asset_addresses_v2_body import (
            AssetServiceV2QueryAssetAddressesV2Body,
        )

        body = AssetServiceV2QueryAssetAddressesV2Body(
            cursor=self._request_cursor_body(req),
            address_type=address_type,
            kyc_status=kyc_status,
        )

        try:
            resp = self._assets_v2_api.asset_service_v2_query_asset_addresses_v2(
                asset_id, body=body
            )
            # Exclusions below never move the cursor.
            page = cursor_page(req.page_size, resp.cursor)

            rows = [
                a for dto in resp.result or [] if (a := asset_address_v2_from_dto(dto)) is not None
            ]
            addresses, excluded = self._verified_asset_addresses(rows)
            return QueryAssetAddressesResult(
                addresses=addresses, page=page, excluded_unverified=excluded
            )
        except Exception as e:
            # WhitelistError too: remapped, a failed whitelist verification would become
            # a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_asset_addresses(
        self, rows: List[AssetAddressV2]
    ) -> Tuple[List[AssetAddressV2], List[ExcludedWhitelistedAddress]]:
        """
        Confirm a page of v2 asset holders through the verified readers::

            INTERNAL    -> address_id             -> verified managed-address read (HSM)
            WHITELISTED -> whitelisted_address_id -> verified whitelist read (6 steps)
                           same address? -- yes -> verified, address from the reader
                                         -- no  -> excluded
            any other type ------------------------> verified False

        A reader is called only when the page has a row of its type. Its call-level
        failures raise; its row-level ones become exclusions here. Rows came back but
        none survived raises, so a filtered page never reads as an empty one.
        """
        internal_ids: Dict[str, None] = {}
        whitelisted_ids: Dict[str, None] = {}
        for row in rows:
            if row.address_type == _INTERNAL and _has_platform_id(row.address_id):
                internal_ids[str(row.address_id)] = None
            elif row.address_type == _WHITELISTED and _has_platform_id(row.whitelisted_address_id):
                whitelisted_ids[str(row.whitelisted_address_id)] = None

        internal: Dict[str, Optional[str]] = {}
        internal_failed: Dict[str, str] = {}
        if internal_ids:
            addresses, internal_failed = self._address_service._verified_addresses_by_id(
                list(internal_ids)
            )
            internal = {i: a.address for i, a in addresses.items()}

        whitelisted: Dict[str, Optional[str]] = {}
        whitelisted_failed: Dict[str, str] = {}
        if whitelisted_ids:
            envelopes, whitelisted_failed = (
                self._whitelisted_address_service._verified_envelopes_by_id(list(whitelisted_ids))
            )
            whitelisted = {
                i: e.verified_whitelisted_address.address
                for i, e in envelopes.items()
                if e.verified_whitelisted_address is not None
            }

        kept: List[AssetAddressV2] = []
        excluded: List[ExcludedWhitelistedAddress] = []
        for row in rows:
            if row.address_type == _INTERNAL:
                exclusion_id, address, reason = _confirm_holder(
                    row.address,
                    row.address_id,
                    "addressID",
                    internal,
                    internal_failed,
                    "managed address",
                )
            elif row.address_type == _WHITELISTED:
                exclusion_id, address, reason = _confirm_holder(
                    row.address,
                    row.whitelisted_address_id,
                    "whitelistedAddressID",
                    whitelisted,
                    whitelisted_failed,
                    "whitelisted address",
                )
            else:
                # An EXTERNAL holder, or an untyped row: on-chain data nothing signs.
                kept.append(row.model_copy(update={"verified": False}))
                continue
            if reason is not None:
                excluded.append(ExcludedWhitelistedAddress(id=exclusion_id, reason=reason))
                continue
            kept.append(row.model_copy(update={"address": address, "verified": True}))

        if rows and not kept:
            raise IntegrityError(
                f"all {len(rows)} asset address(es) failed verification; "
                f"first failure: {excluded[0].reason}"
            )
        return kept, excluded

    def list_asset_operations(
        self,
        asset_id: str,
        *,
        type: Optional[str] = None,
        status: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[AssetOperationV2], CursorPage]:
        """
        List the operations of a v2 asset, one page at a time.

        Args:
            asset_id: The v2 asset ID.
            type: Filter by operation type (e.g. ``ASSET_OPERATION_TYPE_V2_MINT``).
            status: Filter by operation status.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (operations, page).

        Raises:
            ValueError: If asset_id is empty or paging options are invalid.
            APIError: If API request fails.
        """
        self._validate_required(asset_id, "asset_id")
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._assets_v2_api.asset_service_v2_list_asset_operations_v2(
                asset_id,
                type=type,
                status=status,
                **req.query_params(),
            )

            operations = [
                o
                for dto in resp.result or []
                if (o := asset_operation_v2_from_dto(dto)) is not None
            ]
            return operations, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e


def _has_platform_id(platform_id: Optional[str]) -> bool:
    # validatord sends zero, or omits the field, when the address is not managed or not
    # whitelisted.
    return bool(platform_id) and platform_id != "0"


def _confirm_holder(
    row_address: Optional[str],
    platform_id: Optional[str],
    id_field: str,
    verified: Dict[str, Optional[str]],
    failed: Dict[str, str],
    kind: str,
) -> Tuple[Optional[str], Optional[str], Optional[str]]:
    """Check one holder row against its verified reader: (exclusion id, address, reason)."""
    if platform_id is None or not _has_platform_id(platform_id):
        return row_address, None, f"the row carries no {id_field} to verify it by"
    if platform_id not in verified:
        why = failed.get(platform_id)
        if why is not None:
            return platform_id, None, f"the {kind} did not verify: {why}"
        return (
            platform_id,
            None,
            f"the verified {kind} read did not return {id_field} {platform_id}",
        )
    address = verified[platform_id]
    if not address or address != row_address:
        return (
            platform_id,
            None,
            f"the row's address differs from the verified {kind} {platform_id}",
        )
    return platform_id, address, None
