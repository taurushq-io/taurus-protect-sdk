"""Unit tests for fiat mapper functions."""

from types import SimpleNamespace

from taurus_protect._internal.openapi.models import (
    TgvalidatordFiatProviderAccount,
    TgvalidatordFiatProviderEntity,
)
from taurus_protect.services.fiat_service import (
    fiat_currencies_from_dto,
    fiat_currency_from_dto,
    fiat_provider_account_from_dto,
    fiat_provider_accounts_from_dto,
    fiat_provider_entity_from_dto,
)


class TestFiatCurrencyFromDto:
    """Tests for fiat_currency_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="USD",
            code="USD",
            currency=None,
            name="US Dollar",
            symbol="$",
            decimals=2,
            enabled=True,
        )
        result = fiat_currency_from_dto(dto)
        assert result is not None
        assert result.id == "USD"
        assert result.code == "USD"
        assert result.name == "US Dollar"
        assert result.symbol == "$"
        assert result.decimals == 2
        assert result.enabled is True

    def test_returns_none_for_none(self) -> None:
        assert fiat_currency_from_dto(None) is None

    def test_code_fallback_to_currency(self) -> None:
        dto = SimpleNamespace(
            id="EUR",
            code=None,
            currency="EUR",
            name="Euro",
            symbol=None,
            decimals=2,
            enabled=True,
        )
        result = fiat_currency_from_dto(dto)
        assert result is not None
        assert result.code == "EUR"

    def test_code_fallback_to_id(self) -> None:
        dto = SimpleNamespace(
            id="CHF",
            code=None,
            currency=None,
            name="Swiss Franc",
            symbol=None,
            decimals=2,
            enabled=None,
        )
        result = fiat_currency_from_dto(dto)
        assert result is not None
        assert result.code == "CHF"
        # safe_bool(None) returns False
        assert result.enabled is False


class TestFiatCurrenciesFromDto:
    """Tests for fiat_currencies_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="USD",
                code="USD",
                currency=None,
                name="US Dollar",
                symbol="$",
                decimals=2,
                enabled=True,
            ),
            SimpleNamespace(
                id="EUR",
                code="EUR",
                currency=None,
                name="Euro",
                symbol=None,
                decimals=2,
                enabled=True,
            ),
        ]
        result = fiat_currencies_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert fiat_currencies_from_dto(None) == []


class TestFiatProviderAccountFromDto:
    """fiat_provider_account_from_dto reads the generated account's own fields."""

    def test_maps_all_fields(self) -> None:
        dto = TgvalidatordFiatProviderAccount.from_dict(
            {
                "id": "acc-1",
                "provider": "bank-provider",
                "label": "main",
                "accountType": "CURRENT",
                "accountIdentifier": "CH00",
                "accountName": "Main Account",
                "totalBalance": "10000.00",
                "currencyID": "cur-usd",
                "currencyInfo": {"symbol": "USD"},
                "baseCurrencyValuation": "9000",
                "creationDate": "2026-01-02T03:04:05Z",
            }
        )
        result = fiat_provider_account_from_dto(dto)
        assert result is not None
        assert result.id == "acc-1"
        assert result.name == "Main Account"
        assert result.provider == "bank-provider"
        assert result.label == "main"
        assert result.account_type == "CURRENT"
        assert result.account_identifier == "CH00"
        assert result.currency_id == "cur-usd"
        assert result.currency_code == "USD"
        assert result.balance == "10000.00"
        assert result.base_currency_valuation == "9000"
        assert result.created_at is not None
        assert result.enabled is True

    def test_returns_none_for_none(self) -> None:
        assert fiat_provider_account_from_dto(None) is None

    def test_currency_code_absent_without_currency_info(self) -> None:
        dto = TgvalidatordFiatProviderAccount.from_dict({"id": "acc-2", "currencyID": "c"})
        result = fiat_provider_account_from_dto(dto)
        assert result is not None
        assert result.currency_id == "c"
        assert result.currency_code is None
        assert result.balance is None


class TestFiatProviderAccountsFromDto:
    """Tests for fiat_provider_accounts_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [TgvalidatordFiatProviderAccount.from_dict({"id": "acc-1"})]
        result = fiat_provider_accounts_from_dto(dtos)
        assert [a.id for a in result] == ["acc-1"]

    def test_returns_empty_for_none(self) -> None:
        assert fiat_provider_accounts_from_dto(None) == []


class TestFiatProviderEntityFromDto:
    """fiat_provider_entity_from_dto reads the generated entity's own fields."""

    def test_maps_all_fields(self) -> None:
        dto = TgvalidatordFiatProviderEntity.from_dict(
            {
                "id": "ent-1",
                "provider": "p",
                "label": "l",
                "accountIdentifier": "acct",
                "name": "Entity",
                "details": "{}",
                "updateDate": "2026-01-02T03:04:05Z",
            }
        )
        result = fiat_provider_entity_from_dto(dto)
        assert result is not None
        assert result.id == "ent-1"
        assert (result.provider, result.label, result.name) == ("p", "l", "Entity")
        assert result.account_identifier == "acct"
        assert result.details == "{}"
        assert result.updated_at is not None

    def test_returns_none_for_none(self) -> None:
        assert fiat_provider_entity_from_dto(None) is None
