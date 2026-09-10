package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.InternalAddress;
import com.taurushq.sdk.protect.client.model.InternalWallet;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;

import java.nio.charset.StandardCharsets;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Signature;
import java.util.ArrayList;
import java.util.List;

/**
 * Helper class for whitelist signature operations and JSON parsing.
 */
public final class WhitelistHashHelper {

    private static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();

    private WhitelistHashHelper() {
        // Prevent instantiation
    }

    /**
     * Parses a WhitelistedAddress from the verified JSON payload string.
     * The JSON format contains the signed fields (currency, network, address, etc.)
     * that have been cryptographically verified.
     *
     * @param json the verified JSON payload string
     * @return the WhitelistedAddress populated from JSON fields
     * @throws WhitelistException if parsing fails
     */
    public static WhitelistedAddress parseWhitelistedAddressFromJson(String json) throws WhitelistException {
        if (Strings.isNullOrEmpty(json)) {
            throw new WhitelistException("JSON payload cannot be null or empty");
        }
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
        if (Strings.isNullOrEmpty(payloadAsString)) {
            throw new WhitelistException("cannot resolve governance rule key: payload is empty");
        }
        if (payloadAsString.length() > MAX_PAYLOAD_BYTES) {
            throw new WhitelistException(String.format(
                    "cannot resolve governance rule key: payload exceeds %d bytes", MAX_PAYLOAD_BYTES));
        }

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

        return new String[] {blockchain, network};
    }

    private static String readString(final JsonObject obj, final String key) {
        if (!obj.has(key) || obj.get(key).isJsonNull()) {
            return null;
        }
        return obj.get(key).getAsString();
    }
}
