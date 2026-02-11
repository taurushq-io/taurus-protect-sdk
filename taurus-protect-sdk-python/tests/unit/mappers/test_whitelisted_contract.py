"""Unit tests for whitelisted contract mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.services.contract_whitelisting_service import (
    ContractWhitelistingService,
    WhitelistedContract,
)


class TestMapContractFromDto:
    """Tests for ContractWhitelistingService._map_contract_from_dto."""

    def test_maps_all_fields(self) -> None:
        created = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id=42,
            address="0x1234567890abcdef",
            name="USDC Contract",
            blockchain="ETH",
            network="mainnet",
            abi='[{"type":"function","name":"transfer"}]',
            status="APPROVED",
            created_at=created,
            createdAt=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert isinstance(result, WhitelistedContract)
        assert result.id == "42"
        assert result.address == "0x1234567890abcdef"
        assert result.name == "USDC Contract"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.abi is not None
        assert result.status == "APPROVED"
        assert result.created_at == created

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id=1,
            address=None,
            name=None,
            blockchain=None,
            network=None,
            abi=None,
            status=None,
            created_at=None,
            createdAt=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.id == "1"
        assert result.address is None
        assert result.name is None
        assert result.blockchain is None
        assert result.status is None

    def test_created_at_camelcase_fallback(self) -> None:
        created = datetime(2025, 1, 1, 0, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id=2,
            address="0xabc",
            name="Test",
            blockchain="ETH",
            network="goerli",
            abi=None,
            status="PENDING",
            created_at=None,
            createdAt=created,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        assert result.created_at == created

    def test_result_attributes_accessible(self) -> None:
        dto = SimpleNamespace(
            id=3,
            address="0xdef",
            name="Swap Contract",
            blockchain="BSC",
            network="mainnet",
            abi="[]",
            status="ACTIVE",
            created_at=None,
            createdAt=None,
        )
        result = ContractWhitelistingService._map_contract_from_dto(dto)
        # Verify all attributes are directly accessible
        assert result.id == "3"
        assert result.address == "0xdef"
        assert result.name == "Swap Contract"
        assert result.blockchain == "BSC"
