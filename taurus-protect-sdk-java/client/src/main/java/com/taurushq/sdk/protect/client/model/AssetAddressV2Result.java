package com.taurushq.sdk.protect.client.model;

import java.util.Collections;
import java.util.List;

/**
 * One page of addresses of an asset (v2 asset service) from a cursor list. Continue with {@code getPage()}.
 * <p>
 * INTERNAL and WHITELISTED rows that could not be confirmed against their verified
 * counterpart are withheld and named in {@link #getExcludedUnverified()}. Exclusions never
 * move the cursor.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService
 */
public class AssetAddressV2Result extends CursorPagedResult {

    private List<AssetAddressV2> addresses;
    private List<ExcludedWhitelistedAddress> excludedUnverified = Collections.emptyList();

    /**
     * Gets the addresses of an asset of this page.
     *
     * @return the addresses
     */
    public List<AssetAddressV2> getAddresses() {
        return addresses;
    }

    /**
     * Sets the addresses of an asset of this page.
     *
     * @param addresses the addresses
     */
    public void setAddresses(final List<AssetAddressV2> addresses) {
        this.addresses = addresses;
    }

    /**
     * Gets the rows withheld because they could not be verified. The id is the row's
     * addressID or whitelistedAddressID, or its address when it carries neither.
     *
     * @return the excluded rows, never null
     */
    public List<ExcludedWhitelistedAddress> getExcludedUnverified() {
        return excludedUnverified;
    }

    /**
     * Sets the rows withheld because they could not be verified.
     *
     * @param excludedUnverified the excluded rows, null for none
     */
    public void setExcludedUnverified(final List<ExcludedWhitelistedAddress> excludedUnverified) {
        this.excludedUnverified = excludedUnverified == null
                ? Collections.emptyList() : excludedUnverified;
    }
}
