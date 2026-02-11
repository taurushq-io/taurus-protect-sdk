package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.helper.AddressSignatureVerifier;
import com.taurushq.sdk.protect.client.mapper.AddressMapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.CreateAddressRequest;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AddressesApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateAddressAttributeRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateAddressReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateAddressRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAddressReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAddressesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAddressProofOfReserveReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordProofOfReserve;
import com.taurushq.sdk.protect.openapi.model.WalletServiceCreateAddressAttributesBody;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing blockchain addresses in the Taurus Protect system.
 * <p>
 * This service provides operations for creating, retrieving, and managing addresses
 * within wallets. Addresses are the on-chain identifiers used for receiving and
 * sending cryptocurrency.
 * <p>
 * All addresses retrieved through this service are automatically verified for
 * cryptographic integrity using the rules container public keys.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a new address in a wallet
 * CreateAddressRequest request = CreateAddressRequest.builder()
 *     .walletId(walletId)
 *     .label("Customer Deposit")
 *     .comment("Primary deposit address")
 *     .build();
 * Address address = client.getAddressService().createAddress(request);
 *
 * // Retrieve an address with signature verification
 * Address address = client.getAddressService().getAddress(addressId);
 *
 * // List addresses for a wallet
 * List<Address> addresses = client.getAddressService().getAddresses(walletId, 50, 0);
 * }</pre>
 *
 * @see Address
 * @see CreateAddressRequest
 * @see WalletService
 */
public class AddressService {

    /**
     * The underlying OpenAPI client for address operations.
     */
    private final AddressesApi addressesApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Cache for rules container used in address signature verification.
     */
    private final RulesContainerCache rulesContainerCache;

    /**
     * Instantiates a new Address service.
     *
     * @param openApiClient       the open api client
     * @param apiExceptionMapper  the api exception mapper
     * @param rulesContainerCache the rules container cache for signature verification
     */
    public AddressService(final ApiClient openApiClient,
                          final ApiExceptionMapper apiExceptionMapper,
                          final RulesContainerCache rulesContainerCache) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");
        checkNotNull(rulesContainerCache, "rulesContainerCache cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.addressesApi = new AddressesApi(openApiClient);
        this.rulesContainerCache = rulesContainerCache;
    }


    /**
     * Creates an address using a request object.
     * <p>
     * This is the recommended method for creating addresses as it provides
     * a cleaner API through the builder pattern:
     * <pre>{@code
     * CreateAddressRequest request = CreateAddressRequest.builder()
     *     .walletId(123)
     *     .label("Deposit Address")
     *     .comment("Customer deposit")
     *     .customerId("customer-456")
     *     .build();
     *
     * Address address = client.getAddressService().createAddress(request);
     * }</pre>
     *
     * @param request the address creation request
     * @return the created address
     * @throws ApiException the api exception
     */
    public Address createAddress(final CreateAddressRequest request) throws ApiException {
        checkNotNull(request, "request cannot be null");
        return createAddress(
                request.getWalletId(),
                request.getLabel(),
                request.getComment(),
                request.getCustomerId()
        );
    }

    /**
     * Create address.
     *
     * @param walletId   the wallet id
     * @param label      the label
     * @param comment    the comment
     * @param customerId the customer id
     * @return the address
     * @throws ApiException the api exception
     */
    public Address createAddress(final long walletId, final String label, final String comment, final String customerId) throws ApiException {

        checkArgument(walletId > 0, "walletId cannot be zero");
        checkArgument(!Strings.isNullOrEmpty(label), "label cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(comment), "comment cannot be null or empty");


        TgvalidatordCreateAddressRequest request = new TgvalidatordCreateAddressRequest();
        request.setWalletId(String.valueOf(walletId));
        request.setLabel(label);
        request.setComment(comment);
        request.setCustomerId(customerId);
        try {
            TgvalidatordCreateAddressReply reply = addressesApi.walletServiceCreateAddress(request);
            return AddressMapper.INSTANCE.fromDTO(reply.getResult());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets address with mandatory signature verification.
     *
     * @param id the id
     * @return the address (verified)
     * @throws ApiException the api exception
     */
    public Address getAddress(final long id) throws ApiException {

        checkArgument(id > 0, "address id cannot be zero");

        try {
            TgvalidatordGetAddressReply reply = addressesApi.walletServiceGetAddress(String.valueOf(id));
            Address address = AddressMapper.INSTANCE.fromDTO(reply.getResult());

            // Mandatory signature verification
            DecodedRulesContainer rulesContainer = rulesContainerCache.getDecodedRulesContainer();
            AddressSignatureVerifier.verifyAddressSignature(address, rulesContainer);

            return address;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets addresses for a wallet with mandatory signature verification.
     *
     * @param walletId the wallet id
     * @param limit    the maximum number of addresses to return
     * @param offset   the offset for pagination
     * @return the list of addresses (verified)
     * @throws ApiException the api exception
     */
    public List<Address> getAddresses(final long walletId, final int limit, final int offset) throws ApiException {

        checkArgument(walletId > 0, "walletId cannot be zero");
        checkArgument(limit > 0, "limit must be positive");
        checkArgument(offset >= 0, "offset cannot be negative");

        try {
            TgvalidatordGetAddressesReply reply = addressesApi.walletServiceGetAddresses(
                    null,                       // currency
                    null,                       // query
                    String.valueOf(limit),      // limit
                    String.valueOf(offset),     // offset
                    null,                       // scoreProvider
                    null,                       // scoreInBelow
                    null,                       // scoreOutBelow
                    null,                       // scoreExclusive
                    null,                       // onlyPositiveBalance
                    null,                       // sortBy
                    null,                       // sortOrder
                    null,                       // balanceBelow
                    null,                       // balanceAbove
                    String.valueOf(walletId),   // walletId
                    null,                       // customerId
                    null,                       // coinfirmScoreGreater
                    null,                       // chainalysisScoreGreater
                    null,                       // tagIDs
                    null,                       // blockchain
                    null,                       // network
                    null,                       // addressIds
                    null,                       // nfts
                    null,                       // addresses
                    null,                       // scoreFilterScoreProvider
                    null,                       // scoreFilterScorechainFiltersScoreInBelow
                    null,                       // scoreFilterScorechainFiltersScoreOutBelow
                    null,                       // scoreFilterScorechainFiltersScoreExclusive
                    null,                       // scoreFilterCoinfirmFiltersScoreGreater
                    null,                       // scoreFilterChainalysisFiltersScoreGreater
                    null,                       // scoreFilterEllipticFiltersScoreGreater
                    null,                       // scoreFilterMerkleFiltersScoreGreater
                    null,                       // ellipticScoreGreater
                    null                        // merkleScoreGreater
            );

            List<TgvalidatordAddress> result = reply.getResult();
            if (result == null) {
                return Collections.emptyList();
            }

            // Fetch rules container once for all addresses
            DecodedRulesContainer rulesContainer = rulesContainerCache.getDecodedRulesContainer();

            // Map and verify each address
            List<Address> addresses = result.stream()
                    .map(AddressMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());

            // Mandatory signature verification for all addresses
            for (Address address : addresses) {
                AddressSignatureVerifier.verifyAddressSignature(address, rulesContainer);
            }

            return addresses;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Creates an attribute for an address.
     *
     * @param addressId the address id
     * @param key       the attribute key
     * @param value     the attribute value
     * @throws ApiException the api exception
     */
    public void createAddressAttribute(final long addressId, final String key, final String value) throws ApiException {

        checkArgument(addressId > 0, "addressId cannot be zero");
        checkArgument(!Strings.isNullOrEmpty(key), "key cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(value), "value cannot be null or empty");

        try {
            TgvalidatordCreateAddressAttributeRequest attribute = new TgvalidatordCreateAddressAttributeRequest();
            attribute.setKey(key);
            attribute.setValue(value);

            WalletServiceCreateAddressAttributesBody body = new WalletServiceCreateAddressAttributesBody();
            body.addAttributesItem(attribute);

            addressesApi.walletServiceCreateAddressAttributes(String.valueOf(addressId), body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Deletes an attribute from an address.
     *
     * @param addressId   the address id
     * @param attributeId the attribute id
     * @throws ApiException the api exception
     */
    public void deleteAddressAttribute(final long addressId, final long attributeId) throws ApiException {

        checkArgument(addressId > 0, "addressId cannot be zero");
        checkArgument(attributeId > 0, "attributeId cannot be zero");

        try {
            addressesApi.walletServiceDeleteAddressAttribute(
                    String.valueOf(addressId),
                    String.valueOf(attributeId)
            );
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets the proof of reserve for an address.
     *
     * @param addressId the address id
     * @param challenge the challenge string (optional)
     * @return the proof of reserve
     * @throws ApiException the api exception
     */
    public TgvalidatordProofOfReserve getAddressProofOfReserve(final long addressId, final String challenge) throws ApiException {

        checkArgument(addressId > 0, "addressId cannot be zero");

        try {
            TgvalidatordGetAddressProofOfReserveReply reply = addressesApi.walletServiceGetAddressProofOfReserve(
                    String.valueOf(addressId),
                    challenge
            );
            return reply.getResult();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
