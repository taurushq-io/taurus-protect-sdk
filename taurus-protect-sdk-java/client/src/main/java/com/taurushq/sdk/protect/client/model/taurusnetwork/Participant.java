package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a participant in the Taurus Network.
 */
public class Participant {

    private String id;
    private String name;
    private String legalAddress;
    private String country;
    private String publicKey;
    private String shield;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;
    private String outgoingTotalPledgesValuationBaseCurrency;
    private String incomingTotalPledgesValuationBaseCurrency;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(final String name) {
        this.name = name;
    }

    public String getLegalAddress() {
        return legalAddress;
    }

    public void setLegalAddress(final String legalAddress) {
        this.legalAddress = legalAddress;
    }

    public String getCountry() {
        return country;
    }

    public void setCountry(final String country) {
        this.country = country;
    }

    public String getPublicKey() {
        return publicKey;
    }

    public void setPublicKey(final String publicKey) {
        this.publicKey = publicKey;
    }

    public String getShield() {
        return shield;
    }

    public void setShield(final String shield) {
        this.shield = shield;
    }

    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(final OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(final OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    public String getOutgoingTotalPledgesValuationBaseCurrency() {
        return outgoingTotalPledgesValuationBaseCurrency;
    }

    public void setOutgoingTotalPledgesValuationBaseCurrency(final String value) {
        this.outgoingTotalPledgesValuationBaseCurrency = value;
    }

    public String getIncomingTotalPledgesValuationBaseCurrency() {
        return incomingTotalPledgesValuationBaseCurrency;
    }

    public void setIncomingTotalPledgesValuationBaseCurrency(final String value) {
        this.incomingTotalPledgesValuationBaseCurrency = value;
    }
}
