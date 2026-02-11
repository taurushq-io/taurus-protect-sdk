"""Unit tests for statistics, price, and score mappers."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.statistics import (
    portfolio_statistics_from_dto,
    price_from_dto,
    price_history_from_dto,
    price_history_point_from_dto,
    prices_from_dto,
    score_from_dto,
    scores_from_dto,
)


class TestPriceFromDto:
    """Tests for price_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            currency_from="BTC",
            currency_to="USD",
            rate="50000.00",
            blockchain="BTC",
            decimals="8",
            change_percent24_hour="2.5",
            source="CoinGecko",
            creation_date="2024-01-15T10:30:00Z",
            update_date="2024-06-01T10:30:00Z",
        )
        result = price_from_dto(dto)
        assert result is not None
        assert result.currency_from == "BTC"
        assert result.currency_to == "USD"
        assert result.rate == "50000.00"
        assert result.blockchain == "BTC"
        assert result.source == "CoinGecko"

    def test_returns_none_for_none(self) -> None:
        assert price_from_dto(None) is None


class TestPricesFromDto:
    """Tests for prices_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert prices_from_dto(None) == []

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                currency_from="BTC", currency_to="USD", rate="50000",
                blockchain=None, decimals=None, change_percent24_hour=None,
                source=None, creation_date=None, update_date=None,
            ),
        ]
        result = prices_from_dto(dtos)
        assert len(result) == 1


class TestPriceHistoryPointFromDto:
    """Tests for price_history_point_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(
            timestamp="2024-01-15T10:30:00Z",
            rate="50000.00",
        )
        result = price_history_point_from_dto(dto)
        assert result is not None
        assert result.rate == "50000.00"
        assert result.timestamp is not None

    def test_returns_none_for_none(self) -> None:
        assert price_history_point_from_dto(None) is None


class TestPriceHistoryFromDto:
    """Tests for price_history_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert price_history_from_dto(None) == []


class TestPortfolioStatisticsFromDto:
    """Tests for portfolio_statistics_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            addresses_count="100",
            wallets_count="10",
            total_balance="500000",
            total_balance_base_currency="50000000",
            avg_balance_per_address="5000",
        )
        result = portfolio_statistics_from_dto(dto)
        assert result is not None
        assert result.addresses_count == "100"
        assert result.wallets_count == "10"
        assert result.total_balance == "500000"

    def test_returns_none_for_none(self) -> None:
        assert portfolio_statistics_from_dto(None) is None


class TestScoreFromDto:
    """Tests for score_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="score-1",
            provider="Chainalysis",
            type="risk",
            score="85",
            update_date="2024-01-15T10:30:00Z",
        )
        result = score_from_dto(dto)
        assert result is not None
        assert result.id == "score-1"
        assert result.provider == "Chainalysis"
        assert result.score_type == "risk"
        assert result.score == "85"

    def test_returns_none_for_none(self) -> None:
        assert score_from_dto(None) is None


class TestScoresFromDto:
    """Tests for scores_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert scores_from_dto(None) == []
