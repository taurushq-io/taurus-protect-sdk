package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestAssetFilterToDTO(t *testing.T) {
	tests := []struct {
		name   string
		filter *model.AssetFilter
	}{
		{
			name:   "nil input returns empty asset",
			filter: nil,
		},
		{
			name: "currency only",
			filter: &model.AssetFilter{
				Currency: "BTC",
			},
		},
		{
			name: "currency with kind",
			filter: &model.AssetFilter{
				Currency: "ETH",
				Kind:     "NFT",
			},
		},
		{
			name: "currency with NFT filter",
			filter: &model.AssetFilter{
				Currency: "ETH",
				Kind:     "NFT",
				NFT: &model.AssetNFTFilter{
					TokenID: "token-123",
				},
			},
		},
		{
			name: "currency with unknown filter",
			filter: &model.AssetFilter{
				Currency: "Unknown",
				Kind:     "Unknown",
				Unknown: &model.AssetUnknownFilter{
					Blockchain: "ethereum",
					Arg1:       "0x1234",
					Arg2:       "5678",
					Network:    "mainnet",
				},
			},
		},
		{
			name: "unknown filter with partial fields",
			filter: &model.AssetFilter{
				Currency: "Unknown",
				Kind:     "Unknown",
				Unknown: &model.AssetUnknownFilter{
					Blockchain: "polygon",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetFilterToDTO(tt.filter)
			if tt.filter == nil {
				if got.Currency != "" {
					t.Errorf("Currency = %v, want empty", got.Currency)
				}
				return
			}

			// Verify currency
			if got.Currency != tt.filter.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, tt.filter.Currency)
			}

			// Verify kind
			if tt.filter.Kind != "" {
				if got.Kind == nil || *got.Kind != tt.filter.Kind {
					t.Errorf("Kind = %v, want %v", safeString(got.Kind), tt.filter.Kind)
				}
			} else if got.Kind != nil {
				t.Errorf("Kind should be nil, got %v", *got.Kind)
			}

			// Verify NFT filter
			if tt.filter.NFT != nil {
				if got.Nft == nil {
					t.Error("Nft should not be nil")
				} else if tt.filter.NFT.TokenID != "" && safeString(got.Nft.Tokenid) != tt.filter.NFT.TokenID {
					t.Errorf("NFT TokenID = %v, want %v", safeString(got.Nft.Tokenid), tt.filter.NFT.TokenID)
				}
			}

			// Verify unknown filter
			if tt.filter.Unknown != nil {
				if got.Unknown == nil {
					t.Error("Unknown should not be nil")
				} else {
					if tt.filter.Unknown.Blockchain != "" && safeString(got.Unknown.Blockchain) != tt.filter.Unknown.Blockchain {
						t.Errorf("Unknown Blockchain = %v, want %v", safeString(got.Unknown.Blockchain), tt.filter.Unknown.Blockchain)
					}
					if tt.filter.Unknown.Arg1 != "" && safeString(got.Unknown.Arg1) != tt.filter.Unknown.Arg1 {
						t.Errorf("Unknown Arg1 = %v, want %v", safeString(got.Unknown.Arg1), tt.filter.Unknown.Arg1)
					}
					if tt.filter.Unknown.Arg2 != "" && safeString(got.Unknown.Arg2) != tt.filter.Unknown.Arg2 {
						t.Errorf("Unknown Arg2 = %v, want %v", safeString(got.Unknown.Arg2), tt.filter.Unknown.Arg2)
					}
					if tt.filter.Unknown.Network != "" && safeString(got.Unknown.Network) != tt.filter.Unknown.Network {
						t.Errorf("Unknown Network = %v, want %v", safeString(got.Unknown.Network), tt.filter.Unknown.Network)
					}
				}
			}
		})
	}
}

func TestAssetFilterFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAsset
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "currency only",
			dto: &openapi.TgvalidatordAsset{
				Currency: "BTC",
			},
		},
		{
			name: "currency with kind",
			dto: func() *openapi.TgvalidatordAsset {
				kind := "NFT"
				return &openapi.TgvalidatordAsset{
					Currency: "ETH",
					Kind:     &kind,
				}
			}(),
		},
		{
			name: "currency with NFT filter",
			dto: func() *openapi.TgvalidatordAsset {
				kind := "NFT"
				tokenID := "token-456"
				return &openapi.TgvalidatordAsset{
					Currency: "ETH",
					Kind:     &kind,
					Nft: &openapi.TgvalidatordAssetNFT{
						Tokenid: &tokenID,
					},
				}
			}(),
		},
		{
			name: "currency with unknown filter",
			dto: func() *openapi.TgvalidatordAsset {
				kind := "Unknown"
				blockchain := "ethereum"
				arg1 := "0xabcd"
				arg2 := "9999"
				network := "mainnet"
				return &openapi.TgvalidatordAsset{
					Currency: "Unknown",
					Kind:     &kind,
					Unknown: &openapi.AssetUnknown{
						Blockchain: &blockchain,
						Arg1:       &arg1,
						Arg2:       &arg2,
						Network:    &network,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetFilterFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("AssetFilterFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("AssetFilterFromDTO() returned nil for non-nil input")
			}

			// Verify currency
			if got.Currency != tt.dto.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, tt.dto.Currency)
			}

			// Verify kind
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}

			// Verify NFT filter
			if tt.dto.Nft != nil {
				if got.NFT == nil {
					t.Error("NFT should not be nil")
				} else if tt.dto.Nft.Tokenid != nil && got.NFT.TokenID != *tt.dto.Nft.Tokenid {
					t.Errorf("NFT TokenID = %v, want %v", got.NFT.TokenID, *tt.dto.Nft.Tokenid)
				}
			}

			// Verify unknown filter
			if tt.dto.Unknown != nil {
				if got.Unknown == nil {
					t.Error("Unknown should not be nil")
				} else {
					if tt.dto.Unknown.Blockchain != nil && got.Unknown.Blockchain != *tt.dto.Unknown.Blockchain {
						t.Errorf("Unknown Blockchain = %v, want %v", got.Unknown.Blockchain, *tt.dto.Unknown.Blockchain)
					}
					if tt.dto.Unknown.Arg1 != nil && got.Unknown.Arg1 != *tt.dto.Unknown.Arg1 {
						t.Errorf("Unknown Arg1 = %v, want %v", got.Unknown.Arg1, *tt.dto.Unknown.Arg1)
					}
					if tt.dto.Unknown.Arg2 != nil && got.Unknown.Arg2 != *tt.dto.Unknown.Arg2 {
						t.Errorf("Unknown Arg2 = %v, want %v", got.Unknown.Arg2, *tt.dto.Unknown.Arg2)
					}
					if tt.dto.Unknown.Network != nil && got.Unknown.Network != *tt.dto.Unknown.Network {
						t.Errorf("Unknown Network = %v, want %v", got.Unknown.Network, *tt.dto.Unknown.Network)
					}
				}
			}
		})
	}
}

func TestAssetFilterRoundTrip(t *testing.T) {
	// Test that converting to DTO and back preserves data
	original := &model.AssetFilter{
		Currency: "Unknown",
		Kind:     "Unknown",
		Unknown: &model.AssetUnknownFilter{
			Blockchain: "ethereum",
			Arg1:       "contract-123",
			Arg2:       "token-456",
			Network:    "mainnet",
		},
	}

	dto := AssetFilterToDTO(original)
	roundTrip := AssetFilterFromDTO(&dto)

	if roundTrip.Currency != original.Currency {
		t.Errorf("Currency = %v, want %v", roundTrip.Currency, original.Currency)
	}
	if roundTrip.Kind != original.Kind {
		t.Errorf("Kind = %v, want %v", roundTrip.Kind, original.Kind)
	}
	if roundTrip.Unknown == nil {
		t.Fatal("Unknown should not be nil")
	}
	if roundTrip.Unknown.Blockchain != original.Unknown.Blockchain {
		t.Errorf("Unknown.Blockchain = %v, want %v", roundTrip.Unknown.Blockchain, original.Unknown.Blockchain)
	}
	if roundTrip.Unknown.Arg1 != original.Unknown.Arg1 {
		t.Errorf("Unknown.Arg1 = %v, want %v", roundTrip.Unknown.Arg1, original.Unknown.Arg1)
	}
	if roundTrip.Unknown.Arg2 != original.Unknown.Arg2 {
		t.Errorf("Unknown.Arg2 = %v, want %v", roundTrip.Unknown.Arg2, original.Unknown.Arg2)
	}
	if roundTrip.Unknown.Network != original.Unknown.Network {
		t.Errorf("Unknown.Network = %v, want %v", roundTrip.Unknown.Network, original.Unknown.Network)
	}
}

func TestAssetFilterToDTO_EmptyNFTFilter(t *testing.T) {
	filter := &model.AssetFilter{
		Currency: "ETH",
		Kind:     "NFT",
		NFT:      &model.AssetNFTFilter{},
	}

	got := AssetFilterToDTO(filter)
	if got.Nft == nil {
		t.Error("Nft should not be nil")
	}
	// Empty NFT filter should still create the struct
	if got.Nft.Tokenid != nil {
		t.Errorf("Tokenid should be nil for empty filter, got %v", *got.Nft.Tokenid)
	}
}

func TestAssetFilterToDTO_EmptyUnknownFilter(t *testing.T) {
	filter := &model.AssetFilter{
		Currency: "Unknown",
		Kind:     "Unknown",
		Unknown:  &model.AssetUnknownFilter{},
	}

	got := AssetFilterToDTO(filter)
	if got.Unknown == nil {
		t.Error("Unknown should not be nil")
	}
	// Empty unknown filter should still create the struct with nil fields
	if got.Unknown.Blockchain != nil {
		t.Errorf("Blockchain should be nil for empty filter, got %v", *got.Unknown.Blockchain)
	}
	if got.Unknown.Arg1 != nil {
		t.Errorf("Arg1 should be nil for empty filter, got %v", *got.Unknown.Arg1)
	}
	if got.Unknown.Arg2 != nil {
		t.Errorf("Arg2 should be nil for empty filter, got %v", *got.Unknown.Arg2)
	}
	if got.Unknown.Network != nil {
		t.Errorf("Network should be nil for empty filter, got %v", *got.Unknown.Network)
	}
}
