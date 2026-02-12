package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedRequestAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletAttribute;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * The interface Attribute mapper.
 */
@Mapper
public interface AttributeMapper {
    /**
     * The constant INSTANCE.
     */
    AttributeMapper INSTANCE = Mappers.getMapper(AttributeMapper.class);

    /**
     * From dto attribute.
     *
     * @param attribute the attribute
     * @return the attribute
     */
    @Mapping(source = "subtype", target = "subType")
    @Mapping(source = "isfile", target = "isFile")
    Attribute fromDTO(TgvalidatordAddressAttribute attribute);

    /**
     * From dto attribute.
     *
     * @param attribute the attribute
     * @return the attribute
     */
    @Mapping(source = "subtype", target = "subType")
    @Mapping(source = "isfile", target = "isFile")
    Attribute fromDTO(TgvalidatordWalletAttribute attribute);

    /**
     * From dto attribute.
     *
     * @param attribute the attribute
     * @return the attribute
     */
    @Mapping(target = "subType", ignore = true)
    @Mapping(target = "isFile", ignore = true)
    @Mapping(target = "owner", ignore = true)
    @Mapping(target = "type", ignore = true)
    Attribute fromDTO(TgvalidatordRequestAttribute attribute);

    /**
     * From dto attribute.
     *
     * @param attribute the attribute
     * @return the attribute
     */
    @Mapping(target = "subType", ignore = true)
    @Mapping(target = "isFile", ignore = true)
    @Mapping(target = "owner", ignore = true)
    @Mapping(target = "type", ignore = true)
    Attribute fromDTO(TgvalidatordSignedRequestAttribute attribute);
}


