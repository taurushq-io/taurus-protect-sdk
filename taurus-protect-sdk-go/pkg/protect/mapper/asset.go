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
