package com.taurushq.sdk.protect.client.model;

/**
 * Represents a whitelisted contract address (e.g., ERC20 token, NFT collection).
 * <p>
 * Whitelisted contract addresses allow transactions involving specific tokens or NFTs
 * to be processed according to governance rules. This class stores the contract metadata
 * including blockchain, network, token details, and contract identifiers.
 * <p>
 * Common use cases include:
 * <ul>
 *   <li>ERC20 tokens on EVM chains</li>
 *   <li>FA2 tokens on Tezos</li>
 *   <li>NFT collections (ERC721, ERC1155)</li>
 *   <li>SPL tokens on Solana</li>
 * </ul>
 *
 * @see WhitelistedAsset
 * @see SignedWhitelistedAsset
 */
public class WhitelistedContractAddress {

    /**
     * The blockchain this contract is deployed on (e.g., "ETH", "MATIC", "XTZ").
     */
    private String blockchain;

    /**
     * The network identifier (e.g., "mainnet", "goerli", "mumbai").
     */
    private String network;

    /**
     * Human-readable name of the token or NFT collection.
     */
    private String name;

    /**
     * The token symbol (e.g., "USDC", "WETH", "BAYC").
     */
    private String symbol;

    /**
     * Number of decimal places for the token. For NFTs, this is typically 0.
     */
    private int decimals;

    /**
     * The smart contract address for this token or NFT collection.
     */
    private String contractAddress;

    /**
     * The specific token ID for NFTs. For fungible tokens, this is typically null.
     */
    private String tokenId;

    /**
     * The kind/type of contract (e.g., "erc20", "erc721", "fa2").
     */
    private String kind;

    /**
     * Gets the blockchain identifier.
     *
     * @return the blockchain (e.g., "ETH", "MATIC", "XTZ")
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain identifier.
     *
     * @param blockchain the blockchain to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Gets the network identifier.
     *
     * @return the network (e.g., "mainnet", "goerli")
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network identifier.
     *
     * @param network the network to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets the human-readable name.
     *
     * @return the token or collection name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the human-readable name.
     *
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the token symbol.
     *
     * @return the symbol (e.g., "USDC", "WETH")
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
     * Gets the number of decimal places.
     *
     * @return the decimals (0 for NFTs, typically 6-18 for fungible tokens)
     */
    public int getDecimals() {
        return decimals;
    }

    /**
     * Sets the number of decimal places.
     *
     * @param decimals the decimals to set
     */
    public void setDecimals(int decimals) {
        this.decimals = decimals;
    }

    /**
     * Gets the smart contract address.
     *
     * @return the contract address
     */
    public String getContractAddress() {
        return contractAddress;
    }

    /**
     * Sets the smart contract address.
     *
     * @param contractAddress the contract address to set
     */
    public void setContractAddress(String contractAddress) {
        this.contractAddress = contractAddress;
    }

    /**
     * Gets the token ID for NFTs.
     *
     * @return the token ID, or null for fungible tokens
     */
    public String getTokenId() {
        return tokenId;
    }

    /**
     * Sets the token ID.
     *
     * @param tokenId the token ID to set
     */
    public void setTokenId(String tokenId) {
        this.tokenId = tokenId;
    }

    /**
     * Gets the contract kind/type.
     *
     * @return the kind (e.g., "erc20", "erc721", "fa2")
     */
    public String getKind() {
        return kind;
    }

    /**
     * Sets the contract kind/type.
     *
     * @param kind the kind to set
     */
    public void setKind(String kind) {
        this.kind = kind;
    }
}
