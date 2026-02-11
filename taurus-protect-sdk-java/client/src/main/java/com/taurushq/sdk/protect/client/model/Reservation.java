package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a UTXO reservation in the Taurus Protect system.
 * <p>
 * Reservations are used to lock specific UTXOs (Unspent Transaction Outputs)
 * for UTXO-based blockchains like Bitcoin and Litecoin, preventing
 * double-spending during transaction creation.
 *
 * @see Currency
 */
public class Reservation {

    private String id;
    private String amount;
    private OffsetDateTime creationDate;
    private String kind;
    private String comment;
    private String addressId;
    private String address;
    private Currency currencyInfo;
    private String resourceId;
    private String resourceType;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getAmount() {
        return amount;
    }

    public void setAmount(final String amount) {
        this.amount = amount;
    }

    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(final String comment) {
        this.comment = comment;
    }

    public String getAddressId() {
        return addressId;
    }

    public void setAddressId(final String addressId) {
        this.addressId = addressId;
    }

    public String getAddress() {
        return address;
    }

    public void setAddress(final String address) {
        this.address = address;
    }

    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    public void setCurrencyInfo(final Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    public String getResourceId() {
        return resourceId;
    }

    public void setResourceId(final String resourceId) {
        this.resourceId = resourceId;
    }

    public String getResourceType() {
        return resourceType;
    }

    public void setResourceType(final String resourceType) {
        this.resourceType = resourceType;
    }
}
