package com.taurushq.sdk.protect.client.testutil;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/**
 * Loader for the cross-SDK vector files in {@code scripts/resources/}, which every SDK's
 * unit suite consumes. A missing file is a hard failure, never a skip: a skipped gate
 * reports parity that is not there.
 */
public final class SharedVectors {

    private SharedVectors() {
    }

    /**
     * Loads a shared vector file and checks its {@code counts} block.
     *
     * @param fileName the file name under {@code scripts/resources/}
     * @param expected alternating section names and counts, e.g. {@code "offset", 26}; each
     *                 must equal both the file's {@code counts} entry and, when the section
     *                 is an array or object, its real size
     * @return the parsed document
     */
    public static JsonObject load(final String fileName, final Object... expected) {
        String[] candidates = {
                "../../scripts/resources/" + fileName,  // from client/
                "../scripts/resources/" + fileName,     // from sdk-java/
                "scripts/resources/" + fileName,        // from the repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                try {
                    JsonObject doc = JsonParser.parseString(
                            new String(Files.readAllBytes(p), StandardCharsets.UTF_8)).getAsJsonObject();
                    checkCounts(fileName, doc, expected);
                    return doc;
                } catch (IOException e) {
                    throw new UncheckedIOException(e);
                }
            }
        }
        throw new IllegalStateException("Cannot find " + fileName + " under scripts/resources/. "
                + "Working directory: " + System.getProperty("user.dir"));
    }

    private static void checkCounts(final String fileName, final JsonObject doc, final Object... expected) {
        JsonObject counts = doc.getAsJsonObject("counts");
        if (counts == null) {
            throw new IllegalStateException(fileName + " has no counts block");
        }
        if (counts.size() * 2 != expected.length) {
            throw new IllegalStateException(fileName + " counts " + counts.keySet()
                    + " do not match the sections this loader consumes");
        }
        for (int i = 0; i < expected.length; i += 2) {
            String section = (String) expected[i];
            int want = (Integer) expected[i + 1];
            if (!counts.has(section) || counts.get(section).getAsInt() != want) {
                throw new IllegalStateException(fileName + " counts." + section + " is "
                        + counts.get(section) + ", this loader expects " + want
                        + " — update every SDK's loader in lockstep");
            }
            if (doc.has(section)) {
                int actual = doc.get(section).isJsonArray() ? doc.getAsJsonArray(section).size()
                        : doc.getAsJsonObject(section).size();
                if (actual != want) {
                    throw new IllegalStateException(fileName + " section " + section + " has "
                            + actual + " entries but counts says " + want);
                }
            }
        }
    }
}
