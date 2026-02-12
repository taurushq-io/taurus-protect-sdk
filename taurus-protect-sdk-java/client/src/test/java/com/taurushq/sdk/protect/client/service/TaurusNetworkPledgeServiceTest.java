package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TaurusNetworkPledgeServiceTest {

    private TaurusNetworkPledgeService service;

    @BeforeEach
    void setUp() {
        ApiClient apiClient = new ApiClient();
        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        service = new TaurusNetworkPledgeService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_success() {
        assertNotNull(service);
    }

    @Test
    void constructor_nullApiClient() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkPledgeService(null, new ApiExceptionMapper());
        });
    }

    @Test
    void constructor_nullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkPledgeService(new ApiClient(), null);
        });
    }

    @Test
    void get_nullId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.get(null);
        });
    }

    @Test
    void get_emptyId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.get("");
        });
    }

    @Test
    void listWithdrawals_nullPledgeId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.listWithdrawals(null, null, null, null);
        });
    }

    @Test
    void listWithdrawals_emptyPledgeId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.listWithdrawals("", null, null, null);
        });
    }
}
