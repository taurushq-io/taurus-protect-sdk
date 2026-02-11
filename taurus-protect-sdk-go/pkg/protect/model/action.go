package model

import "time"

// Action represents an action envelope in the Taurus-PROTECT system.
// Actions are automated tasks that can be triggered based on certain conditions.
type Action struct {
	// ID is the unique identifier for the action.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Label is a human-readable label for the action.
	Label string `json:"label,omitempty"`
	// Status is the current status of the action.
	Status string `json:"status,omitempty"`
	// CreationDate is when the action was created.
	CreationDate time.Time `json:"creation_date"`
	// UpdateDate is when the action was last updated.
	UpdateDate time.Time `json:"update_date"`
	// LastCheckedDate is when the action was last checked.
	LastCheckedDate time.Time `json:"last_checked_date"`
	// AutoApprove indicates if the action should be auto-approved.
	AutoApprove bool `json:"auto_approve"`
	// ActionDetails contains the trigger and tasks for this action.
	ActionDetails *ActionDetails `json:"action_details,omitempty"`
	// Attributes are custom attributes associated with the action.
	Attributes []*ActionAttribute `json:"attributes,omitempty"`
	// Trails are audit trails for this action.
	Trails []*ActionTrail `json:"trails,omitempty"`
}

// ActionDetails contains the trigger and tasks for an action.
type ActionDetails struct {
	// Trigger defines when the action should be triggered.
	Trigger *ActionTrigger `json:"trigger,omitempty"`
	// Tasks are the tasks to execute when the action is triggered.
	Tasks []*ActionTask `json:"tasks,omitempty"`
}

// ActionTrigger defines the conditions that trigger an action.
type ActionTrigger struct {
	// Kind is the type of trigger.
	Kind string `json:"kind,omitempty"`
	// Balance is the balance trigger configuration.
	Balance *TriggerBalance `json:"balance,omitempty"`
}

// TriggerBalance defines a balance-based trigger.
type TriggerBalance struct {
	// Target specifies what to monitor.
	Target *ActionTarget `json:"target,omitempty"`
	// Comparator specifies how to compare the balance.
	Comparator *ActionComparator `json:"comparator,omitempty"`
	// Amount is the threshold amount for the trigger.
	Amount *ActionAmount `json:"amount,omitempty"`
}

// ActionTarget specifies the target of an action trigger.
type ActionTarget struct {
	// Kind is the type of target.
	Kind string `json:"kind,omitempty"`
	// Address is the address target configuration.
	Address *TargetAddress `json:"address,omitempty"`
	// Wallet is the wallet target configuration.
	Wallet *TargetWallet `json:"wallet,omitempty"`
}

// TargetAddress specifies an address as the target.
type TargetAddress struct {
	// Kind is the type of address target.
	Kind string `json:"kind,omitempty"`
	// AddressID is the ID of the target address.
	AddressID string `json:"address_id,omitempty"`
}

// TargetWallet specifies a wallet as the target.
type TargetWallet struct {
	// Kind is the type of wallet target.
	Kind string `json:"kind,omitempty"`
	// WalletID is the ID of the target wallet.
	WalletID string `json:"wallet_id,omitempty"`
}

// ActionComparator specifies how to compare values.
type ActionComparator struct {
	// Kind is the type of comparator (e.g., "LESS_THAN", "GREATER_THAN").
	Kind string `json:"kind,omitempty"`
}

// ActionAmount specifies an amount for triggers or transfers.
type ActionAmount struct {
	// Kind is the type of amount.
	Kind string `json:"kind,omitempty"`
	// CryptoAmount is the amount in cryptocurrency units.
	CryptoAmount string `json:"crypto_amount,omitempty"`
}

// ActionTask represents a task to execute when an action is triggered.
type ActionTask struct {
	// Kind is the type of task.
	Kind string `json:"kind,omitempty"`
	// Transfer is the transfer task configuration.
	Transfer *TaskTransfer `json:"transfer,omitempty"`
	// Notification is the notification task configuration.
	Notification *TaskNotification `json:"notification,omitempty"`
}

// TaskTransfer defines a transfer task.
type TaskTransfer struct {
	// From specifies the source of the transfer.
	From *ActionSource `json:"from,omitempty"`
	// To specifies the destination of the transfer.
	To *ActionDestination `json:"to,omitempty"`
	// Amount is the amount to transfer.
	Amount *ActionAmount `json:"amount,omitempty"`
	// TopUp indicates if this is a top-up transfer.
	TopUp bool `json:"top_up"`
	// UseAllFunds indicates if all available funds should be transferred.
	UseAllFunds bool `json:"use_all_funds"`
}

// ActionSource specifies the source of a transfer.
type ActionSource struct {
	// Kind is the type of source.
	Kind string `json:"kind,omitempty"`
	// AddressID is the ID of the source address.
	AddressID string `json:"address_id,omitempty"`
	// WalletID is the ID of the source wallet.
	WalletID string `json:"wallet_id,omitempty"`
}

// ActionDestination specifies the destination of a transfer.
type ActionDestination struct {
	// Kind is the type of destination.
	Kind string `json:"kind,omitempty"`
	// AddressID is the ID of the destination address.
	AddressID string `json:"address_id,omitempty"`
	// WhitelistedAddressID is the ID of a whitelisted address.
	WhitelistedAddressID string `json:"whitelisted_address_id,omitempty"`
	// WalletID is the ID of the destination wallet.
	WalletID string `json:"wallet_id,omitempty"`
}

// TaskNotification defines a notification task.
type TaskNotification struct {
	// EmailAddresses are the email addresses to notify.
	EmailAddresses []string `json:"email_addresses,omitempty"`
	// NotificationMessage is the message to send.
	NotificationMessage string `json:"notification_message,omitempty"`
	// NumberOfReminders is the number of reminder notifications to send.
	NumberOfReminders string `json:"number_of_reminders,omitempty"`
}

// ActionAttribute represents a custom attribute on an action.
type ActionAttribute struct {
	// ID is the unique identifier for the attribute.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Key is the attribute key.
	Key string `json:"key,omitempty"`
	// Value is the attribute value.
	Value string `json:"value,omitempty"`
	// ContentType is the MIME type of the value.
	ContentType string `json:"content_type,omitempty"`
}

// ActionTrail represents an audit trail entry for an action.
type ActionTrail struct {
	// ID is the unique identifier for the trail entry.
	ID string `json:"id"`
	// Action is the action that was performed.
	Action string `json:"action,omitempty"`
	// Comment is an optional comment about the action.
	Comment string `json:"comment,omitempty"`
	// Date is when the action occurred.
	Date time.Time `json:"date"`
	// ActionStatus is the status after the action was performed.
	ActionStatus string `json:"action_status,omitempty"`
}

// ListActionsOptions contains options for listing actions.
type ListActionsOptions struct {
	// Limit is the maximum number of actions to return.
	Limit int64
	// Offset is the number of actions to skip.
	Offset int64
	// IDs filters by specific action IDs.
	IDs []string
}

// ListActionsResult contains the result of listing actions.
type ListActionsResult struct {
	// Actions is the list of actions.
	Actions []*Action `json:"actions"`
	// TotalItems is the total number of actions available.
	TotalItems int64 `json:"total_items"`
}
