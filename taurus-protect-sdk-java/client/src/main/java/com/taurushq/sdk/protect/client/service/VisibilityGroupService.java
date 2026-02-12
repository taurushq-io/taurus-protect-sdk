package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.UserMapper;
import com.taurushq.sdk.protect.client.mapper.VisibilityGroupMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.model.VisibilityGroup;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.RestrictedVisibilityGroupsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetUsersByVisibilityGroupIDReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetVisibilityGroupsReply;

import java.util.List;

/**
 * Service for managing visibility groups in the Taurus Protect system.
 * <p>
 * Visibility groups are used to control data access. Users can only see
 * wallets, addresses, and other entities that belong to their assigned
 * visibility groups. This provides fine-grained access control within
 * a tenant.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all visibility groups
 * List<VisibilityGroup> groups = client.getVisibilityGroupService().getVisibilityGroups();
 *
 * // Get users in a specific visibility group
 * List<User> users = client.getVisibilityGroupService().getUsersByVisibilityGroup("vg-123");
 * }</pre>
 *
 * @see VisibilityGroup
 */
public class VisibilityGroupService {

    private final RestrictedVisibilityGroupsApi visibilityGroupsApi;
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Creates a new VisibilityGroupService.
     *
     * @param apiClient          the API client for making requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public VisibilityGroupService(final ApiClient apiClient,
                                   final ApiExceptionMapper apiExceptionMapper) {
        Preconditions.checkNotNull(apiClient, "apiClient must not be null");
        Preconditions.checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.visibilityGroupsApi = new RestrictedVisibilityGroupsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
    }

    /**
     * Retrieves all visibility groups.
     *
     * @return a list of all visibility groups
     * @throws ApiException if the API call fails
     */
    public List<VisibilityGroup> getVisibilityGroups() throws ApiException {
        try {
            TgvalidatordGetVisibilityGroupsReply reply =
                    visibilityGroupsApi.userServiceGetVisibilityGroups();
            return VisibilityGroupMapper.INSTANCE.fromDTOList(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves users assigned to a specific visibility group.
     *
     * @param visibilityGroupId the visibility group ID
     * @return a list of users in the visibility group
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if visibilityGroupId is null or empty
     */
    public List<User> getUsersByVisibilityGroup(final String visibilityGroupId) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(visibilityGroupId),
                "visibilityGroupId must not be null or empty");
        try {
            TgvalidatordGetUsersByVisibilityGroupIDReply reply =
                    visibilityGroupsApi.userServiceGetUsersByVisibilityGroupID(visibilityGroupId);
            return UserMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
