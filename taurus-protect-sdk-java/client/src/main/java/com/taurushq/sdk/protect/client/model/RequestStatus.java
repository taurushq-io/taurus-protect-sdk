package com.taurushq.sdk.protect.client.model;

import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

import static com.google.common.base.Strings.isNullOrEmpty;

/**
 * Represents the status of a transaction request in its lifecycle.
 * <p>
 * Request statuses track the progress of a request from creation through
 * approval, signing, broadcasting, and confirmation on the blockchain.
 * <p>
 * Common status transitions:
 * <ul>
 *   <li>{@link #CREATED} → {@link #PENDING} → {@link #APPROVED} → {@link #BROADCASTED} → {@link #CONFIRMED}</li>
 *   <li>{@link #PENDING} → {@link #REJECTED} (if rejected by an approver)</li>
 *   <li>{@link #BROADCASTED} → {@link #PERMANENT_FAILURE} (if transaction fails)</li>
 * </ul>
 * <p>
 * Example usage:
 * <pre>{@code
 * Request request = client.getRequestService().getRequest(requestId);
 * if (request.getStatus() == RequestStatus.PENDING) {
 *     // Request is awaiting approval
 * } else if (request.getStatus() == RequestStatus.CONFIRMED) {
 *     // Transaction has been confirmed on the blockchain
 * }
 * }</pre>
 *
 * @see Request
 * @see RequestService
 */
public enum RequestStatus {

    /**
     * Request has received secondary approval.
     */
    APPROVED_2("APPROVED_2"),
    /**
     * Request has been approved and is ready for signing.
     */
    APPROVED("APPROVED"),
    /**
     * Request is in the process of being approved.
     */
    APPROVING("APPROVING"),
    /**
     * Auto prepared 2 request status.
     */
    AUTO_PREPARED_2("AUTO_PREPARED_2"),
    /**
     * Auto prepared request status.
     */
    AUTO_PREPARED("AUTO_PREPARED"),
    /**
     * Broadcasting 2 request status.
     */
    BROADCASTING_2("BROADCASTING_2"),
    /**
     * Broadcasting request status.
     */
    BROADCASTING("BROADCASTING"),
    /**
     * Broadcasted request status.
     */
    BROADCASTED("BROADCASTED"),
    /**
     * Bundle approved request status.
     */
    BUNDLE_APPROVED("BUNDLE_APPROVED"),
    /**
     * Bundle broadcasting request status.
     */
    BUNDLE_BROADCASTING("BUNDLE_BROADCASTING"),
    /**
     * Bundle ready request status.
     */
    BUNDLE_READY("BUNDLE_READY"),
    /**
     * Canceled request status.
     */
    CANCELED("CANCELED"),
    /**
     * Confirmed request status.
     */
    CONFIRMED("CONFIRMED"),
    /**
     * Created request status.
     */
    CREATED("CREATED"),
    /**
     * Diem burn mbs approved request status.
     */
    DIEM_BURN_MBS_APPROVED("DIEM_BURN_MBS_APPROVED"),
    /**
     * Diem burn mbs pending request status.
     */
    DIEM_BURN_MBS_PENDING("DIEM_BURN_MBS_PENDING"),
    /**
     * Diem mint mbs approved request status.
     */
    DIEM_MINT_MBS_APPROVED("DIEM_MINT_MBS_APPROVED"),
    /**
     * Diem mint mbs completed request status.
     */
    DIEM_MINT_MBS_COMPLETED("DIEM_MINT_MBS_COMPLETED"),
    /**
     * Diem mint mbs pending request status.
     */
    DIEM_MINT_MBS_PENDING("DIEM_MINT_MBS_PENDING"),
    /**
     * Expired request status.
     */
    EXPIRED("EXPIRED"),
    /**
     * Request was fast-approved (single-step approval).
     */
    FAST_APPROVED("FAST_APPROVED"),
    /**
     * Fast approved 2 request status.
     */
    FAST_APPROVED_2("FAST_APPROVED_2"),
    /**
     * Invalid request status (request was rejected as invalid).
     */
    INVALID("INVALID"),
    /**
     * Hsm failed 2 request status.
     */
    HSM_FAILED_2("HSM_FAILED_2"),
    /**
     * Hsm failed request status.
     */
    HSM_FAILED("HSM_FAILED"),
    /**
     * Hsm ready 2 request status.
     */
    HSM_READY_2("HSM_READY_2"),
    /**
     * Hsm ready request status.
     */
    HSM_READY("HSM_READY"),
    /**
     * Hsm signed 2 request status.
     */
    HSM_SIGNED_2("HSM_SIGNED_2"),
    /**
     * Hsm signed request status.
     */
    HSM_SIGNED("HSM_SIGNED"),
    /**
     * Manual broadcast request status.
     */
    MANUAL_BROADCAST("MANUAL_BROADCAST"),
    /**
     * Mined request status.
     */
    MINED("MINED"),
    /**
     * New request status (initial state for some blockchains).
     */
    NEW("NEW"),
    /**
     * Partially confirmed request status.
     */
    PARTIALLY_CONFIRMED("PARTIALLY_CONFIRMED"),
    /**
     * Pending request status.
     */
    PENDING("PENDING"),
    /**
     * Permanent failure request status.
     */
    PERMANENT_FAILURE("PERMANENT_FAILURE"),
    /**
     * Ready request status.
     */
    READY("READY"),
    /**
     * Rejected request status.
     */
    REJECTED("REJECTED"),
    /**
     * Sent request status.
     */
    SENT("SENT"),
    /**
     * Signet completed request status.
     */
    SIGNET_COMPLETED("SIGNET_COMPLETED"),
    /**
     * Signet pending request status.
     */
    SIGNET_PENDING("SIGNET_PENDING");


    private static final Map<String, RequestStatus> byLabel = new HashMap<>();

    static {
        for (RequestStatus e : values()) {
            byLabel.put(e.label, e);
        }
    }

    /**
     * The Label.
     */
    public final String label;

    RequestStatus(String label) {
        this.label = label;
    }

    /**
     * Value of label request status.
     *
     * @param label the label
     * @return the request status
     */
    public static RequestStatus valueOfLabel(String label) {
        if (isNullOrEmpty(label)) {
            throw new IllegalArgumentException("Request status label must not be null or empty");
        }
        RequestStatus status = byLabel.get(label.toUpperCase(Locale.ENGLISH));
        if (status == null) {
            throw new IllegalArgumentException("Unrecognized request status: '" + label + "'");
        }
        return status;
    }
}

