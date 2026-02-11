package com.taurushq.sdk.protect.client.model;

/**
 * Represents token metadata for ERC-20/ERC-721/ERC-1155 or FA tokens.
 * <p>
 * Contains details about a token including name, description, decimals,
 * and optional media data for NFTs.
 *
 * @see TokenMetadataService
 */
public class TokenMetadata {

    private String name;
    private String symbol;
    private String description;
    private String decimals;
    private String dataType;
    private String base64Data;
    private String uri;

    /**
     * Gets the token name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the token name.
     *
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the token symbol (for FA tokens).
     *
     * @return the symbol
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets the token symbol.
     *
     * @param symbol the symbol to set
     */
    public void setSymbol(String symbol) {
        this.symbol = symbol;
    }

    /**
     * Gets the token description.
     *
     * @return the description
     */
    public String getDescription() {
        return description;
    }

    /**
     * Sets the token description.
     *
     * @param description the description to set
     */
    public void setDescription(String description) {
        this.description = description;
    }

    /**
     * Gets the number of decimals for the token.
     *
     * @return the decimals
     */
    public String getDecimals() {
        return decimals;
    }

    /**
     * Sets the decimals.
     *
     * @param decimals the decimals to set
     */
    public void setDecimals(String decimals) {
        this.decimals = decimals;
    }

    /**
     * Gets the data type (for NFT media).
     *
     * @return the data type (e.g., "image/png")
     */
    public String getDataType() {
        return dataType;
    }

    /**
     * Sets the data type.
     *
     * @param dataType the data type to set
     */
    public void setDataType(String dataType) {
        this.dataType = dataType;
    }

    /**
     * Gets the base64 encoded data (for NFT media).
     *
     * @return the base64 data
     */
    public String getBase64Data() {
        return base64Data;
    }

    /**
     * Sets the base64 data.
     *
     * @param base64Data the base64 data to set
     */
    public void setBase64Data(String base64Data) {
        this.base64Data = base64Data;
    }

    /**
     * Gets the token URI (metadata URL).
     *
     * @return the URI
     */
    public String getUri() {
        return uri;
    }

    /**
     * Sets the URI.
     *
     * @param uri the URI to set
     */
    public void setUri(String uri) {
        this.uri = uri;
    }
}
