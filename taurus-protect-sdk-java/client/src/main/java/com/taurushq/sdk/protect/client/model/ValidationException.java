package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when the API rejects a request due to validation errors.
 * <p>
 * This corresponds to HTTP 400 Bad Request responses. Common causes include:
 * <ul>
 *   <li>Invalid request parameters</li>
 *   <li>Missing required fields</li>
 *   <li>Invalid data formats</li>
 *   <li>Business rule violations</li>
 * </ul>
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getWalletService().createWallet(...);
 * } catch (ValidationException e) {
 *     System.out.println("Invalid request: " + e.getMessage());
 *     System.out.println("Error code: " + e.getErrorCode());
 * }
 * }</pre>
 */
public class ValidationException extends ApiException {

    /**
     * Default constructor.
     */
    public ValidationException() {
        setCode(400);
    }

    /**
     * Constructs a ValidationException with the specified message.
     *
     * @param message the error message
     */
    public ValidationException(String message) {
        super(message, 400);
    }

    /**
     * Constructs a ValidationException with full details.
     *
     * @param message   the error message
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public ValidationException(String message, String error, String errorCode) {
        super(message, 400, error, errorCode);
    }

    @Override
    public boolean isRetryable() {
        return false;
    }
}
