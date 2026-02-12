package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.math.BigInteger;

/**
 * Represents the balance of a wallet or address in the Taurus Protect system.
 * <p>
 * Balances are tracked in three categories:
 * <ul>
 *   <li><b>Total</b> - The complete balance including reserved amounts</li>
 *   <li><b>Available</b> - The balance available for new transactions</li>
 *   <li><b>Reserved</b> - The balance reserved for pending transactions</li>
 * </ul>
 * <p>
 * Each category has both confirmed (on-chain) and unconfirmed (pending) values.
 * All amounts are in the smallest unit of the currency (e.g., wei for ETH).
 * <p>
 * Example usage:
 * <pre>{@code
 * Balance balance = wallet.getBalance();
 * BigInteger available = balance.getAvailableConfirmed();
 * BigInteger pending = balance.getAvailableUnconfirmed();
 * }</pre>
 *
 * @see Wallet
 * @see Address
 */
public class Balance {

    /**
     * The total confirmed balance (on-chain).
     */
    private BigInteger totalConfirmed;

    /**
     * The total unconfirmed balance (pending transactions).
     */
    private BigInteger totalUnconfirmed;

    /**
     * The available confirmed balance (on-chain, not reserved).
     */
    private BigInteger availableConfirmed;

    /**
     * The available unconfirmed balance (pending, not reserved).
     */
    private BigInteger availableUnconfirmed;

    /**
     * The reserved confirmed balance (held for pending outgoing transactions).
     */
    private BigInteger reservedConfirmed;

    /**
     * The reserved unconfirmed balance.
     */
    private BigInteger reservedUnconfirmed;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets total confirmed.
     *
     * @return the total confirmed
     */
    public BigInteger getTotalConfirmed() {
        return totalConfirmed;
    }

    /**
     * Sets total confirmed.
     *
     * @param totalConfirmed the total confirmed
     */
    public void setTotalConfirmed(BigInteger totalConfirmed) {
        this.totalConfirmed = totalConfirmed;
    }

    /**
     * Gets total unconfirmed.
     *
     * @return the total unconfirmed
     */
    public BigInteger getTotalUnconfirmed() {
        return totalUnconfirmed;
    }

    /**
     * Sets total unconfirmed.
     *
     * @param totalUnconfirmed the total unconfirmed
     */
    public void setTotalUnconfirmed(BigInteger totalUnconfirmed) {
        this.totalUnconfirmed = totalUnconfirmed;
    }

    /**
     * Gets available confirmed.
     *
     * @return the available confirmed
     */
    public BigInteger getAvailableConfirmed() {
        return availableConfirmed;
    }

    /**
     * Sets available confirmed.
     *
     * @param availableConfirmed the available confirmed
     */
    public void setAvailableConfirmed(BigInteger availableConfirmed) {
        this.availableConfirmed = availableConfirmed;
    }

    /**
     * Gets available unconfirmed.
     *
     * @return the available unconfirmed
     */
    public BigInteger getAvailableUnconfirmed() {
        return availableUnconfirmed;
    }

    /**
     * Sets available unconfirmed.
     *
     * @param availableUnconfirmed the available unconfirmed
     */
    public void setAvailableUnconfirmed(BigInteger availableUnconfirmed) {
        this.availableUnconfirmed = availableUnconfirmed;
    }

    /**
     * Gets reserved confirmed.
     *
     * @return the reserved confirmed
     */
    public BigInteger getReservedConfirmed() {
        return reservedConfirmed;
    }

    /**
     * Sets reserved confirmed.
     *
     * @param reservedConfirmed the reserved confirmed
     */
    public void setReservedConfirmed(BigInteger reservedConfirmed) {
        this.reservedConfirmed = reservedConfirmed;
    }

    /**
     * Gets reserved unconfirmed.
     *
     * @return the reserved unconfirmed
     */
    public BigInteger getReservedUnconfirmed() {
        return reservedUnconfirmed;
    }

    /**
     * Sets reserved unconfirmed.
     *
     * @param reservedUnconfirmed the reserved unconfirmed
     */
    public void setReservedUnconfirmed(BigInteger reservedUnconfirmed) {
        this.reservedUnconfirmed = reservedUnconfirmed;
    }
}
