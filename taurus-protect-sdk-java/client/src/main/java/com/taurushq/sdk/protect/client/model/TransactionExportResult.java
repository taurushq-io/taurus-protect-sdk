package com.taurushq.sdk.protect.client.model;

/**
 * A transaction export: the exported text plus the server's count of matching transactions.
 * <p>
 * The export cannot page: the server ignores any offset and always exports from the first
 * matching transaction, so a larger limit is the only way to reach more rows. There is
 * therefore no next offset here; compare {@link #getTotalItems()} with the limit sent to
 * tell whether the export was cut short.
 *
 * @see com.taurushq.sdk.protect.client.service.TransactionService
 */
public final class TransactionExportResult {

    private final String content;
    private final long totalItems;

    /**
     * Creates a result.
     *
     * @param content    the exported text, null for none
     * @param totalItems the number of matching transactions the server reported
     */
    public TransactionExportResult(final String content, final long totalItems) {
        this.content = content == null ? "" : content;
        this.totalItems = totalItems;
    }

    /**
     * Builds a result from the reply's fields as received.
     *
     * @param content    the reply's {@code result}, null when absent
     * @param totalItems the reply's {@code totalItems}, null when absent
     * @return the result, never null
     * @throws IntegrityException if the count is not a canonical count in range
     */
    public static TransactionExportResult of(final String content, final String totalItems) {
        return new TransactionExportResult(content, Pagination.parseCount("totalItems", totalItems));
    }

    /**
     * Gets the exported text, in the format asked for or the server's default (JSON).
     *
     * @return the content, empty when nothing matched
     */
    public String getContent() {
        return content;
    }

    /**
     * Gets the number of matching transactions, which can exceed the rows exported.
     *
     * @return the total, 0 when the server reported none
     */
    public long getTotalItems() {
        return totalItems;
    }
}
