package com.taurushq.sdk.protect.client.helper;

/**
 * One backward-compatible rewrite of a signed payload: the exact byte string a
 * pre-schema-change signer covered, together with its hash.
 *
 * <p>Step 4 must carry the PAYLOAD forward, not just the hash. The legacy strips are not
 * injective, so a response-controlling server can append a member the strip removes — a
 * duplicate {@code ,"label":"X"} immediately before the closing brace — to a genuinely
 * signed payload. The residue is then the signed bytes exactly, every signature check
 * passes, and a step 6 that parsed the DELIVERED payload would return the appended value
 * as verified (Gson's {@code JsonObject} keeps the LAST of two duplicate keys). Parsing
 * the MATCHED VARIANT is what makes step 6's contract true: every field came from bytes a
 * counted signature covered.
 *
 * <p>The strips are global, and that does NOT bound the exposure the way it first appears.
 * A row whose DELIVERED payload carries inner {@code linkedInternalAddresses} labels is
 * still exposed, because validatord rebuilds those labels on every read from live DB
 * relations rather than from the signed envelope — so for a row signed before per-object
 * labels existed, removing every label (the inner ones the server added AND the one the
 * attacker appended) lands exactly on the signed bytes. Both injectable members reach the
 * caller on any legacy row: {@code label} at either level, and {@code contractType}.
 *
 * <p>What does bound it is the regex alphabet: {@code [^"]*} cannot contain a quote, so
 * nothing beyond those two string values can be smuggled in.
 *
 * <p>At most one variant can hash-match, and by preimage resistance a matching variant
 * <em>is</em> the byte string governance signed — so this is a sound fix, not a narrowing.
 *
 * @see WhitelistHashHelper#computeLegacyPayloadVariants(String)
 * @see AssetHashHelper#computeAssetLegacyPayloadVariants(String)
 */
public final class LegacyPayloadVariant {

    private final String hash;
    private final String payload;

    /**
     * Constructs a variant.
     *
     * @param hash    the hex SHA-256 of {@code payload}
     * @param payload the rewritten payload whose hash a signer may have covered
     */
    public LegacyPayloadVariant(final String hash, final String payload) {
        this.hash = hash;
        this.payload = payload;
    }

    /**
     * Returns the hex SHA-256 of {@link #getPayload()}.
     *
     * @return the hash
     */
    public String getHash() {
        return hash;
    }

    /**
     * Returns the payload this variant's hash covers.
     *
     * <p>This, not the delivered payload, is what step 6 must parse when this variant is
     * the match. The delivered payload is deliberately left untouched elsewhere: a caller
     * needs it to reproduce {@code metadata.hash}. Read this when the question is "what
     * was actually signed".
     *
     * @return the payload
     */
    public String getPayload() {
        return payload;
    }

    @Override
    public String toString() {
        // Deliberately does not print the payload: it is attacker-influenced until step 5
        // completes, and a log line is not the place to render it.
        return "LegacyPayloadVariant{hash=" + hash + ", payloadLength="
                + (payload == null ? 0 : payload.length()) + "}";
    }
}
