package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of addresses from an offset list, with its {@link OffsetPagination}. Every
 * address string has had its HSM signature verified.
 *
 * @see com.taurushq.sdk.protect.client.service.AddressService
 */
public final class AddressResult extends OffsetPagedResult<Address> {

    /**
     * Creates a result.
     *
     * @param addresses the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public AddressResult(final List<Address> addresses, final OffsetPagination pagination) {
        super(addresses, pagination);
    }

    /**
     * Gets the addresses of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<Address> getAddresses() {
        return items();
    }
}
