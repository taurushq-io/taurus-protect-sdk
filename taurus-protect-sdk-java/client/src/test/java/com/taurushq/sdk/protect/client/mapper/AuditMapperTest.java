package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.AuditTrail;
import com.taurushq.sdk.protect.client.model.AuditTrailResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAuditTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAuditTrailsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AuditMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordAuditTrail dto = new TgvalidatordAuditTrail();
        dto.setId("audit-123");
        dto.setEntity("request");
        dto.setAction("approve");
        dto.setSubAction("signature");
        dto.setDetails("Request approved by user");
        dto.setCreationDate(now);

        AuditTrail auditTrail = AuditMapper.INSTANCE.fromDTO(dto);

        assertEquals("audit-123", auditTrail.getId());
        assertEquals("request", auditTrail.getEntity());
        assertEquals("approve", auditTrail.getAction());
        assertEquals("signature", auditTrail.getSubAction());
        assertEquals("Request approved by user", auditTrail.getDetails());
        assertEquals(now, auditTrail.getCreationDate());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordAuditTrail dto = new TgvalidatordAuditTrail();
        dto.setId("audit-456");

        AuditTrail auditTrail = AuditMapper.INSTANCE.fromDTO(dto);

        assertEquals("audit-456", auditTrail.getId());
        assertNull(auditTrail.getEntity());
        assertNull(auditTrail.getAction());
        assertNull(auditTrail.getSubAction());
        assertNull(auditTrail.getDetails());
        assertNull(auditTrail.getCreationDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        AuditTrail auditTrail = AuditMapper.INSTANCE.fromDTO(null);
        assertNull(auditTrail);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordAuditTrail dto1 = new TgvalidatordAuditTrail();
        dto1.setId("audit-1");
        dto1.setEntity("wallet");
        dto1.setAction("create");

        TgvalidatordAuditTrail dto2 = new TgvalidatordAuditTrail();
        dto2.setId("audit-2");
        dto2.setEntity("address");
        dto2.setAction("delete");

        List<AuditTrail> auditTrails = AuditMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(auditTrails);
        assertEquals(2, auditTrails.size());

        assertEquals("audit-1", auditTrails.get(0).getId());
        assertEquals("wallet", auditTrails.get(0).getEntity());
        assertEquals("create", auditTrails.get(0).getAction());

        assertEquals("audit-2", auditTrails.get(1).getId());
        assertEquals("address", auditTrails.get(1).getEntity());
        assertEquals("delete", auditTrails.get(1).getAction());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<AuditTrail> auditTrails = AuditMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(auditTrails);
        assertTrue(auditTrails.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<AuditTrail> auditTrails = AuditMapper.INSTANCE.fromDTOList(null);
        assertNull(auditTrails);
    }

    @Test
    void fromReply_mapsAuditTrailsAndCursor() {
        TgvalidatordAuditTrail trail1 = new TgvalidatordAuditTrail();
        trail1.setId("audit-1");
        trail1.setEntity("request");
        trail1.setAction("approve");

        TgvalidatordAuditTrail trail2 = new TgvalidatordAuditTrail();
        trail2.setId("audit-2");
        trail2.setEntity("wallet");
        trail2.setAction("create");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-token-xyz");
        cursor.setHasNext(true);

        TgvalidatordGetAuditTrailsReply reply = new TgvalidatordGetAuditTrailsReply();
        reply.setResult(Arrays.asList(trail1, trail2));
        reply.setCursor(cursor);

        AuditTrailResult result = AuditMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getAuditTrails());
        assertEquals(2, result.getAuditTrails().size());

        assertEquals("audit-1", result.getAuditTrails().get(0).getId());
        assertEquals("request", result.getAuditTrails().get(0).getEntity());
        assertEquals("approve", result.getAuditTrails().get(0).getAction());

        assertEquals("audit-2", result.getAuditTrails().get(1).getId());
        assertEquals("wallet", result.getAuditTrails().get(1).getEntity());
        assertEquals("create", result.getAuditTrails().get(1).getAction());

        assertNotNull(result.getCursor());
        assertEquals("page-token-xyz", result.getCursor().getCurrentPage());
        assertTrue(result.hasNext());
    }

    @Test
    void fromReply_handlesEmptyAuditTrails() {
        TgvalidatordGetAuditTrailsReply reply = new TgvalidatordGetAuditTrailsReply();
        reply.setResult(Collections.emptyList());
        reply.setCursor(null);

        AuditTrailResult result = AuditMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getAuditTrails());
        assertTrue(result.getAuditTrails().isEmpty());
        assertFalse(result.hasNext());
    }

    @Test
    void fromReply_handlesNullCursor() {
        TgvalidatordAuditTrail trail = new TgvalidatordAuditTrail();
        trail.setId("audit-single");

        TgvalidatordGetAuditTrailsReply reply = new TgvalidatordGetAuditTrailsReply();
        reply.setResult(Arrays.asList(trail));
        reply.setCursor(null);

        AuditTrailResult result = AuditMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertEquals(1, result.getAuditTrails().size());
        assertNull(result.getCursor());
        assertFalse(result.hasNext());
    }

    @Test
    void fromReply_handlesHasNextFalse() {
        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("last-page");
        cursor.setHasNext(false);

        TgvalidatordGetAuditTrailsReply reply = new TgvalidatordGetAuditTrailsReply();
        reply.setResult(Collections.emptyList());
        reply.setCursor(cursor);

        AuditTrailResult result = AuditMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getCursor());
        assertEquals("last-page", result.getCursor().getCurrentPage());
        assertFalse(result.hasNext());
    }

    @Test
    void fromReply_mapsVariousEntityAndActionTypes() {
        TgvalidatordAuditTrail walletCreate = new TgvalidatordAuditTrail();
        walletCreate.setId("a1");
        walletCreate.setEntity("wallet");
        walletCreate.setAction("create");

        TgvalidatordAuditTrail addressApprove = new TgvalidatordAuditTrail();
        addressApprove.setId("a2");
        addressApprove.setEntity("address");
        addressApprove.setAction("approve");

        TgvalidatordAuditTrail requestReject = new TgvalidatordAuditTrail();
        requestReject.setId("a3");
        requestReject.setEntity("request");
        requestReject.setAction("reject");

        TgvalidatordAuditTrail userLogin = new TgvalidatordAuditTrail();
        userLogin.setId("a4");
        userLogin.setEntity("user");
        userLogin.setAction("login");

        TgvalidatordGetAuditTrailsReply reply = new TgvalidatordGetAuditTrailsReply();
        reply.setResult(Arrays.asList(walletCreate, addressApprove, requestReject, userLogin));

        AuditTrailResult result = AuditMapper.INSTANCE.fromReply(reply);

        assertEquals(4, result.getAuditTrails().size());
        assertEquals("wallet", result.getAuditTrails().get(0).getEntity());
        assertEquals("create", result.getAuditTrails().get(0).getAction());
        assertEquals("address", result.getAuditTrails().get(1).getEntity());
        assertEquals("approve", result.getAuditTrails().get(1).getAction());
        assertEquals("request", result.getAuditTrails().get(2).getEntity());
        assertEquals("reject", result.getAuditTrails().get(2).getAction());
        assertEquals("user", result.getAuditTrails().get(3).getEntity());
        assertEquals("login", result.getAuditTrails().get(3).getAction());
    }
}
