package com.taurushq.sdk.protect.client.testutil;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Loader for the shared {@code scripts/resources/verification-signed-fixtures.json}.
 * <p>
 * Public and in testutil because the two threshold sections are exercised from different
 * packages: the SuperAdmin one against {@code helper.SignatureVerifier}, the per-group one
 * against a package-private service method. Duplicating the loader is how the two would
 * drift on the count assertion, which is the file's own non-vacuity guard.
 */
public final class SignedFixtures {

    private SignedFixtures() {
    }

    /**
     * Reads and validates the fixture file.
     *
     * @return the parsed document
     * @throws IOException if the file cannot be found from any of the working directories
     *                     Maven and IDEs run tests from
     */
    public static JsonObject load() throws IOException {
        String[] candidates = {
                "../../scripts/resources/verification-signed-fixtures.json", // from client/
                "../scripts/resources/verification-signed-fixtures.json",    // from sdk-java/
                "scripts/resources/verification-signed-fixtures.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                JsonObject fixtures = JsonParser.parseString(json).getAsJsonObject();
                assertCounts(fixtures);
                return fixtures;
            }
        }
        throw new IOException("Cannot find verification-signed-fixtures.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    /**
     * A section added to the file but consumed by nobody must fail loudly.
     *
     * @param fixtures the parsed document
     */
    private static void assertCounts(final JsonObject fixtures) {
        JsonObject counts = fixtures.getAsJsonObject("counts");
        assertTrue(counts.size() > 0, "fixture declares no section counts");
        for (String key : counts.keySet()) {
            assertEquals(counts.get(key).getAsInt(), fixtures.getAsJsonArray(key).size(),
                    key + ": vector count disagrees with the count the file declares");
        }
    }
}
