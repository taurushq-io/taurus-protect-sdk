package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddressListResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A row-level integrity failure must be EXCLUDED from a list, not abort it.
 *
 * <p>{@code IntegrityException} extends {@code SecurityException} and is therefore
 * UNCHECKED, while {@code WhitelistException} is checked. The per-row catch in
 * {@link WhitelistedAddressService#verifiedAddresses} caught only the checked one, so
 * every failure raised by step 1 (metadata hash), step 2 (rules-container signatures)
 * and step 4 (hash coverage) escaped the loop and aborted the whole page. One row with a
 * tampered payload denied the entire whitelist listing — the exact failure
 * {@code WhitelistedAddressListResult}'s javadoc says was fixed, and the failure hid its
 * own cause, since listing is how an operator finds the bad row.
 *
 * <p>Go excludes these rows, Python excludes them, TypeScript excludes them. Java was
 * the outlier and nothing tested it: every existing case in
 * {@code WhitelistedAddressApproveTest} uses bare DTOs that fail with the CHECKED
 * "signed address payload is null" and then trip the all-failed guard, so the unchecked
 * path was never exercised.
 *
 * <p>These tests drive the package-private seam directly. That is what it is
 * package-private for — this module has JUnit only, no Mockito and no HTTP stub, so
 * there is no way to drive the public list methods without a live API.
 */
class WhitelistedAddressExclusionTest {

    private static PublicKey superAdminKey;

    @BeforeAll
    static void setUp() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        superAdminKey = generator.generateKeyPair().getPublic();
    }

    private static WhitelistedAddressService service() {
        return new WhitelistedAddressService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(superAdminKey), 1);
    }

    /**
     * A row whose payload does not hash to its stated hash. Step 1 rejects it with the
     * UNCHECKED IntegrityException, which is the class that used to escape the loop.
     *
     * <p>The signedAddress payload and one signature entry are required only to clear
     * {@code initializeEnvelope}'s two preconditions, which throw the CHECKED
     * WhitelistException and would otherwise short-circuit before step 1 — the reason
     * the existing bare-DTO fixtures never exercised the unchecked path.
     */
    private static TgvalidatordSignedWhitelistedAddressEnvelope hashMismatchRow(final String id) {
        TgvalidatordSignedWhitelistedAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedAddressEnvelope();
        dto.setId(id);

        TgvalidatordMetadata metadata = new TgvalidatordMetadata();
        metadata.setPayloadAsString("{\"currency\":\"ETH\",\"address\":\"0xabc\"}");
        // Not the SHA-256 of the payload above.
        metadata.setHash("0000000000000000000000000000000000000000000000000000000000000000");
        dto.setMetadata(metadata);

        TgvalidatordSignedWhitelistedAddress signed = new TgvalidatordSignedWhitelistedAddress();
        signed.setPayload("{\"currency\":\"ETH\"}".getBytes(StandardCharsets.UTF_8));
        signed.addSignaturesItem(new TgvalidatordWhitelistSignature());
        dto.setSignedAddress(signed);

        return dto;
    }

    @Test
    @DisplayName("a hash-mismatch row is excluded by the row loop, not thrown straight out")
    void hashMismatchIsExcludedNotRethrownRaw() {
        // Before the fix this threw the RAW step-1 exception ("metadata hash
        // verification failed") from inside the loop, aborting the page. Now the row is
        // excluded; the page then fails only because NOTHING survived, which is the
        // aggregate guard and a different, weaker statement.
        IntegrityException e = assertThrows(IntegrityException.class,
                () -> service().verifiedAddresses(
                        Collections.singletonList(hashMismatchRow("1")), new HashMap<>(), null));

        assertTrue(e.getMessage().startsWith("all 1 whitelisted address(es) failed verification"),
                "expected the aggregate none-survived message, got: " + e.getMessage());
        // and the row-level reason is carried through rather than lost
        assertTrue(e.getMessage().contains("metadata hash verification failed"), e.getMessage());
    }

    @Test
    @DisplayName("every row is attempted; the loop does not stop at the first bad one")
    void allRowsAreAttempted() {
        // Two bad rows must report "all 2", not "all 1": if the unchecked exception
        // still escaped, the count would never be reached at all.
        IntegrityException e = assertThrows(IntegrityException.class,
                () -> service().verifiedAddresses(
                        Arrays.asList(hashMismatchRow("1"), hashMismatchRow("2")),
                        new HashMap<>(), null));

        assertTrue(e.getMessage().startsWith("all 2 whitelisted address(es) failed verification"),
                e.getMessage());
    }

    @Test
    @DisplayName("totalItems is reduced by the number of rows withheld")
    void totalItemsIsReducedByExclusions() throws Exception {
        // The server counts rows it returned; the caller receives only those that
        // verified. Reporting the raw total makes pagination promise rows that can never
        // be read. Java exposed no total at all, raw or adjusted, while Go reduced it and
        // this SDK's own WhitelistedAssetResult carried one.
        //
        // An empty page keeps the server's total untouched: nothing was withheld.
        WhitelistedAddressListResult result =
                service().verifiedAddresses(Collections.emptyList(), new HashMap<>(), "10");

        assertEquals("10", result.getTotalItems());
    }

    @Test
    @DisplayName("the container hash label matches validatord's convention")
    void containerHashLabelMatchesValidatordConvention() {
        // The label is base64(SHA256(<the base64 container TEXT>)) — over the base64
        // text, not the decoded protobuf, and base64 out, not hex. validatord computes
        // it in one line in its whitelist controller. Hashing the decoded bytes instead,
        // or emitting hex, would reject every container and break every list call, so
        // pin the convention itself rather than only its use.
        //
        // echo -n "abc" | sha256sum -> ba7816bf...; base64 of those raw digest bytes:
        assertEquals("ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0=",
                WhitelistedAddressService.containerHashLabel("abc"));
    }

    @Test
    @DisplayName("an absent or unparseable total stays absent rather than becoming a guess")
    void unusableTotalStaysNull() throws Exception {
        assertNull(service()
                .verifiedAddresses(Collections.emptyList(), new HashMap<>(), null)
                .getTotalItems());
        assertNull(service()
                .verifiedAddresses(Collections.emptyList(), new HashMap<>(), "")
                .getTotalItems());
        assertNull(service()
                .verifiedAddresses(Collections.emptyList(), new HashMap<>(), "not-a-number")
                .getTotalItems());
    }
}
