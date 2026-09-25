package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.AssetBalanceMapper;
import com.taurushq.sdk.protect.client.mapper.BalanceHistoryPointMapper;
import com.taurushq.sdk.protect.client.mapper.WalletMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.AssetBalance;
import com.taurushq.sdk.protect.client.model.BalanceHistoryPoint;
import com.taurushq.sdk.protect.client.model.CreateWalletRequest;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.client.model.WalletResult;
import com.taurushq.sdk.protect.client.model.WalletTokensResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.WalletsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalanceHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWalletAttributeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWalletReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWalletRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWalletBalanceHistoryReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWalletInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWalletTokensReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWalletsInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import com.taurushq.sdk.protect.openapi.model.WalletServiceCreateWalletAttributesBody;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing blockchain wallets in the Taurus Protect system.
 * <p>
 * This service provides operations for creating, retrieving, and managing wallets.
 * A wallet is a container for one or more addresses on a specific blockchain and network.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a new wallet
 * CreateWalletRequest request = CreateWalletRequest.builder()
 *     .blockchain("ETH")
 *     .network("mainnet")
 *     .name("Trading Wallet")
 *     .omnibus(false)
 *     .build();
 * Wallet wallet = client.getWalletService().createWallet(request);
 *
 * // Retrieve wallet information
 * Wallet wallet = client.getWalletService().getWallet(walletId);
 *
 * // List wallets, one page at a time
 * WalletResult page = client.getWalletService().getWallets(20, 0);
 * while (page.getPagination().hasMore()) {
 *     page = client.getWalletService().getWallets(20, page.getPagination().getNextOffset());
 * }
 * }</pre>
 *
 * @see Wallet
 * @see CreateWalletRequest
 * @see AddressService
 */
public class WalletService {

    /**
     * The underlying OpenAPI client for wallet operations.
     */
    private final WalletsApi walletsApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Instantiates a new Wallet service.
     *
     * @param openApiClient      the open api client
     * @param apiExceptionMapper the api exception mapper
     */
    public WalletService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.walletsApi = new WalletsApi(openApiClient);
    }


    /**
     * Creates a wallet using a request object.
     * <p>
     * This is the recommended method for creating wallets as it provides
     * a cleaner API through the builder pattern:
     * <pre>{@code
     * CreateWalletRequest request = CreateWalletRequest.builder()
     *     .blockchain("ETH")
     *     .network("mainnet")
     *     .name("Trading Wallet")
     *     .omnibus(true)
     *     .comment("Primary account")
     *     .build();
     *
     * Wallet wallet = client.getWalletService().createWallet(request);
     * }</pre>
     *
     * @param request the wallet creation request
     * @return the created wallet
     * @throws ApiException the api exception
     */
    public Wallet createWallet(final CreateWalletRequest request) throws ApiException {
        checkNotNull(request, "request cannot be null");
        return createWallet(
                request.getBlockchain(),
                request.getNetwork(),
                request.getName(),
                request.isOmnibus(),
                request.getComment(),
                request.getCustomerId()
        );
    }

    /**
     * Create wallet.
     *
     * @param blockchain the blockchain
     * @param network    the network
     * @param walletName the wallet name
     * @param isOmnibus  the is omnibus
     * @return the wallet
     * @throws ApiException the api exception
     */
    public Wallet createWallet(final String blockchain, final String network, final String walletName, final boolean isOmnibus) throws ApiException {
        return createWallet(blockchain, network, walletName, isOmnibus, "", "");
    }

    /**
     * Create wallet.
     *
     * @param blockchain the blockchain
     * @param network    the network
     * @param walletName the wallet name
     * @param isOmnibus  the is omnibus
     * @param comment    the comment
     * @return the wallet
     * @throws ApiException the api exception
     */
    public Wallet createWallet(final String blockchain, final String network, final String walletName, final boolean isOmnibus, final String comment) throws ApiException {
        return createWallet(blockchain, network, walletName, isOmnibus, comment, "");
    }


    /**
     * Create wallet.
     *
     * @param blockchain the blockchain
     * @param network    the network
     * @param walletName the wallet name
     * @param isOmnibus  the is omnibus
     * @param comment    the comment
     * @param customerId the customer id
     * @return the wallet
     * @throws ApiException the api exception
     */
    public Wallet createWallet(final String blockchain, final String network, final String walletName, final boolean isOmnibus, final String comment, final String customerId) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(blockchain), "blockchain cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(walletName), "walletName cannot be null or empty");


        TgvalidatordCreateWalletRequest request = new TgvalidatordCreateWalletRequest();
        request.setBlockchain(blockchain);
        request.setNetwork(network);
        request.setName(walletName);
        request.setIsOmnibus(isOmnibus);
        request.setComment(comment);
        request.setCustomerId(customerId);
        try {
            TgvalidatordCreateWalletReply reply = walletsApi.walletServiceCreateWallet(request);
            return WalletMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets wallet.
     *
     * @param walletId the wallet id
     * @return the wallet
     * @throws ApiException the api exception
     */
    public Wallet getWallet(final long walletId) throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");

        try {
            TgvalidatordGetWalletInfoReply result = walletsApi.walletServiceGetWalletV2(String.valueOf(walletId));
            return WalletMapper.INSTANCE.fromDTO(result.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }

    }


    /**
     * Gets a page of wallets.
     *
     * @param limit  the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset the offset, 0 for the first page
     * @return the wallets and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public WalletResult getWallets(final int limit, final long offset) throws ApiException {
        return listWallets(null, limit, offset, null);
    }

    /**
     * Gets a page of wallets, optionally hiding disabled ones.
     * <p>
     * By default the server hides wallets whose currency is disabled; {@code excludeDisabled}
     * true hides every disabled wallet.
     *
     * @param limit           the page size, 0 for the default
     * @param offset          the offset, 0 for the first page
     * @param excludeDisabled true to hide every disabled wallet, null for the server default
     * @return the wallets and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if limit or offset is out of range
     */
    public WalletResult getWallets(final int limit, final long offset, final Boolean excludeDisabled)
            throws ApiException {
        return listWallets(null, limit, offset, excludeDisabled);
    }

    /**
     * Gets a page of wallets by name.
     *
     * @param name   the wallet name
     * @param limit  the page size, 0 for the default
     * @param offset the offset, 0 for the first page
     * @return the wallets and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if name is empty or limit or offset is out of range
     */
    public WalletResult getWalletsByName(final String name, final int limit, final long offset)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(name), "name cannot be empty");
        return listWallets(name, limit, offset, null);
    }

    private WalletResult listWallets(final String name, final int limit, final long offset,
                                     final Boolean excludeDisabled) throws ApiException {
        final int size = PagedOperation.WALLETS.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);

        try {
            TgvalidatordGetWalletsInfoReply reply = walletsApi.walletServiceGetWalletsV2(
                    null,                       // currencies
                    null,                       // query
                    String.valueOf(size),       // limit
                    from == 0 ? null : String.valueOf(from), // offset
                    name,                       // name
                    null,                       // sortOrder
                    excludeDisabled,            // excludeDisabled
                    null,                       // tagIDs
                    null,                       // onlyPositiveBalance
                    null,                       // blockchain
                    null,                       // network
                    null                        // ids
            );

            List<TgvalidatordWalletInfo> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            List<Wallet> wallets = rows.stream()
                    .map(WalletMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());
            return new WalletResult(wallets, PagedOperation.WALLETS.offsetPage(
                    size, from, rows.size(), 0, reply.getTotalItems(), reply.getOffset()));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Creates an attribute for a wallet.
     *
     * @param walletId the wallet id
     * @param key      the attribute key
     * @param value    the attribute value
     * @throws ApiException the api exception
     */
    public void createWalletAttribute(final long walletId, final String key, final String value) throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");
        checkArgument(!Strings.isNullOrEmpty(key), "key cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(value), "value cannot be null or empty");

        try {
            TgvalidatordCreateWalletAttributeRequest attribute = new TgvalidatordCreateWalletAttributeRequest();
            attribute.setKey(key);
            attribute.setValue(value);

            WalletServiceCreateWalletAttributesBody body = new WalletServiceCreateWalletAttributesBody();
            body.addAttributesItem(attribute);

            walletsApi.walletServiceCreateWalletAttributes(String.valueOf(walletId), body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets wallet balance history.
     *
     * @param walletId      the wallet id
     * @param intervalHours the interval in hours for balance snapshots
     * @return the list of balance history points
     * @throws ApiException the api exception
     */
    public List<BalanceHistoryPoint> getWalletBalanceHistory(final long walletId, final int intervalHours) throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");
        checkArgument(intervalHours > 0, "intervalHours must be positive");

        try {
            TgvalidatordGetWalletBalanceHistoryReply reply = walletsApi.walletServiceGetWalletBalanceHistory(
                    String.valueOf(walletId),
                    String.valueOf(intervalHours)
            );

            List<TgvalidatordBalanceHistoryPoint> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return BalanceHistoryPointMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets the first page of a wallet's token balances.
     *
     * @param walletId the wallet id
     * @param pageSize the page size, null or 0 for the default
     * @return the balances and their page
     * @throws ApiException the api exception
     */
    public WalletTokensResult getWalletTokens(final long walletId, final Integer pageSize) throws ApiException {
        return getWalletTokens(walletId, pageSize, null);
    }

    /**
     * Gets a page of a wallet's token balances.
     *
     * @param walletId the wallet id
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the balances and their page
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if the page size is out of range or the cursor is not
     *                                  a cursor this list returned
     */
    public WalletTokensResult getWalletTokens(final long walletId, final Integer pageSize,
                                              final String cursor) throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");
        final int size = PagedOperation.WALLET_TOKENS.resolveSize("pageSize", pageSize);
        final byte[] token = CursorRequest.hasToken(cursor) ? CursorRequest.tokenBytes(cursor) : null;

        try {
            TgvalidatordGetWalletTokensReply reply = walletsApi.walletServiceGetWalletTokens(
                    String.valueOf(walletId),
                    String.valueOf(size),
                    token
            );

            List<TgvalidatordAssetBalance> result = reply.getBalances();
            List<AssetBalance> balances = result == null
                    ? Collections.emptyList() : AssetBalanceMapper.INSTANCE.fromDTO(result);
            return new WalletTokensResult(balances, PagedOperation.WALLET_TOKENS.tokenPage(
                    size, CursorRequest.tokenText(reply.getNext()), reply.getTotal()));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

}
