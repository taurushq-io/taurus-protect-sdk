package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AssetBalanceFromDTO converts an OpenAPI AssetBalance to a domain AssetBalance.
func AssetBalanceFromDTO(dto *openapi.TgvalidatordAssetBalance) *model.AssetBalance {
	if dto == nil {
		return nil
	}

	return &model.AssetBalance{
		Asset:   AssetFromDTO(dto.Asset),
		Balance: BalanceFromDTO(dto.Balance),
	}
}

// AssetBalancesFromDTO converts a slice of OpenAPI AssetBalance to domain AssetBalances.
func AssetBalancesFromDTO(dtos []openapi.TgvalidatordAssetBalance) []*model.AssetBalance {
	if dtos == nil {
		return nil
	}
	balances := make([]*model.AssetBalance, len(dtos))
	for i := range dtos {
		balances[i] = AssetBalanceFromDTO(&dtos[i])
	}
	return balances
}

// AssetFromDTO converts an OpenAPI Asset to a domain Asset.
func AssetFromDTO(dto *openapi.TgvalidatordAsset) *model.Asset {
	if dto == nil {
		return nil
	}

	return &model.Asset{
		Currency:     dto.Currency,
		Kind:         safeString(dto.Kind),
		CurrencyInfo: CurrencyInfoFromDTO(dto.CurrencyInfo),
	}
}

// CurrencyInfoFromDTO converts an OpenAPI Currency to a domain CurrencyInfo.
func CurrencyInfoFromDTO(dto *openapi.TgvalidatordCurrency) *model.CurrencyInfo {
	if dto == nil {
		return nil
	}

	return &model.CurrencyInfo{
		ID:              safeString(dto.Id),
		Name:            safeString(dto.Name),
		Symbol:          safeString(dto.Symbol),
		DisplayName:     safeString(dto.DisplayName),
		Blockchain:      safeString(dto.Blockchain),
		Network:         safeString(dto.Network),
		Decimals:        safeString(dto.Decimals),
		ContractAddress: safeString(dto.ContractAddress),
		TokenID:         safeString(dto.TokenID),
		Type:            safeString(dto.Type),
		IsToken:         safeBool(dto.IsToken),
		IsERC20:         safeBool(dto.IsERC20),
		IsNFT:           safeBool(dto.IsNFT),
		IsFiat:          safeBool(dto.IsFiat),
		Enabled:         safeBool(dto.Enabled),
	}
}
