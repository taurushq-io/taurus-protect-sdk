"""Unit tests for WebhookCallService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.webhook_call_service import (
    ApiRequestCursor,
    WebhookCallResult,
    WebhookCallService,
)


class TestGetWebhookCalls:
    """Tests for WebhookCallService.get_webhook_calls()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhook_calls_api = MagicMock()
        service = WebhookCallService(
            api_client=api_client, webhook_calls_api=webhook_calls_api
        )
        return service, webhook_calls_api

    def test_returns_empty_result_when_no_calls(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.calls = None
        resp.result = None
        resp.cursor = None
        api.webhook_service_get_webhook_calls.return_value = resp

        result = service.get_webhook_calls()

        assert isinstance(result, WebhookCallResult)
        assert result.calls == []
        assert result.has_more is False

    def test_passes_filter_parameters(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.calls = []
        resp.cursor = None
        api.webhook_service_get_webhook_calls.return_value = resp

        service.get_webhook_calls(
            event_id="evt-1",
            webhook_id="wh-1",
            status="SUCCESS",
            sort_order="ASC",
        )

        api.webhook_service_get_webhook_calls.assert_called_once_with(
            event_id="evt-1",
            webhook_id="wh-1",
            status="SUCCESS",
            cursor_current_page=None,
            cursor_page_request=None,
            cursor_page_size=None,
            sort_order="ASC",
        )

    def test_uses_cursor_when_provided(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.calls = []
        resp.cursor = None
        api.webhook_service_get_webhook_calls.return_value = resp

        cursor = ApiRequestCursor(current_page="page2", page_request="NEXT", page_size=25)
        service.get_webhook_calls(cursor=cursor)

        api.webhook_service_get_webhook_calls.assert_called_once_with(
            event_id=None,
            webhook_id=None,
            status=None,
            cursor_current_page="page2",
            cursor_page_request="NEXT",
            cursor_page_size="25",
            sort_order=None,
        )


class TestWebhookCallServiceList:
    """Tests for WebhookCallService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhook_calls_api = MagicMock()
        service = WebhookCallService(
            api_client=api_client, webhook_calls_api=webhook_calls_api
        )
        return service, webhook_calls_api

    def test_list_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_returns_empty_when_no_calls(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.calls = None
        resp.result = None
        resp.cursor = None
        api.webhook_service_get_webhook_calls.return_value = resp

        calls, pagination = service.list()

        assert calls == []
        assert pagination is not None


class TestWebhookCallServiceGet:
    """Tests for WebhookCallService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhook_calls_api = MagicMock()
        service = WebhookCallService(
            api_client=api_client, webhook_calls_api=webhook_calls_api
        )
        return service, webhook_calls_api

    def test_get_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="call_id"):
            service.get("")

    def test_get_raises_not_found(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.calls = []
        api.webhook_service_get_webhook_calls.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get("call-missing")
