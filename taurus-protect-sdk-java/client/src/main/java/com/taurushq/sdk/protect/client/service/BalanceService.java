package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.AssetBalanceMapper;
import com.taurushq.sdk.protect.client.mapper.NFTCollectionBalanceMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.BalanceResult;
import com.taurushq.sdk.protect.client.model.NFTCollectionBalanceResult;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.BalancesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBalancesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetNFTCollectionBalancesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordNFTCollectionBalance;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving asset balances in the Taurus Protect system.
 * <p>
 * This service provides operations for querying balances across wallets and addresses,
 * including fungible tokens and NFT collections. Balances can be filtered by currency
 * or blockchain.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all balances, one page at a time
 * BalanceResult result = client.getBalanceService().getBalances(null, 20, null);
 * while (result.getPage().hasMore()) {
 *     result = client.getBalanceService().getBalances(null, 20, result.getPage().getNextCursor());
 * }
 *
 * // Get balances for a specific currency
 * BalanceResult ethBalances = client.getBalanceService().getBalances("ETH", 20, null);
 *
 * // Get NFT collection balances
 * NFTCollectionBalanceResult nfts = client.getBalanceService()
 *     .getNFTCollectionBalances("ETH", "mainnet", 20, null);
 * }</pre>
 *
 * @see BalanceResult
 * @see NFTCollectionBalanceResult
 * @see WalletService
 */
public class BalanceService {

    /**
     * The underlying OpenAPI client for balance operations.
     */
    private final BalancesApi balancesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Balance service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public BalanceService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.balancesApi = new BalancesApi(openApiClient);
    }


    /**
     * Gets a page of balances for all assets, with a low-level request cursor.
     *
     * @param cursor the request cursor, null for the first page with the default size
     * @return the balances and their page, with the server total
     * @throws ApiException the api exception
     */
    public BalanceResult getBalances(final ApiRequestCursor cursor) throws ApiException {
        return getBalances(null, cursor);
    }


    /**
     * Gets a page of balances, optionally for one currency.
     *
     * @param currency the currency ID or symbol to filter by, or null for every asset
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the balances and their page, with the server total
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public BalanceResult getBalances(final String currency, final Integer pageSize, final String cursor)
            throws ApiException {
        return getBalances(currency, Pagination.page(pageSize, cursor));
    }


    /**
     * Gets a page of balances, optionally for one currency, with a low-level request cursor.
     *
     * @param currency the currency ID or symbol to filter by, or null for every asset
     * @param cursor   the request cursor, null for the first page with the default size
     * @return the balances and their page, with the server total
     * @throws ApiException the api exception
     */
    public BalanceResult getBalances(final String currency, final ApiRequestCursor cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetBalancesReply reply = balancesApi.walletServiceGetBalances(
                    currency,                                   // currency
                    null,                                       // limit: legacy, requestCursor only
                    null,                                       // cursor: legacy, requestCursor only
                    null,                                       // tokenId
                    page.currentPage(),                         // requestCursorCurrentPage
                    page.pageRequest(),                         // requestCursorPageRequest
                    page.pageSizeParam()                        // requestCursorPageSize
            );

            BalanceResult result = new BalanceResult();
            List<TgvalidatordAssetBalance> balances = reply.getBalances();
            result.setBalances(balances == null
                    ? Collections.emptyList() : AssetBalanceMapper.INSTANCE.fromDTO(balances));
            return page.complete(result, PagedOperation.BALANCES, reply.getCursor(), reply.getTotal());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a page of NFT collection balances.
     *
     * @param blockchain the blockchain to filter by (optional)
     * @param network    the network to filter by (optional)
     * @param pageSize   the page size, null or 0 for the default
     * @param cursor     a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the NFT collection balances and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range
     */
    public NFTCollectionBalanceResult getNFTCollectionBalances(final String blockchain, final String network,
                                                               final Integer pageSize, final String cursor)
            throws ApiException {
        return getNFTCollectionBalances(blockchain, network, Pagination.page(pageSize, cursor));
    }


    /**
     * Gets a page of NFT collection balances, with a low-level request cursor.
     *
     * @param blockchain the blockchain to filter by (optional)
     * @param network    the network to filter by (optional)
     * @param cursor     the request cursor, null for the first page with the default size
     * @return the NFT collection balances and their page
     * @throws ApiException the api exception
     */
    public NFTCollectionBalanceResult getNFTCollectionBalances(final String blockchain, final String network,
                                                               final ApiRequestCursor cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);

        try {
            TgvalidatordGetNFTCollectionBalancesReply reply = balancesApi.walletServiceGetNFTCollectionBalances(
                    blockchain,                                 // blockchain
                    null,                                       // query
                    page.currentPage(),                         // cursorCurrentPage
                    page.pageRequest(),                         // cursorPageRequest
                    page.pageSizeParam(),                       // cursorPageSize
                    network,                                    // network
                    null                                        // onlyPositiveBalance
            );

            NFTCollectionBalanceResult result = new NFTCollectionBalanceResult();
            List<TgvalidatordNFTCollectionBalance> balances = reply.getBalances();
            result.setBalances(balances == null
                    ? Collections.emptyList() : NFTCollectionBalanceMapper.INSTANCE.fromDTO(balances));
            return page.complete(result, PagedOperation.NFT_COLLECTION_BALANCES, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
