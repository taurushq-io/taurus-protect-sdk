package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TagService provides tag management operations.
type TagService struct {
	api       *openapi.TagsAPIService
	errMapper *ErrorMapper
}

// NewTagService creates a new TagService.
func NewTagService(client *openapi.APIClient) *TagService {
	return &TagService{
		api:       client.TagsAPI,
		errMapper: NewErrorMapper(),
	}
}

// CreateTag creates a new tag and returns its ID.
func (s *TagService) CreateTag(ctx context.Context, req *model.CreateTagRequest) (string, error) {
	if req == nil {
		return "", fmt.Errorf("request cannot be nil")
	}
	if req.Value == "" {
		return "", fmt.Errorf("value is required")
	}

	createReq := openapi.TgvalidatordCreateTagRequest{
		Value: req.Value,
		Color: req.Color,
	}

	resp, httpResp, err := s.api.TagServiceCreateTag(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return "", s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil || resp.Result.Id == nil {
		return "", fmt.Errorf("failed to create tag: no ID returned")
	}

	return *resp.Result.Id, nil
}

// DeleteTag deletes a tag by ID.
func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	_, httpResp, err := s.api.TagServiceDeleteTag(ctx, id).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ListTags retrieves a list of tags with optional filtering.
func (s *TagService) ListTags(ctx context.Context, opts *model.ListTagsOptions) ([]*model.Tag, error) {
	req := s.api.TagServiceGetTags(ctx)

	if opts != nil {
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.TagsFromDTO(resp.Result), nil
}
