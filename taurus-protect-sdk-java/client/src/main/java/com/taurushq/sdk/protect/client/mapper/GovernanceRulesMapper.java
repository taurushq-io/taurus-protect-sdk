package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.GovernanceRulesTrail;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SuperAdminPublicKey;
import com.taurushq.sdk.protect.openapi.model.GetPublicKeysReplyPublicKey;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRuleUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRules;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRulesTrail;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for governance rules related DTOs.
 */
@Mapper
public interface GovernanceRulesMapper {

    /**
     * The constant INSTANCE.
     */
    GovernanceRulesMapper INSTANCE = Mappers.getMapper(GovernanceRulesMapper.class);

    /**
     * Converts a TgvalidatordRuleUserSignature DTO to RuleUserSignature model.
     *
     * @param dto the DTO
     * @return the model
     */
    RuleUserSignature fromDTO(TgvalidatordRuleUserSignature dto);

    /**
     * Converts a list of TgvalidatordRuleUserSignature DTOs to RuleUserSignature models.
     *
     * @param dtos the DTOs
     * @return the models
     */
    List<RuleUserSignature> fromRuleUserSignatureDTOs(List<TgvalidatordRuleUserSignature> dtos);

    /**
     * Converts a TgvalidatordRulesTrail DTO to GovernanceRulesTrail model.
     *
     * @param dto the DTO
     * @return the model
     */
    GovernanceRulesTrail fromDTO(TgvalidatordRulesTrail dto);

    /**
     * Converts a list of TgvalidatordRulesTrail DTOs to GovernanceRulesTrail models.
     *
     * @param dtos the DTOs
     * @return the models
     */
    List<GovernanceRulesTrail> fromRulesTrailDTOs(List<TgvalidatordRulesTrail> dtos);

    /**
     * Converts a TgvalidatordRules DTO to GovernanceRules model.
     *
     * @param dto the DTO
     * @return the model
     */
    GovernanceRules fromDTO(TgvalidatordRules dto) throws ApiException;

    /**
     * Converts a list of TgvalidatordRules DTOs to GovernanceRules models.
     *
     * @param dtos the DTOs
     * @return the models
     */
    List<GovernanceRules> fromRulesDTOs(List<TgvalidatordRules> dtos);

    /**
     * Converts a GetPublicKeysReplyPublicKey DTO to SuperAdminPublicKey model.
     *
     * @param dto the DTO
     * @return the model
     */
    @Mapping(source = "userID", target = "userId")
    SuperAdminPublicKey fromDTO(GetPublicKeysReplyPublicKey dto);

    /**
     * Converts a list of GetPublicKeysReplyPublicKey DTOs to SuperAdminPublicKey models.
     *
     * @param dtos the DTOs
     * @return the models
     */
    List<SuperAdminPublicKey> fromPublicKeyDTOs(List<GetPublicKeysReplyPublicKey> dtos);
}
