package com.taurushq.sdk.protect.client;

import org.junit.jupiter.api.Test;

import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class ProtectClientBearerTest {

    private static PublicKey p256Key() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        return kpg.generateKeyPair().getPublic();
    }

    @Test
    void bearerTokenCredentialsBuildClient() throws Exception {
        ProtectClient client = ProtectClient.create(
                "https://api.example.com", Credentials.bearerToken("session-token"),
                Collections.singletonList(p256Key()), 1);
        assertNotNull(client);
    }

    @Test
    void bearerTokenProviderCredentialsBuildClient() throws Exception {
        ProtectClient client = ProtectClient.create(
                "https://api.example.com", Credentials.bearerTokenProvider(() -> "per-request-token"),
                Collections.singletonList(p256Key()), 1);
        assertNotNull(client);
    }

    @Test
    void builderCredentialsBearerBuildsClient() throws Exception {
        ProtectClient client = ProtectClient.builder()
                .host("https://api.example.com")
                .credentials(Credentials.bearerTokenProvider(() -> "per-request-token"))
                .superAdminKey(p256Key())
                .build();
        assertNotNull(client);
    }

    @Test
    void bearerCredentialsStillRequireSuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(
                        "https://api.example.com", Credentials.bearerTokenProvider(() -> "tok"),
                        Collections.emptyList(), 1));
    }

    @Test
    void apiKeyCredentialsRequireSuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                ProtectClient.create(
                        "https://api.example.com", Credentials.apiKey("key", "abcdef"),
                        Collections.emptyList(), 1));
    }

    @Test
    void bearerTokenRejectsEmptyToken() {
        assertThrows(IllegalArgumentException.class, () -> Credentials.bearerToken(""));
    }

    @Test
    void bearerTokenProviderRejectsNull() {
        assertThrows(NullPointerException.class, () -> Credentials.bearerTokenProvider(null));
    }

    @Test
    void apiKeyRejectsEmptyValues() {
        assertThrows(IllegalArgumentException.class, () -> Credentials.apiKey("", "abcdef"));
        assertThrows(IllegalArgumentException.class, () -> Credentials.apiKey("key", ""));
    }

    @Test
    void builderRequiresCredentials() throws Exception {
        PublicKey key = p256Key();
        assertThrows(IllegalStateException.class, () ->
                ProtectClient.builder()
                        .host("https://api.example.com")
                        .superAdminKey(key)
                        .build());
    }
}
