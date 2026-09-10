package com.taurushq.sdk.protect.client.model;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

/**
 * Cross-SDK role-extraction alignment.
 *
 * <p>Every SDK must derive the same roles from the same server message, so the shared
 * vector file is consumed by the Go/Java/Python/TypeScript suites alike and a wording
 * change is caught in one place for all four.
 */
class AuthorizationErrorVectorsTest {

    private static JsonArray vectors;

    @BeforeAll
    static void loadVectors() throws IOException {
        String[] candidates = {
                "../../scripts/resources/authorization-error-vectors.json", // from client/
                "../scripts/resources/authorization-error-vectors.json",    // from sdk-java/
                "scripts/resources/authorization-error-vectors.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                vectors = JsonParser.parseString(json).getAsJsonArray();
                return;
            }
        }
        throw new IOException("Cannot find authorization-error-vectors.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    @Test
    void everyVectorParsesToTheExpectedRoles() {
        assertTrue(vectors.size() > 0, "shared vectors file is empty");

        for (JsonElement element : vectors) {
            JsonObject vector = element.getAsJsonObject();
            String description = vector.get("description").getAsString();
            String message = vector.get("message").getAsString();

            List<String> expected = new ArrayList<>();
            for (JsonElement role : vector.get("expected_roles").getAsJsonArray()) {
                expected.add(role.getAsString());
            }

            assertEquals(expected, AuthorizationException.parseRequiredRoles(message), description);
        }
    }

    @Test
    void constructorPopulatesRequiredRoles() {
        AuthorizationException e = new AuthorizationException(
                "one of the 'admin - adminreadonly' role is required", "PermissionDenied", "");

        assertEquals(Arrays.asList("admin", "adminreadonly"), e.getRequiredRoles());
    }

    @Test
    void nonRoleDenialHasNoRequiredRoles() {
        AuthorizationException e = new AuthorizationException("This endpoint has been disabled");

        assertTrue(e.getRequiredRoles().isEmpty());
    }

    @Test
    void defaultConstructorHasNoRequiredRoles() {
        assertTrue(new AuthorizationException().getRequiredRoles().isEmpty());
    }

    @Test
    void nullMessageIsTolerated() {
        assertTrue(AuthorizationException.parseRequiredRoles(null).isEmpty());
    }

    @Test
    void requiredRolesCannotBeMutatedByCallers() {
        AuthorizationException e = new AuthorizationException("one of the 'admin' role is required");

        assertThrows(UnsupportedOperationException.class, () -> e.getRequiredRoles().add("superadmin"));
    }
}
