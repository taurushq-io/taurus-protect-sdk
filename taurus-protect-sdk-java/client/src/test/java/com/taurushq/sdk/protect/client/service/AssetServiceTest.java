package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class AssetServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private AssetService assetService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        assetService = new AssetService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AssetService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AssetService(apiClient, null));
    }

    @Test
    void getAssetAddresses_throwsOnNullCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetAddresses(null));
    }

    @Test
    void getAssetAddresses_throwsOnEmptyCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetAddresses(""));
    }

    @Test
    void getAssetWallets_throwsOnNullCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetWallets(null));
    }

    @Test
    void getAssetWallets_throwsOnEmptyCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetWallets(""));
    }
}
