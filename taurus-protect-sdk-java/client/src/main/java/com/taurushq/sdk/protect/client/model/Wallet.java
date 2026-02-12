package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a blockchain wallet in the Taurus Protect system.
 * <p>
 * A wallet is a container for one or more addresses on a specific blockchain
 * and network. Wallets can be either standard (single customer) or omnibus
 * (pooled funds from multiple customers).
 * <p>
 * Example usage:
 * <pre>{@code
 * // Retrieve a wallet by ID
 * Wallet wallet = client.getWalletService().getWallet(walletId);
 * System.out.println("Wallet: " + wallet.getName());
 * System.out.println("Blockchain: " + wallet.getBlockchain());
 * System.out.println("Balance: " + wallet.getBalance().getAvailable());
 * }</pre>
 *
 * @see Address
 * @see WalletService
 * @see Balance
 */
public class Wallet {

    /**
     * The unique identifier of the wallet assigned by Taurus Protect.
     */
    private long id;

    /**
     * The human-readable name of the wallet.
     */
    private String name;

    /**
     * The currency code for the wallet (e.g., "ETH", "BTC").
     */
    private String currency;

    /**
     * The current balance information for the wallet.
     */
    private Balance balance;

    /**
     * The hierarchical deterministic (HD) account path for key derivation.
     */
    private String accountPath;

    /**
     * Indicates whether this is an omnibus wallet (pooling funds from multiple customers).
     */
    private boolean isOmnibus;

    /**
     * The date and time when the wallet was created.
     */
    private OffsetDateTime creationDate;

    /**
     * The date and time when the wallet was last updated.
     */
    private OffsetDateTime updateDate;

    /**
     * An optional customer identifier associated with the wallet.
     */
    private String customerId;

    /**
     * An optional comment or description for the wallet.
     */
    private String comment;

    /**
     * Indicates whether the wallet is disabled.
     */
    private boolean disabled;

    /**
     * The blockchain type (e.g., "ETH", "BTC", "SOL").
     */
    private String blockchain;

    /**
     * The number of addresses associated with this wallet.
     */
    private int addressesCount;

    /**
     * Detailed information about the wallet's currency.
     */
    private Currency currencyInfo;

    /**
     * Custom key-value attributes associated with the wallet.
     */
    private List<Attribute> attributes;

    /**
     * The network identifier (e.g., "mainnet", "testnet").
     */
    private String network;

    /**
     * An optional external identifier for the wallet.
     */
    private String externalWalletId;

    /**
     * The visibility group ID for access control.
     */
    private String visibilityGroupID;


    /**
     * Returns the unique identifier of the wallet.
     *
     * @return the wallet ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the wallet.
     *
     * @param id the wallet ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the human-readable name of the wallet.
     *
     * @return the wallet name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the human-readable name of the wallet.
     *
     * @param name the wallet name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Returns the currency code for the wallet (e.g., "ETH", "BTC").
     *
     * @return the currency code
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code for the wallet.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns the current balance information for the wallet.
     *
     * @return the balance, or {@code null} if not available
     */
    public Balance getBalance() {
        return balance;
    }

    /**
     * Sets the balance information for the wallet.
     *
     * @param balance the balance to set
     */
    public void setBalance(Balance balance) {
        this.balance = balance;
    }

    /**
     * Returns the hierarchical deterministic (HD) account path for key derivation.
     *
     * @return the account path
     */
    public String getAccountPath() {
        return accountPath;
    }

    /**
     * Sets the hierarchical deterministic (HD) account path.
     *
     * @param accountPath the account path to set
     */
    public void setAccountPath(String accountPath) {
        this.accountPath = accountPath;
    }

    /**
     * Returns whether this is an omnibus wallet.
     * <p>
     * An omnibus wallet pools funds from multiple customers under a single wallet.
     *
     * @return {@code true} if this is an omnibus wallet, {@code false} otherwise
     */
    public boolean isOmnibus() {
        return isOmnibus;
    }

    /**
     * Sets whether this is an omnibus wallet.
     *
     * @param omnibus {@code true} for an omnibus wallet, {@code false} otherwise
     */
    public void setOmnibus(boolean omnibus) {
        isOmnibus = omnibus;
    }

    /**
     * Returns the date and time when the wallet was created.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the date and time when the wallet was created.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Returns the date and time when the wallet was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the date and time when the wallet was last updated.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Returns the customer identifier associated with the wallet.
     *
     * @return the customer ID, or {@code null} if not set
     */
    public String getCustomerId() {
        return customerId;
    }

    /**
     * Sets the customer identifier associated with the wallet.
     *
     * @param customerId the customer ID to set
     */
    public void setCustomerId(String customerId) {
        this.customerId = customerId;
    }

    /**
     * Returns the comment or description for the wallet.
     *
     * @return the comment, or {@code null} if not set
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the comment or description for the wallet.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }

    /**
     * Returns whether the wallet is disabled.
     * <p>
     * Disabled wallets cannot be used for new transactions.
     *
     * @return {@code true} if the wallet is disabled, {@code false} otherwise
     */
    public boolean isDisabled() {
        return disabled;
    }

    /**
     * Sets whether the wallet is disabled.
     *
     * @param disabled {@code true} to disable the wallet, {@code false} to enable it
     */
    public void setDisabled(boolean disabled) {
        this.disabled = disabled;
    }

    /**
     * Returns the blockchain type (e.g., "ETH", "BTC", "SOL").
     *
     * @return the blockchain identifier
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain type.
     *
     * @param blockchain the blockchain identifier to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Returns the number of addresses associated with this wallet.
     *
     * @return the address count
     */
    public int getAddressesCount() {
        return addressesCount;
    }

    /**
     * Sets the number of addresses associated with this wallet.
     *
     * @param addressesCount the address count to set
     */
    public void setAddressesCount(int addressesCount) {
        this.addressesCount = addressesCount;
    }

    /**
     * Returns detailed information about the wallet's currency.
     *
     * @return the currency information, or {@code null} if not available
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the detailed currency information for the wallet.
     *
     * @param currencyInfo the currency information to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Returns the custom key-value attributes associated with the wallet.
     *
     * @return the list of attributes, or {@code null} if none are set
     */
    public List<Attribute> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom key-value attributes for the wallet.
     *
     * @param attributes the list of attributes to set
     */
    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }

    /**
     * Returns the network identifier (e.g., "mainnet", "testnet").
     *
     * @return the network identifier
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network identifier.
     *
     * @param network the network identifier to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Returns the external identifier for the wallet.
     *
     * @return the external wallet ID, or {@code null} if not set
     */
    public String getExternalWalletId() {
        return externalWalletId;
    }

    /**
     * Sets the external identifier for the wallet.
     *
     * @param externalWalletId the external wallet ID to set
     */
    public void setExternalWalletId(String externalWalletId) {
        this.externalWalletId = externalWalletId;
    }

    /**
     * Returns the visibility group ID for access control.
     *
     * @return the visibility group ID, or {@code null} if not set
     */
    public String getVisibilityGroupID() {
        return visibilityGroupID;
    }

    /**
     * Sets the visibility group ID for access control.
     *
     * @param visibilityGroupID the visibility group ID to set
     */
    public void setVisibilityGroupID(String visibilityGroupID) {
        this.visibilityGroupID = visibilityGroupID;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

}


