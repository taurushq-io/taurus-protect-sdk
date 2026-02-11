package com.taurushq.sdk.protect.client.validation;

import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.client.model.RequestMetadataAmount;
import com.taurushq.sdk.protect.client.model.RequestStatus;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class RequestValidatorTest {

    private static final String PAYLOAD = "[{\"key\":\"request_id\",\"type\":\"String\",\"value\":\"6286376\",\"column\":\"\"},{\"key\":\"rules_key\",\"type\":\"String\",\"value\":\"XLM\",\"column\":\"\"},{\"key\":\"currency\",\"type\":\"String\",\"value\":\"XLM\",\"column\":\"\"},{\"key\":\"currency_name\",\"type\":\"String\",\"value\":\"Stellar - XLM\",\"column\":\"\"},{\"key\":\"currency_id\",\"type\":\"String\",\"value\":\"1a31194be3de1a85216300c2c0640af1d9a87a2f600a2e8927b91f78740e10ff\",\"column\":\"\"},{\"key\":\"source\",\"type\":\"Source\",\"value\":{\"type\":\"SourceInternalAddress\",\"payload\":{\"id\":\"4\",\"address\":\"GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B\",\"label\":\"nostro XLM\",\"path\":\"m/44'/148'/0'\"}},\"column\":\"source\"},{\"key\":\"destination\",\"type\":\"Destination\",\"value\":{\"type\":\"DestinationInternalAddress\",\"payload\":{\"id\":\"6\",\"address\":\"GDWIH4JIZEDCHIRNQUQAOPPBGVUS23KCVH6JH7NM2D55LIC66QVY4MK6\",\"label\":\"stellar 1\",\"path\":\"m/44'/148'/1'\"}},\"column\":\"destination\"},{\"key\":\"amount\",\"type\":\"Amount\",\"value\":{\"valueFrom\":\"2\",\"valueTo\":\"0.0000\",\"rate\":\"0.08386473412745667\",\"decimals\":\"7\",\"currencyFrom\":\"XLM\",\"currencyTo\":\"CHF\"},\"column\":\"amount\"},{\"key\":\"total_fiat_amount\",\"type\":\"String\",\"value\":\"0.0000\",\"column\":\"\"},{\"key\":\"fee_limits\",\"type\":\"BigIntArray\",\"value\":[\"1500000\"],\"column\":\"\"},{\"key\":\"fee_paid_by_receiver\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"use_unconfirmed_funds\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"is_cancel_request\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"transaction_ids\",\"type\":\"StringArray\",\"value\":[\"4:7671e0c5-623d-482a-9abd-ed90a262f6f6\"],\"column\":\"\"}]";

    private Request request;

    @BeforeEach
    void setUp() {
        request = new Request();
        request.setId(6286376);
        request.setCurrency("XLM");
        request.setStatus(RequestStatus.PENDING);
        request.setType("transfer");

        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(PAYLOAD);
        request.setMetadata(metadata);
    }

    @Test
    void validRequestPassesAllValidations() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectRequestId(6286376)
                        .expectCurrency("XLM")
                        .expectRulesKey("XLM")
                        .expectSourceAddress("GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B")
                        .expectDestinationAddress("GDWIH4JIZEDCHIRNQUQAOPPBGVUS23KCVH6JH7NM2D55LIC66QVY4MK6")
                        .expectAmount(amt -> amt
                                .valueFrom(2)
                                .decimals(7)
                                .currencyFrom("XLM")
                                .currencyTo("CHF")
                                .rate(0.08386473412745667))
                        .validate()
        );
    }

    @Test
    void isValidReturnsTrueForValidRequest() {
        boolean valid = RequestValidator.of(request)
                .expectRequestId(6286376)
                .expectCurrency("XLM")
                .isValid();

        assertTrue(valid);
    }

    @Test
    void isValidReturnsFalseForInvalidRequest() {
        boolean valid = RequestValidator.of(request)
                .expectRequestId(9999999)
                .isValid();

        assertFalse(valid);
    }

    @Test
    void wrongRequestIdThrowsValidationException() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectRequestId(9999999)
                        .validate()
        );

        assertEquals(1, ex.getErrorCount());
        assertTrue(ex.getMessage().contains("requestId"));
    }

    @Test
    void wrongCurrencyThrowsValidationException() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectCurrency("BTC")
                        .validate()
        );

        assertEquals(1, ex.getErrorCount());
        assertTrue(ex.getMessage().contains("currency"));
    }

    @Test
    void wrongSourceAddressThrowsValidationException() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectSourceAddress("WRONG_ADDRESS")
                        .validate()
        );

        assertEquals(1, ex.getErrorCount());
        assertTrue(ex.getMessage().contains("sourceAddress"));
    }

    @Test
    void wrongDestinationAddressThrowsValidationException() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectDestinationAddress("WRONG_ADDRESS")
                        .validate()
        );

        assertEquals(1, ex.getErrorCount());
        assertTrue(ex.getMessage().contains("destinationAddress"));
    }

    @Test
    void multipleErrorsAreCollected() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectRequestId(9999999)
                        .expectCurrency("BTC")
                        .expectSourceAddress("WRONG")
                        .expectDestinationAddress("ALSO_WRONG")
                        .validate()
        );

        assertEquals(4, ex.getErrorCount());
        ValidationResult result = ex.getValidationResult();
        assertTrue(result.hasErrors());
        assertEquals(4, result.getErrors().size());
    }

    @Test
    void checkReturnsResultWithoutThrowing() {
        ValidationResult result = RequestValidator.of(request)
                .expectRequestId(9999999)
                .expectCurrency("BTC")
                .check();

        assertTrue(result.hasErrors());
        assertEquals(2, result.getErrorCount());
        assertFalse(result.isValid());
    }

    @Test
    void amountValidationWorksCorrectly() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt
                                .valueFrom(2)
                                .decimals(7)
                                .currencyFrom("XLM")
                                .currencyTo("CHF")
                                .rate(0.08386473412745667))
                        .validate()
        );
    }

    @Test
    void wrongAmountValueFromThrowsValidationException() {
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueFrom(999))
                        .validate()
        );

        assertEquals(1, ex.getErrorCount());
        assertTrue(ex.getMessage().contains("amount.valueFrom"));
    }

    @Test
    void amountValueFromBetweenValidation() {
        // Should pass - 2 is within [1, 10]
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueFromBetween(1, 10))
                        .validate()
        );

        // Should fail - 2 is not within [10, 100]
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueFromBetween(10, 100))
                        .validate()
        );
        assertTrue(ex.getMessage().contains("out of range"));
    }

    @Test
    void amountRateBetweenValidation() {
        // Should pass - rate is approximately 0.084
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.rateBetween(0.01, 0.1))
                        .validate()
        );

        // Should fail - rate 0.084 is not within [0.5, 1.0]
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.rateBetween(0.5, 1.0))
                        .validate()
        );
        assertTrue(ex.getMessage().contains("rate out of range"));
    }

    @Test
    void customValidationWithRequire() {
        // Pass custom validation - use request.getId() to avoid checked exception
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .require("custom-check", md -> request.getId() > 0)
                        .validate()
        );

        // Fail custom validation
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .require("custom-check", md -> request.getId() < 0)
                        .validate()
        );
        assertTrue(ex.getMessage().contains("custom-check"));
    }

    @Test
    void expectStatusValidation() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectStatus("PENDING")
                        .validate()
        );

        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectStatus("APPROVED")
                        .validate()
        );
        assertTrue(ex.getMessage().contains("status"));
    }

    @Test
    void expectTypeValidation() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectType("transfer")
                        .validate()
        );

        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectType("withdrawal")
                        .validate()
        );
        assertTrue(ex.getMessage().contains("type"));
    }

    @Test
    void validationResultSummaryContainsAllDetails() {
        ValidationResult result = RequestValidator.of(request)
                .expectRequestId(9999999)
                .expectCurrency("BTC")
                .check();

        String summary = result.getSummary();
        assertTrue(summary.contains("Validation failed"));
        assertTrue(summary.contains("2 error(s)"));
        assertTrue(summary.contains("requestId"));
        assertTrue(summary.contains("currency"));
        assertTrue(summary.contains("expected:"));
        assertTrue(summary.contains("actual:"));
    }

    @Test
    void validationErrorContainsExpectedAndActualValues() {
        ValidationResult result = RequestValidator.of(request)
                .expectRequestId(9999999)
                .check();

        ValidationResult.ValidationError error = result.getErrors().get(0);
        assertEquals("requestId", error.getField());
        assertEquals("value mismatch", error.getMessage());
        assertEquals(9999999L, error.getExpected());
        assertEquals(6286376L, error.getActual());
    }

    @Test
    void getRequestReturnsOriginalRequest() {
        RequestValidator validator = RequestValidator.of(request);
        assertEquals(request, validator.getRequest());
    }

    @Test
    void nullRequestThrowsNullPointerException() {
        assertThrows(NullPointerException.class, () ->
                RequestValidator.of(null)
        );
    }

    @Test
    void nullMetadataThrowsNullPointerException() {
        Request requestWithoutMetadata = new Request();
        assertThrows(NullPointerException.class, () ->
                RequestValidator.of(requestWithoutMetadata)
        );
    }

    @Test
    void expectRulesKeyValidation() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectRulesKey("XLM")
                        .validate()
        );

        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectRulesKey("WRONG")
                        .validate()
        );
        assertTrue(ex.getMessage().contains("rulesKey"));
    }

    @Test
    void amountValueToValidation() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueTo(0.0))
                        .validate()
        );

        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueTo(100.0))
                        .validate()
        );
        assertTrue(ex.getMessage().contains("amount.valueTo"));
    }

    @Test
    void amountValueToBetweenValidation() {
        // valueTo is 0.0, should pass
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueToBetween(-1.0, 1.0))
                        .validate()
        );

        // valueTo is 0.0, should fail
        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.valueToBetween(10.0, 100.0))
                        .validate()
        );
        assertTrue(ex.getMessage().contains("amount.valueTo"));
    }

    @Test
    void amountDecimalsValidation() {
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.decimals(7))
                        .validate()
        );

        RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                RequestValidator.of(request)
                        .expectAmount(amt -> amt.decimals(8))
                        .validate()
        );
        assertTrue(ex.getMessage().contains("amount.decimals"));
    }

    @Test
    void complexValidationScenario() {
        // Validate everything in one fluent chain
        assertDoesNotThrow(() ->
                RequestValidator.of(request)
                        .expectRequestId(6286376)
                        .expectCurrency("XLM")
                        .expectRulesKey("XLM")
                        .expectSourceAddress("GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B")
                        .expectDestinationAddress("GDWIH4JIZEDCHIRNQUQAOPPBGVUS23KCVH6JH7NM2D55LIC66QVY4MK6")
                        .expectStatus("PENDING")
                        .expectType("transfer")
                        .expectAmount(amt -> amt
                                .valueFrom(2)
                                .valueFromBetween(1, 100)
                                .valueTo(0.0)
                                .valueToBetween(-10.0, 10.0)
                                .decimals(7)
                                .currencyFrom("XLM")
                                .currencyTo("CHF")
                                .rate(0.08386473412745667)
                                .rateBetween(0.01, 1.0))
                        .require("positive-id", md -> request.getId() > 0)
                        .validate()
        );
    }

    // ==================== ValidationResult Edge Cases ====================

    @Nested
    @DisplayName("ValidationResult Edge Cases")
    class ValidationResultEdgeCases {

        @Test
        void getSummaryReturnsPassedWhenNoErrors() {
            ValidationResult result = new ValidationResult();
            assertEquals("Validation passed", result.getSummary());
            assertTrue(result.isValid());
            assertFalse(result.hasErrors());
            assertEquals(0, result.getErrorCount());
        }

        @Test
        void getSummaryExcludesExpectedActualWhenBothNull() {
            ValidationResult result = new ValidationResult();
            result.addError("field", "custom error message");

            String summary = result.getSummary();
            assertTrue(summary.contains("field"));
            assertTrue(summary.contains("custom error message"));
            assertFalse(summary.contains("expected:"));
            assertFalse(summary.contains("actual:"));
        }

        @Test
        void getSummaryIncludesExpectedActualWhenOnlyOneIsNull() {
            ValidationResult result = new ValidationResult();
            result.addError("field", "message", "expected", null);

            String summary = result.getSummary();
            assertTrue(summary.contains("expected:"));
            assertTrue(summary.contains("actual:"));
        }

        @Test
        void validationErrorToStringFormat() {
            ValidationResult.ValidationError error = new ValidationResult.ValidationError(
                    "testField", "testMessage", "expectedVal", "actualVal");

            String str = error.toString();
            assertTrue(str.contains("testField"));
            assertTrue(str.contains("testMessage"));
            assertTrue(str.contains("expectedVal"));
            assertTrue(str.contains("actualVal"));
            assertTrue(str.contains("ValidationError{"));
        }

        @Test
        void getErrorsReturnsUnmodifiableList() {
            ValidationResult result = new ValidationResult();
            result.addError("field", "message");

            assertThrows(UnsupportedOperationException.class, () ->
                    result.getErrors().add(new ValidationResult.ValidationError("x", "y", null, null))
            );
        }

        @Test
        void addMultipleErrorsIncreasesCount() {
            ValidationResult result = new ValidationResult();
            assertEquals(0, result.getErrorCount());

            result.addError("field1", "msg1");
            assertEquals(1, result.getErrorCount());

            result.addError("field2", "msg2", "exp", "act");
            assertEquals(2, result.getErrorCount());

            result.addError("field3", "msg3");
            assertEquals(3, result.getErrorCount());
        }

        @Test
        void validationErrorWithNullValues() {
            ValidationResult.ValidationError error = new ValidationResult.ValidationError(
                    null, null, null, null);

            assertEquals(null, error.getField());
            assertEquals(null, error.getMessage());
            assertEquals(null, error.getExpected());
            assertEquals(null, error.getActual());
        }
    }

    // ==================== AmountValidator Boundary Tests ====================

    @Nested
    @DisplayName("AmountValidator Boundary Tests")
    class AmountValidatorBoundaryTests {

        @Test
        void valueFromAtExactMinBoundary() {
            // valueFrom is 2, check boundary at 2
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFromBetween(2, 10))
                            .validate()
            );
        }

        @Test
        void valueFromAtExactMaxBoundary() {
            // valueFrom is 2, check boundary at 2
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFromBetween(0, 2))
                            .validate()
            );
        }

        @Test
        void valueFromJustOutsideMinBoundary() {
            // valueFrom is 2, check just outside
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFromBetween(3, 10))
                            .validate()
            );
            assertTrue(ex.getMessage().contains("out of range"));
        }

        @Test
        void valueFromJustOutsideMaxBoundary() {
            // valueFrom is 2, check just outside
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFromBetween(0, 1))
                            .validate()
            );
            assertTrue(ex.getMessage().contains("out of range"));
        }

        @Test
        void valueFromWithZeroExpected() {
            // valueFrom is 2, expecting 0 should fail
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFrom(0))
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void valueFromWithNegativeExpected() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFrom(-100))
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void valueFromBetweenWithNegativeRange() {
            // valueFrom is 2, negative range should fail
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFromBetween(-100, -1))
                            .validate()
            );
            assertTrue(ex.getMessage().contains("out of range"));
        }

        @Test
        void valueFromWithMaxLongExpected() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueFrom(Long.MAX_VALUE))
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void valueToBetweenWithNegativeRange() {
            // valueTo is 0.0, should pass for range containing 0
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.valueToBetween(-100.0, 100.0))
                            .validate()
            );
        }

        @Test
        void rateBetweenAtExactMinBoundary() {
            // rate is ~0.084, test exact min
            double rate = 0.08386473412745667;
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.rateBetween(rate, 1.0))
                            .validate()
            );
        }

        @Test
        void rateBetweenAtExactMaxBoundary() {
            double rate = 0.08386473412745667;
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.rateBetween(0.0, rate))
                            .validate()
            );
        }

        @Test
        void rateBetweenWithZeroRange() {
            // rate is ~0.084, zero range should fail
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.rateBetween(0.0, 0.0))
                            .validate()
            );
            assertTrue(ex.getMessage().contains("rate out of range"));
        }

        @Test
        void currencyFromWithNullExpected() {
            // currencyFrom is "XLM", expecting null should fail
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.currencyFrom(null))
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void currencyToWithNullExpected() {
            // currencyTo is "CHF", expecting null should fail
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt.currencyTo(null))
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void multipleAmountErrorsCollected() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> amt
                                    .valueFrom(999)
                                    .valueTo(999.0)
                                    .decimals(99)
                                    .currencyFrom("WRONG")
                                    .currencyTo("WRONG")
                                    .rate(999.0))
                            .validate()
            );
            assertEquals(6, ex.getErrorCount());
        }

        @Test
        void amountValidatorDirectInstantiation() {
            RequestMetadataAmount amount = new RequestMetadataAmount();
            amount.setValueFrom("100");
            amount.setValueTo("50.5");
            amount.setDecimals(8);
            amount.setCurrencyFrom("BTC");
            amount.setCurrencyTo("USD");
            amount.setRate("45000.0");

            ValidationResult result = new ValidationResult();
            AmountValidator validator = new AmountValidator(amount, result);

            validator.valueFrom(100)
                    .valueTo(50.5)
                    .decimals(8)
                    .currencyFrom("BTC")
                    .currencyTo("USD")
                    .rate(45000.0);

            assertTrue(result.isValid());
        }

        @Test
        void amountValidatorChainingReturnsThis() {
            RequestMetadataAmount amount = new RequestMetadataAmount();
            ValidationResult result = new ValidationResult();
            AmountValidator validator = new AmountValidator(amount, result);

            AmountValidator returned = validator.valueFrom(0);
            assertEquals(validator, returned);

            returned = validator.valueFromBetween(0, 100);
            assertEquals(validator, returned);

            returned = validator.valueTo(0.0);
            assertEquals(validator, returned);

            returned = validator.valueToBetween(0.0, 100.0);
            assertEquals(validator, returned);

            returned = validator.decimals(0);
            assertEquals(validator, returned);

            returned = validator.currencyFrom(null);
            assertEquals(validator, returned);

            returned = validator.currencyTo(null);
            assertEquals(validator, returned);

            returned = validator.rate(0.0);
            assertEquals(validator, returned);

            returned = validator.rateBetween(0.0, 100.0);
            assertEquals(validator, returned);
        }
    }

    // ==================== RequestValidator Edge Cases ====================

    @Nested
    @DisplayName("RequestValidator Edge Cases")
    class RequestValidatorEdgeCases {

        @Test
        void expectCurrencyWithEmptyString() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectCurrency("")
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void requirePredicateThatThrowsRuntimeException() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .require("throwing-predicate", md -> {
                                throw new RuntimeException("intentional test exception");
                            })
                            .validate()
            );
            assertTrue(ex.getMessage().contains("throwing-predicate"));
            assertTrue(ex.getMessage().contains("threw exception"));
        }

        @Test
        void expectStatusWhenStatusIsNull() {
            Request requestWithNullStatus = new Request();
            requestWithNullStatus.setMetadata(request.getMetadata());

            // Expecting null matches null
            assertDoesNotThrow(() ->
                    RequestValidator.of(requestWithNullStatus)
                            .expectStatus(null)
                            .validate()
            );

            // Expecting non-null fails
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(requestWithNullStatus)
                            .expectStatus("PENDING")
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void expectTypeWhenTypeIsNull() {
            Request requestWithNullType = new Request();
            requestWithNullType.setMetadata(request.getMetadata());

            // Expecting null matches null
            assertDoesNotThrow(() ->
                    RequestValidator.of(requestWithNullType)
                            .expectType(null)
                            .validate()
            );

            // Expecting non-null fails
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(requestWithNullType)
                            .expectType("transfer")
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void validateWithNoValidationsAdded() {
            // No validations added should pass
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .validate()
            );

            ValidationResult result = RequestValidator.of(request).check();
            assertTrue(result.isValid());
        }

        @Test
        void validatorChainingReturnsThis() {
            RequestValidator validator = RequestValidator.of(request);

            assertEquals(validator, validator.expectRequestId(0));
            assertEquals(validator, validator.expectCurrency(""));
            assertEquals(validator, validator.expectRulesKey(""));
            assertEquals(validator, validator.expectSourceAddress(""));
            assertEquals(validator, validator.expectDestinationAddress(""));
            assertEquals(validator, validator.expectStatus(""));
            assertEquals(validator, validator.expectType(""));
            assertEquals(validator, validator.expectAmount(amt -> {}));
            assertEquals(validator, validator.require("test", md -> true));
        }

        @Test
        void requireWithNullNameThrows() {
            assertThrows(NullPointerException.class, () ->
                    RequestValidator.of(request)
                            .require(null, md -> true)
            );
        }

        @Test
        void requireWithNullPredicateThrows() {
            assertThrows(NullPointerException.class, () ->
                    RequestValidator.of(request)
                            .require("test", null)
            );
        }

        @Test
        void expectSourceAddressWithSpecialCharacters() {
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectSourceAddress("GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B")
                            .validate()
            );
        }

        @Test
        void expectDestinationAddressWithUnicode() {
            // Unicode in expected value that doesn't match
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectDestinationAddress("地址")
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }

        @Test
        void veryLongExpectedAddressString() {
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < 10000; i++) {
                sb.append("A");
            }
            String longString = sb.toString();
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectSourceAddress(longString)
                            .validate()
            );
            assertEquals(1, ex.getErrorCount());
        }
    }

    // ==================== RequestValidationException Tests ====================

    @Nested
    @DisplayName("RequestValidationException Tests")
    class RequestValidationExceptionTests {

        @Test
        void exceptionWithCustomMessage() {
            ValidationResult result = new ValidationResult();
            result.addError("field", "error");

            RequestValidationException ex = new RequestValidationException("Custom message", result);

            assertEquals("Custom message", ex.getMessage());
            assertEquals(result, ex.getValidationResult());
        }

        @Test
        void exceptionMessageContainsSummary() {
            ValidationResult result = new ValidationResult();
            result.addError("testField", "testError", "expected", "actual");

            RequestValidationException ex = new RequestValidationException(result);

            assertTrue(ex.getMessage().contains("Validation failed"));
            assertTrue(ex.getMessage().contains("testField"));
            assertTrue(ex.getMessage().contains("testError"));
        }

        @Test
        void getValidationResultReturnsOriginal() {
            ValidationResult result = new ValidationResult();
            result.addError("field", "error");

            RequestValidationException ex = new RequestValidationException(result);

            // Same object reference
            assertTrue(result == ex.getValidationResult());
        }

        @Test
        void exceptionIsRuntimeException() {
            ValidationResult result = new ValidationResult();
            RequestValidationException ex = new RequestValidationException(result);

            assertTrue(ex instanceof RuntimeException);
        }

        @Test
        void errorCountDelegatesToValidationResult() {
            ValidationResult result = new ValidationResult();
            result.addError("field1", "error1");
            result.addError("field2", "error2");
            result.addError("field3", "error3");

            RequestValidationException ex = new RequestValidationException(result);

            assertEquals(3, ex.getErrorCount());
            assertEquals(result.getErrorCount(), ex.getErrorCount());
        }
    }

    // ==================== Integration/Complex Scenarios ====================

    @Nested
    @DisplayName("Integration and Complex Scenarios")
    class IntegrationTests {

        @Test
        void allFieldsFailSimultaneously() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectRequestId(0)
                            .expectCurrency("WRONG")
                            .expectRulesKey("WRONG")
                            .expectSourceAddress("WRONG")
                            .expectDestinationAddress("WRONG")
                            .expectStatus("WRONG")
                            .expectType("WRONG")
                            .expectAmount(amt -> amt
                                    .valueFrom(0)
                                    .valueTo(999.0)
                                    .decimals(0)
                                    .currencyFrom("WRONG")
                                    .currencyTo("WRONG")
                                    .rate(0.0))
                            .require("always-fails", md -> false)
                            .validate()
            );

            // 7 request fields + 6 amount fields + 1 custom = 14 errors
            assertEquals(14, ex.getErrorCount());
        }

        @Test
        void eachNewValidatorIsIndependent() {
            // First validator
            ValidationResult result1 = RequestValidator.of(request)
                    .expectRequestId(0)
                    .check();
            assertEquals(1, result1.getErrorCount());

            // Second validator should start fresh
            ValidationResult result2 = RequestValidator.of(request)
                    .expectCurrency("WRONG")
                    .check();
            assertEquals(1, result2.getErrorCount());

            // They are independent
            assertTrue(result1 != result2);
        }

        @Test
        void validateThenCheckReturnsSameResult() {
            RequestValidator validator = RequestValidator.of(request)
                    .expectRequestId(6286376);

            // Check first
            ValidationResult result = validator.check();
            assertTrue(result.isValid());

            // Validate should not throw for same validator
            assertDoesNotThrow(() -> validator.validate());
        }

        @Test
        void multipleCallsToSameValidation() {
            // Calling the same validation multiple times accumulates errors
            RequestValidator validator = RequestValidator.of(request)
                    .expectRequestId(0)
                    .expectRequestId(1)
                    .expectRequestId(2);

            ValidationResult result = validator.check();
            assertEquals(3, result.getErrorCount());
        }

        @Test
        void mixedPassAndFailValidations() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectRequestId(6286376)      // PASS
                            .expectCurrency("XLM")         // PASS
                            .expectSourceAddress("WRONG")   // FAIL
                            .expectRulesKey("XLM")         // PASS
                            .expectDestinationAddress("X") // FAIL
                            .validate()
            );

            assertEquals(2, ex.getErrorCount());
        }

        @Test
        void emptyAmountValidation() {
            // No amount validations should pass
            assertDoesNotThrow(() ->
                    RequestValidator.of(request)
                            .expectAmount(amt -> {})
                            .validate()
            );
        }

        @Test
        void validationWithMinimalPayload() {
            // Create request with minimal JSON payload (empty array)
            Request minimalRequest = new Request();
            minimalRequest.setStatus(RequestStatus.PENDING);
            RequestMetadata metadata = new RequestMetadata();
            metadata.setPayloadAsString("[]");
            minimalRequest.setMetadata(metadata);

            // Fields should report "not found"
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(minimalRequest)
                            .expectRequestId(123)
                            .expectCurrency("BTC")
                            .expectSourceAddress("addr")
                            .expectDestinationAddress("addr2")
                            .expectAmount(amt -> amt.valueFrom(100))
                            .validate()
            );

            // All 5 should fail with "not found"
            assertEquals(5, ex.getErrorCount());
            assertTrue(ex.getMessage().contains("not found"));
        }

        @Test
        void concurrentErrorMessageFormatting() {
            RequestValidationException ex = assertThrows(RequestValidationException.class, () ->
                    RequestValidator.of(request)
                            .expectRequestId(0)
                            .validate()
            );

            String msg = ex.getMessage();
            // Message should be well-formed
            assertTrue(msg.startsWith("Validation failed"));
            assertTrue(msg.contains("error(s)"));
            assertTrue(msg.contains("requestId"));
        }

        @Test
        void isValidAndHasErrorsAreOpposites() {
            ValidationResult validResult = RequestValidator.of(request)
                    .expectRequestId(6286376)
                    .check();
            assertTrue(validResult.isValid());
            assertFalse(validResult.hasErrors());

            ValidationResult invalidResult = RequestValidator.of(request)
                    .expectRequestId(0)
                    .check();
            assertFalse(invalidResult.isValid());
            assertTrue(invalidResult.hasErrors());
        }
    }
}
