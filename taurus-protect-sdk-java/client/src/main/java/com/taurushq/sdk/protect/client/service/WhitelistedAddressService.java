package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveWhitelistedAddressRequest;
import java.util.stream.Collectors;
import java.security.SignatureException;
import java.security.PrivateKey;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.InvalidKeyException;
import com.google.common.base.Strings;
import com.google.gson.Gson;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.helper.WhitelistHashHelper;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.mapper.WhitelistedAddressMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ContainerIntegrityException;
import com.taurushq.sdk.protect.client.model.ExcludedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.WhitelistedAddressListResult;
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

        // Step 4: Verify metadata.hash is in signed hashes list. Returns the hash
        // that was actually covered, which may be a legacy one.
        String verifiedHash = verifyHashInSignedHashes(envelope);

        // Step 5: Verify whitelist signatures are valid per governance rules
        verifyWhitelistSignatures(envelope, rulesContainer, verifiedHash);

        // Step 6: the envelope DERIVES the address from its own signed payloadAsString,
        // plus the non-security trail/attribute fields it already carries. It used to be
        // parsed here and handed to a public setter, which meant the "verified" marker
        // could be flipped with an address the caller chose.
        envelope.markVerified(rulesContainer);
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
    private String verifyHashInSignedHashes(SignedWhitelistedAddressEnvelope envelope)
            throws WhitelistException {
        String metadataHash = envelope.getMetadata().getHash();
        List<WhitelistSignature> signatures = envelope.getSignedAddress().getSignatures();

        // First, try the provided hash directly
        if (SignatureVerifier.verifyHashCoverage(metadataHash, signatures)) {
            return metadataHash;
        }

        // If not found, try alternative hashes for backward compatibility
        // (handles addresses signed before schema changes)
        for (String legacyHash : computeLegacyHashes(envelope.getMetadata().getPayloadAsString())) {
            if (SignatureVerifier.verifyHashCoverage(legacyHash, signatures)) {
                // Returned rather than written back onto the caller's envelope:
                // verification must not mutate its input, and the later steps take
                // the hash as a parameter.
                return legacyHash;
            }
        }

        if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
            // SECURITY: Do not log hash values to prevent information leakage
            LOGGER.warning("Metadata hash not found in any signature's hashes list");
        }
        throw new IntegrityException("metadata hash not found in any signature's hashes list");
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
                                           DecodedRulesContainer rulesContainer,
                                           String metadataHash)
            throws WhitelistException {

        // Which rules judge this address is decided by the SIGNED payload, not by the
        // surrounding response. A DTO with an empty blockchain would select the
        // global-default tier — broader than the rule the address belongs to.
        String[] ruleKey = WhitelistHashHelper.resolveRuleKey(
                envelope.getMetadata().getPayloadAsString(),
                envelope.getBlockchain(), envelope.getNetwork());

        // Find matching address whitelisting rules
        AddressWhitelistingRules whitelistRules = rulesContainer.findAddressWhitelistingRules(
                ruleKey[0], ruleKey[1]);
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
            for (int i = 0; i < rules.getLines().size(); i++) {
                AddressWhitelistingLine line = rules.getLines().get(i);
                // A source cell this SDK could not type might be the one that matches.
                // Falling through to the container defaults would verify the address
                // against a quorum governance never granted it - silently, no error.
                if (lineHasUntypedSource(line)) {
                    throw new ContainerIntegrityException(
                            "address whitelisting rules for blockchain=" + rules.getCurrency()
                                    + " network=" + rules.getNetwork() + " line " + i
                                    + " carry a source cell this SDK version cannot interpret;"
                                    + " refusing to fall back to the container default thresholds,"
                                    + " which may be weaker than the line's. Upgrade the SDK to"
                                    + " match the validatord that signed this container");
                }
                if (matchesWalletPath(line, walletPath)) {
                    return line.getParallelThresholds();
                }
            }
        }

        // Fallback to default thresholds
        return rules.getParallelThresholds();
    }

    /**
     * Returns true when the line's source cell was preserved verbatim rather than typed.
     * {@code matchesWalletPath} reads cell 0, so that is the cell whose meaning must be
     * known before a match/no-match verdict can be trusted.
     * <p>
     * Package-private so it can be tested directly: this project forbids Mockito, so the
     * enclosing private getApplicableThresholds cannot be driven from a service test.
     */
    boolean lineHasUntypedSource(AddressWhitelistingLine line) {
        if (line == null || line.getCells() == null || line.getCells().isEmpty()) {
            return false;
        }
        RuleSource source = line.getCells().get(0);
        return source != null && source.getRaw() != null && !source.getRaw().isEmpty();
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
    // Package-private, not private: the distinct-signer counting below is a security
    // invariant and this is the only way to drive it without a network stub, which this
    // module has none of.
    void verifyGroupThreshold(GroupThreshold groupThreshold,
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

            // Verify signature against JSON-encoded hashes array
            String hashesJson = GSON.toJson(hashes);
            byte[] hashesBytes = hashesJson.getBytes(StandardCharsets.UTF_8);

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
        return getWhitelistedAddressesWithExclusions(
                limit, offset, blockchain, network, rulesContainerNormalized).getEnvelopes();
    }

    /**
     * Same query as {@link #getWhitelistedAddresses(int, int, String, String, boolean)},
     * but also returns the rows that were dropped for failing verification.
     * <p>
     * An overload rather than a changed return type, so existing callers keep compiling.
     * Prefer this one: excluded rows are otherwise only written to the SDK's logger,
     * which a caller cannot read, and a filtered page is then indistinguishable from a
     * complete one.
     *
     * @param limit                    the maximum number of results (max 100)
     * @param offset                   the offset for pagination
     * @param blockchain               filter by blockchain (e.g., "ETH", "BTC")
     * @param network                  filter by network (e.g., "mainnet", "testnet")
     * @param rulesContainerNormalized if true, caches rules containers by hash
     * @return the verified envelopes plus the rows excluded for failing verification
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if verification fails
     */
    public WhitelistedAddressListResult getWhitelistedAddressesWithExclusions(
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
                return new WhitelistedAddressListResult(new ArrayList<>(), new ArrayList<>());
            }

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

                    // The hash is the LABEL a row uses to pick its container, and it
                    // arrives in the same response as the container. Recompute it:
                    // otherwise a server can file container A under container B's label
                    // and steer any row to any other validly-signed container -- an older
                    // ruleset with a weaker group threshold, say. Both pass the SuperAdmin
                    // check, so signature verification alone does not catch it.
                    String computedLabel = containerHashLabel(containerBase64);
                    if (!computedLabel.equals(hashContainer.getHash())) {
                        throw new ContainerIntegrityException(String.format(
                                "rules container hash mismatch: response labelled it %s "
                                        + "but its bytes hash to %s",
                                hashContainer.getHash(), computedLabel));
                    }
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

            // rows -> verify -+- ok    -> returned
            //                  +- fails -> excluded + logged, list survives
            //
            // One unverifiable row used to fail the whole call, which took down the
            // whitelist for every consumer -- and listing is how an operator would
            // find the bad row, so the failure hid its own cause. Excluding stays
            // fail-closed: an omitted destination cannot be selected.
            return verifiedAddresses(
                    reply.getResult(), rulesContainerCache, reply.getTotalItems());

        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets a page of whitelisted addresses awaiting approval, verified exactly as
     * {@link #getWhitelistedAddressesWithExclusions} is.
     *
     * <p>Both this endpoint and {@link #approveWhitelistedAddresses} are generated in
     * all four SDK clients and were wrapped by none of them, so the rows an approver
     * reads before whitelisting a destination were not reachable through the SDK at all
     * — verified or not.
     *
     * @param limit                      the maximum number of results
     * @param offset                     the offset for pagination
     * @param ids                        filter by specific ids, or null
     * @param includeAlreadySignedByUser include rows the caller has already signed,
     *                                   which is how an approver tells "waiting for me"
     *                                   from "waiting for someone else"
     * @return the verified rows and the rows withheld
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if no row survived verification
     */
    public WhitelistedAddressListResult getWhitelistedAddressesForApproval(
            final int limit, final int offset, final List<String> ids,
            final Boolean includeAlreadySignedByUser)
            throws ApiException, WhitelistException {
        try {
            TgvalidatordGetSignedWhitelistedAddressEnvelopesReply reply =
                    whitelistedAddressService.whitelistServiceGetWhitelistedAddressesForApproval(
                            String.valueOf(limit),
                            String.valueOf(offset),
                            ids,
                            null,       // blockchain
                            null,       // addressType
                            null,       // query
                            null,       // network
                            includeAlreadySignedByUser);

            if (reply.getResult() == null) {
                return new WhitelistedAddressListResult(new ArrayList<>(), new ArrayList<>());
            }

            // This endpoint has no normalized-container mode, so the per-row containers
            // are used and the cache starts empty.
            return verifiedAddresses(
                    reply.getResult(), new HashMap<>(), reply.getTotalItems());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Signs and submits an approval for the given whitelisted addresses, all-or-nothing.
     *
     * <p>The batch is re-read through the verifying path and the hashes THOSE rows carry
     * are what gets signed, so the approver's signature covers metadata this SDK checked
     * rather than whatever a caller was handed. Same shape as
     * {@code WhitelistedAssetService.approveWhitelistedAssets}.
     *
     * <pre>
     *   ids -&gt; sort -&gt; ONE filtered verified page -&gt; completeness check -&gt; sign once
     * </pre>
     *
     * <p>Any address that is missing or fails verification aborts the whole call and
     * nothing is signed: one signature covers every hash in the batch, so a partial
     * approval would mean the caller believes they approved more than they did.
     *
     * @param ids        the whitelisted address ids to approve
     * @param privateKey the approver's P-256 private key
     * @param comment    the approval comment
     * @throws ApiException       if the API call fails
     * @throws WhitelistException if any row is missing or fails verification
     */
    public void approveWhitelistedAddresses(final List<Long> ids, final PrivateKey privateKey,
                                            final String comment)
            throws ApiException, WhitelistException {
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");
        checkNotNull(privateKey, "privateKey cannot be null");
        checkArgument(!Strings.isNullOrEmpty(comment), "comment is required");
        ids.forEach(id -> checkArgument(id != null && id > 0,
                "whitelisted address id cannot be zero or negative"));

        // Sorted on a COPY: the endpoint requires ascending order, and this also makes
        // the signed order independent of the order the caller passed without mutating
        // their list.
        List<Long> sorted = new ArrayList<>(ids);
        Collections.sort(sorted);
        List<String> idStrings =
                sorted.stream().map(String::valueOf).collect(Collectors.toList());

        // ONE id-filtered page through the verifying path, not one GET per id.
        Map<String, SignedWhitelistedAddressEnvelope> byId = new HashMap<>();
        for (SignedWhitelistedAddressEnvelope envelope
                : getWhitelistedAddressesForApproval(idStrings.size(), 0, idStrings, Boolean.TRUE)
                        .getEnvelopes()) {
            if (envelope != null) {
                byId.put(String.valueOf(envelope.getId()), envelope);
            }
        }

        List<String> hashes = new ArrayList<>(idStrings.size());
        for (String id : idStrings) {
            SignedWhitelistedAddressEnvelope envelope = byId.get(id);
            if (envelope == null) {
                // A page that silently omits a row must not become an approval of fewer
                // rows than the caller asked for.
                throw new IntegrityException(String.format(
                        "refusing to sign: address %s was not returned by the verified read", id));
            }
            if (envelope.getMetadata() == null
                    || Strings.isNullOrEmpty(envelope.getMetadata().getHash())) {
                throw new IntegrityException(String.format(
                        "refusing to sign: address %s has no metadata hash", id));
            }
            hashes.add(envelope.getMetadata().getHash());
        }

        String toSign = GSON.toJson(hashes);
        TgvalidatordApproveWhitelistedAddressRequest request =
                new TgvalidatordApproveWhitelistedAddressRequest();
        request.setIds(idStrings);
        request.setComment(comment);
        try {
            request.setSignature(CryptoTPV1.calculateBase64Signature(
                    privateKey, toSign.getBytes(StandardCharsets.UTF_8)));
        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException ex) {
            ApiException e = new ApiException();
            e.setCode(400);
            e.setError("ClientInvalidRequest");
            e.setMessage("unable to sign the array of whitelisted address hashes");
            e.setOriginalException(ex);
            throw e;
        }

        try {
            whitelistedAddressService.whitelistServiceApproveWhitelistedAddress(request);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Verifies a page of address envelopes and returns the survivors plus the rows
     * withheld.
     *
     * <pre>
     *   rows -&gt; verify -+- ok    -&gt; returned
     *                   +- fails -&gt; excluded + logged, list survives
     * </pre>
     *
     * <p>LENIENT per row, but rows-returned-but-none-surviving aborts the whole call: it
     * is systemic, and an empty list would look identical to an empty whitelist.
     *
     * <p>Shared by {@code getWhitelistedAddressesWithExclusions} and
     * {@code getWhitelistedAddressesForApproval}. The for-approval endpoint had no
     * verified reader at all, and duplicating this loop into a second one is precisely
     * how the list paths drifted from {@code getWhitelistedAddress} in the first place.
     *
     * <p>Package-private so the exclusion behaviour can be tested directly.
     *
     * @param rows               the DTO rows
     * @param rulesContainerCache pre-verified containers by hash, may be empty
     * @return the verified envelopes and the withheld rows
     * @throws WhitelistException if no row survived
     */
    WhitelistedAddressListResult verifiedAddresses(
            final List<TgvalidatordSignedWhitelistedAddressEnvelope> rows,
            final Map<String, DecodedRulesContainer> rulesContainerCache,
            final String serverTotalItems)
            throws WhitelistException {
        List<SignedWhitelistedAddressEnvelope> envelopes = new ArrayList<>();
        int rowCount = rows.size();
        String firstFailure = null;
        List<ExcludedWhitelistedAddress> excluded = new ArrayList<>();

        for (TgvalidatordSignedWhitelistedAddressEnvelope dto : rows) {
            try {
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
            } catch (ContainerIntegrityException e) {
                // Not a property of this row: this SDK cannot interpret the rules
                // container, which invalidates EVERY row judged against it. Excluding
                // them one by one would empty the whitelist and report success.
                // Must precede the multi-catch below — it is an IntegrityException.
                throw e;
            } catch (WhitelistException | IntegrityException e) {
                // IntegrityException is UNCHECKED (extends SecurityException), so
                // catching only WhitelistException here let steps 1, 2 and 4 escape the
                // loop and abort the whole page: one row with a tampered payload denied
                // the entire whitelist listing, which is the failure this class's javadoc
                // says was fixed. Anything that is NOT one of these two is a defect in
                // this SDK and still propagates, rather than being reported to the caller
                // as "this address failed verification".
                if (firstFailure == null) {
                    firstFailure = e.getMessage();
                }
                excluded.add(new ExcludedWhitelistedAddress(dto.getId(), e.getMessage()));
                if (LOGGER.isLoggable(java.util.logging.Level.WARNING)) {
                    LOGGER.warning("whitelisted address excluded: verification failed (id="
                            + dto.getId() + ")");
                }
            }
        }

        // Rows came back but none survived: a systemic failure, not an empty
        // whitelist, and returning an empty list would look identical to one.
        if (rowCount > 0 && envelopes.isEmpty()) {
            throw new IntegrityException(String.format(
                    "all %d whitelisted address(es) failed verification; first failure: %s",
                    rowCount, firstFailure == null ? "unknown" : firstFailure));
        }
        return new WhitelistedAddressListResult(
                envelopes, excluded, reduceTotal(serverTotalItems, excluded.size()));
    }

    /**
     * Recomputes the label validatord files a normalized rules container under, so a
     * row's {@code rulesContainerHash} resolves only to bytes that really hash to it.
     *
     * <p>The convention is validatord's, in its whitelist controller:
     * {@code base64.StdEncoding.EncodeToString(crypto.Sha256([]byte(e.GetRulesContainer())))}.
     * Note what is hashed: {@code GetRulesContainer()} is ALREADY a base64 string there,
     * so the digest is over the base64 TEXT, not over the decoded protobuf, and the
     * output is base64 rather than hex. Decoding the container first and hashing the
     * protobuf yields a different value and would reject every container, breaking all
     * list calls.
     *
     * <p>Do not confuse it with {@code enforcedRulesHash}, which is
     * {@code base64(SHA256(raw protobuf))} and is a backlink to a ruleset's predecessor
     * rather than its own identity. Both are 44-character base64 SHA-256 digests, so
     * mixing them up is silent.
     *
     * @param containerBase64 the container exactly as the response carried it
     * @return the label those bytes should be filed under
     */
    static String containerHashLabel(final String containerBase64) {
        try {
            // CryptoTPV1 only exposes calculateHexHash; this label is base64 of the raw
            // digest, so the digest is taken directly rather than re-decoding hex.
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] sum = digest.digest(containerBase64.getBytes(StandardCharsets.UTF_8));
            return Base64.getEncoder().encodeToString(sum);
        } catch (NoSuchAlgorithmException e) {
            // SHA-256 is mandatory in every JRE; if it is absent nothing here can work.
            throw new IllegalStateException("SHA-256 is unavailable", e);
        }
    }

    /**
     * Reduces the server's page total by the number of rows withheld.
     *
     * <p>The server counts rows it returned; the caller receives only those that
     * verified, so reporting the raw total makes pagination promise rows that can never
     * be read. Floored at zero, and a total the server did not send or that does not
     * parse stays absent rather than becoming a guess. The wire type is a uint64 the
     * generated client surfaces as a String, hence the re-stringify — Go and TypeScript
     * do the same arithmetic on a numeric field.
     *
     * @param serverTotalItems the total as reported, may be null
     * @param excludedCount    how many rows were withheld
     * @return the adjusted total, or null when there is nothing dependable to report
     */
    private static String reduceTotal(final String serverTotalItems, final int excludedCount) {
        if (serverTotalItems == null || serverTotalItems.isEmpty()) {
            return null;
        }
        try {
            long total = Long.parseLong(serverTotalItems.trim());
            return String.valueOf(Math.max(0L, total - excludedCount));
        } catch (NumberFormatException e) {
            return null;
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

        // Step 4: Verify metadata.hash is in signed hashes list. Returns the hash
        // that was actually covered, which may be a legacy one.
        String verifiedHash = verifyHashInSignedHashes(envelope);

        // Step 5: Verify whitelist signatures are valid per governance rules
        verifyWhitelistSignatures(envelope, cachedRulesContainer, verifiedHash);

        // Step 6: same derivation as the full-data path — one implementation on the
        // envelope, so the cached-container path cannot drift from it.
        envelope.markVerified(cachedRulesContainer);
    }
}
