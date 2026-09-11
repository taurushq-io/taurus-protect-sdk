package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

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

    /**
     * Pins the given ids from this verified read, producing the selection
     * {@code approveWhitelistedAssets} requires.
     *
     * <p>This is the only way a usable {@link WhitelistedAssetApproval} comes into
     * existence, so the content pin cannot be forgotten. See
     * {@link WhitelistedAddressApproval} for the substitution attack it defeats.
     *
     * <p>An id this read did not return is an ERROR rather than a silent omission.
     *
     * @param ids the row ids to pin
     * @return the reviewed selection
     * @throws IntegrityException if an id is absent from this read or carries no metadata
     *                            hash to pin
     */
    public WhitelistedAssetApproval select(final List<Long> ids) {
        if (ids == null || ids.isEmpty()) {
            throw new IntegrityException("cannot select an empty set of ids: an empty pin "
                    + "would silently restore unpinned approval");
        }

        Map<Long, SignedWhitelistedAssetEnvelope> byId = new HashMap<>();
        if (assets != null) {
            for (SignedWhitelistedAssetEnvelope envelope : assets) {
                if (envelope != null) {
                    byId.put(envelope.getId(), envelope);
                }
            }
        }

        Map<Long, String> pinned = new LinkedHashMap<>();
        for (Long id : ids) {
            if (id == null) {
                throw new IntegrityException("cannot select a null whitelisted asset id");
            }
            SignedWhitelistedAssetEnvelope envelope = byId.get(id);
            if (envelope == null) {
                throw new IntegrityException(String.format(
                        "whitelisted asset %d is not in this verified read: it was either "
                                + "excluded as unverifiable or not on this page", id));
            }
            if (envelope.getMetadata() == null
                    || envelope.getMetadata().getHash() == null
                    || envelope.getMetadata().getHash().isEmpty()) {
                throw new IntegrityException(String.format(
                        "whitelisted asset %d carries no metadata hash, so there is nothing "
                                + "to pin the approval to", id));
            }
            pinned.put(id, envelope.getMetadata().getHash());
        }
        return new WhitelistedAssetApproval(pinned);
    }

    /**
     * Pins every row this verified read returned. Approving is all-or-nothing over what is
     * pinned here.
     *
     * @return the reviewed selection
     * @throws IntegrityException if this read returned no verified rows
     */
    public WhitelistedAssetApproval selectAll() {
        List<Long> ids = new ArrayList<>();
        if (assets != null) {
            for (SignedWhitelistedAssetEnvelope envelope : assets) {
                if (envelope != null) {
                    ids.add(envelope.getId());
                }
            }
        }
        if (ids.isEmpty()) {
            throw new IntegrityException(
                    "this read returned no verified whitelisted assets to approve");
        }
        return select(ids);
    }
}
