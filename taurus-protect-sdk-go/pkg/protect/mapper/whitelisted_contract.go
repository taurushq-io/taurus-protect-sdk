package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedContractFromDTO converts an OpenAPI SignedWhitelistedContractAddressEnvelope to a domain WhitelistedContract.
func WhitelistedContractFromDTO(dto *openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope) *model.WhitelistedContract {
	if dto == nil {
		return nil
	}

	contract := &model.WhitelistedContract{
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
		contract.Metadata = WhitelistedAssetMetadataFromDTO(dto.Metadata)
	}

	// Convert signed contract address
	if dto.SignedContractAddress != nil {
		contract.SignedContractAddress = SignedContractAddressFromDTO(dto.SignedContractAddress)
	}

	// Convert approvers
	if dto.Approvers != nil {
		contract.Approvers = ApproversFromDTO(dto.Approvers)
	}

	// Convert attributes
	if dto.Attributes != nil {
		contract.Attributes = WhitelistedContractAttributesFromDTO(dto.Attributes)
	}

	// Convert trails
	if dto.Trails != nil {
		contract.Trails = make([]model.Trail, len(dto.Trails))
		for i, trail := range dto.Trails {
			contract.Trails[i] = TrailFromDTO(&trail)
		}
	}

	return contract
}

// WhitelistedContractsFromDTO converts a slice of OpenAPI SignedWhitelistedContractAddressEnvelope to domain WhitelistedContracts.
func WhitelistedContractsFromDTO(dtos []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope) []*model.WhitelistedContract {
	if dtos == nil {
		return nil
	}
	contracts := make([]*model.WhitelistedContract, len(dtos))
	for i := range dtos {
		contracts[i] = WhitelistedContractFromDTO(&dtos[i])
	}
	return contracts
}

// WhitelistedContractAttributeFromDTO converts an OpenAPI WhitelistedContractAddressAttribute to a domain WhitelistedContractAttribute.
func WhitelistedContractAttributeFromDTO(dto *openapi.TgvalidatordWhitelistedContractAddressAttribute) model.WhitelistedContractAttribute {
	if dto == nil {
		return model.WhitelistedContractAttribute{}
	}
	return model.WhitelistedContractAttribute{
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

// WhitelistedContractAttributesFromDTO converts a slice of OpenAPI WhitelistedContractAddressAttribute to domain WhitelistedContractAttributes.
func WhitelistedContractAttributesFromDTO(dtos []openapi.TgvalidatordWhitelistedContractAddressAttribute) []model.WhitelistedContractAttribute {
	if dtos == nil {
		return nil
	}
	attrs := make([]model.WhitelistedContractAttribute, len(dtos))
	for i := range dtos {
		attrs[i] = WhitelistedContractAttributeFromDTO(&dtos[i])
	}
	return attrs
}
