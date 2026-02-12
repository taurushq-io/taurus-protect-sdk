package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when the API rate limit has been exceeded.
 * <p>
 * This corresponds to HTTP 429 Too Many Requests responses. The client should
 * wait before retrying the request.
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getWalletService().getWallets();
 * } catch (RateLimitException e) {
 *     Thread.sleep(e.getRetryAfterMs());
 *     // Retry the request
 * }
 * }</pre>
 */
public class RateLimitException extends ApiException {

    /**
     * Suggested time to wait before retrying the request, in milliseconds.
     */
    private long retryAfterMs;

    /**
     * Default constructor.
     */
    public RateLimitException() {
        setCode(429);
        this.retryAfterMs = 1000; // Default 1 second
    }

    /**
     * Constructs a RateLimitException with the specified message.
     *
     * @param message the error message
     */
    public RateLimitException(String message) {
        super(message, 429);
        this.retryAfterMs = 1000;
    }

    /**
     * Constructs a RateLimitException with message and retry delay.
     *
     * @param message      the error message
     * @param retryAfterMs the suggested retry delay in milliseconds
     */
    public RateLimitException(String message, long retryAfterMs) {
        super(message, 429);
        this.retryAfterMs = retryAfterMs;
    }

    /**
     * Constructs a RateLimitException with full details.
     *
     * @param message      the error message
     * @param error        the error description
     * @param errorCode    the application-specific error code
     * @param retryAfterMs the suggested retry delay in milliseconds
     */
    public RateLimitException(String message, String error, String errorCode, long retryAfterMs) {
        super(message, 429, error, errorCode);
        this.retryAfterMs = retryAfterMs;
    }

    /**
     * Gets the suggested time to wait before retrying, in milliseconds.
     *
     * @return retry delay in milliseconds
     */
    public long getRetryAfterMs() {
        return retryAfterMs;
    }

    /**
     * Sets the retry delay.
     *
     * @param retryAfterMs the retry delay in milliseconds
     */
    public void setRetryAfterMs(long retryAfterMs) {
        this.retryAfterMs = retryAfterMs;
    }

    @Override
    public boolean isRetryable() {
        return true;
    }

    @Override
    public long getSuggestedRetryDelayMs() {
        return retryAfterMs;
    }
}
