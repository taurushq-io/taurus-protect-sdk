"""Reservation service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.models.pagination import Pagination
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
    ):
        self.id = id
        self.wallet_id = wallet_id
        self.address_id = address_id
        self.currency = currency
        self.amount = amount
        self.status = status
        self.expires_at = expires_at


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
            reply = self._api.reservation_service_get_reservation(str(reservation_id))
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
        wallet_id: Optional[int] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Reservation], Optional[Pagination]]:
        """List reservations."""
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            reply = self._api.reservation_service_get_reservations(
                wallet_id=str(wallet_id) if wallet_id else None,
                limit=str(limit),
                offset=str(offset),
            )

            reservations: List[Reservation] = []
            if reply.result:
                for dto in reply.result:
                    reservations.append(self._map_reservation_from_dto(dto))

            pagination = self._extract_pagination(
                getattr(reply, "total_items", None),
                offset,
                limit,
            )
            return reservations, pagination
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def cancel(self, reservation_id: int) -> None:
        """Cancel a reservation."""
        if reservation_id <= 0:
            raise ValueError("reservation_id must be positive")

        try:
            self._api.reservation_service_cancel_reservation(str(reservation_id))
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    @staticmethod
    def _map_reservation_from_dto(dto: Any) -> Reservation:
        return Reservation(
            id=str(getattr(dto, "id", "")),
            wallet_id=(
                str(getattr(dto, "wallet_id", "")) if getattr(dto, "wallet_id", None) else None
            ),
            address_id=(
                str(getattr(dto, "address_id", "")) if getattr(dto, "address_id", None) else None
            ),
            currency=getattr(dto, "currency", None),
            amount=getattr(dto, "amount", None),
            status=getattr(dto, "status", None),
            expires_at=getattr(dto, "expires_at", None) or getattr(dto, "expiresAt", None),
        )
