package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestBusinessRuleFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordBusinessRule
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns business rule with zero values",
			dto:  &openapi.TgvalidatordBusinessRule{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordBusinessRule {
				id := "rule-123"
				tenantId := "tenant-456"
				currency := "ETH"
				walletId := "wallet-789"
				ruleKey := "TRANSACTIONS_ENABLED"
				ruleValue := "true"
				ruleGroup := "security"
				ruleDescription := "Enable transactions"
				ruleValidation := "^(true|false)$"
				addressId := "address-abc"
				entityType := "global"
				entityID := ""
				return &openapi.TgvalidatordBusinessRule{
					Id:              &id,
					TenantId:        &tenantId,
					Currency:        &currency,
					WalletId:        &walletId,
					RuleKey:         &ruleKey,
					RuleValue:       &ruleValue,
					RuleGroup:       &ruleGroup,
					RuleDescription: &ruleDescription,
					RuleValidation:  &ruleValidation,
					AddressId:       &addressId,
					EntityType:      &entityType,
					EntityID:        &entityID,
				}
			}(),
		},
		{
			name: "DTO with currency info",
			dto: func() *openapi.TgvalidatordBusinessRule {
				id := "rule-with-currency"
				ruleKey := "MAX_AMOUNT"
				ruleValue := "1000"
				currencyId := "currency-123"
				currencyName := "ethereum"
				currencySymbol := "ETH"
				return &openapi.TgvalidatordBusinessRule{
					Id:        &id,
					RuleKey:   &ruleKey,
					RuleValue: &ruleValue,
					CurrencyInfo: &openapi.TgvalidatordCurrency{
						Id:     &currencyId,
						Name:   &currencyName,
						Symbol: &currencySymbol,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BusinessRuleFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("BusinessRuleFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("BusinessRuleFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Currency != nil && got.Currency != *tt.dto.Currency {
				t.Errorf("Currency = %v, want %v", got.Currency, *tt.dto.Currency)
			}
			if tt.dto.WalletId != nil && got.WalletID != *tt.dto.WalletId {
				t.Errorf("WalletID = %v, want %v", got.WalletID, *tt.dto.WalletId)
			}
			if tt.dto.RuleKey != nil && got.RuleKey != *tt.dto.RuleKey {
				t.Errorf("RuleKey = %v, want %v", got.RuleKey, *tt.dto.RuleKey)
			}
			if tt.dto.RuleValue != nil && got.RuleValue != *tt.dto.RuleValue {
				t.Errorf("RuleValue = %v, want %v", got.RuleValue, *tt.dto.RuleValue)
			}
			if tt.dto.RuleGroup != nil && got.RuleGroup != *tt.dto.RuleGroup {
				t.Errorf("RuleGroup = %v, want %v", got.RuleGroup, *tt.dto.RuleGroup)
			}
			if tt.dto.RuleDescription != nil && got.RuleDescription != *tt.dto.RuleDescription {
				t.Errorf("RuleDescription = %v, want %v", got.RuleDescription, *tt.dto.RuleDescription)
			}
			if tt.dto.RuleValidation != nil && got.RuleValidation != *tt.dto.RuleValidation {
				t.Errorf("RuleValidation = %v, want %v", got.RuleValidation, *tt.dto.RuleValidation)
			}
			if tt.dto.AddressId != nil && got.AddressID != *tt.dto.AddressId {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressId)
			}
			if tt.dto.EntityType != nil && got.EntityType != *tt.dto.EntityType {
				t.Errorf("EntityType = %v, want %v", got.EntityType, *tt.dto.EntityType)
			}
			if tt.dto.EntityID != nil && got.EntityID != *tt.dto.EntityID {
				t.Errorf("EntityID = %v, want %v", got.EntityID, *tt.dto.EntityID)
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

func TestBusinessRulesFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordBusinessRule
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordBusinessRule{},
			want: 0,
		},
		{
			name: "converts multiple business rules",
			dtos: func() []openapi.TgvalidatordBusinessRule {
				key1 := "TRANSACTIONS_ENABLED"
				key2 := "MAX_AMOUNT"
				return []openapi.TgvalidatordBusinessRule{
					{RuleKey: &key1},
					{RuleKey: &key2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BusinessRulesFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("BusinessRulesFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("BusinessRulesFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestBusinessRuleFromDTO_NilCurrencyInfo(t *testing.T) {
	ruleKey := "TRANSACTIONS_ENABLED"
	dto := &openapi.TgvalidatordBusinessRule{
		RuleKey:      &ruleKey,
		CurrencyInfo: nil,
	}

	got := BusinessRuleFromDTO(dto)
	if got == nil {
		t.Fatal("BusinessRuleFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo != nil {
		t.Errorf("CurrencyInfo should be nil when DTO currency info is nil, got %v", got.CurrencyInfo)
	}
	if got.RuleKey != "TRANSACTIONS_ENABLED" {
		t.Errorf("RuleKey = %v, want TRANSACTIONS_ENABLED", got.RuleKey)
	}
}

func TestBusinessRuleFromDTO_WithCurrencyInfo(t *testing.T) {
	ruleKey := "MAX_AMOUNT"
	currencyId := "currency-123"
	currencyName := "ethereum"
	currencySymbol := "ETH"
	dto := &openapi.TgvalidatordBusinessRule{
		RuleKey: &ruleKey,
		CurrencyInfo: &openapi.TgvalidatordCurrency{
			Id:     &currencyId,
			Name:   &currencyName,
			Symbol: &currencySymbol,
		},
	}

	got := BusinessRuleFromDTO(dto)
	if got == nil {
		t.Fatal("BusinessRuleFromDTO() returned nil for non-nil input")
	}
	if got.CurrencyInfo == nil {
		t.Fatal("CurrencyInfo should not be nil when DTO has currency info")
	}
	if got.CurrencyInfo.ID != currencyId {
		t.Errorf("CurrencyInfo.ID = %v, want %v", got.CurrencyInfo.ID, currencyId)
	}
	if got.CurrencyInfo.Name != currencyName {
		t.Errorf("CurrencyInfo.Name = %v, want %v", got.CurrencyInfo.Name, currencyName)
	}
	if got.CurrencyInfo.Symbol != currencySymbol {
		t.Errorf("CurrencyInfo.Symbol = %v, want %v", got.CurrencyInfo.Symbol, currencySymbol)
	}
}

func TestBusinessRuleFromDTO_EntityTypeValues(t *testing.T) {
	entityTypes := []string{"global", "currency", "wallet", "address", "exchange", "exchange_account", "tn_participant"}

	for _, entityType := range entityTypes {
		t.Run(entityType, func(t *testing.T) {
			ruleKey := "TEST_RULE"
			et := entityType
			dto := &openapi.TgvalidatordBusinessRule{
				RuleKey:    &ruleKey,
				EntityType: &et,
			}
			got := BusinessRuleFromDTO(dto)
			if got == nil {
				t.Fatal("BusinessRuleFromDTO() returned nil for non-nil input")
			}
			if got.EntityType != entityType {
				t.Errorf("EntityType = %v, want %v", got.EntityType, entityType)
			}
		})
	}
}
