package com.taurushq.sdk.protect.client.model;

/**
 * Represents Cardano (ADA) stake pool information.
 * <p>
 * This model contains information about a Cardano stake pool, including
 * its pledge, margin, costs, and current active stake.
 */
public class ADAStakePoolInfo {

    private String pledge;
    private Float margin;
    private String fixedCost;
    private String url;
    private String activeStake;
    private String epoch;

    /**
     * Gets the pool's pledge amount in lovelace.
     *
     * @return the pledge amount
     */
    public String getPledge() {
        return pledge;
    }

    /**
     * Sets the pool's pledge amount.
     *
     * @param pledge the pledge amount
     */
    public void setPledge(String pledge) {
        this.pledge = pledge;
    }

    /**
     * Gets the pool's margin (fee percentage).
     *
     * @return the margin as a decimal (e.g., 0.05 for 5%)
     */
    public Float getMargin() {
        return margin;
    }

    /**
     * Sets the pool's margin.
     *
     * @param margin the margin
     */
    public void setMargin(Float margin) {
        this.margin = margin;
    }

    /**
     * Gets the pool's fixed cost per epoch in lovelace.
     *
     * @return the fixed cost
     */
    public String getFixedCost() {
        return fixedCost;
    }

    /**
     * Sets the pool's fixed cost.
     *
     * @param fixedCost the fixed cost
     */
    public void setFixedCost(String fixedCost) {
        this.fixedCost = fixedCost;
    }

    /**
     * Gets the pool's metadata URL.
     *
     * @return the metadata URL
     */
    public String getUrl() {
        return url;
    }

    /**
     * Sets the pool's metadata URL.
     *
     * @param url the metadata URL
     */
    public void setUrl(String url) {
        this.url = url;
    }

    /**
     * Gets the pool's active stake in lovelace.
     *
     * @return the active stake
     */
    public String getActiveStake() {
        return activeStake;
    }

    /**
     * Sets the pool's active stake.
     *
     * @param activeStake the active stake
     */
    public void setActiveStake(String activeStake) {
        this.activeStake = activeStake;
    }

    /**
     * Gets the current epoch number.
     *
     * @return the epoch
     */
    public String getEpoch() {
        return epoch;
    }

    /**
     * Sets the current epoch.
     *
     * @param epoch the epoch
     */
    public void setEpoch(String epoch) {
        this.epoch = epoch;
    }

    @Override
    public String toString() {
        return "ADAStakePoolInfo{"
                + "pledge='" + pledge + '\''
                + ", margin=" + margin
                + ", activeStake='" + activeStake + '\''
                + ", epoch='" + epoch + '\''
                + '}';
    }
}
