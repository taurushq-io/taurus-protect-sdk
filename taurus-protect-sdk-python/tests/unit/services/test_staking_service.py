"""Unit tests for StakingService."""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi import StakingApi
from taurus_protect.services.staking_service import StakingService
from tests.unit.transport_stub import StubTransport, api_client


class TestListValidators:
    """list_validators: the ETH endpoint does not page, so every validator comes back."""

    def _service(self) -> StakingService:
        ac = api_client()
        return StakingService(ac, StakingApi(ac))

    def test_raises_on_empty_blockchain(self) -> None:
        with pytest.raises(ValueError, match="blockchain"):
            self._service().list_validators(blockchain="")

    def test_returns_empty_for_unknown_blockchain(self) -> None:
        with StubTransport() as transport:
            validators = self._service().list_validators(blockchain="UNKNOWN")

        assert validators == []
        assert transport.requests == []

    def test_eth_returns_every_validator_from_the_validators_field(self) -> None:
        reply = {"validators": [{"id": f"v{i}", "pubkey": f"0x{i}"} for i in range(25)]}
        with StubTransport(reply) as transport:
            validators = self._service().list_validators(blockchain="ETH")

        assert [v.id for v in validators] == [f"v{i}" for i in range(25)]
        assert transport.last.path == "/api/rest/v1/staking/eth/mainnet/validators"
        assert transport.last.query == []

    def test_ids_reach_the_wire(self) -> None:
        with StubTransport({}) as transport:
            self._service().list_validators(blockchain="ETH", network="holesky", ids=["v1", "v2"])

        assert transport.last.path == "/api/rest/v1/staking/eth/holesky/validators"
        assert transport.last.query == [("ids", "v1"), ("ids", "v2")]


class TestGetStakingInfo:
    """get_staking_info reads the first stake account, asking for a page of one."""

    def _service(self) -> StakingService:
        ac = api_client()
        return StakingService(ac, StakingApi(ac))

    def test_raises_on_invalid_address_id(self) -> None:
        with pytest.raises(ValueError, match="address_id must be positive"):
            self._service().get_staking_info(address_id=0)

    def test_returns_empty_info_when_no_data(self) -> None:
        with StubTransport({}) as transport:
            info = self._service().get_staking_info(address_id=123)

        assert info.address_id == "123"
        assert transport.last.query == [("addressId", "123"), ("cursor.pageSize", "1")]

    def test_reads_the_stake_accounts_field(self) -> None:
        with StubTransport({"stakeAccounts": [{"id": "sa1", "addressId": "42"}]}):
            info = self._service().get_staking_info(address_id=42)

        assert info.address_id == "42"
