package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when SDK configuration is invalid.
 * <p>
 * This exception indicates that the SDK was configured with invalid or
 * incompatible settings. Common causes include:
 * <ul>
 *   <li>Missing required credentials (API key or secret)</li>
 *   <li>Invalid API host URL format</li>
 *   <li>Invalid SuperAdmin key format or encoding</li>
 *   <li>Invalid min_valid_signatures value</li>
 *   <li>Incompatible configuration options</li>
 * </ul>
 * <p>
 * This exception should be handled at application startup and typically
 * indicates a programming or deployment error that needs to be fixed.
 */
public class ConfigurationException extends Exception {

    /**
     * Constructs a new ConfigurationException with the specified message.
     *
     * @param message the detail message describing the configuration error
     */
    public ConfigurationException(String message) {
        super(message);
    }

    /**
     * Constructs a new ConfigurationException with the specified message and cause.
     *
     * @param message the detail message describing the configuration error
     * @param cause   the cause of this exception
     */
    public ConfigurationException(String message, Throwable cause) {
        super(message, cause);
    }
}
