package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a shared address in the Taurus Network.
 */
public class SharedAddress {

    private String id;
    private String internalAddressID;
    private String ownerParticipantId;
    private String targetParticipantId;
    private String blockchain;
    private String network;
    private String address;
    private String status;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getInternalAddressID() {
        return internalAddressID;
    }

    public void setInternalAddressID(final String internalAddressID) {
        this.internalAddressID = internalAddressID;
    }

    public String getOwnerParticipantId() {
        return ownerParticipantId;
    }

    public void setOwnerParticipantId(final String ownerParticipantId) {
        this.ownerParticipantId = ownerParticipantId;
    }

    public String getTargetParticipantId() {
        return targetParticipantId;
    }

    public void setTargetParticipantId(final String targetParticipantId) {
        this.targetParticipantId = targetParticipantId;
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

    public String getAddress() {
        return address;
    }

    public void setAddress(final String address) {
        this.address = address;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
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
