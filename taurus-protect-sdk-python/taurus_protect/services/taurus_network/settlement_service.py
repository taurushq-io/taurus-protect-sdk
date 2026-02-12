"""Settlement service for Taurus Network operations."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


@dataclass
class SettlementAssetTransfer:
    """
    Asset transfer within a settlement leg.

    Attributes:
        shared_address_id: The shared address ID for the transfer.
        currency: The currency/asset being transferred.
        amount: The amount to transfer.
    """

    shared_address_id: str = ""
    currency: str = ""
    amount: str = ""


@dataclass
class SettlementClipTransaction:
    """
    A transaction within a settlement clip.

    Attributes:
        id: Transaction ID.
        status: Transaction status.
    """

    id: str = ""
    status: str = ""


@dataclass
class SettlementClip:
    """
    A clip (execution batch) within a settlement.

    Attributes:
        id: Clip ID.
        status: Clip status.
        transactions: List of transactions in this clip.
    """

    id: str = ""
    status: str = ""
    transactions: List[SettlementClipTransaction] = field(default_factory=list)


@dataclass
class Settlement:
    """
    A Taurus Network settlement.

    Represents a bilateral settlement of funds between two participants.

    Attributes:
        id: The settlement ID.
        creator_participant_id: The participant who created the settlement.
        target_participant_id: The target participant for the settlement.
        first_leg_participant_id: The participant who executes the first leg.
        first_leg_assets: Assets transferred in the first leg.
        second_leg_assets: Assets transferred in the second leg.
        clips: Execution clips for the settlement.
        start_execution_date: When execution should start.
        status: Current settlement status.
        workflow_id: Associated workflow ID.
        created_at: When the settlement was created.
        updated_at: When the settlement was last updated.
    """

    id: str = ""
    creator_participant_id: str = ""
    target_participant_id: str = ""
    first_leg_participant_id: str = ""
    first_leg_assets: List[SettlementAssetTransfer] = field(default_factory=list)
    second_leg_assets: List[SettlementAssetTransfer] = field(default_factory=list)
    clips: List[SettlementClip] = field(default_factory=list)
    start_execution_date: Optional[datetime] = None
    status: str = ""
    workflow_id: str = ""
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class ListSettlementsOptions:
    """
    Options for listing settlements.

    Attributes:
        counter_participant_id: Filter by counter participant ID.
        statuses: Filter by settlement statuses.
        sort_order: Sort order (ASC or DESC).
        page_size: Number of items per page.
        current_page: Current page cursor (base64).
        page_request: Page request direction (FIRST, PREVIOUS, NEXT, LAST).
    """

    counter_participant_id: Optional[str] = None
    statuses: Optional[List[str]] = None
    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class ListSettlementsForApprovalOptions:
    """
    Options for listing settlements pending approval.

    Attributes:
        ids: Filter by settlement IDs.
        sort_order: Sort order (ASC or DESC).
        page_size: Number of items per page.
        current_page: Current page cursor (base64).
        page_request: Page request direction (FIRST, PREVIOUS, NEXT, LAST).
    """

    ids: Optional[List[str]] = None
    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class CreateSettlementRequest:
    """
    Request to create a settlement.

    Attributes:
        target_participant_id: The target participant for the settlement.
        first_leg_participant_id: The participant who executes the first leg.
        first_leg_assets: Assets to transfer in the first leg.
        second_leg_assets: Assets to transfer in the second leg.
        clips: Optional execution clips.
        start_execution_date: When execution should start.
    """

    target_participant_id: str
    first_leg_participant_id: str
    first_leg_assets: List[SettlementAssetTransfer]
    second_leg_assets: List[SettlementAssetTransfer]
    clips: Optional[List[Dict[str, Any]]] = None
    start_execution_date: Optional[datetime] = None


@dataclass
class CursorPagination:
    """
    Cursor-based pagination information.

    Attributes:
        current_page: The current page cursor.
        has_next: Whether there is a next page.
        has_previous: Whether there is a previous page.
    """

    current_page: Optional[str] = None
    has_next: bool = False
    has_previous: bool = False


def _settlement_from_dto(dto: Any) -> Optional[Settlement]:
    """Convert OpenAPI settlement DTO to domain model."""
    if dto is None:
        return None

    first_leg_assets = []
    for asset in getattr(dto, "first_leg_assets", None) or []:
        first_leg_assets.append(
            SettlementAssetTransfer(
                shared_address_id=getattr(asset, "shared_address_id", "") or "",
                currency=getattr(asset, "currency", "") or "",
                amount=getattr(asset, "amount", "") or "",
            )
        )

    second_leg_assets = []
    for asset in getattr(dto, "second_leg_assets", None) or []:
        second_leg_assets.append(
            SettlementAssetTransfer(
                shared_address_id=getattr(asset, "shared_address_id", "") or "",
                currency=getattr(asset, "currency", "") or "",
                amount=getattr(asset, "amount", "") or "",
            )
        )

    clips = []
    for clip in getattr(dto, "clips", None) or []:
        transactions = []
        for tx in getattr(clip, "transactions", None) or []:
            transactions.append(
                SettlementClipTransaction(
                    id=getattr(tx, "id", "") or "",
                    status=getattr(tx, "status", "") or "",
                )
            )
        clips.append(
            SettlementClip(
                id=getattr(clip, "id", "") or "",
                status=getattr(clip, "status", "") or "",
                transactions=transactions,
            )
        )

    return Settlement(
        id=getattr(dto, "id", "") or "",
        creator_participant_id=getattr(dto, "creator_participant_id", "") or "",
        target_participant_id=getattr(dto, "target_participant_id", "") or "",
        first_leg_participant_id=getattr(dto, "first_leg_participant_id", "") or "",
        first_leg_assets=first_leg_assets,
        second_leg_assets=second_leg_assets,
        clips=clips,
        start_execution_date=getattr(dto, "start_execution_date", None),
        status=getattr(dto, "status", "") or "",
        workflow_id=getattr(dto, "workflow_id", "") or "",
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


class SettlementService(BaseService):
    """
    Service for Taurus Network settlement operations.

    Provides methods to create, manage, and query settlements between
    Taurus Network participants.

    Example:
        >>> # List settlements
        >>> settlements, pagination = client.taurus_network.settlements.list_settlements()
        >>> for settlement in settlements:
        ...     print(f"{settlement.id}: {settlement.status}")
        >>>
        >>> # Get single settlement
        >>> settlement = client.taurus_network.settlements.get_settlement("123")
        >>> print(f"Status: {settlement.status}")
    """

    def __init__(self, api_client: Any, settlement_api: Any) -> None:
        """
        Initialize settlement service.

        Args:
            api_client: The OpenAPI client instance.
            settlement_api: The TaurusNetworkSettlementApi service.
        """
        super().__init__(api_client)
        self._settlement_api = settlement_api

    def get_settlement(self, settlement_id: str) -> Settlement:
        """
        Get a settlement by ID.

        Args:
            settlement_id: The settlement ID to retrieve.

        Returns:
            The settlement.

        Raises:
            ValueError: If settlement_id is empty.
            NotFoundError: If settlement not found.
            APIError: If API request fails.
        """
        self._validate_required(settlement_id, "settlement_id")

        try:
            resp = self._settlement_api.taurus_network_service_get_settlement(settlement_id)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Settlement {settlement_id} not found")

            settlement = _settlement_from_dto(result)
            if settlement is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Settlement {settlement_id} not found")

            return settlement
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_settlements(
        self,
        options: Optional[ListSettlementsOptions] = None,
    ) -> Tuple[List[Settlement], Optional[CursorPagination]]:
        """
        List settlements.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (settlements list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListSettlementsOptions()

        try:
            resp = self._settlement_api.taurus_network_service_get_settlements(
                counter_participant_id=opts.counter_participant_id,
                statuses=opts.statuses,
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            settlements = []
            if result:
                for dto in result:
                    settlement = _settlement_from_dto(dto)
                    if settlement:
                        settlements.append(settlement)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return settlements, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list_settlements_for_approval(
        self,
        options: Optional[ListSettlementsForApprovalOptions] = None,
    ) -> Tuple[List[Settlement], Optional[CursorPagination]]:
        """
        List settlements pending approval.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (settlements list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListSettlementsForApprovalOptions()

        try:
            resp = self._settlement_api.taurus_network_service_get_settlements_for_approval(
                ids=opts.ids,
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            settlements = []
            if result:
                for dto in result:
                    settlement = _settlement_from_dto(dto)
                    if settlement:
                        settlements.append(settlement)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return settlements, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def create_settlement(self, request: CreateSettlementRequest) -> str:
        """
        Create a new settlement.

        Args:
            request: Settlement creation parameters.

        Returns:
            The created settlement ID.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.target_participant_id, "target_participant_id")
        self._validate_required(request.first_leg_participant_id, "first_leg_participant_id")
        if not request.first_leg_assets:
            raise ValueError("first_leg_assets cannot be empty")
        if not request.second_leg_assets:
            raise ValueError("second_leg_assets cannot be empty")

        try:
            # Build the request body
            first_leg_assets = [
                {
                    "sharedAddressID": asset.shared_address_id,
                    "currency": asset.currency,
                    "amount": asset.amount,
                }
                for asset in request.first_leg_assets
            ]

            second_leg_assets = [
                {
                    "sharedAddressID": asset.shared_address_id,
                    "currency": asset.currency,
                    "amount": asset.amount,
                }
                for asset in request.second_leg_assets
            ]

            body = {
                "targetParticipantID": request.target_participant_id,
                "firstLegParticipantID": request.first_leg_participant_id,
                "firstLegAssets": first_leg_assets,
                "secondLegAssets": second_leg_assets,
            }

            if request.clips:
                body["clips"] = request.clips
            if request.start_execution_date:
                body["startExecutionDate"] = request.start_execution_date.isoformat()

            resp = self._settlement_api.taurus_network_service_create_settlement(body=body)

            result = getattr(resp, "result", None)
            if result:
                return getattr(result, "id", "") or ""

            # Check for id directly on response
            return getattr(resp, "id", "") or ""
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def cancel_settlement(self, settlement_id: str) -> None:
        """
        Cancel a settlement.

        Args:
            settlement_id: The settlement ID to cancel.

        Raises:
            ValueError: If settlement_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(settlement_id, "settlement_id")

        try:
            self._settlement_api.taurus_network_service_cancel_settlement(settlement_id, body={})
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def replace_settlement(self, settlement_id: str, request: CreateSettlementRequest) -> None:
        """
        Replace a settlement with new attributes.

        This endpoint replaces a settlement with new attributes at target side
        before the target approves the settlement.

        Args:
            settlement_id: The settlement ID to replace.
            request: New settlement parameters.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        self._validate_required(settlement_id, "settlement_id")
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.target_participant_id, "target_participant_id")
        self._validate_required(request.first_leg_participant_id, "first_leg_participant_id")
        if not request.first_leg_assets:
            raise ValueError("first_leg_assets cannot be empty")
        if not request.second_leg_assets:
            raise ValueError("second_leg_assets cannot be empty")

        try:
            # Build the request body
            first_leg_assets = [
                {
                    "sharedAddressID": asset.shared_address_id,
                    "currency": asset.currency,
                    "amount": asset.amount,
                }
                for asset in request.first_leg_assets
            ]

            second_leg_assets = [
                {
                    "sharedAddressID": asset.shared_address_id,
                    "currency": asset.currency,
                    "amount": asset.amount,
                }
                for asset in request.second_leg_assets
            ]

            create_settlement_request = {
                "targetParticipantID": request.target_participant_id,
                "firstLegParticipantID": request.first_leg_participant_id,
                "firstLegAssets": first_leg_assets,
                "secondLegAssets": second_leg_assets,
            }

            if request.clips:
                create_settlement_request["clips"] = request.clips
            if request.start_execution_date:
                create_settlement_request["startExecutionDate"] = (
                    request.start_execution_date.isoformat()
                )

            body = {"createSettlementRequest": create_settlement_request}

            self._settlement_api.taurus_network_service_replace_settlement(settlement_id, body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
