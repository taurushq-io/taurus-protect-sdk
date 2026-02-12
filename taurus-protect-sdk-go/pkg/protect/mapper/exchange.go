package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ExchangeFromDTO converts an OpenAPI TgvalidatordExchange to a domain Exchange.
func ExchangeFromDTO(dto *openapi.TgvalidatordExchange) *model.Exchange {
	if dto == nil {
		return nil
	}

	exchange := &model.Exchange{
		ID:                    safeString(dto.Id),
		Exchange:              safeString(dto.Exchange),
		Account:               safeString(dto.Account),
		Currency:              safeString(dto.Currency),
		Type:                  safeString(dto.Type),
		TotalBalance:          safeString(dto.TotalBalance),
		Status:                safeString(dto.Status),
		Container:             safeString(dto.Container),
		Label:                 safeString(dto.Label),
		DisplayLabel:          safeString(dto.DisplayLabel),
		HasWLA:                safeBool(dto.HasWLA),
		BaseCurrencyValuation: safeString(dto.BaseCurrencyValuation),
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		exchange.CreationDate = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		exchange.UpdateDate = *dto.UpdateDate
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		exchange.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return exchange
}

// ExchangesFromDTO converts a slice of OpenAPI TgvalidatordExchange to domain Exchanges.
func ExchangesFromDTO(dtos []openapi.TgvalidatordExchange) []*model.Exchange {
	if dtos == nil {
		return nil
	}
	exchanges := make([]*model.Exchange, len(dtos))
	for i := range dtos {
		exchanges[i] = ExchangeFromDTO(&dtos[i])
	}
	return exchanges
}

// ExchangeCounterpartyFromDTO converts an OpenAPI TgvalidatordExchangeCounterparty to a domain ExchangeCounterparty.
func ExchangeCounterpartyFromDTO(dto *openapi.TgvalidatordExchangeCounterparty) *model.ExchangeCounterparty {
	if dto == nil {
		return nil
	}

	return &model.ExchangeCounterparty{
		Name:                  safeString(dto.Name),
		BaseCurrencyValuation: safeString(dto.BaseCurrencyValuation),
	}
}

// ExchangeCounterpartiesFromDTO converts a slice of OpenAPI TgvalidatordExchangeCounterparty to domain ExchangeCounterparties.
func ExchangeCounterpartiesFromDTO(dtos []openapi.TgvalidatordExchangeCounterparty) []*model.ExchangeCounterparty {
	if dtos == nil {
		return nil
	}
	counterparties := make([]*model.ExchangeCounterparty, len(dtos))
	for i := range dtos {
		counterparties[i] = ExchangeCounterpartyFromDTO(&dtos[i])
	}
	return counterparties
}
