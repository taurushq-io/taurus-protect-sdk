package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a lending agreement in the Taurus Network.
 */
public class LendingAgreement {

    private String id;
    private String lenderParticipantID;
    private String borrowerParticipantID;
    private String lendingOfferID;
    private String currencyID;
    private String amount;
    private String annualYield;
    private String duration;
    private String status;
    private String workflowID;
    private OffsetDateTime startLoanDate;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getLenderParticipantID() {
        return lenderParticipantID;
    }

    public void setLenderParticipantID(final String lenderParticipantID) {
        this.lenderParticipantID = lenderParticipantID;
    }

    public String getBorrowerParticipantID() {
        return borrowerParticipantID;
    }

    public void setBorrowerParticipantID(final String borrowerParticipantID) {
        this.borrowerParticipantID = borrowerParticipantID;
    }

    public String getLendingOfferID() {
        return lendingOfferID;
    }

    public void setLendingOfferID(final String lendingOfferID) {
        this.lendingOfferID = lendingOfferID;
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

    public String getAnnualYield() {
        return annualYield;
    }

    public void setAnnualYield(final String annualYield) {
        this.annualYield = annualYield;
    }

    public String getDuration() {
        return duration;
    }

    public void setDuration(final String duration) {
        this.duration = duration;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getWorkflowID() {
        return workflowID;
    }

    public void setWorkflowID(final String workflowID) {
        this.workflowID = workflowID;
    }

    public OffsetDateTime getStartLoanDate() {
        return startLoanDate;
    }

    public void setStartLoanDate(final OffsetDateTime startLoanDate) {
        this.startLoanDate = startLoanDate;
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
