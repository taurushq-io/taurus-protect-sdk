package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AddressFromDTO converts an OpenAPI Address to a domain Address.
func AddressFromDTO(dto *openapi.TgvalidatordAddress) *model.Address {
	if dto == nil {
		return nil
	}

	address := &model.Address{
		ID:                safeString(dto.Id),
		WalletID:          safeString(dto.WalletId),
		Address:           safeString(dto.Address),
		AlternateAddress:  safeString(dto.AlternateAddress),
		Label:             safeString(dto.Label),
		Comment:           safeString(dto.Comment),
		Currency:          safeString(dto.Currency),
		CustomerID:        safeString(dto.CustomerId),
		ExternalAddressID: safeString(dto.ExternalAddressId),
		AddressPath:       safeString(dto.AddressPath),
		AddressIndex:      safeString(dto.AddressIndex),
		Nonce:             safeString(dto.Nonce),
		Status:            safeString(dto.Status),
		Signature:         safeString(dto.Signature),
		Disabled:          safeBool(dto.Disabled),
		CanUseAllFunds:    safeBool(dto.CanUseAllFunds),
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		address.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		address.UpdatedAt = *dto.UpdateDate
	}

	// Convert balance
	if dto.Balance != nil {
		address.Balance = BalanceFromDTO(dto.Balance)
	}

	// Convert attributes
	if dto.Attributes != nil {
		address.Attributes = make([]model.AddressAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			address.Attributes[i] = AddressAttributeFromDTO(&attr)
		}
	}

	// Copy linked whitelisted address IDs
	if dto.LinkedWhitelistedAddressIds != nil {
		address.LinkedWhitelistedAddressIDs = make([]string, len(dto.LinkedWhitelistedAddressIds))
		copy(address.LinkedWhitelistedAddressIDs, dto.LinkedWhitelistedAddressIds)
	}

	return address
}

// AddressesFromDTO converts a slice of OpenAPI Address to domain Addresses.
func AddressesFromDTO(dtos []openapi.TgvalidatordAddress) []*model.Address {
	if dtos == nil {
		return nil
	}
	addresses := make([]*model.Address, len(dtos))
	for i := range dtos {
		addresses[i] = AddressFromDTO(&dtos[i])
	}
	return addresses
}

// AddressAttributeFromDTO converts an OpenAPI AddressAttribute to a domain AddressAttribute.
func AddressAttributeFromDTO(dto *openapi.TgvalidatordAddressAttribute) model.AddressAttribute {
	if dto == nil {
		return model.AddressAttribute{}
	}
	return model.AddressAttribute{
		ID:    safeString(dto.Id),
		Key:   safeString(dto.Key),
		Value: safeString(dto.Value),
	}
}

// AddressAttributesFromDTO converts a slice of OpenAPI AddressAttribute to domain AddressAttributes.
func AddressAttributesFromDTO(dtos []openapi.TgvalidatordAddressAttribute) []model.AddressAttribute {
	if dtos == nil {
		return nil
	}
	attrs := make([]model.AddressAttribute, len(dtos))
	for i := range dtos {
		attrs[i] = AddressAttributeFromDTO(&dtos[i])
	}
	return attrs
}

// ProofOfReserveFromDTO converts an OpenAPI ProofOfReserve to a domain ProofOfReserve.
func ProofOfReserveFromDTO(dto *openapi.TgvalidatordProofOfReserve) *model.ProofOfReserve {
	if dto == nil {
		return nil
	}

	result := &model.ProofOfReserve{
		Path:                   safeString(dto.Path),
		Address:                safeString(dto.Address),
		PublicKey:              safeString(dto.PublicKey),
		Challenge:              safeString(dto.Challenge),
		ChallengeResponse:      safeString(dto.ChallengeResponse),
		StakePublicKey:         safeString(dto.StakePublicKey),
		StakeChallengeResponse: safeString(dto.StakeChallengeResponse),
	}

	// Convert curve enum to string
	if dto.Curve != nil {
		result.Curve = string(*dto.Curve)
	}

	// Convert cipher enum to string
	if dto.Cipher != nil {
		result.Cipher = string(*dto.Cipher)
	}

	// Convert type enum to string
	if dto.Type != nil {
		result.Type = string(*dto.Type)
	}

	return result
}
