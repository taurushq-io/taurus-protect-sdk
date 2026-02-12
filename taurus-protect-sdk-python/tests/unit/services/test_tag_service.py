"""Unit tests for TagService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import APIError, NotFoundError
from taurus_protect.services.tag_service import TagService


class TestGet:
    """Tests for TagService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        tags_api = MagicMock()
        service = TagService(api_client=api_client, tags_api=tags_api)
        return service, tags_api

    def test_get_returns_tag(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        api.tag_service_get_tags.return_value = reply

        mock_tag = MagicMock()
        with patch(
            "taurus_protect.services.tag_service.tag_from_dto",
            return_value=mock_tag,
        ):
            result = service.get("tag-1")

        assert result is mock_tag

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="tag_id"):
            service.get("")

    def test_get_raises_not_found_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = []
        api.tag_service_get_tags.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("tag-999")

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.tag_service_get_tags.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("tag-999")


class TestList:
    """Tests for TagService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        tags_api = MagicMock()
        service = TagService(api_client=api_client, tags_api=tags_api)
        return service, tags_api

    def test_list_returns_tags(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock(), MagicMock(), MagicMock()]
        api.tag_service_get_tags.return_value = reply

        with patch(
            "taurus_protect.services.tag_service.tags_from_dto",
            return_value=[MagicMock(), MagicMock(), MagicMock()],
        ):
            result = service.list()

        assert len(result) == 3

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_applies_client_side_pagination(self) -> None:
        service, api = self._make_service()

        tags = [MagicMock() for _ in range(10)]
        reply = MagicMock()
        reply.result = [MagicMock() for _ in range(10)]
        api.tag_service_get_tags.return_value = reply

        with patch(
            "taurus_protect.services.tag_service.tags_from_dto",
            return_value=tags,
        ):
            result = service.list(limit=3, offset=2)

        assert len(result) == 3


class TestCreate:
    """Tests for TagService.create()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        tags_api = MagicMock()
        service = TagService(api_client=api_client, tags_api=tags_api)
        return service, tags_api

    def test_create_raises_for_empty_name(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="name"):
            service.create("", "#FF0000")

    def test_create_raises_for_empty_color(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="color"):
            service.create("Important", "")

    def test_create_returns_tag(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.tag_service_create_tag.return_value = reply

        mock_tag = MagicMock()
        with patch(
            "taurus_protect.services.tag_service.tag_from_dto",
            return_value=mock_tag,
        ):
            result = service.create("Important", "#FF0000")

        assert result is mock_tag


class TestDelete:
    """Tests for TagService.delete()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        tags_api = MagicMock()
        service = TagService(api_client=api_client, tags_api=tags_api)
        return service, tags_api

    def test_delete_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="tag_id"):
            service.delete("")

    def test_delete_calls_api(self) -> None:
        service, api = self._make_service()

        service.delete("tag-1")

        api.tag_service_delete_tag.assert_called_once_with(id="tag-1")
