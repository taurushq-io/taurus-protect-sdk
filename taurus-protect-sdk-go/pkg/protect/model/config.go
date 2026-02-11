package model

// TenantConfig represents the configuration of a tenant.
type TenantConfig struct {
	// TenantID is the unique identifier for the tenant.
	TenantID string `json:"tenant_id"`
	// SuperAdminMinimumSignatures is the minimum number of signatures required for super admin operations.
	SuperAdminMinimumSignatures string `json:"super_admin_minimum_signatures"`
	// BaseCurrency is the base currency for the tenant.
	BaseCurrency string `json:"base_currency"`
	// IsMFAMandatory indicates if multi-factor authentication is mandatory.
	IsMFAMandatory bool `json:"is_mfa_mandatory"`
	// ExcludeContainer indicates if container should be excluded.
	ExcludeContainer bool `json:"exclude_container"`
	// FeeLimitFactor is the factor used to limit transaction fees.
	FeeLimitFactor float32 `json:"fee_limit_factor"`
	// ProtectEngineVersion is the version of the protect engine.
	ProtectEngineVersion string `json:"protect_engine_version"`
	// RestrictSourcesForWhitelistedAddresses indicates if source restriction is enabled for whitelisted addresses.
	RestrictSourcesForWhitelistedAddresses bool `json:"restrict_sources_for_whitelisted_addresses"`
	// NFTMinting contains the NFT minting configuration.
	NFTMinting *NFTMintingConfig `json:"nft_minting,omitempty"`
	// IsProtectEngineCold indicates if the protect engine is in cold mode.
	IsProtectEngineCold bool `json:"is_protect_engine_cold"`
	// IsColdProtectEngineOffline indicates if the cold protect engine is offline.
	IsColdProtectEngineOffline bool `json:"is_cold_protect_engine_offline"`
	// IsPhysicalAirGapEnabled indicates if physical air gap is enabled.
	IsPhysicalAirGapEnabled bool `json:"is_physical_air_gap_enabled"`
}

// NFTMintingConfig represents the NFT minting configuration.
type NFTMintingConfig struct {
	// Enabled indicates if NFT minting is enabled.
	Enabled bool `json:"enabled"`
	// PublicBaseURL is the public base URL for NFT minting.
	PublicBaseURL string `json:"public_base_url"`
}
