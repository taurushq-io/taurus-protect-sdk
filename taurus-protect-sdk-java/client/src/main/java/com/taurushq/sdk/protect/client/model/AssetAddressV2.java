package com.taurushq.sdk.protect.client.model;

/**
 * An address holding a v2 asset.
 * <p>
 * The row itself carries no signature. {@link #isVerified()} is true only for an INTERNAL or
 * WHITELISTED row whose address was confirmed against its verified counterpart, and the
 * address is then the counterpart's. Every other row (EXTERNAL, or no type) is on-chain
 * data that nothing signs, returned with {@code isVerified() == false}.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService#queryAssetAddresses
 */
public class AssetAddressV2 {

    private String address;
    private boolean verified;
    private String kycStatus;
    private String balance;
    private String addressType;
    private String addressId;
    private String whitelistedAddressId;

    /**
     * Gets the blockchain address.
     *
     * @return the blockchain address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the blockchain address.
     *
     * @param address the blockchain address
     */
    public void setAddress(final String address) {
        this.address = address;
    }

    /**
     * Gets the KYC status, as the wire value (e.g. KYC_STATUS_V2_APPROVED).
     *
     * @return the KYC status, as the wire value (e.g. KYC_STATUS_V2_APPROVED)
     */
    public String getKycStatus() {
        return kycStatus;
    }

    /**
     * Sets the KYC status, as the wire value (e.g. KYC_STATUS_V2_APPROVED).
     *
     * @param kycStatus the KYC status, as the wire value (e.g. KYC_STATUS_V2_APPROVED)
     */
    public void setKycStatus(final String kycStatus) {
        this.kycStatus = kycStatus;
    }

    /**
     * Gets the balance.
     *
     * @return the balance
     */
    public String getBalance() {
        return balance;
    }

    /**
     * Sets the balance.
     *
     * @param balance the balance
     */
    public void setBalance(final String balance) {
        this.balance = balance;
    }

    /**
     * Gets the address type, as the wire value (e.g. ADDRESS_TYPE_V2_INTERNAL).
     *
     * @return the address type, as the wire value (e.g. ADDRESS_TYPE_V2_INTERNAL)
     */
    public String getAddressType() {
        return addressType;
    }

    /**
     * Sets the address type, as the wire value (e.g. ADDRESS_TYPE_V2_INTERNAL).
     *
     * @param addressType the address type, as the wire value (e.g. ADDRESS_TYPE_V2_INTERNAL)
     */
    public void setAddressType(final String addressType) {
        this.addressType = addressType;
    }

    /**
     * Gets the internal address id, when the address is internal.
     *
     * @return the internal address id, when the address is internal
     */
    public String getAddressId() {
        return addressId;
    }

    /**
     * Sets the internal address id, when the address is internal.
     *
     * @param addressId the internal address id, when the address is internal
     */
    public void setAddressId(final String addressId) {
        this.addressId = addressId;
    }

    /**
     * Gets the whitelisted address id, when the address is whitelisted.
     *
     * @return the whitelisted address id, when the address is whitelisted
     */
    public String getWhitelistedAddressId() {
        return whitelistedAddressId;
    }

    /**
     * Sets the whitelisted address id, when the address is whitelisted.
     *
     * @param whitelistedAddressId the whitelisted address id, when the address is whitelisted
     */
    public void setWhitelistedAddressId(final String whitelistedAddressId) {
        this.whitelistedAddressId = whitelistedAddressId;
    }

    /**
     * Reports whether the address was confirmed against its verified counterpart: the
     * HSM-verified managed address for an INTERNAL row, the verified whitelisted address for
     * a WHITELISTED one.
     *
     * @return true when the address is verified
     */
    public boolean isVerified() {
        return verified;
    }

    /**
     * Marks this row verified, taking the address from the verified counterpart that
     * confirmed it.
     * <p>
     * The address and the flag are set together so that a row cannot be flagged verified
     * while still carrying the unsigned address it arrived with.
     *
     * @param verifiedAddress the address of the verified managed or whitelisted address
     */
    public void markVerified(final String verifiedAddress) {
        this.address = verifiedAddress;
        this.verified = true;
    }
}
