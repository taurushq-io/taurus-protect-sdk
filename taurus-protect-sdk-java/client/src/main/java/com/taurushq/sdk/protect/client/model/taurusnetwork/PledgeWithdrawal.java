package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a pledge withdrawal in the Taurus Network.
 */
public class PledgeWithdrawal {

    private String id;
    private String pledgeID;
    private String destinationSharedAddressID;
    private String amount;
    private String status;
    private String txHash;
    private String requestID;
    private String initiatorParticipantID;
    private OffsetDateTime createdAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getPledgeID() {
        return pledgeID;
    }

    public void setPledgeID(final String pledgeID) {
        this.pledgeID = pledgeID;
    }

    public String getDestinationSharedAddressID() {
        return destinationSharedAddressID;
    }

    public void setDestinationSharedAddressID(final String destinationSharedAddressID) {
        this.destinationSharedAddressID = destinationSharedAddressID;
    }

    public String getAmount() {
        return amount;
    }

    public void setAmount(final String amount) {
        this.amount = amount;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getTxHash() {
        return txHash;
    }

    public void setTxHash(final String txHash) {
        this.txHash = txHash;
    }

    public String getRequestID() {
        return requestID;
    }

    public void setRequestID(final String requestID) {
        this.requestID = requestID;
    }

    public String getInitiatorParticipantID() {
        return initiatorParticipantID;
    }

    public void setInitiatorParticipantID(final String initiatorParticipantID) {
        this.initiatorParticipantID = initiatorParticipantID;
    }

    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(final OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }
}
