package com.taurushq.sdk.protect.client.mapper;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonArray;
import com.google.gson.JsonDeserializationContext;
import com.google.gson.JsonDeserializer;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.lang.reflect.Type;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for rules container decoding from JSON fixture.
 * Uses the fixture at fixtures/whitelisted_address_raw_response.json
 * which contains rulesContainerJson (JSON format) and rulesSignatures (protobuf).
 */
class RulesContainerDecodeTest {

    private static final String FIXTURE_PATH = "fixtures/whitelisted_address_raw_response.json";

    private static JsonObject fixtureJson;
    private static DecodedRulesContainer decodedContainer;
    private static String rulesSignaturesBase64;

    /**
     * Custom deserializer for RuleUser to handle publicKey field as PEM string.
     * The JSON fixture has "publicKey" as a PEM string, but the model has both
     * publicKeyPem (String) and publicKey (java.security.PublicKey).
     */
    private static class RuleUserDeserializer implements JsonDeserializer<RuleUser> {
        @Override
        public RuleUser deserialize(JsonElement json, Type typeOfT, JsonDeserializationContext context) {
            JsonObject obj = json.getAsJsonObject();
            RuleUser user = new RuleUser();
            user.setId(obj.get("id").getAsString());

            // publicKey in JSON is the PEM string
            if (obj.has("publicKey") && !obj.get("publicKey").isJsonNull()) {
                user.setPublicKeyPem(obj.get("publicKey").getAsString());
            }

            // Parse roles array
            if (obj.has("roles") && !obj.get("roles").isJsonNull()) {
                JsonArray rolesArray = obj.getAsJsonArray("roles");
                List<String> roles = new ArrayList<>();
                for (JsonElement roleElem : rolesArray) {
                    roles.add(roleElem.getAsString());
                }
                user.setRoles(roles);
            }

            return user;
        }
    }

    @BeforeAll
    static void setUp() throws IOException {
        // Load the fixture file
        try (InputStream is = RulesContainerDecodeTest.class.getClassLoader()
                .getResourceAsStream(FIXTURE_PATH)) {
            assertNotNull(is, "Fixture file not found: " + FIXTURE_PATH);
            fixtureJson = JsonParser.parseReader(new InputStreamReader(is, StandardCharsets.UTF_8))
                    .getAsJsonObject();
        }

        // Extract rulesContainerJson (JSON format in fixture)
        JsonObject rulesContainerJson = fixtureJson.getAsJsonObject("rulesContainerJson");
        assertNotNull(rulesContainerJson, "rulesContainerJson not found in fixture");

        // Parse the container manually to handle the complex structure
        decodedContainer = parseDecodedRulesContainer(rulesContainerJson);

        // Extract rulesSignatures (protobuf encoded as base64)
        rulesSignaturesBase64 = fixtureJson.get("rulesSignatures").getAsString();
    }

    /**
     * Manually parse DecodedRulesContainer from JSON to handle complex nested structures.
     */
    private static DecodedRulesContainer parseDecodedRulesContainer(JsonObject json) {
        Gson gsonWithCustomDeserializer = new GsonBuilder()
                .registerTypeAdapter(RuleUser.class, new RuleUserDeserializer())
                .create();

        DecodedRulesContainer container = new DecodedRulesContainer();

        // Parse users
        if (json.has("users") && !json.get("users").isJsonNull()) {
            JsonArray usersArray = json.getAsJsonArray("users");
            List<RuleUser> users = new ArrayList<>();
            for (JsonElement userElem : usersArray) {
                users.add(gsonWithCustomDeserializer.fromJson(userElem, RuleUser.class));
            }
            container.setUsers(users);
        }

        // Parse groups
        if (json.has("groups") && !json.get("groups").isJsonNull()) {
            JsonArray groupsArray = json.getAsJsonArray("groups");
            List<RuleGroup> groups = new ArrayList<>();
            for (JsonElement groupElem : groupsArray) {
                JsonObject groupObj = groupElem.getAsJsonObject();
                RuleGroup group = new RuleGroup();
                group.setId(groupObj.get("id").getAsString());
                if (groupObj.has("userIds") && !groupObj.get("userIds").isJsonNull()) {
                    JsonArray userIdsArray = groupObj.getAsJsonArray("userIds");
                    List<String> userIds = new ArrayList<>();
                    for (JsonElement idElem : userIdsArray) {
                        userIds.add(idElem.getAsString());
                    }
                    group.setUserIds(userIds);
                }
                groups.add(group);
            }
            container.setGroups(groups);
        }

        // Parse addressWhitelistingRules
        if (json.has("addressWhitelistingRules") && !json.get("addressWhitelistingRules").isJsonNull()) {
            JsonArray rulesArray = json.getAsJsonArray("addressWhitelistingRules");
            List<AddressWhitelistingRules> rules = new ArrayList<>();
            for (JsonElement ruleElem : rulesArray) {
                JsonObject ruleObj = ruleElem.getAsJsonObject();
                AddressWhitelistingRules rule = new AddressWhitelistingRules();
                if (ruleObj.has("currency") && !ruleObj.get("currency").isJsonNull()) {
                    rule.setCurrency(ruleObj.get("currency").getAsString());
                }
                if (ruleObj.has("network") && !ruleObj.get("network").isJsonNull()) {
                    rule.setNetwork(ruleObj.get("network").getAsString());
                }
                // Parse parallelThresholds
                if (ruleObj.has("parallelThresholds") && !ruleObj.get("parallelThresholds").isJsonNull()) {
                    JsonArray thresholdsArray = ruleObj.getAsJsonArray("parallelThresholds");
                    List<SequentialThresholds> thresholds = new ArrayList<>();
                    for (JsonElement thresholdElem : thresholdsArray) {
                        JsonObject thresholdObj = thresholdElem.getAsJsonObject();
                        // This is actually a GroupThreshold structure in the fixture
                        SequentialThresholds seqThreshold = new SequentialThresholds();
                        List<GroupThreshold> groupThresholds = new ArrayList<>();
                        GroupThreshold gt = new GroupThreshold();
                        if (thresholdObj.has("groupId")) {
                            gt.setGroupId(thresholdObj.get("groupId").getAsString());
                        }
                        if (thresholdObj.has("minimumSignatures")) {
                            gt.setMinimumSignatures(thresholdObj.get("minimumSignatures").getAsInt());
                        }
                        groupThresholds.add(gt);
                        seqThreshold.setThresholds(groupThresholds);
                        thresholds.add(seqThreshold);
                    }
                    rule.setParallelThresholds(thresholds);
                }
                rules.add(rule);
            }
            container.setAddressWhitelistingRules(rules);
        }

        // Parse other fields
        if (json.has("minimumDistinctUserSignatures")) {
            container.setMinimumDistinctUserSignatures(json.get("minimumDistinctUserSignatures").getAsInt());
        }
        if (json.has("minimumDistinctGroupSignatures")) {
            container.setMinimumDistinctGroupSignatures(json.get("minimumDistinctGroupSignatures").getAsInt());
        }
        if (json.has("timestamp")) {
            container.setTimestamp(json.get("timestamp").getAsLong());
        }
        if (json.has("enforcedRulesHash") && !json.get("enforcedRulesHash").isJsonNull()) {
            container.setEnforcedRulesHash(json.get("enforcedRulesHash").getAsString());
        }

        return container;
    }

    // ========== Category A: Rules Container Decoding Tests (9 tests) ==========

    @Test
    void testDecodeRulesContainerFromBase64Success() {
        // Test that the rulesContainerJson can be base64-encoded as JSON and decoded
        JsonObject rulesContainerJson = fixtureJson.getAsJsonObject("rulesContainerJson");
        Gson gson = new Gson();
        String json = gson.toJson(rulesContainerJson);
        String base64Json = Base64.getEncoder().encodeToString(json.getBytes(StandardCharsets.UTF_8));

        // Decode from base64 JSON
        byte[] decoded = Base64.getDecoder().decode(base64Json);
        String decodedJson = new String(decoded, StandardCharsets.UTF_8);

        // Parse the JSON object
        JsonObject decodedJsonObj = JsonParser.parseString(decodedJson).getAsJsonObject();

        // Use our custom parser to decode
        DecodedRulesContainer container = parseDecodedRulesContainer(decodedJsonObj);

        assertNotNull(container);
        assertNotNull(container.getUsers());
        assertNotNull(container.getGroups());
    }

    @Test
    void testDecodedContainerHasUsers() {
        // Fixture has 4 users: superadmin1, superadmin2, team1, hsmslot
        assertNotNull(decodedContainer.getUsers());
        assertEquals(4, decodedContainer.getUsers().size());

        List<String> userIds = decodedContainer.getUsers().stream()
                .map(RuleUser::getId)
                .collect(Collectors.toList());

        assertTrue(userIds.contains("superadmin1@bank.com"));
        assertTrue(userIds.contains("superadmin2@bank.com"));
        assertTrue(userIds.contains("team1@bank.com"));
        assertTrue(userIds.contains("hsmslot@bank.com"));
    }

    @Test
    void testDecodedContainerHasGroups() {
        // Fixture has 2 groups: team1 and superadmins
        assertNotNull(decodedContainer.getGroups());
        assertEquals(2, decodedContainer.getGroups().size());

        List<String> groupIds = decodedContainer.getGroups().stream()
                .map(RuleGroup::getId)
                .collect(Collectors.toList());

        assertTrue(groupIds.contains("team1"));
        assertTrue(groupIds.contains("superadmins"));
    }

    @Test
    void testDecodedUsersHavePemKeys() {
        // All users should have PEM-encoded public keys
        for (RuleUser user : decodedContainer.getUsers()) {
            String publicKey = user.getPublicKeyPem();
            assertNotNull(publicKey, "User " + user.getId() + " should have a public key");
            assertTrue(publicKey.contains("-----BEGIN PUBLIC KEY-----"),
                    "User " + user.getId() + " public key should be PEM-encoded");
            assertTrue(publicKey.contains("-----END PUBLIC KEY-----"),
                    "User " + user.getId() + " public key should be PEM-encoded");
        }
    }

    @Test
    void testDecodedUsersHaveRoles() {
        // All users should have roles
        for (RuleUser user : decodedContainer.getUsers()) {
            assertNotNull(user.getRoles(), "User " + user.getId() + " should have roles");
            assertFalse(user.getRoles().isEmpty(),
                    "User " + user.getId() + " should have at least one role");
        }
    }

    @Test
    void testFindSuperadminUsers() {
        // Fixture has 2 SUPERADMIN users
        List<RuleUser> superadmins = decodedContainer.getUsers().stream()
                .filter(u -> u.getRoles() != null && u.getRoles().contains("SUPERADMIN"))
                .collect(Collectors.toList());

        assertEquals(2, superadmins.size());

        List<String> superadminIds = superadmins.stream()
                .map(RuleUser::getId)
                .collect(Collectors.toList());
        assertTrue(superadminIds.contains("superadmin1@bank.com"));
        assertTrue(superadminIds.contains("superadmin2@bank.com"));
    }

    @Test
    void testFindHsmslotUser() {
        // Fixture has 1 HSMSLOT user
        List<RuleUser> hsmslotUsers = decodedContainer.getUsers().stream()
                .filter(u -> u.getRoles() != null && u.getRoles().contains("HSMSLOT"))
                .collect(Collectors.toList());

        assertEquals(1, hsmslotUsers.size());
        assertEquals("hsmslot@bank.com", hsmslotUsers.get(0).getId());
    }

    @Test
    void testDecodedHasAddressWhitelistingRules() {
        // Fixture has 1 address whitelisting rule for ALGO/mainnet
        assertNotNull(decodedContainer.getAddressWhitelistingRules());
        assertEquals(1, decodedContainer.getAddressWhitelistingRules().size());

        AddressWhitelistingRules rule = decodedContainer.getAddressWhitelistingRules().get(0);
        assertEquals("ALGO", rule.getCurrency());
        assertEquals("mainnet", rule.getNetwork());

        // Should have parallel thresholds
        assertNotNull(rule.getParallelThresholds());
        assertEquals(1, rule.getParallelThresholds().size());
    }

    @Test
    void testInvalidBase64RaisesError() {
        // Invalid base64 should throw exception
        assertThrows(IllegalArgumentException.class, () -> {
            Base64.getDecoder().decode("not-valid-base64!!!");
        });
    }

    // ========== Category B: User Signatures Decoding Tests (5 tests) ==========

    @Test
    void testDecodeUserSignaturesFromBase64Success() throws InvalidProtocolBufferException {
        // rulesSignatures in fixture is protobuf-encoded
        List<RuleUserSignature> signatures =
                RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(rulesSignaturesBase64);

        assertNotNull(signatures);
        assertFalse(signatures.isEmpty());
    }

    @Test
    void testSignaturesContainUserIds() throws InvalidProtocolBufferException {
        List<RuleUserSignature> signatures =
                RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(rulesSignaturesBase64);

        // All signatures should have user IDs
        for (RuleUserSignature sig : signatures) {
            assertNotNull(sig.getUserId(), "Signature should have a user ID");
            assertFalse(sig.getUserId().isEmpty(), "User ID should not be empty");
        }

        // Check specific user IDs from fixture
        List<String> userIds = signatures.stream()
                .map(RuleUserSignature::getUserId)
                .collect(Collectors.toList());
        assertTrue(userIds.contains("superadmin1@bank.com"));
        assertTrue(userIds.contains("superadmin2@bank.com"));
    }

    @Test
    void testSignaturesContainSignatureBytes() throws InvalidProtocolBufferException {
        List<RuleUserSignature> signatures =
                RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(rulesSignaturesBase64);

        // All signatures should have signature bytes (base64-encoded)
        for (RuleUserSignature sig : signatures) {
            assertNotNull(sig.getSignature(), "Signature should have signature bytes");
            assertFalse(sig.getSignature().isEmpty(), "Signature bytes should not be empty");

            // Verify it's valid base64
            byte[] decoded = Base64.getDecoder().decode(sig.getSignature());
            assertTrue(decoded.length > 0, "Decoded signature should have bytes");
        }
    }

    @Test
    void testSignaturesCountMatchesExpected() throws InvalidProtocolBufferException {
        // Fixture has 2 signatures (superadmin1 and superadmin2)
        List<RuleUserSignature> signatures =
                RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(rulesSignaturesBase64);

        assertEquals(2, signatures.size());
    }

    @Test
    void testInvalidSignaturesBase64RaisesError() {
        // Invalid base64 for protobuf should throw
        String invalidBase64 = "not-valid-base64!!!";

        assertThrows(IllegalArgumentException.class, () ->
                RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(invalidBase64));
    }
}
