package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Fee;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordKeyValue;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FeeMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordKeyValue dto = new TgvalidatordKeyValue();
        dto.setKey("ETH_gas_price");
        dto.setValue("50000000000");

        Fee fee = FeeMapper.INSTANCE.fromDTO(dto);

        assertEquals("ETH_gas_price", fee.getKey());
        assertEquals("50000000000", fee.getValue());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordKeyValue dto = new TgvalidatordKeyValue();

        Fee fee = FeeMapper.INSTANCE.fromDTO(dto);

        assertNull(fee.getKey());
        assertNull(fee.getValue());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Fee fee = FeeMapper.INSTANCE.fromDTO(null);
        assertNull(fee);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordKeyValue dto1 = new TgvalidatordKeyValue();
        dto1.setKey("ETH_gas_price");
        dto1.setValue("50000000000");

        TgvalidatordKeyValue dto2 = new TgvalidatordKeyValue();
        dto2.setKey("BTC_sat_per_byte");
        dto2.setValue("25");

        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(fees);
        assertEquals(2, fees.size());

        assertEquals("ETH_gas_price", fees.get(0).getKey());
        assertEquals("50000000000", fees.get(0).getValue());

        assertEquals("BTC_sat_per_byte", fees.get(1).getKey());
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
        TgvalidatordKeyValue eth = new TgvalidatordKeyValue();
        eth.setKey("ETH");
        eth.setValue("21000");

        TgvalidatordKeyValue btc = new TgvalidatordKeyValue();
        btc.setKey("BTC");
        btc.setValue("10");

        TgvalidatordKeyValue sol = new TgvalidatordKeyValue();
        sol.setKey("SOL");
        sol.setValue("5000");

        List<Fee> fees = FeeMapper.INSTANCE.fromDTOList(Arrays.asList(eth, btc, sol));

        assertEquals(3, fees.size());
        assertEquals("ETH", fees.get(0).getKey());
        assertEquals("BTC", fees.get(1).getKey());
        assertEquals("SOL", fees.get(2).getKey());
    }
}
