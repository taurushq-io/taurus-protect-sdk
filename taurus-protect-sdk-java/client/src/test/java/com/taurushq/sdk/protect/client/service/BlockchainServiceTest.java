package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class BlockchainServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private BlockchainService blockchainService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        blockchainService = new BlockchainService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new BlockchainService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new BlockchainService(apiClient, null));
    }

    // Note: getBlockchains() accepts all nullable parameters for filtering,
    // so no input validation tests are needed. The methods accept:
    // - null blockchain (no filter)
    // - null network (no filter)
    // - null includeBlockHeight (defaults to false)
}
