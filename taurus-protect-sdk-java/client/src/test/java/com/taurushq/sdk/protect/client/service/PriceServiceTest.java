package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class PriceServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private PriceService priceService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        priceService = new PriceService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new PriceService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new PriceService(apiClient, null));
    }

    @Test
    void getPriceHistory_throwsOnNullBase() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.getPriceHistory(null, "USD", 100));
    }

    @Test
    void getPriceHistory_throwsOnEmptyBase() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.getPriceHistory("", "USD", 100));
    }

    @Test
    void getPriceHistory_throwsOnNullQuote() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.getPriceHistory("ETH", null, 100));
    }

    @Test
    void getPriceHistory_throwsOnEmptyQuote() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.getPriceHistory("ETH", "", 100));
    }

    @Test
    void getPriceHistory_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.getPriceHistory("ETH", "USD", 0));
    }

    @Test
    void convert_throwsOnNullCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.convert(null, "1000", Arrays.asList("USD")));
    }

    @Test
    void convert_throwsOnEmptyCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.convert("", "1000", Arrays.asList("USD")));
    }

    @Test
    void convert_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.convert("ETH", null, Arrays.asList("USD")));
    }

    @Test
    void convert_throwsOnEmptyTargetCurrencies() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.convert("ETH", "1000", Collections.emptyList()));
    }

    @Test
    void convert_throwsOnNullTargetCurrencies() {
        assertThrows(IllegalArgumentException.class, () ->
                priceService.convert("ETH", "1000", null));
    }
}
