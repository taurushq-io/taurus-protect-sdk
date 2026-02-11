package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestTransactionFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTransaction
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns transaction with zero values",
			dto:  &openapi.TgvalidatordTransaction{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTransaction {
				id := "tx-123"
				uniqueId := "unique-456"
				transactionId := "txid-789"
				hash := "0xabc123def456"
				direction := "outgoing"
				txType := "transfer"
				currency := "ETH"
				blockchain := "ETH"
				amount := "1000000000000000000"
				amountMainUnit := "1.0"
				fee := "21000000000000"
				feeMainUnit := "0.000021"
				block := "12345678"
				now := time.Now()
				return &openapi.TgvalidatordTransaction{
					Id:               &id,
					UniqueId:         &uniqueId,
					TransactionId:    &transactionId,
					Hash:             &hash,
					Direction:        &direction,
					Type:             &txType,
					Currency:         &currency,
					Blockchain:       &blockchain,
					Amount:           &amount,
					AmountMainUnit:   &amountMainUnit,
					Fee:              &fee,
					FeeMainUnit:      &feeMainUnit,
					Block:            &block,
					ReceptionDate:    &now,
					ConfirmationDate: &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransactionFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TransactionFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TransactionFromDTO() returned nil for non-nil input")
			}
			// Verify key fields
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Hash != nil && got.Hash != *tt.dto.Hash {
				t.Errorf("Hash = %v, want %v", got.Hash, *tt.dto.Hash)
			}
			if tt.dto.Amount != nil && got.Amount != *tt.dto.Amount {
				t.Errorf("Amount = %v, want %v", got.Amount, *tt.dto.Amount)
			}
		})
	}
}

func TestTransactionFromDTO_WithSources(t *testing.T) {
	addr := "0x1234567890abcdef"
	label := "Source Address"
	amount := "500000000000000000"
	dto := &openapi.TgvalidatordTransaction{
		Sources: []openapi.TgvalidatordAddressInfo{
			{Address: &addr, Label: &label, Amount: &amount},
		},
	}

	got := TransactionFromDTO(dto)
	if len(got.Sources) != 1 {
		t.Fatalf("Sources length = %v, want 1", len(got.Sources))
	}
	if got.Sources[0].Address != addr {
		t.Errorf("Source Address = %v, want %v", got.Sources[0].Address, addr)
	}
	if got.Sources[0].Label != label {
		t.Errorf("Source Label = %v, want %v", got.Sources[0].Label, label)
	}
	if got.Sources[0].Amount != amount {
		t.Errorf("Source Amount = %v, want %v", got.Sources[0].Amount, amount)
	}
}

func TestTransactionFromDTO_WithDestinations(t *testing.T) {
	addr := "0xabcdef1234567890"
	label := "Destination Address"
	amount := "500000000000000000"
	internalAddrId := "internal-123"
	whitelistedAddrId := "whitelist-456"
	dto := &openapi.TgvalidatordTransaction{
		Destinations: []openapi.TgvalidatordAddressInfo{
			{
				Address:              &addr,
				Label:                &label,
				Amount:               &amount,
				InternalAddressId:    &internalAddrId,
				WhitelistedAddressId: &whitelistedAddrId,
			},
		},
	}

	got := TransactionFromDTO(dto)
	if len(got.Destinations) != 1 {
		t.Fatalf("Destinations length = %v, want 1", len(got.Destinations))
	}
	dest := got.Destinations[0]
	if dest.Address != addr {
		t.Errorf("Destination Address = %v, want %v", dest.Address, addr)
	}
	if dest.InternalAddressID != internalAddrId {
		t.Errorf("InternalAddressID = %v, want %v", dest.InternalAddressID, internalAddrId)
	}
	if dest.WhitelistedAddressID != whitelistedAddrId {
		t.Errorf("WhitelistedAddressID = %v, want %v", dest.WhitelistedAddressID, whitelistedAddrId)
	}
}

func TestTransactionFromDTO_MultipleSourcesAndDestinations(t *testing.T) {
	addr1 := "0x1111"
	addr2 := "0x2222"
	addr3 := "0x3333"
	dto := &openapi.TgvalidatordTransaction{
		Sources: []openapi.TgvalidatordAddressInfo{
			{Address: &addr1},
			{Address: &addr2},
		},
		Destinations: []openapi.TgvalidatordAddressInfo{
			{Address: &addr3},
		},
	}

	got := TransactionFromDTO(dto)
	if len(got.Sources) != 2 {
		t.Errorf("Sources length = %v, want 2", len(got.Sources))
	}
	if len(got.Destinations) != 1 {
		t.Errorf("Destinations length = %v, want 1", len(got.Destinations))
	}
}

func TestTransactionsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTransaction
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTransaction{},
			want: 0,
		},
		{
			name: "converts multiple transactions",
			dtos: func() []openapi.TgvalidatordTransaction {
				id1 := "tx-1"
				id2 := "tx-2"
				id3 := "tx-3"
				return []openapi.TgvalidatordTransaction{
					{Id: &id1},
					{Id: &id2},
					{Id: &id3},
				}
			}(),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransactionsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("TransactionsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("TransactionsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestAddressInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAddressInfo
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAddressInfo {
				addr := "0x1234567890abcdef"
				label := "Test Address"
				customerId := "customer-123"
				amount := "1000000000000000000"
				amountMainUnit := "1.0"
				addrType := "destination"
				internalAddrId := "internal-456"
				whitelistedAddrId := "whitelist-789"
				return &openapi.TgvalidatordAddressInfo{
					Address:              &addr,
					Label:                &label,
					CustomerId:           &customerId,
					Amount:               &amount,
					AmountMainUnit:       &amountMainUnit,
					Type:                 &addrType,
					InternalAddressId:    &internalAddrId,
					WhitelistedAddressId: &whitelistedAddrId,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddressInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got.Address != "" || got.Label != "" {
					t.Errorf("AddressInfoFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Address != nil && got.Address != *tt.dto.Address {
				t.Errorf("Address = %v, want %v", got.Address, *tt.dto.Address)
			}
			if tt.dto.Label != nil && got.Label != *tt.dto.Label {
				t.Errorf("Label = %v, want %v", got.Label, *tt.dto.Label)
			}
			if tt.dto.Amount != nil && got.Amount != *tt.dto.Amount {
				t.Errorf("Amount = %v, want %v", got.Amount, *tt.dto.Amount)
			}
		})
	}
}

func TestTransactionFromDTO_TimestampConversion(t *testing.T) {
	reception := time.Now().Add(-time.Hour).Truncate(time.Second)
	confirmation := time.Now().Truncate(time.Second)
	dto := &openapi.TgvalidatordTransaction{
		ReceptionDate:    &reception,
		ConfirmationDate: &confirmation,
	}

	got := TransactionFromDTO(dto)
	if !got.ReceptionDate.Equal(reception) {
		t.Errorf("ReceptionDate = %v, want %v", got.ReceptionDate, reception)
	}
	if !got.ConfirmationDate.Equal(confirmation) {
		t.Errorf("ConfirmationDate = %v, want %v", got.ConfirmationDate, confirmation)
	}
}
