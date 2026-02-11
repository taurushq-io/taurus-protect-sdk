package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWhitelistedAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordSignedWhitelistedAddressEnvelope
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns address with zero values",
			dto:  &openapi.TgvalidatordSignedWhitelistedAddressEnvelope{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordSignedWhitelistedAddressEnvelope {
				id := "wla-123"
				tenantId := "tenant-456"
				blockchain := "ETH"
				network := "mainnet"
				status := "APPROVED"
				action := "CREATE"
				rule := "rule-1"
				rulesContainer := "{}"
				rulesContainerHash := "hash-abc"
				rulesSignatures := "sig-xyz"
				visibilityGroupID := "vg-001"
				tnParticipantID := "tn-001"
				return &openapi.TgvalidatordSignedWhitelistedAddressEnvelope{
					Id:                 &id,
					TenantId:           &tenantId,
					Blockchain:         &blockchain,
					Network:            &network,
					Status:             &status,
					Action:             &action,
					Rule:               &rule,
					RulesContainer:     &rulesContainer,
					RulesContainerHash: &rulesContainerHash,
					RulesSignatures:    &rulesSignatures,
					VisibilityGroupID:  &visibilityGroupID,
					TnParticipantID:    &tnParticipantID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WhitelistedAddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WhitelistedAddressFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestWhitelistedAddressFromDTO_WithMetadata(t *testing.T) {
	hash := "metadata-hash"
	payloadAsString := `{"address":"0x123","label":"Test Address"}`
	dto := &openapi.TgvalidatordSignedWhitelistedAddressEnvelope{
		Metadata: &openapi.TgvalidatordMetadata{
			Hash:            &hash,
			PayloadAsString: &payloadAsString,
			Payload: map[string]interface{}{
				"address": "0x123",
				"label":   "Test Address",
			},
		},
	}

	got := WhitelistedAddressFromDTO(dto)
	if got == nil {
		t.Fatal("WhitelistedAddressFromDTO() returned nil")
	}
	if got.Metadata == nil {
		t.Fatal("Metadata is nil")
	}
	if got.Metadata.Hash != hash {
		t.Errorf("Metadata.Hash = %v, want %v", got.Metadata.Hash, hash)
	}
	if got.Address != "0x123" {
		t.Errorf("Address = %v, want %v", got.Address, "0x123")
	}
	if got.Label != "Test Address" {
		t.Errorf("Label = %v, want %v", got.Label, "Test Address")
	}
}

func TestWhitelistedAddressesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordSignedWhitelistedAddressEnvelope
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordSignedWhitelistedAddressEnvelope{},
			want: 0,
		},
		{
			name: "converts multiple addresses",
			dtos: func() []openapi.TgvalidatordSignedWhitelistedAddressEnvelope {
				id1 := "wla-1"
				id2 := "wla-2"
				return []openapi.TgvalidatordSignedWhitelistedAddressEnvelope{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAddressesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WhitelistedAddressesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WhitelistedAddressesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestWhitelistedAssetMetadataFromDTO_WhitelistedAddress(t *testing.T) {
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
				hash := "hash-123"
				payloadAsString := `{"key":"value"}`
				return &openapi.TgvalidatordMetadata{
					Hash:            &hash,
					PayloadAsString: &payloadAsString,
					Payload:         map[string]interface{}{"key": "value"},
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
		})
	}
}

func TestAddressScoreFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordScore
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordScore {
				id := "score-123"
				provider := "chainalysis"
				scoreType := "risk"
				score := "75"
				now := time.Now()
				return &openapi.TgvalidatordScore{
					Id:         &id,
					Provider:   &provider,
					Type:       &scoreType,
					Score:      &score,
					UpdateDate: &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddressScoreFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Provider != "" {
					t.Errorf("AddressScoreFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Provider != nil && got.Provider != *tt.dto.Provider {
				t.Errorf("Provider = %v, want %v", got.Provider, *tt.dto.Provider)
			}
			if tt.dto.Score != nil && got.Score != *tt.dto.Score {
				t.Errorf("Score = %v, want %v", got.Score, *tt.dto.Score)
			}
		})
	}
}

func TestAddressScoresFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordScore
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordScore{},
			want: 0,
		},
		{
			name: "converts multiple scores",
			dtos: func() []openapi.TgvalidatordScore {
				id1 := "score-1"
				id2 := "score-2"
				return []openapi.TgvalidatordScore{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddressScoresFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("AddressScoresFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("AddressScoresFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTrail
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTrail {
				id := "trail-123"
				userId := "user-456"
				externalUserId := "ext-789"
				action := "APPROVE"
				comment := "Approved by admin"
				now := time.Now()
				return &openapi.TgvalidatordTrail{
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
			got := TrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.UserID != "" {
					t.Errorf("TrailFromDTO(nil) should return zero value")
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

func TestTrailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTrail
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTrail{},
			want: 0,
		},
		{
			name: "converts multiple trails",
			dtos: func() []openapi.TgvalidatordTrail {
				id1 := "trail-1"
				id2 := "trail-2"
				return []openapi.TgvalidatordTrail{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrailsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("TrailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("TrailsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestApproversFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordApprovers
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty approvers",
			dto:  &openapi.TgvalidatordApprovers{},
		},
		{
			name: "complete DTO with parallel groups",
			dto: &openapi.TgvalidatordApprovers{
				Parallel: []openapi.TgvalidatordParallelApproversGroups{
					{
						Sequential: []openapi.TgvalidatordApproversGroup{
							{
								ExternalGroupID:   strPtr("group-1"),
								MinimumSignatures: int64Ptr(2),
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproversFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ApproversFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ApproversFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Parallel != nil && len(got.Parallel) != len(tt.dto.Parallel) {
				t.Errorf("Parallel length = %v, want %v", len(got.Parallel), len(tt.dto.Parallel))
			}
		})
	}
}

func TestWhitelistedAddressAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWhitelistedAddressAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWhitelistedAddressAttribute {
				id := "attr-123"
				key := "environment"
				value := "production"
				contentType := "text/plain"
				owner := "user-456"
				attrType := "custom"
				subtype := "env"
				isFile := false
				return &openapi.TgvalidatordWhitelistedAddressAttribute{
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
			got := WhitelistedAddressAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("WhitelistedAddressAttributeFromDTO(nil) should return zero value")
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
		})
	}
}

func TestWhitelistedAddressAttributesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordWhitelistedAddressAttribute
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordWhitelistedAddressAttribute{},
			want: 0,
		},
		{
			name: "converts multiple attributes",
			dtos: func() []openapi.TgvalidatordWhitelistedAddressAttribute {
				id1 := "attr-1"
				id2 := "attr-2"
				return []openapi.TgvalidatordWhitelistedAddressAttribute{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistedAddressAttributesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WhitelistedAddressAttributesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WhitelistedAddressAttributesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSignedWhitelistedAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordSignedWhitelistedAddress
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordSignedWhitelistedAddress {
				payload := "base64payload=="
				userId := "user-123"
				signature := "base64sig=="
				comment := "Signed by user"
				return &openapi.TgvalidatordSignedWhitelistedAddress{
					Payload: &payload,
					Signatures: []openapi.TgvalidatordWhitelistSignature{
						{
							Signature: &openapi.TgvalidatordWhitelistUserSignature{
								UserId:    &userId,
								Signature: &signature,
								Comment:   &comment,
							},
							Hashes: []string{"hash1", "hash2"},
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SignedWhitelistedAddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SignedWhitelistedAddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SignedWhitelistedAddressFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Payload != nil && got.Payload != *tt.dto.Payload {
				t.Errorf("Payload = %v, want %v", got.Payload, *tt.dto.Payload)
			}
			if tt.dto.Signatures != nil && len(got.Signatures) != len(tt.dto.Signatures) {
				t.Errorf("Signatures length = %v, want %v", len(got.Signatures), len(tt.dto.Signatures))
			}
		})
	}
}

func TestWhitelistSignatureFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWhitelistSignature
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWhitelistSignature {
				userId := "user-123"
				signature := "base64sig=="
				comment := "Signed by user"
				return &openapi.TgvalidatordWhitelistSignature{
					Signature: &openapi.TgvalidatordWhitelistUserSignature{
						UserId:    &userId,
						Signature: &signature,
						Comment:   &comment,
					},
					Hashes: []string{"hash1", "hash2"},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WhitelistSignatureFromDTO(tt.dto)
			if tt.dto == nil {
				if got.UserSignature != nil || got.Hashes != nil {
					t.Errorf("WhitelistSignatureFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Signature != nil && tt.dto.Signature.UserId != nil {
				if got.UserSignature == nil {
					t.Error("UserSignature should not be nil")
				} else if got.UserSignature.UserID != *tt.dto.Signature.UserId {
					t.Errorf("UserSignature.UserID = %v, want %v", got.UserSignature.UserID, *tt.dto.Signature.UserId)
				}
			}
			if tt.dto.Hashes != nil && len(got.Hashes) != len(tt.dto.Hashes) {
				t.Errorf("Hashes length = %v, want %v", len(got.Hashes), len(tt.dto.Hashes))
			}
		})
	}
}

func TestApproversGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordApproversGroup
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordApproversGroup {
				externalGroupID := "group-123"
				var minSigs int64 = 3
				return &openapi.TgvalidatordApproversGroup{
					ExternalGroupID:   &externalGroupID,
					MinimumSignatures: &minSigs,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproversGroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ExternalGroupID != "" || got.MinimumSignatures != 0 {
					t.Errorf("ApproversGroupFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.ExternalGroupID != nil && got.ExternalGroupID != *tt.dto.ExternalGroupID {
				t.Errorf("ExternalGroupID = %v, want %v", got.ExternalGroupID, *tt.dto.ExternalGroupID)
			}
			if tt.dto.MinimumSignatures != nil && got.MinimumSignatures != *tt.dto.MinimumSignatures {
				t.Errorf("MinimumSignatures = %v, want %v", got.MinimumSignatures, *tt.dto.MinimumSignatures)
			}
		})
	}
}

func TestParallelApproversGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordParallelApproversGroups
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps sequential groups",
			dto: func() *openapi.TgvalidatordParallelApproversGroups {
				externalGroupID := "group-123"
				var minSigs int64 = 2
				return &openapi.TgvalidatordParallelApproversGroups{
					Sequential: []openapi.TgvalidatordApproversGroup{
						{
							ExternalGroupID:   &externalGroupID,
							MinimumSignatures: &minSigs,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParallelApproversGroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got.Sequential != nil {
					t.Errorf("ParallelApproversGroupFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Sequential != nil && len(got.Sequential) != len(tt.dto.Sequential) {
				t.Errorf("Sequential length = %v, want %v", len(got.Sequential), len(tt.dto.Sequential))
			}
		})
	}
}
