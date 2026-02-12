package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// SharedAddressFromDTO converts an OpenAPI TnSharedAddress to a domain SharedAddress.
func SharedAddressFromDTO(dto *openapi.TgvalidatordTnSharedAddress) *taurusnetwork.SharedAddress {
	if dto == nil {
		return nil
	}

	addr := &taurusnetwork.SharedAddress{
		ID:                   safeString(dto.Id),
		InternalAddressID:    safeString(dto.InternalAddressID),
		WhitelistedAddressID: safeString(dto.WladdressID),
		OwnerParticipantID:   safeString(dto.OwnerParticipantId),
		TargetParticipantID:  safeString(dto.TargetParticipantId),
		Blockchain:           safeString(dto.Blockchain),
		Network:              safeString(dto.Network),
		Address:              safeString(dto.Address),
		OriginLabel:          safeString(dto.OriginLabel),
		Status:               taurusnetwork.SharedAddressStatus(safeString(dto.Status)),
	}

	// Convert timestamps
	if dto.OriginCreationDate != nil {
		addr.OriginCreationDate = *dto.OriginCreationDate
	}
	if dto.OriginDeletionDate != nil {
		addr.OriginDeletionDate = *dto.OriginDeletionDate
	}
	if dto.CreatedAt != nil {
		addr.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		addr.UpdatedAt = *dto.UpdatedAt
	}
	if dto.TargetAcceptedAt != nil {
		addr.TargetAcceptedAt = *dto.TargetAcceptedAt
	}

	// Parse pledges count
	if dto.PledgesCount != nil {
		if count, err := strconv.ParseInt(*dto.PledgesCount, 10, 64); err == nil {
			addr.PledgesCount = count
		}
	}

	// Convert proof of ownership
	if dto.ProofOfOwnership != nil {
		addr.ProofOfOwnership = ProofOfOwnershipFromDTO(dto.ProofOfOwnership)
	}

	// Convert trails
	if dto.Trails != nil {
		addr.Trails = SharedAddressTrailsFromDTO(dto.Trails)
	}

	return addr
}

// SharedAddressesFromDTO converts a slice of OpenAPI TnSharedAddress to domain SharedAddresses.
func SharedAddressesFromDTO(dtos []openapi.TgvalidatordTnSharedAddress) []*taurusnetwork.SharedAddress {
	if dtos == nil {
		return nil
	}
	addresses := make([]*taurusnetwork.SharedAddress, len(dtos))
	for i := range dtos {
		addresses[i] = SharedAddressFromDTO(&dtos[i])
	}
	return addresses
}

// SharedAddressTrailFromDTO converts an OpenAPI TnSharedAddressTrail to a domain SharedAddressTrail.
func SharedAddressTrailFromDTO(dto *openapi.TgvalidatordTnSharedAddressTrail) *taurusnetwork.SharedAddressTrail {
	if dto == nil {
		return nil
	}

	trail := &taurusnetwork.SharedAddressTrail{
		ID:              safeString(dto.Id),
		SharedAddressID: safeString(dto.SharedAddressID),
		AddressStatus:   safeString(dto.AddressStatus),
		Comment:         safeString(dto.Comment),
	}

	if dto.CreatedAt != nil {
		trail.CreatedAt = *dto.CreatedAt
	}

	return trail
}

// SharedAddressTrailsFromDTO converts a slice of OpenAPI TnSharedAddressTrail to domain SharedAddressTrails.
func SharedAddressTrailsFromDTO(dtos []openapi.TgvalidatordTnSharedAddressTrail) []*taurusnetwork.SharedAddressTrail {
	if dtos == nil {
		return nil
	}
	trails := make([]*taurusnetwork.SharedAddressTrail, len(dtos))
	for i := range dtos {
		trails[i] = SharedAddressTrailFromDTO(&dtos[i])
	}
	return trails
}

// ProofOfOwnershipFromDTO converts an OpenAPI ProofOfOwnership to a domain ProofOfOwnership.
func ProofOfOwnershipFromDTO(dto *openapi.TgvalidatordProofOfOwnership) *taurusnetwork.ProofOfOwnership {
	if dto == nil {
		return nil
	}

	return &taurusnetwork.ProofOfOwnership{
		SignedPayloadHash:     safeString(dto.SignedPayloadHash),
		SignedPayloadAsString: safeString(dto.SignedPayloadAsString),
	}
}

// SharedAssetFromDTO converts an OpenAPI TnSharedAsset to a domain SharedAsset.
func SharedAssetFromDTO(dto *openapi.TgvalidatordTnSharedAsset) *taurusnetwork.SharedAsset {
	if dto == nil {
		return nil
	}

	asset := &taurusnetwork.SharedAsset{
		ID:                           safeString(dto.Id),
		WhitelistedContractAddressID: safeString(dto.WlContractAddressID),
		OwnerParticipantID:           safeString(dto.OwnerParticipantId),
		TargetParticipantID:          safeString(dto.TargetParticipantId),
		Blockchain:                   safeString(dto.Blockchain),
		Network:                      safeString(dto.Network),
		Name:                         safeString(dto.Name),
		Symbol:                       safeString(dto.Symbol),
		Decimals:                     safeString(dto.Decimals),
		ContractAddress:              safeString(dto.ContractAddress),
		TokenID:                      safeString(dto.TokenId),
		Kind:                         safeString(dto.Kind),
		Status:                       taurusnetwork.SharedAssetStatus(safeString(dto.Status)),
	}

	// Convert timestamps
	if dto.OriginCreationDate != nil {
		asset.OriginCreationDate = *dto.OriginCreationDate
	}
	if dto.OriginDeletionDate != nil {
		asset.OriginDeletionDate = *dto.OriginDeletionDate
	}
	if dto.CreatedAt != nil {
		asset.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		asset.UpdatedAt = *dto.UpdatedAt
	}
	if dto.TargetAcceptedAt != nil {
		asset.TargetAcceptedAt = *dto.TargetAcceptedAt
	}
	if dto.TargetRejectedAt != nil {
		asset.TargetRejectedAt = *dto.TargetRejectedAt
	}

	// Convert trails
	if dto.Trails != nil {
		asset.Trails = SharedAssetTrailsFromDTO(dto.Trails)
	}

	return asset
}

// SharedAssetsFromDTO converts a slice of OpenAPI TnSharedAsset to domain SharedAssets.
func SharedAssetsFromDTO(dtos []openapi.TgvalidatordTnSharedAsset) []*taurusnetwork.SharedAsset {
	if dtos == nil {
		return nil
	}
	assets := make([]*taurusnetwork.SharedAsset, len(dtos))
	for i := range dtos {
		assets[i] = SharedAssetFromDTO(&dtos[i])
	}
	return assets
}

// SharedAssetTrailFromDTO converts an OpenAPI TnSharedAssetTrail to a domain SharedAssetTrail.
func SharedAssetTrailFromDTO(dto *openapi.TgvalidatordTnSharedAssetTrail) *taurusnetwork.SharedAssetTrail {
	if dto == nil {
		return nil
	}

	trail := &taurusnetwork.SharedAssetTrail{
		ID:            safeString(dto.Id),
		SharedAssetID: safeString(dto.SharedAssetID),
		AssetStatus:   safeString(dto.AssetStatus),
		Comment:       safeString(dto.Comment),
	}

	if dto.CreatedAt != nil {
		trail.CreatedAt = *dto.CreatedAt
	}

	return trail
}

// SharedAssetTrailsFromDTO converts a slice of OpenAPI TnSharedAssetTrail to domain SharedAssetTrails.
func SharedAssetTrailsFromDTO(dtos []openapi.TgvalidatordTnSharedAssetTrail) []*taurusnetwork.SharedAssetTrail {
	if dtos == nil {
		return nil
	}
	trails := make([]*taurusnetwork.SharedAssetTrail, len(dtos))
	for i := range dtos {
		trails[i] = SharedAssetTrailFromDTO(&dtos[i])
	}
	return trails
}

// KeyValueAttributesToDTO converts domain KeyValueAttributes to OpenAPI KeyValues.
func KeyValueAttributesToDTO(attrs []taurusnetwork.KeyValueAttribute) []openapi.TgvalidatordKeyValue {
	if attrs == nil {
		return nil
	}
	result := make([]openapi.TgvalidatordKeyValue, len(attrs))
	for i, attr := range attrs {
		result[i] = openapi.TgvalidatordKeyValue{
			Key:   stringPtr(attr.Key),
			Value: stringPtr(attr.Value),
		}
	}
	return result
}
