"""Unit tests for business rule mapper functions."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.business_rule import (
    business_rule_from_dto,
    business_rules_from_dto,
)
from taurus_protect.models.business_rule import BusinessRule


class TestBusinessRuleFromDto:
    """Tests for business_rule_from_dto function."""

    def test_maps_all_fields(self) -> None:
        currency_dto = SimpleNamespace(
            id="1", symbol="XLM", name="Stellar Lumens",
            blockchain=None, network=None, contract_address=None,
            decimals=None, is_token=None, is_erc20=None,
            is_utxo_based=None, is_account_based=None,
            is_fiat=None, is_fa12=None, is_fa20=None, is_nft=None,
        )
        dto = SimpleNamespace(
            id="42",
            tenant_id="5",
            currency="XLM",
            wallet_id="100",
            address_id="200",
            rule_key="max_outgoing_transaction_per_day",
            rule_value="1000",
            rule_group="limits",
            rule_description="Max outgoing per day",
            rule_validation="^[0-9]+$",
            entity_type="currency",
            entity_id="300",
            currency_info=currency_dto,
        )
        result = business_rule_from_dto(dto)
        assert result is not None
        assert result.id == "42"
        assert result.tenant_id == 5
        assert result.currency == "XLM"
        assert result.wallet_id == "100"
        assert result.address_id == "200"
        assert result.rule_key == "max_outgoing_transaction_per_day"
        assert result.rule_value == "1000"
        assert result.rule_group == "limits"
        assert result.rule_description == "Max outgoing per day"
        assert result.rule_validation == "^[0-9]+$"
        assert result.entity_type == "currency"
        assert result.entity_id == "300"
        assert result.currency_info is not None

    def test_maps_camel_case_fields(self) -> None:
        dto = SimpleNamespace(
            id="10",
            tenantId="3",
            currency="ETH",
            walletId="50",
            addressId="60",
            ruleKey="max_amount",
            ruleValue="500",
            ruleGroup="security",
            ruleDescription="Max amount",
            ruleValidation=None,
            entityType="wallet",
            entityID="70",
            currencyInfo=None,
        )
        result = business_rule_from_dto(dto)
        assert result is not None
        assert result.tenant_id == 3
        assert result.wallet_id == "50"
        assert result.address_id == "60"
        assert result.rule_key == "max_amount"
        assert result.rule_value == "500"
        assert result.rule_group == "security"
        assert result.entity_type == "wallet"
        assert result.entity_id == "70"

    def test_returns_none_for_none(self) -> None:
        assert business_rule_from_dto(None) is None

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="1",
            tenant_id=None,
            currency=None,
            wallet_id=None,
            address_id=None,
            rule_key=None,
            rule_value=None,
            rule_group=None,
            rule_description=None,
            rule_validation=None,
            entity_type=None,
            entity_id=None,
            currency_info=None,
        )
        result = business_rule_from_dto(dto)
        assert result is not None
        assert result.id == "1"
        assert result.tenant_id == 0
        assert result.currency is None
        assert result.rule_key is None

    def test_handles_integer_id(self) -> None:
        dto = SimpleNamespace(
            id=123,
            tenant_id=None,
            currency=None,
            wallet_id=None,
            address_id=None,
            rule_key=None,
            rule_value=None,
            rule_group=None,
            rule_description=None,
            rule_validation=None,
            entity_type=None,
            entity_id=None,
            currency_info=None,
        )
        result = business_rule_from_dto(dto)
        assert result is not None
        assert result.id == "123"


class TestBusinessRulesFromDto:
    """Tests for business_rules_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", tenant_id=None, currency=None, wallet_id=None,
                address_id=None, rule_key="key1", rule_value="val1",
                rule_group=None, rule_description=None, rule_validation=None,
                entity_type=None, entity_id=None, currency_info=None,
            ),
            SimpleNamespace(
                id="2", tenant_id=None, currency=None, wallet_id=None,
                address_id=None, rule_key="key2", rule_value="val2",
                rule_group=None, rule_description=None, rule_validation=None,
                entity_type=None, entity_id=None, currency_info=None,
            ),
        ]
        result = business_rules_from_dto(dtos)
        assert len(result) == 2
        assert result[0].rule_key == "key1"
        assert result[1].rule_key == "key2"

    def test_returns_empty_for_none(self) -> None:
        assert business_rules_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert business_rules_from_dto([]) == []


class TestBusinessRuleModel:
    """Tests for the BusinessRule Pydantic model."""

    def test_creates_with_all_fields(self) -> None:
        rule = BusinessRule(
            id="1",
            tenant_id=5,
            currency="BTC",
            rule_key="max_amount",
            rule_value="1000",
            rule_group="limits",
        )
        assert rule.id == "1"
        assert rule.tenant_id == 5
        assert rule.currency == "BTC"
        assert rule.rule_key == "max_amount"

    def test_default_values(self) -> None:
        rule = BusinessRule()
        assert rule.id == ""
        assert rule.tenant_id == 0
        assert rule.currency is None
        assert rule.rule_key is None
        assert rule.rule_value is None

    def test_is_frozen(self) -> None:
        rule = BusinessRule(id="1")
        with pytest.raises(Exception):
            rule.id = "2"
