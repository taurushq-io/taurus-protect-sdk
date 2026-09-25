package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ReservationService provides fund reservation management operations.
type ReservationService struct {
	api       *openapi.ReservationsAPIService
	errMapper *ErrorMapper
}

// NewReservationService creates a new ReservationService.
func NewReservationService(client *openapi.APIClient) *ReservationService {
	return &ReservationService{
		api:       client.ReservationsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetReservation retrieves a single reservation by ID.
func (s *ReservationService) GetReservation(ctx context.Context, id string) (*model.Reservation, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.WalletServiceGetReservation(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("reservation not found")
	}

	return mapper.ReservationFromDTO(resp.Result), nil
}

// GetReservationUTXO retrieves the UTXO associated with a reservation.
func (s *ReservationService) GetReservationUTXO(ctx context.Context, id string) (*model.ReservationUTXO, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.WalletServiceGetReservationUTXO(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("UTXO not found for reservation")
	}

	return mapper.ReservationUTXOFromDTO(resp.Result), nil
}

// ListReservations retrieves a list of reservations with optional filtering and pagination.
func (s *ReservationService) ListReservations(ctx context.Context, opts *model.ListReservationsOptions) (*model.ListReservationsResult, error) {
	if opts == nil {
		opts = &model.ListReservationsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.WalletServiceGetReservations(ctx), window)
	if len(opts.Kinds) > 0 {
		req = req.Kinds(opts.Kinds)
	}
	if opts.Address != "" {
		req = req.Address(opts.Address)
	}
	if opts.AddressID != "" {
		req = req.AddressId(opts.AddressID)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListReservationsResult{
		Reservations: mapper.ReservationsFromDTO(resp.Result),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}
