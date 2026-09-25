"""Unit tests for AssetService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect._internal.openapi import (
    AddressesApi,
    AddressWhitelistingApi,
    AssetsApi,
    AssetV2Api,
)
from taurus_protect.errors import IntegrityError, NotFoundError
from taurus_protect.models.address import Address
from taurus_protect.models.asset import AssetAddressV2, AssetOperationV2, AssetV2
from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.wallet import Wallet
from taurus_protect.services.address_service import AddressService
from taurus_protect.services.asset_service import AssetService
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService
from tests.unit.transport_stub import StubTransport, api_client


def _service(rules_cache: MagicMock = None) -> AssetService:
    ac = api_client()
    rules_cache = rules_cache or MagicMock()
    return AssetService(
        ac,
        AssetsApi(ac),
        rules_cache,
        AssetV2Api(ac),
        address_service=AddressService(ac, AddressesApi(ac), rules_cache),
        whitelisted_address_service=WhitelistedAddressService(
            ac,
            AddressWhitelistingApi(ac),
            [ec.generate_private_key(ec.SECP256R1()).public_key()],
            1,
        ),
    )


class TestAssetServiceList:
    """list returns the wallets holding the asset (it mapped wallets into Asset)."""

    def test_returns_wallets(self) -> None:
        reply = {
            "wallets": [{"id": "1", "name": "Treasury", "currency": "ETH"}],
            "totalItems": "3",
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            wallets, page = _service().list("ETH", page_size=1)

        assert [type(w) for w in wallets] == [Wallet]
        assert wallets[0].name == "Treasury"
        assert page == CursorPage(page_size=1, next_cursor="n", has_more=True, total_items=3)
        assert transport.last.path == "/api/rest/v1/asset/wallets"
        assert transport.last.body == {
            "asset": {"currency": "ETH"},
            "requestCursor": {"pageSize": "1"},
        }

    def test_raises_on_invalid_page_size(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                _service().list("ETH", page_size=101)

        assert transport.requests == []


class TestAssetServiceGet:
    """get raises on an empty page instead of inventing an Asset."""

    def test_raises_on_empty_id(self) -> None:
        with pytest.raises(ValueError, match="asset_id"):
            _service().get("")

    def test_not_found_on_an_empty_page(self) -> None:
        with StubTransport({}) as transport:
            with pytest.raises(NotFoundError):
                _service().get("ETH")

        assert transport.last.body == {
            "asset": {"currency": "ETH"},
            "requestCursor": {"pageSize": "1"},
        }

    def test_reads_the_asset_from_the_currency_info(self) -> None:
        reply = {
            "addresses": [
                {
                    "id": "1",
                    "currency": "USDC",
                    "currencyInfo": {
                        "id": "cur-usdc",
                        "name": "USD Coin",
                        "symbol": "USDC",
                        "blockchain": "ETH",
                        "network": "mainnet",
                        "decimals": "6",
                        "isToken": True,
                        "contractAddress": "0xa0b8",
                    },
                }
            ]
        }
        with StubTransport(reply):
            asset = _service().get("USDC")

        assert (asset.id, asset.name, asset.symbol) == ("cur-usdc", "USD Coin", "USDC")
        assert (asset.blockchain, asset.network, asset.decimals) == ("ETH", "mainnet", 6)
        assert asset.is_token is True and asset.contract_address == "0xa0b8"


class TestAssetServiceGetWallets:
    """get_wallets: request cursor only, filters in the body."""

    def test_raises_on_empty_currency(self) -> None:
        with pytest.raises(ValueError, match="currency"):
            _service().get_wallets("")

    def test_filters_and_continuation(self) -> None:
        with StubTransport() as transport:
            _service().get_wallets("ETH", cursor="n", wallet_id="7", wallet_name="T")

        assert transport.last.body == {
            "asset": {"currency": "ETH"},
            "walletId": "7",
            "walletName": "T",
            "requestCursor": {"pageSize": "20", "currentPage": "n", "pageRequest": "NEXT"},
        }


class TestAssetServiceGetAddresses:
    """get_addresses: request cursor only, every address verified."""

    def test_raises_on_empty_currency(self) -> None:
        with pytest.raises(ValueError, match="currency"):
            _service().get_addresses("")

    def test_filters_and_total(self) -> None:
        with StubTransport({"totalItems": "4"}) as transport:
            addresses, page = _service().get_addresses(
                "ETH", page_size=2, wallet_id="7", address_id="9", addresses=["0x1"]
            )

        assert addresses == []
        assert page == CursorPage(page_size=2, total_items=4)
        assert transport.last.body == {
            "asset": {"currency": "ETH"},
            "walletId": "7",
            "addressId": "9",
            "addresses": ["0x1"],
            "requestCursor": {"pageSize": "2"},
        }


class TestAssetServiceAddressVerification:
    """get_addresses returns the same Address entity AddressService verifies.

    It used to hand back the raw generated DTOs unverified, so the mandatory
    verification there was reachable around by asking for the same rows here.
    """

    def test_rules_cache_is_mandatory(self) -> None:
        with pytest.raises(ValueError):
            AssetService(api_client=MagicMock(), assets_api=MagicMock(), rules_cache=None)

    def test_every_address_signature_is_verified(self) -> None:
        rules_cache = MagicMock()
        reply = {"addresses": [{"id": "1"}, {"id": "2"}], "totalItems": "2"}

        # The mapping is covered by the address mapper tests; what matters here is that
        # every mapped address reaches the verifier.
        with (
            StubTransport(reply),
            patch(
                "taurus_protect.services.asset_service.address_from_dto",
                side_effect=lambda _dto: Address(
                    id="1", wallet_id="1", address="0x123", signature="sig1"
                ),
            ),
            patch(
                "taurus_protect.services.asset_service.verified_address",
                side_effect=lambda address, _container: address,
            ) as verify,
        ):
            addresses, _ = _service(rules_cache).get_addresses("ETH")

        assert len(addresses) == 2
        # Every mapped address goes through the SHARED seam, which is what stops this
        # path drifting from AddressService's.
        assert verify.call_count == 2
        rules_cache.get_decoded_rules_container.assert_called_once()

    def test_verification_failure_is_not_swallowed(self) -> None:
        with StubTransport({"addresses": [{"id": "1", "address": "0xEVIL"}]}):
            with pytest.raises(IntegrityError):
                _service().get_addresses("ETH")


class TestQueryAssets:
    """query_assets: the v2 asset definitions (body cursor)."""

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "result": [
                {
                    "id": "a1",
                    "symbol": "TKN",
                    "status": "ASSET_VIEW_STATUS_V2_ACTIVE",
                    "attributes": [{"key": "k", "value": "v"}],
                }
            ],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            assets, page = _service().query_assets(
                blockchain="CANTON",
                network="mainnet",
                symbol="TKN",
                contract_address="c",
                label="l",
                currency_name="Token",
                page_size=5,
            )

        assert assets == [
            AssetV2(
                id="a1", symbol="TKN", status="ASSET_VIEW_STATUS_V2_ACTIVE", attributes={"k": "v"}
            )
        ]
        assert page == CursorPage(page_size=5, next_cursor="n", has_more=True)
        assert transport.last.method == "POST"
        assert transport.last.path == "/api/rest/v2/assets/query"
        assert transport.last.body == {
            "cursor": {"pageSize": "5"},
            "blockchain": "CANTON",
            "network": "mainnet",
            "symbol": "TKN",
            "contractAddress": "c",
            "label": "l",
            "currencyName": "Token",
        }


class TestQueryAssetAddresses:
    """query_asset_addresses: the page and its filters; verification is in
    test_asset_holders_verification.py."""

    def test_rows_and_filters(self) -> None:
        reply = {
            "result": [
                {"address": "0x1", "balance": "5", "addressType": "ADDRESS_TYPE_V2_EXTERNAL"}
            ]
        }
        with StubTransport(reply) as transport:
            result = _service().query_asset_addresses(
                "a1", address_type="ADDRESS_TYPE_V2_EXTERNAL", kyc_status="KYC_STATUS_V2_APPROVED"
            )

        assert result.addresses == [
            AssetAddressV2(address="0x1", balance="5", address_type="ADDRESS_TYPE_V2_EXTERNAL")
        ]
        assert result.excluded_unverified == []
        assert result.page == CursorPage(page_size=20)
        assert transport.last.path == "/api/rest/v2/assets/a1/addresses/query"
        assert transport.last.body == {
            "cursor": {"pageSize": "20"},
            "addressType": "ADDRESS_TYPE_V2_EXTERNAL",
            "kycStatus": "KYC_STATUS_V2_APPROVED",
        }

    def test_an_unknown_filter_value_is_sent_verbatim(self) -> None:
        """validatord judges a value it does not know; the SDK neither refuses nor substitutes."""
        with StubTransport({}) as transport:
            _service().query_asset_addresses(
                "a1", address_type="INTERNAL", kyc_status="KYC_STATUS_V2_FUTURE"
            )

        assert transport.last.body == {
            "cursor": {"pageSize": "20"},
            "addressType": "INTERNAL",
            "kycStatus": "KYC_STATUS_V2_FUTURE",
        }


class TestListAssetOperations:
    """list_asset_operations: query cursor, type/status filters."""

    def test_rows_and_filters(self) -> None:
        reply = {
            "result": [{"id": "op1", "assetID": "a1", "type": "ASSET_OPERATION_TYPE_V2_MINT"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            operations, page = _service().list_asset_operations(
                "a1",
                type="ASSET_OPERATION_TYPE_V2_MINT",
                status="ASSET_OPERATION_STATUS_V2_PENDING",
            )

        assert operations == [
            AssetOperationV2(id="op1", asset_id="a1", type="ASSET_OPERATION_TYPE_V2_MINT")
        ]
        assert page == CursorPage(page_size=20, next_cursor="n", has_more=True)
        assert transport.last.path == "/api/rest/v2/assets/a1/operations"
        assert transport.last.query == sorted(
            [
                ("type", "ASSET_OPERATION_TYPE_V2_MINT"),
                ("status", "ASSET_OPERATION_STATUS_V2_PENDING"),
                ("cursor.pageSize", "20"),
            ]
        )
