package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.FeePayer;
import com.taurushq.sdk.protect.client.model.FeePayerEth;
import com.taurushq.sdk.protect.client.model.FeePayerEthLocal;
import com.taurushq.sdk.protect.client.model.FeePayerEthRemote;
import com.taurushq.sdk.protect.client.model.FeePayerInfo;
import com.taurushq.sdk.protect.openapi.model.ETHLocal;
import com.taurushq.sdk.protect.openapi.model.ETHRemote;
import com.taurushq.sdk.protect.openapi.model.FeePayerETH;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFeePayer;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFeePayerEnvelope;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FeePayerMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();

        ETHLocal local = new ETHLocal();
        local.setAddressId("addr-local-123");

        FeePayerETH eth = new FeePayerETH();
        eth.setKind("local");
        eth.setLocal(local);

        TgvalidatordFeePayer feePayerInfo = new TgvalidatordFeePayer();
        feePayerInfo.setBlockchain("ETH");
        feePayerInfo.setEth(eth);

        TgvalidatordFeePayerEnvelope dto = new TgvalidatordFeePayerEnvelope();
        dto.setId("fp-123");
        dto.setTenantId("tenant-456");
        dto.setBlockchain("ETH");
        dto.setNetwork("mainnet");
        dto.setName("Main Fee Payer");
        dto.setCreationDate(creationDate);
        dto.setFeePayer(feePayerInfo);

        FeePayer result = FeePayerMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("fp-123", result.getId());
        assertEquals("tenant-456", result.getTenantId());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("mainnet", result.getNetwork());
        assertEquals("Main Fee Payer", result.getName());
        assertEquals(creationDate, result.getCreationDate());
        assertNotNull(result.getFeePayerInfo());
        assertEquals("ETH", result.getFeePayerInfo().getBlockchain());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordFeePayerEnvelope dto = new TgvalidatordFeePayerEnvelope();
        dto.setId("fp-minimal");

        FeePayer result = FeePayerMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("fp-minimal", result.getId());
        assertNull(result.getTenantId());
        assertNull(result.getBlockchain());
        assertNull(result.getNetwork());
        assertNull(result.getName());
        assertNull(result.getCreationDate());
        assertNull(result.getFeePayerInfo());
    }

    @Test
    void fromDTO_handlesNullDto() {
        FeePayer result = FeePayerMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordFeePayerEnvelope dto1 = new TgvalidatordFeePayerEnvelope();
        dto1.setId("fp-1");
        dto1.setName("Fee Payer 1");

        TgvalidatordFeePayerEnvelope dto2 = new TgvalidatordFeePayerEnvelope();
        dto2.setId("fp-2");
        dto2.setName("Fee Payer 2");

        List<FeePayer> result = FeePayerMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("fp-1", result.get(0).getId());
        assertEquals("Fee Payer 1", result.get(0).getName());
        assertEquals("fp-2", result.get(1).getId());
        assertEquals("Fee Payer 2", result.get(1).getName());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<FeePayer> result = FeePayerMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<FeePayer> result = FeePayerMapper.INSTANCE.fromDTOList(null);
        assertNull(result);
    }

    @Test
    void fromFeePayerDTO_mapsFields() {
        FeePayerETH eth = new FeePayerETH();
        eth.setKind("remote");

        TgvalidatordFeePayer dto = new TgvalidatordFeePayer();
        dto.setBlockchain("ETH");
        dto.setEth(eth);

        FeePayerInfo result = FeePayerMapper.INSTANCE.fromFeePayerDTO(dto);

        assertNotNull(result);
        assertEquals("ETH", result.getBlockchain());
        assertNotNull(result.getEth());
        assertEquals("remote", result.getEth().getKind());
    }

    @Test
    void fromFeePayerDTO_handlesNullDto() {
        FeePayerInfo result = FeePayerMapper.INSTANCE.fromFeePayerDTO(null);
        assertNull(result);
    }

    @Test
    void fromEthDTO_mapsAllFields() {
        ETHLocal local = new ETHLocal();
        local.setAddressId("addr-local");
        local.setAutoApprove(true);

        ETHRemote remote = new ETHRemote();
        remote.setUrl("https://example.com");
        remote.setUsername("user1");
        remote.setFromAddressId("from-addr");

        FeePayerETH dto = new FeePayerETH();
        dto.setKind("local");
        dto.setLocal(local);
        dto.setRemote(remote);

        FeePayerEth result = FeePayerMapper.INSTANCE.fromEthDTO(dto);

        assertNotNull(result);
        assertEquals("local", result.getKind());
        assertNotNull(result.getLocal());
        assertEquals("addr-local", result.getLocal().getAddressId());
        assertTrue(result.getLocal().getAutoApprove());
        assertNotNull(result.getRemote());
        assertEquals("https://example.com", result.getRemote().getUrl());
        assertEquals("user1", result.getRemote().getUsername());
        assertEquals("from-addr", result.getRemote().getFromAddressId());
    }

    @Test
    void fromEthLocalDTO_mapsFields() {
        ETHLocal dto = new ETHLocal();
        dto.setAddressId("addr-123");
        dto.setForwarderAddressId("fwd-456");
        dto.setAutoApprove(false);
        dto.setCreatorAddressId("creator-789");

        FeePayerEthLocal result = FeePayerMapper.INSTANCE.fromEthLocalDTO(dto);

        assertNotNull(result);
        assertEquals("addr-123", result.getAddressId());
        assertEquals("fwd-456", result.getForwarderAddressId());
        assertEquals(false, result.getAutoApprove());
        assertEquals("creator-789", result.getCreatorAddressId());
    }

    @Test
    void fromEthRemoteDTO_mapsFields() {
        ETHRemote dto = new ETHRemote();
        dto.setUrl("https://feepayer.example.com");
        dto.setUsername("admin");
        dto.setFromAddressId("from-addr");
        dto.setForwarderAddress("0xForwarder");
        dto.setCreatorAddress("0xCreator");

        FeePayerEthRemote result = FeePayerMapper.INSTANCE.fromEthRemoteDTO(dto);

        assertNotNull(result);
        assertEquals("https://feepayer.example.com", result.getUrl());
        assertEquals("admin", result.getUsername());
        assertEquals("from-addr", result.getFromAddressId());
        assertEquals("0xForwarder", result.getForwarderAddress());
        assertEquals("0xCreator", result.getCreatorAddress());
    }
}
