package com.taurushq.sdk.protect.client.model;

/**
 * Represents Fantom (FTM) validator information.
 * <p>
 * This model contains information about a Fantom validator node,
 * including its stake amounts and status.
 */
public class FTMValidatorInfo {

    private String validatorId;
    private String address;
    private Boolean active;
    private String totalStake;
    private String selfStake;
    private String deactivatedAtDateUnix;
    private String createdAtDateUnix;

    /**
     * Gets the validator ID.
     *
     * @return the validator ID
     */
    public String getValidatorId() {
        return validatorId;
    }

    /**
     * Sets the validator ID.
     *
     * @param validatorId the validator ID
     */
    public void setValidatorId(String validatorId) {
        this.validatorId = validatorId;
    }

    /**
     * Gets the validator's address.
     *
     * @return the address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the validator's address.
     *
     * @param address the address
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Checks if the validator is currently active.
     *
     * @return true if active, false otherwise
     */
    public Boolean getActive() {
        return active;
    }

    /**
     * Sets the validator's active status.
     *
     * @param active the active status
     */
    public void setActive(Boolean active) {
        this.active = active;
    }

    /**
     * Gets the validator's total stake in wei.
     *
     * @return the total stake
     */
    public String getTotalStake() {
        return totalStake;
    }

    /**
     * Sets the validator's total stake.
     *
     * @param totalStake the total stake
     */
    public void setTotalStake(String totalStake) {
        this.totalStake = totalStake;
    }

    /**
     * Gets the validator's self stake in wei.
     *
     * @return the self stake
     */
    public String getSelfStake() {
        return selfStake;
    }

    /**
     * Sets the validator's self stake.
     *
     * @param selfStake the self stake
     */
    public void setSelfStake(String selfStake) {
        this.selfStake = selfStake;
    }

    /**
     * Gets the Unix timestamp when the validator was deactivated.
     *
     * @return the deactivation timestamp, or null if still active
     */
    public String getDeactivatedAtDateUnix() {
        return deactivatedAtDateUnix;
    }

    /**
     * Sets the deactivation timestamp.
     *
     * @param deactivatedAtDateUnix the deactivation timestamp
     */
    public void setDeactivatedAtDateUnix(String deactivatedAtDateUnix) {
        this.deactivatedAtDateUnix = deactivatedAtDateUnix;
    }

    /**
     * Gets the Unix timestamp when the validator was created.
     *
     * @return the creation timestamp
     */
    public String getCreatedAtDateUnix() {
        return createdAtDateUnix;
    }

    /**
     * Sets the creation timestamp.
     *
     * @param createdAtDateUnix the creation timestamp
     */
    public void setCreatedAtDateUnix(String createdAtDateUnix) {
        this.createdAtDateUnix = createdAtDateUnix;
    }

    @Override
    public String toString() {
        return "FTMValidatorInfo{"
                + "validatorId='" + validatorId + '\''
                + ", address='" + address + '\''
                + ", active=" + active
                + ", totalStake='" + totalStake + '\''
                + ", selfStake='" + selfStake + '\''
                + '}';
    }
}
