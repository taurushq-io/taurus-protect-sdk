package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of whitelisted contract addresses.
 * <p>
 * This class wraps a list of contract address envelopes along with pagination
 * information to support offset-based navigation through large result sets.
 *
 * @see SignedWhitelistedContractAddressEnvelope
 * @see ContractWhitelistingService
 */
public class WhitelistedContractAddressResult {

    /**
     * The list of contract address envelopes in this page.
     */
    private List<SignedWhitelistedContractAddressEnvelope> contracts;

    /**
     * Total number of items matching the query.
     */
    private long totalItems;

    /**
     * Gets the list of contract address envelopes.
     *
     * @return the contracts
     */
    public List<SignedWhitelistedContractAddressEnvelope> getContracts() {
        return contracts;
    }

    /**
     * Sets the list of contract address envelopes.
     *
     * @param contracts the contracts to set
     */
    public void setContracts(List<SignedWhitelistedContractAddressEnvelope> contracts) {
        this.contracts = contracts;
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
     * Checks if there are more results available beyond this page.
     *
     * @param currentOffset the current offset
     * @param pageSize      the page size
     * @return true if more results are available
     */
    public boolean hasMore(int currentOffset, int pageSize) {
        return (currentOffset + pageSize) < totalItems;
    }
}
