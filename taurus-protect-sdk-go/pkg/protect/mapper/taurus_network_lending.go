package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// LendingAgreementFromDTO converts an OpenAPI TgvalidatordLendingAgreement to a domain LendingAgreement.
func LendingAgreementFromDTO(dto *openapi.TgvalidatordLendingAgreement) *taurusnetwork.LendingAgreement {
	if dto == nil {
		return nil
	}

	agreement := &taurusnetwork.LendingAgreement{
		ID:                      safeString(dto.Id),
		BorrowerParticipantID:   safeString(dto.BorrowerParticipantID),
		LenderParticipantID:     safeString(dto.LenderParticipantID),
		LendingOfferID:          safeString(dto.LendingOfferID),
		CurrencyID:              safeString(dto.CurrencyID),
		Amount:                  safeString(dto.Amount),
		AmountMainUnit:          safeString(dto.AmountMainUnit),
		AnnualYield:             safeString(dto.AnnualYield),
		AnnualYieldMainUnit:     safeString(dto.AnnualYieldMainUnit),
		Duration:                safeString(dto.Duration),
		Status:                  safeString(dto.Status),
		WorkflowID:              safeString(dto.WorkflowID),
		BorrowerSharedAddressID: safeString(dto.BorrowerSharedAddressID),
		LenderSharedAddressID:   safeString(dto.LenderSharedAddressID),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		agreement.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		agreement.UpdatedAt = *dto.UpdatedAt
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		agreement.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	// Convert collaterals
	if dto.LendingAgreementCollaterals != nil {
		agreement.Collaterals = LendingAgreementCollateralsFromDTO(dto.LendingAgreementCollaterals)
	}

	// Convert transactions
	if dto.LendingAgreementTransactions != nil {
		agreement.Transactions = LendingAgreementTransactionsFromDTO(dto.LendingAgreementTransactions)
	}

	// Convert start loan date and repayment due date
	if dto.StartLoanDate != nil {
		agreement.StartLoanDate = *dto.StartLoanDate
	}
	if dto.RepaymentDueDate != nil {
		agreement.RepaymentDueDate = *dto.RepaymentDueDate
	}

	return agreement
}

// LendingAgreementsFromDTO converts a slice of OpenAPI TgvalidatordLendingAgreement to domain LendingAgreements.
func LendingAgreementsFromDTO(dtos []openapi.TgvalidatordLendingAgreement) []*taurusnetwork.LendingAgreement {
	if dtos == nil {
		return nil
	}
	agreements := make([]*taurusnetwork.LendingAgreement, len(dtos))
	for i := range dtos {
		agreements[i] = LendingAgreementFromDTO(&dtos[i])
	}
	return agreements
}

// LendingAgreementCollateralFromDTO converts an OpenAPI TgvalidatordLendingAgreementCollateral to a domain LendingAgreementCollateral.
func LendingAgreementCollateralFromDTO(dto *openapi.TgvalidatordLendingAgreementCollateral) taurusnetwork.LendingAgreementCollateral {
	if dto == nil {
		return taurusnetwork.LendingAgreementCollateral{}
	}

	collateral := taurusnetwork.LendingAgreementCollateral{
		ID:                 safeString(dto.Id),
		LendingAgreementID: safeString(dto.LendingAgreementID),
		CurrencyID:         safeString(dto.CurrencyID),
		Amount:             safeString(dto.Amount),
		AmountMainUnit:     safeString(dto.AmountMainUnit),
		Status:             safeString(dto.Status),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		collateral.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		collateral.UpdatedAt = *dto.UpdatedAt
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		collateral.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return collateral
}

// LendingAgreementCollateralsFromDTO converts a slice of OpenAPI collaterals to domain collaterals.
func LendingAgreementCollateralsFromDTO(dtos []openapi.TgvalidatordLendingAgreementCollateral) []taurusnetwork.LendingAgreementCollateral {
	if dtos == nil {
		return nil
	}
	collaterals := make([]taurusnetwork.LendingAgreementCollateral, len(dtos))
	for i := range dtos {
		collaterals[i] = LendingAgreementCollateralFromDTO(&dtos[i])
	}
	return collaterals
}

// LendingAgreementTransactionFromDTO converts an OpenAPI TgvalidatordLendingAgreementTransaction to a domain LendingAgreementTransaction.
func LendingAgreementTransactionFromDTO(dto *openapi.TgvalidatordLendingAgreementTransaction) taurusnetwork.LendingAgreementTransaction {
	if dto == nil {
		return taurusnetwork.LendingAgreementTransaction{}
	}

	transaction := taurusnetwork.LendingAgreementTransaction{
		ID:                     safeString(dto.Id),
		LendingAgreementID:     safeString(dto.LendingAgreementID),
		Amount:                 safeString(dto.Amount),
		AmountMainUnit:         safeString(dto.AmountMainUnit),
		CurrencyID:             safeString(dto.CurrencyID),
		RequestID:              safeString(dto.RequestID),
		TransactionID:          safeString(dto.TransactionID),
		TransactionHash:        safeString(dto.TransactionHash),
		TransactionBlockNumber: safeString(dto.TransactionBlockNumber),
		Type:                   safeString(dto.Type),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		transaction.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		transaction.UpdatedAt = *dto.UpdatedAt
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		transaction.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	return transaction
}

// LendingAgreementTransactionsFromDTO converts a slice of OpenAPI transactions to domain transactions.
func LendingAgreementTransactionsFromDTO(dtos []openapi.TgvalidatordLendingAgreementTransaction) []taurusnetwork.LendingAgreementTransaction {
	if dtos == nil {
		return nil
	}
	transactions := make([]taurusnetwork.LendingAgreementTransaction, len(dtos))
	for i := range dtos {
		transactions[i] = LendingAgreementTransactionFromDTO(&dtos[i])
	}
	return transactions
}

// LendingAgreementAttachmentFromDTO converts an OpenAPI TgvalidatordLendingAgreementAttachment to a domain LendingAgreementAttachment.
func LendingAgreementAttachmentFromDTO(dto *openapi.TgvalidatordLendingAgreementAttachment) *taurusnetwork.LendingAgreementAttachment {
	if dto == nil {
		return nil
	}

	attachment := &taurusnetwork.LendingAgreementAttachment{
		ID:                    safeString(dto.Id),
		LendingAgreementID:    safeString(dto.LendingAgreementID),
		UploaderParticipantID: safeString(dto.UploaderParticipantID),
		Name:                  safeString(dto.Name),
		Type:                  safeString(dto.Type),
		ContentType:           safeString(dto.ContentType),
		Value:                 safeString(dto.Value),
		FileSize:              safeString(dto.FileSize),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		attachment.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		attachment.UpdatedAt = *dto.UpdatedAt
	}

	return attachment
}

// LendingAgreementAttachmentsFromDTO converts a slice of OpenAPI attachments to domain attachments.
func LendingAgreementAttachmentsFromDTO(dtos []openapi.TgvalidatordLendingAgreementAttachment) []*taurusnetwork.LendingAgreementAttachment {
	if dtos == nil {
		return nil
	}
	attachments := make([]*taurusnetwork.LendingAgreementAttachment, len(dtos))
	for i := range dtos {
		attachments[i] = LendingAgreementAttachmentFromDTO(&dtos[i])
	}
	return attachments
}

// LendingOfferFromDTO converts an OpenAPI TgvalidatordTnLendingOffer to a domain LendingOffer.
func LendingOfferFromDTO(dto *openapi.TgvalidatordTnLendingOffer) *taurusnetwork.LendingOffer {
	if dto == nil {
		return nil
	}

	offer := &taurusnetwork.LendingOffer{
		ID:                            safeString(dto.Id),
		ParticipantID:                 safeString(dto.ParticipantID),
		Blockchain:                    safeString(dto.Blockchain),
		Network:                       safeString(dto.Network),
		Arg1:                          safeString(dto.Arg1),
		Arg2:                          safeString(dto.Arg2),
		AnnualPercentageYield:         safeString(dto.AnnualPercentageYield),
		AnnualPercentageYieldMainUnit: safeString(dto.AnnualPercentageYieldMainUnit),
		Duration:                      safeString(dto.Duration),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		offer.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		offer.UpdatedAt = *dto.UpdatedAt
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		offer.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}

	// Convert collateral requirement
	if dto.CollateralRequirement != nil {
		offer.CollateralRequirement = LendingCollateralRequirementFromDTO(dto.CollateralRequirement)
	}

	return offer
}

// LendingOffersFromDTO converts a slice of OpenAPI TgvalidatordTnLendingOffer to domain LendingOffers.
func LendingOffersFromDTO(dtos []openapi.TgvalidatordTnLendingOffer) []*taurusnetwork.LendingOffer {
	if dtos == nil {
		return nil
	}
	offers := make([]*taurusnetwork.LendingOffer, len(dtos))
	for i := range dtos {
		offers[i] = LendingOfferFromDTO(&dtos[i])
	}
	return offers
}

// LendingCollateralRequirementFromDTO converts an OpenAPI TgvalidatordLendingCollateralRequirement to a domain LendingCollateralRequirement.
func LendingCollateralRequirementFromDTO(dto *openapi.TgvalidatordLendingCollateralRequirement) *taurusnetwork.LendingCollateralRequirement {
	if dto == nil {
		return nil
	}

	requirement := &taurusnetwork.LendingCollateralRequirement{}

	// Convert accepted currencies
	if dto.AcceptedCurrencies != nil {
		requirement.AcceptedCurrencies = CurrencyCollateralRequirementsFromDTO(dto.AcceptedCurrencies)
	}

	return requirement
}

// CurrencyCollateralRequirementFromDTO converts an OpenAPI TgvalidatordCurrencyCollateralRequirement to a domain CurrencyCollateralRequirement.
func CurrencyCollateralRequirementFromDTO(dto *openapi.TgvalidatordCurrencyCollateralRequirement) taurusnetwork.CurrencyCollateralRequirement {
	if dto == nil {
		return taurusnetwork.CurrencyCollateralRequirement{}
	}

	req := taurusnetwork.CurrencyCollateralRequirement{
		Blockchain: safeString(dto.Blockchain),
		Network:    safeString(dto.Network),
		Arg1:       safeString(dto.Arg1),
		Arg2:       safeString(dto.Arg2),
		Ratio:      safeString(dto.Ratio),
	}
	if dto.CurrencyInfo != nil {
		req.CurrencyInfo = CurrencyFromDTO(dto.CurrencyInfo)
	}
	return req
}

// CurrencyCollateralRequirementsFromDTO converts a slice of OpenAPI currency collateral requirements to domain requirements.
func CurrencyCollateralRequirementsFromDTO(dtos []openapi.TgvalidatordCurrencyCollateralRequirement) []taurusnetwork.CurrencyCollateralRequirement {
	if dtos == nil {
		return nil
	}
	requirements := make([]taurusnetwork.CurrencyCollateralRequirement, len(dtos))
	for i := range dtos {
		requirements[i] = CurrencyCollateralRequirementFromDTO(&dtos[i])
	}
	return requirements
}
