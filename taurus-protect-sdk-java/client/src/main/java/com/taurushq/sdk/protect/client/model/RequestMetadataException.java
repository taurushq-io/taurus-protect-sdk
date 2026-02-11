package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when request metadata cannot be parsed or extracted.
 * <p>
 * This exception is thrown when attempting to extract data from
 * {@link RequestMetadata} and the expected fields are not found
 * or cannot be parsed.
 * <p>
 * Common causes include:
 * <ul>
 *   <li>Missing required fields in the metadata payload</li>
 *   <li>Malformed JSON structure</li>
 *   <li>Type mismatches when parsing values</li>
 * </ul>
 *
 * @see RequestMetadata
 */
public class RequestMetadataException extends Exception {

    /**
     * Constructs a new exception with the specified detail message.
     *
     * @param message the detail message describing what metadata field
     *                could not be extracted or parsed
     */
    public RequestMetadataException(String message) {
        super(message);
    }

    /**
     * Constructs a new exception with the specified detail message and cause.
     *
     * @param message the detail message describing what metadata field
     *                could not be extracted or parsed
     * @param cause   the underlying cause of the exception (e.g., a JSON parsing error)
     */
    public RequestMetadataException(String message, Throwable cause) {
        super(message, cause);
    }
}
