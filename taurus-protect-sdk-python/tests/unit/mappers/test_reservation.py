"""Unit tests for reservation mapper functions."""

from taurus_protect._internal.openapi.models import TgvalidatordReservation
from taurus_protect.services.reservation_service import (
    Reservation,
    ReservationService,
)


class TestMapReservationFromDto:
    """ReservationService._map_reservation_from_dto reads the generated reservation's fields."""

    def test_maps_all_fields(self) -> None:
        dto = TgvalidatordReservation.from_dict(
            {
                "id": "42",
                "amount": "0.5",
                "creationDate": "2024-12-31T23:59:59Z",
                "kind": "UTXO",
                "comment": "hold",
                "addressid": "a-1",
                "address": "bc1q",
                "currencyInfo": {"symbol": "BTC"},
                "resourceId": "r-1",
                "resourceType": "request",
            }
        )
        result = ReservationService._map_reservation_from_dto(dto)
        assert result.id == "42"
        assert result.address_id == "a-1"
        assert result.address == "bc1q"
        assert result.currency == "BTC"
        assert result.amount == "0.5"
        assert result.kind == "UTXO"
        assert result.comment == "hold"
        assert result.created_at is not None
        assert (result.resource_id, result.resource_type) == ("r-1", "request")

    def test_handles_none_optional_fields(self) -> None:
        result = ReservationService._map_reservation_from_dto(
            TgvalidatordReservation.from_dict({"id": "1"})
        )
        assert result.id == "1"
        assert result.address_id is None
        assert result.currency is None
        assert result.kind is None

    def test_missing_id_maps_to_empty_string(self) -> None:
        result = ReservationService._map_reservation_from_dto(TgvalidatordReservation.from_dict({}))
        assert result.id == ""

    def test_result_is_reservation_instance(self) -> None:
        result = ReservationService._map_reservation_from_dto(
            TgvalidatordReservation.from_dict({"id": "5"})
        )
        assert isinstance(result, Reservation)
