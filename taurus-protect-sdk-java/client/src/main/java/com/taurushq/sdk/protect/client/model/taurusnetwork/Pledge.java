package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a pledge in the Taurus Network.
 */
public class Pledge {

    private String id;
    private String ownerParticipantID;
    private String targetParticipantID;
    private String sharedAddressID;
    private String currencyID;
    private String blockchain;
    private String network;
    private String amount;
    private String status;
    private String pledgeType;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getOwnerParticipantID() {
        return ownerParticipantID;
    }

    public void setOwnerParticipantID(final String ownerParticipantID) {
        this.ownerParticipantID = ownerParticipantID;
    }

    public String getTargetParticipantID() {
        return targetParticipantID;
    }

    public void setTargetParticipantID(final String targetParticipantID) {
        this.targetParticipantID = targetParticipantID;
    }

    public String getSharedAddressID() {
        return sharedAddressID;
    }

    public void setSharedAddressID(final String sharedAddressID) {
        this.sharedAddressID = sharedAddressID;
    }

    public String getCurrencyID() {
        return currencyID;
    }

    public void setCurrencyID(final String currencyID) {
        this.currencyID = currencyID;
    }

    public String getAmount() {
        return amount;
    }

    public void setAmount(final String amount) {
        this.amount = amount;
    }

    public String getBlockchain() {
        return blockchain;
    }

    public void setBlockchain(final String blockchain) {
        this.blockchain = blockchain;
    }

    public String getNetwork() {
        return network;
    }

    public void setNetwork(final String network) {
        this.network = network;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getPledgeType() {
        return pledgeType;
    }

    public void setPledgeType(final String pledgeType) {
        this.pledgeType = pledgeType;
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
}
