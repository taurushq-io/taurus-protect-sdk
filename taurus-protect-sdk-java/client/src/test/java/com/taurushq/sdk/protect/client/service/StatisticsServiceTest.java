package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class StatisticsServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private StatisticsService statisticsService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        statisticsService = new StatisticsService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new StatisticsService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new StatisticsService(apiClient, null));
    }

    // Note: getPortfolioStatistics() has no required parameters,
    // so no input validation tests are needed beyond constructor validation.
    // Integration tests would verify API responses are correctly mapped.
}
