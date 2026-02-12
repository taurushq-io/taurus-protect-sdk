package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAssetFromDTO converts an OpenAPI SignedWhitelistedContractAddressEnvelope to a domain WhitelistedAsset.
func WhitelistedAssetFromDTO(dto *openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope) *model.WhitelistedAsset {
	if dto == nil {
		return nil
	}

	asset := &model.WhitelistedAsset{
		ID:                  safeString(dto.Id),
		TenantID:            safeString(dto.TenantId),
		Status:              safeString(dto.Status),
		Action:              safeString(dto.Action),
		Blockchain:          safeString(dto.Blockchain),
		Network:             safeString(dto.Network),
		Rule:                safeString(dto.Rule),
		RulesContainer:      safeString(dto.RulesContainer),
		RulesSignatures:     safeString(dto.RulesSignatures),
		BusinessRuleEnabled: safeBool(dto.BusinessRuleEnabled),
	}

	// Convert metadata
	if dto.Metadata != nil {
		asset.Metadata = WhitelistedAssetMetadataFromDTO(dto.Metadata)
	}

	// Convert signed contract address
	if dto.SignedContractAddress != nil {
		asset.SignedContractAddress = SignedContractAddressFromDTO(dto.SignedContractAddress)
	}

	// Convert approvers
	if dto.Approvers != nil {
		asset.Approvers = ApproversFromDTO(dto.Approvers)
	}

	// Convert attributes
	if dto.Attributes != nil {
		asset.Attributes = make([]model.WhitelistedAssetAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			asset.Attributes[i] = WhitelistedAssetAttributeFromDTO(&attr)
		}
	}

	// Convert trails
	if dto.Trails != nil {
		asset.Trails = make([]model.Trail, len(dto.Trails))
		for i, trail := range dto.Trails {
			asset.Trails[i] = TrailFromDTO(&trail)
		}
	}

	return asset
}

// WhitelistedAssetsFromDTO converts a slice of OpenAPI SignedWhitelistedContractAddressEnvelope to domain WhitelistedAssets.
func WhitelistedAssetsFromDTO(dtos []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope) []*model.WhitelistedAsset {
	if dtos == nil {
		return nil
	}
	assets := make([]*model.WhitelistedAsset, len(dtos))
	for i := range dtos {
		assets[i] = WhitelistedAssetFromDTO(&dtos[i])
	}
	return assets
}

// WhitelistedAssetMetadataFromDTO converts an OpenAPI Metadata to a domain WhitelistedAssetMetadata.
func WhitelistedAssetMetadataFromDTO(dto *openapi.TgvalidatordMetadata) *model.WhitelistedAssetMetadata {
	if dto == nil {
		return nil
	}
	return &model.WhitelistedAssetMetadata{
		Hash: safeString(dto.Hash),
		// SECURITY: Payload intentionally not mapped - use PayloadAsString only.
		// The raw payload object could be tampered with while payloadAsString
		// remains unchanged (hash still verifies). By extracting from PayloadAsString,
		// we ensure all data comes from the cryptographically verified source.
		PayloadAsString: safeString(dto.PayloadAsString),
	}
}

// SignedContractAddressFromDTO converts an OpenAPI SignedWhitelistedContractAddress to a domain SignedContractAddress.
func SignedContractAddressFromDTO(dto *openapi.TgvalidatordSignedWhitelistedContractAddress) *model.SignedContractAddress {
	if dto == nil {
		return nil
	}

	signed := &model.SignedContractAddress{
		Payload: safeString(dto.Payload),
	}

	if dto.Signatures != nil {
		signed.Signatures = make([]model.WhitelistSignature, len(dto.Signatures))
		for i, sig := range dto.Signatures {
			signed.Signatures[i] = WhitelistSignatureFromDTO(&sig)
		}
	}

	return signed
}

// WhitelistSignatureFromDTO converts an OpenAPI WhitelistSignature to a domain WhitelistSignature.
func WhitelistSignatureFromDTO(dto *openapi.TgvalidatordWhitelistSignature) model.WhitelistSignature {
	if dto == nil {
		return model.WhitelistSignature{}
	}

	sig := model.WhitelistSignature{
		Hashes: dto.Hashes,
	}

	if dto.Signature != nil {
		sig.UserSignature = WhitelistUserSignatureFromDTO(dto.Signature)
	}

	return sig
}

// WhitelistUserSignatureFromDTO converts an OpenAPI WhitelistUserSignature to a domain WhitelistUserSignature.
func WhitelistUserSignatureFromDTO(dto *openapi.TgvalidatordWhitelistUserSignature) *model.WhitelistUserSignature {
	if dto == nil {
		return nil
	}
	return &model.WhitelistUserSignature{
		UserID:    safeString(dto.UserId),
		Signature: safeString(dto.Signature),
		Comment:   safeString(dto.Comment),
	}
}

// ApproversFromDTO converts an OpenAPI Approvers to a domain Approvers.
func ApproversFromDTO(dto *openapi.TgvalidatordApprovers) *model.Approvers {
	if dto == nil {
		return nil
	}

	approvers := &model.Approvers{}

	if dto.Parallel != nil {
		approvers.Parallel = make([]model.ParallelApproversGroup, len(dto.Parallel))
		for i, pg := range dto.Parallel {
			approvers.Parallel[i] = ParallelApproversGroupFromDTO(&pg)
		}
	}

	return approvers
}

// ParallelApproversGroupFromDTO converts an OpenAPI ParallelApproversGroups to a domain ParallelApproversGroup.
func ParallelApproversGroupFromDTO(dto *openapi.TgvalidatordParallelApproversGroups) model.ParallelApproversGroup {
	if dto == nil {
		return model.ParallelApproversGroup{}
	}

	pag := model.ParallelApproversGroup{}

	if dto.Sequential != nil {
		pag.Sequential = make([]model.ApproversGroup, len(dto.Sequential))
		for i, sg := range dto.Sequential {
			pag.Sequential[i] = ApproversGroupFromDTO(&sg)
		}
	}

	return pag
}

// ApproversGroupFromDTO converts an OpenAPI ApproversGroup to a domain ApproversGroup.
func ApproversGroupFromDTO(dto *openapi.TgvalidatordApproversGroup) model.ApproversGroup {
	if dto == nil {
		return model.ApproversGroup{}
	}
	return model.ApproversGroup{
		ExternalGroupID:   safeString(dto.ExternalGroupID),
		MinimumSignatures: safeInt64(dto.MinimumSignatures),
	}
}

// TrailFromDTO converts an OpenAPI Trail to a domain Trail.
func TrailFromDTO(dto *openapi.TgvalidatordTrail) model.Trail {
	if dto == nil {
		return model.Trail{}
	}

	trail := model.Trail{
		ID:             safeString(dto.Id),
		UserID:         safeString(dto.UserId),
		ExternalUserID: safeString(dto.ExternalUserId),
		Action:         safeString(dto.Action),
		Comment:        safeString(dto.Comment),
	}

	if dto.Date != nil {
		trail.Date = *dto.Date
	}

	return trail
}


// TrailsFromDTO converts a slice of OpenAPI Trails to domain Trails.
func TrailsFromDTO(dtos []openapi.TgvalidatordTrail) []model.Trail {
	if dtos == nil {
		return nil
	}
	trails := make([]model.Trail, len(dtos))
	for i := range dtos {
		trails[i] = TrailFromDTO(&dtos[i])
	}
	return trails
}

// WhitelistedAssetEnvelopeFromDTO converts an OpenAPI SignedWhitelistedContractAddressEnvelope to a domain WhitelistedAssetEnvelope.
// The envelope contains the raw data needed for 5-step verification.
func WhitelistedAssetEnvelopeFromDTO(dto *openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope) *model.WhitelistedAssetEnvelope {
	if dto == nil {
		return nil
	}

	envelope := &model.WhitelistedAssetEnvelope{
		ID:                  safeString(dto.Id),
		TenantID:            safeString(dto.TenantId),
		Blockchain:          safeString(dto.Blockchain),
		Network:             safeString(dto.Network),
		Status:              safeString(dto.Status),
		Action:              safeString(dto.Action),
		Rule:                safeString(dto.Rule),
		BusinessRuleEnabled: safeBool(dto.BusinessRuleEnabled),
		RulesContainer:      safeString(dto.RulesContainer),
		RulesSignatures:     safeString(dto.RulesSignatures),
	}

	// Convert metadata
	if dto.Metadata != nil {
		envelope.Metadata = WhitelistedAssetMetadataFromDTO(dto.Metadata)
	}

	// Convert signed contract address
	if dto.SignedContractAddress != nil {
		envelope.SignedContractAddress = SignedContractAddressFromDTO(dto.SignedContractAddress)
	}

	// Convert approvers
	if dto.Approvers != nil {
		envelope.Approvers = ApproversFromDTO(dto.Approvers)
	}

	// Convert attributes
	if dto.Attributes != nil {
		envelope.Attributes = make([]model.WhitelistedAssetAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			envelope.Attributes[i] = WhitelistedAssetAttributeFromDTO(&attr)
		}
	}

	// Convert trails
	if dto.Trails != nil {
		envelope.Trails = make([]model.Trail, len(dto.Trails))
		for i, trail := range dto.Trails {
			envelope.Trails[i] = TrailFromDTO(&trail)
		}
	}

	return envelope
}

// WhitelistedAssetAttributeFromDTO converts an OpenAPI WhitelistedContractAddressAttribute to a domain WhitelistedAssetAttribute.
func WhitelistedAssetAttributeFromDTO(dto *openapi.TgvalidatordWhitelistedContractAddressAttribute) model.WhitelistedAssetAttribute {
	if dto == nil {
		return model.WhitelistedAssetAttribute{}
	}
	return model.WhitelistedAssetAttribute{
		ID:          safeString(dto.Id),
		Key:         safeString(dto.Key),
		Value:       safeString(dto.Value),
		ContentType: safeString(dto.ContentType),
		Owner:       safeString(dto.Owner),
		Type:        safeString(dto.Type),
		Subtype:     safeString(dto.Subtype),
		IsFile:      safeBool(dto.Isfile),
	}
}
