package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.MultiFactorSignatureEntityType;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureInfo;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateMultiFactorSignaturesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMultiFactorSignaturesEntityType;
import org.junit.jupiter.api.Test;

import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class MultiFactorSignatureMapperTest {

    @Test
    void fromDTO_mapsInfoFields() {
        TgvalidatordGetMultiFactorSignatureEntitiesInfoReply dto =
                new TgvalidatordGetMultiFactorSignatureEntitiesInfoReply();
        dto.setId("mfs-123");
        dto.setPayloadToSign(Arrays.asList("payload1", "payload2"));
        dto.setEntityType(TgvalidatordMultiFactorSignaturesEntityType.REQUEST);

        MultiFactorSignatureInfo result = MultiFactorSignatureMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("mfs-123", result.getId());
        assertEquals(2, result.getPayloadToSign().size());
        assertEquals("payload1", result.getPayloadToSign().get(0));
        assertNotNull(result.getEntityType());
    }

    @Test
    void fromDTO_handlesNullInfo() {
        MultiFactorSignatureInfo result = MultiFactorSignatureMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromCreateDTO_mapsResultFields() {
        TgvalidatordCreateMultiFactorSignaturesReply dto =
                new TgvalidatordCreateMultiFactorSignaturesReply();
        dto.setId("mfs-result-123");

        MultiFactorSignatureResult result = MultiFactorSignatureMapper.INSTANCE.fromCreateDTO(dto);

        assertNotNull(result);
        assertEquals("mfs-result-123", result.getId());
    }

    @Test
    void fromCreateDTO_handlesNullResult() {
        MultiFactorSignatureResult result = MultiFactorSignatureMapper.INSTANCE.fromCreateDTO(null);
        assertNull(result);
    }

    @Test
    void fromEntityTypeDTO_mapsFields() {
        MultiFactorSignatureEntityType result = MultiFactorSignatureMapper.INSTANCE
                .fromEntityTypeDTO(TgvalidatordMultiFactorSignaturesEntityType.REQUEST);

        assertNotNull(result);
    }

    @Test
    void fromEntityTypeDTO_handlesNull() {
        MultiFactorSignatureEntityType result = MultiFactorSignatureMapper.INSTANCE
                .fromEntityTypeDTO(null);

        assertNull(result);
    }
}
