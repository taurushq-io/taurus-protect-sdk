package com.taurushq.sdk.protect.client.validation;

import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.client.model.RequestMetadataAmount;
import com.taurushq.sdk.protect.client.model.RequestMetadataException;

import java.util.function.Consumer;
import java.util.function.Predicate;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Fluent validator for {@link Request} metadata fields.
 * <p>
 * This validator provides a declarative API for verifying request metadata
 * before approving a request. It collects all validation errors rather than
 * failing fast, allowing callers to see all issues at once.
 * <p>
 * Example usage:
 * <pre>{@code
 * Request request = client.getRequestService().getRequest(6286272);
 *
 * RequestValidator.of(request)
 *     .expectRequestId(6286272)
 *     .expectCurrency("XLM")
 *     .expectSourceAddress("GBLNAHS...")
 *     .expectDestinationAddress("GDWIH4J...")
 *     .expectAmount(amt -> amt
 *         .valueFrom(2)
 *         .decimals(7)
 *         .currencyTo("CHF")
 *         .rateBetween(0.01, 0.2))
 *     .validate();
 *
 * client.getRequestService().approveRequests(Arrays.asList(request), privateKey);
 * }</pre>
 *
 * @see AmountValidator
 * @see ValidationResult
 * @see RequestValidationException
 */
public final class RequestValidator {

    private final Request request;
    private final RequestMetadata metadata;
    private final ValidationResult result;

    private RequestValidator(Request request) {
        this.request = checkNotNull(request, "request cannot be null");
        this.metadata = checkNotNull(request.getMetadata(), "request metadata cannot be null");
        this.result = new ValidationResult();
    }

    /**
     * Creates a new validator for the given request.
     *
     * @param request the request to validate
     * @return a new validator
     * @throws NullPointerException if request or its metadata is null
     */
    public static RequestValidator of(Request request) {
        return new RequestValidator(request);
    }

    /**
     * Validates that the request ID matches the expected value.
     *
     * @param expected the expected request ID
     * @return this validator for chaining
     */
    public RequestValidator expectRequestId(long expected) {
        try {
            long actual = metadata.getRequestId();
            if (actual != expected) {
                result.addError("requestId", "value mismatch", expected, actual);
            }
        } catch (RequestMetadataException e) {
            result.addError("requestId", "field not found in metadata");
        }
        return this;
    }

    /**
     * Validates that the currency matches the expected value.
     *
     * @param expected the expected currency
     * @return this validator for chaining
     */
    public RequestValidator expectCurrency(String expected) {
        try {
            String actual = metadata.getCurrency();
            if (!equalsNullSafe(expected, actual)) {
                result.addError("currency", "value mismatch", expected, actual);
            }
        } catch (RequestMetadataException e) {
            result.addError("currency", "field not found in metadata");
        }
        return this;
    }

    /**
     * Validates that the rules key matches the expected value.
     *
     * @param expected the expected rules key
     * @return this validator for chaining
     */
    public RequestValidator expectRulesKey(String expected) {
        try {
            String actual = metadata.getRulesKey();
            if (!equalsNullSafe(expected, actual)) {
                result.addError("rulesKey", "value mismatch", expected, actual);
            }
        } catch (RequestMetadataException e) {
            result.addError("rulesKey", "field not found in metadata");
        }
        return this;
    }

    /**
     * Validates that the source address matches the expected value.
     *
     * @param expected the expected source address
     * @return this validator for chaining
     */
    public RequestValidator expectSourceAddress(String expected) {
        try {
            String actual = metadata.getSourceAddress();
            if (!equalsNullSafe(expected, actual)) {
                result.addError("sourceAddress", "value mismatch", expected, actual);
            }
        } catch (RequestMetadataException e) {
            result.addError("sourceAddress", "field not found in metadata");
        }
        return this;
    }

    /**
     * Validates that the destination address matches the expected value.
     *
     * @param expected the expected destination address
     * @return this validator for chaining
     */
    public RequestValidator expectDestinationAddress(String expected) {
        try {
            String actual = metadata.getDestinationAddress();
            if (!equalsNullSafe(expected, actual)) {
                result.addError("destinationAddress", "value mismatch", expected, actual);
            }
        } catch (RequestMetadataException e) {
            result.addError("destinationAddress", "field not found in metadata");
        }
        return this;
    }

    /**
     * Validates amount fields using a nested {@link AmountValidator}.
     *
     * @param amountConfig consumer that configures the amount validator
     * @return this validator for chaining
     */
    public RequestValidator expectAmount(Consumer<AmountValidator> amountConfig) {
        try {
            RequestMetadataAmount amount = metadata.getAmount();
            AmountValidator amountValidator = new AmountValidator(amount, result);
            amountConfig.accept(amountValidator);
        } catch (RequestMetadataException e) {
            result.addError("amount", "field not found in metadata");
        }
        return this;
    }

    /**
     * Adds a custom validation rule.
     * <p>
     * Use this for validations not covered by the built-in methods.
     *
     * @param name  a descriptive name for the validation (used in error messages)
     * @param check predicate that returns true if validation passes
     * @return this validator for chaining
     */
    public RequestValidator require(String name, Predicate<RequestMetadata> check) {
        checkNotNull(name, "validation name cannot be null");
        checkNotNull(check, "validation predicate cannot be null");
        try {
            if (!check.test(metadata)) {
                result.addError(name, "custom validation failed");
            }
        } catch (Exception e) {
            result.addError(name, "validation threw exception: " + e.getMessage());
        }
        return this;
    }

    /**
     * Validates the request status matches the expected value.
     *
     * @param expected the expected status string
     * @return this validator for chaining
     */
    public RequestValidator expectStatus(String expected) {
        String actual = request.getStatus() != null ? request.getStatus().toString() : null;
        if (!equalsNullSafe(expected, actual)) {
            result.addError("status", "value mismatch", expected, actual);
        }
        return this;
    }

    /**
     * Validates the request type matches the expected value.
     *
     * @param expected the expected type
     * @return this validator for chaining
     */
    public RequestValidator expectType(String expected) {
        String actual = request.getType();
        if (!equalsNullSafe(expected, actual)) {
            result.addError("type", "value mismatch", expected, actual);
        }
        return this;
    }

    /**
     * Executes all validations and throws if any failed.
     *
     * @throws RequestValidationException if any validation failed
     */
    public void validate() throws RequestValidationException {
        if (result.hasErrors()) {
            throw new RequestValidationException(result);
        }
    }

    /**
     * Executes all validations and returns the result without throwing.
     *
     * @return the validation result
     */
    public ValidationResult check() {
        return result;
    }

    /**
     * Checks if all validations passed.
     *
     * @return true if no validation errors
     */
    public boolean isValid() {
        return result.isValid();
    }

    /**
     * Gets the request being validated.
     *
     * @return the request
     */
    public Request getRequest() {
        return request;
    }

    private static boolean equalsNullSafe(String expected, String actual) {
        if (expected == null && actual == null) {
            return true;
        }
        if (expected == null || actual == null) {
            return false;
        }
        return expected.equals(actual);
    }
}
