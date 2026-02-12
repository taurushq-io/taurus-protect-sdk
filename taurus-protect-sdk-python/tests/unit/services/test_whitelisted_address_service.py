"""Security tests for WhitelistedAddressService.

These tests verify that security-critical fields are sourced ONLY from the
cryptographically verified payload, not from unverified DTO attributes.
"""

from __future__ import annotations

import json
from typing import Any, List
from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.whitelisted_address_verifier import AddressVerificationResult
from taurus_protect.helpers.whitelist_hash_helper import parse_whitelisted_address_from_json
from taurus_protect.models.governance_rules import DecodedRulesContainer
from taurus_protect.models.whitelisted_address import WhitelistedAddress
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService


class MockDTO:
    """Mock DTO object for testing."""

    def __init__(self, **kwargs: Any) -> None:
        for key, value in kwargs.items():
            setattr(self, key, value)


def create_valid_payload(
    id_val: str = "wa-123",
    address: str = "0xf631ce893edb440e49188a991250051d07968186",
    label: str = "Test Address",
    currency: str = "ETH",
    network: str = "mainnet",
    contract_type: str = "ERC20",
) -> str:
    """Create a valid JSON payload string."""
    return json.dumps(
        {
            "id": id_val,
            "address": address,
            "label": label,
            "currency": currency,
            "network": network,
            "contractType": contract_type,
        },
        separators=(",", ":"),
    )


def create_mock_dto_with_metadata(
    payload_str: str,
    dto_overrides: dict[str, Any] | None = None,
) -> MockDTO:
    """Create a mock DTO with metadata containing the given payload."""
    payload_hash = calculate_hex_hash(payload_str)

    metadata = MockDTO(
        hash=payload_hash,
        payload_as_string=payload_str,
        payloadAsString=payload_str,
    )

    dto_attrs = {
        "id": "wa-123",
        "metadata": metadata,
        "signed_address": None,
        "signedAddress": None,
        "rules_container": None,
        "rulesContainer": None,
        "rules_signatures": None,
        "rulesSignatures": None,
    }

    if dto_overrides:
        dto_attrs.update(dto_overrides)

    return MockDTO(**dto_attrs)


def _mock_verify_side_effect(envelope, **kwargs):
    """Mock verifier that parses the WhitelistedAddress from the envelope payload."""
    payload_str = envelope.metadata.payload_as_string if envelope.metadata else ""
    verified_addr = parse_whitelisted_address_from_json(payload_str)
    return AddressVerificationResult(
        rules_container=DecodedRulesContainer(),
        verified_hash="mocked_hash",
        verified_whitelisted_address=verified_addr,
    )


def _make_service_with_mocked_verifier() -> WhitelistedAddressService:
    """Create a WhitelistedAddressService with a mocked verifier.

    The verifier is mocked to return a successful result, so tests
    can focus on field extraction and DTO-to-model mapping behavior.
    """
    api_client = MagicMock()
    whitelisting_api = MagicMock()
    super_admin_keys: List[Any] = []
    min_valid_signatures = 0

    service = WhitelistedAddressService(
        api_client=api_client,
        whitelisting_api=whitelisting_api,
        super_admin_keys=super_admin_keys,
        min_valid_signatures=min_valid_signatures,
    )

    # Mock the verifier to always succeed, parsing the address from the payload
    mock_verifier = MagicMock()
    mock_verifier.verify_whitelisted_address.side_effect = _mock_verify_side_effect
    service._verifier = mock_verifier

    return service


class TestWhitelistedAddressServiceSecurity:
    """Security tests for WhitelistedAddressService."""

    @pytest.fixture
    def service(self) -> WhitelistedAddressService:
        """Create a WhitelistedAddressService with mocked verifier for field tests."""
        return _make_service_with_mocked_verifier()

    def test_list_verifies_all_addresses(self, service: WhitelistedAddressService) -> None:
        """Test that list() verifies each address (strict mode)."""
        # Create valid payload and DTO
        payload = create_valid_payload()
        dto = create_mock_dto_with_metadata(payload)

        # Mock API response with multiple items
        mock_reply = MagicMock()
        mock_reply.result = [dto, dto]  # Two addresses
        mock_reply.total_items = "2"

        service._api.whitelist_service_get_whitelisted_addresses.return_value = mock_reply

        # Should succeed and return verified addresses
        addresses, pagination = service.list(limit=50, offset=0)

        assert len(addresses) == 2
        assert all(addr.address == "0xf631ce893edb440e49188a991250051d07968186" for addr in addresses)

    def test_list_raises_integrity_error_on_invalid_envelope(self) -> None:
        """Test that list() raises IntegrityError when envelope verification fails."""
        # Use a service with a verifier that raises IntegrityError
        service = _make_service_with_mocked_verifier()
        service._verifier.verify_whitelisted_address.side_effect = IntegrityError(
            "metadata hash verification failed"
        )

        payload = create_valid_payload()
        dto = create_mock_dto_with_metadata(payload)

        mock_reply = MagicMock()
        mock_reply.result = [dto]
        mock_reply.total_items = "1"

        service._api.whitelist_service_get_whitelisted_addresses.return_value = mock_reply

        # Should raise IntegrityError (strict mode)
        with pytest.raises(IntegrityError):
            service.list(limit=50, offset=0)

    def test_address_fields_from_verified_payload_only(
        self, service: WhitelistedAddressService
    ) -> None:
        """Test that address fields come from verified payload, not DTO."""
        # Payload has these values (cryptographically verified)
        payload = create_valid_payload(
            address="0xverified_address_from_payload",
            label="Verified Label",
            currency="ETH",
            network="mainnet",
        )

        # DTO has DIFFERENT values (unverified - should be ignored)
        dto = create_mock_dto_with_metadata(
            payload,
            dto_overrides={
                "address": "0xmalicious_dto_address",
                "label": "Malicious DTO Label",
                "currency": "FAKE",
                "network": "attacker_network",
            },
        )

        mock_reply = MagicMock()
        mock_reply.result = [dto]
        mock_reply.total_items = "1"

        service._api.whitelist_service_get_whitelisted_addresses.return_value = mock_reply

        addresses, _ = service.list(limit=50, offset=0)

        # Address fields should come from payload, not DTO
        assert len(addresses) == 1
        addr = addresses[0]
        assert addr.address == "0xverified_address_from_payload"
        assert addr.label == "Verified Label"
        assert addr.currency == "ETH"
        assert addr.network == "mainnet"

    def test_get_envelope_extracts_address_from_payload(
        self, service: WhitelistedAddressService
    ) -> None:
        """Test that get_envelope() extracts address from verified payload."""
        payload = create_valid_payload(
            address="0xcryptographically_verified",
            label="Verified",
        )
        dto = create_mock_dto_with_metadata(payload)

        mock_reply = MagicMock()
        mock_reply.result = dto

        service._api.whitelist_service_get_whitelisted_address.return_value = mock_reply

        envelope = service.get_envelope(123)

        assert envelope.verified_whitelisted_address is not None
        assert envelope.verified_whitelisted_address.address == "0xcryptographically_verified"
        assert envelope.verified_whitelisted_address.label == "Verified"

    def test_missing_payload_as_string_raises_integrity_error(self) -> None:
        """Test that missing payload_as_string raises IntegrityError."""
        # Use a service with a verifier that raises IntegrityError for empty payload
        service = _make_service_with_mocked_verifier()
        service._verifier.verify_whitelisted_address.side_effect = IntegrityError(
            "payloadAsString is empty"
        )

        dto = create_mock_dto_with_metadata("")
        dto.metadata.payload_as_string = None
        dto.metadata.payloadAsString = None

        mock_reply = MagicMock()
        mock_reply.result = dto

        service._api.whitelist_service_get_whitelisted_address.return_value = mock_reply

        with pytest.raises(IntegrityError) as exc_info:
            service.get_envelope(123)
        assert "empty" in str(exc_info.value).lower()


class TestMapEnvelopeFromDto:
    """Tests for _map_envelope_from_dto static method."""

    def test_map_envelope_extracts_metadata_hash(self) -> None:
        """Test that envelope mapper extracts hash from metadata."""
        payload = create_valid_payload()
        expected_hash = calculate_hex_hash(payload)
        dto = create_mock_dto_with_metadata(payload)

        envelope = WhitelistedAddressService._map_envelope_from_dto(dto)

        assert envelope.metadata is not None
        assert envelope.metadata.hash == expected_hash
        assert envelope.metadata.payload_as_string == payload

    def test_map_envelope_handles_missing_metadata(self) -> None:
        """Test that envelope mapper handles missing metadata gracefully."""
        dto = MockDTO(
            id="wa-123",
            metadata=None,
            signed_address=None,
            rules_container=None,
            rules_signatures=None,
        )

        envelope = WhitelistedAddressService._map_envelope_from_dto(dto)

        assert envelope.metadata is None

    def test_map_envelope_extracts_signatures(self) -> None:
        """Test that envelope mapper extracts signatures."""
        payload = create_valid_payload()
        # Create nested structure matching TgvalidatordWhitelistSignature:
        # - signature: TgvalidatordWhitelistUserSignature (user_id, signature, comment)
        # - hashes: List[str]
        user_sig = MockDTO(
            user_id="user-1",
            signature="sig123",
            comment="test comment",
        )
        sig_dto = MockDTO(
            signature=user_sig,
            hashes=["hash123"],
        )
        signed_address = MockDTO(signatures=[sig_dto])

        dto = create_mock_dto_with_metadata(payload)
        dto.signed_address = signed_address
        dto.signedAddress = signed_address

        envelope = WhitelistedAddressService._map_envelope_from_dto(dto)

        assert len(envelope.signatures) == 1
        assert envelope.signatures[0].user_id == "user-1"
        assert envelope.signatures[0].signature == "sig123"
