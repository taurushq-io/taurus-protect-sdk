package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ReservationMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.Reservation;
import com.taurushq.sdk.protect.client.model.ReservationResult;
import com.taurushq.sdk.protect.client.model.ReservationUtxo;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ReservationsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationUTXOReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationsReply;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing UTXO reservations in the Taurus Protect system.
 * <p>
 * Reservations are used to lock specific UTXOs (Unspent Transaction Outputs)
 * for UTXO-based blockchains like Bitcoin and Litecoin, preventing
 * double-spending during transaction creation.
 * <p>
 * Example usage:
 * <pre>{@code
 * // First page of reservations (default page size)
 * ReservationResult reservations = client.getReservationService().getReservations();
 * // next page: getReservations(null, null, null, null, 20, reservations.getPage().getNextCursor())
 *
 * // Get a specific reservation
 * Reservation reservation = client.getReservationService().getReservation("res-123");
 *
 * // Get UTXO details for a reservation
 * ReservationUtxo utxo = client.getReservationService().getReservationUtxo("res-123");
 * }</pre>
 *
 * @see Reservation
 * @see ReservationUtxo
 */
public class ReservationService {

    private final ReservationsApi reservationsApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final ReservationMapper reservationMapper;

    /**
     * Creates a new ReservationService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public ReservationService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.reservationsApi = new ReservationsApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.reservationMapper = ReservationMapper.INSTANCE;
    }

    /**
     * Retrieves the first page of reservations, with the default page size.
     *
     * @return the reservations and their page
     * @throws ApiException if the API call fails
     */
    public ReservationResult getReservations() throws ApiException {
        return getReservations(null, null, null, null, (ApiRequestCursor) null);
    }

    /**
     * Retrieves a page of reservations with optional filters.
     *
     * @param kind      optional kind to filter by
     * @param address   optional address to filter by
     * @param addressId optional address ID to filter by
     * @param kinds     optional list of kinds to filter by
     * @param pageSize  the page size, null or 0 for the default
     * @param cursor    a previous page's {@code getPage().getNextCursor()}, null for the first page
     * @return the reservations and their page
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ReservationResult getReservations(
            final String kind,
            final String address,
            final String addressId,
            final List<String> kinds,
            final Integer pageSize,
            final String cursor
    ) throws ApiException {
        return getReservations(kind, address, addressId, kinds, Pagination.page(pageSize, cursor));
    }

    /**
     * Retrieves a page of reservations with optional filters, with a low-level request cursor.
     *
     * @param kind      optional kind to filter by
     * @param address   optional address to filter by
     * @param addressId optional address ID to filter by
     * @param kinds     optional list of kinds to filter by
     * @param cursor    the request cursor, null for the first page with the default size
     * @return the reservations and their page
     * @throws ApiException if the API call fails
     */
    public ReservationResult getReservations(
            final String kind,
            final String address,
            final String addressId,
            final List<String> kinds,
            final ApiRequestCursor cursor
    ) throws ApiException {
        final CursorRequest page = CursorRequest.of(cursor);
        try {
            TgvalidatordGetReservationsReply reply = reservationsApi.walletServiceGetReservations(
                    kind, address, addressId,
                    page.currentPage(), page.pageRequest(), page.pageSizeParam(),
                    kinds);
            ReservationResult result = new ReservationResult();
            result.setReservations(reply.getResult() == null
                    ? Collections.emptyList() : reservationMapper.fromDTOList(reply.getResult()));
            return page.complete(result, PagedOperation.RESERVATIONS, reply.getCursor(), null);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves a reservation by ID.
     *
     * @param id the reservation ID
     * @return the reservation
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public Reservation getReservation(final String id) throws ApiException {
        checkArgument(id != null && !id.isEmpty(), "id must not be null or empty");
        try {
            TgvalidatordGetReservationReply reply = reservationsApi.walletServiceGetReservation(id);
            return reservationMapper.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves the UTXO details for a reservation.
     *
     * @param id the reservation ID
     * @return the UTXO details
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public ReservationUtxo getReservationUtxo(final String id) throws ApiException {
        checkArgument(id != null && !id.isEmpty(), "id must not be null or empty");
        try {
            TgvalidatordGetReservationUTXOReply reply = reservationsApi.walletServiceGetReservationUTXO(id);
            return reservationMapper.fromUtxoDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
