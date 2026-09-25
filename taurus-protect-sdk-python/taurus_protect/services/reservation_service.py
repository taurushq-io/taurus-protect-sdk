"""Reservation service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.models.pagination import CursorPage, cursor_page, cursor_request
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.reservations_api import ReservationsApi


class Reservation:
    """A balance reservation for pending transactions."""

    def __init__(
        self,
        id: str,
        wallet_id: Optional[str] = None,
        address_id: Optional[str] = None,
        currency: Optional[str] = None,
        amount: Optional[str] = None,
        status: Optional[str] = None,
        expires_at: Optional[datetime] = None,
        kind: Optional[str] = None,
        comment: Optional[str] = None,
        address: Optional[str] = None,
        created_at: Optional[datetime] = None,
        resource_id: Optional[str] = None,
        resource_type: Optional[str] = None,
    ):
        self.id = id
        self.wallet_id = wallet_id
        self.address_id = address_id
        self.currency = currency
        self.amount = amount
        self.status = status
        self.expires_at = expires_at
        self.kind = kind
        self.comment = comment
        self.address = address
        self.created_at = created_at
        self.resource_id = resource_id
        self.resource_type = resource_type


class ReservationService(BaseService):
    """
    Service for managing balance reservations.

    Reservations lock funds for pending transactions to prevent
    double-spending.
    """

    def __init__(
        self,
        api_client: Any,
        reservations_api: "ReservationsApi",
    ) -> None:
        super().__init__(api_client)
        self._api = reservations_api

    def get(self, reservation_id: int) -> Reservation:
        """Get a reservation by ID."""
        if reservation_id <= 0:
            raise ValueError("reservation_id must be positive")

        try:
            reply = self._api.wallet_service_get_reservation(str(reservation_id))
            result = reply.result
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Reservation {reservation_id} not found")
            return self._map_reservation_from_dto(result)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def list(
        self,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        kind: Optional[str] = None,
        kinds: Optional[List[str]] = None,
        address: Optional[str] = None,
        address_id: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Reservation], CursorPage]:
        """
        List reservations, one page at a time.

        Args:
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            kind: Filter by reservation kind.
            kinds: Filter by any of these kinds.
            address: Filter by blockchain address.
            address_id: Filter by address ID.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (reservations, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If the API call fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            reply = self._api.wallet_service_get_reservations(
                kind=kind,
                kinds=kinds,
                address=address,
                address_id=address_id,
                **req.query_params(),
            )

            reservations = [self._map_reservation_from_dto(dto) for dto in reply.result or []]
            return reservations, cursor_page(req.page_size, reply.cursor)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    @staticmethod
    def _map_reservation_from_dto(dto: Any) -> Reservation:
        """Map a generated ``TgvalidatordReservation``."""
        currency_info = dto.currency_info
        return Reservation(
            id=dto.id or "",
            address_id=dto.addressid,
            currency=currency_info.symbol if currency_info is not None else None,
            amount=dto.amount,
            kind=dto.kind,
            comment=dto.comment,
            address=dto.address,
            created_at=dto.creation_date,
            resource_id=dto.resource_id,
            resource_type=dto.resource_type,
        )
