"""Cross-SDK list-request vectors: the same options must put the same request on the wire.

``scripts/resources/list-request-vectors.json`` names, per paged operation, the endpoint
and the exact decoded query (a multiset) or body a set of CANONICAL options must produce,
or that the options must be refused before any request. Each operation is mapped to the
Python method that wraps it, and each canonical option to that method's parameter, in
the explicit table in ``list_adapters``: an operation or option missing from it FAILS.
Operations this SDK does not wrap are listed, with the reason, in ``NOT_WRAPPED``.

Requests go through ``tests.unit.transport_stub``, so the generated client serializes
them; the reply is always ``{}``.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any, Dict

import pytest

from taurus_protect.models.pagination import CursorPage, Pagination
from tests.unit.services.list_adapters import ADAPTERS, NOT_WRAPPED
from tests.unit.transport_stub import StubTransport, api_client

VECTORS_PATH = (
    Path(__file__).resolve().parents[4] / "scripts" / "resources" / "list-request-vectors.json"
)

EXPECTED_COUNTS = {"methods": 53, "vectors": 278}


def _load() -> Dict[str, Any]:
    if not VECTORS_PATH.is_file():
        raise AssertionError(f"list-request vectors missing: {VECTORS_PATH}")
    with VECTORS_PATH.open(encoding="utf-8") as handle:
        return json.load(handle)


try:
    _VECTORS: Dict[str, Any] = _load()
except AssertionError:
    _VECTORS = {}


def _vectors() -> Dict[str, Any]:
    return _VECTORS or _load()


def _ids() -> list:
    entries = _VECTORS.get("vectors") or []
    ids = []
    for i in range(EXPECTED_COUNTS["vectors"]):
        if i < len(entries):
            ids.append(f"{i:03d}-{entries[i]['operation']}-{entries[i]['description']}")
        else:
            ids.append(f"{i:03d}-missing")
    return ids


def test_counts_match_the_file() -> None:
    data = _vectors()
    assert data["counts"] == EXPECTED_COUNTS
    assert len(data["methods"]) == EXPECTED_COUNTS["methods"]
    assert len(data["vectors"]) == EXPECTED_COUNTS["vectors"]


def test_every_operation_is_mapped_or_declared_unwrapped() -> None:
    methods = set(_vectors()["methods"])
    assert not set(ADAPTERS) & set(NOT_WRAPPED)
    assert methods - set(ADAPTERS) - set(NOT_WRAPPED) == set(), "unmapped operations"
    assert (set(ADAPTERS) | set(NOT_WRAPPED)) - methods == set(), "stale adapter entries"


@pytest.mark.parametrize("index", range(EXPECTED_COUNTS["vectors"]), ids=_ids())
def test_list_request_vector(index: int) -> None:
    data = _vectors()
    vector = data["vectors"][index]
    operation = vector["operation"]
    if operation in NOT_WRAPPED:
        pytest.skip(f"{operation}: {NOT_WRAPPED[operation]}")
    adapter = ADAPTERS.get(operation)
    assert adapter is not None, f"unmapped operation {operation}"

    unmapped = set(vector["options"]) - set(adapter.options)
    assert not unmapped, f"{operation}: unmapped options {sorted(unmapped)}"
    kwargs = {adapter.options[k]: v for k, v in vector["options"].items()}

    ac = api_client()
    service = adapter.build(ac)
    with StubTransport() as transport:
        if vector.get("expect_error"):
            with pytest.raises((ValueError, TypeError)) as excinfo:
                adapter.call(service, vector["path_params"], kwargs)
            if excinfo.type is TypeError:
                # Only a keyword the method does not take counts as refused by name.
                assert any(f"'{kw}'" in str(excinfo.value) for kw in kwargs), excinfo.value
            assert transport.requests == [], "refused options must send nothing"
            return
        result = adapter.call(service, vector["path_params"], kwargs)

    assert len(transport.requests) == 1
    request = transport.last
    method = data["methods"][operation]
    path = method["path"]
    for name, value in vector["path_params"].items():
        path = path.replace("{" + name + "}", value)
    assert request.method == method["method"]
    assert request.path == path

    expect = vector["expect"]
    if "query" in expect:
        assert request.query == sorted(tuple(pair) for pair in expect["query"])
        assert not request.body
    if "body" in expect:
        assert request.body == expect["body"]
        assert request.query == []

    _rows, page = result
    if adapter.page is None:
        assert page is None
    else:
        assert isinstance(page, adapter.page), type(page)
    if isinstance(page, CursorPage):
        has_total = _operation_has_total(operation)
        assert page == CursorPage(
            page_size=page.page_size, total_items=0 if has_total else None
        ), "an empty reply is an empty last page"
    if isinstance(page, Pagination):
        assert page.total_items == 0 and page.has_more is False


_PAGINATION_VECTORS = VECTORS_PATH.with_name("pagination-vectors.json")
_TOTALS: Dict[str, bool] = {}


def _operation_has_total(operation: str) -> bool:
    """Whether the operation reports a total, per the sibling pagination vectors."""
    if not _TOTALS:
        with _PAGINATION_VECTORS.open(encoding="utf-8") as handle:
            operations = json.load(handle)["operations"]
        _TOTALS.update({op: entry["total"] for op, entry in operations.items()})
    return _TOTALS[operation]
