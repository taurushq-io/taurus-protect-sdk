package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.mapper.ChangeMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Change;
import com.taurushq.sdk.protect.client.model.ChangeResult;
import com.taurushq.sdk.protect.client.model.CreateChangeRequest;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ChangesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveChangesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordChange;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateChangeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateChangeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetChangeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetChangesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRejectChangesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing configuration changes in the Taurus Protect system.
 * <p>
 * Changes represent modifications to system configuration that require approval
 * before taking effect. This includes user management, role assignments, and
 * other administrative operations that follow an approval workflow.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get changes pending approval
 * ApiRequestCursor cursor = Pagination.first(50);
 * ChangeResult result = client.getChangeService().getChangesForApproval(cursor);
 *
 * // Approve a change
 * client.getChangeService().approveChange(changeId);
 *
 * // Reject a change
 * client.getChangeService().rejectChange(changeId);
 *
 * // Get changes with filters
 * ChangeResult filtered = client.getChangeService()
 *     .getChanges("user", "pending", cursor);
 * }</pre>
 *
 * @see ChangeResult
 * @see Change
 */
public class ChangeService {

    /**
     * The underlying OpenAPI client for change operations.
     */
    private final ChangesApi changesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Change service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public ChangeService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.changesApi = new ChangesApi(openApiClient);
    }


    /**
     * Creates a change request.
     *
     * @param request the change request with action, entity, changes, etc.
     * @return the created change id
     * @throws ApiException the api exception
     */
    public String createChange(final CreateChangeRequest request) throws ApiException {
        checkNotNull(request, "request cannot be null");
        checkArgument(!Strings.isNullOrEmpty(request.getAction()), "action cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(request.getEntity()), "entity cannot be null or empty");

        try {
            TgvalidatordCreateChangeRequest body = ChangeMapper.INSTANCE.toDTO(request);
            TgvalidatordCreateChangeReply reply = changesApi.changeServiceCreateChange(body);
            return reply.getResult().getId();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a change by id.
     *
     * @param id the change id
     * @return the change
     * @throws ApiException the api exception
     */
    public Change getChange(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetChangeReply reply = changesApi.changeServiceGetChange(id);
            return ChangeMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets changes with filters.
     *
     * @param entity the entity type to filter by
     * @param status the status to filter by
     * @param cursor the request cursor for pagination
     * @return the change result with list and response cursor
     * @throws ApiException the api exception
     */
    public ChangeResult getChanges(final String entity, final String status, final ApiRequestCursor cursor) throws ApiException {
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetChangesReply reply = changesApi.changeServiceGetChanges(
                    entity,                             // entity
                    null,                               // entityId
                    status,                             // status
                    null,                               // creatorId
                    null,                               // sortOrder
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // entityIDs
                    null                                // entityUUIDs
            );

            ChangeResult result = new ChangeResult();

            List<TgvalidatordChange> changes = reply.getResult();
            if (changes == null) {
                result.setChanges(Collections.emptyList());
            } else {
                result.setChanges(ChangeMapper.INSTANCE.fromDTO(changes));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets changes pending approval.
     *
     * @param cursor the request cursor for pagination
     * @return the change result with list and response cursor
     * @throws ApiException the api exception
     */
    public ChangeResult getChangesForApproval(final ApiRequestCursor cursor) throws ApiException {
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetChangesReply reply = changesApi.changeServiceGetChangesForApproval(
                    null,                               // entities
                    null,                               // sortOrder
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // entityIDs
                    null                                // entityUUIDs
            );

            ChangeResult result = new ChangeResult();

            List<TgvalidatordChange> changes = reply.getResult();
            if (changes == null) {
                result.setChanges(Collections.emptyList());
            } else {
                result.setChanges(ChangeMapper.INSTANCE.fromDTO(changes));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Approves a change.
     *
     * @param id the change id
     * @throws ApiException the api exception
     */
    public void approveChange(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            changesApi.changeServiceApproveChange(id, new Object());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Approves multiple changes.
     *
     * @param ids the list of change ids
     * @throws ApiException the api exception
     */
    public void approveChanges(final List<String> ids) throws ApiException {
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");

        try {
            TgvalidatordApproveChangesRequest body = new TgvalidatordApproveChangesRequest();
            body.setIds(ids);
            changesApi.changeServiceApproveChanges(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Rejects a change.
     *
     * @param id the change id
     * @throws ApiException the api exception
     */
    public void rejectChange(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            changesApi.changeServiceRejectChange(id, new Object());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Rejects multiple changes.
     *
     * @param ids the list of change ids
     * @throws ApiException the api exception
     */
    public void rejectChanges(final List<String> ids) throws ApiException {
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");

        try {
            TgvalidatordRejectChangesRequest body = new TgvalidatordRejectChangesRequest();
            body.setIds(ids);
            changesApi.changeServiceRejectChanges(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
