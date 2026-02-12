/**
 * Pagination information for list responses.
 */
export interface Pagination {
  /** Total number of items across all pages */
  readonly totalItems: number;
  /** Current offset (0-based) */
  readonly offset: number;
  /** Number of items per page */
  readonly limit: number;
}

/**
 * Cursor-based pagination for TaurusNetwork services.
 */
export interface CursorPagination {
  /** Cursor for the next page, undefined if no more pages */
  readonly nextCursor: string | undefined;
  /** Whether there are more pages available */
  readonly hasMore: boolean;
}

/**
 * A paginated result containing items and pagination info.
 */
export interface PaginatedResult<T> {
  /** The items in this page */
  readonly items: T[];
  /** Pagination information */
  readonly pagination: Pagination;
}

/**
 * A cursor-paginated result containing items and cursor info.
 */
export interface CursorPaginatedResult<T> {
  /** The items in this page */
  readonly items: T[];
  /** Cursor pagination information */
  readonly pagination: CursorPagination;
}
