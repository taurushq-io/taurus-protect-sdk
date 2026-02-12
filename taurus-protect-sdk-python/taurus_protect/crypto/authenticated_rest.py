"""Authenticated REST client with TPV1 signing (follows Go SDK pattern)."""

from __future__ import annotations

import json
from typing import Any, Optional

from taurus_protect._internal.openapi.rest import RESTClientObject
from taurus_protect.crypto.tpv1 import TPV1Auth


class AuthenticatedRESTClient(RESTClientObject):
    """
    REST client that automatically signs requests with TPV1 authentication.

    Similar to Go SDK's TPV1Transport - wraps the base REST client and
    adds Authorization header to all outgoing requests.
    """

    def __init__(self, configuration: Any, tpv1_auth: TPV1Auth) -> None:
        """
        Initialize the authenticated REST client.

        Args:
            configuration: OpenAPI Configuration object.
            tpv1_auth: TPV1 authentication handler.
        """
        super().__init__(configuration)
        self._tpv1_auth = tpv1_auth

    def request(
        self,
        method: str,
        url: str,
        headers: Optional[dict] = None,
        body: Any = None,
        post_params: Any = None,
        _request_timeout: Any = None,
    ) -> Any:
        """
        Perform an HTTP request with TPV1 authentication.

        Signs the request with TPV1-HMAC-SHA256 before delegating to the
        base REST client.

        Args:
            method: HTTP method (GET, POST, etc.).
            url: Full request URL.
            headers: Request headers.
            body: Request body (dict for JSON, string for raw).
            post_params: Form parameters.
            _request_timeout: Request timeout.

        Returns:
            RESTResponse object.
        """
        # Ensure headers dict exists
        if headers is None:
            headers = {}

        # Parse URL into components for signing
        host, path, query = TPV1Auth.parse_url(url)

        # Get content type from headers
        content_type = headers.get("Content-Type")

        # Serialize body to string for signing
        # The base class will also serialize it, but we need the string for signing
        body_str: Optional[str] = None
        if body is not None:
            if isinstance(body, (str, bytes)):
                body_str = body if isinstance(body, str) else body.decode("utf-8")
            else:
                # JSON serialize for signing (same as base class does)
                body_str = json.dumps(body)

        # Sign the request (like Go's TPV1Auth.SignRequest)
        auth_header = self._tpv1_auth.sign_request(
            method=method,
            host=host,
            path=path,
            query=query,
            content_type=content_type,
            body=body_str,
        )

        # Add Authorization header
        headers["Authorization"] = auth_header

        # Delegate to base REST client
        return super().request(
            method=method,
            url=url,
            headers=headers,
            body=body,
            post_params=post_params,
            _request_timeout=_request_timeout,
        )
