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
