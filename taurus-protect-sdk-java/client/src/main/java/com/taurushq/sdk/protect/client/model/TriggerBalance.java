package com.taurushq.sdk.protect.client.model;

/**
 * Represents a balance trigger condition for an action.
 * <p>
 * Balance triggers allow actions to be executed when an address or wallet
 * balance crosses a specified threshold.
 *
 * @see ActionTrigger
 */
public class TriggerBalance {

    private ActionTarget target;
    private ActionComparator comparator;
    private ActionAmount amount;

    public ActionTarget getTarget() {
        return target;
    }

    public void setTarget(final ActionTarget target) {
        this.target = target;
    }

    public ActionComparator getComparator() {
        return comparator;
    }

    public void setComparator(final ActionComparator comparator) {
        this.comparator = comparator;
    }

    public ActionAmount getAmount() {
        return amount;
    }

    public void setAmount(final ActionAmount amount) {
        this.amount = amount;
    }
}
