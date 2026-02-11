"""Unit tests for AuditService."""

from __future__ import annotations

from datetime import datetime
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.audit_service import AuditService


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
