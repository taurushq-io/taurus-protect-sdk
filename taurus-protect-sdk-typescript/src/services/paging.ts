/**
 * Wire glue between the pagination model and the generated clients.
 *
 * `models/pagination.ts` decides WHAT is sent and what a reply means; this module only
 * spells a resolved request in the generated parameter names. Not exported from the
 * package: the generated names are an implementation detail.
 */

import {
  MAX_PAGE_SIZE,
  buildCursorPage,
  cursorRequest,
  type CursorRequest,
  type OffsetRequest,
  type ReplyCursor,
} from "../models/pagination";

/** A cursor request as the generated `cursor.*` query parameters. */
export function cursorQuery(request: CursorRequest): {
  cursorCurrentPage?: string;
  cursorPageRequest?: string;
  cursorPageSize: string;
} {
  return {
    cursorCurrentPage: request.cursor.currentPage,
    cursorPageRequest: request.cursor.pageRequest,
    cursorPageSize: request.cursor.pageSize,
  };
}

/** A cursor request as the generated `requestCursor.*` query parameters. */
export function requestCursorQuery(request: CursorRequest): {
  requestCursorCurrentPage?: string;
  requestCursorPageRequest?: string;
  requestCursorPageSize: string;
} {
  return {
    requestCursorCurrentPage: request.cursor.currentPage,
    requestCursorPageRequest: request.cursor.pageRequest,
    requestCursorPageSize: request.cursor.pageSize,
  };
}

/** An offset request as the generated `limit` / `offset` strings; offset 0 is not sent. */
export function offsetQuery(request: OffsetRequest): { limit: string; offset?: string } {
  return {
    limit: String(request.limit),
    offset: request.offset > 0 ? String(request.offset) : undefined,
  };
}

/**
 * Walks a cursor list page by page until `match` finds a row, for the get-by-id helpers
 * whose endpoint has no single-row read.
 *
 * Each page is a normal, bounded page; the walk stops at the last page, or when the server
 * hands back a cursor it already gave (no progress), so a looping server cannot hold the
 * caller forever.
 *
 * @param fetchPage - Fetches one page for the given paging request
 * @param match - Selects the row being looked for
 * @returns The first matching row, or undefined when no page holds one
 */
export async function scanPages<R>(
  fetchPage: (request: CursorRequest) => Promise<{
    rows: R[] | undefined;
    cursor: ReplyCursor | undefined;
  }>,
  match: (row: R) => boolean
): Promise<R | undefined> {
  const seen = new Set<string>();
  let cursor: string | undefined;
  for (;;) {
    const request = cursorRequest({ pageSize: MAX_PAGE_SIZE, cursor });
    const reply = await fetchPage(request);
    const found = (reply.rows ?? []).find(match);
    if (found !== undefined) {
      return found;
    }
    const page = buildCursorPage(request.pageSize, reply.cursor);
    if (!page.hasMore || seen.has(page.nextCursor)) {
      return undefined;
    }
    seen.add(page.nextCursor);
    cursor = page.nextCursor;
  }
}
