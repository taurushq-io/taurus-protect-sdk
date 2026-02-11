package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedRequestAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletAttribute;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AttributeMapperTest {

    @Test
    void fromAddressAttribute_mapsAllFields() {
        TgvalidatordAddressAttribute dto = new TgvalidatordAddressAttribute();
        dto.setId("1");
        dto.setKey("department");
        dto.setValue("treasury");
        dto.setContentType("text/plain");
        dto.setOwner("admin");
        dto.setType("custom");
        dto.setSubtype("label");
        dto.setIsfile(true);

        Attribute result = AttributeMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertEquals("department", result.getKey());
        assertEquals("treasury", result.getValue());
        assertEquals("text/plain", result.getContentType());
        assertEquals("admin", result.getOwner());
        assertEquals("custom", result.getType());
        assertEquals("label", result.getSubType());
        assertTrue(result.isFile());
    }

    @Test
    void fromAddressAttribute_handlesNullDto() {
        Attribute result = AttributeMapper.INSTANCE.fromDTO((TgvalidatordAddressAttribute) null);
        assertNull(result);
    }

    @Test
    void fromWalletAttribute_mapsAllFields() {
        TgvalidatordWalletAttribute dto = new TgvalidatordWalletAttribute();
        dto.setId("2");
        dto.setKey("team");
        dto.setValue("engineering");
        dto.setContentType("application/json");
        dto.setOwner("system");
        dto.setType("metadata");
        dto.setSubtype("info");
        dto.setIsfile(false);

        Attribute result = AttributeMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(2L, result.getId());
        assertEquals("team", result.getKey());
        assertEquals("engineering", result.getValue());
        assertEquals("metadata", result.getType());
        assertEquals("info", result.getSubType());
    }

    @Test
    void fromWalletAttribute_handlesNullDto() {
        Attribute result = AttributeMapper.INSTANCE.fromDTO((TgvalidatordWalletAttribute) null);
        assertNull(result);
    }

    @Test
    void fromRequestAttribute_mapsAvailableFields() {
        TgvalidatordRequestAttribute dto = new TgvalidatordRequestAttribute();
        dto.setId("3");
        dto.setKey("note");
        dto.setValue("payment for services");
        dto.setContentType("text/plain");

        Attribute result = AttributeMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(3L, result.getId());
        assertEquals("note", result.getKey());
        assertEquals("payment for services", result.getValue());
        // subType, isFile, owner, type are ignored mappings for request attributes
        assertNull(result.getSubType());
        assertNull(result.getOwner());
        assertNull(result.getType());
    }

    @Test
    void fromRequestAttribute_handlesNullDto() {
        Attribute result = AttributeMapper.INSTANCE.fromDTO((TgvalidatordRequestAttribute) null);
        assertNull(result);
    }

    @Test
    void fromSignedRequestAttribute_mapsAvailableFields() {
        TgvalidatordSignedRequestAttribute dto = new TgvalidatordSignedRequestAttribute();
        dto.setId("4");
        dto.setKey("reference");
        dto.setValue("INV-2024-001");

        Attribute result = AttributeMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(4L, result.getId());
        assertEquals("reference", result.getKey());
        assertEquals("INV-2024-001", result.getValue());
    }

    @Test
    void fromSignedRequestAttribute_handlesNullDto() {
        Attribute result = AttributeMapper.INSTANCE.fromDTO((TgvalidatordSignedRequestAttribute) null);
        assertNull(result);
    }
}
