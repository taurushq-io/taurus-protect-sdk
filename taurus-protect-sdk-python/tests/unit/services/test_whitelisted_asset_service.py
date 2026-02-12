"""Security tests for WhitelistedAssetService.

These tests verify that security-critical fields are sourced ONLY from the
cryptographically verified payload, not from unverified DTO attributes.
"""

from __future__ import annotations

import json
from typing import Any, Dict, List, Optional
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.whitelisted_asset_service import (
    WhitelistedAssetService,
)


class MockDTO:
    """Mock DTO object for testing."""

    def __init__(self, **kwargs: Any) -> None:
        for key, value in kwargs.items():
            setattr(self, key, value)


def create_mock_dto_with_payload(
    payload: Dict[str, Any],
    dto_overrides: Optional[Dict[str, Any]] = None,
) -> MockDTO:
    """Create a mock DTO with metadata containing the given payload."""
    metadata = MockDTO(
        hash="abc123",
        payload=payload,
        payload_as_string=json.dumps(payload, separators=(",", ":")),
        payloadAsString=json.dumps(payload, separators=(",", ":")),
    )

    dto_attrs: Dict[str, Any] = {
        "id": "asset-123",
        "tenant_id": "tenant-1",
        "tenantId": "tenant-1",
        "metadata": metadata,
        "signed_contract_address": None,
        "signedContractAddress": None,
        "rules_container": None,
        "rulesContainer": None,
        "rules_signatures": None,
        "rulesSignatures": None,
        "status": "APPROVED",
        "action": None,
        "rule": None,
        "created_at": None,
        "createdAt": None,
        "business_rule_enabled": False,
        "businessRuleEnabled": False,
        # Default DTO values (should be ignored for security fields)
        "name": None,
        "symbol": None,
        "blockchain": None,
        "network": None,
        "contract_address": None,
        "contractAddress": None,
    }

    if dto_overrides:
        dto_attrs.update(dto_overrides)

    return MockDTO(**dto_attrs)


def _create_service_with_mock_verifier() -> WhitelistedAssetService:
    """Create a WhitelistedAssetService with a mocked verifier for unit tests.

    The verifier is always required, but for field sourcing tests we mock it
    to avoid needing real SuperAdmin keys.
    """
    api_client = MagicMock()
    assets_api = MagicMock()
    mock_keys = [MagicMock()]  # Mock key list (non-empty)
    service = WhitelistedAssetService(
        api_client=api_client,
        assets_api=assets_api,
        super_admin_keys=mock_keys,
        min_valid_signatures=1,
    )
    # Replace the real verifier with a mock to skip actual verification
    service._verifier = MagicMock()
    return service


class TestWhitelistedAssetServiceSecurity:
    """Security tests for field sourcing in WhitelistedAssetService."""

    @pytest.fixture
    def service(self) -> WhitelistedAssetService:
        """Create a WhitelistedAssetService for testing with mocked verifier."""
        return _create_service_with_mock_verifier()

    def test_security_fields_from_payload_not_dto(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that name, symbol, contract_address come from payload, not DTO."""
        # Payload has verified values
        payload = {
            "name": "Verified Token Name",
            "symbol": "VTN",
            "contract_address": "0xverified_contract",
            "blockchain": "ETH",
            "network": "mainnet",
        }

        # DTO has DIFFERENT values (attacker-controlled, should be ignored)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "name": "Malicious Name",
                "symbol": "FAKE",
                "contract_address": "0xmalicious_contract",
                "contractAddress": "0xmalicious_contract",
                "blockchain": "ATTACKER_CHAIN",
                "network": "attacker_net",
            },
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # Security-critical fields MUST come from payload
        assert asset.name == "Verified Token Name"
        assert asset.symbol == "VTN"
        assert asset.contract_address == "0xverified_contract"
        assert asset.blockchain == "ETH"
        assert asset.network == "mainnet"

    def test_payload_missing_field_returns_none_not_dto_value(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that missing payload field returns None, not DTO value."""
        # Payload is missing 'symbol' (simulating incomplete payload)
        payload = {
            "name": "Token Name",
            # symbol is intentionally missing
            "contract_address": "0xcontract",
        }

        # DTO has symbol value (should NOT be used as fallback)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "symbol": "SHOULD_NOT_USE",  # This MUST be ignored
            },
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # Symbol should be None, not the DTO value
        assert asset.symbol is None
        # But name and contract_address should be from payload
        assert asset.name == "Token Name"
        assert asset.contract_address == "0xcontract"

    def test_payload_empty_returns_none_not_dto_value(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that empty payload results in None fields, not DTO values."""
        # Empty payload
        payload: Dict[str, Any] = {}

        # DTO has all values (should NOT be used)
        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "name": "DTO Name",
                "symbol": "DTO",
                "contract_address": "0xdto_contract",
                "contractAddress": "0xdto_contract",
                "blockchain": "DTO_CHAIN",
                "network": "dto_net",
            },
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # All security fields should be None (not from DTO)
        assert asset.name is None
        assert asset.symbol is None
        assert asset.contract_address is None
        assert asset.blockchain is None
        assert asset.network is None

    def test_contract_address_snake_case_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test contract_address is extracted with snake_case key."""
        payload = {
            "contract_address": "0xsnake_case_address",
        }

        dto = create_mock_dto_with_payload(payload)
        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        assert asset.contract_address == "0xsnake_case_address"

    def test_contract_address_camel_case_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test contract_address is extracted with camelCase key."""
        payload = {
            "contractAddress": "0xcamel_case_address",
        }

        dto = create_mock_dto_with_payload(payload)
        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        assert asset.contract_address == "0xcamel_case_address"

    def test_blockchain_case_variations_in_payload(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test blockchain is extracted with different case variations."""
        # Test lowercase
        payload1 = {"blockchain": "eth"}
        dto1 = create_mock_dto_with_payload(payload1)
        asset1 = WhitelistedAssetService._map_asset_from_dto(dto1)
        assert asset1.blockchain == "eth"

        # Test capitalized (Blockchain with capital B)
        payload2 = {"Blockchain": "ETH"}
        dto2 = create_mock_dto_with_payload(payload2)
        asset2 = WhitelistedAssetService._map_asset_from_dto(dto2)
        assert asset2.blockchain == "ETH"

    def test_non_security_fields_can_come_from_dto(
        self, service: WhitelistedAssetService
    ) -> None:
        """Test that non-security fields (status, action, etc.) come from DTO."""
        payload = {"name": "Token"}

        dto = create_mock_dto_with_payload(
            payload,
            dto_overrides={
                "status": "APPROVED",
                "action": "ADD",
                "rule": "rule-1",
            },
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # Non-security fields are fine from DTO
        assert asset.status == "APPROVED"
        assert asset.action == "ADD"
        assert asset.rule == "rule-1"

    def test_list_always_verifies_each_asset(self) -> None:
        """Test that list() calls _verify_asset for each asset."""
        from unittest.mock import patch

        service = _create_service_with_mock_verifier()

        payload = {
            "name": "Test Token",
            "symbol": "TT",
            "contract_address": "0xtest",
        }
        dto = create_mock_dto_with_payload(payload)

        mock_reply = MagicMock()
        mock_reply.result = [dto, dto]
        mock_reply.total_items = "2"

        service._api.whitelist_service_get_whitelisted_contracts.return_value = mock_reply

        # Mock _verify_asset to avoid precondition checks on missing
        # rules_container/signed_contract_address in the test DTO
        with patch.object(service, "_verify_asset") as mock_verify:
            assets, pagination = service.list(limit=50, offset=0)

        # Should return both assets and verify each one
        assert len(assets) == 2
        assert all(a.name == "Test Token" for a in assets)
        # _verify_asset should have been called for each asset
        assert mock_verify.call_count == 2


class TestMapAssetFromDto:
    """Direct tests for _map_asset_from_dto static method."""

    def test_no_metadata_returns_none_security_fields(self) -> None:
        """Test that missing metadata results in None for security fields."""
        dto = MockDTO(
            id="asset-1",
            tenant_id="tenant-1",
            metadata=None,
            signed_contract_address=None,
            rules_container=None,
            rules_signatures=None,
            status="PENDING",
            action=None,
            rule=None,
            created_at=None,
            business_rule_enabled=False,
            # DTO has values but should NOT be used
            name="DTO Name",
            symbol="DTO",
            blockchain="DTO_CHAIN",
            network="dto_net",
            contract_address="0xdto",
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # Without metadata payload, security fields must be None
        assert asset.name is None
        assert asset.symbol is None
        assert asset.contract_address is None
        assert asset.blockchain is None
        assert asset.network is None

    def test_payload_as_string_is_used_not_payload_dict(self) -> None:
        """Test that payload_as_string is used for extraction, not the raw payload dict.

        SECURITY: This test verifies that even if metadata.payload is None or different,
        we extract from payload_as_string which is the cryptographically verified source.
        """
        metadata = MockDTO(
            hash="abc",
            payload=None,  # DTO payload is None, but payload_as_string has data
            payload_as_string='{"name":"test"}',
        )

        dto = MockDTO(
            id="asset-1",
            metadata=metadata,
            signed_contract_address=None,
            rules_container=None,
            rules_signatures=None,
            status="PENDING",
            action=None,
            rule=None,
            created_at=None,
            tenant_id=None,
            business_rule_enabled=False,
        )

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        # SECURITY: Name should come from payload_as_string (verified source),
        # not from the None payload dict (unverified)
        assert asset.name == "test"

    def test_signed_contract_address_mapping(self) -> None:
        """Test that signed_contract_address is correctly mapped."""
        payload = {"name": "Token"}

        user_sig = MockDTO(
            user_id="user-1",
            userId="user-1",
            signature="sig123",
            comment="LGTM",
        )
        sig_entry = MockDTO(
            user_signature=user_sig,
            userSignature=user_sig,
            hashes=["hash1", "hash2"],
        )
        signed = MockDTO(
            payload="signed_payload",
            signatures=[sig_entry],
        )

        dto = create_mock_dto_with_payload(payload)
        dto.signed_contract_address = signed
        dto.signedContractAddress = signed

        asset = WhitelistedAssetService._map_asset_from_dto(dto)

        assert asset.signed_contract_address is not None
        assert len(asset.signed_contract_address.signatures) == 1
        assert asset.signed_contract_address.signatures[0].user_signature.user_id == "user-1"
        assert asset.signed_contract_address.signatures[0].hashes == ["hash1", "hash2"]
