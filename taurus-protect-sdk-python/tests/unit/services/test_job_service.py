"""Unit tests for JobService."""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi import JobsApi
from taurus_protect.services.job_service import JobService
from tests.unit.transport_stub import Reply, StubTransport, api_client


class TestJobServiceList:
    """JobService.list: the endpoint does not page, so every job it returns comes back."""

    def _service(self) -> JobService:
        ac = api_client()
        return JobService(ac, JobsApi(ac))

    def test_returns_every_job_from_the_jobs_field(self) -> None:
        """The list read ``result``, a field the reply does not have, and was always empty."""
        reply = {"jobs": [{"name": f"job{i}"} for i in range(25)]}
        with StubTransport(reply) as transport:
            jobs = self._service().list()

        assert [j.name for j in jobs] == [f"job{i}" for i in range(25)]
        assert [j.id for j in jobs] == [j.name for j in jobs]
        assert transport.last.query == [], "the endpoint takes no paging parameters"


class TestJobServiceGet:
    """JobService.get reads the reply's ``job``."""

    def _service(self) -> JobService:
        ac = api_client()
        return JobService(ac, JobsApi(ac))

    def test_raises_on_empty_id(self) -> None:
        with pytest.raises(ValueError, match="job_id"):
            self._service().get(job_id="")

    def test_returns_the_job(self) -> None:
        with StubTransport({"job": {"name": "sync"}}) as transport:
            job = self._service().get(job_id="sync")

        assert job.name == "sync"
        assert transport.last.path == "/api/rest/v1/jobs/sync"

    def test_raises_not_found_when_absent(self) -> None:
        from taurus_protect.errors import NotFoundError

        with StubTransport({}):
            with pytest.raises(NotFoundError):
                self._service().get(job_id="sync")

    def test_wraps_api_error(self) -> None:
        from taurus_protect.errors import APIError

        with StubTransport(Reply({"message": "down"}, status=503)):
            with pytest.raises(APIError):
                self._service().get(job_id="123")
