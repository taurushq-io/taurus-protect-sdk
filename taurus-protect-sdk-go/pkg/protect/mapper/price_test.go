package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestPriceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrencyPrice
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns price with zero values",
			dto:  &openapi.TgvalidatordCurrencyPrice{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordCurrencyPrice {
				blockchain := "ETH"
				currencyFrom := "ETH"
				currencyTo := "USD"
				decimals := "8"
				rate := "3000.50"
				change := "2.5"
				source := "cryptocompare"
				now := time.Now()
				id := "currency-123"
				name := "Ethereum"
				symbol := "ETH"
				return &openapi.TgvalidatordCurrencyPrice{
					Blockchain:          &blockchain,
					CurrencyFrom:        &currencyFrom,
					CurrencyTo:          &currencyTo,
					Decimals:            &decimals,
					Rate:                &rate,
					ChangePercent24Hour: &change,
					Source:              &source,
					CreationDate:        &now,
					UpdateDate:          &now,
					Signatures: []openapi.TgvalidatordCurrencyPriceSignature{
						{UserId: "user-1", Signature: "sig-1"},
					},
					CurrencyFromInfo: &openapi.TgvalidatordCurrency{
						Id:     &id,
						Name:   &name,
						Symbol: &symbol,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriceFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("PriceFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("PriceFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.CurrencyFrom != nil && got.CurrencyFrom != *tt.dto.CurrencyFrom {
				t.Errorf("CurrencyFrom = %v, want %v", got.CurrencyFrom, *tt.dto.CurrencyFrom)
			}
			if tt.dto.CurrencyTo != nil && got.CurrencyTo != *tt.dto.CurrencyTo {
				t.Errorf("CurrencyTo = %v, want %v", got.CurrencyTo, *tt.dto.CurrencyTo)
			}
			if tt.dto.Rate != nil && got.Rate != *tt.dto.Rate {
				t.Errorf("Rate = %v, want %v", got.Rate, *tt.dto.Rate)
			}
			if tt.dto.Signatures != nil && len(got.Signatures) != len(tt.dto.Signatures) {
				t.Errorf("Signatures length = %v, want %v", len(got.Signatures), len(tt.dto.Signatures))
			}
		})
	}
}

func TestPricesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordCurrencyPrice
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordCurrencyPrice{},
			want: 0,
		},
		{
			name: "converts multiple prices",
			dtos: func() []openapi.TgvalidatordCurrencyPrice {
				currencyFrom := "BTC"
				currencyTo := "USD"
				return []openapi.TgvalidatordCurrencyPrice{
					{CurrencyFrom: &currencyFrom, CurrencyTo: &currencyTo},
					{CurrencyFrom: &currencyFrom, CurrencyTo: &currencyTo},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PricesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("PricesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("PricesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestPriceSignatureFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrencyPriceSignature
	}{
		{
			name: "nil input returns empty signature",
			dto:  nil,
		},
		{
			name: "maps all fields",
			dto: &openapi.TgvalidatordCurrencyPriceSignature{
				UserId:    "user-123",
				Signature: "sig-abc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriceSignatureFromDTO(tt.dto)
			if tt.dto == nil {
				if got.UserID != "" || got.Signature != "" {
					t.Errorf("PriceSignatureFromDTO(nil) = %v, want empty", got)
				}
				return
			}
			if got.UserID != tt.dto.UserId {
				t.Errorf("UserID = %v, want %v", got.UserID, tt.dto.UserId)
			}
			if got.Signature != tt.dto.Signature {
				t.Errorf("Signature = %v, want %v", got.Signature, tt.dto.Signature)
			}
		})
	}
}

func TestPriceSignaturesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordCurrencyPriceSignature
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordCurrencyPriceSignature{},
			want: 0,
		},
		{
			name: "converts multiple signatures",
			dtos: []openapi.TgvalidatordCurrencyPriceSignature{
				{UserId: "user-1", Signature: "sig-1"},
				{UserId: "user-2", Signature: "sig-2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriceSignaturesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("PriceSignaturesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("PriceSignaturesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestPriceHistoryPointFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordPricesHistoryPoint
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns point with zero values",
			dto:  &openapi.TgvalidatordPricesHistoryPoint{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordPricesHistoryPoint {
				blockchain := "ETH"
				currencyFrom := "ETH"
				currencyTo := "USD"
				high := "3100.00"
				low := "2900.00"
				open := "2950.00"
				closeVal := "3050.00"
				volumeFrom := "1000"
				volumeTo := "3050000"
				changePercent := "3.4"
				now := time.Now()
				return &openapi.TgvalidatordPricesHistoryPoint{
					PeriodStartDate: &now,
					Blockchain:      &blockchain,
					CurrencyFrom:    &currencyFrom,
					CurrencyTo:      &currencyTo,
					High:            &high,
					Low:             &low,
					Open:            &open,
					Close:           &closeVal,
					VolumeFrom:      &volumeFrom,
					VolumeTo:        &volumeTo,
					ChangePercent:   &changePercent,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriceHistoryPointFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("PriceHistoryPointFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("PriceHistoryPointFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.High != nil && got.High != *tt.dto.High {
				t.Errorf("High = %v, want %v", got.High, *tt.dto.High)
			}
			if tt.dto.Low != nil && got.Low != *tt.dto.Low {
				t.Errorf("Low = %v, want %v", got.Low, *tt.dto.Low)
			}
			if tt.dto.Open != nil && got.Open != *tt.dto.Open {
				t.Errorf("Open = %v, want %v", got.Open, *tt.dto.Open)
			}
			if tt.dto.Close != nil && got.Close != *tt.dto.Close {
				t.Errorf("Close = %v, want %v", got.Close, *tt.dto.Close)
			}
		})
	}
}

func TestPriceHistoryPointsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordPricesHistoryPoint
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordPricesHistoryPoint{},
			want: 0,
		},
		{
			name: "converts multiple points",
			dtos: func() []openapi.TgvalidatordPricesHistoryPoint {
				high := "100"
				return []openapi.TgvalidatordPricesHistoryPoint{
					{High: &high},
					{High: &high},
					{High: &high},
				}
			}(),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PriceHistoryPointsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("PriceHistoryPointsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("PriceHistoryPointsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestConversionValueFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordConversionValue
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns value with zero values",
			dto:  &openapi.TgvalidatordConversionValue{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordConversionValue {
				symbol := "USD"
				value := "1500000000000000000"
				mainUnitValue := "1.5"
				id := "currency-123"
				name := "US Dollar"
				return &openapi.TgvalidatordConversionValue{
					Symbol:        &symbol,
					Value:         &value,
					MainUnitValue: &mainUnitValue,
					CurrencyInfo: &openapi.TgvalidatordCurrency{
						Id:   &id,
						Name: &name,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConversionValueFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ConversionValueFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ConversionValueFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Symbol != nil && got.Symbol != *tt.dto.Symbol {
				t.Errorf("Symbol = %v, want %v", got.Symbol, *tt.dto.Symbol)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.MainUnitValue != nil && got.MainUnitValue != *tt.dto.MainUnitValue {
				t.Errorf("MainUnitValue = %v, want %v", got.MainUnitValue, *tt.dto.MainUnitValue)
			}
		})
	}
}

func TestConversionValuesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordConversionValue
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordConversionValue{},
			want: 0,
		},
		{
			name: "converts multiple values",
			dtos: func() []openapi.TgvalidatordConversionValue {
				symbol := "USD"
				return []openapi.TgvalidatordConversionValue{
					{Symbol: &symbol},
					{Symbol: &symbol},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConversionValuesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ConversionValuesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ConversionValuesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestConversionResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordConversionReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with zero values",
			dto:  &openapi.TgvalidatordConversionReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordConversionReply {
				currencyFrom := "ETH"
				baseCurrency := "USD"
				symbol := "USD"
				value := "3000.00"
				id := "currency-123"
				name := "Ethereum"
				return &openapi.TgvalidatordConversionReply{
					CurrencyFrom: &currencyFrom,
					BaseCurrency: &baseCurrency,
					Result: []openapi.TgvalidatordConversionValue{
						{Symbol: &symbol, Value: &value},
					},
					FullCurrencyFrom: &openapi.TgvalidatordCurrency{
						Id:   &id,
						Name: &name,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConversionResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ConversionResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ConversionResultFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.CurrencyFrom != nil && got.CurrencyFrom != *tt.dto.CurrencyFrom {
				t.Errorf("CurrencyFrom = %v, want %v", got.CurrencyFrom, *tt.dto.CurrencyFrom)
			}
			if tt.dto.BaseCurrency != nil && got.BaseCurrency != *tt.dto.BaseCurrency {
				t.Errorf("BaseCurrency = %v, want %v", got.BaseCurrency, *tt.dto.BaseCurrency)
			}
			if tt.dto.Result != nil && len(got.Values) != len(tt.dto.Result) {
				t.Errorf("Values length = %v, want %v", len(got.Values), len(tt.dto.Result))
			}
		})
	}
}
