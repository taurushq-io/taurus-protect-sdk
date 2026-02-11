package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.AssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAsset;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AssetBalanceMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordAsset asset = new TgvalidatordAsset();
        asset.setCurrency("ETH");
        asset.setKind("native");

        TgvalidatordBalance balance = new TgvalidatordBalance();
        balance.setTotalConfirmed("1000");
        balance.setAvailableConfirmed("800");

        TgvalidatordAssetBalance dto = new TgvalidatordAssetBalance();
        dto.setAsset(asset);
        dto.setBalance(balance);

        AssetBalance result = AssetBalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getAsset());
        assertNotNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordAssetBalance dto = new TgvalidatordAssetBalance();

        AssetBalance result = AssetBalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNull(result.getAsset());
        assertNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullDto() {
        AssetBalance result = AssetBalanceMapper.INSTANCE.fromDTO((TgvalidatordAssetBalance) null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordAssetBalance dto1 = new TgvalidatordAssetBalance();
        dto1.setAsset(new TgvalidatordAsset());

        TgvalidatordAssetBalance dto2 = new TgvalidatordAssetBalance();
        dto2.setAsset(new TgvalidatordAsset());

        List<AssetBalance> result = AssetBalanceMapper.INSTANCE.fromDTO(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<AssetBalance> result = AssetBalanceMapper.INSTANCE.fromDTO(
                Collections.<TgvalidatordAssetBalance>emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<AssetBalance> result = AssetBalanceMapper.INSTANCE.fromDTO(
                (List<TgvalidatordAssetBalance>) null);
        assertNull(result);
    }
}
