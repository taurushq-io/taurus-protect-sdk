package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.GovernanceRulesMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.GovernanceRulesHistoryResult;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.SuperAdminPublicKey;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.GovernanceRulesApi;
import com.taurushq.sdk.protect.openapi.model.GetPublicKeysReplyPublicKey;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPublicKeysReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRulesHistoryReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetRulesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRules;

import java.security.PublicKey;
import java.util.Collections;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing governance rules.
 */
public class GovernanceRuleService {

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
            return verifyGovernanceRules(GovernanceRulesMapper.INSTANCE.fromDTO(result), minValidSignatures);
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
            return verifyGovernanceRules(GovernanceRulesMapper.INSTANCE.fromDTO(result), minValidSignatures);
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

            List<TgvalidatordRules> rules = reply.getResult();
            if (rules == null) {
                result.setRules(Collections.emptyList());
            } else {
                result.setRules(GovernanceRulesMapper.INSTANCE.fromRulesDTOs(rules));
            }

            result.setCursor(reply.getCursor());
            result.setTotalItems(reply.getTotalItems());

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
    public GovernanceRules verifyGovernanceRules(GovernanceRules rules, int minValidSignatures) {
        SignatureVerifier.verifyGovernanceRules(rules, minValidSignatures, superAdminPublicKeys);
        return rules;
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
}
