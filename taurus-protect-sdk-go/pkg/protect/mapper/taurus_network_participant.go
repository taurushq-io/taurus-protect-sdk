package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TnParticipantFromDTO converts an OpenAPI TnParticipant to a domain TnParticipant.
func TnParticipantFromDTO(dto *openapi.TgvalidatordTnParticipant) *taurusnetwork.TnParticipant {
	if dto == nil {
		return nil
	}

	participant := &taurusnetwork.TnParticipant{
		ID:                                        safeString(dto.Id),
		Name:                                      safeString(dto.Name),
		LegalAddress:                              safeString(dto.LegalAddress),
		Country:                                   safeString(dto.Country),
		LogoBase64:                                safeString(dto.LogoBase64),
		PublicKey:                                 safeString(dto.PublicKey),
		Shield:                                    safeString(dto.Shield),
		OwnedSharedAddressesCount:                 safeString(dto.OwnedSharedAddressesCount),
		TargetedSharedAddressesCount:              safeString(dto.TargetedSharedAddressesCount),
		OutgoingTotalPledgesValuationBaseCurrency: safeString(dto.OutgoingTotalPledgesValuationBaseCurrency),
		IncomingTotalPledgesValuationBaseCurrency: safeString(dto.IncomingTotalPledgesValuationBaseCurrency),
		PublicSubname:                             safeString(dto.PublicSubname),
		LegalEntityIdentifier:                     safeString(dto.LegalEntityIdentifier),
		Status:                                    safeString(dto.Status),
	}

	// Convert timestamps
	if dto.OriginRegistrationDate != nil {
		participant.OriginRegistrationDate = *dto.OriginRegistrationDate
	}
	if dto.OriginDeletionDate != nil {
		participant.OriginDeletionDate = *dto.OriginDeletionDate
	}
	if dto.CreatedAt != nil {
		participant.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		participant.UpdatedAt = *dto.UpdatedAt
	}

	// Convert details
	if dto.Details != nil {
		participant.Details = TnParticipantDetailsFromDTO(dto.Details)
	}

	// Convert block confirmations
	if dto.BlockConfirmations != nil {
		participant.BlockConfirmations = make([]taurusnetwork.TnBlockConfirmations, len(dto.BlockConfirmations))
		for i := range dto.BlockConfirmations {
			participant.BlockConfirmations[i] = TnBlockConfirmationsFromDTO(&dto.BlockConfirmations[i])
		}
	}

	// Convert attributes
	if dto.Attributes != nil {
		participant.Attributes = make([]taurusnetwork.TnParticipantAttribute, len(dto.Attributes))
		for i := range dto.Attributes {
			participant.Attributes[i] = TnParticipantAttributeFromDTO(&dto.Attributes[i])
		}
	}

	return participant
}

// TnParticipantsFromDTO converts a slice of OpenAPI TnParticipant to domain TnParticipants.
func TnParticipantsFromDTO(dtos []openapi.TgvalidatordTnParticipant) []*taurusnetwork.TnParticipant {
	if dtos == nil {
		return nil
	}
	participants := make([]*taurusnetwork.TnParticipant, len(dtos))
	for i := range dtos {
		participants[i] = TnParticipantFromDTO(&dtos[i])
	}
	return participants
}

// TnParticipantDetailsFromDTO converts an OpenAPI TnParticipantDetails to a domain TnParticipantDetails.
func TnParticipantDetailsFromDTO(dto *openapi.TgvalidatordTnParticipantDetails) *taurusnetwork.TnParticipantDetails {
	if dto == nil {
		return nil
	}

	details := &taurusnetwork.TnParticipantDetails{}

	// Convert contact persons
	if dto.ContactPersons != nil {
		details.ContactPersons = make([]taurusnetwork.TnContactPerson, len(dto.ContactPersons))
		for i := range dto.ContactPersons {
			details.ContactPersons[i] = TnContactPersonFromDTO(&dto.ContactPersons[i])
		}
	}

	// Convert attribute specifications
	if dto.AttributesSpecifications != nil {
		details.AttributesSpecifications = make([]taurusnetwork.TnParticipantAttributeSpecification, len(dto.AttributesSpecifications))
		for i := range dto.AttributesSpecifications {
			details.AttributesSpecifications[i] = TnParticipantAttributeSpecificationFromDTO(&dto.AttributesSpecifications[i])
		}
	}

	// Convert supported blockchains
	if dto.SupportedBlockchains != nil {
		details.SupportedBlockchains = make([]taurusnetwork.TnBlockchainEntity, len(dto.SupportedBlockchains))
		for i := range dto.SupportedBlockchains {
			details.SupportedBlockchains[i] = TnBlockchainEntityFromDTO(&dto.SupportedBlockchains[i])
		}
	}

	return details
}

// TnContactPersonFromDTO converts an OpenAPI TnContactPerson to a domain TnContactPerson.
func TnContactPersonFromDTO(dto *openapi.TgvalidatordTnContactPerson) taurusnetwork.TnContactPerson {
	if dto == nil {
		return taurusnetwork.TnContactPerson{}
	}
	return taurusnetwork.TnContactPerson{
		FirstName:   safeString(dto.FirstName),
		LastName:    safeString(dto.LastName),
		PhoneNumber: safeString(dto.PhoneNumber),
		Email:       safeString(dto.Email),
	}
}

// TnParticipantAttributeSpecificationFromDTO converts an OpenAPI TnParticipantAttributeSpecification to a domain TnParticipantAttributeSpecification.
func TnParticipantAttributeSpecificationFromDTO(dto *openapi.TgvalidatordTnParticipantAttributeSpecification) taurusnetwork.TnParticipantAttributeSpecification {
	if dto == nil {
		return taurusnetwork.TnParticipantAttributeSpecification{}
	}
	return taurusnetwork.TnParticipantAttributeSpecification{
		AttributeKey:         safeString(dto.AttributeKey),
		AttributeType:        safeString(dto.AttributeType),
		AttributeDescription: safeString(dto.AttributeDescription),
	}
}

// TnBlockchainEntityFromDTO converts an OpenAPI BlockchainEntity to a domain TnBlockchainEntity.
func TnBlockchainEntityFromDTO(dto *openapi.TgvalidatordBlockchainEntity) taurusnetwork.TnBlockchainEntity {
	if dto == nil {
		return taurusnetwork.TnBlockchainEntity{}
	}
	return taurusnetwork.TnBlockchainEntity{
		Symbol:  safeString(dto.Symbol),
		Name:    safeString(dto.Name),
		Network: safeString(dto.Network),
	}
}

// TnBlockConfirmationsFromDTO converts an OpenAPI BlockConfirmations to a domain TnBlockConfirmations.
func TnBlockConfirmationsFromDTO(dto *openapi.TgvalidatordBlockConfirmations) taurusnetwork.TnBlockConfirmations {
	if dto == nil {
		return taurusnetwork.TnBlockConfirmations{}
	}
	return taurusnetwork.TnBlockConfirmations{
		Blockchain:             safeString(dto.Blockchain),
		Network:                safeString(dto.Network),
		ConfirmationsThreshold: safeString(dto.ConfirmationsThreshold),
	}
}

// TnParticipantAttributeFromDTO converts an OpenAPI TnParticipantAttribute to a domain TnParticipantAttribute.
func TnParticipantAttributeFromDTO(dto *openapi.TgvalidatordTnParticipantAttribute) taurusnetwork.TnParticipantAttribute {
	if dto == nil {
		return taurusnetwork.TnParticipantAttribute{}
	}
	return taurusnetwork.TnParticipantAttribute{
		ID:                    safeString(dto.Id),
		Key:                   safeString(dto.Key),
		Value:                 safeString(dto.Value),
		Owner:                 safeString(dto.Owner),
		Type:                  safeString(dto.Type),
		Subtype:               safeString(dto.Subtype),
		ContentType:           safeString(dto.ContentType),
		IsTaurusNetworkShared: safeBool(dto.IsTaurusNetworkShared),
	}
}

// TnParticipantSettingsFromDTO converts an OpenAPI TnParticipantSettings to a domain TnParticipantSettings.
func TnParticipantSettingsFromDTO(dto *openapi.TgvalidatordTnParticipantSettings) *taurusnetwork.TnParticipantSettings {
	if dto == nil {
		return nil
	}

	settings := &taurusnetwork.TnParticipantSettings{
		Status: safeString(dto.Status),
	}

	// Convert allowed countries
	if dto.InteractingAllowedCountries != nil {
		settings.InteractingAllowedCountries = dto.InteractingAllowedCountries
	}

	// Convert allowed participants
	if dto.InteractingAllowedParticipants != nil {
		settings.InteractingAllowedParticipants = make([]taurusnetwork.TnAllowedParticipant, len(dto.InteractingAllowedParticipants))
		for i := range dto.InteractingAllowedParticipants {
			settings.InteractingAllowedParticipants[i] = TnAllowedParticipantFromDTO(&dto.InteractingAllowedParticipants[i])
		}
	}

	// Convert terms and conditions accepted time
	if dto.TermsAndConditionsAcceptedAt != nil {
		settings.TermsAndConditionsAcceptedAt = *dto.TermsAndConditionsAcceptedAt
	}

	return settings
}

// TnAllowedParticipantFromDTO converts an OpenAPI TnAllowedParticipant to a domain TnAllowedParticipant.
func TnAllowedParticipantFromDTO(dto *openapi.TgvalidatordTnAllowedParticipant) taurusnetwork.TnAllowedParticipant {
	if dto == nil {
		return taurusnetwork.TnAllowedParticipant{}
	}
	return taurusnetwork.TnAllowedParticipant{
		ID:     safeString(dto.Id),
		Name:   safeString(dto.Name),
		Status: safeString(dto.Status),
	}
}

// CreateParticipantAttributeBodyToDTO converts a CreateParticipantAttributeRequest to an OpenAPI body.
func CreateParticipantAttributeBodyToDTO(req *taurusnetwork.CreateParticipantAttributeRequest) openapi.TaurusNetworkServiceCreateParticipantAttributeBody {
	if req == nil {
		return openapi.TaurusNetworkServiceCreateParticipantAttributeBody{}
	}

	attributeData := openapi.TgvalidatordParticipantAttributeData{}
	if req.Key != "" {
		attributeData.Key = stringPtr(req.Key)
	}
	if req.Value != "" {
		attributeData.Value = stringPtr(req.Value)
	}
	if req.ContentType != "" {
		attributeData.ContentType = stringPtr(req.ContentType)
	}
	if req.Type != "" {
		attributeData.Type = stringPtr(req.Type)
	}
	if req.Subtype != "" {
		attributeData.Subtype = stringPtr(req.Subtype)
	}

	body := openapi.TaurusNetworkServiceCreateParticipantAttributeBody{
		AttributeData:                   &attributeData,
		ShareToTaurusNetworkParticipant: boolPtr(req.ShareToTaurusNetworkParticipant),
	}

	return body
}
