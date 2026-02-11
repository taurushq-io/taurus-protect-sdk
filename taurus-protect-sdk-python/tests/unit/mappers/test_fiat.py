"""Unit tests for fiat mapper functions."""

from types import SimpleNamespace

from taurus_protect.services.fiat_service import (
    fiat_currencies_from_dto,
    fiat_currency_from_dto,
    fiat_provider_account_from_dto,
    fiat_provider_accounts_from_dto,
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
                id="USD", code="USD", currency=None, name="US Dollar",
                symbol="$", decimals=2, enabled=True,
            ),
            SimpleNamespace(
                id="EUR", code="EUR", currency=None, name="Euro",
                symbol=None, decimals=2, enabled=True,
            ),
        ]
        result = fiat_currencies_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert fiat_currencies_from_dto(None) == []


class TestFiatProviderAccountFromDto:
    """Tests for fiat_provider_account_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="acc-1",
            name="Main Account",
            provider="bank-provider",
            currency_code="USD",
            currencyCode=None,
            balance="10000.00",
            enabled=True,
        )
        result = fiat_provider_account_from_dto(dto)
        assert result is not None
        assert result.id == "acc-1"
        assert result.name == "Main Account"
        assert result.provider == "bank-provider"
        assert result.currency_code == "USD"
        assert result.balance == "10000.00"
        assert result.enabled is True

    def test_returns_none_for_none(self) -> None:
        assert fiat_provider_account_from_dto(None) is None

    def test_currency_code_camelcase_fallback(self) -> None:
        dto = SimpleNamespace(
            id="acc-2",
            name="Backup",
            provider="other",
            currency_code=None,
            currencyCode="EUR",
            balance=None,
            enabled=None,
        )
        result = fiat_provider_account_from_dto(dto)
        assert result is not None
        assert result.currency_code == "EUR"
        # safe_bool(None) returns False
        assert result.enabled is False


class TestFiatProviderAccountsFromDto:
    """Tests for fiat_provider_accounts_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="acc-1", name="A", provider="p", currency_code="USD",
                currencyCode=None, balance="100", enabled=True,
            ),
        ]
        result = fiat_provider_accounts_from_dto(dtos)
        assert len(result) == 1

    def test_returns_empty_for_none(self) -> None:
        assert fiat_provider_accounts_from_dto(None) == []
