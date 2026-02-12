package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ActionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ActionEnvelope;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ActionsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetActionReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetActionsReply;

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
 * // Get all actions
 * List<ActionEnvelope> actions = client.getActionService().getActions();
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
     * Retrieves all actions.
     *
     * @return a list of all action envelopes
     * @throws ApiException if the API call fails
     */
    public List<ActionEnvelope> getActions() throws ApiException {
        return getActions(null, null, null);
    }

    /**
     * Retrieves actions with optional filters.
     *
     * @param limit  the maximum number of actions to return (optional)
     * @param offset the offset for pagination (optional)
     * @param ids    optional list of action IDs to filter by
     * @return a list of action envelopes
     * @throws ApiException if the API call fails
     */
    public List<ActionEnvelope> getActions(final String limit, final String offset,
                                           final List<String> ids) throws ApiException {
        try {
            TgvalidatordGetActionsReply reply = actionsApi.actionServiceGetActions(limit, offset, ids);
            return ActionMapper.INSTANCE.fromDTOList(reply.getResult());
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
