package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Price signature verification against the PRICEUPDATER role.
 * <p>
 * Rate and decimals feed amount conversion, so an unverified price is a wrong number a
 * caller acts on.
 */
class PriceVerifierTest {

    @BeforeAll
    static void setUpProvider() {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private static KeyPair generateKeyPair() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        return generator.generateKeyPair();
    }

    private static Price price() {
        Price p = new Price();
        p.setBlockchain("ETH");
        p.setCurrencyFrom("ETH");
        p.setCurrencyTo("USD");
        p.setDecimals("18");
        p.setRate("2500.00");
        return p;
    }

    private static DecodedRulesContainer container(String role, PublicKey key) {
        RuleUser user = new RuleUser();
        user.setId("price@bank.com");
        user.setPublicKey(key);
        user.setRoles(Collections.singletonList(role));

        DecodedRulesContainer container = new DecodedRulesContainer();
        container.setUsers(Collections.singletonList(user));
        return container;
    }

    private static Price signed(Price p, KeyPair keyPair) throws Exception {
        PriceSignature sig = new PriceSignature();
        sig.setUserId("price@bank.com");
        sig.setSignature(CryptoTPV1.calculateBase64Signature(
                keyPair.getPrivate(), PriceVerifier.priceSignedBytes(p)));
        p.setSignatures(Collections.singletonList(sig));
        return p;
    }

    // The canonical form is validatord's CurrencyPrice JSON projection. Pinned by value:
    // a reordered or extended form silently changes every signature this SDK accepts.
    @Test
    @DisplayName("signed bytes are the canonical projection")
    void signedBytesAreCanonical() {
        assertEquals(
                "{\"blockchain\":\"ETH\",\"currencyFrom\":\"ETH\",\"currencyTo\":\"USD\","
                        + "\"decimals\":\"18\",\"rate\":\"2500.00\"}",
                new String(PriceVerifier.priceSignedBytes(price()), StandardCharsets.UTF_8));
    }

    @Test
    @DisplayName("a PRICEUPDATER signature verifies")
    void acceptsPriceUpdaterSignature() throws Exception {
        KeyPair keyPair = generateKeyPair();
        Price p = signed(price(), keyPair);
        assertDoesNotThrow(() ->
                PriceVerifier.verifyPrice(p, container("PRICEUPDATER", keyPair.getPublic())));
    }

    @Test
    @DisplayName("a signature from a key without the role does not verify")
    void rejectsNonPriceUpdater() throws Exception {
        KeyPair signer = generateKeyPair();
        KeyPair updater = generateKeyPair();
        Price p = signed(price(), signer);
        assertThrows(IntegrityException.class, () ->
                PriceVerifier.verifyPrice(p, container("PRICEUPDATER", updater.getPublic())));
    }

    @Test
    @DisplayName("a rate altered after signing does not verify")
    void rejectsTamperedRate() throws Exception {
        KeyPair keyPair = generateKeyPair();
        Price p = signed(price(), keyPair);
        // The rate is what a caller converts with.
        p.setRate("1.00");
        assertThrows(IntegrityException.class, () ->
                PriceVerifier.verifyPrice(p, container("PRICEUPDATER", keyPair.getPublic())));
    }

    @Test
    @DisplayName("stripped signatures are rejected when a PRICEUPDATER exists")
    void rejectsStrippedSignatures() throws Exception {
        KeyPair keyPair = generateKeyPair();
        assertThrows(IntegrityException.class, () ->
                PriceVerifier.verifyPrice(price(),
                        container("PRICEUPDATER", keyPair.getPublic())));
    }

    @Test
    @DisplayName("a tenant with no PRICEUPDATER is not forced to sign")
    void passesThroughWithoutPriceUpdater() throws Exception {
        KeyPair keyPair = generateKeyPair();
        assertDoesNotThrow(() ->
                PriceVerifier.verifyPrice(price(),
                        container("REQUESTAPPROVER", keyPair.getPublic())));
    }

    @Test
    @DisplayName("a rules container is required")
    void requiresRulesContainer() {
        assertThrows(IntegrityException.class, () -> PriceVerifier.verifyPrice(price(), null));
    }

    @Test
    @DisplayName("one unverifiable price fails the batch")
    void batchReportsFirstFailure() throws Exception {
        KeyPair keyPair = generateKeyPair();
        Price good = signed(price(), keyPair);
        Price bad = price();
        bad.setCurrencyTo("EUR");

        assertThrows(IntegrityException.class, () ->
                PriceVerifier.verifyPrices(Arrays.asList(good, bad),
                        container("PRICEUPDATER", keyPair.getPublic())));
    }
}
