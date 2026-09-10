package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ContractWhitelistingMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ContractWhitelistingApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveWhitelistedContractAddressRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWhitelistedContractAddressAttributeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateWhitelistedContractAddressRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordDeleteWhitelistedContractAddressRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWhitelistedContractAddressAttributeReply;
import com.taurushq.sdk.protect.openapi.model.WhitelistServiceCreateWhitelistedContractAttributesBody;
import com.taurushq.sdk.protect.openapi.model.WhitelistServiceUpdateWhitelistedContractBody;

import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for whitelisted contract WRITE operations (tokens, NFTs) in Taurus Protect.
 * <p>
 * Creating, approving, updating and deleting whitelisted contract addresses such as ERC20
 * tokens, NFT collections (ERC721/ERC1155), FA2 tokens on Tezos, and other smart
 * contract-based assets.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a new whitelisted contract
 * String id = client.getContractWhitelistingService().createWhitelistedContract(
 *     "ETH", "mainnet", "0x1234...", "USDC", "USD Coin", 6, "erc20", null);
 *
 * // Read them back through the verified reader
 * WhitelistedAssetResult result = client.getWhitelistedAssetService()
 *     .getWhitelistedAssets(50, 0, "ETH", "mainnet", null, null, null, null);
 * }</pre>
 *
 * @see WhitelistedAssetService
 */
// Reads live on WhitelistedAssetService, which is the same server entity
// (/api/rest/v1/whitelists/contracts) verified through the six-step chain. The
// get/getAll/withFilters/forApproval that used to sit here returned the envelope with no
// verification, which made the verified reader avoidable. Do not re-add them.
public class ContractWhitelistingService {

    /**
     * The underlying OpenAPI client for contract whitelisting operations.
     */
    private final ContractWhitelistingApi whitelistApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Mapper for converting DTOs to domain models.
     */
    private final ContractWhitelistingMapper mapper;

    /**
     * Instantiates a new Contract whitelisting service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public ContractWhitelistingService(final ApiClient openApiClient,
                                        final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.whitelistApi = new ContractWhitelistingApi(openApiClient);
        this.mapper = ContractWhitelistingMapper.INSTANCE;
    }

    /**
     * Creates a new whitelisted contract address.
     * <p>
     * The contract will be created in a pending state and require approval
     * according to the configured governance rules.
     *
     * @param blockchain      the blockchain identifier (e.g., "ETH", "MATIC", "XTZ")
     * @param network         the network identifier (e.g., "mainnet", "goerli")
     * @param contractAddress the smart contract address
     * @param symbol          the token symbol (e.g., "USDC", "WETH")
     * @param name            the human-readable name
     * @param decimals        the number of decimals (0 for NFTs)
     * @param kind            the contract kind (e.g., "erc20", "erc721", "fa2")
     * @param tokenId         the token ID for NFTs (null for fungible tokens)
     * @return the ID of the created whitelist entry
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if required parameters are null or empty
     */
    public String createWhitelistedContract(final String blockchain, final String network,
                                            final String contractAddress, final String symbol,
                                            final String name, final int decimals,
                                            final String kind, final String tokenId)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(blockchain), "blockchain cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(symbol), "symbol cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(name), "name cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(kind), "kind cannot be null or empty");

        try {
            TgvalidatordCreateWhitelistedContractAddressRequest request =
                    new TgvalidatordCreateWhitelistedContractAddressRequest();
            request.setBlockchain(blockchain);
            request.setNetwork(network);
            request.setContractAddress(contractAddress);
            request.setSymbol(symbol);
            request.setName(name);
            request.setDecimals(String.valueOf(decimals));
            request.setKind(kind);
            request.setTokenId(tokenId);

            return whitelistApi.whitelistServiceCreateWhitelistedContract(request).getResult().getId();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Approves one or more whitelisted contract addresses.
     * <p>
     * Requires a cryptographic signature computed over the metadata hashes
     * of the contracts being approved.
     *
     * @param ids       the list of contract IDs to approve
     * @param signature the approval signature (base64-encoded)
     * @param comment   the approval comment
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if parameters are invalid
     * @deprecated the signature is an opaque blob over hashes nothing verified, so the
     *     caller cannot know what they signed. Use
     *     {@link WhitelistedAssetService#approveWhitelistedAssets}, which re-reads and
     *     verifies the rows first.
     */
    @Deprecated
    public void approveWhitelistedContracts(final List<String> ids, final String signature,
                                            final String comment)
            throws ApiException {
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");
        checkArgument(!Strings.isNullOrEmpty(signature), "signature cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(comment), "comment cannot be null or empty");

        try {
            TgvalidatordApproveWhitelistedContractAddressRequest request =
                    new TgvalidatordApproveWhitelistedContractAddressRequest();
            request.setIds(ids);
            request.setSignature(signature);
            request.setComment(comment);

            whitelistApi.whitelistServiceApproveWhitelistedContract(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    
    
    
    
    /**
     * Updates an existing whitelisted contract.
     * <p>
     * Only the symbol, name, and decimals can be updated.
     *
     * @param id       the contract ID to update
     * @param symbol   the new symbol
     * @param name     the new name
     * @param decimals the new decimals value
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if required parameters are invalid
     */
    public void updateWhitelistedContract(final String id, final String symbol,
                                          final String name, final int decimals)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(symbol), "symbol cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(name), "name cannot be null or empty");

        try {
            WhitelistServiceUpdateWhitelistedContractBody body =
                    new WhitelistServiceUpdateWhitelistedContractBody();
            body.setSymbol(symbol);
            body.setName(name);
            body.setDecimals(String.valueOf(decimals));

            whitelistApi.whitelistServiceUpdateWhitelistedContract(id, body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Deletes a whitelisted contract.
     * <p>
     * Deletion may require approval depending on governance rules.
     *
     * @param id      the contract ID to delete
     * @param comment the reason for deletion (optional)
     * @return the ID of the deletion request
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public String deleteWhitelistedContract(final String id, final String comment)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordDeleteWhitelistedContractAddressRequest request =
                    new TgvalidatordDeleteWhitelistedContractAddressRequest();
            request.setId(id);
            request.setComment(comment);

            return whitelistApi.whitelistServiceDeleteWhitelistedContract(request).getResult().getId();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Creates an attribute on a whitelisted contract.
     *
     * @param contractId  the contract ID
     * @param key         the attribute key
     * @param value       the attribute value
     * @param contentType the content type (optional)
     * @param type        the attribute type (optional)
     * @param subType     the attribute subtype (optional)
     * @return the list of created attributes
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if required parameters are invalid
     */
    public List<Attribute> createAttribute(final String contractId, final String key,
                                           final String value, final String contentType,
                                           final String type, final String subType)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(contractId), "contractId cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(key), "key cannot be null or empty");

        try {
            TgvalidatordCreateWhitelistedContractAddressAttributeRequest attrRequest =
                    new TgvalidatordCreateWhitelistedContractAddressAttributeRequest();
            attrRequest.setKey(key);
            attrRequest.setValue(value);
            attrRequest.setContentType(contentType);
            attrRequest.setType(type);
            attrRequest.setSubtype(subType);

            WhitelistServiceCreateWhitelistedContractAttributesBody body =
                    new WhitelistServiceCreateWhitelistedContractAttributesBody();
            body.setAttributes(Collections.singletonList(attrRequest));

            return mapper.fromAttributeDTOList(
                    whitelistApi.whitelistServiceCreateWhitelistedContractAttributes(contractId, body)
                            .getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets an attribute from a whitelisted contract.
     *
     * @param contractId  the contract ID
     * @param attributeId the attribute ID
     * @return the attribute
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if parameters are invalid
     */
    public Attribute getAttribute(final String contractId, final String attributeId)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(contractId), "contractId cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(attributeId), "attributeId cannot be null or empty");

        try {
            TgvalidatordGetWhitelistedContractAddressAttributeReply reply =
                    whitelistApi.whitelistServiceGetWhitelistedContractAttribute(
                            contractId, attributeId);
            return mapper.fromAttributeDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
