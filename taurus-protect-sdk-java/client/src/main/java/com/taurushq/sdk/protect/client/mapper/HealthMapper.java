package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.HealthCheckStatus;
import com.taurushq.sdk.protect.client.model.HealthComponent;
import com.taurushq.sdk.protect.client.model.HealthGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealth;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealthComponent;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealthGroup;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting between OpenAPI health DTOs and client model objects.
 */
@Mapper
public interface HealthMapper {

    /**
     * Singleton instance of the mapper.
     */
    HealthMapper INSTANCE = Mappers.getMapper(HealthMapper.class);

    /**
     * Converts an OpenAPI health DTO to a client model HealthCheckStatus.
     *
     * @param dto the OpenAPI DTO
     * @return the client model HealthCheckStatus
     */
    HealthCheckStatus fromDTO(TgvalidatordHealth dto);

    /**
     * Converts a list of OpenAPI health DTOs to client model HealthCheckStatus list.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of client model HealthCheckStatus
     */
    List<HealthCheckStatus> fromDTOList(List<TgvalidatordHealth> dtos);

    /**
     * Converts an OpenAPI health group DTO to a client model HealthGroup.
     *
     * @param dto the OpenAPI DTO
     * @return the client model HealthGroup
     */
    HealthGroup fromGroupDTO(TgvalidatordHealthGroup dto);

    /**
     * Converts an OpenAPI health component DTO to a client model HealthComponent.
     *
     * @param dto the OpenAPI DTO
     * @return the client model HealthComponent
     */
    HealthComponent fromComponentDTO(TgvalidatordHealthComponent dto);
}
