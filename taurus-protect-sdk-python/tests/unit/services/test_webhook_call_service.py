"""Unit tests for WebhookCallService."""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi import WebhookCallsApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.webhook_call_service import (
    ApiRequestCursor,
    WebhookCallResult,
    WebhookCallService,
)
from tests.unit.transport_stub import StubTransport, api_client


class TestGetWebhookCalls:
    """get_webhook_calls: has_more comes from hasNext, and the next cursor is returned."""

    def _service(self) -> WebhookCallService:
        ac = api_client()
        return WebhookCallService(ac, WebhookCallsApi(ac))

    def test_returns_empty_result_when_no_calls(self) -> None:
        with StubTransport({}) as transport:
            result = self._service().get_webhook_calls()

        assert isinstance(result, WebhookCallResult)
        assert result.calls == []
        assert result.page == CursorPage(page_size=20)
        assert transport.last.query == [("cursor.pageSize", "20")]

    def test_last_page_is_not_has_more(self) -> None:
        """has_more was read from the presence of currentPage, which the last page also has."""
        with StubTransport({"calls": [{"id": "c1"}], "cursor": {"currentPage": "last"}}):
            result = self._service().get_webhook_calls()

        assert result.page.has_more is False
        assert result.page.next_cursor == ""

    def test_passes_filter_parameters(self) -> None:
        with StubTransport() as transport:
            self._service().get_webhook_calls(
                event_id="e", webhook_id="w", status="FAILED", sort_order="ASC", page_size=5
            )

        assert transport.last.query == sorted(
            [
                ("eventID", "e"),
                ("webhookID", "w"),
                ("status", "FAILED"),
                ("sortOrder", "ASC"),
                ("cursor.pageSize", "5"),
            ]
        )

    def test_string_cursor_continues(self) -> None:
        with StubTransport() as transport:
            self._service().get_webhook_calls(cursor="n+/=")

        assert transport.last.param("cursor.currentPage") == "n+/="
        assert transport.last.param("cursor.pageRequest") == "NEXT"

    def test_low_level_cursor(self) -> None:
        with StubTransport() as transport:
            self._service().get_webhook_calls(
                cursor=ApiRequestCursor(current_page="p", page_request="PREVIOUS", page_size=7)
            )

        assert transport.last.param("cursor.currentPage") == "p"
        assert transport.last.param("cursor.pageRequest") == "PREVIOUS"
        assert transport.last.param("cursor.pageSize") == "7"


class TestWebhookCallServiceList:
    """WebhookCallService.list: rows, page and continuation."""

    def _service(self) -> WebhookCallService:
        ac = api_client()
        return WebhookCallService(ac, WebhookCallsApi(ac))

    def test_list_raises_on_invalid_page_size(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                self._service().list(page_size=101)

        assert transport.requests == []

    def test_rows_and_next_cursor(self) -> None:
        with StubTransport(
            {"calls": [{"id": "c1"}], "cursor": {"currentPage": "n", "hasNext": True}}
        ):
            calls, page = self._service().list(webhook_id="w")

        assert [c.id for c in calls] == ["c1"]
        assert page == CursorPage(page_size=20, next_cursor="n", has_more=True)


class TestWebhookCallServiceGet:
    """WebhookCallService.get walks the pages."""

    def _service(self) -> WebhookCallService:
        ac = api_client()
        return WebhookCallService(ac, WebhookCallsApi(ac))

    def test_get_raises_on_empty_id(self) -> None:
        with pytest.raises(ValueError, match="call_id"):
            self._service().get("")

    def test_found_on_page_two(self) -> None:
        with StubTransport(
            {"calls": [{"id": "c1"}], "cursor": {"currentPage": "n", "hasNext": True}},
            {"calls": [{"id": "c2"}]},
        ) as transport:
            call = self._service().get("c2")

        assert call.id == "c2"
        assert transport.requests[1].param("cursor.currentPage") == "n"

    def test_get_raises_not_found(self) -> None:
        with StubTransport({"calls": [{"id": "c1"}]}):
            with pytest.raises(NotFoundError):
                self._service().get("c2")
