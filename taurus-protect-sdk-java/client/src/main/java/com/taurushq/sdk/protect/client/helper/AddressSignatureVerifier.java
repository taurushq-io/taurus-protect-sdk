package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;

import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PublicKey;
import java.security.SignatureException;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Helper class for verifying address signatures using HSM slot public keys.
 * <p>
 * Address signatures are verified against the public key of a user
 * with the HSMSLOT role in the rules container.
 */
public final class AddressSignatureVerifier {

    private AddressSignatureVerifier() {
        // Utility class - prevent instantiation
    }

    /**
     * Verifies the address signature using the HSM slot public key from the rules container.
     * <p>
     * The signed data is the blockchain address string only. The public key is found
     * by locating a user with the HSMSLOT role in the rules container.
     *
     * @param address        the address to verify
     * @param rulesContainer the decoded rules container containing HSM public keys
     * @throws IntegrityException if the address has no signature, no user with HSMSLOT role is found,
     *                            or signature verification fails
     */
    public static void verifyAddressSignature(Address address, DecodedRulesContainer rulesContainer)
            throws IntegrityException {

        checkNotNull(address, "address cannot be null");
        checkNotNull(rulesContainer, "rulesContainer cannot be null");

        // Check for signature presence (mandatory)
        if (Strings.isNullOrEmpty(address.getSignature())) {
            throw new IntegrityException("Address " + address.getId() + " has no signature");
        }

        // Check for address string presence
        if (Strings.isNullOrEmpty(address.getAddress())) {
            throw new IntegrityException("Address " + address.getId() + " has no blockchain address to verify");
        }

        // Get the HSM public key (cached in rules container)
        PublicKey hsmPublicKey = rulesContainer.getHsmPublicKey();
        if (hsmPublicKey == null) {
            throw new IntegrityException("No user with HSMSLOT role found in rules container");
        }

        // Verify the signature
        try {
            byte[] addressData = address.getAddress().getBytes(StandardCharsets.UTF_8);
            boolean valid = CryptoTPV1.verifyBase64Signature(hsmPublicKey, addressData, address.getSignature());

            if (!valid) {
                throw new IntegrityException("Address signature verification failed for address " + address.getId());
            }
        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException e) {
            throw new IntegrityException("Address signature verification failed for address " + address.getId(), e);
        }
    }
}
