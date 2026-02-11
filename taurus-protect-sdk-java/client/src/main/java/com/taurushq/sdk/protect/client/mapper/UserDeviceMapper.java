package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.UserDevicePairing;
import com.taurushq.sdk.protect.client.model.UserDevicePairingInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateUserDevicePairingReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUserDevicePairingInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting user device-related DTOs to domain models.
 */
@Mapper
public interface UserDeviceMapper {

    UserDeviceMapper INSTANCE = Mappers.getMapper(UserDeviceMapper.class);

    /**
     * Converts a CreateUserDevicePairing reply to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    @Mapping(source = "pairingID", target = "pairingId")
    UserDevicePairing fromCreateDTO(TgvalidatordCreateUserDevicePairingReply dto);

    /**
     * Converts a UserDevicePairingInfo DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    @Mapping(source = "pairingID", target = "pairingId")
    @Mapping(target = "status", expression = "java(dto.getStatus().getValue())")
    UserDevicePairingInfo fromInfoDTO(TgvalidatordUserDevicePairingInfo dto);
}
