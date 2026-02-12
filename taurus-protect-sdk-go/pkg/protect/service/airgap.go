package service

import (
	"context"
	"fmt"
	"io"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AirGapService provides air-gap operations for cold HSM integration.
type AirGapService struct {
	api       *openapi.AirGapAPIService
	errMapper *ErrorMapper
}

// NewAirGapService creates a new AirGapService.
func NewAirGapService(client *openapi.APIClient) *AirGapService {
	return &AirGapService{
		api:       client.AirGapAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetOutgoingAirGap exports HSM ready requests for cold HSM.
// This method returns the payload to be transmitted to the cold HSM.
func (s *AirGapService) GetOutgoingAirGap(ctx context.Context, req *model.GetOutgoingAirGapRequest) (*model.GetOutgoingAirGapResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if len(req.RequestIDs) == 0 && len(req.AddressIDs) == 0 {
		return nil, fmt.Errorf("either request IDs or address IDs must be provided")
	}

	// Build the OpenAPI request
	apiReq := openapi.TgvalidatordGetOutgoingAirGapRequest{}

	// Set request IDs if provided
	if len(req.RequestIDs) > 0 {
		requests := openapi.TgvalidatordGetOutgoingAirGapRequestRequests{
			Ids: req.RequestIDs,
		}
		if req.RequestSignature != "" {
			requests.Signature = &req.RequestSignature
		}
		apiReq.Requests = &requests
	}

	// Set address IDs if provided
	if len(req.AddressIDs) > 0 {
		addresses := openapi.GetOutgoingAirGapRequestAddresses{
			Ids: req.AddressIDs,
		}
		apiReq.Addresses = &addresses
	}

	// Execute the API call
	file, httpResp, err := s.api.AirGapServiceGetOutgoingAirGap(ctx).
		Body(apiReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// Read the binary data from the file
	if file == nil {
		return nil, fmt.Errorf("no data returned from air-gap export")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read air-gap data: %w", err)
	}

	return &model.GetOutgoingAirGapResult{
		Data: data,
	}, nil
}

// SubmitIncomingAirGap imports signed requests from the cold HSM.
// This method accepts an envelope of signed requests from the cold HSM.
func (s *AirGapService) SubmitIncomingAirGap(ctx context.Context, req *model.SubmitIncomingAirGapRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.Payload == "" {
		return fmt.Errorf("payload is required")
	}
	if req.Signature == "" {
		return fmt.Errorf("signature is required")
	}

	// Build the OpenAPI request
	apiReq := openapi.TgvalidatordSubmitIncomingAirGapRequest{
		Payload:   &req.Payload,
		Signature: &req.Signature,
	}

	// Execute the API call
	_, httpResp, err := s.api.AirGapServiceSubmitIncomingAirGap(ctx).
		Body(apiReq).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
