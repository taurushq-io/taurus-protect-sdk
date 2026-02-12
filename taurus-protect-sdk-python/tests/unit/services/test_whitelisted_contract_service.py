"""Unit tests for ContractWhitelistingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.contract_whitelisting_service import (
    ContractWhitelistingService,
)


class TestContractWhitelistingServiceGet:
    """Tests for ContractWhitelistingService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_invalid_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_id must be positive"):
            service.get(contract_id=0)

    def test_returns_contract(self) -> None:
        service, api = self._make_service()
        dto = MagicMock()
        dto.id = "1"
        dto.address = "0xabc"
        dto.name = "Uniswap Router"
        dto.blockchain = "ETH"
        dto.network = "mainnet"
        dto.abi = "{}"
        dto.status = "APPROVED"
        dto.created_at = None
        dto.createdAt = None
        reply = MagicMock()
        reply.result = dto
        api.whitelist_service_get_whitelisted_contract.return_value = reply

        contract = service.get(contract_id=1)

        assert contract.id == "1"
        assert contract.address == "0xabc"
        assert contract.name == "Uniswap Router"

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.whitelist_service_get_whitelisted_contract.return_value = reply

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(contract_id=1)


class TestContractWhitelistingServiceList:
    """Tests for ContractWhitelistingService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

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
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.whitelist_service_get_whitelisted_contracts.return_value = reply

        contracts, pagination = service.list()

        assert contracts == []

    def test_passes_blockchain_filter(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.whitelist_service_get_whitelisted_contracts.return_value = reply

        service.list(blockchain="ETH", network="mainnet")

        api.whitelist_service_get_whitelisted_contracts.assert_called_once_with(
            blockchain="ETH",
            network="mainnet",
            limit="50",
            offset="0",
        )


class TestContractWhitelistingServiceCreate:
    """Tests for ContractWhitelistingService.create()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_empty_address(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="address"):
            service.create(address="", name="Test", blockchain="ETH")

    def test_raises_on_empty_name(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="name"):
            service.create(address="0xabc", name="", blockchain="ETH")

    def test_raises_on_empty_blockchain(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="blockchain"):
            service.create(address="0xabc", name="Test", blockchain="")


class TestContractWhitelistingServiceDelete:
    """Tests for ContractWhitelistingService.delete()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_invalid_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_id must be positive"):
            service.delete(contract_id=0)

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        api.whitelist_service_delete_whitelisted_contract.return_value = None

        service.delete(contract_id=42)

        api.whitelist_service_delete_whitelisted_contract.assert_called_once_with(
            "42"
        )


class TestContractWhitelistingServiceApprove:
    """Tests for ContractWhitelistingService.approve_whitelisted_contracts()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_empty_ids(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_ids cannot be empty"):
            service.approve_whitelisted_contracts(
                contract_ids=[], signature="sig"
            )

    def test_raises_on_empty_signature(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="signature"):
            service.approve_whitelisted_contracts(
                contract_ids=["1"], signature=""
            )


class TestContractWhitelistingServiceCreateAttribute:
    """Tests for ContractWhitelistingService.create_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_empty_contract_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_id"):
            service.create_attribute(contract_id="", key="k", value="v")

    def test_raises_on_empty_key(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="key"):
            service.create_attribute(contract_id="1", key="", value="v")


class TestContractWhitelistingServiceGetAttribute:
    """Tests for ContractWhitelistingService.get_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        contract_api = MagicMock()
        service = ContractWhitelistingService(
            api_client=api_client, contract_whitelisting_api=contract_api
        )
        return service, contract_api

    def test_raises_on_empty_contract_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_id"):
            service.get_attribute(contract_id="", key="k")

    def test_raises_on_empty_key(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="key"):
            service.get_attribute(contract_id="1", key="")

    def test_returns_none_when_result_is_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.whitelist_service_get_whitelisted_contract_attribute.return_value = reply

        result = service.get_attribute(contract_id="1", key="mykey")

        assert result is None
