package taurusnetwork

import (
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// LendingAgreement represents a Taurus Network lending agreement.
type LendingAgreement struct {
	// ID is the unique identifier for the lending agreement.
	ID string `json:"id"`
	// BorrowerParticipantID is the ID of the borrower participant.
	BorrowerParticipantID string `json:"borrower_participant_id"`
	// LenderParticipantID is the ID of the lender participant.
	LenderParticipantID string `json:"lender_participant_id"`
	// LendingOfferID is the ID of the lending offer this agreement is based on.
	LendingOfferID string `json:"lending_offer_id,omitempty"`
	// CurrencyID is the currency ID for the loan.
	CurrencyID string `json:"currency_id"`
	// Amount is the loan amount in the smallest currency unit.
	Amount string `json:"amount"`
	// AmountMainUnit is the loan amount in the main currency unit.
	AmountMainUnit string `json:"amount_main_unit,omitempty"`
	// AnnualYield is the interest rate.
	AnnualYield string `json:"annual_yield"`
	// AnnualYieldMainUnit is the interest rate in main unit.
	AnnualYieldMainUnit string `json:"annual_yield_main_unit,omitempty"`
	// Duration is the loan duration.
	Duration string `json:"duration"`
	// Status is the current status of the agreement.
	Status string `json:"status"`
	// WorkflowID is the workflow ID for the agreement.
	WorkflowID string `json:"workflow_id,omitempty"`
	// BorrowerSharedAddressID is the borrower's shared address ID for receiving funds.
	BorrowerSharedAddressID string `json:"borrower_shared_address_id,omitempty"`
	// LenderSharedAddressID is the lender's shared address ID for receiving repayment.
	LenderSharedAddressID string `json:"lender_shared_address_id,omitempty"`
	// Collaterals contains the collateral information for the agreement.
	Collaterals []LendingAgreementCollateral `json:"collaterals,omitempty"`
	// Transactions contains the transactions associated with the agreement.
	Transactions []LendingAgreementTransaction `json:"transactions,omitempty"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *model.Currency `json:"currency_info,omitempty"`
	// StartLoanDate is when the loan started.
	StartLoanDate time.Time `json:"start_loan_date,omitempty"`
	// RepaymentDueDate is when the repayment is due.
	RepaymentDueDate time.Time `json:"repayment_due_date,omitempty"`
	// CreatedAt is when the agreement was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the agreement was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingAgreementCollateral represents collateral for a lending agreement.
type LendingAgreementCollateral struct {
	// ID is the unique identifier for the collateral.
	ID string `json:"id"`
	// LendingAgreementID is the ID of the associated lending agreement.
	LendingAgreementID string `json:"lending_agreement_id"`
	// CurrencyID is the currency ID of the collateral.
	CurrencyID string `json:"currency_id"`
	// Amount is the collateral amount in the smallest currency unit.
	Amount string `json:"amount"`
	// AmountMainUnit is the collateral amount in the main currency unit.
	AmountMainUnit string `json:"amount_main_unit,omitempty"`
	// Status is the current status of the collateral.
	Status string `json:"status,omitempty"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *model.Currency `json:"currency_info,omitempty"`
	// CreatedAt is when the collateral was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the collateral was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingAgreementTransaction represents a transaction within a lending agreement.
type LendingAgreementTransaction struct {
	// ID is the unique identifier for the transaction.
	ID string `json:"id"`
	// LendingAgreementID is the ID of the associated lending agreement.
	LendingAgreementID string `json:"lending_agreement_id"`
	// Amount is the transaction amount in the smallest currency unit.
	Amount string `json:"amount"`
	// AmountMainUnit is the transaction amount in the main currency unit.
	AmountMainUnit string `json:"amount_main_unit,omitempty"`
	// CurrencyID is the currency ID of the transaction.
	CurrencyID string `json:"currency_id"`
	// RequestID is the ID of the associated request.
	RequestID string `json:"request_id,omitempty"`
	// TransactionID is the ID of the blockchain transaction.
	TransactionID string `json:"transaction_id,omitempty"`
	// TransactionHash is the blockchain transaction hash.
	TransactionHash string `json:"transaction_hash,omitempty"`
	// TransactionBlockNumber is the block number of the transaction.
	TransactionBlockNumber string `json:"transaction_block_number,omitempty"`
	// Type is the type of transaction (e.g., "LOAN", "REPAYMENT", "COLLATERAL").
	Type string `json:"type"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *model.Currency `json:"currency_info,omitempty"`
	// CreatedAt is when the transaction was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the transaction was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingAgreementAttachment represents an attachment for a lending agreement.
type LendingAgreementAttachment struct {
	// ID is the unique identifier for the attachment.
	ID string `json:"id"`
	// LendingAgreementID is the ID of the associated lending agreement.
	LendingAgreementID string `json:"lending_agreement_id"`
	// UploaderParticipantID is the ID of the participant who uploaded the attachment.
	UploaderParticipantID string `json:"uploader_participant_id"`
	// Name is the name of the attachment.
	Name string `json:"name"`
	// Type is the type of attachment (e.g., "EMBEDDED", "EXTERNAL_LINK").
	Type string `json:"type"`
	// ContentType is the MIME type of the attachment (e.g., "application/pdf").
	ContentType string `json:"content_type,omitempty"`
	// Value is the content (base64 for EMBEDDED, URL for EXTERNAL_LINK).
	Value string `json:"value,omitempty"`
	// FileSize is the size of the file in bytes (only for EMBEDDED type).
	FileSize string `json:"file_size,omitempty"`
	// CreatedAt is when the attachment was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the attachment was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingOffer represents a Taurus Network lending offer.
type LendingOffer struct {
	// ID is the unique identifier for the lending offer.
	ID string `json:"id"`
	// ParticipantID is the ID of the participant who created the offer.
	ParticipantID string `json:"participant_id"`
	// Blockchain is the blockchain for the offer.
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network for the offer.
	Network string `json:"network,omitempty"`
	// Arg1 is additional argument 1.
	Arg1 string `json:"arg1,omitempty"`
	// Arg2 is additional argument 2.
	Arg2 string `json:"arg2,omitempty"`
	// AnnualPercentageYield is the interest rate (5 decimals, e.g., 525000 = 5.25%).
	AnnualPercentageYield string `json:"annual_percentage_yield"`
	// AnnualPercentageYieldMainUnit is the interest rate in main unit.
	AnnualPercentageYieldMainUnit string `json:"annual_percentage_yield_main_unit,omitempty"`
	// Duration is the loan duration.
	Duration string `json:"duration"`
	// CollateralRequirement contains the collateral requirements for the offer.
	CollateralRequirement *LendingCollateralRequirement `json:"collateral_requirement,omitempty"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *model.Currency `json:"currency_info,omitempty"`
	// OriginCreatedAt is the original creation time.
	OriginCreatedAt time.Time `json:"origin_created_at,omitempty"`
	// CreatedAt is when the offer was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the offer was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingCollateralRequirement represents collateral requirements for a lending offer.
type LendingCollateralRequirement struct {
	// AcceptedCurrencies lists the currencies accepted as collateral.
	AcceptedCurrencies []CurrencyCollateralRequirement `json:"accepted_currencies,omitempty"`
}

// CurrencyCollateralRequirement represents a currency requirement for collateral.
type CurrencyCollateralRequirement struct {
	// Blockchain is the blockchain for the collateral.
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network for the collateral.
	Network string `json:"network,omitempty"`
	// Arg1 is additional argument 1.
	Arg1 string `json:"arg1,omitempty"`
	// Arg2 is additional argument 2.
	Arg2 string `json:"arg2,omitempty"`
	// Ratio is the loan-to-value ratio (2 decimals, e.g., 12500 = 125%).
	Ratio string `json:"ratio,omitempty"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *model.Currency `json:"currency_info,omitempty"`
}

// CreateLendingAgreementRequest contains parameters for creating a lending agreement.
type CreateLendingAgreementRequest struct {
	// LendingOfferID is the ID of the lending offer (optional if negotiated directly).
	LendingOfferID string `json:"lending_offer_id,omitempty"`
	// LenderParticipantID is the ID of the lender participant.
	LenderParticipantID string `json:"lender_participant_id,omitempty"`
	// CurrencyID is the currency ID for the loan.
	CurrencyID string `json:"currency_id,omitempty"`
	// Amount is the loan amount in the smallest currency unit.
	Amount string `json:"amount,omitempty"`
	// AnnualPercentageYield is the interest rate (5 decimals).
	AnnualPercentageYield string `json:"annual_percentage_yield,omitempty"`
	// Duration is the loan duration.
	Duration string `json:"duration,omitempty"`
	// BorrowerSharedAddressID is the borrower's shared address for receiving funds.
	BorrowerSharedAddressID string `json:"borrower_shared_address_id,omitempty"`
	// Collaterals contains the collateral to provide.
	Collaterals []CreateLendingCollateralRequest `json:"collaterals,omitempty"`
}

// CreateLendingCollateralRequest contains parameters for creating collateral.
type CreateLendingCollateralRequest struct {
	// CurrencyID is the currency ID of the collateral.
	CurrencyID string `json:"currency_id"`
	// Amount is the collateral amount in the smallest currency unit.
	Amount string `json:"amount"`
}

// UpdateLendingAgreementRequest contains parameters for updating a lending agreement.
type UpdateLendingAgreementRequest struct {
	// LenderSharedAddressID is the lender's shared address for receiving repayment.
	LenderSharedAddressID string `json:"lender_shared_address_id,omitempty"`
}

// RepayLendingAgreementRequest contains parameters for repaying a lending agreement.
type RepayLendingAgreementRequest struct {
	// RepayerSharedAddressID is the shared address to use for repayment.
	RepayerSharedAddressID string `json:"repayer_shared_address_id,omitempty"`
}

// CreateLendingOfferRequest contains parameters for creating a lending offer.
type CreateLendingOfferRequest struct {
	// CurrencyID is the currency ID for the loan.
	CurrencyID string `json:"currency_id,omitempty"`
	// Amount is the loan amount in the smallest currency unit.
	Amount string `json:"amount,omitempty"`
	// AnnualPercentageYield is the interest rate (5 decimals).
	AnnualPercentageYield string `json:"annual_percentage_yield,omitempty"`
	// Duration is the loan duration.
	Duration string `json:"duration,omitempty"`
	// CollateralRequirements contains the collateral requirements.
	CollateralRequirements []CreateCollateralRequirementRequest `json:"collateral_requirements,omitempty"`
}

// CreateCollateralRequirementRequest contains parameters for creating a collateral requirement.
type CreateCollateralRequirementRequest struct {
	// CurrencyID is the currency ID accepted as collateral.
	CurrencyID string `json:"currency_id"`
	// CollateralRatio is the required collateral ratio.
	CollateralRatio string `json:"collateral_ratio,omitempty"`
}

// CreateLendingAgreementAttachmentRequest contains parameters for creating an attachment.
type CreateLendingAgreementAttachmentRequest struct {
	// Name is the name of the attachment.
	Name string `json:"name"`
	// Type is the type of attachment ("EMBEDDED" or "EXTERNAL_LINK").
	Type string `json:"type"`
	// ContentType is the MIME type (e.g., "application/pdf").
	ContentType string `json:"content_type,omitempty"`
	// Value is the content (base64 for EMBEDDED, URL for EXTERNAL_LINK).
	Value string `json:"value"`
}

// ListLendingAgreementsOptions contains options for listing lending agreements.
type ListLendingAgreementsOptions struct {
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
}

// ListLendingAgreementsResult contains the result of listing lending agreements.
type ListLendingAgreementsResult struct {
	// LendingAgreements is the list of lending agreements.
	LendingAgreements []*LendingAgreement `json:"lending_agreements"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListLendingAgreementsForApprovalOptions contains options for listing agreements for approval.
type ListLendingAgreementsForApprovalOptions struct {
	// IDs filters by specific agreement IDs.
	IDs []string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
}

// ListLendingAgreementsForApprovalResult contains the result of listing agreements for approval.
type ListLendingAgreementsForApprovalResult struct {
	// LendingAgreements is the list of lending agreements pending approval.
	LendingAgreements []*LendingAgreement `json:"lending_agreements"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListLendingOffersOptions contains options for listing lending offers.
type ListLendingOffersOptions struct {
	// CurrencyIDs filters by currency IDs.
	CurrencyIDs []string
	// ParticipantID filters by participant ID.
	ParticipantID string
	// Duration filters by duration.
	Duration string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
}

// ListLendingOffersResult contains the result of listing lending offers.
type ListLendingOffersResult struct {
	// LendingOffers is the list of lending offers.
	LendingOffers []*LendingOffer `json:"lending_offers"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
