package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a set of transaction approval rules for a specific action type.
 * <p>
 * Transaction rules define the approval requirements for different types of operations
 * (e.g., transfers, staking, contract calls) on specific blockchains. Rules are organized
 * as a matrix with columns (conditions) and lines (threshold configurations).
 * <p>
 * The key identifies the rule set (e.g., blockchain/action type combination), and the
 * lines define different approval thresholds based on source wallet or other conditions.
 *
 * @see RuleColumn
 * @see RuleLine
 * @see TransactionRuleDetails
 * @see DecodedRulesContainer
 */
public class TransactionRules {

    /**
     * Unique key identifying this rule set (e.g., blockchain/action type).
     */
    private String key;

    /**
     * Column definitions specifying the conditions evaluated in each rule line.
     */
    private List<RuleColumn> columns;

    /**
     * Rule lines, each containing conditions and approval thresholds.
     */
    private List<RuleLine> lines;

    /**
     * Additional details about the transaction rules (e.g., feature flags).
     */
    private TransactionRuleDetails details;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the unique key for this rule set.
     *
     * @return the key
     */
    public String getKey() {
        return key;
    }

    /**
     * Sets the unique key for this rule set.
     *
     * @param key the key
     */
    public void setKey(String key) {
        this.key = key;
    }

    /**
     * Gets the columns definition.
     *
     * @return the columns
     */
    public List<RuleColumn> getColumns() {
        return columns;
    }

    /**
     * Sets the columns definition.
     *
     * @param columns the columns
     */
    public void setColumns(List<RuleColumn> columns) {
        this.columns = columns;
    }

    /**
     * Gets the rule lines.
     *
     * @return the lines
     */
    public List<RuleLine> getLines() {
        return lines;
    }

    /**
     * Sets the rule lines.
     *
     * @param lines the lines
     */
    public void setLines(List<RuleLine> lines) {
        this.lines = lines;
    }

    /**
     * Gets the transaction rule details.
     *
     * @return the details
     */
    public TransactionRuleDetails getDetails() {
        return details;
    }

    /**
     * Sets the transaction rule details.
     *
     * @param details the details
     */
    public void setDetails(TransactionRuleDetails details) {
        this.details = details;
    }
}
