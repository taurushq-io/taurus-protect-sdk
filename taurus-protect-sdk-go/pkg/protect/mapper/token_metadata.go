package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// CryptoPunkMetadataFromDTO converts an OpenAPI CryptoPunkMetadata to a domain CryptoPunkMetadata.
func CryptoPunkMetadataFromDTO(dto *openapi.TgvalidatordCryptoPunkMetadata) *model.CryptoPunkMetadata {
	if dto == nil {
		return nil
	}

	return &model.CryptoPunkMetadata{
		PunkID:     safeString(dto.PunkId),
		Attributes: safeString(dto.PunkAttributes),
		Image:      safeString(dto.Image),
	}
}

// ERCTokenMetadataFromDTO converts an OpenAPI ERCTokenMetadata to a domain ERCTokenMetadata.
func ERCTokenMetadataFromDTO(dto *openapi.TgvalidatordERCTokenMetadata) *model.ERCTokenMetadata {
	if dto == nil {
		return nil
	}

	return &model.ERCTokenMetadata{
		Name:        safeString(dto.Name),
		Description: safeString(dto.Description),
		Decimals:    safeString(dto.Decimals),
		DataType:    safeString(dto.DataType),
		Base64Data:  safeString(dto.Base64Data),
		URI:         safeString(dto.Uri),
	}
}
