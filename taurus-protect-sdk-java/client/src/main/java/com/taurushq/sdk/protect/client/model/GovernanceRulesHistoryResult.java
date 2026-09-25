package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.Collections;
import java.util.List;

/**
 * One page of the governance rules history, with its {@link CursorPage}.
 * <p>
 * The history pages by an opaque token: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true. The page total is the server's count reduced by the
 * entries withheld because their SuperAdmin signatures did not verify.
 *
 * @see GovernanceRules
 */
public class GovernanceRulesHistoryResult {

    /**
     * The list of governance rules in this page of results.
     */
    private List<GovernanceRules> rules;

    /**
     * The page: next cursor, whether more pages exist, and the reduced total.
     */
    private CursorPage page;

    /**
     * Entries withheld because their SuperAdmin signatures did not verify.
     */
    private List<ExcludedRuleset> excludedUnverified = Collections.emptyList();

    /**
     * Gets the entries withheld because their signatures did not verify, so a shortened
     * page cannot read as a complete one.
     *
     * @return the excluded entries, never null
     */
    public List<ExcludedRuleset> getExcludedUnverified() {
        return excludedUnverified;
    }

    /**
     * Sets the entries withheld because their signatures did not verify.
     *
     * @param excludedUnverified the excluded entries
     */
    public void setExcludedUnverified(final List<ExcludedRuleset> excludedUnverified) {
        this.excludedUnverified = excludedUnverified;
    }

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
     * Gets the page: next cursor, whether more pages exist, and the total reduced by the
     * withheld entries.
     *
     * @return the page, never null on a result returned by the service
     */
    public CursorPage getPage() {
        return page;
    }

    /**
     * Sets the page.
     *
     * @param page the page
     */
    public void setPage(final CursorPage page) {
        this.page = page;
    }

    /**
     * Returns true if there are more pages available.
     *
     * @return true if more pages available
     */
    public boolean hasMorePages() {
        return page != null && page.hasMore();
    }
}
