package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.WebhookCall;
import com.taurushq.sdk.protect.client.model.WebhookCallResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhookCallsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWebhookCall;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting webhook call DTOs to domain models.
 *
 * @see WebhookCall
 * @see TgvalidatordWebhookCall
 */
@Mapper
public interface WebhookCallsMapper {

    /**
     * Singleton instance of the mapper.
     */
    WebhookCallsMapper INSTANCE = Mappers.getMapper(WebhookCallsMapper.class);

    /**
     * Converts a webhook call DTO to a domain model.
     *
     * @param dto the OpenAPI webhook call DTO
     * @return the domain model webhook call
     */
    WebhookCall fromDTO(TgvalidatordWebhookCall dto);

    /**
     * Converts a list of webhook call DTOs to domain models.
     *
     * @param dtos the list of OpenAPI webhook call DTOs
     * @return the list of domain model webhook calls
     */
    List<WebhookCall> fromDTOList(List<TgvalidatordWebhookCall> dtos);

    /**
     * Converts a get webhook calls reply to a webhook call result.
     *
     * @param reply the OpenAPI reply
     * @return the webhook call result with pagination
     */
    @Mapping(target = "calls", source = "calls")
    @Mapping(target = "cursor", source = "cursor")
    WebhookCallResult fromReply(TgvalidatordGetWebhookCallsReply reply);

    /**
     * Converts a response cursor DTO to a domain model.
     *
     * @param cursor the OpenAPI response cursor
     * @return the domain model cursor
     */
    ApiResponseCursor fromCursor(TgvalidatordResponseCursor cursor);
}
