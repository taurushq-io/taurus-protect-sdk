package com.taurushq.sdk.protect.client.mapper;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.AuthenticationException;
import com.taurushq.sdk.protect.client.model.AuthorizationException;
import com.taurushq.sdk.protect.client.model.NotFoundException;
import com.taurushq.sdk.protect.client.model.RateLimitException;
import com.taurushq.sdk.protect.client.model.ServerException;
import com.taurushq.sdk.protect.client.model.ValidationException;


/**
 * Maps OpenAPI exceptions to typed SDK exceptions.
 * <p>
 * This mapper creates specific exception types based on HTTP status codes,
 * allowing clients to handle different error categories appropriately.
 */
@SuppressWarnings("PMD.EmptyCatchBlock")
public class ApiExceptionMapper {

    private static final long DEFAULT_RATE_LIMIT_RETRY_MS = 1000;

    private static final Gson GSON = new Gson();

    /**
     * Instantiates a new Api exception mapper.
     */
    public ApiExceptionMapper() {
        // No instance state needed - Gson is shared via static field
    }


    /**
     * Converts an OpenAPI exception to a typed SDK exception.
     * <p>
     * The returned exception type is determined by the HTTP status code:
     * <ul>
     *   <li>400 → {@link ValidationException}</li>
     *   <li>401 → {@link AuthenticationException}</li>
     *   <li>403 → {@link AuthorizationException}</li>
     *   <li>404 → {@link NotFoundException}</li>
     *   <li>429 → {@link RateLimitException}</li>
     *   <li>5xx → {@link ServerException}</li>
     *   <li>Other → {@link ApiException}</li>
     * </ul>
     *
     * @param e the OpenAPI exception
     * @return a typed SDK exception
     */
    @SuppressWarnings("PMD.EmptyCatchBlock")
    public ApiException toApiException(com.taurushq.sdk.protect.openapi.ApiException e) {
        int statusCode = e.getCode();

        // Try to parse the response body for additional error details
        ApiException parsed = parseResponseBody(e);

        // Create the appropriate typed exception
        ApiException result = createTypedException(statusCode, parsed, e);

        // Set the original exception for debugging
        result.setOriginalException(e);

        return result;
    }

    /**
     * Parses the response body to extract error details.
     */
    private ApiException parseResponseBody(com.taurushq.sdk.protect.openapi.ApiException e) {
        ApiException parsed = new ApiException();
        parsed.setError("Unknown");
        parsed.setCode(e.getCode());
        parsed.setMessage(e.getMessage());

        if (!Strings.isNullOrEmpty(e.getResponseBody())) {
            try {
                ApiException fromJson = GSON.fromJson(e.getResponseBody(), ApiException.class);
                if (fromJson != null) {
                    parsed = fromJson;
                    // Ensure code is set from HTTP status if not in body
                    if (parsed.getCode() == 0) {
                        parsed.setCode(e.getCode());
                    }
                }
            } catch (JsonSyntaxException ex) {
                // Keep the defaults if parsing fails
            }
        }

        return parsed;
    }

    /**
     * Creates a typed exception based on the HTTP status code.
     */
    private ApiException createTypedException(int statusCode, ApiException parsed,
            com.taurushq.sdk.protect.openapi.ApiException original) {

        ApiException result;

        switch (statusCode) {
            case 400:
                result = new ValidationException(parsed.getMessage(), parsed.getError(),
                        parsed.getErrorCode());
                break;
            case 401:
                result = new AuthenticationException(parsed.getMessage(), parsed.getError(),
                        parsed.getErrorCode());
                break;
            case 403:
                result = new AuthorizationException(parsed.getMessage(), parsed.getError(),
                        parsed.getErrorCode());
                break;
            case 404:
                result = new NotFoundException(parsed.getMessage(), parsed.getError(),
                        parsed.getErrorCode());
                break;
            case 429:
                long retryAfter = parseRetryAfterHeader(original);
                RateLimitException rle = new RateLimitException(parsed.getMessage(),
                        parsed.getError(), parsed.getErrorCode(), retryAfter);
                result = rle;
                break;
            default:
                if (statusCode >= 500) {
                    result = new ServerException(parsed.getMessage(), statusCode,
                            parsed.getError(), parsed.getErrorCode());
                } else {
                    // For other status codes, use the base ApiException
                    result = parsed;
                }
                break;
        }

        return result;
    }

    /**
     * Parses the Retry-After header to determine wait time.
     */
    private long parseRetryAfterHeader(com.taurushq.sdk.protect.openapi.ApiException e) {
        if (e.getResponseHeaders() != null) {
            java.util.List<String> retryAfterValues = e.getResponseHeaders().get("Retry-After");
            if (retryAfterValues != null && !retryAfterValues.isEmpty()) {
                try {
                    // Retry-After can be seconds or a date; we handle seconds here
                    int seconds = Integer.parseInt(retryAfterValues.get(0));
                    return seconds * 1000L;
                } catch (NumberFormatException ex) {
                    // Could be a date format, use default
                }
            }
        }
        return DEFAULT_RATE_LIMIT_RETRY_MS;
    }
}
