package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a fee payer configuration in the Taurus Protect system.
 * <p>
 * Fee payers are accounts used to pay transaction fees on behalf of other
 * addresses, commonly used for sponsored transactions on EVM-compatible
 * blockchains.
 *
 * @see FeePayerService
 */
public class FeePayer {

    private String id;
    private String tenantId;
    private String blockchain;
    private String network;
    private String name;
    private OffsetDateTime creationDate;
    private FeePayerInfo feePayerInfo;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(final String tenantId) {
        this.tenantId = tenantId;
    }

    public String getBlockchain() {
        return blockchain;
    }

    public void setBlockchain(final String blockchain) {
        this.blockchain = blockchain;
    }

    public String getNetwork() {
        return network;
    }

    public void setNetwork(final String network) {
        this.network = network;
    }

    public String getName() {
        return name;
    }

    public void setName(final String name) {
        this.name = name;
    }

    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    public FeePayerInfo getFeePayerInfo() {
        return feePayerInfo;
    }

    public void setFeePayerInfo(final FeePayerInfo feePayerInfo) {
        this.feePayerInfo = feePayerInfo;
    }
}
