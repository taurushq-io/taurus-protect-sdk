package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonToken;
import com.taurushq.sdk.protect.client.model.InternalAddress;
import com.taurushq.sdk.protect.client.model.InternalWallet;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;

import java.io.IOException;
import java.io.StringReader;
import java.nio.charset.StandardCharsets;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Signature;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.regex.Pattern;

/**
 * Helper class for whitelist signature operations and JSON parsing.
 */
public final class WhitelistHashHelper {

    private static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();

    /** Matches {@code ,"contractType":"..."}. */
    private static final Pattern CONTRACT_TYPE_PATTERN =
            Pattern.compile(",\"contractType\":\"[^\"]*\"");

    /**
     * Matches {@code ,"label":"..."} followed by a closing brace, i.e. a label that is the
     * LAST member of its object. That is normally an inner
     * {@code linkedInternalAddresses[]} label rather than the main address label, which is
     * followed by other fields — but "normally" is exactly the gap this class's fix
     * closes, because a server is free to append the main label as the last member too.
     */
    private static final Pattern LABEL_IN_OBJECT_PATTERN =
            Pattern.compile(",\"label\":\"[^\"]*\"}");

    private WhitelistHashHelper() {
        // Prevent instantiation
    }

    /**
     * Computes the backward-compatible payload rewrites for an ADDRESS, each paired with
     * its hash. Handles addresses signed before schema changes by removing certain fields
     * and recomputing the hash.
     *
     * <p>Strategies, in order, and the order matters because the FIRST variant that
     * hash-matches wins:
     * <ol>
     *   <li>remove {@code contractType} (signed before {@code contractType} existed)</li>
     *   <li>remove trailing {@code label} members (signed after {@code contractType} but
     *       before per-object labels)</li>
     *   <li>remove both (signed before either existed)</li>
     * </ol>
     *
     * <p>This lives here rather than privately inside {@code WhitelistedAddressService}
     * for two reasons. It is a hash helper, alongside the rest of this class; and while it
     * was private, three test files re-implemented these regexes inline with
     * {@code String.replaceAll} — so a fix to the production code left the whole
     * cross-SDK legacy-hash gate asserting the OLD behaviour while reporting green.
     *
     * @param payloadAsString the payload exactly as delivered, may be null
     * @return the variants to try, in strategy order, deduplicated by hash; never null
     */
    public static List<LegacyPayloadVariant> computeLegacyPayloadVariants(
            final String payloadAsString) {
        if (payloadAsString == null) {
            return new ArrayList<>();
        }

        // Insertion-ordered and keyed by hash, so a strategy that collapses onto an
        // earlier one is dropped without disturbing the try order.
        Map<String, LegacyPayloadVariant> unique = new LinkedHashMap<>();

        // Strategy 1: remove contractType only.
        String withoutContractType =
                CONTRACT_TYPE_PATTERN.matcher(payloadAsString).replaceAll("");
        addVariant(unique, payloadAsString, withoutContractType);

        // Strategy 2: remove trailing label members only (keep contractType).
        String withoutLabels = LABEL_IN_OBJECT_PATTERN.matcher(payloadAsString).replaceAll("}");
        addVariant(unique, payloadAsString, withoutLabels);

        // Strategy 3: remove BOTH. Labels first, then contractType — the two orders can
        // differ when a label is the member immediately before contractType.
        String withoutBoth = LABEL_IN_OBJECT_PATTERN.matcher(payloadAsString).replaceAll("}");
        withoutBoth = CONTRACT_TYPE_PATTERN.matcher(withoutBoth).replaceAll("");
        addVariant(unique, payloadAsString, withoutBoth);

        return new ArrayList<>(unique.values());
    }

    /**
     * Computes alternative ADDRESS hashes for backward compatibility, discarding the
     * payload each one came from.
     *
     * <p>Verification must use {@link #computeLegacyPayloadVariants(String)} instead:
     * step 6 needs the payload, not just the hash. This projection remains because the
     * cross-SDK vector oracle ({@code docs/test-vectors/crypto-test-vectors.json}) asserts
     * hashes and counts.
     *
     * @param payloadAsString the payload exactly as delivered, may be null
     * @return the legacy hashes to try, in strategy order; never null
     */
    public static List<String> computeLegacyHashes(final String payloadAsString) {
        return projectHashes(computeLegacyPayloadVariants(payloadAsString));
    }

    /**
     * Adds {@code rewritten} as a variant when the rewrite actually changed something and
     * its hash is new.
     *
     * <p>Package-private so {@link AssetHashHelper} shares one implementation: the two
     * flows use DIFFERENT strips but the same dedupe-and-order rule, and duplicating it
     * is how they would drift.
     *
     * @param unique    the accumulator, keyed by hash
     * @param original  the delivered payload
     * @param rewritten the candidate rewrite
     */
    static void addVariant(final Map<String, LegacyPayloadVariant> unique,
                           final String original, final String rewritten) {
        if (rewritten.equals(original)) {
            return;
        }
        String hash = CryptoTPV1.calculateHexHash(rewritten);
        if (!unique.containsKey(hash)) {
            unique.put(hash, new LegacyPayloadVariant(hash, rewritten));
        }
    }

    /**
     * Projects variants onto their hashes.
     *
     * @param variants the variants
     * @return the hashes, in the same order
     */
    static List<String> projectHashes(final List<LegacyPayloadVariant> variants) {
        List<String> hashes = new ArrayList<>(variants.size());
        for (LegacyPayloadVariant variant : variants) {
            hashes.add(variant.getHash());
        }
        return hashes;
    }

    /**
     * Parses a WhitelistedAddress from the verified JSON payload string.
     * The JSON format contains the signed fields (currency, network, address, etc.)
     * that have been cryptographically verified.
     *
     * <p>Callers must pass the payload the matched signature COVERED — see
     * {@link LegacyPayloadVariant} — not whatever the response delivered.
     *
     * @param json the verified JSON payload string
     * @return the WhitelistedAddress populated from JSON fields
     * @throws WhitelistException if parsing fails
     */
    public static WhitelistedAddress parseWhitelistedAddressFromJson(String json) throws WhitelistException {
        if (Strings.isNullOrEmpty(json)) {
            throw new WhitelistException("JSON payload cannot be null or empty");
        }
        // Bounded and duplicate-checked HERE, not left to the caller's ordering: this
        // method is public and is step 6 of two flows plus a model derivation.
        if (json.length() > MAX_PAYLOAD_BYTES) {
            throw new WhitelistException(String.format(
                    "cannot parse whitelisted address: payload exceeds %d bytes",
                    MAX_PAYLOAD_BYTES));
        }
        rejectDuplicateObjectKeys(json);
        try {
            JsonObject obj = JsonParser.parseString(json).getAsJsonObject();

            WhitelistedAddress addr = new WhitelistedAddress();
            addr.setBlockchain(getStringOrNull(obj, "currency"));
            addr.setNetwork(getStringOrNull(obj, "network"));
            addr.setAddress(getStringOrNull(obj, "address"));
            addr.setMemo(getStringOrNull(obj, "memo"));
            addr.setLabel(getStringOrNull(obj, "label"));
            addr.setCustomerId(getStringOrNull(obj, "customerId"));
            addr.setContractType(getStringOrNull(obj, "contractType"));
            addr.setTnParticipantID(getStringOrNull(obj, "tnParticipantID"));

            // Parse addressType enum
            String addressTypeStr = getStringOrNull(obj, "addressType");
            if (!Strings.isNullOrEmpty(addressTypeStr)) {
                addr.setAddressType(WhitelistedAddress.AddressType.valueOf(addressTypeStr));
            }

            // Parse exchangeAccountId (String in JSON → long in model)
            String exchangeIdStr = getStringOrNull(obj, "exchangeAccountId");
            if (!Strings.isNullOrEmpty(exchangeIdStr)) {
                addr.setExchangeAccountId(Long.parseLong(exchangeIdStr));
            }

            // Parse linkedInternalAddresses
            if (obj.has("linkedInternalAddresses") && !obj.get("linkedInternalAddresses").isJsonNull()) {
                addr.setLinkedInternalAddresses(
                        parseInternalAddresses(obj.getAsJsonArray("linkedInternalAddresses")));
            }

            // Parse linkedWallets
            if (obj.has("linkedWallets") && !obj.get("linkedWallets").isJsonNull()) {
                addr.setLinkedWallets(
                        parseInternalWallets(obj.getAsJsonArray("linkedWallets")));
            }

            return addr;
        } catch (Exception e) {
            throw new WhitelistException("Failed to parse WhitelistedAddress from JSON: " + e.getMessage(), e);
        }
    }

    private static String getStringOrNull(JsonObject obj, String key) {
        if (!obj.has(key) || obj.get(key).isJsonNull()) {
            return null;
        }
        String value = obj.get(key).getAsString();
        return value.isEmpty() ? null : value;
    }

    private static List<InternalAddress> parseInternalAddresses(JsonArray arr) {
        List<InternalAddress> result = new ArrayList<>();
        for (JsonElement elem : arr) {
            JsonObject o = elem.getAsJsonObject();
            InternalAddress ia = new InternalAddress();
            ia.setId(o.has("id") && !o.get("id").isJsonNull() ? o.get("id").getAsLong() : 0);
            ia.setAddress(getStringOrNull(o, "address"));
            ia.setLabel(getStringOrNull(o, "label"));
            result.add(ia);
        }
        return result;
    }

    private static List<InternalWallet> parseInternalWallets(JsonArray arr) {
        List<InternalWallet> result = new ArrayList<>();
        for (JsonElement elem : arr) {
            JsonObject o = elem.getAsJsonObject();
            InternalWallet iw = new InternalWallet();
            iw.setId(o.has("id") && !o.get("id").isJsonNull() ? o.get("id").getAsLong() : 0);
            iw.setName(getStringOrNull(o, "name"));
            iw.setPath(getStringOrNull(o, "path"));
            result.add(iw);
        }
        return result;
    }

    /**
     * Signs a list of hashes with an ECDSA private key.
     * Equivalent to Go's SignHashes function.
     *
     * @param hashes     the list of hashes to sign
     * @param privateKey the private key to sign with
     * @return the signature bytes
     * @throws WhitelistException if signing fails
     */
    public static byte[] signHashes(List<String> hashes, PrivateKey privateKey) throws WhitelistException {
        if (hashes == null || privateKey == null) {
            throw new WhitelistException("hashes and privateKey cannot be null");
        }
        try {
            String json = GSON.toJson(hashes);
            byte[] jsonBytes = json.getBytes(StandardCharsets.UTF_8);

            Signature signer = Signature.getInstance("SHA256withPLAIN-ECDSA");
            signer.initSign(privateKey);
            signer.update(jsonBytes);
            return signer.sign();
        } catch (Exception e) {
            throw new WhitelistException("failed to sign hashes", e);
        }
    }

    /**
     * Verifies a signature of a list of hashes.
     * Equivalent to Go's CheckHashesSignature function.
     *
     * @param hashes    the list of hashes that were signed
     * @param signature the signature to verify
     * @param publicKey the public key to verify against
     * @throws WhitelistException if verification fails
     */
    public static void checkHashesSignature(List<String> hashes, byte[] signature,
                                            PublicKey publicKey) throws WhitelistException {
        if (hashes == null || signature == null || publicKey == null) {
            throw new WhitelistException("hashes, signature, and publicKey cannot be null");
        }

        String json = GSON.toJson(hashes);
        byte[] jsonBytes = json.getBytes(StandardCharsets.UTF_8);

        if (!SignatureVerifier.verifySignature(jsonBytes, signature, publicKey)) {
            throw new WhitelistException(String.format(
                    "invalid signature of hashes %s", json));
        }
    }

    /**
     * Maximum size of a signed payload before it is parsed.
     *
     * <p>The payload is hash-checked in step 1 but not AUTHENTICATED until step 5's
     * signatures verify, so anything parsed in between is still attacker-influenced.
     * A generous ceiling that only stops a hostile response consuming memory.
     */
    public static final int MAX_PAYLOAD_BYTES = 1 << 20;

    /**
     * Maximum nesting the duplicate-key pre-pass will walk.
     *
     * <p>The signed payload is three levels at most (object &rarr;
     * {@code linkedInternalAddresses} &rarr; object), so anything deeper is not a payload
     * this SDK can be looking at — and recursing on server-chosen depth is how a decoder
     * becomes a stack-overflow primitive.
     */
    private static final int MAX_SIGNED_PAYLOAD_DEPTH = 32;

    /**
     * Fails when any object in the JSON document carries the same key twice.
     *
     * <p>Gson's {@code JsonObject} keeps the LAST of two duplicate keys and reports no
     * error, which is a verification bypass on this path rather than a curiosity. The
     * legacy-hash tolerance strips a member the parser would still read, so a server can
     * append {@code ,"label":"X"} before the closing brace of a genuinely signed payload:
     * the strip recovers the signed bytes, every signature check passes, and the parse
     * then returns the attacker's value. Parsing the matched variant (see
     * {@link LegacyPayloadVariant}) closes that; this closes the shapes the strip does not
     * reach, and the third parse in {@link #resolveRuleKey(String, String, String)} where
     * a duplicated {@code currency} or {@code network} would re-point rule selection.
     *
     * <p>It has to be a separate structural pass over the token stream: no Gson setting
     * rejects duplicates, and {@code JsonParser} has already collapsed them by the time a
     * {@code JsonObject} exists.
     *
     * <p>Package-private so {@link AssetHashHelper} uses the same implementation.
     *
     * @param json the document to check
     * <p>It throws the CHECKED {@link WhitelistException}, not
     * {@link IntegrityException}: a bad payload in one row is a row-level failure, and
     * that is the class the whitelist row loops catch and exclude on. An unchecked
     * throw here escapes every existing {@code catch (WhitelistException)} and turns
     * one unparseable row into an aborted listing -- the exact inversion the repo
     * records Java getting wrong once already.
     *
     * @throws WhitelistException if a duplicate key, a malformed document or excessive
     *                            nesting is found
     */
    static void rejectDuplicateObjectKeys(final String json) throws WhitelistException {
        try (JsonReader reader = new JsonReader(new StringReader(json))) {
            // The signed payload is strict JSON; lenient mode accepts unquoted names and
            // trailing garbage, which is the wrong posture for material a signature is
            // supposed to pin.
            reader.setLenient(false);
            checkValueForDuplicates(reader, 0);
            // Trailing content after the top-level value would mean the document the
            // signature covers is not the document that was parsed.
            if (reader.peek() != JsonToken.END_DOCUMENT) {
                throw new WhitelistException(
                        "signed payload carries trailing content after the JSON document");
            }
        } catch (IOException e) {
            throw new WhitelistException(
                    "signed payload is not well-formed JSON: " + e.getMessage(), e);
        }
    }

    private static void checkValueForDuplicates(final JsonReader reader, final int depth)
            throws IOException, WhitelistException {
        if (depth > MAX_SIGNED_PAYLOAD_DEPTH) {
            throw new WhitelistException(String.format(
                    "signed payload nests deeper than %d levels", MAX_SIGNED_PAYLOAD_DEPTH));
        }

        JsonToken token = reader.peek();
        switch (token) {
            case BEGIN_OBJECT:
                reader.beginObject();
                Set<String> seen = new HashSet<>();
                while (reader.hasNext()) {
                    String name = reader.nextName();
                    if (!seen.add(name)) {
                        throw new WhitelistException(String.format(
                                "signed payload carries duplicate key \"%s\"; refusing to "
                                        + "choose between two values for one field", name));
                    }
                    checkValueForDuplicates(reader, depth + 1);
                }
                reader.endObject();
                break;
            case BEGIN_ARRAY:
                reader.beginArray();
                while (reader.hasNext()) {
                    checkValueForDuplicates(reader, depth + 1);
                }
                reader.endArray();
                break;
            default:
                // A scalar (or null): nothing to check, but it must still be consumed or
                // the enclosing loop never advances.
                reader.skipValue();
                break;
        }
    }

    /**
     * Returns the (blockchain, network) pair that selects the governance rules, taken
     * from the SIGNED payload rather than the surrounding DTO.
     *
     * <p>The DTO is free-floating: nothing binds it to the signatures, so a response
     * that set blockchain to empty would steer verification to the global-default rule
     * tier, which is broader than the rule the entity belongs to. The payload pair is
     * at least hash-bound to the material step 5 checks.
     *
     * <p>The chain is always in the payload. The NETWORK is not: governance rules carry
     * a per-rule {@code includeNetworkInPayload} flag, and real signed payloads omit
     * {@code network} when it is off. Requiring it would reject correctly-signed
     * addresses, so the DTO network is used only when the payload has none.
     *
     * <p>Two failures, both closed: the payload omitting the chain (never treated as a
     * wildcard, because an empty value matches every tier), and the payload disagreeing
     * with the DTO on a field the payload DOES carry.
     *
     * <p>Addresses name the chain {@code currency}; assets name it {@code blockchain}.
     * Both are accepted so one helper serves both flows.
     *
     * @param payloadAsString the signed payload
     * @param dtoBlockchain   the blockchain the response claims, may be null or empty
     * @param dtoNetwork      the network the response claims, may be null or empty
     * @return a two-element array: blockchain, then network
     * @throws WhitelistException if the pair cannot be resolved or the two disagree
     */
    public static String[] resolveRuleKey(final String payloadAsString,
                                          final String dtoBlockchain,
                                          final String dtoNetwork) throws WhitelistException {
        ResolvedRuleKey key = resolveRuleKeyWithSource(payloadAsString, dtoBlockchain, dtoNetwork);
        return new String[] {key.getBlockchain(), key.getNetwork()};
    }

    /**
     * {@link #resolveRuleKey(String, String, String)} plus whether the NETWORK was signed.
     *
     * <p>Step 5 needs that fact: when the network came from the unsigned DTO, a single
     * rule lookup lets a response-controlling server choose which rule — and therefore
     * which group quorum — judges the row. See {@link ResolvedRuleKey}.
     *
     * @param payloadAsString the signed payload
     * @param dtoBlockchain   the blockchain the response claims, may be null or empty
     * @param dtoNetwork      the network the response claims, may be null or empty
     * @return the resolved key
     * @throws WhitelistException if the pair cannot be resolved or the two disagree
     */
    public static ResolvedRuleKey resolveRuleKeyWithSource(final String payloadAsString,
                                                           final String dtoBlockchain,
                                                           final String dtoNetwork)
            throws WhitelistException {
        if (Strings.isNullOrEmpty(payloadAsString)) {
            throw new WhitelistException("cannot resolve governance rule key: payload is empty");
        }
        if (payloadAsString.length() > MAX_PAYLOAD_BYTES) {
            throw new WhitelistException(String.format(
                    "cannot resolve governance rule key: payload exceeds %d bytes", MAX_PAYLOAD_BYTES));
        }
        // This is the THIRD parse of the same payload in the address flow (step 5's
        // linked-address read and step 6's derivation are the others), and Gson is
        // last-duplicate-wins in all three. A duplicated "currency" or "network" here
        // re-points rule selection itself, which is strictly worse than a wrong label:
        // it chooses the quorum. Same guard, same reason.
        rejectDuplicateObjectKeys(payloadAsString);

        String blockchain;
        String network;
        try {
            JsonObject obj = JsonParser.parseString(payloadAsString).getAsJsonObject();
            blockchain = readString(obj, "blockchain");
            if (Strings.isNullOrEmpty(blockchain)) {
                blockchain = readString(obj, "currency");
            }
            network = readString(obj, "network");
        } catch (RuntimeException e) {
            throw new WhitelistException("cannot resolve governance rule key: " + e.getMessage(), e);
        }

        if (Strings.isNullOrEmpty(blockchain)) {
            throw new WhitelistException("signed payload does not carry a blockchain; "
                    + "refusing to fall back to a wildcard rule");
        }

        // Only reachable when includeNetworkInPayload is off for this rule.
        boolean networkFromPayload = !Strings.isNullOrEmpty(network);
        if (!networkFromPayload) {
            network = dtoNetwork;
        }

        if (!Strings.isNullOrEmpty(dtoBlockchain) && !dtoBlockchain.equalsIgnoreCase(blockchain)) {
            throw new WhitelistException(String.format(
                    "blockchain disagrees between the signed payload (%s) and the response (%s)",
                    blockchain, dtoBlockchain));
        }
        if (networkFromPayload && !Strings.isNullOrEmpty(dtoNetwork)
                && !dtoNetwork.equalsIgnoreCase(network)) {
            throw new WhitelistException(String.format(
                    "network disagrees between the signed payload (%s) and the response (%s)",
                    network, dtoNetwork));
        }

        return new ResolvedRuleKey(blockchain, network, networkFromPayload);
    }

    private static String readString(final JsonObject obj, final String key) {
        if (!obj.has(key) || obj.get(key).isJsonNull()) {
            return null;
        }
        return obj.get(key).getAsString();
    }
}
