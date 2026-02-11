package com.taurushq.sdk.protect.client.model;

import com.taurushq.sdk.protect.client.mapper.WalletMapper;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class WalletMapperTest {

    @Test
    void fromDTO() {

        TgvalidatordWalletInfo walletInfo = new TgvalidatordWalletInfo();
        walletInfo.setId("10");
        walletInfo.setName("name");
        walletInfo.setCurrency("currency");
        Wallet w = WalletMapper.INSTANCE.fromDTO(walletInfo);

        assertEquals(10, w.getId());
        assertEquals("name", w.getName());
        assertEquals("currency", w.getCurrency());

    }
}