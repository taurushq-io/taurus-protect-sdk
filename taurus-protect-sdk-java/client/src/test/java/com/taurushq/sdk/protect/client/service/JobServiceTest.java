package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class JobServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private JobService jobService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        jobService = new JobService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new JobService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new JobService(apiClient, null));
    }

    @Test
    void getJob_throwsOnNullName() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJob(null));
    }

    @Test
    void getJob_throwsOnEmptyName() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJob(""));
    }

    @Test
    void getJobStatus_throwsOnNullName() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJobStatus(null, "id-123"));
    }

    @Test
    void getJobStatus_throwsOnEmptyName() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJobStatus("", "id-123"));
    }

    @Test
    void getJobStatus_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJobStatus("job-name", null));
    }

    @Test
    void getJobStatus_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                jobService.getJobStatus("job-name", ""));
    }
}
