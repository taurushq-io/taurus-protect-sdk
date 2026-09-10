package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.mapper.RequestMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestResult;
import com.taurushq.sdk.protect.client.model.RequestStatus;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.RequestsApi;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveRequestsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveRequestsRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateIncomingRequestRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateOutgoingCancelRequestRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateOutgoingRequestRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateRequestReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRequestReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRequestsV2Reply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRejectRequestsRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.SignatureException;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.logging.Logger;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing transaction requests in the Taurus Protect system.
 * <p>
 * This service provides operations for creating, approving, rejecting, and
 * querying transaction requests. Requests represent actions to be performed
 * on the blockchain, such as transfers between addresses.
 * <p>
 * The typical workflow for a transfer:
 * <ol>
 *   <li>Create a request using one of the create methods</li>
 *   <li>The request goes through an approval workflow</li>
 *   <li>Approve the request with a private key using {@link #approveRequest}</li>
 *   <li>The system broadcasts the transaction to the blockchain</li>
 * </ol>
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create an internal transfer request
 * Request request = client.getRequestService()
 *     .createInternalTransferRequest(fromAddressId, toAddressId, amount);
 *
 * // Get requests pending approval
 * ApiRequestCursor cursor = Pagination.first(50);
 * RequestResult result = client.getRequestService().getRequestsForApproval(cursor);
 *
 * // Approve a request
 * int signed = client.getRequestService().approveRequest(request, privateKey);
 * }</pre>
 *
 * @see Request
 * @see RequestStatus
 * @see TransactionService
 */
public class RequestService {

    private static final Logger LOGGER = Logger.getLogger(RequestService.class.getName());

    /**
     * JSON serializer for request signing.
     */
    private static final Gson GSON = new Gson();

    /**
     * The underlying OpenAPI client for request operations.
     */
    private final RequestsApi requestsApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Creates a new RequestService instance.
     *
     * @param openApiClient      the OpenAPI client for HTTP communication
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public RequestService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.requestsApi = new RequestsApi(openApiClient);

    }


    /**
     * Create internal transfer request.
     *
     * @param fromAddressId the source address id
     * @param toAddressId   the destination address id
     * @param amount        the amount
     * @return the created request
     * @throws ApiException the api exception
     */
    public Request createInternalTransferRequest(final long fromAddressId, final long toAddressId, final BigInteger amount) throws ApiException {
        checkArgument(fromAddressId > 0, "fromAddressId cannot be zero");
        checkArgument(toAddressId > 0, "toAddressId cannot be zero");
        checkArgument(amount != null && amount.signum() > 0, "amount cannot be null or zero");


        TgvalidatordCreateOutgoingRequestRequest request = new TgvalidatordCreateOutgoingRequestRequest();
        request.setFromAddressId(String.valueOf(fromAddressId));
        request.setToAddressId(String.valueOf(toAddressId));
        request.setAmount(amount.toString());
        try {
            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateOutgoingRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Create internal transfer request from an omnibus wallet.
     *
     * @param fromWalletId the source wallet id where the funds will be taken - the wallet must be an omnibus
     * @param toAddressId  the destination address id
     * @param amount       the amount
     * @return the created request
     * @throws ApiException the api exception
     */
    public Request createInternalTransferFromWalletRequest(final long fromWalletId, final long toAddressId, final BigInteger amount) throws ApiException {
        checkArgument(fromWalletId > 0, "fromWalletId cannot be zero");
        checkArgument(toAddressId > 0, "toAddressId cannot be zero");
        checkArgument(amount != null && amount.signum() > 0, "amount cannot be null or zero");


        TgvalidatordCreateOutgoingRequestRequest request = new TgvalidatordCreateOutgoingRequestRequest();
        request.setFromWalletId(String.valueOf(fromWalletId));
        request.setToAddressId(String.valueOf(toAddressId));
        request.setAmount(amount.toString());
        try {
            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateOutgoingRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Create external transfer request.
     *
     * @param fromAddressId          the source address id
     * @param toWhitelistedAddressId the destination whitelisted address id
     * @param amount                 the amount
     * @return the created request
     * @throws ApiException the api exception
     */
    public Request createExternalTransferRequest(final long fromAddressId, final long toWhitelistedAddressId, final BigInteger amount) throws ApiException {
        checkArgument(fromAddressId > 0, "fromAddressId cannot be zero");
        checkArgument(toWhitelistedAddressId > 0, "toWhitelistedAddressId cannot be zero");
        checkArgument(amount != null && amount.signum() > 0, "amount cannot be null or zero");


        TgvalidatordCreateOutgoingRequestRequest request = new TgvalidatordCreateOutgoingRequestRequest();
        request.setFromAddressId(String.valueOf(fromAddressId));
        request.setToWhitelistedAddressId(String.valueOf(toWhitelistedAddressId));
        request.setAmount(amount.toString());
        try {
            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateOutgoingRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Create external transfer request from an omnibus wallet.
     *
     * @param fromWalletId           the source wallet id where the funds will be taken - the wallet must be an omnibus
     * @param toWhitelistedAddressId the destination whitelisted address id
     * @param amount                 the amount
     * @return the created request
     * @throws ApiException the api exception
     */
    public Request createExternalTransferFromWalletRequest(final long fromWalletId, final long toWhitelistedAddressId, final BigInteger amount) throws ApiException {
        checkArgument(fromWalletId > 0, "fromWalletId cannot be zero");
        checkArgument(toWhitelistedAddressId > 0, "toWhitelistedAddressId cannot be zero");
        checkArgument(amount != null && amount.signum() > 0, "amount cannot be null or zero");


        TgvalidatordCreateOutgoingRequestRequest request = new TgvalidatordCreateOutgoingRequestRequest();
        request.setFromWalletId(String.valueOf(fromWalletId));
        request.setToWhitelistedAddressId(String.valueOf(toWhitelistedAddressId));
        request.setAmount(amount.toString());
        try {
            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateOutgoingRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets request.
     *
     * @param id the id
     * @return the request
     * @throws ApiException the api exception
     */
    public Request getRequest(final long id) throws ApiException {
        checkArgument(id > 0, "request id cannot be zero");

        try {
            TgvalidatordGetRequestReply reply = requestsApi.requestServiceGetRequest(String.valueOf(id));
            // Same seam every other path uses: one implementation, so they cannot drift.
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Approve requests int.
     *
     * @param requests   the requests
     * @param privateKey the private key
     * @return the int
     * @throws ApiException the api exception
     */
    public int approveRequests(final List<Request> requests, final PrivateKey privateKey) throws ApiException {
        return approveRequests(requests, privateKey, null);
    }

    /**
     * Approves requests, recording the caller's comment against the approval.
     *
     * <p>An overload rather than a new parameter on the existing method, so callers
     * keep compiling. Python and TypeScript take the comment as an optional argument;
     * it used to be hardcoded here.
     *
     * @param requests   the requests to approve
     * @param privateKey the key to sign the approval with
     * @param comment    the rationale to record; a default is used when null or empty
     * @return the number of requests approved
     * @throws ApiException if the API call fails
     */
    public int approveRequests(final List<Request> requests, final PrivateKey privateKey,
                               final String comment) throws ApiException {
        checkNotNull(requests, "requests list cannot be null");
        checkArgument(!requests.isEmpty(), "requests list cannot be empty");
        requests.forEach(r -> checkNotNull(r.getMetadata(), "request metadata cannot be null"));
        requests.forEach(r -> checkArgument(!Strings.isNullOrEmpty(r.getMetadata().getHash()), "request metadata hash cannot be null or zero"));
        checkNotNull(privateKey, "privateKey cannot be null");

        // RE-VERIFY each hash this signature will attest to. isHashVerified() is
        // deliberately NOT consulted: Gson and any other reflective deserializer can set
        // the private field, so a Request rebuilt from JSON (a queue, a webhook, a cached
        // blob) can arrive claiming true and get an attacker-chosen hash signed by the
        // approver's real key. Re-verification is a SHA-256 over a string already in hand,
        // and it is the same rule approveWhitelistedAssets gets by re-reading — so both
        // signing paths rest on one rule, not two. Last of the checks: argument validation
        // first, so a bad argument still reports as one, and absent metadata as absent.
        requests.forEach(r -> {
            try {
                verifyMetadataHash(r);
            } catch (IntegrityException e) {
                throw new IntegrityException(String.format(
                        "refusing to sign request %d: %s", r.getId(), e.getMessage()));
            }
        });


        // Sorted on a COPY. Sorting the caller's list in place mutated their argument
        // and threw UnsupportedOperationException on an immutable one; Go, Python and
        // TypeScript all sort a copy.
        List<Request> sorted = new ArrayList<>(requests);
        sorted.sort(Comparator.comparingLong(Request::getId));

        TgvalidatordApproveRequestsRequest request = new TgvalidatordApproveRequestsRequest();
        request.setIds(sorted.stream().map(r -> Long.toString(r.getId())).collect(Collectors.toList()));
        request.setComment(Strings.isNullOrEmpty(comment)
                ? "approving via taurus-protect-sdk-java" : comment);
        String toSign = GSON.toJson(sorted.stream().map(r -> r.getMetadata().getHash()).collect(Collectors.toList()));

        try {
            request.setSignature(CryptoTPV1.calculateBase64Signature(privateKey, toSign.getBytes(StandardCharsets.UTF_8)));
        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException ex) {
            ApiException e = new ApiException();
            e.setCode(400);
            e.setError("ClientInvalidRequest");
            e.setMessage(String.format("unable to sign the array of hashes '%s'", toSign));
            e.setOriginalException(ex);
            throw e;
        }

        try {
            TgvalidatordApproveRequestsReply reply = requestsApi.requestServiceApproveRequests(request);
            if (reply.getSignedRequests() != null) {
                return Integer.parseInt(reply.getSignedRequests());
            }
            return 0;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Approve request int.
     *
     * @param r          the r
     * @param privateKey the private key
     * @return the int
     * @throws ApiException the api exception
     */
    public int approveRequest(final Request r, final PrivateKey privateKey) throws ApiException {
        return approveRequest(r, privateKey, null);
    }

    /**
     * Approves a single request, recording the caller's comment.
     *
     * @param r          the request to approve
     * @param privateKey the key to sign the approval with
     * @param comment    the rationale to record; a default is used when null or empty
     * @return the number of requests approved
     * @throws ApiException if the API call fails
     */
    public int approveRequest(final Request r, final PrivateKey privateKey,
                              final String comment) throws ApiException {
        return approveRequests(Collections.singletonList(r), privateKey, comment);
    }


    /**
     * Keeps only the requests whose metadata integrity verifies.
     *
     * <pre>
     *   rows -&gt; verifyMetadataHash -+- ok    -&gt; kept, marked hashVerified
     *                               +- fails -&gt; excluded + logged
     * </pre>
     *
     * <p>The list endpoints previously returned every row without verifying any of
     * them, while {@code getRequest} verified — so a payload altered in transit
     * reached the caller with no error and no flag. Excluding rather than throwing
     * keeps one bad row from denying access to every good one; logging keeps a
     * shortened list from passing as a complete one.
     *
     * <p>The log carries metadata only: an id and a reason, never the payload.
     *
     * <p>Package-private so the exclusion behaviour can be tested directly; this
     * class has no seam for injecting a stub API. Takes DTOs, not models, so the
     * mapper stays behind {@link #verifiedRequest}.
     *
     * @param dtos the request DTOs as returned by the API
     * @return only those whose metadata verified, plus those with no metadata to verify
     */
    List<Request> verifiedRequests(final List<TgvalidatordRequest> dtos) {
        List<Request> kept = new ArrayList<>(dtos.size());
        for (TgvalidatordRequest dto : dtos) {
            if (dto == null) {
                continue;
            }
            try {
                kept.add(verifiedRequest(dto));
            } catch (IntegrityException e) {
                if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                    LOGGER.warning(String.format(
                            "request excluded: metadata integrity verification failed (request_id=%s, reason=%s)",
                            dto.getId(), e.getMessage()));
                }
            }
        }
        return kept;
    }

    /**
     * Maps one request DTO and verifies its metadata. The ONLY way this class builds
     * a {@link Request}.
     *
     * <pre>
     *   getRequest ----------+
     *   getRequests(ForApproval)-+
     *   create* (x6) --------+--&gt; verifiedRequest -&gt; fromDTO -&gt; verifyAndMaterialise
     *                                    ok &lt;-----+-----&gt; IntegrityException
     *                             (hashVerified set)      (names the request id)
     * </pre>
     *
     * <p>Verification lived in {@code getRequest} alone. The list paths skipped it in
     * one pass and the six create paths in the next — both times because "remember to
     * verify" was a rule rather than the only available construction path. So do NOT
     * call {@code RequestMapper.INSTANCE.fromDTO} anywhere else in this class.
     *
     * <p>The id goes into the exception because a failing create has already succeeded
     * server-side; the caller needs it to reconcile rather than retrying and creating
     * a second request.
     *
     * @param dto the request DTO
     * @return the verified request, marked
     * @throws IntegrityException if the metadata hash does not cover the payload
     */
    Request verifiedRequest(final TgvalidatordRequest dto) {
        Request r = RequestMapper.INSTANCE.fromDTO(dto);
        if (r == null) {
            throw new IntegrityException("request not found");
        }
        try {
            verifyMetadataHash(r);
        } catch (IntegrityException e) {
            throw new IntegrityException(String.format("request %s: %s", r.getId(), e.getMessage()));
        }
        return r;
    }

    /**
     * Verifies that a request's metadata hash covers its payload string.
     *
     * <p>Extracted from {@code getRequest} so both the single-request and list paths
     * enforce the same check. Keeping it inline in one method is what let the list
     * paths skip it.
     *
     * @param r the request to verify
     * @return true if metadata was present and verified, false if there was nothing
     *         to verify (an early-status request has no metadata yet)
     * @throws IntegrityException if the hash does not cover the payload
     */
    private boolean verifyMetadataHash(final Request r) {
        return r.getMetadata() != null && r.getMetadata().verifyAndMaterialise();
    }

    /**
     * Gets requests with filtering.
     *
     * @param from       filter requests created after this date (optional)
     * @param to         filter requests created before this date (optional)
     * @param currencyId filter by currency ID or symbol (optional)
     * @param statuses   filter by request statuses (optional)
     * @param cursor     the request cursor for pagination
     * @return the request result with list and response cursor
     * @throws ApiException the api exception
     */
    public RequestResult getRequests(final OffsetDateTime from, final OffsetDateTime to,
                                     final String currencyId, final List<RequestStatus> statuses,
                                     final ApiRequestCursor cursor) throws ApiException {

        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            List<String> statusStrings = statuses != null
                    ? statuses.stream().map(s -> s.label).collect(Collectors.toList())
                    : null;

            TgvalidatordGetRequestsV2Reply reply = requestsApi.requestServiceGetRequestsV2(
                    from,
                    to,
                    currencyId,
                    statusStrings,
                    null,                               // types
                    null,                               // ids
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // sortOrder
                    null                                // excludeStatuses
            );

            RequestResult result = new RequestResult();

            List<TgvalidatordRequest> requests = reply.getResult();
            if (requests == null) {
                result.setRequests(Collections.emptyList());
            } else {
                result.setRequests(verifiedRequests(requests));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets requests pending approval.
     *
     * @param cursor the request cursor for pagination
     * @return the request result with list and response cursor
     * @throws ApiException the api exception
     */
    public RequestResult getRequestsForApproval(final ApiRequestCursor cursor) throws ApiException {

        checkNotNull(cursor, "cursor cannot be null");

        TgvalidatordRequestCursor requestCursor = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        try {
            TgvalidatordGetRequestsV2Reply reply = requestsApi.requestServiceGetRequestsForApprovalV2(
                    null,                               // currencyID
                    null,                               // types
                    null,                               // ids
                    requestCursor.getCurrentPage(),     // cursorCurrentPage
                    requestCursor.getPageRequest(),     // cursorPageRequest
                    requestCursor.getPageSize(),        // cursorPageSize
                    null,                               // sortOrder
                    null,                               // excludeTypes
                    null                                // statuses
            );

            RequestResult result = new RequestResult();

            List<TgvalidatordRequest> requests = reply.getResult();
            if (requests == null) {
                result.setRequests(Collections.emptyList());
            } else {
                result.setRequests(verifiedRequests(requests));
            }

            result.setCursor(ApiResponseCursorMapper.INSTANCE.fromDTO(reply.getCursor()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Rejects requests.
     *
     * @param requestIds the request ids to reject
     * @param comment    the rejection comment
     * @throws ApiException the api exception
     */
    public void rejectRequests(final List<Long> requestIds, final String comment) throws ApiException {

        checkNotNull(requestIds, "requestIds list cannot be null");
        checkArgument(!requestIds.isEmpty(), "requestIds list cannot be empty");
        checkArgument(!Strings.isNullOrEmpty(comment), "comment cannot be null or empty");

        try {
            TgvalidatordRejectRequestsRequest request = new TgvalidatordRejectRequestsRequest();
            request.setIds(requestIds.stream().map(String::valueOf).collect(Collectors.toList()));
            request.setComment(comment);

            requestsApi.requestServiceRejectRequests(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Rejects a single request.
     *
     * @param requestId the request id to reject
     * @param comment   the rejection comment
     * @throws ApiException the api exception
     */
    public void rejectRequest(final long requestId, final String comment) throws ApiException {
        rejectRequests(Collections.singletonList(requestId), comment);
    }


    /**
     * Creates a cancel request for a pending transaction.
     *
     * @param addressId the address id
     * @param nonce     the nonce of the transaction to cancel
     * @return the created cancel request
     * @throws ApiException the api exception
     */
    public Request createCancelRequest(final long addressId, final long nonce) throws ApiException {

        checkArgument(addressId > 0, "addressId cannot be zero");
        checkArgument(nonce >= 0, "nonce cannot be negative");

        try {
            TgvalidatordCreateOutgoingCancelRequestRequest request = new TgvalidatordCreateOutgoingCancelRequestRequest();
            request.setAddressId(String.valueOf(addressId));
            request.setNonce(String.valueOf(nonce));

            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateOutgoingCancelRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Creates an incoming request to log an incoming transaction from an exchange.
     *
     * @param fromExchangeId the source exchange id
     * @param toAddressId    the destination address id
     * @param amount         the amount
     * @return the created request
     * @throws ApiException the api exception
     */
    public Request createIncomingRequest(final long fromExchangeId, final long toAddressId,
                                          final BigInteger amount) throws ApiException {

        checkArgument(fromExchangeId > 0, "fromExchangeId cannot be zero");
        checkArgument(toAddressId > 0, "toAddressId cannot be zero");
        checkArgument(amount != null && amount.signum() > 0, "amount cannot be null or zero");

        try {
            TgvalidatordCreateIncomingRequestRequest request = new TgvalidatordCreateIncomingRequestRequest();
            request.setFromExchangeId(String.valueOf(fromExchangeId));
            request.setToAddressId(String.valueOf(toAddressId));
            request.setAmount(amount.toString());

            TgvalidatordCreateRequestReply reply = requestsApi.requestServiceCreateIncomingRequest(request);
            return verifiedRequest(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
