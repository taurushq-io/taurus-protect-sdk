package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.WebhookCallsMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Pagination;
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
 * // Get webhook calls, first page
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, null, null, null, 20, null);
 *
 * // Get calls for a specific webhook
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, "webhook-123", null, null, 20, null);
 *
 * // Get failed calls only
 * WebhookCallResult result = client.getWebhookCallsService()
 *     .getWebhookCalls(null, null, "FAILED", null, 20, null);
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
     * Retrieves a page of webhook calls with optional filtering.
     *
     * @param eventId   filter by event ID (optional)
     * @param webhookId filter by webhook ID (optional)
     * @param status    filter by call status (optional, e.g., "SUCCESS", "FAILED")
     * @param sortOrder sort order for results (optional, "ASC" or "DESC", default "DESC")
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the webhook calls and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public WebhookCallResult getWebhookCalls(final String eventId, final String webhookId,
                                              final String status, final String sortOrder,
                                              final Integer pageSize, final String cursor)
            throws ApiException {
        return getWebhookCalls(eventId, webhookId, status, sortOrder, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of webhook calls with optional filtering, with a low-level request cursor.
     *
     * @param eventId   filter by event ID (optional)
     * @param webhookId filter by webhook ID (optional)
     * @param status    filter by call status (optional, e.g., "SUCCESS", "FAILED")
     * @param sortOrder sort order for results (optional, "ASC" or "DESC", default "DESC")
     * @param cursor    the request cursor, null for the first page with the default size
     * @return the webhook calls and their page
     * @throws ApiException if the API call fails
     */
    public WebhookCallResult getWebhookCalls(final String eventId, final String webhookId,
                                              final String status, final String sortOrder,
                                              final ApiRequestCursor cursor)
            throws ApiException {

        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetWebhookCallsReply reply = webhookCallsApi.webhookServiceGetWebhookCalls(
                    eventId,
                    webhookId,
                    status,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    sortOrder
            );
            return page.complete(mapper.fromReply(reply), PagedOperation.WEBHOOK_CALLS);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
