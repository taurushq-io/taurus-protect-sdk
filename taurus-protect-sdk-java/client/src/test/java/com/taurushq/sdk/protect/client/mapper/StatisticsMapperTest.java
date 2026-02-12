package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.PortfolioStatistics;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAggregatedStatsData;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

class StatisticsMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordAggregatedStatsData dto = new TgvalidatordAggregatedStatsData();
        dto.setAvgBalancePerAddress("1000000000000000000");
        dto.setAddressesCount("150");
        dto.setWalletsCount("25");
        dto.setTotalBalance("150000000000000000000");
        dto.setTotalBalanceBaseCurrency("250000.50");

        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(dto);

        assertEquals("1000000000000000000", stats.getAvgBalancePerAddress());
        assertEquals("150", stats.getAddressesCount());
        assertEquals("25", stats.getWalletsCount());
        assertEquals("150000000000000000000", stats.getTotalBalance());
        assertEquals("250000.50", stats.getTotalBalanceBaseCurrency());
    }

    @Test
    void fromDTO_handlesLargeNumbers() {
        TgvalidatordAggregatedStatsData dto = new TgvalidatordAggregatedStatsData();
        dto.setTotalBalance("999999999999999999999999999999");
        dto.setTotalBalanceBaseCurrency("1000000000.99");
        dto.setAddressesCount("1000000");
        dto.setWalletsCount("50000");

        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(dto);

        assertEquals("999999999999999999999999999999", stats.getTotalBalance());
        assertEquals("1000000000.99", stats.getTotalBalanceBaseCurrency());
        assertEquals("1000000", stats.getAddressesCount());
        assertEquals("50000", stats.getWalletsCount());
    }

    @Test
    void fromDTO_handlesZeroValues() {
        TgvalidatordAggregatedStatsData dto = new TgvalidatordAggregatedStatsData();
        dto.setAvgBalancePerAddress("0");
        dto.setAddressesCount("0");
        dto.setWalletsCount("0");
        dto.setTotalBalance("0");
        dto.setTotalBalanceBaseCurrency("0.00");

        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(dto);

        assertEquals("0", stats.getAvgBalancePerAddress());
        assertEquals("0", stats.getAddressesCount());
        assertEquals("0", stats.getWalletsCount());
        assertEquals("0", stats.getTotalBalance());
        assertEquals("0.00", stats.getTotalBalanceBaseCurrency());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordAggregatedStatsData dto = new TgvalidatordAggregatedStatsData();

        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(dto);

        assertNull(stats.getAvgBalancePerAddress());
        assertNull(stats.getAddressesCount());
        assertNull(stats.getWalletsCount());
        assertNull(stats.getTotalBalance());
        assertNull(stats.getTotalBalanceBaseCurrency());
    }

    @Test
    void fromDTO_handlesNullDto() {
        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(null);
        assertNull(stats);
    }

    @Test
    void fromDTO_mapsPartialData() {
        TgvalidatordAggregatedStatsData dto = new TgvalidatordAggregatedStatsData();
        dto.setWalletsCount("10");
        dto.setAddressesCount("50");

        PortfolioStatistics stats = StatisticsMapper.INSTANCE.fromDTO(dto);

        assertEquals("10", stats.getWalletsCount());
        assertEquals("50", stats.getAddressesCount());
        assertNull(stats.getAvgBalancePerAddress());
        assertNull(stats.getTotalBalance());
        assertNull(stats.getTotalBalanceBaseCurrency());
    }
}
