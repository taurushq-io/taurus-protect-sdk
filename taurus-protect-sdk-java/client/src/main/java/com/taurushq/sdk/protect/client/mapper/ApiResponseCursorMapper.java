package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for cursor pagination.
 */
@Mapper
public interface ApiResponseCursorMapper {

    /**
     * The constant INSTANCE.
     */
    ApiResponseCursorMapper INSTANCE = Mappers.getMapper(ApiResponseCursorMapper.class);

    /**
     * Converts a TgvalidatordResponseCursor DTO to ApiResponseCursor model.
     *
     * @param dto the DTO
     * @return the model
     */
    ApiResponseCursor fromDTO(TgvalidatordResponseCursor dto);

    /**
     * Converts an ApiRequestCursor model to TgvalidatordRequestCursor DTO.
     *
     * @param cursor the model
     * @return the DTO
     */
    @Mapping(source = "pageRequest", target = "pageRequest", qualifiedByName = "pageRequestToString")
    @Mapping(source = "pageSize", target = "pageSize", qualifiedByName = "longToString")
    TgvalidatordRequestCursor toDTO(ApiRequestCursor cursor);

    /**
     * Converts PageRequest enum to String.
     *
     * @param pageRequest the enum value
     * @return the string value
     */
    @Named("pageRequestToString")
    default String pageRequestToString(PageRequest pageRequest) {
        return pageRequest != null ? pageRequest.name() : null;
    }

    /**
     * Converts long to String.
     *
     * @param value the long value
     * @return the string value
     */
    @Named("longToString")
    default String longToString(long value) {
        return String.valueOf(value);
    }
}
