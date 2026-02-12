package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewAirGapService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestAirGapService_GetOutgoingAirGap_NilRequest(t *testing.T) {
	svc := &AirGapService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetOutgoingAirGap(nil, nil)
	if err == nil {
		t.Error("GetOutgoingAirGap() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("GetOutgoingAirGap() error = %v, want 'request cannot be nil'", err)
	}
}

func TestAirGapService_GetOutgoingAirGap_EmptyRequestIDs(t *testing.T) {
	svc := &AirGapService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetOutgoingAirGap(nil, &model.GetOutgoingAirGapRequest{})
	if err == nil {
		t.Error("GetOutgoingAirGap() with empty request IDs should return error")
	}
	if err.Error() != "either request IDs or address IDs must be provided" {
		t.Errorf("GetOutgoingAirGap() error = %v, want 'either request IDs or address IDs must be provided'", err)
	}
}

func TestAirGapService_GetOutgoingAirGap_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		request *model.GetOutgoingAirGapRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name:    "empty request",
			request: &model.GetOutgoingAirGapRequest{},
			wantErr: true,
			errMsg:  "either request IDs or address IDs must be provided",
		},
		{
			name: "empty request IDs slice",
			request: &model.GetOutgoingAirGapRequest{
				RequestIDs: []string{},
			},
			wantErr: true,
			errMsg:  "either request IDs or address IDs must be provided",
		},
		{
			name: "empty address IDs slice",
			request: &model.GetOutgoingAirGapRequest{
				AddressIDs: []string{},
			},
			wantErr: true,
			errMsg:  "either request IDs or address IDs must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &AirGapService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}

			_, err := svc.GetOutgoingAirGap(nil, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOutgoingAirGap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("GetOutgoingAirGap() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAirGapService_SubmitIncomingAirGap_NilRequest(t *testing.T) {
	svc := &AirGapService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.SubmitIncomingAirGap(nil, nil)
	if err == nil {
		t.Error("SubmitIncomingAirGap() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("SubmitIncomingAirGap() error = %v, want 'request cannot be nil'", err)
	}
}

func TestAirGapService_SubmitIncomingAirGap_EmptyPayload(t *testing.T) {
	svc := &AirGapService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.SubmitIncomingAirGap(nil, &model.SubmitIncomingAirGapRequest{})
	if err == nil {
		t.Error("SubmitIncomingAirGap() with empty payload should return error")
	}
	if err.Error() != "payload is required" {
		t.Errorf("SubmitIncomingAirGap() error = %v, want 'payload is required'", err)
	}
}

func TestAirGapService_SubmitIncomingAirGap_EmptySignature(t *testing.T) {
	svc := &AirGapService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.SubmitIncomingAirGap(nil, &model.SubmitIncomingAirGapRequest{
		Payload: "base64payload",
	})
	if err == nil {
		t.Error("SubmitIncomingAirGap() with empty signature should return error")
	}
	if err.Error() != "signature is required" {
		t.Errorf("SubmitIncomingAirGap() error = %v, want 'signature is required'", err)
	}
}

func TestAirGapService_SubmitIncomingAirGap_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		request *model.SubmitIncomingAirGapRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name:    "empty request",
			request: &model.SubmitIncomingAirGapRequest{},
			wantErr: true,
			errMsg:  "payload is required",
		},
		{
			name: "missing signature",
			request: &model.SubmitIncomingAirGapRequest{
				Payload: "base64payload",
			},
			wantErr: true,
			errMsg:  "signature is required",
		},
		{
			name: "missing payload",
			request: &model.SubmitIncomingAirGapRequest{
				Signature: "base64signature",
			},
			wantErr: true,
			errMsg:  "payload is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &AirGapService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}

			err := svc.SubmitIncomingAirGap(nil, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmitIncomingAirGap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("SubmitIncomingAirGap() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGetOutgoingAirGapRequest_Fields(t *testing.T) {
	// Test that the request struct fields can be set correctly
	req := &model.GetOutgoingAirGapRequest{
		RequestIDs:       []string{"req-1", "req-2", "req-3"},
		RequestSignature: "base64signature",
		AddressIDs:       []string{"addr-1", "addr-2"},
	}

	if len(req.RequestIDs) != 3 {
		t.Errorf("RequestIDs length = %d, want 3", len(req.RequestIDs))
	}
	if req.RequestSignature != "base64signature" {
		t.Errorf("RequestSignature = %v, want 'base64signature'", req.RequestSignature)
	}
	if len(req.AddressIDs) != 2 {
		t.Errorf("AddressIDs length = %d, want 2", len(req.AddressIDs))
	}
}

func TestSubmitIncomingAirGapRequest_Fields(t *testing.T) {
	// Test that the request struct fields can be set correctly
	req := &model.SubmitIncomingAirGapRequest{
		Payload:   "base64payload",
		Signature: "base64signature",
	}

	if req.Payload != "base64payload" {
		t.Errorf("Payload = %v, want 'base64payload'", req.Payload)
	}
	if req.Signature != "base64signature" {
		t.Errorf("Signature = %v, want 'base64signature'", req.Signature)
	}
}

func TestGetOutgoingAirGapResult_Fields(t *testing.T) {
	// Test that the result struct fields can be set correctly
	result := &model.GetOutgoingAirGapResult{
		Data: []byte{0x01, 0x02, 0x03, 0x04},
	}

	if len(result.Data) != 4 {
		t.Errorf("Data length = %d, want 4", len(result.Data))
	}
	if result.Data[0] != 0x01 || result.Data[3] != 0x04 {
		t.Errorf("Data = %v, want [0x01, 0x02, 0x03, 0x04]", result.Data)
	}
}
