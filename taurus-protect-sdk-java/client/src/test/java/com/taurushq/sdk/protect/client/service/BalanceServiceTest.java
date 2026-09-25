package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertEquals;
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
    void getBalances_nullCursorAsksForTheFirstPageWithTheDefaultSize() throws Exception {
        // It used to throw NullPointerException; an unset cursor is the first page, sent
        // through requestCursor only (the legacy limit/cursor pair is never sent).
        StubTransport stub = StubTransport.replying("{}");
        new BalanceService(stub.client(), apiExceptionMapper).getBalances((ApiRequestCursor) null);
        assertEquals(Collections.singletonList(Arrays.asList("requestCursor.pageSize", "20")),
                stub.only().query());
    }

    @Test
    void getBalancesWithCurrency_sendsTheCurrencyAndTheDefaultSize() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new BalanceService(stub.client(), apiExceptionMapper).getBalances("ETH", (ApiRequestCursor) null);
        assertEquals(Arrays.asList(Arrays.asList("currency", "ETH"), Arrays.asList("requestCursor.pageSize", "20")),
                stub.only().query());
    }
}
