"""Unit tests for webhook call mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.mappers.webhook import webhook_call_from_dto, webhook_calls_from_dto


class TestWebhookCallFromDto:
    """Tests for webhook_call_from_dto function."""

    def test_maps_all_fields(self) -> None:
        created = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        updated = datetime(2024, 6, 15, 12, 1, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="wc-1",
            event_id="evt-100",
            webhook_id="wh-10",
            payload='{"type": "TRANSACTION"}',
            status="SUCCESS",
            status_message="200 OK",
            attempts=1,
            created_at=created,
            updated_at=updated,
        )
        result = webhook_call_from_dto(dto)
        assert result is not None
        assert result.id == "wc-1"
        assert result.event_id == "evt-100"
        assert result.webhook_id == "wh-10"
        assert result.payload == '{"type": "TRANSACTION"}'
        assert result.status == "SUCCESS"
        assert result.status_message == "200 OK"
        assert result.attempts == 1
        assert result.created_at == created
        assert result.updated_at == updated

    def test_returns_none_for_none(self) -> None:
        assert webhook_call_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="wc-2",
            event_id=None,
            webhook_id=None,
            payload=None,
            status=None,
            status_message=None,
            attempts=None,
            created_at=None,
            updated_at=None,
        )
        result = webhook_call_from_dto(dto)
        assert result is not None
        assert result.id == "wc-2"
        assert result.event_id == ""
        assert result.attempts == 0
        assert result.created_at is None

    def test_attempts_as_string(self) -> None:
        dto = SimpleNamespace(
            id="wc-3",
            event_id="e1",
            webhook_id="w1",
            payload="{}",
            status="FAILED",
            status_message="500",
            attempts="3",
            created_at=None,
            updated_at=None,
        )
        result = webhook_call_from_dto(dto)
        assert result is not None
        assert result.attempts == 3


class TestWebhookCallsFromDto:
    """Tests for webhook_calls_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="wc-1", event_id="e1", webhook_id="w1", payload="{}",
                status="SUCCESS", status_message="", attempts=1,
                created_at=None, updated_at=None,
            ),
            SimpleNamespace(
                id="wc-2", event_id="e2", webhook_id="w1", payload="{}",
                status="FAILED", status_message="timeout", attempts=3,
                created_at=None, updated_at=None,
            ),
        ]
        result = webhook_calls_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "wc-1"
        assert result[1].status == "FAILED"

    def test_returns_empty_for_none(self) -> None:
        assert webhook_calls_from_dto(None) == []

    def test_filters_none_dtos(self) -> None:
        dtos = [
            None,
            SimpleNamespace(
                id="wc-1", event_id="e", webhook_id="w", payload="",
                status="OK", status_message="", attempts=0,
                created_at=None, updated_at=None,
            ),
        ]
        result = webhook_calls_from_dto(dtos)
        assert len(result) == 1
