package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a business rule query with cursor-based pagination.
 * <p>
 * Contains a page of business rules and cursor information for fetching
 * additional pages. Use {@link #hasNext()} to check for more pages and
 * {@link #nextCursor(long)} to create a cursor for the next page.
 *
 * @see BusinessRule
 * @see ApiResponseCursor
 */
public class BusinessRuleResult {

    /**
     * The list of business rules in this page of results.
     */
    private List<BusinessRule> rules;

    /**
     * Pagination cursor containing page information.
     */
    private ApiResponseCursor cursor;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the list of business rules.
     *
     * @return the rules
     */
    public List<BusinessRule> getRules() {
        return rules;
    }

    /**
     * Sets the list of business rules.
     *
     * @param rules the rules
     */
    public void setRules(List<BusinessRule> rules) {
        this.rules = rules;
    }

    /**
     * Gets the response cursor for pagination.
     *
     * @return the cursor
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the response cursor for pagination.
     *
     * @param cursor the cursor
     */
    public void setCursor(ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Creates a cursor for the next page, or null if no more pages.
     *
     * @param pageSize the page size
     * @return the next cursor, or null if no more pages
     */
    public ApiRequestCursor nextCursor(long pageSize) {
        if (hasNext()) {
            return new ApiRequestCursor(cursor.getCurrentPage(), PageRequest.NEXT, pageSize);
        }
        return null;
    }

    /**
     * Checks if there is a next page of results.
     *
     * @return true if there is a next page
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }
}
