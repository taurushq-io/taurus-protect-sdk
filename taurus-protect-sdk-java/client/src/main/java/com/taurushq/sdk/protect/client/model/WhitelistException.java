package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when whitelist operations fail.
 * <p>
 * This exception is thrown when there are issues with whitelist operations,
 * such as validation failures, signature verification errors, or when
 * attempting operations on non-existent whitelist entries.
 * <p>
 * Common causes include:
 * <ul>
 *   <li>Invalid address format for the target blockchain</li>
 *   <li>Signature verification failure</li>
 *   <li>Missing required approvals</li>
 *   <li>Whitelist entry not found</li>
 * </ul>
 *
 * @see WhitelistedAddress
 * @see WhitelistSignature
 */
public class WhitelistException extends Exception {

    /**
     * Constructs a new exception with the specified detail message.
     *
     * @param message the detail message describing the whitelist operation failure
     */
    public WhitelistException(String message) {
        super(message);
    }

    /**
     * Constructs a new exception with the specified detail message and cause.
     *
     * @param message the detail message describing the whitelist operation failure
     * @param cause   the underlying cause of the exception
     */
    public WhitelistException(String message, Throwable cause) {
        super(message, cause);
    }
}
