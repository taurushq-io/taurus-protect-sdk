package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ConversionResult;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionValue;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyPrice;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordPricesHistoryPoint;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class PriceMapperTest {

    @Test
    void fromDTO_currencyPrice_mapsAllFields() {
        // Given
        TgvalidatordCurrencyPrice dto = new TgvalidatordCurrencyPrice();
        dto.setBlockchain("ETH");
        dto.setCurrencyFrom("ETH");
        dto.setCurrencyTo("USD");
        dto.setDecimals("18");
        dto.setRate("2500.50");
        dto.setSource("exchange");
        dto.setChangePercent24Hour("2.5");

        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = creationDate.plusHours(1);
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);

        // When
        Price result = PriceMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("ETH", result.getBlockchain());
        assertEquals("ETH", result.getCurrencyFrom());
        assertEquals("USD", result.getCurrencyTo());
        assertEquals("18", result.getDecimals());
        assertEquals("2500.50", result.getRate());
        assertEquals("exchange", result.getSource());
        assertEquals("2.5", result.getChangePercent24Hour());
        assertEquals(creationDate, result.getCreatedAt());
        assertEquals(updateDate, result.getUpdatedAt());
    }

    @Test
    void fromDTO_currencyPrice_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordCurrencyPrice dto = new TgvalidatordCurrencyPrice();
        dto.setCurrencyFrom("BTC");
        // All other fields left null

        // When
        Price result = PriceMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("BTC", result.getCurrencyFrom());
        assertNull(result.getRate());
        assertNull(result.getCreatedAt());
    }

    @Test
    void fromDTO_priceList_mapsMultiplePrices() {
        // Given
        TgvalidatordCurrencyPrice dto1 = new TgvalidatordCurrencyPrice();
        dto1.setCurrencyFrom("ETH");
        dto1.setRate("2500.00");

        TgvalidatordCurrencyPrice dto2 = new TgvalidatordCurrencyPrice();
        dto2.setCurrencyFrom("BTC");
        dto2.setRate("45000.00");

        List<TgvalidatordCurrencyPrice> dtos = Arrays.asList(dto1, dto2);

        // When
        List<Price> results = PriceMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("ETH", results.get(0).getCurrencyFrom());
        assertEquals("2500.00", results.get(0).getRate());
        assertEquals("BTC", results.get(1).getCurrencyFrom());
        assertEquals("45000.00", results.get(1).getRate());
    }

    @Test
    void fromDTO_priceHistoryPoint_mapsAllFields() {
        // Given
        TgvalidatordPricesHistoryPoint dto = new TgvalidatordPricesHistoryPoint();
        dto.setHigh("2600.00");
        dto.setLow("2400.00");
        dto.setOpen("2450.00");
        dto.setClose("2550.00");
        dto.setVolumeFrom("1000000");
        dto.setVolumeTo("2500000000");
        dto.setChangePercent("4.08");
        dto.setBlockchain("ETH");
        dto.setCurrencyFrom("ETH");
        dto.setCurrencyTo("USD");

        OffsetDateTime periodStartDate = OffsetDateTime.now();
        dto.setPeriodStartDate(periodStartDate);

        // When
        PriceHistoryPoint result = PriceMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("2600.00", result.getHigh());
        assertEquals("2400.00", result.getLow());
        assertEquals("2450.00", result.getOpen());
        assertEquals("2550.00", result.getClose());
        assertEquals("1000000", result.getVolumeFrom());
        assertEquals("2500000000", result.getVolumeTo());
        assertEquals("4.08", result.getChangePercent());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("ETH", result.getCurrencyFrom());
        assertEquals("USD", result.getCurrencyTo());
        assertEquals(periodStartDate, result.getPeriodStartDate());
    }

    @Test
    void fromDTOHistory_mapsMultipleHistoryPoints() {
        // Given
        TgvalidatordPricesHistoryPoint dto1 = new TgvalidatordPricesHistoryPoint();
        dto1.setOpen("2450.00");
        dto1.setClose("2500.00");

        TgvalidatordPricesHistoryPoint dto2 = new TgvalidatordPricesHistoryPoint();
        dto2.setOpen("2500.00");
        dto2.setClose("2550.00");

        List<TgvalidatordPricesHistoryPoint> dtos = Arrays.asList(dto1, dto2);

        // When
        List<PriceHistoryPoint> results = PriceMapper.INSTANCE.fromDTOHistory(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("2450.00", results.get(0).getOpen());
        assertEquals("2500.00", results.get(0).getClose());
        assertEquals("2500.00", results.get(1).getOpen());
        assertEquals("2550.00", results.get(1).getClose());
    }

    @Test
    void fromDTO_conversionResult_mapsAllFields() {
        // Given
        TgvalidatordConversionValue dto = new TgvalidatordConversionValue();
        dto.setSymbol("USD");
        dto.setValue("62500.00");
        dto.setMainUnitValue("62500.00");

        // When
        ConversionResult result = PriceMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("USD", result.getSymbol());
        assertEquals("62500.00", result.getValue());
        assertEquals("62500.00", result.getMainUnitValue());
    }

    @Test
    void fromDTOConversion_mapsMultipleConversionResults() {
        // Given
        TgvalidatordConversionValue dto1 = new TgvalidatordConversionValue();
        dto1.setSymbol("USD");
        dto1.setValue("62500.00");

        TgvalidatordConversionValue dto2 = new TgvalidatordConversionValue();
        dto2.setSymbol("EUR");
        dto2.setValue("57500.00");

        List<TgvalidatordConversionValue> dtos = Arrays.asList(dto1, dto2);

        // When
        List<ConversionResult> results = PriceMapper.INSTANCE.fromDTOConversion(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("USD", results.get(0).getSymbol());
        assertEquals("62500.00", results.get(0).getValue());
        assertEquals("EUR", results.get(1).getSymbol());
        assertEquals("57500.00", results.get(1).getValue());
    }

    @Test
    void fromDTOHistory_emptyList_returnsEmptyList() {
        // Given
        List<TgvalidatordPricesHistoryPoint> dtos = Collections.emptyList();

        // When
        List<PriceHistoryPoint> results = PriceMapper.INSTANCE.fromDTOHistory(dtos);

        // Then
        assertNotNull(results);
        assertTrue(results.isEmpty());
    }

    @Test
    void fromDTOConversion_emptyList_returnsEmptyList() {
        // Given
        List<TgvalidatordConversionValue> dtos = Collections.emptyList();

        // When
        List<ConversionResult> results = PriceMapper.INSTANCE.fromDTOConversion(dtos);

        // Then
        assertNotNull(results);
        assertTrue(results.isEmpty());
    }
}
