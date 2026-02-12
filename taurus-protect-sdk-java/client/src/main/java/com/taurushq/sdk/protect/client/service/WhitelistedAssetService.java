package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.helper.AssetHashHelper;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
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
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ContractWhitelistingApi;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;

import java.nio.charset.StandardCharsets;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;
import java.util.logging.Logger;
import java.util.regex.Pattern;

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
        try {
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply =
                    contractWhitelistingApi.whitelistServiceGetWhitelistedContracts(
                            String.valueOf(limit),
                            String.valueOf(offset),
                            null,       // query
                            blockchain, // blockchain
                            null,       // includeForApproval
                            network,    // network
                            null,       // isNFT
                            null,       // whitelistedContractAddressIds
                            null);      // kindTypes

            if (reply.getResult() == null) {
                return new ArrayList<>();
            }

            List<SignedWhitelistedAssetEnvelope> envelopes = new ArrayList<>();
            for (TgvalidatordSignedWhitelistedContractAddressEnvelope dto : reply.getResult()) {
                SignedWhitelistedAssetEnvelope envelope =
                        WhitelistedAssetMapper.INSTANCE.fromDTO(dto);
                initializeEnvelope(envelope);
                envelopes.add(envelope);
            }
            return envelopes;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Initializes and fully verifies the envelope.
     * After this method completes, the envelope's getWhitelistedAsset() will return
     * the verified asset.
     */
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
        verifyHashInSignedHashes(envelope);

        // Step 5: Verify whitelist signatures are valid per governance rules
        verifyWhitelistSignatures(envelope, rulesContainer);

        // Step 6: Parse WhitelistedAsset from verified payloadAsString (signed fields)
        WhitelistedAsset asset = AssetHashHelper.parseWhitelistedAssetFromJson(
                envelope.getMetadata().getPayloadAsString());

        // Store verified data in envelope
        envelope.setVerifiedWhitelistedAsset(asset);
        envelope.setVerifiedRulesContainer(rulesContainer);
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
        int validCount = 0;

        for (RuleUserSignature sig : signatures) {
            if (Strings.isNullOrEmpty(sig.getSignature())) {
                continue;
            }
            if (SignatureVerifier.isValidSignature(rulesData, sig.getSignature(), superAdminPublicKeys)) {
                validCount++;
            }
        }

        if (validCount < minValidSignatures) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                LOGGER.warning("Rules container verification failed: insufficient valid signatures");
            }
            throw new IntegrityException(String.format(
                    "Rules container verification failed: only %d valid signatures found, minimum %d required",
                    validCount, minValidSignatures));
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
    private void verifyHashInSignedHashes(SignedWhitelistedAssetEnvelope envelope)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();
        List<WhitelistSignature> signatures = envelope.getSignedAsset().getSignatures();

        // First, try the provided hash directly
        if (hashExistsInSignatures(metadataHash, signatures)) {
            return; // Found - verification passed
        }

        // If not found, try alternative hashes for backward compatibility
        // (handles assets signed before schema changes)
        for (String legacyHash : computeLegacyHashes(envelope.getMetadata().getPayloadAsString())) {
            if (hashExistsInSignatures(legacyHash, signatures)) {
                // Update the metadata hash so subsequent verification steps use the correct hash
                envelope.getMetadata().setHash(legacyHash);
                return; // Found with legacy hash
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

    private boolean hashExistsInSignatures(String hash, List<WhitelistSignature> signatures) {
        for (WhitelistSignature sig : signatures) {
            if (sig.getHashes() != null && sig.getHashes().contains(hash)) {
                return true;
            }
        }
        return false;
    }

    /**
     * Verifies whitelist signatures according to governance rules threshold requirements.
     */
    private void verifyWhitelistSignatures(SignedWhitelistedAssetEnvelope envelope,
                                           DecodedRulesContainer rulesContainer)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();

        // Find matching contract address whitelisting rules
        ContractAddressWhitelistingRules whitelistRules = rulesContainer.findContractAddressWhitelistingRules(
                envelope.getBlockchain(), envelope.getNetwork());
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
                metadataHash);
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
                                           String metadataHash) {
        List<String> pathFailures = new ArrayList<>();
        for (int i = 0; i < parallelThresholds.size(); i++) {
            SequentialThresholds seqThreshold = parallelThresholds.get(i);
            try {
                verifySequentialThresholds(seqThreshold, rulesContainer, signatures, metadataHash);
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
                                            String metadataHash) {
        List<GroupThreshold> thresholds = seqThreshold.getThresholds();
        if (thresholds == null || thresholds.isEmpty()) {
            throw new IntegrityException("no group thresholds defined");
        }

        // ALL group thresholds must be satisfied (AND logic)
        for (GroupThreshold groupThreshold : thresholds) {
            verifyGroupThreshold(groupThreshold, rulesContainer, signatures, metadataHash);
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
                                      String metadataHash) {
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

        // Convert to set for faster lookup
        Set<String> groupUserIdSet = new HashSet<>(groupUserIds);

        // Count valid signatures from users in this group
        int validCount = 0;
        List<String> skippedReasons = new ArrayList<>();

        for (WhitelistSignature sig : signatures) {
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
            if (hashes == null || !hashes.contains(metadataHash)) {
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

            // Verify signature against JSON-encoded hashes array
            String hashesJson = GSON.toJson(hashes);
            byte[] hashesBytes = hashesJson.getBytes(StandardCharsets.UTF_8);

            if (SignatureVerifier.verifySignature(hashesBytes, userSig.getSignature(),
                    user.getPublicKey())) {
                validCount++;
                if (validCount >= minSigs) {
                    return; // Threshold met
                }
            } else {
                skippedReasons.add(String.format("user '%s' signature verification failed", sigUserId));
            }
        }

        // Threshold not met
        StringBuilder message = new StringBuilder();
        message.append(String.format("group '%s' requires %d signature(s) but only %d valid",
                groupId, minSigs, validCount));
        if (!skippedReasons.isEmpty()) {
            message.append(" [").append(String.join("; ", skippedReasons)).append("]");
        }
        throw new IntegrityException(message.toString());
    }
}
