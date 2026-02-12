"""Unit tests for StakingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.staking_service import StakingService


class TestListValidators:
    """Tests for StakingService.list_validators()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        staking_api = MagicMock()
        service = StakingService(api_client=api_client, staking_api=staking_api)
        return service, staking_api

    def test_raises_on_empty_blockchain(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="blockchain"):
            service.list_validators(blockchain="")

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list_validators(blockchain="ETH", limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list_validators(blockchain="ETH", offset=-1)

    def test_returns_empty_for_unknown_blockchain(self) -> None:
        service, _ = self._make_service()

        validators, pagination = service.list_validators(blockchain="UNKNOWN")

        assert validators == []

    def test_calls_eth_api_for_eth_blockchain(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = []
        resp.validators = None
        api.staking_service_get_eth_validators_info.return_value = resp

        validators, pagination = service.list_validators(blockchain="ETH")

        api.staking_service_get_eth_validators_info.assert_called_once_with(
            network="mainnet",
        )


class TestGetStakingInfo:
    """Tests for StakingService.get_staking_info()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        staking_api = MagicMock()
        service = StakingService(api_client=api_client, staking_api=staking_api)
        return service, staking_api

    def test_raises_on_invalid_address_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="address_id must be positive"):
            service.get_staking_info(address_id=0)

    def test_returns_empty_info_when_no_data(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.stake_accounts = None
        api.staking_service_get_stake_accounts.return_value = resp

        info = service.get_staking_info(address_id=123)

        assert info.address_id == "123"

    def test_calls_api_with_string_address_id(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.stake_accounts = None
        api.staking_service_get_stake_accounts.return_value = resp

        service.get_staking_info(address_id=42)

        api.staking_service_get_stake_accounts.assert_called_once_with(
            address_id="42",
        )
