package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.client.model.GroupUser;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalGroupUser;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting between OpenAPI group DTOs and client model objects.
 */
@Mapper
public interface GroupMapper {

    /**
     * Singleton instance of the mapper.
     */
    GroupMapper INSTANCE = Mappers.getMapper(GroupMapper.class);

    /**
     * Converts an OpenAPI group DTO to a client model Group.
     *
     * @param dto the OpenAPI DTO
     * @return the client model Group
     */
    Group fromDTO(TgvalidatordInternalGroup dto);

    /**
     * Converts a list of OpenAPI group DTOs to client model Groups.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of client model Groups
     */
    List<Group> fromDTOList(List<TgvalidatordInternalGroup> dtos);

    /**
     * Converts an OpenAPI group user DTO to a client model GroupUser.
     *
     * @param dto the OpenAPI DTO
     * @return the client model GroupUser
     */
    GroupUser fromUserDTO(TgvalidatordInternalGroupUser dto);

    /**
     * Converts a list of OpenAPI group user DTOs to client model GroupUsers.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of client model GroupUsers
     */
    List<GroupUser> fromUserDTOList(List<TgvalidatordInternalGroupUser> dtos);
}
