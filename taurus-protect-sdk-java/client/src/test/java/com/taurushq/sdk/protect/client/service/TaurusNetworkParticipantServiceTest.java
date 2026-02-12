package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TaurusNetworkParticipantServiceTest {

    private TaurusNetworkParticipantService service;

    @BeforeEach
    void setUp() {
        ApiClient apiClient = new ApiClient();
        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        service = new TaurusNetworkParticipantService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_success() {
        assertNotNull(service);
    }

    @Test
    void constructor_nullApiClient() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkParticipantService(null, new ApiExceptionMapper());
        });
    }

    @Test
    void constructor_nullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () -> {
            new TaurusNetworkParticipantService(new ApiClient(), null);
        });
    }

    @Test
    void get_nullId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.get(null, false);
        });
    }

    @Test
    void get_emptyId() {
        assertThrows(IllegalArgumentException.class, () -> {
            service.get("", false);
        });
    }
}
