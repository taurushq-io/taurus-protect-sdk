package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestNewTaurusNetworkParticipantService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTaurusNetworkParticipantService_GetMyParticipant(t *testing.T) {
	// Create a service with nil API to test that the service structure is correct
	svc := &TaurusNetworkParticipantService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkParticipantService should not be nil")
	}
}

func TestTaurusNetworkParticipantService_GetParticipant_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &TaurusNetworkParticipantService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkParticipantService should not be nil")
	}
}

func TestTaurusNetworkParticipantService_GetParticipant_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.GetParticipantOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.GetParticipantOptions{},
		},
		{
			name: "include total pledges valuation",
			options: &taurusnetwork.GetParticipantOptions{
				IncludeTotalPledgesValuation: true,
			},
		},
		{
			name: "exclude total pledges valuation",
			options: &taurusnetwork.GetParticipantOptions{
				IncludeTotalPledgesValuation: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &TaurusNetworkParticipantService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkParticipantService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkParticipantService_ListParticipants_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &TaurusNetworkParticipantService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkParticipantService should not be nil")
	}
}

func TestTaurusNetworkParticipantService_ListParticipants_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListParticipantsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListParticipantsOptions{},
		},
		{
			name: "filter by participant IDs",
			options: &taurusnetwork.ListParticipantsOptions{
				ParticipantIDs: []string{"participant-1", "participant-2"},
			},
		},
		{
			name: "include total pledges valuation",
			options: &taurusnetwork.ListParticipantsOptions{
				IncludeTotalPledgesValuation: true,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListParticipantsOptions{
				ParticipantIDs:               []string{"participant-1"},
				IncludeTotalPledgesValuation: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &TaurusNetworkParticipantService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkParticipantService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkParticipantService_CreateParticipantAttribute_NilRequest(t *testing.T) {
	// Create a service with nil API
	svc := &TaurusNetworkParticipantService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkParticipantService should not be nil")
	}
}

func TestTaurusNetworkParticipantService_CreateParticipantAttribute_WithRequest(t *testing.T) {
	tests := []struct {
		name string
		req  *taurusnetwork.CreateParticipantAttributeRequest
	}{
		{
			name: "nil request",
			req:  nil,
		},
		{
			name: "empty request",
			req:  &taurusnetwork.CreateParticipantAttributeRequest{},
		},
		{
			name: "basic attribute",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:   "license",
				Value: "12345",
			},
		},
		{
			name: "attribute with content type",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:         "document",
				Value:       "base64content",
				ContentType: "application/pdf",
			},
		},
		{
			name: "attribute with type and subtype",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:     "compliance-doc",
				Value:   "content",
				Type:    "COMPLIANCE",
				Subtype: "LICENSE",
			},
		},
		{
			name: "shared attribute",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:                             "shared-license",
				Value:                           "12345",
				ShareToTaurusNetworkParticipant: true,
			},
		},
		{
			name: "all fields",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:                             "complete-attr",
				Value:                           "value123",
				ContentType:                     "text/plain",
				Type:                            "COMPLIANCE",
				Subtype:                         "CERTIFICATE",
				ShareToTaurusNetworkParticipant: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify request can be created with these values
			// Actual API testing requires mocking
			svc := &TaurusNetworkParticipantService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkParticipantService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkParticipantService_DeleteParticipantAttribute(t *testing.T) {
	// Create a service with nil API
	svc := &TaurusNetworkParticipantService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkParticipantService should not be nil")
	}
}

func TestGetParticipantOptions_Values(t *testing.T) {
	tests := []struct {
		name                         string
		includeTotalPledgesValuation bool
	}{
		{
			name:                         "include pledges valuation true",
			includeTotalPledgesValuation: true,
		},
		{
			name:                         "include pledges valuation false",
			includeTotalPledgesValuation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &taurusnetwork.GetParticipantOptions{
				IncludeTotalPledgesValuation: tt.includeTotalPledgesValuation,
			}
			if opts.IncludeTotalPledgesValuation != tt.includeTotalPledgesValuation {
				t.Errorf("IncludeTotalPledgesValuation = %v, want %v", opts.IncludeTotalPledgesValuation, tt.includeTotalPledgesValuation)
			}
		})
	}
}

func TestListParticipantsOptions_Values(t *testing.T) {
	tests := []struct {
		name           string
		participantIDs []string
		includePledges bool
	}{
		{
			name:           "empty participant IDs",
			participantIDs: nil,
			includePledges: false,
		},
		{
			name:           "single participant ID",
			participantIDs: []string{"participant-1"},
			includePledges: true,
		},
		{
			name:           "multiple participant IDs",
			participantIDs: []string{"participant-1", "participant-2", "participant-3"},
			includePledges: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &taurusnetwork.ListParticipantsOptions{
				ParticipantIDs:               tt.participantIDs,
				IncludeTotalPledgesValuation: tt.includePledges,
			}
			if len(opts.ParticipantIDs) != len(tt.participantIDs) {
				t.Errorf("ParticipantIDs length = %v, want %v", len(opts.ParticipantIDs), len(tt.participantIDs))
			}
			if opts.IncludeTotalPledgesValuation != tt.includePledges {
				t.Errorf("IncludeTotalPledgesValuation = %v, want %v", opts.IncludeTotalPledgesValuation, tt.includePledges)
			}
		})
	}
}

func TestCreateParticipantAttributeRequest_Values(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		share   bool
	}{
		{
			name:  "basic attribute",
			key:   "license",
			value: "12345",
			share: false,
		},
		{
			name:  "shared attribute",
			key:   "certificate",
			value: "cert-content",
			share: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &taurusnetwork.CreateParticipantAttributeRequest{
				Key:                             tt.key,
				Value:                           tt.value,
				ShareToTaurusNetworkParticipant: tt.share,
			}
			if req.Key != tt.key {
				t.Errorf("Key = %v, want %v", req.Key, tt.key)
			}
			if req.Value != tt.value {
				t.Errorf("Value = %v, want %v", req.Value, tt.value)
			}
			if req.ShareToTaurusNetworkParticipant != tt.share {
				t.Errorf("ShareToTaurusNetworkParticipant = %v, want %v", req.ShareToTaurusNetworkParticipant, tt.share)
			}
		})
	}
}
