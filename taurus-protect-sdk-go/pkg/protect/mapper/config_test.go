package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestTenantConfigFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetConfigTenantReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "nil config returns nil",
			dto:  &openapi.TgvalidatordGetConfigTenantReply{Config: nil},
		},
		{
			name: "empty config returns tenant config with zero values",
			dto: &openapi.TgvalidatordGetConfigTenantReply{
				Config: &openapi.TgvalidatordTenantConfig{},
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetConfigTenantReply {
				tenantId := "tenant-123"
				superAdminMinSigs := "3"
				baseCurrency := "USD"
				isMFAMandatory := true
				excludeContainer := false
				feeLimitFactor := float32(1.5)
				protectEngineVersion := "2.0.0"
				restrictSources := true
				isProtectEngineCold := false
				isColdEngineOffline := false
				isAirGapEnabled := true
				nftEnabled := true
				nftBaseURL := "https://nft.example.com"

				return &openapi.TgvalidatordGetConfigTenantReply{
					Config: &openapi.TgvalidatordTenantConfig{
						TenantId:                               &tenantId,
						SuperAdminMinimumSignatures:            &superAdminMinSigs,
						BaseCurrency:                           &baseCurrency,
						IsMFAMandatory:                         &isMFAMandatory,
						ExcludeContainer:                       &excludeContainer,
						FeeLimitFactor:                         &feeLimitFactor,
						ProtectEngineVersion:                   &protectEngineVersion,
						RestrictSourcesForWhitelistedAddresses: &restrictSources,
						IsProtectEngineCold:                    &isProtectEngineCold,
						IsColdProtectEngineOffline:             &isColdEngineOffline,
						IsPhysicalAirGapEnabled:                &isAirGapEnabled,
						NftMinting: &openapi.TenantConfigNFTMinting{
							Enabled:       &nftEnabled,
							PublicBaseURL: &nftBaseURL,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TenantConfigFromDTO(tt.dto)
			if tt.dto == nil || tt.dto.Config == nil {
				if got != nil {
					t.Errorf("TenantConfigFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TenantConfigFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			config := tt.dto.Config
			if config.TenantId != nil && got.TenantID != *config.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *config.TenantId)
			}
			if config.SuperAdminMinimumSignatures != nil && got.SuperAdminMinimumSignatures != *config.SuperAdminMinimumSignatures {
				t.Errorf("SuperAdminMinimumSignatures = %v, want %v", got.SuperAdminMinimumSignatures, *config.SuperAdminMinimumSignatures)
			}
			if config.BaseCurrency != nil && got.BaseCurrency != *config.BaseCurrency {
				t.Errorf("BaseCurrency = %v, want %v", got.BaseCurrency, *config.BaseCurrency)
			}
			if config.IsMFAMandatory != nil && got.IsMFAMandatory != *config.IsMFAMandatory {
				t.Errorf("IsMFAMandatory = %v, want %v", got.IsMFAMandatory, *config.IsMFAMandatory)
			}
			if config.ExcludeContainer != nil && got.ExcludeContainer != *config.ExcludeContainer {
				t.Errorf("ExcludeContainer = %v, want %v", got.ExcludeContainer, *config.ExcludeContainer)
			}
			if config.FeeLimitFactor != nil && got.FeeLimitFactor != *config.FeeLimitFactor {
				t.Errorf("FeeLimitFactor = %v, want %v", got.FeeLimitFactor, *config.FeeLimitFactor)
			}
			if config.ProtectEngineVersion != nil && got.ProtectEngineVersion != *config.ProtectEngineVersion {
				t.Errorf("ProtectEngineVersion = %v, want %v", got.ProtectEngineVersion, *config.ProtectEngineVersion)
			}
			if config.RestrictSourcesForWhitelistedAddresses != nil && got.RestrictSourcesForWhitelistedAddresses != *config.RestrictSourcesForWhitelistedAddresses {
				t.Errorf("RestrictSourcesForWhitelistedAddresses = %v, want %v", got.RestrictSourcesForWhitelistedAddresses, *config.RestrictSourcesForWhitelistedAddresses)
			}
			if config.IsProtectEngineCold != nil && got.IsProtectEngineCold != *config.IsProtectEngineCold {
				t.Errorf("IsProtectEngineCold = %v, want %v", got.IsProtectEngineCold, *config.IsProtectEngineCold)
			}
			if config.IsColdProtectEngineOffline != nil && got.IsColdProtectEngineOffline != *config.IsColdProtectEngineOffline {
				t.Errorf("IsColdProtectEngineOffline = %v, want %v", got.IsColdProtectEngineOffline, *config.IsColdProtectEngineOffline)
			}
			if config.IsPhysicalAirGapEnabled != nil && got.IsPhysicalAirGapEnabled != *config.IsPhysicalAirGapEnabled {
				t.Errorf("IsPhysicalAirGapEnabled = %v, want %v", got.IsPhysicalAirGapEnabled, *config.IsPhysicalAirGapEnabled)
			}
			// Verify NFT minting is mapped if present
			if config.NftMinting != nil {
				if got.NFTMinting == nil {
					t.Error("NFTMinting should not be nil when DTO has nft minting config")
				}
			}
		})
	}
}

func TestNFTMintingConfigFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TenantConfigNFTMinting
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns NFT minting config with zero values",
			dto:  &openapi.TenantConfigNFTMinting{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TenantConfigNFTMinting {
				enabled := true
				publicBaseURL := "https://nft.example.com"
				return &openapi.TenantConfigNFTMinting{
					Enabled:       &enabled,
					PublicBaseURL: &publicBaseURL,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NFTMintingConfigFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("NFTMintingConfigFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("NFTMintingConfigFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Enabled != nil && got.Enabled != *tt.dto.Enabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, *tt.dto.Enabled)
			}
			if tt.dto.PublicBaseURL != nil && got.PublicBaseURL != *tt.dto.PublicBaseURL {
				t.Errorf("PublicBaseURL = %v, want %v", got.PublicBaseURL, *tt.dto.PublicBaseURL)
			}
		})
	}
}

func TestNFTMintingConfigFromDTO_EnabledField(t *testing.T) {
	tests := []struct {
		name        string
		enabled     *bool
		wantEnabled bool
	}{
		{
			name:        "nil enabled defaults to false",
			enabled:     nil,
			wantEnabled: false,
		},
		{
			name:        "true enabled",
			enabled:     boolPtr(true),
			wantEnabled: true,
		},
		{
			name:        "false enabled",
			enabled:     boolPtr(false),
			wantEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TenantConfigNFTMinting{
				Enabled: tt.enabled,
			}
			got := NFTMintingConfigFromDTO(dto)
			if got.Enabled != tt.wantEnabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, tt.wantEnabled)
			}
		})
	}
}

func TestTenantConfigFromDTO_NilNFTMinting(t *testing.T) {
	baseCurrency := "EUR"
	dto := &openapi.TgvalidatordGetConfigTenantReply{
		Config: &openapi.TgvalidatordTenantConfig{
			BaseCurrency: &baseCurrency,
			NftMinting:   nil,
		},
	}

	got := TenantConfigFromDTO(dto)
	if got == nil {
		t.Fatal("TenantConfigFromDTO() returned nil for non-nil input")
	}
	if got.NFTMinting != nil {
		t.Errorf("NFTMinting should be nil when DTO nft minting is nil, got %v", got.NFTMinting)
	}
	if got.BaseCurrency != "EUR" {
		t.Errorf("BaseCurrency = %v, want EUR", got.BaseCurrency)
	}
}

func TestTenantConfigFromDTO_BooleanFields(t *testing.T) {
	tests := []struct {
		name                 string
		isMFAMandatory       *bool
		wantIsMFAMandatory   bool
		isProtectEngineCold  *bool
		wantProtectEngineCold bool
	}{
		{
			name:                 "nil booleans default to false",
			isMFAMandatory:       nil,
			wantIsMFAMandatory:   false,
			isProtectEngineCold:  nil,
			wantProtectEngineCold: false,
		},
		{
			name:                 "true values",
			isMFAMandatory:       boolPtr(true),
			wantIsMFAMandatory:   true,
			isProtectEngineCold:  boolPtr(true),
			wantProtectEngineCold: true,
		},
		{
			name:                 "false values",
			isMFAMandatory:       boolPtr(false),
			wantIsMFAMandatory:   false,
			isProtectEngineCold:  boolPtr(false),
			wantProtectEngineCold: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordGetConfigTenantReply{
				Config: &openapi.TgvalidatordTenantConfig{
					IsMFAMandatory:      tt.isMFAMandatory,
					IsProtectEngineCold: tt.isProtectEngineCold,
				},
			}
			got := TenantConfigFromDTO(dto)
			if got.IsMFAMandatory != tt.wantIsMFAMandatory {
				t.Errorf("IsMFAMandatory = %v, want %v", got.IsMFAMandatory, tt.wantIsMFAMandatory)
			}
			if got.IsProtectEngineCold != tt.wantProtectEngineCold {
				t.Errorf("IsProtectEngineCold = %v, want %v", got.IsProtectEngineCold, tt.wantProtectEngineCold)
			}
		})
	}
}

func TestSafeFloat32(t *testing.T) {
	tests := []struct {
		name  string
		input *float32
		want  float32
	}{
		{
			name:  "nil returns 0",
			input: nil,
			want:  0,
		},
		{
			name:  "positive value",
			input: float32Ptr(1.5),
			want:  1.5,
		},
		{
			name:  "negative value",
			input: float32Ptr(-2.5),
			want:  -2.5,
		},
		{
			name:  "zero value",
			input: float32Ptr(0),
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeFloat32(tt.input)
			if got != tt.want {
				t.Errorf("safeFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

// float32Ptr returns a pointer to a float32 value.
func float32Ptr(f float32) *float32 {
	return &f
}
