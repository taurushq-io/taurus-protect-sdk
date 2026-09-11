package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * The asset peer of {@link WhitelistedAddressApproval}: the rows an approver reviewed, with
 * the metadata hash each carried at review time.
 *
 * <p>See {@link WhitelistedAddressApproval} for why the approval path needs a content pin
 * rather than bare ids, and why the pin map is private with no public constructor — a
 * usable selection can only be minted from a read that verified, via
 * {@link WhitelistedAssetResult#select(List)} or {@link WhitelistedAssetResult#selectAll()}.
 *
 * @see com.taurushq.sdk.protect.client.service.WhitelistedAssetService
 */
public final class WhitelistedAssetApproval {

    /** Row id to the metadata hash that row carried when it was reviewed. */
    private final Map<Long, String> pinned;

    /**
     * Package-private on purpose: only a verified read may mint a usable selection.
     *
     * @param pinned row id to reviewed metadata hash
     */
    WhitelistedAssetApproval(final Map<Long, String> pinned) {
        this.pinned = new LinkedHashMap<>(pinned);
    }

    /**
     * Returns the pinned row ids. The approval path sorts them itself.
     *
     * @return an unmodifiable list of ids
     */
    public List<Long> getIds() {
        return Collections.unmodifiableList(new ArrayList<>(pinned.keySet()));
    }

    /**
     * Returns the reviewed metadata hash for a row.
     *
     * @param id the row id
     * @return the hash pinned at review time, or {@code null} when this row was not part
     *         of the reviewed selection
     */
    public String pinnedHash(final long id) {
        return pinned.get(id);
    }

    /**
     * Reports whether this selection pins nothing. An empty selection is an ERROR at the
     * approval site, never "approve nothing".
     *
     * @return true when nothing is pinned
     */
    public boolean isEmpty() {
        return pinned.isEmpty();
    }

    /**
     * Returns how many rows are pinned.
     *
     * @return the pin count
     */
    public int size() {
        return pinned.size();
    }

    @Override
    public String toString() {
        return "WhitelistedAssetApproval{ids=" + pinned.keySet() + "}";
    }
}
