package com.taurushq.sdk.protect.client;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.service.ActionService;
import com.taurushq.sdk.protect.client.service.AddressService;
import com.taurushq.sdk.protect.client.service.AirGapService;
import com.taurushq.sdk.protect.client.service.AssetService;
import com.taurushq.sdk.protect.client.service.AuditService;
import com.taurushq.sdk.protect.client.service.BalanceService;
import com.taurushq.sdk.protect.client.service.BlockchainService;
import com.taurushq.sdk.protect.client.service.BusinessRuleService;
import com.taurushq.sdk.protect.client.service.ChangeService;
import com.taurushq.sdk.protect.client.service.ConfigService;
import com.taurushq.sdk.protect.client.service.ContractWhitelistingService;
import com.taurushq.sdk.protect.client.service.CurrencyService;
import com.taurushq.sdk.protect.client.service.ExchangeService;
import com.taurushq.sdk.protect.client.service.FeePayerService;
import com.taurushq.sdk.protect.client.service.FeeService;
import com.taurushq.sdk.protect.client.service.FiatService;
import com.taurushq.sdk.protect.client.service.GovernanceRuleService;
import com.taurushq.sdk.protect.client.service.GroupService;
import com.taurushq.sdk.protect.client.service.HealthService;
import com.taurushq.sdk.protect.client.service.JobService;
import com.taurushq.sdk.protect.client.service.MultiFactorSignatureService;
import com.taurushq.sdk.protect.client.service.PriceService;
import com.taurushq.sdk.protect.client.service.RequestService;
import com.taurushq.sdk.protect.client.service.ReservationService;
import com.taurushq.sdk.protect.client.service.ScoreService;
import com.taurushq.sdk.protect.client.service.StakingService;
import com.taurushq.sdk.protect.client.service.StatisticsService;
import com.taurushq.sdk.protect.client.service.TagService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkLendingService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkParticipantService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkPledgeService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkSettlementService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkSharingService;
import com.taurushq.sdk.protect.client.service.TokenMetadataService;
import com.taurushq.sdk.protect.client.service.TransactionService;
import com.taurushq.sdk.protect.client.service.UserDeviceService;
import com.taurushq.sdk.protect.client.service.UserService;
import com.taurushq.sdk.protect.client.service.VisibilityGroupService;
import com.taurushq.sdk.protect.client.service.WalletService;
import com.taurushq.sdk.protect.client.service.WebhookCallsService;
import com.taurushq.sdk.protect.client.service.WebhookService;
import com.taurushq.sdk.protect.client.service.WhitelistedAddressService;
import com.taurushq.sdk.protect.client.service.WhitelistedAssetService;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Auth;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Exception;
import com.taurushq.sdk.protect.openapi.auth.Authentication;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;

import java.io.IOException;
import java.lang.reflect.Field;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * The main entry point for the Taurus-PROTECT SDK.
 * <p>
 * This client provides access to all Taurus-PROTECT API services including wallets,
 * addresses, requests, transactions, and governance rules.
 * <p>
 * The client implements {@link AutoCloseable} to ensure proper cleanup of sensitive
 * credentials when no longer needed. It is recommended to use try-with-resources:
 * <pre>{@code
 * try (ProtectClient client = ProtectClient.create(host, apiKey, apiSecret, keys, 2)) {
 *     // Use the client
 *     Wallet wallet = client.getWalletService().getWallet(123);
 * }
 * // API secret is securely cleared from memory
 * }</pre>
 */
@SuppressWarnings({"PMD.AvoidUsingVolatile", "PMD.EmptyCatchBlock", "PMD.AvoidAccessibilityAlteration"})
public final class ProtectClient implements AutoCloseable {

    private final ApiClient openApiClient;
    private final List<PublicKey> superAdminPublicKeys;
    private final RulesContainerCache rulesContainerCache;
    private volatile boolean closed = false;
    private final WalletService walletService;
    private final AddressService addressService;
    private final RequestService requestService;
    private final TransactionService transactionService;
    private final CurrencyService currencyService;
    private final ScoreService scoreService;
    private final BalanceService balanceService;
    private final UserService userService;
    private final PriceService priceService;
    private final ChangeService changeService;
    private final BusinessRuleService businessRuleService;
    private final GovernanceRuleService governanceRuleService;
    private final WhitelistedAddressService whitelistedAddressService;
    private final WhitelistedAssetService whitelistedAssetService;
    private final WebhookService webhookService;
    private final StakingService stakingService;
    private final ContractWhitelistingService contractWhitelistingService;
    private final ExchangeService exchangeService;
    private final FeeService feeService;
    private final AuditService auditService;
    private final BlockchainService blockchainService;
    private final TokenMetadataService tokenMetadataService;
    private final StatisticsService statisticsService;
    private final TagService tagService;
    private final GroupService groupService;
    private final HealthService healthService;
    private final JobService jobService;
    private final ReservationService reservationService;
    private final UserDeviceService userDeviceService;
    private final FeePayerService feePayerService;
    private final ConfigService configService;
    private final ActionService actionService;
    private final AssetService assetService;
    private final MultiFactorSignatureService multiFactorSignatureService;
    private final VisibilityGroupService visibilityGroupService;
    private final WebhookCallsService webhookCallsService;
    private final FiatService fiatService;
    private final AirGapService airGapService;
    private final TaurusNetworkClient taurusNetworkClient;


    private ProtectClient(String host, String apiKey, String apiSecret,
                          List<PublicKey> superAdminPublicKeys,
                          int minValidSignatures,
                          long rulesContainerCacheTtlMs) throws ApiKeyTPV1Exception {

        checkArgument(!Strings.isNullOrEmpty(host), "host cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(apiKey), "apiKey cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(apiSecret), "apiSecret cannot be null or empty");
        checkArgument(minValidSignatures > 0, "minValidSignatures must be greater than zero");
        checkArgument(rulesContainerCacheTtlMs > 0, "rulesContainerCacheTtlMs must be positive");

        this.superAdminPublicKeys = superAdminPublicKeys;

        openApiClient = new ApiClient();
        openApiClient.setBasePath(host);
        openApiClient.setApiKeyTPV1(apiKey);
        openApiClient.setApiSecretTPV1(apiSecret);
        openApiClient.setLenientOnJson(false);
        openApiClient.setDebugging(false);

        ApiExceptionMapper apiExceptionMapper = new ApiExceptionMapper();
        this.walletService = new WalletService(openApiClient, apiExceptionMapper);
        this.requestService = new RequestService(openApiClient, apiExceptionMapper);
        this.transactionService = new TransactionService(openApiClient, apiExceptionMapper);
        this.currencyService = new CurrencyService(openApiClient, apiExceptionMapper);
        this.scoreService = new ScoreService(openApiClient, apiExceptionMapper);
        this.balanceService = new BalanceService(openApiClient, apiExceptionMapper);
        this.userService = new UserService(openApiClient, apiExceptionMapper);
        this.priceService = new PriceService(openApiClient, apiExceptionMapper);
        this.changeService = new ChangeService(openApiClient, apiExceptionMapper);
        this.businessRuleService = new BusinessRuleService(openApiClient, apiExceptionMapper);
        this.governanceRuleService = new GovernanceRuleService(openApiClient, apiExceptionMapper, superAdminPublicKeys, minValidSignatures);
        this.whitelistedAddressService = new WhitelistedAddressService(openApiClient, apiExceptionMapper, superAdminPublicKeys, minValidSignatures);
        this.whitelistedAssetService = new WhitelistedAssetService(openApiClient, apiExceptionMapper, superAdminPublicKeys, minValidSignatures);
        this.webhookService = new WebhookService(openApiClient, apiExceptionMapper);
        this.stakingService = new StakingService(openApiClient, apiExceptionMapper);
        this.contractWhitelistingService = new ContractWhitelistingService(openApiClient, apiExceptionMapper);
        this.exchangeService = new ExchangeService(openApiClient, apiExceptionMapper);
        this.feeService = new FeeService(openApiClient, apiExceptionMapper);
        this.auditService = new AuditService(openApiClient, apiExceptionMapper);
        this.blockchainService = new BlockchainService(openApiClient, apiExceptionMapper);
        this.tokenMetadataService = new TokenMetadataService(openApiClient, apiExceptionMapper);
        this.statisticsService = new StatisticsService(openApiClient, apiExceptionMapper);
        this.tagService = new TagService(openApiClient, apiExceptionMapper);
        this.groupService = new GroupService(openApiClient, apiExceptionMapper);
        this.healthService = new HealthService(openApiClient, apiExceptionMapper);
        this.jobService = new JobService(openApiClient, apiExceptionMapper);
        this.reservationService = new ReservationService(openApiClient, apiExceptionMapper);
        this.userDeviceService = new UserDeviceService(openApiClient, apiExceptionMapper);
        this.feePayerService = new FeePayerService(openApiClient, apiExceptionMapper);
        this.configService = new ConfigService(openApiClient, apiExceptionMapper);
        this.actionService = new ActionService(openApiClient, apiExceptionMapper);
        this.assetService = new AssetService(openApiClient, apiExceptionMapper);
        this.multiFactorSignatureService = new MultiFactorSignatureService(openApiClient, apiExceptionMapper);
        this.visibilityGroupService = new VisibilityGroupService(openApiClient, apiExceptionMapper);
        this.webhookCallsService = new WebhookCallsService(openApiClient, apiExceptionMapper);
        this.fiatService = new FiatService(openApiClient, apiExceptionMapper);
        this.airGapService = new AirGapService(openApiClient, apiExceptionMapper);
        this.taurusNetworkClient = new TaurusNetworkClient(
                new TaurusNetworkParticipantService(openApiClient, apiExceptionMapper),
                new TaurusNetworkPledgeService(openApiClient, apiExceptionMapper),
                new TaurusNetworkLendingService(openApiClient, apiExceptionMapper),
                new TaurusNetworkSettlementService(openApiClient, apiExceptionMapper),
                new TaurusNetworkSharingService(openApiClient, apiExceptionMapper)
        );

        // Create rules container cache (depends on governanceRuleService)
        this.rulesContainerCache = new RulesContainerCache(governanceRuleService, rulesContainerCacheTtlMs);

        // Create address service with cache (depends on rulesContainerCache)
        this.addressService = new AddressService(openApiClient, apiExceptionMapper, rulesContainerCache);
    }

    /**
     * Creates a new builder for configuring and creating a ProtectClient.
     * <p>
     * Example usage:
     * <pre>{@code
     * ProtectClient client = ProtectClient.builder()
     *     .host("https://api.protect.taurushq.com")
     *     .credentials(apiKey, apiSecret)
     *     .superAdminKeysPem(pemKeys)
     *     .minValidSignatures(2)
     *     .build();
     * }</pre>
     *
     * @return a new ProtectClientBuilder
     */
    public static ProtectClientBuilder builder() {
        return new ProtectClientBuilder();
    }

    /**
     * Creates a new Protect client with SuperAdmin public keys and default cache TTL.
     *
     * @param host                 the host
     * @param apiKey               the api key
     * @param apiSecret            the api secret
     * @param superAdminPublicKeys the list of SuperAdmin public keys
     * @param minValidSignatures   the minimum number of valid signatures required for governance rules verification
     * @return the protect client
     * @throws ApiKeyTPV1Exception the api key tpv 1 exception
     */
    public static ProtectClient create(String host, String apiKey, String apiSecret,
                                       List<PublicKey> superAdminPublicKeys,
                                       int minValidSignatures) throws ApiKeyTPV1Exception {

        return create(host, apiKey, apiSecret, superAdminPublicKeys, minValidSignatures,
                RulesContainerCache.DEFAULT_CACHE_TTL_MS);
    }

    /**
     * Creates a new Protect client with SuperAdmin public keys and custom cache TTL.
     *
     * @param host                       the host
     * @param apiKey                     the api key
     * @param apiSecret                  the api secret
     * @param superAdminPublicKeys       the list of SuperAdmin public keys
     * @param minValidSignatures         the minimum number of valid signatures required for governance rules verification
     * @param rulesContainerCacheTtlMs   the cache TTL for rules container in milliseconds
     * @return the protect client
     * @throws ApiKeyTPV1Exception the api key tpv 1 exception
     */
    public static ProtectClient create(String host, String apiKey, String apiSecret,
                                       List<PublicKey> superAdminPublicKeys,
                                       int minValidSignatures,
                                       long rulesContainerCacheTtlMs) throws ApiKeyTPV1Exception {

        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys must contain at least 1 key");

        return new ProtectClient(host, apiKey, apiSecret,
                Collections.unmodifiableList(new ArrayList<>(superAdminPublicKeys)), minValidSignatures,
                rulesContainerCacheTtlMs);
    }

    /**
     * Creates a new Protect client with SuperAdmin public keys in PEM format and default cache TTL.
     *
     * @param host                    the host
     * @param apiKey                  the api key
     * @param apiSecret               the api secret
     * @param superAdminPublicKeysPem the list of SuperAdmin public keys in PEM format
     * @param minValidSignatures      the minimum number of valid signatures required for governance rules verification
     * @return the protect client
     * @throws ApiKeyTPV1Exception the api key tpv 1 exception
     * @throws IOException         if a PEM key cannot be decoded
     */
    public static ProtectClient createFromPem(String host, String apiKey, String apiSecret,
                                              List<String> superAdminPublicKeysPem,
                                              int minValidSignatures) throws ApiKeyTPV1Exception, IOException {

        return createFromPem(host, apiKey, apiSecret, superAdminPublicKeysPem, minValidSignatures,
                RulesContainerCache.DEFAULT_CACHE_TTL_MS);
    }

    /**
     * Creates a new Protect client with SuperAdmin public keys in PEM format and custom cache TTL.
     *
     * @param host                       the host
     * @param apiKey                     the api key
     * @param apiSecret                  the api secret
     * @param superAdminPublicKeysPem    the list of SuperAdmin public keys in PEM format
     * @param minValidSignatures         the minimum number of valid signatures required for governance rules verification
     * @param rulesContainerCacheTtlMs   the cache TTL for rules container in milliseconds
     * @return the protect client
     * @throws ApiKeyTPV1Exception the api key tpv 1 exception
     * @throws IOException         if a PEM key cannot be decoded
     */
    public static ProtectClient createFromPem(String host, String apiKey, String apiSecret,
                                              List<String> superAdminPublicKeysPem,
                                              int minValidSignatures,
                                              long rulesContainerCacheTtlMs) throws ApiKeyTPV1Exception, IOException {

        checkNotNull(superAdminPublicKeysPem, "superAdminPublicKeysPem cannot be null");
        checkArgument(!superAdminPublicKeysPem.isEmpty(), "superAdminPublicKeysPem must contain at least 1 key");

        List<PublicKey> keys = new ArrayList<>();
        for (String pem : superAdminPublicKeysPem) {
            keys.add(CryptoTPV1.decodePublicKey(pem));
        }
        return new ProtectClient(host, apiKey, apiSecret, Collections.unmodifiableList(keys), minValidSignatures,
                rulesContainerCacheTtlMs);
    }

    /**
     * returns the internal OpenApi client, allowing for raw API access
     *
     * @return the internal OpenApiClient
     */
    public ApiClient getOpenApiClient() {
        return openApiClient;
    }

    /**
     * Gets wallet service.
     *
     * @return the wallet service
     */
    public WalletService getWalletService() {
        return walletService;
    }

    /**
     * Gets address service.
     *
     * @return the address service
     */
    public AddressService getAddressService() {
        return addressService;
    }

    /**
     * Gets request service.
     *
     * @return the request service
     */
    public RequestService getRequestService() {
        return requestService;
    }

    public TransactionService getTransactionService() {
        return transactionService;
    }

    /**
     * Gets currency service.
     *
     * @return the currency service
     */
    public CurrencyService getCurrencyService() {
        return currencyService;
    }

    /**
     * Gets score service.
     *
     * @return the score service
     */
    public ScoreService getScoreService() {
        return scoreService;
    }

    /**
     * Gets balance service.
     *
     * @return the balance service
     */
    public BalanceService getBalanceService() {
        return balanceService;
    }

    /**
     * Gets user service.
     *
     * @return the user service
     */
    public UserService getUserService() {
        return userService;
    }

    /**
     * Gets price service.
     *
     * @return the price service
     */
    public PriceService getPriceService() {
        return priceService;
    }

    /**
     * Gets change service.
     *
     * @return the change service
     */
    public ChangeService getChangeService() {
        return changeService;
    }

    /**
     * Gets business rule service.
     *
     * @return the business rule service
     */
    public BusinessRuleService getBusinessRuleService() {
        return businessRuleService;
    }

    /**
     * Gets governance rule service.
     *
     * @return the governance rule service
     */
    public GovernanceRuleService getGovernanceRuleService() {
        return governanceRuleService;
    }

    /**
     * Gets the list of SuperAdmin public keys.
     *
     * @return the super admin public keys (unmodifiable)
     */
    public List<PublicKey> getSuperAdminPublicKeys() {
        return superAdminPublicKeys;
    }

    /**
     * Gets whitelisted address service.
     *
     * @return the whitelisted address service
     */
    public WhitelistedAddressService getWhitelistedAddressService() {
        return whitelistedAddressService;
    }

    /**
     * Gets whitelisted asset service (for contract addresses/tokens).
     *
     * @return the whitelisted asset service
     */
    public WhitelistedAssetService getWhitelistedAssetService() {
        return whitelistedAssetService;
    }

    /**
     * Gets webhook service.
     *
     * @return the webhook service
     */
    public WebhookService getWebhookService() {
        return webhookService;
    }

    /**
     * Gets staking service.
     *
     * @return the staking service
     */
    public StakingService getStakingService() {
        return stakingService;
    }

    /**
     * Gets contract whitelisting service.
     *
     * @return the contract whitelisting service
     */
    public ContractWhitelistingService getContractWhitelistingService() {
        return contractWhitelistingService;
    }

    /**
     * Gets exchange service.
     *
     * @return the exchange service
     */
    public ExchangeService getExchangeService() {
        return exchangeService;
    }

    /**
     * Gets fee service.
     *
     * @return the fee service
     */
    public FeeService getFeeService() {
        return feeService;
    }

    /**
     * Gets audit service.
     *
     * @return the audit service
     */
    public AuditService getAuditService() {
        return auditService;
    }

    /**
     * Gets blockchain service.
     *
     * @return the blockchain service
     */
    public BlockchainService getBlockchainService() {
        return blockchainService;
    }

    /**
     * Gets token metadata service.
     *
     * @return the token metadata service
     */
    public TokenMetadataService getTokenMetadataService() {
        return tokenMetadataService;
    }

    /**
     * Gets statistics service.
     *
     * @return the statistics service
     */
    public StatisticsService getStatisticsService() {
        return statisticsService;
    }

    /**
     * Gets the tag service for managing tags.
     *
     * @return the tag service
     */
    public TagService getTagService() {
        return tagService;
    }

    /**
     * Gets the group service for managing user groups.
     *
     * @return the group service
     */
    public GroupService getGroupService() {
        return groupService;
    }

    /**
     * Gets the health service for system health checks.
     *
     * @return the health service
     */
    public HealthService getHealthService() {
        return healthService;
    }

    /**
     * Gets the job service for monitoring background jobs.
     *
     * @return the job service
     */
    public JobService getJobService() {
        return jobService;
    }

    /**
     * Gets the reservation service for UTXO management.
     *
     * @return the reservation service
     */
    public ReservationService getReservationService() {
        return reservationService;
    }

    /**
     * Gets the user device service for device pairing.
     *
     * @return the user device service
     */
    public UserDeviceService getUserDeviceService() {
        return userDeviceService;
    }

    /**
     * Gets the fee payer service for managing fee payers.
     *
     * @return the fee payer service
     */
    public FeePayerService getFeePayerService() {
        return feePayerService;
    }

    /**
     * Gets the config service for tenant configuration.
     *
     * @return the config service
     */
    public ConfigService getConfigService() {
        return configService;
    }

    /**
     * Gets the action service for managing automated actions.
     *
     * @return the action service
     */
    public ActionService getActionService() {
        return actionService;
    }

    /**
     * Gets the asset service for querying addresses and wallets by asset.
     *
     * @return the asset service
     */
    public AssetService getAssetService() {
        return assetService;
    }

    /**
     * Gets the multi-factor signature service for managing multi-party approvals.
     *
     * @return the multi-factor signature service
     */
    public MultiFactorSignatureService getMultiFactorSignatureService() {
        return multiFactorSignatureService;
    }

    /**
     * Gets the visibility group service for managing access control groups.
     *
     * @return the visibility group service
     */
    public VisibilityGroupService getVisibilityGroupService() {
        return visibilityGroupService;
    }

    /**
     * Gets the webhook calls service for retrieving webhook call history.
     *
     * @return the webhook calls service
     */
    public WebhookCallsService getWebhookCallsService() {
        return webhookCallsService;
    }

    /**
     * Gets the fiat service for managing fiat currency operations.
     *
     * @return the fiat service
     */
    public FiatService getFiatService() {
        return fiatService;
    }

    /**
     * Gets the air gap service for cold HSM operations.
     *
     * @return the air gap service
     */
    public AirGapService getAirGapService() {
        return airGapService;
    }

    /**
     * Gets the Taurus Network namespace client providing access to all Taurus Network services.
     * <p>
     * Example usage:
     * <pre>{@code
     * // Get current participant
     * Participant me = client.taurusNetwork().participants().getMyParticipant();
     *
     * // List pledges
     * PledgeResult pledges = client.taurusNetwork().pledges().list(null, null, null, null, null, null);
     *
     * // Get lending agreement
     * LendingAgreement agreement = client.taurusNetwork().lending().getAgreement("agreement-id");
     *
     * // Get settlement
     * Settlement settlement = client.taurusNetwork().settlements().get("settlement-id");
     *
     * // List shared addresses
     * SharedAddressResult addresses = client.taurusNetwork().sharing()
     *     .listSharedAddresses(null, null, null, "ETH", "mainnet", null, null, null);
     * }</pre>
     *
     * @return the Taurus Network namespace client
     */
    public TaurusNetworkClient taurusNetwork() {
        return taurusNetworkClient;
    }

    /**
     * Gets the rules container cache.
     * <p>
     * The cache can be used to manually invalidate the rules container
     * or to check the cache status.
     *
     * @return the rules container cache
     */
    public RulesContainerCache getRulesContainerCache() {
        return rulesContainerCache;
    }

    /**
     * Checks if this client has been closed.
     *
     * @return true if the client has been closed
     */
    public boolean isClosed() {
        return closed;
    }

    /**
     * Closes this client and securely clears the API secret from memory.
     * <p>
     * After calling this method, the client should not be used for API calls.
     * This method is idempotent and can be called multiple times safely.
     * <p>
     * This method zeros out the API secret byte array to minimize the time
     * sensitive credentials remain in memory. While Java's garbage collector
     * may leave copies of the data in memory, this reduces the attack surface.
     */
    @Override
    public synchronized void close() {
        if (closed) {
            return;
        }
        closed = true;
        clearApiSecret();
    }

    /**
     * Clears the API secret from memory using reflection.
     * <p>
     * This method accesses the internal apiSecret byte array in ApiKeyTPV1Auth
     * and zeros it out. This is a best-effort security measure.
     */
    private void clearApiSecret() {
        try {
            for (Authentication auth : openApiClient.getAuthentications().values()) {
                if (auth instanceof ApiKeyTPV1Auth) {
                    clearApiSecretFromAuth((ApiKeyTPV1Auth) auth);
                }
            }
        } catch (Exception e) {
            // Log the error but don't throw - this is a best-effort cleanup
            java.util.logging.Logger.getLogger(ProtectClient.class.getName())
                    .warning("Failed to clear API secret during cleanup: " + e.getMessage()
                            + ". The API secret may remain in memory until garbage collection.");
        }
    }

    private void clearApiSecretFromAuth(ApiKeyTPV1Auth auth) {
        try {
            Field secretField = ApiKeyTPV1Auth.class.getDeclaredField("apiSecret");
            secretField.setAccessible(true);
            byte[] secret = (byte[]) secretField.get(auth);
            if (secret != null) {
                Arrays.fill(secret, (byte) 0);
            }
        } catch (NoSuchFieldException e) {
            // Field structure changed in OpenAPI module - security cleanup failed
            // Log at warning level so operators are aware the secret may remain in memory
            java.util.logging.Logger.getLogger(ProtectClient.class.getName())
                    .warning("Unable to clear API secret: field 'apiSecret' not found in ApiKeyTPV1Auth. "
                            + "The API secret may remain in memory until garbage collection.");
        } catch (IllegalAccessException e) {
            // Security manager denied access - cannot clear secret
            java.util.logging.Logger.getLogger(ProtectClient.class.getName())
                    .warning("Unable to clear API secret: access denied. "
                            + "The API secret may remain in memory until garbage collection.");
        }
    }
}
