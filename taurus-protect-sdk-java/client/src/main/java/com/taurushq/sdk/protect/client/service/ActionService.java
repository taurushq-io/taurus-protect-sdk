package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ActionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ActionEnvelope;
import com.taurushq.sdk.protect.client.model.ActionResult;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ActionsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetActionReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetActionsReply;

import java.util.Collections;
import java.util.List;

/**
 * Service for managing automated actions in the Taurus Protect system.
 * <p>
 * Actions allow automated workflows to be triggered based on specific conditions
 * such as balance thresholds. When conditions are met, tasks like transfers or
 * notifications can be executed automatically.
 * <p>
 * Example usage:
 * <pre>{@code
 * // First page of actions (default page size)
 * ActionResult actions = client.getActionService().getActions();
 *
 * // Get a specific action
 * ActionEnvelope action = client.getActionService().getAction("action-123");
 * }</pre>
 *
 * @see ActionEnvelope
 */
public class ActionService {

    private final ActionsApi actionsApi;
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Creates a new ActionService.
     *
     * @param apiClient          the API client for making requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public ActionService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        Preconditions.checkNotNull(apiClient, "apiClient must not be null");
        Preconditions.checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.actionsApi = new ActionsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
    }

    /**
     * Retrieves the first page of actions, with the default page size.
     *
     * @return the actions and their pagination
     * @throws ApiException if the API call fails
     */
    public ActionResult getActions() throws ApiException {
        return getActions(0, 0, null);
    }

    /**
     * Retrieves a page of actions with optional filters.
     *
     * @param limit  the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset the offset, 0 for the first page
     * @param ids    optional list of action IDs to filter by
     * @return the actions and their pagination
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public ActionResult getActions(final int limit, final long offset,
                                   final List<String> ids) throws ApiException {
        final int size = PagedOperation.ACTIONS.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);
        try {
            TgvalidatordGetActionsReply reply = actionsApi.actionServiceGetActions(
                    String.valueOf(size), from == 0 ? null : String.valueOf(from), ids);
            List<TgvalidatordActionEnvelope> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            return new ActionResult(ActionMapper.INSTANCE.fromDTOList(rows),
                    PagedOperation.ACTIONS.offsetPage(size, from, rows.size(), 0,
                            reply.getTotalItems(), null));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a specific action by its ID.
     *
     * @param actionId the ID of the action to retrieve
     * @return the action envelope
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if actionId is null or empty
     */
    public ActionEnvelope getAction(final String actionId) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(actionId), "actionId must not be null or empty");
        try {
            TgvalidatordGetActionReply reply = actionsApi.actionServiceGetAction(actionId);
            return ActionMapper.INSTANCE.fromDTO(reply.getAction());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
