package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.InternalAddress;
import com.taurushq.sdk.protect.client.model.InternalWallet;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.StringReader;
import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for WhitelistHashHelper.
 *
 * <p><strong>SECURITY NOTICE:</strong> The cryptographic keys below are TEST-ONLY values
 * generated specifically for unit testing purposes. These keys have NO production value
 * and are NOT connected to any real accounts or systems. They are safe to include in
 * a public repository.</p>
 */
class WhitelistHashHelperTest {

    // ==================================================================================
    // TEST KEYS - FOR UNIT TESTING ONLY
    // These are synthetic test keys with no real-world security implications.
    // They are NOT valid for any production system and pose no security risk.
    // ==================================================================================

    private static final String USER1_PRIVATE_KEY_PEM =
            "-----BEGIN EC PRIVATE KEY-----\n"
                    + "MHcCAQEEIOd7BwfDcXGDo0cTF9KczH9/jq27xIUEFk6v7iCeY5n3oAoGCCqGSM49\n"
                    + "AwEHoUQDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMreg6fLAVtHkwKZ6H4L7J9WhRC\n"
                    + "VtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==\n"
                    + "-----END EC PRIVATE KEY-----";

    private static final String USER1_PUBLIC_KEY_PEM =
            "-----BEGIN PUBLIC KEY-----\n"
                    + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMr\n"
                    + "eg6fLAVtHkwKZ6H4L7J9WhRCVtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==\n"
                    + "-----END PUBLIC KEY-----";

    private static PrivateKey user1PrivateKey;
    private static PublicKey user1PublicKey;

    @BeforeAll
    static void setup() throws Exception {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }

        user1PrivateKey = loadPrivateKey(USER1_PRIVATE_KEY_PEM);
        user1PublicKey = loadPublicKey(USER1_PUBLIC_KEY_PEM);
    }

    private static PrivateKey loadPrivateKey(String pem) throws Exception {
        org.bouncycastle.openssl.PEMParser parser =
                new org.bouncycastle.openssl.PEMParser(new StringReader(pem));
        Object object = parser.readObject();
        parser.close();

        if (object instanceof org.bouncycastle.openssl.PEMKeyPair) {
            org.bouncycastle.openssl.PEMKeyPair keyPair = (org.bouncycastle.openssl.PEMKeyPair) object;
            org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter converter =
                    new org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter().setProvider("BC");
            KeyPair kp = converter.getKeyPair(keyPair);
            return kp.getPrivate();
        }
        throw new IllegalArgumentException("Not a PEM key pair");
    }

    private static PublicKey loadPublicKey(String pem) throws Exception {
        org.bouncycastle.openssl.PEMParser parser =
                new org.bouncycastle.openssl.PEMParser(new StringReader(pem));
        Object object = parser.readObject();
        parser.close();

        if (object instanceof org.bouncycastle.asn1.x509.SubjectPublicKeyInfo) {
            org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter converter =
                    new org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter().setProvider("BC");
            return converter.getPublicKey((org.bouncycastle.asn1.x509.SubjectPublicKeyInfo) object);
        }
        throw new IllegalArgumentException("Not a PEM public key");
    }

    // ============ Tests for parseWhitelistedAddressFromJson ============

    @Test
    void testParseWhitelistedAddressFromJson_BasicFields() throws WhitelistException {
        String json = "{\"currency\":\"ETH\",\"network\":\"mainnet\",\"addressType\":\"individual\","
                + "\"address\":\"0x1234\",\"memo\":\"test memo\",\"label\":\"test label\","
                + "\"customerId\":\"cust123\",\"exchangeAccountId\":\"\","
                + "\"linkedInternalAddresses\":[],\"linkedWallets\":[],"
                + "\"contractType\":\"\",\"tnParticipantID\":\"\"}";

        WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(json);

        assertEquals("ETH", addr.getBlockchain());
        assertEquals("mainnet", addr.getNetwork());
        assertEquals(WhitelistedAddress.AddressType.individual, addr.getAddressType());
        assertEquals("0x1234", addr.getAddress());
        assertEquals("test memo", addr.getMemo());
        assertEquals("test label", addr.getLabel());
        assertEquals("cust123", addr.getCustomerId());
    }

    @Test
    void testParseWhitelistedAddressFromJson_WithExchangeAccountId() throws WhitelistException {
        String json = "{\"currency\":\"BTC\",\"network\":\"mainnet\",\"addressType\":\"individual\","
                + "\"address\":\"bc1abc\",\"memo\":\"\",\"label\":\"label\","
                + "\"customerId\":\"\",\"exchangeAccountId\":\"42\","
                + "\"linkedInternalAddresses\":[],\"linkedWallets\":[],"
                + "\"contractType\":\"\",\"tnParticipantID\":\"\"}";

        WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(json);

        assertEquals("BTC", addr.getBlockchain());
        assertEquals(42L, addr.getExchangeAccountId());
    }

    @Test
    void testParseWhitelistedAddressFromJson_WithLinkedInternalAddresses() throws WhitelistException {
        String json = "{\"currency\":\"ETH\",\"network\":\"mainnet\",\"addressType\":\"individual\","
                + "\"address\":\"0x1234\",\"memo\":\"\",\"label\":\"\","
                + "\"customerId\":\"\",\"exchangeAccountId\":\"\","
                + "\"linkedInternalAddresses\":[{\"id\":1,\"address\":\"addr1\",\"label\":\"label1\"},"
                + "{\"id\":2,\"address\":\"addr2\",\"label\":\"label2\"}],"
                + "\"linkedWallets\":[],\"contractType\":\"\",\"tnParticipantID\":\"\"}";

        WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(json);

        assertNotNull(addr.getLinkedInternalAddresses());
        assertEquals(2, addr.getLinkedInternalAddresses().size());

        InternalAddress ia1 = addr.getLinkedInternalAddresses().get(0);
        assertEquals("1", ia1.getId());
        assertEquals("addr1", ia1.getAddress());
        assertEquals("label1", ia1.getLabel());

        InternalAddress ia2 = addr.getLinkedInternalAddresses().get(1);
        assertEquals("2", ia2.getId());
        assertEquals("addr2", ia2.getAddress());
        assertEquals("label2", ia2.getLabel());
    }

    @Test
    void testParseWhitelistedAddressFromJson_WithLinkedWallets() throws WhitelistException {
        String json = "{\"currency\":\"ETH\",\"network\":\"mainnet\",\"addressType\":\"individual\","
                + "\"address\":\"0x1234\",\"memo\":\"\",\"label\":\"\","
                + "\"customerId\":\"\",\"exchangeAccountId\":\"\","
                + "\"linkedInternalAddresses\":[],"
                + "\"linkedWallets\":[{\"id\":10,\"name\":\"Wallet1\",\"path\":\"m/44'/60'/0'/0\"}],"
                + "\"contractType\":\"\",\"tnParticipantID\":\"\"}";

        WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(json);

        assertNotNull(addr.getLinkedWallets());
        assertEquals(1, addr.getLinkedWallets().size());

        InternalWallet w1 = addr.getLinkedWallets().get(0);
        assertEquals(10L, w1.getId());
        assertEquals("Wallet1", w1.getName());
        assertEquals("m/44'/60'/0'/0", w1.getPath());
    }

    @Test
    void testParseWhitelistedAddressFromJson_NullOrEmptyJson() {
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.parseWhitelistedAddressFromJson(null));
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.parseWhitelistedAddressFromJson(""));
    }

    @Test
    void testParseWhitelistedAddressFromJson_InvalidJson() {
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.parseWhitelistedAddressFromJson("not valid json"));
    }

    @Test
    void testParseWhitelistedAddressFromJson_EmptyStringFieldsTreatedAsNull() throws WhitelistException {
        String json = "{\"currency\":\"ETH\",\"network\":\"\",\"addressType\":\"individual\","
                + "\"address\":\"0x1234\",\"memo\":\"\",\"label\":\"\","
                + "\"customerId\":\"\",\"exchangeAccountId\":\"\","
                + "\"linkedInternalAddresses\":[],\"linkedWallets\":[],"
                + "\"contractType\":\"\",\"tnParticipantID\":\"\"}";

        WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(json);

        assertNull(addr.getNetwork());
        assertNull(addr.getMemo());
        assertNull(addr.getLabel());
        assertNull(addr.getCustomerId());
        assertNull(addr.getContractType());
        assertNull(addr.getTnParticipantID());
    }

    // ============ Tests for signHashes and checkHashesSignature ============

    @Test
    void testSignAndVerifyHashes() throws Exception {
        List<String> hashes = Collections.singletonList(
                "abc123def456abc123def456abc123def456abc123def456abc123def456abcd");

        // Sign with private key
        byte[] signature = WhitelistHashHelper.signHashes(hashes, user1PrivateKey);
        assertNotNull(signature);
        assertTrue(signature.length >= 64 && signature.length <= 72);

        // Verify with public key - should not throw
        WhitelistHashHelper.checkHashesSignature(hashes, signature, user1PublicKey);
    }

    @Test
    void testSignAndVerifyMultipleHashes() throws Exception {
        List<String> hashes = Arrays.asList(
                "hash1abc123def456abc123def456abc123def456abc123def456abc123def456",
                "hash2abc123def456abc123def456abc123def456abc123def456abc123def456",
                "hash3abc123def456abc123def456abc123def456abc123def456abc123def456"
        );

        byte[] signature = WhitelistHashHelper.signHashes(hashes, user1PrivateKey);
        assertNotNull(signature);

        // Verify should pass
        WhitelistHashHelper.checkHashesSignature(hashes, signature, user1PublicKey);
    }

    @Test
    void testVerifyWithWrongHashes() throws Exception {
        List<String> originalHashes = Collections.singletonList("originalhash123");
        List<String> wrongHashes = Collections.singletonList("wronghash456");

        byte[] signature = WhitelistHashHelper.signHashes(originalHashes, user1PrivateKey);

        // Verification should fail with wrong hashes
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.checkHashesSignature(wrongHashes, signature, user1PublicKey));
    }

    @Test
    void testSignHashesNullParameters() {
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.signHashes(null, user1PrivateKey));
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.signHashes(Collections.singletonList("hash"), null));
    }

    @Test
    void testCheckHashesSignatureNullParameters() {
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.checkHashesSignature(null, new byte[64], user1PublicKey));
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.checkHashesSignature(Collections.singletonList("hash"), null, user1PublicKey));
        assertThrows(WhitelistException.class, () ->
                WhitelistHashHelper.checkHashesSignature(Collections.singletonList("hash"), new byte[64], null));
    }
}
