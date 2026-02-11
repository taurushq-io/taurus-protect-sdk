package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when integrity verification fails.
 * <p>
 * This exception indicates that a cryptographic signature or hash verification
 * has failed, suggesting potential data tampering or corruption. Common causes include:
 * <ul>
 *   <li>Invalid cryptographic signature</li>
 *   <li>Hash mismatch after data transfer</li>
 *   <li>Tampered request or response data</li>
 *   <li>Corrupted data during transmission</li>
 * </ul>
 * <p>
 * This exception typically indicates a serious security issue that should be
 * investigated and not simply retried.
 *
 * @see com.taurushq.sdk.protect.client.helper.AddressSignatureVerifier
 */
public class IntegrityException extends SecurityException {

    /**
     * Constructs a new IntegrityException with the specified message.
     *
     * @param message the detail message
     */
    public IntegrityException(String message) {
        super(message);
    }

    /**
     * Constructs a new IntegrityException with the specified message and cause.
     *
     * @param message the detail message
     * @param cause   the cause of this exception
     */
    public IntegrityException(String message, Throwable cause) {
        super(message, cause);
    }
}
