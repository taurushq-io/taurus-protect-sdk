package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Attribute mappers only. The envelope mappers that lived here built a
// model.WhitelistedContract straight from the DTO with no verification, which is what
// made the verified reader on WhitelistedAssetService avoidable.

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
