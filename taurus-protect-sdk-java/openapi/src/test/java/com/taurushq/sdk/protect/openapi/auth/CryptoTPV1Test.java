package com.taurushq.sdk.protect.openapi.auth;

import com.google.gson.Gson;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Base64;
import org.apache.commons.codec.binary.Hex;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for CryptoTPV1 cryptographic operations.
 *
 * <p><strong>SECURITY NOTICE:</strong> All cryptographic keys in this test class are TEST-ONLY
 * values generated specifically for unit testing purposes. These keys have NO production value
 * and are NOT connected to any real accounts or systems. They are safe to include in
 * a public repository.</p>
 */
public class CryptoTPV1Test {

    // ==================================================================================
    // TEST KEYS - FOR UNIT TESTING ONLY
    // These are synthetic test keys with no real-world security implications.
    // They are NOT valid for any production system and pose no security risk.
    // ==================================================================================

    PrivateKey user1PrivKey() throws Exception {

        // Test private key - generated for testing, has no production value
        String secretPEMKey = "-----BEGIN EC PRIVATE KEY-----\n" +
                "MHcCAQEEIOd7BwfDcXGDo0cTF9KczH9/jq27xIUEFk6v7iCeY5n3oAoGCCqGSM49\n" +
                "AwEHoUQDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMreg6fLAVtHkwKZ6H4L7J9WhRC\n" +
                "VtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==\n" +
                "-----END EC PRIVATE KEY-----";

        PrivateKey privateKey = CryptoTPV1.decodePrivateKey(secretPEMKey);
        return privateKey;
    }

    PublicKey user1PubKey() throws Exception {
        // Test public key - generated for testing, has no production value
        String pubPEMKey = "-----BEGIN PUBLIC KEY-----\n" +
                "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMr\n" +
                "eg6fLAVtHkwKZ6H4L7J9WhRCVtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==\n" +
                "-----END PUBLIC KEY-----";

        PublicKey publicKey = CryptoTPV1.decodePublicKey(pubPEMKey);
        return publicKey;

    }

    PublicKey superAdmin1PubKey() throws Exception {
        // Test public key - generated for testing, has no production value
        String pubPEMKey = "-----BEGIN PUBLIC KEY-----\n"
                + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW\n"
                + "L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==\n"
                + "-----END PUBLIC KEY-----";

        return CryptoTPV1.decodePublicKey(pubPEMKey);
    }

    PublicKey superAdmin2PubKey() throws Exception {
        String pubPEMKey = "-----BEGIN PUBLIC KEY-----\n"
                + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEY9zGugzNLIfpZuaUrzywEh/8ZdtX\n"
                + "4IIuIpDHLvJ36glFjfxxSZdOG6yHKFFlQh1GX3OCFZxHe+xeOGBJHBgraA==\n"
                + "-----END PUBLIC KEY-----";

        return CryptoTPV1.decodePublicKey(pubPEMKey);
    }

    PublicKey superAdmin3PubKey() throws Exception {
        String pubPEMKey = "-----BEGIN PUBLIC KEY-----\n"
                + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEV2zEDvZjJ/qR+Qh1r0sKNci39w79\n"
                + "B00Ng5MGg4wRM7Ah9pr1GbkgrKIUf4kaHPjrPVaQKIR8zDU3jgm0lWxeHA==\n"
                + "-----END PUBLIC KEY-----";

        return CryptoTPV1.decodePublicKey(pubPEMKey);
    }

    @Test
    void checkHashesSignature() throws Exception {
        String h = "eb57ef103a0b4c7b7bd5b394963eecf03dd5cad5d92ad03f8a05ed6a98540518";
        List<String> hashes = Collections.singletonList(h);
        Gson gson = new Gson();
        String j = gson.toJson(hashes);


        assertEquals("[\"eb57ef103a0b4c7b7bd5b394963eecf03dd5cad5d92ad03f8a05ed6a98540518\"]", j);

        String sig = CryptoTPV1.calculateBase64Signature(user1PrivKey(), j.getBytes(StandardCharsets.UTF_8));

        assertTrue(CryptoTPV1.verifyBase64Signature(user1PubKey(), j.getBytes(StandardCharsets.UTF_8), sig));

        String goBase64Sig = "Yoz0xd0N0OIOcOkBQxUtaYIxYEoyyteH1o07v+c/ycCyNSo/z8tcvXWQ1Uvc8E9dWd970bMQjzGCCTH1mF2L6g==";

        assertTrue(CryptoTPV1.verifyBase64Signature(user1PubKey(), j.getBytes(StandardCharsets.UTF_8), goBase64Sig));

    }

    @Test
    void calculateHexHash() {
        assertEquals("31f7a65e315586ac198bd798b6629ce4903d0899476d5741a9f32e2e521b6a66", CryptoTPV1.calculateHexHash("toto"));
        assertEquals("b2167b0aa7ef7794740b055ac7a880a52934aa67ef1ca6887ad81dccefd5b9de", CryptoTPV1.calculateHexHash("my data"));
    }

    @Test
    void calculateBase64Hmac() throws DecoderException {
        String sig = CryptoTPV1.calculateBase64Hmac("my secret".getBytes(), "my data");
        assertEquals(Base64.encodeBase64String(Hex.decodeHex("049f32be4f33a698205529818c1b676d460c9cb2f6901457d012a6646127ae37")), sig);
        assertEquals("BJ8yvk8zppggVSmBjBtnbUYMnLL2kBRX0BKmZGEnrjc=", sig);
    }

    @Test
    void checkBase64Hmac() {
        assertTrue(CryptoTPV1.checkBase64Hmac("my secret".getBytes(), "my data", "BJ8yvk8zppggVSmBjBtnbUYMnLL2kBRX0BKmZGEnrjc="));
        assertFalse(CryptoTPV1.checkBase64Hmac("my secret".getBytes(), "my dota", "BJ8yvk8zppggVSmBjBtnbUYMnLL2kBRX0BKmZGEnrjc="));
    }

    @Test
    void calculateHexHmac() {
        String sig = CryptoTPV1.calculateHexHmac("my secret", "my data");
        assertEquals("049f32be4f33a698205529818c1b676d460c9cb2f6901457d012a6646127ae37", sig);
    }

    @Test
    void checkHexHmac() {
        assertTrue(CryptoTPV1.checkHexHmac("my secret", "my data", "049f32be4f33a698205529818c1b676d460c9cb2f6901457d012a6646127ae37"));
        assertFalse(CryptoTPV1.checkHexHmac("my secret", "my dota", "049f32be4f33a698205529818c1b676d460c9cb2f6901457d012a6646127ae37"));
    }


    @Test
    void calculateSignedHeader() {
        String h1 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "host", "path", "query", "content-type", "body");
        String sig1 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method host path query content-type body");
        String eh1 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig1;
        assertEquals(eh1, h1);

        String h2 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "host", "path", "query", "content-type", "");
        String sig2 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method host path query content-type");
        String eh2 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig2;
        assertEquals(eh2, h2);

        String h3 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "host", "path", "query", "", "");
        String sig3 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method host path query");
        String eh3 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig3;
        assertEquals(eh3, h3);

        String h4 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "host", "path", "", "", "");
        String sig4 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method host path");
        String eh4 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig4;
        assertEquals(eh4, h4);

        String h5 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "host", "", "", "", "");
        String sig5 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method host");
        String eh5 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig5;
        assertEquals(eh5, h5);

        String h6 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "method", "", "", "", "", "");
        String sig6 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10 method");
        String eh6 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig6;
        assertEquals(eh6, h6);

        String h7 = CryptoTPV1.calculateSignedHeader("api-key", "api-secret".getBytes(), "nonce", 10, "", "", "", "", "", "");
        String sig7 = CryptoTPV1.calculateBase64Hmac("api-secret".getBytes(), "TPV1 api-key nonce 10");
        String eh7 = "TPV1-HMAC-SHA256 ApiKey=api-key Nonce=nonce Timestamp=10 Signature=" + sig7;
        assertEquals(eh7, h7);

    }

    @Test
    void decodeSuperAdminPublicKeys() throws Exception {
        PublicKey key1 = superAdmin1PubKey();
        PublicKey key2 = superAdmin2PubKey();
        PublicKey key3 = superAdmin3PubKey();

        assertNotNull(key1);
        assertNotNull(key2);
        assertNotNull(key3);
        assertEquals("ECDSA", key1.getAlgorithm());
        assertEquals("ECDSA", key2.getAlgorithm());
        assertEquals("ECDSA", key3.getAlgorithm());
    }
}
