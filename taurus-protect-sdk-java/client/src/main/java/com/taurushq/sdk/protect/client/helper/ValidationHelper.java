package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;

/**
 * Utility class for validating input parameters.
 * <p>
 * Provides consistent validation with clear error messages suitable for
 * SDK users. All methods throw {@link IllegalArgumentException} on failure.
 * <p>
 * Example usage:
 * <pre>{@code
 * public void createWallet(String name, long balance) {
 *     ValidationHelper.requireNotBlank(name, "wallet name");
 *     ValidationHelper.requirePositive(balance, "initial balance");
 *     // proceed with creation
 * }
 * }</pre>
 */
public final class ValidationHelper {

    private ValidationHelper() {
        // Static utility class
    }

    /**
     * Validates that a value is not null.
     *
     * @param value     the value to check
     * @param fieldName the name of the field (for error messages)
     * @param <T>       the type of the value
     * @return the non-null value
     * @throws IllegalArgumentException if value is null
     */
    public static <T> T requireNotNull(T value, String fieldName) {
        if (value == null) {
            throw new IllegalArgumentException(fieldName + " cannot be null");
        }
        return value;
    }

    /**
     * Validates that a string is not null or empty.
     *
     * @param value     the string to check
     * @param fieldName the name of the field (for error messages)
     * @return the non-empty string
     * @throws IllegalArgumentException if value is null or empty
     */
    public static String requireNotEmpty(String value, String fieldName) {
        if (Strings.isNullOrEmpty(value)) {
            throw new IllegalArgumentException(fieldName + " cannot be null or empty");
        }
        return value;
    }

    /**
     * Validates that a string is not null, empty, or blank (whitespace only).
     *
     * @param value     the string to check
     * @param fieldName the name of the field (for error messages)
     * @return the non-blank string
     * @throws IllegalArgumentException if value is null, empty, or blank
     */
    public static String requireNotBlank(String value, String fieldName) {
        if (Strings.isNullOrEmpty(value) || value.trim().isEmpty()) {
            throw new IllegalArgumentException(fieldName + " cannot be null, empty, or blank");
        }
        return value;
    }

    /**
     * Validates that a number is positive (greater than zero).
     *
     * @param value     the value to check
     * @param fieldName the name of the field (for error messages)
     * @return the positive value
     * @throws IllegalArgumentException if value is not positive
     */
    public static long requirePositive(long value, String fieldName) {
        if (value <= 0) {
            throw new IllegalArgumentException(fieldName + " must be positive, got: " + value);
        }
        return value;
    }

    /**
     * Validates that an integer is positive (greater than zero).
     *
     * @param value     the value to check
     * @param fieldName the name of the field (for error messages)
     * @return the positive value
     * @throws IllegalArgumentException if value is not positive
     */
    public static int requirePositive(int value, String fieldName) {
        if (value <= 0) {
            throw new IllegalArgumentException(fieldName + " must be positive, got: " + value);
        }
        return value;
    }

    /**
     * Validates that a number is non-negative (zero or greater).
     *
     * @param value     the value to check
     * @param fieldName the name of the field (for error messages)
     * @return the non-negative value
     * @throws IllegalArgumentException if value is negative
     */
    public static long requireNonNegative(long value, String fieldName) {
        if (value < 0) {
            throw new IllegalArgumentException(fieldName + " cannot be negative, got: " + value);
        }
        return value;
    }

    /**
     * Validates that an integer is non-negative (zero or greater).
     *
     * @param value     the value to check
     * @param fieldName the name of the field (for error messages)
     * @return the non-negative value
     * @throws IllegalArgumentException if value is negative
     */
    public static int requireNonNegative(int value, String fieldName) {
        if (value < 0) {
            throw new IllegalArgumentException(fieldName + " cannot be negative, got: " + value);
        }
        return value;
    }

    /**
     * Validates that a number is within a range (inclusive).
     *
     * @param value     the value to check
     * @param min       the minimum allowed value (inclusive)
     * @param max       the maximum allowed value (inclusive)
     * @param fieldName the name of the field (for error messages)
     * @return the value if within range
     * @throws IllegalArgumentException if value is outside the range
     */
    public static long requireInRange(long value, long min, long max, String fieldName) {
        if (value < min || value > max) {
            throw new IllegalArgumentException(
                    fieldName + " must be between " + min + " and " + max + ", got: " + value);
        }
        return value;
    }

    /**
     * Validates that an integer is within a range (inclusive).
     *
     * @param value     the value to check
     * @param min       the minimum allowed value (inclusive)
     * @param max       the maximum allowed value (inclusive)
     * @param fieldName the name of the field (for error messages)
     * @return the value if within range
     * @throws IllegalArgumentException if value is outside the range
     */
    public static int requireInRange(int value, int min, int max, String fieldName) {
        if (value < min || value > max) {
            throw new IllegalArgumentException(
                    fieldName + " must be between " + min + " and " + max + ", got: " + value);
        }
        return value;
    }

    /**
     * Validates that a string does not exceed a maximum length.
     *
     * @param value     the string to check
     * @param maxLength the maximum allowed length
     * @param fieldName the name of the field (for error messages)
     * @return the string if within length
     * @throws IllegalArgumentException if the string exceeds maxLength
     */
    public static String requireMaxLength(String value, int maxLength, String fieldName) {
        if (value != null && value.length() > maxLength) {
            throw new IllegalArgumentException(
                    fieldName + " cannot exceed " + maxLength + " characters, got: " + value.length());
        }
        return value;
    }

    /**
     * Validates that a string has a minimum length.
     *
     * @param value     the string to check
     * @param minLength the minimum required length
     * @param fieldName the name of the field (for error messages)
     * @return the string if sufficient length
     * @throws IllegalArgumentException if the string is shorter than minLength
     */
    public static String requireMinLength(String value, int minLength, String fieldName) {
        if (value == null || value.length() < minLength) {
            int actualLength = value == null ? 0 : value.length();
            throw new IllegalArgumentException(
                    fieldName + " must be at least " + minLength + " characters, got: " + actualLength);
        }
        return value;
    }

    /**
     * Validates that a condition is true.
     *
     * @param condition the condition to check
     * @param message   the error message if condition is false
     * @throws IllegalArgumentException if condition is false
     */
    public static void requireTrue(boolean condition, String message) {
        if (!condition) {
            throw new IllegalArgumentException(message);
        }
    }
}
