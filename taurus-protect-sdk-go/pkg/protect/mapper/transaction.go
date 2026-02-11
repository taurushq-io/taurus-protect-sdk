package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TransactionFromDTO converts an OpenAPI Transaction to a domain Transaction.
func TransactionFromDTO(dto *openapi.TgvalidatordTransaction) *model.Transaction {
	if dto == nil {
		return nil
	}

	transaction := &model.Transaction{
		ID:             safeString(dto.Id),
		UniqueID:       safeString(dto.UniqueId),
		TransactionID:  safeString(dto.TransactionId),
		Hash:           safeString(dto.Hash),
		Direction:      safeString(dto.Direction),
		Type:           safeString(dto.Type),
		Currency:       safeString(dto.Currency),
		Blockchain:     safeString(dto.Blockchain),
		Amount:         safeString(dto.Amount),
		AmountMainUnit: safeString(dto.AmountMainUnit),
		Fee:            safeString(dto.Fee),
		FeeMainUnit:    safeString(dto.FeeMainUnit),
		Block:          safeString(dto.Block),
	}

	// Convert timestamps
	if dto.ReceptionDate != nil {
		transaction.ReceptionDate = *dto.ReceptionDate
	}
	if dto.ConfirmationDate != nil {
		transaction.ConfirmationDate = *dto.ConfirmationDate
	}

	// Convert sources
	if dto.Sources != nil {
		transaction.Sources = make([]model.TransactionAddressInfo, len(dto.Sources))
		for i, src := range dto.Sources {
			transaction.Sources[i] = AddressInfoFromDTO(&src)
		}
	}

	// Convert destinations
	if dto.Destinations != nil {
		transaction.Destinations = make([]model.TransactionAddressInfo, len(dto.Destinations))
		for i, dst := range dto.Destinations {
			transaction.Destinations[i] = AddressInfoFromDTO(&dst)
		}
	}

	return transaction
}

// TransactionsFromDTO converts a slice of OpenAPI Transaction to domain Transactions.
func TransactionsFromDTO(dtos []openapi.TgvalidatordTransaction) []*model.Transaction {
	if dtos == nil {
		return nil
	}
	transactions := make([]*model.Transaction, len(dtos))
	for i := range dtos {
		transactions[i] = TransactionFromDTO(&dtos[i])
	}
	return transactions
}

// AddressInfoFromDTO converts an OpenAPI AddressInfo to a domain TransactionAddressInfo.
func AddressInfoFromDTO(dto *openapi.TgvalidatordAddressInfo) model.TransactionAddressInfo {
	if dto == nil {
		return model.TransactionAddressInfo{}
	}
	return model.TransactionAddressInfo{
		Address:              safeString(dto.Address),
		Label:                safeString(dto.Label),
		CustomerID:           safeString(dto.CustomerId),
		Amount:               safeString(dto.Amount),
		AmountMainUnit:       safeString(dto.AmountMainUnit),
		Type:                 safeString(dto.Type),
		InternalAddressID:    safeString(dto.InternalAddressId),
		WhitelistedAddressID: safeString(dto.WhitelistedAddressId),
	}
}
