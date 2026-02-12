"""Job service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.audit import job_from_dto, jobs_from_dto
from taurus_protect.models.audit import Job
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class JobService(BaseService):
    """
    Service for job operations.

    Provides methods to list and retrieve jobs.

    Example:
        >>> # List jobs
        >>> jobs, pagination = client.jobs.list(limit=50, offset=0)
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

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Job], Optional[Pagination]]:
        """
        List jobs with pagination.

        Args:
            limit: Maximum number of jobs to return (must be positive).
            offset: Number of jobs to skip (must be non-negative).

        Returns:
            Tuple of (jobs list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            # The jobs API returns all jobs (no server-side pagination)
            resp = self._jobs_api.job_service_get_jobs()

            result = getattr(resp, "result", None)
            all_jobs = jobs_from_dto(result) if result else []

            # Apply client-side pagination
            paged_jobs = all_jobs[offset:offset + limit]

            pagination = self._extract_pagination(
                total_items=len(all_jobs),
                offset=offset,
                limit=limit,
            )

            return paged_jobs, pagination
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

            result = getattr(resp, "result", None)
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
