package com.taurushq.sdk.protect.client.model.rulescontainer;

import java.util.HashMap;
import java.util.Map;

/**
 * Enum representing the type of rule source.
 */
public enum RuleSourceType {
    /**
     * Any source (no restriction).
     */
    RuleSourceAny(0),
    /**
     * Internal wallet source.
     */
    RuleSourceInternalWallet(1),
    /**
     * Internal address source.
     */
    RuleSourceInternalAddress(2),
    /**
     * Exchange source.
     */
    RuleSourceExchange(4),
    /**
     * External address source.
     */
    RuleSourceExternalAddress(5),
    /**
     * Unrecognized source type.
     */
    UNRECOGNIZED(-1);

    private static final Map<Integer, RuleSourceType> VALUE_MAP = new HashMap<>();

    static {
        for (RuleSourceType type : values()) {
            VALUE_MAP.put(type.value, type);
        }
    }

    private final int value;

    RuleSourceType(int value) {
        this.value = value;
    }

    /**
     * Gets the numeric value of the type.
     *
     * @return the numeric value
     */
    public int getValue() {
        return value;
    }

    /**
     * Gets the RuleSourceType for a numeric value.
     *
     * @param value the numeric value
     * @return the corresponding RuleSourceType, or UNRECOGNIZED if not found
     */
    public static RuleSourceType fromValue(int value) {
        RuleSourceType type = VALUE_MAP.get(value);
        return type != null ? type : UNRECOGNIZED;
    }
}
