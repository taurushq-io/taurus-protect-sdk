package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * An asset from the v2 asset service (a token or instrument Taurus-PROTECT manages).
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService#queryAssets
 */
public class AssetV2 {

    private String id;
    private String tenantId;
    private String version;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;
    private String label;
    private String assetType;
    private String status;
    private String blockchain;
    private String network;
    private String currencyId;
    private String name;
    private String symbol;
    private String decimals;
    private String contractAddress;
    private List<AssetAttributeV2> attributes;
    private CantonNativeTokenV2 cantonNativeToken;

    /**
     * Gets the asset id.
     *
     * @return the asset id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the asset id.
     *
     * @param id the asset id
     */
    public void setId(final String id) {
        this.id = id;
    }

    /**
     * Gets the tenant id.
     *
     * @return the tenant id
     */
    public String getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant id.
     *
     * @param tenantId the tenant id
     */
    public void setTenantId(final String tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Gets the record version.
     *
     * @return the record version
     */
    public String getVersion() {
        return version;
    }

    /**
     * Sets the record version.
     *
     * @param version the record version
     */
    public void setVersion(final String version) {
        this.version = version;
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
     * Gets the label.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the label.
     *
     * @param label the label
     */
    public void setLabel(final String label) {
        this.label = label;
    }

    /**
     * Gets the asset type.
     *
     * @return the asset type
     */
    public String getAssetType() {
        return assetType;
    }

    /**
     * Sets the asset type.
     *
     * @param assetType the asset type
     */
    public void setAssetType(final String assetType) {
        this.assetType = assetType;
    }

    /**
     * Gets the view status, as the wire value (e.g. ASSET_VIEW_STATUS_V2_ACTIVE).
     *
     * @return the view status, as the wire value (e.g. ASSET_VIEW_STATUS_V2_ACTIVE)
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the view status, as the wire value (e.g. ASSET_VIEW_STATUS_V2_ACTIVE).
     *
     * @param status the view status, as the wire value (e.g. ASSET_VIEW_STATUS_V2_ACTIVE)
     */
    public void setStatus(final String status) {
        this.status = status;
    }

    /**
     * Gets the blockchain.
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain.
     *
     * @param blockchain the blockchain
     */
    public void setBlockchain(final String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the network.
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network.
     *
     * @param network the network
     */
    public void setNetwork(final String network) {
        this.network = network;
    }

    /**
     * Gets the currency id.
     *
     * @return the currency id
     */
    public String getCurrencyId() {
        return currencyId;
    }

    /**
     * Sets the currency id.
     *
     * @param currencyId the currency id
     */
    public void setCurrencyId(final String currencyId) {
        this.currencyId = currencyId;
    }

    /**
     * Gets the name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the name.
     *
     * @param name the name
     */
    public void setName(final String name) {
        this.name = name;
    }

    /**
     * Gets the symbol.
     *
     * @return the symbol
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets the symbol.
     *
     * @param symbol the symbol
     */
    public void setSymbol(final String symbol) {
        this.symbol = symbol;
    }

    /**
     * Gets the decimals.
     *
     * @return the decimals
     */
    public String getDecimals() {
        return decimals;
    }

    /**
     * Sets the decimals.
     *
     * @param decimals the decimals
     */
    public void setDecimals(final String decimals) {
        this.decimals = decimals;
    }

    /**
     * Gets the contract address.
     *
     * @return the contract address
     */
    public String getContractAddress() {
        return contractAddress;
    }

    /**
     * Sets the contract address.
     *
     * @param contractAddress the contract address
     */
    public void setContractAddress(final String contractAddress) {
        this.contractAddress = contractAddress;
    }

    /**
     * Gets the attributes.
     *
     * @return the attributes
     */
    public List<AssetAttributeV2> getAttributes() {
        return attributes;
    }

    /**
     * Sets the attributes.
     *
     * @param attributes the attributes
     */
    public void setAttributes(final List<AssetAttributeV2> attributes) {
        this.attributes = attributes;
    }

    /**
     * Gets the Canton native-token details, null for other assets.
     *
     * @return the Canton native-token details, null for other assets
     */
    public CantonNativeTokenV2 getCantonNativeToken() {
        return cantonNativeToken;
    }

    /**
     * Sets the Canton native-token details, null for other assets.
     *
     * @param cantonNativeToken the Canton native-token details, null for other assets
     */
    public void setCantonNativeToken(final CantonNativeTokenV2 cantonNativeToken) {
        this.cantonNativeToken = cantonNativeToken;
    }
}
