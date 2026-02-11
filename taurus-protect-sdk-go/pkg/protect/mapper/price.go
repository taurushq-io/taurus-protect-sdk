package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// PriceFromDTO converts an OpenAPI CurrencyPrice to a domain Price.
func PriceFromDTO(dto *openapi.TgvalidatordCurrencyPrice) *model.Price {
	if dto == nil {
		return nil
	}

	price := &model.Price{
		Blockchain:          safeString(dto.Blockchain),
		CurrencyFrom:        safeString(dto.CurrencyFrom),
		CurrencyTo:          safeString(dto.CurrencyTo),
		Decimals:            safeString(dto.Decimals),
		Rate:                safeString(dto.Rate),
		ChangePercent24Hour: safeString(dto.ChangePercent24Hour),
		Source:              safeString(dto.Source),
		CurrencyFromInfo:    CurrencyInfoFromDTO(dto.CurrencyFromInfo),
		CurrencyToInfo:      CurrencyInfoFromDTO(dto.CurrencyToInfo),
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		price.CreationDate = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		price.UpdateDate = *dto.UpdateDate
	}

	// Convert signatures
	if dto.Signatures != nil {
		price.Signatures = PriceSignaturesFromDTO(dto.Signatures)
	}

	return price
}

// PricesFromDTO converts a slice of OpenAPI CurrencyPrice to domain Prices.
func PricesFromDTO(dtos []openapi.TgvalidatordCurrencyPrice) []*model.Price {
	if dtos == nil {
		return nil
	}
	prices := make([]*model.Price, len(dtos))
	for i := range dtos {
		prices[i] = PriceFromDTO(&dtos[i])
	}
	return prices
}

// PriceSignatureFromDTO converts an OpenAPI CurrencyPriceSignature to a domain PriceSignature.
func PriceSignatureFromDTO(dto *openapi.TgvalidatordCurrencyPriceSignature) model.PriceSignature {
	if dto == nil {
		return model.PriceSignature{}
	}
	return model.PriceSignature{
		UserID:    dto.UserId,
		Signature: dto.Signature,
	}
}

// PriceSignaturesFromDTO converts a slice of OpenAPI CurrencyPriceSignature to domain PriceSignatures.
func PriceSignaturesFromDTO(dtos []openapi.TgvalidatordCurrencyPriceSignature) []model.PriceSignature {
	if dtos == nil {
		return nil
	}
	signatures := make([]model.PriceSignature, len(dtos))
	for i := range dtos {
		signatures[i] = PriceSignatureFromDTO(&dtos[i])
	}
	return signatures
}

// PriceHistoryPointFromDTO converts an OpenAPI PricesHistoryPoint to a domain PriceHistoryPoint.
func PriceHistoryPointFromDTO(dto *openapi.TgvalidatordPricesHistoryPoint) *model.PriceHistoryPoint {
	if dto == nil {
		return nil
	}

	point := &model.PriceHistoryPoint{
		Blockchain:       safeString(dto.Blockchain),
		CurrencyFrom:     safeString(dto.CurrencyFrom),
		CurrencyTo:       safeString(dto.CurrencyTo),
		High:             safeString(dto.High),
		Low:              safeString(dto.Low),
		Open:             safeString(dto.Open),
		Close:            safeString(dto.Close),
		VolumeFrom:       safeString(dto.VolumeFrom),
		VolumeTo:         safeString(dto.VolumeTo),
		ChangePercent:    safeString(dto.ChangePercent),
		CurrencyFromInfo: CurrencyInfoFromDTO(dto.CurrencyFromInfo),
		CurrencyToInfo:   CurrencyInfoFromDTO(dto.CurrencyToInfo),
	}

	// Convert timestamp
	if dto.PeriodStartDate != nil {
		point.PeriodStartDate = *dto.PeriodStartDate
	}

	return point
}

// PriceHistoryPointsFromDTO converts a slice of OpenAPI PricesHistoryPoint to domain PriceHistoryPoints.
func PriceHistoryPointsFromDTO(dtos []openapi.TgvalidatordPricesHistoryPoint) []*model.PriceHistoryPoint {
	if dtos == nil {
		return nil
	}
	points := make([]*model.PriceHistoryPoint, len(dtos))
	for i := range dtos {
		points[i] = PriceHistoryPointFromDTO(&dtos[i])
	}
	return points
}

// ConversionValueFromDTO converts an OpenAPI ConversionValue to a domain ConversionValue.
func ConversionValueFromDTO(dto *openapi.TgvalidatordConversionValue) *model.ConversionValue {
	if dto == nil {
		return nil
	}
	return &model.ConversionValue{
		Symbol:        safeString(dto.Symbol),
		Value:         safeString(dto.Value),
		MainUnitValue: safeString(dto.MainUnitValue),
		CurrencyInfo:  CurrencyInfoFromDTO(dto.CurrencyInfo),
	}
}

// ConversionValuesFromDTO converts a slice of OpenAPI ConversionValue to domain ConversionValues.
func ConversionValuesFromDTO(dtos []openapi.TgvalidatordConversionValue) []*model.ConversionValue {
	if dtos == nil {
		return nil
	}
	values := make([]*model.ConversionValue, len(dtos))
	for i := range dtos {
		values[i] = ConversionValueFromDTO(&dtos[i])
	}
	return values
}

// ConversionResultFromDTO converts an OpenAPI ConversionReply to a domain ConversionResult.
func ConversionResultFromDTO(dto *openapi.TgvalidatordConversionReply) *model.ConversionResult {
	if dto == nil {
		return nil
	}
	return &model.ConversionResult{
		CurrencyFrom:     safeString(dto.CurrencyFrom),
		BaseCurrency:     safeString(dto.BaseCurrency),
		Values:           ConversionValuesFromDTO(dto.Result),
		FullCurrencyFrom: CurrencyInfoFromDTO(dto.FullCurrencyFrom),
		FullBaseCurrency: CurrencyInfoFromDTO(dto.FullBaseCurrency),
	}
}
