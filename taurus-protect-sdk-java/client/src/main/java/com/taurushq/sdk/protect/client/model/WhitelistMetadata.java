package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents metadata associated with a whitelist entry.
 * <p>
 * Whitelist metadata contains additional information about a whitelisted
 * address, including a cryptographic hash for integrity verification and
 * the raw payload data.
 *
 * <h2>SECURITY DESIGN</h2>
 * <p>
 * The API returns metadata with two representations of the same data:
 * <ul>
 *   <li>{@code payload} - Raw JSON object (UNVERIFIED)</li>
 *   <li>{@code payloadAsString} - JSON string that is cryptographically hashed (VERIFIED)</li>
 * </ul>
 * <p>
 * The security model works as follows:
 * <ol>
 *   <li>The server computes: {@code metadata.hash = SHA256(payloadAsString)}</li>
 *   <li>The hash is signed by governance rules (SuperAdmin keys)</li>
 *   <li>Clients verify: {@code computed_hash(payloadAsString) == metadata.hash}</li>
 * </ol>
 * <p>
 * <strong>CRITICAL:</strong> Security-critical fields (address, blockchain, network, label,
 * memo, contractType) MUST be extracted from {@code payloadAsString} via
 * {@link com.taurushq.sdk.protect.client.helper.WhitelistHashHelper#parseWhitelistedAddressFromJson(String)}.
 * <p>
 * <strong>ATTACK VECTOR (if using raw payload):</strong>
 * An attacker intercepting API responses could modify the payload object
 * (e.g., change destination address) while leaving payloadAsString unchanged.
 * The hash would still verify, but the client would extract tampered data.
 *
 * @see WhitelistedAddress
 * @see SignedWhitelistedAddress
 */
public class WhitelistMetadata {

    /**
     * Cryptographic hash of the metadata payload for integrity verification.
     */
    private String hash;

    // SECURITY: payload field intentionally omitted - use payloadAsString only.
    // The raw payload object could be tampered with by an attacker while
    // payloadAsString remains unchanged (hash still verifies). By not having
    // this field, we enforce that all data extraction uses the verified source.
    // Use WhitelistHashHelper.parseWhitelistedAddressFromJson(payloadAsString) for extraction.

    /**
     * String representation of the payload.
     */
    private String payloadAsString;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the cryptographic hash of the metadata payload.
     *
     * @return the metadata hash
     */
    public String getHash() {
        return hash;
    }

    /**
     * Sets the cryptographic hash of the metadata payload.
     *
     * @param hash the metadata hash to set
     */
    public void setHash(String hash) {
        this.hash = hash;
    }

    // SECURITY: getPayload() and setPayload() intentionally removed.
    // Use WhitelistHashHelper.parseWhitelistedAddressFromJson(payloadAsString) for extraction.

    /**
     * Returns the payload as a string.
     *
     * @return the payload as a string
     */
    public String getPayloadAsString() {
        return payloadAsString;
    }

    /**
     * Sets the payload as a string.
     *
     * @param payloadAsString the payload string to set
     */
    public void setPayloadAsString(String payloadAsString) {
        this.payloadAsString = payloadAsString;
    }
}
