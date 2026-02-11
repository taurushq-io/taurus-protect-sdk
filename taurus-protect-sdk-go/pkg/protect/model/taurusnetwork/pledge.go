package taurusnetwork

import (
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Pledge represents a Taurus Network pledge.
type Pledge struct {
	ID                  string               `json:"id"`
	OwnerParticipantID  string               `json:"ownerParticipantID"`
	TargetParticipantID string               `json:"targetParticipantID"`
	SharedAddressID     string               `json:"sharedAddressID"`
	CurrencyID          string               `json:"currencyID"`
	Blockchain          string               `json:"blockchain"`
	Network             string               `json:"network"`
	Amount              string               `json:"amount"`
	Status              string               `json:"status"`
	PledgeType          string               `json:"pledgeType"`
	Direction           string               `json:"direction"`
	ExternalReferenceID string               `json:"externalReferenceId"`
	ReconciliationNote  string               `json:"reconciliationNote"`
	WladdressID         string               `json:"wladdressID"`
	DurationSetup       *PledgeDurationSetup `json:"durationSetup,omitempty"`
	Attributes          []PledgeAttribute    `json:"attributes,omitempty"`
	Trails              []PledgeTrail        `json:"trails,omitempty"`
	OriginCreationDate  time.Time            `json:"originCreationDate"`
	UnpledgeDate        time.Time            `json:"unpledgeDate"`
	CreatedAt           time.Time            `json:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
}

// PledgeDurationSetup contains the duration configuration for a pledge.
type PledgeDurationSetup struct {
	MinimumDuration          string    `json:"minimumDuration"`
	EndOfMinimumDurationDate time.Time `json:"endOfMinimumDurationDate"`
	NoticePeriodDuration     string    `json:"noticePeriodDuration"`
	EndOfNoticePeriodDate    time.Time `json:"endOfNoticePeriodDate"`
}

// PledgeAttribute represents a key-value attribute on a pledge.
type PledgeAttribute struct {
	ID                   string `json:"id"`
	Key                  string `json:"key"`
	Value                string `json:"value"`
	Owner                string `json:"owner"`
	Type                 string `json:"type"`
	Subtype              string `json:"subtype"`
	ContentType          string `json:"contentType"`
	IsTaurusNetworkShared bool  `json:"isTaurusNetworkShared"`
}

// PledgeTrail represents a trail entry for a pledge.
type PledgeTrail struct {
	ID               string    `json:"id"`
	PledgeID         string    `json:"pledgeID"`
	AddressCommandID string    `json:"addressCommandID"`
	ParticipantID    string    `json:"participantID"`
	PledgeAmount     string    `json:"pledgeAmount"`
	Action           string    `json:"action"`
	Comment          string    `json:"comment"`
	CreatedAt        time.Time `json:"createdAt"`
}

// PledgeWithdrawal represents a withdrawal from a pledge.
type PledgeWithdrawal struct {
	ID                         string                   `json:"id"`
	PledgeID                   string                   `json:"pledgeID"`
	DestinationSharedAddressID string                   `json:"destinationSharedAddressID"`
	Amount                     string                   `json:"amount"`
	Status                     string                   `json:"status"`
	TxHash                     string                   `json:"txHash"`
	TxID                       string                   `json:"txID"`
	RequestID                  string                   `json:"requestID"`
	TxBlockNumber              string                   `json:"txBlockNumber"`
	InitiatorParticipantID     string                   `json:"initiatorParticipantID"`
	ExternalReferenceID        string                   `json:"externalReferenceID"`
	Trails                     []PledgeWithdrawalTrail  `json:"trails,omitempty"`
	CreatedAt                  time.Time                `json:"createdAt"`
}

// PledgeWithdrawalTrail represents a trail entry for a pledge withdrawal.
type PledgeWithdrawalTrail struct {
	ID                 string    `json:"id"`
	PledgeWithdrawalID string    `json:"pledgeWithdrawalID"`
	AddressCommandID   string    `json:"addressCommandID"`
	ParticipantID      string    `json:"participantID"`
	Action             string    `json:"action"`
	Comment            string    `json:"comment"`
	CreatedAt          time.Time `json:"createdAt"`
}

// PledgeAction represents an action on a pledge.
type PledgeAction struct {
	ID                  string              `json:"id"`
	PledgeID            string              `json:"pledgeID"`
	ActionType          string              `json:"actionType"`
	Status              string              `json:"status"`
	Metadata            *model.RequestMetadata    `json:"metadata,omitempty"`
	Rule                string              `json:"rule"`
	Approvers           *model.Approvers          `json:"approvers,omitempty"`
	NeedsApprovalFrom   []string            `json:"needsApprovalFrom,omitempty"`
	Envelope            string              `json:"envelope"`
	PledgeWithdrawalID  string              `json:"pledgeWithdrawalID"`
	Trails              []PledgeActionTrail `json:"trails,omitempty"`
	CreatedAt           time.Time           `json:"createdAt"`
	LastApprovalDate    time.Time           `json:"lastApprovalDate"`
}

// Note: We use model.RequestMetadata for metadata since it has the same structure.

// Note: model.Approvers, model.ParallelApproversGroups, and model.ApproversGroup types
// are defined in the parent model package. We reuse those definitions here.

// PledgeActionTrail represents a trail entry for a pledge action.
type PledgeActionTrail struct {
	ID             string    `json:"id"`
	PledgeActionID string    `json:"pledgeActionID"`
	UserID         string    `json:"userID"`
	ExternalUserID string    `json:"externalUserID"`
	Action         string    `json:"action"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
}

// Note: model.CursorPagination type is defined in the parent model package. We reuse that definition here.

// CreatePledgeRequest represents a request to create a new pledge.
type CreatePledgeRequest struct {
	SharedAddressID      string                          `json:"sharedAddressID"`
	CurrencyID           string                          `json:"currencyID"`
	Amount               string                          `json:"amount"`
	PledgeType           string                          `json:"pledgeType"`
	ExternalReferenceID  string                          `json:"externalReferenceId,omitempty"`
	ReconciliationNote   string                          `json:"reconciliationNote,omitempty"`
	PledgeDurationSetup  *CreatePledgeDurationSetup      `json:"pledgeDurationSetup,omitempty"`
	KeyValueAttributes   []KeyValue                      `json:"keyValueAttributes,omitempty"`
}

// CreatePledgeDurationSetup represents the duration setup for a new pledge.
type CreatePledgeDurationSetup struct {
	MinimumDuration          string     `json:"minimumDuration,omitempty"`
	EndOfMinimumDurationDate *time.Time `json:"endOfMinimumDurationDate,omitempty"`
	NoticePeriodDuration     string     `json:"noticePeriodDuration,omitempty"`
}

// KeyValue represents a key-value pair for pledge attributes.
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreatePledgeResponse represents the response from creating a pledge.
type CreatePledgeResponse struct {
	Pledge         *Pledge `json:"pledge,omitempty"`
	PledgeActionID string  `json:"pledgeActionID"`
}

// UpdatePledgeRequest represents a request to update a pledge.
type UpdatePledgeRequest struct {
	DefaultDestinationSharedAddressID   string `json:"defaultDestinationSharedAddressID,omitempty"`
	DefaultDestinationInternalAddressID string `json:"defaultDestinationInternalAddressID,omitempty"`
}

// AddPledgeCollateralRequest represents a request to add collateral to a pledge.
type AddPledgeCollateralRequest struct {
	Amount string `json:"amount"`
}

// AddPledgeCollateralResponse represents the response from adding collateral.
type AddPledgeCollateralResponse struct {
	PledgeActionID string `json:"pledgeActionID"`
}

// WithdrawPledgeRequest represents a request to withdraw from a pledge.
type WithdrawPledgeRequest struct {
	DestinationSharedAddressID   string `json:"destinationSharedAddressID,omitempty"`
	DestinationInternalAddressID string `json:"destinationInternalAddressID,omitempty"`
	Amount                       string `json:"amount"`
	ExternalReferenceID          string `json:"externalReferenceID,omitempty"`
}

// WithdrawPledgeResponse represents the response from a pledge withdrawal.
type WithdrawPledgeResponse struct {
	PledgeWithdrawalID string `json:"pledgeWithdrawalID"`
	PledgeActionID     string `json:"pledgeActionID"`
}

// InitiateWithdrawPledgeRequest represents a request to initiate a withdrawal.
type InitiateWithdrawPledgeRequest struct {
	DestinationSharedAddressID string `json:"destinationSharedAddressID,omitempty"`
	Amount                     string `json:"amount"`
}

// InitiateWithdrawPledgeResponse represents the response from initiating a withdrawal.
type InitiateWithdrawPledgeResponse struct {
	PledgeWithdrawalID string `json:"pledgeWithdrawalID"`
	PledgeActionID     string `json:"pledgeActionID"`
}

// UnpledgeResponse represents the response from unpledging.
type UnpledgeResponse struct {
	PledgeActionID string `json:"pledgeActionID"`
}

// ListPledgesOptions represents options for listing pledges.
type ListPledgesOptions struct {
	OwnerParticipantID       string
	TargetParticipantID      string
	SharedAddressIDs         []string
	CurrencyID               string
	Statuses                 []string
	SortOrder                string
	CurrentPage              string
	PageRequest              string
	PageSize                 int32
	AttributeFiltersJSON     string
	AttributeFiltersOperator string
}

// ListPledgeActionsOptions represents options for listing pledge actions.
type ListPledgeActionsOptions struct {
	PledgeID    string
	ActionIDs   []string
	SortOrder   string
	CurrentPage string
	PageRequest string
	PageSize    int32
}

// ListPledgeActionsForApprovalOptions represents options for listing pledge actions for approval.
type ListPledgeActionsForApprovalOptions struct {
	ActionIDs   []string
	Types       []string
	SortOrder   string
	CurrentPage string
	PageRequest string
	PageSize    int32
}

// ListPledgeWithdrawalsOptions represents options for listing pledge withdrawals.
type ListPledgeWithdrawalsOptions struct {
	PledgeID         string
	WithdrawalStatus string
	SortOrder        string
	CurrentPage      string
	PageRequest      string
	PageSize         int32
}

// ApprovePledgeActionsRequest represents a request to approve pledge actions.
type ApprovePledgeActionsRequest struct {
	Ids       []string `json:"ids"`
	Signature string   `json:"signature"`
	Comment   string   `json:"comment"`
}

// ApprovePledgeActionsResponse represents the response from approving pledge actions.
type ApprovePledgeActionsResponse struct {
	Signatures string `json:"signatures,omitempty"`
}

// RejectPledgeActionsRequest represents a request to reject pledge actions.
type RejectPledgeActionsRequest struct {
	Ids     []string `json:"ids"`
	Comment string   `json:"comment"`
}

// RejectPledgeRequest represents a request to reject a pledge.
type RejectPledgeRequest struct {
	Comment string `json:"comment,omitempty"`
}
