package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when authentication fails.
 * <p>
 * This corresponds to HTTP 401 Unauthorized responses. Common causes include:
 * <ul>
 *   <li>Invalid API key</li>
 *   <li>Invalid API secret</li>
 *   <li>Expired credentials</li>
 *   <li>Malformed authentication header</li>
 * </ul>
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getWalletService().getWallets();
 * } catch (AuthenticationException e) {
 *     System.out.println("Authentication failed - check your API credentials");
 * }
 * }</pre>
 */
public class AuthenticationException extends ApiException {

    /**
     * Default constructor.
     */
    public AuthenticationException() {
        setCode(401);
    }

    /**
     * Constructs an AuthenticationException with the specified message.
     *
     * @param message the error message
     */
    public AuthenticationException(String message) {
        super(message, 401);
    }

    /**
     * Constructs an AuthenticationException with full details.
     *
     * @param message   the error message
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public AuthenticationException(String message, String error, String errorCode) {
        super(message, 401, error, errorCode);
    }

    @Override
    public boolean isRetryable() {
        return false;
    }
}
