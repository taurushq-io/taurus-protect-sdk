"""Unit tests for price mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.mappers.statistics import (
    price_from_dto,
    price_history_from_dto,
    price_history_point_from_dto,
    prices_from_dto,
)


class TestPriceFromDto:
    """Tests for price_from_dto function."""

    def test_maps_all_fields(self) -> None:
        created = datetime(2024, 6, 1, 0, 0, 0, tzinfo=timezone.utc)
        updated = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            currency_from="BTC",
            currency_to="USD",
            rate="67500.50",
            blockchain="BTC",
            decimals="8",
            change_percent24_hour="-2.5",
            source="exchange",
            creation_date=created,
            update_date=updated,
        )
        result = price_from_dto(dto)
        assert result is not None
        assert result.currency_from == "BTC"
        assert result.currency_to == "USD"
        assert result.rate == "67500.50"
        assert result.blockchain == "BTC"
        assert result.decimals == "8"
        assert result.change_percent_24h == "-2.5"
        assert result.source == "exchange"
        assert result.created_at == created
        assert result.updated_at == updated

    def test_returns_none_for_none(self) -> None:
        assert price_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            currency_from="ETH",
            currency_to="EUR",
            rate="3200",
            blockchain=None,
            decimals=None,
            change_percent24_hour=None,
            source=None,
            creation_date=None,
            update_date=None,
        )
        result = price_from_dto(dto)
        assert result is not None
        assert result.currency_from == "ETH"
        assert result.blockchain is None
        assert result.created_at is None


class TestPricesFromDto:
    """Tests for prices_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                currency_from="BTC", currency_to="USD", rate="67000",
                blockchain=None, decimals=None, change_percent24_hour=None,
                source=None, creation_date=None, update_date=None,
            ),
            SimpleNamespace(
                currency_from="ETH", currency_to="USD", rate="3200",
                blockchain=None, decimals=None, change_percent24_hour=None,
                source=None, creation_date=None, update_date=None,
            ),
        ]
        result = prices_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert prices_from_dto(None) == []


class TestPriceHistoryPointFromDto:
    """Tests for price_history_point_from_dto function."""

    def test_maps_fields(self) -> None:
        ts = datetime(2024, 1, 1, 0, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            timestamp=ts,
            rate="42000.00",
        )
        result = price_history_point_from_dto(dto)
        assert result is not None
        assert result.timestamp == ts
        assert result.rate == "42000.00"

    def test_returns_none_for_none(self) -> None:
        assert price_history_point_from_dto(None) is None


class TestPriceHistoryFromDto:
    """Tests for price_history_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(timestamp=None, rate="100"),
            SimpleNamespace(timestamp=None, rate="200"),
        ]
        result = price_history_from_dto(dtos)
        assert len(result) == 2
        assert result[0].rate == "100"
        assert result[1].rate == "200"

    def test_returns_empty_for_none(self) -> None:
        assert price_history_from_dto(None) == []
