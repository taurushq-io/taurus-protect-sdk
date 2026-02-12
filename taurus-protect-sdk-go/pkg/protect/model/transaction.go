package model

import "time"

// Transaction represents a blockchain transaction.
type Transaction struct {
	// ID is the unique identifier for the transaction.
	ID string `json:"id"`
	// UniqueID is another unique identifier for the transaction.
	UniqueID string `json:"unique_id,omitempty"`
	// TransactionID is an internal transaction ID.
	TransactionID string `json:"transaction_id,omitempty"`
	// Hash is the blockchain transaction hash.
	Hash string `json:"hash,omitempty"`
	// Direction is the transaction direction (incoming, outgoing).
	Direction string `json:"direction"`
	// Type is the transaction type (transfer, burn, approve, etc.).
	Type string `json:"type,omitempty"`
	// Currency is the currency symbol.
	Currency string `json:"currency"`
	// Blockchain is the blockchain name (ETH, BTC, etc.).
	Blockchain string `json:"blockchain,omitempty"`
	// Amount is the transaction amount in smallest currency unit.
	Amount string `json:"amount"`
	// AmountMainUnit is the amount in main currency unit.
	AmountMainUnit string `json:"amount_main_unit,omitempty"`
	// Fee is the transaction fee in smallest currency unit.
	Fee string `json:"fee,omitempty"`
	// FeeMainUnit is the fee in main currency unit.
	FeeMainUnit string `json:"fee_main_unit,omitempty"`
	// Block is the block number.
	Block string `json:"block,omitempty"`
	// Sources are the sending addresses.
	Sources []TransactionAddressInfo `json:"sources,omitempty"`
	// Destinations are the receiving addresses.
	Destinations []TransactionAddressInfo `json:"destinations,omitempty"`
	// ReceptionDate is when the transaction was received by the blockchain.
	ReceptionDate time.Time `json:"reception_date,omitempty"`
	// ConfirmationDate is when the transaction was fully confirmed.
	ConfirmationDate time.Time `json:"confirmation_date,omitempty"`
}

// TransactionAddressInfo represents an address involved in a transaction.
type TransactionAddressInfo struct {
	// Address is the blockchain address string.
	Address string `json:"address"`
	// Label is the address label if known.
	Label string `json:"label,omitempty"`
	// CustomerID is the customer identifier if set.
	CustomerID string `json:"customer_id,omitempty"`
	// Amount is the transaction amount for this address.
	Amount string `json:"amount,omitempty"`
	// AmountMainUnit is the amount in main currency unit.
	AmountMainUnit string `json:"amount_main_unit,omitempty"`
	// Type is the element type (source or destination).
	Type string `json:"type,omitempty"`
	// InternalAddressID is the ID if this is an internal address.
	InternalAddressID string `json:"internal_address_id,omitempty"`
	// WhitelistedAddressID is the ID if this is a whitelisted address.
	WhitelistedAddressID string `json:"whitelisted_address_id,omitempty"`
}

// ListTransactionsOptions contains options for listing transactions.
type ListTransactionsOptions struct {
	// Limit is the maximum number of transactions to return.
	Limit int64
	// Offset is the number of transactions to skip.
	Offset int64
	// Currency filters by currency symbol.
	Currency string
	// Direction filters by transaction direction (incoming, outgoing).
	Direction string
	// Blockchain filters by blockchain.
	Blockchain string
	// Query searches transaction fields.
	Query string
}

// ListTransactionsByAddressOptions contains options for listing transactions by address.
type ListTransactionsByAddressOptions struct {
	// Limit is the maximum number of transactions to return.
	Limit int64
	// Offset is the number of transactions to skip.
	Offset int64
	// Currency filters by currency symbol.
	Currency string
	// Direction filters by transaction direction (incoming, outgoing).
	Direction string
	// Blockchain filters by blockchain.
	Blockchain string
}

// ExportTransactionsOptions contains options for exporting transactions.
type ExportTransactionsOptions struct {
	// From filters transactions after this date.
	From *time.Time
	// To filters transactions before this date.
	To *time.Time
	// Currency filters by currency ID or symbol.
	Currency string
	// Direction filters by transaction direction (incoming, outgoing).
	Direction string
	// Blockchain filters by blockchain.
	Blockchain string
	// Format specifies the export format (csv, json, csv_simple).
	Format string
	// Limit is the maximum number of transactions to export.
	Limit int64
	// Offset is the number of transactions to skip.
	Offset int64
	// Address filters transactions involving this address.
	Address string
	// Query searches transaction fields.
	Query string
}
