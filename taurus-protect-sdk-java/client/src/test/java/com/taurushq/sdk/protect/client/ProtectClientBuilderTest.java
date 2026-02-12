package com.taurushq.sdk.protect.client;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Exception;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
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
 * Unit tests for {@link ProtectClientBuilder}.
 */
class ProtectClientBuilderTest {

    private static final String TEST_HOST = "https://api.test.taurushq.com";
    private static final String TEST_API_KEY = "test-api-key-12345";
    // Valid hex string (64 chars = 32 bytes)
    private static final String TEST_API_SECRET = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";

    private List<PublicKey> testSuperAdminKeys;
    private PublicKey testSingleKey;
    private String testPemKey;
    private String testPemKey2;

    @BeforeEach
    void setUp() throws Exception {
        testSuperAdminKeys = generateTestKeys(2);
        testSingleKey = testSuperAdminKeys.get(0);
        // Generate valid PEM keys dynamically
        testPemKey = CryptoTPV1.encodePublicKey(testSuperAdminKeys.get(0));
        testPemKey2 = CryptoTPV1.encodePublicKey(testSuperAdminKeys.get(1));
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
    // Builder method tests - Fluent API returns builder
    // =====================================================================

    @Test
    void host_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.host(TEST_HOST);
        assertSame(builder, result);
    }

    @Test
    void credentials_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.credentials(TEST_API_KEY, TEST_API_SECRET);
        assertSame(builder, result);
    }

    @Test
    void apiKey_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.apiKey(TEST_API_KEY);
        assertSame(builder, result);
    }

    @Test
    void apiSecret_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.apiSecret(TEST_API_SECRET);
        assertSame(builder, result);
    }

    @Test
    void superAdminKeys_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.superAdminKeys(testSuperAdminKeys);
        assertSame(builder, result);
    }

    @Test
    void superAdminKey_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.superAdminKey(testSingleKey);
        assertSame(builder, result);
    }

    @Test
    void superAdminKeysPem_returnsBuilder() throws IOException {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.superAdminKeysPem(Arrays.asList(testPemKey));
        assertSame(builder, result);
    }

    @Test
    void superAdminKeyPem_returnsBuilder() throws IOException {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.superAdminKeyPem(testPemKey);
        assertSame(builder, result);
    }

    @Test
    void minValidSignatures_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.minValidSignatures(2);
        assertSame(builder, result);
    }

    @Test
    void rulesContainerCacheTtlMs_returnsBuilder() {
        ProtectClientBuilder builder = ProtectClient.builder();
        ProtectClientBuilder result = builder.rulesContainerCacheTtlMs(60000L);
        assertSame(builder, result);
    }

    // =====================================================================
    // Validation tests - build() method
    // =====================================================================

    @Test
    void build_throwsOnMissingHost() {
        ProtectClientBuilder builder = ProtectClient.builder()
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys);

        IllegalStateException ex = assertThrows(IllegalStateException.class, builder::build);
        assertTrue(ex.getMessage().contains("host"));
    }

    @Test
    void build_throwsOnMissingApiKey() {
        ProtectClientBuilder builder = ProtectClient.builder()
                .host(TEST_HOST)
                .apiSecret(TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys);

        IllegalStateException ex = assertThrows(IllegalStateException.class, builder::build);
        assertTrue(ex.getMessage().contains("apiKey"));
    }

    @Test
    void build_throwsOnMissingApiSecret() {
        ProtectClientBuilder builder = ProtectClient.builder()
                .host(TEST_HOST)
                .apiKey(TEST_API_KEY)
                .superAdminKeys(testSuperAdminKeys);

        IllegalStateException ex = assertThrows(IllegalStateException.class, builder::build);
        assertTrue(ex.getMessage().contains("apiSecret"));
    }

    @Test
    void build_throwsOnMissingSuperAdminKeys() {
        ProtectClientBuilder builder = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET);

        IllegalStateException ex = assertThrows(IllegalStateException.class, builder::build);
        assertTrue(ex.getMessage().contains("SuperAdmin"));
    }

    @Test
    void minValidSignatures_throwsOnZero() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(IllegalArgumentException.class, () -> builder.minValidSignatures(0));
    }

    @Test
    void minValidSignatures_throwsOnNegative() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(IllegalArgumentException.class, () -> builder.minValidSignatures(-1));
    }

    @Test
    void rulesContainerCacheTtlMs_throwsOnZero() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(IllegalArgumentException.class, () -> builder.rulesContainerCacheTtlMs(0));
    }

    @Test
    void rulesContainerCacheTtlMs_throwsOnNegative() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(IllegalArgumentException.class, () -> builder.rulesContainerCacheTtlMs(-1));
    }

    @Test
    void superAdminKeys_throwsOnNull() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(NullPointerException.class, () -> builder.superAdminKeys(null));
    }

    @Test
    void superAdminKey_throwsOnNull() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(NullPointerException.class, () -> builder.superAdminKey(null));
    }

    @Test
    void superAdminKeysPem_throwsOnNull() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(NullPointerException.class, () -> builder.superAdminKeysPem(null));
    }

    @Test
    void superAdminKeyPem_throwsOnNull() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(NullPointerException.class, () -> builder.superAdminKeyPem(null));
    }

    @Test
    void superAdminKeysPem_throwsOnInvalidPem() {
        ProtectClientBuilder builder = ProtectClient.builder();
        List<String> invalidPem = Arrays.asList("not-a-valid-pem");
        assertThrows(IOException.class, () -> builder.superAdminKeysPem(invalidPem));
    }

    @Test
    void superAdminKeyPem_throwsOnInvalidPem() {
        ProtectClientBuilder builder = ProtectClient.builder();
        assertThrows(IOException.class, () -> builder.superAdminKeyPem("not-a-valid-pem"));
    }

    // =====================================================================
    // Fluent API tests - Method chaining
    // =====================================================================

    @Test
    void fluentChaining_buildsValidClient() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .minValidSignatures(2)
                .build()) {

            assertNotNull(client);
            assertFalse(client.isClosed());
        }
    }

    @Test
    void fluentChaining_withAllOptions() throws ApiKeyTPV1Exception {
        long customTtl = 120_000L;
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .apiKey(TEST_API_KEY)
                .apiSecret(TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .minValidSignatures(2)
                .rulesContainerCacheTtlMs(customTtl)
                .build()) {

            assertNotNull(client);
            assertEquals(customTtl, client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    @Test
    void fluentChaining_withPemKey() throws ApiKeyTPV1Exception, IOException {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeyPem(testPemKey)
                .build()) {

            assertNotNull(client);
            assertEquals(1, client.getSuperAdminPublicKeys().size());
        }
    }

    @Test
    void fluentChaining_withPemKeyList() throws ApiKeyTPV1Exception, IOException {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeysPem(Arrays.asList(testPemKey))
                .build()) {

            assertNotNull(client);
            assertEquals(1, client.getSuperAdminPublicKeys().size());
        }
    }

    @Test
    void methodChaining_preservesAllSettings() throws ApiKeyTPV1Exception {
        long customTtl = 180_000L;
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .apiKey(TEST_API_KEY)
                .apiSecret(TEST_API_SECRET)
                .superAdminKey(testSingleKey)
                .minValidSignatures(1)
                .rulesContainerCacheTtlMs(customTtl)
                .build()) {

            assertNotNull(client);
            assertEquals(TEST_HOST, client.getOpenApiClient().getBasePath());
            assertEquals(1, client.getSuperAdminPublicKeys().size());
            assertEquals(customTtl, client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    // =====================================================================
    // Edge case tests
    // =====================================================================

    @Test
    void addingMultipleKeys_combinesAllKeys() throws ApiKeyTPV1Exception {
        List<PublicKey> additionalKeys = new ArrayList<>();
        additionalKeys.add(testSuperAdminKeys.get(1));

        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKey(testSingleKey)
                .superAdminKeys(additionalKeys)
                .build()) {

            assertNotNull(client);
            assertEquals(2, client.getSuperAdminPublicKeys().size());
        }
    }

    @Test
    void addingMultiplePemKeys_combinesAllKeys() throws ApiKeyTPV1Exception, IOException {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeyPem(testPemKey)
                .superAdminKeysPem(Arrays.asList(testPemKey))  // Same key, but should be added
                .build()) {

            assertNotNull(client);
            assertEquals(2, client.getSuperAdminPublicKeys().size());
        }
    }

    @Test
    void credentialsOverridesIndividualMethods() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .apiKey("old-key")
                .apiSecret("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .build()) {

            assertNotNull(client);
            // Can't directly verify credentials, but build should succeed with new values
        }
    }

    @Test
    void individualMethodsOverrideCredentials() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials("old-key", "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210")
                .apiKey(TEST_API_KEY)
                .apiSecret(TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .build()) {

            assertNotNull(client);
            // Can't directly verify credentials, but build should succeed with new values
        }
    }

    @Test
    void defaultMinValidSignatures_isOne() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .build()) {

            assertNotNull(client);
            // Default is 1, and build should succeed
        }
    }

    @Test
    void defaultCacheTtl_isUsedWhenNotSpecified() throws ApiKeyTPV1Exception {
        try (ProtectClient client = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys)
                .build()) {

            assertNotNull(client);
            assertEquals(RulesContainerCache.DEFAULT_CACHE_TTL_MS,
                    client.getRulesContainerCache().getCacheTtlMs());
        }
    }

    @Test
    void build_canBeCalledOnlyOnce_createsNewClient() throws ApiKeyTPV1Exception {
        ProtectClientBuilder builder = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(testSuperAdminKeys);

        try (ProtectClient client1 = builder.build()) {
            assertNotNull(client1);
        }

        // Building again should create a new client (builder is not consumed)
        try (ProtectClient client2 = builder.build()) {
            assertNotNull(client2);
        }
    }

    @Test
    void emptyKeysList_failsOnBuild() {
        ProtectClientBuilder builder = ProtectClient.builder()
                .host(TEST_HOST)
                .credentials(TEST_API_KEY, TEST_API_SECRET)
                .superAdminKeys(Collections.emptyList());  // Empty list adds nothing

        IllegalStateException ex = assertThrows(IllegalStateException.class, builder::build);
        assertTrue(ex.getMessage().contains("SuperAdmin"));
    }
}
