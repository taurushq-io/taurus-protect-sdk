package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestReservationFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordReservation
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns reservation with zero values",
			dto:  &openapi.TgvalidatordReservation{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordReservation {
				id := "res-123"
				kind := "PENDING_REQUEST"
				comment := "Test reservation"
				addressId := "addr-456"
				address := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
				amount := "1000000"
				creationDate := time.Now()
				resourceId := "req-789"
				resourceType := "REQUEST"
				return &openapi.TgvalidatordReservation{
					Id:           &id,
					Kind:         &kind,
					Comment:      &comment,
					Addressid:    &addressId,
					Address:      &address,
					Amount:       &amount,
					CreationDate: &creationDate,
					ResourceId:   &resourceId,
					ResourceType: &resourceType,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReservationFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ReservationFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ReservationFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Kind != nil && got.Kind != *tt.dto.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Kind)
			}
			if tt.dto.Comment != nil && got.Comment != *tt.dto.Comment {
				t.Errorf("Comment = %v, want %v", got.Comment, *tt.dto.Comment)
			}
			if tt.dto.Addressid != nil && got.AddressID != *tt.dto.Addressid {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.Addressid)
			}
			if tt.dto.Address != nil && got.Address != *tt.dto.Address {
				t.Errorf("Address = %v, want %v", got.Address, *tt.dto.Address)
			}
			if tt.dto.Amount != nil && got.Amount != *tt.dto.Amount {
				t.Errorf("Amount = %v, want %v", got.Amount, *tt.dto.Amount)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
			if tt.dto.ResourceId != nil && got.ResourceID != *tt.dto.ResourceId {
				t.Errorf("ResourceID = %v, want %v", got.ResourceID, *tt.dto.ResourceId)
			}
			if tt.dto.ResourceType != nil && got.ResourceType != *tt.dto.ResourceType {
				t.Errorf("ResourceType = %v, want %v", got.ResourceType, *tt.dto.ResourceType)
			}
		})
	}
}

func TestReservationsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordReservation
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordReservation{},
			want: 0,
		},
		{
			name: "converts multiple reservations",
			dtos: func() []openapi.TgvalidatordReservation {
				kind1 := "PENDING_REQUEST"
				kind2 := "MINIMUM_BALANCE"
				return []openapi.TgvalidatordReservation{
					{Kind: &kind1},
					{Kind: &kind2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReservationsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ReservationsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ReservationsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestReservationUTXOFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordUTXO
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns UTXO with zero values",
			dto:  &openapi.TgvalidatordUTXO{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordUTXO {
				id := "utxo-123"
				hash := "abc123def456"
				outputIndex := int64(0)
				script := "76a914..."
				value := "100000"
				valueString := "0.001 BTC"
				blockHeight := "800000"
				reservationId := "res-456"
				return &openapi.TgvalidatordUTXO{
					Id:            &id,
					Hash:          &hash,
					OutputIndex:   &outputIndex,
					Script:        &script,
					Value:         &value,
					ValueString:   &valueString,
					BlockHeight:   &blockHeight,
					ReservationId: &reservationId,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReservationUTXOFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ReservationUTXOFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ReservationUTXOFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Hash != nil && got.Hash != *tt.dto.Hash {
				t.Errorf("Hash = %v, want %v", got.Hash, *tt.dto.Hash)
			}
			if tt.dto.OutputIndex != nil && got.OutputIndex != *tt.dto.OutputIndex {
				t.Errorf("OutputIndex = %v, want %v", got.OutputIndex, *tt.dto.OutputIndex)
			}
			if tt.dto.Script != nil && got.Script != *tt.dto.Script {
				t.Errorf("Script = %v, want %v", got.Script, *tt.dto.Script)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.ValueString != nil && got.ValueString != *tt.dto.ValueString {
				t.Errorf("ValueString = %v, want %v", got.ValueString, *tt.dto.ValueString)
			}
			if tt.dto.BlockHeight != nil && got.BlockHeight != *tt.dto.BlockHeight {
				t.Errorf("BlockHeight = %v, want %v", got.BlockHeight, *tt.dto.BlockHeight)
			}
			if tt.dto.ReservationId != nil && got.ReservationID != *tt.dto.ReservationId {
				t.Errorf("ReservationID = %v, want %v", got.ReservationID, *tt.dto.ReservationId)
			}
		})
	}
}

func TestReservationFromDTO_WithCurrencyInfo(t *testing.T) {
	id := "res-123"
	currencyID := "currency-456"
	currencyName := "BTC"
	dto := &openapi.TgvalidatordReservation{
		Id: &id,
		CurrencyInfo: &openapi.TgvalidatordCurrency{
			Id:   &currencyID,
			Name: &currencyName,
		},
	}

	got := ReservationFromDTO(dto)
	if got == nil {
		t.Fatal("ReservationFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo == nil {
		t.Fatal("CurrencyInfo should not be nil")
	}
	if got.CurrencyInfo.ID != currencyID {
		t.Errorf("CurrencyInfo.ID = %v, want %v", got.CurrencyInfo.ID, currencyID)
	}
	if got.CurrencyInfo.Name != currencyName {
		t.Errorf("CurrencyInfo.Name = %v, want %v", got.CurrencyInfo.Name, currencyName)
	}
}

func TestReservationFromDTO_WithAsset(t *testing.T) {
	id := "res-123"
	currency := "ETH"
	dto := &openapi.TgvalidatordReservation{
		Id: &id,
		Asset: &openapi.TgvalidatordAsset{
			Currency: currency,
		},
	}

	got := ReservationFromDTO(dto)
	if got == nil {
		t.Fatal("ReservationFromDTO() returned nil for non-nil input")
	}
	if got.Asset == nil {
		t.Fatal("Asset should not be nil")
	}
	if got.Asset.Currency != currency {
		t.Errorf("Asset.Currency = %v, want %v", got.Asset.Currency, currency)
	}
}

func TestReservationFromDTO_NilCreationDate(t *testing.T) {
	id := "res-123"
	dto := &openapi.TgvalidatordReservation{
		Id:           &id,
		CreationDate: nil,
	}

	got := ReservationFromDTO(dto)
	if got == nil {
		t.Fatal("ReservationFromDTO() returned nil for non-nil input")
	}
	// When creation date is nil, it should be the zero time value
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate should be zero time when nil, got %v", got.CreationDate)
	}
}

func TestReservationUTXOFromDTO_NilOutputIndex(t *testing.T) {
	id := "utxo-123"
	dto := &openapi.TgvalidatordUTXO{
		Id:          &id,
		OutputIndex: nil,
	}

	got := ReservationUTXOFromDTO(dto)
	if got == nil {
		t.Fatal("ReservationUTXOFromDTO() returned nil for non-nil input")
	}
	// When output index is nil, it should be 0
	if got.OutputIndex != 0 {
		t.Errorf("OutputIndex should be 0 when nil, got %v", got.OutputIndex)
	}
}
