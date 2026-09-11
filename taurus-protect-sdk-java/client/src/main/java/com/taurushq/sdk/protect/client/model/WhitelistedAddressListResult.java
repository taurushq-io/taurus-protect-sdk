package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Result of a whitelisted-address list query: the rows that verified, plus the rows
 * that did not.
 * <p>
 * Verification is lenient by design. A row that cannot be verified is excluded rather
 * than failing the whole call, because one bad row used to deny access to every good
 * one — and listing is how an operator finds the bad row. Excluding stays fail-closed,
 * since an omitted destination cannot be selected.
 * <p>
 * The omission is <em>reported</em>, not only logged: a caller cannot read the SDK's
 * logger, and a shortened list must never be mistaken for a complete one. A caller
 * asking "is this destination approved?" would otherwise get a false negative.
 *
 * @see com.taurushq.sdk.protect.client.service.WhitelistedAddressService
 */
public final class WhitelistedAddressListResult {

    private final List<SignedWhitelistedAddressEnvelope> envelopes;
    private final List<ExcludedWhitelistedAddress> excludedUnverified;
    private final String totalItems;

    /**
     * Constructs a result with no page total.
     *
     * @param envelopes          the envelopes that passed verification
     * @param excludedUnverified the rows dropped for failing verification
     */
    public WhitelistedAddressListResult(
            final List<SignedWhitelistedAddressEnvelope> envelopes,
            final List<ExcludedWhitelistedAddress> excludedUnverified) {
        this(envelopes, excludedUnverified, null);
    }

    /**
     * Constructs a result carrying the page total.
     *
     * @param envelopes          the envelopes that passed verification
     * @param excludedUnverified the rows dropped for failing verification
     * @param totalItems         the total already reduced by the exclusion count, or
     *                           {@code null} when the server reported none
     */
    public WhitelistedAddressListResult(
            final List<SignedWhitelistedAddressEnvelope> envelopes,
            final List<ExcludedWhitelistedAddress> excludedUnverified,
            final String totalItems) {
        this.envelopes = envelopes == null
                ? new ArrayList<>() : new ArrayList<>(envelopes);
        this.excludedUnverified = excludedUnverified == null
                ? new ArrayList<>() : new ArrayList<>(excludedUnverified);
        this.totalItems = totalItems;
    }

    /**
     * Returns the page total, already reduced by the number of excluded rows.
     *
     * <p>The server counts the rows it returned; the caller receives only the ones that
     * verified. Reporting the server's total would make this promise rows that can never
     * be read, and would let a filtered page pass for a complete one. This is a
     * {@code String} because the wire type is a {@code uint64} the generated client
     * surfaces as text; Go and TypeScript expose the same value as a number.
     *
     * @return the adjusted total, or {@code null} when the server reported none
     */
    public String getTotalItems() {
        return totalItems;
    }

    /**
     * Returns the envelopes that passed verification.
     *
     * @return an unmodifiable list of verified envelopes
     */
    public List<SignedWhitelistedAddressEnvelope> getEnvelopes() {
        return Collections.unmodifiableList(envelopes);
    }

    /**
     * Returns the rows dropped from the result for failing verification.
     *
     * @return an unmodifiable list of exclusions, empty when nothing was dropped
     */
    public List<ExcludedWhitelistedAddress> getExcludedUnverified() {
        return Collections.unmodifiableList(excludedUnverified);
    }

    /**
     * Pins the given ids from this verified read, producing the selection
     * {@code approveWhitelistedAddresses} requires.
     *
     * <p>This is the only way a usable {@link WhitelistedAddressApproval} comes into
     * existence, which is the point: the approval cannot be reached with bare ids, so the
     * content pin cannot be forgotten. See {@link WhitelistedAddressApproval} for the
     * substitution attack it defeats.
     *
     * <p>An id this read did not return is an ERROR rather than a silent omission: it
     * means the caller is trying to approve something this read did not verify — either it
     * was excluded as unverifiable (see {@link #getExcludedUnverified()}) or it was never
     * on the page — and approving fewer rows than asked for would tell the approver they
     * approved more than they did.
     *
     * @param ids the row ids to pin
     * @return the reviewed selection
     * @throws IntegrityException if an id is absent from this read or carries no metadata
     *                            hash to pin
     */
    public WhitelistedAddressApproval select(final List<Long> ids) {
        if (ids == null || ids.isEmpty()) {
            throw new IntegrityException("cannot select an empty set of ids: an empty pin "
                    + "would silently restore unpinned approval");
        }

        Map<Long, SignedWhitelistedAddressEnvelope> byId = new HashMap<>();
        for (SignedWhitelistedAddressEnvelope envelope : envelopes) {
            if (envelope != null) {
                byId.put(envelope.getId(), envelope);
            }
        }

        Map<Long, String> pinned = new LinkedHashMap<>();
        for (Long id : ids) {
            if (id == null) {
                throw new IntegrityException("cannot select a null whitelisted address id");
            }
            SignedWhitelistedAddressEnvelope envelope = byId.get(id);
            if (envelope == null) {
                throw new IntegrityException(String.format(
                        "whitelisted address %d is not in this verified read: it was either "
                                + "excluded as unverifiable or not on this page", id));
            }
            if (envelope.getMetadata() == null
                    || envelope.getMetadata().getHash() == null
                    || envelope.getMetadata().getHash().isEmpty()) {
                throw new IntegrityException(String.format(
                        "whitelisted address %d carries no metadata hash, so there is nothing "
                                + "to pin the approval to", id));
            }
            pinned.put(id, envelope.getMetadata().getHash());
        }
        return new WhitelistedAddressApproval(pinned);
    }

    /**
     * Pins every row this verified read returned.
     *
     * <p>Note it pins what SURVIVED verification, not what the server sent: rows in
     * {@link #getExcludedUnverified()} are not included, so a caller who wants to know
     * about them must read that list. Approving is all-or-nothing over what is pinned here.
     *
     * @return the reviewed selection
     * @throws IntegrityException if this read returned no verified rows
     */
    public WhitelistedAddressApproval selectAll() {
        List<Long> ids = new ArrayList<>(envelopes.size());
        for (SignedWhitelistedAddressEnvelope envelope : envelopes) {
            if (envelope != null) {
                ids.add(envelope.getId());
            }
        }
        if (ids.isEmpty()) {
            throw new IntegrityException(
                    "this read returned no verified whitelisted addresses to approve");
        }
        return select(ids);
    }
}
