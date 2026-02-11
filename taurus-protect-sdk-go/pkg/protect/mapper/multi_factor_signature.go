package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// MultiFactorSignatureInfoFromDTO converts an OpenAPI reply to domain MultiFactorSignatureInfo.
func MultiFactorSignatureInfoFromDTO(dto *openapi.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply) *model.MultiFactorSignatureInfo {
	if dto == nil {
		return nil
	}

	return &model.MultiFactorSignatureInfo{
		ID:            dto.Id,
		PayloadToSign: dto.PayloadToSign,
		EntityType:    model.MultiFactorSignatureEntityType(dto.EntityType),
	}
}

// MultiFactorSignatureResultFromDTO converts an OpenAPI reply to domain MultiFactorSignatureResult.
func MultiFactorSignatureResultFromDTO(dto *openapi.TgvalidatordCreateMultiFactorSignaturesReply) *model.MultiFactorSignatureResult {
	if dto == nil {
		return nil
	}

	return &model.MultiFactorSignatureResult{
		ID: safeString(dto.Id),
	}
}

// MultiFactorSignatureApprovalResultFromDTO converts an OpenAPI reply to domain MultiFactorSignatureApprovalResult.
func MultiFactorSignatureApprovalResultFromDTO(dto *openapi.TgvalidatordApproveMultiFactorSignatureReply) *model.MultiFactorSignatureApprovalResult {
	if dto == nil {
		return nil
	}

	return &model.MultiFactorSignatureApprovalResult{
		SignatureCount: dto.Signatures,
	}
}

// EntityTypeToDTO converts a domain entity type to OpenAPI entity type.
func EntityTypeToDTO(entityType model.MultiFactorSignatureEntityType) openapi.TgvalidatordMultiFactorSignaturesEntityType {
	return openapi.TgvalidatordMultiFactorSignaturesEntityType(entityType)
}
