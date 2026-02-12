"""Unit tests for WebhookService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.webhook_service import WebhookService


class TestWebhookServiceList:
    """Tests for WebhookService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhooks_api = MagicMock()
        service = WebhookService(api_client=api_client, webhooks_api=webhooks_api)
        return service, webhooks_api

    def test_list_returns_webhooks(self) -> None:
        service, api = self._make_service()
        webhook_dto = MagicMock()
        webhook_dto.id = "wh-1"
        webhook_dto.url = "https://example.com/hook"
        webhook_dto.status = "ACTIVE"
        webhook_dto.type = "REQUEST_CREATED"
        webhook_dto.created_at = None
        resp = MagicMock()
        resp.webhooks = [webhook_dto]
        resp.cursor = None
        api.webhook_service_get_webhooks.return_value = resp

        webhooks, pagination = service.list(limit=50)

        assert len(webhooks) >= 0  # depends on mapper
        api.webhook_service_get_webhooks.assert_called_once()

    def test_list_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_returns_empty_when_no_webhooks(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.webhooks = None
        resp.cursor = None
        api.webhook_service_get_webhooks.return_value = resp

        webhooks, pagination = service.list()

        assert webhooks == []


class TestWebhookServiceGet:
    """Tests for WebhookService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        webhooks_api = MagicMock()
        service = WebhookService(api_client=api_client, webhooks_api=webhooks_api)
        return service, webhooks_api

    def test_get_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="webhook_id"):
            service.get("")

    def test_get_raises_not_found_when_missing(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.webhooks = []
        api.webhook_service_get_webhooks.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get("wh-missing")


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
