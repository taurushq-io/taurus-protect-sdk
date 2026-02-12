package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents address whitelisting rules for a specific blockchain and network.
 * <p>
 * These rules define the approval requirements for whitelisting external addresses.
 * Rules can be configured per blockchain (currency) and network combination, with
 * fallback to global defaults when specific rules are not defined.
 * <p>
 * The approval structure consists of:
 * <ul>
 *   <li><b>Parallel thresholds</b> - Default approval requirements (groups running in parallel)</li>
 *   <li><b>Lines</b> - Source-specific overrides based on source wallet restrictions</li>
 * </ul>
 *
 * @see AddressWhitelistingLine
 * @see SequentialThresholds
 * @see DecodedRulesContainer
 */
public class AddressWhitelistingRules {

    /**
     * Blockchain identifier (e.g., "ETH", "BTC"). Null or empty means global default.
     */
    private String currency;

    /**
     * Network identifier (e.g., "mainnet", "testnet"). Null or empty matches any network.
     */
    private String network;

    /**
     * Default approval thresholds when no source-specific line matches.
     * Each entry represents a parallel approval path.
     */
    private List<SequentialThresholds> parallelThresholds;

    /**
     * Source-specific rule lines with custom thresholds based on source wallet.
     */
    private List<AddressWhitelistingLine> lines;

    /**
     * Whether to include the network in the payload hash during signature verification.
     */
    private boolean includeNetworkInPayload;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the currency (blockchain identifier).
     *
     * @return the currency
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency (blockchain identifier).
     *
     * @param currency the currency
     */
    public void setCurrency(String currency) {
        this.currency = currency;
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

    /**
     * Gets the rule lines for source-specific threshold overrides.
     *
     * @return the lines, or null if no source-specific rules defined
     */
    public List<AddressWhitelistingLine> getLines() {
        return lines;
    }

    /**
     * Sets the rule lines.
     *
     * @param lines the lines
     */
    public void setLines(List<AddressWhitelistingLine> lines) {
        this.lines = lines;
    }

    /**
     * Gets whether network should be included in the payload hash.
     *
     * @return true if network should be included
     */
    public boolean isIncludeNetworkInPayload() {
        return includeNetworkInPayload;
    }

    /**
     * Sets whether network should be included in the payload hash.
     *
     * @param includeNetworkInPayload true to include network
     */
    public void setIncludeNetworkInPayload(boolean includeNetworkInPayload) {
        this.includeNetworkInPayload = includeNetworkInPayload;
    }
}
