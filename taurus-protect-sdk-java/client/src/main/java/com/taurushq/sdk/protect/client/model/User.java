package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a user in the Taurus Protect system.
 * <p>
 * Users are individuals who have access to the Taurus Protect platform with
 * specific roles and permissions. Users can participate in approval workflows,
 * create requests, and manage wallets based on their assigned roles.
 *
 * @see RequestTrail
 * @see WhitelistTrail
 */
public class User {

    /**
     * Unique identifier for the user.
     */
    private String id;

    /**
     * ID of the tenant (organization) this user belongs to.
     */
    private int tenantId;

    /**
     * External user ID from an identity provider (e.g., SSO system).
     */
    private String externalUserId;

    /**
     * User's first name.
     */
    private String firstName;

    /**
     * User's last name.
     */
    private String lastName;

    /**
     * Current status of the user (e.g., "active", "inactive").
     */
    private String status;

    /**
     * User's email address.
     */
    private String email;

    /**
     * User's login username.
     */
    private String username;

    /**
     * User's public key for cryptographic operations (e.g., signing approvals).
     */
    private String publicKey;

    /**
     * List of role names assigned to the user.
     */
    private List<String> roles;

    /**
     * Whether TOTP (Time-based One-Time Password) two-factor authentication is enabled.
     */
    private Boolean totpEnabled;

    /**
     * Whether this user is enforced in governance rules (must be part of approval workflow).
     */
    private Boolean enforcedInRules;

    /**
     * Timestamp when the user was created.
     */
    private OffsetDateTime createdAt;

    /**
     * Timestamp when the user was last updated.
     */
    private OffsetDateTime updatedAt;

    /**
     * Timestamp of the user's last login.
     */
    private OffsetDateTime lastLogin;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier for this user.
     *
     * @return the user ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this user.
     *
     * @param id the user ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Returns the tenant ID this user belongs to.
     *
     * @return the tenant ID
     */
    public int getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID this user belongs to.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(int tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Returns the external user ID from an identity provider.
     *
     * @return the external user ID, or {@code null} if not set
     */
    public String getExternalUserId() {
        return externalUserId;
    }

    /**
     * Sets the external user ID from an identity provider.
     *
     * @param externalUserId the external user ID to set
     */
    public void setExternalUserId(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    /**
     * Returns the user's first name.
     *
     * @return the first name
     */
    public String getFirstName() {
        return firstName;
    }

    /**
     * Sets the user's first name.
     *
     * @param firstName the first name to set
     */
    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    /**
     * Returns the user's last name.
     *
     * @return the last name
     */
    public String getLastName() {
        return lastName;
    }

    /**
     * Sets the user's last name.
     *
     * @param lastName the last name to set
     */
    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    /**
     * Returns the current status of the user.
     *
     * @return the user status (e.g., "active", "inactive")
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the current status of the user.
     *
     * @param status the user status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Returns the user's email address.
     *
     * @return the email address
     */
    public String getEmail() {
        return email;
    }

    /**
     * Sets the user's email address.
     *
     * @param email the email address to set
     */
    public void setEmail(String email) {
        this.email = email;
    }

    /**
     * Returns the user's login username.
     *
     * @return the username
     */
    public String getUsername() {
        return username;
    }

    /**
     * Sets the user's login username.
     *
     * @param username the username to set
     */
    public void setUsername(String username) {
        this.username = username;
    }

    /**
     * Returns the user's public key for cryptographic operations.
     *
     * @return the public key
     */
    public String getPublicKey() {
        return publicKey;
    }

    /**
     * Sets the user's public key for cryptographic operations.
     *
     * @param publicKey the public key to set
     */
    public void setPublicKey(String publicKey) {
        this.publicKey = publicKey;
    }

    /**
     * Returns the list of roles assigned to this user.
     *
     * @return the list of role names
     */
    public List<String> getRoles() {
        return roles;
    }

    /**
     * Sets the list of roles assigned to this user.
     *
     * @param roles the list of role names to set
     */
    public void setRoles(List<String> roles) {
        this.roles = roles;
    }

    /**
     * Returns whether TOTP two-factor authentication is enabled.
     *
     * @return {@code true} if TOTP is enabled, {@code false} otherwise
     */
    public Boolean getTotpEnabled() {
        return totpEnabled;
    }

    /**
     * Sets whether TOTP two-factor authentication is enabled.
     *
     * @param totpEnabled {@code true} to enable TOTP, {@code false} to disable
     */
    public void setTotpEnabled(Boolean totpEnabled) {
        this.totpEnabled = totpEnabled;
    }

    /**
     * Returns whether this user is enforced in governance rules.
     *
     * @return {@code true} if the user must be part of approval workflows
     */
    public Boolean getEnforcedInRules() {
        return enforcedInRules;
    }

    /**
     * Sets whether this user is enforced in governance rules.
     *
     * @param enforcedInRules {@code true} if the user must be part of approval workflows
     */
    public void setEnforcedInRules(Boolean enforcedInRules) {
        this.enforcedInRules = enforcedInRules;
    }

    /**
     * Returns the timestamp when the user was created.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the timestamp when the user was created.
     *
     * @param createdAt the creation timestamp to set
     */
    public void setCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    /**
     * Returns the timestamp when the user was last updated.
     *
     * @return the last update timestamp
     */
    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    /**
     * Sets the timestamp when the user was last updated.
     *
     * @param updatedAt the last update timestamp to set
     */
    public void setUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    /**
     * Returns the timestamp of the user's last login.
     *
     * @return the last login timestamp, or {@code null} if the user has never logged in
     */
    public OffsetDateTime getLastLogin() {
        return lastLogin;
    }

    /**
     * Sets the timestamp of the user's last login.
     *
     * @param lastLogin the last login timestamp to set
     */
    public void setLastLogin(OffsetDateTime lastLogin) {
        this.lastLogin = lastLogin;
    }
}
