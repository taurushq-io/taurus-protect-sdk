package com.taurushq.sdk.protect.openapi.auth;

import com.google.common.base.Strings;
import org.apache.commons.codec.binary.Base64;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.codec.digest.DigestUtils;
import org.apache.commons.codec.digest.HmacUtils;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.PEMKeyPair;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;

import java.io.IOException;
import java.io.StringReader;
import java.io.StringWriter;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.security.Signature;
import java.security.SignatureException;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import static org.apache.commons.codec.digest.HmacAlgorithms.HMAC_SHA_256;
import static org.bouncycastle.jce.provider.BouncyCastleProvider.PROVIDER_NAME;
import static org.bouncycastle.util.Strings.constantTimeAreEqual;

public class CryptoTPV1 {
    static {
        if (Security.getProvider(PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    public static PrivateKey decodePrivateKey(String pemPrivateKey) throws IOException {
        Object parsed = new PEMParser(new StringReader(pemPrivateKey)).readObject();
        return new JcaPEMKeyConverter().getKeyPair((PEMKeyPair) parsed).getPrivate();
    }

    public static PublicKey decodePublicKey(String pemPublicKey) throws IOException {
        Object parsed = new PEMParser(new StringReader(pemPublicKey)).readObject();
        PublicKey publicKey = new JcaPEMKeyConverter().getPublicKey((org.bouncycastle.asn1.x509.SubjectPublicKeyInfo) parsed);

        // Validate P-256 curve for ECDSA keys — weaker curves (e.g. P-192) provide insufficient security
        if (publicKey instanceof java.security.interfaces.ECPublicKey) {
            java.security.spec.ECParameterSpec params = ((java.security.interfaces.ECPublicKey) publicKey).getParams();
            if (params.getOrder().bitLength() != 256) {
                throw new IllegalArgumentException("Only P-256 (secp256r1) curve is supported");
            }
        }

        return publicKey;
    }

    public static String encodePublicKey(PublicKey publicKey) throws IOException {
        StringWriter stringWriter = new StringWriter();
        try (JcaPEMWriter pemWriter = new JcaPEMWriter(stringWriter)) {
            pemWriter.writeObject(publicKey);
        }
        return stringWriter.toString();
    }

    /**
     * Decodes a list of PEM-encoded public keys.
     *
     * @param pemPublicKeys the list of PEM-encoded public keys
     * @return the list of decoded public keys
     * @throws IOException if a PEM key cannot be decoded
     */
    public static List<PublicKey> decodePublicKeys(List<String> pemPublicKeys) throws IOException {
        List<PublicKey> keys = new ArrayList<>();
        for (String pem : pemPublicKeys) {
            keys.add(decodePublicKey(pem));
        }
        return keys;
    }

    public static String calculateBase64Signature(PrivateKey privateKey, byte[] data) throws NoSuchAlgorithmException, InvalidKeyException, SignatureException {
        Signature signer = Signature.getInstance("SHA256withPLAIN-ECDSA");
        signer.initSign(privateKey);
        signer.update(data);
        return Base64.encodeBase64String(signer.sign());
    }

    public static boolean verifyBase64Signature(PublicKey publicKey, byte[] data, String base64Signature) throws NoSuchAlgorithmException, InvalidKeyException, SignatureException {
        Signature verifier = Signature.getInstance("SHA256withPLAIN-ECDSA");
        verifier.initVerify(publicKey);
        verifier.update(data);
        return verifier.verify(Base64.decodeBase64(base64Signature));
    }

    public static String calculateHexHash(String data) {
        return new DigestUtils(DigestUtils.getSha256Digest()).digestAsHex(data);
    }

    public static String calculateBase64Hmac(byte[] secret, String data) {
        return Base64.encodeBase64String(new HmacUtils(HMAC_SHA_256, secret).hmac(data));
    }

    public static boolean checkBase64Hmac(byte[] secret, String data, String base64Hmac) {
        return constantTimeAreEqual(calculateBase64Hmac(secret, data), base64Hmac);
    }

    public static String calculateHexHmac(String secret, String data) {
        return Hex.encodeHexString(new HmacUtils(HMAC_SHA_256, secret).hmac(data));
    }

    public static boolean checkHexHmac(String secret, String data, String hexHmac) {
        return constantTimeAreEqual(calculateHexHmac(secret, data), hexHmac);
    }

    public static String calculateSignedHeader(String apiKey, byte[] apiSecret, String nonce, long timestamp, String method, String host, String path, String query, String contentType, String body) {
        String msg = Stream.of("TPV1", apiKey, nonce, Long.toString(timestamp), method, host, path, query, contentType, body)
                .filter(p -> !Strings.isNullOrEmpty(p))
                .collect(Collectors.joining(" "));

        return String.format("TPV1-HMAC-SHA256 ApiKey=%s Nonce=%s Timestamp=%s Signature=%s", apiKey, nonce, timestamp, calculateBase64Hmac(apiSecret, msg));
    }
}
