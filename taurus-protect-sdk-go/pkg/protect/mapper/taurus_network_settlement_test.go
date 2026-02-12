package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestSettlementFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSettlement
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns settlement with zero values",
			dto:  &openapi.TgvalidatordTnSettlement{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSettlement {
				id := "settlement-123"
				creatorID := "creator-456"
				targetID := "target-789"
				firstLegID := "firstleg-111"
				status := "CREATED"
				workflowID := "workflow-222"
				createdAt := time.Now()
				updatedAt := time.Now().Add(time.Hour)
				startDate := time.Now().Add(24 * time.Hour)
				currencyID := "BTC"
				amount := "1.5"

				return &openapi.TgvalidatordTnSettlement{
					Id:                    &id,
					CreatorParticipantID:  &creatorID,
					TargetParticipantID:   &targetID,
					FirstLegParticipantID: &firstLegID,
					Status:                &status,
					WorkflowID:            &workflowID,
					CreatedAt:             &createdAt,
					UpdatedAt:             &updatedAt,
					StartExecutionDate:    &startDate,
					FirstLegAssets: []openapi.TgvalidatordTnSettlementAssetTransfer{
						{CurrencyID: &currencyID, Amount: &amount},
					},
					SecondLegAssets: []openapi.TgvalidatordTnSettlementAssetTransfer{
						{CurrencyID: &currencyID, Amount: &amount},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SettlementFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SettlementFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.CreatorParticipantID != nil && got.CreatorParticipantID != *tt.dto.CreatorParticipantID {
				t.Errorf("CreatorParticipantID = %v, want %v", got.CreatorParticipantID, *tt.dto.CreatorParticipantID)
			}
			if tt.dto.TargetParticipantID != nil && got.TargetParticipantID != *tt.dto.TargetParticipantID {
				t.Errorf("TargetParticipantID = %v, want %v", got.TargetParticipantID, *tt.dto.TargetParticipantID)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.FirstLegAssets != nil && len(got.FirstLegAssets) != len(tt.dto.FirstLegAssets) {
				t.Errorf("FirstLegAssets length = %v, want %v", len(got.FirstLegAssets), len(tt.dto.FirstLegAssets))
			}
		})
	}
}

func TestSettlementsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTnSettlement
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTnSettlement{},
			want: 0,
		},
		{
			name: "converts multiple settlements",
			dtos: func() []openapi.TgvalidatordTnSettlement {
				id1 := "settlement-1"
				id2 := "settlement-2"
				return []openapi.TgvalidatordTnSettlement{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SettlementsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SettlementsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSettlementAssetTransferFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSettlementAssetTransfer
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns asset transfer with zero values",
			dto:  &openapi.TgvalidatordTnSettlementAssetTransfer{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSettlementAssetTransfer {
				currencyID := "BTC"
				amount := "2.5"
				sourceID := "source-addr-123"
				destID := "dest-addr-456"
				return &openapi.TgvalidatordTnSettlementAssetTransfer{
					CurrencyID:                 &currencyID,
					Amount:                     &amount,
					SourceSharedAddressID:      &sourceID,
					DestinationSharedAddressID: &destID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementAssetTransferFromDTO(tt.dto)
			if tt.dto == nil {
				// nil input should return empty struct
				if got.CurrencyID != "" || got.Amount != "" {
					t.Errorf("SettlementAssetTransferFromDTO(nil) should return empty struct")
				}
				return
			}
			// Verify fields if set
			if tt.dto.CurrencyID != nil && got.CurrencyID != *tt.dto.CurrencyID {
				t.Errorf("CurrencyID = %v, want %v", got.CurrencyID, *tt.dto.CurrencyID)
			}
			if tt.dto.Amount != nil && got.Amount != *tt.dto.Amount {
				t.Errorf("Amount = %v, want %v", got.Amount, *tt.dto.Amount)
			}
			if tt.dto.SourceSharedAddressID != nil && got.SourceSharedAddressID != *tt.dto.SourceSharedAddressID {
				t.Errorf("SourceSharedAddressID = %v, want %v", got.SourceSharedAddressID, *tt.dto.SourceSharedAddressID)
			}
			if tt.dto.DestinationSharedAddressID != nil && got.DestinationSharedAddressID != *tt.dto.DestinationSharedAddressID {
				t.Errorf("DestinationSharedAddressID = %v, want %v", got.DestinationSharedAddressID, *tt.dto.DestinationSharedAddressID)
			}
		})
	}
}

func TestSettlementAssetTransferToDTO(t *testing.T) {
	tests := []struct {
		name     string
		transfer taurusnetwork.SettlementAssetTransfer
	}{
		{
			name:     "empty transfer",
			transfer: taurusnetwork.SettlementAssetTransfer{},
		},
		{
			name: "complete transfer",
			transfer: taurusnetwork.SettlementAssetTransfer{
				CurrencyID:                 "ETH",
				Amount:                     "10.0",
				SourceSharedAddressID:      "source-123",
				DestinationSharedAddressID: "dest-456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementAssetTransferToDTO(tt.transfer)
			if tt.transfer.CurrencyID != "" && (got.CurrencyID == nil || *got.CurrencyID != tt.transfer.CurrencyID) {
				t.Errorf("CurrencyID = %v, want %v", got.CurrencyID, tt.transfer.CurrencyID)
			}
			if tt.transfer.Amount != "" && (got.Amount == nil || *got.Amount != tt.transfer.Amount) {
				t.Errorf("Amount = %v, want %v", got.Amount, tt.transfer.Amount)
			}
		})
	}
}

func TestSettlementAssetTransfersToDTO(t *testing.T) {
	tests := []struct {
		name      string
		transfers []taurusnetwork.SettlementAssetTransfer
		want      int
	}{
		{
			name:      "nil slice returns nil",
			transfers: nil,
			want:      -1,
		},
		{
			name:      "empty slice returns empty slice",
			transfers: []taurusnetwork.SettlementAssetTransfer{},
			want:      0,
		},
		{
			name: "converts multiple transfers",
			transfers: []taurusnetwork.SettlementAssetTransfer{
				{CurrencyID: "BTC", Amount: "1.0"},
				{CurrencyID: "ETH", Amount: "2.0"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementAssetTransfersToDTO(tt.transfers)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("SettlementAssetTransfersToDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("SettlementAssetTransfersToDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSettlementClipFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSettlementClip
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns clip with zero values",
			dto:  &openapi.TgvalidatordTnSettlementClip{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSettlementClip {
				id := "clip-123"
				index := "0"
				status := "PENDING"
				workflowID := "workflow-456"
				return &openapi.TgvalidatordTnSettlementClip{
					Id:         &id,
					Index:      &index,
					Status:     &status,
					WorkflowID: &workflowID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementClipFromDTO(tt.dto)
			if tt.dto == nil {
				// nil input should return empty struct
				if got.ID != "" || got.Index != "" {
					t.Errorf("SettlementClipFromDTO(nil) should return empty struct")
				}
				return
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Index != nil && got.Index != *tt.dto.Index {
				t.Errorf("Index = %v, want %v", got.Index, *tt.dto.Index)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
		})
	}
}

func TestSettlementClipTransactionFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTnSettlementClipTransaction
	}{
		{
			name: "nil input returns empty struct",
			dto:  nil,
		},
		{
			name: "empty DTO returns transaction with zero values",
			dto:  &openapi.TgvalidatordTnSettlementClipTransaction{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTnSettlementClipTransaction {
				id := "tx-123"
				requestID := "req-456"
				txHash := "0xabc123"
				txID := "tx-id-789"
				txBlockNumber := "12345"
				status := "COMPLETED"
				workflowID := "workflow-111"
				createdAt := time.Now()
				currencyID := "BTC"
				amount := "0.5"

				return &openapi.TgvalidatordTnSettlementClipTransaction{
					Id:            &id,
					RequestID:     &requestID,
					TxHash:        &txHash,
					TxID:          &txID,
					TxBlockNumber: &txBlockNumber,
					Status:        &status,
					WorkflowID:    &workflowID,
					CreatedAt:     &createdAt,
					AssetTransfer: &openapi.TgvalidatordTnSettlementAssetTransfer{
						CurrencyID: &currencyID,
						Amount:     &amount,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SettlementClipTransactionFromDTO(tt.dto)
			if tt.dto == nil {
				// nil input should return empty struct
				if got.ID != "" || got.TxHash != "" {
					t.Errorf("SettlementClipTransactionFromDTO(nil) should return empty struct")
				}
				return
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TxHash != nil && got.TxHash != *tt.dto.TxHash {
				t.Errorf("TxHash = %v, want %v", got.TxHash, *tt.dto.TxHash)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.AssetTransfer != nil && got.AssetTransfer == nil {
				t.Error("AssetTransfer should not be nil when DTO has asset transfer")
			}
		})
	}
}

func TestCreateSettlementRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *taurusnetwork.CreateSettlementRequest
	}{
		{
			name: "nil input returns empty DTO",
			req:  nil,
		},
		{
			name: "empty request",
			req:  &taurusnetwork.CreateSettlementRequest{},
		},
		{
			name: "complete request",
			req: &taurusnetwork.CreateSettlementRequest{
				TargetParticipantID:   "target-123",
				FirstLegParticipantID: "firstleg-456",
				FirstLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "BTC", Amount: "1.0"},
				},
				SecondLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "ETH", Amount: "10.0"},
				},
				Clips: []taurusnetwork.CreateSettlementClipRequest{
					{Index: "0"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSettlementRequestToDTO(tt.req)
			if tt.req == nil {
				// nil input should return empty struct
				if got.TargetParticipantID != "" {
					t.Errorf("CreateSettlementRequestToDTO(nil) should return empty DTO")
				}
				return
			}
			if tt.req.TargetParticipantID != "" && got.TargetParticipantID != tt.req.TargetParticipantID {
				t.Errorf("TargetParticipantID = %v, want %v", got.TargetParticipantID, tt.req.TargetParticipantID)
			}
			if tt.req.FirstLegAssets != nil && len(got.FirstLegAssets) != len(tt.req.FirstLegAssets) {
				t.Errorf("FirstLegAssets length = %v, want %v", len(got.FirstLegAssets), len(tt.req.FirstLegAssets))
			}
			if tt.req.Clips != nil && len(got.Clips) != len(tt.req.Clips) {
				t.Errorf("Clips length = %v, want %v", len(got.Clips), len(tt.req.Clips))
			}
		})
	}
}

func TestCreateSettlementClipRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  taurusnetwork.CreateSettlementClipRequest
	}{
		{
			name: "empty request",
			req:  taurusnetwork.CreateSettlementClipRequest{},
		},
		{
			name: "complete request",
			req: taurusnetwork.CreateSettlementClipRequest{
				Index: "1",
				FirstLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "BTC", Amount: "0.5"},
				},
				SecondLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "ETH", Amount: "5.0"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateSettlementClipRequestToDTO(tt.req)
			if tt.req.Index != "" && (got.Index == nil || *got.Index != tt.req.Index) {
				t.Errorf("Index = %v, want %v", got.Index, tt.req.Index)
			}
			if tt.req.FirstLegAssets != nil && len(got.FirstLegAssets) != len(tt.req.FirstLegAssets) {
				t.Errorf("FirstLegAssets length = %v, want %v", len(got.FirstLegAssets), len(tt.req.FirstLegAssets))
			}
		})
	}
}

func TestReplaceSettlementRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *taurusnetwork.ReplaceSettlementRequest
	}{
		{
			name: "nil input returns empty DTO",
			req:  nil,
		},
		{
			name: "nil create request returns empty DTO",
			req:  &taurusnetwork.ReplaceSettlementRequest{CreateSettlementRequest: nil},
		},
		{
			name: "complete request",
			req: &taurusnetwork.ReplaceSettlementRequest{
				CreateSettlementRequest: &taurusnetwork.CreateSettlementRequest{
					TargetParticipantID:   "new-target",
					FirstLegParticipantID: "new-firstleg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceSettlementRequestToDTO(tt.req)
			if tt.req == nil || tt.req.CreateSettlementRequest == nil {
				if got.CreateSettlementRequest != nil {
					t.Errorf("ReplaceSettlementRequestToDTO should return DTO with nil CreateSettlementRequest")
				}
				return
			}
			if got.CreateSettlementRequest == nil {
				t.Error("CreateSettlementRequest should not be nil")
				return
			}
			if got.CreateSettlementRequest.TargetParticipantID != tt.req.CreateSettlementRequest.TargetParticipantID {
				t.Errorf("TargetParticipantID = %v, want %v", got.CreateSettlementRequest.TargetParticipantID, tt.req.CreateSettlementRequest.TargetParticipantID)
			}
		})
	}
}

func TestSettlementFromDTO_WithClips(t *testing.T) {
	clipID := "clip-123"
	clipIndex := "0"
	txID := "tx-456"
	txStatus := "COMPLETED"

	dto := &openapi.TgvalidatordTnSettlement{
		Clips: []openapi.TgvalidatordTnSettlementClip{
			{
				Id:    &clipID,
				Index: &clipIndex,
				FirstLegTransactions: []openapi.TgvalidatordTnSettlementClipTransaction{
					{Id: &txID, Status: &txStatus},
				},
			},
		},
	}

	got := SettlementFromDTO(dto)
	if got == nil {
		t.Fatal("SettlementFromDTO() returned nil for non-nil input")
	}
	if len(got.Clips) != 1 {
		t.Errorf("Clips length = %v, want 1", len(got.Clips))
	}
	if got.Clips[0].ID != clipID {
		t.Errorf("Clip ID = %v, want %v", got.Clips[0].ID, clipID)
	}
	if len(got.Clips[0].FirstLegTransactions) != 1 {
		t.Errorf("FirstLegTransactions length = %v, want 1", len(got.Clips[0].FirstLegTransactions))
	}
	if got.Clips[0].FirstLegTransactions[0].ID != txID {
		t.Errorf("Transaction ID = %v, want %v", got.Clips[0].FirstLegTransactions[0].ID, txID)
	}
}

func TestSettlementFromDTO_NilTimestamps(t *testing.T) {
	id := "settlement-123"
	dto := &openapi.TgvalidatordTnSettlement{
		Id:                 &id,
		CreatedAt:          nil,
		UpdatedAt:          nil,
		StartExecutionDate: nil,
	}

	got := SettlementFromDTO(dto)
	if got == nil {
		t.Fatal("SettlementFromDTO() returned nil for non-nil input")
	}
	if !got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero time when nil, got %v", got.CreatedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
	if !got.StartExecutionDate.IsZero() {
		t.Errorf("StartExecutionDate should be zero time when nil, got %v", got.StartExecutionDate)
	}
}
