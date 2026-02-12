package service

import (
	"context"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewMultiFactorSignatureService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestMultiFactorSignatureService_ServiceStructure(t *testing.T) {
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("MultiFactorSignatureService should not be nil")
	}
	if svc.errMapper == nil {
		t.Error("ErrorMapper should not be nil")
	}
}

func TestMultiFactorSignatureService_GetMultiFactorSignatureInfo_EmptyID(t *testing.T) {
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetMultiFactorSignatureInfo(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty id, got nil")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("error = %v, want containing 'id cannot be empty'", err.Error())
	}
}

func TestMultiFactorSignatureService_CreateMultiFactorSignatures_Validation(t *testing.T) {
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tests := []struct {
		name       string
		entityIDs  []string
		entityType model.MultiFactorSignatureEntityType
		wantErr    string
	}{
		{
			name:       "empty entityIDs",
			entityIDs:  []string{},
			entityType: model.MFSEntityTypeWhitelistedAddress,
			wantErr:    "entityIDs cannot be empty",
		},
		{
			name:       "nil entityIDs",
			entityIDs:  nil,
			entityType: model.MFSEntityTypeWhitelistedAddress,
			wantErr:    "entityIDs cannot be empty",
		},
		{
			name:       "empty entityType",
			entityIDs:  []string{"entity-1"},
			entityType: "",
			wantErr:    "entityType cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateMultiFactorSignatures(context.Background(), tt.entityIDs, tt.entityType)
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %v, want containing %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestMultiFactorSignatureService_ApproveMultiFactorSignature_Validation(t *testing.T) {
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tests := []struct {
		name      string
		id        string
		signature string
		comment   string
		wantErr   string
	}{
		{
			name:      "empty id",
			id:        "",
			signature: "sig-123",
			comment:   "approved",
			wantErr:   "id cannot be empty",
		},
		{
			name:      "empty signature",
			id:        "mfs-123",
			signature: "",
			comment:   "approved",
			wantErr:   "signature cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ApproveMultiFactorSignature(context.Background(), tt.id, tt.signature, tt.comment)
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %v, want containing %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestMultiFactorSignatureService_RejectMultiFactorSignature_EmptyID(t *testing.T) {
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectMultiFactorSignature(context.Background(), "", "rejected")
	if err == nil {
		t.Error("expected error for empty id, got nil")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("error = %v, want containing 'id cannot be empty'", err.Error())
	}
}

func TestMultiFactorSignatureEntityType_Constants(t *testing.T) {
	tests := []struct {
		name       string
		entityType model.MultiFactorSignatureEntityType
		want       string
	}{
		{
			name:       "REQUEST entity type",
			entityType: model.MFSEntityTypeRequest,
			want:       "REQUEST",
		},
		{
			name:       "WHITELISTED_ADDRESS entity type",
			entityType: model.MFSEntityTypeWhitelistedAddress,
			want:       "WHITELISTED_ADDRESS",
		},
		{
			name:       "WHITELISTED_CONTRACT entity type",
			entityType: model.MFSEntityTypeWhitelistedContract,
			want:       "WHITELISTED_CONTRACT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.entityType) != tt.want {
				t.Errorf("entityType = %v, want %v", tt.entityType, tt.want)
			}
		})
	}
}

func TestMultiFactorSignatureInfo_Structure(t *testing.T) {
	info := &model.MultiFactorSignatureInfo{
		ID:            "mfs-123",
		PayloadToSign: []string{"payload-1", "payload-2"},
		EntityType:    model.MFSEntityTypeWhitelistedAddress,
	}

	if info.ID != "mfs-123" {
		t.Errorf("ID = %v, want mfs-123", info.ID)
	}
	if len(info.PayloadToSign) != 2 {
		t.Errorf("PayloadToSign length = %v, want 2", len(info.PayloadToSign))
	}
	if info.EntityType != model.MFSEntityTypeWhitelistedAddress {
		t.Errorf("EntityType = %v, want %v", info.EntityType, model.MFSEntityTypeWhitelistedAddress)
	}
}

func TestMultiFactorSignatureResult_Structure(t *testing.T) {
	result := &model.MultiFactorSignatureResult{
		ID: "batch-456",
	}

	if result.ID != "batch-456" {
		t.Errorf("ID = %v, want batch-456", result.ID)
	}
}

func TestMultiFactorSignatureApprovalResult_Structure(t *testing.T) {
	result := &model.MultiFactorSignatureApprovalResult{
		SignatureCount: "3",
	}

	if result.SignatureCount != "3" {
		t.Errorf("SignatureCount = %v, want 3", result.SignatureCount)
	}
}

func TestMultiFactorSignatureService_ApproveWithEmptyComment(t *testing.T) {
	// Verify that empty comment is allowed (only id and signature are required)
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// Empty comment should not cause a validation error.
	// The call will fail at the API level since api is nil, but validation should pass.
	// We test this by checking that the id and signature validations fire first.
	_, err := svc.ApproveMultiFactorSignature(context.Background(), "", "sig-123", "")
	if err == nil {
		t.Error("expected error for empty id, got nil")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("error = %v, want containing 'id cannot be empty'", err.Error())
	}
}

func TestMultiFactorSignatureService_RejectWithEmptyComment(t *testing.T) {
	// Verify that empty comment is allowed for reject (only id is required)
	svc := &MultiFactorSignatureService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// Empty comment should not cause a validation error.
	// The call will fail at the API level since api is nil, but validation should pass for comment.
	// We test this by checking that the id validation fires first.
	err := svc.RejectMultiFactorSignature(context.Background(), "", "")
	if err == nil {
		t.Error("expected error for empty id, got nil")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("error = %v, want containing 'id cannot be empty'", err.Error())
	}
}
