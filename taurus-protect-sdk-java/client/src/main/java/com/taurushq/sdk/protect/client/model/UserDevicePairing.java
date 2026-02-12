package com.taurushq.sdk.protect.client.model;

/**
 * Represents a user device pairing in the Taurus Protect system.
 * <p>
 * Device pairing is used to register and authenticate user devices
 * for secure access to the Taurus Protect platform.
 *
 * @see UserDevicePairingInfo
 */
public class UserDevicePairing {

    private String pairingId;

    public String getPairingId() {
        return pairingId;
    }

    public void setPairingId(final String pairingId) {
        this.pairingId = pairingId;
    }
}
