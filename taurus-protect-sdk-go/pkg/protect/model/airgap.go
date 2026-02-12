package model

// GetOutgoingAirGapRequest contains parameters for exporting HSM ready requests for cold HSM.
type GetOutgoingAirGapRequest struct {
	// RequestIDs is the list of request IDs that have been approved to be signed by the cold HSM.
	RequestIDs []string `json:"request_ids,omitempty"`
	// RequestSignature is the signature of the requests (optional).
	// Format: base64(ecdsa_sign(sha256([hex(sha256(req1_metadata)),hex(sha256(req2_metadata)),...hex(sha256(reqN_metadata))])))
	RequestSignature string `json:"request_signature,omitempty"`
	// AddressIDs is the list of address IDs to be signed by the cold HSM (optional).
	AddressIDs []string `json:"address_ids,omitempty"`
}

// GetOutgoingAirGapResult contains the result of exporting HSM ready requests.
type GetOutgoingAirGapResult struct {
	// Data contains the binary payload to be transmitted to the cold HSM.
	Data []byte `json:"data"`
}

// SubmitIncomingAirGapRequest contains parameters for importing signed requests from cold HSM.
type SubmitIncomingAirGapRequest struct {
	// Payload is the base64-encoded signed payload from the air-gap importer.
	Payload string `json:"payload"`
	// Signature is the base64-encoded signature of the air-gap importer.
	Signature string `json:"signature"`
}
