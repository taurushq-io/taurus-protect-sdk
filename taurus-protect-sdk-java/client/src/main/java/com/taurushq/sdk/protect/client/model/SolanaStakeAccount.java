package com.taurushq.sdk.protect.client.model;

/**
 * Represents Solana-specific stake account details.
 * <p>
 * This model contains the state and delegation information for a
 * Solana stake account, including its staked balance and validator.
 */
public class SolanaStakeAccount {

    private String derivationIndex;
    private String state;
    private String validatorAddress;
    private String activeBalance;
    private String inactiveBalance;
    private Boolean allowMerge;

    /**
     * Gets the derivation index used to generate this stake account.
     *
     * @return the derivation index
     */
    public String getDerivationIndex() {
        return derivationIndex;
    }

    /**
     * Sets the derivation index.
     *
     * @param derivationIndex the derivation index
     */
    public void setDerivationIndex(String derivationIndex) {
        this.derivationIndex = derivationIndex;
    }

    /**
     * Gets the stake account state (e.g., "ACTIVE", "INACTIVE", "ACTIVATING", "DEACTIVATING").
     *
     * @return the state
     */
    public String getState() {
        return state;
    }

    /**
     * Sets the stake account state.
     *
     * @param state the state
     */
    public void setState(String state) {
        this.state = state;
    }

    /**
     * Gets the delegated validator's address.
     *
     * @return the validator address
     */
    public String getValidatorAddress() {
        return validatorAddress;
    }

    /**
     * Sets the validator address.
     *
     * @param validatorAddress the validator address
     */
    public void setValidatorAddress(String validatorAddress) {
        this.validatorAddress = validatorAddress;
    }

    /**
     * Gets the active stake balance in lamports.
     *
     * @return the active balance
     */
    public String getActiveBalance() {
        return activeBalance;
    }

    /**
     * Sets the active balance.
     *
     * @param activeBalance the active balance
     */
    public void setActiveBalance(String activeBalance) {
        this.activeBalance = activeBalance;
    }

    /**
     * Gets the inactive stake balance in lamports.
     *
     * @return the inactive balance
     */
    public String getInactiveBalance() {
        return inactiveBalance;
    }

    /**
     * Sets the inactive balance.
     *
     * @param inactiveBalance the inactive balance
     */
    public void setInactiveBalance(String inactiveBalance) {
        this.inactiveBalance = inactiveBalance;
    }

    /**
     * Checks if this stake account can be merged with others.
     *
     * @return true if merge is allowed, false otherwise
     */
    public Boolean getAllowMerge() {
        return allowMerge;
    }

    /**
     * Sets whether merge is allowed.
     *
     * @param allowMerge the allow merge flag
     */
    public void setAllowMerge(Boolean allowMerge) {
        this.allowMerge = allowMerge;
    }

    @Override
    public String toString() {
        return "SolanaStakeAccount{"
                + "derivationIndex='" + derivationIndex + '\''
                + ", state='" + state + '\''
                + ", validatorAddress='" + validatorAddress + '\''
                + ", activeBalance='" + activeBalance + '\''
                + ", inactiveBalance='" + inactiveBalance + '\''
                + '}';
    }
}
