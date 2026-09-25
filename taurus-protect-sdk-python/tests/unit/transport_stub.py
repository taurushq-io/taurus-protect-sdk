"""A transport stub for service tests.

Only ``RESTClientObject.request`` is replaced, so the real generated operation runs: its
signature validation, query and body serialization, and reply deserializer. A bare
``MagicMock`` API accepts any keyword and answers any attribute, which is how list
methods came to pass invalid arguments and read reply fields that do not exist.
"""

from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Any, List, Optional, Tuple, Union
from unittest.mock import patch
from urllib.parse import parse_qsl, urlsplit

import urllib3

from taurus_protect._internal.openapi import ApiClient, Configuration
from taurus_protect._internal.openapi.rest import RESTClientObject, RESTResponse

HOST = "https://api.example.test"


def api_client() -> ApiClient:
    """A generated ApiClient configured as ProtectClient configures its own."""
    config = Configuration(host=HOST)
    config.datetime_format = "%Y-%m-%dT%H:%M:%S.%fZ"
    return ApiClient(configuration=config)


@dataclass
class Reply:
    """One canned HTTP reply."""

    body: Union[dict, list, str] = "{}"
    status: int = 200

    def encoded(self) -> bytes:
        if isinstance(self.body, (dict, list)):
            return json.dumps(self.body).encode("utf-8")
        return self.body.encode("utf-8")


@dataclass
class RecordedRequest:
    """One request the generated client sent."""

    method: str
    url: str
    body: Any

    @property
    def path(self) -> str:
        return urlsplit(self.url).path

    @property
    def query(self) -> List[Tuple[str, str]]:
        """Decoded query pairs, sorted: compare as an exact multiset."""
        return sorted(parse_qsl(urlsplit(self.url).query, keep_blank_values=True))

    @property
    def raw_query(self) -> str:
        return urlsplit(self.url).query

    def param(self, name: str) -> Optional[str]:
        values = [v for k, v in self.query if k == name]
        assert len(values) <= 1, f"{name} sent {len(values)} times"
        return values[0] if values else None


class StubTransport:
    """
    Patch the transport for the duration of a ``with`` block.

    With no replies every request is answered ``{}``. With replies they are served in
    order and one request too many fails the test, so a walk that keeps going after
    ``has_more`` turned false is caught.
    """

    def __init__(self, *replies: Union[Reply, dict, list, str]) -> None:
        self._replies = [r if isinstance(r, Reply) else Reply(r) for r in replies]
        self._strict = bool(replies)
        self.requests: List[RecordedRequest] = []
        self._patch = patch.object(RESTClientObject, "request", self._request)

    def __enter__(self) -> "StubTransport":
        self._patch.start()
        return self

    def __exit__(self, *exc: Any) -> None:
        self._patch.stop()

    @property
    def last(self) -> RecordedRequest:
        assert self.requests, "no request was sent"
        return self.requests[-1]

    def _request(
        self,
        method: str,
        url: str,
        headers: Any = None,
        body: Any = None,
        post_params: Any = None,
        _request_timeout: Any = None,
    ) -> RESTResponse:
        self.requests.append(RecordedRequest(method=method, url=url, body=body))
        index = len(self.requests) - 1
        if self._strict:
            assert index < len(self._replies), f"unexpected request #{index + 1}: {method} {url}"
            reply = self._replies[index]
        else:
            reply = Reply()
        raw = urllib3.HTTPResponse(
            body=reply.encoded(),
            status=reply.status,
            headers={"content-type": "application/json"},
            preload_content=True,
        )
        return RESTResponse(raw)
