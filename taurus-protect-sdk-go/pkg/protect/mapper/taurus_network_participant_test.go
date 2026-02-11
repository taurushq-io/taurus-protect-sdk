package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestTnParticipantFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnParticipant
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns participant with zero values",
			dto:  &openapi.TgvalidatordTnParticipant{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnParticipant {
				id := "participant-123"
				name := "Test Participant"
				legalAddress := "123 Test St"
				country := "CH"
				logoBase64 := "base64logo"
				publicKey := "pubkey123"
				shield := "shield-config"
				ownedCount := "10"
				targetedCount := "5"
				outgoingVal := "1000.00"
				incomingVal := "500.00"
				publicSubname := "test-subname"
				lei := "LEI123456"
				status := "ACTIVE"
				now := time.Now()
				return &openapi.TgvalidatordTnParticipant{
					Id:                                        &id,
					Name:                                      &name,
					LegalAddress:                              &legalAddress,
					Country:                                   &country,
					LogoBase64:                                &logoBase64,
					PublicKey:                                 &publicKey,
					Shield:                                    &shield,
					OriginRegistrationDate:                    &now,
					OriginDeletionDate:                        &now,
					CreatedAt:                                 &now,
					UpdatedAt:                                 &now,
					OwnedSharedAddressesCount:                 &ownedCount,
					TargetedSharedAddressesCount:              &targetedCount,
					OutgoingTotalPledgesValuationBaseCurrency: &outgoingVal,
					IncomingTotalPledgesValuationBaseCurrency: &incomingVal,
					PublicSubname:                             &publicSubname,
					LegalEntityIdentifier:                     &lei,
					Status:                                    &status,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TnParticipantFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TnParticipantFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.LegalAddress != nil && got.LegalAddress != *tt.dto.LegalAddress {
				t.Errorf("LegalAddress = %v, want %v", got.LegalAddress, *tt.dto.LegalAddress)
			}
			if tt.dto.Country != nil && got.Country != *tt.dto.Country {
				t.Errorf("Country = %v, want %v", got.Country, *tt.dto.Country)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestTnParticipantsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnParticipant
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnParticipant{},
			want: 0,
		},
		{
			name: "converts multiple participants",
			dtos: func() []openapi.TgvalidatordTnParticipant {
				name1 := "Participant 1"
				name2 := "Participant 2"
				return []openapi.TgvalidatordTnParticipant{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("TnParticipantsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("TnParticipantsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestTnParticipantDetailsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnParticipantDetails
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty details",
			dto:  &openapi.TgvalidatordTnParticipantDetails{},
		},
		{
			name: "DTO with contact persons",
			dto: func() *openapi.TgvalidatordTnParticipantDetails {
				firstName := "John"
				lastName := "Doe"
				return &openapi.TgvalidatordTnParticipantDetails{
					ContactPersons: []openapi.TgvalidatordTnContactPerson{
						{FirstName: &firstName, LastName: &lastName},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantDetailsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TnParticipantDetailsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TnParticipantDetailsFromDTO() returned nil for non-nil input")
			}
			if tt.dto.ContactPersons != nil && len(got.ContactPersons) != len(tt.dto.ContactPersons) {
				t.Errorf("ContactPersons length = %v, want %v", len(got.ContactPersons), len(tt.dto.ContactPersons))
			}
		})
	}
}

func TestTnContactPersonFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnContactPerson
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordTnContactPerson{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnContactPerson {
				firstName := "John"
				lastName := "Doe"
				phone := "+41123456789"
				email := "john.doe@example.com"
				return &openapi.TgvalidatordTnContactPerson{
					FirstName:   &firstName,
					LastName:    &lastName,
					PhoneNumber: &phone,
					Email:       &email,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnContactPersonFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnContactPerson{}) {
					t.Errorf("TnContactPersonFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.FirstName != nil && got.FirstName != *tt.dto.FirstName {
				t.Errorf("FirstName = %v, want %v", got.FirstName, *tt.dto.FirstName)
			}
			if tt.dto.LastName != nil && got.LastName != *tt.dto.LastName {
				t.Errorf("LastName = %v, want %v", got.LastName, *tt.dto.LastName)
			}
			if tt.dto.Email != nil && got.Email != *tt.dto.Email {
				t.Errorf("Email = %v, want %v", got.Email, *tt.dto.Email)
			}
		})
	}
}

func TestTnParticipantAttributeFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnParticipantAttribute
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordTnParticipantAttribute{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnParticipantAttribute {
				id := "attr-123"
				key := "license"
				value := "12345"
				owner := "participant-1"
				attrType := "COMPLIANCE"
				subtype := "LICENSE"
				contentType := "text/plain"
				isShared := true
				return &openapi.TgvalidatordTnParticipantAttribute{
					Id:                    &id,
					Key:                   &key,
					Value:                 &value,
					Owner:                 &owner,
					Type:                  &attrType,
					Subtype:               &subtype,
					ContentType:           &contentType,
					IsTaurusNetworkShared: &isShared,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantAttributeFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnParticipantAttribute{}) {
					t.Errorf("TnParticipantAttributeFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Key != nil && got.Key != *tt.dto.Key {
				t.Errorf("Key = %v, want %v", got.Key, *tt.dto.Key)
			}
			if tt.dto.Value != nil && got.Value != *tt.dto.Value {
				t.Errorf("Value = %v, want %v", got.Value, *tt.dto.Value)
			}
			if tt.dto.IsTaurusNetworkShared != nil && got.IsTaurusNetworkShared != *tt.dto.IsTaurusNetworkShared {
				t.Errorf("IsTaurusNetworkShared = %v, want %v", got.IsTaurusNetworkShared, *tt.dto.IsTaurusNetworkShared)
			}
		})
	}
}

func TestTnBlockConfirmationsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordBlockConfirmations
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordBlockConfirmations{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordBlockConfirmations {
				blockchain := "ETH"
				network := "mainnet"
				threshold := "12"
				return &openapi.TgvalidatordBlockConfirmations{
					Blockchain:             &blockchain,
					Network:                &network,
					ConfirmationsThreshold: &threshold,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnBlockConfirmationsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnBlockConfirmations{}) {
					t.Errorf("TnBlockConfirmationsFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.ConfirmationsThreshold != nil && got.ConfirmationsThreshold != *tt.dto.ConfirmationsThreshold {
				t.Errorf("ConfirmationsThreshold = %v, want %v", got.ConfirmationsThreshold, *tt.dto.ConfirmationsThreshold)
			}
		})
	}
}

func TestTnParticipantSettingsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnParticipantSettings
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns settings with zero values",
			dto:  &openapi.TgvalidatordTnParticipantSettings{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnParticipantSettings {
				status := "ACTIVE"
				now := time.Now()
				return &openapi.TgvalidatordTnParticipantSettings{
					InteractingAllowedCountries:   []string{"CH", "DE", "US"},
					Status:                        &status,
					TermsAndConditionsAcceptedAt:  &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantSettingsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TnParticipantSettingsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TnParticipantSettingsFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.InteractingAllowedCountries != nil && len(got.InteractingAllowedCountries) != len(tt.dto.InteractingAllowedCountries) {
				t.Errorf("InteractingAllowedCountries length = %v, want %v", len(got.InteractingAllowedCountries), len(tt.dto.InteractingAllowedCountries))
			}
		})
	}
}

func TestTnAllowedParticipantFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnAllowedParticipant
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordTnAllowedParticipant{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnAllowedParticipant {
				id := "allowed-123"
				name := "Allowed Participant"
				status := "APPROVED"
				return &openapi.TgvalidatordTnAllowedParticipant{
					Id:     &id,
					Name:   &name,
					Status: &status,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnAllowedParticipantFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnAllowedParticipant{}) {
					t.Errorf("TnAllowedParticipantFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestCreateParticipantAttributeBodyToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *taurusnetwork.CreateParticipantAttributeRequest
	}{
		{
			name: "nil input returns empty body",
			req:  nil,
		},
		{
			name: "empty request returns body with empty attribute data",
			req:  &taurusnetwork.CreateParticipantAttributeRequest{},
		},
		{
			name: "complete request maps all fields",
			req: &taurusnetwork.CreateParticipantAttributeRequest{
				Key:                             "license",
				Value:                           "12345",
				ContentType:                     "text/plain",
				Type:                            "COMPLIANCE",
				Subtype:                         "LICENSE",
				ShareToTaurusNetworkParticipant: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateParticipantAttributeBodyToDTO(tt.req)
			if tt.req == nil {
				if got.AttributeData != nil || got.ShareToTaurusNetworkParticipant != nil {
					t.Errorf("CreateParticipantAttributeBodyToDTO() should return empty body for nil input")
				}
				return
			}
			if tt.req.Key != "" && (got.AttributeData == nil || got.AttributeData.Key == nil || *got.AttributeData.Key != tt.req.Key) {
				t.Errorf("Key mapping failed")
			}
			if got.ShareToTaurusNetworkParticipant == nil || *got.ShareToTaurusNetworkParticipant != tt.req.ShareToTaurusNetworkParticipant {
				t.Errorf("ShareToTaurusNetworkParticipant = %v, want %v", got.ShareToTaurusNetworkParticipant, tt.req.ShareToTaurusNetworkParticipant)
			}
		})
	}
}

func TestTnParticipantFromDTO_WithDetails(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	dto := &openapi.TgvalidatordTnParticipant{
		Details: &openapi.TgvalidatordTnParticipantDetails{
			ContactPersons: []openapi.TgvalidatordTnContactPerson{
				{FirstName: &firstName, LastName: &lastName},
			},
		},
	}

	got := TnParticipantFromDTO(dto)
	if got == nil {
		t.Fatal("TnParticipantFromDTO() returned nil for non-nil input")
	}
	if got.Details == nil {
		t.Fatal("Details should not be nil")
	}
	if len(got.Details.ContactPersons) != 1 {
		t.Errorf("ContactPersons length = %v, want 1", len(got.Details.ContactPersons))
	}
	if got.Details.ContactPersons[0].FirstName != "John" {
		t.Errorf("ContactPersons[0].FirstName = %v, want John", got.Details.ContactPersons[0].FirstName)
	}
}

func TestTnParticipantFromDTO_WithAttributes(t *testing.T) {
	attrID := "attr-1"
	attrKey := "license"
	dto := &openapi.TgvalidatordTnParticipant{
		Attributes: []openapi.TgvalidatordTnParticipantAttribute{
			{Id: &attrID, Key: &attrKey},
		},
	}

	got := TnParticipantFromDTO(dto)
	if got == nil {
		t.Fatal("TnParticipantFromDTO() returned nil for non-nil input")
	}
	if len(got.Attributes) != 1 {
		t.Errorf("Attributes length = %v, want 1", len(got.Attributes))
	}
	if got.Attributes[0].ID != "attr-1" {
		t.Errorf("Attributes[0].ID = %v, want attr-1", got.Attributes[0].ID)
	}
}

func TestTnParticipantFromDTO_WithBlockConfirmations(t *testing.T) {
	blockchain := "ETH"
	network := "mainnet"
	threshold := "12"
	dto := &openapi.TgvalidatordTnParticipant{
		BlockConfirmations: []openapi.TgvalidatordBlockConfirmations{
			{Blockchain: &blockchain, Network: &network, ConfirmationsThreshold: &threshold},
		},
	}

	got := TnParticipantFromDTO(dto)
	if got == nil {
		t.Fatal("TnParticipantFromDTO() returned nil for non-nil input")
	}
	if len(got.BlockConfirmations) != 1 {
		t.Errorf("BlockConfirmations length = %v, want 1", len(got.BlockConfirmations))
	}
	if got.BlockConfirmations[0].Blockchain != "ETH" {
		t.Errorf("BlockConfirmations[0].Blockchain = %v, want ETH", got.BlockConfirmations[0].Blockchain)
	}
}

func TestTnParticipantSettingsFromDTO_WithAllowedParticipants(t *testing.T) {
	id := "allowed-1"
	name := "Allowed Participant 1"
	dto := &openapi.TgvalidatordTnParticipantSettings{
		InteractingAllowedParticipants: []openapi.TgvalidatordTnAllowedParticipant{
			{Id: &id, Name: &name},
		},
	}

	got := TnParticipantSettingsFromDTO(dto)
	if got == nil {
		t.Fatal("TnParticipantSettingsFromDTO() returned nil for non-nil input")
	}
	if len(got.InteractingAllowedParticipants) != 1 {
		t.Errorf("InteractingAllowedParticipants length = %v, want 1", len(got.InteractingAllowedParticipants))
	}
	if got.InteractingAllowedParticipants[0].ID != "allowed-1" {
		t.Errorf("InteractingAllowedParticipants[0].ID = %v, want allowed-1", got.InteractingAllowedParticipants[0].ID)
	}
}

func TestTnParticipantAttributeSpecificationFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnParticipantAttributeSpecification
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordTnParticipantAttributeSpecification{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnParticipantAttributeSpecification {
				key := "license"
				attrType := "COMPLIANCE"
				desc := "License number"
				return &openapi.TgvalidatordTnParticipantAttributeSpecification{
					AttributeKey:         &key,
					AttributeType:        &attrType,
					AttributeDescription: &desc,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnParticipantAttributeSpecificationFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnParticipantAttributeSpecification{}) {
					t.Errorf("TnParticipantAttributeSpecificationFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.AttributeKey != nil && got.AttributeKey != *tt.dto.AttributeKey {
				t.Errorf("AttributeKey = %v, want %v", got.AttributeKey, *tt.dto.AttributeKey)
			}
			if tt.dto.AttributeType != nil && got.AttributeType != *tt.dto.AttributeType {
				t.Errorf("AttributeType = %v, want %v", got.AttributeType, *tt.dto.AttributeType)
			}
			if tt.dto.AttributeDescription != nil && got.AttributeDescription != *tt.dto.AttributeDescription {
				t.Errorf("AttributeDescription = %v, want %v", got.AttributeDescription, *tt.dto.AttributeDescription)
			}
		})
	}
}

func TestTnBlockchainEntityFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordBlockchainEntity
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns empty struct",
			dto:  &openapi.TgvalidatordBlockchainEntity{},
		},
		{
			name: "complete DTO maps relevant fields",
			dto: func() *openapi.TgvalidatordBlockchainEntity {
				symbol := "ETH"
				name := "Ethereum"
				network := "mainnet"
				return &openapi.TgvalidatordBlockchainEntity{
					Symbol:  &symbol,
					Name:    &name,
					Network: &network,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TnBlockchainEntityFromDTO(tt.dto)
			if tt.dto == nil {
				if got != (taurusnetwork.TnBlockchainEntity{}) {
					t.Errorf("TnBlockchainEntityFromDTO() = %v, want empty struct", got)
				}
				return
			}
			if tt.dto.Symbol != nil && got.Symbol != *tt.dto.Symbol {
				t.Errorf("Symbol = %v, want %v", got.Symbol, *tt.dto.Symbol)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
		})
	}
}
