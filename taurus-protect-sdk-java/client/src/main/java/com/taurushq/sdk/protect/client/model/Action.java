package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents an automated action configuration containing a trigger and tasks.
 * <p>
 * An action defines what should happen (tasks) when certain conditions are met (trigger).
 *
 * @see ActionTrigger
 * @see ActionTask
 */
public class Action {

    private ActionTrigger trigger;
    private List<ActionTask> tasks;

    public ActionTrigger getTrigger() {
        return trigger;
    }

    public void setTrigger(final ActionTrigger trigger) {
        this.trigger = trigger;
    }

    public List<ActionTask> getTasks() {
        return tasks;
    }

    public void setTasks(final List<ActionTask> tasks) {
        this.tasks = tasks;
    }
}
