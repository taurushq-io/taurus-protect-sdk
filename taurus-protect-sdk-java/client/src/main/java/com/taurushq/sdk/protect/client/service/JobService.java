package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.JobMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Job;
import com.taurushq.sdk.protect.client.model.JobStatus;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.JobsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetJobReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetJobStatusReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetJobsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for monitoring jobs in the Taurus Protect system.
 * <p>
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all jobs
 * List<Job> jobs = client.getJobService().getJobs();
 *
 * // Get a specific job
 * Job job = client.getJobService().getJob("balance-sync");
 *
 * // Get job execution status
 * JobStatus status = client.getJobService().getJobStatus("balance-sync", "exec-123");
 * }</pre>
 *
 * @see Job
 * @see JobStatus
 */
public class JobService {

    private final JobsApi jobsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final JobMapper jobMapper;

    /**
     * Creates a new JobService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public JobService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.jobsApi = new JobsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.jobMapper = JobMapper.INSTANCE;
    }

    /**
     * Retrieves all jobs.
     *
     * @return the list of jobs
     * @throws ApiException if the API call fails
     */
    public List<Job> getJobs() throws ApiException {
        try {
            TgvalidatordGetJobsReply reply = jobsApi.jobServiceGetJobs();
            return jobMapper.fromDTOList(reply.getJobs());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a job by name.
     *
     * @param name the job name
     * @return the job
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if name is null or empty
     */
    public Job getJob(final String name) throws ApiException {
        checkArgument(name != null && !name.isEmpty(), "name must not be null or empty");
        try {
            TgvalidatordGetJobReply reply = jobsApi.jobServiceGetJob(name);
            return jobMapper.fromDTO(reply.getJob());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves the status of a specific job execution.
     *
     * @param name the job name
     * @param id   the job execution ID
     * @return the job status
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if name or id is null or empty
     */
    public JobStatus getJobStatus(final String name, final String id) throws ApiException {
        checkArgument(name != null && !name.isEmpty(), "name must not be null or empty");
        checkArgument(id != null && !id.isEmpty(), "id must not be null or empty");
        try {
            TgvalidatordGetJobStatusReply reply = jobsApi.jobServiceGetJobStatus(name, id);
            return jobMapper.fromStatusDTO(reply.getStatus());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
