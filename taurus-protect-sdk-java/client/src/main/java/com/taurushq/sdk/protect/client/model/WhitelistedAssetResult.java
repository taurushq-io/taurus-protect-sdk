package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of verified whitelisted asset envelopes.
 * <p>
 * The list-returning overloads on the service carry no page total, so a caller could not
 * tell a full page from the last one.
 *
 * @see SignedWhitelistedAssetEnvelope
 */
public class WhitelistedAssetResult {

    /**
     * The verified asset envelopes in this page.
     */
    private List<SignedWhitelistedAssetEnvelope> assets;

    /**
     * Total number of items matching the query.
     */
    private long totalItems;

    /**
     * Gets the verified asset envelopes.
     *
     * @return the assets
     */
    public List<SignedWhitelistedAssetEnvelope> getAssets() {
        return assets;
    }

    /**
     * Sets the verified asset envelopes.
     *
     * @param assets the assets to set
     */
    public void setAssets(List<SignedWhitelistedAssetEnvelope> assets) {
        this.assets = assets;
    }

    /**
     * Gets the total number of items.
     *
     * @return the total items count
     */
    public long getTotalItems() {
        return totalItems;
    }

    /**
     * Sets the total number of items.
     *
     * @param totalItems the total items count to set
     */
    public void setTotalItems(long totalItems) {
        this.totalItems = totalItems;
    }

    /**
     * Checks whether more results exist beyond this page.
     * <p>
     * Overflow-safe: currentOffset + pageSize can wrap on caller-supplied values.
     *
     * @param currentOffset the current offset
     * @param pageSize      the page size
     * @return true if more results are available
     */
    public boolean hasMore(int currentOffset, int pageSize) {
        return totalItems > currentOffset && totalItems - currentOffset > pageSize;
    }
}
