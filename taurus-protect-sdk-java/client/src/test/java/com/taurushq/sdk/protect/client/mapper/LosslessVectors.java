package com.taurushq.sdk.protect.client.mapper;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonParser;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

/**
 * Loader for the cross-SDK lossless/parity vectors.
 *
 * <p>The scenarios in governance-lossless-vectors.json cannot be produced by this SDK's
 * encoder — they are deliberately non-canonical or schema-newer wire bytes — so they are
 * not in governance-cell-vectors.json, which the Go SDK generates. They used to live as
 * base64 literals hand-copied into all four SDKs' round-trip suites, where nothing
 * compared the copies and they could drift silently.
 */
final class LosslessVectors {

    /**
     * Asserted so a vector added to the shared file without being consumed here fails
     * loudly instead of being silently ignored by this SDK.
     */
    static final int VECTOR_COUNT = 9;

    private static final String[] CANDIDATES = {
            "../../scripts/resources/governance-lossless-vectors.json", // from client/
            "../scripts/resources/governance-lossless-vectors.json",    // from sdk-java/
            "scripts/resources/governance-lossless-vectors.json",       // from repo root
    };

    private LosslessVectors() {
    }

    static JsonArray load() {
        for (String candidate : CANDIDATES) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                try {
                    String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                    JsonArray vectors = JsonParser.parseString(json).getAsJsonArray();
                    if (vectors.size() != VECTOR_COUNT) {
                        throw new IllegalStateException(
                                "shared lossless vectors file has " + vectors.size()
                                        + " entries, expected " + VECTOR_COUNT
                                        + " — update every SDK's suite in lockstep");
                    }
                    return vectors;
                } catch (IOException e) {
                    throw new UncheckedIOException(e);
                }
            }
        }
        throw new IllegalStateException(
                "Cannot find governance-lossless-vectors.json. Working directory: "
                        + System.getProperty("user.dir"));
    }

    /** All vectors for a scenario, in file order. */
    static List<String> forScenario(String scenario) {
        List<String> out = new ArrayList<>();
        List<String> descriptions = new ArrayList<>();
        for (JsonElement e : load()) {
            if (scenario.equals(e.getAsJsonObject().get("scenario").getAsString())) {
                out.add(e.getAsJsonObject().get("wire_base64").getAsString());
                descriptions.add(e.getAsJsonObject().get("description").getAsString());
            }
        }
        if (out.isEmpty()) {
            throw new IllegalStateException("no shared lossless vector for scenario " + scenario);
        }
        return out;
    }

    /** Descriptions for a scenario, in the same order as {@link #forScenario}. */
    static List<String> descriptionsForScenario(String scenario) {
        List<String> out = new ArrayList<>();
        for (JsonElement e : load()) {
            if (scenario.equals(e.getAsJsonObject().get("scenario").getAsString())) {
                out.add(e.getAsJsonObject().get("description").getAsString());
            }
        }
        return out;
    }

    /** The single vector for a scenario. */
    static String one(String scenario) {
        List<String> all = forScenario(scenario);
        if (all.size() != 1) {
            throw new IllegalStateException(
                    "scenario " + scenario + " has " + all.size() + " vectors, expected 1");
        }
        return all.get(0);
    }
}
