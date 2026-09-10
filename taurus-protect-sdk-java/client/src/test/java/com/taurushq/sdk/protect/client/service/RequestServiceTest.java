package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertThrows;

class RequestServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private RequestService requestService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        requestService = new RequestService(apiClient, apiExceptionMapper);
    }

    // --- Constructor validation ---

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new RequestService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new RequestService(apiClient, null));
    }

    // --- createInternalTransferRequest validation ---

    @Test
    void createInternalTransferRequest_throwsOnZeroFromAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferRequest(0, 1, BigInteger.ONE));
    }

    @Test
    void createInternalTransferRequest_throwsOnZeroToAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferRequest(1, 0, BigInteger.ONE));
    }

    @Test
    void createInternalTransferRequest_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferRequest(1, 2, null));
    }

    @Test
    void createInternalTransferRequest_throwsOnZeroAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferRequest(1, 2, BigInteger.ZERO));
    }

    @Test
    void createInternalTransferRequest_throwsOnNegativeAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferRequest(1, 2, BigInteger.valueOf(-1)));
    }

    // --- createExternalTransferRequest validation ---

    @Test
    void createExternalTransferRequest_throwsOnZeroFromAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferRequest(0, 1, BigInteger.ONE));
    }

    @Test
    void createExternalTransferRequest_throwsOnZeroToWhitelistedAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferRequest(1, 0, BigInteger.ONE));
    }

    @Test
    void createExternalTransferRequest_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferRequest(1, 2, null));
    }

    // --- createInternalTransferFromWalletRequest validation ---

    @Test
    void createInternalTransferFromWalletRequest_throwsOnZeroFromWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferFromWalletRequest(0, 1, BigInteger.ONE));
    }

    @Test
    void createInternalTransferFromWalletRequest_throwsOnZeroToAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferFromWalletRequest(1, 0, BigInteger.ONE));
    }

    @Test
    void createInternalTransferFromWalletRequest_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createInternalTransferFromWalletRequest(1, 2, null));
    }

    // --- createExternalTransferFromWalletRequest validation ---

    @Test
    void createExternalTransferFromWalletRequest_throwsOnZeroFromWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferFromWalletRequest(0, 1, BigInteger.ONE));
    }

    @Test
    void createExternalTransferFromWalletRequest_throwsOnZeroToWhitelistedAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferFromWalletRequest(1, 0, BigInteger.ONE));
    }

    @Test
    void createExternalTransferFromWalletRequest_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createExternalTransferFromWalletRequest(1, 2, null));
    }

    // --- getRequest validation ---

    @Test
    void getRequest_throwsOnZeroId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.getRequest(0));
    }

    @Test
    void getRequest_throwsOnNegativeId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.getRequest(-1));
    }

    // --- approveRequest / approveRequests validation ---

    @Test
    void approveRequests_throwsOnNullRequestList() {
        assertThrows(NullPointerException.class, () ->
                requestService.approveRequests(null, null));
    }

    @Test
    void approveRequests_throwsOnEmptyRequestList() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.approveRequests(Collections.emptyList(), null));
    }

    @Test
    void approveRequests_throwsOnNullMetadata() {
        Request request = new Request();
        request.setId(1);
        // metadata is null by default
        List<Request> requests = Collections.singletonList(request);
        assertThrows(NullPointerException.class, () ->
                requestService.approveRequests(requests, null));
    }

    @Test
    void approveRequests_throwsOnNullMetadataHash() {
        Request request = new Request();
        request.setId(1);
        RequestMetadata metadata = new RequestMetadata();
        metadata.setHash(null);
        request.setMetadata(metadata);
        List<Request> requests = Collections.singletonList(request);
        assertThrows(IllegalArgumentException.class, () ->
                requestService.approveRequests(requests, null));
    }

    @Test
    void approveRequests_throwsOnEmptyMetadataHash() {
        Request request = new Request();
        request.setId(1);
        RequestMetadata metadata = new RequestMetadata();
        metadata.setHash("");
        request.setMetadata(metadata);
        List<Request> requests = Collections.singletonList(request);
        assertThrows(IllegalArgumentException.class, () ->
                requestService.approveRequests(requests, null));
    }

    @Test
    void approveRequests_throwsOnNullPrivateKey() {
        Request request = new Request();
        request.setId(1);
        RequestMetadata metadata = new RequestMetadata();
        metadata.setHash("abc123");
        request.setMetadata(metadata);
        List<Request> requests = Collections.singletonList(request);
        assertThrows(NullPointerException.class, () ->
                requestService.approveRequests(requests, null));
    }

    // The signature attests to the metadata hash, so signing an unverified one means
    // attesting to a payload nothing checked. These checks run before any API call.
    @Test
    void approveRequests_refusesUnverifiedMetadata() throws Exception {
        Request request = new Request();
        request.setId(1);
        RequestMetadata metadata = new RequestMetadata();
        metadata.setHash("abc123");
        request.setMetadata(metadata);
        List<Request> requests = Collections.singletonList(request);

        PrivateKey key = generateKey();
        assertThrows(IntegrityException.class, () ->
                requestService.approveRequests(requests, key));
    }

    @Test
    void approveRequests_refusesBatchWhenOneRowIsUnverified() throws Exception {
        Request verified = new Request();
        verified.setId(1);
        RequestMetadata verifiedMetadata = new RequestMetadata();
        String verifiedPayload = "[{\"key\":\"currency\",\"value\":\"ETH\"}]";
        verifiedMetadata.setPayloadAsString(verifiedPayload);
        verifiedMetadata.setHash(CryptoTPV1.calculateHexHash(verifiedPayload));
        verifiedMetadata.verifyAndMaterialise();
        verified.setMetadata(verifiedMetadata);

        Request unverified = new Request();
        unverified.setId(2);
        RequestMetadata unverifiedMetadata = new RequestMetadata();
        unverifiedMetadata.setHash("bbb");
        unverified.setMetadata(unverifiedMetadata);

        List<Request> requests = Arrays.asList(verified, unverified);
        PrivateKey key = generateKey();
        assertThrows(IntegrityException.class, () ->
                requestService.approveRequests(requests, key));
    }

    private static PrivateKey generateKey() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        return generator.generateKeyPair().getPrivate();
    }

    // --- getRequests validation ---

    @Test
    void getRequests_throwsOnNullCursor() {
        assertThrows(NullPointerException.class, () ->
                requestService.getRequests(null, null, null, null, null));
    }

    // --- getRequestsForApproval validation ---

    @Test
    void getRequestsForApproval_throwsOnNullCursor() {
        assertThrows(NullPointerException.class, () ->
                requestService.getRequestsForApproval(null));
    }

    // --- rejectRequests validation ---

    @Test
    void rejectRequests_throwsOnNullRequestIds() {
        assertThrows(NullPointerException.class, () ->
                requestService.rejectRequests(null, "comment"));
    }

    @Test
    void rejectRequests_throwsOnEmptyRequestIds() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.rejectRequests(Collections.emptyList(), "comment"));
    }

    @Test
    void rejectRequests_throwsOnNullComment() {
        List<Long> ids = Collections.singletonList(1L);
        assertThrows(IllegalArgumentException.class, () ->
                requestService.rejectRequests(ids, null));
    }

    @Test
    void rejectRequests_throwsOnEmptyComment() {
        List<Long> ids = Collections.singletonList(1L);
        assertThrows(IllegalArgumentException.class, () ->
                requestService.rejectRequests(ids, ""));
    }

    // --- createCancelRequest validation ---

    @Test
    void createCancelRequest_throwsOnZeroAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createCancelRequest(0, 0));
    }

    @Test
    void createCancelRequest_throwsOnNegativeNonce() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createCancelRequest(1, -1));
    }

    // --- createIncomingRequest validation ---

    @Test
    void createIncomingRequest_throwsOnZeroFromExchangeId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createIncomingRequest(0, 1, BigInteger.ONE));
    }

    @Test
    void createIncomingRequest_throwsOnZeroToAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createIncomingRequest(1, 0, BigInteger.ONE));
    }

    @Test
    void createIncomingRequest_throwsOnNullAmount() {
        assertThrows(IllegalArgumentException.class, () ->
                requestService.createIncomingRequest(1, 2, null));
    }

    // --- list-path metadata verification ---
    //
    // The defect: getRequests and getRequestsForApproval returned every row without
    // verifying any of them, while getRequest verified. A payload altered in transit
    // reached the caller with no error and no flag.

    private Request requestWith(final long id, final String hash, final String payload) {
        Request r = new Request();
        r.setId(id);
        RequestMetadata md = new RequestMetadata();
        md.setHash(hash);
        md.setPayloadAsString(payload);
        // Deliberately NOT verified here: these are raw mapped rows, some tampered, and
        // verifying them is what the service under test is supposed to do.
        r.setMetadata(md);
        return r;
    }

    /**
     * The DTO shape the API returns. verifiedRequests takes DTOs, not models, so the
     * mapper stays behind the single verification seam — a test that hands it models
     * would be exercising a path the service no longer has.
     */
    private static com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest dtoWith(
            final String id, final String hash, final String payload) {
        com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest dto =
                new com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest();
        dto.setId(id);
        // The mapper rejects a null/empty status label, so a DTO fixture needs one even
        // when the test is only about metadata.
        dto.setStatus("APPROVING");
        if (hash != null || payload != null) {
            com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata md =
                    new com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata();
            md.setHash(hash);
            md.setPayloadAsString(payload);
            dto.setMetadata(md);
        }
        return dto;
    }

    /**
     * Why approveRequests re-verifies instead of trusting {@code isHashVerified()}: the
     * field is private with no setter, but any reflective deserializer — Gson included,
     * which this SDK already uses — sets it happily. So this is what a Request rebuilt
     * from JSON (a queue, a webhook, a cached blob) can look like, and trusting the flag
     * gets an attacker-chosen hash signed by the approver's real key.
     */
    @Test
    void approveRequests_refusesForgedHashVerified() throws Exception {
        String good = "[{\"key\":\"currency\",\"value\":\"BTC\"}]";
        String goodHash = com.taurushq.sdk.protect.openapi.auth.CryptoTPV1.calculateHexHash(good);

        Request forged = new Request();
        forged.setId(9L);
        RequestMetadata md = new RequestMetadata();
        md.setHash(goodHash);
        md.setPayloadAsString("[{\"key\":\"currency\",\"value\":\"ETH\"}]");
        forged.setMetadata(md);

        // Reflection stands in for any reflective deserializer. Gson would do the same
        // thing to this field, but it cannot build a Request under JDK 17 (the model
        // carries OffsetDateTime, whose internals are not open), so set it directly.
        try {
            java.lang.reflect.Field f = RequestMetadata.class.getDeclaredField("hashVerified");
            f.setAccessible(true);
            f.setBoolean(md, true);
        } catch (ReflectiveOperationException ex) {
            org.junit.jupiter.api.Assertions.fail(
                    "hashVerified is no longer a field named that; update this test", ex);
        }

        org.junit.jupiter.api.Assertions.assertTrue(forged.getMetadata().isHashVerified(),
                "fixture is not exercising the forgery: hashVerified was not set");

        IntegrityException e = org.junit.jupiter.api.Assertions.assertThrows(
                IntegrityException.class,
                // A REAL key: argument validation runs first, so a null key would trip
                // checkNotNull and this test would pass without reaching the check it is about.
                () -> requestService.approveRequests(Collections.singletonList(forged), generateKey()));
        org.junit.jupiter.api.Assertions.assertTrue(
                e.getMessage().contains("refusing to sign request 9"), e.getMessage());
    }

    @Test
    void verifiedRequestsDropsTamperedRowAndMarksTheRest() {
        String good = "[{\"key\":\"currency\",\"value\":\"BTC\"}]";
        String goodHash = com.taurushq.sdk.protect.openapi.auth.CryptoTPV1.calculateHexHash(good);

        List<Request> kept = requestService.verifiedRequests(java.util.Arrays.asList(
                dtoWith("1", goodHash, good),
                // payload altered, hash left covering the original
                dtoWith("2", goodHash, "[{\"key\":\"currency\",\"value\":\"ETH\"}]"),
                // an early-status request has no metadata yet: not a failure
                dtoWith("3", null, null)));

        org.junit.jupiter.api.Assertions.assertEquals(2, kept.size(),
                "the tampered row must be excluded");
        org.junit.jupiter.api.Assertions.assertEquals(1L, kept.get(0).getId());
        org.junit.jupiter.api.Assertions.assertEquals(3L, kept.get(1).getId());
        org.junit.jupiter.api.Assertions.assertTrue(kept.get(0).getMetadata().isHashVerified(),
                "a verified row must be marked");
    }

    @Test
    void verifiedRequestsKeepsRowsWithNothingToVerify() {
        List<Request> kept = requestService.verifiedRequests(
                Collections.singletonList(dtoWith("3", null, null)));

        org.junit.jupiter.api.Assertions.assertEquals(1, kept.size());
        org.junit.jupiter.api.Assertions.assertFalse(kept.get(0).getMetadata() != null
                && kept.get(0).getMetadata().isHashVerified());
    }

}
