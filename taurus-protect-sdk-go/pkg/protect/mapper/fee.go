package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FeeFromDTO converts an OpenAPI KeyValue to a domain Fee.
// This is used for the deprecated v1 GetFees endpoint.
func FeeFromDTO(dto *openapi.TgvalidatordKeyValue) *model.Fee {
	if dto == nil {
		return nil
	}

	return &model.Fee{
		Key:   safeString(dto.Key),
		Value: safeString(dto.Value),
	}
}

// FeesFromDTO converts a slice of OpenAPI KeyValue to domain Fees.
func FeesFromDTO(dtos []openapi.TgvalidatordKeyValue) []*model.Fee {
	if dtos == nil {
		return nil
	}
	fees := make([]*model.Fee, len(dtos))
	for i := range dtos {
		fees[i] = FeeFromDTO(&dtos[i])
	}
	return fees
}

// FeeV2FromDTO converts an OpenAPI Fee to a domain FeeV2.
// This is used for the v2 GetFees endpoint.
func FeeV2FromDTO(dto *openapi.TgvalidatordFee) *model.FeeV2 {
	if dto == nil {
		return nil
	}

	fee := &model.FeeV2{
		CurrencyID:   safeString(dto.CurrencyId),
		Value:        safeString(dto.Value),
		Denom:        safeString(dto.Denom),
		CurrencyInfo: CurrencyInfoFromDTO(dto.CurrencyInfo),
	}

	if dto.UpdateDate != nil {
		fee.UpdateDate = *dto.UpdateDate
	}

	return fee
}

// FeesV2FromDTO converts a slice of OpenAPI Fee to domain FeeV2s.
func FeesV2FromDTO(dtos []openapi.TgvalidatordFee) []*model.FeeV2 {
	if dtos == nil {
		return nil
	}
	fees := make([]*model.FeeV2, len(dtos))
	for i := range dtos {
		fees[i] = FeeV2FromDTO(&dtos[i])
	}
	return fees
}
