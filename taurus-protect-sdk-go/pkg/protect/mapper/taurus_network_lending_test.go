package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestLendingAgreementFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordLendingAgreement
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns agreement with zero values",
			dto:  &openapi.TgvalidatordLendingAgreement{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordLendingAgreement {
				now := time.Now()
				return &openapi.TgvalidatordLendingAgreement{
					Id:                      strPtr("agreement-123"),
					BorrowerParticipantID:   strPtr("borrower-1"),
					LenderParticipantID:     strPtr("lender-1"),
					LendingOfferID:          strPtr("offer-456"),
					CurrencyID:              strPtr("ETH"),
					Amount:                  strPtr("1000000000000000000"),
					AmountMainUnit:          strPtr("1.0"),
					AnnualYield:             strPtr("500"),
					AnnualYieldMainUnit:     strPtr("5.00"),
					Duration:                strPtr("30d"),
					Status:                  strPtr("ACTIVE"),
					WorkflowID:              strPtr("wf-789"),
					BorrowerSharedAddressID: strPtr("bsa-1"),
					LenderSharedAddressID:   strPtr("lsa-1"),
					CreatedAt:               &now,
					UpdatedAt:               &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingAgreementFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("LendingAgreementFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("LendingAgreementFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.BorrowerParticipantID != nil && got.BorrowerParticipantID != *tt.dto.BorrowerParticipantID {
				t.Errorf("BorrowerParticipantID = %v, want %v", got.BorrowerParticipantID, *tt.dto.BorrowerParticipantID)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestLendingAgreementsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordLendingAgreement
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordLendingAgreement{},
			want: 0,
		},
		{
			name: "converts multiple agreements",
			dtos: func() []openapi.TgvalidatordLendingAgreement {
				return []openapi.TgvalidatordLendingAgreement{
					{Id: strPtr("agreement-1")},
					{Id: strPtr("agreement-2")},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingAgreementsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("LendingAgreementsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("LendingAgreementsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestLendingAgreementCollateralFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordLendingAgreementCollateral
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordLendingAgreementCollateral {
				return &openapi.TgvalidatordLendingAgreementCollateral{
					Id:                 strPtr("col-123"),
					LendingAgreementID: strPtr("agreement-456"),
					CurrencyID:         strPtr("BTC"),
					Amount:             strPtr("100000000"),
					AmountMainUnit:     strPtr("1.0"),
					Status:             strPtr("LOCKED"),
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingAgreementCollateralFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" {
					t.Errorf("LendingAgreementCollateralFromDTO(nil) should return zero value, got ID=%v", got.ID)
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestLendingAgreementTransactionFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordLendingAgreementTransaction
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordLendingAgreementTransaction {
				return &openapi.TgvalidatordLendingAgreementTransaction{
					Id:                     strPtr("tx-123"),
					LendingAgreementID:     strPtr("agreement-456"),
					Amount:                 strPtr("500000000"),
					CurrencyID:             strPtr("ETH"),
					TransactionHash:        strPtr("0xabc123"),
					TransactionBlockNumber: strPtr("12345678"),
					Type:                   strPtr("REPAYMENT"),
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingAgreementTransactionFromDTO(tt.dto)
			if tt.dto == nil {
				if got.ID != "" {
					t.Errorf("LendingAgreementTransactionFromDTO(nil) should return zero value, got ID=%v", got.ID)
				}
				return
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TransactionHash != nil && got.TransactionHash != *tt.dto.TransactionHash {
				t.Errorf("TransactionHash = %v, want %v", got.TransactionHash, *tt.dto.TransactionHash)
			}
		})
	}
}

func TestLendingAgreementAttachmentFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordLendingAgreementAttachment
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordLendingAgreementAttachment {
				return &openapi.TgvalidatordLendingAgreementAttachment{
					Id:                    strPtr("att-123"),
					LendingAgreementID:    strPtr("agreement-456"),
					UploaderParticipantID: strPtr("participant-1"),
					Name:                  strPtr("contract.pdf"),
					Type:                  strPtr("DOCUMENT"),
					ContentType:           strPtr("application/pdf"),
					FileSize:              strPtr("1024"),
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingAgreementAttachmentFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("LendingAgreementAttachmentFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("LendingAgreementAttachmentFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
		})
	}
}

func TestLendingOfferFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnLendingOffer
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns offer with zero values",
			dto:  &openapi.TgvalidatordTnLendingOffer{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnLendingOffer {
				now := time.Now()
				return &openapi.TgvalidatordTnLendingOffer{
					Id:                            strPtr("offer-123"),
					ParticipantID:                 strPtr("participant-1"),
					Blockchain:                    strPtr("ETH"),
					Network:                       strPtr("mainnet"),
					AnnualPercentageYield:         strPtr("500"),
					AnnualPercentageYieldMainUnit: strPtr("5.00"),
					Duration:                      strPtr("30d"),
					CreatedAt:                     &now,
					UpdatedAt:                     &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingOfferFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("LendingOfferFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("LendingOfferFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
		})
	}
}

func TestLendingCollateralRequirementFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordLendingCollateralRequirement
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns requirement with nil accepted currencies",
			dto:  &openapi.TgvalidatordLendingCollateralRequirement{},
		},
		{
			name: "complete DTO with accepted currencies",
			dto: &openapi.TgvalidatordLendingCollateralRequirement{
				AcceptedCurrencies: []openapi.TgvalidatordCurrencyCollateralRequirement{
					{
						Blockchain: strPtr("ETH"),
						Network:    strPtr("mainnet"),
						Ratio:      strPtr("1.5"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LendingCollateralRequirementFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("LendingCollateralRequirementFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("LendingCollateralRequirementFromDTO() returned nil for non-nil input")
			}
			if tt.dto.AcceptedCurrencies != nil && len(got.AcceptedCurrencies) != len(tt.dto.AcceptedCurrencies) {
				t.Errorf("AcceptedCurrencies length = %v, want %v", len(got.AcceptedCurrencies), len(tt.dto.AcceptedCurrencies))
			}
		})
	}
}

func TestCurrencyCollateralRequirementFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordCurrencyCollateralRequirement
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "complete DTO maps all fields",
			dto: &openapi.TgvalidatordCurrencyCollateralRequirement{
				Blockchain: strPtr("ETH"),
				Network:    strPtr("mainnet"),
				Arg1:       strPtr("arg1"),
				Arg2:       strPtr("arg2"),
				Ratio:      strPtr("1.5"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CurrencyCollateralRequirementFromDTO(tt.dto)
			if tt.dto == nil {
				if got.Blockchain != "" || got.Network != "" {
					t.Errorf("CurrencyCollateralRequirementFromDTO(nil) should return zero value")
				}
				return
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Ratio != nil && got.Ratio != *tt.dto.Ratio {
				t.Errorf("Ratio = %v, want %v", got.Ratio, *tt.dto.Ratio)
			}
		})
	}
}
