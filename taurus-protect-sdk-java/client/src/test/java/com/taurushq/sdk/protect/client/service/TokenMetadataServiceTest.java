package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class TokenMetadataServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private TokenMetadataService tokenMetadataService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        tokenMetadataService = new TokenMetadataService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new TokenMetadataService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new TokenMetadataService(apiClient, null));
    }

    @Test
    void getERCTokenMetadata_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getERCTokenMetadata(null, "0x123", null, false, "ETH"));
    }

    @Test
    void getERCTokenMetadata_throwsOnEmptyNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getERCTokenMetadata("", "0x123", null, false, "ETH"));
    }

    @Test
    void getERCTokenMetadata_throwsOnNullContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getERCTokenMetadata("mainnet", null, null, false, "ETH"));
    }

    @Test
    void getERCTokenMetadata_throwsOnEmptyContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getERCTokenMetadata("mainnet", "", null, false, "ETH"));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata(null, "0x123", null, false, "MATIC"));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnEmptyNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata("", "0x123", null, false, "MATIC"));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnNullContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata("mainnet", null, null, false, "MATIC"));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnEmptyContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata("mainnet", "", null, false, "MATIC"));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnNullBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata("mainnet", "0x123", null, false, null));
    }

    @Test
    void getEVMERCTokenMetadata_throwsOnEmptyBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getEVMERCTokenMetadata("mainnet", "0x123", null, false, ""));
    }

    @Test
    void getFATokenMetadata_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getFATokenMetadata(null, "KT1xxx", "0", false));
    }

    @Test
    void getFATokenMetadata_throwsOnEmptyNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getFATokenMetadata("", "KT1xxx", "0", false));
    }

    @Test
    void getFATokenMetadata_throwsOnNullContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getFATokenMetadata("mainnet", null, "0", false));
    }

    @Test
    void getFATokenMetadata_throwsOnEmptyContract() {
        assertThrows(IllegalArgumentException.class, () ->
                tokenMetadataService.getFATokenMetadata("mainnet", "", "0", false));
    }
}
