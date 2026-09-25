package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a change query with cursor-based pagination.
 * <p>
 * Contains a page of change audit records and cursor information for fetching
 * additional pages: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true.
 *
 * @see Change
 */
public class ChangeResult extends CursorPagedResult {

    private List<Change> changes;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the changes of this page.
     *
     * @return the changes
     */
    public List<Change> getChanges() {
        return changes;
    }

    /**
     * Sets the changes of this page.
     *
     * @param changes the changes
     */
    public void setChanges(final List<Change> changes) {
        this.changes = changes;
    }
}
