package com.taurushq.sdk.protect.client.model;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Strings.isNullOrEmpty;

/**
 * Request object for creating a new wallet.
 * <p>
 * Use the builder pattern for a cleaner API:
 * <pre>{@code
 * CreateWalletRequest request = CreateWalletRequest.builder()
 *     .blockchain("ETH")
 *     .network("mainnet")
 *     .name("Trading Wallet")
 *     .omnibus(true)
 *     .comment("Primary trading account")
 *     .customerId("customer-123")
 *     .build();
 *
 * Wallet wallet = client.getWalletService().createWallet(request);
 * }</pre>
 */
public final class CreateWalletRequest {

    /**
     * The blockchain type (e.g., "ETH", "BTC", "SOL").
     */
    private final String blockchain;

    /**
     * The network identifier (e.g., "mainnet", "testnet").
     */
    private final String network;

    /**
     * The human-readable name for the wallet.
     */
    private final String name;

    /**
     * Indicates whether this is an omnibus wallet.
     */
    private final boolean omnibus;

    /**
     * An optional comment or description for the wallet.
     */
    private final String comment;

    /**
     * An optional customer identifier associated with the wallet.
     */
    private final String customerId;

    private CreateWalletRequest(Builder builder) {
        this.blockchain = builder.blockchain;
        this.network = builder.network;
        this.name = builder.name;
        this.omnibus = builder.omnibus;
        this.comment = builder.comment;
        this.customerId = builder.customerId;
    }

    /**
     * Creates a new builder for CreateWalletRequest.
     *
     * @return a new builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Gets the blockchain identifier (e.g., "ETH", "BTC", "SOL").
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Gets the network identifier (e.g., "mainnet", "testnet").
     *
     * @return the network
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Gets the wallet name.
     *
     * @return the wallet name
     */
    public String getName() {
        return name;
    }

    /**
     * Returns whether this is an omnibus wallet.
     *
     * @return true if omnibus wallet
     */
    public boolean isOmnibus() {
        return omnibus;
    }

    /**
     * Gets the optional comment.
     *
     * @return the comment, or empty string if not set
     */
    public String getComment() {
        return comment;
    }

    /**
     * Gets the optional customer ID.
     *
     * @return the customer ID, or empty string if not set
     */
    public String getCustomerId() {
        return customerId;
    }

    /**
     * Builder for constructing {@link CreateWalletRequest} instances.
     * <p>
     * Required fields:
     * <ul>
     *   <li>{@link #blockchain(String)} - the blockchain type</li>
     *   <li>{@link #network(String)} - the network identifier</li>
     *   <li>{@link #name(String)} - the wallet name</li>
     * </ul>
     */
    public static final class Builder {

        /**
         * The blockchain type.
         */
        private String blockchain;

        /**
         * The network identifier.
         */
        private String network;

        /**
         * The wallet name.
         */
        private String name;

        /**
         * Whether this is an omnibus wallet.
         */
        private boolean omnibus;

        /**
         * The optional comment.
         */
        private String comment = "";

        /**
         * The optional customer ID.
         */
        private String customerId = "";

        private Builder() {
        }

        /**
         * Sets the blockchain identifier (required).
         *
         * @param blockchain the blockchain (e.g., "ETH", "BTC", "SOL")
         * @return this builder
         */
        public Builder blockchain(String blockchain) {
            this.blockchain = blockchain;
            return this;
        }

        /**
         * Sets the network identifier (required).
         *
         * @param network the network (e.g., "mainnet", "testnet")
         * @return this builder
         */
        public Builder network(String network) {
            this.network = network;
            return this;
        }

        /**
         * Sets the wallet name (required).
         *
         * @param name the wallet name
         * @return this builder
         */
        public Builder name(String name) {
            this.name = name;
            return this;
        }

        /**
         * Sets whether this is an omnibus wallet.
         * <p>
         * An omnibus wallet pools funds from multiple customers.
         *
         * @param omnibus true for omnibus wallet
         * @return this builder
         */
        public Builder omnibus(boolean omnibus) {
            this.omnibus = omnibus;
            return this;
        }

        /**
         * Sets an optional comment for the wallet.
         *
         * @param comment the comment
         * @return this builder
         */
        public Builder comment(String comment) {
            this.comment = comment != null ? comment : "";
            return this;
        }

        /**
         * Sets an optional customer ID for the wallet.
         *
         * @param customerId the customer ID
         * @return this builder
         */
        public Builder customerId(String customerId) {
            this.customerId = customerId != null ? customerId : "";
            return this;
        }

        /**
         * Builds the CreateWalletRequest.
         *
         * @return the request object
         * @throws IllegalArgumentException if required fields are missing
         */
        public CreateWalletRequest build() {
            checkArgument(!isNullOrEmpty(blockchain), "blockchain is required");
            checkArgument(!isNullOrEmpty(network), "network is required");
            checkArgument(!isNullOrEmpty(name), "name is required");
            return new CreateWalletRequest(this);
        }
    }
}
