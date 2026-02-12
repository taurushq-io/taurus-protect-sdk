package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Base64;

import java.security.PublicKey;
import java.util.List;
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
     * Verifies that governance rules have enough valid SuperAdmin signatures.
     *
     * @param rules                the governance rules to verify
     * @param minValidSignatures   the minimum number of valid signatures required
     * @param superAdminPublicKeys the list of SuperAdmin public keys for verification
     * @throws IntegrityException if verification fails or not enough valid signatures
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

        byte[] rulesData = Base64.decodeBase64(rules.getRulesContainer());
        int validCount = 0;

        for (RuleUserSignature sig : signatures) {
            if (isValidSignature(rulesData, sig.getSignature(), superAdminPublicKeys)) {
                validCount++;
            }
        }

        if (validCount < minValidSignatures) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                LOGGER.warning("Governance rules verification failed: insufficient valid signatures");
            }
            throw createVerificationException(
                    String.format("only %d valid signatures found, minimum %d required",
                            validCount, minValidSignatures));
        }
        if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
            LOGGER.fine("Governance rules signature verification succeeded");
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
        for (PublicKey publicKey : superAdminPublicKeys) {
            try {
                if (CryptoTPV1.verifyBase64Signature(publicKey, data, signature)) {
                    if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
                        LOGGER.fine("Signature verification succeeded for SuperAdmin key");
                    }
                    return true;
                }
            } catch (java.security.SignatureException e) {
                // Signature verification failed for this key, continue trying other keys
            } catch (java.security.InvalidKeyException e) {
                // Key incompatible with algorithm, continue trying other keys
            } catch (java.security.NoSuchAlgorithmException e) {
                // Algorithm not available - this is a configuration error, fail fast
                throw new IllegalStateException("Signature algorithm not available", e);
            }
        }
        return false;
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
        } catch (java.security.SignatureException e) {
            // Invalid signature format or verification failed
            return false;
        } catch (java.security.InvalidKeyException e) {
            // Key incompatible with algorithm
            return false;
        } catch (java.security.NoSuchAlgorithmException e) {
            // Algorithm not available - this is a configuration error, fail fast
            throw new IllegalStateException("Signature algorithm not available", e);
        }
    }

    private static IntegrityException createVerificationException(String message) {
        return new IntegrityException("Governance rules verification failed: " + message);
    }
}
