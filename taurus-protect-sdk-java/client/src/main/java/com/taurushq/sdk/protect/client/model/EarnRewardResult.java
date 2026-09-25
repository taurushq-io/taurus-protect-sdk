package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of earn rewards from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.EarnService
 */
public class EarnRewardResult extends CursorPagedResult {

    private List<EarnReward> rewards;

    /**
     * Gets the earn rewards of this page.
     *
     * @return the rewards
     */
    public List<EarnReward> getRewards() {
        return rewards;
    }

    /**
     * Sets the earn rewards of this page.
     *
     * @param rewards the rewards
     */
    public void setRewards(final List<EarnReward> rewards) {
        this.rewards = rewards;
    }
}
