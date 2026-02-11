package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.GovernanceRulesTrail;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SuperAdminPublicKey;
import com.taurushq.sdk.protect.openapi.model.GetPublicKeysReplyPublicKey;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRuleUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRules;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRulesTrail;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class GovernanceRulesMapperTest {

    @Test
    void fromRuleUserSignatureDTO_mapsAllFields() {
        TgvalidatordRuleUserSignature dto = new TgvalidatordRuleUserSignature();
        dto.setUserId("user-123");
        dto.setSignature("base64sig==");

        RuleUserSignature result = GovernanceRulesMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("user-123", result.getUserId());
        assertEquals("base64sig==", result.getSignature());
    }

    @Test
    void fromRuleUserSignatureDTO_handlesNullDto() {
        RuleUserSignature result = GovernanceRulesMapper.INSTANCE.fromDTO(
                (TgvalidatordRuleUserSignature) null);
        assertNull(result);
    }

    @Test
    void fromRuleUserSignatureDTOs_mapsList() {
        TgvalidatordRuleUserSignature dto1 = new TgvalidatordRuleUserSignature();
        dto1.setUserId("user-1");
        dto1.setSignature("sig1");

        TgvalidatordRuleUserSignature dto2 = new TgvalidatordRuleUserSignature();
        dto2.setUserId("user-2");
        dto2.setSignature("sig2");

        List<RuleUserSignature> result = GovernanceRulesMapper.INSTANCE.fromRuleUserSignatureDTOs(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("user-1", result.get(0).getUserId());
        assertEquals("sig1", result.get(0).getSignature());
        assertEquals("user-2", result.get(1).getUserId());
    }

    @Test
    void fromRuleUserSignatureDTOs_handlesEmptyList() {
        List<RuleUserSignature> result = GovernanceRulesMapper.INSTANCE.fromRuleUserSignatureDTOs(
                Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromRuleUserSignatureDTOs_handlesNullList() {
        List<RuleUserSignature> result = GovernanceRulesMapper.INSTANCE.fromRuleUserSignatureDTOs(null);
        assertNull(result);
    }

    @Test
    void fromRulesTrailDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordRulesTrail dto = new TgvalidatordRulesTrail();
        dto.setId("trail-123");
        dto.setUserId("user-456");
        dto.setExternalUserId("ext-789");
        dto.setAction("approve");
        dto.setComment("Approved changes");
        dto.setDate(now);

        GovernanceRulesTrail result = GovernanceRulesMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("trail-123", result.getId());
        assertEquals("user-456", result.getUserId());
        assertEquals("ext-789", result.getExternalUserId());
        assertEquals("approve", result.getAction());
        assertEquals("Approved changes", result.getComment());
        assertEquals(now, result.getDate());
    }

    @Test
    void fromRulesTrailDTO_handlesNullDto() {
        GovernanceRulesTrail result = GovernanceRulesMapper.INSTANCE.fromDTO(
                (TgvalidatordRulesTrail) null);
        assertNull(result);
    }

    @Test
    void fromRulesTrailDTOs_mapsList() {
        TgvalidatordRulesTrail dto1 = new TgvalidatordRulesTrail();
        dto1.setId("t1");
        dto1.setAction("create");

        TgvalidatordRulesTrail dto2 = new TgvalidatordRulesTrail();
        dto2.setId("t2");
        dto2.setAction("approve");

        List<GovernanceRulesTrail> result = GovernanceRulesMapper.INSTANCE.fromRulesTrailDTOs(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("create", result.get(0).getAction());
        assertEquals("approve", result.get(1).getAction());
    }

    @Test
    void fromRulesTrailDTOs_handlesNullList() {
        List<GovernanceRulesTrail> result = GovernanceRulesMapper.INSTANCE.fromRulesTrailDTOs(null);
        assertNull(result);
    }

    @Test
    void fromRulesDTO_mapsAllFields() throws Exception {
        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = OffsetDateTime.now().plusHours(1);

        TgvalidatordRuleUserSignature sig = new TgvalidatordRuleUserSignature();
        sig.setUserId("admin-1");
        sig.setSignature("rulesig==");

        TgvalidatordRulesTrail trail = new TgvalidatordRulesTrail();
        trail.setId("t-1");
        trail.setAction("create");

        TgvalidatordRules dto = new TgvalidatordRules();
        dto.setRulesContainer("base64rulescontainer==");
        dto.setRulesSignatures(Arrays.asList(sig));
        dto.setLocked(true);
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setTrails(Arrays.asList(trail));

        GovernanceRules result = GovernanceRulesMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("base64rulescontainer==", result.getRulesContainer());
        assertNotNull(result.getRulesSignatures());
        assertEquals(1, result.getRulesSignatures().size());
        assertEquals("admin-1", result.getRulesSignatures().get(0).getUserId());
        assertTrue(result.getLocked());
        assertEquals(creationDate, result.getCreationDate());
        assertEquals(updateDate, result.getUpdateDate());
        assertNotNull(result.getTrails());
        assertEquals(1, result.getTrails().size());
        assertEquals("create", result.getTrails().get(0).getAction());
    }

    @Test
    void fromRulesDTO_handlesNullDto() throws Exception {
        GovernanceRules result = GovernanceRulesMapper.INSTANCE.fromDTO(
                (TgvalidatordRules) null);
        assertNull(result);
    }

    @Test
    void fromRulesDTO_handlesNullFields() throws Exception {
        TgvalidatordRules dto = new TgvalidatordRules();

        GovernanceRules result = GovernanceRulesMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNull(result.getRulesContainer());
        // MapStruct maps null list to empty list
        assertNotNull(result.getRulesSignatures());
        assertNull(result.getLocked());
    }

    @Test
    void fromPublicKeyDTO_mapsAllFields() {
        GetPublicKeysReplyPublicKey dto = new GetPublicKeysReplyPublicKey();
        dto.setUserID("user-sa-1");
        dto.setPublicKey("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...");

        SuperAdminPublicKey result = GovernanceRulesMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("user-sa-1", result.getUserId());
        assertEquals("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...", result.getPublicKey());
    }

    @Test
    void fromPublicKeyDTO_handlesNullDto() {
        SuperAdminPublicKey result = GovernanceRulesMapper.INSTANCE.fromDTO(
                (GetPublicKeysReplyPublicKey) null);
        assertNull(result);
    }

    @Test
    void fromPublicKeyDTOs_mapsList() {
        GetPublicKeysReplyPublicKey dto1 = new GetPublicKeysReplyPublicKey();
        dto1.setUserID("sa-1");
        dto1.setPublicKey("key1");

        GetPublicKeysReplyPublicKey dto2 = new GetPublicKeysReplyPublicKey();
        dto2.setUserID("sa-2");
        dto2.setPublicKey("key2");

        List<SuperAdminPublicKey> result = GovernanceRulesMapper.INSTANCE.fromPublicKeyDTOs(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("sa-1", result.get(0).getUserId());
        assertEquals("key1", result.get(0).getPublicKey());
        assertEquals("sa-2", result.get(1).getUserId());
    }

    @Test
    void fromPublicKeyDTOs_handlesNullList() {
        List<SuperAdminPublicKey> result = GovernanceRulesMapper.INSTANCE.fromPublicKeyDTOs(null);
        assertNull(result);
    }
}
