"""Credentials - authentication mechanism for ProtectClient."""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Callable, Optional

from taurus_protect.crypto.tpv1 import TPV1Auth
from taurus_protect.errors import ConfigurationError


class Credentials(ABC):
    """The authentication mechanism for a ProtectClient: either static TPV1-HMAC
    (api key + secret) or a Bearer token (static, or resolved per request for
    rotating/per-caller tokens).

    Build one with :meth:`api_key`, :meth:`bearer_token`, or
    :meth:`bearer_token_provider` and pass it to ``ProtectClient.create``.
    SuperAdmin keys are required regardless of the mechanism.
    """

    @staticmethod
    def api_key(api_key: str, api_secret: str) -> "Credentials":
        """TPV1-HMAC credentials (a static shared-service identity).

        Args:
            api_key: The API key.
            api_secret: The API secret as a hex-encoded string.
        """
        return _ApiKeyCredentials(api_key, api_secret)

    @staticmethod
    def bearer_token(token: str) -> "Credentials":
        """A single static Bearer token (the "Authorization: Bearer" header)."""
        if not token:
            raise ConfigurationError("bearer_token cannot be empty")
        return _BearerCredentials(lambda: token)

    @staticmethod
    def bearer_token_provider(provider: Callable[[], str]) -> "Credentials":
        """A Bearer token resolved per request from ``provider()`` (for
        rotating/per-caller tokens)."""
        if provider is None:
            raise ConfigurationError("bearer_token_provider cannot be None")
        return _BearerCredentials(provider)

    @abstractmethod
    def _build_rest_client(self, config: Any) -> Any:
        """Build the authenticating OpenAPI REST client from the configuration."""

    def close(self) -> None:
        """Wipe any sensitive material held by the credentials. No-op by default."""


class _ApiKeyCredentials(Credentials):
    def __init__(self, api_key: str, api_secret: str) -> None:
        if not api_key:
            raise ConfigurationError("api_key cannot be empty")
        if not api_secret:
            raise ConfigurationError("api_secret cannot be empty")
        try:
            self._auth: Optional[TPV1Auth] = TPV1Auth(api_key, api_secret)
        except ValueError as e:
            raise ConfigurationError(str(e)) from e

    def _build_rest_client(self, config: Any) -> Any:
        from taurus_protect.crypto.authenticated_rest import AuthenticatedRESTClient

        return AuthenticatedRESTClient(config, self._auth)

    def close(self) -> None:
        if self._auth:
            self._auth.close()


class _BearerCredentials(Credentials):
    def __init__(self, provider: Callable[[], str]) -> None:
        self._provider = provider

    def _build_rest_client(self, config: Any) -> Any:
        from taurus_protect.crypto.bearer_rest import BearerRESTClient

        return BearerRESTClient(config, self._provider)
