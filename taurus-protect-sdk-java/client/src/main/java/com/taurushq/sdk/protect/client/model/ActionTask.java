package com.taurushq.sdk.protect.client.model;

/**
 * Represents a task to be executed when an action is triggered.
 * <p>
 * Tasks can be transfers, notifications, etc.
 *
 * @see Action
 */
public class ActionTask {

    private String kind;
    private TaskTransfer transfer;
    private TaskNotification notification;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }

    public TaskTransfer getTransfer() {
        return transfer;
    }

    public void setTransfer(final TaskTransfer transfer) {
        this.transfer = transfer;
    }

    public TaskNotification getNotification() {
        return notification;
    }

    public void setNotification(final TaskNotification notification) {
        this.notification = notification;
    }
}
