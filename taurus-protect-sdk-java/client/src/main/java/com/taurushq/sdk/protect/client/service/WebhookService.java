package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.WebhookMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WebhookStatus;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.WebhooksApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWebhookReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWebhookRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhooksReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUpdateWebhookStatusReply;
import com.taurushq.sdk.protect.openapi.model.WebhookServiceUpdateWebhookStatusBody;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing webhooks in the Taurus Protect system.
 * <p>
 * Webhooks allow external systems to receive real-time notifications
 * about events such as transactions, request status changes, and other
 * platform events.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a webhook for transaction events
 * Webhook webhook = client.getWebhookService().createWebhook(
 *     "https://example.com/webhook",
 *     "TRANSACTION",
 *     "my-webhook-secret"
 * );
 *
 * // List all webhooks
 * WebhookResult result = client.getWebhookService().getWebhooks(null, null, null);
 * for (Webhook wh : result.getWebhooks()) {
 *     System.out.println(wh.getUrl() + " - " + wh.getStatus());
 * }
 *
 * // Disable a webhook
 * client.getWebhookService().updateWebhookStatus(webhookId, WebhookStatus.DISABLED);
 *
 * // Delete a webhook
 * client.getWebhookService().deleteWebhook(webhookId);
 * }</pre>
 *
 * @see Webhook
 * @see WebhookStatus
 * @see WebhookResult
 */
public class WebhookService {

    /**
     * The underlying OpenAPI client for webhook operations.
     */
    private final WebhooksApi webhooksApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Mapper for converting webhook DTOs to domain models.
     */
    private final WebhookMapper webhookMapper;

    /**
     * Instantiates a new Webhook service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public WebhookService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.webhooksApi = new WebhooksApi(openApiClient);
        this.webhookMapper = WebhookMapper.INSTANCE;
    }

    /**
     * Creates a new webhook configuration.
     * <p>
     * The webhook will be created in an enabled state by default.
     * The secret is used to sign webhook payloads so you can verify
     * their authenticity.
     *
     * @param url    the URL to receive webhook notifications (must be HTTPS)
     * @param type   the type of events to subscribe to (e.g., "TRANSACTION", "REQUEST")
     * @param secret the secret key for signing webhook payloads
     * @return the created webhook with its assigned ID
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if url or type is null or empty
     */
    public Webhook createWebhook(final String url, final String type, final String secret) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(url), "url cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(type), "type cannot be null or empty");

        TgvalidatordCreateWebhookRequest request = new TgvalidatordCreateWebhookRequest();
        request.setUrl(url);
        request.setType(type);
        request.setSecret(secret);

        try {
            TgvalidatordCreateWebhookReply reply = webhooksApi.webhookServiceCreateWebhook(request);
            return webhookMapper.fromDTO(reply.getWebhook());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a paginated list of webhooks.
     * <p>
     * Results can be filtered by type and URL. Use the cursor for pagination.
     *
     * @param type   filter by webhook type (optional)
     * @param url    filter by webhook URL (optional)
     * @param cursor pagination cursor (optional, null for first page)
     * @return a paginated result containing webhooks
     * @throws ApiException if the API call fails
     */
    public WebhookResult getWebhooks(final String type, final String url, final ApiRequestCursor cursor)
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
            TgvalidatordGetWebhooksReply reply = webhooksApi.webhookServiceGetWebhooks(
                    type,
                    url,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    null  // sortOrder - use default
            );
            return webhookMapper.fromReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Deletes a webhook configuration.
     * <p>
     * Once deleted, the webhook will no longer receive notifications.
     *
     * @param webhookId the ID of the webhook to delete
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if webhookId is null or empty
     */
    public void deleteWebhook(final String webhookId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(webhookId), "webhookId cannot be null or empty");

        try {
            webhooksApi.webhookServiceDeleteWebhook(webhookId);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Updates the status of a webhook.
     * <p>
     * Use this method to enable or disable a webhook without deleting it.
     *
     * @param webhookId the ID of the webhook to update
     * @param status    the new status (ENABLED or DISABLED)
     * @return the updated webhook
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if webhookId is null/empty or status is null
     */
    public Webhook updateWebhookStatus(final String webhookId, final WebhookStatus status) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(webhookId), "webhookId cannot be null or empty");
        checkNotNull(status, "status cannot be null");

        WebhookServiceUpdateWebhookStatusBody body = new WebhookServiceUpdateWebhookStatusBody();
        body.setStatus(status.getValue());

        try {
            TgvalidatordUpdateWebhookStatusReply reply = webhooksApi.webhookServiceUpdateWebhookStatus(webhookId, body);
            return webhookMapper.fromDTO(reply.getWebhook());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
