"""Unit tests for TaurusNetwork SharingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.taurus_network.sharing_service import SharingService


class TestListSharedAddresses:
    """Tests for SharingService.list_shared_addresses()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        shared_api = MagicMock()
        service = SharingService(api_client=api_client, shared_api=shared_api)
        return service, shared_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_shared_addresses.return_value = resp

        addresses, pagination = service.list_shared_addresses()

        assert addresses == []


class TestShareAddress:
    """Tests for SharingService.share_address()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        shared_api = MagicMock()
        service = SharingService(api_client=api_client, shared_api=shared_api)
        return service, shared_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.share_address(request=None)


class TestUnshareAddress:
    """Tests for SharingService.unshare_address()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        shared_api = MagicMock()
        service = SharingService(api_client=api_client, shared_api=shared_api)
        return service, shared_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="shared_address_id"):
            service.unshare_address(shared_address_id="")


class TestListSharedAssets:
    """Tests for SharingService.list_shared_assets()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        shared_api = MagicMock()
        service = SharingService(api_client=api_client, shared_api=shared_api)
        return service, shared_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_shared_whitelisted_assets.return_value = resp

        assets, pagination = service.list_shared_assets()

        assert assets == []


class TestUnshareWhitelistedAsset:
    """Tests for SharingService.unshare_whitelisted_asset()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        shared_api = MagicMock()
        service = SharingService(api_client=api_client, shared_api=shared_api)
        return service, shared_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="shared_asset_id"):
            service.unshare_whitelisted_asset(shared_asset_id="")
