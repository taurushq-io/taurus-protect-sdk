package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when the authenticated user lacks permission for the requested operation.
 * <p>
 * This corresponds to HTTP 403 Forbidden responses. Common causes include:
 * <ul>
 *   <li>Insufficient user permissions</li>
 *   <li>Operation not allowed for the user's role</li>
 *   <li>Resource access denied</li>
 *   <li>Visibility group restrictions</li>
 * </ul>
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getRequestService().approveRequests(requests, privateKey);
 * } catch (AuthorizationException e) {
 *     System.out.println("You don't have permission to approve requests");
 * }
 * }</pre>
 */
public class AuthorizationException extends ApiException {

    /**
     * Default constructor.
     */
    public AuthorizationException() {
        setCode(403);
    }

    /**
     * Constructs an AuthorizationException with the specified message.
     *
     * @param message the error message
     */
    public AuthorizationException(String message) {
        super(message, 403);
    }

    /**
     * Constructs an AuthorizationException with full details.
     *
     * @param message   the error message
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public AuthorizationException(String message, String error, String errorCode) {
        super(message, 403, error, errorCode);
    }

    @Override
    public boolean isRetryable() {
        return false;
    }
}
