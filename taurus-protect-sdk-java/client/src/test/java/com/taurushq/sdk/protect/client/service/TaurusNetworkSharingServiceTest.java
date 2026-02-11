package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TaurusNetworkSharingServiceTest {

    private TaurusNetworkSharingService service;

    @BeforeEach
    void setUp() {
        ApiClient apiClient = new ApiClient();
        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        service = new TaurusNetworkSharingService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_success() {
        assertNotNull(service);
    }

    @Test
    void constructor_nullApiClient() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkSharingService(null, new ApiExceptionMapper());
        });
    }

    @Test
    void constructor_nullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkSharingService(new ApiClient(), null);
        });
    }
}
