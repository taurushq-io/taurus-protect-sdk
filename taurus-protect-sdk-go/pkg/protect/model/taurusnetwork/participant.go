package taurusnetwork

import "time"

// TnParticipant represents a Taurus-NETWORK participant.
type TnParticipant struct {
	// ID is the unique identifier of the participant.
	ID string `json:"id"`
	// Name is the name of the participant.
	Name string `json:"name"`
	// LegalAddress is the legal address of the participant.
	LegalAddress string `json:"legal_address,omitempty"`
	// Country is the country code of the participant.
	Country string `json:"country,omitempty"`
	// LogoBase64 is the base64-encoded logo of the participant.
	LogoBase64 string `json:"logo_base64,omitempty"`
	// PublicKey is the public key of the participant.
	PublicKey string `json:"public_key,omitempty"`
	// Shield is the shield configuration of the participant.
	Shield string `json:"shield,omitempty"`
	// OriginRegistrationDate is when the participant was originally registered.
	OriginRegistrationDate time.Time `json:"origin_registration_date,omitempty"`
	// OriginDeletionDate is when the participant was originally deleted (if applicable).
	OriginDeletionDate time.Time `json:"origin_deletion_date,omitempty"`
	// CreatedAt is when the participant was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt is when the participant was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Details contains additional details about the participant.
	Details *TnParticipantDetails `json:"details,omitempty"`
	// BlockConfirmations contains the block confirmation settings.
	BlockConfirmations []TnBlockConfirmations `json:"block_confirmations,omitempty"`
	// OwnedSharedAddressesCount is the count of shared addresses owned by the participant.
	OwnedSharedAddressesCount string `json:"owned_shared_addresses_count,omitempty"`
	// TargetedSharedAddressesCount is the count of shared addresses where the participant is the target.
	TargetedSharedAddressesCount string `json:"targeted_shared_addresses_count,omitempty"`
	// OutgoingTotalPledgesValuationBaseCurrency is the total valuation of outgoing pledges in base currency.
	OutgoingTotalPledgesValuationBaseCurrency string `json:"outgoing_total_pledges_valuation_base_currency,omitempty"`
	// IncomingTotalPledgesValuationBaseCurrency is the total valuation of incoming pledges in base currency.
	IncomingTotalPledgesValuationBaseCurrency string `json:"incoming_total_pledges_valuation_base_currency,omitempty"`
	// PublicSubname is the public subname of the participant.
	PublicSubname string `json:"public_subname,omitempty"`
	// LegalEntityIdentifier is the legal entity identifier (LEI) of the participant.
	LegalEntityIdentifier string `json:"legal_entity_identifier,omitempty"`
	// Attributes are the attributes of the participant.
	Attributes []TnParticipantAttribute `json:"attributes,omitempty"`
	// Status is the status of the participant.
	Status string `json:"status,omitempty"`
}

// TnParticipantDetails contains detailed information about a participant.
type TnParticipantDetails struct {
	// ContactPersons contains the contact persons for the participant.
	ContactPersons []TnContactPerson `json:"contact_persons,omitempty"`
	// AttributesSpecifications contains the attribute specifications for the participant.
	AttributesSpecifications []TnParticipantAttributeSpecification `json:"attributes_specifications,omitempty"`
	// SupportedBlockchains contains the supported blockchains for the participant.
	SupportedBlockchains []TnBlockchainEntity `json:"supported_blockchains,omitempty"`
}

// TnContactPerson represents a contact person for a participant.
type TnContactPerson struct {
	// FirstName is the first name of the contact person.
	FirstName string `json:"first_name,omitempty"`
	// LastName is the last name of the contact person.
	LastName string `json:"last_name,omitempty"`
	// PhoneNumber is the phone number of the contact person.
	PhoneNumber string `json:"phone_number,omitempty"`
	// Email is the email address of the contact person.
	Email string `json:"email,omitempty"`
}

// TnParticipantAttributeSpecification describes an attribute specification for a participant.
type TnParticipantAttributeSpecification struct {
	// AttributeKey is the key of the attribute.
	AttributeKey string `json:"attribute_key,omitempty"`
	// AttributeType is the type of the attribute.
	AttributeType string `json:"attribute_type,omitempty"`
	// AttributeDescription is the description of the attribute.
	AttributeDescription string `json:"attribute_description,omitempty"`
}

// TnBlockchainEntity represents a supported blockchain for a participant.
type TnBlockchainEntity struct {
	// Symbol is the symbol of the blockchain (e.g., BTC, ETH).
	Symbol string `json:"symbol,omitempty"`
	// Name is the name of the blockchain.
	Name string `json:"name,omitempty"`
	// Network is the network name.
	Network string `json:"network,omitempty"`
}

// TnBlockConfirmations contains block confirmation settings for a blockchain.
type TnBlockConfirmations struct {
	// Blockchain is the blockchain name.
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network name.
	Network string `json:"network,omitempty"`
	// ConfirmationsThreshold is the required number of confirmations.
	ConfirmationsThreshold string `json:"confirmations_threshold,omitempty"`
}

// TnParticipantAttribute represents an attribute of a participant.
type TnParticipantAttribute struct {
	// ID is the unique identifier of the attribute.
	ID string `json:"id"`
	// Key is the attribute key.
	Key string `json:"key,omitempty"`
	// Value is the attribute value.
	Value string `json:"value,omitempty"`
	// Owner is the owner of the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the type of the attribute.
	Type string `json:"type,omitempty"`
	// Subtype is the subtype of the attribute.
	Subtype string `json:"subtype,omitempty"`
	// ContentType is the content type of the attribute.
	ContentType string `json:"content_type,omitempty"`
	// IsTaurusNetworkShared indicates if the attribute is shared on Taurus-NETWORK.
	IsTaurusNetworkShared bool `json:"is_taurus_network_shared"`
}

// TnParticipantSettings contains settings for the current participant.
type TnParticipantSettings struct {
	// InteractingAllowedCountries contains the list of allowed countries for interaction.
	InteractingAllowedCountries []string `json:"interacting_allowed_countries,omitempty"`
	// Status is the status of the participant settings.
	Status string `json:"status,omitempty"`
	// InteractingAllowedParticipants contains the list of allowed participants for interaction.
	InteractingAllowedParticipants []TnAllowedParticipant `json:"interacting_allowed_participants,omitempty"`
	// TermsAndConditionsAcceptedAt is when the terms and conditions were accepted.
	TermsAndConditionsAcceptedAt time.Time `json:"terms_and_conditions_accepted_at,omitempty"`
}

// TnAllowedParticipant represents a participant allowed for interaction.
type TnAllowedParticipant struct {
	// ID is the unique identifier of the participant.
	ID string `json:"id"`
	// Name is the name of the participant.
	Name string `json:"name,omitempty"`
	// Status is the status of the allowed participant relationship.
	Status string `json:"status,omitempty"`
}

// GetMyParticipantResult contains the result of getting the current participant.
type GetMyParticipantResult struct {
	// Participant is the current participant.
	Participant *TnParticipant `json:"participant,omitempty"`
	// Settings contains the settings for the current participant.
	Settings *TnParticipantSettings `json:"settings,omitempty"`
}

// GetParticipantOptions contains options for getting a participant.
type GetParticipantOptions struct {
	// IncludeTotalPledgesValuation includes the total pledges valuation if true.
	IncludeTotalPledgesValuation bool
}

// ListParticipantsOptions contains options for listing participants.
type ListParticipantsOptions struct {
	// ParticipantIDs filters participants by the specified IDs.
	ParticipantIDs []string
	// IncludeTotalPledgesValuation includes the total pledges valuation if true.
	IncludeTotalPledgesValuation bool
}

// ListParticipantsResult contains the result of listing participants.
type ListParticipantsResult struct {
	// Participants is the list of participants.
	Participants []*TnParticipant `json:"participants"`
}

// CreateParticipantAttributeRequest contains the data for creating a participant attribute.
type CreateParticipantAttributeRequest struct {
	// Key is the attribute key.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the content type of the attribute.
	ContentType string `json:"content_type,omitempty"`
	// Type is the type of the attribute.
	Type string `json:"type,omitempty"`
	// Subtype is the subtype of the attribute.
	Subtype string `json:"subtype,omitempty"`
	// ShareToTaurusNetworkParticipant indicates if the attribute should be shared.
	ShareToTaurusNetworkParticipant bool `json:"share_to_taurus_network_participant"`
}
