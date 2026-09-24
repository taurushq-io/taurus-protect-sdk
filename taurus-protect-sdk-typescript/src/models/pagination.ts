/**
 * Pagination for list results.
 *
 * Implements the cross-SDK contract in the repo-root CLAUDE.md § "Pagination (cross-SDK)",
 * pinned in every SDK by `scripts/resources/pagination-vectors.json`. Every list method
 * builds its paging request and its paging result through the functions in this file, so
 * page sizes are validated, cursors forwarded and next offsets computed in one place.
 *
 * validatord omits zero values, so an absent `totalItems` / `offset` / `hasNext` / `next`
 * means 0 / false / none, and `{}` is a valid empty page.
 */

import { PaginationError, ValidationError } from "../errors";

/** Page size sent when the caller sets none, or 0. */
export const DEFAULT_PAGE_SIZE = 20;

/** Largest page size a list method accepts. */
export const MAX_PAGE_SIZE = 100;

/** Largest `limit` the price history accepts: a year of daily points. */
export const MAX_PRICE_HISTORY_LIMIT = 365;

/** Largest count a reply may carry: 2^53 - 1, the largest integer a number holds exactly. */
export const MAX_COUNT = Number.MAX_SAFE_INTEGER;

/**
 * Offset pagination of a list result. Present on every successful offset list, the empty
 * page included.
 */
export interface Pagination {
  /** The page size that was sent. */
  readonly limit: number;
  /** The offset that was sent — the request's, never the reply's. */
  readonly offset: number;
  /** The server's total, reduced by rows this SDK withheld from the page (never below 0). */
  readonly totalItems: number;
  /** Offset of the next page: pass it back as `offset` while `hasMore` is true. */
  readonly nextOffset: number;
  /** Whether another page exists. */
  readonly hasMore: boolean;
}

/**
 * Cursor pagination of a list result. Present on every successful cursor list, the empty
 * page included.
 */
export interface CursorPage {
  /** The page size that was sent. */
  readonly pageSize: number;
  /** Cursor of the next page: pass it back as `cursor` while `hasMore` is true; "" otherwise. */
  readonly nextCursor: string;
  /** Whether another page exists. */
  readonly hasMore: boolean;
  /** The server's total where the endpoint returns one (0 when omitted); undefined otherwise. */
  readonly totalItems?: number;
}

/** A page of an offset list. */
export interface PaginatedResult<T> {
  /** The rows of this page. */
  readonly items: T[];
  /** Where this page sits and how to fetch the next one. */
  readonly pagination: Pagination;
}

/** Direction of the low-level cursor navigation options. */
export type PageRequest = "FIRST" | "PREVIOUS" | "NEXT" | "LAST";

/** Paging options of an offset list. */
export interface OffsetPageOptions {
  /** Page size, 1 to 100 (default 20). */
  readonly limit?: number;
  /** Rows to skip (default 0): a previous result's `pagination.nextOffset`. */
  readonly offset?: number;
}

/** Paging options of a cursor list. */
export interface CursorPageOptions {
  /** Page size, 1 to 100 (default 20). */
  readonly pageSize?: number;
  /** A previous result's `pagination.nextCursor`: fetches the page after it. */
  readonly cursor?: string;
}

/** Paging options of a cursor list that also offers low-level cursor navigation. */
export interface CursorNavigationOptions extends CursorPageOptions {
  /** Low-level: the page token to navigate from. Cannot be combined with `cursor`. */
  readonly currentPage?: string;
  /** Low-level: the direction to navigate from `currentPage`. Cannot be combined with `cursor`. */
  readonly pageRequest?: PageRequest;
}

/**
 * Which bound a page size is checked against: a pageable list (max 100), the price history
 * (max 365), or an export that cannot page (no maximum — a larger limit is the only way to
 * reach more rows).
 */
export type PageSizeKind = "page" | "price_history" | "export";

/**
 * How an offset endpoint's next page is found (the contract's rule table).
 *
 * - `reply_offset`: the reply's `offset`, else offset + rows (wallets, addresses).
 * - `plus_rows`: offset + rows (transactions, fee payers, actions).
 * - `plus_min_rows_limit`: offset + min(rows, limit) — the server may append a synthetic
 *   daemon user / technical group beyond `limit` (users, groups).
 * - `plus_server_rows`: offset + the rows the SERVER returned, before this SDK withheld
 *   any (whitelisted addresses).
 * - `plus_limit`: offset + limit — skipped rows keep their SQL slot (whitelisted contracts).
 */
export type OffsetRule =
  | "reply_offset"
  | "plus_rows"
  | "plus_min_rows_limit"
  | "plus_server_rows"
  | "plus_limit";

/** The paging request actually sent to an offset endpoint. */
export interface OffsetRequest {
  readonly limit: number;
  readonly offset: number;
}

/** An offset reply's paging fields, exactly as the generated reply types carry them. */
export interface OffsetReply {
  readonly totalItems?: string | number;
  readonly offset?: string | number;
}

/** The paging request actually sent to a cursor endpoint. */
export interface CursorRequest {
  /** The page size that is sent. */
  readonly pageSize: number;
  /**
   * The cursor in the shape of the generated `TgvalidatordRequestCursor`, which is the
   * body form; the `cursor.*` and `requestCursor.*` query forms are spelled from it.
   */
  readonly cursor: {
    readonly currentPage?: string;
    readonly pageRequest?: string;
    readonly pageSize: string;
  };
}

/** The paging request actually sent to a token endpoint: `limit` plus the opaque token. */
export interface TokenRequest {
  /** The page size that is sent. */
  readonly pageSize: number;
  /** The page size as the `limit` query string. */
  readonly limit: string;
  /** The token to continue from, when the caller passed one. */
  readonly cursor?: string;
}

/** A reply cursor's fields, exactly as the generated `TgvalidatordResponseCursor` carries them. */
export interface ReplyCursor {
  readonly currentPage?: string;
  readonly hasNext?: boolean;
}

/** The total of a cursor or token page, for the endpoints that return one. */
export interface PageTotal {
  /** The reply's total as the generated reply type carries it; absent means 0. */
  readonly total: string | number | undefined;
  /** Rows this SDK withheld from the page; they reduce `totalItems` (never below 0). */
  readonly excluded?: number;
}

/**
 * Resolves a requested page size: unset or 0 means the default, anything else must be a
 * whole number in bounds. The one page-size validator every list method goes through.
 *
 * @param value - The page size the caller passed
 * @param option - The option's name, for the error message
 * @param kind - `page` (max 100), `price_history` (max 365) or `export` (no maximum)
 * @returns The page size to send
 * @throws {@link ValidationError} If the value is negative, fractional or above the maximum
 */
export function resolvePageSize(
  value: number | null | undefined,
  option = "pageSize",
  kind: PageSizeKind = "page"
): number {
  if (value === undefined || value === null || value === 0) {
    return DEFAULT_PAGE_SIZE;
  }
  if (typeof value !== "number" || !Number.isSafeInteger(value)) {
    throw new ValidationError(`${option} must be a whole number, got ${String(value)}`);
  }
  if (value < 0) {
    throw new ValidationError(`${option} must not be negative, got ${value}`);
  }
  const max =
    kind === "page" ? MAX_PAGE_SIZE : kind === "price_history" ? MAX_PRICE_HISTORY_LIMIT : undefined;
  if (max !== undefined && value > max) {
    throw new ValidationError(`${option} must be at most ${max}, got ${value}`);
  }
  return value;
}

/**
 * Resolves a requested offset: unset means 0, anything else must be a whole number >= 0.
 *
 * @param value - The offset the caller passed
 * @param option - The option's name, for the error message
 * @returns The offset to send
 * @throws {@link ValidationError} If the value is negative or fractional
 */
export function resolveOffset(value: number | null | undefined, option = "offset"): number {
  if (value === undefined || value === null) {
    return 0;
  }
  if (typeof value !== "number" || !Number.isSafeInteger(value)) {
    throw new ValidationError(`${option} must be a whole number, got ${String(value)}`);
  }
  if (value < 0) {
    throw new ValidationError(`${option} must not be negative, got ${value}`);
  }
  return value;
}

/**
 * Resolves the paging request of an offset list.
 *
 * @param options - The caller's `limit` / `offset`
 * @returns The limit and offset to send
 * @throws {@link ValidationError} If either is out of bounds
 */
export function offsetRequest(options: OffsetPageOptions | undefined): OffsetRequest {
  return {
    limit: resolvePageSize(options?.limit, "limit"),
    offset: resolveOffset(options?.offset, "offset"),
  };
}

/**
 * Resolves the paging request of a cursor list.
 *
 * With `cursor` set the request continues after that page: `currentPage=<cursor>` and
 * `pageRequest=NEXT`. Without it the low-level `currentPage` / `pageRequest` pass through
 * unchanged. The page size is always sent.
 *
 * @param options - The caller's paging options
 * @returns The page size and the cursor to send
 * @throws {@link ValidationError} If the page size is out of bounds, or `cursor` is combined
 *   with `currentPage` or `pageRequest`
 */
export function cursorRequest(options: CursorNavigationOptions | undefined): CursorRequest {
  const pageSize = resolvePageSize(options?.pageSize, "pageSize");
  const cursor = continuationToken(options?.cursor);
  if (cursor !== undefined) {
    if (options?.currentPage !== undefined) {
      throw new ValidationError(
        "cursor and currentPage cannot both be set: cursor already names the page to continue from"
      );
    }
    if (options?.pageRequest !== undefined) {
      throw new ValidationError(
        "cursor and pageRequest cannot both be set: a cursor always continues with the NEXT page"
      );
    }
    return {
      pageSize,
      cursor: { currentPage: cursor, pageRequest: "NEXT", pageSize: String(pageSize) },
    };
  }
  return {
    pageSize,
    cursor: {
      currentPage: options?.currentPage,
      pageRequest: options?.pageRequest,
      pageSize: String(pageSize),
    },
  };
}

/**
 * Resolves the paging request of a token list (`limit` plus an opaque `cursor` token).
 *
 * @param options - The caller's page size and token
 * @returns The page size, limit and token to send
 * @throws {@link ValidationError} If the page size is out of bounds
 */
export function tokenRequest(options: CursorPageOptions | undefined): TokenRequest {
  const pageSize = resolvePageSize(options?.pageSize, "pageSize");
  return { pageSize, limit: String(pageSize), cursor: continuationToken(options?.cursor) };
}

/**
 * Builds the pagination of an offset list page — the one place a next offset is computed.
 *
 * `offset` is the request's. `nextOffset` follows the endpoint's rule, and `hasMore` is
 * `nextOffset > offset && nextOffset < total`: no progress ends a walk, and some totals are
 * upper bounds. `totalItems` is the server total less the rows withheld, never below 0;
 * `nextOffset` and `hasMore` always use the server's numbers.
 *
 * @param rule - The endpoint's next-offset rule
 * @param request - The limit and offset that were sent
 * @param reply - The reply's `totalItems` / `offset`, as the generated reply carries them
 * @param servedRows - How many rows the server returned, before this SDK withheld any
 * @param sdkExcluded - How many of those rows this SDK withheld
 * @returns The page's pagination
 * @throws {@link PaginationError} If the reply's total or offset is not a canonical count
 */
export function buildOffsetPagination(
  rule: OffsetRule,
  request: OffsetRequest,
  reply: OffsetReply,
  servedRows: number,
  sdkExcluded = 0
): Pagination {
  const total = parseReplyCount(reply.totalItems, "totalItems");
  const { limit, offset } = request;
  let nextOffset: number;
  switch (rule) {
    case "reply_offset":
      nextOffset =
        reply.offset === undefined || reply.offset === null
          ? offset + servedRows
          : parseReplyCount(reply.offset, "offset");
      break;
    case "plus_rows":
    case "plus_server_rows":
      nextOffset = offset + servedRows;
      break;
    case "plus_min_rows_limit":
      nextOffset = offset + Math.min(servedRows, limit);
      break;
    case "plus_limit":
      nextOffset = offset + limit;
      break;
    default:
      throw new ValidationError(`unknown offset rule: ${describe(rule)}`);
  }
  return {
    limit,
    offset,
    totalItems: Math.max(0, total - sdkExcluded),
    nextOffset,
    hasMore: nextOffset > offset && nextOffset < total,
  };
}

/**
 * Builds the pagination of a cursor or token list page — the one place a next cursor is
 * derived.
 *
 * `next` is the reply's cursor object for a cursor list (`hasMore` = `hasNext`, and the
 * next cursor is its `currentPage`), or the reply's token for a token list (`hasMore` =
 * the token is present, and the next cursor is the token). Either way the next cursor is
 * "" when there is no next page, and is never derived any other way: sending NEXT past the
 * end is a 400 on the v2 lists and silently returns page 1 on the keyset lists.
 *
 * @param pageSize - The page size that was sent
 * @param next - The reply's cursor object, or its token
 * @param total - The reply's total, for the endpoints that return one
 * @returns The page's pagination
 * @throws {@link PaginationError} If the cursor reports a next page without naming it, or
 *   the total is not a canonical count
 */
export function buildCursorPage(
  pageSize: number,
  next: ReplyCursor | string | undefined | null,
  total?: PageTotal
): CursorPage {
  let nextCursor = "";
  let hasMore = false;
  if (typeof next === "string") {
    hasMore = next.length > 0;
    nextCursor = next;
  } else if (next !== undefined && next !== null) {
    hasMore = next.hasNext === true;
    if (hasMore) {
      if (typeof next.currentPage !== "string" || next.currentPage.length === 0) {
        throw new PaginationError(
          "reply cursor reports a next page (hasNext) but carries no currentPage to continue from"
        );
      }
      nextCursor = next.currentPage;
    }
  }
  if (total === undefined) {
    return { pageSize, nextCursor, hasMore };
  }
  const count = parseReplyCount(total.total, "total");
  return {
    pageSize,
    nextCursor,
    hasMore,
    totalItems: Math.max(0, count - (total.excluded ?? 0)),
  };
}

/**
 * Parses a count off a reply: a canonical decimal in [0, 2^53 - 1]. Absent means 0.
 *
 * Anything else raises rather than reading as 0, which would end a walk early and pass a
 * truncated list off as complete.
 *
 * @param value - The count as the generated reply type carries it
 * @param field - The reply field's name, for the error message
 * @returns The count
 * @throws {@link PaginationError} If the value is not a canonical count
 */
export function parseReplyCount(value: string | number | undefined | null, field: string): number {
  if (value === undefined || value === null) {
    return 0;
  }
  if (typeof value === "number") {
    if (Number.isSafeInteger(value) && value >= 0) {
      return value;
    }
    throw new PaginationError(`reply ${field} is not a count in [0, 2^53 - 1]: ${value}`);
  }
  if (typeof value !== "string" || !/^(0|[1-9][0-9]*)$/.test(value)) {
    throw new PaginationError(`reply ${field} is not a canonical count: ${describe(value)}`);
  }
  const count = Number(value);
  if (!Number.isSafeInteger(count)) {
    throw new PaginationError(`reply ${field} ${describe(value)} is above 2^53 - 1`);
  }
  return count;
}

/** Normalises a caller's continuation token: "" and unset both mean "start at page 1". */
function continuationToken(cursor: string | undefined): string | undefined {
  if (cursor === undefined || cursor === null) {
    return undefined;
  }
  if (typeof cursor !== "string") {
    throw new ValidationError("cursor must be a string: a previous result's nextCursor");
  }
  return cursor.length === 0 ? undefined : cursor;
}

/** Quotes a server value for an error message, bounded so a large reply cannot flood it. */
function describe(value: unknown): string {
  const text = JSON.stringify(value) ?? String(value);
  return text.length > 40 ? `${text.slice(0, 37)}...` : text;
}
