package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.CursorPagedResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;

import java.util.Base64;

/**
 * The wire form of one cursor-list page request: page size (always), plus currentPage and
 * pageRequest when the caller continues a walk.
 * <p>
 * Every cursor list goes through {@link #of(ApiRequestCursor)}, whatever form the operation
 * takes on the wire ({@code cursor.*} query, {@code requestCursor.*} query, or a body
 * cursor), so a page size is always sent.
 */
final class CursorRequest {

    private final int pageSize;
    private final String currentPage;
    private final String pageRequest;

    private CursorRequest(final int pageSize, final String currentPage, final String pageRequest) {
        this.pageSize = pageSize;
        this.currentPage = currentPage;
        this.pageRequest = pageRequest;
    }

    /**
     * Converts a request cursor. A null cursor, or one whose page size is unset, sends
     * {@link Pagination#DEFAULT_PAGE_SIZE}.
     *
     * @param cursor the caller's cursor, may be null
     * @return the wire form
     * @throws IllegalArgumentException if the page size is out of range
     */
    static CursorRequest of(final ApiRequestCursor cursor) {
        if (cursor == null) {
            return new CursorRequest(Pagination.DEFAULT_PAGE_SIZE, null, null);
        }
        long size = cursor.getPageSize();
        if (size < 0 || size > Pagination.MAX_PAGE_SIZE) {
            throw new IllegalArgumentException("pageSize must be between 1 and "
                    + Pagination.MAX_PAGE_SIZE + ", got " + size);
        }
        String current = cursor.getCurrentPage();
        return new CursorRequest(Pagination.resolvePageSize("pageSize", (int) size),
                current == null || current.isEmpty() ? null : current,
                cursor.getPageRequest() == null ? null : cursor.getPageRequest().name());
    }

    /**
     * Returns the page size sent.
     *
     * @return the page size
     */
    int pageSize() {
        return pageSize;
    }

    /**
     * Returns the page size as the query/body value.
     *
     * @return the page size text
     */
    String pageSizeParam() {
        return String.valueOf(pageSize);
    }

    /**
     * Returns the current page token.
     *
     * @return the token, or null on a first page
     */
    String currentPage() {
        return currentPage;
    }

    /**
     * Returns the page request.
     *
     * @return the page request name, or null when none is sent
     */
    String pageRequest() {
        return pageRequest;
    }

    /**
     * Returns the body-cursor form.
     *
     * @return a new request-cursor DTO
     */
    TgvalidatordRequestCursor toDTO() {
        TgvalidatordRequestCursor dto = new TgvalidatordRequestCursor();
        dto.setCurrentPage(currentPage);
        dto.setPageRequest(pageRequest);
        dto.setPageSize(pageSizeParam());
        return dto;
    }

    /**
     * Completes a cursor result with the reply's cursor and the page built from it.
     *
     * @param result      the result, rows already set
     * @param op          the operation, which decides whether the reply carries a total
     * @param replyCursor the reply's cursor, may be null
     * @param totalItems  the reply's total as received, ignored when the operation has none
     * @param <R>         the result type
     * @return the same result
     */
    <R extends CursorPagedResult> R complete(final R result, final PagedOperation op,
                                             final TgvalidatordResponseCursor replyCursor,
                                             final String totalItems) {
        ApiResponseCursor cursor = ApiResponseCursorMapper.INSTANCE.fromDTO(replyCursor);
        result.setCursor(cursor);
        result.setPage(op.cursorPage(pageSize, cursor, totalItems));
        return result;
    }

    /**
     * Completes a cursor result whose cursor the reply mapper already set.
     *
     * @param result the mapped result
     * @param op     the operation
     * @param <R>    the result type
     * @return the same result
     */
    <R extends CursorPagedResult> R complete(final R result, final PagedOperation op) {
        result.setPage(op.cursorPage(pageSize, result.getCursor(), null));
        return result;
    }

    /**
     * Reports whether a caller passed a token, i.e. is continuing a token-paged walk.
     *
     * @param cursor the caller's cursor, may be null
     * @return true when there is a token to send
     */
    static boolean hasToken(final String cursor) {
        return cursor != null && !cursor.isEmpty();
    }

    /**
     * Converts a caller's token (a previous next cursor) to the bytes the generated client
     * takes for the token operations, which re-encodes them as standard base64 on the wire.
     *
     * @param cursor the token, non-empty (see {@link #hasToken})
     * @return the token bytes
     * @throws IllegalArgumentException if the token is not base64
     */
    static byte[] tokenBytes(final String cursor) {
        try {
            boolean urlSafe = cursor.indexOf('-') >= 0 || cursor.indexOf('_') >= 0;
            return (urlSafe ? Base64.getUrlDecoder() : Base64.getDecoder()).decode(cursor);
        } catch (IllegalArgumentException e) {
            throw new IllegalArgumentException(
                    "cursor must be the nextCursor of a previous page (base64 text)", e);
        }
    }

    /**
     * Renders a reply token as the opaque cursor text callers pass back.
     *
     * @param token the reply's token bytes, may be null
     * @return the standard-base64 text, empty when there is no token
     */
    static String tokenText(final byte[] token) {
        return token == null || token.length == 0 ? "" : Base64.getEncoder().encodeToString(token);
    }
}
