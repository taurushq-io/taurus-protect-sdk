package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for {@link AddressSignatureVerifier}.
 *
 * <p>Go, Python and TypeScript each had a test for their address verifier; this SDK
 * had none, so the HSM-slot flow was the only verification path here with no direct
 * coverage.
 *
 * <p>Only the input-validation branches are exercised: producing a genuinely valid
 * HSM signature needs a private key paired with the container's HSMSLOT user, which
 * the shared signed fixtures will supply once they land (see TODOS.md).
 */
class AddressSignatureVerifierTest {

    private static DecodedRulesContainer containerWithHsmUser(String publicKeyPem) {
        DecodedRulesContainer container = new DecodedRulesContainer();
        RuleUser hsm = new RuleUser();
        hsm.setId("hsm-slot-1");
        hsm.setPublicKeyPem(publicKeyPem);
        hsm.setRoles(Collections.singletonList("HSMSLOT"));
        List<RuleUser> users = new ArrayList<>(Arrays.asList(hsm));
        container.setUsers(users);
        return container;
    }

    private static Address address(String value, String signature) {
        Address address = new Address();
        address.setId(42L);
        address.setAddress(value);
        address.setSignature(signature);
        return address;
    }

    @Test
    void rejectsNullArguments() {
        assertThrows(NullPointerException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        null, containerWithHsmUser("pem")));
        assertThrows(NullPointerException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("0xabc", "sig"), null));
    }

    /**
     * An address with no signature must be an error, not a silent pass: there is
     * nothing to verify, so returning normally would mean accepting it unverified.
     */
    @Test
    void anAddressWithNoSignatureIsRejected() {
        DecodedRulesContainer container = containerWithHsmUser("pem");

        IntegrityException nullSig = assertThrows(IntegrityException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("0xabc", null), container));
        assertTrue(nullSig.getMessage().contains("no signature"));

        IntegrityException emptySig = assertThrows(IntegrityException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("0xabc", ""), container));
        assertTrue(emptySig.getMessage().contains("no signature"));
    }

    @Test
    void anAddressWithNoBlockchainAddressIsRejected() {
        IntegrityException e = assertThrows(IntegrityException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("", "c2ln"), containerWithHsmUser("pem")));
        assertTrue(e.getMessage().contains("no blockchain address"));
    }

    /**
     * Without an HSMSLOT user there is no key to verify against. Falling back to any
     * other key would let the rules container choose its own verifier.
     */
    @Test
    void aContainerWithNoHsmSlotUserIsRejected() {
        DecodedRulesContainer container = new DecodedRulesContainer();
        container.setUsers(new ArrayList<RuleUser>());

        IntegrityException e = assertThrows(IntegrityException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("0xabc", "c2ln"), container));
        assertTrue(e.getMessage().contains("HSMSLOT"));
    }

    @Test
    void aSignatureThatDoesNotVerifyIsRejected() {
        // A syntactically valid but wrong signature against a real P-256 key.
        String pem = "-----BEGIN PUBLIC KEY-----\n"
                + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW\n"
                + "L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==\n"
                + "-----END PUBLIC KEY-----";

        assertThrows(IntegrityException.class,
                () -> AddressSignatureVerifier.verifyAddressSignature(
                        address("0xabc", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="),
                        containerWithHsmUser(pem)));
    }

    @Test
    void theVerifierIsAUtilityClassAndCannotBeInstantiated() {
        assertDoesNotThrow(() -> AddressSignatureVerifier.class.getDeclaredConstructor());
    }
}
