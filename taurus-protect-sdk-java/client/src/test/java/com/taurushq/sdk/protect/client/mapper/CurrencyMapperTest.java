package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Currency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class CurrencyMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordCurrency dto = new TgvalidatordCurrency();
        dto.setId("ETH");
        dto.setName("Ethereum");
        dto.setSymbol("ETH");
        dto.setCoinTypeIndex("60");
        dto.setBlockchain("ETH");
        dto.setIsToken(false);
        dto.setIsERC20(false);
        dto.setDecimals("18");
        dto.setContractAddress(null);
        dto.setHasStaking(true);
        dto.setIsUTXOBased(false);
        dto.setIsAccountBased(true);
        dto.setIsFiat(false);
        dto.setIsFA12(false);
        dto.setIsFA20(false);
        dto.setIsNFT(false);
        dto.setEnabled(true);
        dto.setDisplayName("Ethereum (ETH)");
        dto.setType("NATIVE");

        Currency result = CurrencyMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("ETH", result.getId());
        assertEquals("Ethereum", result.getName());
        assertEquals("ETH", result.getSymbol());
        assertEquals("60", result.getCoinTypeIndex());
        assertEquals("ETH", result.getBlockchain());
        assertFalse(result.isToken());
        assertFalse(result.isERC20());
        assertEquals(18, result.getDecimals());
        assertNull(result.getContractAddress());
        assertTrue(result.hasStaking());
        assertFalse(result.isUTXOBased());
        assertTrue(result.isAccountBased());
        assertFalse(result.isFiat());
        assertFalse(result.isFA12());
        assertFalse(result.isFA20());
        assertFalse(result.isNFT());
        assertTrue(result.isEnabled());
        assertEquals("Ethereum (ETH)", result.getDisplayName());
        assertEquals("NATIVE", result.getType());
    }

    @Test
    void fromDTO_mapsERC20Token() {
        TgvalidatordCurrency dto = new TgvalidatordCurrency();
        dto.setId("USDC");
        dto.setName("USD Coin");
        dto.setSymbol("USDC");
        dto.setBlockchain("ETH");
        dto.setIsToken(true);
        dto.setIsERC20(true);
        dto.setDecimals("6");
        dto.setContractAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48");
        dto.setIsUTXOBased(false);
        dto.setIsAccountBased(true);
        dto.setIsFiat(false);
        dto.setIsNFT(false);

        Currency result = CurrencyMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertTrue(result.isToken());
        assertTrue(result.isERC20());
        assertEquals(6, result.getDecimals());
        assertEquals("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", result.getContractAddress());
    }

    @Test
    void fromDTO_mapsBooleanFieldsCorrectly() {
        TgvalidatordCurrency dto = new TgvalidatordCurrency();
        dto.setId("XTZ_FA12");
        dto.setIsToken(true);
        dto.setIsFA12(true);
        dto.setIsFA20(false);
        dto.setIsNFT(false);
        dto.setIsERC20(false);
        dto.setIsUTXOBased(false);
        dto.setIsAccountBased(true);
        dto.setIsFiat(false);

        Currency result = CurrencyMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertTrue(result.isToken());
        assertTrue(result.isFA12());
        assertFalse(result.isFA20());
        assertFalse(result.isNFT());
    }

    @Test
    void fromDTO_handlesNullBooleans() {
        TgvalidatordCurrency dto = new TgvalidatordCurrency();
        dto.setId("BTC");

        Currency result = CurrencyMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("BTC", result.getId());
        assertFalse(result.isToken());
        assertFalse(result.isERC20());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Currency result = CurrencyMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }
}
