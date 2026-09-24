package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TransactionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.Transaction;
import com.taurushq.sdk.protect.client.model.TransactionExportResult;
import com.taurushq.sdk.protect.client.model.TransactionResult;
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
 * // Get recent transactions, one page at a time
 * TransactionResult page = client.getTransactionService()
 *     .getTransactions(null, null, "ETH", null, 20, 0);
 * // next page: .getTransactions(null, null, "ETH", null, 20, page.getPagination().getNextOffset())
 *
 * // Get transactions for a specific address
 * TransactionResult addrTx = client.getTransactionService()
 *     .getTransactionsByAddress("0x...", 100, 0);
 *
 * // Get a transaction by its blockchain hash
 * Transaction tx = client.getTransactionService()
 *     .getTransactionByHash("0x1234...");
 *
 * // Export up to 1000 transactions as CSV
 * TransactionExportResult csv = client.getTransactionService()
 *     .exportTransactions(startDate, endDate, "ETH", "outgoing", null, null, "csv", 1000);
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
     * Gets a page of transactions.
     *
     * @param from      filter transactions after this date (optional)
     * @param to        filter transactions before this date (optional)
     * @param currency  filter by currency ID or symbol (optional)
     * @param direction filter by direction: "incoming" or "outgoing" (optional)
     * @param limit     the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset    the offset, 0 for the first page
     * @return the transactions and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public TransactionResult getTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                             final String currency, final String direction,
                                             final int limit, final long offset) throws ApiException {
        return getTransactions(from, to, currency, direction, null, null, limit, offset);
    }

    /**
     * Gets a page of transactions, including by chain.
     *
     * <p>The API accepts blockchain and network on this endpoint and the Go and TypeScript
     * SDKs expose both; this SDK used to hardcode them to null, so the filters were
     * unreachable and a caller silently got every chain back.
     *
     * @param from       filter transactions after this date (optional)
     * @param to         filter transactions before this date (optional)
     * @param currency   filter by currency ID or symbol (optional)
     * @param direction  filter by direction: "incoming" or "outgoing" (optional)
     * @param blockchain filter by blockchain, e.g. "ETH" (optional)
     * @param network    filter by network, e.g. "mainnet" (optional)
     * @param limit      the page size, 0 for the default
     * @param offset     the offset, 0 for the first page
     * @return the transactions and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public TransactionResult getTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                             final String currency, final String direction,
                                             final String blockchain, final String network,
                                             final int limit, final long offset) throws ApiException {

        final int size = PagedOperation.TRANSACTIONS.resolveSize("limit", limit);
        final long start = Pagination.resolveOffset("offset", offset);

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    currency,                   // currency
                    direction,                  // direction
                    null,                       // query
                    String.valueOf(size),       // limit
                    start == 0 ? null : String.valueOf(start), // offset
                    from,                       // from
                    to,                         // to
                    null,                       // transactionIds
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    blockchain,                 // blockchain
                    network,                    // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    null,                       // hashes
                    null,                       // address
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null                        // customerId
            );
            return toResult(reply, size, start);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a page of the transactions of an address.
     *
     * @param address the address string (blockchain address)
     * @param limit   the page size, 0 for the default
     * @param offset  the offset, 0 for the first page
     * @return the transactions and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if address is empty or limit or offset is out of range
     */
    public TransactionResult getTransactionsByAddress(final String address, final int limit,
                                                      final long offset) throws ApiException {

        checkNotNull(address, "address cannot be null");
        checkArgument(!address.isEmpty(), "address cannot be empty");
        final int size = PagedOperation.TRANSACTIONS.resolveSize("limit", limit);
        final long start = Pagination.resolveOffset("offset", offset);

        try {
            TgvalidatordGetTransactionsReply reply = transactionsApi.transactionServiceGetTransactions(
                    null,                       // currency
                    null,                       // direction
                    null,                       // query
                    String.valueOf(size),       // limit
                    start == 0 ? null : String.valueOf(start), // offset
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
            return toResult(reply, size, start);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    private static TransactionResult toResult(final TgvalidatordGetTransactionsReply reply,
                                              final int limit, final long offset) {
        List<TgvalidatordTransaction> rows = reply.getResult() == null
                ? Collections.emptyList() : reply.getResult();
        return new TransactionResult(TransactionMapper.INSTANCE.fromDTO(rows),
                PagedOperation.TRANSACTIONS.offsetPage(limit, offset, rows.size(), 0,
                        reply.getTotalItems(), null));
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
     * Exports transactions in the server's default format (JSON).
     * <p>
     * The export cannot page (the server ignores any offset); a larger limit is the only way
     * to export more rows. Compare {@link TransactionExportResult#getTotalItems()} with the
     * limit to tell whether the export was cut short.
     *
     * @param from      filter transactions after this date (optional)
     * @param to        filter transactions before this date (optional)
     * @param currency  filter by currency ID or symbol (optional)
     * @param direction filter by direction: "incoming" or "outgoing" (optional)
     * @param limit     how many transactions to export, 0 for the default
     *                  ({@link Pagination#DEFAULT_PAGE_SIZE}); no SDK maximum
     * @return the exported text and the number of matching transactions
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit is negative
     */
    public TransactionExportResult exportTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                                      final String currency, final String direction,
                                                      final int limit) throws ApiException {
        return exportTransactions(from, to, currency, direction, null, null, null, limit);
    }

    /**
     * Exports transactions in the server's default format (JSON), including by chain.
     *
     * @param from       filter transactions after this date (optional)
     * @param to         filter transactions before this date (optional)
     * @param currency   filter by currency ID or symbol (optional)
     * @param direction  filter by direction: "incoming" or "outgoing" (optional)
     * @param blockchain filter by blockchain, e.g. "ETH" (optional)
     * @param network    filter by network, e.g. "mainnet" (optional)
     * @param limit      how many transactions to export, 0 for the default; no SDK maximum
     * @return the exported text and the number of matching transactions
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit is negative
     */
    public TransactionExportResult exportTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                                      final String currency, final String direction,
                                                      final String blockchain, final String network,
                                                      final int limit) throws ApiException {
        return exportTransactions(from, to, currency, direction, blockchain, network, null, limit);
    }

    /**
     * Exports transactions in the given format, including by chain.
     *
     * @param from       filter transactions after this date (optional)
     * @param to         filter transactions before this date (optional)
     * @param currency   filter by currency ID or symbol (optional)
     * @param direction  filter by direction: "incoming" or "outgoing" (optional)
     * @param blockchain filter by blockchain, e.g. "ETH" (optional)
     * @param network    filter by network, e.g. "mainnet" (optional)
     * @param format     "json", "csv" or "csv_simple", null for the server default (JSON)
     * @param limit      how many transactions to export, 0 for the default; no SDK maximum
     * @return the exported text and the number of matching transactions
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit is negative
     */
    public TransactionExportResult exportTransactions(final OffsetDateTime from, final OffsetDateTime to,
                                                      final String currency, final String direction,
                                                      final String blockchain, final String network,
                                                      final String format, final int limit)
            throws ApiException {

        final int size = PagedOperation.TRANSACTION_EXPORT.resolveSize("limit", limit);

        try {
            TgvalidatordExportTransactionsReply reply = transactionsApi.transactionServiceExportTransactions(
                    currency,                   // currency
                    direction,                  // direction
                    null,                       // query
                    String.valueOf(size),       // limit
                    null,                       // offset: ignored by the server
                    from,                       // from
                    to,                         // to
                    null,                       // transactionIds
                    format,                     // format
                    null,                       // type
                    null,                       // source
                    null,                       // destination
                    null,                       // ids
                    blockchain,                 // blockchain
                    network,                    // network
                    null,                       // fromBlockNumber
                    null,                       // toBlockNumber
                    null,                       // amountAbove
                    null,                       // excludeUnknownSourceDestination
                    null,                       // hashes
                    null                        // address
            );

            return TransactionExportResult.of(reply.getResult(), reply.getTotalItems());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
