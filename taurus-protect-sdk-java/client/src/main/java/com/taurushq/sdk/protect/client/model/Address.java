package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a blockchain address within a wallet in the Taurus Protect system.
 * <p>
 * An address is a unique identifier on a blockchain that can receive and send
 * cryptocurrency. Each address belongs to a specific wallet and has an associated
 * balance, label, and optional compliance scores.
 * <p>
 * Addresses are cryptographically signed by the Taurus Protect system to ensure
 * integrity. The signature can be verified using the rules container public keys.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a new address in a wallet
 * Address address = client.getAddressService().createAddress(walletId, "My Address", "Customer deposit");
 * System.out.println("Address: " + address.getAddress());
 * System.out.println("Balance: " + address.getBalance().getAvailable());
 * }</pre>
 *
 * @see Wallet
 * @see AddressService
 * @see Score
 */
public class Address {

    /**
     * The unique identifier of the address assigned by Taurus Protect.
     */
    private long id;

    /**
     * The ID of the wallet that contains this address.
     */
    private long walletId;

    /**
     * Indicates whether the address is disabled.
     */
    private boolean disabled;

    /**
     * The currency code for the address (e.g., "ETH", "BTC").
     */
    private String currency;

    /**
     * Detailed information about the address's currency.
     */
    private Currency currencyInfo;

    /**
     * The hierarchical deterministic (HD) derivation path for this address.
     */
    private String addressPath;

    /**
     * The blockchain address string (e.g., "0x..." for Ethereum).
     */
    private String address;

    /**
     * An optional comment or description for the address.
     */
    private String comment;

    /**
     * A human-readable label for the address.
     */
    private String label;

    /**
     * The status of address creation, as reported by the API.
     * <p>
     * One of {@code created}, {@code creating}, {@code signed}, {@code observed} or
     * {@code confirmed}. {@code creating} is used for asynchronous address creation:
     * the row exists but the blockchain address string is still empty and the HSM has
     * not signed anything yet.
     * <p>
     * This field is what lets the SDK tell "not signed yet" apart from "signature
     * stripped": both arrive with no signature, but only the first has an empty address
     * string, and only the second is an attack. It is a non-security field taken from
     * the DTO — never rely on it to decide whether a signature was checked.
     */
    private String status;

    /**
     * The cryptographic signature of the address for integrity verification.
     */
    private String signature;

    /**
     * The date and time when the address was created.
     */
    private OffsetDateTime creationDate;

    /**
     * The date and time when the address was last updated.
     */
    private OffsetDateTime updateDate;

    /**
     * Information about the wallet containing this address.
     */
    private Wallet walletInfo;

    /**
     * The current balance information for the address.
     */
    private Balance balance;

    /**
     * Compliance scores from various risk assessment providers.
     */
    private List<Score> scores;

    /**
     * Custom key-value attributes associated with the address.
     */
    private List<Attribute> attributes;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the address.
     *
     * @return the address ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the address.
     *
     * @param id the address ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the ID of the wallet that contains this address.
     *
     * @return the wallet ID
     */
    public long getWalletId() {
        return walletId;
    }

    /**
     * Sets the ID of the wallet that contains this address.
     *
     * @param walletId the wallet ID to set
     */
    public void setWalletId(long walletId) {
        this.walletId = walletId;
    }

    /**
     * Returns whether the address is disabled.
     * <p>
     * Disabled addresses cannot be used for new transactions.
     *
     * @return {@code true} if the address is disabled, {@code false} otherwise
     */
    public boolean isDisabled() {
        return disabled;
    }

    /**
     * Sets whether the address is disabled.
     *
     * @param disabled {@code true} to disable the address, {@code false} to enable it
     */
    public void setDisabled(boolean disabled) {
        this.disabled = disabled;
    }

    /**
     * Returns the currency code for the address (e.g., "ETH", "BTC").
     *
     * @return the currency code
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code for the address.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns detailed information about the address's currency.
     *
     * @return the currency information, or {@code null} if not available
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the detailed currency information for the address.
     *
     * @param currencyInfo the currency information to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Returns the hierarchical deterministic (HD) derivation path for this address.
     *
     * @return the address derivation path
     */
    public String getAddressPath() {
        return addressPath;
    }

    /**
     * Sets the hierarchical deterministic (HD) derivation path for this address.
     *
     * @param addressPath the address derivation path to set
     */
    public void setAddressPath(String addressPath) {
        this.addressPath = addressPath;
    }

    /**
     * Returns the blockchain address string (e.g., "0x..." for Ethereum).
     *
     * @return the blockchain address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the blockchain address string.
     *
     * @param address the blockchain address to set
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Returns the comment or description for the address.
     *
     * @return the comment, or {@code null} if not set
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the comment or description for the address.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }

    /**
     * Returns the human-readable label for the address.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the human-readable label for the address.
     *
     * @param label the label to set
     */
    public void setLabel(String label) {
        this.label = label;
    }

    /**
     * Returns the status of address creation.
     * <p>
     * One of {@code created}, {@code creating}, {@code signed}, {@code observed} or
     * {@code confirmed}; {@code creating} means asynchronous creation is still in
     * flight, so {@link #getAddress()} is empty and no signature exists yet.
     *
     * @return the creation status, or {@code null} if the API did not report one
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the status of address creation.
     * <p>
     * The setter is named {@code setStatus} on purpose: MapStruct follows JavaBean
     * conventions and silently skips a property whose setter is not {@code set*},
     * which is how {@code Address.disabled} and {@code Wallet.omnibus} were once
     * dropped without any build failure.
     *
     * @param status the creation status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Returns the cryptographic signature of the address for integrity verification.
     * <p>
     * This signature is generated by the Taurus Protect system and can be verified
     * using the rules container public keys.
     *
     * @return the address signature
     */
    public String getSignature() {
        return signature;
    }

    /**
     * Sets the cryptographic signature of the address.
     *
     * @param signature the signature to set
     */
    public void setSignature(String signature) {
        this.signature = signature;
    }

    /**
     * Returns the date and time when the address was created.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the date and time when the address was created.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Returns the date and time when the address was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the date and time when the address was last updated.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Returns information about the wallet containing this address.
     *
     * @return the wallet information, or {@code null} if not populated
     */
    public Wallet getWalletInfo() {
        return walletInfo;
    }

    /**
     * Sets the wallet information for this address.
     *
     * @param walletInfo the wallet information to set
     */
    public void setWalletInfo(Wallet walletInfo) {
        this.walletInfo = walletInfo;
    }

    /**
     * Returns the current balance information for the address.
     *
     * @return the balance, or {@code null} if not available
     */
    public Balance getBalance() {
        return balance;
    }

    /**
     * Sets the balance information for the address.
     *
     * @param balance the balance to set
     */
    public void setBalance(Balance balance) {
        this.balance = balance;
    }

    /**
     * Returns the compliance scores from various risk assessment providers.
     * <p>
     * Scores are provided by third-party compliance services such as Chainalysis,
     * Coinfirm, or Elliptic.
     *
     * @return the list of scores, or {@code null} if not available
     */
    public List<Score> getScores() {
        return scores;
    }

    /**
     * Sets the compliance scores for the address.
     *
     * @param scores the list of scores to set
     */
    public void setScores(List<Score> scores) {
        this.scores = scores;
    }

    /**
     * Returns the custom key-value attributes associated with the address.
     *
     * @return the list of attributes, or {@code null} if none are set
     */
    public List<Attribute> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom key-value attributes for the address.
     *
     * @param attributes the list of attributes to set
     */
    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }
}
