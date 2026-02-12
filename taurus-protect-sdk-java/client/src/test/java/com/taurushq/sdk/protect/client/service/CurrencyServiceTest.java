package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class CurrencyServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private CurrencyService currencyService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        currencyService = new CurrencyService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new CurrencyService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new CurrencyService(apiClient, null));
    }

    @Test
    void getCurrency_throwsOnNullCurrencyId() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrency(null));
    }

    @Test
    void getCurrency_throwsOnEmptyCurrencyId() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrency(""));
    }

    @Test
    void getCurrencyByBlockchain_throwsOnNullBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrencyByBlockchain(null, "mainnet"));
    }

    @Test
    void getCurrencyByBlockchain_throwsOnEmptyBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrencyByBlockchain("", "mainnet"));
    }

    @Test
    void getCurrencyByBlockchain_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrencyByBlockchain("ETH", null));
    }

    @Test
    void getCurrencyByBlockchain_throwsOnEmptyNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                currencyService.getCurrencyByBlockchain("ETH", ""));
    }
}
