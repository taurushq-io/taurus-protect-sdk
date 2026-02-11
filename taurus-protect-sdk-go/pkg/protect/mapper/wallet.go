// Package mapper provides functions to convert between OpenAPI DTOs and domain models.
package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WalletFromDTO converts an OpenAPI WalletInfo to a domain Wallet.
func WalletFromDTO(dto *openapi.TgvalidatordWalletInfo) *model.Wallet {
	if dto == nil {
		return nil
	}

	wallet := &model.Wallet{
		ID:                safeString(dto.Id),
		Name:              safeString(dto.Name),
		Currency:          safeString(dto.Currency),
		Blockchain:        safeString(dto.Blockchain),
		Network:           safeString(dto.Network),
		IsOmnibus:         safeBool(dto.IsOmnibus),
		Disabled:          safeBool(dto.Disabled),
		Comment:           safeString(dto.Comment),
		CustomerID:        safeString(dto.CustomerId),
		ExternalWalletID:  safeString(dto.ExternalWalletId),
		VisibilityGroupID: safeString(dto.VisibilityGroupID),
		AccountPath:       safeString(dto.AccountPath),
	}

	// Parse addresses count
	if dto.AddressesCount != nil {
		if count, err := strconv.ParseInt(*dto.AddressesCount, 10, 64); err == nil {
			wallet.AddressesCount = count
		}
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		wallet.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		wallet.UpdatedAt = *dto.UpdateDate
	}

	// Convert balance
	if dto.Balance != nil {
		wallet.Balance = BalanceFromDTO(dto.Balance)
	}

	// Convert attributes
	if dto.Attributes != nil {
		wallet.Attributes = make([]model.WalletAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			wallet.Attributes[i] = WalletAttributeFromDTO(&attr)
		}
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		wallet.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return wallet
}

// WalletFromCreateDTO converts an OpenAPI Wallet (from create response) to a domain Wallet.
func WalletFromCreateDTO(dto *openapi.TgvalidatordWallet) *model.Wallet {
	if dto == nil {
		return nil
	}

	wallet := &model.Wallet{
		ID:               safeString(dto.Id),
		Name:             safeString(dto.Name),
		Currency:         safeString(dto.Currency),
		Blockchain:       safeString(dto.Blockchain),
		IsOmnibus:        safeBool(dto.IsOmnibus),
		Disabled:         safeBool(dto.Disabled),
		Comment:          safeString(dto.Comment),
		CustomerID:       safeString(dto.CustomerId),
		ExternalWalletID: safeString(dto.ExternalWalletId),
		AccountPath:      safeString(dto.AccountPath),
	}

	// Parse addresses count
	if dto.AddressesCount != nil {
		if count, err := strconv.ParseInt(*dto.AddressesCount, 10, 64); err == nil {
			wallet.AddressesCount = count
		}
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		wallet.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		wallet.UpdatedAt = *dto.UpdateDate
	}

	// Convert balance
	if dto.Balance != nil {
		wallet.Balance = BalanceFromDTO(dto.Balance)
	}

	// Convert attributes
	if dto.Attributes != nil {
		wallet.Attributes = make([]model.WalletAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			wallet.Attributes[i] = WalletAttributeFromDTO(&attr)
		}
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		wallet.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return wallet
}

// WalletsFromDTO converts a slice of OpenAPI WalletInfo to domain Wallets.
func WalletsFromDTO(dtos []openapi.TgvalidatordWalletInfo) []*model.Wallet {
	if dtos == nil {
		return nil
	}
	wallets := make([]*model.Wallet, len(dtos))
	for i := range dtos {
		wallets[i] = WalletFromDTO(&dtos[i])
	}
	return wallets
}

// WalletAttributeFromDTO converts an OpenAPI WalletAttribute to a domain WalletAttribute.
func WalletAttributeFromDTO(dto *openapi.TgvalidatordWalletAttribute) model.WalletAttribute {
	if dto == nil {
		return model.WalletAttribute{}
	}
	return model.WalletAttribute{
		ID:          safeString(dto.Id),
		Key:         safeString(dto.Key),
		Value:       safeString(dto.Value),
		ContentType: safeString(dto.ContentType),
		Owner:       safeString(dto.Owner),
		Type:        safeString(dto.Type),
		Subtype:     safeString(dto.Subtype),
		IsFile:      safeBool(dto.Isfile),
	}
}

// BalanceFromDTO converts an OpenAPI Balance to a domain Balance.
func BalanceFromDTO(dto *openapi.TgvalidatordBalance) *model.Balance {
	if dto == nil {
		return nil
	}
	return &model.Balance{
		TotalConfirmed:       safeString(dto.TotalConfirmed),
		TotalUnconfirmed:     safeString(dto.TotalUnconfirmed),
		AvailableConfirmed:   safeString(dto.AvailableConfirmed),
		AvailableUnconfirmed: safeString(dto.AvailableUnconfirmed),
		ReservedConfirmed:    safeString(dto.ReservedConfirmed),
		ReservedUnconfirmed:  safeString(dto.ReservedUnconfirmed),
	}
}

// BalanceHistoryPointFromDTO converts an OpenAPI BalanceHistoryPoint to a domain BalanceHistoryPoint.
func BalanceHistoryPointFromDTO(dto *openapi.TgvalidatordBalanceHistoryPoint) *model.BalanceHistoryPoint {
	if dto == nil {
		return nil
	}
	point := &model.BalanceHistoryPoint{}
	if dto.PointDate != nil {
		point.PointDate = *dto.PointDate
	}
	if dto.Balance != nil {
		point.Balance = BalanceFromDTO(dto.Balance)
	}
	return point
}

// BalanceHistoryPointsFromDTO converts a slice of OpenAPI BalanceHistoryPoint to domain BalanceHistoryPoints.
func BalanceHistoryPointsFromDTO(dtos []openapi.TgvalidatordBalanceHistoryPoint) []*model.BalanceHistoryPoint {
	if dtos == nil {
		return nil
	}
	points := make([]*model.BalanceHistoryPoint, len(dtos))
	for i := range dtos {
		points[i] = BalanceHistoryPointFromDTO(&dtos[i])
	}
	return points
}

