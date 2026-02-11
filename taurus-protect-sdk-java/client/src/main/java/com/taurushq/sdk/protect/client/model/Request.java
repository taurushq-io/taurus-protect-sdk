package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a transaction request in the Taurus Protect system.
 * <p>
 * A request represents an action to be performed on the blockchain, such as a transfer,
 * staking operation, or contract call. Requests go through an approval workflow and
 * can have various statuses (pending, approved, rejected, broadcast, confirmed, etc.).
 * <p>
 * The request lifecycle typically follows these stages:
 * <ol>
 *   <li>Created - The request is submitted to the system</li>
 *   <li>Pending Approval - Waiting for required approvers</li>
 *   <li>Approved - All required approvals received</li>
 *   <li>Broadcast - Transaction submitted to the blockchain</li>
 *   <li>Confirmed - Transaction confirmed on the blockchain</li>
 * </ol>
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a transfer request
 * Request request = client.getRequestService().createTransferRequest(
 *     sourceAddressId, destinationAddress, amount, comment);
 *
 * // Check the request status
 * if (request.getStatus() == RequestStatus.PENDING_APPROVAL) {
 *     System.out.println("Awaiting approval from: " + request.getNeedsApprovalFrom());
 * }
 * }</pre>
 *
 * @see RequestStatus
 * @see RequestService
 * @see SignedRequest
 */
public class Request {

    /**
     * The unique identifier of the request assigned by Taurus Protect.
     */
    private long id;

    /**
     * The tenant identifier for multi-tenant environments.
     */
    private int tenantId;

    /**
     * The currency code for the request (e.g., "ETH", "BTC").
     */
    private String currency;

    /**
     * The serialized transaction envelope containing the unsigned transaction data.
     */
    private String envelope;

    /**
     * The current status of the request in the approval workflow.
     */
    private RequestStatus status;

    /**
     * The audit trail of actions performed on the request.
     */
    private List<RequestTrail> trails;

    /**
     * The list of signed transaction requests after approval.
     */
    private List<SignedRequest> signedRequests;

    /**
     * The date and time when the request was created.
     */
    private OffsetDateTime creationDate;

    /**
     * The date and time when the request was last updated.
     */
    private OffsetDateTime updateDate;

    /**
     * Metadata containing details about the transaction (source, destination, amount, etc.).
     */
    private RequestMetadata metadata;

    /**
     * Custom key-value attributes associated with the request.
     */
    private List<Attribute> attributes;

    /**
     * The business rule that was applied to this request.
     */
    private String rule;

    /**
     * Information about the approval configuration for this request.
     */
    private RequestApprovers approvers;

    /**
     * The type of request (e.g., "transfer", "stake", "delegate").
     */
    private String type;

    /**
     * Detailed information about the request's currency.
     */
    private Currency currencyInfo;

    /**
     * The list of user IDs or roles that still need to approve the request.
     */
    private List<String> needsApprovalFrom;

    /**
     * The bundle ID if this request is part of a request bundle.
     */
    private String requestBundleId;

    /**
     * An external identifier for the request provided by the client application.
     */
    private String externalRequestId;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the request.
     *
     * @return the request ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the request.
     *
     * @param id the request ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the tenant identifier for multi-tenant environments.
     *
     * @return the tenant ID
     */
    public int getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant identifier.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(int tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Returns the currency code for the request (e.g., "ETH", "BTC").
     *
     * @return the currency code
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code for the request.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns the serialized transaction envelope containing the unsigned transaction data.
     *
     * @return the transaction envelope
     */
    public String getEnvelope() {
        return envelope;
    }

    /**
     * Sets the transaction envelope.
     *
     * @param envelope the transaction envelope to set
     */
    public void setEnvelope(String envelope) {
        this.envelope = envelope;
    }

    /**
     * Returns the current status of the request in the approval workflow.
     *
     * @return the request status
     * @see RequestStatus
     */
    public RequestStatus getStatus() {
        return status;
    }

    /**
     * Sets the current status of the request.
     *
     * @param status the request status to set
     */
    public void setStatus(RequestStatus status) {
        this.status = status;
    }

    /**
     * Returns the audit trail of actions performed on the request.
     *
     * @return the list of request trails, or {@code null} if not populated
     */
    public List<RequestTrail> getTrails() {
        return trails;
    }

    /**
     * Sets the audit trail for the request.
     *
     * @param trails the list of request trails to set
     */
    public void setTrails(List<RequestTrail> trails) {
        this.trails = trails;
    }

    /**
     * Returns the list of signed transaction requests after approval.
     * <p>
     * Signed requests contain the cryptographically signed transaction data
     * ready for broadcast to the blockchain.
     *
     * @return the list of signed requests, or {@code null} if not yet signed
     */
    public List<SignedRequest> getSignedRequests() {
        return signedRequests;
    }

    /**
     * Sets the list of signed requests.
     *
     * @param signedRequests the list of signed requests to set
     */
    public void setSignedRequests(List<SignedRequest> signedRequests) {
        this.signedRequests = signedRequests;
    }

    /**
     * Returns the date and time when the request was created.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the date and time when the request was created.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Returns the date and time when the request was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the date and time when the request was last updated.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Returns the metadata containing details about the transaction.
     * <p>
     * Metadata includes information such as source address, destination address,
     * amount, and other transaction-specific details.
     *
     * @return the request metadata, or {@code null} if not available
     */
    public RequestMetadata getMetadata() {
        return metadata;
    }

    /**
     * Sets the request metadata.
     *
     * @param metadata the request metadata to set
     */
    public void setMetadata(RequestMetadata metadata) {
        this.metadata = metadata;
    }

    /**
     * Returns the custom key-value attributes associated with the request.
     *
     * @return the list of attributes, or {@code null} if none are set
     */
    public List<Attribute> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom key-value attributes for the request.
     *
     * @param attributes the list of attributes to set
     */
    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }

    /**
     * Returns the business rule that was applied to this request.
     *
     * @return the rule name, or {@code null} if no rule was applied
     */
    public String getRule() {
        return rule;
    }

    /**
     * Sets the business rule for the request.
     *
     * @param rule the rule name to set
     */
    public void setRule(String rule) {
        this.rule = rule;
    }

    /**
     * Returns information about the approval configuration for this request.
     *
     * @return the approvers configuration, or {@code null} if not available
     */
    public RequestApprovers getApprovers() {
        return approvers;
    }

    /**
     * Sets the approvers configuration for the request.
     *
     * @param approvers the approvers configuration to set
     */
    public void setApprovers(RequestApprovers approvers) {
        this.approvers = approvers;
    }

    /**
     * Returns the type of request (e.g., "transfer", "stake", "delegate").
     *
     * @return the request type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the type of request.
     *
     * @param type the request type to set
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Returns detailed information about the request's currency.
     *
     * @return the currency information, or {@code null} if not available
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the detailed currency information for the request.
     *
     * @param currencyInfo the currency information to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Returns the list of user IDs or roles that still need to approve the request.
     * <p>
     * This list is populated when the request is in a pending approval state and
     * indicates who can take action on the request.
     *
     * @return the list of pending approvers, or {@code null} if not applicable
     */
    public List<String> getNeedsApprovalFrom() {
        return needsApprovalFrom;
    }

    /**
     * Sets the list of pending approvers.
     *
     * @param needsApprovalFrom the list of pending approvers to set
     */
    public void setNeedsApprovalFrom(List<String> needsApprovalFrom) {
        this.needsApprovalFrom = needsApprovalFrom;
    }

    /**
     * Returns the bundle ID if this request is part of a request bundle.
     * <p>
     * Request bundles allow grouping multiple requests for batch processing.
     *
     * @return the request bundle ID, or {@code null} if not part of a bundle
     */
    public String getRequestBundleId() {
        return requestBundleId;
    }

    /**
     * Sets the request bundle ID.
     *
     * @param requestBundleId the request bundle ID to set
     */
    public void setRequestBundleId(String requestBundleId) {
        this.requestBundleId = requestBundleId;
    }

    /**
     * Returns the external identifier for the request provided by the client application.
     * <p>
     * This ID can be used to correlate requests with external systems.
     *
     * @return the external request ID, or {@code null} if not set
     */
    public String getExternalRequestId() {
        return externalRequestId;
    }

    /**
     * Sets the external identifier for the request.
     *
     * @param externalRequestId the external request ID to set
     */
    public void setExternalRequestId(String externalRequestId) {
        this.externalRequestId = externalRequestId;
    }
}



