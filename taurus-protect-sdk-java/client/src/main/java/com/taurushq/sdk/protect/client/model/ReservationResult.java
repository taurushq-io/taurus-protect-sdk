package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of reservations from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.ReservationService
 */
public class ReservationResult extends CursorPagedResult {

    private List<Reservation> reservations;

    /**
     * Gets the reservations of this page.
     *
     * @return the reservations
     */
    public List<Reservation> getReservations() {
        return reservations;
    }

    /**
     * Sets the reservations of this page.
     *
     * @param reservations the reservations
     */
    public void setReservations(final List<Reservation> reservations) {
        this.reservations = reservations;
    }
}
