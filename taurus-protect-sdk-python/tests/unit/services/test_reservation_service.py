"""Unit tests for ReservationService."""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi import ReservationsApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.reservation_service import Reservation, ReservationService
from tests.unit.transport_stub import StubTransport, api_client


def _service() -> ReservationService:
    ac = api_client()
    return ReservationService(ac, ReservationsApi(ac))


class TestReservationServiceGet:
    """ReservationService.get calls the real generated operation."""

    def test_raises_on_invalid_id(self) -> None:
        with pytest.raises(ValueError, match="reservation_id must be positive"):
            _service().get(0)

    def test_returns_reservation(self) -> None:
        reply = {"result": {"id": "42", "amount": "0.5", "addressid": "a-1", "kind": "UTXO"}}
        with StubTransport(reply) as transport:
            result = _service().get(42)

        assert isinstance(result, Reservation)
        assert (result.id, result.amount, result.address_id, result.kind) == (
            "42",
            "0.5",
            "a-1",
            "UTXO",
        )
        assert transport.last.path == "/api/rest/v1/reservations/42"

    def test_raises_not_found_when_absent(self) -> None:
        with StubTransport({}):
            with pytest.raises(NotFoundError):
                _service().get(42)


class TestReservationServiceList:
    """ReservationService.list: a cursor list (it called an operation that does not exist)."""

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "result": [{"id": "1"}, {"id": "2"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            reservations, page = _service().list(
                page_size=2, kind="UTXO", kinds=["A", "B"], address="bc1q", address_id="7"
            )

        assert [r.id for r in reservations] == ["1", "2"]
        assert page == CursorPage(page_size=2, next_cursor="n", has_more=True)
        assert transport.last.path == "/api/rest/v1/reservations"
        assert transport.last.query == sorted(
            [
                ("kind", "UTXO"),
                ("kinds", "A"),
                ("kinds", "B"),
                ("address", "bc1q"),
                ("addressId", "7"),
                ("cursor.pageSize", "2"),
            ]
        )

    def test_page_two_is_reachable(self) -> None:
        with StubTransport() as transport:
            _service().list(cursor="n")

        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"

    def test_the_wallet_filter_does_not_exist(self) -> None:
        with pytest.raises(TypeError, match="wallet_id"):
            _service().list(wallet_id=1)  # type: ignore[call-arg]

    def test_invalid_page_size_sends_nothing(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                _service().list(page_size=101)

        assert transport.requests == []

    def test_there_is_no_cancel(self) -> None:
        """No endpoint cancels a reservation; the method called an operation that does not exist."""
        assert not hasattr(ReservationService, "cancel")
