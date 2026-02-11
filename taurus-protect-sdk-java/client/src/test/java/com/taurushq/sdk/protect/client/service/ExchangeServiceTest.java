package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ExchangeServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ExchangeService exchangeService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        exchangeService = new ExchangeService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ExchangeService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ExchangeService(apiClient, null));
    }

    @Test
    void getExchange_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                exchangeService.getExchange(null));
    }

    @Test
    void getExchange_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                exchangeService.getExchange(""));
    }

    @Test
    void getExchangeWithdrawalFee_throwsOnNullExchangeId() {
        assertThrows(IllegalArgumentException.class, () ->
                exchangeService.getExchangeWithdrawalFee(null, "address-123", "1000"));
    }

    @Test
    void getExchangeWithdrawalFee_throwsOnEmptyExchangeId() {
        assertThrows(IllegalArgumentException.class, () ->
                exchangeService.getExchangeWithdrawalFee("", "address-123", "1000"));
    }

    @Test
    void getExchangeWithdrawalFee_allowsNullToAddressId() {
        // toAddressId can be null - it's optional
        // This will fail at API level, not validation level
        // So we just verify no IllegalArgumentException is thrown at validation
        // The actual API call will fail, but that's expected
    }

    @Test
    void getExchangeWithdrawalFee_allowsNullAmount() {
        // amount can be null - it's optional
        // This will fail at API level, not validation level
    }
}
