package com.taurushq.sdk.protect.client.model;

/**
 * Represents the withdrawal fee for an exchange transfer.
 * <p>
 * Contains the fee amount that will be charged for withdrawing
 * assets from an exchange to an external address.
 *
 * @see ExchangeService
 */
public class ExchangeWithdrawalFee {

    private String fee;

    /**
     * Gets the withdrawal fee amount.
     *
     * @return the fee in the smallest currency unit
     */
    public String getFee() {
        return fee;
    }

    /**
     * Sets the withdrawal fee amount.
     *
     * @param fee the fee to set
     */
    public void setFee(String fee) {
        this.fee = fee;
    }
}
