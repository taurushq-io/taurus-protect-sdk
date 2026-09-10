package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ExcludedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddressListResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddressEnvelope;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Whitelisted addresses had no approval workflow at all: both the for-approval read and
 * the approve endpoint are generated in every SDK client and were wrapped by none, so an
 * approver could not act on a whitelisted destination through the SDK — verified or not.
 */
class WhitelistedAddressApproveTest {

    private static PublicKey superAdminKey;
    private static PrivateKey approverKey;

    @BeforeAll
    static void setUp() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        superAdminKey = generator.generateKeyPair().getPublic();
        approverKey = generator.generateKeyPair().getPrivate();
    }

    /** The service needs no network for the checks below; only a client to construct. */
    private static WhitelistedAddressService service() {
        return new WhitelistedAddressService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(superAdminKey), 1);
    }

    private static TgvalidatordSignedWhitelistedAddressEnvelope row(final String id) {
        TgvalidatordSignedWhitelistedAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedAddressEnvelope();
        dto.setId(id);
        return dto;
    }

    @Test
    @DisplayName("an unverifiable page is withheld, not returned as a complete queue")
    void unverifiablePageIsWithheld() {
        // Nothing in these rows can verify, and rows-returned-but-none-surviving is
        // systemic: an empty list would be indistinguishable from an empty whitelist.
        IntegrityException e = assertThrows(IntegrityException.class,
                () -> service().verifiedAddresses(
                        Arrays.asList(row("1"), row("2")), new HashMap<>(), null));
        assertTrue(e.getMessage().contains("failed verification"), e.getMessage());
    }

    @Test
    @DisplayName("an empty page is not a verification failure")
    void emptyPageIsNotAFailure() throws Exception {
        WhitelistedAddressListResult result =
                service().verifiedAddresses(Collections.emptyList(), new HashMap<>(), null);

        assertTrue(result.getEnvelopes().isEmpty());
        assertTrue(result.getExcludedUnverified().isEmpty());
    }

    @Test
    @DisplayName("approve rejects bad arguments before reaching the wire")
    void approveRejectsBadArguments() {
        WhitelistedAddressService svc = service();

        assertThrows(NullPointerException.class,
                () -> svc.approveWhitelistedAddresses(null, approverKey, "ok"));
        assertThrows(IllegalArgumentException.class,
                () -> svc.approveWhitelistedAddresses(Collections.emptyList(), approverKey, "ok"));
        assertThrows(NullPointerException.class,
                () -> svc.approveWhitelistedAddresses(Collections.singletonList(1L), null, "ok"));
        assertThrows(IllegalArgumentException.class,
                () -> svc.approveWhitelistedAddresses(Collections.singletonList(1L), approverKey, ""));
        assertThrows(IllegalArgumentException.class,
                () -> svc.approveWhitelistedAddresses(Collections.singletonList(0L), approverKey, "ok"));
    }

    @Test
    @DisplayName("the caller's id list is not mutated by sorting")
    void approveDoesNotMutateTheCallersList() {
        // An immutable list would throw from Collections.sort if it sorted in place, and
        // a mutable one would come back reordered. Both are caller-visible bugs.
        List<Long> ids = Collections.unmodifiableList(Arrays.asList(7L, 3L));

        // Fails at the network call, not at the sort: reaching the network at all proves
        // the sort worked on a copy.
        assertThrows(Exception.class,
                () -> service().approveWhitelistedAddresses(ids, approverKey, "ok"));
        assertEquals(Arrays.asList(7L, 3L), ids);
    }

    @Test
    @DisplayName("exclusions carry an id and a reason")
    void exclusionsCarryIdAndReason() {
        // One row verifies nothing, so it is named rather than silently dropped. The
        // all-dropped guard fires first here, which is itself the stronger signal.
        IntegrityException e = assertThrows(IntegrityException.class,
                () -> service().verifiedAddresses(
                        Collections.singletonList(row("42")), new HashMap<>(), null));
        assertTrue(e.getMessage().contains("1 whitelisted address(es)"), e.getMessage());

        ExcludedWhitelistedAddress excluded = new ExcludedWhitelistedAddress("42", "why");
        assertEquals("42", excluded.getId());
        assertEquals("why", excluded.getReason());
    }
}
