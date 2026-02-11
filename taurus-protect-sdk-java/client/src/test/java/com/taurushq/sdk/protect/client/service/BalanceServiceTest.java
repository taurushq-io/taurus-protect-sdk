package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class BalanceServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private BalanceService balanceService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        balanceService = new BalanceService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new BalanceService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new BalanceService(apiClient, null));
    }

    @Test
    void getBalances_throwsOnNullCursor() {
        assertThrows(NullPointerException.class, () ->
                balanceService.getBalances(null));
    }

    @Test
    void getBalancesWithCurrency_throwsOnNullCursor() {
        assertThrows(NullPointerException.class, () ->
                balanceService.getBalances("ETH", null));
    }
}
