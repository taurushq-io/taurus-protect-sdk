package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of users from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.UserService
 */
public final class UserResult extends OffsetPagedResult<User> {

    /**
     * Creates a result.
     *
     * @param users the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public UserResult(final List<User> users, final OffsetPagination pagination) {
        super(users, pagination);
    }

    /**
     * Gets the users of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<User> getUsers() {
        return items();
    }
}
