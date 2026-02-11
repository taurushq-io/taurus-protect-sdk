package taurusnetwork

import "time"

// Settlement represents a Taurus-NETWORK settlement.
type Settlement struct {
	// ID is the unique identifier for the settlement.
	ID string `json:"id"`
	// CreatorParticipantID is the ID of the participant who created the settlement.
	CreatorParticipantID string `json:"creator_participant_id"`
	// TargetParticipantID is the ID of the target participant.
	TargetParticipantID string `json:"target_participant_id"`
	// FirstLegParticipantID is the ID of the participant who executes the first leg.
	FirstLegParticipantID string `json:"first_leg_participant_id"`
	// FirstLegAssets contains the asset transfers for the first leg.
	FirstLegAssets []SettlementAssetTransfer `json:"first_leg_assets,omitempty"`
	// SecondLegAssets contains the asset transfers for the second leg.
	SecondLegAssets []SettlementAssetTransfer `json:"second_leg_assets,omitempty"`
	// Clips contains the settlement clips.
	Clips []SettlementClip `json:"clips,omitempty"`
	// StartExecutionDate is the date when the settlement execution should start.
	StartExecutionDate time.Time `json:"start_execution_date,omitempty"`
	// Status is the current status of the settlement.
	Status string `json:"status"`
	// WorkflowID is the ID of the workflow associated with this settlement.
	WorkflowID string `json:"workflow_id,omitempty"`
	// CreatedAt is when the settlement was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the settlement was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// SettlementAssetTransfer represents an asset transfer within a settlement.
type SettlementAssetTransfer struct {
	// CurrencyID is the ID of the currency being transferred.
	CurrencyID string `json:"currency_id"`
	// Amount is the amount being transferred.
	Amount string `json:"amount"`
	// SourceSharedAddressID is the ID of the source shared address.
	SourceSharedAddressID string `json:"source_shared_address_id,omitempty"`
	// DestinationSharedAddressID is the ID of the destination shared address.
	DestinationSharedAddressID string `json:"destination_shared_address_id,omitempty"`
}

// SettlementClip represents a clip within a settlement.
type SettlementClip struct {
	// ID is the unique identifier for the clip.
	ID string `json:"id"`
	// Index is the index of the clip within the settlement.
	Index string `json:"index"`
	// FirstLegTransactions contains the transactions for the first leg.
	FirstLegTransactions []SettlementClipTransaction `json:"first_leg_transactions,omitempty"`
	// SecondLegTransactions contains the transactions for the second leg.
	SecondLegTransactions []SettlementClipTransaction `json:"second_leg_transactions,omitempty"`
	// Status is the current status of the clip.
	Status string `json:"status"`
	// WorkflowID is the ID of the workflow associated with this clip.
	WorkflowID string `json:"workflow_id,omitempty"`
}

// SettlementClipTransaction represents a transaction within a settlement clip.
type SettlementClipTransaction struct {
	// ID is the unique identifier for the transaction.
	ID string `json:"id"`
	// AssetTransfer contains the asset transfer details.
	AssetTransfer *SettlementAssetTransfer `json:"asset_transfer,omitempty"`
	// RequestID is the ID of the associated request.
	RequestID string `json:"request_id,omitempty"`
	// TxHash is the transaction hash on the blockchain.
	TxHash string `json:"tx_hash,omitempty"`
	// TxID is the transaction ID.
	TxID string `json:"tx_id,omitempty"`
	// TxBlockNumber is the block number of the transaction.
	TxBlockNumber string `json:"tx_block_number,omitempty"`
	// Status is the current status of the transaction.
	Status string `json:"status"`
	// WorkflowID is the ID of the workflow associated with this transaction.
	WorkflowID string `json:"workflow_id,omitempty"`
	// CreatedAt is when the transaction was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// CreateSettlementRequest represents a request to create a settlement.
type CreateSettlementRequest struct {
	// TargetParticipantID is the ID of the target participant.
	TargetParticipantID string `json:"target_participant_id"`
	// FirstLegParticipantID is the ID of the participant who executes the first leg.
	FirstLegParticipantID string `json:"first_leg_participant_id"`
	// FirstLegAssets contains the asset transfers for the first leg.
	FirstLegAssets []SettlementAssetTransfer `json:"first_leg_assets"`
	// SecondLegAssets contains the asset transfers for the second leg.
	SecondLegAssets []SettlementAssetTransfer `json:"second_leg_assets"`
	// Clips contains the optional clip requests.
	Clips []CreateSettlementClipRequest `json:"clips,omitempty"`
	// StartExecutionDate is the optional date when the settlement execution should start.
	StartExecutionDate *time.Time `json:"start_execution_date,omitempty"`
}

// CreateSettlementClipRequest represents a clip request when creating a settlement.
type CreateSettlementClipRequest struct {
	// Index is the index of the clip.
	Index string `json:"index,omitempty"`
	// FirstLegAssets contains the asset transfers for the first leg of this clip.
	FirstLegAssets []SettlementAssetTransfer `json:"first_leg_assets,omitempty"`
	// SecondLegAssets contains the asset transfers for the second leg of this clip.
	SecondLegAssets []SettlementAssetTransfer `json:"second_leg_assets,omitempty"`
}

// CreateSettlementResult contains the result of creating a settlement.
type CreateSettlementResult struct {
	// SettlementID is the ID of the created settlement.
	SettlementID string `json:"settlement_id"`
}

// ReplaceSettlementRequest represents a request to replace a settlement.
type ReplaceSettlementRequest struct {
	// CreateSettlementRequest contains the new settlement attributes.
	CreateSettlementRequest *CreateSettlementRequest `json:"create_settlement_request,omitempty"`
}

// ListSettlementsOptions contains options for listing settlements.
type ListSettlementsOptions struct {
	// CounterParticipantID filters settlements where this ID is either the creator or target.
	CounterParticipantID string
	// Statuses filters by settlement statuses.
	Statuses []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListSettlementsResult contains the result of listing settlements.
type ListSettlementsResult struct {
	// Settlements is the list of settlements.
	Settlements []*Settlement `json:"settlements"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListSettlementsForApprovalOptions contains options for listing settlements for approval.
type ListSettlementsForApprovalOptions struct {
	// IDs filters by settlement IDs.
	IDs []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListSettlementsForApprovalResult contains the result of listing settlements for approval.
type ListSettlementsForApprovalResult struct {
	// Settlements is the list of settlements pending approval.
	Settlements []*Settlement `json:"settlements"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// Settlement status constants.
const (
	SettlementStatusCreating           = "CREATING"
	SettlementStatusCreated            = "CREATED"
	SettlementStatusRejectedByCreator  = "REJECTED_BY_CREATOR"
	SettlementStatusApprovedByCreator  = "APPROVED_BY_CREATOR"
	SettlementStatusReceived           = "RECEIVED"
	SettlementStatusRejectedByTarget   = "REJECTED_BY_TARGET"
	SettlementStatusAcceptedByTarget   = "ACCEPTED_BY_TARGET"
	SettlementStatusPending            = "PENDING"
	SettlementStatusCompleted          = "COMPLETED"
	SettlementStatusFailed             = "FAILED"
)
