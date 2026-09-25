package com.taurushq.sdk.protect.client.testutil;

import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.protobuf.ByteString;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.proto.v1.RequestReply;

import java.nio.charset.StandardCharsets;
import java.security.GeneralSecurityException;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Whitelisted-address envelopes that pass the SDK's six-step verification for real: fresh
 * keys, a SuperAdmin-signed rules container, and rows signed by the one approver its
 * threshold requires. A tampered row carries a payload its hash does not cover, so it fails
 * step 1 and is excluded.
 * <pre>{@code
 * SignedWhitelist whitelist = new SignedWhitelist();
 * new WhitelistedAddressService(client, mapper, whitelist.superAdminKeys(), 1);
 * stub.route("/api/rest/v1/whitelists/addresses", r -> SignedWhitelist.page(rows, r.values("ids")));
 * }</pre>
 */
public final class SignedWhitelist {

    private static final Gson GSON = new Gson();

    private final KeyPair superAdmin;
    private final KeyPair team;
    private final String container;
    private final String signatures;

    /**
     * Creates the keys and the signed rules container.
     *
     * @throws Exception if key generation or signing fails
     */
    public SignedWhitelist() throws Exception {
        superAdmin = newKeyPair();
        team = newKeyPair();
        RequestReply.RulesContainer rules = RequestReply.RulesContainer.newBuilder()
                .addUsers(RequestReply.User.newBuilder().setId("team1@bank.com")
                        .setPublicKey(CryptoTPV1.encodePublicKey(team.getPublic())))
                .addGroups(RequestReply.Group.newBuilder().setId("team1").addUserIds("team1@bank.com"))
                .addAddressWhitelistingRules(RequestReply.RulesContainer.AddressWhitelistingRules.newBuilder()
                        .setCurrency("ALGO")
                        .setNetwork("mainnet")
                        .addParallelThresholds(RequestReply.SequentialThresholds.newBuilder()
                                .addThresholds(RequestReply.GroupThreshold.newBuilder()
                                        .setGroupId("team1").setMinimumSignatures(1))))
                .build();
        byte[] raw = rules.toByteArray();
        byte[] signature = Base64.getDecoder().decode(
                CryptoTPV1.calculateBase64Signature(superAdmin.getPrivate(), raw));
        container = Base64.getEncoder().encodeToString(raw);
        signatures = Base64.getEncoder().encodeToString(RequestReply.UserSignatures.newBuilder()
                .addSignatures(RequestReply.UserSignature.newBuilder()
                        .setUserId("superadmin@bank.com").setSignature(ByteString.copyFrom(signature)))
                .build().toByteArray());
    }

    /**
     * Generates a P-256 key pair.
     *
     * @return the key pair
     * @throws GeneralSecurityException if the curve is unavailable
     */
    public static KeyPair newKeyPair() throws GeneralSecurityException {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        return generator.generateKeyPair();
    }

    /**
     * Returns the SuperAdmin key pair that signed the rules container.
     *
     * @return the key pair
     */
    public KeyPair superAdmin() {
        return superAdmin;
    }

    /**
     * Returns the SuperAdmin public keys a verifying service must be configured with.
     *
     * @return the keys
     */
    public List<PublicKey> superAdminKeys() {
        return Collections.singletonList(superAdmin.getPublic());
    }

    /**
     * Returns the metadata hash of row {@code id}'s untampered payload.
     *
     * @param id the row id
     * @return the hex hash
     */
    public static String hashOf(final String id) {
        return CryptoTPV1.calculateHexHash(payload(id));
    }

    private static String payload(final String id) {
        return "{\"currency\":\"ALGO\",\"network\":\"mainnet\",\"addressType\":\"individual\","
                + "\"address\":\"ADDR" + id + "\",\"memo\":\"\",\"label\":\"row " + id + "\",\"customerId\":\"\"}";
    }

    /**
     * Builds a signed envelope whose address is ADDR{id}, carrying its own rules container.
     *
     * @param id       the row id
     * @param tampered true for a payload the signed hash does not cover
     * @return the envelope JSON
     * @throws Exception if signing fails
     */
    public JsonObject row(final String id, final boolean tampered) throws Exception {
        String payload = payload(id);
        String hash = CryptoTPV1.calculateHexHash(payload);
        List<String> hashes = Collections.singletonList(hash);
        String signature = CryptoTPV1.calculateBase64Signature(team.getPrivate(),
                GSON.toJson(hashes).getBytes(StandardCharsets.UTF_8));
        if (tampered) {
            payload = payload.replace("ADDR", "EVIL");
        }

        JsonObject userSignature = new JsonObject();
        userSignature.addProperty("userId", "team1@bank.com");
        userSignature.addProperty("signature", signature);
        JsonObject entry = new JsonObject();
        entry.add("signature", userSignature);
        entry.add("hashes", GSON.toJsonTree(hashes));
        JsonArray entries = new JsonArray();
        entries.add(entry);
        JsonObject signedAddress = new JsonObject();
        signedAddress.addProperty("payload",
                Base64.getEncoder().encodeToString(payload.getBytes(StandardCharsets.UTF_8)));
        signedAddress.add("signatures", entries);
        JsonObject metadata = new JsonObject();
        metadata.addProperty("hash", hash);
        metadata.addProperty("payloadAsString", payload);

        JsonObject row = new JsonObject();
        row.addProperty("id", id);
        row.addProperty("blockchain", "ALGO");
        row.addProperty("network", "mainnet");
        row.add("metadata", metadata);
        row.addProperty("rulesContainer", container);
        row.addProperty("rulesSignatures", signatures);
        row.add("signedAddress", signedAddress);
        return row;
    }

    /**
     * Answers an id-filtered list request: the rows it asks for that exist, in request order.
     *
     * @param rows the rows by id
     * @param ids  the ids the request asks for
     * @return the reply body
     */
    public static String page(final Map<String, JsonObject> rows, final List<String> ids) {
        List<JsonObject> found = new ArrayList<>();
        for (String id : ids) {
            if (rows.containsKey(id)) {
                found.add(rows.get(id));
            }
        }
        return page(found, String.valueOf(found.size()));
    }

    /**
     * Builds a list reply.
     *
     * @param rows       the rows
     * @param totalItems the server total
     * @return the reply body
     */
    public static String page(final List<JsonObject> rows, final String totalItems) {
        JsonArray result = new JsonArray();
        rows.forEach(result::add);
        JsonObject reply = new JsonObject();
        reply.add("result", result);
        reply.addProperty("totalItems", totalItems);
        return reply.toString();
    }
}
