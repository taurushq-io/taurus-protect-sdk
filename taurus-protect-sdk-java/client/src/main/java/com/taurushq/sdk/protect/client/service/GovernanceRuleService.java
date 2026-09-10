package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.helper.StrictBase64;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.GovernanceRulesMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ExcludedRuleset;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.GovernanceRulesHistoryResult;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.SuperAdminPublicKey;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.GovernanceRulesApi;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.openapi.model.GetPublicKeysReplyPublicKey;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveRulesProposalRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPublicKeysReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRulesHistoryReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRulesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRejectRulesProposalRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRules;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUpdateRulesProposalRequest;

import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SignatureException;
import java.util.Base64;
import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing governance rules.
 */
public class GovernanceRuleService {

    /** Bounds the verification memo. */
    private static final int MAX_VERIFIED_RULESETS = 256;

    /**
     * Rulesets whose SuperAdmin signatures already checked out, keyed by
     * {@link #rulesetVerificationKey}. See {@link #verifiedRuleset}.
     */
    private final java.util.Set<String> verified =
            java.util.Collections.newSetFromMap(new java.util.concurrent.ConcurrentHashMap<>());


    private final GovernanceRulesApi governanceRulesApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final List<PublicKey> superAdminPublicKeys;
    private final int minValidSignatures;

    /**
     * Instantiates a new Governance rule service.
     *
     * @param openApiClient        the open api client
     * @param apiExceptionMapper   the api exception mapper
     * @param superAdminPublicKeys the list of SuperAdmin public keys for verification
     * @param minValidSignatures   the minimum number of valid signatures required for verification
     */
    public GovernanceRuleService(final ApiClient openApiClient,
                                 final ApiExceptionMapper apiExceptionMapper,
                                 final List<PublicKey> superAdminPublicKeys,
                                 final int minValidSignatures) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");
        checkNotNull(superAdminPublicKeys, "superAdminPublicKeys cannot be null");
        checkArgument(!superAdminPublicKeys.isEmpty(), "superAdminPublicKeys cannot be empty");
        checkArgument(minValidSignatures > 0, "minValidSignatures must be greater than zero");

        this.apiExceptionMapper = apiExceptionMapper;
        this.governanceRulesApi = new GovernanceRulesApi(openApiClient);
        this.superAdminPublicKeys = superAdminPublicKeys;
        this.minValidSignatures = minValidSignatures;
    }

    /**
     * Gets the currently enforced governance rules.
     *
     * @return the governance rules
     * @throws ApiException the api exception
     */
    public GovernanceRules getRules() throws ApiException {
        try {
            TgvalidatordGetRulesReply reply = governanceRulesApi.ruleServiceGetRules();
            TgvalidatordRules result = reply.getResult();
            if (result == null) {
                return null;
            }
            return verifiedRuleset(GovernanceRulesMapper.INSTANCE.fromDTO(result));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets a governance ruleset by its ID.
     *
     * @param id the ruleset id
     * @return the governance rules
     * @throws ApiException the api exception
     */
    public GovernanceRules getRulesById(final String id) throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(id), "id cannot be null or empty");

        try {
            TgvalidatordGetRulesReply reply = governanceRulesApi.ruleServiceGetRulesByID(id);
            TgvalidatordRules result = reply.getResult();
            if (result == null) {
                return null;
            }
            return verifiedRuleset(GovernanceRulesMapper.INSTANCE.fromDTO(result));
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets governance rules history with cursor-based pagination.
     *
     * @param pageSize the page size
     * @param cursor   the cursor from previous response (null for first page)
     * @return the governance rules history result with rules and pagination cursor
     * @throws ApiException the api exception
     */
    public GovernanceRulesHistoryResult getRulesHistory(final int pageSize, final byte[] cursor) throws ApiException {
        checkArgument(pageSize > 0, "pageSize must be positive");

        try {
            TgvalidatordGetRulesHistoryReply reply = governanceRulesApi.ruleServiceGetRulesHistory(
                    String.valueOf(pageSize),
                    cursor
            );

            GovernanceRulesHistoryResult result = new GovernanceRulesHistoryResult();

            // Every historical ruleset is a past ENFORCED document, so each is verified.
            // The memo makes this one ECDSA pass per distinct container, not one per row.
            //
            //   entries -> verifiedRuleset -+- ok    -> kept
            //                               +- error -> excluded + named in the result
            //
            // LENIENT, unlike the whitelist lists, and it does NOT throw when nothing
            // survives: a SuperAdmin key rotation makes every pre-rotation ruleset
            // unverifiable, so aborting would deny the whole audit trail from then on.
            List<TgvalidatordRules> rules = reply.getResult();
            List<GovernanceRules> entries = rules == null
                    ? Collections.emptyList()
                    : GovernanceRulesMapper.INSTANCE.fromRulesDTOs(rules);
            List<ExcludedRuleset> excluded = new java.util.ArrayList<>();
            result.setRules(verifiedHistoryEntries(entries, excluded));
            result.setExcludedUnverified(excluded);

            result.setCursor(reply.getCursor());
            // Reduced by the exclusions: the server counts rows it returned, the caller
            // receives only those that verified.
            result.setTotalItems(reducedTotal(reply.getTotalItems(), excluded.size()));

            return result;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets governance rules history (first page).
     *
     * @param pageSize the page size
     * @return the governance rules history result with rules and pagination cursor
     * @throws ApiException the api exception
     */
    public GovernanceRulesHistoryResult getRulesHistory(final int pageSize) throws ApiException {
        return getRulesHistory(pageSize, null);
    }

    /**
     * Gets the proposed governance rules.
     * Requires SuperAdmin or SuperAdminReadOnly role.
     *
     * @return the proposed governance rules
     * @throws ApiException the api exception
     */
    public GovernanceRules getRulesProposal() throws ApiException {
        try {
            TgvalidatordGetRulesReply reply = governanceRulesApi.ruleServiceGetRulesProposal();
            TgvalidatordRules result = reply.getResult();
            if (result == null) {
                return null;
            }
            return GovernanceRulesMapper.INSTANCE.fromDTO(result);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets the list of superadmin public keys.
     *
     * @return the list of public keys
     * @throws ApiException the api exception
     */
    public List<SuperAdminPublicKey> getPublicKeys() throws ApiException {
        try {
            TgvalidatordGetPublicKeysReply reply = governanceRulesApi.ruleServiceGetPublicKeys();
            List<GetPublicKeysReplyPublicKey> publicKeys = reply.getPublicKeys();
            if (publicKeys == null) {
                return Collections.emptyList();
            }
            return GovernanceRulesMapper.INSTANCE.fromPublicKeyDTOs(publicKeys);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Verifies that governance rules have enough valid SuperAdmin signatures.
     *
     * @param rules              the governance rules to verify
     * @param minValidSignatures the minimum number of valid signatures required
     * @throws IntegrityException if verification fails or not enough valid signatures
     */
    /**
     * Keeps the history entries whose SuperAdmin signatures verify, collecting the rest
     * into {@code excluded}.
     *
     * <pre>
     *   entries -&gt; verifiedRuleset -+- ok    -&gt; kept
     *                               +- error -&gt; excluded + named in the result
     * </pre>
     *
     * <p>LENIENT, unlike the whitelist lists, and it does NOT throw when nothing
     * survives: a SuperAdmin key rotation makes every pre-rotation ruleset unverifiable,
     * so aborting the page would deny access to the whole audit trail from then on.
     *
     * <p>Package-private so the exclusion behaviour can be tested directly; this class
     * has no seam for injecting a stub API.
     *
     * @param entries  the mapped history entries
     * @param excluded collects the entries that did not verify
     * @return only those that verified
     */
    List<GovernanceRules> verifiedHistoryEntries(final List<GovernanceRules> entries,
                                                 final List<ExcludedRuleset> excluded) {
        List<GovernanceRules> kept = new java.util.ArrayList<>(entries.size());
        for (GovernanceRules entry : entries) {
            try {
                kept.add(verifiedRuleset(entry));
            } catch (IntegrityException e) {
                ExcludedRuleset ex = new ExcludedRuleset();
                ex.setCreationDate(entry.getCreationDate());
                ex.setReason(e.getMessage());
                excluded.add(ex);
            }
        }
        return kept;
    }

    /**
     * Reduces a server-reported total by the number of excluded rows. A non-numeric
     * total is the server's to explain, so it passes through rather than inventing one.
     *
     * @param total      the server-reported total, may be null
     * @param excluded   how many rows were withheld
     * @return the reduced total, or the original when it is absent or non-numeric
     */
    private static String reducedTotal(final String total, final int excluded) {
        if (total == null || excluded == 0) {
            return total;
        }
        try {
            return String.valueOf(Math.max(0, Long.parseLong(total) - excluded));
        } catch (NumberFormatException e) {
            return total;
        }
    }

    /**
     * Verifies a ruleset's SuperAdmin signatures, memoising the outcome so the same
     * document is not re-verified on every read.
     *
     * <pre>
     *   getRules / getRulesById / getRulesHistory -&gt; verifiedRuleset
     *                                                 |- memo hit -&gt; return
     *                                                 +- miss -&gt; ECDSA -&gt; memoise
     * </pre>
     *
     * <p>The container itself is already cached for the address/asset/price paths by
     * {@code RulesContainerCache}, which fetches through {@code getDecodedRulesContainer}.
     * This memo covers only the reads that return the RAW ruleset. Do not add a third
     * cache of the same bytes.
     *
     * <p>Only SUCCESSES are memoised: a failure must resurface its error on every call.
     *
     * @param rules the ruleset to verify
     * @return the same ruleset, verified
     */
    private GovernanceRules verifiedRuleset(final GovernanceRules rules) {
        String key = rulesetVerificationKey(rules);
        if (key != null && verified.contains(key)) {
            return rules;
        }
        verifyGovernanceRules(rules, minValidSignatures);
        // An undecodable container has no stable identity to memoise on; verification
        // above has already had the final say.
        if (key == null) {
            return rules;
        }
        // Bounded so walking a long history cannot grow it without limit. Containers
        // change rarely, so a plain clear beats LRU bookkeeping.
        if (verified.size() >= MAX_VERIFIED_RULESETS) {
            verified.clear();
        }
        verified.add(key);
        return rules;
    }

    /**
     * Writes an 8-byte big-endian length, then the bytes, so a sequence of fields cannot
     * be re-partitioned into a different sequence that encodes identically.
     *
     * @param digest the digest to update
     * @param data   the field bytes
     */
    private static void updateLengthPrefixed(final java.security.MessageDigest digest,
                                             final byte[] data) {
        digest.update(bigEndian(data.length));
        digest.update(data);
    }

    /**
     * Encodes a length as 8 big-endian bytes.
     *
     * @param value the length
     * @return the 8-byte encoding
     */
    private static byte[] bigEndian(final long value) {
        byte[] out = new byte[8];
        long n = value;
        for (int i = 7; i >= 0; i--) {
            out[i] = (byte) (n & 0xFF);
            n >>>= 8;
        }
        return out;
    }

    /**
     * Identifies the exact document verification cleared: the DECODED container bytes --
     * the same bytes {@link SignatureVerifier#verifyGovernanceRules} checks the signatures
     * over -- plus every signature over them. A container whose signature set changed is a
     * memo MISS rather than a stale hit; signature order is server-controlled, so it is
     * sorted out of the key.
     *
     * <p>Every field is LENGTH-PREFIXED, and that is load-bearing. A plain concatenation
     * is not injective: with {@code sha256(container || sig || ...)} the boundary between
     * the container and the signature list is not committed, so a response-controlling
     * attacker can shift bytes across it and make a MODIFIED container collide with a
     * genuine one's key -- inheriting its "already verified" status and skipping ECDSA
     * entirely. Interleaving a {@code 0x00} separator does NOT fix this; only length
     * prefixes do.
     *
     * <p>{@code userId} is deliberately excluded: verification never reads it (it consumes
     * only the signature), so a server-controlled field that verification ignores must not
     * be able to influence a verification-SKIP decision.
     *
     * <p>Keying on the decoded bytes rather than the base64 text also removes base64
     * leniency as a lever, and guarantees the key cannot identify something other than
     * what was verified.
     *
     * @param rules the ruleset
     * @return a hex digest key, or null when the container does not decode -- in which
     *     case there is nothing worth memoising and the caller must fall through to
     *     verification
     */
    static String rulesetVerificationKey(final GovernanceRules rules) {
        final byte[] data;
        try {
            data = StrictBase64.decode(rules.getRulesContainer());
        } catch (IllegalArgumentException e) {
            // A container we cannot decode has no stable identity to memoise on.
            return null;
        }

        try {
            java.security.MessageDigest digest = java.security.MessageDigest.getInstance("SHA-256");
            updateLengthPrefixed(digest, data);

            java.util.List<String> signatures = new java.util.ArrayList<>();
            if (rules.getRulesSignatures() != null) {
                for (RuleUserSignature sig : rules.getRulesSignatures()) {
                    signatures.add(sig.getSignature() == null ? "" : sig.getSignature());
                }
            }
            java.util.Collections.sort(signatures);

            digest.update(bigEndian(signatures.size()));
            for (String signature : signatures) {
                updateLengthPrefixed(digest,
                        signature.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            }

            StringBuilder hex = new StringBuilder();
            for (byte b : digest.digest()) {
                hex.append(String.format("%02x", b));
            }
            return hex.toString();
        } catch (java.security.NoSuchAlgorithmException e) {
            // SHA-256 is mandated by every JRE; if it is absent nothing here can work.
            throw new IllegalStateException("SHA-256 unavailable", e);
        }
    }

    public GovernanceRules verifyGovernanceRules(GovernanceRules rules, int minValidSignatures) {
        SignatureVerifier.verifyGovernanceRules(rules, minValidSignatures, superAdminPublicKeys);
        return rules;
    }

    /**
     * Verifies governance rules against the threshold this service was configured with.
     *
     * <p>The single-argument form is the cross-SDK shape (Go, Python and TypeScript all
     * take only the rules and read the threshold from the service). Prefer it: passing a
     * threshold that differs from the configured one silently verifies against something
     * other than what the rest of the service enforces.
     *
     * @param rules the governance rules to verify
     * @return the verified rules
     * @throws IntegrityException if verification fails or there are too few valid signatures
     */
    public GovernanceRules verifyGovernanceRules(GovernanceRules rules) {
        return verifyGovernanceRules(rules, minValidSignatures);
    }

    /**
     * Gets the decoded rules container from governance rules.
     * Verifies signatures and decodes the rules container using the configured keys.
     *
     * @param rules the governance rules
     * @return the decoded rules container
     * @throws IntegrityException if signature verification fails
     */
    public DecodedRulesContainer getDecodedRulesContainer(GovernanceRules rules) throws IntegrityException {
        checkNotNull(rules, "rules cannot be null");
        return rules.getDecodedRulesContainer(superAdminPublicKeys, minValidSignatures);
    }

    /**
     * Gets the configured SuperAdmin public keys.
     *
     * @return the list of SuperAdmin public keys
     */
    public List<PublicKey> getSuperAdminPublicKeys() {
        return Collections.unmodifiableList(superAdminPublicKeys);
    }

    /**
     * Gets the configured minimum valid signatures required.
     *
     * @return the minimum valid signatures
     */
    public int getMinValidSignatures() {
        return minValidSignatures;
    }

    /**
     * Submits a new governance rules proposal.
     *
     * <p>The container is encoded to protobuf (server-controlled fields are omitted) and
     * submitted as the pending proposal. This does not return the stored proposal: the
     * update endpoint returns no body, and a read-back would be unreliable because the
     * enforced rules are cached server-side. Use {@link #getRulesProposal()} to inspect
     * the pending proposal afterwards.
     *
     * @param container the rules container to propose
     * @throws ApiException if the submission fails
     */
    public void updateRulesProposal(final DecodedRulesContainer container) throws ApiException {
        checkNotNull(container, "container cannot be null");

        TgvalidatordUpdateRulesProposalRequest body = new TgvalidatordUpdateRulesProposalRequest();
        body.setRulesContainer(RulesContainerMapper.INSTANCE.toBase64String(container));
        try {
            governanceRulesApi.ruleServiceUpdateRulesProposal(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Returns the canonical SHA-256 hex digest of a ruleset's decoded rules container.
     *
     * <p>This is the value to pass as {@code expectedContainerHash} to
     * {@link #approveRulesProposal}, and it is what pins the approval to reviewed content.
     *
     * <p>It digests the DECODED bytes, not the base64 text: base64 has multiple encodings
     * of the same bytes, so hashing the text would let a re-encoded but byte-identical
     * container read as a mismatch. Those decoded bytes are also exactly what gets signed.
     *
     * @param rules the ruleset whose container to digest
     * @return the hex digest
     * @throws IntegrityException if the ruleset carries no container, or it is not valid base64
     */
    public String proposalContainerHash(final GovernanceRules rules) throws IntegrityException {
        if (rules == null || Strings.isNullOrEmpty(rules.getRulesContainer())) {
            throw new IntegrityException("ruleset carries no rules container");
        }
        final byte[] data;
        try {
            data = StrictBase64.decode(rules.getRulesContainer());
        } catch (IllegalArgumentException e) {
            throw new IntegrityException("failed to decode rules container: " + e.getMessage(), e);
        }
        try {
            java.security.MessageDigest digest = java.security.MessageDigest.getInstance("SHA-256");
            StringBuilder hex = new StringBuilder();
            for (byte b : digest.digest(data)) {
                hex.append(String.format("%02x", b));
            }
            return hex.toString();
        } catch (java.security.NoSuchAlgorithmException e) {
            // SHA-256 is mandated by every JRE; if it is absent nothing here can work.
            throw new IllegalStateException("SHA-256 unavailable", e);
        }
    }

    /**
     * Decodes a PENDING rules proposal so a SuperAdmin can inspect it before approving.
     *
     * <p>UNVERIFIED BY DESIGN. A pending proposal legitimately carries 0..N signatures --
     * signing IS the approval step -- so there is no threshold to check yet, and
     * {@link #getDecodedRulesContainer} (which verifies unconditionally) always fails on
     * one. That is why this is a separate, explicitly-named entry point rather than a
     * flag: the absence of verification has to be visible at the call site.
     *
     * <p>Pair it with {@link #proposalContainerHash} and pass that digest to
     * {@link #approveRulesProposal}, so the bytes signed are the bytes reviewed.
     *
     * @param rules the pending proposal, from {@link #getRulesProposal}
     * @return the decoded rules container, NOT verified
     * @throws IntegrityException if the ruleset carries no container, or decoding fails
     */
    public DecodedRulesContainer decodeProposalForReview(final GovernanceRules rules)
            throws IntegrityException {
        if (rules == null || Strings.isNullOrEmpty(rules.getRulesContainer())) {
            throw new IntegrityException("ruleset carries no rules container");
        }
        try {
            return RulesContainerMapper.INSTANCE.fromBase64String(rules.getRulesContainer());
        } catch (com.google.protobuf.InvalidProtocolBufferException e) {
            throw new IntegrityException("unable to decode the rules container from proto", e);
        }
    }

    /**
     * Approves the pending governance rules proposal.
     *
     * <p>{@code expectedContainerHash} PINS the content being approved: it must be the
     * {@link #proposalContainerHash} of the proposal the caller actually reviewed. The
     * proposal is re-fetched here, and if its container does not match that digest the
     * call aborts WITHOUT signing.
     *
     * <p>That pin is the security control. Without it a server able to shape responses
     * could serve the benign proposal to the review call and a different container to the
     * re-fetch inside this method, obtaining a GENUINE SuperAdmin signature over bytes of
     * its choosing -- and such a container then verifies clean everywhere, including in
     * independent and air-gapped verifiers that never trusted that server. Nothing else
     * binds the two calls: the wire format carries no proposal id, hash or version, and a
     * pending proposal has no signatures to verify against, so re-verification cannot
     * substitute for pinning.
     *
     * <p>The signature is computed over the decoded bytes of the pending proposal's rules
     * container (SHA-256 + P-256 ECDSA, base64 raw r||s).
     *
     * @param privateKey            the SuperAdmin private key to sign with
     * @param comment               an optional approval comment (may be null)
     * @param expectedContainerHash {@link #proposalContainerHash} of the reviewed proposal
     * @throws ApiException       if there is no pending proposal, signing fails, or the submission fails
     * @throws IntegrityException if the pending container differs from the pin
     */
    public void approveRulesProposal(final PrivateKey privateKey, final String comment,
                                     final String expectedContainerHash) throws ApiException {
        checkNotNull(privateKey, "privateKey cannot be null");
        if (Strings.isNullOrEmpty(expectedContainerHash)) {
            throw new IllegalArgumentException("expectedContainerHash is required: it pins the "
                    + "approval to the container you reviewed (see proposalContainerHash)");
        }

        GovernanceRules proposal = getRulesProposal();
        if (proposal == null || Strings.isNullOrEmpty(proposal.getRulesContainer())) {
            ApiException e = new ApiException();
            e.setCode(404);
            e.setError("ClientInvalidRequest");
            e.setMessage("no pending rules proposal to approve");
            throw e;
        }

        String actualHash = proposalContainerHash(proposal);
        // Constant-time, matching how every other hash comparison in this SDK is done.
        if (!java.security.MessageDigest.isEqual(
                actualHash.getBytes(java.nio.charset.StandardCharsets.US_ASCII),
                expectedContainerHash.getBytes(java.nio.charset.StandardCharsets.US_ASCII))) {
            throw new IntegrityException(String.format(
                    "refusing to sign: the pending rules proposal changed since it was reviewed "
                            + "(reviewed %s, pending %s)", expectedContainerHash, actualHash));
        }

        byte[] toSign = StrictBase64.decode(proposal.getRulesContainer());
        TgvalidatordApproveRulesProposalRequest body = new TgvalidatordApproveRulesProposalRequest();
        try {
            body.setSignature(CryptoTPV1.calculateBase64Signature(privateKey, toSign));
        } catch (NoSuchAlgorithmException | InvalidKeyException | SignatureException ex) {
            ApiException e = new ApiException();
            e.setCode(400);
            e.setError("ClientInvalidRequest");
            e.setMessage("unable to sign the pending rules proposal");
            e.setOriginalException(ex);
            throw e;
        }
        if (comment != null) {
            body.setComment(comment);
        }
        try {
            governanceRulesApi.ruleServiceApproveRulesProposal(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Rejects the pending governance rules proposal.
     *
     * @param comment an optional rejection comment (may be null)
     * @throws ApiException if the submission fails
     */
    public void rejectRulesProposal(final String comment) throws ApiException {
        TgvalidatordRejectRulesProposalRequest body = new TgvalidatordRejectRulesProposalRequest();
        if (comment != null) {
            body.setComment(comment);
        }
        try {
            governanceRulesApi.ruleServiceRejectRulesProposal(body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
