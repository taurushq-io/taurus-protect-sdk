"""Bearer-token REST client (no TPV1 signing)."""

from __future__ import annotations

from typing import Any, Callable, Optional

from taurus_protect._internal.openapi.rest import RESTClientObject
from taurus_protect.errors import ConfigurationError


class BearerRESTClient(RESTClientObject):
    """
    REST client that authenticates each request with a bearer token from a provider.

    The provider is called once per request, so a single client can serve rotating
    or per-caller tokens (e.g. a provider reading a contextvar). Unlike
    AuthenticatedRESTClient it neither reads nor signs the body. Mirrors the Go
    SDK's per-request bearer provider.
    """

    def __init__(self, configuration: Any, token_provider: Callable[[], str]) -> None:
        """
        Initialize the bearer REST client.

        Args:
            configuration: OpenAPI Configuration object.
            token_provider: Callable returning the bearer token for a request.
        """
        super().__init__(configuration)
        self._token_provider = token_provider

    def request(
        self,
        method: str,
        url: str,
        headers: Optional[dict] = None,
        body: Any = None,
        post_params: Any = None,
        _request_timeout: Any = None,
    ) -> Any:
        """Perform an HTTP request with a Bearer Authorization header."""
        if headers is None:
            headers = {}
        # A refresh that silently yields nothing would send "Bearer None" and surface
        # as an opaque 401 rather than the real cause.
        token = self._token_provider()
        if not token:
            raise ConfigurationError("bearer token provider returned an empty token")
        headers["Authorization"] = f"Bearer {token}"
        return super().request(
            method=method,
            url=url,
            headers=headers,
            body=body,
            post_params=post_params,
            _request_timeout=_request_timeout,
        )
