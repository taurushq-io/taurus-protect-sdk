package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class AirGapServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private AirGapService airGapService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        airGapService = new AirGapService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AirGapService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AirGapService(apiClient, null));
    }

    @Test
    void constructor_createsServiceSuccessfully() {
        assertNotNull(airGapService);
    }

    @Test
    void getOutgoingAirGap_throwsOnNullRequestIds() {
        assertThrows(NullPointerException.class, () ->
                airGapService.getOutgoingAirGap(null));
    }

    @Test
    void getOutgoingAirGap_throwsOnEmptyRequestIds() {
        assertThrows(IllegalArgumentException.class, () ->
                airGapService.getOutgoingAirGap(Collections.emptyList()));
    }

    @Test
    void submitIncomingAirGap_throwsOnNullPayload() {
        assertThrows(NullPointerException.class, () ->
                airGapService.submitIncomingAirGap(null));
    }

    @Test
    void submitIncomingAirGap_throwsOnEmptyPayload() {
        assertThrows(IllegalArgumentException.class, () ->
                airGapService.submitIncomingAirGap(""));
    }
}
