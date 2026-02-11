package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.NFTCollectionBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordNFTCollectionBalance;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class NFTCollectionBalanceMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordCurrency currencyInfo = new TgvalidatordCurrency();
        currencyInfo.setId("BAYC");
        currencyInfo.setName("Bored Ape Yacht Club");
        currencyInfo.setBlockchain("ETH");

        TgvalidatordBalance balance = new TgvalidatordBalance();
        balance.setTotalConfirmed("10");
        balance.setAvailableConfirmed("10");

        TgvalidatordNFTCollectionBalance dto = new TgvalidatordNFTCollectionBalance();
        dto.setCurrencyInfo(currencyInfo);
        dto.setBalance(balance);

        NFTCollectionBalance result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getCurrencyInfo());
        assertEquals("BAYC", result.getCurrencyInfo().getId());
        assertEquals("Bored Ape Yacht Club", result.getCurrencyInfo().getName());
        assertNotNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordNFTCollectionBalance dto = new TgvalidatordNFTCollectionBalance();

        NFTCollectionBalance result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNull(result.getCurrencyInfo());
        assertNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullDto() {
        NFTCollectionBalance result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(
                (TgvalidatordNFTCollectionBalance) null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordNFTCollectionBalance dto1 = new TgvalidatordNFTCollectionBalance();
        TgvalidatordCurrency curr1 = new TgvalidatordCurrency();
        curr1.setId("BAYC");
        dto1.setCurrencyInfo(curr1);

        TgvalidatordNFTCollectionBalance dto2 = new TgvalidatordNFTCollectionBalance();
        TgvalidatordCurrency curr2 = new TgvalidatordCurrency();
        curr2.setId("PUNK");
        dto2.setCurrencyInfo(curr2);

        List<NFTCollectionBalance> result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("BAYC", result.get(0).getCurrencyInfo().getId());
        assertEquals("PUNK", result.get(1).getCurrencyInfo().getId());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<NFTCollectionBalance> result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(
                Collections.<TgvalidatordNFTCollectionBalance>emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<NFTCollectionBalance> result = NFTCollectionBalanceMapper.INSTANCE.fromDTO(
                (List<TgvalidatordNFTCollectionBalance>) null);
        assertNull(result);
    }
}
