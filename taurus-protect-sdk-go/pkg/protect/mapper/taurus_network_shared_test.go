package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestSharedAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSharedAddress
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns shared address with zero values",
			dto:  &openapi.TgvalidatordTnSharedAddress{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSharedAddress {
				id := "shared-addr-123"
				internalAddressID := "internal-addr-456"
				wlAddressID := "wl-addr-789"
				ownerParticipantID := "owner-123"
				targetParticipantID := "target-456"
				blockchain := "ethereum"
				network := "mainnet"
				address := "0x1234567890abcdef"
				originLabel := "Test Address"
				status := "accepted"
				pledgesCount := "5"
				now := time.Now()
				return &openapi.TgvalidatordTnSharedAddress{
					Id:                  &id,
					InternalAddressID:   &internalAddressID,
					WladdressID:         &wlAddressID,
					OwnerParticipantId:  &ownerParticipantID,
					TargetParticipantId: &targetParticipantID,
					Blockchain:          &blockchain,
					Network:             &network,
					Address:             &address,
					OriginLabel:         &originLabel,
					Status:              &status,
					PledgesCount:        &pledgesCount,
					OriginCreationDate:  &now,
					OriginDeletionDate:  &now,
					CreatedAt:           &now,
					UpdatedAt:           &now,
					TargetAcceptedAt:    &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SharedAddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SharedAddressFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.InternalAddressID != nil && got.InternalAddressID != *tt.dto.InternalAddressID {
				t.Errorf("InternalAddressID = %v, want %v", got.InternalAddressID, *tt.dto.InternalAddressID)
			}
			if tt.dto.WladdressID != nil && got.WhitelistedAddressID != *tt.dto.WladdressID {
				t.Errorf("WhitelistedAddressID = %v, want %v", got.WhitelistedAddressID, *tt.dto.WladdressID)
			}
			if tt.dto.OwnerParticipantId != nil && got.OwnerParticipantID != *tt.dto.OwnerParticipantId {
				t.Errorf("OwnerParticipantID = %v, want %v", got.OwnerParticipantID, *tt.dto.OwnerParticipantId)
			}
			if tt.dto.TargetParticipantId != nil && got.TargetParticipantID != *tt.dto.TargetParticipantId {
				t.Errorf("TargetParticipantID = %v, want %v", got.TargetParticipantID, *tt.dto.TargetParticipantId)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Address != nil && got.Address != *tt.dto.Address {
				t.Errorf("Address = %v, want %v", got.Address, *tt.dto.Address)
			}
			if tt.dto.Status != nil && string(got.Status) != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.PledgesCount != nil && got.PledgesCount != 5 {
				t.Errorf("PledgesCount = %v, want 5", got.PledgesCount)
			}
		})
	}
}

func TestSharedAddressesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnSharedAddress
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnSharedAddress{},
			want: 0,
		},
		{
			name: "converts multiple shared addresses",
			dtos: func() []openapi.TgvalidatordTnSharedAddress {
				addr1 := "0x1234"
				addr2 := "0x5678"
				return []openapi.TgvalidatordTnSharedAddress{
					{Address: &addr1},
					{Address: &addr2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAddressesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SharedAddressesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SharedAddressesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSharedAddressTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSharedAddressTrail
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns trail with zero values",
			dto:  &openapi.TgvalidatordTnSharedAddressTrail{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSharedAddressTrail {
				id := "trail-123"
				sharedAddressID := "shared-addr-456"
				addressStatus := "accepted"
				comment := "Approved by admin"
				now := time.Now()
				return &openapi.TgvalidatordTnSharedAddressTrail{
					Id:              &id,
					SharedAddressID: &sharedAddressID,
					AddressStatus:   &addressStatus,
					Comment:         &comment,
					CreatedAt:       &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAddressTrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SharedAddressTrailFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SharedAddressTrailFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.SharedAddressID != nil && got.SharedAddressID != *tt.dto.SharedAddressID {
				t.Errorf("SharedAddressID = %v, want %v", got.SharedAddressID, *tt.dto.SharedAddressID)
			}
			if tt.dto.AddressStatus != nil && got.AddressStatus != *tt.dto.AddressStatus {
				t.Errorf("AddressStatus = %v, want %v", got.AddressStatus, *tt.dto.AddressStatus)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
		})
	}
}

func TestSharedAddressTrailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnSharedAddressTrail
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnSharedAddressTrail{},
			want: 0,
		},
		{
			name: "converts multiple trails",
			dtos: func() []openapi.TgvalidatordTnSharedAddressTrail {
				status1 := "pending"
				status2 := "accepted"
				return []openapi.TgvalidatordTnSharedAddressTrail{
					{AddressStatus: &status1},
					{AddressStatus: &status2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAddressTrailsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SharedAddressTrailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SharedAddressTrailsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestProofOfOwnershipFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordProofOfOwnership
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns proof with zero values",
			dto:  &openapi.TgvalidatordProofOfOwnership{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordProofOfOwnership {
				hash := "abc123hash"
				payload := "signed payload string"
				return &openapi.TgvalidatordProofOfOwnership{
					SignedPayloadHash:     &hash,
					SignedPayloadAsString: &payload,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProofOfOwnershipFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ProofOfOwnershipFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ProofOfOwnershipFromDTO() returned nil for non-nil input")
			}
			if tt.dto.SignedPayloadHash != nil && got.SignedPayloadHash != *tt.dto.SignedPayloadHash {
				t.Errorf("SignedPayloadHash = %v, want %v", got.SignedPayloadHash, *tt.dto.SignedPayloadHash)
			}
			if tt.dto.SignedPayloadAsString != nil && got.SignedPayloadAsString != *tt.dto.SignedPayloadAsString {
				t.Errorf("SignedPayloadAsString = %v, want %v", got.SignedPayloadAsString, *tt.dto.SignedPayloadAsString)
			}
		})
	}
}

func TestSharedAssetFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSharedAsset
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns shared asset with zero values",
			dto:  &openapi.TgvalidatordTnSharedAsset{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSharedAsset {
				id := "shared-asset-123"
				wlContractAddressID := "wl-contract-456"
				ownerParticipantID := "owner-123"
				targetParticipantID := "target-456"
				blockchain := "ethereum"
				network := "mainnet"
				name := "Test Token"
				symbol := "TST"
				decimals := "18"
				contractAddress := "0xabcdef1234567890"
				tokenID := "1"
				kind := "ERC20"
				status := "accepted"
				now := time.Now()
				return &openapi.TgvalidatordTnSharedAsset{
					Id:                  &id,
					WlContractAddressID: &wlContractAddressID,
					OwnerParticipantId:  &ownerParticipantID,
					TargetParticipantId: &targetParticipantID,
					Blockchain:          &blockchain,
					Network:             &network,
					Name:                &name,
					Symbol:              &symbol,
					Decimals:            &decimals,
					ContractAddress:     &contractAddress,
					TokenId:             &tokenID,
					Kind:                &kind,
					Status:              &status,
					OriginCreationDate:  &now,
					OriginDeletionDate:  &now,
					CreatedAt:           &now,
					UpdatedAt:           &now,
					TargetAcceptedAt:    &now,
					TargetRejectedAt:    &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAssetFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SharedAssetFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SharedAssetFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.WlContractAddressID != nil && got.WhitelistedContractAddressID != *tt.dto.WlContractAddressID {
				t.Errorf("WhitelistedContractAddressID = %v, want %v", got.WhitelistedContractAddressID, *tt.dto.WlContractAddressID)
			}
			if tt.dto.OwnerParticipantId != nil && got.OwnerParticipantID != *tt.dto.OwnerParticipantId {
				t.Errorf("OwnerParticipantID = %v, want %v", got.OwnerParticipantID, *tt.dto.OwnerParticipantId)
			}
			if tt.dto.TargetParticipantId != nil && got.TargetParticipantID != *tt.dto.TargetParticipantId {
				t.Errorf("TargetParticipantID = %v, want %v", got.TargetParticipantID, *tt.dto.TargetParticipantId)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Symbol != nil && got.Symbol != *tt.dto.Symbol {
				t.Errorf("Symbol = %v, want %v", got.Symbol, *tt.dto.Symbol)
			}
			if tt.dto.Decimals != nil && got.Decimals != *tt.dto.Decimals {
				t.Errorf("Decimals = %v, want %v", got.Decimals, *tt.dto.Decimals)
			}
			if tt.dto.ContractAddress != nil && got.ContractAddress != *tt.dto.ContractAddress {
				t.Errorf("ContractAddress = %v, want %v", got.ContractAddress, *tt.dto.ContractAddress)
			}
			if tt.dto.TokenId != nil && got.TokenID != *tt.dto.TokenId {
				t.Errorf("TokenID = %v, want %v", got.TokenID, *tt.dto.TokenId)
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.Status != nil && string(got.Status) != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestSharedAssetsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnSharedAsset
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnSharedAsset{},
			want: 0,
		},
		{
			name: "converts multiple shared assets",
			dtos: func() []openapi.TgvalidatordTnSharedAsset {
				name1 := "Token1"
				name2 := "Token2"
				return []openapi.TgvalidatordTnSharedAsset{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAssetsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SharedAssetsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SharedAssetsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSharedAssetTrailFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSharedAssetTrail
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns trail with zero values",
			dto:  &openapi.TgvalidatordTnSharedAssetTrail{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSharedAssetTrail {
				id := "trail-123"
				sharedAssetID := "shared-asset-456"
				assetStatus := "accepted"
				comment := "Approved by admin"
				now := time.Now()
				return &openapi.TgvalidatordTnSharedAssetTrail{
					Id:            &id,
					SharedAssetID: &sharedAssetID,
					AssetStatus:   &assetStatus,
					Comment:       &comment,
					CreatedAt:     &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAssetTrailFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SharedAssetTrailFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SharedAssetTrailFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.SharedAssetID != nil && got.SharedAssetID != *tt.dto.SharedAssetID {
				t.Errorf("SharedAssetID = %v, want %v", got.SharedAssetID, *tt.dto.SharedAssetID)
			}
			if tt.dto.AssetStatus != nil && got.AssetStatus != *tt.dto.AssetStatus {
				t.Errorf("AssetStatus = %v, want %v", got.AssetStatus, *tt.dto.AssetStatus)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
		})
	}
}

func TestSharedAssetTrailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnSharedAssetTrail
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnSharedAssetTrail{},
			want: 0,
		},
		{
			name: "converts multiple trails",
			dtos: func() []openapi.TgvalidatordTnSharedAssetTrail {
				status1 := "pending"
				status2 := "accepted"
				return []openapi.TgvalidatordTnSharedAssetTrail{
					{AssetStatus: &status1},
					{AssetStatus: &status2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SharedAssetTrailsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SharedAssetTrailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SharedAssetTrailsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestKeyValueAttributesToDTO(t *testing.T) {
	tests := []struct {
		name  string
		attrs []taurusnetwork.KeyValueAttribute
		want  int
	}{
		{
			name:  "nil slice returns nil",
			attrs: nil,
			want:  -1,
		},
		{
			name:  "empty slice returns empty slice",
			attrs: []taurusnetwork.KeyValueAttribute{},
			want:  0,
		},
		{
			name: "converts multiple attributes",
			attrs: []taurusnetwork.KeyValueAttribute{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeyValueAttributesToDTO(tt.attrs)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("KeyValueAttributesToDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("KeyValueAttributesToDTO() length = %v, want %v", len(got), tt.want)
			}
			// Verify values for non-empty case
			if tt.want > 0 {
				for i, attr := range tt.attrs {
					if got[i].Key == nil || *got[i].Key != attr.Key {
						t.Errorf("Key[%d] = %v, want %v", i, got[i].Key, attr.Key)
					}
					if got[i].Value == nil || *got[i].Value != attr.Value {
						t.Errorf("Value[%d] = %v, want %v", i, got[i].Value, attr.Value)
					}
				}
			}
		})
	}
}

func TestSharedAddressFromDTO_WithProofOfOwnership(t *testing.T) {
	hash := "signedPayloadHash123"
	payload := "signedPayloadAsString456"
	dto := &openapi.TgvalidatordTnSharedAddress{
		ProofOfOwnership: &openapi.TgvalidatordProofOfOwnership{
			SignedPayloadHash:     &hash,
			SignedPayloadAsString: &payload,
		},
	}

	got := SharedAddressFromDTO(dto)
	if got == nil {
		t.Fatal("SharedAddressFromDTO() returned nil for non-nil input")
	}
	if got.ProofOfOwnership == nil {
		t.Fatal("ProofOfOwnership should not be nil")
	}
	if got.ProofOfOwnership.SignedPayloadHash != hash {
		t.Errorf("SignedPayloadHash = %v, want %v", got.ProofOfOwnership.SignedPayloadHash, hash)
	}
	if got.ProofOfOwnership.SignedPayloadAsString != payload {
		t.Errorf("SignedPayloadAsString = %v, want %v", got.ProofOfOwnership.SignedPayloadAsString, payload)
	}
}

func TestSharedAddressFromDTO_WithTrails(t *testing.T) {
	trailID := "trail-123"
	status := "accepted"
	dto := &openapi.TgvalidatordTnSharedAddress{
		Trails: []openapi.TgvalidatordTnSharedAddressTrail{
			{Id: &trailID, AddressStatus: &status},
		},
	}

	got := SharedAddressFromDTO(dto)
	if got == nil {
		t.Fatal("SharedAddressFromDTO() returned nil for non-nil input")
	}
	if len(got.Trails) != 1 {
		t.Fatalf("Trails length = %v, want 1", len(got.Trails))
	}
	if got.Trails[0].ID != trailID {
		t.Errorf("Trail ID = %v, want %v", got.Trails[0].ID, trailID)
	}
	if got.Trails[0].AddressStatus != status {
		t.Errorf("Trail AddressStatus = %v, want %v", got.Trails[0].AddressStatus, status)
	}
}

func TestSharedAssetFromDTO_WithTrails(t *testing.T) {
	trailID := "trail-456"
	status := "pending"
	dto := &openapi.TgvalidatordTnSharedAsset{
		Trails: []openapi.TgvalidatordTnSharedAssetTrail{
			{Id: &trailID, AssetStatus: &status},
		},
	}

	got := SharedAssetFromDTO(dto)
	if got == nil {
		t.Fatal("SharedAssetFromDTO() returned nil for non-nil input")
	}
	if len(got.Trails) != 1 {
		t.Fatalf("Trails length = %v, want 1", len(got.Trails))
	}
	if got.Trails[0].ID != trailID {
		t.Errorf("Trail ID = %v, want %v", got.Trails[0].ID, trailID)
	}
	if got.Trails[0].AssetStatus != status {
		t.Errorf("Trail AssetStatus = %v, want %v", got.Trails[0].AssetStatus, status)
	}
}

func TestSharedAddressFromDTO_PledgesCountParsing(t *testing.T) {
	tests := []struct {
		name         string
		pledgesCount *string
		want         int64
	}{
		{
			name:         "nil pledges count",
			pledgesCount: nil,
			want:         0,
		},
		{
			name:         "valid pledges count",
			pledgesCount: stringPtr("10"),
			want:         10,
		},
		{
			name:         "invalid pledges count",
			pledgesCount: stringPtr("invalid"),
			want:         0,
		},
		{
			name:         "zero pledges count",
			pledgesCount: stringPtr("0"),
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordTnSharedAddress{
				PledgesCount: tt.pledgesCount,
			}
			got := SharedAddressFromDTO(dto)
			if got.PledgesCount != tt.want {
				t.Errorf("PledgesCount = %v, want %v", got.PledgesCount, tt.want)
			}
		})
	}
}
