package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestFiatProviderFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordFiatProvider
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns fiat provider with zero values",
			dto:  &openapi.TgvalidatordFiatProvider{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordFiatProvider {
				provider := "circle"
				label := "main-account"
				valuation := "1000000.50"
				return &openapi.TgvalidatordFiatProvider{
					Provider:              &provider,
					Label:                 &label,
					BaseCurrencyValuation: &valuation,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FiatProviderFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FiatProviderFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FiatProviderFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Provider != nil && got.Provider != *tt.dto.Provider {
				t.Errorf("Provider = %v, want %v", got.Provider, *tt.dto.Provider)
			}
			if tt.dto.Label != nil && got.Label != *tt.dto.Label {
				t.Errorf("Label = %v, want %v", got.Label, *tt.dto.Label)
			}
			if tt.dto.BaseCurrencyValuation != nil && got.BaseCurrencyValuation != *tt.dto.BaseCurrencyValuation {
				t.Errorf("BaseCurrencyValuation = %v, want %v", got.BaseCurrencyValuation, *tt.dto.BaseCurrencyValuation)
			}
		})
	}
}

func TestFiatProvidersFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordFiatProvider
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordFiatProvider{},
			want: 0,
		},
		{
			name: "converts multiple fiat providers",
			dtos: func() []openapi.TgvalidatordFiatProvider {
				provider1 := "circle"
				provider2 := "cubnet"
				return []openapi.TgvalidatordFiatProvider{
					{Provider: &provider1},
					{Provider: &provider2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FiatProvidersFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("FiatProvidersFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("FiatProvidersFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFiatProviderAccountFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordFiatProviderAccount
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns account with zero values",
			dto:  &openapi.TgvalidatordFiatProviderAccount{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordFiatProviderAccount {
				id := "account-123"
				provider := "circle"
				label := "main-account"
				accountType := "wallet"
				accountIdentifier := "wallet-456"
				accountName := "Main Wallet"
				totalBalance := "1000000"
				currencyID := "USD"
				valuation := "1000000.00"
				creationDate := time.Now().Add(-24 * time.Hour)
				updateDate := time.Now()
				return &openapi.TgvalidatordFiatProviderAccount{
					Id:                    &id,
					Provider:              &provider,
					Label:                 &label,
					AccountType:           &accountType,
					AccountIdentifier:     &accountIdentifier,
					AccountName:           &accountName,
					TotalBalance:          &totalBalance,
					CurrencyID:            &currencyID,
					BaseCurrencyValuation: &valuation,
					CreationDate:          &creationDate,
					UpdateDate:            &updateDate,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FiatProviderAccountFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FiatProviderAccountFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FiatProviderAccountFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Provider != nil && got.Provider != *tt.dto.Provider {
				t.Errorf("Provider = %v, want %v", got.Provider, *tt.dto.Provider)
			}
			if tt.dto.Label != nil && got.Label != *tt.dto.Label {
				t.Errorf("Label = %v, want %v", got.Label, *tt.dto.Label)
			}
			if tt.dto.AccountType != nil && got.AccountType != *tt.dto.AccountType {
				t.Errorf("AccountType = %v, want %v", got.AccountType, *tt.dto.AccountType)
			}
			if tt.dto.AccountIdentifier != nil && got.AccountIdentifier != *tt.dto.AccountIdentifier {
				t.Errorf("AccountIdentifier = %v, want %v", got.AccountIdentifier, *tt.dto.AccountIdentifier)
			}
			if tt.dto.AccountName != nil && got.AccountName != *tt.dto.AccountName {
				t.Errorf("AccountName = %v, want %v", got.AccountName, *tt.dto.AccountName)
			}
			if tt.dto.TotalBalance != nil && got.TotalBalance != *tt.dto.TotalBalance {
				t.Errorf("TotalBalance = %v, want %v", got.TotalBalance, *tt.dto.TotalBalance)
			}
			if tt.dto.CurrencyID != nil && got.CurrencyID != *tt.dto.CurrencyID {
				t.Errorf("CurrencyID = %v, want %v", got.CurrencyID, *tt.dto.CurrencyID)
			}
			if tt.dto.BaseCurrencyValuation != nil && got.BaseCurrencyValuation != *tt.dto.BaseCurrencyValuation {
				t.Errorf("BaseCurrencyValuation = %v, want %v", got.BaseCurrencyValuation, *tt.dto.BaseCurrencyValuation)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
			if tt.dto.UpdateDate != nil && !got.UpdateDate.Equal(*tt.dto.UpdateDate) {
				t.Errorf("UpdateDate = %v, want %v", got.UpdateDate, *tt.dto.UpdateDate)
			}
		})
	}
}

func TestFiatProviderAccountFromDTO_WithCurrencyInfo(t *testing.T) {
	currencyId := "usd"
	currencyName := "US Dollar"
	dto := &openapi.TgvalidatordFiatProviderAccount{
		CurrencyInfo: &openapi.TgvalidatordCurrency{
			Id:   &currencyId,
			Name: &currencyName,
		},
	}

	got := FiatProviderAccountFromDTO(dto)
	if got == nil {
		t.Fatal("FiatProviderAccountFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo == nil {
		t.Error("CurrencyInfo should not be nil when DTO has currency info")
	}
	if got.CurrencyInfo.ID != currencyId {
		t.Errorf("CurrencyInfo.ID = %v, want %v", got.CurrencyInfo.ID, currencyId)
	}
	if got.CurrencyInfo.Name != currencyName {
		t.Errorf("CurrencyInfo.Name = %v, want %v", got.CurrencyInfo.Name, currencyName)
	}
}

func TestFiatProviderAccountFromDTO_NilCurrencyInfo(t *testing.T) {
	id := "account-123"
	dto := &openapi.TgvalidatordFiatProviderAccount{
		Id:           &id,
		CurrencyInfo: nil,
	}

	got := FiatProviderAccountFromDTO(dto)
	if got == nil {
		t.Fatal("FiatProviderAccountFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo != nil {
		t.Errorf("CurrencyInfo should be nil when DTO currency info is nil, got %v", got.CurrencyInfo)
	}
	if got.ID != "account-123" {
		t.Errorf("ID = %v, want account-123", got.ID)
	}
}

func TestFiatProviderAccountFromDTO_NilDates(t *testing.T) {
	id := "account-123"
	dto := &openapi.TgvalidatordFiatProviderAccount{
		Id:           &id,
		CreationDate: nil,
		UpdateDate:   nil,
	}

	got := FiatProviderAccountFromDTO(dto)
	if got == nil {
		t.Fatal("FiatProviderAccountFromDTO() returned nil for non-nil input")
	}
	// When dates are nil, they should be zero time values
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate should be zero time when nil, got %v", got.CreationDate)
	}
	if !got.UpdateDate.IsZero() {
		t.Errorf("UpdateDate should be zero time when nil, got %v", got.UpdateDate)
	}
}

func TestFiatProviderAccountsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordFiatProviderAccount
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordFiatProviderAccount{},
			want: 0,
		},
		{
			name: "converts multiple accounts",
			dtos: func() []openapi.TgvalidatordFiatProviderAccount {
				id1 := "account-1"
				id2 := "account-2"
				return []openapi.TgvalidatordFiatProviderAccount{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FiatProviderAccountsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("FiatProviderAccountsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("FiatProviderAccountsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}
