package com.taurushq.sdk.protect.client.model;

/**
 * Represents the blockchain-specific information for a fee payer.
 * <p>
 * Contains the blockchain type and the associated configuration details.
 *
 * @see FeePayer
 */
public class FeePayerInfo {

    private String blockchain;
    private FeePayerEth eth;

    public String getBlockchain() {
        return blockchain;
    }

    public void setBlockchain(final String blockchain) {
        this.blockchain = blockchain;
    }

    public FeePayerEth getEth() {
        return eth;
    }

    public void setEth(final FeePayerEth eth) {
        this.eth = eth;
    }
}
