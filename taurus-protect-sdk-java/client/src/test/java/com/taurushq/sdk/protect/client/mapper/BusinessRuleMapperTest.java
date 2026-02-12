package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BusinessRule;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBusinessRule;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class BusinessRuleMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordBusinessRule dto = new TgvalidatordBusinessRule();
        dto.setId("rule-123");
        dto.setTenantId("42");
        dto.setCurrency("ETH");
        dto.setWalletId("wallet-456");
        dto.setAddressId("address-789");
        dto.setRuleKey("MAX_AMOUNT");
        dto.setRuleValue("1000000");
        dto.setRuleGroup("LIMITS");
        dto.setRuleDescription("Maximum transfer amount");
        dto.setRuleValidation("amount <= 1000000");
        dto.setEntityType("WALLET");
        dto.setEntityID("entity-111");

        // When
        BusinessRule result = BusinessRuleMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("rule-123", result.getId());
        assertEquals(42, result.getTenantId());
        assertEquals("ETH", result.getCurrency());
        assertEquals("wallet-456", result.getWalletId());
        assertEquals("address-789", result.getAddressId());
        assertEquals("MAX_AMOUNT", result.getRuleKey());
        assertEquals("1000000", result.getRuleValue());
        assertEquals("LIMITS", result.getRuleGroup());
        assertEquals("Maximum transfer amount", result.getRuleDescription());
        assertEquals("amount <= 1000000", result.getRuleValidation());
        assertEquals("WALLET", result.getEntityType());
        assertEquals("entity-111", result.getEntityID());
    }

    @Test
    void fromDTO_withNullTenantId_defaultsToZero() {
        // Given
        TgvalidatordBusinessRule dto = new TgvalidatordBusinessRule();
        dto.setId("rule-1");
        dto.setTenantId(null);

        // When
        BusinessRule result = BusinessRuleMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(0, result.getTenantId());
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordBusinessRule dto = new TgvalidatordBusinessRule();
        dto.setId("rule-1");
        // All other fields left null

        // When
        BusinessRule result = BusinessRuleMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("rule-1", result.getId());
        assertNull(result.getCurrency());
        assertNull(result.getRuleKey());
        assertNull(result.getRuleValue());
    }

    @Test
    void fromDTO_list_mapsMultipleRules() {
        // Given
        TgvalidatordBusinessRule dto1 = new TgvalidatordBusinessRule();
        dto1.setId("rule-1");
        dto1.setRuleKey("MAX_AMOUNT");

        TgvalidatordBusinessRule dto2 = new TgvalidatordBusinessRule();
        dto2.setId("rule-2");
        dto2.setRuleKey("MIN_AMOUNT");

        List<TgvalidatordBusinessRule> dtos = Arrays.asList(dto1, dto2);

        // When
        List<BusinessRule> results = BusinessRuleMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("rule-1", results.get(0).getId());
        assertEquals("MAX_AMOUNT", results.get(0).getRuleKey());
        assertEquals("rule-2", results.get(1).getId());
        assertEquals("MIN_AMOUNT", results.get(1).getRuleKey());
    }

    @Test
    void fromDTO_withCurrencyInfo_mapsCurrency() {
        // Given
        TgvalidatordBusinessRule dto = new TgvalidatordBusinessRule();
        dto.setId("rule-1");

        TgvalidatordCurrency currencyInfo = new TgvalidatordCurrency();
        currencyInfo.setId("ETH");
        currencyInfo.setName("Ethereum");
        dto.setCurrencyInfo(currencyInfo);

        // When
        BusinessRule result = BusinessRuleMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getCurrencyInfo());
        assertEquals("ETH", result.getCurrencyInfo().getId());
        assertEquals("Ethereum", result.getCurrencyInfo().getName());
    }
}
