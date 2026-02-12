package com.taurushq.sdk.protect.client.validation;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Collects validation errors from a {@link RequestValidator}.
 * <p>
 * This class aggregates all validation failures rather than failing fast,
 * allowing callers to see all issues at once.
 */
public class ValidationResult {

    private final List<ValidationError> errors;

    /**
     * Creates a new empty validation result.
     */
    public ValidationResult() {
        this.errors = new ArrayList<>();
    }

    /**
     * Adds a validation error.
     *
     * @param field    the field that failed validation
     * @param message  the error message
     * @param expected the expected value (may be null)
     * @param actual   the actual value (may be null)
     */
    public void addError(String field, String message, Object expected, Object actual) {
        errors.add(new ValidationError(field, message, expected, actual));
    }

    /**
     * Adds a validation error without expected/actual values.
     *
     * @param field   the field that failed validation
     * @param message the error message
     */
    public void addError(String field, String message) {
        errors.add(new ValidationError(field, message, null, null));
    }

    /**
     * Checks if there are any validation errors.
     *
     * @return true if there are no errors
     */
    public boolean isValid() {
        return errors.isEmpty();
    }

    /**
     * Checks if there are validation errors.
     *
     * @return true if there are errors
     */
    public boolean hasErrors() {
        return !errors.isEmpty();
    }

    /**
     * Gets all validation errors.
     *
     * @return unmodifiable list of errors
     */
    public List<ValidationError> getErrors() {
        return Collections.unmodifiableList(errors);
    }

    /**
     * Gets the number of validation errors.
     *
     * @return error count
     */
    public int getErrorCount() {
        return errors.size();
    }

    /**
     * Builds a summary message of all validation errors.
     *
     * @return summary message
     */
    public String getSummary() {
        if (errors.isEmpty()) {
            return "Validation passed";
        }
        StringBuilder sb = new StringBuilder();
        sb.append("Validation failed with ").append(errors.size()).append(" error(s):");
        for (ValidationError error : errors) {
            sb.append("\n  - ").append(error.getField()).append(": ").append(error.getMessage());
            if (error.getExpected() != null || error.getActual() != null) {
                sb.append(" (expected: ").append(error.getExpected());
                sb.append(", actual: ").append(error.getActual()).append(")");
            }
        }
        return sb.toString();
    }

    /**
     * Represents a single validation error.
     */
    public static class ValidationError {
        private final String field;
        private final String message;
        private final Object expected;
        private final Object actual;

        /**
         * Creates a new validation error.
         *
         * @param field    the field name
         * @param message  the error message
         * @param expected the expected value
         * @param actual   the actual value
         */
        public ValidationError(String field, String message, Object expected, Object actual) {
            this.field = field;
            this.message = message;
            this.expected = expected;
            this.actual = actual;
        }

        /**
         * Gets the field name.
         *
         * @return the field name
         */
        public String getField() {
            return field;
        }

        /**
         * Gets the error message.
         *
         * @return the error message
         */
        public String getMessage() {
            return message;
        }

        /**
         * Gets the expected value.
         *
         * @return the expected value, may be null
         */
        public Object getExpected() {
            return expected;
        }

        /**
         * Gets the actual value.
         *
         * @return the actual value, may be null
         */
        public Object getActual() {
            return actual;
        }

        @Override
        public String toString() {
            return "ValidationError{field='" + field + "', message='" + message
                    + "', expected=" + expected + ", actual=" + actual + "}";
        }
    }
}
