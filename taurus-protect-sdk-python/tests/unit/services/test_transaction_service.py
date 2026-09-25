"""Unit tests for TransactionService."""

from __future__ import annotations

from datetime import datetime, timezone
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import TransactionsApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.transaction import TransactionExport
from taurus_protect.services.transaction_service import TransactionService
from tests.unit.transport_stub import StubTransport, api_client


class TestGet:
    """Tests for TransactionService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(api_client=api_client, transactions_api=transactions_api)
        return service, transactions_api

    def test_get_returns_transaction(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        api.transaction_service_get_transactions.return_value = reply

        mock_tx = MagicMock()
        with patch(
            "taurus_protect.services.transaction_service.map_transaction",
            return_value=mock_tx,
        ):
            result = service.get(1)

        assert result is mock_tx

    def test_get_raises_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="transaction_id must be positive"):
            service.get(0)

    def test_get_raises_not_found_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = []
        api.transaction_service_get_transactions.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get(999)


class TestGetByHash:
    """Tests for TransactionService.get_by_hash()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(api_client=api_client, transactions_api=transactions_api)
        return service, transactions_api

    def test_get_by_hash_raises_for_empty_hash(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="tx_hash"):
            service.get_by_hash("")

    def test_get_by_hash_returns_transaction(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        api.transaction_service_get_transactions.return_value = reply

        mock_tx = MagicMock()
        with patch(
            "taurus_protect.services.transaction_service.map_transaction",
            return_value=mock_tx,
        ):
            result = service.get_by_hash("0x1234")

        assert result is mock_tx

    def test_get_by_hash_raises_not_found(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = []
        api.transaction_service_get_transactions.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get_by_hash("0xnonexistent")


class TestList:
    """TransactionService.list over the real generated client (next offset = offset + rows)."""

    def _service(self) -> TransactionService:
        ac = api_client()
        return TransactionService(ac, TransactionsApi(ac))

    def test_walk_ends_on_a_trailing_empty_page(self) -> None:
        """The total can be an upper bound: an empty page with no progress ends the walk."""
        rows = [{"id": str(i)} for i in range(2)]
        with StubTransport(
            {"result": rows, "totalItems": "5"},
            {"totalItems": "5"},
        ) as transport:
            transactions, pagination = self._service().list(limit=2)
            assert len(transactions) == 2
            assert pagination == Pagination(
                limit=2, offset=0, total_items=5, next_offset=2, has_more=True
            )
            transactions, pagination = self._service().list(limit=2, offset=2)

        assert transactions == []
        assert pagination == Pagination(
            limit=2, offset=2, total_items=5, next_offset=2, has_more=False
        )
        assert transport.requests[1].query == [("limit", "2"), ("offset", "2")]

    def test_filters_reach_the_wire(self) -> None:
        start = datetime(2026, 1, 2, 3, 4, 5, tzinfo=timezone.utc)
        with StubTransport() as transport:
            self._service().list(
                from_date=start,
                to_date=start,
                currency="ETH",
                direction="incoming",
                blockchain="ETH",
                network="mainnet",
            )

        query = dict(transport.last.query)
        assert query["currency"] == "ETH"
        assert query["direction"] == "incoming"
        assert query["blockchain"] == "ETH"
        assert query["network"] == "mainnet"
        assert query["from"].startswith("2026-01-02T03:04:05")
        assert query["to"].startswith("2026-01-02T03:04:05")
        assert query["limit"] == "20"

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []


class TestListByAddress:
    """Tests for TransactionService.list_by_address()."""

    def _service(self) -> TransactionService:
        ac = api_client()
        return TransactionService(ac, TransactionsApi(ac))

    def test_list_by_address_raises_for_empty_address(self) -> None:
        with pytest.raises(ValueError, match="address"):
            self._service().list_by_address("")

    def test_address_and_page_window_reach_the_wire(self) -> None:
        with StubTransport({"result": [{"id": "1"}], "totalItems": "1"}) as transport:
            transactions, pagination = self._service().list_by_address("0xabc123", offset=40)

        assert [t.id for t in transactions] == ["1"]
        assert pagination.next_offset == 41 and pagination.has_more is False
        assert transport.last.query == [("address", "0xabc123"), ("limit", "20"), ("offset", "40")]


class TestExport:
    """The export cannot page: the server always exports from the first matching row."""

    def _service(self) -> TransactionService:
        ac = api_client()
        return TransactionService(ac, TransactionsApi(ac))

    def test_returns_the_text_and_the_server_total(self) -> None:
        with StubTransport({"result": "id,amount\n1,5\n", "totalItems": "40"}) as transport:
            export = self._service().export(limit=5000, blockchain="ETH", network="mainnet")

        assert export == TransactionExport(content="id,amount\n1,5\n", total_items=40)
        assert transport.last.path == "/api/rest/v1/transactions/export"
        assert transport.last.query == sorted(
            [("limit", "5000"), ("blockchain", "ETH"), ("network", "mainnet")]
        )

    def test_empty_reply(self) -> None:
        with StubTransport({}) as transport:
            export = self._service().export()

        assert export == TransactionExport(content="", total_items=0)
        assert transport.last.query == [("limit", "20")], "format only when the caller passes it"

    def test_export_csv_asks_for_csv(self) -> None:
        with StubTransport() as transport:
            self._service().export_csv(currency="ETH")

        assert transport.last.query == sorted(
            [("currency", "ETH"), ("format", "csv"), ("limit", "20")]
        )

    def test_there_is_no_offset(self) -> None:
        with pytest.raises(TypeError, match="offset"):
            self._service().export(offset=20)  # type: ignore[call-arg]

    def test_negative_limit_sends_nothing(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="limit"):
                self._service().export_csv(limit=-1)

        assert transport.requests == []
