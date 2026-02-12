package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Job;
import com.taurushq.sdk.protect.client.model.JobStatistics;
import com.taurushq.sdk.protect.client.model.JobStatus;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJob;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJobStatistics;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJobStatus;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class JobMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordJobStatistics statsDto = new TgvalidatordJobStatistics();
        statsDto.setPending("5");
        statsDto.setSuccesses("100");
        statsDto.setFailures("2");

        TgvalidatordJob dto = new TgvalidatordJob();
        dto.setName("balance-sync");
        dto.setStatistics(statsDto);

        Job job = JobMapper.INSTANCE.fromDTO(dto);

        assertEquals("balance-sync", job.getName());
        assertNotNull(job.getStatistics());
        assertEquals("5", job.getStatistics().getPending());
        assertEquals("100", job.getStatistics().getSuccesses());
        assertEquals("2", job.getStatistics().getFailures());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordJob dto = new TgvalidatordJob();
        dto.setName("test-job");

        Job job = JobMapper.INSTANCE.fromDTO(dto);

        assertEquals("test-job", job.getName());
        assertNull(job.getStatistics());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Job job = JobMapper.INSTANCE.fromDTO(null);
        assertNull(job);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordJob dto1 = new TgvalidatordJob();
        dto1.setName("job-1");

        TgvalidatordJob dto2 = new TgvalidatordJob();
        dto2.setName("job-2");

        List<Job> jobs = JobMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(jobs);
        assertEquals(2, jobs.size());
        assertEquals("job-1", jobs.get(0).getName());
        assertEquals("job-2", jobs.get(1).getName());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Job> jobs = JobMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(jobs);
        assertTrue(jobs.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Job> jobs = JobMapper.INSTANCE.fromDTOList(null);
        assertNull(jobs);
    }

    @Test
    void fromStatusDTO_mapsAllFields() {
        OffsetDateTime startedAt = OffsetDateTime.now();
        OffsetDateTime updatedAt = OffsetDateTime.now().plusMinutes(5);
        OffsetDateTime timeoutAt = OffsetDateTime.now().plusHours(1);

        TgvalidatordJobStatus dto = new TgvalidatordJobStatus();
        dto.setId("status-123");
        dto.setStartedAt(startedAt);
        dto.setUpdatedAt(updatedAt);
        dto.setTimeoutAt(timeoutAt);
        dto.setMessage("Processing...");
        dto.setStatus("running");

        JobStatus status = JobMapper.INSTANCE.fromStatusDTO(dto);

        assertEquals("status-123", status.getId());
        assertEquals(startedAt, status.getStartedAt());
        assertEquals(updatedAt, status.getUpdatedAt());
        assertEquals(timeoutAt, status.getTimeoutAt());
        assertEquals("Processing...", status.getMessage());
        assertEquals("running", status.getStatus());
    }

    @Test
    void fromStatusDTO_handlesNullFields() {
        TgvalidatordJobStatus dto = new TgvalidatordJobStatus();
        dto.setId("minimal");

        JobStatus status = JobMapper.INSTANCE.fromStatusDTO(dto);

        assertEquals("minimal", status.getId());
        assertNull(status.getStartedAt());
        assertNull(status.getMessage());
        assertNull(status.getStatus());
    }

    @Test
    void fromStatusDTO_handlesNullDto() {
        JobStatus status = JobMapper.INSTANCE.fromStatusDTO(null);
        assertNull(status);
    }

    @Test
    void fromStatisticsDTO_mapsAllFields() {
        TgvalidatordJobStatus lastSuccess = new TgvalidatordJobStatus();
        lastSuccess.setId("success-1");

        TgvalidatordJobStatus lastFailure = new TgvalidatordJobStatus();
        lastFailure.setId("failure-1");

        TgvalidatordJobStatistics dto = new TgvalidatordJobStatistics();
        dto.setPending("10");
        dto.setSuccesses("500");
        dto.setFailures("5");
        dto.setLastSuccess(lastSuccess);
        dto.setLastFailure(lastFailure);
        dto.setAvgDuration("1000ms");
        dto.setMaxDuration("5000ms");

        JobStatistics stats = JobMapper.INSTANCE.fromStatisticsDTO(dto);

        assertEquals("10", stats.getPending());
        assertEquals("500", stats.getSuccesses());
        assertEquals("5", stats.getFailures());
        assertNotNull(stats.getLastSuccess());
        assertEquals("success-1", stats.getLastSuccess().getId());
        assertNotNull(stats.getLastFailure());
        assertEquals("failure-1", stats.getLastFailure().getId());
        assertEquals("1000ms", stats.getAvgDuration());
        assertEquals("5000ms", stats.getMaxDuration());
    }

    @Test
    void fromStatisticsDTO_handlesNullDto() {
        JobStatistics stats = JobMapper.INSTANCE.fromStatisticsDTO(null);
        assertNull(stats);
    }
}
