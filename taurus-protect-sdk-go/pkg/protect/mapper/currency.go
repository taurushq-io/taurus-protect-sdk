package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// CurrencyFromDTO converts an OpenAPI TgvalidatordCurrency to a domain Currency.
func CurrencyFromDTO(dto *openapi.TgvalidatordCurrency) *model.Currency {
	if dto == nil {
		return nil
	}

	currency := &model.Currency{
		ID:              safeString(dto.Id),
		Name:            safeString(dto.Name),
		Symbol:          safeString(dto.Symbol),
		DisplayName:     safeString(dto.DisplayName),
		Type:            safeString(dto.Type),
		Blockchain:      safeString(dto.Blockchain),
		Network:         safeString(dto.Network),
		ContractAddress: safeString(dto.ContractAddress),
		TokenID:         safeString(dto.TokenID),
		CoinTypeIndex:   safeString(dto.CoinTypeIndex),
		Logo:            safeString(dto.Logo),
		IsToken:         safeBool(dto.IsToken),
		IsERC20:         safeBool(dto.IsERC20),
		IsFA12:          safeBool(dto.IsFA12),
		IsFA20:          safeBool(dto.IsFA20),
		IsNFT:           safeBool(dto.IsNFT),
		IsFiat:          safeBool(dto.IsFiat),
		IsUTXOBased:     safeBool(dto.IsUTXOBased),
		IsAccountBased:  safeBool(dto.IsAccountBased),
		HasStaking:      safeBool(dto.HasStaking),
		Enabled:         safeBool(dto.Enabled),
	}

	// Parse decimals from string to int64
	if dto.Decimals != nil {
		if decimals, err := strconv.ParseInt(*dto.Decimals, 10, 64); err == nil {
			currency.Decimals = decimals
		}
	}

	// Parse wlca ID from string to int64
	if dto.WlcaId != nil {
		if wlcaID, err := strconv.ParseInt(*dto.WlcaId, 10, 64); err == nil {
			currency.WlcaID = wlcaID
		}
	}

	return currency
}

// CurrenciesFromDTO converts a slice of OpenAPI TgvalidatordCurrency to domain Currencies.
func CurrenciesFromDTO(dtos []openapi.TgvalidatordCurrency) []*model.Currency {
	if dtos == nil {
		return nil
	}
	currencies := make([]*model.Currency, len(dtos))
	for i := range dtos {
		currencies[i] = CurrencyFromDTO(&dtos[i])
	}
	return currencies
}
