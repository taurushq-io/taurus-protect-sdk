"""Unit tests for TaurusNetwork SharingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import TaurusNetworkSharedAddressAssetApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.taurus_network.sharing import (
    ListSharedAddressesOptions,
    ListSharedAssetsOptions,
)
from taurus_protect.services.taurus_network.sharing_service import SharingService
from tests.unit.transport_stub import StubTransport, api_client


class TestListSharedAddresses:
    """list_shared_addresses reads ``sharedAddresses`` (it read ``result``, always empty)."""

    def _service(self) -> SharingService:
        ac = api_client()
        return SharingService(ac, TaurusNetworkSharedAddressAssetApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "sharedAddresses": [{"id": "sa1"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        options = ListSharedAddressesOptions(
            participant_id="p",
            owner_participant_id="o",
            target_participant_id="t",
            blockchain="ETH",
            network="mainnet",
            ids=["1"],
            statuses=["accepted"],
            sort_order="ASC",
        )
        with StubTransport(reply) as transport:
            addresses, page = self._service().list_shared_addresses(options)

        assert [a.id for a in addresses] == ["sa1"]
        assert page == CursorPage(page_size=20, next_cursor="n", has_more=True)
        assert transport.last.query == sorted(
            [
                ("participantID", "p"),
                ("ownerParticipantID", "o"),
                ("targetParticipantID", "t"),
                ("blockchain", "ETH"),
                ("network", "mainnet"),
                ("ids", "1"),
                ("statuses", "accepted"),
                ("sortOrder", "ASC"),
                ("cursor.pageSize", "20"),
            ]
        )


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
    """list_shared_assets reads ``sharedAssets`` (it read ``result``, always empty)."""

    def _service(self) -> SharingService:
        ac = api_client()
        return SharingService(ac, TaurusNetworkSharedAddressAssetApi(ac))

    def test_rows_and_continuation(self) -> None:
        with StubTransport({"sharedAssets": [{"id": "as1"}]}) as transport:
            assets, page = self._service().list_shared_assets(ListSharedAssetsOptions(cursor="n"))

        assert [a.id for a in assets] == ["as1"]
        assert page == CursorPage(page_size=20)
        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"


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
