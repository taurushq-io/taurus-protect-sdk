package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TagMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Tag;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TagsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateTagReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateTagRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetTagsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing tags in the Taurus Protect system.
 * <p>
 * Tags are used to label and categorize various resources such as
 * wallets, addresses, and transactions for easier organization and filtering.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all tags
 * List<Tag> tags = client.getTagService().getTags();
 *
 * // Create a new tag
 * Tag newTag = client.getTagService().createTag("Production", "#FF0000");
 *
 * // Delete a tag
 * client.getTagService().deleteTag("tag-123");
 * }</pre>
 *
 * @see Tag
 */
public class TagService {

    private final TagsApi tagsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TagMapper tagMapper;

    /**
     * Creates a new TagService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public TagService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.tagsApi = new TagsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.tagMapper = TagMapper.INSTANCE;
    }

    /**
     * Retrieves all tags.
     *
     * @return the list of tags
     * @throws ApiException if the API call fails
     */
    public List<Tag> getTags() throws ApiException {
        return getTags(null, null);
    }

    /**
     * Retrieves tags with optional filters.
     *
     * @param ids   optional list of tag IDs to filter by
     * @param query optional query string to filter tags by value
     * @return the list of tags matching the filters
     * @throws ApiException if the API call fails
     */
    public List<Tag> getTags(final List<String> ids, final String query) throws ApiException {
        try {
            TgvalidatordGetTagsReply reply = tagsApi.tagServiceGetTags(ids, query);
            return tagMapper.fromDTOList(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Creates a new tag.
     *
     * @param value the tag value/label
     * @param color the tag color (e.g., hex code like "#FF0000")
     * @return the created tag
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if value is null or empty
     */
    public Tag createTag(final String value, final String color) throws ApiException {
        checkArgument(value != null && !value.isEmpty(), "value must not be null or empty");
        try {
            TgvalidatordCreateTagRequest request = new TgvalidatordCreateTagRequest();
            request.setValue(value);
            request.setColor(color);
            TgvalidatordCreateTagReply reply = tagsApi.tagServiceCreateTag(request);
            return tagMapper.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Deletes a tag by ID.
     * <p>
     * This will also remove all assignments of this tag from all entities.
     *
     * @param tagId the ID of the tag to delete
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if tagId is null or empty
     */
    public void deleteTag(final String tagId) throws ApiException {
        checkArgument(tagId != null && !tagId.isEmpty(), "tagId must not be null or empty");
        try {
            tagsApi.tagServiceDeleteTag(tagId);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
