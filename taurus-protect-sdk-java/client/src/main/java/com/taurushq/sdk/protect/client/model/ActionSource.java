package com.taurushq.sdk.protect.client.model;

/**
 * Represents the source of a transfer in an action task.
 */
public class ActionSource {

    private String kind;
    private String addressId;
    private String walletId;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public String getAddressId() {
        return addressId;
    }

    public void setAddressId(final String addressId) {
        this.addressId = addressId;
    }

    public String getWalletId() {
        return walletId;
    }

    public void setWalletId(final String walletId) {
        this.walletId = walletId;
    }
}
