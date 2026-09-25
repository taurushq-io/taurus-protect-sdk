"""The pagination builders, at the edges the shared vectors do not reach."""

from __future__ import annotations

import pytest

from taurus_protect._internal.openapi.models import TgvalidatordResponseCursor
from taurus_protect.errors import APIError
from taurus_protect.models.pagination import (
    DEFAULT_PAGE_SIZE,
    MAX_COUNT,
    MAX_PAGE_SIZE,
    PLUS_ROWS,
    CursorListOptions,
    CursorPage,
    CursorRequest,
    Pagination,
    cursor_page,
    cursor_request,
    offset_pagination,
    offset_query,
    parse_count,
    resolve_offset,
    resolve_page_size,
)


class TestResolvePageSize:
    def test_constants(self) -> None:
        assert (DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE) == (20, 100)

    def test_error_names_the_option_and_the_maximum(self) -> None:
        with pytest.raises(ValueError, match=r"limit must be at most 100, got 101"):
            resolve_page_size(101, "limit")
        with pytest.raises(ValueError, match=r"page_size must be between 1 and 100, got -1"):
            resolve_page_size(-1)

    @pytest.mark.parametrize("value", [True, False, "20", 2.0])
    def test_non_integers_are_refused(self, value: object) -> None:
        with pytest.raises(ValueError, match="must be an integer"):
            resolve_page_size(value)  # type: ignore[arg-type]

    def test_no_maximum(self) -> None:
        assert resolve_page_size(5000, "limit", maximum=None) == 5000
        with pytest.raises(ValueError, match="limit cannot be negative"):
            resolve_page_size(-5, "limit", maximum=None)


class TestResolveOffset:
    def test_unset_and_valid(self) -> None:
        assert resolve_offset(None) == 0
        assert resolve_offset(7) == 7

    @pytest.mark.parametrize("value", [True, "3"])
    def test_non_integers_are_refused(self, value: object) -> None:
        with pytest.raises(ValueError, match="offset must be an integer"):
            resolve_offset(value)  # type: ignore[arg-type]


def test_offset_zero_is_not_sent() -> None:
    assert offset_query(20, 0) == {"limit": "20", "offset": None}
    assert offset_query(20, 40) == {"limit": "20", "offset": "40"}


class TestParseCount:
    def test_ints_and_bounds(self) -> None:
        assert parse_count(None, "totalItems") == 0
        assert parse_count(MAX_COUNT, "totalItems") == MAX_COUNT
        with pytest.raises(APIError, match="out of range"):
            parse_count(MAX_COUNT + 1, "totalItems")
        with pytest.raises(APIError, match="out of range"):
            parse_count(-1, "totalItems")

    @pytest.mark.parametrize("raw", [True, "+1", " 1", "1.0", "١"])
    def test_non_canonical_counts_are_refused(self, raw: object) -> None:
        with pytest.raises(APIError, match="totalItems"):
            parse_count(raw, "totalItems")  # type: ignore[arg-type]

    def test_the_error_is_not_retryable(self) -> None:
        with pytest.raises(APIError) as excinfo:
            parse_count("x", "totalItems")
        assert excinfo.value.is_retryable() is False


def test_unknown_offset_rule_is_refused() -> None:
    with pytest.raises(ValueError, match="unknown pagination rule"):
        offset_pagination("plus_whatever", limit=1, offset=0, served_rows=0, total_items=None)


def test_exclusions_never_move_the_next_offset() -> None:
    page = offset_pagination(
        PLUS_ROWS, limit=10, offset=0, served_rows=10, total_items="12", excluded=3
    )
    assert page == Pagination(limit=10, offset=0, total_items=9, next_offset=10, has_more=True)


class TestCursorRequest:
    def test_cursor_sends_next(self) -> None:
        assert cursor_request(5, "c+/=") == CursorRequest(5, "c+/=", "NEXT")

    def test_empty_cursor_is_the_first_page(self) -> None:
        assert cursor_request(None, "") == CursorRequest(DEFAULT_PAGE_SIZE)

    @pytest.mark.parametrize("low_level", [{"current_page": "p"}, {"page_request": "PREVIOUS"}])
    def test_cursor_and_the_low_level_pair_are_exclusive(self, low_level: dict) -> None:
        with pytest.raises(ValueError, match="cursor cannot be combined"):
            cursor_request(None, "c", **low_level)

    def test_low_level_pair_passes_through(self) -> None:
        assert cursor_request(3, current_page="p", page_request="LAST") == CursorRequest(
            3, "p", "LAST"
        )

    def test_cursor_must_be_a_string(self) -> None:
        with pytest.raises(ValueError, match="cursor must be a string"):
            cursor_request(None, b"c")  # type: ignore[arg-type]

    def test_query_params_prefix(self) -> None:
        assert CursorRequest(7, "c", "NEXT").query_params("request_cursor") == {
            "request_cursor_current_page": "c",
            "request_cursor_page_request": "NEXT",
            "request_cursor_page_size": "7",
        }

    def test_options_resolve_through_the_same_builder(self) -> None:
        assert CursorListOptions(page_size=4, cursor="c").to_cursor_request() == CursorRequest(
            4, "c", "NEXT"
        )
        with pytest.raises(ValueError, match="page_size"):
            CursorListOptions(page_size=101).to_cursor_request()


class TestCursorPage:
    def test_bytes_token_is_decoded(self) -> None:
        page = cursor_page(20, token=b"YWJj", total="3", has_total=True)
        assert page == CursorPage(page_size=20, next_cursor="YWJj", has_more=True, total_items=3)

    def test_empty_token_ends_the_walk(self) -> None:
        assert cursor_page(20, token="").has_more is False

    def test_cursor_and_token_are_exclusive(self) -> None:
        with pytest.raises(ValueError, match="not both"):
            cursor_page(20, TgvalidatordResponseCursor(), token="t")

    def test_exclusions_reduce_the_total_and_clamp(self) -> None:
        page = cursor_page(20, token="t", total="2", has_total=True, excluded=5)
        assert page.total_items == 0
        assert page.has_more is True

    def test_no_total_unless_the_endpoint_reports_one(self) -> None:
        assert cursor_page(20, TgvalidatordResponseCursor(has_next=False)).total_items is None
