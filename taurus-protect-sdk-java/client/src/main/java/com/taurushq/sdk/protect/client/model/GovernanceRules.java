package com.taurushq.sdk.protect.client.model;

import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.apache.commons.codec.binary.Base64;
import org.apache.commons.lang3.builder.ToStringBuilder;

import java.security.PublicKey;
import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a governance ruleset in the Taurus Protect system.
 * <p>
 * Governance rules define the approval workflows and policies for transaction
 * requests. The rules are stored in a signed, protobuf-encoded container that
 * must be verified before use. Changes to governance rules require multi-signature
 * approval from super administrators.
 * <p>
 * Example usage:
 * <pre>{@code
 * GovernanceRules rules = client.getGovernanceService().getCurrentRules();
 * DecodedRulesContainer decoded = rules.getDecodedRulesContainer(publicKeys, 2);
 * }</pre>
 *
 * @see RuleUserSignature
 * @see GovernanceRulesTrail
 * @see com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer
 */
public class GovernanceRules {

    /**
     * Base64-encoded protobuf container holding the governance rules.
     */
    private String rulesContainer;

    /**
     * List of super administrator signatures authorizing these rules.
     */
    private List<RuleUserSignature> rulesSignatures;

    /**
     * Cached decoded rules container after signature verification.
     */
    private DecodedRulesContainer decodedRulesContainer;

    /**
     * Lock object for thread-safe lazy initialization of decoded rules container.
     */
    private final Object decodedLock = new Object();

    /**
     * Whether these rules are locked (cannot be modified).
     */
    private Boolean locked;

    /**
     * Timestamp when these rules were created.
     */
    private OffsetDateTime creationDate;

    /**
     * Timestamp when these rules were last updated.
     */
    private OffsetDateTime updateDate;

    /**
     * Audit trail of changes made to these rules.
     */
    private List<GovernanceRulesTrail> trails;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets rules container (base64 encoded).
     *
     * @return the rules container
     */
    public String getRulesContainer() {
        return rulesContainer;
    }

    /**
     * Sets rules container (base64 encoded).
     *
     * @param rulesContainer the rules container
     */
    public void setRulesContainer(String rulesContainer) {
        this.rulesContainer = rulesContainer;
        this.decodedRulesContainer = null; // Reset cache
    }

    /**
     * Verifies signatures and returns the decoded rules container.
     * The rules container is decoded and cached on first successful verification.
     * This method is thread-safe.
     *
     * @param superAdminPublicKeys the list of SuperAdmin public keys for verification
     * @param minValidSignatures   the minimum number of valid signatures required
     * @return the decoded rules container, or null if rulesContainer is null
     * @throws IntegrityException if signature verification fails
     */
    public DecodedRulesContainer getDecodedRulesContainer(
            List<PublicKey> superAdminPublicKeys,
            int minValidSignatures) throws IntegrityException {
        synchronized (decodedLock) {
            if (decodedRulesContainer == null && rulesContainer != null) {
                // Verify signatures first
                SignatureVerifier.verifyGovernanceRules(this, minValidSignatures, superAdminPublicKeys);

                // Decode after verification passes
                byte[] bytes = Base64.decodeBase64(rulesContainer);
                try {
                    RequestReply.RulesContainer proto = RequestReply.RulesContainer.parseFrom(bytes);
                    this.decodedRulesContainer = RulesContainerMapper.INSTANCE.fromProto(proto);
                } catch (InvalidProtocolBufferException e) {
                    throw new IntegrityException("unable to decode the rules container from proto", e);
                }
            }
            return this.decodedRulesContainer;
        }
    }

    /**
     * Returns the list of super administrator signatures authorizing these rules.
     *
     * @return the list of rule signatures
     */
    public List<RuleUserSignature> getRulesSignatures() {
        return rulesSignatures;
    }

    /**
     * Sets the list of super administrator signatures authorizing these rules.
     *
     * @param rulesSignatures the list of rule signatures to set
     */
    public void setRulesSignatures(List<RuleUserSignature> rulesSignatures) {
        this.rulesSignatures = rulesSignatures;
    }

    /**
     * Returns whether these rules are locked (cannot be modified).
     *
     * @return {@code true} if the rules are locked, {@code false} otherwise
     */
    public Boolean getLocked() {
        return locked;
    }

    /**
     * Sets whether these rules are locked (cannot be modified).
     *
     * @param locked {@code true} to lock the rules, {@code false} to unlock
     */
    public void setLocked(Boolean locked) {
        this.locked = locked;
    }

    /**
     * Returns the timestamp when these rules were created.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the timestamp when these rules were created.
     *
     * @param creationDate the creation timestamp to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Returns the timestamp when these rules were last updated.
     *
     * @return the last update timestamp
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the timestamp when these rules were last updated.
     *
     * @param updateDate the last update timestamp to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Returns the audit trail of changes made to these rules.
     *
     * @return the list of audit trail entries
     */
    public List<GovernanceRulesTrail> getTrails() {
        return trails;
    }

    /**
     * Sets the audit trail of changes made to these rules.
     *
     * @param trails the list of audit trail entries to set
     */
    public void setTrails(List<GovernanceRulesTrail> trails) {
        this.trails = trails;
    }


}
