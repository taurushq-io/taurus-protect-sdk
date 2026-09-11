package com.taurushq.sdk.protect.client.mapper;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.gson.JsonParseException;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.AuthenticationException;
import com.taurushq.sdk.protect.client.model.AuthorizationException;
import com.taurushq.sdk.protect.client.model.NotFoundException;
import com.taurushq.sdk.protect.client.model.RateLimitException;
import com.taurushq.sdk.protect.client.model.ServerException;
import com.taurushq.sdk.protect.client.model.ValidationException;

import java.util.logging.Level;
import java.util.logging.Logger;


/**
 * Maps OpenAPI exceptions to typed SDK exceptions.
 * <p>
 * This mapper creates specific exception types based on HTTP status codes,
 * allowing clients to handle different error categories appropriately.
 *
 * <h2>Threat model: the API server is the adversary</h2>
 *
 * <p>The HTTP error body is fully attacker-controlled, and this class is the one place the
 * SDK hands it to a reflective JSON parser. Three properties follow from that and must not
 * be relaxed:
 *
 * <ul>
 *   <li><b>Nothing deserializes into a {@link Throwable}.</b> This mapper used to call
 *       {@code Gson.fromJson(body, ApiException.class)}, and
 *       {@code ApiException extends Exception}, so Gson reflected into
 *       {@link Throwable}'s own fields: {@code cause} (a {@code Throwable}),
 *       {@code suppressedExceptions} (a {@code List<Throwable>}), {@code stackTrace}, plus
 *       the SDK's own {@code originalException}. Every one of those is a recursion or graph
 *       sink — a body of {@code {"cause":{"cause":{"cause":...}}}} drives Gson's reflective
 *       adapter into unbounded mutual recursion, and {@code suppressedExceptions} adds
 *       branching on top. Parsing now targets {@link ErrorBody}, a flat non-{@code Throwable}
 *       DTO, so no {@code Throwable} field is reachable from the wire at all.</li>
 *   <li><b>Every parse failure is contained.</b> See {@link #parseResponseBody}.</li>
 *   <li><b>The transport status always decides the error taxonomy.</b> See
 *       {@link #toApiException}.</li>
 * </ul>
 *
 * <p>The {@code Throwable} deserialization was also a live availability defect, not only a
 * hardening concern: JDK 16+ denies deep reflection into {@code java.base/java.lang} by
 * default, so {@code Gson.fromJson(body, ApiException.class)} threw
 * {@code JsonIOException: Failed making field 'java.lang.Throwable#detailMessage'
 * accessible} while building the adapter — before even looking at the body. That escaped
 * the old {@code catch (JsonSyntaxException)} (see {@link #parseResponseBody}), so on a
 * modern JDK every API error carrying a response body surfaced to callers as a raw Gson
 * exception instead of the typed SDK exception. It was invisible in this repo only because
 * a surefire {@code add-opens} argument opened {@code java.lang} for the test JVM; the
 * shipped artifact carries no such flag.
 *
 * <p>The returned exception type is determined by the HTTP status code:
 * <ul>
 *   <li>400 → {@link ValidationException}</li>
 *   <li>401 → {@link AuthenticationException}</li>
 *   <li>403 → {@link AuthorizationException}</li>
 *   <li>404 → {@link NotFoundException}</li>
 *   <li>429 → {@link RateLimitException}</li>
 *   <li>5xx → {@link ServerException}</li>
 *   <li>Other → {@link ApiException}</li>
 * </ul>
 */
@SuppressWarnings("PMD.EmptyCatchBlock")
public class ApiExceptionMapper {

    private static final Logger LOGGER = Logger.getLogger(ApiExceptionMapper.class.getName());

    private static final long DEFAULT_RATE_LIMIT_RETRY_MS = 1000;

    /**
     * Ceiling on the error body handed to the JSON parser: 64 KiB.
     *
     * <p>An oversized body is REJECTED, never truncated — the transport status and message
     * are already known and trustworthy, whereas a truncated attacker-chosen document is
     * still an attacker-chosen document, just a shorter one.
     *
     * <p>64 KiB is roughly a hundredfold headroom: the four recognised fields are scalars,
     * and the largest legitimate one is {@code message}, which even for validatord's
     * wrapped gRPC status plus a required-role list runs to a few hundred bytes. Anything
     * above the ceiling is a hostile response trying to make the SDK do work, not an error
     * report. Measured with {@code String.length()} (UTF-16 chars), like the
     * {@code WhitelistHashHelper.MAX_PAYLOAD_BYTES} precedent it follows; the aim is to
     * bound the work handed to the parser, not to account for bytes exactly.
     */
    static final int MAX_ERROR_BODY_BYTES = 1 << 16;

    /**
     * Ceiling on a transport-derived fallback message. Generous enough to keep the
     * status line and a useful body excerpt, small enough that a hostile body cannot
     * turn one failed call into megabytes in an integrator's log.
     */
    static final int MAX_MESSAGE_CHARS = 2048;

    private static final Gson GSON = new Gson();


    /**
     * Converts an OpenAPI exception to a typed SDK exception.
     *
     * <p>The status code that selects the type, and the {@code code} carried on the result,
     * come from the TRANSPORT status and never from the response body. A body is free to
     * claim {@code {"code":200}} on a 503; honouring that would invert the error taxonomy,
     * turning a retryable server error into a permanent client error (or the reverse, which
     * is worse: a 4xx reported as retryable makes a caller retry a request that can never
     * succeed). {@link ErrorBody#code} is therefore read only to log the discrepancy.
     *
     * @param e the OpenAPI exception
     * @return a typed SDK exception
     */
    public ApiException toApiException(com.taurushq.sdk.protect.openapi.ApiException e) {
        int statusCode = e.getCode();

        // Try to parse the response body for additional error details
        ErrorBody body = parseResponseBody(e);

        // A body disagreeing with the transport status is an attempted taxonomy inversion.
        // Record the claim at FINE and discard it; the message itself is server-controlled
        // text and is deliberately not echoed here.
        if (LOGGER.isLoggable(Level.FINE) && body.code != 0 && body.code != statusCode) {
            LOGGER.fine(String.format(
                    "error response body claimed HTTP code %d on a transport status of %d;"
                            + " the transport status wins",
                    body.code, statusCode));
        }

        // Create the appropriate typed exception
        ApiException result = createTypedException(statusCode, body, e);

        // Set the original exception for debugging
        result.setOriginalException(e);

        return result;
    }

    /**
     * Caps a transport-derived fallback message.
     *
     * <p>Not a formatting nicety: the value is server-controlled and, via the generated
     * exception's {@code getMessage()}, unbounded. A truncation marker is appended so a
     * reader can tell a cut message from a short one.
     *
     * @param message the transport message, possibly null
     * @return the message, capped at {@link #MAX_MESSAGE_CHARS}
     */
    private static String truncateForMessage(final String message) {
        if (message == null || message.length() <= MAX_MESSAGE_CHARS) {
            return message;
        }
        return message.substring(0, MAX_MESSAGE_CHARS) + "… [truncated]";
    }

    /**
     * Parses the response body into the plain DTO, falling back to transport-derived
     * defaults whenever it cannot be trusted.
     *
     * <p>The catch is {@link JsonParseException} rather than {@code JsonSyntaxException}
     * because {@code JsonIOException} is its SIBLING, not its subclass — both extend
     * {@code JsonParseException}. Catching only the syntax half let the I/O half escape,
     * which is exactly what happened with the reflective-access failure described on the
     * class javadoc. The shared parent covers both, and any future sibling.
     *
     * <p>{@link StackOverflowError} is caught as a backstop. With {@link ErrorBody} the
     * reflective descent is bounded to depth one BY CONSTRUCTION — all four fields are
     * scalars, so a nested value under a recognised key is a type mismatch rather than a
     * descent, and an unrecognised key is skipped by {@code JsonReader.skipValue()}, which
     * is iterative. That reasoning depends on Gson internals, so the catch keeps the
     * guarantee true across a Gson upgrade that made either of those recursive.
     *
     * <p>Deliberately NOT caught: {@code OutOfMemoryError} (bounded by
     * {@link #MAX_ERROR_BODY_BYTES}, and not recoverable) and any unchecked exception from
     * a genuine defect in this SDK — reporting a {@code NullPointerException} as "the body
     * was unparseable" would hide it.
     */
    private ErrorBody parseResponseBody(com.taurushq.sdk.protect.openapi.ApiException e) {
        ErrorBody defaults = new ErrorBody();
        defaults.error = "Unknown";
        // BOUNDED, because the generated exception's getMessage() interpolates the ENTIRE
        // response body ("HTTP response body: %s"). Without this, the MAX_ERROR_BODY_BYTES
        // ceiling below stops an unbounded body being PARSED and then copies it into the
        // SDK exception's message anyway — so the one case the ceiling exists for is the one
        // case it did not cover, and an integrator logging the exception logs the whole
        // thing. The status-line prefix carries the diagnostic value; the body tail does not.
        defaults.message = truncateForMessage(e.getMessage());

        String responseBody = e.getResponseBody();
        if (Strings.isNullOrEmpty(responseBody)) {
            return defaults;
        }

        if (responseBody.length() > MAX_ERROR_BODY_BYTES) {
            if (LOGGER.isLoggable(Level.FINE)) {
                LOGGER.fine(String.format(
                        "error response body of %d chars exceeds the %d-char ceiling and was"
                                + " not parsed; using transport-derived defaults",
                        responseBody.length(), MAX_ERROR_BODY_BYTES));
            }
            return defaults;
        }

        try {
            ErrorBody fromJson = GSON.fromJson(responseBody, ErrorBody.class);
            // A body of "null" parses successfully to null; keep the defaults for it.
            return fromJson == null ? defaults : fromJson;
        } catch (JsonParseException | StackOverflowError ex) {
            if (LOGGER.isLoggable(Level.FINE)) {
                // The parser's message can embed the server-controlled body; do not echo it.
                LOGGER.fine("error response body is not a parseable JSON error object;"
                        + " using transport-derived defaults");
            }
            return defaults;
        }
    }

    /**
     * Creates a typed exception based on the HTTP status code.
     *
     * <p>Every branch CONSTRUCTS a fresh typed exception from the DTO's fields. None of them
     * returns the parsed object itself: the {@code default} branch used to do exactly that,
     * which handed the caller a server-deserialized object and was the one path by which the
     * body's {@code code} reached {@link ApiException#getCode()} — and through it
     * {@link ApiException#isRetryable()}.
     */
    private ApiException createTypedException(int statusCode, ErrorBody body,
            com.taurushq.sdk.protect.openapi.ApiException original) {

        ApiException result;

        switch (statusCode) {
            case 400:
                result = new ValidationException(body.message, body.error, body.errorCode);
                break;
            case 401:
                result = new AuthenticationException(body.message, body.error, body.errorCode);
                break;
            case 403:
                result = new AuthorizationException(body.message, body.error, body.errorCode);
                break;
            case 404:
                result = new NotFoundException(body.message, body.error, body.errorCode);
                break;
            case 429:
                long retryAfter = parseRetryAfterHeader(original);
                RateLimitException rle = new RateLimitException(body.message,
                        body.error, body.errorCode, retryAfter);
                result = rle;
                break;
            default:
                if (statusCode >= 500) {
                    result = new ServerException(body.message, statusCode,
                            body.error, body.errorCode);
                } else {
                    // For other status codes, use the base ApiException — built here from
                    // the DTO, with the TRANSPORT status as the code.
                    result = new ApiException(body.message, statusCode,
                            body.error, body.errorCode);
                }
                break;
        }

        return result;
    }

    /**
     * Parses the Retry-After header to determine wait time.
     */
    private long parseRetryAfterHeader(com.taurushq.sdk.protect.openapi.ApiException e) {
        if (e.getResponseHeaders() != null) {
            java.util.List<String> retryAfterValues = e.getResponseHeaders().get("Retry-After");
            if (retryAfterValues != null && !retryAfterValues.isEmpty()) {
                try {
                    // Retry-After can be seconds or a date; we handle seconds here
                    int seconds = Integer.parseInt(retryAfterValues.get(0));
                    return seconds * 1000L;
                } catch (NumberFormatException ex) {
                    // Could be a date format, use default
                }
            }
        }
        return DEFAULT_RATE_LIMIT_RETRY_MS;
    }

    /**
     * The four fields the SDK recognises in an API error body — and nothing else.
     *
     * <p>This type exists so that a hostile response body never reaches a {@link Throwable}
     * through a reflective deserializer. It is a plain, flat, non-{@code Throwable} carrier:
     * no {@code cause}, no {@code suppressedExceptions}, no {@code stackTrace}, no
     * {@code originalException}, and no field whose type could pull another object graph in.
     * Anything else in the body is skipped by the parser and discarded.
     *
     * <p>Its values are COPIED into a freshly-constructed typed exception by
     * {@link #createTypedException}; the DTO instance itself is never handed to a caller.
     *
     * <p>Keep every field a scalar. Adding one of a nested or collection type would
     * reintroduce the attacker-controlled recursion depth this type removes.
     */
    private static final class ErrorBody {

        /** Short error label, e.g. {@code "NotFound"}. */
        private String error;

        /** Human-readable error message. Server-controlled text — do not echo it into logs. */
        private String message;

        /**
         * The HTTP code the BODY claims. Accepted from the wire so the claim can be
         * observed, and then deliberately dropped: the transport status is the only source
         * of {@link ApiException#code}. Never copy this into the returned exception.
         */
        private int code;

        /** Application-specific error code, e.g. {@code "ERR_001"}. */
        private String errorCode;
    }
}
