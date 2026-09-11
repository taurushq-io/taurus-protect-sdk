package com.taurushq.sdk.protect.client.crypto;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.helper.AssetHashHelper;
import com.taurushq.sdk.protect.client.helper.WhitelistHashHelper;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Hex;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

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

            List<String> legacyHashes = WhitelistHashHelper.computeLegacyHashes(payload);

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

            List<String> legacyHashes = AssetHashHelper.computeAssetLegacyHashes(payload);

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

    // No local re-implementation of the legacy-hash regexes lives here any more.
    //
    // It used to: two private methods mirroring the (then private) production ones with
    // String.replaceAll. That made this file — the CROSS-SDK legacy-hash oracle — assert
    // the behaviour of a copy rather than of the SDK, so the injection fix could land in
    // production while this gate stayed green against the old semantics. The production
    // functions were moved to WhitelistHashHelper / AssetHashHelper and made public
    // precisely so this test can call them.

    @Test
    @DisplayName("the TPV1 canonical string matches every cross-SDK vector")
    void canonicalStringMatchesEveryVector() {
        // Nothing pinned the canonical MESSAGE before 2026-09-10, and that is how a real
        // interop break shipped: the `hmac_sha256` group HMACs a hardcoded string that
        // merely LOOKS like a canonical message, and no consumer routed through
        // calculateSignedHeader — so Python and TypeScript upper-cased the HTTP method
        // while Java and Go signed it verbatim, leaving a caller who issued a lowercase
        // `get` unable to authenticate against one of the two families. This section was
        // consumed by the Go suite alone until now.
        JsonObject group = vectors.getAsJsonObject("canonical_string");
        byte[] secret;
        try {
            secret = Hex.decodeHex(group.get("secret_hex").getAsString().toCharArray());
        } catch (org.apache.commons.codec.DecoderException ex) {
            throw new AssertionError("canonical_string.secret_hex is not hex", ex);
        }
        JsonArray cases = group.getAsJsonArray("cases");

        // A case added and consumed by nobody must fail loudly.
        assertEquals(group.get("count").getAsInt(), cases.size(),
                "canonical_string: case count disagrees with the file's own declaration");

        for (JsonElement e : cases) {
            JsonObject vec = e.getAsJsonObject();
            String description = vec.get("description").getAsString();

            // Asserting the SIGNATURE is what pins the MESSAGE: the secret is fixed, so a
            // match means the exact byte string was signed. Rebuilding the message here
            // would assert a copy of the implementation instead — the trap this file
            // already carries a note about for the legacy-hash regexes.
            String header = CryptoTPV1.calculateSignedHeader(
                    vec.get("api_key").getAsString(),
                    secret,
                    vec.get("nonce").getAsString(),
                    vec.get("timestamp").getAsLong(),
                    vec.get("method").getAsString(),
                    vec.get("host").getAsString(),
                    vec.get("path").getAsString(),
                    vec.get("query").getAsString(),
                    vec.get("content_type").getAsString(),
                    vec.get("body").getAsString());

            assertTrue(
                    header.contains("Signature=" + vec.get("expected_signature").getAsString()),
                    description + ": expected a signature over "
                            + vec.get("expected_message").getAsString() + ", got " + header);
        }
    }
}
