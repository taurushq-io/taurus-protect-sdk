package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

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
}
