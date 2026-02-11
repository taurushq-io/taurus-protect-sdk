"""Unit tests for JobService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.job_service import JobService


class TestJobServiceList:
    """Tests for JobService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        jobs_api = MagicMock()
        service = JobService(api_client=api_client, jobs_api=jobs_api)
        return service, jobs_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.job_service_get_jobs.return_value = resp

        jobs, pagination = service.list()

        assert jobs == []

    def test_calls_api_without_pagination_params(self) -> None:
        """The jobs API takes no params; pagination is applied client-side."""
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.job_service_get_jobs.return_value = resp

        service.list(limit=25, offset=10)

        api.job_service_get_jobs.assert_called_once_with()

    def test_client_side_pagination(self) -> None:
        """Verify client-side pagination slices the full result set."""
        service, api = self._make_service()
        resp = MagicMock()
        # Simulate 5 jobs returned by the API
        mock_jobs = [MagicMock() for _ in range(5)]
        for i, j in enumerate(mock_jobs):
            j.id = str(i + 1)
            j.type = "test"
            j.description = f"Job {i + 1}"
            j.status = None
            j.created_at = None
        resp.result = mock_jobs
        api.job_service_get_jobs.return_value = resp

        jobs, pagination = service.list(limit=2, offset=1)

        # Should get jobs at indices 1 and 2 (0-based slice [1:3])
        assert len(jobs) == 2
        assert pagination is not None
        assert pagination.total_items == 5


class TestJobServiceGet:
    """Tests for JobService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        jobs_api = MagicMock()
        service = JobService(api_client=api_client, jobs_api=jobs_api)
        return service, jobs_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="job_id"):
            service.get(job_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.job_service_get_job.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(job_id="123")

    def test_wraps_api_error(self) -> None:
        service, api = self._make_service()
        error = Exception("connection refused")
        error.status = 503
        error.body = None
        error.headers = {}
        api.job_service_get_job.side_effect = error

        from taurus_protect.errors import APIError

        with pytest.raises(APIError):
            service.get(job_id="123")
