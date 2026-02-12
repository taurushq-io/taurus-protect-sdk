package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collections;

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
    void getUsers_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                userService.getUsers(0, 0));
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
