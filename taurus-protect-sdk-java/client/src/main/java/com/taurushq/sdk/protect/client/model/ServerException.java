package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when the API encounters an internal server error.
 * <p>
 * This corresponds to HTTP 5xx responses. These errors are typically
 * transient and may succeed on retry.
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getWalletService().getWallets();
 * } catch (ServerException e) {
 *     if (e.isRetryable()) {
 *         Thread.sleep(e.getSuggestedRetryDelayMs());
 *         // Retry the request
 *     }
 * }
 * }</pre>
 */
public class ServerException extends ApiException {

    /**
     * Default constructor.
     */
    public ServerException() {
        setCode(500);
    }

    /**
     * Constructs a ServerException with the specified message and code.
     *
     * @param message the error message
     * @param code    the HTTP status code (5xx)
     */
    public ServerException(String message, int code) {
        super(message, code);
    }

    /**
     * Constructs a ServerException with full details.
     *
     * @param message   the error message
     * @param code      the HTTP status code (5xx)
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public ServerException(String message, int code, String error, String errorCode) {
        super(message, code, error, errorCode);
    }

    @Override
    public boolean isRetryable() {
        return true;
    }

    @Override
    public long getSuggestedRetryDelayMs() {
        return 5000; // 5 seconds for server errors
    }
}
