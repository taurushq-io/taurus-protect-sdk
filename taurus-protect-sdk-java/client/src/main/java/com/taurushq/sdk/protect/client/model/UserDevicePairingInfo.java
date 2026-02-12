package com.taurushq.sdk.protect.client.model;

/**
 * Represents the status and information of a user device pairing.
 * <p>
 * Contains the current status of the pairing process and, upon success,
 * the API key for the paired device.
 *
 * @see UserDevicePairing
 */
public class UserDevicePairingInfo {

    private String status;
    private String pairingId;
    private String apiKey;

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getPairingId() {
        return pairingId;
    }

    public void setPairingId(final String pairingId) {
        this.pairingId = pairingId;
    }

    public String getApiKey() {
        return apiKey;
    }

    public void setApiKey(final String apiKey) {
        this.apiKey = apiKey;
    }
}
