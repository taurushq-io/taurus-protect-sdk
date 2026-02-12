package com.taurushq.sdk.protect.client;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Exception;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;

import java.io.IOException;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;
import static com.google.common.base.Strings.isNullOrEmpty;

/**
 * Builder for creating {@link ProtectClient} instances with a fluent API.
 * <p>
 * This builder provides a clearer and more discoverable way to configure
 * the ProtectClient compared to the static factory methods.
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
 * @see ProtectClient
 */
public final class ProtectClientBuilder {

    private String host;
    private String apiKey;
    private String apiSecret;
    private final List<PublicKey> superAdminPublicKeys = new ArrayList<>();
    private int minValidSignatures = 1;
    private long rulesContainerCacheTtlMs = RulesContainerCache.DEFAULT_CACHE_TTL_MS;

    /**
     * Creates a new builder instance.
     */
    ProtectClientBuilder() {
        // Package-private constructor for use by ProtectClient.builder()
    }

    /**
     * Sets the host URL for the Taurus-PROTECT API.
     *
     * @param host the API host URL (e.g., "https://api.protect.taurushq.com")
     * @return this builder for chaining
     */
    public ProtectClientBuilder host(String host) {
        this.host = host;
        return this;
    }

    /**
     * Sets the API credentials.
     *
     * @param apiKey    the API key
     * @param apiSecret the API secret (hex-encoded)
     * @return this builder for chaining
     */
    public ProtectClientBuilder credentials(String apiKey, String apiSecret) {
        this.apiKey = apiKey;
        this.apiSecret = apiSecret;
        return this;
    }

    /**
     * Sets the API key.
     *
     * @param apiKey the API key
     * @return this builder for chaining
     */
    public ProtectClientBuilder apiKey(String apiKey) {
        this.apiKey = apiKey;
        return this;
    }

    /**
     * Sets the API secret.
     *
     * @param apiSecret the API secret (hex-encoded)
     * @return this builder for chaining
     */
    public ProtectClientBuilder apiSecret(String apiSecret) {
        this.apiSecret = apiSecret;
        return this;
    }

    /**
     * Adds SuperAdmin public keys for governance rule verification.
     *
     * @param publicKeys the public keys to add
     * @return this builder for chaining
     */
    public ProtectClientBuilder superAdminKeys(List<PublicKey> publicKeys) {
        checkNotNull(publicKeys, "publicKeys cannot be null");
        this.superAdminPublicKeys.addAll(publicKeys);
        return this;
    }

    /**
     * Adds a single SuperAdmin public key for governance rule verification.
     *
     * @param publicKey the public key to add
     * @return this builder for chaining
     */
    public ProtectClientBuilder superAdminKey(PublicKey publicKey) {
        checkNotNull(publicKey, "publicKey cannot be null");
        this.superAdminPublicKeys.add(publicKey);
        return this;
    }

    /**
     * Adds SuperAdmin public keys from PEM-formatted strings.
     *
     * @param pemKeys the PEM-formatted public keys
     * @return this builder for chaining
     * @throws IOException if a PEM key cannot be decoded
     */
    public ProtectClientBuilder superAdminKeysPem(List<String> pemKeys) throws IOException {
        checkNotNull(pemKeys, "pemKeys cannot be null");
        for (String pem : pemKeys) {
            this.superAdminPublicKeys.add(CryptoTPV1.decodePublicKey(pem));
        }
        return this;
    }

    /**
     * Adds a single SuperAdmin public key from a PEM-formatted string.
     *
     * @param pemKey the PEM-formatted public key
     * @return this builder for chaining
     * @throws IOException if the PEM key cannot be decoded
     */
    public ProtectClientBuilder superAdminKeyPem(String pemKey) throws IOException {
        checkNotNull(pemKey, "pemKey cannot be null");
        this.superAdminPublicKeys.add(CryptoTPV1.decodePublicKey(pemKey));
        return this;
    }

    /**
     * Sets the minimum number of valid signatures required for governance rule verification.
     * <p>
     * This determines the threshold for multi-signature verification of governance rules.
     *
     * @param minValidSignatures the minimum number of signatures (must be greater than 0)
     * @return this builder for chaining
     */
    public ProtectClientBuilder minValidSignatures(int minValidSignatures) {
        checkArgument(minValidSignatures > 0, "minValidSignatures must be greater than 0");
        this.minValidSignatures = minValidSignatures;
        return this;
    }

    /**
     * Sets the cache TTL for the rules container.
     * <p>
     * The rules container cache stores governance rules to avoid fetching them
     * on every request. Set a lower value if rules change frequently.
     *
     * @param cacheTtlMs the cache TTL in milliseconds (must be positive)
     * @return this builder for chaining
     */
    public ProtectClientBuilder rulesContainerCacheTtlMs(long cacheTtlMs) {
        checkArgument(cacheTtlMs > 0, "cacheTtlMs must be positive");
        this.rulesContainerCacheTtlMs = cacheTtlMs;
        return this;
    }

    /**
     * Builds the ProtectClient with the configured settings.
     *
     * @return a new ProtectClient instance
     * @throws ApiKeyTPV1Exception  if the API secret is invalid
     * @throws IllegalStateException if required configuration is missing
     */
    public ProtectClient build() throws ApiKeyTPV1Exception {
        validate();
        return ProtectClient.create(host, apiKey, apiSecret, superAdminPublicKeys,
                minValidSignatures, rulesContainerCacheTtlMs);
    }

    private void validate() {
        if (isNullOrEmpty(host)) {
            throw new IllegalStateException("host is required. Call host(\"https://...\")");
        }
        if (isNullOrEmpty(apiKey)) {
            throw new IllegalStateException("apiKey is required. Call credentials(key, secret) or apiKey(key)");
        }
        if (isNullOrEmpty(apiSecret)) {
            throw new IllegalStateException("apiSecret is required. Call credentials(key, secret) or apiSecret(secret)");
        }
        if (superAdminPublicKeys.isEmpty()) {
            throw new IllegalStateException(
                    "At least one SuperAdmin public key is required. Call superAdminKeys(...) or superAdminKeysPem(...)");
        }
    }
}
