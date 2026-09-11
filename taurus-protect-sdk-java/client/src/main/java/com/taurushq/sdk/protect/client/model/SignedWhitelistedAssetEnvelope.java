package com.taurushq.sdk.protect.client.model;

import com.taurushq.sdk.protect.client.helper.AssetHashHelper;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * The type SignedWhitelistedAssetEnvelope.
 * Contains a signed whitelisted asset (contract address) with all metadata required for verification.
 *
 * <p>When obtained via WhitelistedAssetService, the envelope is fully verified:
 * <ul>
 *   <li>Metadata hash verified</li>
 *   <li>Rules container signatures verified (SuperAdmin)</li>
 *   <li>Whitelist signatures verified (per governance rules)</li>
 * </ul>
 *
 * <p>Use {@link #getWhitelistedAsset()} to get the verified whitelisted asset.
 */
public class SignedWhitelistedAssetEnvelope {

    /**
     * Flag indicating whether the envelope has been verified and initialized.
     */
    private final AtomicBoolean isInitialized = new AtomicBoolean(false);

    /**
     * The verified and decoded whitelisted asset.
     * Set by WhitelistedAssetService after full cryptographic verification.
     */
    private WhitelistedAsset verifiedWhitelistedAsset;

    /**
     * The verified and decoded governance rules container.
     * Set by WhitelistedAssetService after SuperAdmin signature verification.
     */
    private DecodedRulesContainer verifiedRulesContainer;

    /**
     * Unique identifier for this whitelisted asset envelope.
     */
    private long id;

    /**
     * The tenant (organization) ID this asset belongs to.
     */
    private long tenantId;

    /**
     * Audit trail of all actions taken on this whitelisted asset.
     * Includes creation, approvals, rejections, and status changes.
     */
    private List<WhitelistTrail> trails = new ArrayList<>();

    /**
     * Custom attributes associated with this whitelisted asset.
     */
    private List<Attribute> attributes = new ArrayList<>();

    /**
     * The signed asset containing payload and cryptographic signatures.
     */
    private SignedWhitelistedAsset signedAsset;

    /**
     * Metadata about this whitelisted asset request including hash and timestamps.
     */
    private WhitelistMetadata metadata;

    /**
     * The action type for this whitelisting (e.g., "create", "delete").
     */
    private String action;

    /**
     * Base64-encoded governance rules container that was active when this asset was whitelisted.
     */
    private String rulesContainer;

    /**
     * The specific rule applied to this whitelisting request.
     */
    private String rule;

    /**
     * The approvers structure showing who needs to approve and who has approved.
     */
    private Approvers approvers;

    /**
     * The blockchain this asset contract is deployed on (e.g., "ETH", "MATIC").
     */
    private String blockchain;

    /**
     * Current status of this whitelisting request (e.g., "pending", "approved", "rejected").
     */
    private String status;

    /**
     * The network identifier (e.g., "mainnet", "goerli").
     */
    private String network;

    /**
     * Whether business rules are enabled for this asset.
     */
    private Boolean businessRuleEnabled;

    /**
     * Base64-encoded signatures for the rules container from SuperAdmin.
     */
    private String rulesSignatures;

    /**
     * Gets the verified and decoded rules container.
     * Only available after verification via WhitelistedAssetService.
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

    /**
     * Gets the audit trail for this whitelisted asset.
     *
     * @return the list of trail entries
     */
    public List<WhitelistTrail> getTrails() {
        return trails;
    }

    /**
     * Sets the audit trail.
     *
     * @param trails the trail entries to set
     */
    public void setTrails(List<WhitelistTrail> trails) {
        this.trails = trails;
    }

    /**
     * Gets the custom attributes associated with this asset.
     *
     * @return the list of attributes
     */
    public List<Attribute> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom attributes.
     *
     * @param attributes the attributes to set
     */
    public void setAttributes(List<Attribute> attributes) {
        this.attributes = attributes;
    }

    /**
     * Gets the unique identifier.
     *
     * @return the ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier.
     *
     * @param id the ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Gets the tenant ID.
     *
     * @return the tenant ID
     */
    public long getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(long tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Gets the signed asset containing payload and signatures.
     *
     * @return the signed asset
     */
    public SignedWhitelistedAsset getSignedAsset() {
        return signedAsset;
    }

    /**
     * Sets the signed asset.
     *
     * @param signedAsset the signed asset to set
     */
    public void setSignedAsset(SignedWhitelistedAsset signedAsset) {
        this.signedAsset = signedAsset;
    }

    /**
     * Gets the verified whitelisted asset.
     * The envelope must have been obtained via WhitelistedAssetService for this to work.
     *
     * @return the verified and decoded WhitelistedAsset
     * @throws IllegalStateException if the envelope has not been verified yet
     */
    public WhitelistedAsset getWhitelistedAsset() {
        if (!isInitialized.get()) {
            throw new IllegalStateException(
                    "Envelope not initialized - obtain via WhitelistedAssetService");
        }
        return this.verifiedWhitelistedAsset;
    }

    /**
     * Marks this envelope verified, DERIVING the asset from its own signed payload.
     *
     * <p>It does not accept an asset. The previous
     * {@code setVerifiedWhitelistedAsset(WhitelistedAsset)} was public and flipped the
     * same {@code isInitialized} gate {@link #getWhitelistedAsset()} checks, so a caller
     * could construct an envelope, inject a fabricated asset, and read it back with no
     * exception — the "verified" marker returned attacker-chosen data. Deriving instead
     * of accepting gives Java the property Go's {@code helper.VerifiedAsset} has:
     * forging the marker is possible but useless, because the data still comes from the
     * envelope's own payload.
     *
     * <p>Intended for WhitelistedAssetService; it must be public because that service is
     * in another package. That is also why the full fix is the verifier extraction in
     * {@code TODOS.md} — a witness type the verifier alone can mint. This closes the
     * exploitable half without it.
     *
     * <p><b>{@code verifiedPayload} is the payload the matched signature COVERED, which
     * is not always the payload the server delivered.</b> Step 4 accepts
     * backward-compatible rewrites of the delivered payload ({@code isNFT},
     * {@code kindType}), and those strips are not injective, so parsing the DELIVERED
     * text can return members no signature covered. The address flow is where that is
     * directly exploitable today; the asset flow takes the same parameter so the two
     * cannot drift, and so a schema change that makes a stripped asset field readable
     * does not silently reopen it here. {@code metadata.getPayloadAsString()} is left
     * untouched, because a caller needs it to reproduce {@code metadata.hash}.
     *
     * @param rulesContainer  the verified rules container
     * @param verifiedPayload the payload the matched signature covered — step 4's
     *                        {@code LegacyPayloadVariant.getPayload()}, NOT
     *                        {@code metadata.payloadAsString}
     * @throws WhitelistException if the signed payload is absent or unparseable
     */
    public void markVerified(final DecodedRulesContainer rulesContainer,
                             final String verifiedPayload)
            throws WhitelistException {
        if (this.metadata == null || Strings.isNullOrEmpty(this.metadata.getPayloadAsString())) {
            throw new WhitelistException(
                    "cannot mark verified: the envelope carries no signed payload");
        }
        if (Strings.isNullOrEmpty(verifiedPayload)) {
            // Not a defaulting opportunity: falling back to metadata.payloadAsString is
            // exactly the defect this parameter exists to close.
            throw new WhitelistException("cannot mark verified: no verified payload was "
                    + "supplied; step 4 must hand over the payload it matched");
        }
        this.verifiedWhitelistedAsset =
                AssetHashHelper.parseWhitelistedAssetFromJson(verifiedPayload);
        this.verifiedRulesContainer = rulesContainer;
        this.isInitialized.set(true);
    }

    /**
     * Gets the metadata for this whitelisting request.
     *
     * @return the metadata
     */
    public WhitelistMetadata getMetadata() {
        return metadata;
    }

    /**
     * Sets the metadata.
     *
     * @param metadata the metadata to set
     */
    public void setMetadata(WhitelistMetadata metadata) {
        this.metadata = metadata;
    }

    /**
     * Gets the action type for this whitelisting.
     *
     * @return the action (e.g., "create", "delete")
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the action type.
     *
     * @param action the action to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Gets the base64-encoded rules container.
     *
     * @return the encoded rules container
     */
    public String getRulesContainer() {
        return rulesContainer;
    }

    /**
     * Sets the rules container.
     *
     * @param rulesContainer the encoded rules container to set
     */
    public void setRulesContainer(String rulesContainer) {
        this.rulesContainer = rulesContainer;
    }

    /**
     * Gets the rule applied to this whitelisting.
     *
     * @return the rule identifier
     */
    public String getRule() {
        return rule;
    }

    /**
     * Sets the rule.
     *
     * @param rule the rule to set
     */
    public void setRule(String rule) {
        this.rule = rule;
    }

    /**
     * Gets the approvers structure.
     *
     * @return the approvers
     */
    public Approvers getApprovers() {
        return approvers;
    }

    /**
     * Sets the approvers.
     *
     * @param approvers the approvers to set
     */
    public void setApprovers(Approvers approvers) {
        this.approvers = approvers;
    }

    /**
     * Gets the blockchain identifier.
     *
     * @return the blockchain (e.g., "ETH", "MATIC")
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain.
     *
     * @param blockchain the blockchain to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the current status.
     *
     * @return the status (e.g., "pending", "approved", "rejected")
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the status.
     *
     * @param status the status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Gets the network identifier.
     *
     * @return the network (e.g., "mainnet", "goerli")
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network.
     *
     * @param network the network to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets whether business rules are enabled.
     *
     * @return true if business rules are enabled, false otherwise
     */
    public Boolean getBusinessRuleEnabled() {
        return businessRuleEnabled;
    }

    /**
     * Sets whether business rules are enabled.
     *
     * @param businessRuleEnabled the business rule enabled flag to set
     */
    public void setBusinessRuleEnabled(Boolean businessRuleEnabled) {
        this.businessRuleEnabled = businessRuleEnabled;
    }

    /**
     * Gets the base64-encoded rules signatures.
     *
     * @return the rules signatures
     */
    public String getRulesSignatures() {
        return rulesSignatures;
    }

    /**
     * Sets the rules signatures.
     *
     * @param rulesSignatures the rules signatures to set
     */
    public void setRulesSignatures(String rulesSignatures) {
        this.rulesSignatures = rulesSignatures;
    }
}
