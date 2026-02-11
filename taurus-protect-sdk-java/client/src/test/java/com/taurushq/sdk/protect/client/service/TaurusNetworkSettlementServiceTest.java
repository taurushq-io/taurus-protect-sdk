package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TaurusNetworkSettlementServiceTest {

    private TaurusNetworkSettlementService service;

    @BeforeEach
    void setUp() {
        ApiClient apiClient = new ApiClient();
        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        service = new TaurusNetworkSettlementService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_success() {
        assertNotNull(service);
    }

    @Test
    void constructor_nullApiClient() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkSettlementService(null, new ApiExceptionMapper());
        });
    }

    @Test
    void constructor_nullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkSettlementService(new ApiClient(), null);
        });
    }

    @Test
    void getSettlement_nullId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getSettlement(null);
        });
    }

    @Test
    void getSettlement_emptyId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getSettlement("");
        });
    }
}
