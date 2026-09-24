package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * An operation on a v2 asset (create, update, import, mint, burn, pause, unpause, account
 * pause/unpause, set KYC).
 * <p>
 * The operation-specific details are flattened: only the fields of this operation's
 * {@link #getType() type} are set.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService#listAssetOperations
 */
public class AssetOperationV2 {

    private String id;
    private String assetId;
    private String type;
    private String status;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;
    private String initiatedByAddressId;
    private String failureReason;
    private String blockingReason;
    private String label;
    private String price;
    private String decimals;
    private String blockchain;
    private String network;
    private String assetType;
    private String address;
    private String instrumentId;
    private String name;
    private String symbol;
    private Boolean requireCredentials;
    private String amount;
    private AddressTargetV2 destination;
    private List<String> nftTokenIds;
    private List<byte[]> nftMetadata;
    private AddressTargetV2 target;
    private String kycStatus;

    /**
     * Gets the operation id.
     *
     * @return the operation id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the operation id.
     *
     * @param id the operation id
     */
    public void setId(final String id) {
        this.id = id;
    }

    /**
     * Gets the asset id.
     *
     * @return the asset id
     */
    public String getAssetId() {
        return assetId;
    }

    /**
     * Sets the asset id.
     *
     * @param assetId the asset id
     */
    public void setAssetId(final String assetId) {
        this.assetId = assetId;
    }

    /**
     * Gets the operation type, as the wire value (e.g. ASSET_OPERATION_TYPE_V2_MINT).
     *
     * @return the operation type, as the wire value (e.g. ASSET_OPERATION_TYPE_V2_MINT)
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the operation type, as the wire value (e.g. ASSET_OPERATION_TYPE_V2_MINT).
     *
     * @param type the operation type, as the wire value (e.g. ASSET_OPERATION_TYPE_V2_MINT)
     */
    public void setType(final String type) {
        this.type = type;
    }

    /**
     * Gets the operation status, as the wire value (e.g. ASSET_OPERATION_STATUS_V2_COMPLETED).
     *
     * @return the operation status, as the wire value (e.g. ASSET_OPERATION_STATUS_V2_COMPLETED)
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the operation status, as the wire value (e.g. ASSET_OPERATION_STATUS_V2_COMPLETED).
     *
     * @param status the operation status, as the wire value (e.g. ASSET_OPERATION_STATUS_V2_COMPLETED)
     */
    public void setStatus(final String status) {
        this.status = status;
    }

    /**
     * Gets the creation date.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the creation date.
     *
     * @param createdAt the creation date
     */
    public void setCreatedAt(final OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    /**
     * Gets the last update date.
     *
     * @return the last update date
     */
    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    /**
     * Sets the last update date.
     *
     * @param updatedAt the last update date
     */
    public void setUpdatedAt(final OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    /**
     * Gets the id of the address that initiated the operation.
     *
     * @return the id of the address that initiated the operation
     */
    public String getInitiatedByAddressId() {
        return initiatedByAddressId;
    }

    /**
     * Sets the id of the address that initiated the operation.
     *
     * @param initiatedByAddressId the id of the address that initiated the operation
     */
    public void setInitiatedByAddressId(final String initiatedByAddressId) {
        this.initiatedByAddressId = initiatedByAddressId;
    }

    /**
     * Gets the failure reason, as the wire value, when the operation failed.
     *
     * @return the failure reason, as the wire value, when the operation failed
     */
    public String getFailureReason() {
        return failureReason;
    }

    /**
     * Sets the failure reason, as the wire value, when the operation failed.
     *
     * @param failureReason the failure reason, as the wire value, when the operation failed
     */
    public void setFailureReason(final String failureReason) {
        this.failureReason = failureReason;
    }

    /**
     * Gets the blocking reason, as the wire value, when the operation is blocked.
     *
     * @return the blocking reason, as the wire value, when the operation is blocked
     */
    public String getBlockingReason() {
        return blockingReason;
    }

    /**
     * Sets the blocking reason, as the wire value, when the operation is blocked.
     *
     * @param blockingReason the blocking reason, as the wire value, when the operation is blocked
     */
    public void setBlockingReason(final String blockingReason) {
        this.blockingReason = blockingReason;
    }

    /**
     * Gets the label (create, update, import).
     *
     * @return the label (create, update, import)
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the label (create, update, import).
     *
     * @param label the label (create, update, import)
     */
    public void setLabel(final String label) {
        this.label = label;
    }

    /**
     * Gets the price (create, update, import).
     *
     * @return the price (create, update, import)
     */
    public String getPrice() {
        return price;
    }

    /**
     * Sets the price (create, update, import).
     *
     * @param price the price (create, update, import)
     */
    public void setPrice(final String price) {
        this.price = price;
    }

    /**
     * Gets the decimals (create, import).
     *
     * @return the decimals (create, import)
     */
    public String getDecimals() {
        return decimals;
    }

    /**
     * Sets the decimals (create, import).
     *
     * @param decimals the decimals (create, import)
     */
    public void setDecimals(final String decimals) {
        this.decimals = decimals;
    }

    /**
     * Gets the blockchain (create, import).
     *
     * @return the blockchain (create, import)
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain (create, import).
     *
     * @param blockchain the blockchain (create, import)
     */
    public void setBlockchain(final String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the network (create, import).
     *
     * @return the network (create, import)
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network (create, import).
     *
     * @param network the network (create, import)
     */
    public void setNetwork(final String network) {
        this.network = network;
    }

    /**
     * Gets the asset type (create).
     *
     * @return the asset type (create)
     */
    public String getAssetType() {
        return assetType;
    }

    /**
     * Sets the asset type (create).
     *
     * @param assetType the asset type (create)
     */
    public void setAssetType(final String assetType) {
        this.assetType = assetType;
    }

    /**
     * Gets the imported contract address (import).
     *
     * @return the imported contract address (import)
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the imported contract address (import).
     *
     * @param address the imported contract address (import)
     */
    public void setAddress(final String address) {
        this.address = address;
    }

    /**
     * Gets the Canton instrument id (create).
     *
     * @return the Canton instrument id (create)
     */
    public String getInstrumentId() {
        return instrumentId;
    }

    /**
     * Sets the Canton instrument id (create).
     *
     * @param instrumentId the Canton instrument id (create)
     */
    public void setInstrumentId(final String instrumentId) {
        this.instrumentId = instrumentId;
    }

    /**
     * Gets the Canton token name (create).
     *
     * @return the Canton token name (create)
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the Canton token name (create).
     *
     * @param name the Canton token name (create)
     */
    public void setName(final String name) {
        this.name = name;
    }

    /**
     * Gets the Canton token symbol (create).
     *
     * @return the Canton token symbol (create)
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets the Canton token symbol (create).
     *
     * @param symbol the Canton token symbol (create)
     */
    public void setSymbol(final String symbol) {
        this.symbol = symbol;
    }

    /**
     * Gets the Canton credential requirement (update).
     *
     * @return the Canton credential requirement (update)
     */
    public Boolean getRequireCredentials() {
        return requireCredentials;
    }

    /**
     * Sets the Canton credential requirement (update).
     *
     * @param requireCredentials the Canton credential requirement (update)
     */
    public void setRequireCredentials(final Boolean requireCredentials) {
        this.requireCredentials = requireCredentials;
    }

    /**
     * Gets the amount (mint, burn).
     *
     * @return the amount (mint, burn)
     */
    public String getAmount() {
        return amount;
    }

    /**
     * Sets the amount (mint, burn).
     *
     * @param amount the amount (mint, burn)
     */
    public void setAmount(final String amount) {
        this.amount = amount;
    }

    /**
     * Gets the destination (mint, burn).
     *
     * @return the destination (mint, burn)
     */
    public AddressTargetV2 getDestination() {
        return destination;
    }

    /**
     * Sets the destination (mint, burn).
     *
     * @param destination the destination (mint, burn)
     */
    public void setDestination(final AddressTargetV2 destination) {
        this.destination = destination;
    }

    /**
     * Gets the NFT token ids (burn).
     *
     * @return the NFT token ids (burn)
     */
    public List<String> getNftTokenIds() {
        return nftTokenIds;
    }

    /**
     * Sets the NFT token ids (burn).
     *
     * @param nftTokenIds the NFT token ids (burn)
     */
    public void setNftTokenIds(final List<String> nftTokenIds) {
        this.nftTokenIds = nftTokenIds;
    }

    /**
     * Gets the NFT metadata (mint).
     *
     * @return the NFT metadata (mint)
     */
    public List<byte[]> getNftMetadata() {
        return nftMetadata;
    }

    /**
     * Sets the NFT metadata (mint).
     *
     * @param nftMetadata the NFT metadata (mint)
     */
    public void setNftMetadata(final List<byte[]> nftMetadata) {
        this.nftMetadata = nftMetadata;
    }

    /**
     * Gets the target address (account pause, account unpause, set KYC).
     *
     * @return the target address (account pause, account unpause, set KYC)
     */
    public AddressTargetV2 getTarget() {
        return target;
    }

    /**
     * Sets the target address (account pause, account unpause, set KYC).
     *
     * @param target the target address (account pause, account unpause, set KYC)
     */
    public void setTarget(final AddressTargetV2 target) {
        this.target = target;
    }

    /**
     * Gets the KYC status set, as the wire value (set KYC).
     *
     * @return the KYC status set, as the wire value (set KYC)
     */
    public String getKycStatus() {
        return kycStatus;
    }

    /**
     * Sets the KYC status set, as the wire value (set KYC).
     *
     * @param kycStatus the KYC status set, as the wire value (set KYC)
     */
    public void setKycStatus(final String kycStatus) {
        this.kycStatus = kycStatus;
    }
}
