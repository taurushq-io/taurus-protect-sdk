package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a source-specific rule line in address whitelisting rules.
 * <p>
 * Each line contains source restrictions (cells) that define when this line applies,
 * and the approval thresholds to use when the source matches. Lines allow different
 * approval requirements for whitelisting addresses from different source wallets.
 *
 * @see AddressWhitelistingRules
 * @see RuleSource
 * @see SequentialThresholds
 */
public class AddressWhitelistingLine {

    /**
     * Source restriction cells defining when this line applies.
     * Typically contains a single RuleSource with the source wallet restriction.
     */
    private List<RuleSource> cells;

    /**
     * Approval thresholds to use when this line's source restriction matches.
     * Multiple entries represent parallel approval paths.
     */
    private List<SequentialThresholds> parallelThresholds;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the source restriction cells.
     * For address whitelisting, typically contains a single RuleSource defining the source restriction.
     *
     * @return the cells
     */
    public List<RuleSource> getCells() {
        return cells;
    }

    /**
     * Sets the source restriction cells.
     *
     * @param cells the cells
     */
    public void setCells(List<RuleSource> cells) {
        this.cells = cells;
    }

    /**
     * Gets the parallel thresholds to use when this line matches.
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
