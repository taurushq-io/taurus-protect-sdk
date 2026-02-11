"""Unit tests for AddressService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.address_service import AddressService


class TestConstructor:
    """Tests for AddressService constructor."""

    def test_raises_for_none_rules_cache(self) -> None:
        with pytest.raises(ValueError, match="rules_cache cannot be None"):
            AddressService(
                api_client=MagicMock(),
                addresses_api=MagicMock(),
                rules_cache=None,
            )

    def test_accepts_valid_rules_cache(self) -> None:
        service = AddressService(
            api_client=MagicMock(),
            addresses_api=MagicMock(),
            rules_cache=MagicMock(),
        )
        assert service is not None


class TestGet:
    """Tests for AddressService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        # Mock signature verification to prevent it from running
        service._verify_address_signature = MagicMock()
        return service, addresses_api

    def test_get_returns_address(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_get_address.return_value = reply

        mock_address = MagicMock()
        with patch(
            "taurus_protect.services.address_service.address_from_dto",
            return_value=mock_address,
        ):
            result = service.get(1)

        assert result is mock_address
        api.wallet_service_get_address.assert_called_once_with("1")

    def test_get_raises_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address_id must be positive"):
            service.get(0)

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.wallet_service_get_address.return_value = reply

        with pytest.raises(NotFoundError):
            service.get(1)

    def test_get_calls_signature_verification(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_get_address.return_value = reply

        mock_address = MagicMock()
        with patch(
            "taurus_protect.services.address_service.address_from_dto",
            return_value=mock_address,
        ):
            service.get(1)

        service._verify_address_signature.assert_called_once()


class TestList:
    """Tests for AddressService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        service._verify_address_signature = MagicMock()
        return service, addresses_api, rules_cache

    def test_list_returns_addresses_and_pagination(self) -> None:
        service, api, rules_cache = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "5"
        reply.offset = "0"
        api.wallet_service_get_addresses.return_value = reply

        with patch(
            "taurus_protect.services.address_service.addresses_from_dto",
            return_value=[MagicMock()],
        ):
            addresses, pagination = service.list(wallet_id=1)

        assert len(addresses) == 1

    def test_list_raises_for_non_positive_wallet_id(self) -> None:
        service, _, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.list(wallet_id=0)

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(wallet_id=1, limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(wallet_id=1, offset=-1)


class TestCreateAddress:
    """Tests for AddressService.create_address()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        return service, addresses_api

    def test_create_address_raises_for_non_positive_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.create_address(wallet_id=0, label="test", comment="test")

    def test_create_address_raises_for_empty_label(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="label"):
            service.create_address(wallet_id=1, label="", comment="test")

    def test_create_address_returns_address(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_create_address.return_value = reply

        mock_address = MagicMock()
        with patch(
            "taurus_protect.services.address_service.address_from_dto",
            return_value=mock_address,
        ):
            result = service.create_address(wallet_id=1, label="test", comment="comment")

        assert result is mock_address

    def test_create_raises_for_none_request(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="request cannot be None"):
            service.create(None)


class TestCreateAttribute:
    """Tests for AddressService.create_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        return service, addresses_api

    def test_create_attribute_raises_for_non_positive_address_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address_id must be positive"):
            service.create_attribute(0, "key", "value")

    def test_create_attribute_raises_for_empty_key(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="key"):
            service.create_attribute(1, "", "value")


class TestDeleteAttribute:
    """Tests for AddressService.delete_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        return service, addresses_api

    def test_delete_attribute_raises_for_non_positive_address_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address_id must be positive"):
            service.delete_attribute(0, 1)

    def test_delete_attribute_raises_for_non_positive_attribute_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="attribute_id must be positive"):
            service.delete_attribute(1, 0)

    def test_delete_attribute_calls_api(self) -> None:
        service, api = self._make_service()

        service.delete_attribute(1, 2)

        api.wallet_service_delete_address_attribute.assert_called_once_with("1", "2")


class TestGetProofOfReserve:
    """Tests for AddressService.get_proof_of_reserve()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        addresses_api = MagicMock()
        rules_cache = MagicMock()
        service = AddressService(
            api_client=api_client,
            addresses_api=addresses_api,
            rules_cache=rules_cache,
        )
        return service, addresses_api

    def test_get_proof_of_reserve_raises_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="address_id must be positive"):
            service.get_proof_of_reserve(0)
