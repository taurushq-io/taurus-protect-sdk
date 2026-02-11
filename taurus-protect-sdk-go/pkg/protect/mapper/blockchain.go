package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// BlockchainFromDTO converts an OpenAPI TgvalidatordBlockchainEntity to a domain Blockchain.
func BlockchainFromDTO(dto *openapi.TgvalidatordBlockchainEntity) *model.Blockchain {
	if dto == nil {
		return nil
	}

	blockchain := &model.Blockchain{
		Symbol:           safeString(dto.Symbol),
		Name:             safeString(dto.Name),
		Network:          safeString(dto.Network),
		ChainID:          safeString(dto.ChainId),
		BlackholeAddress: safeString(dto.BlackholeAddress),
		Confirmations:    safeString(dto.Confirmations),
		BlockHeight:      safeString(dto.BlockHeight),
		IsLayer2Chain:    safeBool(dto.IsLayer2Chain),
		Layer1Network:    safeString(dto.Layer1Network),
	}

	// Convert base currency
	if dto.BaseCurrency != nil {
		blockchain.BaseCurrency = CurrencyFromDTO(dto.BaseCurrency)
	}

	// Convert EVM info
	if dto.EthInfo != nil {
		blockchain.EVMInfo = EVMBlockchainInfoFromDTO(dto.EthInfo)
	}

	// Convert DOT info
	if dto.DotInfo != nil {
		blockchain.DOTInfo = DOTBlockchainInfoFromDTO(dto.DotInfo)
	}

	// Convert XTZ info
	if dto.XtzInfo != nil {
		blockchain.XTZInfo = XTZBlockchainInfoFromDTO(dto.XtzInfo)
	}

	return blockchain
}

// BlockchainsFromDTO converts a slice of OpenAPI TgvalidatordBlockchainEntity to domain Blockchains.
func BlockchainsFromDTO(dtos []openapi.TgvalidatordBlockchainEntity) []*model.Blockchain {
	if dtos == nil {
		return nil
	}
	blockchains := make([]*model.Blockchain, len(dtos))
	for i := range dtos {
		blockchains[i] = BlockchainFromDTO(&dtos[i])
	}
	return blockchains
}

// EVMBlockchainInfoFromDTO converts an OpenAPI TgvalidatordEVMBlockchainInfo to a domain EVMBlockchainInfo.
func EVMBlockchainInfoFromDTO(dto *openapi.TgvalidatordEVMBlockchainInfo) *model.EVMBlockchainInfo {
	if dto == nil {
		return nil
	}

	return &model.EVMBlockchainInfo{
		ChainID: safeString(dto.ChainId),
	}
}

// DOTBlockchainInfoFromDTO converts an OpenAPI TgvalidatordDOTBlockchainInfo to a domain DOTBlockchainInfo.
func DOTBlockchainInfoFromDTO(dto *openapi.TgvalidatordDOTBlockchainInfo) *model.DOTBlockchainInfo {
	if dto == nil {
		return nil
	}

	info := &model.DOTBlockchainInfo{
		CurrentEra:     safeString(dto.CurrentEra),
		MaxNominations: safeString(dto.MaxNominations),
		ForkNumber:     safeString(dto.ForkNumber),
	}

	if dto.ForkMigratedAt != nil {
		info.ForkMigratedAt = dto.ForkMigratedAt
	}

	return info
}

// XTZBlockchainInfoFromDTO converts an OpenAPI TgvalidatordXTZBlockchainInfo to a domain XTZBlockchainInfo.
func XTZBlockchainInfoFromDTO(dto *openapi.TgvalidatordXTZBlockchainInfo) *model.XTZBlockchainInfo {
	if dto == nil {
		return nil
	}

	return &model.XTZBlockchainInfo{
		CurrentCycle: safeString(dto.CurrentCycle),
	}
}
