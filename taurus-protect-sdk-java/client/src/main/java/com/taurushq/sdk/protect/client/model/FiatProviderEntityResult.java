package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of fiat provider entities from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.FiatService
 */
public class FiatProviderEntityResult extends CursorPagedResult {

    private List<FiatProviderEntity> entities;

    /**
     * Gets the fiat provider entities of this page.
     *
     * @return the entities
     */
    public List<FiatProviderEntity> getEntities() {
        return entities;
    }

    /**
     * Sets the fiat provider entities of this page.
     *
     * @param entities the entities
     */
    public void setEntities(final List<FiatProviderEntity> entities) {
        this.entities = entities;
    }
}
