package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestGovernanceRulesetFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRules
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns rules with zero values",
			dto:  &openapi.TgvalidatordRules{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRules {
				rulesContainer := "base64encodedcontainer"
				locked := true
				now := time.Now()
				userId := "user-123"
				signature := "sig-abc"
				return &openapi.TgvalidatordRules{
					RulesContainer: &rulesContainer,
					Locked:         &locked,
					CreationDate:   &now,
					UpdateDate:     &now,
					RulesSignatures: []openapi.TgvalidatordRuleUserSignature{
						{UserId: &userId, Signature: &signature},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GovernanceRulesetFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("GovernanceRulesetFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("GovernanceRulesetFromDTO() returned nil for non-nil input")
			}
			if tt.dto.RulesContainer != nil && got.RulesContainer != *tt.dto.RulesContainer {
				t.Errorf("RulesContainer = %v, want %v", got.RulesContainer, *tt.dto.RulesContainer)
			}
			if tt.dto.Locked != nil && got.Locked != *tt.dto.Locked {
				t.Errorf("Locked = %v, want %v", got.Locked, *tt.dto.Locked)
			}
			if tt.dto.RulesSignatures != nil && len(got.Signatures) != len(tt.dto.RulesSignatures) {
				t.Errorf("Signatures length = %v, want %v", len(got.Signatures), len(tt.dto.RulesSignatures))
			}
		})
	}
}

func TestGovernanceRulesetsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordRules
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordRules{},
			want: 0,
		},
		{
			name: "converts multiple rules",
			dtos: func() []openapi.TgvalidatordRules {
				container1 := "container1"
				container2 := "container2"
				return []openapi.TgvalidatordRules{
					{RulesContainer: &container1},
					{RulesContainer: &container2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GovernanceRulesetsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("GovernanceRulesetsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("GovernanceRulesetsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestRuleUserSignatureFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRuleUserSignature
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRuleUserSignature {
				userId := "user-123"
				signature := "sig-abc"
				return &openapi.TgvalidatordRuleUserSignature{
					UserId:    &userId,
					Signature: &signature,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RuleUserSignatureFromDTO(tt.dto)
			if tt.dto == nil {
				if got.UserID != "" || got.Signature != "" {
					t.Errorf("RuleUserSignatureFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.UserId != nil && got.UserID != *tt.dto.UserId {
				t.Errorf("UserID = %v, want %v", got.UserID, *tt.dto.UserId)
			}
			if tt.dto.Signature != nil && got.Signature != *tt.dto.Signature {
				t.Errorf("Signature = %v, want %v", got.Signature, *tt.dto.Signature)
			}
		})
	}
}

func TestRulesTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRulesTrail
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRulesTrail {
				id := "trail-123"
				userId := "user-456"
				externalUserId := "ext-789"
				action := "APPROVE"
				comment := "Approved by admin"
				now := time.Now()
				return &openapi.TgvalidatordRulesTrail{
					Id:             &id,
					UserId:         &userId,
					ExternalUserId: &externalUserId,
					Action:         &action,
					Comment:        &comment,
					Date:           &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RulesTrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.UserID != "" || got.Action != "" {
					t.Errorf("RulesTrailFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.UserId != nil && got.UserID != *tt.dto.UserId {
				t.Errorf("UserID = %v, want %v", got.UserID, *tt.dto.UserId)
			}
			if tt.dto.Action != nil && got.Action != *tt.dto.Action {
				t.Errorf("Action = %v, want %v", got.Action, *tt.dto.Action)
			}
		})
	}
}

func TestSuperAdminPublicKeyFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.GetPublicKeysReplyPublicKey
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.GetPublicKeysReplyPublicKey {
				userId := "user-123"
				publicKey := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n-----END PUBLIC KEY-----"
				return &openapi.GetPublicKeysReplyPublicKey{
					UserID:    &userId,
					PublicKey: &publicKey,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuperAdminPublicKeyFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SuperAdminPublicKeyFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SuperAdminPublicKeyFromDTO() returned nil for non-nil input")
			}
			if tt.dto.UserID != nil && got.UserID != *tt.dto.UserID {
				t.Errorf("UserID = %v, want %v", got.UserID, *tt.dto.UserID)
			}
			if tt.dto.PublicKey != nil && got.PublicKey != *tt.dto.PublicKey {
				t.Errorf("PublicKey = %v, want %v", got.PublicKey, *tt.dto.PublicKey)
			}
		})
	}
}

func TestSuperAdminPublicKeysFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.GetPublicKeysReplyPublicKey
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.GetPublicKeysReplyPublicKey{},
			want: 0,
		},
		{
			name: "converts multiple keys",
			dtos: func() []openapi.GetPublicKeysReplyPublicKey {
				userId1 := "user-1"
				userId2 := "user-2"
				return []openapi.GetPublicKeysReplyPublicKey{
					{UserID: &userId1},
					{UserID: &userId2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuperAdminPublicKeysFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SuperAdminPublicKeysFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SuperAdminPublicKeysFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}
