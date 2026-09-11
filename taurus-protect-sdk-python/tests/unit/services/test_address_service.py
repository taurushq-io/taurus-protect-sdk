"""Unit tests for AddressService."""

from __future__ import annotations

from types import SimpleNamespace
from typing import Optional
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.errors import IntegrityError
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
        # Verification now lives inside the ONE construction seam (_verified_address)
        # rather than in a private method a test could neuter, so that is what these
        # plumbing tests stub. The seam's own behaviour is covered by
        # TestVerifiedAddressSeam below.
        service._verified_address = MagicMock()
        return service, addresses_api

    def test_get_returns_address(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_get_address.return_value = reply

        mock_address = MagicMock()
        service._verified_address.return_value = mock_address

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

        service.get(1)

        # get() must not construct an Address by any route other than the seam.
        service._verified_address.assert_called_once()


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
        service._verified_address = MagicMock()
        return service, addresses_api, rules_cache

    def test_list_returns_addresses_and_pagination(self) -> None:
        service, api, rules_cache = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "5"
        reply.offset = "0"
        api.wallet_service_get_addresses.return_value = reply

        service._verified_address.return_value = MagicMock()

        addresses, pagination = service.list(wallet_id=1)

        assert len(addresses) == 1
        # One row in, one trip through the seam: the list path must not map rows by
        # any other route.
        service._verified_address.assert_called_once()

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
        # Same seam stub as the read tests: create now shares their construction path.
        service._verified_address = MagicMock()
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
        service._verified_address.return_value = mock_address

        result = service.create_address(wallet_id=1, label="test", comment="comment")

        # The point of the fix: create goes through the SAME seam as every read, so it
        # cannot hand back an address the HSM signature has not cleared. This test used
        # to stub a reply with no signature at all and assert the Address came back --
        # i.e. it pinned the vulnerability. See TestVerifiedAddressSeam for the refusal.
        assert result is mock_address
        service._verified_address.assert_called_once()

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


class TestVerifiedAddressSeam:
    """
    The ONE construction seam for an Address, driven directly.

    ``create_address`` was the last path that returned an Address without checking the
    HSM signature, and it is the highest-value moment for substitution: its caller is
    about to publish or fund a fresh deposit address. ``get_addresses`` had already been
    fixed for exactly this and the create path was missed anyway, which is why the fix is
    a seam rather than a per-path check.
    """

    def _make_service(self) -> AddressService:
        return AddressService(
            api_client=MagicMock(),
            addresses_api=MagicMock(),
            rules_cache=MagicMock(),
        )

    @staticmethod
    def _dto(address: str, signature: Optional[str], status: str) -> SimpleNamespace:
        """
        A plain attribute holder rather than a MagicMock.

        The mapper reads nested DTOs (``balance``) and hands them to pydantic models with
        typed str fields, so a MagicMock -- whose every attribute is another MagicMock --
        fails validation before the seam is reached, and the test would pass or fail for
        the wrong reason.
        """
        return SimpleNamespace(
            id="42",
            wallet_id="1",
            address=address,
            signature=signature,
            status=status,
            balance=None,
            attributes=None,
            linked_whitelisted_address_ids=None,
        )

    def test_refuses_an_address_string_with_no_signature(self) -> None:
        """
        A non-empty address with no signature is NOT a silent skip. Returning it would
        hand the caller an attacker-controllable destination in the same type as a
        verified one, which is the whole defect.
        """
        service = self._make_service()

        dto = self._dto("0xATTACKERCONTROLLED", None, "created")

        with pytest.raises(IntegrityError) as exc:
            service._verified_address(dto)

        assert "signature" in str(exc.value).lower()

    def test_allows_a_pending_address_with_no_address_string(self) -> None:
        """
        Asynchronous creation is legitimate: ``creating`` with no address yet carries no
        destination, so there is nothing to verify and nothing to misuse. This is why the
        seam branches on the address STRING rather than on the server-controlled status.
        """
        service = self._make_service()

        dto = self._dto("", None, "creating")

        result = service._verified_address(dto)

        assert result is not None
        assert result.address == ""
        # The status is surfaced so the caller knows to re-read.
        assert result.status == "creating"
