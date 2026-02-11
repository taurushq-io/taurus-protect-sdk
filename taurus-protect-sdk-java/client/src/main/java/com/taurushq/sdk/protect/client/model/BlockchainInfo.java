package com.taurushq.sdk.protect.client.model;

/**
 * Represents blockchain network information in the Taurus Protect system.
 * <p>
 * Contains details about a supported blockchain including its symbol, name,
 * network type, chain ID, and various blockchain-specific configuration.
 *
 * @see BlockchainService
 */
public class BlockchainInfo {

    private String symbol;
    private String name;
    private String network;
    private String chainId;
    private String confirmations;
    private String blockHeight;
    private String blackholeAddress;
    private Boolean isLayer2Chain;
    private String layer1Network;
    private Currency baseCurrency;

    /**
     * Gets the blockchain symbol (e.g., "BTC", "ETH", "SOL").
     *
     * @return the symbol
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets the blockchain symbol.
     *
     * @param symbol the symbol to set
     */
    public void setSymbol(String symbol) {
        this.symbol = symbol;
    }

    /**
     * Gets the blockchain name (e.g., "Bitcoin", "Ethereum").
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the blockchain name.
     *
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the network type (e.g., "mainnet", "testnet").
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the network type.
     *
     * @param network the network to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Gets the chain ID (for EVM-compatible chains).
     *
     * @return the chain ID
     */
    public String getChainId() {
        return chainId;
    }

    /**
     * Sets the chain ID.
     *
     * @param chainId the chain ID to set
     */
    public void setChainId(String chainId) {
        this.chainId = chainId;
    }

    /**
     * Gets the number of confirmations required for transactions.
     *
     * @return the confirmations
     */
    public String getConfirmations() {
        return confirmations;
    }

    /**
     * Sets the confirmations.
     *
     * @param confirmations the confirmations to set
     */
    public void setConfirmations(String confirmations) {
        this.confirmations = confirmations;
    }

    /**
     * Gets the current block height of the blockchain.
     *
     * @return the block height
     */
    public String getBlockHeight() {
        return blockHeight;
    }

    /**
     * Sets the block height.
     *
     * @param blockHeight the block height to set
     */
    public void setBlockHeight(String blockHeight) {
        this.blockHeight = blockHeight;
    }

    /**
     * Gets the blackhole/burn address for this blockchain.
     *
     * @return the blackhole address
     */
    public String getBlackholeAddress() {
        return blackholeAddress;
    }

    /**
     * Sets the blackhole address.
     *
     * @param blackholeAddress the blackhole address to set
     */
    public void setBlackholeAddress(String blackholeAddress) {
        this.blackholeAddress = blackholeAddress;
    }

    /**
     * Checks if this is a Layer 2 chain.
     *
     * @return true if Layer 2 chain
     */
    public Boolean getIsLayer2Chain() {
        return isLayer2Chain;
    }

    /**
     * Sets whether this is a Layer 2 chain.
     *
     * @param isLayer2Chain the Layer 2 flag
     */
    public void setIsLayer2Chain(Boolean isLayer2Chain) {
        this.isLayer2Chain = isLayer2Chain;
    }

    /**
     * Gets the Layer 1 network for Layer 2 chains.
     *
     * @return the Layer 1 network
     */
    public String getLayer1Network() {
        return layer1Network;
    }

    /**
     * Sets the Layer 1 network.
     *
     * @param layer1Network the Layer 1 network to set
     */
    public void setLayer1Network(String layer1Network) {
        this.layer1Network = layer1Network;
    }

    /**
     * Gets the base currency information for this blockchain.
     *
     * @return the base currency
     */
    public Currency getBaseCurrency() {
        return baseCurrency;
    }

    /**
     * Sets the base currency.
     *
     * @param baseCurrency the base currency to set
     */
    public void setBaseCurrency(Currency baseCurrency) {
        this.baseCurrency = baseCurrency;
    }
}
