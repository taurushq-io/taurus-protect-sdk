package com.taurushq.sdk.protect.client.model;

/**
 * Statistics for a job in the Taurus Protect system.
 * <p>
 * Contains metrics about job execution including success/failure counts
 * and duration statistics.
 *
 * @see Job
 * @see JobStatus
 */
public class JobStatistics {

    private String pending;
    private String successes;
    private String failures;
    private JobStatus lastSuccess;
    private JobStatus lastFailure;
    private String avgDuration;
    private String maxDuration;

    public String getPending() {
        return pending;
    }

    public void setPending(final String pending) {
        this.pending = pending;
    }

    public String getSuccesses() {
        return successes;
    }

    public void setSuccesses(final String successes) {
        this.successes = successes;
    }

    public String getFailures() {
        return failures;
    }

    public void setFailures(final String failures) {
        this.failures = failures;
    }

    public JobStatus getLastSuccess() {
        return lastSuccess;
    }

    public void setLastSuccess(final JobStatus lastSuccess) {
        this.lastSuccess = lastSuccess;
    }

    public JobStatus getLastFailure() {
        return lastFailure;
    }

    public void setLastFailure(final JobStatus lastFailure) {
        this.lastFailure = lastFailure;
    }

    public String getAvgDuration() {
        return avgDuration;
    }

    public void setAvgDuration(final String avgDuration) {
        this.avgDuration = avgDuration;
    }

    public String getMaxDuration() {
        return maxDuration;
    }

    public void setMaxDuration(final String maxDuration) {
        this.maxDuration = maxDuration;
    }
}
