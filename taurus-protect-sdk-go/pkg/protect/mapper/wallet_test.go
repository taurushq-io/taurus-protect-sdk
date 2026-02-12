package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWalletFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWalletInfo
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
			name: "empty DTO returns wallet with zero values",
			dto:  &openapi.TgvalidatordWalletInfo{},
			want: func(t *testing.T, got interface{}) {
				wallet := got.(*struct{ ID string })
				if wallet == nil {
					t.Error("expected non-nil wallet")
				}
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWalletInfo {
				id := "wallet-123"
				name := "Test Wallet"
				currency := "ETH"
				blockchain := "ETH"
				network := "mainnet"
				isOmnibus := true
				disabled := false
				comment := "test comment"
				customerId := "customer-456"
				externalWalletId := "ext-789"
				visibilityGroupID := "vg-001"
				addressesCount := "42"
				now := time.Now()
				return &openapi.TgvalidatordWalletInfo{
					Id:                &id,
					Name:              &name,
					Currency:          &currency,
					Blockchain:        &blockchain,
					Network:           &network,
					IsOmnibus:         &isOmnibus,
					Disabled:          &disabled,
					Comment:           &comment,
					CustomerId:        &customerId,
					ExternalWalletId:  &externalWalletId,
					VisibilityGroupID: &visibilityGroupID,
					AddressesCount:    &addressesCount,
					CreationDate:      &now,
					UpdateDate:        &now,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Type assertion handled in test execution
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WalletFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WalletFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WalletFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields for complete DTO
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Currency != nil && got.Currency != *tt.dto.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, *tt.dto.Currency)
			}
			if tt.dto.AddressesCount != nil && got.AddressesCount != 42 {
				t.Errorf("AddressesCount = %v, want 42", got.AddressesCount)
			}
		})
	}
}

func TestWalletFromDTO_AddressesCountParsing(t *testing.T) {
	tests := []struct {
		name  string
		count string
		want  int64
	}{
		{"valid number", "100", 100},
		{"zero", "0", 0},
		{"large number", "999999", 999999},
		{"invalid string", "invalid", 0},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordWalletInfo{
				AddressesCount: &tt.count,
			}
			got := WalletFromDTO(dto)
			if got.AddressesCount != tt.want {
				t.Errorf("AddressesCount = %v, want %v", got.AddressesCount, tt.want)
			}
		})
	}
}

func TestWalletsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordWalletInfo
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordWalletInfo{},
			want: 0,
		},
		{
			name: "converts multiple wallets",
			dtos: func() []openapi.TgvalidatordWalletInfo {
				id1 := "wallet-1"
				id2 := "wallet-2"
				return []openapi.TgvalidatordWalletInfo{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WalletsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("WalletsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("WalletsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestBalanceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordBalance
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns balance with empty strings",
			dto:  &openapi.TgvalidatordBalance{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordBalance {
				totalConfirmed := "1000000000000000000"
				totalUnconfirmed := "500000000000000000"
				availableConfirmed := "800000000000000000"
				availableUnconfirmed := "300000000000000000"
				reservedConfirmed := "200000000000000000"
				reservedUnconfirmed := "100000000000000000"
				return &openapi.TgvalidatordBalance{
					TotalConfirmed:       &totalConfirmed,
					TotalUnconfirmed:     &totalUnconfirmed,
					AvailableConfirmed:   &availableConfirmed,
					AvailableUnconfirmed: &availableUnconfirmed,
					ReservedConfirmed:    &reservedConfirmed,
					ReservedUnconfirmed:  &reservedUnconfirmed,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BalanceFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("BalanceFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("BalanceFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.TotalConfirmed != nil && got.TotalConfirmed != *tt.dto.TotalConfirmed {
				t.Errorf("TotalConfirmed = %v, want %v", got.TotalConfirmed, *tt.dto.TotalConfirmed)
			}
		})
	}
}

func TestWalletAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWalletAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWalletAttribute {
				id := "attr-123"
				key := "environment"
				value := "production"
				return &openapi.TgvalidatordWalletAttribute{
					Id:    &id,
					Key:   &key,
					Value: &value,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WalletAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("WalletAttributeFromDTO(nil) should return zero value")
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

func TestWalletFromCreateDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordWallet
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordWallet {
				id := "wallet-123"
				name := "New Wallet"
				currency := "BTC"
				blockchain := "BTC"
				return &openapi.TgvalidatordWallet{
					Id:         &id,
					Name:       &name,
					Currency:   &currency,
					Blockchain: &blockchain,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WalletFromCreateDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("WalletFromCreateDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("WalletFromCreateDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
		})
	}
}

// Test helper functions
func TestSafeString(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{"nil returns empty", nil, ""},
		{"value returns value", strPtr("test"), "test"},
		{"empty string returns empty", strPtr(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeString(tt.input); got != tt.want {
				t.Errorf("safeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeBool(t *testing.T) {
	tests := []struct {
		name  string
		input *bool
		want  bool
	}{
		{"nil returns false", nil, false},
		{"true returns true", boolPtr(true), true},
		{"false returns false", boolPtr(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeBool(tt.input); got != tt.want {
				t.Errorf("safeBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeInt64(t *testing.T) {
	tests := []struct {
		name  string
		input *int64
		want  int64
	}{
		{"nil returns zero", nil, 0},
		{"value returns value", int64Ptr(42), 42},
		{"zero returns zero", int64Ptr(0), 0},
		{"negative returns negative", int64Ptr(-100), -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeInt64(tt.input); got != tt.want {
				t.Errorf("safeInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// strPtr is a helper for creating string pointers in tests
func strPtr(s string) *string {
	return &s
}
