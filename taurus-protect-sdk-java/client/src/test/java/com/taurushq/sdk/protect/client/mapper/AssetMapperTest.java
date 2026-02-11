package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.Asset;
import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAsset;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AssetMapperTest {

    @Test
    void fromDTO_mapsAssetFields() {
        TgvalidatordAsset dto = new TgvalidatordAsset();
        dto.setCurrency("ETH");
        dto.setKind("native");

        Asset result = AssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("ETH", result.getCurrency());
        assertEquals("native", result.getKind());
    }

    @Test
    void fromDTO_handlesNullAsset() {
        Asset result = AssetMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromAddressDTOList_mapsList() {
        TgvalidatordAddress dto1 = new TgvalidatordAddress();
        dto1.setId("1");
        dto1.setAddress("0x123");
        dto1.setCurrency("ETH");

        TgvalidatordAddress dto2 = new TgvalidatordAddress();
        dto2.setId("2");
        dto2.setAddress("0x456");
        dto2.setCurrency("ETH");

        List<Address> result = AssetMapper.INSTANCE.fromAddressDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("0x123", result.get(0).getAddress());
        assertEquals("0x456", result.get(1).getAddress());
    }

    @Test
    void fromAddressDTOList_handlesEmptyList() {
        List<Address> result = AssetMapper.INSTANCE.fromAddressDTOList(Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromAddressDTOList_handlesNullList() {
        List<Address> result = AssetMapper.INSTANCE.fromAddressDTOList(null);
        assertNull(result);
    }

    @Test
    void fromWalletInfoDTOList_mapsList() {
        TgvalidatordWalletInfo dto1 = new TgvalidatordWalletInfo();
        dto1.setId("1");
        dto1.setName("Wallet 1");
        dto1.setCurrency("ETH");

        TgvalidatordWalletInfo dto2 = new TgvalidatordWalletInfo();
        dto2.setId("2");
        dto2.setName("Wallet 2");
        dto2.setCurrency("BTC");

        List<Wallet> result = AssetMapper.INSTANCE.fromWalletInfoDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("Wallet 1", result.get(0).getName());
        assertEquals("Wallet 2", result.get(1).getName());
    }

    @Test
    void fromWalletInfoDTOList_handlesEmptyList() {
        List<Wallet> result = AssetMapper.INSTANCE.fromWalletInfoDTOList(Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromWalletInfoDTOList_handlesNullList() {
        List<Wallet> result = AssetMapper.INSTANCE.fromWalletInfoDTOList(null);
        assertNull(result);
    }
}
