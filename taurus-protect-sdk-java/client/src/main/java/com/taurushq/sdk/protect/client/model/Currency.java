package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a cryptocurrency or token in the Taurus Protect system.
 * <p>
 * Currency contains metadata about supported cryptocurrencies and tokens,
 * including information about the blockchain, token standards, and display
 * configuration.
 * <p>
 * Key currency types:
 * <ul>
 *   <li><b>Native currencies</b> - ETH, BTC, SOL, etc.</li>
 *   <li><b>ERC-20 tokens</b> - Ethereum-based fungible tokens</li>
 *   <li><b>FA1.2/FA2.0 tokens</b> - Tezos token standards</li>
 *   <li><b>NFTs</b> - Non-fungible tokens</li>
 * </ul>
 *
 * @see Wallet
 * @see Address
 * @see Transaction
 */
public class Currency {

    /**
     * The unique identifier of the currency.
     */
    private String id;

    /**
     * The full name of the currency (e.g., "Ethereum").
     */
    private String name;

    /**
     * The ticker symbol (e.g., "ETH", "BTC").
     */
    private String symbol;

    /**
     * The BIP-44 coin type index.
     */
    private String coinTypeIndex;

    /**
     * The blockchain identifier (e.g., "ETH", "BTC").
     */
    private String blockchain;

    /**
     * Indicates whether this is a token (vs. native currency).
     */
    private boolean isToken;

    /**
     * Indicates whether this is an ERC-20 token.
     */
    private boolean isERC20;

    /**
     * The number of decimal places for the currency.
     */
    private int decimals;

    /**
     * The smart contract address for tokens.
     */
    private String contractAddress;

    /**
     * Indicates whether staking is supported for this currency.
     */
    private boolean hasStaking;

    /**
     * Indicates whether this is a UTXO-based currency (e.g., Bitcoin).
     */
    private boolean isUTXOBased;

    /**
     * Indicates whether this is an account-based currency (e.g., Ethereum).
     */
    private boolean isAccountBased;

    /**
     * Indicates whether this is a fiat currency.
     */
    private boolean isFiat;

    /**
     * Indicates whether this is a Tezos FA1.2 token.
     */
    private boolean isFA12;

    /**
     * Indicates whether this is a Tezos FA2.0 token.
     */
    private boolean isFA20;

    /**
     * Indicates whether this is a non-fungible token (NFT).
     */
    private boolean isNFT;

    /**
     * Indicates whether this currency is enabled in the system.
     */
    private boolean enabled;

    /**
     * The display name for the currency.
     */
    private String displayName;

    /**
     * The currency type classification.
     */
    private String type;

    /**
     * The whitelisted contract address ID.
     */
    private long wlcaId;

    /**
     * The network identifier (e.g., "mainnet", "testnet").
     */
    private String network;

    /**
     * The token ID for NFTs or multi-token contracts.
     */
    private String tokenId;

    /**
     * The URL to the currency's logo image.
     */
    private String logo;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets id.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets id.
     *
     * @param id the id
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets name.
     *
     * @param name the name
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets symbol.
     *
     * @return the symbol
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets symbol.
     *
     * @param symbol the symbol
     */
    public void setSymbol(String symbol) {
        this.symbol = symbol;
    }

    /**
     * Gets coin type index.
     *
     * @return the coin type index
     */
    public String getCoinTypeIndex() {
        return coinTypeIndex;
    }

    /**
     * Sets coin type index.
     *
     * @param coinTypeIndex the coin type index
     */
    public void setCoinTypeIndex(String coinTypeIndex) {
        this.coinTypeIndex = coinTypeIndex;
    }

    /**
     * Gets blockchain.
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets blockchain.
     *
     * @param blockchain the blockchain
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Is token boolean.
     *
     * @return the boolean
     */
    public boolean isToken() {
        return isToken;
    }

    /**
     * Is token.
     *
     * @param token the token
     */
    public void setToken(boolean token) {
        isToken = token;
    }

    /**
     * Is erc 20 boolean.
     *
     * @return the boolean
     */
    public boolean isERC20() {
        return isERC20;
    }

    /**
     * Is erc 20.
     *
     * @param isERC20 the is erc 20
     */
    public void setERC20(boolean isERC20) {
        this.isERC20 = isERC20;
    }

    /**
     * Gets decimals.
     *
     * @return the decimals
     */
    public int getDecimals() {
        return decimals;
    }

    /**
     * Sets decimals.
     *
     * @param decimals the decimals
     */
    public void setDecimals(int decimals) {
        this.decimals = decimals;
    }

    /**
     * Gets contract address.
     *
     * @return the contract address
     */
    public String getContractAddress() {
        return contractAddress;
    }

    /**
     * Sets contract address.
     *
     * @param contractAddress the contract address
     */
    public void setContractAddress(String contractAddress) {
        this.contractAddress = contractAddress;
    }

    /**
     * Has staking boolean.
     *
     * @return the boolean
     */
    public boolean hasStaking() {
        return hasStaking;
    }

    /**
     * Sets has staking.
     *
     * @param hasStaking the has staking
     */
    public void setHasStaking(boolean hasStaking) {
        this.hasStaking = hasStaking;
    }

    /**
     * Is utxo based boolean.
     *
     * @return the boolean
     */
    public boolean isUTXOBased() {
        return isUTXOBased;
    }

    /**
     * Is utxo based.
     *
     * @param isUTXOBased the is utxo based
     */
    public void setUTXOBased(boolean isUTXOBased) {
        this.isUTXOBased = isUTXOBased;
    }

    /**
     * Is account based boolean.
     *
     * @return the boolean
     */
    public boolean isAccountBased() {
        return isAccountBased;
    }

    /**
     * Is account based.
     *
     * @param accountBased the account based
     */
    public void setAccountBased(boolean accountBased) {
        isAccountBased = accountBased;
    }

    /**
     * Is fiat boolean.
     *
     * @return the boolean
     */
    public boolean isFiat() {
        return isFiat;
    }

    /**
     * Is fiat.
     *
     * @param fiat the fiat
     */
    public void setFiat(boolean fiat) {
        isFiat = fiat;
    }

    /**
     * Is fa 12 boolean.
     *
     * @return the boolean
     */
    public boolean isFA12() {
        return isFA12;
    }

    /**
     * Is fa 12.
     *
     * @param isFA12 the is fa 12
     */
    public void setFA12(boolean isFA12) {
        this.isFA12 = isFA12;
    }

    /**
     * Is fa 20 boolean.
     *
     * @return the boolean
     */
    public boolean isFA20() {
        return isFA20;
    }

    /**
     * Is fa 20.
     *
     * @param isFA20 the is fa 20
     */
    public void setFA20(boolean isFA20) {
        this.isFA20 = isFA20;
    }

    /**
     * Is nft boolean.
     *
     * @return the boolean
     */
    public boolean isNFT() {
        return isNFT;
    }

    /**
     * Is nft.
     *
     * @param isNFT the is nft
     */
    public void setNFT(boolean isNFT) {
        this.isNFT = isNFT;
    }

    /**
     * Is enabled boolean.
     *
     * @return the boolean
     */
    public boolean isEnabled() {
        return enabled;
    }

    /**
     * Is enabled.
     *
     * @param enabled the enabled
     */
    public void setEnabled(boolean enabled) {
        this.enabled = enabled;
    }

    /**
     * Gets display name.
     *
     * @return the display name
     */
    public String getDisplayName() {
        return displayName;
    }

    /**
     * Sets display name.
     *
     * @param displayName the display name
     */
    public void setDisplayName(String displayName) {
        this.displayName = displayName;
    }

    /**
     * Gets type.
     *
     * @return the type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets type.
     *
     * @param type the type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets wlca id.
     *
     * @return the wlca id
     */
    public long getWlcaId() {
        return wlcaId;
    }

    /**
     * Sets wlca id.
     *
     * @param wlcaId the wlca id
     */
    public void setWlcaId(long wlcaId) {
        this.wlcaId = wlcaId;
    }

    /**
     * Gets network.
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets network.
     *
     * @param network the network
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets token id.
     *
     * @return the token id
     */
    public String getTokenId() {
        return tokenId;
    }

    /**
     * Sets token id.
     *
     * @param tokenId the token id
     */
    public void setTokenId(String tokenId) {
        this.tokenId = tokenId;
    }

    /**
     * Gets logo.
     *
     * @return the logo
     */
    public String getLogo() {
        return logo;
    }

    /**
     * Sets logo.
     *
     * @param logo the logo
     */
    public void setLogo(String logo) {
        this.logo = logo;
    }
}
