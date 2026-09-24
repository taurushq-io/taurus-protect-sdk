package com.taurushq.sdk.protect.client.model;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.JsonPrimitive;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.service.AddressService;
import com.taurushq.sdk.protect.client.service.AssetService;
import com.taurushq.sdk.protect.client.service.GovernanceRuleService;
import com.taurushq.sdk.protect.client.service.MultiFactorSignatureService;
import com.taurushq.sdk.protect.client.service.WhitelistedAddressService;
import com.taurushq.sdk.protect.client.testutil.SignedWhitelist;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.JSON;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBalancesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMultiFactorSignaturesEntityType;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccountType;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTokenType;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

/**
 * Cross-SDK decode tolerance.
 *
 * <p>A field or an enum value the client was not generated with must never fail a decode, and
 * the four SDKs must treat it the same way: unknown fields land in the generated model's
 * additional properties and are written back, unknown enum values keep their raw value all the
 * way into the domain model, and caller-supplied values are sent verbatim. The shared vector
 * file is consumed by the Go/Java/Python/TypeScript suites alike.
 */
class DecodeToleranceVectorsTest {

    private static final Map<String, Class<?>> MODELS = new HashMap<>();
    private static final Map<String, Class<?>> ENUMS = new HashMap<>();

    static {
        MODELS.put("tgvalidatordGetBalancesReply", TgvalidatordGetBalancesReply.class);
        MODELS.put("tgvalidatordStakeAccount", TgvalidatordStakeAccount.class);
        MODELS.put("tgvalidatordGetMultiFactorSignatureEntitiesInfoReply",
                TgvalidatordGetMultiFactorSignatureEntitiesInfoReply.class);
        ENUMS.put("tgvalidatordTokenType", TgvalidatordTokenType.class);
        ENUMS.put("tgvalidatordMultiFactorSignaturesEntityType", TgvalidatordMultiFactorSignaturesEntityType.class);
        ENUMS.put("tgvalidatordStakeAccountType", TgvalidatordStakeAccountType.class);
    }

    private static JsonArray vectors;

    @BeforeAll
    static void loadVectors() throws IOException {
        String[] candidates = {
                "../../scripts/resources/decode-tolerance-vectors.json", // from client/
                "../scripts/resources/decode-tolerance-vectors.json",    // from sdk-java/
                "scripts/resources/decode-tolerance-vectors.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                vectors = JsonParser.parseString(json).getAsJsonArray();
                return;
            }
        }
        throw new IOException("Cannot find decode-tolerance-vectors.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    private static List<JsonObject> ofKind(final String kind) {
        List<JsonObject> matching = new ArrayList<>();
        for (JsonElement element : vectors) {
            JsonObject vector = element.getAsJsonObject();
            if (kind.equals(vector.get("kind").getAsString())) {
                matching.add(vector);
            }
        }
        assertFalse(matching.isEmpty(), "no '" + kind + "' vectors");
        return matching;
    }

    @Test
    void everyVectorIsOfAKindThisSuiteRuns() {
        assertTrue(vectors.size() > 0, "shared vectors file is empty");
        for (JsonElement element : vectors) {
            String kind = element.getAsJsonObject().get("kind").getAsString();
            assertTrue(kind.equals("model") || kind.equals("enum") || kind.equals("service"),
                    "unknown vector kind " + kind);
        }
    }

    @Test
    void generatedModelsKeepUnknownFieldsAndRoundTripLosslessly() throws Exception {
        for (JsonObject vector : ofKind("model")) {
            String name = vector.get("name").getAsString();
            Class<?> type = MODELS.get(vector.get("model").getAsString());
            assertNotNull(type, name + ": no generated class mapped for " + vector.get("model"));
            JsonObject expected = vector.getAsJsonObject("expected");
            String wire = vector.get("wire").toString();

            Object decoded;
            try {
                decoded = JSON.getGson().fromJson(wire, type);
            } catch (RuntimeException e) {
                throw new AssertionError(name + ": decode failed: " + e.getMessage(), e);
            }

            Object extras = type.getMethod("getAdditionalProperties").invoke(decoded);
            assertEquals(expected.get("rootAdditionalProperties"), JSON.getGson().toJsonTree(extras),
                    name + ": root additional properties");

            if (expected.get("lossless").getAsBoolean()) {
                assertEquals(JsonParser.parseString(wire), JsonParser.parseString(JSON.serialize(decoded)),
                        name + ": re-serialization is not lossless");
            }
        }
    }

    @Test
    void generatedEnumsKeepUnknownValuesRaw() throws Exception {
        for (JsonObject vector : ofKind("enum")) {
            String name = vector.get("name").getAsString();
            Class<?> type = ENUMS.get(vector.get("enum").getAsString());
            assertNotNull(type, name + ": no generated enum mapped for " + vector.get("enum"));
            JsonObject expected = vector.getAsJsonObject("expected");

            Object decoded;
            try {
                decoded = JSON.getGson().fromJson(new JsonPrimitive(vector.get("wire").getAsString()), type);
            } catch (RuntimeException e) {
                throw new AssertionError(name + ": decode failed: " + e.getMessage(), e);
            }

            assertEquals(expected.get("value").getAsString(), type.getMethod("getValue").invoke(decoded), name);
            boolean known = false;
            for (Object constant : (Object[]) type.getMethod("values").invoke(null)) {
                known |= constant == decoded;
            }
            assertEquals(expected.get("known").getAsBoolean(), known, name + ": known");
        }
    }

    @Test
    void servicesPassUnknownValuesThroughBothWays() throws Exception {
        for (JsonObject vector : ofKind("service")) {
            String name = vector.get("name").getAsString();
            String operation = vector.get("operation").getAsString();
            JsonObject request = vector.getAsJsonObject("request");
            JsonObject expected = vector.getAsJsonObject("expected");
            String reply = vector.get("reply").toString();
            switch (operation) {
                case "getMultiFactorSignatureInfo":
                    getMultiFactorSignatureInfo(name, request, expected, reply);
                    break;
                case "createMultiFactorSignatures":
                    createMultiFactorSignatures(name, request, expected, reply);
                    break;
                case "queryAssetAddresses":
                    queryAssetAddresses(name, request, expected, reply);
                    break;
                default:
                    fail(name + ": unhandled operation " + operation);
            }
        }
    }

    private static void getMultiFactorSignatureInfo(final String name, final JsonObject request,
                                                    final JsonObject expected, final String reply)
            throws Exception {
        StubTransport stub = StubTransport.replying(reply);
        MultiFactorSignatureService service = new MultiFactorSignatureService(stub.client(), new ApiExceptionMapper());

        MultiFactorSignatureInfo info = service.getMultiFactorSignatureInfo(request.get("id").getAsString());

        assertEquals(expected.get("entityType").getAsString(), info.getEntityType().getKind(), name);
    }

    private static void createMultiFactorSignatures(final String name, final JsonObject request,
                                                    final JsonObject expected, final String reply)
            throws Exception {
        StubTransport stub = StubTransport.replying(reply);
        MultiFactorSignatureService service = new MultiFactorSignatureService(stub.client(), new ApiExceptionMapper());
        List<String> ids = new ArrayList<>();
        for (JsonElement id : request.getAsJsonArray("entityIds")) {
            ids.add(id.getAsString());
        }

        service.createMultiFactorSignatures(ids,
                TgvalidatordMultiFactorSignaturesEntityType.fromValue(request.get("entityType").getAsString()));

        assertSentFields(name, expected.getAsJsonObject("requestBody"), stub.only().body());
    }

    private static void queryAssetAddresses(final String name, final JsonObject request,
                                            final JsonObject expected, final String reply)
            throws Exception {
        String assetId = request.get("assetId").getAsString();
        String path = "/api/rest/v2/assets/" + assetId + "/addresses/query";
        StubTransport stub = StubTransport.replying().route(path, r -> reply);
        ApiClient client = stub.client();
        ApiExceptionMapper mapper = new ApiExceptionMapper();
        List<PublicKey> superAdminKeys = new SignedWhitelist().superAdminKeys();
        RulesContainerCache cache = new RulesContainerCache(
                new GovernanceRuleService(client, mapper, superAdminKeys, 1));
        AssetService assets = new AssetService(client, mapper, cache, new AddressService(client, mapper, cache),
                new WhitelistedAddressService(client, mapper, superAdminKeys, 1));

        AssetAddressV2Result result = assets.queryAssetAddresses(assetId,
                request.get("addressType").getAsString(), null, null, null);

        assertEquals(1, stub.requests().size(), name + ": only the holders page may be read");
        assertSentFields(name, expected.getAsJsonObject("requestBody"), stub.requestsTo(path).get(0).body());
        JsonArray rows = expected.getAsJsonArray("rows");
        assertEquals(rows.size(), result.getAddresses().size(), name + ": rows returned");
        for (int i = 0; i < rows.size(); i++) {
            JsonObject row = rows.get(i).getAsJsonObject();
            AssetAddressV2 got = result.getAddresses().get(i);
            assertEquals(row.get("address").getAsString(), got.getAddress(), name);
            assertEquals(row.get("addressType").getAsString(), got.getAddressType(), name);
            assertEquals(row.get("verified").getAsBoolean(), got.isVerified(), name);
        }
        assertEquals(expected.get("excludedCount").getAsInt(), result.getExcludedUnverified().size(),
                name + ": excluded");
    }

    private static void assertSentFields(final String name, final JsonObject expected, final String body) {
        assertNotNull(body, name + ": the request had no body");
        JsonObject sent = JsonParser.parseString(body).getAsJsonObject();
        for (Map.Entry<String, JsonElement> field : expected.entrySet()) {
            assertEquals(field.getValue(), sent.get(field.getKey()), name + ": sent " + field.getKey());
        }
    }
}
