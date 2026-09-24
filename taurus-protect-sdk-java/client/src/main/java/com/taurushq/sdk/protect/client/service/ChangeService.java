package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ChangeMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Change;
import com.taurushq.sdk.protect.client.model.ChangeResult;
import com.taurushq.sdk.protect.client.model.CreateChangeRequest;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ChangesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveChangesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordChange;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateChangeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateChangeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetChangeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetChangesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRejectChangesRequest;

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
 * // Get changes pending approval, one page at a time
 * ChangeResult result = client.getChangeService().getChangesForApproval(20, null);
 * // next page: getChangesForApproval(20, result.getPage().getNextCursor())
 *
 * // Approve a change
 * client.getChangeService().approveChange(changeId);
 *
 * // Reject a change
 * client.getChangeService().rejectChange(changeId);
 *
 * // Get changes with filters
 * ChangeResult filtered = client.getChangeService()
 *     .getChanges("user", "pending", 20, null);
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
     * Gets a page of changes with filters.
     *
     * @param entity   the entity type to filter by (optional)
     * @param status   the status to filter by (optional)
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the changes and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ChangeResult getChanges(final String entity, final String status, final Integer pageSize,
                                   final String cursor) throws ApiException {
        return getChanges(entity, status, Pagination.page(pageSize, cursor));
    }

    /**
     * Gets a page of changes with filters, with a low-level request cursor.
     *
     * @param entity the entity type to filter by (optional)
     * @param status the status to filter by (optional)
     * @param cursor the request cursor, null for the first page with the default size
     * @return the changes and their page
     * @throws ApiException the api exception
     */
    public ChangeResult getChanges(final String entity, final String status, final ApiRequestCursor cursor)
            throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetChangesReply reply = changesApi.changeServiceGetChanges(
                    entity,                             // entity
                    null,                               // entityId
                    status,                             // status
                    null,                               // creatorId
                    null,                               // sortOrder
                    page.currentPage(),                 // cursorCurrentPage
                    page.pageRequest(),                 // cursorPageRequest
                    page.pageSizeParam(),               // cursorPageSize
                    null,                               // entityIDs
                    null                                // entityUUIDs
            );

            return page.complete(toResult(reply.getResult()), PagedOperation.CHANGES, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a page of changes pending approval.
     *
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the changes and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ChangeResult getChangesForApproval(final Integer pageSize, final String cursor) throws ApiException {
        return getChangesForApproval(Pagination.page(pageSize, cursor));
    }

    /**
     * Gets a page of changes pending approval, with a low-level request cursor.
     *
     * @param cursor the request cursor, null for the first page with the default size
     * @return the changes and their page
     * @throws ApiException the api exception
     */
    public ChangeResult getChangesForApproval(final ApiRequestCursor cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetChangesReply reply = changesApi.changeServiceGetChangesForApproval(
                    null,                               // entities
                    null,                               // sortOrder
                    page.currentPage(),                 // cursorCurrentPage
                    page.pageRequest(),                 // cursorPageRequest
                    page.pageSizeParam(),               // cursorPageSize
                    null,                               // entityIDs
                    null                                // entityUUIDs
            );

            return page.complete(toResult(reply.getResult()), PagedOperation.CHANGES_FOR_APPROVAL,
                    reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    private static ChangeResult toResult(final List<TgvalidatordChange> changes) {
        ChangeResult result = new ChangeResult();
        result.setChanges(changes == null ? Collections.emptyList() : ChangeMapper.INSTANCE.fromDTO(changes));
        return result;
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
