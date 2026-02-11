package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ApiExceptionMapperTest {

    private ApiExceptionMapper mapper;

    @BeforeEach
    void setUp() {
        mapper = new ApiExceptionMapper();
    }

    @Test
    void toApiException_withValidJsonResponseBody_parsesAllFields() {
        // Given
        String jsonBody = "{\"error\":\"NotFound\",\"message\":\"Resource not found\",\"code\":404,\"errorCode\":\"ERR_001\"}";
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(404, "Not Found", null, jsonBody);

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("NotFound", result.getError());
        assertEquals("Resource not found", result.getMessage());
        assertEquals(404, result.getCode());
        assertEquals("ERR_001", result.getErrorCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withNullResponseBody_usesDefaultValues() {
        // Given
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(500, "Internal Server Error");

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("Unknown", result.getError());
        // The message includes the full OpenAPI exception message format
        assertTrue(result.getMessage().contains("Internal Server Error"));
        assertEquals(500, result.getCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withEmptyResponseBody_usesDefaultValues() {
        // Given
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(400, "Bad Request", null, "");

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("Unknown", result.getError());
        assertTrue(result.getMessage().contains("Bad Request"));
        assertEquals(400, result.getCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withInvalidJson_fallsBackToDefaults() {
        // Given
        String invalidJson = "this is not valid json";
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(500, "Server Error", null, invalidJson);

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("Unknown", result.getError());
        assertTrue(result.getMessage().contains("Server Error"));
        assertEquals(500, result.getCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withPartialJson_mapsAvailableFields() {
        // Given
        String partialJson = "{\"error\":\"ValidationError\",\"message\":\"Invalid input\"}";
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(422, "Validation Failed", null, partialJson);

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("ValidationError", result.getError());
        assertEquals("Invalid input", result.getMessage());
        assertNull(result.getErrorCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withHtmlResponseBody_fallsBackToDefaults() {
        // Given - sometimes servers return HTML error pages
        String htmlBody = "<html><body><h1>500 Internal Server Error</h1></body></html>";
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(500, "Internal Server Error", null, htmlBody);

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("Unknown", result.getError());
        assertTrue(result.getMessage().contains("Internal Server Error"));
        assertEquals(500, result.getCode());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_preservesOriginalException() {
        // Given
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(401, "Unauthorized");

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertNotNull(result.getOriginalException());
        assertSame(openApiException, result.getOriginalException());
    }

    @Test
    void toApiException_withJsonArrayResponseBody_fallsBackToDefaults() {
        // Given - response body is valid JSON but wrong structure
        String jsonArrayBody = "[\"error1\", \"error2\"]";
        com.taurushq.sdk.protect.openapi.ApiException openApiException =
                new com.taurushq.sdk.protect.openapi.ApiException(400, "Bad Request", null, jsonArrayBody);

        // When
        ApiException result = mapper.toApiException(openApiException);

        // Then
        assertEquals("Unknown", result.getError());
        assertTrue(result.getMessage().contains("Bad Request"));
        assertSame(openApiException, result.getOriginalException());
    }
}
