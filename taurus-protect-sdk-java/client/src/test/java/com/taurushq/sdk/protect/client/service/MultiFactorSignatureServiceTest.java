package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMultiFactorSignaturesEntityType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class MultiFactorSignatureServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private MultiFactorSignatureService multiFactorSignatureService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        multiFactorSignatureService = new MultiFactorSignatureService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new MultiFactorSignatureService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new MultiFactorSignatureService(apiClient, null));
    }

    @Test
    void getMultiFactorSignatureInfo_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.getMultiFactorSignatureInfo(null));
    }

    @Test
    void getMultiFactorSignatureInfo_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.getMultiFactorSignatureInfo(""));
    }

    @Test
    void createMultiFactorSignatures_throwsOnNullEntityIds() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.createMultiFactorSignatures(
                        null, TgvalidatordMultiFactorSignaturesEntityType.REQUEST));
    }

    @Test
    void createMultiFactorSignatures_throwsOnEmptyEntityIds() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.createMultiFactorSignatures(
                        Collections.emptyList(), TgvalidatordMultiFactorSignaturesEntityType.REQUEST));
    }

    @Test
    void createMultiFactorSignatures_throwsOnNullEntityType() {
        assertThrows(NullPointerException.class, () ->
                multiFactorSignatureService.createMultiFactorSignatures(
                        Arrays.asList("id1", "id2"), null));
    }

    @Test
    void approveMultiFactorSignature_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.approveMultiFactorSignature(null, "sig", "comment"));
    }

    @Test
    void approveMultiFactorSignature_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.approveMultiFactorSignature("", "sig", "comment"));
    }

    @Test
    void approveMultiFactorSignature_throwsOnNullSignature() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.approveMultiFactorSignature("id", null, "comment"));
    }

    @Test
    void approveMultiFactorSignature_throwsOnEmptySignature() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.approveMultiFactorSignature("id", "", "comment"));
    }

    @Test
    void rejectMultiFactorSignature_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.rejectMultiFactorSignature(null, "comment"));
    }

    @Test
    void rejectMultiFactorSignature_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                multiFactorSignatureService.rejectMultiFactorSignature("", "comment"));
    }
}
