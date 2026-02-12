package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// BusinessRuleFromDTO converts an OpenAPI TgvalidatordBusinessRule to a domain BusinessRule.
func BusinessRuleFromDTO(dto *openapi.TgvalidatordBusinessRule) *model.BusinessRule {
	if dto == nil {
		return nil
	}

	rule := &model.BusinessRule{
		ID:              safeString(dto.Id),
		TenantID:        safeString(dto.TenantId),
		Currency:        safeString(dto.Currency),
		WalletID:        safeString(dto.WalletId),
		RuleKey:         safeString(dto.RuleKey),
		RuleValue:       safeString(dto.RuleValue),
		RuleGroup:       safeString(dto.RuleGroup),
		RuleDescription: safeString(dto.RuleDescription),
		RuleValidation:  safeString(dto.RuleValidation),
		AddressID:       safeString(dto.AddressId),
		EntityType:      safeString(dto.EntityType),
		EntityID:        safeString(dto.EntityID),
	}

	// Convert currency info if present
	if dto.CurrencyInfo != nil {
		rule.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return rule
}

// BusinessRulesFromDTO converts a slice of OpenAPI TgvalidatordBusinessRule to domain BusinessRules.
func BusinessRulesFromDTO(dtos []openapi.TgvalidatordBusinessRule) []*model.BusinessRule {
	if dtos == nil {
		return nil
	}
	rules := make([]*model.BusinessRule, len(dtos))
	for i := range dtos {
		rules[i] = BusinessRuleFromDTO(&dtos[i])
	}
	return rules
}
