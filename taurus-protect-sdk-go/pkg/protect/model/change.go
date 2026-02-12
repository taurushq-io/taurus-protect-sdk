package model

import "time"

// ChangeStatus represents the status of a change request.
type ChangeStatus string

const (
	// ChangeStatusCreated indicates the change is pending approval.
	ChangeStatusCreated ChangeStatus = "Created"
	// ChangeStatusApproved indicates the change has been approved.
	ChangeStatusApproved ChangeStatus = "Approved"
	// ChangeStatusRejected indicates the change has been rejected.
	ChangeStatusRejected ChangeStatus = "Rejected"
	// ChangeStatusCanceled indicates the change has been canceled.
	ChangeStatusCanceled ChangeStatus = "Canceled"
)

// ChangeAction represents the type of action being performed on an entity.
type ChangeAction string

const (
	// ChangeActionCreate creates a new entity.
	ChangeActionCreate ChangeAction = "create"
	// ChangeActionUpdate updates an existing entity.
	ChangeActionUpdate ChangeAction = "update"
	// ChangeActionDelete deletes an entity.
	ChangeActionDelete ChangeAction = "delete"
	// ChangeActionResetPassword resets a user's password.
	ChangeActionResetPassword ChangeAction = "resetpassword"
	// ChangeActionResetTOTP resets a user's TOTP configuration.
	ChangeActionResetTOTP ChangeAction = "resettotp"
	// ChangeActionResetKeyContainer resets a user's key container.
	ChangeActionResetKeyContainer ChangeAction = "resetkeycontainer"
	// ChangeActionAssign assigns an entity to another.
	ChangeActionAssign ChangeAction = "assign"
	// ChangeActionUnassign unassigns an entity from another.
	ChangeActionUnassign ChangeAction = "unassign"
)

// ChangeEntity represents the type of entity being changed.
type ChangeEntity string

const (
	ChangeEntityUser                     ChangeEntity = "user"
	ChangeEntityGroup                    ChangeEntity = "group"
	ChangeEntityUserGroup                ChangeEntity = "usergroup"
	ChangeEntityBusinessRule             ChangeEntity = "businessrule"
	ChangeEntityExchange                 ChangeEntity = "exchange"
	ChangeEntityPrice                    ChangeEntity = "price"
	ChangeEntityAction                   ChangeEntity = "action"
	ChangeEntityFeePayer                 ChangeEntity = "feepayer"
	ChangeEntityUserAPIKey               ChangeEntity = "userapikey"
	ChangeEntitySecurityDomain           ChangeEntity = "securitydomain"
	ChangeEntityTaurusNetworkParticipant ChangeEntity = "taurusnetworkparticipant"
	ChangeEntityVisibilityGroup          ChangeEntity = "visibilitygroup"
	ChangeEntityUserVisibilityGroup      ChangeEntity = "uservisibilitygroup"
	ChangeEntityManualAccountFreeze      ChangeEntity = "manualaccountbalancefreeze"
	ChangeEntityManualUTXOFreeze         ChangeEntity = "manualutxofreeze"
	ChangeEntityWallet                   ChangeEntity = "wallet"
	ChangeEntityWhitelistedAddress       ChangeEntity = "whitelistedaddress"
	ChangeEntityAutoTransferEventHandler ChangeEntity = "autotransfereventhandler"
)

// Change represents a change request in the system.
type Change struct {
	// ID is the unique identifier for the change.
	ID string `json:"id"`
	// TenantID is the tenant this change belongs to.
	TenantID string `json:"tenant_id,omitempty"`
	// CreatorID is the ID of the user who created the change.
	CreatorID string `json:"creator_id,omitempty"`
	// CreatorExternalID is the external user ID of the creator.
	CreatorExternalID string `json:"creator_external_id,omitempty"`
	// Action is the type of action being performed (create, update, delete, etc.).
	Action string `json:"action"`
	// EntityID is the ID of the entity being changed.
	EntityID string `json:"entity_id,omitempty"`
	// EntityUUID is the UUID of the entity being changed (for entities that use UUID).
	EntityUUID string `json:"entity_uuid,omitempty"`
	// Entity is the type of entity being changed (user, group, wallet, etc.).
	Entity string `json:"entity"`
	// Changes contains the actual changes to be applied.
	Changes map[string]string `json:"changes,omitempty"`
	// Comment is an optional description of the change.
	Comment string `json:"comment,omitempty"`
	// CreationDate is when the change was created.
	CreationDate time.Time `json:"creation_date"`
}

// ListChangesOptions contains options for listing changes.
type ListChangesOptions struct {
	// Entity filters by entity type (user, group, wallet, etc.).
	Entity string
	// EntityID filters by entity ID (deprecated: use EntityIDs instead).
	EntityID string
	// EntityIDs filters by multiple entity IDs.
	EntityIDs []string
	// EntityUUIDs filters by multiple entity UUIDs.
	EntityUUIDs []string
	// Status filters by change status.
	Status string
	// CreatorID filters by the creator's user ID.
	CreatorID string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// PageSize is the number of items per page.
	PageSize int64
	// CurrentPage is the current page cursor for pagination.
	CurrentPage string
	// PageRequest specifies which page to fetch (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
}

// ListChangesResult contains the result of a ListChanges call.
type ListChangesResult struct {
	// Changes is the list of changes.
	Changes []*Change `json:"changes"`
	// Cursor contains pagination information.
	Cursor *CursorPagination `json:"cursor,omitempty"`
}

// CursorPagination represents cursor-based pagination information.
type CursorPagination struct {
	// CurrentPage is the current page token.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListChangesForApprovalOptions contains options for listing changes pending approval.
type ListChangesForApprovalOptions struct {
	// Entities filters by entity types.
	Entities []string
	// EntityIDs filters by entity IDs (valid when only one entity type is given).
	EntityIDs []string
	// EntityUUIDs filters by entity UUIDs (valid when only one entity type is given).
	EntityUUIDs []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// PageSize is the number of items per page.
	PageSize int64
	// CurrentPage is the current page cursor for pagination.
	CurrentPage string
	// PageRequest specifies which page to fetch (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
}

// CreateChangeRequest contains parameters for creating a change.
type CreateChangeRequest struct {
	// Action is the type of action to perform (required).
	Action string `json:"action"`
	// Entity is the type of entity being changed (required).
	Entity string `json:"entity"`
	// EntityID is the ID of the entity to change.
	EntityID string `json:"entity_id,omitempty"`
	// EntityUUID is the UUID of the entity to change.
	EntityUUID string `json:"entity_uuid,omitempty"`
	// Changes contains the changes to apply (required for create/update actions).
	Changes map[string]string `json:"changes,omitempty"`
	// Comment is an optional description of the change.
	Comment string `json:"comment,omitempty"`
}

// CreateChangeResult contains the result of a CreateChange call.
type CreateChangeResult struct {
	// ID is the unique identifier of the created change.
	ID string `json:"id"`
}
