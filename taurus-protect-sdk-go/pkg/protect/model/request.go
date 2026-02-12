package model

import (
	"encoding/json"
	"strconv"
	"time"
)

// RequestStatus represents the status of a transaction request.
//
// Request statuses track the progress of a request from creation through
// approval, signing, broadcasting, and confirmation on the blockchain.
//
// Common status transitions:
//   - CREATED -> PENDING -> APPROVED -> BROADCASTED -> CONFIRMED
//   - PENDING -> REJECTED (if rejected by an approver)
//   - BROADCASTED -> PERMANENT_FAILURE (if transaction fails)
type RequestStatus string

const (
	// RequestStatusApproved2 indicates request has received secondary approval.
	RequestStatusApproved2 RequestStatus = "APPROVED_2"
	// RequestStatusApproved indicates request has been approved and is ready for signing.
	RequestStatusApproved RequestStatus = "APPROVED"
	// RequestStatusApproving indicates request is in the process of being approved.
	RequestStatusApproving RequestStatus = "APPROVING"
	// RequestStatusAutoPrepared2 is auto prepared secondary status.
	RequestStatusAutoPrepared2 RequestStatus = "AUTO_PREPARED_2"
	// RequestStatusAutoPrepared is auto prepared status.
	RequestStatusAutoPrepared RequestStatus = "AUTO_PREPARED"
	// RequestStatusBroadcasting2 is broadcasting secondary status.
	RequestStatusBroadcasting2 RequestStatus = "BROADCASTING_2"
	// RequestStatusBroadcasting is broadcasting status.
	RequestStatusBroadcasting RequestStatus = "BROADCASTING"
	// RequestStatusBroadcasted indicates broadcasted to the blockchain.
	RequestStatusBroadcasted RequestStatus = "BROADCASTED"
	// RequestStatusBundleApproved indicates bundle has been approved.
	RequestStatusBundleApproved RequestStatus = "BUNDLE_APPROVED"
	// RequestStatusBundleBroadcasting indicates bundle is being broadcast.
	RequestStatusBundleBroadcasting RequestStatus = "BUNDLE_BROADCASTING"
	// RequestStatusBundleReady indicates bundle is ready.
	RequestStatusBundleReady RequestStatus = "BUNDLE_READY"
	// RequestStatusCanceled indicates request was canceled.
	RequestStatusCanceled RequestStatus = "CANCELED"
	// RequestStatusConfirmed indicates transaction confirmed on blockchain.
	RequestStatusConfirmed RequestStatus = "CONFIRMED"
	// RequestStatusCreated indicates request has been created.
	RequestStatusCreated RequestStatus = "CREATED"
	// RequestStatusDiemBurnMbsApproved is Diem burn MBS approved.
	RequestStatusDiemBurnMbsApproved RequestStatus = "DIEM_BURN_MBS_APPROVED"
	// RequestStatusDiemBurnMbsPending is Diem burn MBS pending.
	RequestStatusDiemBurnMbsPending RequestStatus = "DIEM_BURN_MBS_PENDING"
	// RequestStatusDiemMintMbsApproved is Diem mint MBS approved.
	RequestStatusDiemMintMbsApproved RequestStatus = "DIEM_MINT_MBS_APPROVED"
	// RequestStatusDiemMintMbsCompleted is Diem mint MBS completed.
	RequestStatusDiemMintMbsCompleted RequestStatus = "DIEM_MINT_MBS_COMPLETED"
	// RequestStatusDiemMintMbsPending is Diem mint MBS pending.
	RequestStatusDiemMintMbsPending RequestStatus = "DIEM_MINT_MBS_PENDING"
	// RequestStatusExpired indicates request has expired.
	RequestStatusExpired RequestStatus = "EXPIRED"
	// RequestStatusFastApproved2 is fast approved secondary status.
	RequestStatusFastApproved2 RequestStatus = "FAST_APPROVED_2"
	// RequestStatusHsmFailed2 indicates HSM signing failed (secondary).
	RequestStatusHsmFailed2 RequestStatus = "HSM_FAILED_2"
	// RequestStatusHsmFailed indicates HSM signing failed.
	RequestStatusHsmFailed RequestStatus = "HSM_FAILED"
	// RequestStatusHsmReady2 indicates HSM ready for signing (secondary).
	RequestStatusHsmReady2 RequestStatus = "HSM_READY_2"
	// RequestStatusHsmReady indicates HSM ready for signing.
	RequestStatusHsmReady RequestStatus = "HSM_READY"
	// RequestStatusHsmSigned2 indicates HSM signed (secondary).
	RequestStatusHsmSigned2 RequestStatus = "HSM_SIGNED_2"
	// RequestStatusHsmSigned indicates HSM has signed the transaction.
	RequestStatusHsmSigned RequestStatus = "HSM_SIGNED"
	// RequestStatusManualBroadcast indicates manual broadcast required.
	RequestStatusManualBroadcast RequestStatus = "MANUAL_BROADCAST"
	// RequestStatusMined indicates transaction has been mined.
	RequestStatusMined RequestStatus = "MINED"
	// RequestStatusPartiallyConfirmed indicates transaction partially confirmed.
	RequestStatusPartiallyConfirmed RequestStatus = "PARTIALLY_CONFIRMED"
	// RequestStatusPending indicates request is pending approval.
	RequestStatusPending RequestStatus = "PENDING"
	// RequestStatusPermanentFailure indicates transaction permanently failed.
	RequestStatusPermanentFailure RequestStatus = "PERMANENT_FAILURE"
	// RequestStatusReady indicates request is ready.
	RequestStatusReady RequestStatus = "READY"
	// RequestStatusRejected indicates request was rejected.
	RequestStatusRejected RequestStatus = "REJECTED"
	// RequestStatusSent indicates transaction has been sent.
	RequestStatusSent RequestStatus = "SENT"
	// RequestStatusSignetCompleted indicates Signet transaction completed.
	RequestStatusSignetCompleted RequestStatus = "SIGNET_COMPLETED"
	// RequestStatusSignetPending indicates Signet transaction pending.
	RequestStatusSignetPending RequestStatus = "SIGNET_PENDING"
	// RequestStatusUnknown indicates unknown status.
	RequestStatusUnknown RequestStatus = "UNKNOWN"
)

// String returns the string representation of the RequestStatus.
func (s RequestStatus) String() string {
	return string(s)
}

// RequestStatusFromString converts a string to RequestStatus, returning UNKNOWN for unrecognized values.
func RequestStatusFromString(s string) RequestStatus {
	switch s {
	case "APPROVED_2":
		return RequestStatusApproved2
	case "APPROVED":
		return RequestStatusApproved
	case "APPROVING":
		return RequestStatusApproving
	case "AUTO_PREPARED_2":
		return RequestStatusAutoPrepared2
	case "AUTO_PREPARED":
		return RequestStatusAutoPrepared
	case "BROADCASTING_2":
		return RequestStatusBroadcasting2
	case "BROADCASTING":
		return RequestStatusBroadcasting
	case "BROADCASTED":
		return RequestStatusBroadcasted
	case "BUNDLE_APPROVED":
		return RequestStatusBundleApproved
	case "BUNDLE_BROADCASTING":
		return RequestStatusBundleBroadcasting
	case "BUNDLE_READY":
		return RequestStatusBundleReady
	case "CANCELED":
		return RequestStatusCanceled
	case "CONFIRMED":
		return RequestStatusConfirmed
	case "CREATED":
		return RequestStatusCreated
	case "DIEM_BURN_MBS_APPROVED":
		return RequestStatusDiemBurnMbsApproved
	case "DIEM_BURN_MBS_PENDING":
		return RequestStatusDiemBurnMbsPending
	case "DIEM_MINT_MBS_APPROVED":
		return RequestStatusDiemMintMbsApproved
	case "DIEM_MINT_MBS_COMPLETED":
		return RequestStatusDiemMintMbsCompleted
	case "DIEM_MINT_MBS_PENDING":
		return RequestStatusDiemMintMbsPending
	case "EXPIRED":
		return RequestStatusExpired
	case "FAST_APPROVED_2":
		return RequestStatusFastApproved2
	case "HSM_FAILED_2":
		return RequestStatusHsmFailed2
	case "HSM_FAILED":
		return RequestStatusHsmFailed
	case "HSM_READY_2":
		return RequestStatusHsmReady2
	case "HSM_READY":
		return RequestStatusHsmReady
	case "HSM_SIGNED_2":
		return RequestStatusHsmSigned2
	case "HSM_SIGNED":
		return RequestStatusHsmSigned
	case "MANUAL_BROADCAST":
		return RequestStatusManualBroadcast
	case "MINED":
		return RequestStatusMined
	case "PARTIALLY_CONFIRMED":
		return RequestStatusPartiallyConfirmed
	case "PENDING":
		return RequestStatusPending
	case "PERMANENT_FAILURE":
		return RequestStatusPermanentFailure
	case "READY":
		return RequestStatusReady
	case "REJECTED":
		return RequestStatusRejected
	case "SENT":
		return RequestStatusSent
	case "SIGNET_COMPLETED":
		return RequestStatusSignetCompleted
	case "SIGNET_PENDING":
		return RequestStatusSignetPending
	default:
		return RequestStatusUnknown
	}
}

// SignedRequest represents a signed blockchain transaction associated with a request.
type SignedRequest struct {
	// ID is the unique identifier for the signed request.
	ID string `json:"id"`
	// SignedRequest is the signed transaction data.
	SignedRequest string `json:"signed_request,omitempty"`
	// Status is the status of the signed request.
	Status RequestStatus `json:"status"`
	// Hash is the blockchain transaction hash.
	Hash string `json:"hash,omitempty"`
	// Block is the block number.
	Block int64 `json:"block,omitempty"`
	// Details contains additional details.
	Details string `json:"details,omitempty"`
	// CreationDate is when the signed request was created.
	CreationDate time.Time `json:"creation_date"`
	// UpdateDate is when the signed request was last updated.
	UpdateDate time.Time `json:"update_date"`
	// BroadcastDate is when the signed request was broadcast.
	BroadcastDate time.Time `json:"broadcast_date"`
	// ConfirmationDate is when the signed request was confirmed.
	ConfirmationDate time.Time `json:"confirmation_date"`
}

// Request represents a transaction request (withdrawal, transfer, etc.).
type Request struct {
	// ID is the unique identifier for the request.
	ID string `json:"id"`
	// TenantID is the ID of the tenant where the request was created.
	TenantID string `json:"tenant_id,omitempty"`
	// Currency is the currency symbol.
	Currency string `json:"currency"`
	// Status is the request status (CREATED, APPROVING, HSM_SIGNED, BROADCASTING, BROADCASTED, CONFIRMED, REJECTED).
	Status string `json:"status"`
	// Type is the request type (transfer, staking, unstaking, etc.).
	Type string `json:"type,omitempty"`
	// Rule is the rule applied to this request.
	Rule string `json:"rule,omitempty"`
	// RequestBundleID is the ID of the request bundle if applicable.
	RequestBundleID string `json:"request_bundle_id,omitempty"`
	// ExternalRequestID is an optional external identifier.
	ExternalRequestID string `json:"external_request_id,omitempty"`
	// NeedsApprovalFrom lists groups that need to approve the request.
	NeedsApprovalFrom []string `json:"needs_approval_from,omitempty"`
	// Metadata contains request metadata.
	Metadata *RequestMetadata `json:"metadata,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []RequestAttribute `json:"attributes,omitempty"`
	// SignedRequests contains signed blockchain transactions associated with this request.
	SignedRequests []SignedRequest `json:"signed_requests,omitempty"`
	// CreatedAt is when the request was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the request was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// RequestMetadata contains metadata about a request.
type RequestMetadata struct {
	// Hash is the metadata hash.
	Hash string `json:"hash,omitempty"`
	// SECURITY: Payload field intentionally omitted - use PayloadAsString only.
	// See ParsePayloadEntries() for secure data extraction from verified source.
	// PayloadAsString is the payload serialized as a string.
	PayloadAsString string `json:"payload_as_string,omitempty"`
}

// RequestAttribute represents a custom attribute on a request.
type RequestAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
}

// CreateOutgoingRequest contains parameters for creating an outgoing (withdrawal) request.
type CreateOutgoingRequest struct {
	// Amount is the transaction amount in smallest currency unit (required).
	Amount string `json:"amount"`
	// FromAddressID is the source address ID (either FromAddressID or FromWalletID required).
	FromAddressID string `json:"from_address_id,omitempty"`
	// FromWalletID is the source wallet ID (for omnibus wallets).
	FromWalletID string `json:"from_wallet_id,omitempty"`
	// ToAddressID is the destination address ID (either ToAddressID or ToWhitelistedAddressID required).
	ToAddressID string `json:"to_address_id,omitempty"`
	// ToWhitelistedAddressID is the destination whitelisted address ID.
	ToWhitelistedAddressID string `json:"to_whitelisted_address_id,omitempty"`
	// FeeLimit is the maximum fee amount.
	FeeLimit string `json:"fee_limit,omitempty"`
	// GasLimit is the maximum gas for the transaction.
	GasLimit string `json:"gas_limit,omitempty"`
	// Comment is a reconciliation note.
	Comment string `json:"comment,omitempty"`
	// TransactionComment is an external reference (max 256 chars).
	TransactionComment string `json:"transaction_comment,omitempty"`
	// ExternalRequestID is an optional external identifier.
	ExternalRequestID string `json:"external_request_id,omitempty"`
	// UseUnconfirmedFunds allows use of unconfirmed funds.
	UseUnconfirmedFunds bool `json:"use_unconfirmed_funds,omitempty"`
	// FeePaidByReceiver indicates if fees are deducted from amount.
	FeePaidByReceiver bool `json:"fee_paid_by_receiver,omitempty"`
	// UseAllFunds sends all funds from the address/wallet.
	UseAllFunds bool `json:"use_all_funds,omitempty"`
}

// ListRequestsOptions contains options for listing requests.
type ListRequestsOptions struct {
	// PageSize is the maximum number of requests per page.
	PageSize int
	// Cursor is the base64-encoded cursor for the current page (from RequestResult.NextCursor).
	Cursor string
	// Status filters by request status.
	Status string
	// Currency filters by currency.
	Currency string
}

// RequestResult contains the result of a paginated request list query.
type RequestResult struct {
	// Requests is the list of requests in the current page.
	Requests []*Request
	// NextCursor is the cursor to use for fetching the next page. Empty if no more pages.
	NextCursor string
	// HasNext indicates whether more pages are available.
	HasNext bool
}

// CreateIncomingRequest contains parameters for creating an incoming request from an exchange.
type CreateIncomingRequest struct {
	// FromExchangeID is the source exchange ID (required).
	FromExchangeID string `json:"from_exchange_id"`
	// ToAddressID is the destination address ID (required).
	ToAddressID string `json:"to_address_id"`
	// Amount is the transaction amount (required).
	Amount string `json:"amount"`
	// Comment is an optional reconciliation note.
	Comment string `json:"comment,omitempty"`
	// ExternalRequestID is an optional external identifier.
	ExternalRequestID string `json:"external_request_id,omitempty"`
}

// PayloadEntry represents a single key-value entry in the metadata payload array.
// The API returns payload as an array of these entries:
//
//	[{"key": "source", "value": {...}}, {"key": "destination", "value": {...}}, ...]
type PayloadEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// PayloadEntryValue provides typed access to nested payload values.
type PayloadEntryValue struct {
	raw interface{}
}

// NewPayloadEntryValue wraps a raw value for typed access.
func NewPayloadEntryValue(v interface{}) *PayloadEntryValue {
	if v == nil {
		return nil
	}
	return &PayloadEntryValue{raw: v}
}

// AsMap returns the value as a map if possible, nil otherwise.
func (v *PayloadEntryValue) AsMap() map[string]interface{} {
	if v == nil || v.raw == nil {
		return nil
	}
	if m, ok := v.raw.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// GetString retrieves a nested string value by key path.
func (v *PayloadEntryValue) GetString(keys ...string) string {
	if v == nil {
		return ""
	}
	current := v.raw
	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = m[key]
	}
	if s, ok := current.(string); ok {
		return s
	}
	return ""
}

// GetInt64 retrieves a nested int64 value by key path.
func (v *PayloadEntryValue) GetInt64(keys ...string) int64 {
	if v == nil {
		return 0
	}
	current := v.raw
	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return 0
		}
		current = m[key]
	}
	// JSON numbers unmarshal as float64
	if f, ok := current.(float64); ok {
		return int64(f)
	}
	return 0
}

// GetFloat64 retrieves a nested float64 value by key path.
func (v *PayloadEntryValue) GetFloat64(keys ...string) float64 {
	if v == nil {
		return 0
	}
	current := v.raw
	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return 0
		}
		current = m[key]
	}
	if f, ok := current.(float64); ok {
		return f
	}
	return 0
}

// ParsePayloadEntries parses the PayloadAsString into structured entries.
// The API returns metadata payload as a JSON array:
//
//	[{"key": "source", "value": {...}}, {"key": "destination", "value": {...}}, ...]
//
// This method parses that array into PayloadEntry structs for easy access.
// Returns nil if PayloadAsString is empty or not a valid JSON array.
//
// SECURITY NOTE:
// This method parses from PayloadAsString, which is the cryptographically
// verified source. The metadata.hash field is a SHA-256 hash of PayloadAsString,
// and this hash is signed by governance rules (SuperAdmin keys). By extracting
// data only from PayloadAsString (not from the raw Payload field), we ensure that:
//   - All extracted data has been integrity-verified
//   - Any tampering with the raw payload object is ignored
//   - The verification chain is preserved: payloadAsString → hash → signature
//
// ATTACK VECTOR (if using raw Payload):
// An attacker intercepting API responses could modify the payload object
// (e.g., change a destination address) while leaving payloadAsString unchanged.
// The hash would still verify, but the client would extract tampered data.
func (m *RequestMetadata) ParsePayloadEntries() ([]PayloadEntry, error) {
	if m == nil || m.PayloadAsString == "" {
		return nil, nil
	}

	var entries []PayloadEntry
	if err := json.Unmarshal([]byte(m.PayloadAsString), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// GetPayloadValue retrieves a value from the payload by key.
// Returns nil if the key is not found or if parsing fails.
//
// SECURITY NOTE: This method uses ParsePayloadEntries() internally, which
// extracts data from PayloadAsString (the cryptographically verified source).
// See ParsePayloadEntries() documentation for security rationale.
func (m *RequestMetadata) GetPayloadValue(key string) *PayloadEntryValue {
	entries, err := m.ParsePayloadEntries()
	if err != nil || entries == nil {
		return nil
	}
	for _, entry := range entries {
		if entry.Key == key {
			return NewPayloadEntryValue(entry.Value)
		}
	}
	return nil
}

// GetSourceAddress extracts the source blockchain address from the metadata payload.
// Returns an empty string if the source address is not found.
func (m *RequestMetadata) GetSourceAddress() string {
	val := m.GetPayloadValue("source")
	if val == nil {
		return ""
	}
	return val.GetString("payload", "address")
}

// GetDestinationAddress extracts the destination blockchain address from the metadata payload.
// Returns an empty string if the destination address is not found.
func (m *RequestMetadata) GetDestinationAddress() string {
	val := m.GetPayloadValue("destination")
	if val == nil {
		return ""
	}
	return val.GetString("payload", "address")
}

// GetMetadataCurrency extracts the currency from the metadata payload.
// Returns an empty string if the currency is not found.
func (m *RequestMetadata) GetMetadataCurrency() string {
	val := m.GetPayloadValue("currency")
	if val == nil {
		return ""
	}
	// Currency value is directly a string, not a nested object
	if s, ok := val.raw.(string); ok {
		return s
	}
	return ""
}

// GetMetadataRequestID extracts the request ID from the metadata payload.
// Returns 0 if the request ID is not found.
func (m *RequestMetadata) GetMetadataRequestID() int64 {
	val := m.GetPayloadValue("request_id")
	if val == nil {
		return 0
	}
	// Request ID is directly a number, not a nested object
	if f, ok := val.raw.(float64); ok {
		return int64(f)
	}
	return 0
}

// RequestMetadataAmount contains amount details from the metadata payload.
type RequestMetadataAmount struct {
	// ValueFrom is the amount in the source currency (smallest unit).
	// Stored as string to preserve arbitrary-precision values that may exceed int64.
	ValueFrom string
	// ValueTo is the amount in the target currency.
	// Stored as string to preserve arbitrary-precision values.
	ValueTo string
	// Rate is the conversion rate applied.
	// Stored as string to preserve precision.
	Rate string
	// Decimals is the decimal precision.
	Decimals int
	// CurrencyFrom is the source currency code.
	CurrencyFrom string
	// CurrencyTo is the target currency code.
	CurrencyTo string
}

// jsonValueToString converts a JSON value to string, handling both string and float64 types.
// The API returns amount fields as strings (to support arbitrary-precision values),
// but we also handle float64 for backward compatibility with any JSON using numbers.
func jsonValueToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	}
	return ""
}

// GetAmount extracts the amount details from the metadata payload.
// Returns nil if the amount is not found.
func (m *RequestMetadata) GetAmount() *RequestMetadataAmount {
	val := m.GetPayloadValue("amount")
	if val == nil {
		return nil
	}
	v := val.AsMap()
	if v == nil {
		return nil
	}

	amount := &RequestMetadataAmount{}

	amount.ValueFrom = jsonValueToString(v["valueFrom"])
	amount.ValueTo = jsonValueToString(v["valueTo"])
	amount.Rate = jsonValueToString(v["rate"])

	// decimals: handle both string and float64 representations
	if s, ok := v["decimals"].(string); ok {
		amount.Decimals, _ = strconv.Atoi(s)
	} else if f, ok := v["decimals"].(float64); ok {
		amount.Decimals = int(f)
	}

	if s, ok := v["currencyFrom"].(string); ok {
		amount.CurrencyFrom = s
	}
	if s, ok := v["currencyTo"].(string); ok {
		amount.CurrencyTo = s
	}

	return amount
}
