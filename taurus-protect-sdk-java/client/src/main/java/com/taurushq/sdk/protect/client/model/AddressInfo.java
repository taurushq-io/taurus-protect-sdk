package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.math.BigInteger;
import java.util.List;

/**
 * Represents address information within a transaction or request.
 * <p>
 * AddressInfo provides details about source or destination addresses in transactions,
 * including both internal (managed by Taurus Protect) and external (whitelisted) addresses.
 *
 * @see Transaction
 * @see Request
 * @see Score
 */
public class AddressInfo {

    /**
     * The blockchain address string.
     */
    private String address;

    /**
     * Human-readable label for the address.
     */
    private String label;

    /**
     * The container or wallet path this address belongs to.
     */
    private String container;

    /**
     * Customer identifier associated with this address.
     */
    private String customerId;

    /**
     * The transaction amount in smallest units (e.g., wei, satoshi).
     */
    private BigInteger amount;

    /**
     * The transaction amount in main units (e.g., ETH, BTC).
     */
    private double amountMainUnit;

    /**
     * The type of address (e.g., "source", "destination").
     */
    private String type;

    /**
     * The derivation index for this address within its wallet.
     */
    private int idx;

    /**
     * The internal address ID if this is a managed address.
     */
    private long internalAddressId;

    /**
     * The whitelisted address ID if this is an external whitelisted address.
     */
    private long whitelistedAddressId;

    /**
     * Risk assessment scores for this address (AML/KYC compliance).
     */
    private List<Score> scores;

    /**
     * Returns the blockchain address string.
     *
     * @return the address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the blockchain address string.
     *
     * @param address the address to set
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Returns the human-readable label for this address.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the human-readable label for this address.
     *
     * @param label the label to set
     */
    public void setLabel(String label) {
        this.label = label;
    }

    /**
     * Returns the container or wallet path this address belongs to.
     *
     * @return the container path
     */
    public String getContainer() {
        return container;
    }

    /**
     * Sets the container or wallet path this address belongs to.
     *
     * @param container the container path to set
     */
    public void setContainer(String container) {
        this.container = container;
    }

    /**
     * Returns the customer identifier associated with this address.
     *
     * @return the customer ID
     */
    public String getCustomerId() {
        return customerId;
    }

    /**
     * Sets the customer identifier associated with this address.
     *
     * @param customerId the customer ID to set
     */
    public void setCustomerId(String customerId) {
        this.customerId = customerId;
    }

    /**
     * Returns the transaction amount in smallest units.
     *
     * @return the amount in smallest units (e.g., wei, satoshi)
     */
    public BigInteger getAmount() {
        return amount;
    }

    /**
     * Sets the transaction amount in smallest units.
     *
     * @param amount the amount to set
     */
    public void setAmount(BigInteger amount) {
        this.amount = amount;
    }

    /**
     * Returns the transaction amount in main units.
     *
     * @return the amount in main units (e.g., ETH, BTC)
     */
    public double getAmountMainUnit() {
        return amountMainUnit;
    }

    /**
     * Sets the transaction amount in main units.
     *
     * @param amountMainUnit the amount to set
     */
    public void setAmountMainUnit(double amountMainUnit) {
        this.amountMainUnit = amountMainUnit;
    }

    /**
     * Returns the type of address in the transaction context.
     *
     * @return the type (e.g., "source", "destination")
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the type of address in the transaction context.
     *
     * @param type the type to set
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Returns the derivation index for this address.
     *
     * @return the derivation index
     */
    public int getIdx() {
        return idx;
    }

    /**
     * Sets the derivation index for this address.
     *
     * @param idx the derivation index to set
     */
    public void setIdx(int idx) {
        this.idx = idx;
    }

    /**
     * Returns the internal address ID if this is a managed address.
     *
     * @return the internal address ID, or 0 if not applicable
     */
    public long getInternalAddressId() {
        return internalAddressId;
    }

    /**
     * Sets the internal address ID.
     *
     * @param internalAddressId the internal address ID to set
     */
    public void setInternalAddressId(long internalAddressId) {
        this.internalAddressId = internalAddressId;
    }

    /**
     * Returns the whitelisted address ID if this is an external address.
     *
     * @return the whitelisted address ID, or 0 if not applicable
     */
    public long getWhitelistedAddressId() {
        return whitelistedAddressId;
    }

    /**
     * Sets the whitelisted address ID.
     *
     * @param whitelistedAddressId the whitelisted address ID to set
     */
    public void setWhitelistedAddressId(long whitelistedAddressId) {
        this.whitelistedAddressId = whitelistedAddressId;
    }

    /**
     * Returns the risk assessment scores for this address.
     *
     * @return the list of AML/KYC compliance scores
     */
    public List<Score> getScores() {
        return scores;
    }

    /**
     * Sets the risk assessment scores for this address.
     *
     * @param scores the list of scores to set
     */
    public void setScores(List<Score> scores) {
        this.scores = scores;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }
}
