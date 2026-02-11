package helper

import (
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// VerifyAddressSignature verifies the address signature using the HSM slot public key.
// The signed data is the blockchain address string. The public key is found by locating
// a user with the HSMSLOT role in the rules container.
//
// Returns nil if verification succeeds, IntegrityError otherwise.
func VerifyAddressSignature(address *model.Address, rulesContainer *model.DecodedRulesContainer) error {
	if address == nil {
		return &model.IntegrityError{Message: "address cannot be nil"}
	}
	if rulesContainer == nil {
		return &model.IntegrityError{Message: "rulesContainer cannot be nil"}
	}

	// Check for signature presence (mandatory)
	if address.Signature == "" {
		return &model.IntegrityError{
			Message: fmt.Sprintf("Address %s has no signature", address.ID),
		}
	}

	// Check for address string presence
	if address.Address == "" {
		return &model.IntegrityError{
			Message: fmt.Sprintf("Address %s has no blockchain address to verify", address.ID),
		}
	}

	// Get the HSM public key (cached in rules container)
	hsmPublicKey := rulesContainer.GetHsmPublicKey()
	if hsmPublicKey == nil {
		return &model.IntegrityError{
			Message: "No user with HSMSLOT role found in rules container",
		}
	}

	// Verify the signature
	addressData := []byte(address.Address)
	valid, err := crypto.VerifySignature(hsmPublicKey, addressData, address.Signature)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("Address signature verification failed for address %s: %v", address.ID, err),
		}
	}

	if !valid {
		return &model.IntegrityError{
			Message: fmt.Sprintf("Address signature verification failed for address %s", address.ID),
		}
	}

	return nil
}

// VerifyAddressSignatures verifies signatures for multiple addresses.
// Returns the first verification error encountered, or nil if all signatures are valid.
func VerifyAddressSignatures(addresses []*model.Address, rulesContainer *model.DecodedRulesContainer) error {
	for _, address := range addresses {
		if err := VerifyAddressSignature(address, rulesContainer); err != nil {
			return err
		}
	}
	return nil
}
