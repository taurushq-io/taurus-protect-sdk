package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// UserDeviceService provides user device pairing operations.
type UserDeviceService struct {
	api       *openapi.UserDeviceAPIService
	errMapper *ErrorMapper
}

// NewUserDeviceService creates a new UserDeviceService.
func NewUserDeviceService(client *openapi.APIClient) *UserDeviceService {
	return &UserDeviceService{
		api:       client.UserDeviceAPI,
		errMapper: NewErrorMapper(),
	}
}

// CreatePairing creates a new user device pairing request (Step 1).
// This initiates the pairing process and returns a pairingID that should be used
// in subsequent steps to complete the pairing.
func (s *UserDeviceService) CreatePairing(ctx context.Context) (*model.CreatePairingResult, error) {
	resp, httpResp, err := s.api.UserDeviceServiceCreateUserDevicePairing(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.CreatePairingResultFromDTO(resp), nil
}

// StartPairing starts a user device pairing request (Step 2).
// This step provides the nonce and public key from the user's device.
// The nonce is a 6-digit number, and the public key is an ECDSA key encoded in base64.
func (s *UserDeviceService) StartPairing(ctx context.Context, pairingID string, req *model.StartPairingRequest) error {
	if pairingID == "" {
		return fmt.Errorf("pairingID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.Nonce == "" {
		return fmt.Errorf("nonce is required")
	}
	if req.PublicKey == "" {
		return fmt.Errorf("publicKey is required")
	}

	body := mapper.StartPairingRequestToDTO(req)

	_, httpResp, err := s.api.UserDeviceServiceStartUserDevicePairing(ctx, pairingID).
		Body(*body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ApprovePairing approves a user device pairing request (Step 3).
// This is the final step that completes the pairing process.
// The nonce must match the one used in the start pairing step.
func (s *UserDeviceService) ApprovePairing(ctx context.Context, pairingID string, req *model.ApprovePairingRequest) error {
	if pairingID == "" {
		return fmt.Errorf("pairingID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.Nonce == "" {
		return fmt.Errorf("nonce is required")
	}

	body := mapper.ApprovePairingRequestToDTO(req)

	_, httpResp, err := s.api.UserDeviceServiceApproveUserDevicePairing(ctx, pairingID).
		Body(*body).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// GetPairingStatus retrieves the status of a user device pairing request.
// The nonce parameter is required and must be the same 6-digit number used to start the pairing.
func (s *UserDeviceService) GetPairingStatus(ctx context.Context, pairingID string, nonce string) (*model.PairingStatusResult, error) {
	if pairingID == "" {
		return nil, fmt.Errorf("pairingID cannot be empty")
	}
	if nonce == "" {
		return nil, fmt.Errorf("nonce cannot be empty")
	}

	resp, httpResp, err := s.api.UserDeviceServiceGetUserDevicePairingStatus(ctx, pairingID).
		Nonce(nonce).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.PairingStatusResultFromDTO(resp), nil
}
