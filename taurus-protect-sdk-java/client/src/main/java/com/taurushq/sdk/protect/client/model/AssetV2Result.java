package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of assets (v2 asset service) from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService
 */
public class AssetV2Result extends CursorPagedResult {

    private List<AssetV2> assets;

    /**
     * Gets the assets of this page.
     *
     * @return the assets
     */
    public List<AssetV2> getAssets() {
        return assets;
    }

    /**
     * Sets the assets of this page.
     *
     * @param assets the assets
     */
    public void setAssets(final List<AssetV2> assets) {
        this.assets = assets;
    }
}
