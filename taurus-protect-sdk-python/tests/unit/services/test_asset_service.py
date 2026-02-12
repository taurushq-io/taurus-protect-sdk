"""Unit tests for AssetService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.asset_service import AssetService


class TestAssetServiceList:
    """Tests for AssetService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        assets_api = MagicMock()
        service = AssetService(api_client=api_client, assets_api=assets_api)
        return service, assets_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.wallets = None
        resp.total_items = None
        resp.totalItems = None
        resp.offset = None
        api.wallet_service_get_asset_wallets.return_value = resp

        assets, pagination = service.list()

        assert assets == []


class TestAssetServiceGet:
    """Tests for AssetService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        assets_api = MagicMock()
        service = AssetService(api_client=api_client, assets_api=assets_api)
        return service, assets_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="asset_id"):
            service.get(asset_id="")

    def test_returns_basic_asset_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.addresses = None
        api.wallet_service_get_asset_addresses.return_value = resp

        asset = service.get("BTC")

        assert asset.id == "BTC"
        assert asset.symbol == "BTC"


class TestAssetServiceGetWallets:
    """Tests for AssetService.get_wallets()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        assets_api = MagicMock()
        service = AssetService(api_client=api_client, assets_api=assets_api)
        return service, assets_api

    def test_raises_on_empty_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="currency"):
            service.get_wallets(currency="")

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.get_wallets(currency="BTC", limit=0)


class TestAssetServiceGetAddresses:
    """Tests for AssetService.get_addresses()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        assets_api = MagicMock()
        service = AssetService(api_client=api_client, assets_api=assets_api)
        return service, assets_api

    def test_raises_on_empty_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="currency"):
            service.get_addresses(currency="")

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.get_addresses(currency="BTC", offset=-1)
