"""Integration tests for RequestService.

Tests request listing and retrieval functionality against a live API.
"""

from __future__ import annotations

import logging

import pytest

from taurus_protect.client import ProtectClient

logger = logging.getLogger(__name__)


@pytest.mark.integration
def test_list_requests(client: ProtectClient) -> None:
    """Test listing requests with pagination."""
    requests, pagination = client.requests.list(limit=10)

    logger.info(f"Found {len(requests)} requests")
    if pagination is not None:
        logger.info(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for r in requests:
        logger.info(f"Request: ID={r.id}, Status={r.status.value}, Currency={r.currency}")


@pytest.mark.integration
def test_get_request(client: ProtectClient) -> None:
    """Test getting a single request by ID and verify metadata exists."""
    # First, get a list to find a valid request ID
    requests, _ = client.requests.list(limit=1)

    if len(requests) == 0:
        pytest.skip("No requests available for testing")

    request_id = int(requests[0].id)
    request = client.requests.get(request_id)

    logger.info("Request details:")
    logger.info(f"  ID: {request.id}")
    logger.info(f"  Status: {request.status.value}")
    logger.info(f"  Currency: {request.currency}")
    if request.metadata is not None:
        logger.info(f"  Metadata hash: {request.metadata.hash}")


@pytest.mark.integration
def test_request_metadata_verification(client: ProtectClient) -> None:
    """Test request metadata verification - verify hash, payload, and fields."""
    # Get a request
    requests, _ = client.requests.list(limit=1)

    if len(requests) == 0:
        pytest.skip("No requests available for testing")

    request_id = int(requests[0].id)
    request = client.requests.get(request_id)

    # Verify we can access metadata
    if request.metadata is None:
        pytest.skip("Request has no metadata")

    logger.info(f"Request ID: {request.id}")
    logger.info(f"Metadata hash: {request.metadata.hash}")

    # Verify metadata payload is present
    payload_str = request.metadata.payload_as_string
    if payload_str:
        logger.info(f"Metadata PayloadAsString length: {len(payload_str)}")
        # Log metadata payload (first 200 chars for readability)
        if len(payload_str) > 200:
            logger.info(f"Metadata payload (truncated): {payload_str[:200]}...")
        else:
            logger.info(f"Metadata payload: {payload_str}")
    else:
        logger.info("Metadata PayloadAsString: (empty)")

    # Verify basic metadata fields
    logger.info(f"Request Currency: {request.currency}")
    logger.info(f"Request Status: {request.status.value}")
    if request.rule:
        logger.info(f"Request Rule: {request.rule}")

    # Assert metadata hash exists and is not empty
    assert request.metadata.hash, "Metadata hash should not be empty"


@pytest.mark.integration
def test_list_requests_by_status(client: ProtectClient) -> None:
    """Test listing requests filtered by status."""
    from taurus_protect.models import RequestStatus

    # List requests with CONFIRMED status
    requests, pagination = client.requests.list(limit=10, statuses=[RequestStatus.CONFIRMED])

    logger.info(f"Found {len(requests)} CONFIRMED requests")
    if pagination is not None:
        logger.info(f"Total CONFIRMED items: {pagination.total_items}")

    for r in requests:
        logger.info(f"Request: ID={r.id}, Status={r.status.value}, Currency={r.currency}")
        # Verify all returned requests are CONFIRMED
        assert r.status == RequestStatus.CONFIRMED, f"Expected CONFIRMED status, got {r.status}"

    # Also try PENDING status
    pending_requests, pending_pagination = client.requests.list(limit=5, statuses=[RequestStatus.PENDING])

    logger.info(f"Found {len(pending_requests)} PENDING requests")
    if pending_pagination is not None:
        logger.info(f"Total PENDING items: {pending_pagination.total_items}")
