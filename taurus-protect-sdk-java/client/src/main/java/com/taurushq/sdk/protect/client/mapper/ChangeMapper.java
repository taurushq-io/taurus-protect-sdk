package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Change;
import com.taurushq.sdk.protect.client.model.CreateChangeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordChange;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateChangeRequest;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface Change mapper.
 */
@Mapper
public interface ChangeMapper {
    /**
     * The constant INSTANCE.
     */
    ChangeMapper INSTANCE = Mappers.getMapper(ChangeMapper.class);

    /**
     * From dto change.
     *
     * @param change the change
     * @return the change
     */
    @Mapping(source = "creationDate", target = "createdAt")
    @Mapping(target = "tenantId", expression = "java(change.getTenantId() != null ? Integer.parseInt(change.getTenantId()) : 0)")
    Change fromDTO(TgvalidatordChange change);

    /**
     * From dto list.
     *
     * @param changes the changes
     * @return the list
     */
    List<Change> fromDTO(List<TgvalidatordChange> changes);

    /**
     * Converts an SDK CreateChangeRequest to the OpenAPI DTO.
     *
     * @param request the SDK create change request
     * @return the OpenAPI DTO
     */
    @Mapping(source = "comment", target = "changeComment")
    TgvalidatordCreateChangeRequest toDTO(CreateChangeRequest request);
}
