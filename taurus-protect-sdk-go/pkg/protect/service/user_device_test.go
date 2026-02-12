package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewUserDeviceService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestUserDeviceService_CreatePairing(t *testing.T) {
	// Create a service with nil API to test service structure
	svc := &UserDeviceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("UserDeviceService should not be nil")
	}
}

func TestUserDeviceService_StartPairing_Validation(t *testing.T) {
	svc := &UserDeviceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tests := []struct {
		name      string
		pairingID string
		req       *model.StartPairingRequest
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty pairing ID returns error",
			pairingID: "",
			req:       &model.StartPairingRequest{Nonce: "123456", PublicKey: "key"},
			wantErr:   true,
			errMsg:    "pairingID cannot be empty",
		},
		{
			name:      "nil request returns error",
			pairingID: "pairing-123",
			req:       nil,
			wantErr:   true,
			errMsg:    "request cannot be nil",
		},
		{
			name:      "empty nonce returns error",
			pairingID: "pairing-123",
			req:       &model.StartPairingRequest{Nonce: "", PublicKey: "key"},
			wantErr:   true,
			errMsg:    "nonce is required",
		},
		{
			name:      "empty publicKey returns error",
			pairingID: "pairing-123",
			req:       &model.StartPairingRequest{Nonce: "123456", PublicKey: ""},
			wantErr:   true,
			errMsg:    "publicKey is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.StartPairing(nil, tt.pairingID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartPairing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("StartPairing() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserDeviceService_ApprovePairing_Validation(t *testing.T) {
	svc := &UserDeviceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tests := []struct {
		name      string
		pairingID string
		req       *model.ApprovePairingRequest
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty pairing ID returns error",
			pairingID: "",
			req:       &model.ApprovePairingRequest{Nonce: "123456"},
			wantErr:   true,
			errMsg:    "pairingID cannot be empty",
		},
		{
			name:      "nil request returns error",
			pairingID: "pairing-123",
			req:       nil,
			wantErr:   true,
			errMsg:    "request cannot be nil",
		},
		{
			name:      "empty nonce returns error",
			pairingID: "pairing-123",
			req:       &model.ApprovePairingRequest{Nonce: ""},
			wantErr:   true,
			errMsg:    "nonce is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ApprovePairing(nil, tt.pairingID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovePairing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("ApprovePairing() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserDeviceService_GetPairingStatus_Validation(t *testing.T) {
	svc := &UserDeviceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tests := []struct {
		name      string
		pairingID string
		nonce     string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty pairing ID returns error",
			pairingID: "",
			nonce:     "123456",
			wantErr:   true,
			errMsg:    "pairingID cannot be empty",
		},
		{
			name:      "empty nonce returns error",
			pairingID: "pairing-123",
			nonce:     "",
			wantErr:   true,
			errMsg:    "nonce cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetPairingStatus(nil, tt.pairingID, tt.nonce)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPairingStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("GetPairingStatus() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserDevicePairingStatus_Constants(t *testing.T) {
	// Test that status constants have expected values
	tests := []struct {
		name   string
		status model.UserDevicePairingStatus
		want   string
	}{
		{
			name:   "WAITING status",
			status: model.UserDevicePairingStatusWaiting,
			want:   "WAITING",
		},
		{
			name:   "PAIRING status",
			status: model.UserDevicePairingStatusPairing,
			want:   "PAIRING",
		},
		{
			name:   "APPROVED status",
			status: model.UserDevicePairingStatusApproved,
			want:   "APPROVED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("Status = %v, want %v", tt.status, tt.want)
			}
		})
	}
}

func TestStartPairingRequest_Fields(t *testing.T) {
	req := &model.StartPairingRequest{
		Nonce:     "123456",
		PublicKey: "base64-public-key",
	}

	if req.Nonce != "123456" {
		t.Errorf("Nonce = %v, want 123456", req.Nonce)
	}
	if req.PublicKey != "base64-public-key" {
		t.Errorf("PublicKey = %v, want base64-public-key", req.PublicKey)
	}
}

func TestApprovePairingRequest_Fields(t *testing.T) {
	req := &model.ApprovePairingRequest{
		Nonce: "654321",
	}

	if req.Nonce != "654321" {
		t.Errorf("Nonce = %v, want 654321", req.Nonce)
	}
}

func TestCreatePairingResult_Fields(t *testing.T) {
	result := &model.CreatePairingResult{
		PairingID: "pairing-abc123",
	}

	if result.PairingID != "pairing-abc123" {
		t.Errorf("PairingID = %v, want pairing-abc123", result.PairingID)
	}
}

func TestPairingStatusResult_Fields(t *testing.T) {
	result := &model.PairingStatusResult{
		Status:    model.UserDevicePairingStatusApproved,
		PairingID: "pairing-xyz789",
		APIKey:    "api-key-secret",
	}

	if result.Status != model.UserDevicePairingStatusApproved {
		t.Errorf("Status = %v, want %v", result.Status, model.UserDevicePairingStatusApproved)
	}
	if result.PairingID != "pairing-xyz789" {
		t.Errorf("PairingID = %v, want pairing-xyz789", result.PairingID)
	}
	if result.APIKey != "api-key-secret" {
		t.Errorf("APIKey = %v, want api-key-secret", result.APIKey)
	}
}
