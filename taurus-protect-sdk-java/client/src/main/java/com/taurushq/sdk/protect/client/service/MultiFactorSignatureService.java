package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Preconditions;
import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.MultiFactorSignatureMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureApprovalResult;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureInfo;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.MultiFactorSignatureApi;
import com.taurushq.sdk.protect.openapi.model.MultiFactorSignatureServiceApproveMultiFactorSignatureBody;
import com.taurushq.sdk.protect.openapi.model.MultiFactorSignatureServiceRejectMultiFactorSignatureBody;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveMultiFactorSignatureReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateMultiFactorSignaturesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateMultiFactorSignaturesRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMultiFactorSignaturesEntityType;

import java.util.List;

/**
 * Service for managing multi-factor signature operations.
 * <p>
 * Multi-factor signatures are used for operations that require approval from
 * multiple parties, such as critical governance changes, high-value transactions,
 * or sensitive administrative actions.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get info about a pending multi-factor signature
 * MultiFactorSignatureInfo info = client.getMultiFactorSignatureService()
 *     .getMultiFactorSignatureInfo("mfs-123");
 *
 * // Approve the multi-factor signature
 * client.getMultiFactorSignatureService().approveMultiFactorSignature("mfs-123", signatures);
 * }</pre>
 *
 * @see MultiFactorSignatureInfo
 */
public class MultiFactorSignatureService {

    private final MultiFactorSignatureApi multiFactorSignatureApi;
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Creates a new MultiFactorSignatureService.
     *
     * @param apiClient          the API client for making requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public MultiFactorSignatureService(final ApiClient apiClient,
                                        final ApiExceptionMapper apiExceptionMapper) {
        Preconditions.checkNotNull(apiClient, "apiClient must not be null");
        Preconditions.checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.multiFactorSignatureApi = new MultiFactorSignatureApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
    }

    /**
     * Retrieves information about a multi-factor signature request.
     *
     * @param id the multi-factor signature ID
     * @return the multi-factor signature info
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public MultiFactorSignatureInfo getMultiFactorSignatureInfo(final String id) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(id), "id must not be null or empty");
        try {
            TgvalidatordGetMultiFactorSignatureEntitiesInfoReply reply =
                    multiFactorSignatureApi.multiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo(id);
            return MultiFactorSignatureMapper.INSTANCE.fromDTO(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Creates a batch of multi-factor signature requests.
     *
     * @param entityIds  the list of entity IDs to create signatures for
     * @param entityType the type of entities (REQUEST, WHITELISTED_ADDRESS, or WHITELISTED_CONTRACT)
     * @return the result containing the created signature IDs
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if entityIds is null or empty, or entityType is null
     */
    public MultiFactorSignatureResult createMultiFactorSignatures(
            final List<String> entityIds,
            final TgvalidatordMultiFactorSignaturesEntityType entityType) throws ApiException {
        Preconditions.checkArgument(entityIds != null && !entityIds.isEmpty(),
                "entityIds must not be null or empty");
        Preconditions.checkNotNull(entityType, "entityType must not be null");
        try {
            TgvalidatordCreateMultiFactorSignaturesRequest request =
                    new TgvalidatordCreateMultiFactorSignaturesRequest();
            request.setEntityIDs(entityIds);
            request.setEntityType(entityType);

            TgvalidatordCreateMultiFactorSignaturesReply reply =
                    multiFactorSignatureApi.multiFactorSignatureServiceCreateMultiFactorSignatureBatch(request);
            return MultiFactorSignatureMapper.INSTANCE.fromCreateDTO(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Approves a multi-factor signature request.
     *
     * @param id        the multi-factor signature ID
     * @param signature the signature for approval
     * @param comment   optional comment for the approval
     * @return the approval result
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty, or signature is null or empty
     */
    public MultiFactorSignatureApprovalResult approveMultiFactorSignature(final String id,
                                                                           final String signature,
                                                                           final String comment)
            throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(id), "id must not be null or empty");
        Preconditions.checkArgument(!Strings.isNullOrEmpty(signature), "signature must not be null or empty");
        try {
            MultiFactorSignatureServiceApproveMultiFactorSignatureBody body =
                    new MultiFactorSignatureServiceApproveMultiFactorSignatureBody();
            body.setSignature(signature);
            body.setComment(comment);

            TgvalidatordApproveMultiFactorSignatureReply reply =
                    multiFactorSignatureApi.multiFactorSignatureServiceApproveMultiFactorSignature(id, body);
            return MultiFactorSignatureMapper.INSTANCE.fromApprovalDTO(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Rejects a multi-factor signature request.
     *
     * @param id      the multi-factor signature ID
     * @param comment optional comment explaining the rejection
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if id is null or empty
     */
    public void rejectMultiFactorSignature(final String id, final String comment) throws ApiException {
        Preconditions.checkArgument(!Strings.isNullOrEmpty(id), "id must not be null or empty");
        try {
            MultiFactorSignatureServiceRejectMultiFactorSignatureBody body =
                    new MultiFactorSignatureServiceRejectMultiFactorSignatureBody();
            body.setComment(comment);

            multiFactorSignatureApi.multiFactorSignatureServiceRejectMultiFactorSignature(id, body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
