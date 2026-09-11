package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistedAssetApproval;
import com.taurushq.sdk.protect.client.model.WhitelistedAssetResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class WhitelistedAssetServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private PublicKey testPublicKey;
    private PrivateKey testPrivateKey;

    @BeforeEach
    void setUp() throws Exception {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();

        // Generate a real EC P-256 key to satisfy constructor validation
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair kp = kpg.generateKeyPair();
        testPublicKey = kp.getPublic();
        testPrivateKey = kp.getPrivate();
    }

    private WhitelistedAssetService service() {
        return new WhitelistedAssetService(apiClient, apiExceptionMapper,
                Collections.singletonList(testPublicKey), 1);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(null, apiExceptionMapper,
                        Collections.singletonList(testPublicKey), 1));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(apiClient, null,
                        Collections.singletonList(testPublicKey), 1));
    }

    @Test
    void constructor_throwsOnNullSuperAdminKeys() {
        assertThrows(NullPointerException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper, null, 1));
    }

    @Test
    void constructor_throwsOnEmptySuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper,
                        Collections.emptyList(), 1));
    }

    @Test
    void constructor_throwsOnZeroMinValidSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                new WhitelistedAssetService(apiClient, apiExceptionMapper,
                        Collections.singletonList(testPublicKey), 0));
    }

    @Test
    void getWhitelistedAsset_throwsOnZeroId() {
        WhitelistedAssetService service = new WhitelistedAssetService(
                apiClient, apiExceptionMapper,
                Collections.singletonList(testPublicKey), 1);
        assertThrows(IllegalArgumentException.class, () ->
                service.getWhitelistedAsset(0));
    }

    @Test
    void getWhitelistedAsset_throwsOnNegativeId() {
        WhitelistedAssetService service = new WhitelistedAssetService(
                apiClient, apiExceptionMapper,
                Collections.singletonList(testPublicKey), 1);
        assertThrows(IllegalArgumentException.class, () ->
                service.getWhitelistedAsset(-1));
    }

    // Approval is all-or-nothing: one signature covers every hash in the batch, so a row
    // this SDK could not verify has to stop the call before anything is signed. Argument
    // validation is what is reachable here -- the project forbids mocking, so the
    // verify-then-sign path is covered by the Go/Python/TypeScript suites.

    /**
     * Mints the content pin the way a caller does: off the result of a verified read.
     *
     * <p>{@code WhitelistedAssetApproval} has no public constructor, so this is the only
     * route to one -- which is the point of the type. A hand-built value would pin nothing
     * and the approval path refuses it.
     */
    private static WhitelistedAssetApproval reviewed(final Long... ids) {
        List<SignedWhitelistedAssetEnvelope> envelopes = new ArrayList<>();
        for (Long id : ids) {
            SignedWhitelistedAssetEnvelope envelope = new SignedWhitelistedAssetEnvelope();
            envelope.setId(id);
            WhitelistMetadata metadata = new WhitelistMetadata();
            metadata.setHash("hash-" + id);
            envelope.setMetadata(metadata);
            envelopes.add(envelope);
        }
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setAssets(envelopes);
        return result.select(Arrays.asList(ids));
    }

    @Test
    void approveWhitelistedAssets_throwsOnNullSelection() {
        assertThrows(NullPointerException.class, () ->
                service().approveWhitelistedAssets(null, testPrivateKey, "c"));
    }

    @Test
    void approveWhitelistedAssets_refusesAnEmptyPin() {
        // An empty pin must not silently restore unpinned approval, so `select` refuses to
        // mint one at all.
        WhitelistedAssetResult empty = new WhitelistedAssetResult();
        empty.setAssets(Collections.emptyList());

        IntegrityException e = assertThrows(IntegrityException.class, empty::selectAll);
        assertTrue(e.getMessage().contains("no verified"), e.getMessage());
    }

    @Test
    void approveWhitelistedAssets_throwsOnNullPrivateKey() {
        assertThrows(NullPointerException.class, () ->
                service().approveWhitelistedAssets(reviewed(1L), null, "c"));
    }

    @Test
    void approveWhitelistedAssets_throwsOnMissingComment() {
        assertThrows(IllegalArgumentException.class, () ->
                service().approveWhitelistedAssets(reviewed(1L), testPrivateKey, ""));
    }

    // A zero or negative id would be sorted into the signed array with no row behind it.
    @Test
    void approveWhitelistedAssets_throwsOnNonPositiveId() {
        assertThrows(IllegalArgumentException.class, () ->
                service().approveWhitelistedAssets(reviewed(1L, 0L), testPrivateKey, "c"));
        assertThrows(IllegalArgumentException.class, () ->
                service().approveWhitelistedAssets(reviewed(1L, -3L), testPrivateKey, "c"));
    }
}
