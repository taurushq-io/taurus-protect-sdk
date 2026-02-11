package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Tag;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTag;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class TagMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();

        TgvalidatordTag dto = new TgvalidatordTag();
        dto.setId("tag-123");
        dto.setValue("Production");
        dto.setColor("#FF0000");
        dto.setCreationDate(creationDate);

        Tag tag = TagMapper.INSTANCE.fromDTO(dto);

        assertEquals("tag-123", tag.getId());
        assertEquals("Production", tag.getValue());
        assertEquals("#FF0000", tag.getColor());
        assertEquals(creationDate, tag.getCreationDate());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordTag dto = new TgvalidatordTag();
        dto.setId("tag-456");

        Tag tag = TagMapper.INSTANCE.fromDTO(dto);

        assertEquals("tag-456", tag.getId());
        assertNull(tag.getValue());
        assertNull(tag.getColor());
        assertNull(tag.getCreationDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Tag tag = TagMapper.INSTANCE.fromDTO(null);
        assertNull(tag);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordTag dto1 = new TgvalidatordTag();
        dto1.setId("tag-1");
        dto1.setValue("Development");

        TgvalidatordTag dto2 = new TgvalidatordTag();
        dto2.setId("tag-2");
        dto2.setValue("Testing");

        List<Tag> tags = TagMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(tags);
        assertEquals(2, tags.size());
        assertEquals("tag-1", tags.get(0).getId());
        assertEquals("Development", tags.get(0).getValue());
        assertEquals("tag-2", tags.get(1).getId());
        assertEquals("Testing", tags.get(1).getValue());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Tag> tags = TagMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(tags);
        assertTrue(tags.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Tag> tags = TagMapper.INSTANCE.fromDTOList(null);
        assertNull(tags);
    }

    @Test
    void fromDTO_mapsTagWithColor() {
        TgvalidatordTag dto = new TgvalidatordTag();
        dto.setId("color-tag");
        dto.setValue("Important");
        dto.setColor("#00FF00");

        Tag tag = TagMapper.INSTANCE.fromDTO(dto);

        assertEquals("color-tag", tag.getId());
        assertEquals("Important", tag.getValue());
        assertEquals("#00FF00", tag.getColor());
    }
}
