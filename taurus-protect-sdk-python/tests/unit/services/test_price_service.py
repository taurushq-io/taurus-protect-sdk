"""Unit tests for PriceService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi import PricesApi
from taurus_protect.errors import IntegrityError
from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.price_service import PriceService
from tests.unit.transport_stub import Reply, StubTransport, api_client


class TestGetCurrent:
    """get_current lists prices through QueryPricesV2 (body cursor), each row verified."""

    def _service(self, rules_cache=None) -> PriceService:
        ac = api_client()
        return PriceService(ac, PricesApi(ac), rules_cache=rules_cache or MagicMock())

    def test_uses_query_prices_v2(self) -> None:
        with StubTransport({}) as transport:
            prices, page = self._service().get_current()

        assert prices == []
        assert page == CursorPage(page_size=20)
        assert transport.last.method == "POST"
        assert transport.last.path == "/api/rest/v2/prices/query"
        assert transport.last.body == {"cursor": {"pageSize": "20"}}

    @pytest.mark.parametrize(
        "kwargs,body",
        [
            ({"from_currency_id": "c1"}, {"from": {"currencyFromId": "c1"}}),
            ({"to_currency_ids": ["c2"]}, {"to": {"currencyToIds": ["c2"]}}),
            (
                {"from_currency_id": "c1", "to_currency_ids": ["c2", "c3"]},
                {"fromTo": {"currencyFromId": "c1", "currencyToIds": ["c2", "c3"]}},
            ),
        ],
    )
    def test_currency_filter_is_the_matching_oneof(self, kwargs: dict, body: dict) -> None:
        with StubTransport() as transport:
            self._service().get_current(only_primary=False, sort_order="ASC", **kwargs)

        assert transport.last.body == {
            "onlyPrimary": False,
            "sortOrder": "ASC",
            "cursor": {"pageSize": "20"},
            **body,
        }

    def test_walk_sends_the_cursor_back(self) -> None:
        row = {"currencyFrom": "BTC", "currencyTo": "USD", "rate": "1"}
        with StubTransport(
            {"result": [row], "cursor": {"currentPage": "n+/=", "hasNext": True}},
            {"result": [row], "cursor": {"currentPage": "z"}},
        ) as transport:
            prices, page = self._service().get_current(page_size=1)
            assert page == CursorPage(page_size=1, next_cursor="n+/=", has_more=True)
            prices, page = self._service().get_current(page_size=1, cursor=page.next_cursor)

        assert page.has_more is False
        assert transport.requests[1].body["cursor"] == {
            "pageSize": "1",
            "currentPage": "n+/=",
            "pageRequest": "NEXT",
        }

    def test_every_row_is_verified(self) -> None:
        """Same signed CurrencyPrice rows as v1: an unsigned price fails under a PRICEUPDATER."""
        pub = ec.generate_private_key(ec.SECP256R1()).public_key()
        pem = pub.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        ).decode()
        container = DecodedRulesContainer(
            users=[RuleUser(id="p@bank", public_key_pem=pem, roles=["PRICEUPDATER"])], groups=[]
        )
        rules_cache = MagicMock()
        rules_cache.get_decoded_rules_container.return_value = container
        reply = {"result": [{"currencyFrom": "BTC", "currencyTo": "USD", "rate": "1"}]}

        with StubTransport(reply):
            with pytest.raises(IntegrityError, match="signature"):
                self._service(rules_cache).get_current()

    def test_the_old_currency_argument_is_refused(self) -> None:
        with pytest.raises(TypeError):
            self._service().get_current("BTC")  # type: ignore[misc]

    def test_wraps_api_error(self) -> None:
        from taurus_protect.errors import APIError

        with StubTransport(Reply({"message": "boom"}, status=500)):
            with pytest.raises(APIError):
                self._service().get_current()


class TestGetHistorical:
    """get_historical: cannot page, so the limit is the whole series (default 20, max 365)."""

    def _service(self) -> PriceService:
        ac = api_client()
        return PriceService(ac, PricesApi(ac), rules_cache=MagicMock())

    def test_raises_on_empty_base_currency(self) -> None:
        with pytest.raises(ValueError, match="base_currency"):
            self._service().get_historical(base_currency="", quote_currency="USD")

    def test_raises_on_empty_quote_currency(self) -> None:
        with pytest.raises(ValueError, match="quote_currency"):
            self._service().get_historical(base_currency="BTC", quote_currency="")

    def test_default_limit_is_sent(self) -> None:
        with StubTransport({}) as transport:
            result = self._service().get_historical(base_currency="BTC", quote_currency="USD")

        assert result == []
        assert transport.last.path == "/api/rest/v1/prices/BTC/USD/history"
        assert transport.last.query == [("limit", "20")]

    def test_limit_up_to_365(self) -> None:
        with StubTransport() as transport:
            self._service().get_historical("BTC", "USD", limit=365)
            with pytest.raises(ValueError, match="limit"):
                self._service().get_historical("BTC", "USD", limit=366)

        assert [r.param("limit") for r in transport.requests] == ["365"]
