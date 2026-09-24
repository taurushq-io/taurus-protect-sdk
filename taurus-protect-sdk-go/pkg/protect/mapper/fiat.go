package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FiatProviderFromDTO converts an OpenAPI TgvalidatordFiatProvider to a domain FiatProvider.
func FiatProviderFromDTO(dto *openapi.TgvalidatordFiatProvider) *model.FiatProvider {
	if dto == nil {
		return nil
	}

	return &model.FiatProvider{
		Provider:              safeString(dto.Provider),
		Label:                 safeString(dto.Label),
		BaseCurrencyValuation: safeString(dto.BaseCurrencyValuation),
	}
}

// FiatProvidersFromDTO converts a slice of OpenAPI TgvalidatordFiatProvider to domain FiatProviders.
func FiatProvidersFromDTO(dtos []openapi.TgvalidatordFiatProvider) []*model.FiatProvider {
	if dtos == nil {
		return nil
	}
	providers := make([]*model.FiatProvider, len(dtos))
	for i := range dtos {
		providers[i] = FiatProviderFromDTO(&dtos[i])
	}
	return providers
}

// FiatProviderAccountFromDTO converts an OpenAPI TgvalidatordFiatProviderAccount to a domain FiatProviderAccount.
func FiatProviderAccountFromDTO(dto *openapi.TgvalidatordFiatProviderAccount) *model.FiatProviderAccount {
	if dto == nil {
		return nil
	}

	account := &model.FiatProviderAccount{
		ID:                    safeString(dto.Id),
		Provider:              safeString(dto.Provider),
		Label:                 safeString(dto.Label),
		AccountType:           safeString(dto.AccountType),
		AccountIdentifier:     safeString(dto.AccountIdentifier),
		AccountName:           safeString(dto.AccountName),
		TotalBalance:          safeString(dto.TotalBalance),
		CurrencyID:            safeString(dto.CurrencyID),
		BaseCurrencyValuation: safeString(dto.BaseCurrencyValuation),
	}

	// Convert dates
	if dto.CreationDate != nil {
		account.CreationDate = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		account.UpdateDate = *dto.UpdateDate
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		account.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return account
}

// FiatProviderAccountsFromDTO converts a slice of OpenAPI TgvalidatordFiatProviderAccount to domain FiatProviderAccounts.
func FiatProviderAccountsFromDTO(dtos []openapi.TgvalidatordFiatProviderAccount) []*model.FiatProviderAccount {
	if dtos == nil {
		return nil
	}
	accounts := make([]*model.FiatProviderAccount, len(dtos))
	for i := range dtos {
		accounts[i] = FiatProviderAccountFromDTO(&dtos[i])
	}
	return accounts
}

// FiatProviderEntityFromDTO converts an OpenAPI FiatProviderEntity to a domain FiatProviderEntity.
func FiatProviderEntityFromDTO(dto *openapi.TgvalidatordFiatProviderEntity) *model.FiatProviderEntity {
	if dto == nil {
		return nil
	}
	entity := &model.FiatProviderEntity{
		ID:                safeString(dto.Id),
		Provider:          safeString(dto.Provider),
		Label:             safeString(dto.Label),
		AccountIdentifier: safeString(dto.AccountIdentifier),
		Name:              safeString(dto.Name),
		Details:           safeString(dto.Details),
	}
	if dto.CreationDate != nil {
		entity.CreationDate = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		entity.UpdateDate = *dto.UpdateDate
	}
	return entity
}

// FiatProviderEntitiesFromDTO converts a slice of OpenAPI FiatProviderEntity to domain entities.
func FiatProviderEntitiesFromDTO(dtos []openapi.TgvalidatordFiatProviderEntity) []*model.FiatProviderEntity {
	if dtos == nil {
		return nil
	}
	entities := make([]*model.FiatProviderEntity, len(dtos))
	for i := range dtos {
		entities[i] = FiatProviderEntityFromDTO(&dtos[i])
	}
	return entities
}
