package com.taurushq.sdk.protect.client.model;

import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * The type SignedWhitelistedAddressEnvelope.
 * Contains a signed whitelisted address with all metadata required for verification.
 *
 * <p>When obtained via WhitelistedAddressService, the envelope is fully verified:
 * <ul>
 *   <li>Metadata hash verified</li>
 *   <li>Rules container signatures verified (SuperAdmin)</li>
 *   <li>Whitelist signatures verified (per governance rules)</li>
 * </ul>
 *
 * <p>Use {@link #getWhitelistedAddress()} to get the verified whitelisted address.
 */
public class SignedWhitelistedAddressEnvelope {

    private final AtomicBoolean isInitialized = new AtomicBoolean(false);

    // Verified data - set by WhitelistedAddressService after full verification
    private WhitelistedAddress verifiedWhitelistedAddress;
    private DecodedRulesContainer verifiedRulesContainer;

    /**
     * Unique identifier for this whitelisted address.
     */
    private long id;

    /**
     * ID of the tenant (organization) this address belongs to.
     */
    private long tenantId;

    /**
     * Risk assessment scores for this address (e.g., AML compliance scores).
     */
    private List<Score> scores = new ArrayList<>();

    /**
     * Audit trail of actions taken on this whitelisted address.
     */
    private List<WhitelistTrail> trails = new ArrayList<>();

    /**
     * Custom attributes associated with this address.
     */
    private List<Attribute> attributes = new ArrayList<>();

    /**
     * The signed address containing payload and signatures.
     */
    private SignedWhitelistedAddress signedAddress;

    /**
     * Metadata including hash and payload information.
     */
    private WhitelistMetadata metadata;

    /**
     * The action being performed (e.g., "create", "update", "delete").
     */
    private String action;

    /**
     * Base64-encoded governance rules container.
     */
    private String rulesContainer;

    /**
     * The governance rule applied to this whitelisting.
     */
    private String rule;

    /**
     * Approvers configuration for this whitelisting request.
     */
    private Approvers approvers;

    /**
     * Hash of the rules container for integrity verification.
     */
    private String rulesContainerHash;

    /**
     * The blockchain this address belongs to (e.g., "ethereum", "bitcoin").
     */
    private String blockchain;

    /**
     * Current status of this whitelisted address (e.g., "pending", "active").
     */
    private String status;

    /**
     * The network within the blockchain (e.g., "mainnet", "testnet").
     */
    private String network;

    /**
     * Visibility group ID controlling access to this address.
     */
    private String visibilityGroupID;

    /**
     * Taurus Network participant ID if applicable.
     */
    private String tnParticipantID;

    /**
     * Signatures on the rules container.
     */
    private String rulesSignatures;


    /**
     * Gets the verified and decoded rules container.
     * Only available after verification via WhitelistedAddressService.
     *
     * @return the decoded rules container, or null if not verified yet
     */
    public DecodedRulesContainer getDecodedRulesContainer() {
        return verifiedRulesContainer;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    public List<Score> getScores() {
        return scores;
    }

    public void setScores(List<Score> scores) {
        this.scores = scores;
    }

    public List<WhitelistTrail> getTrails() {
        return trails;
    }

    public void setTrails(List<WhitelistTrail> trails) {
        this.trails = trails;
    }

    public List<Attribute> getAttributes() {
        return attributes;
    }

    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }

    public long getId() {
        return id;
    }

    public void setId(long id) {
        this.id = id;
    }

    public long getTenantId() {
        return tenantId;
    }

    public void setTenantId(long tenantId) {
        this.tenantId = tenantId;
    }

    public SignedWhitelistedAddress getSignedAddress() {
        return signedAddress;
    }

    public void setSignedAddress(SignedWhitelistedAddress signedAddress) {
        this.signedAddress = signedAddress;
    }

    /**
     * Gets the verified whitelisted address.
     * The envelope must have been obtained via WhitelistedAddressService for this to work.
     *
     * @return the verified and decoded WhitelistedAddress
     * @throws IllegalStateException if the envelope has not been verified yet
     */
    public WhitelistedAddress getWhitelistedAddress() {
        if (!isInitialized.get()) {
            throw new IllegalStateException(
                    "Envelope not initialized - obtain via WhitelistedAddressService");
        }
        return this.verifiedWhitelistedAddress;
    }

    /**
     * Sets the verified whitelisted address.
     * This method is intended for internal use by WhitelistedAddressService.
     * Users should not call this method directly.
     *
     * @param address the verified whitelisted address
     */
    public void setVerifiedWhitelistedAddress(WhitelistedAddress address) {
        this.verifiedWhitelistedAddress = address;
        this.isInitialized.set(true);
    }

    /**
     * Sets the verified rules container.
     * This method is intended for internal use by WhitelistedAddressService.
     * Users should not call this method directly.
     *
     * @param rulesContainer the verified rules container
     */
    public void setVerifiedRulesContainer(DecodedRulesContainer rulesContainer) {
        this.verifiedRulesContainer = rulesContainer;
    }

    public WhitelistMetadata getMetadata() {
        return metadata;
    }

    public void setMetadata(WhitelistMetadata metadata) {
        this.metadata = metadata;
    }

    public String getAction() {
        return action;
    }

    public void setAction(String action) {
        this.action = action;
    }

    public String getRulesContainer() {
        return rulesContainer;
    }

    public void setRulesContainer(String rulesContainer) {
        this.rulesContainer = rulesContainer;
    }

    public String getRule() {
        return rule;
    }

    public void setRule(String rule) {
        this.rule = rule;
    }

    public Approvers getApprovers() {
        return approvers;
    }

    public void setApprovers(Approvers approvers) {
        this.approvers = approvers;
    }

    public String getRulesContainerHash() {
        return rulesContainerHash;
    }

    public void setRulesContainerHash(String rulesContainerHash) {
        this.rulesContainerHash = rulesContainerHash;
    }

    public String getBlockchain() {
        return blockchain;
    }

    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getNetwork() {
        return network;
    }

    public void setNetwork(String network) {
        this.network = network;
    }

    public String getVisibilityGroupID() {
        return visibilityGroupID;
    }

    public void setVisibilityGroupID(String visibilityGroupID) {
        this.visibilityGroupID = visibilityGroupID;
    }

    public String getTnParticipantID() {
        return tnParticipantID;
    }

    public void setTnParticipantID(String tnParticipantID) {
        this.tnParticipantID = tnParticipantID;
    }

    public String getRulesSignatures() {
        return rulesSignatures;
    }

    public void setRulesSignatures(String rulesSignatures) {
        this.rulesSignatures = rulesSignatures;
    }
}
