package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TransactionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Transaction;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TransactionsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExportTransactionsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetTransactionsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTransaction;

import java.time.OffsetDateTime;
import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving and exporting blockchain transactions.
 * <p>
 * This service provides operations for querying transactions recorded in the
 * Taurus Protect system. Transactions can be filtered by date range, currency,
 * direction, address, or blockchain hash.
 * <p>
 * Transactions represent the movement of cryptocurrency on the blockchain and
 * can be either incoming (received) or outgoing (sent).
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get recent transactions
 * List<Transaction> transactions = client.getTransactionService()
 *     .getTransactions(null, null, "ETH", null, 50, 0);
 *
 * // Get transactions for a specific address
 * List<Transaction> addrTx = client.getTransactionService()
 *     .getTransactionsByAddress("0x...", 100, 0);
 *
 * // Get a transaction by its blockchain hash
 * Transaction tx = client.getTransactionService()
 *     .getTransactionByHash("0x1234...");
 *
 * // Export transactions to CSV
 * String csv = client.getTransactionService()
 *     .exportTransactions(startDate, endDate, "ETH", "outgoing", 1000, 0);
 * }</pre>
 *
 * @see Transaction
 * @see RequestService
 */
public class TransactionService {

    /**
     * The underlying OpenAPI client for transaction operations.
     */
    private final TransactionsApi transactionsApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Creates a new TransactionService instance.
     *
     * @param openApiClient      the OpenAPI client for HTTP communication
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public TransactionService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.transactionsApi = new TransactionsApi(openApiClient);
    }


    /**
     * Gets a single transaction by ID.
     *
     * @param id the transaction id
     * @return the transaction
     * @throws ApiException the api exception (including if not found)
     */
    public Transaction getTransactionById(final long id) throws ApiException {
        checkArgument(id > 0, "transaction id cannot be zero");

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    null,                               // currency
                    null,                               // direction
                    null,                               // query
                    "1",                                // limit
                    "0",                                // offset
                    null,                               // from
                    null,                               // to
                    null,                               // transactionIds
                    null,                               // type
                    null,                               // source
                    null,                               // destination
                    Collections.singletonList(String.valueOf(id)), // ids
                    null,                               // blockchain
                    null,                               // network
                    null,                               // fromBlockNumber
                    null,                               // toBlockNumber
                    null,                               // hashes
                    null,                               // address
                    null,                               // amountAbove
                    null,                               // excludeUnknownSourceDestination
                    null                                // customerId
            );

            List<TgvalidatordTransaction> result = reply.getResult();
            if (result == null || result.isEmpty()) {
                ApiException e = new ApiException();
                e.setCode(404);
                e.setError("NotFound");
                e.setMessage(String.format("Transaction with id '%d' not found", id));
                throw e;
            }
            return TransactionMapper.INSTANCE.fromDTO(result.get(0));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets transactions with filtering.
     *
     * @param from      filter transactions after this date (optional)
     * @param to        filter transactions before this date (optional)
     * @param currency  filter by currency ID or symbol (optional)
     * @param direction filter by direction: "incoming" or "outgoing" (optional)
     * @param limit     the maximum number of transactions to return
     * @param offset    the offset for pagination
     * @return the list of transactions
     * @throws ApiException the api exception
     */
    public List<Transaction> getTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                             final String currency, final String direction,
                                             final int limit, final int offset) throws ApiException {

        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    currency,                   // currency
                    direction,                  // direction
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    from,                       // from
                    to,                         // to
                    null,                       // transactionIds
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    null,                       // blockchain
                    null,                       // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    null,                       // hashes
                    null,                       // address
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null                        // customerId
            );

            if (reply.getResult() == null) {
                return Collections.emptyList();
            }
            return TransactionMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets transactions for a specific address.
     *
     * @param address the address string (blockchain address)
     * @param limit   the maximum number of transactions to return
     * @param offset  the offset for pagination
     * @return the list of transactions
     * @throws ApiException the api exception
     */
    public List<Transaction> getTransactionsByAddress(final String address, final int limit, final int offset) throws ApiException {

        checkNotNull(address, "address cannot be null");
        checkArgument(!address.isEmpty(), "address cannot be empty");
        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    null,                       // currency
                    null,                       // direction
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    null,                       // from
                    null,                       // to
                    null,                       // transactionIds
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    null,                       // blockchain
                    null,                       // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    null,                       // hashes
                    address,                    // address
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null                        // customerId
            );

            if (reply.getResult() == null) {
                return Collections.emptyList();
            }
            return TransactionMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a transaction by its blockchain hash.
     *
     * @param hash the transaction hash
     * @return the transaction
     * @throws ApiException the api exception (including if not found)
     */
    public Transaction getTransactionByHash(final String hash) throws ApiException {

        checkNotNull(hash, "hash cannot be null");
        checkArgument(!hash.isEmpty(), "hash cannot be empty");

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    null,                       // currency
                    null,                       // direction
                    null,                       // query
                    "1",                        // limit
                    "0",                        // offset
                    null,                       // from
                    null,                       // to
                    null,                       // transactionIds
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    null,                       // blockchain
                    null,                       // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    Collections.singletonList(hash), // hashes
                    null,                       // address
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null                        // customerId
            );

            List<TgvalidatordTransaction> result = reply.getResult();
            if (result == null || result.isEmpty()) {
                ApiException e = new ApiException();
                e.setCode(404);
                e.setError("NotFound");
                e.setMessage(String.format("Transaction with hash '%s' not found", hash));
                throw e;
            }
            return TransactionMapper.INSTANCE.fromDTO(result.get(0));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Exports transactions to CSV format.
     *
     * @param from      filter transactions after this date (optional)
     * @param to        filter transactions before this date (optional)
     * @param currency  filter by currency ID or symbol (optional)
     * @param direction filter by direction: "incoming" or "outgoing" (optional)
     * @param limit     the maximum number of transactions to export
     * @return the CSV content as a string
     * @throws ApiException the api exception
     */
    public String exportTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                      final String currency, final String direction,
                                      final int limit, final int offset) throws ApiException {

        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordExportTransactionsReply reply = transactionsApi.transactionServiceExportTransactions(
                    currency,                   // currency
                    direction,                  // direction
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    from,                       // from
                    to,                         // to
                    null,                       // transactionIds
                    "csv",                      // format
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    null,                       // blockchain
                    null,                       // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null,                       // hashes
                    null                        // address
            );

            return reply.getResult();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
