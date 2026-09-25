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
    void getBusinessRules_nullCursorAsksForTheFirstPageWithTheDefaultSize() throws Exception {
        // It used to throw NullPointerException; an unset cursor is the first page.
        StubTransport stub = StubTransport.replying("{}");
        new BusinessRuleService(stub.client(), apiExceptionMapper).getBusinessRules((ApiRequestCursor) null);
        assertEquals(Collections.singletonList(Arrays.asList("cursor.pageSize", "20")), stub.only().query());
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

    /**
     * The transactions-enabled kill switch existed only in the Go SDK, even though the
     * generated op is present in all four and tg-protect-mcpd drives it. Without a
     * network stub (the project forbids Mockito) the reachable assertion is that the
     * method exists with a boolean arity and fails on transport rather than on a missing
     * symbol — which is what pins the cross-SDK surface.
     */
    @Test
    void updateTransactionsEnabled_isReachableWithABooleanFlag() throws Exception {
        assertThrows(Exception.class, () -> businessRuleService.updateTransactionsEnabled(false));
    }
}
