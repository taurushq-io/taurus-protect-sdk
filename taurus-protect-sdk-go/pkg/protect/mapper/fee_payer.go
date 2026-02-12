package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FeePayerFromDTO converts an OpenAPI FeePayerEnvelope to a domain FeePayer.
func FeePayerFromDTO(dto *openapi.TgvalidatordFeePayerEnvelope) *model.FeePayer {
	if dto == nil {
		return nil
	}

	feePayer := &model.FeePayer{
		ID:         safeString(dto.Id),
		TenantID:   safeString(dto.TenantId),
		Blockchain: safeString(dto.Blockchain),
		Network:    safeString(dto.Network),
		Name:       safeString(dto.Name),
	}

	// Convert creation date
	if dto.CreationDate != nil {
		feePayer.CreationDate = *dto.CreationDate
	}

	// Convert nested fee payer details
	if dto.FeePayer != nil {
		feePayer.ETH = FeePayerETHFromDTO(dto.FeePayer)
	}

	return feePayer
}

// FeePayersFromDTO converts a slice of OpenAPI FeePayerEnvelope to domain FeePayers.
func FeePayersFromDTO(dtos []openapi.TgvalidatordFeePayerEnvelope) []*model.FeePayer {
	if dtos == nil {
		return nil
	}
	feePayers := make([]*model.FeePayer, len(dtos))
	for i := range dtos {
		feePayers[i] = FeePayerFromDTO(&dtos[i])
	}
	return feePayers
}

// FeePayerETHFromDTO converts an OpenAPI TgvalidatordFeePayer to a domain FeePayerETH.
func FeePayerETHFromDTO(dto *openapi.TgvalidatordFeePayer) *model.FeePayerETH {
	if dto == nil {
		return nil
	}

	eth := &model.FeePayerETH{
		Blockchain: safeString(dto.Blockchain),
	}

	// Convert ETH-specific configuration
	if dto.Eth != nil {
		eth.Kind = safeString(dto.Eth.Kind)
		eth.RemoteEncrypted = safeString(dto.Eth.RemoteEncrypted)

		if dto.Eth.Local != nil {
			eth.Local = ETHLocalConfigFromDTO(dto.Eth.Local)
		}

		if dto.Eth.Remote != nil {
			eth.Remote = ETHRemoteConfigFromDTO(dto.Eth.Remote)
		}
	}

	return eth
}

// ETHLocalConfigFromDTO converts an OpenAPI ETHLocal to a domain ETHLocalConfig.
func ETHLocalConfigFromDTO(dto *openapi.ETHLocal) *model.ETHLocalConfig {
	if dto == nil {
		return nil
	}

	config := &model.ETHLocalConfig{
		AddressID:          safeString(dto.AddressId),
		ForwarderAddressID: safeString(dto.ForwarderAddressId),
		AutoApprove:        safeBool(dto.AutoApprove),
		CreatorAddressID:   safeString(dto.CreatorAddressId),
		DomainSeparator:    safeString(dto.DomainSeparator),
	}

	// Convert forwarder kind
	if dto.ForwarderKind != nil {
		config.ForwarderKind = string(*dto.ForwarderKind)
	}

	return config
}

// ETHRemoteConfigFromDTO converts an OpenAPI ETHRemote to a domain ETHRemoteConfig.
func ETHRemoteConfigFromDTO(dto *openapi.ETHRemote) *model.ETHRemoteConfig {
	if dto == nil {
		return nil
	}

	config := &model.ETHRemoteConfig{
		URL:                safeString(dto.Url),
		Username:           safeString(dto.Username),
		Password:           safeString(dto.Password),
		PrivateKey:         safeString(dto.PrivateKey),
		FromAddressID:      safeString(dto.FromAddressId),
		ForwarderAddress:   safeString(dto.ForwarderAddress),
		ForwarderAddressID: safeString(dto.ForwarderAddressId),
		CreatorAddress:     safeString(dto.CreatorAddress),
		CreatorAddressID:   safeString(dto.CreatorAddressId),
		DomainSeparator:    safeString(dto.DomainSeparator),
	}

	// Convert forwarder kind
	if dto.ForwarderKind != nil {
		config.ForwarderKind = string(*dto.ForwarderKind)
	}

	return config
}

// ChecksumRequestToDTO converts a domain ChecksumRequest to an OpenAPI TgvalidatordGetChecksumRequest.
func ChecksumRequestToDTO(req *model.ChecksumRequest) openapi.TgvalidatordGetChecksumRequest {
	dto := openapi.TgvalidatordGetChecksumRequest{}
	if req != nil {
		dto.Data = stringPtr(req.Data)
	}
	return dto
}
