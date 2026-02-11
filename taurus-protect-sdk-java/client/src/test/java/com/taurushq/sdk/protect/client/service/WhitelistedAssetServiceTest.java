package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class WhitelistedAssetServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private PublicKey testPublicKey;

    @BeforeEach
    void setUp() throws Exception {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();

        // Generate a real EC P-256 key to satisfy constructor validation
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair kp = kpg.generateKeyPair();
        testPublicKey = kp.getPublic();
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(null, apiExceptionMapper,
                        Collections.singletonList(testPublicKey), 1));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(apiClient, null,
                        Collections.singletonList(testPublicKey), 1));
    }

    @Test
    void constructor_throwsOnNullSuperAdminKeys() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper, null, 1));
    }

    @Test
    void constructor_throwsOnEmptySuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper,
                        Collections.emptyList(), 1));
    }

    @Test
    void constructor_throwsOnZeroMinValidSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper,
                        Collections.singletonList(testPublicKey), 0));
    }

    @Test
    void getWhitelistedAsset_throwsOnZeroId() {
        WhitelistedAssetService service = new WhitelistedAssetService(
                apiClient, apiExceptionMapper,
                Collections.singletonList(testPublicKey), 1);
        assertThrows(IllegalArgumentException.class, () ->
                service.getWhitelistedAsset(0));
    }

    @Test
    void getWhitelistedAsset_throwsOnNegativeId() {
        WhitelistedAssetService service = new WhitelistedAssetService(
                apiClient, apiExceptionMapper,
                Collections.singletonList(testPublicKey), 1);
        assertThrows(IllegalArgumentException.class, () ->
                service.getWhitelistedAsset(-1));
    }
}
