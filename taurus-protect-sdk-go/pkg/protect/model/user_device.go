package model

// UserDevicePairingStatus represents the status of a user device pairing.
type UserDevicePairingStatus string

const (
	// UserDevicePairingStatusWaiting indicates the pairing is waiting to be started.
	UserDevicePairingStatusWaiting UserDevicePairingStatus = "WAITING"
	// UserDevicePairingStatusPairing indicates the pairing is in progress.
	UserDevicePairingStatusPairing UserDevicePairingStatus = "PAIRING"
	// UserDevicePairingStatusApproved indicates the pairing has been approved.
	UserDevicePairingStatusApproved UserDevicePairingStatus = "APPROVED"
)

// CreatePairingResult contains the result of creating a user device pairing request.
// This is Step 1 of the pairing process.
type CreatePairingResult struct {
	// PairingID is the unique identifier for the pairing request.
	// This ID should be used in subsequent steps to complete the pairing process.
	PairingID string `json:"pairing_id"`
}

// StartPairingRequest contains the parameters for starting a user device pairing.
// This is Step 2 of the pairing process.
type StartPairingRequest struct {
	// Nonce is a 6-digit number used to verify the pairing.
	Nonce string `json:"nonce"`
	// PublicKey is the ECDSA public key encoded in base64 format.
	// This key will be used for device authentication.
	PublicKey string `json:"public_key"`
}

// ApprovePairingRequest contains the parameters for approving a user device pairing.
// This is Step 3 of the pairing process.
type ApprovePairingRequest struct {
	// Nonce is a 6-digit number used to verify the pairing.
	// Must match the nonce used in the start pairing step.
	Nonce string `json:"nonce"`
}

// PairingStatusResult contains the result of getting the pairing status.
type PairingStatusResult struct {
	// Status is the current status of the pairing request.
	Status UserDevicePairingStatus `json:"status"`
	// PairingID is the unique identifier for the pairing request.
	PairingID string `json:"pairing_id"`
	// APIKey is the API key for the paired device (only set when status is APPROVED).
	APIKey string `json:"api_key,omitempty"`
}
