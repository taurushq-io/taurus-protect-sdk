"""Unit tests for ReservationService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.reservation_service import ReservationService


class TestReservationServiceGet:
    """Tests for ReservationService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        reservations_api = MagicMock()
        service = ReservationService(
            api_client=api_client, reservations_api=reservations_api
        )
        return service, reservations_api

    def test_raises_on_invalid_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="reservation_id must be positive"):
            service.get(reservation_id=0)

    def test_returns_reservation(self) -> None:
        service, api = self._make_service()
        dto = MagicMock()
        dto.id = "10"
        dto.wallet_id = "5"
        dto.address_id = None
        dto.currency = "ETH"
        dto.amount = "1.5"
        dto.status = "ACTIVE"
        dto.expires_at = None
        dto.expiresAt = None
        reply = MagicMock()
        reply.result = dto
        api.reservation_service_get_reservation.return_value = reply

        res = service.get(reservation_id=10)

        assert res.id == "10"
        assert res.currency == "ETH"
        assert res.amount == "1.5"

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.reservation_service_get_reservation.return_value = reply

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(reservation_id=1)


class TestReservationServiceList:
    """Tests for ReservationService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        reservations_api = MagicMock()
        service = ReservationService(
            api_client=api_client, reservations_api=reservations_api
        )
        return service, reservations_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.reservation_service_get_reservations.return_value = reply

        reservations, pagination = service.list()

        assert reservations == []

    def test_passes_wallet_id_filter(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.reservation_service_get_reservations.return_value = reply

        service.list(wallet_id=5)

        api.reservation_service_get_reservations.assert_called_once_with(
            wallet_id="5",
            limit="50",
            offset="0",
        )


class TestReservationServiceCancel:
    """Tests for ReservationService.cancel()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        reservations_api = MagicMock()
        service = ReservationService(
            api_client=api_client, reservations_api=reservations_api
        )
        return service, reservations_api

    def test_raises_on_invalid_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="reservation_id must be positive"):
            service.cancel(reservation_id=0)

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        api.reservation_service_cancel_reservation.return_value = None

        service.cancel(reservation_id=42)

        api.reservation_service_cancel_reservation.assert_called_once_with("42")
