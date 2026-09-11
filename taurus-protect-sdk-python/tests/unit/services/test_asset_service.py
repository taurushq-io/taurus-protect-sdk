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
        service = AssetService(
            api_client=api_client, assets_api=assets_api, rules_cache=MagicMock()
        )
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
        service = AssetService(
            api_client=api_client, assets_api=assets_api, rules_cache=MagicMock()
        )
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
        service = AssetService(
            api_client=api_client, assets_api=assets_api, rules_cache=MagicMock()
        )
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
        service = AssetService(
            api_client=api_client, assets_api=assets_api, rules_cache=MagicMock()
        )
        return service, assets_api

    def test_raises_on_empty_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="currency"):
            service.get_addresses(currency="")

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.get_addresses(currency="BTC", offset=-1)


class TestAssetServiceAddressVerification:
    """get_addresses returns the same Address entity AddressService verifies.

    It used to hand back the raw generated DTOs unverified, so the mandatory
    verification there was reachable around by asking for the same rows here.
    """

    def test_rules_cache_is_mandatory(self) -> None:
        with pytest.raises(ValueError):
            AssetService(
                api_client=MagicMock(), assets_api=MagicMock(), rules_cache=None
            )

    def test_every_address_signature_is_verified(self) -> None:
        from unittest.mock import patch

        from taurus_protect.models.address import Address

        assets_api = MagicMock()
        reply = MagicMock()
        reply.result = [MagicMock(), MagicMock()]
        reply.total_items = "2"
        assets_api.wallet_service_get_asset_addresses.return_value = reply

        rules_cache = MagicMock()
        service = AssetService(
            api_client=MagicMock(), assets_api=assets_api, rules_cache=rules_cache
        )

        # The mapping is covered by the address mapper tests; what matters here is that
        # every mapped address reaches the verifier.
        with patch(
            "taurus_protect.services.asset_service.address_from_dto",
            side_effect=lambda _dto: Address(
                id="1", wallet_id="1", address="0x123", signature="sig1"
            ),
        ), patch(
            "taurus_protect.services.asset_service.verified_address",
            side_effect=lambda address, _container: address,
        ) as verify:
            addresses, _ = service.get_addresses("ETH")

        assert len(addresses) == 2
        # Every mapped address goes through the SHARED seam, which is what stops this
        # path drifting from AddressService's.
        assert verify.call_count == 2
        rules_cache.get_decoded_rules_container.assert_called_once()

    def test_verification_failure_is_not_swallowed(self) -> None:
        from unittest.mock import patch

        from taurus_protect.errors import IntegrityError

        from taurus_protect.models.address import Address

        assets_api = MagicMock()
        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "1"
        assets_api.wallet_service_get_asset_addresses.return_value = reply

        service = AssetService(
            api_client=MagicMock(), assets_api=assets_api, rules_cache=MagicMock()
        )

        with patch(
            "taurus_protect.services.asset_service.address_from_dto",
            side_effect=lambda _dto: Address(id="1", wallet_id="1", address="0xEVIL"),
        ), patch(
            "taurus_protect.services.asset_service.verified_address",
            side_effect=IntegrityError("address signature verification failed"),
        ):
            with pytest.raises(IntegrityError):
                service.get_addresses("ETH")
