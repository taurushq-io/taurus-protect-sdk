package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.UserMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.model.UserResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.UsersApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetUserReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetUsersReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalUser;
import com.taurushq.sdk.protect.openapi.model.UserServiceCreateAttributeBody;

import java.util.ArrayList;
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
 * // List users, one page at a time
 * UserResult page = client.getUserService().getUsers(20, 0);
 * // next page: getUsers(20, page.getPagination().getNextOffset())
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
     * Gets current user details, including whether the user, its group memberships and its
     * public key are enforced in the governance rules.
     *
     * @return the current user
     * @throws ApiException the api exception
     */
    public User getMe() throws ApiException {
        try {
            TgvalidatordGetMeReply reply = usersApi.userServiceGetMe(
                    false,  // includeKeyContainer
                    true    // checkEnforcedInRules: the only way GetMe computes the flags
            );
            return EnforcedInRules.computed(UserMapper.INSTANCE.fromDTO(reply.getResult()));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets a user by ID.
     * <p>
     * This endpoint does not compute the enforced-in-rules flags: they are {@code null} unless
     * the reply carries them.
     *
     * @param userId the user ID
     * @return the user
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if userId is null or empty
     */
    public User getUser(final String userId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(userId), "userId cannot be null or empty");

        try {
            TgvalidatordGetUserReply reply = usersApi.userServiceGetUser(userId);
            return UserMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a page of users.
     *
     * @param limit  the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset the offset, 0 for the first page
     * @return the users and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public UserResult getUsers(final int limit, final long offset) throws ApiException {
        return listUsers(null, limit, offset);
    }

    /**
     * Gets the users with the given email addresses, walking every page.
     *
     * @param emails the list of email addresses
     * @return the matching users
     * @throws ApiException the api exception
     */
    public List<User> getUsersByEmail(final List<String> emails) throws ApiException {
        checkNotNull(emails, "emails cannot be null");
        checkArgument(!emails.isEmpty(), "emails cannot be empty");

        List<User> users = new ArrayList<>();
        long offset = 0;
        UserResult page;
        do {
            page = listUsers(emails, Pagination.MAX_PAGE_SIZE, offset);
            users.addAll(page.getUsers());
            offset = page.getPagination().getNextOffset();
        } while (page.getPagination().hasMore());
        return users;
    }

    private UserResult listUsers(final List<String> emails, final int limit, final long offset)
            throws ApiException {
        final int size = PagedOperation.USERS.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);

        try {
            TgvalidatordGetUsersReply reply = usersApi.userServiceGetUsers(
                    String.valueOf(size),       // limit
                    from == 0 ? null : String.valueOf(from), // offset
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

            List<TgvalidatordInternalUser> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            return new UserResult(EnforcedInRules.computedUsers(UserMapper.INSTANCE.fromDTO(rows)),
                    PagedOperation.USERS.offsetPage(size, from, rows.size(), 0,
                            reply.getTotalItems(), null));
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
