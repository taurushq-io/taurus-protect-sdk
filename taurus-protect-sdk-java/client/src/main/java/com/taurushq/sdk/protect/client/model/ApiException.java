package com.taurushq.sdk.protect.client.model;

/**
 * Base exception for all Taurus-PROTECT API errors.
 * <p>
 * This exception captures HTTP status codes, error messages, and provides
 * helper methods for common error handling patterns.
 * <p>
 * Specific subclasses provide more context for different error categories:
 * <ul>
 *   <li>{@link ValidationException} - 400 Bad Request errors</li>
 *   <li>{@link AuthenticationException} - 401 Unauthorized errors</li>
 *   <li>{@link AuthorizationException} - 403 Forbidden errors</li>
 *   <li>{@link NotFoundException} - 404 Not Found errors</li>
 *   <li>{@link RateLimitException} - 429 Too Many Requests errors</li>
 * </ul>
 */
public class ApiException extends Exception {

    private String error;
    private String message;
    private int code;
    private String errorCode;
    private Exception originalException;

    /**
     * Default constructor.
     */
    public ApiException() {
        // Default constructor
    }

    /**
     * Constructs an ApiException with the specified message.
     *
     * @param message the error message
     */
    public ApiException(String message) {
        super(message);
        this.message = message;
    }

    /**
     * Constructs an ApiException with the specified message and HTTP status code.
     *
     * @param message the error message
     * @param code    the HTTP status code
     */
    public ApiException(String message, int code) {
        super(message);
        this.message = message;
        this.code = code;
    }

    /**
     * Constructs an ApiException with full details.
     *
     * @param message   the error message
     * @param code      the HTTP status code
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public ApiException(String message, int code, String error, String errorCode) {
        super(message);
        this.message = message;
        this.code = code;
        this.error = error;
        this.errorCode = errorCode;
    }

    /**
     * Determines if this error is potentially retryable.
     * <p>
     * Returns true for:
     * <ul>
     *   <li>429 Too Many Requests (rate limited)</li>
     *   <li>500+ Server errors (transient failures)</li>
     * </ul>
     *
     * @return true if the request might succeed on retry
     */
    public boolean isRetryable() {
        return code == 429 || code >= 500;
    }

    /**
     * Determines if this is a client error (4xx).
     *
     * @return true if this is a client error
     */
    public boolean isClientError() {
        return code >= 400 && code < 500;
    }

    /**
     * Determines if this is a server error (5xx).
     *
     * @return true if this is a server error
     */
    public boolean isServerError() {
        return code >= 500;
    }

    /**
     * Gets the suggested retry delay in milliseconds.
     * <p>
     * For rate limit errors, this may return a specific delay.
     * For server errors, returns a default backoff value.
     * For non-retryable errors, returns 0.
     *
     * @return suggested delay in milliseconds, or 0 if not retryable
     */
    public long getSuggestedRetryDelayMs() {
        if (code == 429) {
            return 1000; // Default 1 second for rate limits
        }
        if (code >= 500) {
            return 5000; // Default 5 seconds for server errors
        }
        return 0;
    }

    /**
     * Gets original exception.
     *
     * @return the original exception
     */
    public Exception getOriginalException() {
        return originalException;
    }

    /**
     * Sets original exception.
     *
     * @param originalException the original exception
     */
    public void setOriginalException(Exception originalException) {
        this.originalException = originalException;
    }

    /**
     * Gets error.
     *
     * @return the error
     */
    public String getError() {
        return error;
    }

    /**
     * Sets error.
     *
     * @param error the error
     */
    public void setError(String error) {
        this.error = error;
    }

    @Override
    public String getMessage() {
        return message;
    }

    /**
     * Sets message.
     *
     * @param message the message
     */
    public void setMessage(String message) {
        this.message = message;
    }

    /**
     * Gets code.
     *
     * @return the code
     */
    public int getCode() {
        return code;
    }

    /**
     * Sets code.
     *
     * @param code the code
     */
    public void setCode(int code) {
        this.code = code;
    }

    /**
     * Gets error code.
     *
     * @return the error code
     */
    public String getErrorCode() {
        return errorCode;
    }

    /**
     * Sets error code.
     *
     * @param errorCode the error code
     */
    public void setErrorCode(String errorCode) {
        this.errorCode = errorCode;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append(getClass().getSimpleName());
        sb.append("{code=").append(code);
        if (errorCode != null) {
            sb.append(", errorCode='").append(errorCode).append("'");
        }
        if (message != null) {
            sb.append(", message='").append(message).append("'");
        }
        if (error != null) {
            sb.append(", error='").append(error).append("'");
        }
        sb.append("}");
        return sb.toString();
    }
}
