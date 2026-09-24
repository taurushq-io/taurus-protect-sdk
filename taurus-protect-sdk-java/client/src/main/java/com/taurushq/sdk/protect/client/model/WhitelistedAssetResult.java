package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Represents a paginated result of verified whitelisted asset envelopes.
 * <p>
 * Every list method on the service returns this, with the page's {@link OffsetPagination}.
 *
 * @see SignedWhitelistedAssetEnvelope
 */
public class WhitelistedAssetResult {

    /**
     * The verified asset envelopes in this page.
     */
    private List<SignedWhitelistedAssetEnvelope> assets;

    /**
     * The page's pagination.
     */
    private OffsetPagination pagination;

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
     * Gets the page's pagination: limit and offset sent, total, next offset, has-more.
     * Skipped rows keep their SQL slot on this endpoint, so a short page is not the end:
     * continue while {@code hasMore()}.
     *
     * @return the pagination, never null on a result returned by the service
     */
    public OffsetPagination getPagination() {
        return pagination;
    }

    /**
     * Sets the page's pagination.
     *
     * @param pagination the pagination
     */
    public void setPagination(final OffsetPagination pagination) {
        this.pagination = pagination;
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
