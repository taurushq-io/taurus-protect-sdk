package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when a requested resource is not found.
 * <p>
 * This corresponds to HTTP 404 Not Found responses. Common causes include:
 * <ul>
 *   <li>Invalid wallet ID</li>
 *   <li>Invalid address ID</li>
 *   <li>Invalid request ID</li>
 *   <li>Resource was deleted</li>
 * </ul>
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     Wallet wallet = client.getWalletService().getWallet(walletId);
 * } catch (NotFoundException e) {
 *     System.out.println("Wallet not found: " + walletId);
 * }
 * }</pre>
 */
public class NotFoundException extends ApiException {

    /**
     * Default constructor.
     */
    public NotFoundException() {
        setCode(404);
    }

    /**
     * Constructs a NotFoundException with the specified message.
     *
     * @param message the error message
     */
    public NotFoundException(String message) {
        super(message, 404);
    }

    /**
     * Constructs a NotFoundException with full details.
     *
     * @param message   the error message
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public NotFoundException(String message, String error, String errorCode) {
        super(message, 404, error, errorCode);
    }

    @Override
    public boolean isRetryable() {
        return false;
    }
}
