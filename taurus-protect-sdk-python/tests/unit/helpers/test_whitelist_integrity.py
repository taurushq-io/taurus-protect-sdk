"""Tests for whitelist integrity and hash computation utilities."""

from __future__ import annotations

import json

import pytest

from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.whitelist_hash_helper import (
    compute_legacy_hashes,
    parse_whitelisted_address_from_json,
)
from taurus_protect.helpers.whitelist_integrity_helper import (
    verify_envelope_field_match,
)
from taurus_protect.models.whitelisted_address import (
    InternalAddress,
    SignedWhitelistedAddressEnvelope,
    WhitelistMetadata,
    WhitelistedAddress,
)


# Test data for parse_whitelisted_address_from_json tests
VALID_WHITELISTED_ADDRESS_JSON = json.dumps(
    {
        "id": "wa-123",
        "address": "0xf631ce893edb440e49188a991250051d07968186",
        "label": "My ETH Address",
        "currency": "ETH",
        "network": "mainnet",
        "contractType": "ERC20",
        "linkedInternalAddresses": [
            {"address": "0x1111111111111111111111111111111111111111"},
            {"address": "0x2222222222222222222222222222222222222222"},
        ],
    },
    separators=(",", ":"),
)

# JSON with contractType field for legacy hash testing
JSON_WITH_CONTRACT_TYPE = (
    '{"address":"0x123","contractType":"ERC20","currency":"ETH","label":"test"}'
)

# JSON with labels in linkedInternalAddresses for legacy hash testing
JSON_WITH_LABELS = (
    '{"address":"0x123","currency":"ETH","label":"test",'
    '"linkedInternalAddresses":[{"address":"0xabc","label":"internal1"}]}'
)

# JSON with both contractType and labels
JSON_WITH_BOTH = (
    '{"address":"0x123","contractType":"ERC20","currency":"ETH","label":"test",'
    '"linkedInternalAddresses":[{"address":"0xabc","label":"internal1"}]}'
)

# Compute the hash of valid JSON for envelope tests
from taurus_protect.crypto.hashing import calculate_hex_hash

VALID_JSON_HASH = calculate_hex_hash(VALID_WHITELISTED_ADDRESS_JSON)


class TestParseWhitelistedAddressFromJson:
    """Tests for parse_whitelisted_address_from_json function."""

    def test_parse_whitelisted_address_valid_match(self) -> None:
        """Parse valid JSON and verify all fields are correctly mapped."""
        result = parse_whitelisted_address_from_json(VALID_WHITELISTED_ADDRESS_JSON)

        assert result.id == "wa-123"
        assert result.address == "0xf631ce893edb440e49188a991250051d07968186"
        assert result.label == "My ETH Address"
        assert result.currency == "ETH"
        assert result.network == "mainnet"
        assert result.contract_type == "ERC20"
        assert len(result.linked_internal_addresses) == 2
        assert (
            result.linked_internal_addresses[0].address
            == "0x1111111111111111111111111111111111111111"
        )
        assert (
            result.linked_internal_addresses[1].address
            == "0x2222222222222222222222222222222222222222"
        )

    def test_parse_whitelisted_address_empty_payload(self) -> None:
        """Empty string should raise WhitelistError."""
        with pytest.raises(WhitelistError) as exc_info:
            parse_whitelisted_address_from_json("")
        assert "null or empty" in str(exc_info.value).lower()

    def test_parse_whitelisted_address_invalid_json(self) -> None:
        """Invalid JSON should raise WhitelistError."""
        with pytest.raises(WhitelistError) as exc_info:
            parse_whitelisted_address_from_json("{invalid json}")
        assert "parse" in str(exc_info.value).lower() or "json" in str(
            exc_info.value
        ).lower()

    def test_parse_whitelisted_address_missing_optional_fields(self) -> None:
        """JSON with only required fields should parse successfully."""
        minimal_json = json.dumps(
            {
                "id": "wa-456",
                "address": "0x1234567890abcdef1234567890abcdef12345678",
                "currency": "BTC",
                "label": "Test Wallet",
            },
            separators=(",", ":"),
        )
        result = parse_whitelisted_address_from_json(minimal_json)

        assert result.id == "wa-456"
        assert result.address == "0x1234567890abcdef1234567890abcdef12345678"
        assert result.currency == "BTC"
        assert result.label == "Test Wallet"
        assert result.network is None
        assert result.contract_type is None
        assert result.linked_internal_addresses == []


class TestComputeLegacyHashes:
    """Tests for compute_legacy_hashes function."""

    def test_compute_legacy_hashes_contract_type_only(self) -> None:
        """Verify hash computed without contractType field."""
        hashes = compute_legacy_hashes(JSON_WITH_CONTRACT_TYPE)

        # Should have at least one hash (without contractType)
        assert len(hashes) >= 1

        # The computed hash should be for the payload without contractType
        expected_without_contract = '{"address":"0x123","currency":"ETH","label":"test"}'
        expected_hash = calculate_hex_hash(expected_without_contract)
        assert expected_hash in hashes

    def test_compute_legacy_hashes_labels_only(self) -> None:
        """Verify hash computed without labels in linkedInternalAddresses."""
        hashes = compute_legacy_hashes(JSON_WITH_LABELS)

        # Should have at least one hash (without labels)
        assert len(hashes) >= 1

        # The computed hash should be for the payload without labels
        expected_without_labels = (
            '{"address":"0x123","currency":"ETH","label":"test",'
            '"linkedInternalAddresses":[{"address":"0xabc"}]}'
        )
        expected_hash = calculate_hex_hash(expected_without_labels)
        assert expected_hash in hashes

    def test_compute_legacy_hashes_both_removed(self) -> None:
        """Verify hash computed without both contractType and labels."""
        hashes = compute_legacy_hashes(JSON_WITH_BOTH)

        # Should have multiple hashes
        # Strategy 1: without contractType
        # Strategy 2: without labels
        # Strategy 3: without both
        assert len(hashes) >= 1

        # The hash without both should be present
        expected_without_both = (
            '{"address":"0x123","currency":"ETH","label":"test",'
            '"linkedInternalAddresses":[{"address":"0xabc"}]}'
        )
        expected_hash = calculate_hex_hash(expected_without_both)
        assert expected_hash in hashes

    def test_compute_legacy_hashes_unique(self) -> None:
        """Verify no duplicate hashes in results."""
        hashes = compute_legacy_hashes(JSON_WITH_BOTH)

        # All hashes should be unique
        assert len(hashes) == len(set(hashes))

    def test_compute_legacy_hashes_empty_payload(self) -> None:
        """Empty payload should return empty list."""
        hashes = compute_legacy_hashes("")
        assert hashes == []

    def test_compute_legacy_hashes_no_removable_fields(self) -> None:
        """Payload without contractType or removable labels returns empty list.

        Note: The regex patterns only match:
        - contractType field: ,"contractType":"..."
        - label before closing brace: ,"label":"..."}

        A payload where label is NOT before a closing brace won't be modified.
        """
        # Use JSON where label is followed by another field, not a closing brace
        simple_json = '{"address":"0x123","currency":"ETH","label":"test","memo":"note"}'
        hashes = compute_legacy_hashes(simple_json)
        # No fields to remove (contractType absent, label not before closing brace)
        assert hashes == []


class TestVerifyEnvelopeFieldMatch:
    """Tests for verify_envelope_field_match function."""

    @pytest.fixture
    def valid_envelope(self) -> SignedWhitelistedAddressEnvelope:
        """Create a valid envelope for testing."""
        metadata = WhitelistMetadata(
            hash=VALID_JSON_HASH,
            payload_as_string=VALID_WHITELISTED_ADDRESS_JSON,
        )
        return SignedWhitelistedAddressEnvelope(metadata=metadata)

    @pytest.fixture
    def valid_db_address(self) -> WhitelistedAddress:
        """Create a valid database address that matches the envelope."""
        return WhitelistedAddress(
            id="wa-123",
            address="0xf631ce893edb440e49188a991250051d07968186",
            label="My ETH Address",
            currency="ETH",
            network="mainnet",
            contract_type="ERC20",
            linked_internal_addresses=[
                InternalAddress(address="0x1111111111111111111111111111111111111111"),
                InternalAddress(address="0x2222222222222222222222222222222222222222"),
            ],
        )

    def test_verify_envelope_field_match_valid(
        self, valid_envelope: SignedWhitelistedAddressEnvelope, valid_db_address: WhitelistedAddress
    ) -> None:
        """Matching fields should pass verification without error."""
        # Should not raise
        verify_envelope_field_match(valid_db_address, valid_envelope)

    def test_verify_envelope_field_match_address_mismatch(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """Different address should raise IntegrityError."""
        db_address = WhitelistedAddress(
            id="wa-123",
            address="0x0000000000000000000000000000000000000000",  # Different
            label="My ETH Address",
            currency="ETH",
            network="mainnet",
            contract_type="ERC20",
        )

        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(db_address, valid_envelope)
        assert "Address" in str(exc_info.value)

    def test_verify_envelope_field_match_label_mismatch(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """Different label should raise IntegrityError."""
        db_address = WhitelistedAddress(
            id="wa-123",
            address="0xf631ce893edb440e49188a991250051d07968186",
            label="Different Label",  # Different
            currency="ETH",
            network="mainnet",
            contract_type="ERC20",
        )

        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(db_address, valid_envelope)
        assert "Label" in str(exc_info.value)

    def test_verify_envelope_field_match_currency_mismatch(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """Different currency should raise IntegrityError."""
        db_address = WhitelistedAddress(
            id="wa-123",
            address="0xf631ce893edb440e49188a991250051d07968186",
            label="My ETH Address",
            currency="BTC",  # Different
            network="mainnet",
            contract_type="ERC20",
        )

        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(db_address, valid_envelope)
        assert "Currency" in str(exc_info.value)

    def test_verify_envelope_field_match_network_mismatch(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """Different network should raise IntegrityError."""
        db_address = WhitelistedAddress(
            id="wa-123",
            address="0xf631ce893edb440e49188a991250051d07968186",
            label="My ETH Address",
            currency="ETH",
            network="testnet",  # Different
            contract_type="ERC20",
        )

        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(db_address, valid_envelope)
        assert "Network" in str(exc_info.value)

    def test_verify_envelope_field_match_contract_type_mismatch(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """Different contract_type should raise IntegrityError."""
        db_address = WhitelistedAddress(
            id="wa-123",
            address="0xf631ce893edb440e49188a991250051d07968186",
            label="My ETH Address",
            currency="ETH",
            network="mainnet",
            contract_type="ERC721",  # Different
        )

        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(db_address, valid_envelope)
        assert "ContractType" in str(exc_info.value)

    def test_verify_envelope_field_match_none_address(
        self, valid_envelope: SignedWhitelistedAddressEnvelope
    ) -> None:
        """None db_address should raise IntegrityError."""
        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(None, valid_envelope)  # type: ignore[arg-type]
        assert "None" in str(exc_info.value) or "cannot be" in str(exc_info.value)

    def test_verify_envelope_field_match_none_envelope(
        self, valid_db_address: WhitelistedAddress
    ) -> None:
        """None envelope should raise IntegrityError."""
        with pytest.raises(IntegrityError) as exc_info:
            verify_envelope_field_match(valid_db_address, None)  # type: ignore[arg-type]
        assert "None" in str(exc_info.value) or "cannot be" in str(exc_info.value)

    def test_verify_envelope_field_match_empty_optional_fields(self) -> None:
        """Empty optional fields should match None/empty in envelope."""
        # Create JSON without optional fields
        simple_json = json.dumps(
            {
                "id": "wa-simple",
                "address": "0x1234567890abcdef1234567890abcdef12345678",
                "currency": "ETH",
                "label": "Simple Address",
            },
            separators=(",", ":"),
        )
        simple_hash = calculate_hex_hash(simple_json)

        metadata = WhitelistMetadata(
            hash=simple_hash,
            payload_as_string=simple_json,
        )
        envelope = SignedWhitelistedAddressEnvelope(metadata=metadata)

        db_address = WhitelistedAddress(
            id="wa-simple",
            address="0x1234567890abcdef1234567890abcdef12345678",
            label="Simple Address",
            currency="ETH",
            network=None,  # Matches envelope's missing network
            contract_type=None,  # Matches envelope's missing contract_type
        )

        # Should not raise
        verify_envelope_field_match(db_address, envelope)

    def test_verify_envelope_field_match_empty_string_vs_none(self) -> None:
        """Empty string should be treated as equivalent to None."""
        simple_json = json.dumps(
            {
                "id": "wa-test",
                "address": "0xabcdef1234567890abcdef1234567890abcdef12",
                "currency": "ETH",
                "label": "Test",
                "network": "",  # Empty string in JSON
            },
            separators=(",", ":"),
        )
        simple_hash = calculate_hex_hash(simple_json)

        metadata = WhitelistMetadata(
            hash=simple_hash,
            payload_as_string=simple_json,
        )
        envelope = SignedWhitelistedAddressEnvelope(metadata=metadata)

        db_address = WhitelistedAddress(
            id="wa-test",
            address="0xabcdef1234567890abcdef1234567890abcdef12",
            label="Test",
            currency="ETH",
            network=None,  # None should match empty string
        )

        # Should not raise - empty string and None are equivalent
        verify_envelope_field_match(db_address, envelope)
