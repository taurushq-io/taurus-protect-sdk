package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents rules for whitelisting smart contract addresses.
 * <p>
 * Contract whitelisting rules define the approval requirements for adding new
 * smart contracts that can be interacted with. These are separate from regular
 * address whitelisting rules because contract interactions typically require
 * different governance.
 *
 * @see AddressWhitelistingRules
 * @see SequentialThresholds
 * @see DecodedRulesContainer
 */
public class ContractAddressWhitelistingRules {

    /**
     * Blockchain identifier (e.g., "ETH", "BTC"). Null or empty means global default.
     */
    private String blockchain;

    /**
     * Network identifier (e.g., "mainnet", "testnet"). Null or empty matches any network.
     */
    private String network;

    /**
     * Approval thresholds for contract whitelisting.
     * Multiple entries represent parallel approval paths.
     */
    private List<SequentialThresholds> parallelThresholds;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
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
     * Gets the parallel thresholds.
     *
     * @return the parallel thresholds
     */
    public List<SequentialThresholds> getParallelThresholds() {
        return parallelThresholds;
    }

    /**
     * Sets the parallel thresholds.
     *
     * @param parallelThresholds the parallel thresholds
     */
    public void setParallelThresholds(List<SequentialThresholds> parallelThresholds) {
        this.parallelThresholds = parallelThresholds;
    }
}
