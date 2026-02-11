package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BalanceHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalanceHistoryPoint;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class BalanceHistoryPointMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordBalance balance = new TgvalidatordBalance();
        balance.setTotalConfirmed("1000");
        balance.setAvailableConfirmed("800");

        TgvalidatordBalanceHistoryPoint dto = new TgvalidatordBalanceHistoryPoint();
        dto.setPointDate(now);
        dto.setBalance(balance);

        BalanceHistoryPoint result = BalanceHistoryPointMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(now, result.getPointDate());
        assertNotNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordBalanceHistoryPoint dto = new TgvalidatordBalanceHistoryPoint();

        BalanceHistoryPoint result = BalanceHistoryPointMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNull(result.getPointDate());
        assertNull(result.getBalance());
    }

    @Test
    void fromDTO_handlesNullDto() {
        BalanceHistoryPoint result = BalanceHistoryPointMapper.INSTANCE.fromDTO((TgvalidatordBalanceHistoryPoint) null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordBalanceHistoryPoint dto1 = new TgvalidatordBalanceHistoryPoint();
        dto1.setPointDate(OffsetDateTime.now());

        TgvalidatordBalanceHistoryPoint dto2 = new TgvalidatordBalanceHistoryPoint();
        dto2.setPointDate(OffsetDateTime.now().plusHours(1));

        List<BalanceHistoryPoint> result = BalanceHistoryPointMapper.INSTANCE.fromDTO(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<BalanceHistoryPoint> result = BalanceHistoryPointMapper.INSTANCE.fromDTO(
                Collections.<TgvalidatordBalanceHistoryPoint>emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<BalanceHistoryPoint> result = BalanceHistoryPointMapper.INSTANCE.fromDTO(
                (List<TgvalidatordBalanceHistoryPoint>) null);
        assertNull(result);
    }
}
