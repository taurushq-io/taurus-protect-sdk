package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class UserServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private UserService userService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        userService = new UserService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new UserService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new UserService(apiClient, null));
    }

    @Test
    void getUsers_zeroLimitSendsTheDefault() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new UserService(stub.client(), new ApiExceptionMapper()).getUsers(0, 0);
        assertEquals("20", stub.only().param("limit"));
    }

    @Test
    void getUser_throwsOnEmptyIdWithoutARequest() {
        StubTransport stub = StubTransport.replying("{}");
        UserService service = new UserService(stub.client(), new ApiExceptionMapper());
        assertThrows(IllegalArgumentException.class, () -> service.getUser(""));
        assertThrows(IllegalArgumentException.class, () -> service.getUser(null));
        assertEquals(0, stub.requests().size());
    }

    @Test
    void getUser_readsTheUserById() throws Exception {
        StubTransport stub = StubTransport.replying("{\"result\":{\"id\":\"8\"}}");
        User user = new UserService(stub.client(), new ApiExceptionMapper()).getUser("8");
        assertEquals("8", user.getId());
        assertEquals("/api/rest/v1/users/8", stub.only().path());
    }

    @Test
    void getUsers_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                userService.getUsers(10, -1));
    }

    @Test
    void getUsersByEmail_throwsOnNullEmails() {
        assertThrows(NullPointerException.class, () ->
                userService.getUsersByEmail(null));
    }

    @Test
    void getUsersByEmail_throwsOnEmptyEmails() {
        assertThrows(IllegalArgumentException.class, () ->
                userService.getUsersByEmail(Collections.emptyList()));
    }

    @Test
    void createUserAttribute_throwsOnZeroUserId() {
        assertThrows(IllegalArgumentException.class, () ->
                userService.createUserAttribute(0, "key", "value"));
    }

    @Test
    void createUserAttribute_throwsOnEmptyKey() {
        assertThrows(IllegalArgumentException.class, () ->
                userService.createUserAttribute(1, "", "value"));
    }

    @Test
    void createUserAttribute_throwsOnNullValue() {
        assertThrows(NullPointerException.class, () ->
                userService.createUserAttribute(1, "key", null));
    }
}
