package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class TagServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private TagService tagService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        tagService = new TagService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new TagService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new TagService(apiClient, null));
    }

    @Test
    void createTag_throwsOnNullValue() {
        assertThrows(IllegalArgumentException.class, () ->
                tagService.createTag(null, "#FF0000"));
    }

    @Test
    void createTag_throwsOnEmptyValue() {
        assertThrows(IllegalArgumentException.class, () ->
                tagService.createTag("", "#FF0000"));
    }

    @Test
    void deleteTag_throwsOnNullTagId() {
        assertThrows(IllegalArgumentException.class, () ->
                tagService.deleteTag(null));
    }

    @Test
    void deleteTag_throwsOnEmptyTagId() {
        assertThrows(IllegalArgumentException.class, () ->
                tagService.deleteTag(""));
    }
}
