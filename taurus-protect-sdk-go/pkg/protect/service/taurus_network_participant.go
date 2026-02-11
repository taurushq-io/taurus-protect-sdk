package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TaurusNetworkParticipantService provides Taurus-NETWORK participant management operations.
type TaurusNetworkParticipantService struct {
	api       *openapi.TaurusNetworkParticipantAPIService
	errMapper *ErrorMapper
}

// NewTaurusNetworkParticipantService creates a new TaurusNetworkParticipantService.
func NewTaurusNetworkParticipantService(client *openapi.APIClient) *TaurusNetworkParticipantService {
	return &TaurusNetworkParticipantService{
		api:       client.TaurusNetworkParticipantAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetMyParticipant returns the current participant linked to the tenant with additional Taurus-NETWORK settings.
func (s *TaurusNetworkParticipantService) GetMyParticipant(ctx context.Context) (*taurusnetwork.GetMyParticipantResult, error) {
	req := s.api.TaurusNetworkServiceGetMyParticipant(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.GetMyParticipantResult{}

	if resp.Result != nil {
		result.Participant = mapper.TnParticipantFromDTO(resp.Result)
	}

	if resp.Settings != nil {
		result.Settings = mapper.TnParticipantSettingsFromDTO(resp.Settings)
	}

	return result, nil
}

// GetParticipant returns a Taurus-NETWORK participant based on the provided ID.
func (s *TaurusNetworkParticipantService) GetParticipant(ctx context.Context, participantID string, opts *taurusnetwork.GetParticipantOptions) (*taurusnetwork.TnParticipant, error) {
	req := s.api.TaurusNetworkServiceGetParticipant(ctx, participantID)

	if opts != nil {
		if opts.IncludeTotalPledgesValuation {
			req = req.IncludeTotalPledgesValuation(opts.IncludeTotalPledgesValuation)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.TnParticipantFromDTO(resp.Result), nil
}

// ListParticipants returns a list of visible Taurus-NETWORK participants connected to the current participant.
func (s *TaurusNetworkParticipantService) ListParticipants(ctx context.Context, opts *taurusnetwork.ListParticipantsOptions) (*taurusnetwork.ListParticipantsResult, error) {
	req := s.api.TaurusNetworkServiceGetParticipants(ctx)

	if opts != nil {
		if len(opts.ParticipantIDs) > 0 {
			req = req.ParticipantIDs(opts.ParticipantIDs)
		}
		if opts.IncludeTotalPledgesValuation {
			req = req.IncludeTotalPledgesValuation(opts.IncludeTotalPledgesValuation)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListParticipantsResult{
		Participants: mapper.TnParticipantsFromDTO(resp.Result),
	}

	return result, nil
}

// CreateParticipantAttribute creates an attribute for a given participant.
// The attribute can also be shared to the linked Taurus-NETWORK participant.
// Required role: Admin.
func (s *TaurusNetworkParticipantService) CreateParticipantAttribute(ctx context.Context, participantID string, req *taurusnetwork.CreateParticipantAttributeRequest) error {
	body := mapper.CreateParticipantAttributeBodyToDTO(req)

	apiReq := s.api.TaurusNetworkServiceCreateParticipantAttribute(ctx, participantID)
	apiReq = apiReq.Body(body)

	_, httpResp, err := apiReq.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// DeleteParticipantAttribute deletes an attribute for a given participant.
// If the attribute was shared to the linked Taurus-NETWORK participant, a message will be sent to the participant.
// Required role: Admin.
func (s *TaurusNetworkParticipantService) DeleteParticipantAttribute(ctx context.Context, participantID string, attributeID string) error {
	req := s.api.TaurusNetworkServiceDeleteParticipantAttribute(ctx, participantID, attributeID)

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
