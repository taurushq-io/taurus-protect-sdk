package com.taurushq.sdk.protect.client.model;

/**
 * The Canton native-token details of a v2 asset.
 */
public class CantonNativeTokenV2 {

    private String instrumentId;
    private String cid;
    private Boolean requireCredentials;
    private Boolean paused;
    private String operator;

    /**
     * Gets the instrument id.
     *
     * @return the instrument id
     */
    public String getInstrumentId() {
        return instrumentId;
    }

    /**
     * Sets the instrument id.
     *
     * @param instrumentId the instrument id
     */
    public void setInstrumentId(final String instrumentId) {
        this.instrumentId = instrumentId;
    }

    /**
     * Gets the instrument configuration contract id.
     *
     * @return the instrument configuration contract id
     */
    public String getCid() {
        return cid;
    }

    /**
     * Sets the instrument configuration contract id.
     *
     * @param cid the instrument configuration contract id
     */
    public void setCid(final String cid) {
        this.cid = cid;
    }

    /**
     * Gets the whether holders need credentials.
     *
     * @return the whether holders need credentials
     */
    public Boolean getRequireCredentials() {
        return requireCredentials;
    }

    /**
     * Sets the whether holders need credentials.
     *
     * @param requireCredentials the whether holders need credentials
     */
    public void setRequireCredentials(final Boolean requireCredentials) {
        this.requireCredentials = requireCredentials;
    }

    /**
     * Gets the whether the instrument is paused.
     *
     * @return the whether the instrument is paused
     */
    public Boolean getPaused() {
        return paused;
    }

    /**
     * Sets the whether the instrument is paused.
     *
     * @param paused the whether the instrument is paused
     */
    public void setPaused(final Boolean paused) {
        this.paused = paused;
    }

    /**
     * Gets the instrument operator.
     *
     * @return the instrument operator
     */
    public String getOperator() {
        return operator;
    }

    /**
     * Sets the instrument operator.
     *
     * @param operator the instrument operator
     */
    public void setOperator(final String operator) {
        this.operator = operator;
    }
}
