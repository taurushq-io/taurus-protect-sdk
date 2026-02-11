package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestExchangeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordExchange
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns exchange with zero values",
			dto:  &openapi.TgvalidatordExchange{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordExchange {
				id := "exchange-123"
				exchangeName := "Binance"
				account := "acc-456"
				currency := "BTC"
				typ := "spot"
				totalBalance := "1000000000"
				status := "active"
				container := "container-789"
				label := "My Exchange"
				displayLabel := "My Binance Account"
				hasWLA := true
				baseCurrencyValuation := "50000.00"
				creationDate := time.Now().Add(-24 * time.Hour)
				updateDate := time.Now()
				currencyID := "btc-mainnet"
				currencyName := "Bitcoin"
				return &openapi.TgvalidatordExchange{
					Id:                    &id,
					Exchange:              &exchangeName,
					Account:               &account,
					Currency:              &currency,
					Type:                  &typ,
					TotalBalance:          &totalBalance,
					Status:                &status,
					Container:             &container,
					Label:                 &label,
					DisplayLabel:          &displayLabel,
					HasWLA:                &hasWLA,
					BaseCurrencyValuation: &baseCurrencyValuation,
					CreationDate:          &creationDate,
					UpdateDate:            &updateDate,
					CurrencyInfo: &openapi.TgvalidatordCurrency{
						Id:   &currencyID,
						Name: &currencyName,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExchangeFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ExchangeFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ExchangeFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Exchange != nil && got.Exchange != *tt.dto.Exchange {
				t.Errorf("Exchange = %v, want %v", got.Exchange, *tt.dto.Exchange)
			}
			if tt.dto.Account != nil && got.Account != *tt.dto.Account {
				t.Errorf("Account = %v, want %v", got.Account, *tt.dto.Account)
			}
			if tt.dto.Currency != nil && got.Currency != *tt.dto.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, *tt.dto.Currency)
			}
			if tt.dto.Type != nil && got.Type != *tt.dto.Type {
				t.Errorf("Type = %v, want %v", got.Type, *tt.dto.Type)
			}
			if tt.dto.TotalBalance != nil && got.TotalBalance != *tt.dto.TotalBalance {
				t.Errorf("TotalBalance = %v, want %v", got.TotalBalance, *tt.dto.TotalBalance)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.Container != nil && got.Container != *tt.dto.Container {
				t.Errorf("Container = %v, want %v", got.Container, *tt.dto.Container)
			}
			if tt.dto.Label != nil && got.Label != *tt.dto.Label {
				t.Errorf("Label = %v, want %v", got.Label, *tt.dto.Label)
			}
			if tt.dto.DisplayLabel != nil && got.DisplayLabel != *tt.dto.DisplayLabel {
				t.Errorf("DisplayLabel = %v, want %v", got.DisplayLabel, *tt.dto.DisplayLabel)
			}
			if tt.dto.HasWLA != nil && got.HasWLA != *tt.dto.HasWLA {
				t.Errorf("HasWLA = %v, want %v", got.HasWLA, *tt.dto.HasWLA)
			}
			if tt.dto.BaseCurrencyValuation != nil && got.BaseCurrencyValuation != *tt.dto.BaseCurrencyValuation {
				t.Errorf("BaseCurrencyValuation = %v, want %v", got.BaseCurrencyValuation, *tt.dto.BaseCurrencyValuation)
			}
			// Verify currency info is mapped if present
			if tt.dto.CurrencyInfo != nil {
				if got.CurrencyInfo == nil {
					t.Error("CurrencyInfo should not be nil when DTO has currency info")
				}
			}
		})
	}
}

func TestExchangesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordExchange
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordExchange{},
			want: 0,
		},
		{
			name: "converts multiple exchanges",
			dtos: func() []openapi.TgvalidatordExchange {
				exchange1 := "Binance"
				exchange2 := "Coinbase"
				return []openapi.TgvalidatordExchange{
					{Exchange: &exchange1},
					{Exchange: &exchange2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExchangesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ExchangesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ExchangesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestExchangeCounterpartyFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordExchangeCounterparty
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns counterparty with zero values",
			dto:  &openapi.TgvalidatordExchangeCounterparty{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordExchangeCounterparty {
				name := "Binance"
				valuation := "1000000.00"
				return &openapi.TgvalidatordExchangeCounterparty{
					Name:                  &name,
					BaseCurrencyValuation: &valuation,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExchangeCounterpartyFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ExchangeCounterpartyFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ExchangeCounterpartyFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.BaseCurrencyValuation != nil && got.BaseCurrencyValuation != *tt.dto.BaseCurrencyValuation {
				t.Errorf("BaseCurrencyValuation = %v, want %v", got.BaseCurrencyValuation, *tt.dto.BaseCurrencyValuation)
			}
		})
	}
}

func TestExchangeCounterpartiesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordExchangeCounterparty
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordExchangeCounterparty{},
			want: 0,
		},
		{
			name: "converts multiple counterparties",
			dtos: func() []openapi.TgvalidatordExchangeCounterparty {
				name1 := "Binance"
				name2 := "Coinbase"
				return []openapi.TgvalidatordExchangeCounterparty{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExchangeCounterpartiesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ExchangeCounterpartiesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ExchangeCounterpartiesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestExchangeFromDTO_NilCurrencyInfo(t *testing.T) {
	exchangeName := "Binance"
	dto := &openapi.TgvalidatordExchange{
		Exchange:     &exchangeName,
		CurrencyInfo: nil,
	}

	got := ExchangeFromDTO(dto)
	if got == nil {
		t.Fatal("ExchangeFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo != nil {
		t.Errorf("CurrencyInfo should be nil when DTO currency info is nil, got %v", got.CurrencyInfo)
	}
	if got.Exchange != "Binance" {
		t.Errorf("Exchange = %v, want Binance", got.Exchange)
	}
}

func TestExchangeFromDTO_NilDates(t *testing.T) {
	exchangeName := "Binance"
	dto := &openapi.TgvalidatordExchange{
		Exchange:     &exchangeName,
		CreationDate: nil,
		UpdateDate:   nil,
	}

	got := ExchangeFromDTO(dto)
	if got == nil {
		t.Fatal("ExchangeFromDTO() returned nil for non-nil input")
	}
	// When dates are nil, they should be the zero time value
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate should be zero time when nil, got %v", got.CreationDate)
	}
	if !got.UpdateDate.IsZero() {
		t.Errorf("UpdateDate should be zero time when nil, got %v", got.UpdateDate)
	}
}

func TestExchangeFromDTO_HasWLAField(t *testing.T) {
	tests := []struct {
		name       string
		hasWLA     *bool
		wantHasWLA bool
	}{
		{
			name:       "nil hasWLA defaults to false",
			hasWLA:     nil,
			wantHasWLA: false,
		},
		{
			name:       "true hasWLA",
			hasWLA:     boolPtr(true),
			wantHasWLA: true,
		},
		{
			name:       "false hasWLA",
			hasWLA:     boolPtr(false),
			wantHasWLA: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordExchange{
				HasWLA: tt.hasWLA,
			}
			got := ExchangeFromDTO(dto)
			if got.HasWLA != tt.wantHasWLA {
				t.Errorf("HasWLA = %v, want %v", got.HasWLA, tt.wantHasWLA)
			}
		})
	}
}
