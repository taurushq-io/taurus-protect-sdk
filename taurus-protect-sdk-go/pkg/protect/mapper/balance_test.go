package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestAssetBalanceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAssetBalance
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns asset balance with nil fields",
			dto:  &openapi.TgvalidatordAssetBalance{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAssetBalance {
				totalConfirmed := "1000000000000000000"
				return &openapi.TgvalidatordAssetBalance{
					Asset: &openapi.TgvalidatordAsset{
						Currency: "ETH",
					},
					Balance: &openapi.TgvalidatordBalance{
						TotalConfirmed: &totalConfirmed,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetBalanceFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("AssetBalanceFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("AssetBalanceFromDTO() returned nil for non-nil input")
			}
			// Verify asset is mapped if present
			if tt.dto.Asset != nil {
				if got.Asset == nil {
					t.Error("Asset should not be nil when DTO has asset")
				} else if got.Asset.Currency != tt.dto.Asset.Currency {
					t.Errorf("Asset.Currency = %v, want %v", got.Asset.Currency, tt.dto.Asset.Currency)
				}
			}
			// Verify balance is mapped if present
			if tt.dto.Balance != nil && tt.dto.Balance.TotalConfirmed != nil {
				if got.Balance == nil {
					t.Error("Balance should not be nil when DTO has balance")
				} else if got.Balance.TotalConfirmed != *tt.dto.Balance.TotalConfirmed {
					t.Errorf("Balance.TotalConfirmed = %v, want %v", got.Balance.TotalConfirmed, *tt.dto.Balance.TotalConfirmed)
				}
			}
		})
	}
}

func TestAssetBalancesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordAssetBalance
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordAssetBalance{},
			want: 0,
		},
		{
			name: "converts multiple balances",
			dtos: func() []openapi.TgvalidatordAssetBalance {
				return []openapi.TgvalidatordAssetBalance{
					{Asset: &openapi.TgvalidatordAsset{Currency: "ETH"}},
					{Asset: &openapi.TgvalidatordAsset{Currency: "BTC"}},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetBalancesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("AssetBalancesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("AssetBalancesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestAssetFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAsset
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "maps currency field",
			dto: &openapi.TgvalidatordAsset{
				Currency: "ETH",
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAsset {
				kind := "token"
				id := "currency-123"
				name := "Ethereum"
				symbol := "ETH"
				displayName := "Ethereum"
				blockchain := "ETH"
				network := "mainnet"
				decimals := "18"
				isToken := false
				enabled := true
				return &openapi.TgvalidatordAsset{
					Currency: "ETH",
					Kind:     &kind,
					CurrencyInfo: &openapi.TgvalidatordCurrency{
						Id:          &id,
						Name:        &name,
						Symbol:      &symbol,
						DisplayName: &displayName,
						Blockchain:  &blockchain,
						Network:     &network,
						Decimals:    &decimals,
						IsToken:     &isToken,
						Enabled:     &enabled,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssetFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("AssetFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("AssetFromDTO() returned nil for non-nil input")
			}
			if got.Currency != tt.dto.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, tt.dto.Currency)
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
		})
	}
}

func TestCurrencyInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrency
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns currency info with zero values",
			dto:  &openapi.TgvalidatordCurrency{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordCurrency {
				id := "currency-123"
				name := "Ethereum"
				symbol := "ETH"
				displayName := "Ethereum"
				blockchain := "ETH"
				network := "mainnet"
				decimals := "18"
				contractAddress := "0x1234567890abcdef"
				tokenID := "token-001"
				currType := "native"
				isToken := false
				isERC20 := false
				isNFT := false
				isFiat := false
				enabled := true
				return &openapi.TgvalidatordCurrency{
					Id:              &id,
					Name:            &name,
					Symbol:          &symbol,
					DisplayName:     &displayName,
					Blockchain:      &blockchain,
					Network:         &network,
					Decimals:        &decimals,
					ContractAddress: &contractAddress,
					TokenID:         &tokenID,
					Type:            &currType,
					IsToken:         &isToken,
					IsERC20:         &isERC20,
					IsNFT:           &isNFT,
					IsFiat:          &isFiat,
					Enabled:         &enabled,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrencyInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CurrencyInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CurrencyInfoFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Symbol != nil && got.Symbol != *tt.dto.Symbol {
				t.Errorf("Symbol = %v, want %v", got.Symbol, *tt.dto.Symbol)
			}
			if tt.dto.DisplayName != nil && got.DisplayName != *tt.dto.DisplayName {
				t.Errorf("DisplayName = %v, want %v", got.DisplayName, *tt.dto.DisplayName)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Decimals != nil && got.Decimals != *tt.dto.Decimals {
				t.Errorf("Decimals = %v, want %v", got.Decimals, *tt.dto.Decimals)
			}
			if tt.dto.ContractAddress != nil && got.ContractAddress != *tt.dto.ContractAddress {
				t.Errorf("ContractAddress = %v, want %v", got.ContractAddress, *tt.dto.ContractAddress)
			}
			if tt.dto.TokenID != nil && got.TokenID != *tt.dto.TokenID {
				t.Errorf("TokenID = %v, want %v", got.TokenID, *tt.dto.TokenID)
			}
			if tt.dto.Type != nil && got.Type != *tt.dto.Type {
				t.Errorf("Type = %v, want %v", got.Type, *tt.dto.Type)
			}
			if tt.dto.IsToken != nil && got.IsToken != *tt.dto.IsToken {
				t.Errorf("IsToken = %v, want %v", got.IsToken, *tt.dto.IsToken)
			}
			if tt.dto.IsERC20 != nil && got.IsERC20 != *tt.dto.IsERC20 {
				t.Errorf("IsERC20 = %v, want %v", got.IsERC20, *tt.dto.IsERC20)
			}
			if tt.dto.IsNFT != nil && got.IsNFT != *tt.dto.IsNFT {
				t.Errorf("IsNFT = %v, want %v", got.IsNFT, *tt.dto.IsNFT)
			}
			if tt.dto.IsFiat != nil && got.IsFiat != *tt.dto.IsFiat {
				t.Errorf("IsFiat = %v, want %v", got.IsFiat, *tt.dto.IsFiat)
			}
			if tt.dto.Enabled != nil && got.Enabled != *tt.dto.Enabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, *tt.dto.Enabled)
			}
		})
	}
}

func TestCurrencyInfoFromDTO_BooleanFields(t *testing.T) {
	tests := []struct {
		name      string
		isToken   *bool
		isERC20   *bool
		isNFT     *bool
		isFiat    *bool
		enabled   *bool
		wantToken bool
		wantERC20 bool
		wantNFT   bool
		wantFiat  bool
		wantEnabled bool
	}{
		{
			name:        "nil booleans default to false",
			isToken:     nil,
			isERC20:     nil,
			isNFT:       nil,
			isFiat:      nil,
			enabled:     nil,
			wantToken:   false,
			wantERC20:   false,
			wantNFT:     false,
			wantFiat:    false,
			wantEnabled: false,
		},
		{
			name:        "true values",
			isToken:     boolPtr(true),
			isERC20:     boolPtr(true),
			isNFT:       boolPtr(true),
			isFiat:      boolPtr(true),
			enabled:     boolPtr(true),
			wantToken:   true,
			wantERC20:   true,
			wantNFT:     true,
			wantFiat:    true,
			wantEnabled: true,
		},
		{
			name:        "false values",
			isToken:     boolPtr(false),
			isERC20:     boolPtr(false),
			isNFT:       boolPtr(false),
			isFiat:      boolPtr(false),
			enabled:     boolPtr(false),
			wantToken:   false,
			wantERC20:   false,
			wantNFT:     false,
			wantFiat:    false,
			wantEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordCurrency{
				IsToken: tt.isToken,
				IsERC20: tt.isERC20,
				IsNFT:   tt.isNFT,
				IsFiat:  tt.isFiat,
				Enabled: tt.enabled,
			}
			got := CurrencyInfoFromDTO(dto)
			if got.IsToken != tt.wantToken {
				t.Errorf("IsToken = %v, want %v", got.IsToken, tt.wantToken)
			}
			if got.IsERC20 != tt.wantERC20 {
				t.Errorf("IsERC20 = %v, want %v", got.IsERC20, tt.wantERC20)
			}
			if got.IsNFT != tt.wantNFT {
				t.Errorf("IsNFT = %v, want %v", got.IsNFT, tt.wantNFT)
			}
			if got.IsFiat != tt.wantFiat {
				t.Errorf("IsFiat = %v, want %v", got.IsFiat, tt.wantFiat)
			}
			if got.Enabled != tt.wantEnabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, tt.wantEnabled)
			}
		})
	}
}
