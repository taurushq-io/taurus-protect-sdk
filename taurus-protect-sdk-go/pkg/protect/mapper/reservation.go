package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ReservationFromDTO converts an OpenAPI Reservation to a domain Reservation.
func ReservationFromDTO(dto *openapi.TgvalidatordReservation) *model.Reservation {
	if dto == nil {
		return nil
	}

	reservation := &model.Reservation{
		ID:           safeString(dto.Id),
		Kind:         safeString(dto.Kind),
		Comment:      safeString(dto.Comment),
		AddressID:    safeString(dto.Addressid),
		Address:      safeString(dto.Address),
		Amount:       safeString(dto.Amount),
		ResourceID:   safeString(dto.ResourceId),
		ResourceType: safeString(dto.ResourceType),
	}

	// Convert creation date
	if dto.CreationDate != nil {
		reservation.CreationDate = *dto.CreationDate
	}

	// Convert currency info
	if dto.CurrencyInfo != nil {
		reservation.CurrencyInfo = CurrencyInfoFromDTO(dto.CurrencyInfo)
	}

	// Convert asset
	if dto.Asset != nil {
		reservation.Asset = AssetFromDTO(dto.Asset)
	}

	return reservation
}

// ReservationsFromDTO converts a slice of OpenAPI Reservation to domain Reservations.
func ReservationsFromDTO(dtos []openapi.TgvalidatordReservation) []*model.Reservation {
	if dtos == nil {
		return nil
	}
	reservations := make([]*model.Reservation, len(dtos))
	for i := range dtos {
		reservations[i] = ReservationFromDTO(&dtos[i])
	}
	return reservations
}

// ReservationUTXOFromDTO converts an OpenAPI UTXO to a domain ReservationUTXO.
func ReservationUTXOFromDTO(dto *openapi.TgvalidatordUTXO) *model.ReservationUTXO {
	if dto == nil {
		return nil
	}

	return &model.ReservationUTXO{
		ID:            safeString(dto.Id),
		Hash:          safeString(dto.Hash),
		OutputIndex:   safeInt64(dto.OutputIndex),
		Script:        safeString(dto.Script),
		Value:         safeString(dto.Value),
		ValueString:   safeString(dto.ValueString),
		BlockHeight:   safeString(dto.BlockHeight),
		ReservationID: safeString(dto.ReservationId),
	}
}
