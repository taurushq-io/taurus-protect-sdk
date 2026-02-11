package openapi

import (
	"encoding/json"
)

// UnmarshalJSON provides a custom JSON unmarshaler for TgvalidatordMetadata
// that handles the API's array-based payload format.
//
// SECURITY DESIGN:
// ================
// The API returns metadata with two representations of the same data:
//   - payload: Raw JSON object/array (UNVERIFIED)
//   - payloadAsString: JSON string that is cryptographically hashed (VERIFIED)
//
// The security model works as follows:
//  1. The server computes: metadata.hash = SHA256(payloadAsString)
//  2. The hash is signed by governance rules (SuperAdmin keys)
//  3. Clients verify: computed_hash(payloadAsString) == metadata.hash
//
// ATTACK VECTOR IF USING PAYLOAD DIRECTLY:
// An attacker intercepting API responses could:
//  1. Modify the payload object (e.g., change destination address)
//  2. Leave payloadAsString unchanged (hash still verifies)
//  3. Client extracts data from modified payload → SECURITY BYPASS
//
// SOLUTION:
// We intentionally set Payload to nil and force all data extraction
// through PayloadAsString. This ensures:
//   - All extracted data comes from the cryptographically verified source
//   - Any tampering with the raw payload is ignored
//   - The integrity chain: payloadAsString → hash → signature is preserved
func (o *TgvalidatordMetadata) UnmarshalJSON(data []byte) error {
	type metadataAlias struct {
		Hash            *string `json:"hash,omitempty"`
		PayloadAsString *string `json:"payloadAsString,omitempty"`
	}

	var alias metadataAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	o.Hash = alias.Hash
	o.PayloadAsString = alias.PayloadAsString

	// SECURITY: Payload is intentionally left nil.
	// All data extraction MUST use PayloadAsString (the verified source).
	// See function documentation for detailed security rationale.
	o.Payload = nil

	return nil
}
