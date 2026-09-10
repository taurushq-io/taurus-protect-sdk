package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Base64;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.PublicKey;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.logging.Logger;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Static helper class for verifying signatures using SuperAdmin public keys.
 */
public final class SignatureVerifier {

    private static final Logger LOGGER = Logger.getLogger(SignatureVerifier.class.getName());

    private SignatureVerifier() {
        // Prevent instantiation
    }

    /**
     * Verifies that governance rules are signed by enough DISTINCT SuperAdmin keys.
     *
     * <p>minValidSignatures counts distinct signing keys, not signature entries: ECDSA is
     * randomized, so counting entries would let a single key produce as many valid
     * signatures as any threshold demands, making 2-of-N no stronger than 1-of-N.
     *
     * @param rules                the governance rules to verify
     * @param minValidSignatures   the minimum number of distinct valid signers required
     * @param superAdminPublicKeys the list of SuperAdmin public keys for verification
     * @throws IntegrityException if verification fails or not enough distinct signers
     */
    public static void verifyGovernanceRules(GovernanceRules rules, int minValidSignatures,
                                              List<PublicKey> superAdminPublicKeys) throws IntegrityException {
        checkNotNull(rules, "rules cannot be null");
        checkArgument(minValidSignatures > 0, "minValidSignatures must be positive");
        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys cannot be empty");

        if (rules.getRulesContainer() == null) {
            throw createVerificationException("rulesContainer is null");
        }

        List<RuleUserSignature> signatures = rules.getRulesSignatures();
        if (signatures == null || signatures.isEmpty()) {
            throw createVerificationException("no signatures present");
        }

        // Strict, so the bytes a signature is checked over are unambiguous: see StrictBase64.
        byte[] rulesData = StrictBase64.decode(rules.getRulesContainer());
        verifyGovernanceRulesSignatures(rulesData, signatures, superAdminPublicKeys, minValidSignatures);

        if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
            LOGGER.fine("Governance rules signature verification succeeded");
        }
    }

    /**
     * Verifies that the rules container is signed by enough DISTINCT SuperAdmin keys.
     *
     * <p>This is the single place the threshold is evaluated — every caller routes here.
     * minValidSignatures counts distinct signing keys, not signature entries: ECDSA is
     * randomized, so counting entries would let a single key produce as many valid
     * signatures as any threshold demands, making 2-of-N no stronger than 1-of-N.
     *
     * @param rulesContainerData   the raw signed rules container bytes
     * @param signatures           the signature entries to check
     * @param superAdminPublicKeys the list of SuperAdmin public keys for verification
     * @param minValidSignatures   the minimum number of distinct valid signers required
     * @throws IntegrityException if too few distinct SuperAdmin keys signed
     */
    public static void verifyGovernanceRulesSignatures(byte[] rulesContainerData,
                                                       List<RuleUserSignature> signatures,
                                                       List<PublicKey> superAdminPublicKeys,
                                                       int minValidSignatures) throws IntegrityException {
        checkArgument(minValidSignatures > 0, "minValidSignatures must be positive");
        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys cannot be empty");

        if (rulesContainerData == null || rulesContainerData.length == 0) {
            throw createVerificationException("rules container data cannot be empty");
        }
        if (signatures == null || signatures.isEmpty()) {
            throw createVerificationException("no signatures provided");
        }

        Set<String> signers = new HashSet<String>();

        for (RuleUserSignature sig : signatures) {
            // Skip absent signatures explicitly, as the other three SDKs do. Without this,
            // not-counting relies on the JCE provider throwing SignatureException (which
            // matchingKey catches) rather than NPE for a null signature array.
            if (sig == null || sig.getSignature() == null || sig.getSignature().isEmpty()) {
                continue;
            }

            PublicKey key = matchingKey(rulesContainerData, sig.getSignature(), superAdminPublicKeys);
            if (key != null) {
                signers.add(keyFingerprint(key));
            }
        }

        if (signers.size() < minValidSignatures) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                LOGGER.warning("Governance rules verification failed: insufficient distinct valid signers");
            }
            throw createVerificationException(
                    String.format("only %d distinct valid signers found, minimum %d required",
                            signers.size(), minValidSignatures));
        }
    }

    /**
     * Verifies a signature against the provided SuperAdmin public keys.
     *
     * @param data                 the data that was signed
     * @param signature            the base64-encoded signature
     * @param superAdminPublicKeys the list of SuperAdmin public keys to verify against
     * @return true if the signature is valid for any of the SuperAdmin keys
     */
    public static boolean isValidSignature(byte[] data, String signature,
                                            List<PublicKey> superAdminPublicKeys) {
        return matchingKey(data, signature, superAdminPublicKeys) != null;
    }

    /**
     * Returns the first SuperAdmin key that verifies the signature, or null.
     */
    private static PublicKey matchingKey(byte[] data, String signature,
                                         List<PublicKey> superAdminPublicKeys) {
        for (PublicKey publicKey : superAdminPublicKeys) {
            try {
                if (CryptoTPV1.verifyBase64Signature(publicKey, data, signature)) {
                    if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
                        LOGGER.fine("Signature verification succeeded for SuperAdmin key");
                    }
                    return publicKey;
                }
            } catch (java.security.SignatureException | java.security.InvalidKeyException e) {
                // This key did not verify the signature (bad signature, or a key
                // incompatible with the algorithm) — try the next configured key.
                continue;
            } catch (java.security.NoSuchAlgorithmException e) {
                // Algorithm not available - this is a configuration error, fail fast
                throw new IllegalStateException("Signature algorithm not available", e);
            }
        }
        return null;
    }

    /**
     * Identifies a key by its encoded bytes, so the same key configured twice counts
     * as one signer.
     * <p>
     * Public because the per-group step-5 threshold in the whitelist services counts
     * distinct signers the same way, and a second copy of this is how the
     * hash-comparison helpers drifted.
     *
     * @param publicKey the key to fingerprint
     * @return a stable fingerprint of the encoded key
     */
    public static String keyFingerprint(PublicKey publicKey) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            return Base64.encodeBase64String(digest.digest(publicKey.getEncoded()));
        } catch (java.security.NoSuchAlgorithmException e) {
            // SHA-256 is mandated by the JCA spec - absence is a broken JRE
            throw new IllegalStateException("SHA-256 not available", e);
        }
    }

    /**
     * Verifies a raw signature against a single public key.
     *
     * @param data      the data that was signed
     * @param signature the raw signature bytes
     * @param publicKey the public key to verify against
     * @return true if the signature is valid
     */
    public static boolean verifySignature(byte[] data, byte[] signature, PublicKey publicKey) {
        checkNotNull(data, "data cannot be null");
        checkNotNull(signature, "signature cannot be null");
        checkNotNull(publicKey, "publicKey cannot be null");
        try {
            java.security.Signature verifier = java.security.Signature.getInstance("SHA256withPLAIN-ECDSA");
            verifier.initVerify(publicKey);
            verifier.update(data);
            return verifier.verify(signature);
        } catch (java.security.SignatureException | java.security.InvalidKeyException e) {
            // Invalid signature format, verification failed, or a key incompatible with
            // the algorithm — all mean "this signature does not verify".
            return false;
        } catch (java.security.NoSuchAlgorithmException e) {
            // Algorithm not available - this is a configuration error, fail fast
            throw new IllegalStateException("Signature algorithm not available", e);
        }
    }

    private static IntegrityException createVerificationException(String message) {
        return new IntegrityException("Governance rules verification failed: " + message);
    }

    /**
     * Reports whether a metadata hash is covered by any of the supplied whitelist
     * signatures.
     *
     * <p>The peer of Go VerifyHashCoverage, Python verify_hash_coverage and TS
     * verifyHashCoverage — this SDK was the only one without it, so a caller had to walk
     * the signature list and compare hashes itself, which is exactly where a non
     * constant-time comparison creeps in.
     *
     * <p>Comparison is constant-time and the loop deliberately does not break early, so
     * the time taken does not reveal which signature covered the hash.
     *
     * @param hash       the metadata hash to look for
     * @param signatures the whitelist signatures to search
     * @return true when at least one signature covers the hash
     */
    public static boolean verifyHashCoverage(final String hash,
                                             final List<WhitelistSignature> signatures) {
        if (hash == null || hash.isEmpty() || signatures == null) {
            return false;
        }

        boolean found = false;
        for (WhitelistSignature signature : signatures) {
            List<String> hashes = signature == null ? null : signature.getHashes();
            if (containsHash(hashes, hash)) {
                found = true;
                // No break: returning as soon as a match is found would leak which
                // signature matched through timing.
            }
        }
        return found;
    }

    /**
     * Reports whether one signature's hashes list covers the given hash.
     *
     * <p>The per-signature half of the pair — {@link #verifyHashCoverage} asks the
     * same question across every signature. These two are the only places this SDK
     * compares hash strings; the services used to carry their own
     * {@code List.contains} copies, which compare with {@code String.equals} and
     * return on the first match.
     *
     * <p>Comparison is constant-time and the loop does not break early.
     *
     * @param hashes the hashes a signature covers, may be null
     * @param hash   the hash to look for
     * @return true when the list covers the hash
     */
    public static boolean containsHash(final List<String> hashes, final String hash) {
        if (hashes == null || hash == null || hash.isEmpty()) {
            return false;
        }

        byte[] expected = hash.getBytes(StandardCharsets.UTF_8);
        boolean found = false;
        for (String candidate : hashes) {
            if (candidate == null) {
                continue;
            }
            if (MessageDigest.isEqual(candidate.getBytes(StandardCharsets.UTF_8), expected)) {
                found = true;
                // No break, as above.
            }
        }
        return found;
    }
}
