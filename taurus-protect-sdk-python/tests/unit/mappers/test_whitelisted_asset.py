"""Unit tests for whitelisted asset mapper functions."""

import json
from types import SimpleNamespace

from taurus_protect.models.whitelisted_address import WhitelistedAsset
from taurus_protect.services.whitelisted_asset_service import WhitelistedAssetService


class TestMapAssetFromDto:
    """Tests for WhitelistedAssetService._map_asset_from_dto."""

    def test_maps_fields_from_payload(self) -> None:
        payload = {
            "name": "Tether",
            "symbol": "USDT",
            "blockchain": "ETH",
            "network": "mainnet",
            "contractAddress": "0xdac17f958d2ee523a2206206994597c13d831ec7",
        }
        dto = SimpleNamespace(
            id="asset-1",
            tenant_id="t-1",
            tenantId=None,
            status="APPROVED",
            action="APPROVE",
            rule="global",
            created_at="2024-01-01",
            createdAt=None,
            metadata=SimpleNamespace(
                hash="sha256hash",
                payload_as_string=json.dumps(payload),
                payloadAsString=None,
            ),
            rules_container="base64rc",
            rulesContainer=None,
            rules_signatures="base64rs",
            rulesSignatures=None,
            signed_contract_address=None,
            signedContractAddress=None,
            business_rule_enabled=True,
            businessRuleEnabled=None,
        )
        result = WhitelistedAssetService._map_asset_from_dto(dto)
        assert isinstance(result, WhitelistedAsset)
        assert result.id == "asset-1"
        assert result.tenant_id == "t-1"
        # Security-critical fields from payload only
        assert result.name == "Tether"
        assert result.symbol == "USDT"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.contract_address == "0xdac17f958d2ee523a2206206994597c13d831ec7"
        # Non-security fields from DTO
        assert result.status == "APPROVED"
        assert result.action == "APPROVE"
        assert result.business_rule_enabled is True

    def test_security_fields_none_when_no_payload(self) -> None:
        """Security-critical fields must be None when payload is missing."""
        dto = SimpleNamespace(
            id="asset-2",
            tenant_id=None,
            tenantId=None,
            status="PENDING",
            action=None,
            rule=None,
            created_at=None,
            createdAt=None,
            metadata=None,
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_contract_address=None,
            signedContractAddress=None,
            business_rule_enabled=False,
            businessRuleEnabled=False,
        )
        result = WhitelistedAssetService._map_asset_from_dto(dto)
        assert result.name is None
        assert result.symbol is None
        assert result.blockchain is None
        assert result.network is None
        assert result.contract_address is None

    def test_maps_signed_contract_address(self) -> None:
        user_sig_dto = SimpleNamespace(
            user_id="u-1",
            userId=None,
            signature="sig-data",
            comment="ok",
        )
        sig_dto = SimpleNamespace(
            user_signature=user_sig_dto,
            userSignature=None,
            signature=None,
            hashes=["hash-1"],
        )
        signed_dto = SimpleNamespace(
            payload="signed-payload-base64",
            signatures=[sig_dto],
        )
        payload = {"name": "Token", "symbol": "TKN", "blockchain": "ETH"}
        dto = SimpleNamespace(
            id="asset-3",
            tenant_id=None,
            tenantId=None,
            status=None,
            action=None,
            rule=None,
            created_at=None,
            createdAt=None,
            metadata=SimpleNamespace(
                hash="h",
                payload_as_string=json.dumps(payload),
                payloadAsString=None,
            ),
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_contract_address=signed_dto,
            signedContractAddress=None,
            business_rule_enabled=False,
            businessRuleEnabled=False,
        )
        result = WhitelistedAssetService._map_asset_from_dto(dto)
        assert result.signed_contract_address is not None
        assert result.signed_contract_address.payload == "signed-payload-base64"
        assert len(result.signed_contract_address.signatures) == 1
        entry = result.signed_contract_address.signatures[0]
        assert entry.user_signature is not None
        assert entry.user_signature.user_id == "u-1"
        assert entry.hashes == ["hash-1"]

    def test_camelcase_fallbacks(self) -> None:
        payload = {"name": "Test", "symbol": "TST"}
        dto = SimpleNamespace(
            id="asset-4",
            tenant_id=None,
            tenantId="camel-tenant",
            status=None,
            action=None,
            rule=None,
            created_at=None,
            createdAt="2024-06-01",
            metadata=SimpleNamespace(
                hash="h",
                payload_as_string=None,
                payloadAsString=json.dumps(payload),
            ),
            rules_container=None,
            rulesContainer="camelRC",
            rules_signatures=None,
            rulesSignatures="camelRS",
            signed_contract_address=None,
            signedContractAddress=None,
            business_rule_enabled=False,
            businessRuleEnabled=True,
        )
        result = WhitelistedAssetService._map_asset_from_dto(dto)
        assert result.tenant_id == "camel-tenant"
        # Pydantic coerces "2024-06-01" string to datetime
        from datetime import datetime
        assert result.created_at == datetime(2024, 6, 1, 0, 0)
        assert result.rules_container == "camelRC"
        assert result.rules_signatures == "camelRS"
        assert result.business_rule_enabled is True

    def test_invalid_payload_json_yields_none_fields(self) -> None:
        dto = SimpleNamespace(
            id="asset-5",
            tenant_id=None,
            tenantId=None,
            status=None,
            action=None,
            rule=None,
            created_at=None,
            createdAt=None,
            metadata=SimpleNamespace(
                hash="h",
                payload_as_string="not-valid-json",
                payloadAsString=None,
            ),
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_contract_address=None,
            signedContractAddress=None,
            business_rule_enabled=False,
            businessRuleEnabled=False,
        )
        result = WhitelistedAssetService._map_asset_from_dto(dto)
        # Invalid JSON means payload dict is empty, security fields are None
        assert result.name is None
        assert result.symbol is None
        assert result.blockchain is None
