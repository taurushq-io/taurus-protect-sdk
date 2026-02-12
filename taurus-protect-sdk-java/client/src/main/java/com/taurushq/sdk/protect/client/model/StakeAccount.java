package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a stake account in the Taurus Protect system.
 * <p>
 * Stake accounts are used to track staked assets, particularly on
 * blockchains like Solana where staking uses separate stake accounts.
 */
public class StakeAccount {

    private String id;
    private String addressId;
    private String accountAddress;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;
    private String updatedAtBlock;
    private String accountType;
    private SolanaStakeAccount solanaStakeAccount;

    /**
     * Gets the stake account ID.
     *
     * @return the account ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the stake account ID.
     *
     * @param id the account ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the associated address ID.
     *
     * @return the address ID
     */
    public String getAddressId() {
        return addressId;
    }

    /**
     * Sets the address ID.
     *
     * @param addressId the address ID
     */
    public void setAddressId(String addressId) {
        this.addressId = addressId;
    }

    /**
     * Gets the stake account's on-chain address.
     *
     * @return the account address
     */
    public String getAccountAddress() {
        return accountAddress;
    }

    /**
     * Sets the account address.
     *
     * @param accountAddress the account address
     */
    public void setAccountAddress(String accountAddress) {
        this.accountAddress = accountAddress;
    }

    /**
     * Gets the creation timestamp.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the creation timestamp.
     *
     * @param createdAt the creation timestamp
     */
    public void setCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    /**
     * Gets the last update timestamp.
     *
     * @return the update timestamp
     */
    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    /**
     * Sets the update timestamp.
     *
     * @param updatedAt the update timestamp
     */
    public void setUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    /**
     * Gets the block number at which this account was last updated.
     *
     * @return the block number
     */
    public String getUpdatedAtBlock() {
        return updatedAtBlock;
    }

    /**
     * Sets the block number.
     *
     * @param updatedAtBlock the block number
     */
    public void setUpdatedAtBlock(String updatedAtBlock) {
        this.updatedAtBlock = updatedAtBlock;
    }

    /**
     * Gets the account type (e.g., "SOLANA_STAKE").
     *
     * @return the account type
     */
    public String getAccountType() {
        return accountType;
    }

    /**
     * Sets the account type.
     *
     * @param accountType the account type
     */
    public void setAccountType(String accountType) {
        this.accountType = accountType;
    }

    /**
     * Gets Solana-specific stake account details.
     *
     * @return the Solana stake account info, or null if not applicable
     */
    public SolanaStakeAccount getSolanaStakeAccount() {
        return solanaStakeAccount;
    }

    /**
     * Sets the Solana stake account details.
     *
     * @param solanaStakeAccount the Solana stake account info
     */
    public void setSolanaStakeAccount(SolanaStakeAccount solanaStakeAccount) {
        this.solanaStakeAccount = solanaStakeAccount;
    }

    @Override
    public String toString() {
        return "StakeAccount{"
                + "id='" + id + '\''
                + ", addressId='" + addressId + '\''
                + ", accountAddress='" + accountAddress + '\''
                + ", accountType='" + accountType + '\''
                + '}';
    }
}
