package com.taurushq.sdk.protect.client.model;

/**
 * Represents Tezos (XTZ) staking rewards information.
 * <p>
 * This model contains information about staking rewards received
 * by a Tezos address over a specified time period.
 */
public class XTZStakingRewards {

    private String receivedRewardsAmount;

    /**
     * Gets the total rewards received in mutez.
     *
     * @return the received rewards amount
     */
    public String getReceivedRewardsAmount() {
        return receivedRewardsAmount;
    }

    /**
     * Sets the received rewards amount.
     *
     * @param receivedRewardsAmount the received rewards amount
     */
    public void setReceivedRewardsAmount(String receivedRewardsAmount) {
        this.receivedRewardsAmount = receivedRewardsAmount;
    }

    @Override
    public String toString() {
        return "XTZStakingRewards{"
                + "receivedRewardsAmount='" + receivedRewardsAmount + '\''
                + '}';
    }
}
