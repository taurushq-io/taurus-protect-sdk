package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.BlockchainMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.BlockchainInfo;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.BlockchainApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetBlockchainsReply;

import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving blockchain information in the Taurus Protect system.
 * <p>
 * This service provides access to supported blockchain networks and their
 * configuration, including chain IDs, block heights, and other metadata.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all supported blockchains
 * List<BlockchainInfo> blockchains = client.getBlockchainService().getBlockchains();
 * for (BlockchainInfo bc : blockchains) {
 *     System.out.println(bc.getSymbol() + " (" + bc.getNetwork() + "): " + bc.getName());
 * }
 *
 * // Get blockchains with block height info
 * List<BlockchainInfo> withHeight = client.getBlockchainService()
 *     .getBlockchains("ETH", "mainnet", true);
 * }</pre>
 *
 * @see BlockchainInfo
 */
public class BlockchainService {

    private final BlockchainApi blockchainApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final BlockchainMapper blockchainMapper;

    /**
     * Instantiates a new Blockchain service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public BlockchainService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.blockchainApi = new BlockchainApi(openApiClient);
        this.blockchainMapper = BlockchainMapper.INSTANCE;
    }

    /**
     * Retrieves all supported blockchains.
     *
     * @return the list of blockchain information
     * @throws ApiException if the API call fails
     */
    public List<BlockchainInfo> getBlockchains() throws ApiException {
        return getBlockchains(null, null, null);
    }

    /**
     * Retrieves blockchains with optional filtering.
     *
     * @param blockchain         filter by blockchain symbol (optional)
     * @param network            filter by network type (optional)
     * @param includeBlockHeight whether to include current block height (optional)
     * @return the list of blockchain information
     * @throws ApiException if the API call fails
     */
    public List<BlockchainInfo> getBlockchains(final String blockchain,
                                                final String network,
                                                final Boolean includeBlockHeight) throws ApiException {
        try {
            TgvalidatordGetBlockchainsReply reply = blockchainApi.blockchainServiceGetBlockchains(
                    blockchain, network, includeBlockHeight);
            return blockchainMapper.fromDTOList(reply.getBlockchains());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
