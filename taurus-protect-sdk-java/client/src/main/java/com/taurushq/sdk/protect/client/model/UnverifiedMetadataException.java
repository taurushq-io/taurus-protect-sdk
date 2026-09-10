package com.taurushq.sdk.protect.client.model;

/**
 * Thrown when a request metadata payload is read before verification has cleared it.
 *
 * <p>Distinct from a plain {@link RequestMetadataException}, which means "this key is
 * not in the payload". Those are different facts: one says the field is absent, the
 * other says nothing about the payload can be trusted yet. Collapsing them lets a
 * verification failure read as an absent source address, which a caller might act on.
 *
 * <p>It extends {@code RequestMetadataException} so existing catch blocks keep working;
 * catch this type first when the distinction matters.
 */
public class UnverifiedMetadataException extends RequestMetadataException {

    private static final long serialVersionUID = 1L;

    /**
     * Instantiates a new unverified metadata exception.
     *
     * @param message the detail message
     */
    public UnverifiedMetadataException(final String message) {
        super(message);
    }
}
