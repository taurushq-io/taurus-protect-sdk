package com.taurushq.sdk.protect.client.model;

/**
 * Represents an amount threshold for action triggers.
 */
public class ActionAmount {

    private String kind;
    private String cryptoAmount;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public String getCryptoAmount() {
        return cryptoAmount;
    }

    public void setCryptoAmount(final String cryptoAmount) {
        this.cryptoAmount = cryptoAmount;
    }
}
