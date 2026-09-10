package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class AssetServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private RulesContainerCache rulesContainerCache;
    private AssetService assetService;

    @BeforeAll
    static void setUpProvider() {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    @BeforeEach
    void setUp() throws Exception {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        rulesContainerCache = new RulesContainerCache(governanceService());
        assetService = new AssetService(apiClient, apiExceptionMapper, rulesContainerCache);
    }

    private GovernanceRuleService governanceService() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        PublicKey key = generator.generateKeyPair().getPublic();
        return new GovernanceRuleService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(key), 1);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AssetService(null, apiExceptionMapper, rulesContainerCache));
    }

    @Test
    void constructor_throwsOnNullApiExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AssetService(apiClient, null, rulesContainerCache));
    }

    // Asset addresses are the same entity AddressService verifies, so the cache that
    // supplies the HSM key is mandatory here too.
    @Test
    void constructor_throwsOnNullRulesContainerCache() {
        assertThrows(NullPointerException.class, () ->
                new AssetService(apiClient, apiExceptionMapper, null));
    }

    @Test
    void getAssetAddresses_throwsOnNullCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetAddresses(null));
    }

    @Test
    void getAssetAddresses_throwsOnEmptyCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetAddresses(""));
    }

    @Test
    void getAssetWallets_throwsOnNullCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetWallets(null));
    }

    @Test
    void getAssetWallets_throwsOnEmptyCurrency() {
        assertThrows(IllegalArgumentException.class, () ->
                assetService.getAssetWallets(""));
    }
}
