package com.taurushq.sdk.protect.client.model;

/**
 * Represents a remote (externally managed) Ethereum fee payer.
 *
 * @see FeePayerEth
 */
public class FeePayerEthRemote {

    private String url;
    private String username;
    private String fromAddressId;
    private String forwarderAddress;
    private String forwarderAddressId;
    private String creatorAddress;
    private String creatorAddressId;
    private String forwarderKind;

    public String getUrl() {
        return url;
    }

    public void setUrl(final String url) {
        this.url = url;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(final String username) {
        this.username = username;
    }

    public String getFromAddressId() {
        return fromAddressId;
    }

    public void setFromAddressId(final String fromAddressId) {
        this.fromAddressId = fromAddressId;
    }

    public String getForwarderAddress() {
        return forwarderAddress;
    }

    public void setForwarderAddress(final String forwarderAddress) {
        this.forwarderAddress = forwarderAddress;
    }

    public String getForwarderAddressId() {
        return forwarderAddressId;
    }

    public void setForwarderAddressId(final String forwarderAddressId) {
        this.forwarderAddressId = forwarderAddressId;
    }

    public String getCreatorAddress() {
        return creatorAddress;
    }

    public void setCreatorAddress(final String creatorAddress) {
        this.creatorAddress = creatorAddress;
    }

    public String getCreatorAddressId() {
        return creatorAddressId;
    }

    public void setCreatorAddressId(final String creatorAddressId) {
        this.creatorAddressId = creatorAddressId;
    }

    public String getForwarderKind() {
        return forwarderKind;
    }

    public void setForwarderKind(final String forwarderKind) {
        this.forwarderKind = forwarderKind;
    }
}
