package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Exchange;
import com.taurushq.sdk.protect.client.model.ExchangeCounterparty;
import com.taurushq.sdk.protect.client.model.ExchangeWithdrawalFee;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExchange;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExchangeCounterparty;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetExchangeWithdrawalFeeReply;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ExchangeMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordCurrency currency = new TgvalidatordCurrency();
        currency.setId("BTC");
        currency.setName("Bitcoin");

        TgvalidatordExchange dto = new TgvalidatordExchange();
        dto.setId("exchange-123");
        dto.setExchange("binance");
        dto.setAccount("main-account");
        dto.setCurrency("BTC");
        dto.setType("SPOT");
        dto.setTotalBalance("1000000000");
        dto.setStatus("ACTIVE");
        dto.setContainer("default");
        dto.setLabel("Main BTC");
        dto.setDisplayLabel("Binance Main BTC");
        dto.setBaseCurrencyValuation("50000.00");
        dto.setHasWLA(true);
        dto.setCurrencyInfo(currency);
        dto.setCreationDate(now);
        dto.setUpdateDate(now);

        Exchange exchange = ExchangeMapper.INSTANCE.fromDTO(dto);

        assertEquals("exchange-123", exchange.getId());
        assertEquals("binance", exchange.getExchange());
        assertEquals("main-account", exchange.getAccount());
        assertEquals("BTC", exchange.getCurrency());
        assertEquals("SPOT", exchange.getType());
        assertEquals("1000000000", exchange.getTotalBalance());
        assertEquals("ACTIVE", exchange.getStatus());
        assertEquals("default", exchange.getContainer());
        assertEquals("Main BTC", exchange.getLabel());
        assertEquals("Binance Main BTC", exchange.getDisplayLabel());
        assertEquals("50000.00", exchange.getBaseCurrencyValuation());
        assertEquals(true, exchange.getHasWLA());
        assertNotNull(exchange.getCurrencyInfo());
        assertEquals("BTC", exchange.getCurrencyInfo().getId());
        assertEquals(now, exchange.getCreationDate());
        assertEquals(now, exchange.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordExchange dto = new TgvalidatordExchange();
        dto.setId("exchange-456");

        Exchange exchange = ExchangeMapper.INSTANCE.fromDTO(dto);

        assertEquals("exchange-456", exchange.getId());
        assertNull(exchange.getExchange());
        assertNull(exchange.getAccount());
        assertNull(exchange.getCurrency());
        assertNull(exchange.getType());
        assertNull(exchange.getTotalBalance());
        assertNull(exchange.getStatus());
        assertNull(exchange.getContainer());
        assertNull(exchange.getLabel());
        assertNull(exchange.getDisplayLabel());
        assertNull(exchange.getBaseCurrencyValuation());
        assertNull(exchange.getHasWLA());
        assertNull(exchange.getCurrencyInfo());
        assertNull(exchange.getCreationDate());
        assertNull(exchange.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Exchange exchange = ExchangeMapper.INSTANCE.fromDTO(null);
        assertNull(exchange);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordExchange dto1 = new TgvalidatordExchange();
        dto1.setId("exchange-1");
        dto1.setExchange("binance");

        TgvalidatordExchange dto2 = new TgvalidatordExchange();
        dto2.setId("exchange-2");
        dto2.setExchange("coinbase");

        List<Exchange> exchanges = ExchangeMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(exchanges);
        assertEquals(2, exchanges.size());
        assertEquals("exchange-1", exchanges.get(0).getId());
        assertEquals("binance", exchanges.get(0).getExchange());
        assertEquals("exchange-2", exchanges.get(1).getId());
        assertEquals("coinbase", exchanges.get(1).getExchange());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Exchange> exchanges = ExchangeMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(exchanges);
        assertTrue(exchanges.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Exchange> exchanges = ExchangeMapper.INSTANCE.fromDTOList(null);
        assertNull(exchanges);
    }

    @Test
    void fromCounterpartyDTO_mapsAllFields() {
        TgvalidatordExchangeCounterparty dto = new TgvalidatordExchangeCounterparty();
        dto.setName("binance");
        dto.setBaseCurrencyValuation("100000.50");

        ExchangeCounterparty counterparty = ExchangeMapper.INSTANCE.fromCounterpartyDTO(dto);

        assertEquals("binance", counterparty.getName());
        assertEquals("100000.50", counterparty.getBaseCurrencyValuation());
    }

    @Test
    void fromCounterpartyDTO_handlesNullDto() {
        ExchangeCounterparty counterparty = ExchangeMapper.INSTANCE.fromCounterpartyDTO(null);
        assertNull(counterparty);
    }

    @Test
    void fromCounterpartyDTOList_mapsList() {
        TgvalidatordExchangeCounterparty dto1 = new TgvalidatordExchangeCounterparty();
        dto1.setName("binance");
        dto1.setBaseCurrencyValuation("100000.00");

        TgvalidatordExchangeCounterparty dto2 = new TgvalidatordExchangeCounterparty();
        dto2.setName("coinbase");
        dto2.setBaseCurrencyValuation("50000.00");

        List<ExchangeCounterparty> counterparties =
                ExchangeMapper.INSTANCE.fromCounterpartyDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(counterparties);
        assertEquals(2, counterparties.size());
        assertEquals("binance", counterparties.get(0).getName());
        assertEquals("100000.00", counterparties.get(0).getBaseCurrencyValuation());
        assertEquals("coinbase", counterparties.get(1).getName());
        assertEquals("50000.00", counterparties.get(1).getBaseCurrencyValuation());
    }

    @Test
    void fromCounterpartyDTOList_handlesEmptyList() {
        List<ExchangeCounterparty> counterparties =
                ExchangeMapper.INSTANCE.fromCounterpartyDTOList(Collections.emptyList());
        assertNotNull(counterparties);
        assertTrue(counterparties.isEmpty());
    }

    @Test
    void fromWithdrawalFeeReply_mapsFee() {
        TgvalidatordGetExchangeWithdrawalFeeReply reply = new TgvalidatordGetExchangeWithdrawalFeeReply();
        reply.setResult("0.0001");

        ExchangeWithdrawalFee fee = ExchangeMapper.INSTANCE.fromWithdrawalFeeReply(reply);

        assertNotNull(fee);
        assertEquals("0.0001", fee.getFee());
    }

    @Test
    void fromWithdrawalFeeReply_handlesNullFee() {
        TgvalidatordGetExchangeWithdrawalFeeReply reply = new TgvalidatordGetExchangeWithdrawalFeeReply();
        reply.setResult(null);

        ExchangeWithdrawalFee fee = ExchangeMapper.INSTANCE.fromWithdrawalFeeReply(reply);

        assertNotNull(fee);
        assertNull(fee.getFee());
    }

    @Test
    void fromWithdrawalFeeReply_handlesNullReply() {
        ExchangeWithdrawalFee fee = ExchangeMapper.INSTANCE.fromWithdrawalFeeReply(null);
        assertNull(fee);
    }
}
