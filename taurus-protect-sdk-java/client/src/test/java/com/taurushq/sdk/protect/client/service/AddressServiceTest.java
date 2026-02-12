package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class AddressServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private RulesContainerCache rulesContainerCache;
    private AddressService addressService;

    @BeforeEach
    void setUp() throws Exception {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();

        // Generate a real EC P-256 key to satisfy GovernanceRuleService constructor
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair kp = kpg.generateKeyPair();
        PublicKey publicKey = kp.getPublic();

        GovernanceRuleService governanceRuleService = new GovernanceRuleService(
                apiClient, apiExceptionMapper, Collections.singletonList(publicKey), 1);
        rulesContainerCache = new RulesContainerCache(governanceRuleService);
        addressService = new AddressService(apiClient, apiExceptionMapper, rulesContainerCache);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(null, apiExceptionMapper, rulesContainerCache));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(apiClient, null, rulesContainerCache));
    }

    @Test
    void constructor_throwsOnNullRulesContainerCache() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(apiClient, apiExceptionMapper, null));
    }

    @Test
    void createAddress_throwsOnNullRequest() {
        assertThrows(NullPointerException.class, () ->
                addressService.createAddress(null));
    }

    @Test
    void getAddress_throwsOnZeroId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddress(0));
    }

    @Test
    void getAddress_throwsOnNegativeId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddress(-1));
    }

    @Test
    void getAddresses_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(0, 10, 0));
    }

    @Test
    void getAddresses_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(1, 0, 0));
    }

    @Test
    void getAddresses_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(1, 10, -1));
    }

    @Test
    void createAddressAttribute_throwsOnZeroAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.createAddressAttribute(0, "key", "value"));
    }

    @Test
    void createAddressAttribute_throwsOnEmptyKey() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.createAddressAttribute(1, "", "value"));
    }
}
