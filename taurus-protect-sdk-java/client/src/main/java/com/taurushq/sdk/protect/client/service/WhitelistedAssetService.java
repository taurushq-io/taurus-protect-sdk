package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.helper.AssetHashHelper;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.helper.WhitelistHashHelper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.mapper.WhitelistedAssetMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistedAsset;
import com.taurushq.sdk.protect.client.model.WhitelistedAssetResult;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ContractWhitelistingApi;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveWhitelistedContractAddressRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;

import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SignatureException;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;
import static org.bouncycastle.util.Strings.constantTimeAreEqual;

/**
 * Service for retrieving and verifying whitelisted assets (contract addresses).
 * Performs full cryptographic verification including:
 * <ul>
 *   <li>Metadata hash verification</li>
 *   <li>Rules container signature verification (SuperAdmin)</li>
 *   <li>Whitelist signature verification (per governance rules)</li>
 * </ul>
 */
public class WhitelistedAssetService {

    private static final Logger LOGGER = Logger.getLogger(WhitelistedAssetService.class.getName());
    private static final Gson GSON = new Gson();
    private static final Pattern IS_NFT_TRAILING_PATTERN =
            Pattern.compile(",\"isNFT\":(true|false)");
    private static final Pattern IS_NFT_LEADING_PATTERN =
            Pattern.compile("\"isNFT\":(true|false),");
    private static final Pattern KIND_TYPE_TRAILING_PATTERN =
            Pattern.compile(",\"kindType\":\"[^\"]*\"");
    private static final Pattern KIND_TYPE_LEADING_PATTERN =
            Pattern.compile("\"kindType\":\"[^\"]*\",");

    private final ContractWhitelistingApi contractWhitelistingApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final List<PublicKey> superAdminPublicKeys;
    private final int minValidSignatures;

    /**
     * Instantiates a new Whitelisted Asset service.
     *
     * @param openApiClient        the OpenAPI client
     * @param apiExceptionMapper   the API exception mapper
     * @param superAdminPublicKeys the list of SuperAdmin public keys for rules verification
     * @param minValidSignatures   the minimum number of valid signatures required for rules
     */
    public WhitelistedAssetService(final ApiClient openApiClient,
                                   final ApiExceptionMapper apiExceptionMapper,
                                   final List<PublicKey> superAdminPublicKeys,
                                   final int minValidSignatures) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");
        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys cannot be empty");
        checkArgument(minValidSignatures > 0, "minValidSignatures must be positive");

        this.apiExceptionMapper = apiExceptionMapper;
        this.contractWhitelistingApi = new ContractWhitelistingApi(openApiClient);
        this.superAdminPublicKeys = superAdminPublicKeys;
        this.minValidSignatures = minValidSignatures;
    }

    /**
     * Gets a whitelisted asset by ID with full signature verification.
     *
     * @param id the whitelisted asset (contract address) ID
     * @return the verified and decoded WhitelistedAsset
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if signature verification fails
     */
    public WhitelistedAsset getWhitelistedAsset(final long id) throws ApiException, WhitelistException {
        SignedWhitelistedAssetEnvelope envelope = getWhitelistedAssetEnvelope(id);
        return envelope.getWhitelistedAsset();
    }

    /**
     * Gets the signed whitelisted asset envelope by ID.
     * Performs full verification including metadata hash, rules container signatures,
     * and whitelist signatures.
     *
     * @param id the whitelisted asset (contract address) ID
     * @return the signed whitelisted asset envelope with verified WhitelistedAsset
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public SignedWhitelistedAssetEnvelope getWhitelistedAssetEnvelope(final long id)
            throws ApiException, WhitelistException {
        checkArgument(id > 0, "whitelisted asset id cannot be zero");

        try {
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopeReply reply =
                    contractWhitelistingApi.whitelistServiceGetWhitelistedContract(String.valueOf(id));
            SignedWhitelistedAssetEnvelope envelope =
                    WhitelistedAssetMapper.INSTANCE.fromDTO(reply.getResult());
            initializeEnvelope(envelope);
            return envelope;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets a list of whitelisted asset envelopes with pagination.
     *
     * @param limit  the maximum number of results
     * @param offset the offset for pagination
     * @return the list of signed whitelisted asset envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAssetEnvelope> getWhitelistedAssets(int limit, int offset)
            throws ApiException, WhitelistException {
        return getWhitelistedAssets(limit, offset, null, null);
    }

    /**
     * Gets a list of whitelisted asset envelopes filtered by blockchain.
     *
     * @param limit      the maximum number of results
     * @param offset     the offset for pagination
     * @param blockchain filter by blockchain (e.g., "ETH", "BTC")
     * @return the list of signed whitelisted asset envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAssetEnvelope> getWhitelistedAssets(int limit, int offset,
                                                                      String blockchain)
            throws ApiException, WhitelistException {
        return getWhitelistedAssets(limit, offset, blockchain, null);
    }

    /**
     * Gets a list of whitelisted asset envelopes filtered by blockchain and network.
     *
     * @param limit      the maximum number of results
     * @param offset     the offset for pagination
     * @param blockchain filter by blockchain (e.g., "ETH", "BTC")
     * @param network    filter by network (e.g., "mainnet", "testnet")
     * @return the list of signed whitelisted asset envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAssetEnvelope> getWhitelistedAssets(int limit, int offset,
                                                                      String blockchain, String network)
            throws ApiException, WhitelistException {
        return getWhitelistedAssets(limit, offset, blockchain, network,
                null, null, null, null).getAssets();
    }

    /**
     * Gets a page of whitelisted asset envelopes with the full filter set and the page total.
     *
     * @param limit              the maximum number of results
     * @param offset             the offset for pagination
     * @param blockchain         filter by blockchain (e.g., "ETH", "BTC"), or null
     * @param network            filter by network (e.g., "mainnet", "testnet"), or null
     * @param query              search across address, symbol and name, or null
     * @param includeForApproval include assets pending approval, or null
     * @param kindTypes          filter by contract kind ("nft", "token"), or null
     * @param ids                filter by specific whitelisted asset IDs, or null
     * @return the page of verified asset envelopes and the total item count
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public WhitelistedAssetResult getWhitelistedAssets(int limit, int offset,
                                                        String blockchain, String network,
                                                        String query, Boolean includeForApproval,
                                                        List<String> kindTypes, List<String> ids)
            throws ApiException, WhitelistException {
        try {
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply =
                    contractWhitelistingApi.whitelistServiceGetWhitelistedContracts(
                            String.valueOf(limit),
                            String.valueOf(offset),
                            query,
                            blockchain,
                            includeForApproval,
                            network,
                            null,       // isNFT — deprecated, superseded by kindTypes
                            ids,
                            kindTypes);

            return toVerifiedResult(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Re-reads a batch of assets through the verifying list path, filtered by id, and
     * returns them keyed by id.
     *
     * <pre>
     *   ids -&gt; ONE filtered page -&gt; verify every row -&gt; map by id
     * </pre>
     *
     * <p>One round trip and one rules-container fetch for the whole batch. The per-id
     * GET this replaced cost both per id, so a 50-id approval was 50 sequential round
     * trips each running the full verification chain.
     *
     * <p>Package-private so the completeness behaviour can be tested directly.
     *
     * @param ids the ids to re-read
     * @return the verified envelopes keyed by id
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    java.util.Map<String, SignedWhitelistedAssetEnvelope> verifiedAssetsById(final List<String> ids)
            throws ApiException, WhitelistException {
        WhitelistedAssetResult result = getWhitelistedAssets(
                ids.size(), 0, null, null, null, Boolean.TRUE, null, ids);

        java.util.Map<String, SignedWhitelistedAssetEnvelope> byId = new java.util.HashMap<>();
        for (SignedWhitelistedAssetEnvelope envelope : result.getAssets()) {
            // The id lives on the ENVELOPE: getWhitelistedAsset() is gated on
            // verification having run, so reading through it here would couple the key
            // to that gate for no reason.
            if (envelope != null) {
                byId.put(String.valueOf(envelope.getId()), envelope);
            }
        }
        return byId;
    }

    /**
     * Gets a page of whitelisted assets awaiting approval, verified as in
     * {@link #getWhitelistedAssets(int, int)}.
     * <p>
     * Without this the only reader of the for-approval endpoint was the unverified
     * contract-whitelisting service, so the rows an approver inspects before whitelisting a
     * contract address were never checked against governance.
     *
     * @param limit  the maximum number of results
     * @param offset the offset for pagination
     * @param ids    filter by specific whitelisted asset IDs, or null
     * @return the page of verified asset envelopes and the total item count
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public WhitelistedAssetResult getWhitelistedAssetsForApproval(int limit, int offset,
                                                                   List<String> ids)
            throws ApiException, WhitelistException {
        try {
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply =
                    contractWhitelistingApi.whitelistServiceGetWhitelistedContractsForApproval(
                            String.valueOf(limit),
                            String.valueOf(offset),
                            ids);

            return toVerifiedResult(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Signs and submits an approval for the given whitelisted assets, all-or-nothing.
     * <p>
     * Each asset is re-read through the verified path and the hashes THOSE rows carry are
     * what gets signed, so the approver's signature covers metadata this SDK checked rather
     * than whatever a caller was handed.
     * {@link ContractWhitelistingService#approveWhitelistedContracts} takes an opaque
     * signature over hashes nothing verified.
     * <p>
     * Any asset that is missing or fails verification aborts the whole call and nothing is
     * signed: one signature covers every hash in the batch, so a partial approval would mean
     * the caller believes they approved more than they did.
     *
     * @param ids        the whitelisted asset IDs to approve
     * @param privateKey the approver's P-256 private key
     * @param comment    the approval comment
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if any asset is missing or fails verification
     */
    public void approveWhitelistedAssets(final List<Long> ids, final PrivateKey privateKey,
                                         final String comment) throws ApiException, WhitelistException {
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");
        checkNotNull(privateKey, "privateKey cannot be null");
        checkArgument(!Strings.isNullOrEmpty(comment), "comment is required");
        ids.forEach(id -> checkArgument(id != null && id > 0,
                "whitelisted asset id cannot be zero or negative"));

        // Sorted on a COPY, as the request-approval path does, so the signed order is
        // independent of the order the caller passed and their list is not mutated.
        List<Long> sorted = new ArrayList<>(ids);
        Collections.sort(sorted);

        // ONE id-filtered page through the verifying list path, not one GET per id. The
        // list path verifies every row and fetches the rules container once per call, so
        // a 50-id approval costs one round trip and one container fetch instead of fifty
        // of each. includeForApproval is required: the rows being approved are pending,
        // so the default list does not return them.
        List<String> idStrings = sorted.stream().map(String::valueOf).collect(Collectors.toList());
        java.util.Map<String, SignedWhitelistedAssetEnvelope> byId =
                verifiedAssetsById(idStrings);

        List<String> hashes = new ArrayList<>(sorted.size());
        for (String id : idStrings) {
            SignedWhitelistedAssetEnvelope envelope = byId.get(id);
            if (envelope == null) {
                // A page that silently omits a row must not become an approval of fewer
                // rows than the caller asked for.
                throw new IntegrityException(String.format(
                        "refusing to sign: asset %s was not returned by the verified read", id));
            }
            if (envelope.getMetadata() == null
                    || Strings.isNullOrEmpty(envelope.getMetadata().getHash())) {
                throw new IntegrityException(String.format(
                        "refusing to sign: asset %s has no metadata hash", id));
            }
            hashes.add(envelope.getMetadata().getHash());
        }

        String toSign = GSON.toJson(hashes);
        TgvalidatordApproveWhitelistedContractAddressRequest request =
                new TgvalidatordApproveWhitelistedContractAddressRequest();
        request.setIds(sorted.stream().map(String::valueOf).collect(Collectors.toList()));
        request.setComment(comment);
        try {
            request.setSignature(CryptoTPV1.calculateBase64Signature(
                    privateKey, toSign.getBytes(StandardCharsets.UTF_8)));
        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException ex) {
            ApiException e = new ApiException();
            e.setCode(400);
            e.setError("ClientInvalidRequest");
            e.setMessage("unable to sign the array of whitelisted asset hashes");
            e.setOriginalException(ex);
            throw e;
        }

        try {
            contractWhitelistingApi.whitelistServiceApproveWhitelistedContract(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Maps and fully verifies every row of a reply.
     */
    private WhitelistedAssetResult toVerifiedResult(
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply)
            throws WhitelistException {
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        List<SignedWhitelistedAssetEnvelope> envelopes = new ArrayList<>();

        if (reply.getResult() != null) {
            for (TgvalidatordSignedWhitelistedContractAddressEnvelope dto : reply.getResult()) {
                SignedWhitelistedAssetEnvelope envelope =
                        WhitelistedAssetMapper.INSTANCE.fromDTO(dto);
                initializeEnvelope(envelope);
                envelopes.add(envelope);
            }
        }

        result.setAssets(envelopes);
        if (reply.getTotalItems() != null) {
            try {
                result.setTotalItems(Long.parseLong(reply.getTotalItems()));
            } catch (NumberFormatException e) {
                result.setTotalItems(0L);
            }
        }
        return result;
    }

    /**
     * Initializes and fully verifies the envelope.
     * After this method completes, the envelope's getWhitelistedAsset() will return
     * the verified asset.
     */
    // 5-step verification for a whitelisted asset, plus the parse that makes it usable.
    //
    //   envelope ──▶ 1 hash ──▶ 2 SuperAdmin sigs ──▶ 3 decode rules
    //                                                      │
    //                6 parse VERIFIED payload ◀── 5 thresholds ◀── 4 coverage
    //                           │                       │                │
    //                           ▼                       │      returns the hash it
    //                    returned to caller             │      matched, which is
    //                                                   │      what step 5 checks
    //                                                   ▼
    //                                per group: DISTINCT signers, keyed on the
    //                                container's public key, never the entry's userId
    //
    // Steps 1-5 prove the envelope is authentic; step 6 is what stops an unsigned
    // value reaching the caller. Assets use ContractAddressWhitelistingRules and the
    // ASSET legacy hashes (isNFT / kindType) — not the address ones.
    //
    // All six run here in initializeEnvelope, which ends by calling
    // AssetHashHelper.parseWhitelistedAssetFromJson.
    private void initializeEnvelope(SignedWhitelistedAssetEnvelope envelope) throws WhitelistException {
        // Precondition checks
        if (envelope.getSignedAsset() == null || envelope.getSignedAsset().getPayload() == null) {
            throw new WhitelistException("signed asset payload is null");
        }
        if (envelope.getSignedAsset().getSignatures() == null
                || envelope.getSignedAsset().getSignatures().isEmpty()) {
            throw new WhitelistException("signatures is null or empty");
        }

        // Step 1: Verify computed hash of payloadAsString equals received hash
        verifyMetadataHash(envelope);

        // Step 2: Verify rulesContainer signatures (SuperAdmin)
        verifyRulesContainerSignatures(envelope);

        // Step 3: Decode rulesContainer
        DecodedRulesContainer rulesContainer = decodeRulesContainer(envelope);

        // Step 4: Verify metadata.hash is in signed hashes list
        String verifiedHash = verifyHashInSignedHashes(envelope);

        // Step 5: Verify whitelist signatures are valid per governance rules
        verifyWhitelistSignatures(envelope, rulesContainer, verifiedHash);

        // Step 6: the envelope DERIVES the asset from its own signed payloadAsString.
        // It used to be parsed here and handed to a public setter, which meant the
        // "verified" marker could be flipped with an asset the caller chose.
        envelope.markVerified(rulesContainer);
    }

    /**
     * Verifies that the computed hash of payloadAsString equals the provided hash.
     */
    private void verifyMetadataHash(SignedWhitelistedAssetEnvelope envelope) {
        if (envelope.getMetadata() == null) {
            throw new IntegrityException("metadata is null");
        }
        if (Strings.isNullOrEmpty(envelope.getMetadata().getPayloadAsString())) {
            throw new IntegrityException("metadata payloadAsString is null or empty");
        }
        if (Strings.isNullOrEmpty(envelope.getMetadata().getHash())) {
            throw new IntegrityException("metadata hash is null or empty");
        }

        String computedHash = CryptoTPV1.calculateHexHash(envelope.getMetadata().getPayloadAsString());
        String providedHash = envelope.getMetadata().getHash();

        if (!constantTimeAreEqual(computedHash, providedHash)) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                // SECURITY: Do not log hash values to prevent information leakage
                LOGGER.warning("Metadata hash verification failed for whitelisted asset");
            }
            throw new IntegrityException("metadata hash verification failed");
        }
    }

    /**
     * Verifies the rulesContainer signatures against SuperAdmin public keys.
     */
    private void verifyRulesContainerSignatures(SignedWhitelistedAssetEnvelope envelope)
            throws WhitelistException {
        if (Strings.isNullOrEmpty(envelope.getRulesSignatures())) {
            throw new WhitelistException("rules signatures is null or empty");
        }
        if (Strings.isNullOrEmpty(envelope.getRulesContainer())) {
            throw new WhitelistException("rules container is null or empty");
        }

        // Decode the rulesSignatures protobuf (base64-encoded UserSignatures)
        List<RuleUserSignature> signatures;
        try {
            signatures = RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(
                    envelope.getRulesSignatures());
        } catch (IllegalArgumentException | InvalidProtocolBufferException e) {
            throw new WhitelistException("unable to decode rules signatures", e);
        }

        if (signatures.isEmpty()) {
            throw new WhitelistException("no rules signatures present");
        }

        // Verify signatures against the rulesContainer bytes
        byte[] rulesData = Base64.getDecoder().decode(envelope.getRulesContainer());
        try {
            SignatureVerifier.verifyGovernanceRulesSignatures(rulesData, signatures,
                    superAdminPublicKeys, minValidSignatures);
        } catch (IntegrityException e) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                LOGGER.warning("Rules container verification failed: insufficient valid signatures");
            }
            throw new IntegrityException(
                    "Rules container verification failed: " + e.getMessage(), e);
        }
        if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
            LOGGER.fine("Rules container signature verification succeeded");
        }
    }

    /**
     * Decodes the rulesContainer protobuf and maps it to the client model.
     */
    private DecodedRulesContainer decodeRulesContainer(SignedWhitelistedAssetEnvelope envelope)
            throws WhitelistException {
        try {
            return RulesContainerMapper.INSTANCE.fromBase64String(envelope.getRulesContainer());
        } catch (IllegalArgumentException | InvalidProtocolBufferException e) {
            throw new WhitelistException("unable to decode rules container", e);
        }
    }

    /**
     * Verifies that the metadata hash is present in at least one signature's hashes list.
     * For backward compatibility, also tries alternative hashes for assets signed
     * before certain fields were added to the schema.
     */
    private String verifyHashInSignedHashes(SignedWhitelistedAssetEnvelope envelope)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();
        List<WhitelistSignature> signatures = envelope.getSignedAsset().getSignatures();

        // First, try the provided hash directly
        if (SignatureVerifier.verifyHashCoverage(metadataHash, signatures)) {
            return metadataHash;
        }

        // If not found, try alternative hashes for backward compatibility
        // (handles assets signed before schema changes)
        for (String legacyHash : computeLegacyHashes(envelope.getMetadata().getPayloadAsString())) {
            if (SignatureVerifier.verifyHashCoverage(legacyHash, signatures)) {
                // Returned rather than written back onto the caller's envelope:
                // verification must not mutate its input.
                return legacyHash;
            }
        }

        if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
            LOGGER.warning("Metadata hash not found in any signature's hashes list");
        }
        throw new IntegrityException("metadata hash not found in any signature's hashes list");
    }

    /**
     * Computes legacy hashes by applying transformation combinations to handle schema evolution.
     * Returns a list of possible legacy hashes to try.
     *
     * <p>Strategies cover schema evolution scenarios for contract addresses:
     * <ul>
     *   <li>Strategy 1: Remove optional fields that may have been added after signing</li>
     * </ul>
     *
     * @param payloadAsString the current payload string
     * @return list of legacy hashes to try (may be empty if no transformations apply)
     */
    private List<String> computeLegacyHashes(String payloadAsString) {
        if (payloadAsString == null) {
            return Collections.emptyList();
        }

        Set<String> uniqueHashes = new LinkedHashSet<>();

        // Strategy 1: Remove optional fields that might not have existed when signed
        // E.g., remove "isNFT" field if it was added later
        String withoutIsNFT = IS_NFT_TRAILING_PATTERN.matcher(payloadAsString).replaceAll("");
        withoutIsNFT = IS_NFT_LEADING_PATTERN.matcher(withoutIsNFT).replaceAll("");
        if (!withoutIsNFT.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutIsNFT));
        }

        // Strategy 2: Remove "kindType" field if it was added later
        String withoutKindType = KIND_TYPE_TRAILING_PATTERN.matcher(payloadAsString).replaceAll("");
        withoutKindType = KIND_TYPE_LEADING_PATTERN.matcher(withoutKindType).replaceAll("");
        if (!withoutKindType.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutKindType));
        }

        // Strategy 3: Remove both isNFT and kindType
        String withoutBoth = IS_NFT_TRAILING_PATTERN.matcher(payloadAsString).replaceAll("");
        withoutBoth = IS_NFT_LEADING_PATTERN.matcher(withoutBoth).replaceAll("");
        withoutBoth = KIND_TYPE_TRAILING_PATTERN.matcher(withoutBoth).replaceAll("");
        withoutBoth = KIND_TYPE_LEADING_PATTERN.matcher(withoutBoth).replaceAll("");
        if (!withoutBoth.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutBoth));
        }

        return new ArrayList<>(uniqueHashes);
    }

    /**
     * Encodes each signature's hashes array to its signed JSON form, once.
     *
     * <p>Keyed by position in {@code signatures}, mirroring the Go, Python and
     * TypeScript verifiers. A signature whose hashes cannot be encoded is simply
     * absent from the map and is skipped by the caller.
     */
    private Map<Integer, byte[]> precomputeHashesJson(List<WhitelistSignature> signatures) {
        Map<Integer, byte[]> encoded = new HashMap<>();
        if (signatures == null) {
            return encoded;
        }
        for (int i = 0; i < signatures.size(); i++) {
            WhitelistSignature sig = signatures.get(i);
            if (sig == null || sig.getHashes() == null) {
                continue;
            }
            encoded.put(i, GSON.toJson(sig.getHashes()).getBytes(StandardCharsets.UTF_8));
        }
        return encoded;
    }

    /**
     * Verifies whitelist signatures according to governance rules threshold requirements.
     */
    private void verifyWhitelistSignatures(SignedWhitelistedAssetEnvelope envelope,
                                           DecodedRulesContainer rulesContainer,
                                           String metadataHash)
            throws WhitelistException {
        // Serialised once here, not per signature inside the group loop: that loop
        // runs for every group of every parallel path, so the same array was being
        // re-encoded path x group x signature times right next to the ECDSA verify.
        // Go, Python and TypeScript all precompute this map.
        Map<Integer, byte[]> hashesJsonBySignature =
                precomputeHashesJson(envelope.getSignedAsset().getSignatures());

        // Keyed off the SIGNED payload, not the response. See the address service.
        String[] ruleKey = WhitelistHashHelper.resolveRuleKey(
                envelope.getMetadata().getPayloadAsString(),
                envelope.getBlockchain(), envelope.getNetwork());

        // Find matching contract address whitelisting rules
        ContractAddressWhitelistingRules whitelistRules = rulesContainer.findContractAddressWhitelistingRules(
                ruleKey[0], ruleKey[1]);
        if (whitelistRules == null) {
            throw new WhitelistException("no contract address whitelisting rules found for blockchain="
                    + envelope.getBlockchain() + " network=" + envelope.getNetwork());
        }

        // Contract whitelisting uses parallelThresholds directly (no rule lines matching)
        List<SequentialThresholds> parallelThresholds = whitelistRules.getParallelThresholds();
        if (parallelThresholds == null || parallelThresholds.isEmpty()) {
            throw new WhitelistException("no threshold rules defined");
        }

        // Try to verify all paths
        List<String> pathFailures = tryVerifyAllPaths(
                parallelThresholds, rulesContainer, envelope.getSignedAsset().getSignatures(),
                metadataHash, hashesJsonBySignature);
        if (!pathFailures.isEmpty()) {
            throw new WhitelistException("signature verification failed for whitelisted asset (ID: "
                    + envelope.getId() + ") : no approval path satisfied the threshold requirements. "
                    + String.join("; ", pathFailures));
        }
    }

    /**
     * Tries to verify all parallel threshold paths.
     *
     * @return empty list if verification passed, or list of failure messages if all paths failed
     */
    private List<String> tryVerifyAllPaths(List<SequentialThresholds> parallelThresholds,
                                           DecodedRulesContainer rulesContainer,
                                           List<WhitelistSignature> signatures,
                                           String metadataHash,
                                           Map<Integer, byte[]> hashesJsonBySignature) {
        List<String> pathFailures = new ArrayList<>();
        for (int i = 0; i < parallelThresholds.size(); i++) {
            SequentialThresholds seqThreshold = parallelThresholds.get(i);
            try {
                verifySequentialThresholds(seqThreshold, rulesContainer, signatures, metadataHash,
                        hashesJsonBySignature);
                return Collections.emptyList();  // Verification passed
            } catch (IntegrityException e) {
                pathFailures.add("Path " + (i + 1) + ": " + sanitizeVerificationError(e));
            }
        }
        return pathFailures;
    }

    /**
     * Sanitizes exception messages for external exposure.
     * IntegrityException messages are designed to be safe, other exceptions get generic messages.
     */
    private String sanitizeVerificationError(Exception e) {
        if (e instanceof IntegrityException) {
            return e.getMessage();
        }
        return "verification failed";
    }

    /**
     * Verifies all group thresholds in a sequential threshold path.
     *
     * @throws IntegrityException if any group threshold is not met
     */
    private void verifySequentialThresholds(SequentialThresholds seqThreshold,
                                            DecodedRulesContainer rulesContainer,
                                            List<WhitelistSignature> signatures,
                                            String metadataHash,
                                            Map<Integer, byte[]> hashesJsonBySignature) {
        List<GroupThreshold> thresholds = seqThreshold.getThresholds();
        if (thresholds == null || thresholds.isEmpty()) {
            throw new IntegrityException("no group thresholds defined");
        }

        // ALL group thresholds must be satisfied (AND logic)
        for (GroupThreshold groupThreshold : thresholds) {
            verifyGroupThreshold(groupThreshold, rulesContainer, signatures, metadataHash,
                    hashesJsonBySignature);
        }
    }

    /**
     * Verifies that a group threshold is met.
     *
     * @throws IntegrityException if the threshold is not met, with detailed reason
     */
    private void verifyGroupThreshold(GroupThreshold groupThreshold,
                                      DecodedRulesContainer rulesContainer,
                                      List<WhitelistSignature> signatures,
                                      String metadataHash,
                                      Map<Integer, byte[]> hashesJsonBySignature) {
        String groupId = groupThreshold.getGroupId();
        int minSigs = groupThreshold.getMinimumSignatures();

        RuleGroup group = rulesContainer.findGroupById(groupId);
        if (group == null) {
            throw new IntegrityException(
                    String.format("group '%s' not found in rules container", groupId));
        }

        List<String> groupUserIds = group.getUserIds();
        if (groupUserIds == null || groupUserIds.isEmpty()) {
            if (minSigs > 0) {
                throw new IntegrityException(
                        String.format("group '%s' has no users but requires %d signature(s)",
                                groupId, minSigs));
            }
            return; // minSignatures == 0, so empty group is OK
        }

        // A populated group with a zero threshold is a malformed container, not a
        // group anyone may satisfy. There is no post-loop threshold check -- the
        // only success exit is inside the loop after an increment -- so a zero here
        // silently means "one signature suffices", turning a 2-of-N group into
        // 1-of-N. Fail closed.
        if (minSigs <= 0) {
            throw new IntegrityException(
                    String.format("group '%s' has %d user(s) but requires 0 signature(s): "
                            + "minimumSignatures must be positive", groupId, groupUserIds.size()));
        }

        // Convert to set for faster lookup
        Set<String> groupUserIdSet = new HashSet<>(groupUserIds);

        // Count DISTINCT signers, not signature entries.
        //
        // The entries come from the server-supplied userSignatures blob, so counting them
        // let a duplicated entry from one group member satisfy an N-of-M group. Keyed on
        // the container-resolved public key rather than the server-supplied userId, so one
        // compromised key shared by two IDs counts once.
        Set<String> signers = new HashSet<>();
        List<String> skippedReasons = new ArrayList<>();

        for (int i = 0; i < signatures.size(); i++) {
            WhitelistSignature sig = signatures.get(i);
            WhitelistUserSignature userSig = sig.getSignature();
            if (userSig == null) {
                skippedReasons.add("signature has null userSig");
                continue;
            }

            String sigUserId = userSig.getUserId();
            if (!groupUserIdSet.contains(sigUserId)) {
                continue; // Signer not in this group - not an error, just not relevant
            }

            // Check that metadata hash is covered by this signature
            List<String> hashes = sig.getHashes();
            if (!SignatureVerifier.containsHash(hashes, metadataHash)) {
                skippedReasons.add(String.format(
                        "user '%s' signature does not cover metadata hash '%s' (signed hashes=%s)",
                        sigUserId, metadataHash, hashes));
                continue;
            }

            RuleUser user = rulesContainer.findUserById(sigUserId);
            if (user == null) {
                skippedReasons.add(String.format("user '%s' not found in rules container", sigUserId));
                continue;
            }
            if (user.getPublicKey() == null) {
                skippedReasons.add(String.format("user '%s' has no public key", sigUserId));
                continue;
            }

            // Verify signature against the JSON-encoded hashes array
            byte[] hashesBytes = hashesJsonBySignature.get(i);
            if (hashesBytes == null) {
                skippedReasons.add(String.format("failed to encode hashes for user '%s'", sigUserId));
                continue;
            }

            if (SignatureVerifier.verifySignature(hashesBytes, userSig.getSignature(),
                    user.getPublicKey())) {
                signers.add(SignatureVerifier.keyFingerprint(user.getPublicKey()));
                if (signers.size() >= minSigs) {
                    return; // Threshold met
                }
            } else {
                skippedReasons.add(String.format("user '%s' signature verification failed", sigUserId));
            }
        }

        // Threshold not met
        StringBuilder message = new StringBuilder();
        message.append(String.format("group '%s' requires %d distinct signer(s) but only %d valid",
                groupId, minSigs, signers.size()));
        if (!skippedReasons.isEmpty()) {
            message.append(" [").append(String.join("; ", skippedReasons)).append("]");
        }
        throw new IntegrityException(message.toString());
    }
}
