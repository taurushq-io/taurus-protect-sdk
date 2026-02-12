package model

import "time"

// WebhookCall represents a webhook call record.
type WebhookCall struct {
	// ID is the unique identifier for the webhook call.
	ID string `json:"id"`
	// EventID is the ID of the event that triggered this webhook call.
	EventID string `json:"event_id"`
	// WebhookID is the ID of the webhook that was called.
	WebhookID string `json:"webhook_id"`
	// Payload is the JSON payload sent to the webhook.
	Payload string `json:"payload,omitempty"`
	// Status is the status of the webhook call (e.g., "SUCCESS", "FAILED", "PENDING").
	Status string `json:"status"`
	// StatusMessage provides additional details about the status.
	StatusMessage string `json:"status_message,omitempty"`
	// Attempts is the number of attempts made to deliver this webhook.
	Attempts int64 `json:"attempts"`
	// CreatedAt is when the webhook call was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the webhook call was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// ListWebhookCallsOptions contains options for listing webhook calls.
type ListWebhookCallsOptions struct {
	// EventID filters by the event ID that triggered the webhook calls.
	EventID string
	// WebhookID filters by the webhook ID.
	WebhookID string
	// Status filters by the webhook call status (e.g., "SUCCESS", "FAILED", "PENDING").
	Status string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortOrder specifies the sort order (ASC or DESC). Default is DESC.
	SortOrder string
}

// ListWebhookCallsResult contains the result of listing webhook calls.
type ListWebhookCallsResult struct {
	// WebhookCalls is the list of webhook call records.
	WebhookCalls []*WebhookCall `json:"webhook_calls"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
