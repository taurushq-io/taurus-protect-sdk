package com.taurushq.sdk.protect.client.model;

/**
 * Represents the target of an action trigger.
 * <p>
 * A target can be either an address or a wallet.
 */
public class ActionTarget {

    private String kind;
    private TargetAddress address;
    private TargetWallet wallet;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public TargetAddress getAddress() {
        return address;
    }

    public void setAddress(final TargetAddress address) {
        this.address = address;
    }

    public TargetWallet getWallet() {
        return wallet;
    }

    public void setWallet(final TargetWallet wallet) {
        this.wallet = wallet;
    }
}
