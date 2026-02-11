package com.taurushq.sdk.protect.client.model;

/**
 * Represents a transfer task for automated actions.
 * <p>
 * A transfer task specifies a source, destination, and amount to transfer.
 */
public class TaskTransfer {

    private ActionSource from;
    private ActionDestination to;
    private ActionAmount amount;
    private Boolean topUp;
    private Boolean useAllFunds;

    public ActionSource getFrom() {
        return from;
    }

    public void setFrom(final ActionSource from) {
        this.from = from;
    }

    public ActionDestination getTo() {
        return to;
    }

    public void setTo(final ActionDestination to) {
        this.to = to;
    }

    public ActionAmount getAmount() {
        return amount;
    }

    public void setAmount(final ActionAmount amount) {
        this.amount = amount;
    }

    public Boolean getTopUp() {
        return topUp;
    }

    public void setTopUp(final Boolean topUp) {
        this.topUp = topUp;
    }

    public Boolean getUseAllFunds() {
        return useAllFunds;
    }

    public void setUseAllFunds(final Boolean useAllFunds) {
        this.useAllFunds = useAllFunds;
    }
}
