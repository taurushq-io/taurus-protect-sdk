package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestFeeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordKeyValue
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns fee with empty values",
			dto:  &openapi.TgvalidatordKeyValue{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordKeyValue {
				key := "ETH:mainnet"
				value := "21000"
				return &openapi.TgvalidatordKeyValue{
					Key:   &key,
					Value: &value,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeeFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FeeFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FeeFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Key != nil && got.Key != *tt.dto.Key {
				t.Errorf("Key = %v, want %v", got.Key, *tt.dto.Key)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
		})
	}
}

func TestFeesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordKeyValue
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordKeyValue{},
			want: 0,
		},
		{
			name: "converts multiple fees",
			dtos: func() []openapi.TgvalidatordKeyValue {
				key1 := "ETH:mainnet"
				value1 := "21000"
				key2 := "BTC:mainnet"
				value2 := "1000"
				return []openapi.TgvalidatordKeyValue{
					{Key: &key1, Value: &value1},
					{Key: &key2, Value: &value2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("FeesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("FeesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFeeV2FromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordFee
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns fee with empty values",
			dto:  &openapi.TgvalidatordFee{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordFee {
				currencyID := "ETH"
				value := "21000000000000"
				denom := "wei"
				id := "currency-123"
				name := "Ethereum"
				symbol := "ETH"
				blockchain := "ETH"
				network := "mainnet"
				updateDate := time.Now()
				return &openapi.TgvalidatordFee{
					CurrencyId: &currencyID,
					Value:      &value,
					Denom:      &denom,
					CurrencyInfo: &openapi.TgvalidatordCurrency{
						Id:         &id,
						Name:       &name,
						Symbol:     &symbol,
						Blockchain: &blockchain,
						Network:    &network,
					},
					UpdateDate: &updateDate,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeeV2FromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FeeV2FromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FeeV2FromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.CurrencyId != nil && got.CurrencyID != *tt.dto.CurrencyId {
				t.Errorf("CurrencyID = %v, want %v", got.CurrencyID, *tt.dto.CurrencyId)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.Denom != nil && got.Denom != *tt.dto.Denom {
				t.Errorf("Denom = %v, want %v", got.Denom, *tt.dto.Denom)
			}
			if tt.dto.CurrencyInfo != nil {
				if got.CurrencyInfo == nil {
					t.Error("CurrencyInfo should not be nil when DTO has currency info")
				} else if tt.dto.CurrencyInfo.Name != nil && got.CurrencyInfo.Name != *tt.dto.CurrencyInfo.Name {
					t.Errorf("CurrencyInfo.Name = %v, want %v", got.CurrencyInfo.Name, *tt.dto.CurrencyInfo.Name)
				}
			}
			if tt.dto.UpdateDate != nil && !got.UpdateDate.Equal(*tt.dto.UpdateDate) {
				t.Errorf("UpdateDate = %v, want %v", got.UpdateDate, *tt.dto.UpdateDate)
			}
		})
	}
}

func TestFeesV2FromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordFee
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordFee{},
			want: 0,
		},
		{
			name: "converts multiple fees",
			dtos: func() []openapi.TgvalidatordFee {
				currencyID1 := "ETH"
				value1 := "21000000000000"
				currencyID2 := "BTC"
				value2 := "1000"
				return []openapi.TgvalidatordFee{
					{CurrencyId: &currencyID1, Value: &value1},
					{CurrencyId: &currencyID2, Value: &value2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeesV2FromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("FeesV2FromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("FeesV2FromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFeeV2FromDTO_NilUpdateDate(t *testing.T) {
	currencyID := "ETH"
	value := "21000000000000"

	dto := &openapi.TgvalidatordFee{
		CurrencyId: &currencyID,
		Value:      &value,
		UpdateDate: nil,
	}

	got := FeeV2FromDTO(dto)
	if got == nil {
		t.Fatal("FeeV2FromDTO() returned nil for non-nil input")
	}

	// UpdateDate should be zero time when nil
	if !got.UpdateDate.IsZero() {
		t.Errorf("UpdateDate should be zero time when nil, got %v", got.UpdateDate)
	}
}

func TestFeeV2FromDTO_NilCurrencyInfo(t *testing.T) {
	currencyID := "ETH"
	value := "21000000000000"

	dto := &openapi.TgvalidatordFee{
		CurrencyId:   &currencyID,
		Value:        &value,
		CurrencyInfo: nil,
	}

	got := FeeV2FromDTO(dto)
	if got == nil {
		t.Fatal("FeeV2FromDTO() returned nil for non-nil input")
	}

	if got.CurrencyInfo != nil {
		t.Error("CurrencyInfo should be nil when DTO has nil currency info")
	}
}
