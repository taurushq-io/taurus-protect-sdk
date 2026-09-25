package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Fee;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFee;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FeeMapperTest {

    private static TgvalidatordFee fee(final String currencyId, final String value) {
        TgvalidatordFee dto = new TgvalidatordFee();
        dto.setCurrencyId(currencyId);
        dto.setValue(value);
        return dto;
    }

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime updated = OffsetDateTime.parse("2026-09-24T10:00:00Z");
        TgvalidatordCurrency currency = new TgvalidatordCurrency();
        currency.setId("ETH");
        currency.setSymbol("ETH");
        TgvalidatordFee dto = fee("ETH", "50000000000");
        dto.setDenom("wei");
        dto.setCurrencyInfo(currency);
        dto.setUpdateDate(updated);

        Fee fee = FeeMapper.INSTANCE.fromDTO(dto);

        assertEquals("ETH", fee.getCurrencyId());
        assertEquals("50000000000", fee.getValue());
        assertEquals("wei", fee.getDenom());
        assertEquals(updated, fee.getUpdateDate());
        assertNotNull(fee.getCurrencyInfo());
        assertEquals("ETH", fee.getCurrencyInfo().getSymbol());
    }

    @Test
    void fromDTO_handlesNullFields() {
        Fee fee = FeeMapper.INSTANCE.fromDTO(new TgvalidatordFee());

        assertNull(fee.getCurrencyId());
        assertNull(fee.getValue());
        assertNull(fee.getDenom());
        assertNull(fee.getCurrencyInfo());
        assertNull(fee.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Fee fee = FeeMapper.INSTANCE.fromDTO(null);
        assertNull(fee);
    }

    @Test
    void fromDTOList_mapsList() {
        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(
                Arrays.asList(fee("ETH", "50000000000"), fee("BTC", "25")));

        assertNotNull(fees);
        assertEquals(2, fees.size());

        assertEquals("ETH", fees.get(0).getCurrencyId());
        assertEquals("50000000000", fees.get(0).getValue());

        assertEquals("BTC", fees.get(1).getCurrencyId());
        assertEquals("25", fees.get(1).getValue());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(fees);
        assertTrue(fees.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(null);
        assertNull(fees);
    }

    @Test
    void fromDTOList_mapsMultipleFees() {
        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(
                Arrays.asList(fee("ETH", "21000"), fee("BTC", "10"), fee("SOL", "5000")));

        assertEquals(3, fees.size());
        assertEquals("ETH", fees.get(0).getCurrencyId());
        assertEquals("BTC", fees.get(1).getCurrencyId());
        assertEquals("SOL", fees.get(2).getCurrencyId());
    }
}
