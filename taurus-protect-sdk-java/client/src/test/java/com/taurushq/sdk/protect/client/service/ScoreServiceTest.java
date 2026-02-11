package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ScoreServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ScoreService scoreService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        scoreService = new ScoreService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ScoreService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ScoreService(apiClient, null));
    }

    @Test
    void refreshAddressScore_throwsOnZeroAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                scoreService.refreshAddressScore(0, "chainalysis"));
    }

    @Test
    void refreshAddressScore_throwsOnNegativeAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                scoreService.refreshAddressScore(-1, "chainalysis"));
    }

    @Test
    void refreshAddressScore_throwsOnNullScoreProvider() {
        assertThrows(NullPointerException.class, () ->
                scoreService.refreshAddressScore(1, null));
    }

    @Test
    void refreshWhitelistedAddressScore_throwsOnZeroAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                scoreService.refreshWhitelistedAddressScore(0, "chainalysis"));
    }

    @Test
    void refreshWhitelistedAddressScore_throwsOnNullScoreProvider() {
        assertThrows(NullPointerException.class, () ->
                scoreService.refreshWhitelistedAddressScore(1, null));
    }
}
