package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class VisibilityGroupServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private VisibilityGroupService visibilityGroupService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        visibilityGroupService = new VisibilityGroupService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new VisibilityGroupService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new VisibilityGroupService(apiClient, null));
    }

    @Test
    void getUsersByVisibilityGroup_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                visibilityGroupService.getUsersByVisibilityGroup(null));
    }

    @Test
    void getUsersByVisibilityGroup_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                visibilityGroupService.getUsersByVisibilityGroup(""));
    }
}
