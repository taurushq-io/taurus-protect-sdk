package com.taurushq.sdk.protect.client.model;

/**
 * Exception thrown when the governance rules container carries something this SDK
 * version cannot interpret, so no verification decision taken against it can be trusted.
 * <p>
 * Deliberately distinct from a per-row failure. A caller listing whitelisted addresses
 * excludes a row that fails its own integrity check and keeps the rest, but an
 * uninterpretable container invalidates EVERY row judged against it, so the whole call
 * must fail instead. Collapsing the two lets a schema-newer container empty a whitelist
 * while reporting success.
 * <p>
 * Extends {@link IntegrityException} so existing handlers keep working, while a caller
 * that wants to distinguish "this SDK is too old for this container" from "this row was
 * tampered with" can catch this type specifically.
 *
 * @see com.taurushq.sdk.protect.client.service.WhitelistedAddressService
 */
public class ContainerIntegrityException extends IntegrityException {

    /**
     * Constructs a new ContainerIntegrityException with the specified message.
     *
     * @param message the detail message
     */
    public ContainerIntegrityException(String message) {
        super(message);
    }

    /**
     * Constructs a new ContainerIntegrityException with the specified message and cause.
     *
     * @param message the detail message
     * @param cause   the cause of this exception
     */
    public ContainerIntegrityException(String message, Throwable cause) {
        super(message, cause);
    }
}
