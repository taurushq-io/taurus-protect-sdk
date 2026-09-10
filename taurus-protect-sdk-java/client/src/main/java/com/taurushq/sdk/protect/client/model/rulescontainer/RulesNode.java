package com.taurushq.sdk.protect.client.model.rulescontainer;

import com.google.protobuf.ByteString;

/**
 * Base for every governance-rules model node backed by a protobuf message.
 *
 * <p>Carries the node's unknown protobuf fields verbatim. When a rules container
 * encoded by a newer schema version is decoded, fields this SDK version does not
 * know about are captured here and re-attached on encode, so a decode/encode
 * round-trip never silently drops data. This mirrors the Go, Python and
 * TypeScript SDKs.
 */
public abstract class RulesNode {

    private ByteString unknownFields;

    /**
     * Gets the verbatim unknown protobuf fields captured at decode.
     *
     * @return the unknown-field bytes, or null/empty when none were present
     */
    public ByteString getUnknownFields() {
        return unknownFields;
    }

    /**
     * Sets the verbatim unknown protobuf fields to re-attach on encode.
     *
     * @param unknownFields the unknown-field bytes
     */
    public void setUnknownFields(ByteString unknownFields) {
        this.unknownFields = unknownFields;
    }

    /**
     * Reports whether this node carries protobuf fields unknown to this SDK version.
     *
     * @return true if unknown fields were preserved
     */
    public boolean hasUnknownFields() {
        return unknownFields != null && !unknownFields.isEmpty();
    }
}
