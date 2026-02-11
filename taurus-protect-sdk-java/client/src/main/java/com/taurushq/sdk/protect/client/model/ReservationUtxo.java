package com.taurushq.sdk.protect.client.model;

/**
 * Represents UTXO details for a reservation in the Taurus Protect system.
 * <p>
 * Contains the UTXO (Unspent Transaction Output) data associated with a reservation.
 *
 * @see Reservation
 */
public class ReservationUtxo {

    private String id;
    private String hash;
    private String script;
    private String value;
    private String blockHeight;
    private String reservedByRequestId;
    private String reservationId;
    private String valueString;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getHash() {
        return hash;
    }

    public void setHash(final String hash) {
        this.hash = hash;
    }

    public String getScript() {
        return script;
    }

    public void setScript(final String script) {
        this.script = script;
    }

    public String getValue() {
        return value;
    }

    public void setValue(final String value) {
        this.value = value;
    }

    public String getBlockHeight() {
        return blockHeight;
    }

    public void setBlockHeight(final String blockHeight) {
        this.blockHeight = blockHeight;
    }

    public String getReservedByRequestId() {
        return reservedByRequestId;
    }

    public void setReservedByRequestId(final String reservedByRequestId) {
        this.reservedByRequestId = reservedByRequestId;
    }

    public String getReservationId() {
        return reservationId;
    }

    public void setReservationId(final String reservationId) {
        this.reservationId = reservationId;
    }

    public String getValueString() {
        return valueString;
    }

    public void setValueString(final String valueString) {
        this.valueString = valueString;
    }
}
