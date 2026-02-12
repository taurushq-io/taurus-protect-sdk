package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TaurusNetworkLendingServiceTest {

    private TaurusNetworkLendingService service;

    @BeforeEach
    void setUp() {
        ApiClient apiClient = new ApiClient();
        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        service = new TaurusNetworkLendingService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_success() {
        assertNotNull(service);
    }

    @Test
    void constructor_nullApiClient() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkLendingService(null, new ApiExceptionMapper());
        });
    }

    @Test
    void constructor_nullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkLendingService(new ApiClient(), null);
        });
    }

    @Test
    void getLendingOffer_nullId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getLendingOffer(null);
        });
    }

    @Test
    void getLendingOffer_emptyId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getLendingOffer("");
        });
    }

    @Test
    void getLendingAgreement_nullId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getLendingAgreement(null);
        });
    }

    @Test
    void getLendingAgreement_emptyId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.getLendingAgreement("");
        });
    }
}
