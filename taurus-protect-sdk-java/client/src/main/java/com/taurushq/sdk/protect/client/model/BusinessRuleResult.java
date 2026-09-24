package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a business rule query with cursor-based pagination.
 * <p>
 * Contains a page of business rules and cursor information for fetching
 * additional pages: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true.
 *
 * @see BusinessRule
 */
public class BusinessRuleResult extends CursorPagedResult {

    private List<BusinessRule> rules;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the business rules of this page.
     *
     * @return the rules
     */
    public List<BusinessRule> getRules() {
        return rules;
    }

    /**
     * Sets the business rules of this page.
     *
     * @param rules the rules
     */
    public void setRules(final List<BusinessRule> rules) {
        this.rules = rules;
    }
}
