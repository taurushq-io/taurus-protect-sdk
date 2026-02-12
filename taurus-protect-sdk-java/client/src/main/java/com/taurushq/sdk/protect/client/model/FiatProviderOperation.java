package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a fiat provider operation in the Taurus Protect system.
 *
 * @see FiatService
 */
public class FiatProviderOperation {

    private String id;
    private String provider;
    private String label;
    private String operationType;
    private String operationIdentifier;
    private String operationDirection;
    private String status;
    private String amount;
    private String currencyID;
    private String fromAccountID;
    private String toAccountID;
    private String fromDetails;
    private String toDetails;
    private String comment;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(final String provider) {
        this.provider = provider;
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(final String label) {
        this.label = label;
    }

    public String getOperationType() {
        return operationType;
    }

    public void setOperationType(final String operationType) {
        this.operationType = operationType;
    }

    public String getOperationIdentifier() {
        return operationIdentifier;
    }

    public void setOperationIdentifier(final String operationIdentifier) {
        this.operationIdentifier = operationIdentifier;
    }

    public String getOperationDirection() {
        return operationDirection;
    }

    public void setOperationDirection(final String operationDirection) {
        this.operationDirection = operationDirection;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getAmount() {
        return amount;
    }

    public void setAmount(final String amount) {
        this.amount = amount;
    }

    public String getCurrencyID() {
        return currencyID;
    }

    public void setCurrencyID(final String currencyID) {
        this.currencyID = currencyID;
    }

    public String getFromAccountID() {
        return fromAccountID;
    }

    public void setFromAccountID(final String fromAccountID) {
        this.fromAccountID = fromAccountID;
    }

    public String getToAccountID() {
        return toAccountID;
    }

    public void setToAccountID(final String toAccountID) {
        this.toAccountID = toAccountID;
    }

    public String getFromDetails() {
        return fromDetails;
    }

    public void setFromDetails(final String fromDetails) {
        this.fromDetails = fromDetails;
    }

    public String getToDetails() {
        return toDetails;
    }

    public void setToDetails(final String toDetails) {
        this.toDetails = toDetails;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(final String comment) {
        this.comment = comment;
    }

    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    public void setUpdateDate(final OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
