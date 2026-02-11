"""Integration tests for UserService."""

from __future__ import annotations

import logging

import pytest

from taurus_protect.client import ProtectClient

logger = logging.getLogger(__name__)


@pytest.mark.integration
def test_get_current_user(client: ProtectClient) -> None:
    """Test getting the current authenticated user."""
    user = client.users.get_current()

    logger.info("Current user:")
    logger.info("  ID: %s", user.id)
    logger.info("  Username: %s", user.username)
    logger.info("  Email: %s", user.email)
    logger.info("  Status: %s", user.status)

    assert user.id is not None
    assert user.email is not None


@pytest.mark.integration
def test_list_users(client: ProtectClient) -> None:
    """Test listing users with pagination."""
    users, pagination = client.users.list(limit=10)

    logger.info("Found %d users", len(users))
    if pagination:
        logger.info("Total items: %d, HasMore: %s", pagination.total_items, pagination.has_more)

    for user in users:
        logger.info("User: ID=%s, Username=%s, Email=%s", user.id, user.username, user.email)

    assert isinstance(users, list)
