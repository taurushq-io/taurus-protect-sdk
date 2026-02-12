package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.AddressInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;

import static org.junit.jupiter.api.Assertions.*;

class AddressMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordAddress dto = new TgvalidatordAddress();
        dto.setId("12345");
        dto.setWalletId("67890");
        dto.setDisabled(false);
        dto.setCurrency("ETH");
        dto.setAddressPath("m/44'/60'/0'/0/0");
        dto.setAddress("0x1234567890abcdef");
        dto.setComment("Test comment");
        dto.setLabel("My Address");
        dto.setSignature("sig123");

        OffsetDateTime now = OffsetDateTime.now();
        dto.setCreationDate(now);
        dto.setUpdateDate(now);

        // When
        Address result = AddressMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(12345L, result.getId());
        assertEquals(67890L, result.getWalletId());
        assertFalse(result.isDisabled());
        assertEquals("ETH", result.getCurrency());
        assertEquals("m/44'/60'/0'/0/0", result.getAddressPath());
        assertEquals("0x1234567890abcdef", result.getAddress());
        assertEquals("Test comment", result.getComment());
        assertEquals("My Address", result.getLabel());
        assertEquals("sig123", result.getSignature());
        assertEquals(now, result.getCreationDate());
        assertEquals(now, result.getUpdateDate());
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordAddress dto = new TgvalidatordAddress();
        dto.setId("1");
        dto.setWalletId("2");
        // All other fields left null

        // When
        Address result = AddressMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertEquals(2L, result.getWalletId());
        assertNull(result.getCurrency());
        assertNull(result.getAddress());
        assertNull(result.getLabel());
    }

    @Test
    void fromDTO_withBalance_mapsBalance() {
        // Given
        TgvalidatordAddress dto = new TgvalidatordAddress();
        dto.setId("1");
        dto.setWalletId("2");

        TgvalidatordBalance balance = new TgvalidatordBalance();
        balance.setAvailableConfirmed("1000000000");
        dto.setBalance(balance);

        // When
        Address result = AddressMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getBalance());
    }

    @Test
    void fromDTO_addressInfo_mapsCorrectly() {
        // Given
        TgvalidatordAddressInfo dto = new TgvalidatordAddressInfo();
        dto.setAddress("0xabcdef123456");

        // When
        AddressInfo result = AddressMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("0xabcdef123456", result.getAddress());
    }

    @Test
    void fromDTO_withDisabledTrue_mapsCorrectly() {
        // Given
        TgvalidatordAddress dto = new TgvalidatordAddress();
        dto.setId("1");
        dto.setWalletId("2");
        dto.setDisabled(true);

        // When
        Address result = AddressMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertTrue(result.isDisabled());
    }
}
