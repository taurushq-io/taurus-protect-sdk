package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a sequence of group thresholds that must be satisfied in order.
 * <p>
 * Sequential thresholds define approval requirements where each group's threshold
 * must be met before proceeding to the next group. This enables multi-stage
 * approval workflows (e.g., first get 2 from Finance, then 1 from Compliance).
 *
 * @see GroupThreshold
 * @see AddressWhitelistingRules
 * @see TransactionRules
 */
public class SequentialThresholds {

    /**
     * Ordered list of group thresholds to satisfy sequentially.
     */
    private List<GroupThreshold> thresholds;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the thresholds.
     *
     * @return the thresholds
     */
    public List<GroupThreshold> getThresholds() {
        return thresholds;
    }

    /**
     * Sets the thresholds.
     *
     * @param thresholds the thresholds
     */
    public void setThresholds(List<GroupThreshold> thresholds) {
        this.thresholds = thresholds;
    }
}
