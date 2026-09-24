"""Unit tests for BalanceService."""

from __future__ import annotations

from taurus_protect._internal.openapi import BalancesApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.balance_service import BalanceService
from tests.unit.transport_stub import StubTransport, api_client


class TestList:
    """BalanceService.list: request cursor only, with the server total."""

    def _service(self) -> BalanceService:
        ac = api_client()
        return BalanceService(ac, BalancesApi(ac))

    def test_only_the_request_cursor_is_sent(self) -> None:
        """The list sent both the legacy limit and the request cursor, and read no total."""
        reply = {
            "balances": [{"asset": {"currency": "ETH"}}],
            "total": "7",
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            balances, page = self._service().list(currency="ETH", page_size=1, token_id="t")

        assert len(balances) == 1
        assert page == CursorPage(page_size=1, next_cursor="n", has_more=True, total_items=7)
        assert transport.last.query == sorted(
            [("currency", "ETH"), ("tokenId", "t"), ("requestCursor.pageSize", "1")]
        )

    def test_page_two_is_reachable(self) -> None:
        with StubTransport() as transport:
            self._service().list(cursor="n")

        assert transport.last.query == sorted(
            [
                ("requestCursor.currentPage", "n"),
                ("requestCursor.pageRequest", "NEXT"),
                ("requestCursor.pageSize", "20"),
            ]
        )

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            balances, page = self._service().list()

        assert balances == []
        assert page == CursorPage(page_size=20, total_items=0)


class TestListNftCollections:
    """list_nft_collections reads the reply's ``balances`` (it read fields that do not exist)."""

    def _service(self) -> BalanceService:
        ac = api_client()
        return BalanceService(ac, BalancesApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "balances": [{"name": "c1"}, {"name": "c2"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            collections, page = self._service().list_nft_collections(
                "ETH", "mainnet", page_size=2, query="apes", only_positive_balance=True
            )

        assert len(collections) == 2
        assert page == CursorPage(page_size=2, next_cursor="n", has_more=True)
        assert transport.last.query == sorted(
            [
                ("blockchain", "ETH"),
                ("network", "mainnet"),
                ("query", "apes"),
                ("onlyPositiveBalance", "true"),
                ("cursor.pageSize", "2"),
            ]
        )

    def test_blockchain_and_network_are_optional(self) -> None:
        with StubTransport() as transport:
            self._service().list_nft_collections(cursor="n")

        assert transport.last.query == sorted(
            [("cursor.currentPage", "n"), ("cursor.pageRequest", "NEXT"), ("cursor.pageSize", "20")]
        )
