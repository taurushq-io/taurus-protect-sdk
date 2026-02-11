"""Unit tests for BusinessRuleService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.models.business_rule import BusinessRuleResult
from taurus_protect.services.business_rule_service import BusinessRuleService


class TestBusinessRuleServiceList:
    """Tests for BusinessRuleService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        business_rules_api = MagicMock()
        service = BusinessRuleService(
            api_client=api_client, business_rules_api=business_rules_api
        )
        return service, business_rules_api

    def test_returns_business_rule_result(self) -> None:
        from types import SimpleNamespace

        service, api = self._make_service()
        dto = SimpleNamespace(
            id="42",
            tenant_id="5",
            currency="XLM",
            wallet_id=None,
            address_id=None,
            rule_key="max_amount",
            rule_value="1000",
            rule_group="limits",
            rule_description="Max amount",
            rule_validation=None,
            entity_type="currency",
            entity_id="10",
            currency_info=None,
        )

        cursor_obj = MagicMock()
        cursor_obj.current_page = "page1"
        cursor_obj.has_next = True
        reply = MagicMock()
        reply.result = [dto]
        reply.cursor = cursor_obj
        api.rule_service_get_business_rules_v2.return_value = reply

        result = service.list(page_size=50)

        assert isinstance(result, BusinessRuleResult)
        assert len(result.rules) == 1
        assert result.rules[0].id == "42"
        assert result.rules[0].rule_key == "max_amount"
        assert result.has_next is True
        assert result.current_page == "page1"

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.cursor = None
        api.rule_service_get_business_rules_v2.return_value = reply

        result = service.list()

        assert result.rules == []
        assert result.has_next is False

    def test_calls_v2_api(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.cursor = None
        api.rule_service_get_business_rules_v2.return_value = reply

        service.list(page_size=25, page_request="NEXT", current_page="abc")

        api.rule_service_get_business_rules_v2.assert_called_once_with(
            ids=None,
            rule_keys=None,
            rule_groups=None,
            wallet_ids=None,
            currency_ids=None,
            address_ids=None,
            level=None,
            cursor_current_page="abc",
            cursor_page_request="NEXT",
            cursor_page_size="25",
            entity_type=None,
            entity_ids=None,
        )


class TestBusinessRuleServiceListByWallet:
    """Tests for BusinessRuleService.list_by_wallet()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        business_rules_api = MagicMock()
        service = BusinessRuleService(
            api_client=api_client, business_rules_api=business_rules_api
        )
        return service, business_rules_api

    def test_raises_on_invalid_wallet_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.list_by_wallet(wallet_id=0)

    def test_passes_wallet_id(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.cursor = None
        api.rule_service_get_business_rules_v2.return_value = reply

        service.list_by_wallet(wallet_id=42)

        call_kwargs = api.rule_service_get_business_rules_v2.call_args
        assert call_kwargs.kwargs["wallet_ids"] == ["42"]


class TestBusinessRuleServiceListByCurrency:
    """Tests for BusinessRuleService.list_by_currency()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        business_rules_api = MagicMock()
        service = BusinessRuleService(
            api_client=api_client, business_rules_api=business_rules_api
        )
        return service, business_rules_api

    def test_raises_on_empty_currency_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="currency_id"):
            service.list_by_currency(currency_id="")

    def test_passes_currency_id(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.cursor = None
        api.rule_service_get_business_rules_v2.return_value = reply

        service.list_by_currency(currency_id="ETH")

        call_kwargs = api.rule_service_get_business_rules_v2.call_args
        assert call_kwargs.kwargs["currency_ids"] == ["ETH"]
