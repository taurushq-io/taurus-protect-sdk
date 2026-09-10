package com.taurushq.sdk.protect.client.model.rulescontainer;

import java.util.Map;

import com.google.protobuf.ByteString;

/**
 * Base for governance-rules model nodes that also carry a free-form
 * {@code properties} map in the protobuf schema.
 *
 * @see RulesNode
 */
public abstract class RulesNodeWithProperties extends RulesNode {

    private Map<String, ByteString> properties;

    /**
     * Gets the free-form properties map.
     *
     * @return the properties, or null when none were set
     */
    public Map<String, ByteString> getProperties() {
        return properties;
    }

    /**
     * Sets the free-form properties map.
     *
     * @param properties the properties
     */
    public void setProperties(Map<String, ByteString> properties) {
        this.properties = properties;
    }
}
