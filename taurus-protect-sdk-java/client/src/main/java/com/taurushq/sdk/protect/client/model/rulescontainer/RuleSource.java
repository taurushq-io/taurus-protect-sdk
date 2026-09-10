package com.taurushq.sdk.protect.client.model.rulescontainer;

import com.google.protobuf.ByteString;
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

    /**
     * Internal address restriction details (populated when type is RuleSourceInternalAddress).
     */
    private RuleSourceInternalAddress internalAddress;

    /**
     * Exchange restriction details (populated when type is RuleSourceExchange).
     */
    private RuleSourceExchange exchange;

    /**
     * External address restriction details (populated when type is RuleSourceExternalAddress).
     */
    private RuleSourceExternalAddress externalAddress;

    /**
     * Verbatim cell bytes, retained when the source type is not recognized by this SDK
     * version so it can be re-encoded losslessly. Null for recognized types.
     */
    private ByteString raw;

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

    /**
     * Gets the internal address restriction.
     *
     * @return the internal address restriction, or null if not applicable
     */
    public RuleSourceInternalAddress getInternalAddress() {
        return internalAddress;
    }

    /**
     * Sets the internal address restriction.
     *
     * @param internalAddress the internal address restriction
     */
    public void setInternalAddress(RuleSourceInternalAddress internalAddress) {
        this.internalAddress = internalAddress;
    }

    /**
     * Gets the exchange restriction.
     *
     * @return the exchange restriction, or null if not applicable
     */
    public RuleSourceExchange getExchange() {
        return exchange;
    }

    /**
     * Sets the exchange restriction.
     *
     * @param exchange the exchange restriction
     */
    public void setExchange(RuleSourceExchange exchange) {
        this.exchange = exchange;
    }

    /**
     * Gets the external address restriction.
     *
     * @return the external address restriction, or null if not applicable
     */
    public RuleSourceExternalAddress getExternalAddress() {
        return externalAddress;
    }

    /**
     * Sets the external address restriction.
     *
     * @param externalAddress the external address restriction
     */
    public void setExternalAddress(RuleSourceExternalAddress externalAddress) {
        this.externalAddress = externalAddress;
    }

    /**
     * Gets the verbatim cell bytes retained for an unrecognized source type.
     *
     * @return the raw cell bytes, or null for recognized types
     */
    public ByteString getRaw() {
        return raw;
    }

    /**
     * Sets the verbatim cell bytes for an unrecognized source type.
     *
     * @param raw the raw cell bytes
     */
    public void setRaw(ByteString raw) {
        this.raw = raw;
    }
}
