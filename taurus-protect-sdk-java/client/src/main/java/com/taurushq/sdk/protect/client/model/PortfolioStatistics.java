package com.taurushq.sdk.protect.client.model;

/**
 * Represents portfolio statistics in the Taurus Protect system.
 * <p>
 * Contains aggregated statistics about the portfolio including
 * total balance, address counts, and wallet counts.
 *
 * @see StatisticsService
 */
public class PortfolioStatistics {

    private String avgBalancePerAddress;
    private String addressesCount;
    private String walletsCount;
    private String totalBalance;
    private String totalBalanceBaseCurrency;

    /**
     * Gets the average balance per address in the smallest currency unit.
     *
     * @return the average balance per address
     */
    public String getAvgBalancePerAddress() {
        return avgBalancePerAddress;
    }

    /**
     * Sets the average balance per address.
     *
     * @param avgBalancePerAddress the average balance to set
     */
    public void setAvgBalancePerAddress(String avgBalancePerAddress) {
        this.avgBalancePerAddress = avgBalancePerAddress;
    }

    /**
     * Gets the total number of addresses.
     *
     * @return the addresses count
     */
    public String getAddressesCount() {
        return addressesCount;
    }

    /**
     * Sets the addresses count.
     *
     * @param addressesCount the count to set
     */
    public void setAddressesCount(String addressesCount) {
        this.addressesCount = addressesCount;
    }

    /**
     * Gets the total number of wallets.
     *
     * @return the wallets count
     */
    public String getWalletsCount() {
        return walletsCount;
    }

    /**
     * Sets the wallets count.
     *
     * @param walletsCount the count to set
     */
    public void setWalletsCount(String walletsCount) {
        this.walletsCount = walletsCount;
    }

    /**
     * Gets the total balance in the smallest currency unit.
     *
     * @return the total balance
     */
    public String getTotalBalance() {
        return totalBalance;
    }

    /**
     * Sets the total balance.
     *
     * @param totalBalance the total balance to set
     */
    public void setTotalBalance(String totalBalance) {
        this.totalBalance = totalBalance;
    }

    /**
     * Gets the total balance valuation in the base currency (CHF, EUR, USD, etc.).
     *
     * @return the total balance in base currency
     */
    public String getTotalBalanceBaseCurrency() {
        return totalBalanceBaseCurrency;
    }

    /**
     * Sets the total balance in base currency.
     *
     * @param totalBalanceBaseCurrency the total balance to set
     */
    public void setTotalBalanceBaseCurrency(String totalBalanceBaseCurrency) {
        this.totalBalanceBaseCurrency = totalBalanceBaseCurrency;
    }
}
