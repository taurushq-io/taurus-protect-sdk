package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestCurrencyFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrency
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
			name: "empty DTO returns currency with zero values",
			dto:  &openapi.TgvalidatordCurrency{},
			want: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Error("expected non-nil currency")
				}
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordCurrency {
				id := "currency-123"
				name := "ethereum"
				symbol := "ETH"
				displayName := "Ethereum"
				currencyType := "native"
				blockchain := "ETH"
				network := "mainnet"
				decimals := "18"
				coinTypeIndex := "60"
				contractAddress := "0x1234567890abcdef"
				tokenID := "token-001"
				wlcaId := "42"
				logo := "data:image/png;base64,abc123"
				isToken := false
				isERC20 := false
				isFA12 := false
				isFA20 := false
				isNFT := false
				isFiat := false
				isUTXOBased := false
				isAccountBased := true
				hasStaking := true
				enabled := true
				return &openapi.TgvalidatordCurrency{
					Id:              &id,
					Name:            &name,
					Symbol:          &symbol,
					DisplayName:     &displayName,
					Type:            &currencyType,
					Blockchain:      &blockchain,
					Network:         &network,
					Decimals:        &decimals,
					CoinTypeIndex:   &coinTypeIndex,
					ContractAddress: &contractAddress,
					TokenID:         &tokenID,
					WlcaId:          &wlcaId,
					Logo:            &logo,
					IsToken:         &isToken,
					IsERC20:         &isERC20,
					IsFA12:          &isFA12,
					IsFA20:          &isFA20,
					IsNFT:           &isNFT,
					IsFiat:          &isFiat,
					IsUTXOBased:     &isUTXOBased,
					IsAccountBased:  &isAccountBased,
					HasStaking:      &hasStaking,
					Enabled:         &enabled,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Validated in main test body
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrencyFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("CurrencyFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("CurrencyFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
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
			if tt.dto.Type != nil && got.Type != *tt.dto.Type {
				t.Errorf("Type = %v, want %v", got.Type, *tt.dto.Type)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.ContractAddress != nil && got.ContractAddress != *tt.dto.ContractAddress {
				t.Errorf("ContractAddress = %v, want %v", got.ContractAddress, *tt.dto.ContractAddress)
			}
			if tt.dto.TokenID != nil && got.TokenID != *tt.dto.TokenID {
				t.Errorf("TokenID = %v, want %v", got.TokenID, *tt.dto.TokenID)
			}
			if tt.dto.WlcaId != nil && got.WlcaID != 42 {
				t.Errorf("WlcaID = %v, want %v", got.WlcaID, 42)
			}
			if tt.dto.Logo != nil && got.Logo != *tt.dto.Logo {
				t.Errorf("Logo = %v, want %v", got.Logo, *tt.dto.Logo)
			}
			if tt.dto.IsToken != nil && got.IsToken != *tt.dto.IsToken {
				t.Errorf("IsToken = %v, want %v", got.IsToken, *tt.dto.IsToken)
			}
			if tt.dto.IsERC20 != nil && got.IsERC20 != *tt.dto.IsERC20 {
				t.Errorf("IsERC20 = %v, want %v", got.IsERC20, *tt.dto.IsERC20)
			}
			if tt.dto.IsFA12 != nil && got.IsFA12 != *tt.dto.IsFA12 {
				t.Errorf("IsFA12 = %v, want %v", got.IsFA12, *tt.dto.IsFA12)
			}
			if tt.dto.IsFA20 != nil && got.IsFA20 != *tt.dto.IsFA20 {
				t.Errorf("IsFA20 = %v, want %v", got.IsFA20, *tt.dto.IsFA20)
			}
			if tt.dto.IsNFT != nil && got.IsNFT != *tt.dto.IsNFT {
				t.Errorf("IsNFT = %v, want %v", got.IsNFT, *tt.dto.IsNFT)
			}
			if tt.dto.IsFiat != nil && got.IsFiat != *tt.dto.IsFiat {
				t.Errorf("IsFiat = %v, want %v", got.IsFiat, *tt.dto.IsFiat)
			}
			if tt.dto.IsUTXOBased != nil && got.IsUTXOBased != *tt.dto.IsUTXOBased {
				t.Errorf("IsUTXOBased = %v, want %v", got.IsUTXOBased, *tt.dto.IsUTXOBased)
			}
			if tt.dto.IsAccountBased != nil && got.IsAccountBased != *tt.dto.IsAccountBased {
				t.Errorf("IsAccountBased = %v, want %v", got.IsAccountBased, *tt.dto.IsAccountBased)
			}
			if tt.dto.HasStaking != nil && got.HasStaking != *tt.dto.HasStaking {
				t.Errorf("HasStaking = %v, want %v", got.HasStaking, *tt.dto.HasStaking)
			}
			if tt.dto.Enabled != nil && got.Enabled != *tt.dto.Enabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, *tt.dto.Enabled)
			}
		})
	}
}

func TestCurrencyFromDTO_DecimalsParsing(t *testing.T) {
	tests := []struct {
		name     string
		decimals string
		want     int64
	}{
		{"valid number", "18", 18},
		{"zero", "0", 0},
		{"large number", "24", 24},
		{"invalid string", "invalid", 0},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordCurrency{
				Decimals: &tt.decimals,
			}
			got := CurrencyFromDTO(dto)
			if got.Decimals != tt.want {
				t.Errorf("Decimals = %v, want %v", got.Decimals, tt.want)
			}
		})
	}
}

func TestCurrencyFromDTO_CoinTypeIndexPassthrough(t *testing.T) {
	tests := []struct {
		name          string
		coinTypeIndex string
		want          string
	}{
		{"bitcoin", "0", "0"},
		{"ethereum", "60", "60"},
		{"large number", "999999", "999999"},
		{"non-numeric string", "invalid", "invalid"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordCurrency{
				CoinTypeIndex: &tt.coinTypeIndex,
			}
			got := CurrencyFromDTO(dto)
			if got.CoinTypeIndex != tt.want {
				t.Errorf("CoinTypeIndex = %v, want %v", got.CoinTypeIndex, tt.want)
			}
		})
	}
}

func TestCurrencyFromDTO_WlcaIDParsing(t *testing.T) {
	tests := []struct {
		name   string
		wlcaId string
		want   int64
	}{
		{"valid number", "123", 123},
		{"zero", "0", 0},
		{"large number", "999999", 999999},
		{"invalid string", "invalid", 0},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordCurrency{
				WlcaId: &tt.wlcaId,
			}
			got := CurrencyFromDTO(dto)
			if got.WlcaID != tt.want {
				t.Errorf("WlcaID = %v, want %v", got.WlcaID, tt.want)
			}
		})
	}
}

func TestCurrenciesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordCurrency
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordCurrency{},
			want: 0,
		},
		{
			name: "converts multiple currencies",
			dtos: func() []openapi.TgvalidatordCurrency {
				id1 := "currency-1"
				id2 := "currency-2"
				return []openapi.TgvalidatordCurrency{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrenciesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("CurrenciesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("CurrenciesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestCurrenciesFromDTO_PreservesOrder(t *testing.T) {
	id1 := "first"
	id2 := "second"
	id3 := "third"
	dtos := []openapi.TgvalidatordCurrency{
		{Id: &id1},
		{Id: &id2},
		{Id: &id3},
	}

	got := CurrenciesFromDTO(dtos)

	if len(got) != 3 {
		t.Fatalf("CurrenciesFromDTO() length = %v, want 3", len(got))
	}
	if got[0].ID != "first" {
		t.Errorf("got[0].ID = %v, want 'first'", got[0].ID)
	}
	if got[1].ID != "second" {
		t.Errorf("got[1].ID = %v, want 'second'", got[1].ID)
	}
	if got[2].ID != "third" {
		t.Errorf("got[2].ID = %v, want 'third'", got[2].ID)
	}
}

func TestCurrencyFromDTO_BooleanFlags(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrency
		want map[string]bool
	}{
		{
			name: "all true flags",
			dto: &openapi.TgvalidatordCurrency{
				IsToken:        &trueVal,
				IsERC20:        &trueVal,
				IsFA12:         &trueVal,
				IsFA20:         &trueVal,
				IsNFT:          &trueVal,
				IsFiat:         &trueVal,
				IsUTXOBased:    &trueVal,
				IsAccountBased: &trueVal,
				HasStaking:     &trueVal,
				Enabled:        &trueVal,
			},
			want: map[string]bool{
				"IsToken":        true,
				"IsERC20":        true,
				"IsFA12":         true,
				"IsFA20":         true,
				"IsNFT":          true,
				"IsFiat":         true,
				"IsUTXOBased":    true,
				"IsAccountBased": true,
				"HasStaking":     true,
				"Enabled":        true,
			},
		},
		{
			name: "all false flags",
			dto: &openapi.TgvalidatordCurrency{
				IsToken:        &falseVal,
				IsERC20:        &falseVal,
				IsFA12:         &falseVal,
				IsFA20:         &falseVal,
				IsNFT:          &falseVal,
				IsFiat:         &falseVal,
				IsUTXOBased:    &falseVal,
				IsAccountBased: &falseVal,
				HasStaking:     &falseVal,
				Enabled:        &falseVal,
			},
			want: map[string]bool{
				"IsToken":        false,
				"IsERC20":        false,
				"IsFA12":         false,
				"IsFA20":         false,
				"IsNFT":          false,
				"IsFiat":         false,
				"IsUTXOBased":    false,
				"IsAccountBased": false,
				"HasStaking":     false,
				"Enabled":        false,
			},
		},
		{
			name: "nil flags default to false",
			dto:  &openapi.TgvalidatordCurrency{},
			want: map[string]bool{
				"IsToken":        false,
				"IsERC20":        false,
				"IsFA12":         false,
				"IsFA20":         false,
				"IsNFT":          false,
				"IsFiat":         false,
				"IsUTXOBased":    false,
				"IsAccountBased": false,
				"HasStaking":     false,
				"Enabled":        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrencyFromDTO(tt.dto)

			if got.IsToken != tt.want["IsToken"] {
				t.Errorf("IsToken = %v, want %v", got.IsToken, tt.want["IsToken"])
			}
			if got.IsERC20 != tt.want["IsERC20"] {
				t.Errorf("IsERC20 = %v, want %v", got.IsERC20, tt.want["IsERC20"])
			}
			if got.IsFA12 != tt.want["IsFA12"] {
				t.Errorf("IsFA12 = %v, want %v", got.IsFA12, tt.want["IsFA12"])
			}
			if got.IsFA20 != tt.want["IsFA20"] {
				t.Errorf("IsFA20 = %v, want %v", got.IsFA20, tt.want["IsFA20"])
			}
			if got.IsNFT != tt.want["IsNFT"] {
				t.Errorf("IsNFT = %v, want %v", got.IsNFT, tt.want["IsNFT"])
			}
			if got.IsFiat != tt.want["IsFiat"] {
				t.Errorf("IsFiat = %v, want %v", got.IsFiat, tt.want["IsFiat"])
			}
			if got.IsUTXOBased != tt.want["IsUTXOBased"] {
				t.Errorf("IsUTXOBased = %v, want %v", got.IsUTXOBased, tt.want["IsUTXOBased"])
			}
			if got.IsAccountBased != tt.want["IsAccountBased"] {
				t.Errorf("IsAccountBased = %v, want %v", got.IsAccountBased, tt.want["IsAccountBased"])
			}
			if got.HasStaking != tt.want["HasStaking"] {
				t.Errorf("HasStaking = %v, want %v", got.HasStaking, tt.want["HasStaking"])
			}
			if got.Enabled != tt.want["Enabled"] {
				t.Errorf("Enabled = %v, want %v", got.Enabled, tt.want["Enabled"])
			}
		})
	}
}
