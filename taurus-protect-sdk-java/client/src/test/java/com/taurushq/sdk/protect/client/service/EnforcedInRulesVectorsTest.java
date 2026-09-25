package com.taurushq.sdk.protect.client.service;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.fail;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.client.model.GroupUser;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.model.UserGroup;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestFactory;

/**
 * Cross-SDK enforced-in-rules flags: getMe asks validatord to compute them, an endpoint that
 * computes them reads an absent enforcedInRules as false, and any other absent flag stays null.
 * The shared vector file is consumed by the Go/Java/Python/TypeScript suites alike.
 */
class EnforcedInRulesVectorsTest {

    private static final int VECTOR_COUNT = 7;

    private static JsonArray vectors;

    @BeforeAll
    static void loadVectors() throws IOException {
        String[] candidates = {
                "../../scripts/resources/enforced-in-rules-vectors.json", // from client/
                "../scripts/resources/enforced-in-rules-vectors.json",    // from sdk-java/
                "scripts/resources/enforced-in-rules-vectors.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                vectors = JsonParser.parseString(json).getAsJsonArray();
                return;
            }
        }
        throw new IOException("Cannot find enforced-in-rules-vectors.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    @Test
    void sharedFileHoldsEveryVector() {
        assertEquals(VECTOR_COUNT, vectors.size(), "enforced-in-rules vector count");
    }

    @TestFactory
    List<DynamicTest> everyVectorRunsThroughTheRealService() {
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : vectors) {
            JsonObject vector = element.getAsJsonObject();
            tests.add(DynamicTest.dynamicTest(vector.get("name").getAsString(), () -> run(vector)));
        }
        return tests;
    }

    private static void run(final JsonObject vector) throws Exception {
        String name = vector.get("name").getAsString();
        JsonObject expected = vector.getAsJsonObject("expected");
        StubTransport stub = StubTransport.replying(vector.get("reply").toString());
        ApiExceptionMapper errors = new ApiExceptionMapper();

        String operation = vector.get("operation").getAsString();
        switch (operation) {
            case "getMe":
                assertUsers(name, expected, Collections.singletonList(
                        new UserService(stub.client(), errors).getMe()));
                break;
            case "getUser":
                String id = vector.getAsJsonObject("request").get("id").getAsString();
                assertUsers(name, expected, Collections.singletonList(
                        new UserService(stub.client(), errors).getUser(id)));
                break;
            case "listUsers":
                assertUsers(name, expected, new UserService(stub.client(), errors).getUsers(0, 0).getUsers());
                break;
            case "listGroups":
                assertGroups(name, expected, new GroupService(stub.client(), errors).getGroups().getGroups());
                break;
            default:
                fail(name + ": no runner registered for operation " + operation);
        }

        if (expected.has("query")) {
            for (Map.Entry<String, JsonElement> param : expected.getAsJsonObject("query").entrySet()) {
                assertEquals(param.getValue().getAsString(), stub.only().param(param.getKey()),
                        name + ": query parameter " + param.getKey());
            }
        }
    }

    private static void assertUsers(final String name, final JsonObject expected, final List<User> users) {
        JsonArray want = expected.getAsJsonArray("users");
        assertEquals(want.size(), users.size(), name + ": user count");
        for (int i = 0; i < want.size(); i++) {
            JsonObject w = want.get(i).getAsJsonObject();
            User user = users.get(i);
            assertEquals(flag(w.get("enforcedInRules")), user.getEnforcedInRules(),
                    name + ": user " + i + " enforcedInRules");
            assertEquals(flag(w.get("publicKeyEnforcedInRules")), user.getPublicKeyEnforcedInRules(),
                    name + ": user " + i + " publicKeyEnforcedInRules");

            JsonArray wantGroups = w.getAsJsonArray("groupsEnforcedInRules");
            List<UserGroup> groups = user.getGroups();
            assertEquals(wantGroups.size(), groups.size(), name + ": user " + i + " membership count");
            for (int g = 0; g < wantGroups.size(); g++) {
                assertEquals(flag(wantGroups.get(g)), groups.get(g).getEnforcedInRules(),
                        name + ": user " + i + " membership " + g + " enforcedInRules");
            }
        }
    }

    private static void assertGroups(final String name, final JsonObject expected, final List<Group> groups) {
        JsonArray want = expected.getAsJsonArray("groups");
        assertEquals(want.size(), groups.size(), name + ": group count");
        for (int i = 0; i < want.size(); i++) {
            JsonObject w = want.get(i).getAsJsonObject();
            Group group = groups.get(i);
            assertEquals(flag(w.get("enforcedInRules")), group.getEnforcedInRules(),
                    name + ": group " + i + " enforcedInRules");

            JsonArray wantUsers = w.getAsJsonArray("usersEnforcedInRules");
            List<GroupUser> users = group.getUsers();
            assertEquals(wantUsers.size(), users.size(), name + ": group " + i + " user count");
            for (int u = 0; u < wantUsers.size(); u++) {
                assertEquals(flag(wantUsers.get(u)), users.get(u).getEnforcedInRules(),
                        name + ": group " + i + " user " + u + " enforcedInRules");
            }
        }
    }

    // A JSON null is Java null: the flag was not computed.
    private static Boolean flag(final JsonElement value) {
        return value == null || value.isJsonNull() ? null : value.getAsBoolean();
    }
}
