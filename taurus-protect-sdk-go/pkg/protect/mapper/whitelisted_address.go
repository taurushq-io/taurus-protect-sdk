package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAddressFromDTO converts an OpenAPI SignedWhitelistedAddressEnvelope to a domain WhitelistedAddress.
func WhitelistedAddressFromDTO(dto *openapi.TgvalidatordSignedWhitelistedAddressEnvelope) *model.WhitelistedAddress {
	if dto == nil {
		return nil
	}

	addr := &model.WhitelistedAddress{
		ID:                 safeString(dto.Id),
		TenantID:           safeString(dto.TenantId),
		Blockchain:         safeString(dto.Blockchain),
		Network:            safeString(dto.Network),
		Status:             safeString(dto.Status),
		Action:             safeString(dto.Action),
		Rule:               safeString(dto.Rule),
		RulesContainer:     safeString(dto.RulesContainer),
		RulesContainerHash: safeString(dto.RulesContainerHash),
		RulesSignatures:    safeString(dto.RulesSignatures),
		VisibilityGroupID:  safeString(dto.VisibilityGroupID),
		TnParticipantID:    safeString(dto.TnParticipantID),
	}

	// Convert metadata
	if dto.Metadata != nil {
		addr.Metadata = WhitelistedAssetMetadataFromDTO(dto.Metadata)
		// SECURITY: Extract security-critical fields from PayloadAsString only.
		// IMPORTANT: Service MUST verify hash(PayloadAsString) == Metadata.Hash BEFORE calling this mapper.
		// The raw payload object could be tampered with while payloadAsString remains unchanged.
		if dto.Metadata.PayloadAsString != nil && *dto.Metadata.PayloadAsString != "" {
			parsed, err := helper.ParseWhitelistedAddressFromJSON(*dto.Metadata.PayloadAsString)
			if err == nil && parsed != nil {
				addr.Address = parsed.Address
				addr.Label = parsed.Label
				addr.Memo = parsed.Memo
				addr.CustomerId = parsed.CustomerId
				addr.ContractType = parsed.ContractType
				addr.AddressType = parsed.AddressType
				addr.ExchangeAccountId = parsed.ExchangeAccountId
				addr.LinkedInternalAddresses = parsed.LinkedInternalAddresses
				addr.LinkedWallets = parsed.LinkedWallets
			}
		}
	}

	// Convert signed address
	if dto.SignedAddress != nil {
		addr.SignedAddress = SignedWhitelistedAddressFromDTO(dto.SignedAddress)
	}

	// Convert scores
	if dto.Scores != nil {
		addr.Scores = AddressScoresFromDTO(dto.Scores)
	}

	// Convert trails - reuse existing TrailFromDTO
	if dto.Trails != nil {
		addr.Trails = make([]model.Trail, len(dto.Trails))
		for i, trail := range dto.Trails {
			addr.Trails[i] = TrailFromDTO(&trail)
		}
		// Extract createdAt from trails (find "created" action)
		for i := range addr.Trails {
			if addr.Trails[i].Action == "created" {
				createdAt := addr.Trails[i].Date
				addr.CreatedAt = &createdAt
				break
			}
		}
	}

	// Convert approvers - reuse existing ApproversFromDTO
	if dto.Approvers != nil {
		addr.Approvers = ApproversFromDTO(dto.Approvers)
	}

	// Convert attributes
	if dto.Attributes != nil {
		addr.Attributes = WhitelistedAddressAttributesFromDTO(dto.Attributes)
	}

	return addr
}

// WhitelistedAddressesFromDTO converts a slice of OpenAPI SignedWhitelistedAddressEnvelopes to domain WhitelistedAddresses.
func WhitelistedAddressesFromDTO(dtos []openapi.TgvalidatordSignedWhitelistedAddressEnvelope) []*model.WhitelistedAddress {
	if dtos == nil {
		return nil
	}
	addresses := make([]*model.WhitelistedAddress, len(dtos))
	for i := range dtos {
		addresses[i] = WhitelistedAddressFromDTO(&dtos[i])
	}
	return addresses
}

// SignedWhitelistedAddressFromDTO converts an OpenAPI SignedWhitelistedAddress to a domain SignedWhitelistedAddress.
func SignedWhitelistedAddressFromDTO(dto *openapi.TgvalidatordSignedWhitelistedAddress) *model.SignedWhitelistedAddress {
	if dto == nil {
		return nil
	}
	result := &model.SignedWhitelistedAddress{
		Payload: safeString(dto.Payload),
	}
	if dto.Signatures != nil {
		result.Signatures = make([]model.WhitelistSignature, len(dto.Signatures))
		for i, sig := range dto.Signatures {
			result.Signatures[i] = WhitelistSignatureFromDTO(&sig)
		}
	}
	return result
}

// AddressScoresFromDTO converts a slice of OpenAPI Scores to domain Scores.
func AddressScoresFromDTO(dtos []openapi.TgvalidatordScore) []model.Score {
	if dtos == nil {
		return nil
	}
	scores := make([]model.Score, len(dtos))
	for i := range dtos {
		scores[i] = AddressScoreFromDTO(&dtos[i])
	}
	return scores
}

// AddressScoreFromDTO converts an OpenAPI Score to a domain Score.
func AddressScoreFromDTO(dto *openapi.TgvalidatordScore) model.Score {
	if dto == nil {
		return model.Score{}
	}
	score := model.Score{
		ID:       safeString(dto.Id),
		Provider: safeString(dto.Provider),
		Type:     safeString(dto.Type),
		Score:    safeString(dto.Score),
	}
	if dto.UpdateDate != nil {
		score.UpdateDate = dto.UpdateDate.Format("2006-01-02T15:04:05Z")
	}
	return score
}

// WhitelistedAddressAttributesFromDTO converts a slice of OpenAPI WhitelistedAddressAttributes to domain WhitelistedAddressAttributes.
func WhitelistedAddressAttributesFromDTO(dtos []openapi.TgvalidatordWhitelistedAddressAttribute) []model.WhitelistedAddressAttribute {
	if dtos == nil {
		return nil
	}
	attrs := make([]model.WhitelistedAddressAttribute, len(dtos))
	for i := range dtos {
		attrs[i] = WhitelistedAddressAttributeFromDTO(&dtos[i])
	}
	return attrs
}

// WhitelistedAddressAttributeFromDTO converts an OpenAPI WhitelistedAddressAttribute to a domain WhitelistedAddressAttribute.
func WhitelistedAddressAttributeFromDTO(dto *openapi.TgvalidatordWhitelistedAddressAttribute) model.WhitelistedAddressAttribute {
	if dto == nil {
		return model.WhitelistedAddressAttribute{}
	}
	return model.WhitelistedAddressAttribute{
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

// WhitelistedAddressEnvelopeFromDTO converts an OpenAPI SignedWhitelistedAddressEnvelope to a domain WhitelistedAddressEnvelope.
// The envelope contains the raw data needed for 6-step verification.
func WhitelistedAddressEnvelopeFromDTO(dto *openapi.TgvalidatordSignedWhitelistedAddressEnvelope) *model.WhitelistedAddressEnvelope {
	if dto == nil {
		return nil
	}

	envelope := &model.WhitelistedAddressEnvelope{
		ID:              safeString(dto.Id),
		Blockchain:      safeString(dto.Blockchain),
		Network:         safeString(dto.Network),
		RulesContainer:  safeString(dto.RulesContainer),
		RulesSignatures: safeString(dto.RulesSignatures),
	}

	// Convert metadata
	if dto.Metadata != nil {
		envelope.Metadata = WhitelistedAssetMetadataFromDTO(dto.Metadata)
		// SECURITY: Extract from PayloadAsString only.
		// IMPORTANT: Service MUST verify hash BEFORE calling this mapper.
		if dto.Metadata.PayloadAsString != nil && *dto.Metadata.PayloadAsString != "" {
			parsed, err := helper.ParseWhitelistedAddressFromJSON(*dto.Metadata.PayloadAsString)
			if err == nil && parsed != nil {
				envelope.LinkedInternalAddresses = parsed.LinkedInternalAddresses
				envelope.LinkedWallets = parsed.LinkedWallets
			}
		}
	}

	// Convert signed address
	if dto.SignedAddress != nil {
		envelope.SignedAddress = SignedWhitelistedAddressFromDTO(dto.SignedAddress)
	}

	return envelope
}
