package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Transaction;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTransaction;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class TransactionMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordTransaction dto = new TgvalidatordTransaction();
        dto.setId("12345");
        dto.setRequestId("67890");
        dto.setDirection("OUTGOING");
        dto.setNetwork("mainnet");
        dto.setBlockchain("ETH");
        dto.setCurrency("ETH");
        dto.setAmount("1000000000000000000");
        dto.setAmountMainUnit("1.0");
        dto.setFee("21000");
        dto.setFeeMainUnit("0.000021");
        dto.setHash("0xabc123def456");
        dto.setBlock("1000000");
        dto.setConfirmationBlock("1000006");
        dto.setTransactionId("tx-123");
        dto.setType("transfer");
        dto.setUniqueId("unique-456");
        dto.setArg1("arg1-value");
        dto.setArg2("arg2-value");
        dto.setRequestVisible(true);

        OffsetDateTime receptionDate = OffsetDateTime.now();
        OffsetDateTime confirmationDate = receptionDate.plusHours(1);
        dto.setReceptionDate(receptionDate);
        dto.setConfirmationDate(confirmationDate);

        // When
        Transaction result = TransactionMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(12345L, result.getId());
        assertEquals(67890L, result.getRequestId());
        assertEquals("OUTGOING", result.getDirection());
        assertEquals("mainnet", result.getNetwork());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("ETH", result.getCurrency());
        assertEquals(new BigInteger("1000000000000000000"), result.getAmount());
        assertEquals(1.0, result.getAmountMainUnit(), 0.001);
        assertEquals(new BigInteger("21000"), result.getFee());
        assertEquals(0.000021, result.getFeeMainUnit(), 0.0000001);
        assertEquals("0xabc123def456", result.getHash());
        assertEquals(1000000L, result.getBlock());
        assertEquals(1000006L, result.getConfirmationBlock());
        assertEquals(receptionDate, result.getReceptionDate());
        assertEquals(confirmationDate, result.getConfirmationDate());
        assertEquals("tx-123", result.getTransactionId());
        assertEquals("transfer", result.getType());
        assertEquals("unique-456", result.getUniqueId());
        assertEquals("arg1-value", result.getArg1());
        assertEquals("arg2-value", result.getArg2());
        assertTrue(result.isRequestVisible());
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordTransaction dto = new TgvalidatordTransaction();
        dto.setId("1");
        // All other fields left null

        // When
        Transaction result = TransactionMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertNull(result.getDirection());
        assertNull(result.getHash());
        assertNull(result.getCurrency());
        assertNull(result.getAmount());
    }

    @Test
    void fromDTO_withSources_mapsAddressList() {
        // Given
        TgvalidatordTransaction dto = new TgvalidatordTransaction();
        dto.setId("1");

        TgvalidatordAddressInfo source1 = new TgvalidatordAddressInfo();
        source1.setAddress("0x1234567890abcdef");

        TgvalidatordAddressInfo source2 = new TgvalidatordAddressInfo();
        source2.setAddress("0xfedcba0987654321");

        dto.setSources(Arrays.asList(source1, source2));

        // When
        Transaction result = TransactionMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getSources());
        assertEquals(2, result.getSources().size());
    }

    @Test
    void fromDTO_withDestinations_mapsAddressList() {
        // Given
        TgvalidatordTransaction dto = new TgvalidatordTransaction();
        dto.setId("1");

        TgvalidatordAddressInfo dest = new TgvalidatordAddressInfo();
        dest.setAddress("0xdestination123");

        dto.setDestinations(Collections.singletonList(dest));

        // When
        Transaction result = TransactionMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getDestinations());
        assertEquals(1, result.getDestinations().size());
    }

    @Test
    void fromDTO_withCurrencyInfo_mapsCurrency() {
        // Given
        TgvalidatordTransaction dto = new TgvalidatordTransaction();
        dto.setId("1");

        TgvalidatordCurrency currencyInfo = new TgvalidatordCurrency();
        currencyInfo.setId("ETH");
        currencyInfo.setName("Ethereum");
        currencyInfo.setDecimals("18");
        dto.setCurrencyInfo(currencyInfo);

        // When
        Transaction result = TransactionMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getCurrencyInfo());
        assertEquals("ETH", result.getCurrencyInfo().getId());
        assertEquals("Ethereum", result.getCurrencyInfo().getName());
    }

    @Test
    void fromDTO_list_mapsMultipleTransactions() {
        // Given
        TgvalidatordTransaction dto1 = new TgvalidatordTransaction();
        dto1.setId("1");
        dto1.setHash("hash1");

        TgvalidatordTransaction dto2 = new TgvalidatordTransaction();
        dto2.setId("2");
        dto2.setHash("hash2");

        List<TgvalidatordTransaction> dtos = Arrays.asList(dto1, dto2);

        // When
        List<Transaction> results = TransactionMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals(1L, results.get(0).getId());
        assertEquals("hash1", results.get(0).getHash());
        assertEquals(2L, results.get(1).getId());
        assertEquals("hash2", results.get(1).getHash());
    }

    @Test
    void fromDTO_emptyList_returnsEmptyList() {
        // Given
        List<TgvalidatordTransaction> dtos = Collections.emptyList();

        // When
        List<Transaction> results = TransactionMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertTrue(results.isEmpty());
    }
}
