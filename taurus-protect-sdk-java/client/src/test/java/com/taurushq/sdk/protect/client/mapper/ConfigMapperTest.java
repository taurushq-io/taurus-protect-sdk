package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.TenantConfig;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTenantConfig;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.assertFalse;

class ConfigMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordTenantConfig dto = new TgvalidatordTenantConfig();
        dto.setTenantId("tenant-123");
        dto.setSuperAdminMinimumSignatures("2");
        dto.setBaseCurrency("USD");
        dto.setIsMFAMandatory(true);
        dto.setExcludeContainer(false);
        dto.setFeeLimitFactor(1.5f);
        dto.setProtectEngineVersion("2.0.0");
        dto.setRestrictSourcesForWhitelistedAddresses(true);
        dto.setIsProtectEngineCold(false);
        dto.setIsColdProtectEngineOffline(false);
        dto.setIsPhysicalAirGapEnabled(true);

        TenantConfig result = ConfigMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("tenant-123", result.getTenantId());
        assertEquals("2", result.getSuperAdminMinimumSignatures());
        assertEquals("USD", result.getBaseCurrency());
        assertTrue(result.getMFAMandatory());
        assertFalse(result.getExcludeContainer());
        assertEquals(1.5f, result.getFeeLimitFactor());
        assertEquals("2.0.0", result.getProtectEngineVersion());
        assertTrue(result.getRestrictSourcesForWhitelistedAddresses());
        assertFalse(result.getProtectEngineCold());
        assertFalse(result.getColdProtectEngineOffline());
        assertTrue(result.getPhysicalAirGapEnabled());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordTenantConfig dto = new TgvalidatordTenantConfig();
        dto.setTenantId("minimal-tenant");

        TenantConfig result = ConfigMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("minimal-tenant", result.getTenantId());
        assertNull(result.getSuperAdminMinimumSignatures());
        assertNull(result.getBaseCurrency());
        assertNull(result.getMFAMandatory());
        assertNull(result.getExcludeContainer());
        assertNull(result.getFeeLimitFactor());
        assertNull(result.getProtectEngineVersion());
        assertNull(result.getRestrictSourcesForWhitelistedAddresses());
        assertNull(result.getProtectEngineCold());
        assertNull(result.getColdProtectEngineOffline());
        assertNull(result.getPhysicalAirGapEnabled());
    }

    @Test
    void fromDTO_handlesNullDto() {
        TenantConfig result = ConfigMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromDTO_mapsBooleanFieldsCorrectly() {
        TgvalidatordTenantConfig dto = new TgvalidatordTenantConfig();
        dto.setTenantId("bool-test");
        dto.setIsMFAMandatory(false);
        dto.setIsProtectEngineCold(true);
        dto.setIsColdProtectEngineOffline(true);
        dto.setIsPhysicalAirGapEnabled(false);

        TenantConfig result = ConfigMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertFalse(result.getMFAMandatory());
        assertTrue(result.getProtectEngineCold());
        assertTrue(result.getColdProtectEngineOffline());
        assertFalse(result.getPhysicalAirGapEnabled());
    }
}
