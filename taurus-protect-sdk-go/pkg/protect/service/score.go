package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ScoreService provides score management operations.
type ScoreService struct {
	api       *openapi.ScoresAPIService
	errMapper *ErrorMapper
}

// NewScoreService creates a new ScoreService.
func NewScoreService(client *openapi.APIClient) *ScoreService {
	return &ScoreService{
		api:       client.ScoresAPI,
		errMapper: NewErrorMapper(),
	}
}

// RefreshAddressScore refreshes the risk score for an address using the specified provider.
func (s *ScoreService) RefreshAddressScore(ctx context.Context, addressID string, provider string) (*model.RefreshScoreResult, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	body := openapi.ScoreServiceRefreshAddressScoreBody{}
	if provider != "" {
		body.SetScoreProvider(provider)
	}

	resp, httpResp, err := s.api.ScoreServiceRefreshAddressScore(ctx, addressID).
		Body(body).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &model.RefreshScoreResult{
		Scores: mapper.AddressScoresFromDTO(resp.Scores),
	}, nil
}

// RefreshWLAScore refreshes the risk score for a whitelisted address using the specified provider.
func (s *ScoreService) RefreshWLAScore(ctx context.Context, addressID string, provider string) (*model.RefreshScoreResult, error) {
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	body := openapi.ScoreServiceRefreshWLAScoreBody{}
	if provider != "" {
		body.SetScoreProvider(provider)
	}

	resp, httpResp, err := s.api.ScoreServiceRefreshWLAScore(ctx, addressID).
		Body(body).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &model.RefreshScoreResult{
		Scores: mapper.AddressScoresFromDTO(resp.Scores),
	}, nil
}
