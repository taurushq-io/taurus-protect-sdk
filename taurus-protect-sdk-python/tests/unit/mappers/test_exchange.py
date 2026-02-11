"""Unit tests for exchange mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

import pytest

from taurus_protect.services.exchange_service import exchange_from_dto, exchanges_from_dto


class TestExchangeFromDto:
    """Tests for exchange_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="ex-1",
            name="Main Binance Account",
            exchange_label="Binance",
            status="ACTIVE",
            currency_id="BTC",
            currency="BTC",
            balance="10.5",
            pending_balance="0.5",
            enabled=True,
            created_at="2024-01-15T10:30:00Z",
            updated_at="2024-06-01T12:00:00Z",
        )
        result = exchange_from_dto(dto)
        assert result is not None
        assert result.id == "ex-1"
        assert result.name == "Main Binance Account"
        assert result.exchange_label == "Binance"
        assert result.status == "ACTIVE"
        assert result.currency_id == "BTC"
        assert result.currency == "BTC"
        assert result.balance == "10.5"
        assert result.pending_balance == "0.5"
        assert result.enabled is True
        assert result.created_at is not None
        assert result.updated_at is not None

    def test_returns_none_for_none(self) -> None:
        assert exchange_from_dto(None) is None

    def test_handles_camel_case_field_names(self) -> None:
        dto = SimpleNamespace(
            id="ex-2",
            name="Kraken Account",
            exchangeLabel="Kraken",
            status="ACTIVE",
            currencyId="ETH",
            currency="ETH",
            balance="100.0",
            pendingBalance="5.0",
            enabled=True,
            createdAt="2024-02-01T00:00:00Z",
            updatedAt="2024-06-01T00:00:00Z",
        )
        result = exchange_from_dto(dto)
        assert result is not None
        assert result.exchange_label == "Kraken"
        assert result.currency_id == "ETH"
        assert result.pending_balance == "5.0"
        assert result.created_at is not None
        assert result.updated_at is not None

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="ex-3",
            name=None,
            exchange_label=None,
            status=None,
            currency_id=None,
            currency=None,
            balance=None,
            pending_balance=None,
            enabled=None,
            created_at=None,
            updated_at=None,
        )
        result = exchange_from_dto(dto)
        assert result is not None
        assert result.id == "ex-3"
        assert result.name is None
        assert result.exchange_label is None
        assert result.enabled is False

    def test_enabled_defaults_when_missing(self) -> None:
        dto = SimpleNamespace(
            id="ex-4",
            name="Account",
            exchange_label=None,
            status="ACTIVE",
            currency_id=None,
            currency=None,
            balance=None,
            pending_balance=None,
            created_at=None,
            updated_at=None,
        )
        result = exchange_from_dto(dto)
        assert result is not None
        assert result.enabled is True


class TestExchangesFromDto:
    """Tests for exchanges_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="ex-1", name="Account 1", exchange_label="Binance",
                status="ACTIVE", currency_id=None, currency=None,
                balance=None, pending_balance=None, enabled=True,
                created_at=None, updated_at=None,
            ),
            SimpleNamespace(
                id="ex-2", name="Account 2", exchange_label="Kraken",
                status="ACTIVE", currency_id=None, currency=None,
                balance=None, pending_balance=None, enabled=True,
                created_at=None, updated_at=None,
            ),
        ]
        result = exchanges_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "ex-1"
        assert result[1].id == "ex-2"

    def test_returns_empty_for_none(self) -> None:
        assert exchanges_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert exchanges_from_dto([]) == []

    def test_filters_out_none_dtos(self) -> None:
        dtos = [
            SimpleNamespace(
                id="ex-1", name="Account 1", exchange_label="Binance",
                status="ACTIVE", currency_id=None, currency=None,
                balance=None, pending_balance=None, enabled=True,
                created_at=None, updated_at=None,
            ),
            None,
        ]
        result = exchanges_from_dto(dtos)
        assert len(result) == 1
