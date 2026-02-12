package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestCreatePairingResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCreateUserDevicePairingReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with empty pairing ID",
			dto:  &openapi.TgvalidatordCreateUserDevicePairingReply{},
		},
		{
			name: "complete DTO maps pairing ID",
			dto: &openapi.TgvalidatordCreateUserDevicePairingReply{
				PairingID: "pairing-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreatePairingResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CreatePairingResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CreatePairingResultFromDTO() returned nil for non-nil input")
			}
			if got.PairingID != tt.dto.PairingID {
				t.Errorf("PairingID = %v, want %v", got.PairingID, tt.dto.PairingID)
			}
		})
	}
}

func TestPairingStatusResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordUserDevicePairingInfo
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "waiting status",
			dto: &openapi.TgvalidatordUserDevicePairingInfo{
				Status:    openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_WAITING,
				PairingID: "pairing-waiting",
			},
		},
		{
			name: "pairing status",
			dto: &openapi.TgvalidatordUserDevicePairingInfo{
				Status:    openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_PAIRING,
				PairingID: "pairing-in-progress",
			},
		},
		{
			name: "approved status with API key",
			dto: func() *openapi.TgvalidatordUserDevicePairingInfo {
				apiKey := "api-key-123"
				return &openapi.TgvalidatordUserDevicePairingInfo{
					Status:    openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_APPROVED,
					PairingID: "pairing-approved",
					ApiKey:    &apiKey,
				}
			}(),
		},
		{
			name: "approved status without API key",
			dto: &openapi.TgvalidatordUserDevicePairingInfo{
				Status:    openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_APPROVED,
				PairingID: "pairing-approved-no-key",
				ApiKey:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PairingStatusResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("PairingStatusResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("PairingStatusResultFromDTO() returned nil for non-nil input")
			}
			if string(got.Status) != string(tt.dto.Status) {
				t.Errorf("Status = %v, want %v", got.Status, tt.dto.Status)
			}
			if got.PairingID != tt.dto.PairingID {
				t.Errorf("PairingID = %v, want %v", got.PairingID, tt.dto.PairingID)
			}
			if tt.dto.ApiKey != nil {
				if got.APIKey != *tt.dto.ApiKey {
					t.Errorf("APIKey = %v, want %v", got.APIKey, *tt.dto.ApiKey)
				}
			} else {
				if got.APIKey != "" {
					t.Errorf("APIKey should be empty when DTO has nil ApiKey, got %v", got.APIKey)
				}
			}
		})
	}
}

func TestStartPairingRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *model.StartPairingRequest
	}{
		{
			name: "nil input returns nil",
			req:  nil,
		},
		{
			name: "empty request",
			req:  &model.StartPairingRequest{},
		},
		{
			name: "complete request maps all fields",
			req: &model.StartPairingRequest{
				Nonce:     "123456",
				PublicKey: "base64-encoded-public-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartPairingRequestToDTO(tt.req)
			if tt.req == nil {
				if got != nil {
					t.Errorf("StartPairingRequestToDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("StartPairingRequestToDTO() returned nil for non-nil input")
			}
			if got.Nonce != tt.req.Nonce {
				t.Errorf("Nonce = %v, want %v", got.Nonce, tt.req.Nonce)
			}
			if got.PublicKey != tt.req.PublicKey {
				t.Errorf("PublicKey = %v, want %v", got.PublicKey, tt.req.PublicKey)
			}
		})
	}
}

func TestApprovePairingRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *model.ApprovePairingRequest
	}{
		{
			name: "nil input returns nil",
			req:  nil,
		},
		{
			name: "empty request",
			req:  &model.ApprovePairingRequest{},
		},
		{
			name: "complete request maps nonce",
			req: &model.ApprovePairingRequest{
				Nonce: "654321",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApprovePairingRequestToDTO(tt.req)
			if tt.req == nil {
				if got != nil {
					t.Errorf("ApprovePairingRequestToDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ApprovePairingRequestToDTO() returned nil for non-nil input")
			}
			if got.Nonce != tt.req.Nonce {
				t.Errorf("Nonce = %v, want %v", got.Nonce, tt.req.Nonce)
			}
		})
	}
}

func TestPairingStatusResultFromDTO_StatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		status     openapi.TgvalidatordUserDevicePairingInfoStatus
		wantStatus model.UserDevicePairingStatus
	}{
		{
			name:       "WAITING maps correctly",
			status:     openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_WAITING,
			wantStatus: model.UserDevicePairingStatusWaiting,
		},
		{
			name:       "PAIRING maps correctly",
			status:     openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_PAIRING,
			wantStatus: model.UserDevicePairingStatusPairing,
		},
		{
			name:       "APPROVED maps correctly",
			status:     openapi.TGVALIDATORDUSERDEVICEPAIRINGINFOSTATUS_APPROVED,
			wantStatus: model.UserDevicePairingStatusApproved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordUserDevicePairingInfo{
				Status:    tt.status,
				PairingID: "test-pairing",
			}
			got := PairingStatusResultFromDTO(dto)
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", got.Status, tt.wantStatus)
			}
		})
	}
}
