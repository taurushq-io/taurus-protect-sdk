package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WebhookStatus;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhooksReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWebhook;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting webhook DTOs to domain models.
 * <p>
 * This mapper transforms OpenAPI-generated webhook objects into the SDK's
 * domain model objects.
 *
 * @see Webhook
 * @see TgvalidatordWebhook
 */
@Mapper
public interface WebhookMapper {

    /**
     * Singleton instance of the mapper.
     */
    WebhookMapper INSTANCE = Mappers.getMapper(WebhookMapper.class);

    /**
     * Converts a webhook DTO to a domain model.
     *
     * @param webhook the OpenAPI webhook DTO
     * @return the domain model webhook
     */
    @Mapping(target = "status", source = "status", qualifiedByName = "toWebhookStatus")
    Webhook fromDTO(TgvalidatordWebhook webhook);

    /**
     * Converts a list of webhook DTOs to domain models.
     *
     * @param webhooks the list of OpenAPI webhook DTOs
     * @return the list of domain model webhooks
     */
    List<Webhook> fromDTOList(List<TgvalidatordWebhook> webhooks);

    /**
     * Converts a get webhooks reply to a webhook result.
     *
     * @param reply the OpenAPI reply
     * @return the webhook result with pagination
     */
    @Mapping(target = "webhooks", source = "webhooks")
    @Mapping(target = "cursor", source = "cursor")
    WebhookResult fromReply(TgvalidatordGetWebhooksReply reply);

    /**
     * Converts a response cursor DTO to a domain model.
     *
     * @param cursor the OpenAPI response cursor
     * @return the domain model cursor
     */
    ApiResponseCursor fromCursor(TgvalidatordResponseCursor cursor);

    /**
     * Converts a status string to a WebhookStatus enum.
     *
     * @param status the status string
     * @return the WebhookStatus enum value
     */
    @Named("toWebhookStatus")
    default WebhookStatus toWebhookStatus(String status) {
        return WebhookStatus.fromValue(status);
    }
}
