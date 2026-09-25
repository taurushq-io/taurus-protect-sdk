package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ChangeService provides change management operations.
type ChangeService struct {
	api       *openapi.ChangesAPIService
	errMapper *ErrorMapper
}

// NewChangeService creates a new ChangeService.
func NewChangeService(client *openapi.APIClient) *ChangeService {
	return &ChangeService{
		api:       client.ChangesAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetChange retrieves a change by ID.
func (s *ChangeService) GetChange(ctx context.Context, changeID string) (*model.Change, error) {
	if changeID == "" {
		return nil, fmt.Errorf("changeID cannot be empty")
	}

	resp, httpResp, err := s.api.ChangeServiceGetChange(ctx, changeID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("change not found")
	}

	return mapper.ChangeFromDTO(resp.Result), nil
}

// ListChanges retrieves one page of changes. Continue with Page.NextCursor until
// Page.HasMore is false.
func (s *ChangeService) ListChanges(ctx context.Context, opts *model.ListChangesOptions) (*model.ListChangesResult, error) {
	if opts == nil {
		opts = &model.ListChangesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.ChangeServiceGetChanges(ctx), window)
	if opts.Entity != "" {
		req = req.Entity(opts.Entity)
	}
	if opts.EntityID != "" {
		req = req.EntityId(opts.EntityID)
	}
	if len(opts.EntityIDs) > 0 {
		req = req.EntityIDs(opts.EntityIDs)
	}
	if len(opts.EntityUUIDs) > 0 {
		req = req.EntityUUIDs(opts.EntityUUIDs)
	}
	if opts.Status != "" {
		req = req.Status(opts.Status)
	}
	if opts.CreatorID != "" {
		req = req.CreatorId(opts.CreatorID)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListChangesResult{
		Changes: mapper.ChangesFromDTO(resp.Result),
		Page:    page,
	}, nil
}

// ListChangesForApproval retrieves one page of the changes awaiting the current user's
// approval. Continue with Page.NextCursor until Page.HasMore is false.
func (s *ChangeService) ListChangesForApproval(ctx context.Context, opts *model.ListChangesForApprovalOptions) (*model.ListChangesResult, error) {
	if opts == nil {
		opts = &model.ListChangesForApprovalOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.ChangeServiceGetChangesForApproval(ctx), window)
	if len(opts.Entities) > 0 {
		req = req.Entities(opts.Entities)
	}
	if len(opts.EntityIDs) > 0 {
		req = req.EntityIDs(opts.EntityIDs)
	}
	if len(opts.EntityUUIDs) > 0 {
		req = req.EntityUUIDs(opts.EntityUUIDs)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListChangesResult{
		Changes: mapper.ChangesFromDTO(resp.Result),
		Page:    page,
	}, nil
}

// CreateChange creates a new change request.
func (s *ChangeService) CreateChange(ctx context.Context, req *model.CreateChangeRequest) (*model.CreateChangeResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Action == "" {
		return nil, fmt.Errorf("action is required")
	}
	if req.Entity == "" {
		return nil, fmt.Errorf("entity is required")
	}

	createReq := openapi.TgvalidatordCreateChangeRequest{
		Action: req.Action,
		Entity: req.Entity,
	}

	if req.EntityID != "" {
		createReq.EntityId = &req.EntityID
	}
	if req.EntityUUID != "" {
		createReq.EntityUUID = &req.EntityUUID
	}
	if req.Comment != "" {
		createReq.ChangeComment = &req.Comment
	}
	if len(req.Changes) > 0 {
		createReq.Changes = &req.Changes
	}

	resp, httpResp, err := s.api.ChangeServiceCreateChange(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.CreateChangeResultFromDTO(resp.Result), nil
}

// ApproveChange approves a single change by ID.
func (s *ChangeService) ApproveChange(ctx context.Context, changeID string) error {
	if changeID == "" {
		return fmt.Errorf("changeID cannot be empty")
	}

	// The API requires an empty body for single change approval
	_, httpResp, err := s.api.ChangeServiceApproveChange(ctx, changeID).
		Body(map[string]interface{}{}).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ApproveChanges approves multiple changes by their IDs.
func (s *ChangeService) ApproveChanges(ctx context.Context, changeIDs []string) error {
	if len(changeIDs) == 0 {
		return fmt.Errorf("changeIDs cannot be empty")
	}

	req := openapi.TgvalidatordApproveChangesRequest{
		Ids: changeIDs,
	}

	_, httpResp, err := s.api.ChangeServiceApproveChanges(ctx).
		Body(req).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// RejectChange rejects a single change by ID.
func (s *ChangeService) RejectChange(ctx context.Context, changeID string) error {
	if changeID == "" {
		return fmt.Errorf("changeID cannot be empty")
	}

	// The API requires an empty body for single change rejection
	_, httpResp, err := s.api.ChangeServiceRejectChange(ctx, changeID).
		Body(map[string]interface{}{}).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// RejectChanges rejects multiple changes by their IDs.
func (s *ChangeService) RejectChanges(ctx context.Context, changeIDs []string) error {
	if len(changeIDs) == 0 {
		return fmt.Errorf("changeIDs cannot be empty")
	}

	req := openapi.TgvalidatordRejectChangesRequest{
		Ids: changeIDs,
	}

	_, httpResp, err := s.api.ChangeServiceRejectChanges(ctx).
		Body(req).
		Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
