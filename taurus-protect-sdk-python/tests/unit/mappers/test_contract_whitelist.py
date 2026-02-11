"""Unit tests for contract whitelisting mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

import pytest

from taurus_protect.services.contract_whitelisting_service import (
    ContractWhitelistingService,
    WhitelistedContract,
)


class TestMapContractFromDto:
    """Tests for ContractWhitelistingService._map_contract_from_dto static method."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="42",
            address="0x1234567890abcdef1234567890abcdef12345678",
            name="Uniswap V3 Router",
            blockchain="ETH",
            network="mainnet",
            abi='[{"inputs":[],"name":"factory","outputs":[]}]',
            status="APPROVED",
            created_at="2024-01-15T10:30:00Z",
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.id == "42"
        assert result.address == "0x1234567890abcdef1234567890abcdef12345678"
        assert result.name == "Uniswap V3 Router"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.abi is not None
        assert result.status == "APPROVED"
        assert result.created_at == "2024-01-15T10:30:00Z"

    def test_handles_integer_id(self) -> None:
        dto = SimpleNamespace(
            id=99,
            address="0xabc",
            name="Test Contract",
            blockchain="ETH",
            network="mainnet",
            abi=None,
            status=None,
            created_at=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.id == "99"

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="1",
            address=None,
            name=None,
            blockchain=None,
            network=None,
            abi=None,
            status=None,
            created_at=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.id == "1"
        assert result.address is None
        assert result.name is None
        assert result.blockchain is None
        assert result.network is None
        assert result.abi is None
        assert result.status is None
        assert result.created_at is None

    def test_handles_camel_case_created_at(self) -> None:
        dto = SimpleNamespace(
            id="5",
            address="0xdef",
            name="Contract",
            blockchain="SOL",
            network="devnet",
            abi=None,
            status="PENDING",
            createdAt="2024-03-20T08:00:00Z",
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.created_at == "2024-03-20T08:00:00Z"

    def test_handles_missing_id(self) -> None:
        dto = SimpleNamespace(
            address="0x123",
            name="No ID Contract",
            blockchain="ETH",
            network="mainnet",
            abi=None,
            status=None,
            created_at=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.id == ""


class TestWhitelistedContractModel:
    """Tests for the WhitelistedContract model itself."""

    def test_creates_with_all_fields(self) -> None:
        contract = WhitelistedContract(
            id="1",
            address="0xabc",
            name="Test",
            blockchain="ETH",
            network="mainnet",
            abi="[]",
            status="APPROVED",
            created_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
        )
        assert contract.id == "1"
        assert contract.address == "0xabc"
        assert contract.name == "Test"
        assert contract.blockchain == "ETH"
        assert contract.status == "APPROVED"

    def test_default_values(self) -> None:
        contract = WhitelistedContract(id="1")
        assert contract.address is None
        assert contract.name is None
        assert contract.blockchain is None
        assert contract.network is None
        assert contract.abi is None
        assert contract.status is None
        assert contract.created_at is None
