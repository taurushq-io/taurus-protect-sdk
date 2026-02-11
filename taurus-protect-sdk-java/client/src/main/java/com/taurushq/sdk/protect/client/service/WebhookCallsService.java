package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.WebhookCallsMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.WebhookCallResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.WebhookCallsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhookCallsReply;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving webhook call history in the Taurus Protect system.
 * <p>
 * This service provides access to the history of webhook invocations,
 * including their delivery status and payload information.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all webhook calls
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, null, null, null, null);
 *
 * // Get calls for a specific webhook
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, "webhook-123", null, null, null);
 *
 * // Get failed calls only
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, null, "FAILED", null, null);
 * }</pre>
 *
 * @see WebhookCallResult
 */
public class WebhookCallsService {

    /**
     * The underlying OpenAPI client for webhook calls operations.
     */
    private final WebhookCallsApi webhookCallsApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Mapper for converting webhook calls DTOs to domain models.
     */
    private final WebhookCallsMapper mapper;

    /**
     * Instantiates a new Webhook calls service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public WebhookCallsService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.webhookCallsApi = new WebhookCallsApi(openApiClient);
        this.mapper = WebhookCallsMapper.INSTANCE;
    }

    /**
     * Retrieves webhook call history with optional filtering.
     * <p>
     * Returns a paginated list of webhook calls that can be filtered by
     * event ID, webhook ID, or status.
     *
     * @param eventId   filter by event ID (optional)
     * @param webhookId filter by webhook ID (optional)
     * @param status    filter by call status (optional, e.g., "SUCCESS", "FAILED")
     * @param sortOrder sort order for results (optional, "ASC" or "DESC", default "DESC")
     * @param cursor    pagination cursor (optional, null for first page)
     * @return a paginated result containing webhook calls
     * @throws ApiException if the API call fails
     */
    public WebhookCallResult getWebhookCalls(final String eventId, final String webhookId,
                                              final String status, final String sortOrder,
                                              final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetWebhookCallsReply reply = webhookCallsApi.webhookServiceGetWebhookCalls(
                    eventId,
                    webhookId,
                    status,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    sortOrder
            );
            return mapper.fromReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
