package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Score;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordScore;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTransactionScore;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ScoreMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordScore dto = new TgvalidatordScore();
        dto.setId("123");
        dto.setProvider("chainalysis");
        dto.setType("in");
        dto.setScore("85");
        dto.setUpdateDate(now);

        Score result = ScoreMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(123L, result.getId());
        assertEquals("chainalysis", result.getProvider());
        assertEquals("in", result.getType());
        assertEquals("85", result.getScore());
        assertEquals(now, result.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordScore dto = new TgvalidatordScore();
        dto.setId("1");

        Score result = ScoreMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertNull(result.getProvider());
        assertNull(result.getType());
        assertNull(result.getScore());
        assertNull(result.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Score result = ScoreMapper.INSTANCE.fromDTO((TgvalidatordScore) null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordScore dto1 = new TgvalidatordScore();
        dto1.setId("1");
        dto1.setProvider("chainalysis");

        TgvalidatordScore dto2 = new TgvalidatordScore();
        dto2.setId("2");
        dto2.setProvider("elliptic");

        List<Score> result = ScoreMapper.INSTANCE.fromDTO(Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("chainalysis", result.get(0).getProvider());
        assertEquals("elliptic", result.get(1).getProvider());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Score> result = ScoreMapper.INSTANCE.fromDTO(Collections.<TgvalidatordScore>emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Score> result = ScoreMapper.INSTANCE.fromDTO((List<TgvalidatordScore>) null);
        assertNull(result);
    }

    @Test
    void fromTransactionScoreDTO_mapsFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordTransactionScore dto = new TgvalidatordTransactionScore();
        dto.setId("456");
        dto.setProvider("scorechain");
        dto.setType("out");
        dto.setScore("72");
        dto.setUpdateDate(now);

        Score result = ScoreMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(456L, result.getId());
        assertEquals("scorechain", result.getProvider());
        assertEquals("out", result.getType());
        assertEquals("72", result.getScore());
        assertEquals(now, result.getUpdateDate());
    }

    @Test
    void fromTransactionScoreDTO_handlesNullDto() {
        Score result = ScoreMapper.INSTANCE.fromDTO((TgvalidatordTransactionScore) null);
        assertNull(result);
    }
}
