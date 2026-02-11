package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.util.ArrayList;
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
}
