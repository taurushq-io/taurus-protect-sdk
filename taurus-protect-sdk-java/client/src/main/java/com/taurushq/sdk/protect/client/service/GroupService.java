package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.GroupMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.GroupResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.GroupsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetGroupsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalGroup;

import java.util.Collections;
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
 * // First page of groups (default page size)
 * GroupResult groups = client.getGroupService().getGroups();
 *
 * // A page of 10 groups, then the next one
 * GroupResult page = client.getGroupService().getGroups(10, 0, null, null, null);
 * page = client.getGroupService().getGroups(10, page.getPagination().getNextOffset(), null, null, null);
 *
 * // Search for groups
 * GroupResult admins = client.getGroupService().getGroups(0, 0, null, null, "admin");
 * }</pre>
 *
 * @see com.taurushq.sdk.protect.client.model.Group
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
     * Retrieves the first page of groups, with the default page size.
     *
     * @return the groups and their pagination
     * @throws ApiException if the API call fails
     */
    public GroupResult getGroups() throws ApiException {
        return getGroups(0, 0, null, null, null);
    }

    /**
     * Retrieves a page of groups with optional filters.
     *
     * @param limit            the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset           the offset, 0 for the first page
     * @param ids              optional list of group IDs to filter by
     * @param externalGroupIds optional list of external group IDs to filter by
     * @param query            optional query string to search groups
     * @return the groups and their pagination
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public GroupResult getGroups(final int limit,
                                 final long offset,
                                 final List<String> ids,
                                 final List<String> externalGroupIds,
                                 final String query) throws ApiException {
        final int size = PagedOperation.GROUPS.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);
        try {
            TgvalidatordGetGroupsReply reply = groupsApi.userServiceGetGroups(
                    String.valueOf(size), from == 0 ? null : String.valueOf(from),
                    ids, externalGroupIds, query);
            List<TgvalidatordInternalGroup> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            return new GroupResult(EnforcedInRules.computedGroups(groupMapper.fromDTOList(rows)),
                    PagedOperation.GROUPS.offsetPage(size, from, rows.size(), 0,
                            reply.getTotalItems(), null));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
