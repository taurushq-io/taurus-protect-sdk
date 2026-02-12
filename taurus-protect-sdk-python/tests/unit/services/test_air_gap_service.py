"""Unit tests for AirGapService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.air_gap_service import AirGapService


class TestGetUnsignedPayload:
    """Tests for AirGapService.get_unsigned_payload()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        air_gap_api = MagicMock()
        service = AirGapService(api_client=api_client, air_gap_api=air_gap_api)
        return service, air_gap_api

    def test_raises_on_invalid_request_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request_id must be positive"):
            service.get_unsigned_payload(request_id=0)

    def test_returns_empty_bytes_when_none(self) -> None:
        service, api = self._make_service()
        api.air_gap_service_get_outgoing_air_gap.return_value = None

        result = service.get_unsigned_payload(request_id=1)

        assert result == bytes()

    def test_returns_bytes_from_bytearray(self) -> None:
        service, api = self._make_service()
        api.air_gap_service_get_outgoing_air_gap.return_value = bytearray(b"\x01\x02\x03")

        result = service.get_unsigned_payload(request_id=1)

        assert result == b"\x01\x02\x03"
        assert isinstance(result, bytes)

    def test_returns_bytes_from_bytes(self) -> None:
        service, api = self._make_service()
        api.air_gap_service_get_outgoing_air_gap.return_value = b"\xaa\xbb"

        result = service.get_unsigned_payload(request_id=1)

        assert result == b"\xaa\xbb"

    def test_calls_api_with_request_ids(self) -> None:
        service, api = self._make_service()
        api.air_gap_service_get_outgoing_air_gap.return_value = b""

        service.get_unsigned_payload(request_id=42)

        api.air_gap_service_get_outgoing_air_gap.assert_called_once_with(
            body={"request_ids": ["42"]}
        )


class TestSubmitSignedPayload:
    """Tests for AirGapService.submit_signed_payload()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        air_gap_api = MagicMock()
        service = AirGapService(api_client=api_client, air_gap_api=air_gap_api)
        return service, air_gap_api

    def test_raises_on_invalid_request_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request_id must be positive"):
            service.submit_signed_payload(request_id=0, signed_payload=b"\x01")

    def test_raises_on_empty_payload(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="signed_payload cannot be empty"):
            service.submit_signed_payload(request_id=1, signed_payload=b"")

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        api.air_gap_service_submit_incoming_air_gap.return_value = None

        signed = b"\x01\x02\x03"
        service.submit_signed_payload(request_id=1, signed_payload=signed)

        api.air_gap_service_submit_incoming_air_gap.assert_called_once_with(
            body={"payload": signed}
        )
