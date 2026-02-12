package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// MaxETHValidatorIDs is the maximum number of ETH validator IDs that can be requested at once.
const MaxETHValidatorIDs = 500

// StakingService provides staking-related operations.
type StakingService struct {
	api       *openapi.StakingAPIService
	errMapper *ErrorMapper
}

// NewStakingService creates a new StakingService.
func NewStakingService(client *openapi.APIClient) *StakingService {
	return &StakingService{
		api:       client.StakingAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetADAStakePoolInfo retrieves information about an ADA stake pool.
func (s *StakingService) GetADAStakePoolInfo(ctx context.Context, network, stakePoolID string) (*model.ADAStakePoolInfo, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if stakePoolID == "" {
		return nil, fmt.Errorf("stakePoolID cannot be empty")
	}

	resp, httpResp, err := s.api.StakingServiceGetADAStakePoolInfo(ctx, network, stakePoolID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ADAStakePoolInfoFromDTO(resp), nil
}

// GetETHValidatorsInfo retrieves information about ETH validators.
// Maximum 500 IDs can be provided.
func (s *StakingService) GetETHValidatorsInfo(ctx context.Context, network string, ids []string) ([]*model.ETHValidatorInfo, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if len(ids) > MaxETHValidatorIDs {
		return nil, fmt.Errorf("maximum %d validator IDs allowed, got %d", MaxETHValidatorIDs, len(ids))
	}

	req := s.api.StakingServiceGetETHValidatorsInfo(ctx, network)
	if len(ids) > 0 {
		req = req.Ids(ids)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ETHValidatorsInfoFromDTO(resp.Validators), nil
}

// GetFTMValidatorInfo retrieves information about an FTM validator.
func (s *StakingService) GetFTMValidatorInfo(ctx context.Context, network, validatorAddress string) (*model.FTMValidatorInfo, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if validatorAddress == "" {
		return nil, fmt.Errorf("validatorAddress cannot be empty")
	}

	resp, httpResp, err := s.api.StakingServiceGetFTMValidatorInfo(ctx, network, validatorAddress).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.FTMValidatorInfoFromDTO(resp), nil
}

// GetICPNeuronInfo retrieves information about an ICP neuron.
func (s *StakingService) GetICPNeuronInfo(ctx context.Context, network, neuronID string) (*model.ICPNeuronInfo, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if neuronID == "" {
		return nil, fmt.Errorf("neuronID cannot be empty")
	}

	resp, httpResp, err := s.api.StakingServiceGetICPNeuronInfo(ctx, network, neuronID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ICPNeuronInfoFromDTO(resp), nil
}

// GetNEARValidatorInfo retrieves information about a NEAR validator.
func (s *StakingService) GetNEARValidatorInfo(ctx context.Context, network, validatorAddress string) (*model.NEARValidatorInfo, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if validatorAddress == "" {
		return nil, fmt.Errorf("validatorAddress cannot be empty")
	}

	resp, httpResp, err := s.api.StakingServiceGetNEARValidatorInfo(ctx, network, validatorAddress).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.NEARValidatorInfoFromDTO(resp), nil
}

// ListStakeAccounts retrieves a list of stake accounts with optional filtering and pagination.
func (s *StakingService) ListStakeAccounts(ctx context.Context, opts *model.ListStakeAccountsOptions) (*model.ListStakeAccountsResult, error) {
	req := s.api.StakingServiceGetStakeAccounts(ctx)

	if opts != nil {
		if opts.AddressID != "" {
			req = req.AddressId(opts.AddressID)
		}
		if opts.AccountType != "" {
			req = req.AccountType(opts.AccountType)
		}
		if opts.AccountAddress != "" {
			req = req.AccountAddress(opts.AccountAddress)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListStakeAccountsResult{
		StakeAccounts: mapper.StakeAccountsFromDTO(resp.StakeAccounts),
	}

	// Parse cursor pagination info
	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.CurrentPage = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasPrevious != nil {
			result.HasPrevious = *resp.Cursor.HasPrevious
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// GetXTZAddressStakingRewards retrieves staking rewards for an XTZ address.
func (s *StakingService) GetXTZAddressStakingRewards(ctx context.Context, network, addressID string, opts *model.GetXTZStakingRewardsOptions) (*model.XTZStakingReward, error) {
	if network == "" {
		return nil, fmt.Errorf("network cannot be empty")
	}
	if addressID == "" {
		return nil, fmt.Errorf("addressID cannot be empty")
	}

	req := s.api.StakingServiceGetXTZAddressStakingRewards(ctx, network, addressID)

	if opts != nil {
		if opts.From != nil {
			req = req.From(*opts.From)
		}
		if opts.To != nil {
			req = req.To(*opts.To)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.XTZStakingRewardFromDTO(resp), nil
}
