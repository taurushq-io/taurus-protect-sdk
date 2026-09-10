"""Cross-SDK role-extraction alignment.

Every SDK must derive the same roles from the same server message, so the shared
vector file is consumed by the Go/Java/Python/TypeScript suites alike and a wording
change is caught in one place for all four.
"""

import json
from pathlib import Path

import pytest

from taurus_protect.errors import AuthorizationError, map_http_error, parse_required_roles

VECTORS_PATH = (
    Path(__file__).resolve().parents[3]
    / "scripts"
    / "resources"
    / "authorization-error-vectors.json"
)


def load_vectors() -> list:
    with VECTORS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


VECTORS = load_vectors()


def test_vectors_file_is_not_empty() -> None:
    assert VECTORS


@pytest.mark.parametrize("vector", VECTORS, ids=[v["description"] for v in VECTORS])
def test_vector_parses_to_expected_roles(vector: dict) -> None:
    assert parse_required_roles(vector["message"]) == vector["expected_roles"]


class TestAuthorizationErrorRequiredRoles:
    """Tests for role extraction on AuthorizationError."""

    def test_populated_from_message(self) -> None:
        err = AuthorizationError("one of the 'admin - adminreadonly' role is required")

        assert err.required_roles == ["admin", "adminreadonly"]

    def test_empty_for_non_role_denial(self) -> None:
        assert AuthorizationError("This endpoint has been disabled").required_roles == []

    def test_populated_through_map_http_error(self) -> None:
        err = map_http_error(403, "one of the 'tpuser' role is required")

        assert isinstance(err, AuthorizationError)
        assert err.required_roles == ["tpuser"]

    def test_none_message_is_tolerated(self) -> None:
        assert parse_required_roles(None) == []
