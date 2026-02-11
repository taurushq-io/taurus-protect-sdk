package com.taurushq.sdk.protect.client.model;

/**
 * Represents a job in the Taurus Protect system.
 * <p>
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 *
 * @see JobStatistics
 */
public class Job {

    private String name;
    private JobStatistics statistics;

    public String getName() {
        return name;
    }

    public void setName(final String name) {
        this.name = name;
    }

    public JobStatistics getStatistics() {
        return statistics;
    }

    public void setStatistics(final JobStatistics statistics) {
        this.statistics = statistics;
    }
}
