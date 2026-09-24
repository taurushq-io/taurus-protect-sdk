package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of actions from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.ActionService
 */
public final class ActionResult extends OffsetPagedResult<ActionEnvelope> {

    /**
     * Creates a result.
     *
     * @param actions the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public ActionResult(final List<ActionEnvelope> actions, final OffsetPagination pagination) {
        super(actions, pagination);
    }

    /**
     * Gets the actions of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<ActionEnvelope> getActions() {
        return items();
    }
}
