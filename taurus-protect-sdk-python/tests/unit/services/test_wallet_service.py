"""Unit tests for WalletService."""

from __future__ import annotations

import base64
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import WalletsApi
from taurus_protect.errors import APIError, NotFoundError
from taurus_protect.models.pagination import CursorPage, Pagination
from taurus_protect.models.wallet import ListWalletsOptions
from taurus_protect.services.wallet_service import WalletService
from tests.unit.transport_stub import StubTransport, api_client


class TestGet:
    """Tests for WalletService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_returns_wallet(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_get_wallet_v2.return_value = reply

        mock_wallet = MagicMock()
        with patch(
            "taurus_protect.services.wallet_service.wallet_from_dto",
            return_value=mock_wallet,
        ):
            result = service.get(1)

        assert result is mock_wallet
        api.wallet_service_get_wallet_v2.assert_called_once_with("1")

    def test_get_raises_value_error_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.get(0)

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.wallet_service_get_wallet_v2.return_value = reply

        with pytest.raises(NotFoundError):
            service.get(1)


class TestList:
    """WalletService.list / list_with_options over the real generated client."""

    def _service(self) -> WalletService:
        ac = api_client()
        return WalletService(ac, WalletsApi(ac))

    def test_reply_offset_is_the_next_offset(self) -> None:
        """The reply offset is offset + rows: reading it as the current offset lost the last page."""
        page1 = {"result": [{"id": str(i)} for i in range(20)], "totalItems": "25", "offset": "20"}
        page2 = {
            "result": [{"id": str(i)} for i in range(20, 25)],
            "totalItems": "25",
            "offset": "25",
        }
        service = self._service()

        with StubTransport(page1, page2) as transport:
            wallets, pagination = service.list()
            assert len(wallets) == 20
            assert pagination == Pagination(
                limit=20, offset=0, total_items=25, next_offset=20, has_more=True
            )
            wallets, pagination = service.list(offset=pagination.next_offset)

        assert len(wallets) == 5
        assert pagination.has_more is False
        assert pagination.next_offset == 25
        assert transport.requests[1].param("offset") == "20"

    def test_empty_reply_is_a_complete_empty_page(self) -> None:
        with StubTransport({}):
            wallets, pagination = self._service().list()

        assert wallets == []
        assert pagination == Pagination(limit=20, offset=0, next_offset=0)

    def test_every_filter_reaches_the_wire(self) -> None:
        """list_with_options used to drop name, sort_order, tag_ids, balance, chain and ids."""
        options = ListWalletsOptions(
            limit=5,
            offset=10,
            currency="ETH",
            query="q",
            name="treasury",
            sort_order="ASC",
            exclude_disabled=True,
            tag_ids=["t1", "t2"],
            only_positive_balance=True,
            blockchain="ETH",
            network="mainnet",
            ids=["1", "2"],
        )
        with StubTransport() as transport:
            self._service().list_with_options(options)

        assert transport.last.path == "/api/rest/v2/wallets"
        assert transport.last.query == sorted(
            [
                ("limit", "5"),
                ("offset", "10"),
                ("currencies", "ETH"),
                ("query", "q"),
                ("name", "treasury"),
                ("sortOrder", "ASC"),
                ("excludeDisabled", "true"),
                ("tagIDs", "t1"),
                ("tagIDs", "t2"),
                ("onlyPositiveBalance", "true"),
                ("blockchain", "ETH"),
                ("network", "mainnet"),
                ("ids", "1"),
                ("ids", "2"),
            ]
        )

    def test_exclude_disabled_false_reaches_the_wire(self) -> None:
        with StubTransport() as transport:
            self._service().list(exclude_disabled=False)

        assert transport.last.param("excludeDisabled") == "false"

    @pytest.mark.parametrize(
        "kwargs,name",
        [({"limit": 101}, "limit"), ({"limit": -1}, "limit"), ({"offset": -1}, "offset")],
    )
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []


class TestCreate:
    """Tests for WalletService.create()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_create_raises_for_none_request(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="request cannot be None"):
            service.create(None)

    def test_create_wallet_raises_for_empty_blockchain(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain"):
            service.create_wallet(blockchain="", network="mainnet", name="test")

    def test_create_wallet_raises_for_empty_name(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="name"):
            service.create_wallet(blockchain="ETH", network="mainnet", name="")

    def test_create_wallet_returns_wallet(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_create_wallet.return_value = reply

        mock_wallet = MagicMock()
        with patch(
            "taurus_protect.services.wallet_service.wallet_from_create_dto",
            return_value=mock_wallet,
        ):
            result = service.create_wallet("ETH", "mainnet", "Test Wallet")

        assert result is mock_wallet


class TestGetByName:
    """Tests for WalletService.get_by_name()."""

    def _service(self) -> WalletService:
        ac = api_client()
        return WalletService(ac, WalletsApi(ac))

    def test_get_by_name_raises_for_empty_name(self) -> None:
        with pytest.raises(ValueError, match="name"):
            self._service().get_by_name("")

    def test_name_and_exclude_disabled_reach_the_wire(self) -> None:
        """exclude_disabled was hard-coded to None on this path."""
        with StubTransport(
            {"result": [{"id": "1", "name": "Test"}], "totalItems": "1", "offset": "1"}
        ) as transport:
            wallets, pagination = self._service().get_by_name("Test", exclude_disabled=True)

        assert [w.name for w in wallets] == ["Test"]
        assert pagination.has_more is False
        assert transport.last.query == sorted(
            [("name", "Test"), ("excludeDisabled", "true"), ("limit", "20")]
        )


class TestCreateAttribute:
    """Tests for WalletService.create_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_create_attribute_raises_for_non_positive_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.create_attribute(0, "key", "value")

    def test_create_attribute_raises_for_empty_key(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="key"):
            service.create_attribute(1, "", "value")

    def test_create_attribute_calls_api(self) -> None:
        service, api = self._make_service()

        service.create_attribute(1, "key", "value")

        api.wallet_service_create_wallet_attributes.assert_called_once()


class TestGetBalanceHistory:
    """Tests for WalletService.get_balance_history()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_balance_history_validates_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.get_balance_history(0, 24)

    def test_get_balance_history_validates_interval(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="interval_hours must be positive"):
            service.get_balance_history(1, 0)

    def test_get_balance_history_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.wallet_service_get_wallet_balance_history.return_value = reply

        result = service.get_balance_history(1, 24)
        assert result == []


class TestGetTokens:
    """WalletService.get_tokens: a token-only list (request cursor + limit, reply next + total)."""

    def _service(self) -> WalletService:
        ac = api_client()
        return WalletService(ac, WalletsApi(ac))

    def test_get_tokens_validates_wallet_id(self) -> None:
        with pytest.raises(ValueError, match="wallet_id must be positive"):
            self._service().get_tokens(0)

    def test_get_tokens_validates_page_size(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                self._service().get_tokens(1, page_size=101)
        assert transport.requests == []

    def test_empty_reply(self) -> None:
        with StubTransport({}) as transport:
            balances, page = self._service().get_tokens(1)

        assert balances == []
        assert page == CursorPage(page_size=20, total_items=0)
        assert transport.last.query == [("limit", "20")]

    def test_the_next_token_reaches_page_two_verbatim(self) -> None:
        """Page 2 was unreachable: the reply token was never exposed or sent back."""
        token = base64.b64encode(b"\xfb\xef\xbe\xff").decode()
        row = {"asset": {"currency": "ETH"}}
        with StubTransport(
            {"balances": [row, row], "next": token, "total": "3"},
            {"balances": [row], "total": "3"},
        ) as transport:
            balances, page = self._service().get_tokens(7, page_size=2)
            assert len(balances) == 2
            assert page == CursorPage(page_size=2, next_cursor=token, has_more=True, total_items=3)
            balances, page = self._service().get_tokens(7, page_size=2, cursor=page.next_cursor)

        assert len(balances) == 1
        assert page.has_more is False and page.next_cursor == ""
        assert transport.requests[1].param("cursor") == token
        assert "%2B" in transport.requests[1].raw_query
