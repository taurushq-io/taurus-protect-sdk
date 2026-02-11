package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a notification task for automated actions.
 * <p>
 * A notification task sends emails to specified addresses when the action triggers.
 */
public class TaskNotification {

    private List<String> emailAddresses;
    private String notificationMessage;
    private String numberOfReminders;

    public List<String> getEmailAddresses() {
        return emailAddresses;
    }

    public void setEmailAddresses(final List<String> emailAddresses) {
        this.emailAddresses = emailAddresses;
    }

    public String getNotificationMessage() {
        return notificationMessage;
    }

    public void setNotificationMessage(final String notificationMessage) {
        this.notificationMessage = notificationMessage;
    }

    public String getNumberOfReminders() {
        return numberOfReminders;
    }

    public void setNumberOfReminders(final String numberOfReminders) {
        this.numberOfReminders = numberOfReminders;
    }
}
