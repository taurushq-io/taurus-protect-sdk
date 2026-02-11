package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ActionServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ActionService actionService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        actionService = new ActionService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ActionService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ActionService(apiClient, null));
    }

    @Test
    void getAction_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                actionService.getAction(null));
    }

    @Test
    void getAction_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                actionService.getAction(""));
    }
}
