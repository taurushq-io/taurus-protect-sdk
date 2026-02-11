package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents the balance of NFTs (Non-Fungible Tokens) in a specific collection.
 * <p>
 * NFT collection balances track how many NFTs from a particular collection
 * are held by an address or wallet. The currency info provides metadata
 * about the collection (contract address, name, etc.) while the balance
 * indicates the count of NFTs held.
 *
 * @see Currency
 * @see Balance
 * @see NFTCollectionBalanceResult
 */
public class NFTCollectionBalance {

    /**
     * Metadata about the NFT collection (contract address, name, blockchain, etc.).
     */
    private Currency currencyInfo;

    /**
     * The balance representing the count of NFTs held from this collection.
     */
    private Balance balance;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the NFT collection metadata.
     *
     * @return the currency info containing collection details
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the NFT collection metadata.
     *
     * @param currencyInfo the collection metadata to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Returns the NFT count balance for this collection.
     *
     * @return the balance representing NFT count held
     */
    public Balance getBalance() {
        return balance;
    }

    /**
     * Sets the NFT count balance for this collection.
     *
     * @param balance the balance to set
     */
    public void setBalance(Balance balance) {
        this.balance = balance;
    }
}
