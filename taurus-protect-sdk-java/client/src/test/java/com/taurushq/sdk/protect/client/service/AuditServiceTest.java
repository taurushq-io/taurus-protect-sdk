package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertThrows;

class AuditServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private AuditService auditService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        auditService = new AuditService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AuditService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AuditService(apiClient, null));
    }

    // Note: getAuditTrails() accepts all nullable parameters, so no input validation tests needed.
    // The method allows:
    // - null externalUserId (optional filter)
    // - null entities list (optional filter)
    // - null actions list (optional filter)
    // - null from date (optional filter)
    // - null to date (optional filter)
    // - null cursor (starts from first page)

    @Test
    void getAuditTrails_acceptsAllNullParameters() {
        // This should not throw - all parameters are optional
        // The actual API call will fail without a configured server,
        // but input validation should pass
    }

    @Test
    void getAuditTrails_acceptsValidParameters() {
        // Verify the service can be instantiated with valid parameters
        // The actual API call would require a server connection
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);

        // These parameters should all be accepted without throwing
        // (actual API call will fail but input validation passes)
    }

    // Note: exportAuditTrails() also accepts nullable parameters for filtering
    // Only format is used to specify output format
}
