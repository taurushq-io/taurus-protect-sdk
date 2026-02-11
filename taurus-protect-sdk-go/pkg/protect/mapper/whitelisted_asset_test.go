package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWhitelistedAssetFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope
		want func(t *testing.T, got interface{})
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
			want: func(t *testing.T, got interface{}) {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			name: "empty DTO returns asset with zero values",
			dto:  &openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{},
			want: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Error("expected non-nil asset")
				}
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope {
				id := "asset-123"
				tenantId := "tenant-456"
				status := "approved"
				action := "create"
				blockchain := "ETH"
				network := "mainnet"
				rule := "rule-1"
				rulesContainer := "container-1"
				rulesSignatures := "sig-1"
				businessRuleEnabled := true
				return &openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{
					Id:                  &id,
					TenantId:            &tenantId,
					Status:              &status,
					Action:              &action,
					Blockchain:          &blockchain,
					Network:             &network,
					Rule:                &rule,
					RulesContainer:      &rulesContainer,
					RulesSignatures:     &rulesSignatures,
					BusinessRuleEnabled: &businessRuleEnabled,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Type assertion handled in test execution
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAssetFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WhitelistedAssetFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WhitelistedAssetFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.BusinessRuleEnabled != nil && got.BusinessRuleEnabled != *tt.dto.BusinessRuleEnabled {
				t.Errorf("BusinessRuleEnabled = %v, want %v", got.BusinessRuleEnabled, *tt.dto.BusinessRuleEnabled)
			}
		})
	}
}

func TestWhitelistedAssetsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{},
			want: 0,
		},
		{
			name: "converts multiple assets",
			dtos: func() []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope {
				id1 := "asset-1"
				id2 := "asset-2"
				return []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAssetsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WhitelistedAssetsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WhitelistedAssetsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestWhitelistedAssetMetadataFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordMetadata
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordMetadata {
				hash := "abc123"
				payloadAsString := `{"key":"value"}`
				return &openapi.TgvalidatordMetadata{
					Hash:            &hash,
					Payload:         map[string]interface{}{"key": "value"},
					PayloadAsString: &payloadAsString,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAssetMetadataFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WhitelistedAssetMetadataFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WhitelistedAssetMetadataFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Hash != nil && got.Hash != *tt.dto.Hash {
				t.Errorf("Hash = %v, want %v", got.Hash, *tt.dto.Hash)
			}
			if tt.dto.PayloadAsString != nil && got.PayloadAsString != *tt.dto.PayloadAsString {
				t.Errorf("PayloadAsString = %v, want %v", got.PayloadAsString, *tt.dto.PayloadAsString)
			}
		})
	}
}

func TestSignedContractAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordSignedWhitelistedContractAddress
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordSignedWhitelistedContractAddress {
				payload := "base64payload=="
				return &openapi.TgvalidatordSignedWhitelistedContractAddress{
					Payload:    &payload,
					Signatures: []openapi.TgvalidatordWhitelistSignature{},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SignedContractAddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SignedContractAddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SignedContractAddressFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Payload != nil && got.Payload != *tt.dto.Payload {
				t.Errorf("Payload = %v, want %v", got.Payload, *tt.dto.Payload)
			}
		})
	}
}

func TestWhitelistUserSignatureFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWhitelistUserSignature
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWhitelistUserSignature {
				userId := "user-123"
				sig := "signature=="
				comment := "approved by user"
				return &openapi.TgvalidatordWhitelistUserSignature{
					UserId:    &userId,
					Signature: &sig,
					Comment:   &comment,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistUserSignatureFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WhitelistUserSignatureFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WhitelistUserSignatureFromDTO() returned nil for non-nil input")
			}
			if tt.dto.UserId != nil && got.UserID != *tt.dto.UserId {
				t.Errorf("UserID = %v, want %v", got.UserID, *tt.dto.UserId)
			}
			if tt.dto.Signature != nil && got.Signature != *tt.dto.Signature {
				t.Errorf("Signature = %v, want %v", got.Signature, *tt.dto.Signature)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
		})
	}
}

func TestWhitelistedAssetAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWhitelistedContractAddressAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWhitelistedContractAddressAttribute {
				id := "attr-123"
				key := "environment"
				value := "production"
				contentType := "text/plain"
				owner := "admin"
				attrType := "string"
				subtype := "config"
				isFile := false
				return &openapi.TgvalidatordWhitelistedContractAddressAttribute{
					Id:          &id,
					Key:         &key,
					Value:       &value,
					ContentType: &contentType,
					Owner:       &owner,
					Type:        &attrType,
					Subtype:     &subtype,
					Isfile:      &isFile,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAssetAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("WhitelistedAssetAttributeFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Key != nil && got.Key != *tt.dto.Key {
				t.Errorf("Key = %v, want %v", got.Key, *tt.dto.Key)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.ContentType != nil && got.ContentType != *tt.dto.ContentType {
				t.Errorf("ContentType = %v, want %v", got.ContentType, *tt.dto.ContentType)
			}
			if tt.dto.Owner != nil && got.Owner != *tt.dto.Owner {
				t.Errorf("Owner = %v, want %v", got.Owner, *tt.dto.Owner)
			}
			if tt.dto.Isfile != nil && got.IsFile != *tt.dto.Isfile {
				t.Errorf("IsFile = %v, want %v", got.IsFile, *tt.dto.Isfile)
			}
		})
	}
}

func TestSafeTime(t *testing.T) {
	tests := []struct {
		name  string
		input *time.Time
		want  bool // true if should be zero
	}{
		{"nil returns zero time", nil, true},
		{"value returns value", func() *time.Time { t := time.Now(); return &t }(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeTime(tt.input)
			if tt.want && !got.IsZero() {
				t.Errorf("safeTime() = %v, want zero time", got)
			}
			if !tt.want && got.IsZero() {
				t.Errorf("safeTime() should not be zero")
			}
		})
	}
}
