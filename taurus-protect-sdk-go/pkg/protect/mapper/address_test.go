package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestAddressFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAddress
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns address with zero values",
			dto:  &openapi.TgvalidatordAddress{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAddress {
				id := "addr-123"
				walletId := "wallet-456"
				address := "0x1234567890abcdef"
				alternateAddress := "0xabcdef1234567890"
				label := "Primary Address"
				comment := "Main deposit address"
				currency := "ETH"
				customerId := "customer-789"
				externalAddressId := "ext-addr-001"
				addressPath := "m/44'/60'/0'/0/0"
				addressIndex := "0"
				nonce := "42"
				status := "created"
				disabled := false
				canUseAllFunds := true
				now := time.Now()
				linkedIds := []string{"wl-1", "wl-2"}
				return &openapi.TgvalidatordAddress{
					Id:                          &id,
					WalletId:                    &walletId,
					Address:                     &address,
					AlternateAddress:            &alternateAddress,
					Label:                       &label,
					Comment:                     &comment,
					Currency:                    &currency,
					CustomerId:                  &customerId,
					ExternalAddressId:           &externalAddressId,
					AddressPath:                 &addressPath,
					AddressIndex:                &addressIndex,
					Nonce:                       &nonce,
					Status:                      &status,
					Disabled:                    &disabled,
					CanUseAllFunds:              &canUseAllFunds,
					CreationDate:                &now,
					UpdateDate:                  &now,
					LinkedWhitelistedAddressIds: linkedIds,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddressFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("AddressFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("AddressFromDTO() returned nil for non-nil input")
			}
			// Verify key fields
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Address != nil && got.Address != *tt.dto.Address {
				t.Errorf("Address = %v, want %v", got.Address, *tt.dto.Address)
			}
			if tt.dto.LinkedWhitelistedAddressIds != nil {
				if len(got.LinkedWhitelistedAddressIDs) != len(tt.dto.LinkedWhitelistedAddressIds) {
					t.Errorf("LinkedWhitelistedAddressIDs length = %v, want %v",
						len(got.LinkedWhitelistedAddressIDs), len(tt.dto.LinkedWhitelistedAddressIds))
				}
			}
		})
	}
}

func TestAddressFromDTO_WithBalance(t *testing.T) {
	totalConfirmed := "1000000000000000000"
	dto := &openapi.TgvalidatordAddress{
		Balance: &openapi.TgvalidatordBalance{
			TotalConfirmed: &totalConfirmed,
		},
	}

	got := AddressFromDTO(dto)
	if got.Balance == nil {
		t.Fatal("Balance should not be nil")
	}
	if got.Balance.TotalConfirmed != totalConfirmed {
		t.Errorf("Balance.TotalConfirmed = %v, want %v", got.Balance.TotalConfirmed, totalConfirmed)
	}
}

func TestAddressFromDTO_WithAttributes(t *testing.T) {
	attrId := "attr-1"
	attrKey := "tag"
	attrValue := "important"
	dto := &openapi.TgvalidatordAddress{
		Attributes: []openapi.TgvalidatordAddressAttribute{
			{Id: &attrId, Key: &attrKey, Value: &attrValue},
		},
	}

	got := AddressFromDTO(dto)
	if len(got.Attributes) != 1 {
		t.Fatalf("Attributes length = %v, want 1", len(got.Attributes))
	}
	if got.Attributes[0].ID != attrId {
		t.Errorf("Attribute ID = %v, want %v", got.Attributes[0].ID, attrId)
	}
	if got.Attributes[0].Key != attrKey {
		t.Errorf("Attribute Key = %v, want %v", got.Attributes[0].Key, attrKey)
	}
}

func TestAddressesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordAddress
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordAddress{},
			want: 0,
		},
		{
			name: "converts multiple addresses",
			dtos: func() []openapi.TgvalidatordAddress {
				id1 := "addr-1"
				id2 := "addr-2"
				id3 := "addr-3"
				return []openapi.TgvalidatordAddress{
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
			got := AddressesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("AddressesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("AddressesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestAddressAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAddressAttribute
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAddressAttribute {
				id := "attr-123"
				key := "category"
				value := "deposit"
				return &openapi.TgvalidatordAddressAttribute{
					Id:    &id,
					Key:   &key,
					Value: &value,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddressAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" || got.Key != "" || got.Value != "" {
					t.Errorf("AddressAttributeFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
		})
	}
}

func TestAddressFromDTO_TimestampConversion(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	dto := &openapi.TgvalidatordAddress{
		CreationDate: &now,
		UpdateDate:   &now,
	}

	got := AddressFromDTO(dto)
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
	}
}

func TestAddressFromDTO_LinkedWhitelistedAddressIDsCopy(t *testing.T) {
	// Verify that the slice is copied, not referenced
	linkedIds := []string{"wl-1", "wl-2", "wl-3"}
	dto := &openapi.TgvalidatordAddress{
		LinkedWhitelistedAddressIds: linkedIds,
	}

	got := AddressFromDTO(dto)

	// Modify the original slice
	linkedIds[0] = "modified"

	// The converted address should still have the original value
	if got.LinkedWhitelistedAddressIDs[0] != "wl-1" {
		t.Errorf("LinkedWhitelistedAddressIDs was not properly copied, got %v", got.LinkedWhitelistedAddressIDs[0])
	}
}
