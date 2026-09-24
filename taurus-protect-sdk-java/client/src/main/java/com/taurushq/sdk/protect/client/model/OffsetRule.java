package com.taurushq.sdk.protect.client.model;

/**
 * How an offset list's next offset follows from a reply. Each offset endpoint uses exactly
 * one rule; {@link OffsetPagination#of} applies it.
 */
public enum OffsetRule {

    /**
     * The reply's {@code offset} is the NEXT page's offset (offset + rows); when it is
     * absent, offset + rows. Wallets and addresses.
     */
    REPLY_OFFSET,

    /**
     * offset + rows. Transactions, the transaction export, fee payers, actions.
     */
    PLUS_ROWS,

    /**
     * offset + min(rows, limit): a synthetic daemon user or technical group can be appended
     * beyond the limit. Users and groups.
     */
    PLUS_MIN_ROWS_LIMIT,

    /**
     * offset + the rows the SERVER returned, before SDK exclusions; rows the server drops for
     * a bad signature are gone for good. Whitelisted addresses.
     */
    PLUS_SERVER_ROWS,

    /**
     * offset + limit: rows the server skips keep their SQL slot, so a short page is not the
     * end. Whitelisted contracts (assets).
     */
    PLUS_LIMIT
}
