package com.taurushq.sdk.protect.client.model;

/**
 * Represents a trigger condition for an automated action.
 * <p>
 * Triggers define when an action should be executed, such as when
 * a balance threshold is crossed.
 *
 * @see Action
 */
public class ActionTrigger {

    private String kind;
    private TriggerBalance balance;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public TriggerBalance getBalance() {
        return balance;
    }

    public void setBalance(final TriggerBalance balance) {
        this.balance = balance;
    }
}
