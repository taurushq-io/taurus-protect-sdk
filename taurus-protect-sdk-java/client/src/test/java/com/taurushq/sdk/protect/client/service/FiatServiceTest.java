package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class FiatServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private FiatService fiatService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        fiatService = new FiatService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new FiatService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new FiatService(apiClient, null));
    }

    @Test
    void constructor_createsServiceSuccessfully() {
        assertNotNull(fiatService);
    }

    @Test
    void getFiatProviderAccount_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderAccount(null));
    }

    @Test
    void getFiatProviderAccount_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderAccount(""));
    }

    @Test
    void getFiatProviderCounterpartyAccount_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderCounterpartyAccount(null));
    }

    @Test
    void getFiatProviderCounterpartyAccount_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderCounterpartyAccount(""));
    }

    @Test
    void getFiatProviderOperation_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderOperation(null));
    }

    @Test
    void getFiatProviderOperation_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                fiatService.getFiatProviderOperation(""));
    }
}
