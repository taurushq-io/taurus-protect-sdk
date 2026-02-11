package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Balance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class BalanceMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordBalance dto = new TgvalidatordBalance();
        dto.setTotalConfirmed("1000000000000000000");
        dto.setTotalUnconfirmed("500000000000000000");
        dto.setAvailableConfirmed("800000000000000000");
        dto.setAvailableUnconfirmed("300000000000000000");
        dto.setReservedConfirmed("200000000000000000");
        dto.setReservedUnconfirmed("200000000000000000");

        Balance result = BalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(new BigInteger("1000000000000000000"), result.getTotalConfirmed());
        assertEquals(new BigInteger("500000000000000000"), result.getTotalUnconfirmed());
        assertEquals(new BigInteger("800000000000000000"), result.getAvailableConfirmed());
        assertEquals(new BigInteger("300000000000000000"), result.getAvailableUnconfirmed());
        assertEquals(new BigInteger("200000000000000000"), result.getReservedConfirmed());
        assertEquals(new BigInteger("200000000000000000"), result.getReservedUnconfirmed());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordBalance dto = new TgvalidatordBalance();

        Balance result = BalanceMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNull(result.getTotalConfirmed());
        assertNull(result.getTotalUnconfirmed());
        assertNull(result.getAvailableConfirmed());
        assertNull(result.getAvailableUnconfirmed());
        assertNull(result.getReservedConfirmed());
        assertNull(result.getReservedUnconfirmed());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Balance result = BalanceMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }
}
