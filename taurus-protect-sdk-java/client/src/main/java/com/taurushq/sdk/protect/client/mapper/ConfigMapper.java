package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.TenantConfig;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTenantConfig;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting configuration-related DTOs to domain models.
 */
@Mapper
public interface ConfigMapper {

    ConfigMapper INSTANCE = Mappers.getMapper(ConfigMapper.class);

    /**
     * Converts a TenantConfig DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    @Mapping(source = "isMFAMandatory", target = "MFAMandatory")
    @Mapping(source = "isProtectEngineCold", target = "protectEngineCold")
    @Mapping(source = "isColdProtectEngineOffline", target = "coldProtectEngineOffline")
    @Mapping(source = "isPhysicalAirGapEnabled", target = "physicalAirGapEnabled")
    TenantConfig fromDTO(TgvalidatordTenantConfig dto);
}
