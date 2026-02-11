package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.mapper.AssetBalanceMapper;
import com.taurushq.sdk.protect.client.mapper.NFTCollectionBalanceMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.BalanceResult;
import com.taurushq.sdk.protect.client.model.NFTCollectionBalanceResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.BalancesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBalancesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetNFTCollectionBalancesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordNFTCollectionBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;

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
 * // Get all balances with pagination
 * ApiRequestCursor cursor = Pagination.first(50);
 * BalanceResult result = client.getBalanceService().getBalances(cursor);
 *
 * // Get balances for a specific currency
 * BalanceResult ethBalances = client.getBalanceService().getBalances("ETH", cursor);
 *
 * // Get NFT collection balances
 * NFTCollectionBalanceResult nfts = client.getBalanceService()
 *     .getNFTCollectionBalances("ETH", "mainnet", cursor);
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
     * Gets balances for all assets.
     *
     * @param cursor the request cursor for pagination
     * @return the balance result with list and response cursor
     * @throws ApiException the api exception
     */
    public BalanceResult getBalances(final ApiRequestCursor cursor) throws ApiException {
        return getBalances(null, cursor);
    }


    /**
     * Gets balances for a specific currency.
     *
     * @param currency the currency ID or symbol to filter by
     * @param cursor   the request cursor for pagination
     * @return the balance result with list and response cursor
     * @throws ApiException the api exception
     */
    public BalanceResult getBalances(final String currency, final ApiRequestCursor cursor) throws ApiException {
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetBalancesReply reply = balancesApi.walletServiceGetBalances(
                    currency,                                   // currency
                    null,                                       // limit (deprecated)
                    null,                                       // cursor (deprecated)
                    null,                                       // tokenId
                    requestCursor.getCurrentPage(),             // requestCursorCurrentPage
                    requestCursor.getPageRequest(),             // requestCursorPageRequest
                    requestCursor.getPageSize()                 // requestCursorPageSize
            );

            BalanceResult result = new BalanceResult();

            List<TgvalidatordAssetBalance> balances = reply.getBalances();
            if (balances == null) {
                result.setBalances(Collections.emptyList());
            } else {
                result.setBalances(AssetBalanceMapper.INSTANCE.fromDTO(balances));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets NFT collection balances.
     *
     * @param blockchain the blockchain to filter by
     * @param network    the network to filter by
     * @param cursor     the request cursor for pagination
     * @return the NFT collection balance result with list and response cursor
     * @throws ApiException the api exception
     */
    public NFTCollectionBalanceResult getNFTCollectionBalances(final String blockchain, final String network,
                                                               final ApiRequestCursor cursor) throws ApiException {
        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetNFTCollectionBalancesReply reply = balancesApi.walletServiceGetNFTCollectionBalances(
                    blockchain,                                 // blockchain
                    null,                                       // query
                    requestCursor.getCurrentPage(),             // cursorCurrentPage
                    requestCursor.getPageRequest(),             // cursorPageRequest
                    requestCursor.getPageSize(),                // cursorPageSize
                    network,                                    // network
                    null                                        // onlyPositiveBalance
            );


            NFTCollectionBalanceResult result = new NFTCollectionBalanceResult();

            List<TgvalidatordNFTCollectionBalance> balances = reply.getBalances();
            if (balances == null) {
                result.setBalances(Collections.emptyList());
            } else {
                result.setBalances(NFTCollectionBalanceMapper.INSTANCE.fromDTO(balances));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
