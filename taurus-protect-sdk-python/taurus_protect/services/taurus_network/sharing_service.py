"""Sharing service for Taurus Network shared address/asset operations."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.services._base import BaseService
from taurus_protect.services.taurus_network.settlement_service import CursorPagination

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


@dataclass
class SharedAddressTrail:
    """
    Trail entry for a shared address status change.

    Attributes:
        status: The status at this point in the trail.
        changed_at: When the status changed.
    """

    status: str = ""
    changed_at: Optional[datetime] = None


@dataclass
class SharedAddress:
    """
    A shared address in Taurus Network.

    Represents an address shared between two participants for settlement
    and collateral operations.

    Attributes:
        id: The shared address ID.
        owner_participant_id: The participant who owns/shared the address.
        target_participant_id: The participant the address is shared with.
        address_id: The underlying address ID.
        address: The blockchain address string.
        blockchain: The blockchain type (e.g., ETH, BTC).
        network: The network (e.g., mainnet, testnet).
        status: Current sharing status.
        key_value_attributes: Key-value attributes attached to the shared address.
        trail: Status change trail.
        created_at: When the sharing was created.
        updated_at: When the sharing was last updated.
    """

    id: str = ""
    owner_participant_id: str = ""
    target_participant_id: str = ""
    address_id: str = ""
    address: str = ""
    blockchain: str = ""
    network: str = ""
    status: str = ""
    key_value_attributes: List[Dict[str, str]] = field(default_factory=list)
    trail: List[SharedAddressTrail] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class SharedAsset:
    """
    A shared whitelisted asset in Taurus Network.

    Represents a whitelisted contract/asset shared between two participants.

    Attributes:
        id: The shared asset ID.
        owner_participant_id: The participant who owns/shared the asset.
        target_participant_id: The participant the asset is shared with.
        whitelisted_contract_id: The underlying whitelisted contract ID.
        blockchain: The blockchain type.
        network: The network.
        contract_address: The contract address.
        status: Current sharing status.
        created_at: When the sharing was created.
        updated_at: When the sharing was last updated.
    """

    id: str = ""
    owner_participant_id: str = ""
    target_participant_id: str = ""
    whitelisted_contract_id: str = ""
    blockchain: str = ""
    network: str = ""
    contract_address: str = ""
    status: str = ""
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class ListSharedAddressesOptions:
    """
    Options for listing shared addresses.

    Attributes:
        participant_id: Filter by participant (owner or target).
        owner_participant_id: Filter by owner participant.
        target_participant_id: Filter by target participant.
        blockchain: Filter by blockchain.
        network: Filter by network (requires blockchain).
        ids: Filter by specific shared address IDs.
        statuses: Filter by statuses (new, pending, rejected, accepted, unshared).
        sort_order: Sort order (ASC or DESC).
        page_size: Number of items per page.
        current_page: Current page cursor (base64).
        page_request: Page request direction (FIRST, PREVIOUS, NEXT, LAST).
    """

    participant_id: Optional[str] = None
    owner_participant_id: Optional[str] = None
    target_participant_id: Optional[str] = None
    blockchain: Optional[str] = None
    network: Optional[str] = None
    ids: Optional[List[str]] = None
    statuses: Optional[List[str]] = None
    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class ListSharedAssetsOptions:
    """
    Options for listing shared assets.

    Attributes:
        participant_id: Filter by participant (owner or target).
        owner_participant_id: Filter by owner participant.
        target_participant_id: Filter by target participant.
        blockchain: Filter by blockchain.
        network: Filter by network (requires blockchain).
        ids: Filter by specific shared asset IDs.
        statuses: Filter by statuses.
        sort_order: Sort order (ASC or DESC).
        page_size: Number of items per page.
        current_page: Current page cursor (base64).
        page_request: Page request direction (FIRST, PREVIOUS, NEXT, LAST).
    """

    participant_id: Optional[str] = None
    owner_participant_id: Optional[str] = None
    target_participant_id: Optional[str] = None
    blockchain: Optional[str] = None
    network: Optional[str] = None
    ids: Optional[List[str]] = None
    statuses: Optional[List[str]] = None
    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class ShareAddressRequest:
    """
    Request to share an address with a participant.

    Attributes:
        to_participant_id: The participant to share with.
        address_id: The internal address ID to share.
        key_value_attributes: Optional key-value attributes to attach.
    """

    to_participant_id: str
    address_id: str
    key_value_attributes: Optional[List[Dict[str, str]]] = None


@dataclass
class ShareWhitelistedAssetRequest:
    """
    Request to share a whitelisted asset with a participant.

    Attributes:
        to_participant_id: The participant to share with.
        whitelisted_contract_id: The whitelisted contract ID to share.
    """

    to_participant_id: str
    whitelisted_contract_id: str


def _shared_address_from_dto(dto: Any) -> Optional[SharedAddress]:
    """Convert OpenAPI shared address DTO to domain model."""
    if dto is None:
        return None

    key_value_attributes = []
    for attr in getattr(dto, "key_value_attributes", None) or []:
        key_value_attributes.append(
            {
                "key": getattr(attr, "key", "") or "",
                "value": getattr(attr, "value", "") or "",
            }
        )

    trail = []
    for entry in getattr(dto, "trail", None) or []:
        trail.append(
            SharedAddressTrail(
                status=getattr(entry, "status", "") or "",
                changed_at=getattr(entry, "changed_at", None),
            )
        )

    return SharedAddress(
        id=getattr(dto, "id", "") or "",
        owner_participant_id=getattr(dto, "owner_participant_id", "") or "",
        target_participant_id=getattr(dto, "target_participant_id", "") or "",
        address_id=getattr(dto, "address_id", "") or "",
        address=getattr(dto, "address", "") or "",
        blockchain=getattr(dto, "blockchain", "") or "",
        network=getattr(dto, "network", "") or "",
        status=getattr(dto, "status", "") or "",
        key_value_attributes=key_value_attributes,
        trail=trail,
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


def _shared_asset_from_dto(dto: Any) -> Optional[SharedAsset]:
    """Convert OpenAPI shared asset DTO to domain model."""
    if dto is None:
        return None

    return SharedAsset(
        id=getattr(dto, "id", "") or "",
        owner_participant_id=getattr(dto, "owner_participant_id", "") or "",
        target_participant_id=getattr(dto, "target_participant_id", "") or "",
        whitelisted_contract_id=getattr(dto, "whitelisted_contract_id", "") or "",
        blockchain=getattr(dto, "blockchain", "") or "",
        network=getattr(dto, "network", "") or "",
        contract_address=getattr(dto, "contract_address", "") or "",
        status=getattr(dto, "status", "") or "",
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


class SharingService(BaseService):
    """
    Service for Taurus Network shared address and asset operations.

    Provides methods to share, unshare, and list shared addresses and assets
    between Taurus Network participants.

    Example:
        >>> # List shared addresses
        >>> addresses, pagination = client.taurus_network.sharing.list_shared_addresses()
        >>> for addr in addresses:
        ...     print(f"{addr.id}: {addr.address} ({addr.status})")
        >>>
        >>> # Share an address
        >>> request = ShareAddressRequest(
        ...     to_participant_id="participant-123",
        ...     address_id="address-456",
        ... )
        >>> client.taurus_network.sharing.share_address(request)
    """

    def __init__(self, api_client: Any, shared_api: Any) -> None:
        """
        Initialize sharing service.

        Args:
            api_client: The OpenAPI client instance.
            shared_api: The TaurusNetworkSharedAddressAssetApi service.
        """
        super().__init__(api_client)
        self._shared_api = shared_api

    # =========================================================================
    # Shared Address Operations
    # =========================================================================

    def list_shared_addresses(
        self,
        options: Optional[ListSharedAddressesOptions] = None,
    ) -> Tuple[List[SharedAddress], Optional[CursorPagination]]:
        """
        List shared addresses.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (shared addresses list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListSharedAddressesOptions()

        try:
            resp = self._shared_api.taurus_network_service_get_shared_addresses(
                participant_id=opts.participant_id,
                owner_participant_id=opts.owner_participant_id,
                target_participant_id=opts.target_participant_id,
                blockchain=opts.blockchain,
                network=opts.network,
                ids=opts.ids,
                statuses=opts.statuses,
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            addresses = []
            if result:
                for dto in result:
                    addr = _shared_address_from_dto(dto)
                    if addr:
                        addresses.append(addr)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return addresses, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def share_address(self, request: ShareAddressRequest) -> None:
        """
        Share an address with a Taurus Network participant.

        This will automatically create a whitelisted address to be
        approved/rejected on the target participant side.

        Args:
            request: Address sharing parameters.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.to_participant_id, "to_participant_id")
        self._validate_required(request.address_id, "address_id")

        try:
            body = {
                "toParticipantID": request.to_participant_id,
                "addressID": request.address_id,
            }

            if request.key_value_attributes:
                body["keyValueAttributes"] = request.key_value_attributes

            self._shared_api.taurus_network_service_share_address(body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def unshare_address(self, shared_address_id: str) -> None:
        """
        Unshare an address with a Taurus Network participant.

        The address must be shared with the participant to be unshared.
        Unsharing will update the status of the shared address but not
        delete it from the registry.

        Args:
            shared_address_id: The shared address ID to unshare.

        Raises:
            ValueError: If shared_address_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(shared_address_id, "shared_address_id")

        try:
            self._shared_api.taurus_network_service_unshare_address(
                tn_shared_address_id=shared_address_id, body={}
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    # =========================================================================
    # Shared Asset Operations
    # =========================================================================

    def list_shared_assets(
        self,
        options: Optional[ListSharedAssetsOptions] = None,
    ) -> Tuple[List[SharedAsset], Optional[CursorPagination]]:
        """
        List shared whitelisted assets.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (shared assets list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListSharedAssetsOptions()

        try:
            resp = self._shared_api.taurus_network_service_get_shared_assets(
                participant_id=opts.participant_id,
                owner_participant_id=opts.owner_participant_id,
                target_participant_id=opts.target_participant_id,
                blockchain=opts.blockchain,
                network=opts.network,
                ids=opts.ids,
                statuses=opts.statuses,
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            assets = []
            if result:
                for dto in result:
                    asset = _shared_asset_from_dto(dto)
                    if asset:
                        assets.append(asset)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return assets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def share_whitelisted_asset(self, request: ShareWhitelistedAssetRequest) -> None:
        """
        Share a whitelisted asset with a Taurus Network participant.

        Args:
            request: Asset sharing parameters.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.to_participant_id, "to_participant_id")
        self._validate_required(request.whitelisted_contract_id, "whitelisted_contract_id")

        try:
            body = {
                "toParticipantID": request.to_participant_id,
                "whitelistedContractID": request.whitelisted_contract_id,
            }

            self._shared_api.taurus_network_service_share_whitelisted_asset(body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def unshare_whitelisted_asset(self, shared_asset_id: str) -> None:
        """
        Unshare a whitelisted asset with a Taurus Network participant.

        The asset must be shared with the participant to be unshared.
        Unsharing will update the status of the shared asset but not
        delete it from the registry.

        Args:
            shared_asset_id: The shared asset ID to unshare.

        Raises:
            ValueError: If shared_asset_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(shared_asset_id, "shared_asset_id")

        try:
            self._shared_api.taurus_network_service_unshare_whitelisted_asset(
                tn_shared_asset_id=shared_asset_id, body={}
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
