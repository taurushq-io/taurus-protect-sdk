package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Every non-2xx response in the SDK is converted here, on a body the server fully controls.
 * The pre-existing {@code ApiExceptionMapperTest} covers well-formed and malformed JSON; this
 * one covers the body being HOSTILE rather than merely wrong, which is what the mapper was
 * found not to survive.
 *
 * <p>Two defects, both reachable from any endpoint:
 *
 * <ol>
 *   <li><b>The Gson binding target was {@code ApiException}, which extends
 *   {@code java.lang.Exception}.</b> Gson's reflective adapter therefore bound the inherited
 *   {@code Throwable} fields — {@code cause} ({@code Throwable}),
 *   {@code suppressedExceptions} ({@code List<Throwable>}), {@code stackTrace} — plus the
 *   declared {@code originalException}. All of those have public no-arg constructors, so a
 *   body of nested {@code {"cause":{"cause":…}}} recursed once per level with no depth bound
 *   and terminated in a {@code StackOverflowError}: a {@code java.lang.Error}, not caught by
 *   the {@code JsonSyntaxException} handler, propagating past every
 *   {@code catch (ApiException)} and {@code catch (Exception)} an integrator wrote.</li>
 *   <li><b>The body could override the HTTP status.</b> A body {@code {"code":200}} on a 503
 *   yielded {@code getCode() == 200} and {@code isRetryable() == false} — the same
 *   error-taxonomy inversion as reporting an integrity failure as a retryable 500, reached
 *   from the other side.</li>
 * </ol>
 *
 * <p>The fix is structural rather than a bigger catch: the target is a plain private DTO whose
 * every field is a scalar, which removes the recursion sinks outright, and the transport status
 * is the only source of {@code code}.
 */
class ApiExceptionMapperHostileBodyTest {

    private ApiExceptionMapper mapper;

    @BeforeEach
    void setUp() {
        mapper = new ApiExceptionMapper();
    }

    /** An openapi exception carrying a server-chosen status and body. */
    private static com.taurushq.sdk.protect.openapi.ApiException wire(final int code,
            final String body) {
        return new com.taurushq.sdk.protect.openapi.ApiException(code, null, body);
    }

    // ---------------------------------------------------------------- recursion

    @Test
    @DisplayName("a deeply nested `cause` chain does not raise StackOverflowError")
    void deeplyNestedCauseChainIsSurvived() {
        // ~4000 levels in a few tens of kilobytes. Against the old Throwable-typed target
        // this recursed per level and died; against a scalar-only DTO the nesting is simply
        // unknown fields Gson skips.
        StringBuilder body = new StringBuilder();
        int depth = 4000;
        for (int i = 0; i < depth; i++) {
            body.append("{\"cause\":");
        }
        body.append("null");
        for (int i = 0; i < depth; i++) {
            body.append('}');
        }

        ApiException result = mapper.toApiException(wire(500, body.toString()));

        assertNotNull(result, "the mapper must return a typed exception, not throw an Error");
        assertEquals(500, result.getCode());
    }

    @Test
    @DisplayName("a deeply nested `originalException` chain is survived too")
    void deeplyNestedOriginalExceptionChainIsSurvived() {
        // The declared field, rather than an inherited one — it was the other recursion sink.
        StringBuilder body = new StringBuilder();
        int depth = 4000;
        for (int i = 0; i < depth; i++) {
            body.append("{\"originalException\":");
        }
        body.append("null");
        for (int i = 0; i < depth; i++) {
            body.append('}');
        }

        ApiException result = mapper.toApiException(wire(503, body.toString()));

        assertNotNull(result);
        assertEquals(503, result.getCode());
    }

    @Test
    @DisplayName("a deeply nested `suppressedExceptions` array is survived too")
    void deeplyNestedSuppressedExceptionsIsSurvived() {
        StringBuilder body = new StringBuilder();
        int depth = 3000;
        for (int i = 0; i < depth; i++) {
            body.append("{\"suppressedExceptions\":[");
        }
        for (int i = 0; i < depth; i++) {
            body.append("]}");
        }

        ApiException result = mapper.toApiException(wire(500, body.toString()));

        assertNotNull(result);
        assertEquals(500, result.getCode());
    }

    @Test
    @DisplayName("an over-ceiling body is not parsed, and the message it falls back to is bounded")
    void anOverCeilingBodyIsNotParsedAndTheFallbackMessageIsBounded() {
        // The size ceiling is the cheap first line: an unbounded body is not worth handing
        // to a parser however careful the parser is.
        //
        // But the ceiling alone was not enough, and this test is what found it. The
        // GENERATED exception's getMessage() interpolates the entire response body
        // ("HTTP response body: %s"), and the fallback copied that verbatim — so the one
        // case the ceiling exists for was the one case it did not cover: the body was not
        // parsed and then landed in the SDK exception's message in full, for an integrator
        // to log. A body excerpt is diagnostically useful, so the fix is a bound, not
        // removal; what must not happen is the bound being absent.
        StringBuilder body = new StringBuilder("{\"message\":\"");
        for (int i = 0; i < ApiExceptionMapper.MAX_ERROR_BODY_BYTES + 16; i++) {
            body.append('x');
        }
        body.append("\"}");

        ApiException result = mapper.toApiException(wire(500, body.toString()));

        assertNotNull(result);
        assertEquals(500, result.getCode());
        assertNotNull(result.getMessage());
        assertTrue(result.getMessage().length() <= ApiExceptionMapper.MAX_MESSAGE_CHARS + 32,
                "the fallback message must be bounded, was " + result.getMessage().length()
                        + " chars for a body of " + body.length());
        assertTrue(result.getMessage().contains("truncated"),
                "a cut message must say so, or a reader cannot tell it from a short one");
    }

    @Test
    @DisplayName("a short fallback message is passed through untouched")
    void aShortFallbackMessageIsNotTruncated() {
        // The bound must not cost the ordinary unparseable-body case its diagnostics.
        ApiException result = mapper.toApiException(wire(500, "not json at all"));

        assertNotNull(result.getMessage());
        assertFalse(result.getMessage().contains("truncated"));
        assertTrue(result.getMessage().contains("not json at all"),
                "a short body is still worth showing");
    }

    // ------------------------------------------------------- taxonomy inversion

    @Test
    @DisplayName("the body cannot lower a 503 to a non-retryable 200")
    void bodyCannotOverrideAServerErrorStatus() {
        ApiException result = mapper.toApiException(
                wire(503, "{\"code\":200,\"message\":\"all good\",\"error\":\"OK\"}"));

        assertEquals(503, result.getCode(), "the transport status is the only source of code");
        assertTrue(result.isRetryable(),
                "a 503 must stay retryable however the body labels itself");
    }

    @Test
    @DisplayName("the body cannot raise a 404 to a retryable 500")
    void bodyCannotOverrideAClientErrorStatus() {
        // The inversion in the other direction, and the worse one: a caller told to retry a
        // request that can never succeed retries it forever.
        ApiException result = mapper.toApiException(
                wire(404, "{\"code\":500,\"message\":\"nope\",\"error\":\"NotFound\"}"));

        assertEquals(404, result.getCode());
        assertFalse(result.isRetryable(), "a 404 must not become retryable");
    }

    @Test
    @DisplayName("a status with no typed branch still takes its code from the transport")
    void anUnmappedStatusTakesTheTransportCode() {
        // 409 falls to the `default` branch, which used to return the Gson-built object
        // itself — the one path by which the body's `code` reached getCode().
        ApiException result = mapper.toApiException(
                wire(409, "{\"code\":200,\"message\":\"conflict\",\"error\":\"Conflict\"}"));

        assertEquals(409, result.getCode());
        assertFalse(result.isRetryable());
    }

    // -------------------------------------------------- no server-supplied state

    @Test
    @DisplayName("a body-supplied `cause` never becomes the exception's cause")
    void aBodySuppliedCauseIsNotAdopted() {
        // Not an integrity impact on its own — the server already controls the message text
        // printed beside it — but a server-authored cause and stack trace in an integrator's
        // logs is state nobody should be able to inject, and it is the same reflection that
        // carried the recursion.
        ApiException result = mapper.toApiException(wire(500,
                "{\"message\":\"real\",\"cause\":{\"detailMessage\":\"injected\"},"
                        + "\"stackTrace\":[],\"suppressedExceptions\":[]}"));

        assertNull(result.getCause(), "the body must not supply a cause");
        assertEquals("real", result.getMessage());
    }

    @Test
    @DisplayName("the original transport exception is still attached")
    void theTransportExceptionIsStillAttached() {
        // What SHOULD be there: the real exception, for debugging. Asserted so the fix above
        // cannot be "achieved" by dropping the attachment entirely.
        com.taurushq.sdk.protect.openapi.ApiException original =
                wire(500, "{\"message\":\"boom\"}");

        ApiException result = mapper.toApiException(original);

        assertEquals(original, result.getOriginalException());
    }

    @Test
    @DisplayName("a well-formed body still maps its scalar fields")
    void aWellFormedBodyStillMaps() {
        // The hostile-input handling must not have cost the ordinary case.
        ApiException result = mapper.toApiException(wire(404,
                "{\"error\":\"NotFound\",\"message\":\"wallet 7 not found\","
                        + "\"errorCode\":\"ERR_404\"}"));

        assertEquals(404, result.getCode());
        assertEquals("wallet 7 not found", result.getMessage());
        assertEquals("NotFound", result.getError());
        assertEquals("ERR_404", result.getErrorCode());
    }
}
