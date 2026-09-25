package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.AssetMapper;
import com.taurushq.sdk.protect.client.mapper.AssetV2Mapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.AssetAddressV2;
import com.taurushq.sdk.protect.client.model.AssetAddressV2Result;
import com.taurushq.sdk.protect.client.model.AssetAddressesResult;
import com.taurushq.sdk.protect.client.model.AssetOperationV2Result;
import com.taurushq.sdk.protect.client.model.AssetV2Result;
import com.taurushq.sdk.protect.client.model.AssetWalletsResult;
import com.taurushq.sdk.protect.client.model.ExcludedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AssetV2Api;
import com.taurushq.sdk.protect.openapi.api.AssetsApi;
import com.taurushq.sdk.protect.openapi.model.AssetServiceV2QueryAssetAddressesV2Body;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressTypeV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAsset;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetAddressesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetAddressesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetWalletsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetWalletsRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetsReplyV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAssetsRequestV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordKYCStatusV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordListAssetOperationsReplyV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordQueryAssetAddressesReplyV2;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Service for querying asset balances at address and wallet levels.
 * <p>
 * The AssetService provides methods to retrieve addresses and wallets that hold
 * a specific asset (cryptocurrency or token). This is useful for portfolio management,
 * compliance reporting, and understanding asset distribution across the organization.
 * <p>
 * Example usage:
 * <pre>{@code
 * // First page of the addresses holding ETH
 * AssetAddressesResult ethAddresses = client.getAssetService().getAssetAddresses("ETH");
 *
 * // First page of the wallets holding a specific token
 * AssetWalletsResult usdcWallets = client.getAssetService().getAssetWallets("USDC");
 *
 * // Addresses holding BTC in one wallet, 20 at a time
 * AssetAddressesResult page = client.getAssetService().getAssetAddresses("BTC", walletId, null, 20, null);
 * // next page: getAssetAddresses("BTC", walletId, null, 20, page.getPage().getNextCursor())
 *
 * // v2 asset service: query assets, then an asset's addresses and operations
 * AssetV2Result assets = client.getAssetService().queryAssets("ETH", "mainnet", null, null, null, null, 20, null);
 * }</pre>
 *
 * @see com.taurushq.sdk.protect.client.model.Address
 * @see com.taurushq.sdk.protect.client.model.Wallet
 */
public class AssetService {

    private static final String INTERNAL = TgvalidatordAddressTypeV2.INTERNAL.getValue();
    private static final String WHITELISTED = TgvalidatordAddressTypeV2.WHITELISTED.getValue();

    private final AssetsApi assetsApi;
    private final AssetV2Api assetV2Api;
    private final ApiExceptionMapper apiExceptionMapper;
    private final RulesContainerCache rulesContainerCache;
    private final AddressService addressService;
    private final WhitelistedAddressService whitelistedAddressService;

    /**
     * Creates a new AssetService.
     * <p>
     * Address signature verification is mandatory: getAssetAddresses returns the same
     * Address entity AddressService does, signature and all, and used to hand it over
     * unverified — so the mandatory verification there could be walked around by asking
     * for the same rows here.
     * <p>
     * queryAssetAddresses rows carry no signature at all, so they are confirmed through the
     * two verified readers passed in, rather than through a second copy of the HSM or
     * whitelist verification.
     *
     * @param apiClient                 the API client for making requests
     * @param apiExceptionMapper        the mapper for converting API exceptions
     * @param rulesContainerCache       the cache supplying the HSM key, required
     * @param addressService            the verified managed-address reader, required
     * @param whitelistedAddressService the verified whitelisted-address reader, required
     * @throws NullPointerException if any parameter is null
     */
    public AssetService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper,
                        final RulesContainerCache rulesContainerCache,
                        final AddressService addressService,
                        final WhitelistedAddressService whitelistedAddressService) {
        Preconditions.checkNotNull(apiClient, "apiClient must not be null");
        Preconditions.checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        Preconditions.checkNotNull(rulesContainerCache,
                "rulesContainerCache must not be null - address signature verification is mandatory");
        Preconditions.checkNotNull(addressService,
                "addressService must not be null - asset holders are verified through it");
        Preconditions.checkNotNull(whitelistedAddressService,
                "whitelistedAddressService must not be null - asset holders are verified through it");
        this.assetsApi = new AssetsApi(apiClient);
        this.assetV2Api = new AssetV2Api(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.rulesContainerCache = rulesContainerCache;
        this.addressService = addressService;
        this.whitelistedAddressService = whitelistedAddressService;
    }

    /**
     * Retrieves the first page of the addresses that hold an asset, with the default page size.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @return the verified addresses and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public AssetAddressesResult getAssetAddresses(final String currency) throws ApiException {
        return getAssetAddresses(currency, null, null, (ApiRequestCursor) null);
    }

    /**
     * Retrieves a page of the addresses that hold an asset.
     *
     * @param currency  the currency code (e.g., "ETH", "BTC", "USDC")
     * @param walletId  optional wallet ID to filter addresses
     * @param addressId optional address ID to filter
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the verified addresses and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty, or the page size is out of range
     */
    public AssetAddressesResult getAssetAddresses(final String currency, final String walletId,
                                                  final String addressId, final Integer pageSize,
                                                  final String cursor) throws ApiException {
        return getAssetAddresses(currency, walletId, addressId, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of the addresses that hold an asset, with a low-level request cursor.
     *
     * @param currency  the currency code (e.g., "ETH", "BTC", "USDC")
     * @param walletId  optional wallet ID to filter addresses
     * @param addressId optional address ID to filter
     * @param cursor    the request cursor, null for the first page with the default size
     * @return the verified addresses and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public AssetAddressesResult getAssetAddresses(final String currency, final String walletId,
                                                  final String addressId, final ApiRequestCursor cursor)
            throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(currency), "currency must not be null or empty");
        final CursorRequest page = CursorRequest.of(cursor);
        try {
            TgvalidatordAsset asset = new TgvalidatordAsset();
            asset.setCurrency(currency);

            TgvalidatordGetAssetAddressesRequest request = new TgvalidatordGetAssetAddressesRequest();
            request.setAsset(asset);
            request.setWalletId(walletId);
            request.setAddressId(addressId);
            request.setAddresses(null);
            request.setRequestCursor(page.toDTO());

            TgvalidatordGetAssetAddressesReply reply = assetsApi.walletServiceGetAssetAddresses(request);
            List<Address> addresses = reply.getAddresses() == null
                    ? Collections.emptyList() : AssetMapper.INSTANCE.fromAddressDTOList(reply.getAddresses());

            // Through AddressService's seam, not a copy of it. This service reads the
            // same Address entity, so a second implementation of "when is an address
            // string trustworthy" is a second place for it to drift — which is exactly
            // how createAddress ended up with no verification at all. The seam also
            // handles the asynchronous-creation row (empty address, status "creating")
            // that the previous inline loop rejected outright.
            AssetAddressesResult result = new AssetAddressesResult();
            result.setAddresses(AddressService.verifiedAddresses(addresses,
                    rulesContainerCache::getDecodedRulesContainer));
            return page.complete(result, PagedOperation.ASSET_ADDRESSES, reply.getCursor(),
                    reply.getTotalItems());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves the first page of the wallets that hold an asset, with the default page size.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @return the wallets and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public AssetWalletsResult getAssetWallets(final String currency) throws ApiException {
        return getAssetWallets(currency, (ApiRequestCursor) null);
    }

    /**
     * Retrieves a page of the wallets that hold an asset.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the wallets and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty, or the page size is out of range
     */
    public AssetWalletsResult getAssetWallets(final String currency, final Integer pageSize,
                                              final String cursor) throws ApiException {
        return getAssetWallets(currency, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of the wallets that hold an asset, with a low-level request cursor.
     *
     * @param currency the currency code (e.g., "ETH", "BTC", "USDC")
     * @param cursor   the request cursor, null for the first page with the default size
     * @return the wallets and their page, with the server total
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if currency is null or empty
     */
    public AssetWalletsResult getAssetWallets(final String currency, final ApiRequestCursor cursor)
            throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(currency), "currency must not be null or empty");
        final CursorRequest page = CursorRequest.of(cursor);
        try {
            TgvalidatordAsset asset = new TgvalidatordAsset();
            asset.setCurrency(currency);

            TgvalidatordGetAssetWalletsRequest request = new TgvalidatordGetAssetWalletsRequest();
            request.setAsset(asset);
            request.setRequestCursor(page.toDTO());

            TgvalidatordGetAssetWalletsReply reply = assetsApi.walletServiceGetAssetWallets(request);
            AssetWalletsResult result = new AssetWalletsResult();
            result.setWallets(reply.getWallets() == null
                    ? Collections.emptyList() : AssetMapper.INSTANCE.fromWalletInfoDTOList(reply.getWallets()));
            return page.complete(result, PagedOperation.ASSET_WALLETS, reply.getCursor(), reply.getTotalItems());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Queries a page of assets (v2 asset service).
     *
     * @param blockchain      filter by blockchain (optional)
     * @param network         filter by network (optional)
     * @param symbol          filter by symbol (optional)
     * @param contractAddress filter by contract address (optional)
     * @param label           filter by label (optional)
     * @param currencyName    filter by currency name (optional)
     * @param pageSize        the page size, null or 0 for the default
     * @param cursor          a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the assets and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public AssetV2Result queryAssets(final String blockchain, final String network, final String symbol,
                                     final String contractAddress, final String label,
                                     final String currencyName, final Integer pageSize,
                                     final String cursor) throws ApiException {
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));
        try {
            TgvalidatordGetAssetsRequestV2 request = new TgvalidatordGetAssetsRequestV2();
            request.setCursor(page.toDTO());
            request.setBlockchain(blockchain);
            request.setNetwork(network);
            request.setSymbol(symbol);
            request.setContractAddress(contractAddress);
            request.setLabel(label);
            request.setCurrencyName(currencyName);

            TgvalidatordGetAssetsReplyV2 reply = assetV2Api.assetServiceV2QueryAssetsV2(request);
            AssetV2Result result = new AssetV2Result();
            result.setAssets(AssetV2Mapper.INSTANCE.fromAssetDTOList(reply.getResult()));
            return page.complete(result, PagedOperation.ASSETS_V2, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Queries a page of the addresses of an asset (v2 asset service).
     * <p>
     * The rows carry no signature, so each INTERNAL and WHITELISTED row is confirmed against
     * its verified counterpart and returned with {@code isVerified() == true} and the
     * counterpart's address; a row that cannot be confirmed is withheld and named in
     * {@code getExcludedUnverified()}. Every other row (EXTERNAL, or no type) is returned
     * with {@code isVerified() == false}. Exclusions never move the cursor.
     *
     * @param assetId     the asset id
     * @param addressType filter by address type, as its wire value (ADDRESS_TYPE_V2_INTERNAL,
     *                    ADDRESS_TYPE_V2_WHITELISTED, ADDRESS_TYPE_V2_EXTERNAL), or null; sent
     *                    verbatim, so the server decides on a value this SDK does not know
     * @param kycStatus   filter by KYC status, as its wire value (KYC_STATUS_V2_APPROVED,
     *                    KYC_STATUS_V2_REVOKED), or null; sent verbatim like addressType
     * @param pageSize    the page size, null or 0 for the default
     * @param cursor      a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the asset's addresses, the rows withheld, and their page
     * @throws ApiException             if the API call or a verified re-read fails
     * @throws WhitelistException       if a whitelisted-address rules container fails verification
     * @throws IntegrityException       if rows came back and none could be verified, or a rules
     *                                  container cannot verify any address
     * @throws IllegalArgumentException if assetId is empty or the page size is out of range
     */
    public AssetAddressV2Result queryAssetAddresses(final String assetId, final String addressType,
                                                    final String kycStatus, final Integer pageSize,
                                                    final String cursor)
            throws ApiException, WhitelistException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(assetId), "assetId must not be null or empty");
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));
        AssetServiceV2QueryAssetAddressesV2Body body = new AssetServiceV2QueryAssetAddressesV2Body();
        body.setCursor(page.toDTO());
        body.setAddressType(TgvalidatordAddressTypeV2.fromValue(addressType));
        body.setKycStatus(TgvalidatordKYCStatusV2.fromValue(kycStatus));
        TgvalidatordQueryAssetAddressesReplyV2 reply;
        try {
            reply = assetV2Api.assetServiceV2QueryAssetAddressesV2(assetId, body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
        // Exclusions below never move the cursor.
        AssetAddressV2Result result = page.complete(new AssetAddressV2Result(),
                PagedOperation.ASSET_ADDRESSES_V2, reply.getCursor(), null);
        List<AssetAddressV2> rows = AssetV2Mapper.INSTANCE.fromAssetAddressDTOList(reply.getResult());
        List<ExcludedWhitelistedAddress> excluded = new ArrayList<>();
        result.setAddresses(verifiedHolders(rows == null ? Collections.emptyList() : rows, excluded));
        result.setExcludedUnverified(excluded);
        return result;
    }

    /**
     * Confirms a page of v2 asset holders through the verified readers.
     * <pre>
     *   INTERNAL    -&gt; addressID            -&gt; verified managed-address read (HSM, &lt;= 50 ids) --+
     *   WHITELISTED -&gt; whitelistedAddressID -&gt; verified whitelist read (6 steps, &lt;= 100 ids) --+
     *                                    same address string? -- yes -&gt; verified, address from the reader
     *                                                          +- no  -&gt; excluded
     *   any other type -----------------------------------------------&gt; verified false
     * </pre>
     * A reader is called only when the page has a row of its type. Its call-level failures
     * propagate; its row-level ones become exclusions here. Rows came back but none survived
     * is an error, so a filtered page never reads as an empty one.
     *
     * @param rows     the mapped holder rows
     * @param excluded receives every withheld row
     * @return the kept rows
     * @throws ApiException       if a verified re-read fails
     * @throws WhitelistException if a whitelisted-address rules container fails verification
     */
    private List<AssetAddressV2> verifiedHolders(final List<AssetAddressV2> rows,
                                                 final List<ExcludedWhitelistedAddress> excluded)
            throws ApiException, WhitelistException {
        Set<String> internalIds = new LinkedHashSet<>();
        Set<String> whitelistedIds = new LinkedHashSet<>();
        for (AssetAddressV2 row : rows) {
            if (INTERNAL.equals(row.getAddressType()) && hasPlatformId(row.getAddressId())) {
                internalIds.add(row.getAddressId());
            } else if (WHITELISTED.equals(row.getAddressType())
                    && hasPlatformId(row.getWhitelistedAddressId())) {
                whitelistedIds.add(row.getWhitelistedAddressId());
            }
        }

        Map<String, String> internal = new HashMap<>();
        Map<String, String> internalFailed = new HashMap<>();
        if (!internalIds.isEmpty()) {
            for (Map.Entry<String, Address> entry
                    : addressService.verifiedAddressesById(new ArrayList<>(internalIds), internalFailed).entrySet()) {
                internal.put(entry.getKey(), entry.getValue().getAddress());
            }
        }
        Map<String, String> whitelisted = new HashMap<>();
        Map<String, String> whitelistedFailed = new HashMap<>();
        if (!whitelistedIds.isEmpty()) {
            List<ExcludedWhitelistedAddress> failures = new ArrayList<>();
            Map<String, SignedWhitelistedAddressEnvelope> envelopes =
                    whitelistedAddressService.verifiedEnvelopesById(new ArrayList<>(whitelistedIds), null, failures);
            for (ExcludedWhitelistedAddress failure : failures) {
                whitelistedFailed.put(failure.getId(), failure.getReason());
            }
            for (Map.Entry<String, SignedWhitelistedAddressEnvelope> entry : envelopes.entrySet()) {
                String address = entry.getValue().getWhitelistedAddress().getAddress();
                if (Strings.isNullOrEmpty(address)) {
                    whitelistedFailed.put(entry.getKey(), "its signed payload carries no address");
                } else {
                    whitelisted.put(entry.getKey(), address);
                }
            }
        }

        List<AssetAddressV2> kept = new ArrayList<>(rows.size());
        for (AssetAddressV2 row : rows) {
            ExcludedWhitelistedAddress exclusion;
            if (INTERNAL.equals(row.getAddressType())) {
                exclusion = confirm(row, row.getAddressId(), "addressID", internal, internalFailed,
                        "managed address");
            } else if (WHITELISTED.equals(row.getAddressType())) {
                exclusion = confirm(row, row.getWhitelistedAddressId(), "whitelistedAddressID",
                        whitelisted, whitelistedFailed, "whitelisted address");
            } else {
                // An EXTERNAL holder, or an untyped row: on-chain data nothing signs.
                kept.add(row);
                continue;
            }
            if (exclusion == null) {
                kept.add(row);
            } else {
                excluded.add(exclusion);
            }
        }

        if (!rows.isEmpty() && kept.isEmpty()) {
            throw new IntegrityException(String.format(
                    "all %d asset address(es) failed verification; first failure: %s",
                    rows.size(), excluded.get(0).getReason()));
        }
        return kept;
    }

    /**
     * Confirms one INTERNAL or WHITELISTED holder against the addresses its verified reader
     * returned, marking it verified with the reader's address.
     *
     * @param row        the holder row
     * @param platformId the row's addressID or whitelistedAddressID
     * @param idField    the name of that field, for the reason
     * @param verified   the reader's verified addresses by id
     * @param failed     why each row the reader rejected failed, by id
     * @param kind       what the reader reads, for the reason
     * @return null when the row is confirmed, otherwise why it is withheld
     */
    private static ExcludedWhitelistedAddress confirm(final AssetAddressV2 row, final String platformId,
                                                      final String idField,
                                                      final Map<String, String> verified,
                                                      final Map<String, String> failed,
                                                      final String kind) {
        if (!hasPlatformId(platformId)) {
            return new ExcludedWhitelistedAddress(row.getAddress(),
                    "the row carries no " + idField + " to verify it by");
        }
        String address = verified.get(platformId);
        if (address == null) {
            String why = failed.get(platformId);
            return new ExcludedWhitelistedAddress(platformId, why == null
                    ? "the verified " + kind + " read did not return " + idField + " " + platformId
                    : "the " + kind + " did not verify: " + why);
        }
        if (!address.equals(row.getAddress())) {
            return new ExcludedWhitelistedAddress(platformId,
                    "the row's address differs from the verified " + kind + " " + platformId);
        }
        row.markVerified(address);
        return null;
    }

    /**
     * Reports whether a holder row carries a platform id; validatord sends zero, or omits
     * the field, when the address is not managed or not whitelisted.
     */
    private static boolean hasPlatformId(final String id) {
        return !Strings.isNullOrEmpty(id) && !"0".equals(id);
    }

    /**
     * Lists a page of the operations on an asset (v2 asset service).
     *
     * @param assetId  the asset id
     * @param type     filter by operation type, as its wire value (e.g.
     *                 ASSET_OPERATION_TYPE_V2_MINT), or null
     * @param status   filter by operation status, as its wire value (e.g.
     *                 ASSET_OPERATION_STATUS_V2_COMPLETED), or null
     * @param pageSize the page size, null or 0 for the default
     * @param cursor   a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the operations and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if assetId is empty or the page size is out of range
     */
    public AssetOperationV2Result listAssetOperations(final String assetId, final String type,
                                                      final String status, final Integer pageSize,
                                                      final String cursor) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(assetId), "assetId must not be null or empty");
        final CursorRequest page = CursorRequest.of(Pagination.page(pageSize, cursor));
        try {
            TgvalidatordListAssetOperationsReplyV2 reply = assetV2Api.assetServiceV2ListAssetOperationsV2(
                    assetId,
                    page.currentPage(),
                    page.pageRequest(),
                    page.pageSizeParam(),
                    type,
                    status
            );
            AssetOperationV2Result result = new AssetOperationV2Result();
            result.setOperations(AssetV2Mapper.INSTANCE.fromAssetOperationDTOList(reply.getResult()));
            return page.complete(result, PagedOperation.ASSET_OPERATIONS_V2, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
