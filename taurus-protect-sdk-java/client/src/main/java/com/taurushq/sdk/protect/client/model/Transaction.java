package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.math.BigInteger;
import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a blockchain transaction in the Taurus Protect system.
 * <p>
 * A transaction records the movement of cryptocurrency between addresses on a blockchain.
 * Transactions can be incoming (received) or outgoing (sent), and they contain information
 * about the source, destination, amount, fees, and blockchain confirmation status.
 * <p>
 * Transactions are linked to requests when they are initiated through the Taurus Protect
 * system. Incoming transactions may not have an associated request.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get transactions for an address
 * List<Transaction> transactions = client.getTransactionService().getTransactionsByAddress(addressId, 100, 0);
 * for (Transaction tx : transactions) {
 *     System.out.println("Hash: " + tx.getHash());
 *     System.out.println("Amount: " + tx.getAmountMainUnit() + " " + tx.getCurrency());
 *     System.out.println("Direction: " + tx.getDirection());
 * }
 * }</pre>
 *
 * @see TransactionService
 * @see Request
 * @see AddressInfo
 */
public class Transaction {

    /**
     * The unique identifier of the transaction assigned by Taurus Protect.
     */
    private long id;

    /**
     * The ID of the request associated with this transaction, if any.
     */
    private long requestId;

    /**
     * The direction of the transaction ("incoming" or "outgoing").
     */
    private String direction;

    /**
     * The network identifier (e.g., "mainnet", "testnet").
     */
    private String network;

    /**
     * The blockchain type (e.g., "ETH", "BTC", "SOL").
     */
    private String blockchain;

    /**
     * The currency code for the transaction (e.g., "ETH", "BTC").
     */
    private String currency;

    /**
     * Detailed information about the transaction's currency.
     */
    private Currency currencyInfo;

    /**
     * The transaction amount in the smallest unit (e.g., wei for ETH, satoshi for BTC).
     */
    private BigInteger amount;

    /**
     * The transaction amount in the main unit (e.g., ETH, BTC).
     */
    private double amountMainUnit;

    /**
     * The list of source addresses for the transaction.
     */
    private List<AddressInfo> sources;

    /**
     * The list of destination addresses for the transaction.
     */
    private List<AddressInfo> destinations;

    /**
     * The transaction fee in the smallest unit (e.g., wei for ETH).
     */
    private BigInteger fee;

    /**
     * The transaction fee in the main unit (e.g., ETH).
     */
    private double feeMainUnit;

    /**
     * The blockchain transaction hash.
     */
    private String hash;

    /**
     * The block number where the transaction was included.
     */
    private long block;

    /**
     * The block number at which the transaction reached the required confirmations.
     */
    private long confirmationBlock;

    /**
     * The date and time when the transaction was first detected.
     */
    private OffsetDateTime receptionDate;

    /**
     * The date and time when the transaction was confirmed on the blockchain.
     */
    private OffsetDateTime confirmationDate;

    /**
     * An alternative transaction identifier.
     */
    private String transactionId;

    /**
     * The type of transaction (e.g., "transfer", "stake", "contract_call").
     */
    private String type;

    /**
     * A unique identifier for the transaction across the system.
     */
    private String uniqueId;

    /**
     * An optional argument field for additional transaction data.
     */
    private String arg1;

    /**
     * An optional argument field for additional transaction data.
     */
    private String arg2;

    /**
     * Indicates whether the associated request is visible to the current user.
     */
    private boolean requestVisible;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the transaction.
     *
     * @return the transaction ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the transaction.
     *
     * @param id the transaction ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the ID of the request associated with this transaction.
     *
     * @return the request ID, or 0 if no request is associated
     */
    public long getRequestId() {
        return requestId;
    }

    /**
     * Sets the ID of the request associated with this transaction.
     *
     * @param requestId the request ID to set
     */
    public void setRequestId(long requestId) {
        this.requestId = requestId;
    }

    /**
     * Returns the direction of the transaction.
     *
     * @return "incoming" for received transactions, "outgoing" for sent transactions
     */
    public String getDirection() {
        return direction;
    }

    /**
     * Sets the direction of the transaction.
     *
     * @param direction the direction to set ("incoming" or "outgoing")
     */
    public void setDirection(String direction) {
        this.direction = direction;
    }

    /**
     * Returns the network identifier (e.g., "mainnet", "testnet").
     *
     * @return the network identifier
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network identifier.
     *
     * @param network the network identifier to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Returns the blockchain type (e.g., "ETH", "BTC", "SOL").
     *
     * @return the blockchain identifier
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain type.
     *
     * @param blockchain the blockchain identifier to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Returns the currency code for the transaction (e.g., "ETH", "BTC").
     *
     * @return the currency code
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code for the transaction.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns detailed information about the transaction's currency.
     *
     * @return the currency information, or {@code null} if not available
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the detailed currency information for the transaction.
     *
     * @param currencyInfo the currency information to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Returns the transaction amount in the smallest unit.
     * <p>
     * For example, wei for ETH, satoshi for BTC.
     *
     * @return the amount in the smallest unit
     */
    public BigInteger getAmount() {
        return amount;
    }

    /**
     * Sets the transaction amount in the smallest unit.
     *
     * @param amount the amount to set
     */
    public void setAmount(BigInteger amount) {
        this.amount = amount;
    }

    /**
     * Returns the transaction amount in the main unit.
     * <p>
     * For example, ETH instead of wei, BTC instead of satoshi.
     *
     * @return the amount in the main unit
     */
    public double getAmountMainUnit() {
        return amountMainUnit;
    }

    /**
     * Sets the transaction amount in the main unit.
     *
     * @param amountMainUnit the amount to set
     */
    public void setAmountMainUnit(double amountMainUnit) {
        this.amountMainUnit = amountMainUnit;
    }

    /**
     * Returns the list of source addresses for the transaction.
     *
     * @return the list of source addresses, or {@code null} if not available
     */
    public List<AddressInfo> getSources() {
        return sources;
    }

    /**
     * Sets the list of source addresses for the transaction.
     *
     * @param sources the list of source addresses to set
     */
    public void setSources(List<AddressInfo> sources) {
        this.sources = sources;
    }

    /**
     * Returns the list of destination addresses for the transaction.
     *
     * @return the list of destination addresses, or {@code null} if not available
     */
    public List<AddressInfo> getDestinations() {
        return destinations;
    }

    /**
     * Sets the list of destination addresses for the transaction.
     *
     * @param destinations the list of destination addresses to set
     */
    public void setDestinations(List<AddressInfo> destinations) {
        this.destinations = destinations;
    }

    /**
     * Returns the transaction fee in the smallest unit.
     *
     * @return the fee in the smallest unit (e.g., wei for ETH)
     */
    public BigInteger getFee() {
        return fee;
    }

    /**
     * Sets the transaction fee in the smallest unit.
     *
     * @param fee the fee to set
     */
    public void setFee(BigInteger fee) {
        this.fee = fee;
    }

    /**
     * Returns the transaction fee in the main unit.
     *
     * @return the fee in the main unit (e.g., ETH)
     */
    public double getFeeMainUnit() {
        return feeMainUnit;
    }

    /**
     * Sets the transaction fee in the main unit.
     *
     * @param feeMainUnit the fee to set
     */
    public void setFeeMainUnit(double feeMainUnit) {
        this.feeMainUnit = feeMainUnit;
    }

    /**
     * Returns the blockchain transaction hash.
     * <p>
     * This is the unique identifier of the transaction on the blockchain.
     *
     * @return the transaction hash
     */
    public String getHash() {
        return hash;
    }

    /**
     * Sets the blockchain transaction hash.
     *
     * @param hash the transaction hash to set
     */
    public void setHash(String hash) {
        this.hash = hash;
    }

    /**
     * Returns the block number where the transaction was included.
     *
     * @return the block number, or 0 if the transaction is not yet included in a block
     */
    public long getBlock() {
        return block;
    }

    /**
     * Sets the block number where the transaction was included.
     *
     * @param block the block number to set
     */
    public void setBlock(long block) {
        this.block = block;
    }

    /**
     * Returns the block number at which the transaction reached required confirmations.
     *
     * @return the confirmation block number, or 0 if not yet confirmed
     */
    public long getConfirmationBlock() {
        return confirmationBlock;
    }

    /**
     * Sets the confirmation block number.
     *
     * @param confirmationBlock the confirmation block number to set
     */
    public void setConfirmationBlock(long confirmationBlock) {
        this.confirmationBlock = confirmationBlock;
    }

    /**
     * Returns the date and time when the transaction was first detected.
     *
     * @return the reception date
     */
    public OffsetDateTime getReceptionDate() {
        return receptionDate;
    }

    /**
     * Sets the date and time when the transaction was first detected.
     *
     * @param receptionDate the reception date to set
     */
    public void setReceptionDate(OffsetDateTime receptionDate) {
        this.receptionDate = receptionDate;
    }

    /**
     * Returns the date and time when the transaction was confirmed on the blockchain.
     *
     * @return the confirmation date, or {@code null} if not yet confirmed
     */
    public OffsetDateTime getConfirmationDate() {
        return confirmationDate;
    }

    /**
     * Sets the date and time when the transaction was confirmed.
     *
     * @param confirmationDate the confirmation date to set
     */
    public void setConfirmationDate(OffsetDateTime confirmationDate) {
        this.confirmationDate = confirmationDate;
    }

    /**
     * Returns an alternative transaction identifier.
     *
     * @return the transaction ID, or {@code null} if not set
     */
    public String getTransactionId() {
        return transactionId;
    }

    /**
     * Sets the alternative transaction identifier.
     *
     * @param transactionId the transaction ID to set
     */
    public void setTransactionId(String transactionId) {
        this.transactionId = transactionId;
    }

    /**
     * Returns the type of transaction (e.g., "transfer", "stake", "contract_call").
     *
     * @return the transaction type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the type of transaction.
     *
     * @param type the transaction type to set
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Returns a unique identifier for the transaction across the system.
     *
     * @return the unique ID
     */
    public String getUniqueId() {
        return uniqueId;
    }

    /**
     * Sets the unique identifier for the transaction.
     *
     * @param uniqueId the unique ID to set
     */
    public void setUniqueId(String uniqueId) {
        this.uniqueId = uniqueId;
    }

    /**
     * Returns the first optional argument field for additional transaction data.
     *
     * @return the first argument, or {@code null} if not set
     */
    public String getArg1() {
        return arg1;
    }

    /**
     * Sets the first optional argument field.
     *
     * @param arg1 the first argument to set
     */
    public void setArg1(String arg1) {
        this.arg1 = arg1;
    }

    /**
     * Returns the second optional argument field for additional transaction data.
     *
     * @return the second argument, or {@code null} if not set
     */
    public String getArg2() {
        return arg2;
    }

    /**
     * Sets the second optional argument field.
     *
     * @param arg2 the second argument to set
     */
    public void setArg2(String arg2) {
        this.arg2 = arg2;
    }

    /**
     * Returns whether the associated request is visible to the current user.
     *
     * @return {@code true} if the request is visible, {@code false} otherwise
     */
    public boolean isRequestVisible() {
        return requestVisible;
    }

    /**
     * Sets whether the associated request is visible to the current user.
     *
     * @param requestVisible {@code true} if the request should be visible
     */
    public void setRequestVisible(boolean requestVisible) {
        this.requestVisible = requestVisible;
    }
}
