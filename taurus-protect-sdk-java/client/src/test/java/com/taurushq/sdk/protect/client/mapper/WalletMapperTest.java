package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class WalletMapperTest {

    @Test
    void fromWalletInfoDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = OffsetDateTime.now().plusHours(1);

        TgvalidatordBalance balance = new TgvalidatordBalance();
        balance.setTotalConfirmed("1000");

        TgvalidatordCurrency currencyInfo = new TgvalidatordCurrency();
        currencyInfo.setId("ETH");

        TgvalidatordWalletInfo dto = new TgvalidatordWalletInfo();
        dto.setId("123");
        dto.setName("My Wallet");
        dto.setCurrency("ETH");
        dto.setBalance(balance);
        dto.setAccountPath("m/44'/60'/0'");
        dto.setIsOmnibus(true);
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setCustomerId("cust-456");
        dto.setComment("Test wallet");
        dto.setDisabled(false);
        dto.setBlockchain("ETH");
        dto.setAddressesCount("5");
        dto.setCurrencyInfo(currencyInfo);
        dto.setNetwork("mainnet");
        dto.setVisibilityGroupID("vg-1");
        dto.setExternalWalletId("ext-123");

        Wallet result = WalletMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(123L, result.getId());
        assertEquals("My Wallet", result.getName());
        assertEquals("ETH", result.getCurrency());
        assertNotNull(result.getBalance());
        assertEquals("m/44'/60'/0'", result.getAccountPath());
        assertTrue(result.isOmnibus());
        assertEquals(creationDate, result.getCreationDate());
        assertEquals(updateDate, result.getUpdateDate());
        assertEquals("cust-456", result.getCustomerId());
        assertEquals("Test wallet", result.getComment());
        assertFalse(result.isDisabled());
        assertEquals("ETH", result.getBlockchain());
        assertEquals(5, result.getAddressesCount());
        assertNotNull(result.getCurrencyInfo());
        assertEquals("mainnet", result.getNetwork());
        assertEquals("vg-1", result.getVisibilityGroupID());
        assertEquals("ext-123", result.getExternalWalletId());
    }

    @Test
    void fromWalletInfoDTO_handlesNullFields() {
        TgvalidatordWalletInfo dto = new TgvalidatordWalletInfo();
        dto.setId("1");
        dto.setName("W");

        Wallet result = WalletMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertEquals("W", result.getName());
        assertNull(result.getCurrency());
        assertNull(result.getBalance());
        assertNull(result.getAccountPath());
        assertNull(result.getCurrencyInfo());
    }

    @Test
    void fromWalletInfoDTO_handlesNullDto() {
        Wallet result = WalletMapper.INSTANCE.fromDTO((TgvalidatordWalletInfo) null);
        assertNull(result);
    }

    @Test
    void fromWalletDTO_mapsFields() {
        TgvalidatordWallet dto = new TgvalidatordWallet();
        dto.setId("456");
        dto.setName("BTC Wallet");
        dto.setCurrency("BTC");
        dto.setIsOmnibus(false);
        dto.setDisabled(true);

        Wallet result = WalletMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(456L, result.getId());
        assertEquals("BTC Wallet", result.getName());
        assertEquals("BTC", result.getCurrency());
        assertFalse(result.isOmnibus());
        assertTrue(result.isDisabled());
        // network and visibilityGroupID are ignored for TgvalidatordWallet
        assertNull(result.getNetwork());
        assertNull(result.getVisibilityGroupID());
    }

    @Test
    void fromWalletDTO_handlesNullDto() {
        Wallet result = WalletMapper.INSTANCE.fromDTO((TgvalidatordWallet) null);
        assertNull(result);
    }

    @Test
    void fromWalletInfoDTO_mapsOmnibusTrue() {
        TgvalidatordWalletInfo dto = new TgvalidatordWalletInfo();
        dto.setId("1");
        dto.setIsOmnibus(true);

        Wallet result = WalletMapper.INSTANCE.fromDTO(dto);

        assertTrue(result.isOmnibus());
    }

    @Test
    void fromWalletInfoDTO_mapsOmnibusFalse() {
        TgvalidatordWalletInfo dto = new TgvalidatordWalletInfo();
        dto.setId("1");
        dto.setIsOmnibus(false);

        Wallet result = WalletMapper.INSTANCE.fromDTO(dto);

        assertFalse(result.isOmnibus());
    }
}
