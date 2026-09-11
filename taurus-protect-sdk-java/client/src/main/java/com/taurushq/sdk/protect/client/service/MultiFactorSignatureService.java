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
     * <p><b>SECURITY — {@code payloadToSign} IS UNVERIFIED SERVER DATA.</b> The caller must bind
     * it to a verified entity before signing it.
     *
     * <p>This is the one read path in the SDK that returns bytes intended for a signing key
     * without verifying them, and it is deliberate-but-unresolved rather than an oversight (see
     * TODOS.md, "MultiFactorSignatureService.approve takes an opaque caller signature"). The
     * reason it cannot be fixed here: the reply carries only
     * {@code {id, payloadToSign[], entityType}}. There is <b>no entity id</b>, singular or
     * plural, and {@code payloadToSign} is a bare string list with no per-element id or hash — so
     * nothing in the reply can be joined back to the entities the request was created for, and
     * this method has nothing to check against. Note also that {@code entityType} is a bare kind
     * enum: {@link com.taurushq.sdk.protect.client.model.MultiFactorSignatureEntityType} declares
     * {@code id} and {@code kind}, and MapStruct's enum-source mapping leaves both null.
     *
     * <p>What that means for a second-factor signer: a malicious or compromised API server can
     * answer this call with the metadata hash of an entity of its choosing — a withdrawal to an
     * attacker address, an attacker-controlled whitelisted address — under the expected kind.
     * Sign it and the server holds a valid MobileAppSigner approval over an entity nobody
     * reviewed. That signature is precisely the artefact a compromised server cannot forge on its
     * own. Contrast {@code RequestService.approveRequests}, which re-verifies every metadata hash
     * before the approver's key signs.
     *
     * <p>The workable client-side binding, for when the decision lands:
     * {@link #createMultiFactorSignatures} takes the entity IDs, so a caller who keeps them can
     * re-read those entities through the verifying reader for that entity type and require each
     * {@code payloadToSign} element to equal a locally recomputed, verified metadata hash. Until
     * the SDK does that for you, do it yourself.
     *
     * @param id the multi-factor signature ID
     * @return the multi-factor signature info, whose {@code payloadToSign} is UNVERIFIED
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
     * <p><b>SECURITY — THE SIGNATURE IS OPAQUE TO THIS SDK.</b> It is forwarded after a non-empty
     * check and nothing here knows or checks what it covers, so this method cannot enforce that
     * only verified payloads are approved.
     *
     * <p>That is why this site is absent from {@code scripts/signing-sites/manifest.json}: the
     * manifest enumerates sites where the <i>SDK</i> signs, and here the caller does. It is also
     * why its absence is not reassuring — the gate that would have forced a {@code verifies} /
     * {@code signs-own-bytes} classification simply does not see this path.
     *
     * <p>The caller is responsible for having bound the bytes they signed to a verified entity.
     * See {@link #getMultiFactorSignatureInfo} for why the SDK cannot do it for them yet, and
     * TODOS.md for the open decision.
     *
     * @param id        the multi-factor signature ID
     * @param signature the signature for approval. Assumed to cover bytes the CALLER verified;
     *                  this method does not check it
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
