package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ReservationMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Reservation;
import com.taurushq.sdk.protect.client.model.ReservationUtxo;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ReservationsApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationUTXOReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetReservationsReply;

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
 * // Get all reservations
 * List<Reservation> reservations = client.getReservationService().getReservations();
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
     * Retrieves all reservations.
     *
     * @return the list of reservations
     * @throws ApiException if the API call fails
     */
    public List<Reservation> getReservations() throws ApiException {
        return getReservations(null, null, null, null, (String) null);
    }

    /**
     * Retrieves reservations with optional filters.
     *
     * @param kind      optional kind to filter by
     * @param address   optional address to filter by
     * @param addressId optional address ID to filter by
     * @param kinds     optional list of kinds to filter by
     * @param cursorCurrentPage optional cursor current page for pagination
     * @return the list of reservations matching the filters
     * @throws ApiException if the API call fails
     */
    public List<Reservation> getReservations(
            final String kind,
            final String address,
            final String addressId,
            final List<String> kinds,
            final String cursorCurrentPage
    ) throws ApiException {
        try {
            TgvalidatordGetReservationsReply reply = reservationsApi.walletServiceGetReservations(
                    kind, address, addressId, cursorCurrentPage, null, null, kinds);
            return reservationMapper.fromDTOList(reply.getResult());
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
