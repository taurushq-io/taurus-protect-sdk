"""Unit tests for user device mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect._internal.openapi.models.tgvalidatord_user_device_pairing_info_status import (  # noqa: E501
    TgvalidatordUserDevicePairingInfoStatus,
)
from taurus_protect.mappers.user_device import (
    user_device_pairing_from_dto,
    user_device_pairing_info_from_dto,
)


class TestUserDevicePairingFromDto:
    """Tests for user_device_pairing_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            pairing_id="pair-1",
            status="PENDING",
            creation_date="2024-01-15T10:30:00Z",
            expiration_date="2024-01-16T10:30:00Z",
        )
        result = user_device_pairing_from_dto(dto)
        assert result is not None
        assert result.pairing_id == "pair-1"
        assert result.status == "PENDING"
        assert result.created_at is not None
        assert result.expires_at is not None

    def test_returns_none_for_none(self) -> None:
        assert user_device_pairing_from_dto(None) is None


class TestUserDevicePairingInfoFromDto:
    """Tests for user_device_pairing_info_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            pairing_id="pair-1",
            user_id="user-1",
            status="ACTIVE",
            device_name="iPhone 15",
            device_type="IOS",
            encryption_key="enc-key-base64",
            creation_date="2024-01-15T10:30:00Z",
            expiration_date="2024-01-16T10:30:00Z",
        )
        result = user_device_pairing_info_from_dto(dto)
        assert result is not None
        assert result.pairing_id == "pair-1"
        assert result.user_id == "user-1"
        assert result.device_name == "iPhone 15"
        assert result.device_type == "IOS"
        assert result.encryption_key == "enc-key-base64"

    @pytest.mark.parametrize("wire", ["APPROVED", "FUTURE_STATUS"], ids=["known", "unknown"])
    def test_status_is_the_wire_string_not_the_enum_member(self, wire: str) -> None:
        dto = SimpleNamespace(status=TgvalidatordUserDevicePairingInfoStatus(wire))

        result = user_device_pairing_info_from_dto(dto)

        assert type(result.status) is str
        assert result.status == wire

    def test_returns_none_for_none(self) -> None:
        assert user_device_pairing_info_from_dto(None) is None
