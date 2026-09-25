package com.taurushq.sdk.protect.client.model;

/**
 * An earn reward credited to an address.
 * <p>
 * The Merkl token reward fields are flattened; they are set when the reward type is
 * {@code RewardTypeMerklToken}.
 *
 * @see com.taurushq.sdk.protect.client.service.EarnService
 */
public class EarnReward {

    private String id;
    private String recipientAddressId;
    private String recipientAddress;
    private String rewardType;
    private String amount;
    private String claimed;
    private String pending;
    private String tokenAddress;
    private String tokenSymbol;
    private String tokenAssetId;

    /**
     * Gets the reward id.
     *
     * @return the reward id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the reward id.
     *
     * @param id the reward id
     */
    public void setId(final String id) {
        this.id = id;
    }

    /**
     * Gets the recipient address id.
     *
     * @return the recipient address id
     */
    public String getRecipientAddressId() {
        return recipientAddressId;
    }

    /**
     * Sets the recipient address id.
     *
     * @param recipientAddressId the recipient address id
     */
    public void setRecipientAddressId(final String recipientAddressId) {
        this.recipientAddressId = recipientAddressId;
    }

    /**
     * Gets the recipient blockchain address.
     *
     * @return the recipient blockchain address
     */
    public String getRecipientAddress() {
        return recipientAddress;
    }

    /**
     * Sets the recipient blockchain address.
     *
     * @param recipientAddress the recipient blockchain address
     */
    public void setRecipientAddress(final String recipientAddress) {
        this.recipientAddress = recipientAddress;
    }

    /**
     * Gets the reward type, as the wire value (e.g. RewardTypeMerklToken).
     *
     * @return the reward type, as the wire value (e.g. RewardTypeMerklToken)
     */
    public String getRewardType() {
        return rewardType;
    }

    /**
     * Sets the reward type, as the wire value (e.g. RewardTypeMerklToken).
     *
     * @param rewardType the reward type, as the wire value (e.g. RewardTypeMerklToken)
     */
    public void setRewardType(final String rewardType) {
        this.rewardType = rewardType;
    }

    /**
     * Gets the reward amount.
     *
     * @return the reward amount
     */
    public String getAmount() {
        return amount;
    }

    /**
     * Sets the reward amount.
     *
     * @param amount the reward amount
     */
    public void setAmount(final String amount) {
        this.amount = amount;
    }

    /**
     * Gets the amount already claimed.
     *
     * @return the amount already claimed
     */
    public String getClaimed() {
        return claimed;
    }

    /**
     * Sets the amount already claimed.
     *
     * @param claimed the amount already claimed
     */
    public void setClaimed(final String claimed) {
        this.claimed = claimed;
    }

    /**
     * Gets the amount pending.
     *
     * @return the amount pending
     */
    public String getPending() {
        return pending;
    }

    /**
     * Sets the amount pending.
     *
     * @param pending the amount pending
     */
    public void setPending(final String pending) {
        this.pending = pending;
    }

    /**
     * Gets the reward token contract address.
     *
     * @return the reward token contract address
     */
    public String getTokenAddress() {
        return tokenAddress;
    }

    /**
     * Sets the reward token contract address.
     *
     * @param tokenAddress the reward token contract address
     */
    public void setTokenAddress(final String tokenAddress) {
        this.tokenAddress = tokenAddress;
    }

    /**
     * Gets the reward token symbol.
     *
     * @return the reward token symbol
     */
    public String getTokenSymbol() {
        return tokenSymbol;
    }

    /**
     * Sets the reward token symbol.
     *
     * @param tokenSymbol the reward token symbol
     */
    public void setTokenSymbol(final String tokenSymbol) {
        this.tokenSymbol = tokenSymbol;
    }

    /**
     * Gets the reward token asset id.
     *
     * @return the reward token asset id
     */
    public String getTokenAssetId() {
        return tokenAssetId;
    }

    /**
     * Sets the reward token asset id.
     *
     * @param tokenAssetId the reward token asset id
     */
    public void setTokenAssetId(final String tokenAssetId) {
        this.tokenAssetId = tokenAssetId;
    }
}
