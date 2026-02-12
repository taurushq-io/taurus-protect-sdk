package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a source restriction in governance rules.
 * <p>
 * Rule sources define conditions based on where a transaction originates.
 * The type determines what kind of source restriction applies, and the
 * corresponding payload provides the specific restriction details.
 *
 * @see RuleSourceType
 * @see RuleSourceInternalWallet
 * @see AddressWhitelistingLine
 */
public class RuleSource {

    /**
     * The type of source restriction.
     */
    private RuleSourceType type;

    /**
     * Internal wallet restriction details (populated when type is RuleSourceInternalWallet).
     */
    private RuleSourceInternalWallet internalWallet;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the source type.
     *
     * @return the source type
     */
    public RuleSourceType getType() {
        return type;
    }

    /**
     * Sets the source type.
     *
     * @param type the source type
     */
    public void setType(RuleSourceType type) {
        this.type = type;
    }

    /**
     * Gets the internal wallet restriction.
     * Only populated when type is {@link RuleSourceType#RuleSourceInternalWallet}.
     *
     * @return the internal wallet restriction, or null if not applicable
     */
    public RuleSourceInternalWallet getInternalWallet() {
        return internalWallet;
    }

    /**
     * Sets the internal wallet restriction.
     *
     * @param internalWallet the internal wallet restriction
     */
    public void setInternalWallet(RuleSourceInternalWallet internalWallet) {
        this.internalWallet = internalWallet;
    }
}
