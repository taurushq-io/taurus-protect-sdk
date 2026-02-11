"""Unit tests for reservation mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.services.reservation_service import (
    Reservation,
    ReservationService,
)


class TestMapReservationFromDto:
    """Tests for ReservationService._map_reservation_from_dto."""

    def test_maps_all_fields(self) -> None:
        expires = datetime(2024, 12, 31, 23, 59, 59, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id=42,
            wallet_id="w-1",
            address_id="a-1",
            currency="BTC",
            amount="0.5",
            status="ACTIVE",
            expires_at=expires,
            expiresAt=None,
        )
        result = ReservationService._map_reservation_from_dto(dto)
        assert result.id == "42"
        assert result.wallet_id == "w-1"
        assert result.address_id == "a-1"
        assert result.currency == "BTC"
        assert result.amount == "0.5"
        assert result.status == "ACTIVE"
        assert result.expires_at == expires

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id=1,
            wallet_id=None,
            address_id=None,
            currency=None,
            amount=None,
            status=None,
            expires_at=None,
            expiresAt=None,
        )
        result = ReservationService._map_reservation_from_dto(dto)
        assert result.id == "1"
        assert result.wallet_id is None
        assert result.address_id is None
        assert result.currency is None

    def test_expires_at_camelcase_fallback(self) -> None:
        expires = datetime(2025, 1, 1, 0, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id=2,
            wallet_id=None,
            address_id=None,
            currency=None,
            amount=None,
            status=None,
            expires_at=None,
            expiresAt=expires,
        )
        result = ReservationService._map_reservation_from_dto(dto)
        assert result.expires_at == expires

    def test_result_is_reservation_instance(self) -> None:
        dto = SimpleNamespace(
            id=5,
            wallet_id="w-2",
            address_id=None,
            currency="ETH",
            amount="1.0",
            status="PENDING",
            expires_at=None,
            expiresAt=None,
        )
        result = ReservationService._map_reservation_from_dto(dto)
        assert isinstance(result, Reservation)
