package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a cryptographically signed transaction request.
 * <p>
 * A signed request contains the transaction data that has been signed by the
 * Taurus Protect system after the request has been approved. This signed data
 * is ready to be broadcast to the blockchain network.
 * <p>
 * The lifecycle of a signed request:
 * <ol>
 *   <li>Created - Signature generated after request approval</li>
 *   <li>Broadcast - Submitted to the blockchain network</li>
 *   <li>Confirmed - Included in a block on the blockchain</li>
 * </ol>
 *
 * @see Request
 * @see RequestStatus
 */
public class SignedRequest {

    /**
     * The unique identifier of the signed request.
     */
    private long id;

    /**
     * The signed transaction data in the blockchain-specific format.
     */
    private String signedRequest;

    /**
     * The current status of the signed request.
     */
    private RequestStatus status;

    /**
     * The blockchain transaction hash after broadcast.
     */
    private String hash;

    /**
     * The block number where the transaction was included.
     */
    private long block;

    /**
     * Additional details about the signed request or any errors encountered.
     */
    private String details;

    /**
     * The date and time when the signed request was created.
     */
    private OffsetDateTime creationDate;

    /**
     * The date and time when the signed request was last updated.
     */
    private OffsetDateTime updateDate;

    /**
     * The date and time when the transaction was broadcast to the network.
     */
    private OffsetDateTime broadcastDate;

    /**
     * The date and time when the transaction was confirmed on the blockchain.
     */
    private OffsetDateTime confirmationDate;

    /**
     * Custom key-value attributes associated with the signed request.
     */
    private List<Attribute> attributes;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the signed request.
     *
     * @return the signed request ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the signed request.
     *
     * @param id the signed request ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the signed transaction data in the blockchain-specific format.
     * <p>
     * This data can be submitted directly to the blockchain network for execution.
     *
     * @return the signed transaction data
     */
    public String getSignedRequest() {
        return signedRequest;
    }

    /**
     * Sets the signed transaction data.
     *
     * @param signedRequest the signed transaction data to set
     */
    public void setSignedRequest(String signedRequest) {
        this.signedRequest = signedRequest;
    }

    /**
     * Returns the current status of the signed request.
     *
     * @return the signed request status
     * @see RequestStatus
     */
    public RequestStatus getStatus() {
        return status;
    }

    /**
     * Sets the current status of the signed request.
     *
     * @param status the status to set
     */
    public void setStatus(RequestStatus status) {
        this.status = status;
    }

    /**
     * Returns the blockchain transaction hash after broadcast.
     *
     * @return the transaction hash, or {@code null} if not yet broadcast
     */
    public String getHash() {
        return hash;
    }

    /**
     * Sets the blockchain transaction hash.
     *
     * @param hash the transaction hash to set
     */
    public void setHash(String hash) {
        this.hash = hash;
    }

    /**
     * Returns the block number where the transaction was included.
     *
     * @return the block number, or 0 if not yet included in a block
     */
    public long getBlock() {
        return block;
    }

    /**
     * Sets the block number where the transaction was included.
     *
     * @param block the block number to set
     */
    public void setBlock(long block) {
        this.block = block;
    }

    /**
     * Returns additional details about the signed request or any errors encountered.
     *
     * @return the details string, or {@code null} if not set
     */
    public String getDetails() {
        return details;
    }

    /**
     * Sets additional details about the signed request.
     *
     * @param details the details to set
     */
    public void setDetails(String details) {
        this.details = details;
    }

    /**
     * Returns the date and time when the signed request was created.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the date and time when the signed request was created.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Returns the date and time when the signed request was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the date and time when the signed request was last updated.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Returns the date and time when the transaction was broadcast to the network.
     *
     * @return the broadcast date, or {@code null} if not yet broadcast
     */
    public OffsetDateTime getBroadcastDate() {
        return broadcastDate;
    }

    /**
     * Sets the date and time when the transaction was broadcast.
     *
     * @param broadcastDate the broadcast date to set
     */
    public void setBroadcastDate(OffsetDateTime broadcastDate) {
        this.broadcastDate = broadcastDate;
    }

    /**
     * Returns the date and time when the transaction was confirmed on the blockchain.
     *
     * @return the confirmation date, or {@code null} if not yet confirmed
     */
    public OffsetDateTime getConfirmationDate() {
        return confirmationDate;
    }

    /**
     * Sets the date and time when the transaction was confirmed.
     *
     * @param confirmationDate the confirmation date to set
     */
    public void setConfirmationDate(OffsetDateTime confirmationDate) {
        this.confirmationDate = confirmationDate;
    }

    /**
     * Returns the custom key-value attributes associated with the signed request.
     *
     * @return the list of attributes, or {@code null} if none are set
     */
    public List<Attribute> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom key-value attributes for the signed request.
     *
     * @param attributes the list of attributes to set
     */
    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }
}
