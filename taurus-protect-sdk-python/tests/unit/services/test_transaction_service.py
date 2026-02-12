"""Unit tests for TransactionService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.transaction_service import TransactionService


class TestGet:
    """Tests for TransactionService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(
            api_client=api_client, transactions_api=transactions_api
        )
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
        service = TransactionService(
            api_client=api_client, transactions_api=transactions_api
        )
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
    """Tests for TransactionService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(
            api_client=api_client, transactions_api=transactions_api
        )
        return service, transactions_api

    def test_list_returns_transactions_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "10"
        api.transaction_service_get_transactions.return_value = reply

        with patch(
            "taurus_protect.services.transaction_service.map_transactions",
            return_value=[MagicMock()],
        ):
            transactions, pagination = service.list(limit=50)

        assert len(transactions) == 1

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)


class TestListByAddress:
    """Tests for TransactionService.list_by_address()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(
            api_client=api_client, transactions_api=transactions_api
        )
        return service, transactions_api

    def test_list_by_address_raises_for_empty_address(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address"):
            service.list_by_address("")

    def test_list_by_address_returns_transactions(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "1"
        api.transaction_service_get_transactions.return_value = reply

        with patch(
            "taurus_protect.services.transaction_service.map_transactions",
            return_value=[MagicMock()],
        ):
            transactions, _ = service.list_by_address("0xabc123")

        assert len(transactions) == 1


class TestExportCsv:
    """Tests for TransactionService.export_csv()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        transactions_api = MagicMock()
        service = TransactionService(
            api_client=api_client, transactions_api=transactions_api
        )
        return service, transactions_api

    def test_export_csv_returns_string(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = "col1,col2\nval1,val2"
        api.transaction_service_export_transactions.return_value = reply

        result = service.export_csv()
        assert result == "col1,col2\nval1,val2"

    def test_export_csv_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.transaction_service_export_transactions.return_value = reply

        result = service.export_csv()
        assert result == ""

    def test_export_csv_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.export_csv(limit=0)

    def test_export_csv_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.export_csv(offset=-1)
