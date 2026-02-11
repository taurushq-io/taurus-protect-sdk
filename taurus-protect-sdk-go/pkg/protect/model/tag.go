package model

import "time"

// Tag represents a tag that can be assigned to entities.
type Tag struct {
	// ID is the unique identifier for the tag.
	ID string `json:"id"`
	// Value is the tag value/name.
	Value string `json:"value"`
	// Color is the hex color code for the tag.
	Color string `json:"color,omitempty"`
	// CreationDate is when the tag was created.
	CreationDate time.Time `json:"creation_date"`
}

// CreateTagRequest contains parameters for creating a tag.
type CreateTagRequest struct {
	// Value is the tag value/name (required).
	Value string `json:"value"`
	// Color is the hex color code for the tag (required).
	Color string `json:"color"`
}

// ListTagsOptions contains options for listing tags.
type ListTagsOptions struct {
	// IDs filters tags by their IDs.
	IDs []string
	// Query filters tags where the value contains this substring.
	Query string
}
