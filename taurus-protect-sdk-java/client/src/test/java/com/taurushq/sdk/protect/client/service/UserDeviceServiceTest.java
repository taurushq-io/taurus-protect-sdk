package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class UserDeviceServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private UserDeviceService userDeviceService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        userDeviceService = new UserDeviceService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new UserDeviceService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new UserDeviceService(apiClient, null));
    }

    @Test
    void getPairingStatus_throwsOnNullPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.getPairingStatus(null, "nonce"));
    }

    @Test
    void getPairingStatus_throwsOnEmptyPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.getPairingStatus("", "nonce"));
    }

    @Test
    void getPairingStatus_throwsOnNullNonce() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.getPairingStatus("pairing-123", null));
    }

    @Test
    void getPairingStatus_throwsOnEmptyNonce() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.getPairingStatus("pairing-123", ""));
    }

    @Test
    void startPairing_throwsOnNullPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.startPairing(null, "nonce", "publicKey"));
    }

    @Test
    void startPairing_throwsOnEmptyPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.startPairing("", "nonce", "publicKey"));
    }

    @Test
    void approvePairing_throwsOnNullPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.approvePairing(null, "nonce"));
    }

    @Test
    void approvePairing_throwsOnEmptyPairingId() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.approvePairing("", "nonce"));
    }

    @Test
    void approvePairing_throwsOnNullNonce() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.approvePairing("pairing-123", null));
    }

    @Test
    void approvePairing_throwsOnEmptyNonce() {
        assertThrows(IllegalArgumentException.class, () ->
                userDeviceService.approvePairing("pairing-123", ""));
    }
}
