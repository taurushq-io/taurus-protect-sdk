package com.taurushq.sdk.protect.client.model;

/**
 * Represents NEAR Protocol validator information.
 * <p>
 * This model contains information about a NEAR validator node,
 * including its staked balance and fee structure.
 */
public class NEARValidatorInfo {

    private String validatorAddress;
    private String ownerId;
    private String totalStakedBalance;
    private Float rewardFeeFraction;
    private String stakingKey;
    private Boolean stakingPaused;

    /**
     * Gets the validator's contract address.
     *
     * @return the validator address
     */
    public String getValidatorAddress() {
        return validatorAddress;
    }

    /**
     * Sets the validator's contract address.
     *
     * @param validatorAddress the validator address
     */
    public void setValidatorAddress(String validatorAddress) {
        this.validatorAddress = validatorAddress;
    }

    /**
     * Gets the owner's account ID.
     *
     * @return the owner ID
     */
    public String getOwnerId() {
        return ownerId;
    }

    /**
     * Sets the owner's account ID.
     *
     * @param ownerId the owner ID
     */
    public void setOwnerId(String ownerId) {
        this.ownerId = ownerId;
    }

    /**
     * Gets the total staked balance in yoctoNEAR.
     *
     * @return the total staked balance
     */
    public String getTotalStakedBalance() {
        return totalStakedBalance;
    }

    /**
     * Sets the total staked balance.
     *
     * @param totalStakedBalance the total staked balance
     */
    public void setTotalStakedBalance(String totalStakedBalance) {
        this.totalStakedBalance = totalStakedBalance;
    }

    /**
     * Gets the reward fee fraction (percentage of rewards taken by validator).
     *
     * @return the fee fraction as a decimal (e.g., 0.10 for 10%)
     */
    public Float getRewardFeeFraction() {
        return rewardFeeFraction;
    }

    /**
     * Sets the reward fee fraction.
     *
     * @param rewardFeeFraction the fee fraction
     */
    public void setRewardFeeFraction(Float rewardFeeFraction) {
        this.rewardFeeFraction = rewardFeeFraction;
    }

    /**
     * Gets the staking public key.
     *
     * @return the staking key
     */
    public String getStakingKey() {
        return stakingKey;
    }

    /**
     * Sets the staking public key.
     *
     * @param stakingKey the staking key
     */
    public void setStakingKey(String stakingKey) {
        this.stakingKey = stakingKey;
    }

    /**
     * Checks if staking is currently paused.
     *
     * @return true if staking is paused, false otherwise
     */
    public Boolean getStakingPaused() {
        return stakingPaused;
    }

    /**
     * Sets the staking paused status.
     *
     * @param stakingPaused the paused status
     */
    public void setStakingPaused(Boolean stakingPaused) {
        this.stakingPaused = stakingPaused;
    }

    @Override
    public String toString() {
        return "NEARValidatorInfo{"
                + "validatorAddress='" + validatorAddress + '\''
                + ", ownerId='" + ownerId + '\''
                + ", totalStakedBalance='" + totalStakedBalance + '\''
                + ", stakingPaused=" + stakingPaused
                + '}';
    }
}
