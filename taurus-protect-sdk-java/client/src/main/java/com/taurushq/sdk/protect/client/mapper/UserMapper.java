package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalUser;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface User mapper.
 */
@Mapper
public interface UserMapper {
    /**
     * The constant INSTANCE.
     */
    UserMapper INSTANCE = Mappers.getMapper(UserMapper.class);

    /**
     * From dto user.
     *
     * @param user the user
     * @return the user
     */
    @Mapping(source = "creationDate", target = "createdAt")
    @Mapping(source = "updateDate", target = "updatedAt")
    @Mapping(target = "tenantId", expression = "java(user.getTenantId() != null ? Integer.parseInt(user.getTenantId()) : 0)")
    User fromDTO(TgvalidatordInternalUser user);

    /**
     * From dto list.
     *
     * @param users the users
     * @return the list
     */
    List<User> fromDTO(List<TgvalidatordInternalUser> users);
}
