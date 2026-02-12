package com.taurushq.sdk.protect.client.crypto;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Hex;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

import static org.bouncycastle.util.Strings.constantTimeAreEqual;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Cross-SDK cryptographic test vectors.
 *
 * These tests verify that Java SDK cryptographic functions produce
 * identical results to the Go, Python, and TypeScript SDKs. All SDKs
 * read the same test vectors from docs/test-vectors/crypto-test-vectors.json.
 */
class CrossSdkCryptoVectorTest {

    private static JsonObject vectors;

    @BeforeAll
    static void loadVectors() throws IOException {
        // Try multiple paths to find the vectors file
        Path vectorsPath = findVectorsFile();
        String json = new String(Files.readAllBytes(vectorsPath), StandardCharsets.UTF_8);
        vectors = JsonParser.parseString(json).getAsJsonObject().getAsJsonObject("vectors");
    }

    private static Path findVectorsFile() {
        // Try relative to working directory (varies by build tool)
        String[] candidates = {
                "../../docs/test-vectors/crypto-test-vectors.json",  // From client/
                "../docs/test-vectors/crypto-test-vectors.json",     // From sdk-java/
                "docs/test-vectors/crypto-test-vectors.json",        // From repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                return p;
            }
        }
        throw new RuntimeException("Cannot find crypto-test-vectors.json. "
                + "Working directory: " + System.getProperty("user.dir"));
    }

    // ============ SHA-256 Hex Hash Tests ============

    @Test
    void testHexHashVectors() {
        JsonArray hexHash = vectors.getAsJsonArray("hex_hash");
        for (JsonElement elem : hexHash) {
            JsonObject vec = elem.getAsJsonObject();
            String input = vec.get("input").getAsString();
            String expected = vec.get("expected").getAsString();
            String description = vec.get("description").getAsString();

            String result = CryptoTPV1.calculateHexHash(input);
            assertEquals(expected, result,
                    "SHA-256 mismatch for: " + description);
        }
    }

    // ============ HMAC-SHA256 Tests ============

    @Test
    void testHmacSha256Vectors() throws Exception {
        JsonArray hmacVectors = vectors.getAsJsonArray("hmac_sha256");
        for (JsonElement elem : hmacVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String keyHex = vec.get("key_hex").getAsString();
            String data = vec.get("data").getAsString();
            String expectedBase64 = vec.get("expected_base64").getAsString();
            String description = vec.get("description").getAsString();

            byte[] key = Hex.decodeHex(keyHex.toCharArray());
            String result = CryptoTPV1.calculateBase64Hmac(key, data);
            assertEquals(expectedBase64, result,
                    "HMAC-SHA256 mismatch for: " + description);
        }
    }

    // ============ Constant-Time Compare Tests ============

    @Test
    void testConstantTimeCompareVectors() {
        JsonArray ctcVectors = vectors.getAsJsonArray("constant_time_compare");
        for (JsonElement elem : ctcVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String a = vec.get("a").getAsString();
            String b = vec.get("b").getAsString();
            boolean expected = vec.get("expected").getAsBoolean();
            String description = vec.get("description").getAsString();

            boolean result = constantTimeAreEqual(a, b);
            assertEquals(expected, result,
                    "Constant-time compare mismatch for: " + description);
        }
    }

    // ============ Legacy Address Hash Tests ============

    @Test
    void testLegacyAddressHashOriginals() {
        JsonArray addressVectors = vectors.getAsJsonArray("legacy_hash_address");
        for (JsonElement elem : addressVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String payload = vec.get("payload").getAsString();
            String expectedOriginal = vec.get("original_hash").getAsString();
            String description = vec.get("description").getAsString();

            String result = CryptoTPV1.calculateHexHash(payload);
            assertEquals(expectedOriginal, result,
                    "Original hash mismatch for: " + description);
        }
    }

    @Test
    void testLegacyAddressHashStrategies() {
        JsonArray addressVectors = vectors.getAsJsonArray("legacy_hash_address");
        for (JsonElement elem : addressVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String payload = vec.get("payload").getAsString();
            int expectedCount = vec.get("expected_legacy_count").getAsInt();
            String description = vec.get("description").getAsString();

            List<String> legacyHashes = computeAddressLegacyHashes(payload);

            assertEquals(expectedCount, legacyHashes.size(),
                    "Legacy hash count mismatch for: " + description);

            if (expectedCount > 0) {
                String expectedWithoutCT = vec.get("expected_without_contract_type").getAsString();
                String expectedWithoutLabels = vec.get("expected_without_labels").getAsString();
                String expectedWithoutBoth = vec.get("expected_without_both").getAsString();

                assertTrue(legacyHashes.contains(expectedWithoutCT),
                        "Missing without_contract_type for: " + description);
                assertTrue(legacyHashes.contains(expectedWithoutLabels),
                        "Missing without_labels for: " + description);
                assertTrue(legacyHashes.contains(expectedWithoutBoth),
                        "Missing without_both for: " + description);
            }
        }
    }

    // ============ Legacy Asset Hash Tests ============

    @Test
    void testLegacyAssetHashOriginals() {
        JsonArray assetVectors = vectors.getAsJsonArray("legacy_hash_asset");
        for (JsonElement elem : assetVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String payload = vec.get("payload").getAsString();
            String expectedOriginal = vec.get("original_hash").getAsString();
            String description = vec.get("description").getAsString();

            String result = CryptoTPV1.calculateHexHash(payload);
            assertEquals(expectedOriginal, result,
                    "Original hash mismatch for: " + description);
        }
    }

    @Test
    void testLegacyAssetHashStrategies() {
        JsonArray assetVectors = vectors.getAsJsonArray("legacy_hash_asset");
        for (JsonElement elem : assetVectors) {
            JsonObject vec = elem.getAsJsonObject();
            String payload = vec.get("payload").getAsString();
            int expectedCount = vec.get("expected_legacy_count").getAsInt();
            String description = vec.get("description").getAsString();

            List<String> legacyHashes = computeAssetLegacyHashes(payload);

            assertEquals(expectedCount, legacyHashes.size(),
                    "Legacy hash count mismatch for: " + description);

            if (expectedCount > 0) {
                String expectedWithoutNFT = vec.get("expected_without_is_nft").getAsString();
                String expectedWithoutKT = vec.get("expected_without_kind_type").getAsString();
                String expectedWithoutBoth = vec.get("expected_without_both").getAsString();

                assertTrue(legacyHashes.contains(expectedWithoutNFT),
                        "Missing without_is_nft for: " + description);
                assertTrue(legacyHashes.contains(expectedWithoutKT),
                        "Missing without_kind_type for: " + description);
                assertTrue(legacyHashes.contains(expectedWithoutBoth),
                        "Missing without_both for: " + description);
            }
        }
    }

    // ============ Helper Methods ============
    // These mirror the private methods in WhitelistedAddressService and WhitelistedAssetService

    private static List<String> computeAddressLegacyHashes(String payloadAsString) {
        Set<String> uniqueHashes = new LinkedHashSet<>();

        // Strategy 1: Remove contractType only
        String withoutContractType = payloadAsString.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        if (!withoutContractType.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutContractType));
        }

        // Strategy 2: Remove labels from linkedInternalAddresses objects only
        String withoutLabels = payloadAsString.replaceAll(",\"label\":\"[^\"]*\"}", "}");
        if (!withoutLabels.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutLabels));
        }

        // Strategy 3: Remove both
        String withoutBoth = payloadAsString.replaceAll(",\"label\":\"[^\"]*\"}", "}");
        withoutBoth = withoutBoth.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        if (!withoutBoth.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutBoth));
        }

        return new ArrayList<>(uniqueHashes);
    }

    private static List<String> computeAssetLegacyHashes(String payloadAsString) {
        Set<String> uniqueHashes = new LinkedHashSet<>();

        // Strategy 1: Remove isNFT only
        String withoutIsNFT = payloadAsString.replaceAll(",\"isNFT\":(true|false)", "");
        withoutIsNFT = withoutIsNFT.replaceAll("\"isNFT\":(true|false),", "");
        if (!withoutIsNFT.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutIsNFT));
        }

        // Strategy 2: Remove kindType only
        String withoutKindType = payloadAsString.replaceAll(",\"kindType\":\"[^\"]*\"", "");
        withoutKindType = withoutKindType.replaceAll("\"kindType\":\"[^\"]*\",", "");
        if (!withoutKindType.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutKindType));
        }

        // Strategy 3: Remove both (isNFT first, then kindType — matches Java SDK order)
        String withoutBoth = payloadAsString.replaceAll(",\"isNFT\":(true|false)", "");
        withoutBoth = withoutBoth.replaceAll("\"isNFT\":(true|false),", "");
        withoutBoth = withoutBoth.replaceAll(",\"kindType\":\"[^\"]*\"", "");
        withoutBoth = withoutBoth.replaceAll("\"kindType\":\"[^\"]*\",", "");
        if (!withoutBoth.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutBoth));
        }

        return new ArrayList<>(uniqueHashes);
    }
}
