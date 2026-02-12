package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TenantConfigFromDTO converts an OpenAPI TgvalidatordGetConfigTenantReply to a domain TenantConfig.
func TenantConfigFromDTO(dto *openapi.TgvalidatordGetConfigTenantReply) *model.TenantConfig {
	if dto == nil || dto.Config == nil {
		return nil
	}

	config := dto.Config
	return &model.TenantConfig{
		TenantID:                               safeString(config.TenantId),
		SuperAdminMinimumSignatures:            safeString(config.SuperAdminMinimumSignatures),
		BaseCurrency:                           safeString(config.BaseCurrency),
		IsMFAMandatory:                         safeBool(config.IsMFAMandatory),
		ExcludeContainer:                       safeBool(config.ExcludeContainer),
		FeeLimitFactor:                         safeFloat32(config.FeeLimitFactor),
		ProtectEngineVersion:                   safeString(config.ProtectEngineVersion),
		RestrictSourcesForWhitelistedAddresses: safeBool(config.RestrictSourcesForWhitelistedAddresses),
		NFTMinting:                             NFTMintingConfigFromDTO(config.NftMinting),
		IsProtectEngineCold:                    safeBool(config.IsProtectEngineCold),
		IsColdProtectEngineOffline:             safeBool(config.IsColdProtectEngineOffline),
		IsPhysicalAirGapEnabled:                safeBool(config.IsPhysicalAirGapEnabled),
	}
}

// NFTMintingConfigFromDTO converts an OpenAPI TenantConfigNFTMinting to a domain NFTMintingConfig.
func NFTMintingConfigFromDTO(dto *openapi.TenantConfigNFTMinting) *model.NFTMintingConfig {
	if dto == nil {
		return nil
	}

	return &model.NFTMintingConfig{
		Enabled:       safeBool(dto.Enabled),
		PublicBaseURL: safeString(dto.PublicBaseURL),
	}
}

