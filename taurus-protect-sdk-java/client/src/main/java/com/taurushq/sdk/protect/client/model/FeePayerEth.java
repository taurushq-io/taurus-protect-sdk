package com.taurushq.sdk.protect.client.model;

/**
 * Represents the Ethereum-specific fee payer configuration.
 * <p>
 * Contains details about how the fee payer is configured for
 * EVM-compatible blockchains.
 *
 * @see FeePayerInfo
 */
public class FeePayerEth {

    private String kind;
    private FeePayerEthLocal local;
    private FeePayerEthRemote remote;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public FeePayerEthLocal getLocal() {
        return local;
    }

    public void setLocal(final FeePayerEthLocal local) {
        this.local = local;
    }

    public FeePayerEthRemote getRemote() {
        return remote;
    }

    public void setRemote(final FeePayerEthRemote remote) {
        this.remote = remote;
    }
}
