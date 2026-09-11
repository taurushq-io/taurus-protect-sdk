package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * The set of whitelisted-address rows an approver reviewed, carrying the metadata hash
 * each one had AT REVIEW TIME. It is the content pin the approval path signs against.
 *
 * <p><b>Why this type exists rather than a plain list of ids.</b> The approval API accepts
 * only ids: the SDK re-reads them and signs whatever the server returns under those ids.
 * Nothing bound the approver's intent to the bytes signed, so a response-controlling server
 * could answer the id-filtered re-read with a DIFFERENT row — one whose existing signatures
 * already satisfy the container it presents — and harvest a genuine approver signature over
 * content the approver never saw. This is the same shape
 * {@code GovernanceRuleService.approveRulesProposal} was hardened against with its
 * mandatory {@code expectedContainerHash}.
 *
 * <p><b>It cannot be forged.</b> The pin map is a private field and there is no public
 * constructor, so a value can only be minted inside this package — in practice by
 * {@link WhitelistedAddressListResult#select(List)} or
 * {@link WhitelistedAddressListResult#selectAll()}, i.e. only from a read that verified.
 * Go's peer settles for "forgeable but useless" (its unexported fields still permit a
 * zero value from another package); Java can enforce it outright, so it does.
 *
 * <p>The pinned value is the row's CURRENT {@code metadata.hash}, because that is what the
 * approval signs and what validatord independently recomputes and checks. A mismatch at
 * approval time is a REAL signal, not a false positive: the server recomputes the metadata
 * on every read from the immutable envelope plus the row's live linked-address and
 * linked-wallet rows, so it moves when a linked address is renamed. Either way the content
 * changed since review — re-read, re-review, re-approve.
 *
 * @see WhitelistedAddressListResult#select(List)
 * @see com.taurushq.sdk.protect.client.service.WhitelistedAddressService
 */
public final class WhitelistedAddressApproval {

    /**
     * Row id to the metadata hash that row carried when it was reviewed. Private, with no
     * accessor that hands the map out, so the pin cannot be edited after minting.
     */
    private final Map<Long, String> pinned;

    /**
     * Package-private on purpose: only a verified read may mint a usable selection.
     *
     * @param pinned row id to reviewed metadata hash
     */
    WhitelistedAddressApproval(final Map<Long, String> pinned) {
        this.pinned = new LinkedHashMap<>(pinned);
    }

    /**
     * Returns the pinned row ids.
     *
     * <p>The approval path sorts them itself: the endpoint requires ascending order, and
     * the signed hash array must not depend on the order the caller happened to select in.
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
     * Reports whether this selection pins nothing.
     *
     * <p>An empty selection must be an ERROR at the approval site rather than "approve
     * nothing" — silently restoring unpinned behaviour is exactly the failure mode
     * {@code approveRulesProposal}'s mandatory pin also guards against.
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
        // Ids only. The hashes are not secret, but a log line is not where a reviewer
        // should be reading the pin from.
        return "WhitelistedAddressApproval{ids=" + pinned.keySet() + "}";
    }
}
