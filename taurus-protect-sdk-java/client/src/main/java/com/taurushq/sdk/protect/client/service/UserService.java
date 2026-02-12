package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.UserMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.UsersApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetUsersReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalUser;
import com.taurushq.sdk.protect.openapi.model.UserServiceCreateAttributeBody;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing users in the Taurus Protect system.
 * <p>
 * This service provides operations for retrieving user information and managing
 * user attributes. Users are individuals who can access the system and participate
 * in approval workflows.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get current authenticated user
 * User me = client.getUserService().getMe();
 *
 * // List users with pagination
 * List<User> users = client.getUserService().getUsers(50, 0);
 *
 * // Find users by email
 * List<User> found = client.getUserService()
 *     .getUsersByEmail(Arrays.asList("user@example.com"));
 *
 * // Create a user attribute
 * client.getUserService().createUserAttribute(userId, "department", "Finance");
 * }</pre>
 *
 * @see User
 * @see GovernanceRuleService
 */
public class UserService {

    /**
     * The underlying OpenAPI client for user operations.
     */
    private final UsersApi usersApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new User service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public UserService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.usersApi = new UsersApi(openApiClient);
    }


    /**
     * Gets current user details.
     *
     * @return the current user
     * @throws ApiException the api exception
     */
    public User getMe() throws ApiException {
        try {
            TgvalidatordGetMeReply reply = usersApi.userServiceGetMe(
                    false,  // includeKeyContainer
                    false   // checkEnforcedInRules
            );
            return UserMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets users with pagination.
     *
     * @param limit  the limit
     * @param offset the offset
     * @return the list of users
     * @throws ApiException the api exception
     */
    public List<User> getUsers(final int limit, final int offset) throws ApiException {
        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetUsersReply reply = usersApi.userServiceGetUsers(
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    null,                       // ids
                    null,                       // externalUserIds
                    null,                       // emails
                    null,                       // query
                    null,                       // publicKey
                    null,                       // excludeTechnicalUsers
                    null,                       // roles
                    null,                       // excludeIds
                    null,                       // nonTechnical
                    null                        // groupIds
            );

            List<TgvalidatordInternalUser> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return UserMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets users by email addresses.
     *
     * @param emails the list of email addresses
     * @return the list of users
     * @throws ApiException the api exception
     */
    public List<User> getUsersByEmail(final List<String> emails) throws ApiException {
        checkNotNull(emails, "emails cannot be null");
        checkArgument(!emails.isEmpty(), "emails cannot be empty");

        try {
            TgvalidatordGetUsersReply reply = usersApi.userServiceGetUsers(
                    null,                       // limit
                    null,                       // offset
                    null,                       // ids
                    null,                       // externalUserIds
                    emails,                     // emails
                    null,                       // query
                    null,                       // publicKey
                    null,                       // excludeTechnicalUsers
                    null,                       // roles
                    null,                       // excludeIds
                    null,                       // nonTechnical
                    null                        // groupIds
            );

            List<TgvalidatordInternalUser> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return UserMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Creates a user attribute.
     *
     * @param userId the user id
     * @param key    the attribute key
     * @param value  the attribute value
     * @throws ApiException the api exception
     */
    public void createUserAttribute(final long userId, final String key, final String value) throws ApiException {
        checkArgument(userId > 0, "userId must be positive");
        checkArgument(!Strings.isNullOrEmpty(key), "key cannot be null or empty");
        checkNotNull(value, "value cannot be null");

        try {
            UserServiceCreateAttributeBody body = new UserServiceCreateAttributeBody();
            body.setKey(key);
            body.setValue(value);

            usersApi.userServiceCreateAttribute(String.valueOf(userId), body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
