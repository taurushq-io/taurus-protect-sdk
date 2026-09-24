"""Job service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List

from taurus_protect.mappers.audit import job_from_dto, jobs_from_dto
from taurus_protect.models.audit import Job
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class JobService(BaseService):
    """
    Service for job operations.

    Provides methods to list and retrieve jobs.

    Example:
        >>> # List jobs
        >>> jobs = client.jobs.list()
        >>> for job in jobs:
        ...     print(f"{job.id}: {job.description}")
        >>>
        >>> # Get single job
        >>> job = client.jobs.get("123")
        >>> print(f"Type: {job.type}")
    """

    def __init__(self, api_client: Any, jobs_api: Any) -> None:
        """
        Initialize job service.

        Args:
            api_client: The OpenAPI client instance.
            jobs_api: The JobsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._jobs_api = jobs_api

    def list(self) -> List[Job]:
        """
        List jobs: every job the endpoint returns, which does not page.

        Returns:
            Every job.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._jobs_api.job_service_get_jobs()

            return jobs_from_dto(resp.jobs or [])
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, job_id: str) -> Job:
        """
        Get a job by ID.

        Args:
            job_id: The job ID to retrieve.

        Returns:
            The job.

        Raises:
            ValueError: If job_id is invalid.
            NotFoundError: If job not found.
            APIError: If API request fails.
        """
        self._validate_required(job_id, "job_id")

        try:
            resp = self._jobs_api.job_service_get_job(job_id)

            result = resp.job
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Job {job_id} not found")

            job = job_from_dto(result)
            if job is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Job {job_id} not found")

            return job
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e
