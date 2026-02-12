package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents detailed metadata for a transaction rule set.
 * <p>
 * These details specify what type of transactions the rule applies to,
 * including the action domain, blockchain, network, and any contract-specific
 * configurations.
 *
 * @see TransactionRules
 * @see EvmCallContract
 * @see XtzCallContract
 */
public class TransactionRuleDetails {

    /**
     * The rule domain (e.g., "RuleDomainTransfer", "RuleDomainStaking", "RuleDomainContract").
     */
    private String domain;

    /**
     * The rule sub-domain for more specific categorization.
     */
    private String subDomain;

    /**
     * Blockchain identifier (e.g., "ETH", "BTC"). Null means applies to all blockchains.
     */
    private String blockchain;

    /**
     * Network identifier (e.g., "mainnet", "testnet"). Null means applies to all networks.
     */
    private String network;

    /**
     * EVM-specific contract call configuration (for Ethereum-compatible chains).
     */
    private EvmCallContract evmCallContract;

    /**
     * Tezos-specific contract call configuration.
     */
    private XtzCallContract xtzCallContract;

    /**
     * Cash settlement configuration for fiat-related operations.
     */
    private CashSettlement cashSettlement;

    /**
     * List of Cosmos message type signatures this rule applies to.
     */
    private List<String> cosmosMethodSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the rule domain.
     *
     * @return the domain (e.g., RuleDomainTransfer, RuleDomainStaking)
     */
    public String getDomain() {
        return domain;
    }

    /**
     * Sets the rule domain.
     *
     * @param domain the domain
     */
    public void setDomain(String domain) {
        this.domain = domain;
    }

    /**
     * Gets the rule sub-domain.
     *
     * @return the sub-domain
     */
    public String getSubDomain() {
        return subDomain;
    }

    /**
     * Sets the rule sub-domain.
     *
     * @param subDomain the sub-domain
     */
    public void setSubDomain(String subDomain) {
        this.subDomain = subDomain;
    }

    /**
     * Gets the blockchain.
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain.
     *
     * @param blockchain the blockchain
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the network.
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network.
     *
     * @param network the network
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets the EVM call contract details.
     *
     * @return the EVM call contract
     */
    public EvmCallContract getEvmCallContract() {
        return evmCallContract;
    }

    /**
     * Sets the EVM call contract details.
     *
     * @param evmCallContract the EVM call contract
     */
    public void setEvmCallContract(EvmCallContract evmCallContract) {
        this.evmCallContract = evmCallContract;
    }

    /**
     * Gets the XTZ call contract details.
     *
     * @return the XTZ call contract
     */
    public XtzCallContract getXtzCallContract() {
        return xtzCallContract;
    }

    /**
     * Sets the XTZ call contract details.
     *
     * @param xtzCallContract the XTZ call contract
     */
    public void setXtzCallContract(XtzCallContract xtzCallContract) {
        this.xtzCallContract = xtzCallContract;
    }

    /**
     * Gets the cash settlement details.
     *
     * @return the cash settlement
     */
    public CashSettlement getCashSettlement() {
        return cashSettlement;
    }

    /**
     * Sets the cash settlement details.
     *
     * @param cashSettlement the cash settlement
     */
    public void setCashSettlement(CashSettlement cashSettlement) {
        this.cashSettlement = cashSettlement;
    }

    /**
     * Gets the Cosmos method signatures.
     *
     * @return the Cosmos method signatures
     */
    public List<String> getCosmosMethodSignatures() {
        return cosmosMethodSignatures;
    }

    /**
     * Sets the Cosmos method signatures.
     *
     * @param cosmosMethodSignatures the Cosmos method signatures
     */
    public void setCosmosMethodSignatures(List<String> cosmosMethodSignatures) {
        this.cosmosMethodSignatures = cosmosMethodSignatures;
    }
}
