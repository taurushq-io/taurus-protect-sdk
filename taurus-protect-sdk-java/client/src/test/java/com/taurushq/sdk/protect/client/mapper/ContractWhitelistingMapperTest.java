package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

/**
 * Attribute mapper tests.
 * <p>
 * The envelope, signed-contract and result mappers these tests used to cover built domain
 * objects from the DTO with no verification, so their assertions pinned the bypass behaving
 * consistently. Verified reads are covered by the whitelisted-asset suites.
 */
class ContractWhitelistingMapperTest {

    @Test
    void fromAttributeDTO_mapsAllFields() {
        TgvalidatordWhitelistedContractAddressAttribute dto =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto.setId("123");
        dto.setKey("description");
        dto.setValue("A test token");
        dto.setContentType("text/plain");
        dto.setOwner("user-123");
        dto.setType("metadata");
        dto.setSubtype("basic");
        dto.setIsfile(false);

        Attribute attr = ContractWhitelistingMapper.INSTANCE.fromAttributeDTO(dto);

        assertEquals("description", attr.getKey());
        assertEquals("A test token", attr.getValue());
        assertEquals("text/plain", attr.getContentType());
        assertEquals("user-123", attr.getOwner());
        assertEquals("metadata", attr.getType());
        assertEquals("basic", attr.getSubType());
        assertFalse(attr.isFile());
    }

    @Test
    void fromAttributeDTOList_mapsList() {
        TgvalidatordWhitelistedContractAddressAttribute dto1 =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto1.setId("1");
        dto1.setKey("key1");

        TgvalidatordWhitelistedContractAddressAttribute dto2 =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto2.setId("2");
        dto2.setKey("key2");

        List<Attribute> attrs = ContractWhitelistingMapper.INSTANCE
                .fromAttributeDTOList(Arrays.asList(dto1, dto2));

        assertEquals(2, attrs.size());
        assertEquals("key1", attrs.get(0).getKey());
        assertEquals("key2", attrs.get(1).getKey());
    }
}
