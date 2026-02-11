package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BusinessRule;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBusinessRule;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface Business rule mapper.
 */
@Mapper(uses = {CurrencyMapper.class})
public interface BusinessRuleMapper {
    /**
     * The constant INSTANCE.
     */
    BusinessRuleMapper INSTANCE = Mappers.getMapper(BusinessRuleMapper.class);

    /**
     * From dto business rule.
     *
     * @param businessRule the business rule
     * @return the business rule
     */
    @Mapping(target = "tenantId", expression = "java(businessRule.getTenantId() != null ? Integer.parseInt(businessRule.getTenantId()) : 0)")
    BusinessRule fromDTO(TgvalidatordBusinessRule businessRule);

    /**
     * From dto list.
     *
     * @param businessRules the business rules
     * @return the list
     */
    List<BusinessRule> fromDTO(List<TgvalidatordBusinessRule> businessRules);
}
