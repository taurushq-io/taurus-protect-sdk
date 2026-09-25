package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AssetFilterToDTO converts a domain AssetFilter to an OpenAPI TgvalidatordAsset.
func AssetFilterToDTO(filter *model.AssetFilter) openapi.TgvalidatordAsset {
	if filter == nil {
		return openapi.TgvalidatordAsset{}
	}

	asset := openapi.TgvalidatordAsset{
		Currency: filter.Currency,
	}

	if filter.Kind != "" {
		asset.Kind = &filter.Kind
	}

	if filter.NFT != nil {
		nft := openapi.TgvalidatordAssetNFT{}
		if filter.NFT.TokenID != "" {
			nft.Tokenid = &filter.NFT.TokenID
		}
		asset.Nft = &nft
	}

	if filter.Unknown != nil {
		unknown := openapi.AssetUnknown{}
		if filter.Unknown.Blockchain != "" {
			unknown.Blockchain = &filter.Unknown.Blockchain
		}
		if filter.Unknown.Arg1 != "" {
			unknown.Arg1 = &filter.Unknown.Arg1
		}
		if filter.Unknown.Arg2 != "" {
			unknown.Arg2 = &filter.Unknown.Arg2
		}
		if filter.Unknown.Network != "" {
			unknown.Network = &filter.Unknown.Network
		}
		asset.Unknown = &unknown
	}

	return asset
}

// AssetFilterFromDTO converts an OpenAPI TgvalidatordAsset to a domain AssetFilter.
func AssetFilterFromDTO(dto *openapi.TgvalidatordAsset) *model.AssetFilter {
	if dto == nil {
		return nil
	}

	filter := &model.AssetFilter{
		Currency: dto.Currency,
		Kind:     safeString(dto.Kind),
	}

	if dto.Nft != nil {
		filter.NFT = &model.AssetNFTFilter{
			TokenID: safeString(dto.Nft.Tokenid),
		}
	}

	if dto.Unknown != nil {
		filter.Unknown = &model.AssetUnknownFilter{
			Blockchain: safeString(dto.Unknown.Blockchain),
			Arg1:       safeString(dto.Unknown.Arg1),
			Arg2:       safeString(dto.Unknown.Arg2),
			Network:    safeString(dto.Unknown.Network),
		}
	}

	return filter
}

// AssetResourceFromDTO converts a v2 asset to a domain AssetResource.
func AssetResourceFromDTO(dto *openapi.TgvalidatordAssetResourceV2) *model.AssetResource {
	if dto == nil {
		return nil
	}
	asset := &model.AssetResource{
		ID:              safeString(dto.Id),
		TenantID:        safeString(dto.TenantID),
		Version:         safeString(dto.Version),
		CreatedAt:       safeTime(dto.CreatedAt),
		UpdatedAt:       safeTime(dto.UpdatedAt),
		Label:           safeString(dto.Label),
		AssetType:       safeString(dto.AssetType),
		Blockchain:      safeString(dto.Blockchain),
		Network:         safeString(dto.Network),
		CurrencyID:      safeString(dto.CurrencyID),
		Name:            safeString(dto.Name),
		Symbol:          safeString(dto.Symbol),
		Decimals:        safeString(dto.Decimals),
		ContractAddress: safeString(dto.ContractAddress),
	}
	if dto.Status != nil {
		asset.Status = string(*dto.Status)
	}
	for _, attr := range dto.Attributes {
		asset.Attributes = append(asset.Attributes, model.AssetAttribute{
			Key:   safeString(attr.Key),
			Value: safeString(attr.Value),
		})
	}
	if dto.BlockchainAsset != nil && dto.BlockchainAsset.CantonNativeTokenAsset != nil {
		canton := dto.BlockchainAsset.CantonNativeTokenAsset
		instrument := &model.CantonInstrument{InstrumentID: safeString(canton.InstrumentID)}
		if cfg := canton.Configuration; cfg != nil {
			instrument.ContractID = safeString(cfg.Cid)
			instrument.RequireCredentials = safeBool(cfg.RequireCredentials)
			instrument.Paused = safeBool(cfg.Paused)
			instrument.Operator = safeString(cfg.Operator)
		}
		asset.CantonInstrument = instrument
	}
	return asset
}

// AssetResourcesFromDTO converts a page of v2 assets.
func AssetResourcesFromDTO(dtos []openapi.TgvalidatordAssetResourceV2) []*model.AssetResource {
	if dtos == nil {
		return nil
	}
	assets := make([]*model.AssetResource, len(dtos))
	for i := range dtos {
		assets[i] = AssetResourceFromDTO(&dtos[i])
	}
	return assets
}

// AssetAddressFromDTO converts a v2 asset address to a domain AssetAddress.
func AssetAddressFromDTO(dto *openapi.TgvalidatordAssetAddressV2) *model.AssetAddress {
	if dto == nil {
		return nil
	}
	address := &model.AssetAddress{
		Address:              safeString(dto.Address),
		Balance:              safeString(dto.Balance),
		AddressID:            safeString(dto.AddressID),
		WhitelistedAddressID: safeString(dto.WhitelistedAddressID),
	}
	if dto.KycStatus != nil {
		address.KYCStatus = string(*dto.KycStatus)
	}
	if dto.AddressType != nil {
		address.AddressType = string(*dto.AddressType)
	}
	return address
}

// AssetAddressesFromDTO converts a page of v2 asset addresses.
func AssetAddressesFromDTO(dtos []openapi.TgvalidatordAssetAddressV2) []*model.AssetAddress {
	if dtos == nil {
		return nil
	}
	addresses := make([]*model.AssetAddress, len(dtos))
	for i := range dtos {
		addresses[i] = AssetAddressFromDTO(&dtos[i])
	}
	return addresses
}

// AssetOperationFromDTO converts a v2 asset operation to a domain AssetOperation, flattening
// the per-type details.
func AssetOperationFromDTO(dto *openapi.TgvalidatordAssetOperationV2) *model.AssetOperation {
	if dto == nil {
		return nil
	}
	op := &model.AssetOperation{
		ID:                   safeString(dto.Id),
		AssetID:              safeString(dto.AssetID),
		CreatedAt:            safeTime(dto.CreatedAt),
		UpdatedAt:            safeTime(dto.UpdatedAt),
		InitiatedByAddressID: safeString(dto.InitiatedByAddressID),
	}
	if dto.Type != nil {
		op.Type = string(*dto.Type)
	}
	if dto.Status != nil {
		op.Status = string(*dto.Status)
	}
	if dto.FailureReason != nil {
		op.FailureReason = string(*dto.FailureReason)
	}
	if dto.BlockingReason != nil {
		op.BlockingReason = string(*dto.BlockingReason)
	}
	if c := dto.Create; c != nil {
		op.Label, op.Price, op.Decimals = safeString(c.Label), safeString(c.Price), safeString(c.Decimals)
		op.Blockchain, op.Network, op.AssetType = safeString(c.Blockchain), safeString(c.Network), safeString(c.AssetType)
	}
	if u := dto.Update; u != nil {
		op.Label, op.Price = safeString(u.Label), safeString(u.Price)
	}
	if i := dto.Import; i != nil {
		op.Label, op.Price, op.Decimals = safeString(i.Label), safeString(i.Price), safeString(i.Decimals)
		op.Blockchain, op.Network, op.Address = safeString(i.Blockchain), safeString(i.Network), safeString(i.Address)
	}
	if m := dto.Mint; m != nil {
		op.Amount, op.NFTMetadata = safeString(m.Amount), m.NftMetadata
		op.Target = assetOperationTargetFromDTO(m.Destination)
	}
	if b := dto.Burn; b != nil {
		op.Amount, op.NFTTokenIDs = safeString(b.Amount), b.NftTokenIDs
		op.Target = assetOperationTargetFromDTO(b.Destination)
	}
	if p := dto.PauseAccount; p != nil {
		op.Target = assetOperationTargetFromDTO(p.Target)
	}
	if p := dto.UnpauseAccount; p != nil {
		op.Target = assetOperationTargetFromDTO(p.Target)
	}
	if k := dto.SetKyc; k != nil {
		op.Target = assetOperationTargetFromDTO(k.Target)
		if k.Status != nil {
			op.KYCStatus = string(*k.Status)
		}
	}
	return op
}

func assetOperationTargetFromDTO(dto *openapi.TgvalidatordAddressTargetV2) *model.AssetOperationTarget {
	if dto == nil {
		return nil
	}
	return &model.AssetOperationTarget{
		AddressID:            safeString(dto.AddressID),
		WhitelistedAddressID: safeString(dto.WhitelistedAddressID),
	}
}

// AssetOperationsFromDTO converts a page of v2 asset operations.
func AssetOperationsFromDTO(dtos []openapi.TgvalidatordAssetOperationV2) []*model.AssetOperation {
	if dtos == nil {
		return nil
	}
	ops := make([]*model.AssetOperation, len(dtos))
	for i := range dtos {
		ops[i] = AssetOperationFromDTO(&dtos[i])
	}
	return ops
}
