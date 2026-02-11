package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// SettlementFromDTO converts an OpenAPI TnSettlement to a domain Settlement.
func SettlementFromDTO(dto *openapi.TgvalidatordTnSettlement) *taurusnetwork.Settlement {
	if dto == nil {
		return nil
	}

	settlement := &taurusnetwork.Settlement{
		ID:                    safeString(dto.Id),
		CreatorParticipantID:  safeString(dto.CreatorParticipantID),
		TargetParticipantID:   safeString(dto.TargetParticipantID),
		FirstLegParticipantID: safeString(dto.FirstLegParticipantID),
		Status:                safeString(dto.Status),
		WorkflowID:            safeString(dto.WorkflowID),
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		settlement.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		settlement.UpdatedAt = *dto.UpdatedAt
	}
	if dto.StartExecutionDate != nil {
		settlement.StartExecutionDate = *dto.StartExecutionDate
	}

	// Convert first leg assets
	if dto.FirstLegAssets != nil {
		settlement.FirstLegAssets = make([]taurusnetwork.SettlementAssetTransfer, len(dto.FirstLegAssets))
		for i, asset := range dto.FirstLegAssets {
			settlement.FirstLegAssets[i] = SettlementAssetTransferFromDTO(&asset)
		}
	}

	// Convert second leg assets
	if dto.SecondLegAssets != nil {
		settlement.SecondLegAssets = make([]taurusnetwork.SettlementAssetTransfer, len(dto.SecondLegAssets))
		for i, asset := range dto.SecondLegAssets {
			settlement.SecondLegAssets[i] = SettlementAssetTransferFromDTO(&asset)
		}
	}

	// Convert clips
	if dto.Clips != nil {
		settlement.Clips = make([]taurusnetwork.SettlementClip, len(dto.Clips))
		for i, clip := range dto.Clips {
			settlement.Clips[i] = SettlementClipFromDTO(&clip)
		}
	}

	return settlement
}

// SettlementsFromDTO converts a slice of OpenAPI TnSettlement to domain Settlements.
func SettlementsFromDTO(dtos []openapi.TgvalidatordTnSettlement) []*taurusnetwork.Settlement {
	if dtos == nil {
		return nil
	}
	settlements := make([]*taurusnetwork.Settlement, len(dtos))
	for i := range dtos {
		settlements[i] = SettlementFromDTO(&dtos[i])
	}
	return settlements
}

// SettlementAssetTransferFromDTO converts an OpenAPI TnSettlementAssetTransfer to a domain SettlementAssetTransfer.
func SettlementAssetTransferFromDTO(dto *openapi.TgvalidatordTnSettlementAssetTransfer) taurusnetwork.SettlementAssetTransfer {
	if dto == nil {
		return taurusnetwork.SettlementAssetTransfer{}
	}
	return taurusnetwork.SettlementAssetTransfer{
		CurrencyID:                 safeString(dto.CurrencyID),
		Amount:                     safeString(dto.Amount),
		SourceSharedAddressID:      safeString(dto.SourceSharedAddressID),
		DestinationSharedAddressID: safeString(dto.DestinationSharedAddressID),
	}
}

// SettlementAssetTransferToDTO converts a domain SettlementAssetTransfer to an OpenAPI TnSettlementAssetTransfer.
func SettlementAssetTransferToDTO(transfer taurusnetwork.SettlementAssetTransfer) openapi.TgvalidatordTnSettlementAssetTransfer {
	return openapi.TgvalidatordTnSettlementAssetTransfer{
		CurrencyID:                 stringPtr(transfer.CurrencyID),
		Amount:                     stringPtr(transfer.Amount),
		SourceSharedAddressID:      stringPtr(transfer.SourceSharedAddressID),
		DestinationSharedAddressID: stringPtr(transfer.DestinationSharedAddressID),
	}
}

// SettlementAssetTransfersToDTO converts a slice of domain SettlementAssetTransfer to OpenAPI TnSettlementAssetTransfer.
func SettlementAssetTransfersToDTO(transfers []taurusnetwork.SettlementAssetTransfer) []openapi.TgvalidatordTnSettlementAssetTransfer {
	if transfers == nil {
		return nil
	}
	result := make([]openapi.TgvalidatordTnSettlementAssetTransfer, len(transfers))
	for i, transfer := range transfers {
		result[i] = SettlementAssetTransferToDTO(transfer)
	}
	return result
}

// SettlementClipFromDTO converts an OpenAPI TnSettlementClip to a domain SettlementClip.
func SettlementClipFromDTO(dto *openapi.TgvalidatordTnSettlementClip) taurusnetwork.SettlementClip {
	if dto == nil {
		return taurusnetwork.SettlementClip{}
	}

	clip := taurusnetwork.SettlementClip{
		ID:         safeString(dto.Id),
		Index:      safeString(dto.Index),
		Status:     safeString(dto.Status),
		WorkflowID: safeString(dto.WorkflowID),
	}

	// Convert first leg transactions
	if dto.FirstLegTransactions != nil {
		clip.FirstLegTransactions = make([]taurusnetwork.SettlementClipTransaction, len(dto.FirstLegTransactions))
		for i, tx := range dto.FirstLegTransactions {
			clip.FirstLegTransactions[i] = SettlementClipTransactionFromDTO(&tx)
		}
	}

	// Convert second leg transactions
	if dto.SecondLegTransactions != nil {
		clip.SecondLegTransactions = make([]taurusnetwork.SettlementClipTransaction, len(dto.SecondLegTransactions))
		for i, tx := range dto.SecondLegTransactions {
			clip.SecondLegTransactions[i] = SettlementClipTransactionFromDTO(&tx)
		}
	}

	return clip
}

// SettlementClipTransactionFromDTO converts an OpenAPI TnSettlementClipTransaction to a domain SettlementClipTransaction.
func SettlementClipTransactionFromDTO(dto *openapi.TgvalidatordTnSettlementClipTransaction) taurusnetwork.SettlementClipTransaction {
	if dto == nil {
		return taurusnetwork.SettlementClipTransaction{}
	}

	tx := taurusnetwork.SettlementClipTransaction{
		ID:            safeString(dto.Id),
		RequestID:     safeString(dto.RequestID),
		TxHash:        safeString(dto.TxHash),
		TxID:          safeString(dto.TxID),
		TxBlockNumber: safeString(dto.TxBlockNumber),
		Status:        safeString(dto.Status),
		WorkflowID:    safeString(dto.WorkflowID),
	}

	if dto.CreatedAt != nil {
		tx.CreatedAt = *dto.CreatedAt
	}

	if dto.AssetTransfer != nil {
		assetTransfer := SettlementAssetTransferFromDTO(dto.AssetTransfer)
		tx.AssetTransfer = &assetTransfer
	}

	return tx
}

// CreateSettlementRequestToDTO converts a domain CreateSettlementRequest to an OpenAPI TgvalidatordCreateSettlementRequest.
func CreateSettlementRequestToDTO(req *taurusnetwork.CreateSettlementRequest) openapi.TgvalidatordCreateSettlementRequest {
	if req == nil {
		return openapi.TgvalidatordCreateSettlementRequest{}
	}

	dto := openapi.TgvalidatordCreateSettlementRequest{
		TargetParticipantID:   req.TargetParticipantID,
		FirstLegParticipantID: req.FirstLegParticipantID,
		FirstLegAssets:        SettlementAssetTransfersToDTO(req.FirstLegAssets),
		SecondLegAssets:       SettlementAssetTransfersToDTO(req.SecondLegAssets),
	}

	if req.StartExecutionDate != nil {
		dto.StartExecutionDate = req.StartExecutionDate
	}

	if req.Clips != nil {
		dto.Clips = make([]openapi.CreateSettlementRequestClipRequest, len(req.Clips))
		for i, clip := range req.Clips {
			dto.Clips[i] = CreateSettlementClipRequestToDTO(clip)
		}
	}

	return dto
}

// CreateSettlementClipRequestToDTO converts a domain CreateSettlementClipRequest to an OpenAPI CreateSettlementRequestClipRequest.
func CreateSettlementClipRequestToDTO(req taurusnetwork.CreateSettlementClipRequest) openapi.CreateSettlementRequestClipRequest {
	return openapi.CreateSettlementRequestClipRequest{
		Index:           stringPtr(req.Index),
		FirstLegAssets:  SettlementAssetTransfersToDTO(req.FirstLegAssets),
		SecondLegAssets: SettlementAssetTransfersToDTO(req.SecondLegAssets),
	}
}

// ReplaceSettlementRequestToDTO converts a domain ReplaceSettlementRequest to an OpenAPI TaurusNetworkServiceReplaceSettlementBody.
func ReplaceSettlementRequestToDTO(req *taurusnetwork.ReplaceSettlementRequest) openapi.TaurusNetworkServiceReplaceSettlementBody {
	if req == nil || req.CreateSettlementRequest == nil {
		return openapi.TaurusNetworkServiceReplaceSettlementBody{}
	}

	createReq := CreateSettlementRequestToDTO(req.CreateSettlementRequest)
	return openapi.TaurusNetworkServiceReplaceSettlementBody{
		CreateSettlementRequest: &createReq,
	}
}
