package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.AssetMapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AssetsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAsset;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetAddressesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetAddressesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetWalletsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetWalletsRequest;

import java.util.List;

/**
 * Service for querying asset balances at address and wallet levels.
 * <p>
 * The AssetService provides methods to retrieve addresses and wallets that hold
 * a specific asset (cryptocurrency or token). This is useful for portfolio management,
 * compliance reporting, and understanding asset distribution across the organization.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all addresses holding ETH
 * List<Address> ethAddresses = client.getAssetService().getAssetAddresses("ETH");
 *
 * // Get all wallets holding a specific token
 * List<Wallet> usdcWallets = client.getAssetService().getAssetWallets("USDC");
 *
 * // Get addresses for a specific asset filtered by wallet
 * List<Address> addresses = client.getAssetService().getAssetAddresses("BTC", walletId);
 * }</pre>
 *
 * @see Address
 * @see Wallet
 */
public class AssetService {

    private final AssetsApi assetsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final RulesContainerCache rulesContainerCache;

    /**
     * Creates a new AssetService.
     * <p>
     * Address signature verification is mandatory: getAssetAddresses returns the same
     * Address entity AddressService does, signature and all, and used to hand it over
     * unverified — so the mandatory verification there could be walked around by asking
     * for the same rows here.
     *
     * @param apiClient           the API client for making requests
     * @param apiExceptionMapper  the mapper for converting API exceptions
     * @param rulesContainerCache the cache supplying the HSM key, required
     * @throws NullPointerException if any parameter is null
     */
    public AssetService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper,
                        final RulesContainerCache rulesContainerCache) {
        Preconditions.checkNotNull(apiClient, "apiClient must not be null");
        Preconditions.checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        Preconditions.checkNotNull(rulesContainerCache,
                "rulesContainerCache must not be null - address signature verification is mandatory");
        this.assetsApi = new AssetsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.rulesContainerCache = rulesContainerCache;
    }

    /**
     * Retrieves addresses that hold a specific asset.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @return a list of addresses holding the asset
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public List<Address> getAssetAddresses(final String currency) throws ApiException {
        return getAssetAddresses(currency, null, null, null);
    }

    /**
     * Retrieves addresses that hold a specific asset with optional filters.
     *
     * @param currency  the currency code (e.g., "ETH", "BTC", "USDC")
     * @param walletId  optional wallet ID to filter addresses
     * @param addressId optional address ID to filter
     * @param limit     optional maximum number of addresses to return
     * @return a list of addresses holding the asset
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public List<Address> getAssetAddresses(final String currency, final String walletId,
                                            final String addressId, final String limit) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(currency), "currency must not be null or empty");
        try {
            TgvalidatordAsset asset = new TgvalidatordAsset();
            asset.setCurrency(currency);

            TgvalidatordGetAssetAddressesRequest request = new TgvalidatordGetAssetAddressesRequest();
            request.setAsset(asset);
            request.setWalletId(walletId);
            request.setAddressId(addressId);
            request.setLimit(limit);

            TgvalidatordGetAssetAddressesReply reply = assetsApi.walletServiceGetAssetAddresses(request);
            List<Address> addresses = AssetMapper.INSTANCE.fromAddressDTOList(reply.getAddresses());

            // Through AddressService's seam, not a copy of it. This service reads the
            // same Address entity, so a second implementation of "when is an address
            // string trustworthy" is a second place for it to drift — which is exactly
            // how createAddress ended up with no verification at all. The seam also
            // handles the asynchronous-creation row (empty address, status "creating")
            // that the previous inline loop rejected outright.
            return AddressService.verifiedAddresses(addresses,
                    rulesContainerCache::getDecodedRulesContainer);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves wallets that hold a specific asset.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @return a list of wallets holding the asset
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public List<Wallet> getAssetWallets(final String currency) throws ApiException {
        return getAssetWallets(currency, null);
    }

    /**
     * Retrieves wallets that hold a specific asset with optional filters.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @param limit    optional maximum number of wallets to return
     * @return a list of wallets holding the asset
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public List<Wallet> getAssetWallets(final String currency, final String limit) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(currency), "currency must not be null or empty");
        try {
            TgvalidatordAsset asset = new TgvalidatordAsset();
            asset.setCurrency(currency);

            TgvalidatordGetAssetWalletsRequest request = new TgvalidatordGetAssetWalletsRequest();
            request.setAsset(asset);
            request.setLimit(limit);

            TgvalidatordGetAssetWalletsReply reply = assetsApi.walletServiceGetAssetWallets(request);
            return AssetMapper.INSTANCE.fromWalletInfoDTOList(reply.getWallets());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
