package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ReservationServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ReservationService reservationService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        reservationService = new ReservationService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ReservationService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ReservationService(apiClient, null));
    }

    @Test
    void getReservation_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                reservationService.getReservation(null));
    }

    @Test
    void getReservation_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                reservationService.getReservation(""));
    }

    @Test
    void getReservationUtxo_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                reservationService.getReservationUtxo(null));
    }

    @Test
    void getReservationUtxo_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                reservationService.getReservationUtxo(""));
    }

    // Note: getReservations() accepts all nullable parameters for filtering,
    // so no input validation tests are needed for the method parameters.
}
