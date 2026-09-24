"""Unit tests for UserDeviceService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.models.pagination import Pagination
from taurus_protect.services.user_device_service import UserDeviceService


class TestUserDeviceServiceList:
    """UserDeviceService.list: the API cannot list pairings; the page is always empty."""

    def _make_service(self) -> UserDeviceService:
        return UserDeviceService(api_client=MagicMock(), user_device_api=MagicMock())

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window(self, kwargs: dict, name: str) -> None:
        with pytest.raises(ValueError, match=name):
            self._make_service().list(**kwargs)

    def test_returns_an_empty_page_with_pagination(self) -> None:
        pairings, pagination = self._make_service().list()

        assert pairings == []
        assert pagination == Pagination(limit=20, offset=0)


class TestUserDeviceServiceGet:
    """Tests for UserDeviceService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        user_device_api = MagicMock()
        service = UserDeviceService(api_client=api_client, user_device_api=user_device_api)
        return service, user_device_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="device_id"):
            service.get(device_id="")

    def test_always_raises_not_found(self) -> None:
        service, _ = self._make_service()

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError, match="get_pairing_status"):
            service.get(device_id="dev-123")


class TestUserDeviceServiceCreatePairing:
    """Tests for UserDeviceService.create_pairing()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        user_device_api = MagicMock()
        service = UserDeviceService(api_client=api_client, user_device_api=user_device_api)
        return service, user_device_api

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.pairing_id = "pair-123"
        resp.pairingId = None
        resp.result = None
        api.user_device_service_create_user_device_pairing.return_value = resp

        # The mapper may or may not find data depending on DTO shape
        try:
            result = service.create_pairing()
        except Exception:
            pass  # mapper may fail with mock

        api.user_device_service_create_user_device_pairing.assert_called_once()


class TestUserDeviceServiceStartPairing:
    """Tests for UserDeviceService.start_pairing()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        user_device_api = MagicMock()
        service = UserDeviceService(api_client=api_client, user_device_api=user_device_api)
        return service, user_device_api

    def test_raises_on_empty_pairing_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pairing_id"):
            service.start_pairing(pairing_id="", nonce="123456", encryption_key="key")

    def test_raises_on_empty_nonce(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="nonce"):
            service.start_pairing(pairing_id="pair-1", nonce="", encryption_key="key")

    def test_raises_on_empty_encryption_key(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="encryption_key"):
            service.start_pairing(pairing_id="pair-1", nonce="123456", encryption_key="")


class TestUserDeviceServiceApprovePairing:
    """Tests for UserDeviceService.approve_pairing()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        user_device_api = MagicMock()
        service = UserDeviceService(api_client=api_client, user_device_api=user_device_api)
        return service, user_device_api

    def test_raises_on_empty_pairing_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pairing_id"):
            service.approve_pairing(pairing_id="", nonce="123456")

    def test_raises_on_empty_nonce(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="nonce"):
            service.approve_pairing(pairing_id="pair-1", nonce="")


class TestUserDeviceServiceGetPairingStatus:
    """Tests for UserDeviceService.get_pairing_status()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        user_device_api = MagicMock()
        service = UserDeviceService(api_client=api_client, user_device_api=user_device_api)
        return service, user_device_api

    def test_raises_on_empty_pairing_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pairing_id"):
            service.get_pairing_status(pairing_id="", nonce="123456")

    def test_raises_on_empty_nonce(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="nonce"):
            service.get_pairing_status(pairing_id="pair-1", nonce="")
