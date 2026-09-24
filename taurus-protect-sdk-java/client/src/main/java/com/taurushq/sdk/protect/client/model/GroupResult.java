package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of groups from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.GroupService
 */
public final class GroupResult extends OffsetPagedResult<Group> {

    /**
     * Creates a result.
     *
     * @param groups the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public GroupResult(final List<Group> groups, final OffsetPagination pagination) {
        super(groups, pagination);
    }

    /**
     * Gets the groups of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<Group> getGroups() {
        return items();
    }
}
