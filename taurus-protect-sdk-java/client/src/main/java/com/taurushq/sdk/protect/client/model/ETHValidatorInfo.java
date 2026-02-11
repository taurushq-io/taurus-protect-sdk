package com.taurushq.sdk.protect.client.model;

/**
 * Represents Ethereum validator information.
 * <p>
 * This model contains information about an Ethereum 2.0 validator,
 * including its public key, balance, and status.
 */
public class ETHValidatorInfo {

    private String id;
    private String pubkey;
    private String status;
    private String balance;
    private String network;
    private String provider;
    private String addressId;

    /**
     * Gets the validator ID.
     *
     * @return the validator ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the validator ID.
     *
     * @param id the validator ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the validator's public key.
     *
     * @return the public key
     */
    public String getPubkey() {
        return pubkey;
    }

    /**
     * Sets the validator's public key.
     *
     * @param pubkey the public key
     */
    public void setPubkey(String pubkey) {
        this.pubkey = pubkey;
    }

    /**
     * Gets the validator's status (e.g., "active_ongoing", "pending_queued").
     *
     * @return the status
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the validator's status.
     *
     * @param status the status
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Gets the validator's current balance in Gwei.
     *
     * @return the balance
     */
    public String getBalance() {
        return balance;
    }

    /**
     * Sets the validator's balance.
     *
     * @param balance the balance
     */
    public void setBalance(String balance) {
        this.balance = balance;
    }

    /**
     * Gets the network (e.g., "mainnet", "goerli").
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network.
     *
     * @param network the network
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets the staking provider.
     *
     * @return the provider
     */
    public String getProvider() {
        return provider;
    }

    /**
     * Sets the staking provider.
     *
     * @param provider the provider
     */
    public void setProvider(String provider) {
        this.provider = provider;
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

    @Override
    public String toString() {
        return "ETHValidatorInfo{"
                + "id='" + id + '\''
                + ", pubkey='" + pubkey + '\''
                + ", balance='" + balance + '\''
                + ", status='" + status + '\''
                + ", network='" + network + '\''
                + '}';
    }
}
