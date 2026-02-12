"""Shared helper methods for integration and E2E tests."""

from __future__ import annotations

from typing import Optional

import pytest

from taurus_protect.client import ProtectClient

from tests.testutil.config import get_config


def skip_if_not_enabled() -> None:
    """Skip the current test if tests are not enabled."""
    config = get_config()
    if not config.is_enabled():
        pytest.skip(
            "Skipping test. Set PROTECT_INTEGRATION_TEST=true or configure test.properties."
        )


def get_test_client(identity_index: int = 1) -> ProtectClient:
    """Create a ProtectClient for the identity at the given 1-based index."""
    config = get_config()
    identity = config.get_identity(identity_index)

    return ProtectClient.create(
        host=config.host,
        api_key=identity.api_key,
        api_secret=identity.api_secret,
        super_admin_keys_pem=config.get_super_admin_keys(),
        min_valid_signatures=config.get_min_valid_signatures(),
    )


def get_private_key(identity_index: int = 1) -> Optional[bytes]:
    """Return the raw PEM private key string for the identity at the given 1-based index.

    Returns None if the identity has no private key configured.
    """
    config = get_config()
    identity = config.get_identity(identity_index)
    if not identity.has_private_key():
        return None
    return identity.private_key.encode("utf-8")


def skip_if_insufficient_identities(required: int) -> None:
    """Skip the current test if fewer than required identities are configured."""
    config = get_config()
    if config.identity_count < required:
        pytest.skip(
            f"Skipping: need {required} identities but only {config.identity_count} configured."
        )
