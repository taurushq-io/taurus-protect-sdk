package com.taurushq.sdk.protect.client;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Auth;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Exception;

import java.util.function.Supplier;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * The authentication mechanism for a {@link ProtectClient}: either static
 * TPV1-HMAC (api key + secret) or a Bearer token (static, or resolved per request
 * for rotating/per-caller tokens).
 * <p>
 * Build one with {@link #apiKey}, {@link #bearerToken}, or
 * {@link #bearerTokenProvider} and pass it to
 * {@link ProtectClient#create(String, Credentials, java.util.List, int)} or
 * {@link ProtectClientBuilder#credentials(Credentials)}.
 */
public abstract class Credentials {

    private Credentials() {
    }

    /**
     * TPV1-HMAC credentials (a static shared-service identity). SuperAdmin keys are
     * required, since services verify governance rules signatures.
     *
     * @param apiKey    the API key
     * @param apiSecret the API secret (hex-encoded)
     * @return the credentials
     */
    public static Credentials apiKey(String apiKey, String apiSecret) {
        checkArgument(!Strings.isNullOrEmpty(apiKey), "apiKey cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(apiSecret), "apiSecret cannot be null or empty");
        return new ApiKeyCredentials(apiKey, apiSecret);
    }

    /**
     * A single static Bearer token (the "Authorization: Bearer" header, no TPV1
     * signing). The token carries the identity to the API.
     *
     * @param token the bearer token
     * @return the credentials
     */
    public static Credentials bearerToken(String token) {
        checkArgument(!Strings.isNullOrEmpty(token), "bearerToken cannot be null or empty");
        return new BearerCredentials(() -> token);
    }

    /**
     * A Bearer token resolved per request from a provider (for rotating/per-caller
     * tokens).
     *
     * @param provider supplies the bearer token for each request
     * @return the credentials
     */
    public static Credentials bearerTokenProvider(Supplier<String> provider) {
        checkNotNull(provider, "bearerTokenProvider cannot be null");
        return new BearerCredentials(provider);
    }

    abstract void applyTo(ApiClient openApiClient) throws ApiKeyTPV1Exception;

    private static final class ApiKeyCredentials extends Credentials {
        private final String apiKey;
        private final String apiSecret;

        ApiKeyCredentials(String apiKey, String apiSecret) {
            this.apiKey = apiKey;
            this.apiSecret = apiSecret;
        }

        @Override
        void applyTo(ApiClient openApiClient) throws ApiKeyTPV1Exception {
            openApiClient.setApiKeyTPV1(apiKey);
            openApiClient.setApiSecretTPV1(apiSecret);
        }
    }

    private static final class BearerCredentials extends Credentials {
        private final Supplier<String> provider;

        BearerCredentials(Supplier<String> provider) {
            this.provider = provider;
        }

        @Override
        void applyTo(ApiClient openApiClient) {
            ((ApiKeyTPV1Auth) openApiClient.getAuthentication("ApiKeyTPV1")).setBearerTokenProvider(provider);
        }
    }
}
