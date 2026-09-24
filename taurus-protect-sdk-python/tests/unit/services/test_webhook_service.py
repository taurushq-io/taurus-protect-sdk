"""Unit tests for WebhookService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import WebhooksApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.webhook_service import WebhookService
from tests.unit.transport_stub import StubTransport, api_client


class TestWebhookServiceList:
    """WebhookService.list: a cursor list (the offset it took was never sent)."""

    def _service(self) -> WebhookService:
        ac = api_client()
        return WebhookService(ac, WebhooksApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "webhooks": [{"id": "w1", "url": "https://hook.example"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            webhooks, page = self._service().list(
                page_size=5, type="TRANSACTION", url="https://hook.example", sort_order="ASC"
            )

        assert [w.id for w in webhooks] == ["w1"]
        assert page == CursorPage(page_size=5, next_cursor="n", has_more=True)
        assert transport.last.query == sorted(
            [
                ("type", "TRANSACTION"),
                ("url", "https://hook.example"),
                ("sortOrder", "ASC"),
                ("cursor.pageSize", "5"),
            ]
        )

    def test_page_two_is_reachable(self) -> None:
        with StubTransport() as transport:
            self._service().list(cursor="n")

        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            webhooks, page = self._service().list()

        assert webhooks == []
        assert page == CursorPage(page_size=20)


class TestWebhookServiceGet:
    """WebhookService.get walks the pages instead of reading one page of 100."""

    def _service(self) -> WebhookService:
        ac = api_client()
        return WebhookService(ac, WebhooksApi(ac))

    def test_get_raises_on_empty_id(self) -> None:
        with pytest.raises(ValueError, match="webhook_id"):
            self._service().get("")

    def test_found_on_page_two(self) -> None:
        with StubTransport(
            {"webhooks": [{"id": "w1"}], "cursor": {"currentPage": "n", "hasNext": True}},
            {"webhooks": [{"id": "w2"}]},
        ) as transport:
            webhook = self._service().get("w2")

        assert webhook.id == "w2"
        assert transport.requests[1].param("cursor.currentPage") == "n"
        assert [r.param("cursor.pageSize") for r in transport.requests] == ["100", "100"]

    def test_not_found_after_the_last_page(self) -> None:
        with StubTransport({"webhooks": [{"id": "w1"}]}):
            with pytest.raises(NotFoundError):
                self._service().get("w2")


class TestWebhookServiceCreate:
    """Tests for WebhookService.create()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhooks_api = MagicMock()
        service = WebhookService(api_client=api_client, webhooks_api=webhooks_api)
        return service, webhooks_api

    def test_create_raises_on_empty_url(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="url"):
            service.create(url="", events=["REQUEST_CREATED"])

    def test_create_raises_on_empty_events(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="events cannot be empty"):
            service.create(url="https://example.com", events=[])


class TestWebhookServiceDelete:
    """Tests for WebhookService.delete()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhooks_api = MagicMock()
        service = WebhookService(api_client=api_client, webhooks_api=webhooks_api)
        return service, webhooks_api

    def test_delete_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="webhook_id"):
            service.delete("")

    def test_delete_calls_api(self) -> None:
        service, api = self._make_service()
        api.webhook_service_delete_webhook.return_value = None

        service.delete("wh-123")

        api.webhook_service_delete_webhook.assert_called_once_with(id="wh-123")
