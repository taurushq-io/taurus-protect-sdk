package com.taurushq.sdk.protect.client.validation;

/**
 * Exception thrown when request metadata validation fails.
 * <p>
 * This is an unchecked exception because validation failures typically
 * indicate a programming error or security issue that should not be silently ignored.
 * <p>
 * Use {@link #getValidationResult()} to access all validation errors.
 */
public class RequestValidationException extends RuntimeException {

    private static final long serialVersionUID = 1L;

    private final ValidationResult validationResult;

    /**
     * Creates a new validation exception.
     *
     * @param validationResult the validation result containing all errors
     */
    public RequestValidationException(ValidationResult validationResult) {
        super(validationResult.getSummary());
        this.validationResult = validationResult;
    }

    /**
     * Creates a new validation exception with a custom message.
     *
     * @param message          the exception message
     * @param validationResult the validation result containing all errors
     */
    public RequestValidationException(String message, ValidationResult validationResult) {
        super(message);
        this.validationResult = validationResult;
    }

    /**
     * Gets the validation result containing all errors.
     *
     * @return the validation result
     */
    public ValidationResult getValidationResult() {
        return validationResult;
    }

    /**
     * Gets the number of validation errors.
     *
     * @return error count
     */
    public int getErrorCount() {
        return validationResult.getErrorCount();
    }
}
