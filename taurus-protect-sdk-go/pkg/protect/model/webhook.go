package model

import "time"

// Webhook represents a webhook configuration in Taurus-PROTECT.
type Webhook struct {
	// ID is the unique identifier for the webhook.
	ID string `json:"id"`
	// Type is the webhook type (e.g., "request", "transaction").
	Type string `json:"type"`
	// URL is the endpoint URL where webhook events are sent.
	URL string `json:"url"`
	// Status is the webhook status (e.g., "enabled", "disabled").
	Status string `json:"status"`
	// Secret is the webhook secret (only returned on creation).
	Secret string `json:"secret,omitempty"`
	// TimeoutUntil is when the webhook timeout expires (if timed out).
	TimeoutUntil time.Time `json:"timeout_until,omitempty"`
	// CreatedAt is when the webhook was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the webhook was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateWebhookRequest contains the parameters for creating a webhook.
type CreateWebhookRequest struct {
	// Type is the webhook type (e.g., "request", "transaction").
	Type string `json:"type"`
	// URL is the endpoint URL where webhook events will be sent.
	URL string `json:"url"`
}

// CreateWebhookResult contains the result of creating a webhook.
type CreateWebhookResult struct {
	// Webhook is the created webhook.
	Webhook *Webhook `json:"webhook"`
	// Secret is the webhook secret (only provided once on creation).
	Secret string `json:"secret"`
}

// ListWebhooksOptions contains options for listing webhooks.
type ListWebhooksOptions struct {
	// Type filters webhooks by type.
	Type string
	// URL filters webhooks by URL.
	URL string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortOrder specifies the sort order (ASC or DESC). Default is DESC.
	SortOrder string
}

// ListWebhooksResult contains the result of listing webhooks.
type ListWebhooksResult struct {
	// Webhooks is the list of webhooks.
	Webhooks []*Webhook `json:"webhooks"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
