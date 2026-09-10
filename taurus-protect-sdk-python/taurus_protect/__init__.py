"""
Taurus-PROTECT SDK for Python.

A Python client library for the Taurus-PROTECT cryptocurrency custody
and transaction management API.

Example:
    >>> from taurus_protect import ProtectClient
    >>>
    >>> with ProtectClient.create(
    ...     host="https://api.protect.taurushq.com",
    ...     api_key="your-api-key",
    ...     api_secret="your-api-secret-hex",
    ... ) as client:
    ...     wallets, _ = client.wallets.list()
    ...     for wallet in wallets:
    ...         print(f"{wallet.name}: {wallet.currency}")
"""

from taurus_protect.client import ProtectClient
from taurus_protect.credentials import Credentials
from taurus_protect.errors import (
    APIError,
    AuthenticationError,
    AuthorizationError,
    parse_required_roles,
    ConfigurationError,
    IntegrityError,
    NotFoundError,
    RateLimitError,
    RequestMetadataError,
    UnverifiedMetadataError,
    ValidationError,
    WhitelistError,
)

__version__ = "1.0.0"

__all__ = [
    "ProtectClient",
    "Credentials",
    "APIError",
    "AuthenticationError",
    "AuthorizationError",
    "parse_required_roles",
    "ConfigurationError",
    "IntegrityError",
    "NotFoundError",
    "RateLimitError",
    "RequestMetadataError",
    "UnverifiedMetadataError",
    "ValidationError",
    "WhitelistError",
]
