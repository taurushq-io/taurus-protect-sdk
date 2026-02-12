package com.taurushq.sdk.protect.client.testutil;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Exception;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.Assumptions;

import java.io.IOException;
import java.security.PrivateKey;

/**
 * Shared helper methods for tests.
 */
public final class TestHelper {

    private TestHelper() {
        // Utility class
    }

    /**
     * Skips the current test if tests are not enabled.
     * Call this at the beginning of each test class setup.
     */
    public static void skipIfNotEnabled() {
        Assumptions.assumeTrue(
                TestConfig.isEnabled(),
                "Skipping test. Set PROTECT_INTEGRATION_TEST=true or configure test.properties."
        );
    }

    /**
     * Creates a ProtectClient for the identity at the given 1-based index.
     *
     * @param identityIndex 1-based identity index
     * @return configured ProtectClient
     * @throws ApiKeyTPV1Exception if API key configuration fails
     * @throws IOException         if client creation fails
     */
    public static ProtectClient getTestClient(int identityIndex) throws ApiKeyTPV1Exception, IOException {
        TestConfig.Identity identity = TestConfig.getIdentity(identityIndex);
        return ProtectClient.createFromPem(
                TestConfig.getHost(),
                identity.getApiKey(),
                identity.getApiSecret(),
                TestConfig.getSuperAdminKeys(),
                TestConfig.getMinValidSignatures()
        );
    }

    /**
     * Returns the decoded PrivateKey for the identity at the given 1-based index.
     * Returns null if the identity has no private key configured.
     *
     * @param identityIndex 1-based identity index
     * @return the decoded PrivateKey, or null if no private key is configured
     * @throws IOException if key decoding fails
     */
    public static PrivateKey getPrivateKey(int identityIndex) throws IOException {
        TestConfig.Identity identity = TestConfig.getIdentity(identityIndex);
        if (!identity.hasPrivateKey()) {
            return null;
        }
        return CryptoTPV1.decodePrivateKey(identity.getPrivateKey());
    }

    /**
     * Skips the current test if fewer than requiredCount identities are configured.
     *
     * @param requiredCount minimum number of identities needed
     */
    public static void skipIfInsufficientIdentities(int requiredCount) {
        Assumptions.assumeTrue(
                TestConfig.getIdentityCount() >= requiredCount,
                "Skipping: need " + requiredCount + " identities but only "
                        + TestConfig.getIdentityCount() + " configured."
        );
    }
}
