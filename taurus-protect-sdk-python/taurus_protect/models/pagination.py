"""Pagination values and the builders every list method goes through.

Offset lists return a :class:`Pagination`, cursor lists a :class:`CursorPage`. Both are
built here and nowhere else, so a list cannot drift from the shared contract (repo-root
``CLAUDE.md`` -> "Pagination (cross-SDK)").

validatord omits zero values, so an absent ``totalItems``/``offset``/``hasNext``/``next``
means 0/false/none and ``{}`` is a valid empty page.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from typing import Any, Dict, Optional, Union

from pydantic import BaseModel, Field

from taurus_protect.errors import APIError

DEFAULT_PAGE_SIZE = 20
MAX_PAGE_SIZE = 100

# Price history cannot page: the whole series comes back in one reply.
PRICE_HISTORY_MAX_LIMIT = 365

# The largest count every SDK represents exactly (TypeScript numbers are doubles).
MAX_COUNT = 2**53 - 1

# How each offset endpoint computes the NEXT page's offset.
REPLY_OFFSET = "reply_offset"
PLUS_ROWS = "plus_rows"
PLUS_MIN_ROWS_LIMIT = "plus_min_rows_limit"
PLUS_SERVER_ROWS = "plus_server_rows"
PLUS_LIMIT = "plus_limit"

OFFSET_RULES = (REPLY_OFFSET, PLUS_ROWS, PLUS_MIN_ROWS_LIMIT, PLUS_SERVER_ROWS, PLUS_LIMIT)

# Operations that cannot page: one reply, bounded by a limit.
LIMIT_ONLY = "limit_only"

_PAGE_REQUEST_NEXT = "NEXT"
_CANONICAL_COUNT = re.compile(r"(?:0|[1-9][0-9]*)\Z")


class Pagination(BaseModel):
    """
    The page window of an offset list. Never None on success.

    Continue with ``offset=next_offset`` while ``has_more`` is true.

    Attributes:
        limit: The page size that was sent.
        offset: The offset that was sent (never the reply's).
        total_items: The server's total, reduced by rows the SDK withheld.
        next_offset: The offset of the next page, in the server's numbering.
        has_more: Whether another page exists.
    """

    limit: int = Field(default=0, description="Page size that was sent")
    offset: int = Field(default=0, description="Offset that was sent")
    total_items: int = Field(default=0, description="Server total minus withheld rows")
    next_offset: int = Field(default=0, description="Offset of the next page")
    has_more: bool = Field(default=False, description="Whether another page exists")

    model_config = {"frozen": True}


class CursorPage(BaseModel):
    """
    The page window of a cursor list. Never None on success.

    Continue with ``cursor=next_cursor`` while ``has_more`` is true.

    Attributes:
        page_size: The page size that was sent.
        next_cursor: Opaque cursor of the next page; empty when ``has_more`` is false.
        has_more: Whether another page exists.
        total_items: The server's total, only on endpoints that report one.
    """

    page_size: int = Field(default=0, description="Page size that was sent")
    next_cursor: str = Field(default="", description="Cursor of the next page, or empty")
    has_more: bool = Field(default=False, description="Whether another page exists")
    total_items: Optional[int] = Field(
        default=None, description="Server total, only where the endpoint reports one"
    )

    model_config = {"frozen": True}


@dataclass(frozen=True)
class CursorRequest:
    """The cursor fields one page request sends. Build it with :func:`cursor_request`."""

    page_size: int
    current_page: Optional[str] = None
    page_request: Optional[str] = None

    def query_params(self, prefix: str = "cursor") -> Dict[str, Optional[str]]:
        """Keyword arguments for a generated operation taking ``<prefix>_*`` parameters."""
        return {
            f"{prefix}_current_page": self.current_page,
            f"{prefix}_page_request": self.page_request,
            f"{prefix}_page_size": str(self.page_size),
        }


class CursorListOptions(BaseModel):
    """
    The page window of a cursor list's options.

    Pass ``cursor`` (a previous ``next_cursor``) to continue a walk, or the low-level
    ``current_page``/``page_request`` pair; combining them is refused.
    """

    page_size: Optional[int] = Field(default=None, description="Page size (default 20, max 100)")
    cursor: Optional[str] = Field(default=None, description="A previous next_cursor")
    current_page: Optional[str] = Field(default=None, description="Low-level page token")
    page_request: Optional[str] = Field(
        default=None, description="Low-level page direction (FIRST, PREVIOUS, NEXT, LAST)"
    )

    def to_cursor_request(self) -> "CursorRequest":
        """Resolve these options into the cursor fields one request sends."""
        return cursor_request(
            self.page_size,
            self.cursor,
            current_page=self.current_page,
            page_request=self.page_request,
        )


def resolve_page_size(
    value: Optional[int],
    name: str = "page_size",
    maximum: Optional[int] = MAX_PAGE_SIZE,
) -> int:
    """
    Resolve a page size to the value sent: unset or 0 means the default.

    Args:
        value: The caller's page size.
        name: The option name, used in the error.
        maximum: The largest accepted value, or None for no SDK maximum.

    Returns:
        The page size to send.

    Raises:
        ValueError: If the value is negative, above ``maximum``, or not an integer.
    """
    if value is None:
        return DEFAULT_PAGE_SIZE
    if isinstance(value, bool) or not isinstance(value, int):
        raise ValueError(f"{name} must be an integer, got {value!r}")
    if value == 0:
        return DEFAULT_PAGE_SIZE
    if value < 0:
        if maximum is None:
            raise ValueError(f"{name} cannot be negative, got {value}")
        raise ValueError(f"{name} must be between 1 and {maximum}, got {value}")
    if maximum is not None and value > maximum:
        raise ValueError(f"{name} must be at most {maximum}, got {value}")
    return value


def resolve_offset(value: Optional[int], name: str = "offset") -> int:
    """
    Resolve an offset: unset means 0.

    Raises:
        ValueError: If the value is negative or not an integer.
    """
    if value is None:
        return 0
    if isinstance(value, bool) or not isinstance(value, int):
        raise ValueError(f"{name} must be an integer, got {value!r}")
    if value < 0:
        raise ValueError(f"{name} cannot be negative, got {value}")
    return value


def offset_query(limit: int, offset: int) -> Dict[str, Optional[str]]:
    """The ``limit``/``offset`` query arguments; offset 0 is the server default and is omitted."""
    return {"limit": str(limit), "offset": str(offset) if offset else None}


def parse_count(raw: Optional[Union[str, int]], field: str) -> int:
    """
    Parse a count the server sent; absent means 0.

    A count that is not a canonical decimal in ``[0, 2**53 - 1]`` is refused rather than
    read as 0: a silent 0 would end every walk early.

    Raises:
        APIError: If the count is malformed or out of range.
    """
    if raw is None:
        return 0
    if isinstance(raw, bool):
        raise _invalid_reply(f"{field} is not a count: {raw!r}")
    if isinstance(raw, int):
        value = raw
    elif isinstance(raw, str) and _CANONICAL_COUNT.match(raw):
        value = int(raw)
    else:
        raise _invalid_reply(f"{field} is not a canonical count: {raw!r}")
    if value < 0 or value > MAX_COUNT:
        raise _invalid_reply(f"{field} is out of range: {raw!r}")
    return value


def offset_pagination(
    rule: str,
    *,
    limit: int,
    offset: int,
    served_rows: int,
    total_items: Optional[Union[str, int]],
    reply_offset: Optional[Union[str, int]] = None,
    excluded: int = 0,
) -> Pagination:
    """
    Build the :class:`Pagination` of one offset-list reply.

    Args:
        rule: How the endpoint computes the next offset (one of ``OFFSET_RULES``).
        limit: The page size that was sent.
        offset: The offset that was sent.
        served_rows: Rows the SERVER returned, before any SDK exclusion.
        total_items: The reply's ``totalItems``, as received.
        reply_offset: The reply's ``offset`` (``reply_offset`` rule only), as received.
        excluded: Rows the SDK withheld from the caller.

    Returns:
        The page window. ``next_offset`` and ``has_more`` use the server's numbers;
        only ``total_items`` is reduced by ``excluded``.

    Raises:
        APIError: If a count in the reply is malformed.
        ValueError: If ``rule`` is unknown.
    """
    total = parse_count(total_items, "totalItems")
    if rule == REPLY_OFFSET:
        # validatord's reply offset is already offset + rows: it is the NEXT page.
        if reply_offset is None:
            next_offset = offset + served_rows
        else:
            next_offset = parse_count(reply_offset, "offset")
    elif rule in (PLUS_ROWS, PLUS_SERVER_ROWS):
        next_offset = offset + served_rows
    elif rule == PLUS_MIN_ROWS_LIMIT:
        # A synthetic daemon user / tech group can be appended beyond the limit.
        next_offset = offset + min(served_rows, limit)
    elif rule == PLUS_LIMIT:
        # Skipped rows keep their SQL slot, so a short page is not the end.
        next_offset = offset + limit
    else:
        raise ValueError(f"unknown pagination rule {rule!r}")

    return Pagination(
        limit=limit,
        offset=offset,
        total_items=max(0, total - excluded),
        # No progress ends a walk; some totals are upper bounds.
        next_offset=next_offset,
        has_more=offset < next_offset < total,
    )


def cursor_request(
    page_size: Optional[int] = None,
    cursor: Optional[str] = None,
    *,
    current_page: Optional[str] = None,
    page_request: Optional[str] = None,
    name: str = "page_size",
) -> CursorRequest:
    """
    Resolve the cursor fields of one cursor-list request.

    ``cursor`` is a previous :attr:`CursorPage.next_cursor` and sends
    ``currentPage=<cursor>`` + ``pageRequest=NEXT``. ``current_page``/``page_request``
    are the low-level alternative; combining them with ``cursor`` is refused.

    Raises:
        ValueError: If the page size is invalid, or ``cursor`` is combined with the
            low-level options.
    """
    size = resolve_page_size(page_size, name)
    if cursor is not None and not isinstance(cursor, str):
        raise ValueError(f"cursor must be a string, got {type(cursor).__name__}")
    if cursor:
        if current_page or page_request:
            raise ValueError(
                "cursor cannot be combined with current_page or page_request: "
                "pass a previous next_cursor, or the low-level pair, not both"
            )
        return CursorRequest(page_size=size, current_page=cursor, page_request=_PAGE_REQUEST_NEXT)
    return CursorRequest(
        page_size=size,
        current_page=current_page or None,
        page_request=page_request or None,
    )


def cursor_page(
    page_size: int,
    reply_cursor: Any = None,
    *,
    token: Optional[Union[str, bytes]] = None,
    total: Optional[Union[str, int]] = None,
    has_total: bool = False,
    excluded: int = 0,
) -> CursorPage:
    """
    Build the :class:`CursorPage` of one cursor-list reply.

    Cursor lists pass the reply's generated ``cursor`` (``currentPage``/``hasNext``);
    the two token-only operations (wallet tokens ``next``, rules history ``cursor``)
    pass ``token`` instead.

    Args:
        page_size: The page size that was sent.
        reply_cursor: The reply's ``cursor`` model, or None when the reply omitted it.
        token: The reply's continuation token (token-only operations).
        total: The reply's total, as received (only where ``has_total``).
        has_total: Whether this endpoint reports a total.
        excluded: Rows the SDK withheld, subtracted from the total.

    Returns:
        The page window; ``next_cursor`` is the reply's ``currentPage`` only when
        ``hasNext``. A cursor derived any other way re-reads page 1 or is a 400.

    Raises:
        APIError: If the reply's total is malformed, or ``hasNext`` has no cursor.
    """
    total_items: Optional[int] = None
    if has_total:
        total_items = max(0, parse_count(total, "total") - excluded)

    if token is not None and reply_cursor is not None:
        raise ValueError("pass the reply cursor or the reply token, not both")

    if token is not None:
        text = token.decode("utf-8") if isinstance(token, bytes) else token
        return CursorPage(
            page_size=page_size,
            next_cursor=text,
            has_more=bool(text),
            total_items=total_items,
        )

    has_more = bool(reply_cursor is not None and reply_cursor.has_next)
    next_cursor = ""
    if has_more:
        next_cursor = reply_cursor.current_page or ""
        if not next_cursor:
            raise _invalid_reply("cursor.hasNext is true but cursor.currentPage is empty")
    return CursorPage(
        page_size=page_size,
        next_cursor=next_cursor,
        has_more=has_more,
        total_items=total_items,
    )


def _invalid_reply(detail: str) -> APIError:
    return APIError(f"invalid list reply: {detail}", description="Invalid Reply")
