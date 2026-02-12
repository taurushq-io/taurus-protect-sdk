"""Unit tests for whitelisted address mapper functions."""

import json
from types import SimpleNamespace

from taurus_protect.models.whitelisted_address import (
    SignedWhitelistedAddressEnvelope,
    WhitelistMetadata,
)
from taurus_protect.services.whitelisted_address_service import WhitelistedAddressService


class TestMapEnvelopeFromDto:
    """Tests for WhitelistedAddressService._map_envelope_from_dto."""

    def test_maps_basic_envelope_fields(self) -> None:
        dto = SimpleNamespace(
            metadata=SimpleNamespace(
                hash="abc123",
                payload_as_string='{"address":"0x1234","label":"Test"}',
                payloadAsString=None,
            ),
            blockchain="ETH",
            network="mainnet",
            rules_container="base64rc==",
            rulesContainer=None,
            rules_signatures="base64rs==",
            rulesSignatures=None,
            signed_address=None,
            signedAddress=None,
            attributes=None,
        )
        result = WhitelistedAddressService._map_envelope_from_dto(dto)
        assert isinstance(result, SignedWhitelistedAddressEnvelope)
        assert result.metadata is not None
        assert result.metadata.hash == "abc123"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.rules_container == "base64rc=="
        assert result.rules_signatures == "base64rs=="

    def test_maps_envelope_without_metadata(self) -> None:
        dto = SimpleNamespace(
            metadata=None,
            blockchain="BTC",
            network=None,
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_address=None,
            signedAddress=None,
            attributes=None,
        )
        result = WhitelistedAddressService._map_envelope_from_dto(dto)
        assert result.metadata is None
        assert result.blockchain == "BTC"

    def test_maps_signed_address_with_signatures(self) -> None:
        user_sig_dto = SimpleNamespace(
            user_id="u-1",
            signature="sig-base64",
            comment="approved",
        )
        sig_dto = SimpleNamespace(
            signature=user_sig_dto,
            hashes=["hash1", "hash2"],
        )
        signed_addr_dto = SimpleNamespace(
            signatures=[sig_dto],
        )
        dto = SimpleNamespace(
            metadata=SimpleNamespace(
                hash="h1",
                payload_as_string="{}",
                payloadAsString=None,
            ),
            blockchain="ETH",
            network="mainnet",
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_address=signed_addr_dto,
            signedAddress=None,
            attributes=None,
        )
        result = WhitelistedAddressService._map_envelope_from_dto(dto)
        assert result.signed_address is not None
        assert len(result.signed_address.signatures) == 1
        entry = result.signed_address.signatures[0]
        assert entry.user_signature is not None
        assert entry.user_signature.user_id == "u-1"
        assert entry.user_signature.signature == "sig-base64"
        assert entry.hashes == ["hash1", "hash2"]
        # Flat signatures
        assert len(result.signatures) == 1
        assert result.signatures[0].user_id == "u-1"
        assert result.signatures[0].hash == "hash1"
        assert result.signatures[0].hashes == ["hash1", "hash2"]

    def test_maps_linked_wallets_from_payload(self) -> None:
        payload = {
            "address": "0xtest",
            "linkedWallets": [
                {"id": 1, "path": "m/44'/60'/0'"},
                {"id": 2, "path": "m/44'/60'/1'"},
            ],
        }
        dto = SimpleNamespace(
            metadata=SimpleNamespace(
                hash="h1",
                payload_as_string=json.dumps(payload),
                payloadAsString=None,
            ),
            blockchain="ETH",
            network="mainnet",
            rules_container=None,
            rulesContainer=None,
            rules_signatures=None,
            rulesSignatures=None,
            signed_address=None,
            signedAddress=None,
            attributes=None,
        )
        result = WhitelistedAddressService._map_envelope_from_dto(dto)
        assert len(result.linked_wallets) == 2
        assert result.linked_wallets[0].id == 1
        assert result.linked_wallets[0].path == "m/44'/60'/0'"

    def test_camelcase_fallbacks(self) -> None:
        dto = SimpleNamespace(
            metadata=SimpleNamespace(
                hash="h1",
                payload_as_string=None,
                payloadAsString='{"address":"0x1"}',
            ),
            blockchain="ETH",
            network="mainnet",
            rules_container=None,
            rulesContainer="camelRC",
            rules_signatures=None,
            rulesSignatures="camelRS",
            signed_address=None,
            signedAddress=None,
            attributes=None,
        )
        result = WhitelistedAddressService._map_envelope_from_dto(dto)
        assert result.metadata is not None
        assert result.metadata.payload_as_string == '{"address":"0x1"}'
        assert result.rules_container == "camelRC"
        assert result.rules_signatures == "camelRS"
