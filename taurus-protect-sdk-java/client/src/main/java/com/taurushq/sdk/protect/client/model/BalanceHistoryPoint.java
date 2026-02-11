package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents a single data point in a balance history timeline.
 * <p>
 * Balance history points capture the state of a balance at a specific moment in time,
 * enabling historical analysis of wallet or address balances over time. These points
 * are typically used for generating balance charts, audit trails, and portfolio
 * tracking.
 *
 * @see Balance
 * @see BalanceResult
 */
public class BalanceHistoryPoint {

    /**
     * The timestamp when this balance snapshot was recorded.
     */
    private OffsetDateTime pointDate;

    /**
     * The balance values at the recorded point in time.
     */
    private Balance balance;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the timestamp when this balance snapshot was recorded.
     *
     * @return the timestamp of the balance snapshot
     */
    public OffsetDateTime getPointDate() {
        return pointDate;
    }

    /**
     * Sets the timestamp when this balance snapshot was recorded.
     *
     * @param pointDate the timestamp to set
     */
    public void setPointDate(OffsetDateTime pointDate) {
        this.pointDate = pointDate;
    }

    /**
     * Returns the balance values at the recorded point in time.
     *
     * @return the balance at this point in time
     */
    public Balance getBalance() {
        return balance;
    }

    /**
     * Sets the balance values at this point in time.
     *
     * @param balance the balance to set
     */
    public void setBalance(Balance balance) {
        this.balance = balance;
    }
}