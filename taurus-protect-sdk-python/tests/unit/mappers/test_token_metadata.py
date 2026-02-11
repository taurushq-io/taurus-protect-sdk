"""Unit tests for token metadata mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.token_metadata import (
    crypto_punk_metadata_from_dto,
    fa_token_metadata_from_dto,
    token_metadata_from_dto,
)


class TestTokenMetadataFromDto:
    """Tests for token_metadata_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            name="My Token",
            description="A test token",
            decimals="18",
            uri="https://example.com/token.json",
            data_type="image/png",
            base64_data="aGVsbG8=",
        )
        result = token_metadata_from_dto(dto)
        assert result is not None
        assert result.name == "My Token"
        assert result.description == "A test token"
        assert result.decimals == "18"
        assert result.uri == "https://example.com/token.json"

    def test_returns_none_for_none(self) -> None:
        assert token_metadata_from_dto(None) is None

    def test_handles_empty_fields(self) -> None:
        dto = SimpleNamespace(
            name=None, description=None, decimals=None,
            uri=None, data_type=None, base64_data=None,
        )
        result = token_metadata_from_dto(dto)
        assert result is not None
        assert result.name == ""
        assert result.description == ""


class TestFATokenMetadataFromDto:
    """Tests for fa_token_metadata_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            name="FA Token",
            symbol="FAT",
            decimals="6",
            description="A Tezos FA token",
            uri="https://example.com/fa.json",
            data_type="application/json",
            base64_data="e30=",
        )
        result = fa_token_metadata_from_dto(dto)
        assert result is not None
        assert result.name == "FA Token"
        assert result.symbol == "FAT"
        assert result.decimals == "6"

    def test_returns_none_for_none(self) -> None:
        assert fa_token_metadata_from_dto(None) is None


class TestCryptoPunkMetadataFromDto:
    """Tests for crypto_punk_metadata_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            punk_index="42",
            image_url="https://example.com/punk42.png",
            attributes={"trait": "mohawk"},
        )
        result = crypto_punk_metadata_from_dto(dto)
        assert result is not None
        assert result.punk_index == "42"
        assert result.image_url == "https://example.com/punk42.png"
        assert result.attributes == {"trait": "mohawk"}

    def test_returns_none_for_none(self) -> None:
        assert crypto_punk_metadata_from_dto(None) is None

    def test_handles_none_attributes(self) -> None:
        dto = SimpleNamespace(
            punk_index="1",
            image_url="url",
            attributes=None,
        )
        result = crypto_punk_metadata_from_dto(dto)
        assert result is not None
        assert result.attributes is None
