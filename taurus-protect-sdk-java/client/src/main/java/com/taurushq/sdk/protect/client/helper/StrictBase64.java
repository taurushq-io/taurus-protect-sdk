package com.taurushq.sdk.protect.client.helper;

import java.nio.charset.StandardCharsets;
import java.util.Base64;

/**
 * Strict base64 decoding for security-critical paths.
 *
 * <p>Lenient base64 decoders — commons-codec's {@code Base64.decodeBase64} among them —
 * silently DISCARD characters outside the base64 alphabet and absorb the rest. That makes
 * the decoding non-injective: two different strings can decode to the same bytes, and,
 * worse, a string carrying embedded separators decodes to a <em>valid container with
 * attacker-chosen bytes appended</em>. Protobuf treats concatenation as merge, so appended
 * bytes ADD entries to repeated fields such as {@code users} — which is how an
 * attacker-controlled HSM or PRICEUPDATER key would reach the governance trust root.
 *
 * <p>Every governance path that turns a base64 string into bytes must agree byte-for-byte
 * on what that string means: the verification memo key, the container decode, and the
 * signature verification that decides what a signature actually covers. Go uses strict
 * {@code base64.StdEncoding} throughout, and that is precisely why Go was not exploitable
 * by this route.
 *
 * <p>Throws {@link IllegalArgumentException} (what {@code Base64.Decoder} already throws)
 * and lets each caller wrap it in the appropriate SDK exception type.
 */
public final class StrictBase64 {

    private StrictBase64() {
        // Utility class.
    }

    /**
     * Decodes base64, rejecting any character outside the base64 alphabet.
     *
     * <p>ASCII whitespace is stripped before decoding: line-wrapped base64 is legitimate
     * and introduces no ambiguity about the decoded bytes, whereas every other
     * out-of-alphabet byte is rejected rather than silently dropped.
     *
     * @param data base64 text; null and the empty string decode to an empty array
     * @return the decoded bytes
     * @throws IllegalArgumentException if the input is not well-formed base64
     */
    public static byte[] decode(final String data) {
        if (data == null || data.isEmpty()) {
            return new byte[0];
        }
        String compact = data.replaceAll("[ \t\n\r\f]+", "");
        return Base64.getDecoder().decode(compact.getBytes(StandardCharsets.US_ASCII));
    }
}
