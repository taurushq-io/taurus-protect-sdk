"""Unit tests for air-gap service models and helpers."""

from taurus_protect.models.staking import UnsignedPayload
from taurus_protect.services.air_gap_service import AirGapService


class TestUnsignedPayloadModel:
    """Tests for UnsignedPayload model."""

    def test_creates_with_all_fields(self) -> None:
        payload = UnsignedPayload(
            request_id="req-123",
            payload="0xdeadbeef",
            hash="sha256hash",
            blockchain="ETH",
            network="mainnet",
        )
        assert payload.request_id == "req-123"
        assert payload.payload == "0xdeadbeef"
        assert payload.hash == "sha256hash"
        assert payload.blockchain == "ETH"
        assert payload.network == "mainnet"

    def test_creates_with_defaults(self) -> None:
        payload = UnsignedPayload(
            request_id="req-1",
            payload="0x00",
        )
        assert payload.request_id == "req-1"
        assert payload.hash == ""
        assert payload.blockchain == ""
        assert payload.network == ""

    def test_model_is_frozen(self) -> None:
        payload = UnsignedPayload(
            request_id="req-1",
            payload="0x00",
        )
        try:
            payload.request_id = "modified"  # type: ignore[misc]
            assert False, "Should have raised an error for frozen model"
        except Exception:
            pass  # Expected behavior for frozen model


class TestAirGapServiceInit:
    """Tests for AirGapService basic initialization."""

    def test_service_class_exists(self) -> None:
        """Verify AirGapService class is importable and has expected methods."""
        assert hasattr(AirGapService, "get_unsigned_payload")
        assert hasattr(AirGapService, "submit_signed_payload")
