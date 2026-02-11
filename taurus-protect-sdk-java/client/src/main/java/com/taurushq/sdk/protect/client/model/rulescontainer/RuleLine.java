package com.taurushq.sdk.protect.client.model.rulescontainer;

import com.google.protobuf.ByteString;
import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a rule line in a transaction rule table.
 * <p>
 * Each line contains condition cells (inputs) that are matched against transaction
 * attributes, and parallel thresholds (outputs) that define the approval requirements
 * when the line matches. Lines with higher priority are evaluated first.
 *
 * @see RuleColumn
 * @see SequentialThresholds
 * @see TransactionRules
 */
public class RuleLine {

    /**
     * Condition cells containing encoded values to match against the corresponding columns.
     */
    private List<ByteString> cells;

    /**
     * Approval thresholds to apply when this line matches.
     * Multiple entries represent parallel approval paths (any path can satisfy the rule).
     */
    private List<SequentialThresholds> parallelThresholds;

    /**
     * Priority level for rule evaluation (higher values = higher priority).
     */
    private int priority;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the cells (rule inputs).
     *
     * @return the cells
     */
    public List<ByteString> getCells() {
        return cells;
    }

    /**
     * Sets the cells (rule inputs).
     *
     * @param cells the cells
     */
    public void setCells(List<ByteString> cells) {
        this.cells = cells;
    }

    /**
     * Gets the parallel thresholds (rule outputs).
     *
     * @return the parallel thresholds
     */
    public List<SequentialThresholds> getParallelThresholds() {
        return parallelThresholds;
    }

    /**
     * Sets the parallel thresholds (rule outputs).
     *
     * @param parallelThresholds the parallel thresholds
     */
    public void setParallelThresholds(List<SequentialThresholds> parallelThresholds) {
        this.parallelThresholds = parallelThresholds;
    }

    /**
     * Gets the priority.
     *
     * @return the priority (0 is lowest)
     */
    public int getPriority() {
        return priority;
    }

    /**
     * Sets the priority.
     *
     * @param priority the priority
     */
    public void setPriority(int priority) {
        this.priority = priority;
    }
}
