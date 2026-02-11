package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class BusinessRuleServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private BusinessRuleService businessRuleService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        businessRuleService = new BusinessRuleService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new BusinessRuleService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new BusinessRuleService(apiClient, null));
    }

    @Test
    void getBusinessRules_throwsOnNullCursor() {
        assertThrows(NullPointerException.class, () ->
                businessRuleService.getBusinessRules(null));
    }

    @Test
    void getBusinessRulesByWallet_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                businessRuleService.getBusinessRulesByWallet(0, null));
    }

    @Test
    void getBusinessRulesByCurrency_throwsOnNullCurrencyId() {
        assertThrows(IllegalArgumentException.class, () ->
                businessRuleService.getBusinessRulesByCurrency(null, null));
    }

    @Test
    void getBusinessRulesByCurrency_throwsOnEmptyCurrencyId() {
        assertThrows(IllegalArgumentException.class, () ->
                businessRuleService.getBusinessRulesByCurrency("", null));
    }
}
