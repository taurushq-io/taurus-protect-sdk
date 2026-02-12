package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.helper.WhitelistHashHelper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.mapper.WhitelistedAddressMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.InternalWallet;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalWallet;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AddressWhitelistingApi;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedAddressEnvelopeReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedAddressEnvelopesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHashRulesContainer;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddressEnvelope;

import java.nio.charset.StandardCharsets;
import java.security.PublicKey;
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

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;
import static org.bouncycastle.util.Strings.constantTimeAreEqual;

/**
 * Service for retrieving and verifying whitelisted addresses.
 * Performs full cryptographic verification including:
 * - Metadata hash verification
 * - Rules container signature verification (SuperAdmin)
 * - Whitelist signature verification (per governance rules)
 */
public class WhitelistedAddressService {

    private static final Logger LOGGER = Logger.getLogger(WhitelistedAddressService.class.getName());
    private static final Gson GSON = new Gson();
    private static final Pattern CONTRACT_TYPE_PATTERN =
            Pattern.compile(",\"contractType\":\"[^\"]*\"");
    private static final Pattern LABEL_IN_OBJECT_PATTERN =
            Pattern.compile(",\"label\":\"[^\"]*\"}");

    private final AddressWhitelistingApi whitelistedAddressService;
    private final ApiExceptionMapper apiExceptionMapper;
    private final List<PublicKey> superAdminPublicKeys;
    private final int minValidSignatures;

    /**
     * Instantiates a new Whitelisted Address service.
     *
     * @param openApiClient        the open api client
     * @param apiExceptionMapper   the api exception mapper
     * @param superAdminPublicKeys the list of SuperAdmin public keys for rules verification
     * @param minValidSignatures   the minimum number of valid signatures required for rules
     */
    public WhitelistedAddressService(final ApiClient openApiClient,
                                     final ApiExceptionMapper apiExceptionMapper,
                                     final List<PublicKey> superAdminPublicKeys,
                                     final int minValidSignatures) {

        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");
        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys cannot be empty");
        checkArgument(minValidSignatures > 0, "minValidSignatures must be positive");

        this.apiExceptionMapper = apiExceptionMapper;
        this.whitelistedAddressService = new AddressWhitelistingApi(openApiClient);
        this.superAdminPublicKeys = superAdminPublicKeys;
        this.minValidSignatures = minValidSignatures;
    }

    /**
     * Gets a whitelisted address by ID with full signature verification.
     *
     * @param id the whitelisted address ID
     * @return the verified and decoded WhitelistedAddress
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if signature verification fails
     */
    public WhitelistedAddress getWhitelistedAddress(final long id) throws ApiException, WhitelistException {
        SignedWhitelistedAddressEnvelope envelope = getWhitelistedAddressEnvelope(id);
        return envelope.getWhitelistedAddress();
    }

    /**
     * Gets the signed whitelisted address envelope by ID.
     * Performs full verification including metadata hash, rules container signatures,
     * and whitelist signatures.
     *
     * @param id the whitelisted address ID
     * @return the signed whitelisted address envelope with verified WhitelistedAddress
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public SignedWhitelistedAddressEnvelope getWhitelistedAddressEnvelope(final long id)
            throws ApiException, WhitelistException {
        checkArgument(id > 0, "whitelisted address id cannot be zero");

        try {
            TgvalidatordGetSignedWhitelistedAddressEnvelopeReply reply =
                    whitelistedAddressService.whitelistServiceGetWhitelistedAddress(String.valueOf(id));
            SignedWhitelistedAddressEnvelope envelope =
                    WhitelistedAddressMapper.INSTANCE.fromDTO(reply.getResult());
            initializeEnvelope(envelope);
            return envelope;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Initializes and fully verifies the envelope.
     * After this method completes, the envelope's getWhitelistedAddress() will return
     * the verified address.
     */
    private void initializeEnvelope(SignedWhitelistedAddressEnvelope envelope) throws WhitelistException {
        // Precondition checks
        if (envelope.getSignedAddress() == null || envelope.getSignedAddress().getPayload() == null) {
            throw new WhitelistException("signed address payload is null");
        }
        if (envelope.getSignedAddress().getSignatures() == null
                || envelope.getSignedAddress().getSignatures().isEmpty()) {
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

        // Step 6: Parse WhitelistedAddress from verified payloadAsString (signed fields)
        WhitelistedAddress address = WhitelistHashHelper.parseWhitelistedAddressFromJson(
                envelope.getMetadata().getPayloadAsString());

        // Extract createdAt from trails (find "created" action)
        if (envelope.getTrails() != null) {
            for (WhitelistTrail trail : envelope.getTrails()) {
                if ("created".equals(trail.getAction())) {
                    address.setCreatedAt(trail.getDate());
                    break;
                }
            }
        }

        // Extract attributes from envelope
        if (envelope.getAttributes() != null) {
            Map<String, Object> attrs = new HashMap<>();
            for (Attribute attr : envelope.getAttributes()) {
                if (attr.getKey() != null) {
                    attrs.put(attr.getKey(), attr.getValue());
                }
            }
            address.setAttributes(attrs);
        }

        // Store verified data in envelope
        envelope.setVerifiedWhitelistedAddress(address);
        envelope.setVerifiedRulesContainer(rulesContainer);
    }

    /**
     * Verifies that the computed hash of payloadAsString equals the provided hash.
     */
    private void verifyMetadataHash(SignedWhitelistedAddressEnvelope envelope) {
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
                LOGGER.warning("Metadata hash verification failed for whitelisted address");
            }
            throw new IntegrityException("metadata hash verification failed");
        }
    }

    /**
     * Verifies the rulesContainer signatures against SuperAdmin public keys.
     */
    private void verifyRulesContainerSignatures(SignedWhitelistedAddressEnvelope envelope)
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
    private DecodedRulesContainer decodeRulesContainer(SignedWhitelistedAddressEnvelope envelope)
            throws WhitelistException {
        try {
            return RulesContainerMapper.INSTANCE.fromBase64String(envelope.getRulesContainer());
        } catch (IllegalArgumentException | InvalidProtocolBufferException e) {
            throw new WhitelistException("unable to decode rules container", e);
        }
    }

    /**
     * Verifies rules container signatures and decodes to DecodedRulesContainer.
     * Used for processing the rulesContainers array in normalized responses.
     *
     * @param rulesContainerBase64   the base64-encoded rules container
     * @param rulesSignaturesBase64  the base64-encoded rules signatures
     * @return the decoded and verified rules container
     * @throws WhitelistException if verification or decoding fails
     */
    private DecodedRulesContainer verifyAndDecodeRulesContainer(
            String rulesContainerBase64, String rulesSignaturesBase64) throws WhitelistException {

        if (Strings.isNullOrEmpty(rulesSignaturesBase64)) {
            throw new WhitelistException("rules signatures is null or empty");
        }
        if (Strings.isNullOrEmpty(rulesContainerBase64)) {
            throw new WhitelistException("rules container is null or empty");
        }

        // Decode signatures
        List<RuleUserSignature> signatures;
        try {
            signatures = RulesContainerMapper.INSTANCE.userSignaturesFromBase64String(
                    rulesSignaturesBase64);
        } catch (IllegalArgumentException | InvalidProtocolBufferException e) {
            throw new WhitelistException("unable to decode rules signatures", e);
        }

        if (signatures.isEmpty()) {
            throw new WhitelistException("no rules signatures present");
        }

        // Verify signatures
        byte[] rulesData = Base64.getDecoder().decode(rulesContainerBase64);
        int validCount = 0;
        for (RuleUserSignature sig : signatures) {
            if (Strings.isNullOrEmpty(sig.getSignature())) {
                continue;
            }
            if (SignatureVerifier.isValidSignature(rulesData, sig.getSignature(),
                    superAdminPublicKeys)) {
                validCount++;
            }
        }

        if (validCount < minValidSignatures) {
            if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                LOGGER.warning("Rules container verification failed: insufficient valid signatures");
            }
            throw new IntegrityException(String.format(
                    "Rules container verification failed: only %d valid signatures found, "
                            + "minimum %d required", validCount, minValidSignatures));
        }
        if (LOGGER.isLoggable(java.util.logging.Level.FINE)) {
            LOGGER.fine("Rules container signature verification succeeded");
        }

        // Decode rules container
        try {
            return RulesContainerMapper.INSTANCE.fromBase64String(rulesContainerBase64);
        } catch (IllegalArgumentException | InvalidProtocolBufferException e) {
            throw new WhitelistException("unable to decode rules container", e);
        }
    }

    /**
     * Verifies that the metadata hash is present in at least one signature's hashes list.
     * For backward compatibility, also tries alternative hashes for addresses signed
     * before certain fields (like contractType, labels in linkedInternalAddresses) were added.
     */
    private void verifyHashInSignedHashes(SignedWhitelistedAddressEnvelope envelope)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();
        List<WhitelistSignature> signatures = envelope.getSignedAddress().getSignatures();

        // First, try the provided hash directly
        if (hashExistsInSignatures(metadataHash, signatures)) {
            return; // Found - verification passed
        }

        // If not found, try alternative hashes for backward compatibility
        // (handles addresses signed before schema changes)
        for (String legacyHash : computeLegacyHashes(envelope.getMetadata().getPayloadAsString())) {
            if (hashExistsInSignatures(legacyHash, signatures)) {
                // Update the metadata hash so subsequent verification steps use the correct hash
                envelope.getMetadata().setHash(legacyHash);
                return; // Found with legacy hash
            }
        }

        if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
            // SECURITY: Do not log hash values to prevent information leakage
            LOGGER.warning("Metadata hash not found in any signature's hashes list");
        }
        throw new IntegrityException("metadata hash not found in any signature's hashes list");
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
     * Computes legacy hashes by applying ALL transformation combinations to handle schema evolution.
     * Returns a list of possible legacy hashes to try.
     *
     * <p>Strategies cover all possible schema evolution scenarios:
     * <ul>
     *   <li>Strategy 1: contractType added after signing</li>
     *   <li>Strategy 2: labels added to linkedInternalAddresses after signing (but contractType existed)</li>
     *   <li>Strategy 3: both contractType AND labels added after signing</li>
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

        // Strategy 1: Remove contractType only
        // Handles addresses signed before contractType was added to schema
        String withoutContractType = CONTRACT_TYPE_PATTERN.matcher(payloadAsString).replaceAll("");
        if (!withoutContractType.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutContractType));
        }

        // Strategy 2: Remove labels from linkedInternalAddresses objects only (keep contractType)
        // Handles addresses signed after contractType was added but before labels were added
        // Pattern ,"label":"[^"]*"} matches ONLY labels inside objects (followed by closing brace)
        // This does NOT match the main address label which is followed by ,"customerId":
        String withoutLabels = LABEL_IN_OBJECT_PATTERN.matcher(payloadAsString).replaceAll("}");
        if (!withoutLabels.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutLabels));
        }

        // Strategy 3: Remove BOTH contractType AND labels from linkedInternalAddresses
        // Handles addresses signed before both fields were added
        String withoutBoth = LABEL_IN_OBJECT_PATTERN.matcher(payloadAsString).replaceAll("}");
        withoutBoth = CONTRACT_TYPE_PATTERN.matcher(withoutBoth).replaceAll("");
        if (!withoutBoth.equals(payloadAsString)) {
            uniqueHashes.add(CryptoTPV1.calculateHexHash(withoutBoth));
        }

        return new ArrayList<>(uniqueHashes);
    }

    /**
     * Verifies whitelist signatures according to governance rules threshold requirements.
     */
    private void verifyWhitelistSignatures(SignedWhitelistedAddressEnvelope envelope,
                                           DecodedRulesContainer rulesContainer)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();

        // Find matching address whitelisting rules
        AddressWhitelistingRules whitelistRules = rulesContainer.findAddressWhitelistingRules(
                envelope.getBlockchain(), envelope.getNetwork());
        if (whitelistRules == null) {
            throw new WhitelistException("no address whitelisting rules found for blockchain="
                    + envelope.getBlockchain() + " network=" + envelope.getNetwork());
        }

        // Parse the whitelisted address to check linked addresses/wallets
        WhitelistedAddress wla = WhitelistHashHelper.parseWhitelistedAddressFromJson(
                envelope.getMetadata().getPayloadAsString());

        // Determine which thresholds to use based on rule lines matching
        List<SequentialThresholds> parallelThresholds = getApplicableThresholds(whitelistRules, wla);
        if (parallelThresholds == null || parallelThresholds.isEmpty()) {
            throw new WhitelistException("no threshold rules defined");
        }

        // Try to verify all paths
        List<String> pathFailures = tryVerifyAllPaths(
                parallelThresholds, rulesContainer, envelope.getSignedAddress().getSignatures(),
                metadataHash);
        if (!pathFailures.isEmpty()) {
            throw new WhitelistException("signature verification failed of whitelisted address (ID: " + envelope.getId() + ") : "
                    + "no approval path satisfied the threshold requirements. "
                    + String.join("; ", pathFailures));
        }
    }

    /**
     * Determines which thresholds to use based on rule lines matching.
     * Checks rule lines only when: NO linked addresses AND exactly 1 linked wallet.
     * Otherwise falls back to default thresholds.
     */
    private List<SequentialThresholds> getApplicableThresholds(
            AddressWhitelistingRules rules, WhitelistedAddress wla) {

        boolean hasLinkedAddresses = wla.getLinkedInternalAddresses() != null
                && !wla.getLinkedInternalAddresses().isEmpty();
        List<InternalWallet> linkedWallets = wla.getLinkedWallets();
        int walletCount = (linkedWallets == null) ? 0 : linkedWallets.size();

        // Check rule lines only if: no linked addresses AND exactly 1 linked wallet
        boolean shouldCheckRuleLines = !hasLinkedAddresses && walletCount == 1;

        if (shouldCheckRuleLines && rules.getLines() != null && !rules.getLines().isEmpty()) {
            String walletPath = linkedWallets.get(0).getPath();

            // Find matching line by wallet path
            for (AddressWhitelistingLine line : rules.getLines()) {
                if (matchesWalletPath(line, walletPath)) {
                    return line.getParallelThresholds();
                }
            }
        }

        // Fallback to default thresholds
        return rules.getParallelThresholds();
    }

    /**
     * Checks if a rule line matches the given wallet path.
     */
    private boolean matchesWalletPath(AddressWhitelistingLine line, String walletPath) {
        if (line.getCells() == null || line.getCells().isEmpty()) {
            return false;
        }

        RuleSource source = line.getCells().get(0);
        if (source.getType() != RuleSourceType.RuleSourceInternalWallet) {
            return false;  // Only support internal wallet for now
        }

        RuleSourceInternalWallet internalWallet = source.getInternalWallet();
        return internalWallet != null
                && walletPath != null
                && walletPath.equals(internalWallet.getPath());
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

    /**
     * Gets a list of whitelisted address envelopes with pagination.
     * Uses normalized rules containers by default for better performance.
     *
     * @param limit  the maximum number of results (max 100)
     * @param offset the offset for pagination
     * @return the list of signed whitelisted address envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(
            int limit, int offset) throws ApiException, WhitelistException {
        return getWhitelistedAddresses(limit, offset, null, null, true);
    }

    /**
     * Gets a list of whitelisted address envelopes filtered by blockchain.
     * Uses normalized rules containers by default for better performance.
     *
     * @param limit      the maximum number of results (max 100)
     * @param offset     the offset for pagination
     * @param blockchain filter by blockchain (e.g., "ETH", "BTC")
     * @return the list of signed whitelisted address envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(
            int limit, int offset, String blockchain) throws ApiException, WhitelistException {
        return getWhitelistedAddresses(limit, offset, blockchain, null, true);
    }

    /**
     * Gets a list of whitelisted address envelopes filtered by blockchain and network.
     * Uses normalized rules containers by default for better performance.
     *
     * @param limit      the maximum number of results (max 100)
     * @param offset     the offset for pagination
     * @param blockchain filter by blockchain (e.g., "ETH", "BTC")
     * @param network    filter by network (e.g., "mainnet", "testnet")
     * @return the list of signed whitelisted address envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(
            int limit, int offset, String blockchain, String network)
            throws ApiException, WhitelistException {
        return getWhitelistedAddresses(limit, offset, blockchain, network, true);
    }

    /**
     * Gets a list of whitelisted address envelopes filtered by blockchain and network.
     *
     * @param limit                      the maximum number of results (max 100)
     * @param offset                     the offset for pagination
     * @param blockchain                 filter by blockchain (e.g., "ETH", "BTC")
     * @param network                    filter by network (e.g., "mainnet", "testnet")
     * @param rulesContainerNormalized   if true, caches rules containers by hash to avoid
     *                                   redundant verification of identical containers
     * @return the list of signed whitelisted address envelopes
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public List<SignedWhitelistedAddressEnvelope> getWhitelistedAddresses(
            int limit, int offset, String blockchain, String network,
            boolean rulesContainerNormalized) throws ApiException, WhitelistException {
        try {
            TgvalidatordGetSignedWhitelistedAddressEnvelopesReply reply =
                    whitelistedAddressService.whitelistServiceGetWhitelistedAddresses(
                            String.valueOf(limit),
                            String.valueOf(offset),
                            null,       // exchangeAccountId
                            null,       // addressType
                            null,       // query
                            null,       // currency (deprecated)
                            null,       // scoreProvider
                            null,       // scoreInBelow
                            null,       // scoreOutBelow
                            null,       // scoreExclusive
                            rulesContainerNormalized,  // rulesContainerNormalized
                            null,       // exchangeAccountIds
                            null,       // coinfirmScoreGreater
                            null,       // tagIDs
                            null,       // chainalysisScoreGreater
                            null,       // contractType
                            null,       // allowedForAddressId
                            null,       // allowedForWalletId
                            blockchain, // blockchain
                            null,       // includeForApproval
                            null,       // addresses
                            network,    // network
                            null,       // ids
                            null,       // tnParticipantID
                            null, null, null, null, null, null, null, null,
                            null, null, null, null);

            if (reply.getResult() == null) {
                return new ArrayList<>();
            }

            List<SignedWhitelistedAddressEnvelope> envelopes = new ArrayList<>();
            Map<String, DecodedRulesContainer> rulesContainerCache = new HashMap<>();

            // Check if rulesContainers array exists - if so, pre-populate cache
            List<TgvalidatordHashRulesContainer> rulesContainers = reply.getRulesContainers();
            if (rulesContainers != null && !rulesContainers.isEmpty()) {
                // Deduplicate by base64 container string to avoid re-verifying identical containers
                Map<String, DecodedRulesContainer> verifiedContainers = new HashMap<>();

                for (TgvalidatordHashRulesContainer hashContainer : rulesContainers) {
                    if (hashContainer.getHash() == null
                            || hashContainer.getRulesContainer() == null) {
                        continue;
                    }
                    String containerBase64 = hashContainer.getRulesContainer();
                    DecodedRulesContainer decoded = verifiedContainers.get(containerBase64);
                    if (decoded == null) {
                        decoded = verifyAndDecodeRulesContainer(
                                containerBase64,
                                hashContainer.getRulesSignatures());
                        verifiedContainers.put(containerBase64, decoded);
                    }
                    rulesContainerCache.put(hashContainer.getHash(), decoded);
                }
            }

            // Process all envelopes
            for (TgvalidatordSignedWhitelistedAddressEnvelope dto : reply.getResult()) {
                SignedWhitelistedAddressEnvelope envelope =
                        WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

                // Check if we have a cached rules container for this envelope
                DecodedRulesContainer cached = null;
                if (envelope.getRulesContainerHash() != null) {
                    cached = rulesContainerCache.get(envelope.getRulesContainerHash());
                }

                if (cached != null) {
                    // Use cached rules container
                    initializeEnvelopeWithCachedRules(envelope, cached);
                } else {
                    // Full data mode: envelope has its own rulesContainer
                    initializeEnvelope(envelope);
                }
                envelopes.add(envelope);
            }
            return envelopes;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Initializes and verifies the envelope using a pre-verified cached rules container.
     * This skips rules container signature verification since it was already done.
     */
    private void initializeEnvelopeWithCachedRules(SignedWhitelistedAddressEnvelope envelope,
                                                    DecodedRulesContainer cachedRulesContainer)
            throws WhitelistException {
        // Precondition checks
        if (envelope.getSignedAddress() == null || envelope.getSignedAddress().getPayload() == null) {
            throw new WhitelistException("signed address payload is null");
        }
        if (envelope.getSignedAddress().getSignatures() == null
                || envelope.getSignedAddress().getSignatures().isEmpty()) {
            throw new WhitelistException("signatures is null or empty");
        }

        // Step 1: Verify computed hash of payloadAsString equals received hash
        verifyMetadataHash(envelope);

        // Step 2: Skip rulesContainer signature verification - already done for cached container

        // Step 3: Use cached rulesContainer instead of decoding again

        // Step 4: Verify metadata.hash is in signed hashes list
        verifyHashInSignedHashes(envelope);

        // Step 5: Verify whitelist signatures are valid per governance rules
        verifyWhitelistSignatures(envelope, cachedRulesContainer);

        // Step 6: Parse WhitelistedAddress from verified payloadAsString (signed fields)
        WhitelistedAddress address = WhitelistHashHelper.parseWhitelistedAddressFromJson(
                envelope.getMetadata().getPayloadAsString());

        // Extract createdAt from trails (find "created" action)
        if (envelope.getTrails() != null) {
            for (WhitelistTrail trail : envelope.getTrails()) {
                if ("created".equals(trail.getAction())) {
                    address.setCreatedAt(trail.getDate());
                    break;
                }
            }
        }

        // Extract attributes from envelope
        if (envelope.getAttributes() != null) {
            Map<String, Object> attrs = new HashMap<>();
            for (Attribute attr : envelope.getAttributes()) {
                if (attr.getKey() != null) {
                    attrs.put(attr.getKey(), attr.getValue());
                }
            }
            address.setAttributes(attrs);
        }

        // Store verified data in envelope
        envelope.setVerifiedWhitelistedAddress(address);
        envelope.setVerifiedRulesContainer(cachedRulesContainer);
    }
}
