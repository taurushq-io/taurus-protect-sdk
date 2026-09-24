"""Unit tests for AuditService."""

from __future__ import annotations

from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import AuditApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.audit_service import AuditService
from tests.unit.transport_stub import StubTransport, api_client


class TestExportAuditTrails:
    """Tests for AuditService.export_audit_trails()."""

    def _make_service(self) -> tuple:
        """Create an AuditService with a mocked audit API."""
        api_client = MagicMock()
        audit_api = MagicMock()
        service = AuditService(api_client=api_client, audit_api=audit_api)
        return service, audit_api

    def test_export_returns_result(self) -> None:
        service, audit_api = self._make_service()
        reply = MagicMock()
        reply.result = "col1,col2\nval1,val2"
        audit_api.audit_service_export_audit_trails.return_value = reply

        result = service.export_audit_trails(format="csv")

        assert result == "col1,col2\nval1,val2"
        audit_api.audit_service_export_audit_trails.assert_called_once_with(
            external_user_id=None,
            entities=None,
            actions=None,
            creation_date_from=None,
            creation_date_to=None,
            format="csv",
        )

    def test_export_passes_all_parameters(self) -> None:
        service, audit_api = self._make_service()
        reply = MagicMock()
        reply.result = "{}"
        audit_api.audit_service_export_audit_trails.return_value = reply

        from_dt = datetime(2025, 1, 1)
        to_dt = datetime(2025, 6, 30)

        service.export_audit_trails(
            external_user_id="user-42",
            entities=["wallet", "address"],
            actions=["create", "delete"],
            from_date=from_dt,
            to_date=to_dt,
            format="json",
        )

        audit_api.audit_service_export_audit_trails.assert_called_once_with(
            external_user_id="user-42",
            entities=["wallet", "address"],
            actions=["create", "delete"],
            creation_date_from=from_dt,
            creation_date_to=to_dt,
            format="json",
        )

    def test_export_returns_empty_string_when_result_is_none(self) -> None:
        service, audit_api = self._make_service()
        reply = MagicMock()
        reply.result = None
        audit_api.audit_service_export_audit_trails.return_value = reply

        result = service.export_audit_trails()

        assert result == ""

    def test_export_with_no_result_attribute(self) -> None:
        service, audit_api = self._make_service()
        reply = MagicMock(spec=[])  # No attributes
        audit_api.audit_service_export_audit_trails.return_value = reply

        result = service.export_audit_trails()

        assert result == ""

    def test_export_wraps_api_error(self) -> None:
        service, audit_api = self._make_service()
        error = Exception("connection refused")
        error.status = 503
        error.body = None
        error.headers = {}
        audit_api.audit_service_export_audit_trails.side_effect = error

        from taurus_protect.errors import APIError

        with pytest.raises(APIError):
            service.export_audit_trails()

    def test_export_propagates_value_error(self) -> None:
        service, audit_api = self._make_service()
        audit_api.audit_service_export_audit_trails.side_effect = ValueError("bad")

        with pytest.raises(ValueError, match="bad"):
            service.export_audit_trails()


class TestList:
    """AuditService.list over the real generated client (cursor)."""

    def _service(self) -> AuditService:
        ac = api_client()
        return AuditService(ac, AuditApi(ac))

    def test_rows_page_and_filters(self) -> None:
        start = datetime(2026, 1, 2, tzinfo=timezone.utc)
        reply = {"result": [{"id": "1"}], "cursor": {"currentPage": "n", "hasNext": True}}
        with StubTransport(reply) as transport:
            audits, page = self._service().list(
                page_size=5,
                external_user_id="u",
                entities=["wallet"],
                actions=["create"],
                from_date=start,
                sort_by=["creationDate"],
                sort_order="ASC",
            )

        assert [a.id for a in audits] == ["1"]
        assert page == CursorPage(page_size=5, next_cursor="n", has_more=True)
        query = dict(transport.last.query)
        assert query["externalUserId"] == "u"
        assert query["entities"] == "wallet"
        assert query["actions"] == "create"
        assert query["creationDateFrom"].startswith("2026-01-02T00:00:00")
        assert query["sorting.sortBy"] == "creationDate"
        assert query["sorting.sortOrder"] == "ASC"
        assert query["cursor.pageSize"] == "5"

    def test_page_two_is_reachable(self) -> None:
        """The list could not send a cursor at all."""
        with StubTransport() as transport:
            self._service().list(cursor="n")

        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"


class TestGet:
    """AuditService.get walks the pages instead of asking for 1000 rows at once."""

    def _service(self) -> AuditService:
        ac = api_client()
        return AuditService(ac, AuditApi(ac))

    def test_walks_pages_until_found(self) -> None:
        with StubTransport(
            {"result": [{"id": "1"}], "cursor": {"currentPage": "n", "hasNext": True}},
            {"result": [{"id": "2"}], "cursor": {"currentPage": "m"}},
        ) as transport:
            audit = self._service().get("2")

        assert audit.id == "2"
        assert [r.param("cursor.pageSize") for r in transport.requests] == ["100", "100"]
        assert transport.requests[1].param("cursor.currentPage") == "n"

    def test_not_found_after_the_last_page(self) -> None:
        with StubTransport({"result": [{"id": "1"}]}) as transport:
            with pytest.raises(NotFoundError):
                self._service().get("2")

        assert len(transport.requests) == 1
