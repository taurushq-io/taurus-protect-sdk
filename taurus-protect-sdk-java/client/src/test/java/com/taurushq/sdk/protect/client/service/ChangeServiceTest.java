package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.CreateChangeRequest;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ChangeServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ChangeService changeService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        changeService = new ChangeService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ChangeService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ChangeService(apiClient, null));
    }

    @Test
    void getChange_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.getChange(null));
    }

    @Test
    void getChange_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.getChange(""));
    }

    @Test
    void approveChange_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.approveChange(null));
    }

    @Test
    void approveChange_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.approveChange(""));
    }

    @Test
    void approveChanges_throwsOnNullIds() {
        assertThrows(NullPointerException.class, () ->
                changeService.approveChanges(null));
    }

    @Test
    void approveChanges_throwsOnEmptyIds() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.approveChanges(Collections.emptyList()));
    }

    @Test
    void rejectChange_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.rejectChange(null));
    }

    @Test
    void rejectChange_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.rejectChange(""));
    }

    @Test
    void rejectChanges_throwsOnNullIds() {
        assertThrows(NullPointerException.class, () ->
                changeService.rejectChanges(null));
    }

    @Test
    void rejectChanges_throwsOnEmptyIds() {
        assertThrows(IllegalArgumentException.class, () ->
                changeService.rejectChanges(Collections.emptyList()));
    }

    @Test
    void createChange_throwsOnNullRequest() {
        assertThrows(NullPointerException.class, () ->
                changeService.createChange(null));
    }

    @Test
    void createChange_throwsOnNullAction() {
        CreateChangeRequest request = new CreateChangeRequest();
        request.setEntity("businessrule");
        assertThrows(IllegalArgumentException.class, () ->
                changeService.createChange(request));
    }

    @Test
    void createChange_throwsOnEmptyAction() {
        CreateChangeRequest request = new CreateChangeRequest();
        request.setAction("");
        request.setEntity("businessrule");
        assertThrows(IllegalArgumentException.class, () ->
                changeService.createChange(request));
    }

    @Test
    void createChange_throwsOnNullEntity() {
        CreateChangeRequest request = new CreateChangeRequest();
        request.setAction("update");
        assertThrows(IllegalArgumentException.class, () ->
                changeService.createChange(request));
    }

    @Test
    void createChange_throwsOnEmptyEntity() {
        CreateChangeRequest request = new CreateChangeRequest();
        request.setAction("update");
        request.setEntity("");
        assertThrows(IllegalArgumentException.class, () ->
                changeService.createChange(request));
    }
}
