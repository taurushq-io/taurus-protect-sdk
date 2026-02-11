package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.TokenMetadataMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.TokenMetadata;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.TokenMetadataApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetERCTokenMetadataReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFATokenMetadataReply;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving token metadata in the Taurus Protect system.
 * <p>
 * This service provides access to token metadata for various token standards
 * including ERC-20, ERC-721, ERC-1155 (Ethereum), and FA tokens (Tezos).
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get ERC-20 token metadata
 * TokenMetadata usdc = client.getTokenMetadataService().getERCTokenMetadata(
 *     "mainnet",
 *     "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",  // USDC contract
 *     null,  // token ID (for ERC-721/1155)
 *     false, // withData
 *     "ETH"
 * );
 * System.out.println("Token: " + usdc.getName() + ", Decimals: " + usdc.getDecimals());
 *
 * // Get FA token metadata (Tezos)
 * TokenMetadata tzBtc = client.getTokenMetadataService().getFATokenMetadata(
 *     "mainnet",
 *     "KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn", // tzBTC contract
 *     "0",  // token ID
 *     false // withData
 * );
 * }</pre>
 *
 * @see TokenMetadata
 */
public class TokenMetadataService {

    private final TokenMetadataApi tokenMetadataApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final TokenMetadataMapper tokenMetadataMapper;

    /**
     * Instantiates a new Token Metadata service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public TokenMetadataService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.tokenMetadataApi = new TokenMetadataApi(openApiClient);
        this.tokenMetadataMapper = TokenMetadataMapper.INSTANCE;
    }

    /**
     * Retrieves ERC token metadata (ERC-20, ERC-721, ERC-1155).
     *
     * @param network    the network (e.g., "mainnet", "goerli")
     * @param contract   the contract address
     * @param tokenId    the token ID (required for ERC-721/1155, null for ERC-20)
     * @param withData   whether to include base64 data (for NFTs)
     * @param blockchain the blockchain symbol (e.g., "ETH")
     * @return the token metadata
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or contract is null/empty
     */
    public TokenMetadata getERCTokenMetadata(final String network,
                                              final String contract,
                                              final String tokenId,
                                              final Boolean withData,
                                              final String blockchain) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(contract), "contract cannot be null or empty");

        try {
            TgvalidatordGetERCTokenMetadataReply reply = tokenMetadataApi.tokenMetadataServiceGetERCTokenMetadata(
                    network, contract, tokenId, withData, blockchain);
            return tokenMetadataMapper.fromERCDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves ERC token metadata for EVM-compatible chains.
     *
     * @param network    the network (e.g., "mainnet")
     * @param contract   the contract address
     * @param tokenId    the token ID (required for ERC-721/1155, null for ERC-20)
     * @param withData   whether to include base64 data (for NFTs)
     * @param blockchain the blockchain symbol (e.g., "MATIC", "AVAX")
     * @return the token metadata
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network, contract, or blockchain is null/empty
     */
    public TokenMetadata getEVMERCTokenMetadata(final String network,
                                                 final String contract,
                                                 final String tokenId,
                                                 final Boolean withData,
                                                 final String blockchain) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(contract), "contract cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(blockchain), "blockchain cannot be null or empty");

        try {
            TgvalidatordGetERCTokenMetadataReply reply = tokenMetadataApi.tokenMetadataServiceGetEVMERCTokenMetadata(
                    network, contract, tokenId, withData, blockchain);
            return tokenMetadataMapper.fromERCDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves FA token metadata (Tezos FA1.2/FA2 standards).
     *
     * @param network  the network (e.g., "mainnet")
     * @param contract the contract address
     * @param tokenId  the token ID
     * @param withData whether to include base64 data
     * @return the token metadata
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or contract is null/empty
     */
    public TokenMetadata getFATokenMetadata(final String network,
                                             final String contract,
                                             final String tokenId,
                                             final Boolean withData) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(contract), "contract cannot be null or empty");

        try {
            TgvalidatordGetFATokenMetadataReply reply = tokenMetadataApi.tokenMetadataServiceGetFATokenMetadata(
                    network, contract, tokenId, withData);
            return tokenMetadataMapper.fromFADTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
