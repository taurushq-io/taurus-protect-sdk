package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.helper.AddressSignatureVerifier;
import com.taurushq.sdk.protect.client.mapper.AddressMapper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.AddressResult;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ContainerIntegrityException;
import com.taurushq.sdk.protect.client.model.CreateAddressRequest;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.Pagination;
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
import java.util.HashMap;
import java.util.List;
import java.util.Map;
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
 * // List a wallet's addresses, one page at a time
 * AddressResult page = client.getAddressService().getAddresses(walletId, 20, 0);
 * for (Address a : page.getAddresses()) { ... }
 * // next page: getAddresses(walletId, 20, page.getPagination().getNextOffset())
 * }</pre>
 *
 * @see Address
 * @see CreateAddressRequest
 * @see WalletService
 */
public class AddressService {

    /**
     * The most address ids one by-id re-read sends, validatord's cap on addressIds.
     */
    static final int MAX_ADDRESS_IDS_PER_REQUEST = 50;

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
     * Supplies the SuperAdmin-verified rules container on demand.
     * <p>
     * The seam below takes a source rather than a container so that it can decide
     * whether a container is needed at all: an in-flight asynchronous creation has no
     * address string to check, and must not be made to depend on a governance fetch.
     * It is also what lets a unit test drive the seam with a hand-built container —
     * {@link RulesContainerCache} can only be fed through a live {@code ApiClient}
     * (test scope here is JUnit only: no Mockito, no HTTP stub).
     */
    @FunctionalInterface
    interface RulesContainerSource {

        /**
         * Returns the verified rules container.
         *
         * @return the decoded, SuperAdmin-verified rules container
         * @throws ApiException if fetching the governance rules fails
         */
        DecodedRulesContainer get() throws ApiException;
    }

    /**
     * The ONE verification seam every {@link Address} this SDK returns passes through.
     * <p>
     * The threat model is a response-controlling API server. An {@code Address} carries
     * the deposit destination a caller is about to publish or send funds to, and the
     * only thing binding that string to Taurus Protect is the HSM signature over it. So
     * the invariant is: <b>never hand back a non-empty address string that has not been
     * verified.</b> Three arms, in this order:
     * <ol>
     *   <li><b>No address string</b> (asynchronous creation, status {@code creating}) —
     *       return it. There is no destination to misuse and nothing has been signed
     *       yet, so no container is fetched.</li>
     *   <li><b>Address string but no signature</b> — refuse. This is the arm that was
     *       missing: {@code createAddress} returned the mapped DTO directly, so a
     *       hostile server could answer a creation with an attacker-chosen deposit
     *       address and no signature at all, and the caller received it in exactly the
     *       same {@code Address} type that {@code getAddress} verifies. Handing the
     *       string back "because there was nothing to check" is the bug.</li>
     *   <li><b>Address string and a signature</b> — verify it against the HSMSLOT key
     *       from the verified rules container; a failure is an
     *       {@link IntegrityException}.</li>
     * </ol>
     * A seam rather than a third copy of the verify block is the point: this defect
     * exists because {@code getAddress} and {@code getAddresses} each grew their own
     * inline verification and {@code createAddress} inherited neither.
     * {@code AssetService.getAssetAddresses} — the only other reader of this entity —
     * routes through here too.
     *
     * @param address         the mapped address, may be null (an empty API result)
     * @param containerSource supplies the rules container if a signature must be checked
     * @return the same address, once it satisfies the invariant
     * @throws IntegrityException if the address string cannot be shown to be authentic
     * @throws ApiException       if the rules container could not be fetched
     */
    static Address verifiedAddress(final Address address,
                                   final RulesContainerSource containerSource) throws ApiException {

        if (address == null || !hasAddressString(address)) {
            // Arm 1: nothing to verify, and deliberately no container fetch.
            return address;
        }

        if (Strings.isNullOrEmpty(address.getSignature())) {
            // Arm 2: refuse, and do it before touching the network — an unsigned address
            // is not a transient condition to retry, it is a response we will not trust.
            // The status is named because it is the only way for the caller to tell an
            // in-flight row (which would have no address string) from a stripped
            // signature, and the message points at the verifying getter so the caller
            // does not "work around" this by keeping the value.
            throw new IntegrityException("Address " + address.getId() + " (status "
                    + describeStatus(address) + ") carries a blockchain address but no HSM"
                    + " signature; refusing to return an unverified address. Re-read it with"
                    + " getAddress(" + address.getId() + ") once the status advances past"
                    + " creation.");
        }

        // Arm 3: the real check.
        AddressSignatureVerifier.verifyAddressSignature(address, containerSource.get());
        return address;
    }

    /**
     * Applies {@link #verifiedAddress} to a page of addresses, fail-fast.
     * <p>
     * One unverifiable address is not a row to skip past when the caller is choosing
     * where funds go, so — unlike the whitelist list paths, which exclude and report —
     * this aborts the call.
     * <p>
     * The container is fetched at most once for the whole page: per-row fetching would
     * repeat the round trip on a cold cache, and a TTL expiry mid-page could judge two
     * rows of one response against two different rulesets. A page in which no row has
     * an address string yet fetches nothing.
     *
     * @param addresses       the mapped addresses, may be null or empty
     * @param containerSource supplies the rules container if a signature must be checked
     * @return the same list, once every element satisfies the invariant
     * @throws IntegrityException if any address string cannot be shown to be authentic
     * @throws ApiException       if the rules container could not be fetched
     */
    static List<Address> verifiedAddresses(final List<Address> addresses,
                                           final RulesContainerSource containerSource) throws ApiException {

        if (addresses == null || addresses.isEmpty()) {
            return addresses;
        }

        final DecodedRulesContainer container =
                anyAddressToVerify(addresses) ? containerSource.get() : null;

        for (Address address : addresses) {
            verifiedAddress(address, () -> container);
        }
        return addresses;
    }

    /**
     * Reports whether this row carries a blockchain address string at all.
     * <p>
     * An empty string is the asynchronous-creation case, not a tampered row: there is
     * no destination a caller could act on.
     */
    private static boolean hasAddressString(final Address address) {
        return !Strings.isNullOrEmpty(address.getAddress());
    }

    /**
     * Reports whether any row in the page needs the rules container.
     * <p>
     * Shares {@link #hasAddressString} with the seam so the two cannot disagree about
     * which rows require a signature check.
     */
    private static boolean anyAddressToVerify(final List<Address> addresses) {
        for (Address address : addresses) {
            if (address != null && hasAddressString(address)) {
                return true;
            }
        }
        return false;
    }

    /**
     * Renders the server-reported status for an error message, never null or empty.
     */
    private static String describeStatus(final Address address) {
        return Strings.isNullOrEmpty(address.getStatus()) ? "unreported" : address.getStatus();
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

            // Through the same seam as the read paths. A creation reply carries a
            // deposit address in the very type getAddress verifies, so returning the
            // mapper output raw here made the mandatory verification on the getters
            // avoidable: ask the server to create an address and it can name one.
            return verifiedAddress(AddressMapper.INSTANCE.fromDTO(reply.getResult()),
                    rulesContainerCache::getDecodedRulesContainer);
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

            // Mandatory signature verification, through the shared seam
            return verifiedAddress(AddressMapper.INSTANCE.fromDTO(reply.getResult()),
                    rulesContainerCache::getDecodedRulesContainer);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }


    /**
     * Gets a page of a wallet's addresses, with mandatory signature verification.
     *
     * @param walletId the wallet id
     * @param limit    the page size, 0 for the default ({@link Pagination#DEFAULT_PAGE_SIZE})
     * @param offset   the offset, 0 for the first page
     * @return the verified addresses and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if walletId, limit or offset is out of range
     */
    public AddressResult getAddresses(final long walletId, final int limit, final long offset)
            throws ApiException {
        checkArgument(walletId > 0, "walletId cannot be zero");
        return getAddresses(Long.valueOf(walletId), limit, offset, null);
    }

    /**
     * Gets a page of addresses, with mandatory signature verification.
     *
     * @param walletId        the wallet to list, or null for every wallet
     * @param limit           the page size, 0 for the default
     * @param offset          the offset, 0 for the first page
     * @param excludeDisabled true to hide disabled addresses; null or false keeps them
     * @return the verified addresses and their pagination
     * @throws ApiException             the api exception
     * @throws IllegalArgumentException if walletId, limit or offset is out of range
     */
    public AddressResult getAddresses(final Long walletId, final int limit, final long offset,
                                      final Boolean excludeDisabled) throws ApiException {
        checkArgument(walletId == null || walletId > 0, "walletId cannot be zero");
        final int size = PagedOperation.ADDRESSES.resolveSize("limit", limit);
        final long from = Pagination.resolveOffset("offset", offset);

        try {
            TgvalidatordGetAddressesReply reply = fetchAddresses(size, from, walletId, null,
                    Boolean.TRUE.equals(excludeDisabled) ? "exclude" : null);

            List<TgvalidatordAddress> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            List<Address> addresses = rows.stream()
                    .map(AddressMapper.INSTANCE::fromDTO)
                    .collect(Collectors.toList());

            // Mandatory signature verification for all addresses, through the shared
            // seam (which fetches the rules container once for the whole page)
            return new AddressResult(
                    verifiedAddresses(addresses, rulesContainerCache::getDecodedRulesContainer),
                    PagedOperation.ADDRESSES.offsetPage(size, from, rows.size(), 0,
                            reply.getTotalItems(), reply.getOffset()));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Re-reads managed addresses by id through the verifying seam, at most
     * {@link #MAX_ADDRESS_IDS_PER_REQUEST} ids per request.
     * <p>
     * Returns the rows whose address string verified, keyed by id; every other returned row
     * is reported in {@code failures} (id to reason) rather than failing the call, so the
     * caller decides what to do without it. A row with no address string yet has nothing
     * the HSM signed and is reported too. A rules container without an HSMSLOT key cannot
     * verify any address, so it fails the call, as does any request error.
     *
     * @param ids      the address ids, each sent once
     * @param failures receives the id and reason of every returned row that did not verify
     * @return the verified addresses by id
     * @throws ApiException                if a request or the rules container fetch fails
     * @throws ContainerIntegrityException if the rules container has no HSMSLOT key
     */
    Map<String, Address> verifiedAddressesById(final List<String> ids,
                                               final Map<String, String> failures) throws ApiException {
        Map<String, Address> verified = new HashMap<>();
        if (ids.isEmpty()) {
            return verified;
        }
        final DecodedRulesContainer container = rulesContainerCache.getDecodedRulesContainer();
        if (container == null || container.getHsmPublicKey() == null) {
            throw new ContainerIntegrityException(
                    "the rules container has no HSMSLOT key, so no address can be verified");
        }
        for (int from = 0; from < ids.size(); from += MAX_ADDRESS_IDS_PER_REQUEST) {
            List<String> chunk = ids.subList(from, Math.min(ids.size(), from + MAX_ADDRESS_IDS_PER_REQUEST));
            TgvalidatordGetAddressesReply reply;
            try {
                reply = fetchAddresses(chunk.size(), 0, null, chunk, null);
            } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
                throw apiExceptionMapper.toApiException(e);
            }
            List<TgvalidatordAddress> rows = reply.getResult() == null
                    ? Collections.emptyList() : reply.getResult();
            for (TgvalidatordAddress row : rows) {
                Address address = AddressMapper.INSTANCE.fromDTO(row);
                String id = String.valueOf(address.getId());
                if (!hasAddressString(address)) {
                    failures.put(id, "address " + id + " (status " + describeStatus(address)
                            + ") carries no address string to verify");
                    continue;
                }
                try {
                    verified.put(id, verifiedAddress(address, () -> container));
                } catch (ContainerIntegrityException e) {
                    throw e;
                } catch (IntegrityException e) {
                    failures.put(id, e.getMessage());
                }
            }
        }
        return verified;
    }

    /**
     * The one call into the generated address list, so the page read and the by-id re-read
     * cannot drift on which parameters they send.
     *
     * @param size             the page size to send
     * @param from             the offset, 0 for none
     * @param walletId         the wallet filter, or null
     * @param addressIds       the id filter, or null
     * @param includeDisabled  the includeDisabledAddresses value, or null for the server default
     * @return the raw reply
     * @throws com.taurushq.sdk.protect.openapi.ApiException if the request fails
     */
    private TgvalidatordGetAddressesReply fetchAddresses(final int size, final long from, final Long walletId,
                                                         final List<String> addressIds,
                                                         final String includeDisabled)
            throws com.taurushq.sdk.protect.openapi.ApiException {
        return addressesApi.walletServiceGetAddresses(
                null,                       // currency
                null,                       // query
                String.valueOf(size),       // limit
                from == 0 ? null : String.valueOf(from), // offset
                null,                       // scoreProvider
                null,                       // scoreInBelow
                null,                       // scoreOutBelow
                null,                       // scoreExclusive
                null,                       // onlyPositiveBalance
                null,                       // sortBy
                null,                       // sortOrder
                null,                       // balanceBelow
                null,                       // balanceAbove
                walletId == null ? null : String.valueOf(walletId), // walletId
                null,                       // customerId
                null,                       // coinfirmScoreGreater
                null,                       // chainalysisScoreGreater
                null,                       // tagIDs
                null,                       // blockchain
                null,                       // network
                addressIds,                 // addressIds
                null,                       // nfts
                null,                       // addresses
                null,                       // scoreFilterScoreProvider
                null,                       // scoreFilterScorechainFiltersScoreInBelow
                null,                       // scoreFilterScorechainFiltersScoreOutBelow
                null,                       // scoreFilterScorechainFiltersScoreExclusive
                null,                       // scoreFilterCoinfirmFiltersScoreGreater
                null,                       // scoreFilterChainalysisFiltersScoreGreater
                null,                       // scoreFilterEllipticFiltersScoreGreater
                null,                       // scoreFilterTrmlabsFiltersScoreGreater
                null,                       // attributeFiltersJson
                null,                       // attributeFiltersOperator
                includeDisabled             // includeDisabledAddresses
        );
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
