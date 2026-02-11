package com.taurushq.sdk.protect.client.validation;

import com.taurushq.sdk.protect.client.model.RequestMetadataAmount;

import java.math.BigDecimal;

/**
 * Fluent validator for {@link RequestMetadataAmount} fields.
 * <p>
 * This validator is used as a nested builder within {@link RequestValidator}
 * to validate amount-related fields such as value, decimals, currencies, and rate.
 * <p>
 * Amount fields (valueFrom, valueTo, rate) are string representations to preserve
 * arbitrary precision. Validation methods accept numeric types for convenience and
 * compare using {@link BigDecimal} for precision-safe comparisons.
 * <p>
 * Example usage:
 * <pre>{@code
 * RequestValidator.of(request)
 *     .expectAmount(amt -> amt
 *         .valueFrom(1000)
 *         .decimals(8)
 *         .currencyFrom("BTC")
 *         .currencyTo("USD")
 *         .rateBetween(30000, 50000))
 *     .validate();
 * }</pre>
 */
public class AmountValidator {

    private final RequestMetadataAmount amount;
    private final ValidationResult result;

    /**
     * Creates a new amount validator.
     *
     * @param amount the amount to validate
     * @param result the validation result to collect errors into
     */
    AmountValidator(RequestMetadataAmount amount, ValidationResult result) {
        this.amount = amount;
        this.result = result;
    }

    /**
     * Validates that valueFrom equals the expected value.
     *
     * @param expected the expected valueFrom
     * @return this validator for chaining
     */
    public AmountValidator valueFrom(long expected) {
        return valueFrom(String.valueOf(expected));
    }

    /**
     * Validates that valueFrom equals the expected string value.
     *
     * @param expected the expected valueFrom as a string
     * @return this validator for chaining
     */
    public AmountValidator valueFrom(String expected) {
        if (!compareBigDecimal(expected, amount.getValueFrom())) {
            result.addError("amount.valueFrom", "value mismatch", expected, amount.getValueFrom());
        }
        return this;
    }

    /**
     * Validates that valueFrom is within the specified range (inclusive).
     *
     * @param min minimum value (inclusive)
     * @param max maximum value (inclusive)
     * @return this validator for chaining
     */
    public AmountValidator valueFromBetween(long min, long max) {
        BigDecimal actual = toBigDecimal(amount.getValueFrom());
        if (actual == null || actual.compareTo(BigDecimal.valueOf(min)) < 0
                || actual.compareTo(BigDecimal.valueOf(max)) > 0) {
            result.addError("amount.valueFrom", "value out of range [" + min + ", " + max + "]",
                    min + "-" + max, amount.getValueFrom());
        }
        return this;
    }

    /**
     * Validates that valueTo equals the expected value.
     *
     * @param expected the expected valueTo
     * @return this validator for chaining
     */
    public AmountValidator valueTo(double expected) {
        return valueTo(BigDecimal.valueOf(expected).toPlainString());
    }

    /**
     * Validates that valueTo equals the expected string value.
     *
     * @param expected the expected valueTo as a string
     * @return this validator for chaining
     */
    public AmountValidator valueTo(String expected) {
        if (!compareBigDecimal(expected, amount.getValueTo())) {
            result.addError("amount.valueTo", "value mismatch", expected, amount.getValueTo());
        }
        return this;
    }

    /**
     * Validates that valueTo is within the specified range (inclusive).
     *
     * @param min minimum value (inclusive)
     * @param max maximum value (inclusive)
     * @return this validator for chaining
     */
    public AmountValidator valueToBetween(double min, double max) {
        BigDecimal actual = toBigDecimal(amount.getValueTo());
        if (actual == null || actual.compareTo(BigDecimal.valueOf(min)) < 0
                || actual.compareTo(BigDecimal.valueOf(max)) > 0) {
            result.addError("amount.valueTo", "value out of range [" + min + ", " + max + "]",
                    min + "-" + max, amount.getValueTo());
        }
        return this;
    }

    /**
     * Validates that decimals equals the expected value.
     *
     * @param expected the expected decimals
     * @return this validator for chaining
     */
    public AmountValidator decimals(int expected) {
        if (amount.getDecimals() != expected) {
            result.addError("amount.decimals", "value mismatch", expected, amount.getDecimals());
        }
        return this;
    }

    /**
     * Validates that currencyFrom equals the expected value.
     *
     * @param expected the expected currencyFrom
     * @return this validator for chaining
     */
    public AmountValidator currencyFrom(String expected) {
        String actual = amount.getCurrencyFrom();
        if (!equalsNullSafe(expected, actual)) {
            result.addError("amount.currencyFrom", "value mismatch", expected, actual);
        }
        return this;
    }

    /**
     * Validates that currencyTo equals the expected value.
     *
     * @param expected the expected currencyTo
     * @return this validator for chaining
     */
    public AmountValidator currencyTo(String expected) {
        String actual = amount.getCurrencyTo();
        if (!equalsNullSafe(expected, actual)) {
            result.addError("amount.currencyTo", "value mismatch", expected, actual);
        }
        return this;
    }

    /**
     * Validates that rate equals the expected value.
     *
     * @param expected the expected rate
     * @return this validator for chaining
     */
    public AmountValidator rate(double expected) {
        return rate(BigDecimal.valueOf(expected).toPlainString());
    }

    /**
     * Validates that rate equals the expected string value.
     *
     * @param expected the expected rate as a string
     * @return this validator for chaining
     */
    public AmountValidator rate(String expected) {
        if (!compareBigDecimal(expected, amount.getRate())) {
            result.addError("amount.rate", "value mismatch", expected, amount.getRate());
        }
        return this;
    }

    /**
     * Validates that rate is within the specified range (inclusive).
     *
     * @param min minimum rate (inclusive)
     * @param max maximum rate (inclusive)
     * @return this validator for chaining
     */
    public AmountValidator rateBetween(double min, double max) {
        BigDecimal actual = toBigDecimal(amount.getRate());
        if (actual == null || actual.compareTo(BigDecimal.valueOf(min)) < 0
                || actual.compareTo(BigDecimal.valueOf(max)) > 0) {
            result.addError("amount.rate", "rate out of range [" + min + ", " + max + "]",
                    min + "-" + max, amount.getRate());
        }
        return this;
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

    private static BigDecimal toBigDecimal(String value) {
        if (value == null) {
            return null;
        }
        try {
            return new BigDecimal(value);
        } catch (NumberFormatException e) {
            return null;
        }
    }

    private static boolean compareBigDecimal(String expected, String actual) {
        if (expected == null && actual == null) {
            return true;
        }
        if (expected == null || actual == null) {
            return false;
        }
        BigDecimal expectedBd = toBigDecimal(expected);
        BigDecimal actualBd = toBigDecimal(actual);
        if (expectedBd == null || actualBd == null) {
            return expected.equals(actual);
        }
        return expectedBd.compareTo(actualBd) == 0;
    }
}
