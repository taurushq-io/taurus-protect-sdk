"""Unit tests for BusinessRuleService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import BusinessRulesApi
from taurus_protect.models.business_rule import BusinessRuleResult
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.business_rule_service import BusinessRuleService
from tests.unit.transport_stub import StubTransport, api_client


class TestBusinessRuleServiceList:
    """BusinessRuleService.list over the real generated client (v2, cursor)."""

    def _service(self) -> BusinessRuleService:
        ac = api_client()
        return BusinessRuleService(ac, BusinessRulesApi(ac))

    def test_rows_and_page(self) -> None:
        reply = {
            "result": [{"id": "r1", "ruleKey": "k"}],
            "cursor": {"currentPage": "p2", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            result = self._service().list(rule_keys=["k"], page_size=5)

        assert isinstance(result, BusinessRuleResult)
        assert [r.id for r in result.rules] == ["r1"]
        assert result.page == CursorPage(page_size=5, next_cursor="p2", has_more=True)
        assert transport.last.path == "/api/rest/v2/businessrules"
        assert transport.last.query == [("cursor.pageSize", "5"), ("ruleKeys", "k")]

    def test_last_page_has_no_cursor(self) -> None:
        """The raw currentPage of the last page is not a cursor to continue with."""
        with StubTransport({"cursor": {"currentPage": "p9"}}):
            result = self._service().list()

        assert result.page == CursorPage(page_size=20)

    def test_every_filter_reaches_the_wire(self) -> None:
        with StubTransport() as transport:
            self._service().list(
                wallet_ids=["1"],
                currency_ids=["c"],
                entity_type="wallet",
                entity_ids=["e"],
                ids=["i"],
                rule_groups=["g"],
                address_ids=["a"],
                level="tenant",
            )

        assert dict(transport.last.query) == {
            "walletIds": "1",
            "currencyIds": "c",
            "entityType": "wallet",
            "entityIDs": "e",
            "ids": "i",
            "ruleGroups": "g",
            "addressIds": "a",
            "level": "tenant",
            "cursor.pageSize": "20",
        }

    def test_cursor_and_current_page_are_exclusive(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="cursor"):
                self._service().list(cursor="abc", current_page="def")

        assert transport.requests == []

    def test_low_level_page_request_passes_through(self) -> None:
        with StubTransport() as transport:
            self._service().list(current_page="p", page_request="PREVIOUS")

        assert transport.last.param("cursor.currentPage") == "p"
        assert transport.last.param("cursor.pageRequest") == "PREVIOUS"


class TestBusinessRuleServiceListByWallet:
    """Tests for BusinessRuleService.list_by_wallet()."""

    def _service(self) -> BusinessRuleService:
        ac = api_client()
        return BusinessRuleService(ac, BusinessRulesApi(ac))

    def test_raises_on_invalid_wallet_id(self) -> None:
        with pytest.raises(ValueError, match="wallet_id must be positive"):
            self._service().list_by_wallet(wallet_id=0)

    def test_passes_wallet_id_and_cursor(self) -> None:
        with StubTransport() as transport:
            self._service().list_by_wallet(wallet_id=42, cursor="c+/=")

        assert transport.last.param("walletIds") == "42"
        assert transport.last.param("cursor.currentPage") == "c+/="
        assert transport.last.param("cursor.pageRequest") == "NEXT"


class TestBusinessRuleServiceListByCurrency:
    """Tests for BusinessRuleService.list_by_currency()."""

    def _service(self) -> BusinessRuleService:
        ac = api_client()
        return BusinessRuleService(ac, BusinessRulesApi(ac))

    def test_raises_on_empty_currency_id(self) -> None:
        with pytest.raises(ValueError, match="currency_id"):
            self._service().list_by_currency(currency_id="")

    def test_passes_currency_id(self) -> None:
        with StubTransport() as transport:
            self._service().list_by_currency(currency_id="ETH")

        assert transport.last.param("currencyIds") == "ETH"


class TestUpdateTransactionsEnabled:
    """The transactions-enabled kill switch existed only in the Go SDK, even though the
    generated op is present in all four and tg-protect-mcpd drives it."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        business_rules_api = MagicMock()
        service = BusinessRuleService(api_client=api_client, business_rules_api=business_rules_api)
        return service, business_rules_api

    def test_sends_disabled_flag(self) -> None:
        service, api = self._make_service()

        service.update_transactions_enabled(False)

        body = api.rule_service_update_transactions_enabled_business_rule.call_args.kwargs["body"]
        assert body.enabled is False

    def test_sends_enabled_flag(self) -> None:
        service, api = self._make_service()

        service.update_transactions_enabled(True)

        body = api.rule_service_update_transactions_enabled_business_rule.call_args.kwargs["body"]
        assert body.enabled is True
