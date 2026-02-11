package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a governance rules history query with cursor-based pagination.
 * <p>
 * This class wraps the paginated response from querying the history of governance
 * rules changes. Use the cursor to fetch subsequent pages of results.
 *
 * @see GovernanceRules
 */
public class GovernanceRulesHistoryResult {

    /**
     * The list of governance rules in this page of results.
     */
    private List<GovernanceRules> rules;

    /**
     * Opaque cursor for fetching the next page of results.
     */
    private byte[] cursor;

    /**
     * The total number of items available across all pages.
     */
    private String totalItems;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the list of governance rules.
     *
     * @return the rules
     */
    public List<GovernanceRules> getRules() {
        return rules;
    }

    /**
     * Sets the list of governance rules.
     *
     * @param rules the rules
     */
    public void setRules(List<GovernanceRules> rules) {
        this.rules = rules;
    }

    /**
     * Gets the cursor for the next page.
     *
     * @return the cursor for next page, or null if no more pages
     */
    public byte[] getCursor() {
        return cursor;
    }

    /**
     * Sets the cursor for the next page.
     *
     * @param cursor the cursor
     */
    public void setCursor(byte[] cursor) {
        this.cursor = cursor;
    }

    /**
     * Gets total items count.
     *
     * @return the total items
     */
    public String getTotalItems() {
        return totalItems;
    }

    /**
     * Sets total items count.
     *
     * @param totalItems the total items
     */
    public void setTotalItems(String totalItems) {
        this.totalItems = totalItems;
    }

    /**
     * Returns true if there are more pages available.
     *
     * @return true if more pages available
     */
    public boolean hasMorePages() {
        return cursor != null && cursor.length > 0;
    }
}
