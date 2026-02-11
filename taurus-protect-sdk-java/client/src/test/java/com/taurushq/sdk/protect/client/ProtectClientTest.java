package com.taurushq.sdk.protect.client;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
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
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.lang.reflect.Field;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Unit tests for {@link ProtectClient}.
 */
class ProtectClientTest {

    private static final String TEST_HOST = "https://api.test.taurushq.com";
    private static final String TEST_API_KEY = "test-api-key-12345";
    // Valid hex string (64 chars = 32 bytes)
    private static final String TEST_API_SECRET = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";

    private List<PublicKey> testSuperAdminKeys;
    private String testPemKey;

    @BeforeEach
    void setUp() throws Exception {
        testSuperAdminKeys = generateTestKeys(2);
        // Generate a valid PEM key dynamically
        testPemKey = CryptoTPV1.encodePublicKey(testSuperAdminKeys.get(0));
    }

    private List<PublicKey> generateTestKeys(int count) throws Exception {
        List<PublicKey> keys = new ArrayList<>();
        KeyPairGenerator keyGen = KeyPairGenerator.getInstance("EC");
        keyGen.initialize(new ECGenParameterSpec("secp256r1"));
        for (int i = 0; i < count; i++) {
            KeyPair keyPair = keyGen.generateKeyPair();
            keys.add(keyPair.getPublic());
        }
        return keys;
    }

    // =====================================================================
    // Construction Tests - create() method
    // =====================================================================

    @Test
    void create_throwsOnNullHost() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(null, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnEmptyHost() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create("", TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnNullApiKey() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, null, TEST_API_SECRET, testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnEmptyApiKey() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, "", TEST_API_SECRET, testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnNullApiSecret() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, null, testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnEmptyApiSecret() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, "", testSuperAdminKeys, 1));
    }

    @Test
    void create_throwsOnNullSuperAdminKeys() {
        assertThrows(NullPointerException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, null, 1));
    }

    @Test
    void create_throwsOnEmptySuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, Collections.emptyList(), 1));
    }

    @Test
    void create_throwsOnZeroMinValidSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 0));
    }

    @Test
    void create_throwsOnNegativeMinValidSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, -1));
    }

    @Test
    void create_throwsOnZeroRulesContainerCacheTtl() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1, 0L));
    }

    @Test
    void create_throwsOnNegativeRulesContainerCacheTtl() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1, -1L));
    }

    @Test
    void create_withValidParametersSucceeds() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client);
            assertFalse(client.isClosed());
        }
    }

    @Test
    void create_withCustomCacheTtlSucceeds() throws ApiKeyTPV1Exception {
        long customTtl = 60_000L; // 1 minute
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1, customTtl)) {
            assertNotNull(client);
            assertEquals(customTtl, client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    // =====================================================================
    // Construction Tests - createFromPem() method
    // =====================================================================

    @Test
    void createFromPem_throwsOnNullHost() {
        List<String> pemKeys = Arrays.asList(testPemKey);
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.createFromPem(null, TEST_API_KEY, TEST_API_SECRET, pemKeys, 1));
    }

    @Test
    void createFromPem_throwsOnEmptyHost() {
        List<String> pemKeys = Arrays.asList(testPemKey);
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.createFromPem("", TEST_API_KEY, TEST_API_SECRET, pemKeys, 1));
    }

    @Test
    void createFromPem_throwsOnNullSuperAdminKeysPem() {
        assertThrows(NullPointerException.class, () ->
                ProtectClient.createFromPem(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, null, 1));
    }

    @Test
    void createFromPem_throwsOnEmptySuperAdminKeysPem() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.createFromPem(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, Collections.emptyList(), 1));
    }

    @Test
    void createFromPem_throwsOnInvalidPemKey() {
        List<String> invalidPemKeys = Arrays.asList("not-a-valid-pem-key");
        assertThrows(IOException.class, () ->
                ProtectClient.createFromPem(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, invalidPemKeys, 1));
    }

    @Test
    void createFromPem_withValidPemKeySucceeds() throws ApiKeyTPV1Exception, IOException {
        List<String> pemKeys = Arrays.asList(testPemKey);
        try (ProtectClient client = ProtectClient.createFromPem(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, pemKeys, 1)) {
            assertNotNull(client);
            assertEquals(1, client.getSuperAdminPublicKeys().size());
        }
    }

    @Test
    void createFromPem_withCustomCacheTtlSucceeds() throws ApiKeyTPV1Exception, IOException {
        List<String> pemKeys = Arrays.asList(testPemKey);
        long customTtl = 120_000L; // 2 minutes
        try (ProtectClient client = ProtectClient.createFromPem(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, pemKeys, 1, customTtl)) {
            assertNotNull(client);
            assertEquals(customTtl, client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    // =====================================================================
    // Service Getter Tests - Core Services (38 services)
    // =====================================================================

    @Test
    void getWalletService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getWalletService());
        }
    }

    @Test
    void getWalletService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            WalletService first = client.getWalletService();
            WalletService second = client.getWalletService();
            assertSame(first, second);
        }
    }

    @Test
    void getAddressService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getAddressService());
        }
    }

    @Test
    void getAddressService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            AddressService first = client.getAddressService();
            AddressService second = client.getAddressService();
            assertSame(first, second);
        }
    }

    @Test
    void getRequestService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getRequestService());
        }
    }

    @Test
    void getRequestService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            RequestService first = client.getRequestService();
            RequestService second = client.getRequestService();
            assertSame(first, second);
        }
    }

    @Test
    void getTransactionService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getTransactionService());
        }
    }

    @Test
    void getTransactionService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TransactionService first = client.getTransactionService();
            TransactionService second = client.getTransactionService();
            assertSame(first, second);
        }
    }

    @Test
    void getCurrencyService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getCurrencyService());
        }
    }

    @Test
    void getCurrencyService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            CurrencyService first = client.getCurrencyService();
            CurrencyService second = client.getCurrencyService();
            assertSame(first, second);
        }
    }

    @Test
    void getScoreService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getScoreService());
        }
    }

    @Test
    void getScoreService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ScoreService first = client.getScoreService();
            ScoreService second = client.getScoreService();
            assertSame(first, second);
        }
    }

    @Test
    void getBalanceService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getBalanceService());
        }
    }

    @Test
    void getBalanceService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            BalanceService first = client.getBalanceService();
            BalanceService second = client.getBalanceService();
            assertSame(first, second);
        }
    }

    @Test
    void getUserService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getUserService());
        }
    }

    @Test
    void getUserService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            UserService first = client.getUserService();
            UserService second = client.getUserService();
            assertSame(first, second);
        }
    }

    @Test
    void getPriceService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getPriceService());
        }
    }

    @Test
    void getPriceService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            PriceService first = client.getPriceService();
            PriceService second = client.getPriceService();
            assertSame(first, second);
        }
    }

    @Test
    void getChangeService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getChangeService());
        }
    }

    @Test
    void getChangeService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ChangeService first = client.getChangeService();
            ChangeService second = client.getChangeService();
            assertSame(first, second);
        }
    }

    @Test
    void getBusinessRuleService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getBusinessRuleService());
        }
    }

    @Test
    void getBusinessRuleService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            BusinessRuleService first = client.getBusinessRuleService();
            BusinessRuleService second = client.getBusinessRuleService();
            assertSame(first, second);
        }
    }

    @Test
    void getGovernanceRuleService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getGovernanceRuleService());
        }
    }

    @Test
    void getGovernanceRuleService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            GovernanceRuleService first = client.getGovernanceRuleService();
            GovernanceRuleService second = client.getGovernanceRuleService();
            assertSame(first, second);
        }
    }

    @Test
    void getWhitelistedAddressService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getWhitelistedAddressService());
        }
    }

    @Test
    void getWhitelistedAddressService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            WhitelistedAddressService first = client.getWhitelistedAddressService();
            WhitelistedAddressService second = client.getWhitelistedAddressService();
            assertSame(first, second);
        }
    }

    @Test
    void getWhitelistedAssetService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getWhitelistedAssetService());
        }
    }

    @Test
    void getWhitelistedAssetService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            WhitelistedAssetService first = client.getWhitelistedAssetService();
            WhitelistedAssetService second = client.getWhitelistedAssetService();
            assertSame(first, second);
        }
    }

    @Test
    void getWebhookService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getWebhookService());
        }
    }

    @Test
    void getWebhookService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            WebhookService first = client.getWebhookService();
            WebhookService second = client.getWebhookService();
            assertSame(first, second);
        }
    }

    @Test
    void getStakingService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getStakingService());
        }
    }

    @Test
    void getStakingService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            StakingService first = client.getStakingService();
            StakingService second = client.getStakingService();
            assertSame(first, second);
        }
    }

    @Test
    void getContractWhitelistingService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getContractWhitelistingService());
        }
    }

    @Test
    void getContractWhitelistingService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ContractWhitelistingService first = client.getContractWhitelistingService();
            ContractWhitelistingService second = client.getContractWhitelistingService();
            assertSame(first, second);
        }
    }

    @Test
    void getExchangeService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getExchangeService());
        }
    }

    @Test
    void getExchangeService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ExchangeService first = client.getExchangeService();
            ExchangeService second = client.getExchangeService();
            assertSame(first, second);
        }
    }

    @Test
    void getFeeService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getFeeService());
        }
    }

    @Test
    void getFeeService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            FeeService first = client.getFeeService();
            FeeService second = client.getFeeService();
            assertSame(first, second);
        }
    }

    @Test
    void getAuditService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getAuditService());
        }
    }

    @Test
    void getAuditService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            AuditService first = client.getAuditService();
            AuditService second = client.getAuditService();
            assertSame(first, second);
        }
    }

    @Test
    void getBlockchainService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getBlockchainService());
        }
    }

    @Test
    void getBlockchainService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            BlockchainService first = client.getBlockchainService();
            BlockchainService second = client.getBlockchainService();
            assertSame(first, second);
        }
    }

    @Test
    void getTokenMetadataService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getTokenMetadataService());
        }
    }

    @Test
    void getTokenMetadataService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TokenMetadataService first = client.getTokenMetadataService();
            TokenMetadataService second = client.getTokenMetadataService();
            assertSame(first, second);
        }
    }

    @Test
    void getStatisticsService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getStatisticsService());
        }
    }

    @Test
    void getStatisticsService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            StatisticsService first = client.getStatisticsService();
            StatisticsService second = client.getStatisticsService();
            assertSame(first, second);
        }
    }

    @Test
    void getTagService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getTagService());
        }
    }

    @Test
    void getTagService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TagService first = client.getTagService();
            TagService second = client.getTagService();
            assertSame(first, second);
        }
    }

    @Test
    void getGroupService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getGroupService());
        }
    }

    @Test
    void getGroupService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            GroupService first = client.getGroupService();
            GroupService second = client.getGroupService();
            assertSame(first, second);
        }
    }

    @Test
    void getHealthService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getHealthService());
        }
    }

    @Test
    void getHealthService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            HealthService first = client.getHealthService();
            HealthService second = client.getHealthService();
            assertSame(first, second);
        }
    }

    @Test
    void getJobService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getJobService());
        }
    }

    @Test
    void getJobService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            JobService first = client.getJobService();
            JobService second = client.getJobService();
            assertSame(first, second);
        }
    }

    @Test
    void getReservationService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getReservationService());
        }
    }

    @Test
    void getReservationService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ReservationService first = client.getReservationService();
            ReservationService second = client.getReservationService();
            assertSame(first, second);
        }
    }

    @Test
    void getUserDeviceService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getUserDeviceService());
        }
    }

    @Test
    void getUserDeviceService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            UserDeviceService first = client.getUserDeviceService();
            UserDeviceService second = client.getUserDeviceService();
            assertSame(first, second);
        }
    }

    @Test
    void getFeePayerService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getFeePayerService());
        }
    }

    @Test
    void getFeePayerService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            FeePayerService first = client.getFeePayerService();
            FeePayerService second = client.getFeePayerService();
            assertSame(first, second);
        }
    }

    @Test
    void getConfigService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getConfigService());
        }
    }

    @Test
    void getConfigService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ConfigService first = client.getConfigService();
            ConfigService second = client.getConfigService();
            assertSame(first, second);
        }
    }

    @Test
    void getActionService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getActionService());
        }
    }

    @Test
    void getActionService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            ActionService first = client.getActionService();
            ActionService second = client.getActionService();
            assertSame(first, second);
        }
    }

    @Test
    void getAssetService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getAssetService());
        }
    }

    @Test
    void getAssetService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            AssetService first = client.getAssetService();
            AssetService second = client.getAssetService();
            assertSame(first, second);
        }
    }

    @Test
    void getMultiFactorSignatureService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getMultiFactorSignatureService());
        }
    }

    @Test
    void getMultiFactorSignatureService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            MultiFactorSignatureService first = client.getMultiFactorSignatureService();
            MultiFactorSignatureService second = client.getMultiFactorSignatureService();
            assertSame(first, second);
        }
    }

    @Test
    void getVisibilityGroupService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getVisibilityGroupService());
        }
    }

    @Test
    void getVisibilityGroupService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            VisibilityGroupService first = client.getVisibilityGroupService();
            VisibilityGroupService second = client.getVisibilityGroupService();
            assertSame(first, second);
        }
    }

    @Test
    void getWebhookCallsService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getWebhookCallsService());
        }
    }

    @Test
    void getWebhookCallsService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            WebhookCallsService first = client.getWebhookCallsService();
            WebhookCallsService second = client.getWebhookCallsService();
            assertSame(first, second);
        }
    }

    @Test
    void getFiatService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getFiatService());
        }
    }

    @Test
    void getFiatService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            FiatService first = client.getFiatService();
            FiatService second = client.getFiatService();
            assertSame(first, second);
        }
    }

    @Test
    void getAirGapService_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getAirGapService());
        }
    }

    @Test
    void getAirGapService_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            AirGapService first = client.getAirGapService();
            AirGapService second = client.getAirGapService();
            assertSame(first, second);
        }
    }

    // =====================================================================
    // TaurusNetwork Client and Sub-Services Tests
    // =====================================================================

    @Test
    void taurusNetwork_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork());
        }
    }

    @Test
    void taurusNetwork_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkClient first = client.taurusNetwork();
            TaurusNetworkClient second = client.taurusNetwork();
            assertSame(first, second);
        }
    }

    @Test
    void taurusNetwork_participants_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork().participants());
        }
    }

    @Test
    void taurusNetwork_participants_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkParticipantService first = client.taurusNetwork().participants();
            TaurusNetworkParticipantService second = client.taurusNetwork().participants();
            assertSame(first, second);
        }
    }

    @Test
    void taurusNetwork_pledges_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork().pledges());
        }
    }

    @Test
    void taurusNetwork_pledges_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkPledgeService first = client.taurusNetwork().pledges();
            TaurusNetworkPledgeService second = client.taurusNetwork().pledges();
            assertSame(first, second);
        }
    }

    @Test
    void taurusNetwork_lending_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork().lending());
        }
    }

    @Test
    void taurusNetwork_lending_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkLendingService first = client.taurusNetwork().lending();
            TaurusNetworkLendingService second = client.taurusNetwork().lending();
            assertSame(first, second);
        }
    }

    @Test
    void taurusNetwork_settlements_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork().settlements());
        }
    }

    @Test
    void taurusNetwork_settlements_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkSettlementService first = client.taurusNetwork().settlements();
            TaurusNetworkSettlementService second = client.taurusNetwork().settlements();
            assertSame(first, second);
        }
    }

    @Test
    void taurusNetwork_sharing_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.taurusNetwork().sharing());
        }
    }

    @Test
    void taurusNetwork_sharing_returnsSameInstance() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            TaurusNetworkSharingService first = client.taurusNetwork().sharing();
            TaurusNetworkSharingService second = client.taurusNetwork().sharing();
            assertSame(first, second);
        }
    }

    // =====================================================================
    // Lifecycle Tests
    // =====================================================================

    @Test
    void isClosed_initiallyFalse() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertFalse(client.isClosed());
        }
    }

    @Test
    void close_setsClosedToTrue() throws ApiKeyTPV1Exception {
        ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1);
        assertFalse(client.isClosed());
        client.close();
        assertTrue(client.isClosed());
    }

    @Test
    void close_isIdempotent() throws ApiKeyTPV1Exception {
        ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1);
        client.close();
        assertTrue(client.isClosed());

        // Second close should not throw
        client.close();
        assertTrue(client.isClosed());

        // Third close for good measure
        client.close();
        assertTrue(client.isClosed());
    }

    @Test
    @SuppressWarnings("PMD.AvoidAccessibilityAlteration")
    void close_clearsApiSecret() throws Exception {
        ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1);

        // Get reference to API secret before close
        ApiClient apiClient = client.getOpenApiClient();
        ApiKeyTPV1Auth auth = null;
        for (Object authObj : apiClient.getAuthentications().values()) {
            if (authObj instanceof ApiKeyTPV1Auth) {
                auth = (ApiKeyTPV1Auth) authObj;
                break;
            }
        }
        assertNotNull(auth, "Should find ApiKeyTPV1Auth");

        // Access the apiSecret field via reflection
        Field secretField = ApiKeyTPV1Auth.class.getDeclaredField("apiSecret");
        secretField.setAccessible(true);
        byte[] secretBefore = (byte[]) secretField.get(auth);
        assertNotNull(secretBefore, "API secret should exist before close");

        // Close the client
        client.close();

        // Verify secret is cleared (all zeros)
        byte[] secretAfter = (byte[]) secretField.get(auth);
        if (secretAfter != null) {
            boolean allZeros = true;
            for (byte b : secretAfter) {
                if (b != 0) {
                    allZeros = false;
                    break;
                }
            }
            assertTrue(allZeros, "API secret should be cleared to all zeros after close");
        }
    }

    // =====================================================================
    // Reflection Target Validation Tests
    // =====================================================================

    @Test
    void secretCleanupReflectionTarget_fieldExists() throws Exception {
        // Validate that the field we rely on for secret cleanup in ProtectClient.close()
        // actually exists in ApiKeyTPV1Auth. If OpenAPI regeneration changes the class
        // internals, this test will catch the breakage immediately.
        Class<?> authClass = Class.forName("com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Auth");
        Field secretField = authClass.getDeclaredField("apiSecret");
        assertNotNull(secretField);
        assertEquals(byte[].class, secretField.getType());
    }

    // =====================================================================
    // Configuration Tests
    // =====================================================================

    @Test
    void getSuperAdminPublicKeys_returnsImmutableList() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            List<PublicKey> keys = client.getSuperAdminPublicKeys();
            assertNotNull(keys);
            assertEquals(testSuperAdminKeys.size(), keys.size());

            // Verify the list is unmodifiable
            assertThrows(UnsupportedOperationException.class, () ->
                    keys.add(testSuperAdminKeys.get(0)));
        }
    }

    @Test
    void getOpenApiClient_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getOpenApiClient());
        }
    }

    @Test
    void getRulesContainerCache_returnsNonNull() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertNotNull(client.getRulesContainerCache());
        }
    }

    @Test
    void getOpenApiClient_hasCorrectBasePath() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertEquals(TEST_HOST, client.getOpenApiClient().getBasePath());
        }
    }

    @Test
    void getRulesContainerCache_hasDefaultTtlWhenNotSpecified() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.create(TEST_HOST, TEST_API_KEY, TEST_API_SECRET, testSuperAdminKeys, 1)) {
            assertEquals(RulesContainerCache.DEFAULT_CACHE_TTL_MS, client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    // =====================================================================
    // Builder Access Test
    // =====================================================================

    @Test
    void builder_returnsNonNullBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertNotNull(builder);
    }
}
