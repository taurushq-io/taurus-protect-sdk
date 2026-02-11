package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.security.PublicKey;
import java.util.List;

/**
 * Represents a user defined in the governance rules container.
 * <p>
 * Each user has a unique ID, a public key for signature verification, and a set
 * of roles that determine their permissions within the system. Special roles include:
 * <ul>
 *   <li><b>HSMSLOT</b> - Indicates an HSM engine slot (not a human user)</li>
 *   <li><b>SuperAdmin</b> - Administrative privileges for rules management</li>
 * </ul>
 *
 * @see RuleGroup
 * @see DecodedRulesContainer
 */
public class RuleUser {

    /**
     * Unique identifier for the user within the rules container.
     */
    private String id;

    /**
     * PEM-encoded public key for signature verification.
     */
    private String publicKeyPem;

    /**
     * Decoded Java PublicKey object (populated after parsing the PEM string).
     */
    private PublicKey publicKey;

    /**
     * List of role names assigned to this user (e.g., "SuperAdmin", "HSMSLOT").
     */
    private List<String> roles;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the user id.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the user id.
     *
     * @param id the id
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the public key in PEM format.
     *
     * @return the public key PEM string
     */
    public String getPublicKeyPem() {
        return publicKeyPem;
    }

    /**
     * Sets the public key in PEM format.
     *
     * @param publicKeyPem the public key PEM string
     */
    public void setPublicKeyPem(String publicKeyPem) {
        this.publicKeyPem = publicKeyPem;
    }

    /**
     * Gets the decoded public key.
     *
     * @return the public key
     */
    public PublicKey getPublicKey() {
        return publicKey;
    }

    /**
     * Sets the decoded public key.
     *
     * @param publicKey the public key
     */
    public void setPublicKey(PublicKey publicKey) {
        this.publicKey = publicKey;
    }

    /**
     * Gets the user roles.
     *
     * @return the roles
     */
    public List<String> getRoles() {
        return roles;
    }

    /**
     * Sets the user roles.
     *
     * @param roles the roles
     */
    public void setRoles(List<String> roles) {
        this.roles = roles;
    }
}
