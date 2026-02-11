package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.User;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for UserService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class UserIntegrationTest {

    private ProtectClient client;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        client = TestHelper.getTestClient(1);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
    }

    @Test
    void getMe() throws ApiException {
        User me = client.getUserService().getMe();

        System.out.println("Current user:");
        System.out.println("  ID: " + me.getId());
        System.out.println("  Email: " + me.getEmail());
        System.out.println("  Name: " + me.getFirstName() + " " + me.getLastName());

        assertNotNull(me);
        assertNotNull(me.getId());
    }

    @Test
    void listUsers() throws ApiException {
        List<User> users = client.getUserService().getUsers(10, 0);

        System.out.println("Found " + users.size() + " users");
        for (User u : users) {
            System.out.println("  " + u.getId() + ": " + u.getFirstName() + " " + u.getLastName() + " <" + u.getEmail() + ">");
        }

        assertNotNull(users);
    }
}
