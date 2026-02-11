package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a complete whitelisted contract address envelope with all associated data.
 * <p>
 * The envelope contains the signed contract address data along with metadata, approval
 * trails, attributes, and the current status. This is the primary model returned when
 * querying whitelisted contract addresses through the API.
 * <p>
 * Key components:
 * <ul>
 *   <li>signedContractAddress - The signed contract with cryptographic signatures</li>
 *   <li>metadata - Additional metadata associated with the entry</li>
 *   <li>trails - Audit trail of all actions taken on this entry</li>
 *   <li>approvers - Required approvers for this whitelist entry</li>
 *   <li>attributes - Custom attributes attached to the entry</li>
 * </ul>
 *
 * @see SignedWhitelistedContractAddress
 * @see WhitelistedContractAddress
 * @see ContractWhitelistingService
 */
public class SignedWhitelistedContractAddressEnvelope {

    /**
     * Unique identifier for this whitelist entry.
     */
    private String id;

    /**
     * Tenant identifier.
     */
    private String tenantId;

    /**
     * The signed contract address with signatures and payload.
     */
    private SignedWhitelistedContractAddress signedContractAddress;

    /**
     * Metadata associated with this whitelist entry.
     */
    private WhitelistMetadata metadata;

    /**
     * The action type for this entry (e.g., "create", "update", "delete").
     */
    private String action;

    /**
     * Audit trail of actions taken on this entry.
     */
    private List<WhitelistTrail> trails;

    /**
     * Rules container identifier.
     */
    private String rulesContainer;

    /**
     * Rule identifier.
     */
    private String rule;

    /**
     * Approvers configuration for this entry.
     */
    private Approvers approvers;

    /**
     * Custom attributes attached to this entry.
     */
    private List<Attribute> attributes;

    /**
     * Current status of the whitelist entry (e.g., "pending", "approved", "rejected").
     */
    private String status;

    /**
     * Blockchain identifier (e.g., "ETH", "MATIC").
     */
    private String blockchain;

    /**
     * Network identifier (e.g., "mainnet", "goerli").
     */
    private String network;

    /**
     * Whether the business rule is enabled for this contract.
     */
    private Boolean businessRuleEnabled;

    /**
     * Signatures of the rules container for verification.
     */
    private String rulesSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the unique identifier.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier.
     *
     * @param id the id to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the tenant identifier.
     *
     * @return the tenant id
     */
    public String getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant identifier.
     *
     * @param tenantId the tenant id to set
     */
    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Gets the signed contract address.
     *
     * @return the signed contract address
     */
    public SignedWhitelistedContractAddress getSignedContractAddress() {
        return signedContractAddress;
    }

    /**
     * Sets the signed contract address.
     *
     * @param signedContractAddress the signed contract address to set
     */
    public void setSignedContractAddress(SignedWhitelistedContractAddress signedContractAddress) {
        this.signedContractAddress = signedContractAddress;
    }

    /**
     * Gets the metadata.
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
     * Gets the action type.
     *
     * @return the action
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
     * Gets the audit trails.
     *
     * @return the trails
     */
    public List<WhitelistTrail> getTrails() {
        return trails;
    }

    /**
     * Sets the audit trails.
     *
     * @param trails the trails to set
     */
    public void setTrails(List<WhitelistTrail> trails) {
        this.trails = trails;
    }

    /**
     * Gets the rules container identifier.
     *
     * @return the rules container
     */
    public String getRulesContainer() {
        return rulesContainer;
    }

    /**
     * Sets the rules container identifier.
     *
     * @param rulesContainer the rules container to set
     */
    public void setRulesContainer(String rulesContainer) {
        this.rulesContainer = rulesContainer;
    }

    /**
     * Gets the rule identifier.
     *
     * @return the rule
     */
    public String getRule() {
        return rule;
    }

    /**
     * Sets the rule identifier.
     *
     * @param rule the rule to set
     */
    public void setRule(String rule) {
        this.rule = rule;
    }

    /**
     * Gets the approvers configuration.
     *
     * @return the approvers
     */
    public Approvers getApprovers() {
        return approvers;
    }

    /**
     * Sets the approvers configuration.
     *
     * @param approvers the approvers to set
     */
    public void setApprovers(Approvers approvers) {
        this.approvers = approvers;
    }

    /**
     * Gets the custom attributes.
     *
     * @return the attributes
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
     * Gets the current status.
     *
     * @return the status
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the current status.
     *
     * @param status the status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Gets the blockchain identifier.
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain identifier.
     *
     * @param blockchain the blockchain to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the network identifier.
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network identifier.
     *
     * @param network the network to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets whether the business rule is enabled.
     *
     * @return true if business rule is enabled
     */
    public Boolean getBusinessRuleEnabled() {
        return businessRuleEnabled;
    }

    /**
     * Sets whether the business rule is enabled.
     *
     * @param businessRuleEnabled the business rule enabled flag
     */
    public void setBusinessRuleEnabled(Boolean businessRuleEnabled) {
        this.businessRuleEnabled = businessRuleEnabled;
    }

    /**
     * Gets the rules signatures.
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
