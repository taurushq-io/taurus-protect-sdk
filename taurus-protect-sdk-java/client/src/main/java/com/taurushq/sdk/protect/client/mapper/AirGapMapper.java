package com.taurushq.sdk.protect.client.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for air gap operations.
 * <p>
 * The AirGap API primarily deals with file transfers (binary data),
 * so this mapper provides minimal mapping functionality.
 */
@Mapper
@SuppressWarnings("PMD.ConstantsInInterface")
public interface AirGapMapper {

    /**
     * Singleton instance of the mapper.
     */
    AirGapMapper INSTANCE = Mappers.getMapper(AirGapMapper.class);
}
