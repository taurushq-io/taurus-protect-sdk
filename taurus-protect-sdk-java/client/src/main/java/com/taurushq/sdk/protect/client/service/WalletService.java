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
import com.taurushq.sdk.protect.client.model.Wallet;
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
 * // List wallets with pagination
 * List<Wallet> wallets = client.getWalletService().getWallets(50, 0);
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
     * Gets wallets with pagination.
     *
     * @param limit  the maximum number of wallets to return
     * @param offset the offset for pagination
     * @return the list of wallets
     * @throws ApiException the api exception
     */
    public List<Wallet> getWallets(final int limit, final int offset) throws ApiException {
        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetWalletsInfoReply reply = walletsApi.walletServiceGetWalletsV2(
                    null,                       // currencies
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    null,                       // name
                    null,                       // sortOrder
                    null,                       // excludeDisabled
                    null,                       // tagIDs
                    null,                       // onlyPositiveBalance
                    null,                       // blockchain
                    null,                       // network
                    null                        // ids
            );

            List<TgvalidatordWalletInfo> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return result.stream()
                    .map(WalletMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets wallets by name/query with pagination.
     *
     * @param name  the wallet name
     * @param limit  the maximum number of wallets to return
     * @param offset the offset for pagination
     * @return the list of wallets
     * @throws ApiException the api exception
     */
    public List<Wallet> getWalletsByName(final String name, final int limit, final int offset) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(name), "name cannot be empty");
        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetWalletsInfoReply reply = walletsApi.walletServiceGetWalletsV2(
                    null,                       // currencies
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    name,                       // name
                    null,                       // sortOrder
                    null,                       // excludeDisabled
                    null,                       // tagIDs
                    null,                       // onlyPositiveBalance
                    null,                       // blockchain
                    null,                       // network
                    null                        // ids
            );

            List<TgvalidatordWalletInfo> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }
            return result.stream()
                    .map(WalletMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());
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
     * Gets wallet tokens (asset balances).
     *
     * @param walletId the wallet id
     * @param limit    the maximum number of tokens to return
     * @return the list of asset balances
     * @throws ApiException the api exception
     */
    public List<AssetBalance> getWalletTokens(final long walletId, final int limit) throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");
        checkArgument(limit > 0, "limit must be positive");

        try {
            TgvalidatordGetWalletTokensReply reply = walletsApi.walletServiceGetWalletTokens(
                    String.valueOf(walletId),
                    String.valueOf(limit),
                    null                        // cursor
            );

            List<TgvalidatordAssetBalance> result = reply.getBalances();
            if (result == null) {
                return Collections.emptyList();
            }
            return AssetBalanceMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

}
