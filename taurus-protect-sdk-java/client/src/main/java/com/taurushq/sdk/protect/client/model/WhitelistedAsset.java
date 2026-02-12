package com.taurushq.sdk.protect.client.model;

import java.util.HashMap;
import java.util.Map;

/**
 * Represents a whitelisted contract address (asset/token).
 * <p>
 * This is the domain model for whitelisted tokens, NFTs, and other contract-based assets.
 * Whitelisting a contract address allows transactions involving that token or NFT
 * to be processed according to the organization's governance rules.
 * <p>
 * Different asset kinds are supported across various blockchains:
 * <ul>
 *   <li>EVM chains: ERC721, ERC1155, CryptoPunks NFTs</li>
 *   <li>Tezos: FA2 tokens and NFTs</li>
 *   <li>Solana: SPL tokens (Token and Token2022 standards)</li>
 *   <li>Hedera: Native tokens and NFTs</li>
 * </ul>
 *
 * @see SignedWhitelistedAsset
 * @see SignedWhitelistedAssetEnvelope
 * @see AssetKind
 */
public class WhitelistedAsset {

    /**
     * The blockchain this asset contract is deployed on (e.g., "ETH", "MATIC", "XTZ").
     */
    private String blockchain;

    /**
     * The token symbol (e.g., "USDC", "BAYC").
     */
    private String symbol;

    /**
     * The smart contract address for this token or NFT collection.
     */
    private String contractAddress;

    /**
     * Human-readable name of the asset (e.g., "USD Coin", "Bored Ape Yacht Club").
     */
    private String name;

    /**
     * Number of decimal places for the token (e.g., 18 for most ERC20 tokens, 6 for USDC).
     * For NFTs, this is typically 0.
     */
    private int decimals;

    /**
     * The specific token ID for NFTs. For fungible tokens, this is typically null or empty.
     */
    private String tokenId;

    /**
     * The kind/type of this asset, determining how it should be handled.
     */
    private AssetKind kind;

    /**
     * The network identifier (e.g., "mainnet", "goerli", "mumbai").
     */
    private String network;

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
     * Gets the token symbol.
     *
     * @return the symbol (e.g., "USDC", "BAYC")
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
     * Gets the human-readable name.
     *
     * @return the name
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
     * Gets the token ID for NFTs.
     *
     * @return the token ID, or null/empty for fungible tokens
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
     * Gets the asset kind.
     *
     * @return the kind
     */
    public AssetKind getKind() {
        return kind;
    }

    /**
     * Sets the asset kind.
     *
     * @param kind the kind to set
     */
    public void setKind(AssetKind kind) {
        this.kind = kind;
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
     * Enum representing the kind/type of whitelisted contract address.
     * <p>
     * Different blockchains use different token standards, and this enum
     * identifies which standard applies to a whitelisted asset. This is used
     * to properly decode and validate token transfers.
     * <p>
     * Maps to protobuf WhitelistedContractAddressKind.
     */
    public enum AssetKind {
        /**
         * Default (empty string representation).
         */
        DEFAULT(0),
        /**
         * Generic NFT - backend will find best match.
         */
        NFT_AUTO(1),
        /**
         * Tezos FA2 NFT.
         */
        NFT_XTZ_FA2(2),
        /**
         * EVM ERC721 NFT.
         */
        NFT_EVM_ERC721(5),
        /**
         * EVM ERC1155 NFT.
         */
        NFT_EVM_ERC1155(6),
        /**
         * Solana Token.
         */
        SOL_TOKEN(7),
        /**
         * Solana Token 2022.
         */
        SOL_TOKEN2022(8),
        /**
         * EVM CryptoPunks NFT.
         */
        NFT_EVM_CRYPTOPUNKS(9),
        /**
         * Hedera native token.
         */
        HEDERA_NATIVE_TOKEN(10),
        /**
         * Hedera native NFT.
         */
        HEDERA_NATIVE_NFT(11),
        /**
         * Unrecognized kind.
         */
        UNRECOGNIZED(-1);

        private static final Map<Integer, AssetKind> MAP = new HashMap<>();

        static {
            for (AssetKind kind : values()) {
                MAP.put(kind.value, kind);
            }
        }

        /**
         * The integer value corresponding to this kind in the protobuf definition.
         */
        private final int value;

        /**
         * Constructs an AssetKind with the specified integer value.
         *
         * @param value the integer value from the protobuf definition
         */
        AssetKind(int value) {
            this.value = value;
        }

        /**
         * Gets the AssetKind from its integer value.
         *
         * @param value the integer value
         * @return the corresponding AssetKind, or UNRECOGNIZED if not found
         */
        public static AssetKind valueOf(int value) {
            AssetKind kind = MAP.get(value);
            return kind != null ? kind : UNRECOGNIZED;
        }

        /**
         * Gets the AssetKind from its string representation.
         *
         * @param kindStr the string representation (e.g., "NFT_AUTO", "SOL_TOKEN")
         * @return the corresponding AssetKind, or DEFAULT if null/empty, or UNRECOGNIZED if not found
         */
        public static AssetKind fromString(String kindStr) {
            if (kindStr == null || kindStr.isEmpty()) {
                return DEFAULT;
            }
            try {
                return AssetKind.valueOf(kindStr);
            } catch (IllegalArgumentException e) {
                return UNRECOGNIZED;
            }
        }

        /**
         * Gets the integer value for this kind.
         *
         * @return the integer value
         */
        public int getValue() {
            return value;
        }
    }
}
