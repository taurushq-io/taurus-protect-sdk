package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class FeePayerServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private FeePayerService feePayerService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        feePayerService = new FeePayerService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new FeePayerService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new FeePayerService(apiClient, null));
    }

    @Test
    void getFeePayer_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                feePayerService.getFeePayer(null));
    }

    @Test
    void getFeePayer_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                feePayerService.getFeePayer(""));
    }
}
