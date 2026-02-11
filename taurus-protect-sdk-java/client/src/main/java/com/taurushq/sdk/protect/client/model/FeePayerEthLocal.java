package com.taurushq.sdk.protect.client.model;

/**
 * Represents a local (internally managed) Ethereum fee payer.
 *
 * @see FeePayerEth
 */
public class FeePayerEthLocal {

    private String addressId;
    private String forwarderAddressId;
    private Boolean autoApprove;
    private String creatorAddressId;
    private String forwarderKind;

    public String getAddressId() {
        return addressId;
    }

    public void setAddressId(final String addressId) {
        this.addressId = addressId;
    }

    public String getForwarderAddressId() {
        return forwarderAddressId;
    }

    public void setForwarderAddressId(final String forwarderAddressId) {
        this.forwarderAddressId = forwarderAddressId;
    }

    public Boolean getAutoApprove() {
        return autoApprove;
    }

    public void setAutoApprove(final Boolean autoApprove) {
        this.autoApprove = autoApprove;
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
