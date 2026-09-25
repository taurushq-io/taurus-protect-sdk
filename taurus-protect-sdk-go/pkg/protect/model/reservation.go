package model

import "time"

// Reservation represents a fund reservation in the system.
type Reservation struct {
	// ID is the unique identifier for the reservation.
	ID string `json:"id"`
	// Kind is the type of reservation.
	Kind string `json:"kind,omitempty"`
	// Comment is an optional comment about the reservation.
	Comment string `json:"comment,omitempty"`
	// AddressID references the internal address stored in protect.
	AddressID string `json:"address_id,omitempty"`
	// Address is the blockchain address string/hash.
	Address string `json:"address,omitempty"`
	// Amount is the reserved amount.
	Amount string `json:"amount,omitempty"`
	// CreationDate is when the reservation was created.
	CreationDate time.Time `json:"creation_date,omitempty"`
	// CurrencyInfo contains information about the currency.
	CurrencyInfo *CurrencyInfo `json:"currency_info,omitempty"`
	// Asset contains information about the asset.
	Asset *Asset `json:"asset,omitempty"`
	// ResourceID is the resource ID associated with the reservation.
	ResourceID string `json:"resource_id,omitempty"`
	// ResourceType is the type of resource associated with the reservation.
	ResourceType string `json:"resource_type,omitempty"`
}

// ReservationUTXO represents a UTXO associated with a reservation.
type ReservationUTXO struct {
	// ID is the unique identifier for the UTXO.
	ID string `json:"id"`
	// Hash is the transaction hash.
	Hash string `json:"hash,omitempty"`
	// OutputIndex is the index of the output in the transaction.
	OutputIndex int64 `json:"output_index,omitempty"`
	// Script is the locking script.
	Script string `json:"script,omitempty"`
	// Value is the UTXO value in satoshis or smallest unit.
	Value string `json:"value,omitempty"`
	// ValueString is the human-readable value.
	ValueString string `json:"value_string,omitempty"`
	// BlockHeight is the block height where the UTXO was created.
	BlockHeight string `json:"block_height,omitempty"`
	// ReservationID is the ID of the reservation associated with this UTXO.
	ReservationID string `json:"reservation_id,omitempty"`
}

// ListReservationsOptions contains options for listing reservations.
type ListReservationsOptions struct {
	// Kinds filters by reservation types (e.g., ["PENDING_REQUEST", "MINIMUM_BALANCE"]).
	// Takes precedence over Kind (deprecated).
	Kinds []string
	// Address filters by blockchain address.
	Address string
	// AddressID filters by internal address ID.
	AddressID string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
	// CurrentPage and PageRequest (FIRST, PREVIOUS, NEXT, LAST) page by hand; neither can be
	// combined with Cursor.
	CurrentPage string
	PageRequest string
}

// ListReservationsResult contains the result of listing reservations.
type ListReservationsResult struct {
	// Reservations is the list of reservations.
	Reservations []*Reservation `json:"reservations"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}
