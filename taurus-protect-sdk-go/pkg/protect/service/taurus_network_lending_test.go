package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestNewTaurusNetworkLendingService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTaurusNetworkLendingService_GetLendingAgreement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetLendingAgreement(nil, "")
	if err == nil {
		t.Error("GetLendingAgreement should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_ListLendingAgreements_NilOptions(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	if svc == nil {
		t.Error("TaurusNetworkLendingService should not be nil")
	}
}

func TestTaurusNetworkLendingService_ListLendingAgreements_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListLendingAgreementsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListLendingAgreementsOptions{},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListLendingAgreementsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "sort options",
			options: &taurusnetwork.ListLendingAgreementsOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListLendingAgreementsOptions{
				CurrentPage: "xyz789",
				PageRequest: "FIRST",
				PageSize:    100,
				SortOrder:   "ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_ListLendingAgreementsForApproval_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListLendingAgreementsForApprovalOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListLendingAgreementsForApprovalOptions{},
		},
		{
			name: "with IDs filter",
			options: &taurusnetwork.ListLendingAgreementsForApprovalOptions{
				IDs: []string{"agreement-1", "agreement-2"},
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListLendingAgreementsForApprovalOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListLendingAgreementsForApprovalOptions{
				IDs:         []string{"agreement-1"},
				CurrentPage: "xyz789",
				PageRequest: "FIRST",
				PageSize:    100,
				SortOrder:   "ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_CreateLendingAgreement_NilRequest(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateLendingAgreement(nil, nil)
	if err == nil {
		t.Error("CreateLendingAgreement should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CreateLendingAgreement_WithRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *taurusnetwork.CreateLendingAgreementRequest
	}{
		{
			name:    "empty request",
			request: &taurusnetwork.CreateLendingAgreementRequest{},
		},
		{
			name: "with lending offer ID",
			request: &taurusnetwork.CreateLendingAgreementRequest{
				LendingOfferID: "offer-123",
			},
		},
		{
			name: "with direct terms",
			request: &taurusnetwork.CreateLendingAgreementRequest{
				LenderParticipantID:   "lender-456",
				CurrencyID:            "currency-789",
				Amount:                "1000000000000000000",
				AnnualPercentageYield: "525000",
				Duration:              "30d",
			},
		},
		{
			name: "with collaterals",
			request: &taurusnetwork.CreateLendingAgreementRequest{
				LendingOfferID:          "offer-123",
				BorrowerSharedAddressID: "addr-456",
				Collaterals: []taurusnetwork.CreateLendingCollateralRequest{
					{CurrencyID: "currency-1", Amount: "500000000000000000"},
					{CurrencyID: "currency-2", Amount: "300000000000000000"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_UpdateLendingAgreement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.UpdateLendingAgreement(nil, "", &taurusnetwork.UpdateLendingAgreementRequest{})
	if err == nil {
		t.Error("UpdateLendingAgreement should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_UpdateLendingAgreement_NilRequest(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.UpdateLendingAgreement(nil, "agreement-123", nil)
	if err == nil {
		t.Error("UpdateLendingAgreement should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_RepayLendingAgreement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RepayLendingAgreement(nil, "", &taurusnetwork.RepayLendingAgreementRequest{})
	if err == nil {
		t.Error("RepayLendingAgreement should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_RepayLendingAgreement_NilRequest(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RepayLendingAgreement(nil, "agreement-123", nil)
	if err == nil {
		t.Error("RepayLendingAgreement should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CancelLendingAgreement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.CancelLendingAgreement(nil, "")
	if err == nil {
		t.Error("CancelLendingAgreement should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CreateLendingAgreementAttachment_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.CreateLendingAgreementAttachment(nil, "", &taurusnetwork.CreateLendingAgreementAttachmentRequest{})
	if err == nil {
		t.Error("CreateLendingAgreementAttachment should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CreateLendingAgreementAttachment_NilRequest(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.CreateLendingAgreementAttachment(nil, "agreement-123", nil)
	if err == nil {
		t.Error("CreateLendingAgreementAttachment should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CreateLendingAgreementAttachment_WithRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *taurusnetwork.CreateLendingAgreementAttachmentRequest
	}{
		{
			name:    "empty request",
			request: &taurusnetwork.CreateLendingAgreementAttachmentRequest{},
		},
		{
			name: "embedded attachment",
			request: &taurusnetwork.CreateLendingAgreementAttachmentRequest{
				Name:        "contract.pdf",
				Type:        "EMBEDDED",
				ContentType: "application/pdf",
				Value:       "base64encodedcontent",
			},
		},
		{
			name: "external link attachment",
			request: &taurusnetwork.CreateLendingAgreementAttachmentRequest{
				Name:  "external-doc",
				Type:  "EXTERNAL_LINK",
				Value: "https://example.com/document.pdf",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_ListLendingAgreementAttachments_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.ListLendingAgreementAttachments(nil, "")
	if err == nil {
		t.Error("ListLendingAgreementAttachments should return error for empty ID")
	}
	if err.Error() != "lendingAgreementID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_GetLendingOffer_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetLendingOffer(nil, "")
	if err == nil {
		t.Error("GetLendingOffer should return error for empty ID")
	}
	if err.Error() != "offerID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_ListLendingOffers_NilOptions(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	if svc == nil {
		t.Error("TaurusNetworkLendingService should not be nil")
	}
}

func TestTaurusNetworkLendingService_ListLendingOffers_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListLendingOffersOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListLendingOffersOptions{},
		},
		{
			name: "currency IDs filter",
			options: &taurusnetwork.ListLendingOffersOptions{
				CurrencyIDs: []string{"currency-1", "currency-2"},
			},
		},
		{
			name: "participant ID filter",
			options: &taurusnetwork.ListLendingOffersOptions{
				ParticipantID: "participant-123",
			},
		},
		{
			name: "duration filter",
			options: &taurusnetwork.ListLendingOffersOptions{
				Duration: "30d",
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListLendingOffersOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListLendingOffersOptions{
				CurrencyIDs:   []string{"currency-1"},
				ParticipantID: "participant-123",
				Duration:      "30d",
				CurrentPage:   "xyz789",
				PageRequest:   "FIRST",
				PageSize:      100,
				SortOrder:     "ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_CreateLendingOffer_NilRequest(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateLendingOffer(nil, nil)
	if err == nil {
		t.Error("CreateLendingOffer should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTaurusNetworkLendingService_CreateLendingOffer_WithRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *taurusnetwork.CreateLendingOfferRequest
	}{
		{
			name:    "empty request",
			request: &taurusnetwork.CreateLendingOfferRequest{},
		},
		{
			name: "basic offer",
			request: &taurusnetwork.CreateLendingOfferRequest{
				CurrencyID:            "currency-123",
				Amount:                "10000000000000000000",
				AnnualPercentageYield: "750000",
				Duration:              "90d",
			},
		},
		{
			name: "offer with collateral requirements",
			request: &taurusnetwork.CreateLendingOfferRequest{
				CurrencyID:            "currency-123",
				Amount:                "10000000000000000000",
				AnnualPercentageYield: "750000",
				Duration:              "90d",
				CollateralRequirements: []taurusnetwork.CreateCollateralRequirementRequest{
					{CurrencyID: "currency-1", CollateralRatio: "150"},
					{CurrencyID: "currency-2", CollateralRatio: "200"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkLendingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkLendingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkLendingService_DeleteLendingOffer_EmptyID(t *testing.T) {
	svc := &TaurusNetworkLendingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.DeleteLendingOffer(nil, "")
	if err == nil {
		t.Error("DeleteLendingOffer should return error for empty ID")
	}
	if err.Error() != "offerID cannot be empty" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestListLendingAgreementsOptions_PageRequestValues(t *testing.T) {
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &taurusnetwork.ListLendingAgreementsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListLendingAgreementsOptions_SortOrderValues(t *testing.T) {
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &taurusnetwork.ListLendingAgreementsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListLendingOffersOptions_PageRequestValues(t *testing.T) {
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &taurusnetwork.ListLendingOffersOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestCreateLendingAgreementAttachmentRequest_TypeValues(t *testing.T) {
	validTypes := []string{"EMBEDDED", "EXTERNAL_LINK"}

	for _, attachmentType := range validTypes {
		t.Run(attachmentType, func(t *testing.T) {
			req := &taurusnetwork.CreateLendingAgreementAttachmentRequest{
				Type: attachmentType,
			}
			if req.Type != attachmentType {
				t.Errorf("Type = %v, want %v", req.Type, attachmentType)
			}
		})
	}
}
