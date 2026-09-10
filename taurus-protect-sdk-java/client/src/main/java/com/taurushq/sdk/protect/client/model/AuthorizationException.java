package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Exception thrown when the authenticated user lacks permission for the requested operation.
 * <p>
 * This corresponds to HTTP 403 Forbidden responses. Common causes include:
 * <ul>
 *   <li>Insufficient user permissions</li>
 *   <li>Operation not allowed for the user's role</li>
 *   <li>Resource access denied</li>
 *   <li>Visibility group restrictions</li>
 * </ul>
 * <p>
 * Example:
 * <pre>{@code
 * try {
 *     client.getRequestService().approveRequests(requests, privateKey);
 * } catch (AuthorizationException e) {
 *     System.out.println("You don't have permission to approve requests");
 * }
 * }</pre>
 */
public class AuthorizationException extends ApiException {

    /**
     * Matches how Taurus-PROTECT reports a failed role check, for both its all-of and
     * any-of checks.
     * <p>
     * Unanchored on purpose: the server wraps the gRPC status, so this arrives as
     * "pre-filter failed: … desc = one of the '…' role is required". Anchoring it
     * matches nothing.
     */
    private static final Pattern REQUIRED_ROLES =
            Pattern.compile("one of the '([^']*)' role is required");

    private List<String> requiredRoles = Collections.emptyList();

    /**
     * Default constructor.
     */
    public AuthorizationException() {
        setCode(403);
    }

    /**
     * Constructs an AuthorizationException with the specified message.
     *
     * @param message the error message
     */
    public AuthorizationException(String message) {
        super(message, 403);
        this.requiredRoles = parseRequiredRoles(message);
    }

    /**
     * Constructs an AuthorizationException with full details.
     *
     * @param message   the error message
     * @param error     the error description
     * @param errorCode the application-specific error code
     */
    public AuthorizationException(String message, String error, String errorCode) {
        super(message, 403, error, errorCode);
        this.requiredRoles = parseRequiredRoles(message);
    }

    /**
     * Returns the roles that would satisfy the failed check, letting a caller say which
     * role to ask for instead of just "forbidden".
     *
     * @return the roles named by the server; empty when the denial was not role-based.
     *         One entry is a required role; several mean any one of them suffices.
     */
    public List<String> getRequiredRoles() {
        return Collections.unmodifiableList(requiredRoles);
    }

    /**
     * Extracts the roles named in a 403 message. Role names are lowercase alphanumeric,
     * so " - " is an unambiguous separator.
     *
     * @param message the server error message
     * @return the roles named, or an empty list when the message is not a role check
     */
    public static List<String> parseRequiredRoles(String message) {
        if (message == null) {
            return Collections.emptyList();
        }

        Matcher matcher = REQUIRED_ROLES.matcher(message);
        if (!matcher.find()) {
            return Collections.emptyList();
        }

        List<String> roles = new ArrayList<>();
        for (String role : matcher.group(1).split(" - ")) {
            String trimmed = role.trim();
            if (!trimmed.isEmpty()) {
                roles.add(trimmed);
            }
        }

        return roles;
    }

    @Override
    public boolean isRetryable() {
        return false;
    }
}
