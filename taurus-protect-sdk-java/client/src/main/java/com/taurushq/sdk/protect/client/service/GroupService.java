package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.GroupMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.GroupsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetGroupsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing user groups in the Taurus Protect system.
 * <p>
 * Groups are used to organize users and define access permissions
 * and approval workflows.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all groups
 * List<Group> groups = client.getGroupService().getGroups();
 *
 * // Get groups with pagination
 * List<Group> groups = client.getGroupService().getGroups("10", "0", null, null, null);
 *
 * // Search for groups
 * List<Group> groups = client.getGroupService().getGroups(null, null, null, null, "admin");
 * }</pre>
 *
 * @see Group
 */
public class GroupService {

    private final GroupsApi groupsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final GroupMapper groupMapper;

    /**
     * Creates a new GroupService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public GroupService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.groupsApi = new GroupsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.groupMapper = GroupMapper.INSTANCE;
    }

    /**
     * Retrieves all groups.
     *
     * @return the list of groups
     * @throws ApiException if the API call fails
     */
    public List<Group> getGroups() throws ApiException {
        return getGroups(null, null, null, null, null);
    }

    /**
     * Retrieves groups with optional filters and pagination.
     *
     * @param limit            maximum number of results to return
     * @param offset           number of results to skip
     * @param ids              optional list of group IDs to filter by
     * @param externalGroupIds optional list of external group IDs to filter by
     * @param query            optional query string to search groups
     * @return the list of groups matching the filters
     * @throws ApiException if the API call fails
     */
    public List<Group> getGroups(final String limit,
                                  final String offset,
                                  final List<String> ids,
                                  final List<String> externalGroupIds,
                                  final String query) throws ApiException {
        try {
            TgvalidatordGetGroupsReply reply = groupsApi.userServiceGetGroups(
                    limit, offset, ids, externalGroupIds, query);
            return groupMapper.fromDTOList(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
