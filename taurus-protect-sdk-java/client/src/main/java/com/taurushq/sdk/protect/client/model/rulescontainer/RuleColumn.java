package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a column definition in a transaction rule table.
 * <p>
 * Columns define the conditions that can be evaluated in rule lines.
 * Common column types include:
 * <ul>
 *   <li><b>RuleSource</b> - Source wallet/address restriction</li>
 *   <li><b>RuleDestination</b> - Destination address restriction</li>
 *   <li><b>RuleFiatAmount</b> - Transaction amount threshold in fiat currency</li>
 *   <li><b>RuleMetadata</b> - Custom metadata-based conditions</li>
 * </ul>
 *
 * @see RuleLine
 * @see TransactionRules
 */
public class RuleColumn {

    /**
     * The column type (e.g., "RuleSource", "RuleDestination", "RuleFiatAmount").
     */
    private String type;

    /**
     * Human-readable name for the column.
     */
    private String name;

    /**
     * Key for metadata-based columns (used when type is "RuleMetadata").
     */
    private String metadataKey;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the column type.
     *
     * @return the type (e.g., RuleSource, RuleDestination, RuleFiatAmount)
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the column type.
     *
     * @param type the type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets the column name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the column name.
     *
     * @param name the name
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the metadata key.
     *
     * @return the metadata key
     */
    public String getMetadataKey() {
        return metadataKey;
    }

    /**
     * Sets the metadata key.
     *
     * @param metadataKey the metadata key
     */
    public void setMetadataKey(String metadataKey) {
        this.metadataKey = metadataKey;
    }
}
