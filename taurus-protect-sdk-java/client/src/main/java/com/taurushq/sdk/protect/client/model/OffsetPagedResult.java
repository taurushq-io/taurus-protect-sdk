package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;

/**
 * Base of every offset-list result: one page of rows plus its {@link OffsetPagination}.
 *
 * @param <T> the row type
 */
public abstract class OffsetPagedResult<T> {

    private final List<T> items;
    private final OffsetPagination pagination;

    /**
     * Creates a result.
     *
     * @param items      the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    protected OffsetPagedResult(final List<T> items, final OffsetPagination pagination) {
        this.items = items == null ? Collections.<T>emptyList()
                : Collections.unmodifiableList(new ArrayList<>(items));
        this.pagination = Objects.requireNonNull(pagination, "pagination cannot be null");
    }

    /**
     * Gets the page's pagination: limit and offset sent, total, next offset, has-more.
     *
     * @return the pagination, never null
     */
    public OffsetPagination getPagination() {
        return pagination;
    }

    /**
     * Gets the rows of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    protected List<T> items() {
        return items;
    }
}
