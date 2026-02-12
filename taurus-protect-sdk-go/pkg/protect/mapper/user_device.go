package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// CreatePairingResultFromDTO converts an OpenAPI CreateUserDevicePairingReply to a domain CreatePairingResult.
func CreatePairingResultFromDTO(dto *openapi.TgvalidatordCreateUserDevicePairingReply) *model.CreatePairingResult {
	if dto == nil {
		return nil
	}

	return &model.CreatePairingResult{
		PairingID: dto.PairingID,
	}
}

// PairingStatusResultFromDTO converts an OpenAPI UserDevicePairingInfo to a domain PairingStatusResult.
func PairingStatusResultFromDTO(dto *openapi.TgvalidatordUserDevicePairingInfo) *model.PairingStatusResult {
	if dto == nil {
		return nil
	}

	result := &model.PairingStatusResult{
		Status:    model.UserDevicePairingStatus(dto.Status),
		PairingID: dto.PairingID,
	}

	if dto.ApiKey != nil {
		result.APIKey = *dto.ApiKey
	}

	return result
}

// StartPairingRequestToDTO converts a domain StartPairingRequest to an OpenAPI request body.
func StartPairingRequestToDTO(req *model.StartPairingRequest) *openapi.UserDeviceServiceStartUserDevicePairingBody {
	if req == nil {
		return nil
	}

	return &openapi.UserDeviceServiceStartUserDevicePairingBody{
		Nonce:     req.Nonce,
		PublicKey: req.PublicKey,
	}
}

// ApprovePairingRequestToDTO converts a domain ApprovePairingRequest to an OpenAPI request body.
func ApprovePairingRequestToDTO(req *model.ApprovePairingRequest) *openapi.UserDeviceServiceApproveUserDevicePairingBody {
	if req == nil {
		return nil
	}

	return &openapi.UserDeviceServiceApproveUserDevicePairingBody{
		Nonce: req.Nonce,
	}
}
