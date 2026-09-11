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
        // assertNotNull alone was the whole assertion, and that is why the all-null bean
        // below went unnoticed: the kind is what tells a caller WHICH verifying reader to
        // check payloadToSign against, so a null one leaves the payload uncheckable.
        assertEquals("REQUEST", result.getEntityType().getKind());
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
    void fromEntityTypeDTO_mapsTheKindForEveryValue() {
        // MapStruct's default bean mapping from an ENUM source sets no target properties at
        // all, so this returned a bean with both fields null. The mapper is hand-written now.
        for (TgvalidatordMultiFactorSignaturesEntityType dto
                : TgvalidatordMultiFactorSignaturesEntityType.values()) {
            MultiFactorSignatureEntityType result =
                    MultiFactorSignatureMapper.INSTANCE.fromEntityTypeDTO(dto);

            assertNotNull(result);
            assertEquals(dto.getValue(), result.getKind(),
                    "the kind must survive the mapping for " + dto);
        }
    }

    @Test
    void fromEntityTypeDTO_leavesIdNullBecauseTheReplyCarriesNone() {
        // Not an oversight to be "fixed" by inventing an id: the reply has no entity id, and
        // that absence is exactly why the SDK cannot bind payloadToSign to a verified entity.
        MultiFactorSignatureEntityType result = MultiFactorSignatureMapper.INSTANCE
                .fromEntityTypeDTO(TgvalidatordMultiFactorSignaturesEntityType.REQUEST);

        assertNull(result.getId());
    }

    @Test
    void fromEntityTypeDTO_handlesNull() {
        MultiFactorSignatureEntityType result = MultiFactorSignatureMapper.INSTANCE
                .fromEntityTypeDTO(null);

        assertNull(result);
    }
}
