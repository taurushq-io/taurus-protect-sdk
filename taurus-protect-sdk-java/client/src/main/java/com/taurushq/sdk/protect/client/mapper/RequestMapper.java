package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestApprovers;
import com.taurushq.sdk.protect.client.model.RequestApproversGroup;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.client.model.RequestParallelApproversGroups;
import com.taurushq.sdk.protect.client.model.RequestTrail;
import com.taurushq.sdk.protect.client.model.SignedRequest;
import com.taurushq.sdk.protect.openapi.model.RequestSignedRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproversGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordParallelApproversGroups;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestTrail;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;


/**
 * The interface Request mapper.
 */
@Mapper(uses = {BalanceMapper.class, AttributeMapper.class, CurrencyMapper.class})
public interface RequestMapper {
    /**
     * The constant INSTANCE.
     */
    RequestMapper INSTANCE = Mappers.getMapper(RequestMapper.class);

    /**
     * From dto request.
     *
     * @param request the request
     * @return the request
     */
    @Mapping(target = "status", expression = "java(com.taurushq.sdk.protect.client.model.RequestStatus.valueOfLabel(request.getStatus()))")
    Request fromDTO(TgvalidatordRequest request);

    /**
     * From dto signed request.
     *
     * @param signedRequest the signed request
     * @return the signed request
     */
    @Mapping(target = "status", expression = "java(com.taurushq.sdk.protect.client.model.RequestStatus.valueOfLabel(signedRequest.getStatus()))")
    SignedRequest fromDTO(RequestSignedRequest signedRequest);

    /**
     * From dto request trail.
     *
     * @param trail the trail
     * @return the request trail
     */
    RequestTrail fromDTO(TgvalidatordRequestTrail trail);

    /**
     * From dto request metadata.
     *
     * @param metadata the metadata
     * @return the request metadata
     */
    RequestMetadata fromDTO(TgvalidatordMetadata metadata);

    /**
     * From dto request approvers.
     *
     * @param approvers the approvers
     * @return the request approvers
     */
    RequestApprovers fromDTO(TgvalidatordApprovers approvers);

    /**
     * From dto request parallel approvers groups.
     *
     * @param parallelApproversGroups the parallel approvers groups
     * @return the request parallel approvers groups
     */
    RequestParallelApproversGroups fromDTO(TgvalidatordParallelApproversGroups parallelApproversGroups);

    /**
     * From dto request approvers group.
     *
     * @param approversGroup the approvers group
     * @return the request approvers group
     */
    RequestApproversGroup fromDTO(TgvalidatordApproversGroup approversGroup);

}


