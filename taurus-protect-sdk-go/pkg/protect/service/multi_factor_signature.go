package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// MultiFactorSignatureService provides multi-factor signature management operations.
// Multi-factor signatures are used for operations that require approval from multiple parties,
// such as critical governance changes, high-value transactions, or sensitive administrative actions.
type MultiFactorSignatureService struct {
	api       *openapi.MultiFactorSignatureAPIService
	errMapper *ErrorMapper
}

// NewMultiFactorSignatureService creates a new MultiFactorSignatureService.
func NewMultiFactorSignatureService(client *openapi.APIClient) *MultiFactorSignatureService {
	return &MultiFactorSignatureService{
		api:       client.MultiFactorSignatureAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetMultiFactorSignatureInfo retrieves information about a multi-factor signature request.
//
// SECURITY — PayloadToSign IS UNVERIFIED SERVER DATA. The caller must bind it to a verified
// entity before signing it.
//
// This is the one read path in the SDK that returns bytes intended for a signing key without
// verifying them, and it is deliberate-but-unresolved rather than an oversight (TODOS.md,
// "MultiFactorSignatureService.approve takes an opaque caller signature"). The reason it cannot
// be fixed here: the reply carries only {id, payloadToSign[], entityType}. There is **no entity
// id**, singular or plural, and payloadToSign is a bare string array with no per-element id or
// hash — so nothing in the reply can be joined back to the entities the request was created for,
// and this method has nothing to check against.
//
// What that means for a second-factor signer: a malicious or compromised API server can answer
// this call with the metadata hash of an entity of its choosing — a withdrawal to an attacker
// address, an attacker-controlled whitelisted address — under the expected kind. Sign it and the
// server holds a valid MobileAppSigner approval over an entity nobody reviewed. That signature
// is precisely the artefact a compromised server cannot forge on its own.
//
// The workable client-side binding, for when the decision lands: CreateMultiFactorSignatures
// takes the entity IDs, so a caller who keeps them can re-read those entities through the
// verifying reader for that entityType (RequestService, WhitelistedAddressService) and require
// each PayloadToSign element to equal a locally recomputed, verified metadata hash. Until the
// SDK does that for you, do it yourself.
func (s *MultiFactorSignatureService) GetMultiFactorSignatureInfo(ctx context.Context, id string) (*model.MultiFactorSignatureInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.MultiFactorSignatureServiceGetMultiFactorSignatureEntitiesInfo(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.MultiFactorSignatureInfoFromDTO(resp), nil
}

// CreateMultiFactorSignatures creates a batch of multi-factor signature requests.
func (s *MultiFactorSignatureService) CreateMultiFactorSignatures(
	ctx context.Context,
	entityIDs []string,
	entityType model.MultiFactorSignatureEntityType,
) (*model.MultiFactorSignatureResult, error) {
	if len(entityIDs) == 0 {
		return nil, fmt.Errorf("entityIDs cannot be empty")
	}
	if entityType == "" {
		return nil, fmt.Errorf("entityType cannot be empty")
	}

	createReq := openapi.TgvalidatordCreateMultiFactorSignaturesRequest{
		EntityIDs:  entityIDs,
		EntityType: mapper.EntityTypeToDTO(entityType),
	}

	resp, httpResp, err := s.api.MultiFactorSignatureServiceCreateMultiFactorSignatureBatch(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.MultiFactorSignatureResultFromDTO(resp), nil
}

// ApproveMultiFactorSignature approves a multi-factor signature request.
//
// SECURITY — THE SIGNATURE IS OPAQUE TO THIS SDK. It is forwarded after a non-empty check and
// nothing here knows or checks what it covers, so this method cannot enforce that only verified
// payloads are approved.
//
// That is why this site is absent from scripts/signing-sites/manifest.json: the manifest
// enumerates sites where the *SDK* signs, and here the caller does. It is also why its absence
// is not reassuring — the gate that would have forced a `verifies` / `signs-own-bytes`
// classification simply does not see this path.
//
// The caller is responsible for having bound the bytes they signed to a verified entity. See
// GetMultiFactorSignatureInfo for why the SDK cannot do it for them yet, and TODOS.md for the
// open decision.
func (s *MultiFactorSignatureService) ApproveMultiFactorSignature(
	ctx context.Context,
	id string,
	signature string,
	comment string,
) (*model.MultiFactorSignatureApprovalResult, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	if signature == "" {
		return nil, fmt.Errorf("signature cannot be empty")
	}

	body := openapi.MultiFactorSignatureServiceApproveMultiFactorSignatureBody{
		Signature: signature,
		Comment:   comment,
	}

	resp, httpResp, err := s.api.MultiFactorSignatureServiceApproveMultiFactorSignature(ctx, id).
		Body(body).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.MultiFactorSignatureApprovalResultFromDTO(resp), nil
}

// RejectMultiFactorSignature rejects a multi-factor signature request.
func (s *MultiFactorSignatureService) RejectMultiFactorSignature(ctx context.Context, id string, comment string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	body := openapi.MultiFactorSignatureServiceRejectMultiFactorSignatureBody{
		Comment: comment,
	}

	_, httpResp, err := s.api.MultiFactorSignatureServiceRejectMultiFactorSignature(ctx, id).
		Body(body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
