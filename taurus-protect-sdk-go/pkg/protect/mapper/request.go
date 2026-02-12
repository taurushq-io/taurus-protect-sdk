package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// RequestFromDTO converts an OpenAPI Request to a domain Request.
func RequestFromDTO(dto *openapi.TgvalidatordRequest) *model.Request {
	if dto == nil {
		return nil
	}

	request := &model.Request{
		ID:                safeString(dto.Id),
		TenantID:          safeString(dto.TenantId),
		Currency:          safeString(dto.Currency),
		Status:            safeString(dto.Status),
		Type:              safeString(dto.Type),
		Rule:              safeString(dto.Rule),
		RequestBundleID:   safeString(dto.RequestBundleId),
		ExternalRequestID: safeString(dto.ExternalRequestId),
	}

	// Convert timestamps
	if dto.CreationDate != nil {
		request.CreatedAt = *dto.CreationDate
	}
	if dto.UpdateDate != nil {
		request.UpdatedAt = *dto.UpdateDate
	}

	// Convert metadata
	if dto.Metadata != nil {
		request.Metadata = MetadataFromDTO(dto.Metadata)
	}

	// Convert attributes
	if dto.Attributes != nil {
		request.Attributes = make([]model.RequestAttribute, len(dto.Attributes))
		for i, attr := range dto.Attributes {
			request.Attributes[i] = RequestAttributeFromDTO(&attr)
		}
	}

	// Copy needs approval from
	if dto.NeedsApprovalFrom != nil {
		request.NeedsApprovalFrom = make([]string, len(dto.NeedsApprovalFrom))
		copy(request.NeedsApprovalFrom, dto.NeedsApprovalFrom)
	}

	// Convert signed requests
	request.SignedRequests = SignedRequestsFromDTO(dto.SignedRequests)

	return request
}

// RequestsFromDTO converts a slice of OpenAPI Request to domain Requests.
func RequestsFromDTO(dtos []openapi.TgvalidatordRequest) []*model.Request {
	if dtos == nil {
		return nil
	}
	requests := make([]*model.Request, len(dtos))
	for i := range dtos {
		requests[i] = RequestFromDTO(&dtos[i])
	}
	return requests
}

// MetadataFromDTO converts an OpenAPI Metadata to a domain RequestMetadata.
func MetadataFromDTO(dto *openapi.TgvalidatordMetadata) *model.RequestMetadata {
	if dto == nil {
		return nil
	}
	metadata := &model.RequestMetadata{
		Hash:            safeString(dto.Hash),
		PayloadAsString: safeString(dto.PayloadAsString),
		// SECURITY: Payload intentionally not mapped - use PayloadAsString only.
		// The raw payload object could be tampered with while payloadAsString
		// remains unchanged (hash still verifies).
	}
	return metadata
}

// SignedRequestFromDTO maps an OpenAPI RequestSignedRequest to a domain SignedRequest.
func SignedRequestFromDTO(dto openapi.RequestSignedRequest) model.SignedRequest {
	return model.SignedRequest{
		ID:               safeString(dto.Id),
		SignedRequest:    safeString(dto.SignedRequest),
		Status:           model.RequestStatusFromString(safeString(dto.Status)),
		Hash:             safeString(dto.Hash),
		Block:            safeInt64(dto.Block),
		Details:          safeString(dto.Details),
		CreationDate:     safeTime(dto.CreationDate),
		UpdateDate:       safeTime(dto.UpdateDate),
		BroadcastDate:    safeTime(dto.BroadcastDate),
		ConfirmationDate: safeTime(dto.ConfirmationDate),
	}
}

// SignedRequestsFromDTO maps a slice of OpenAPI RequestSignedRequest to domain SignedRequests.
func SignedRequestsFromDTO(dtos []openapi.RequestSignedRequest) []model.SignedRequest {
	if len(dtos) == 0 {
		return nil
	}
	result := make([]model.SignedRequest, len(dtos))
	for i, dto := range dtos {
		result[i] = SignedRequestFromDTO(dto)
	}
	return result
}

// RequestAttributeFromDTO converts an OpenAPI RequestAttribute to a domain RequestAttribute.
func RequestAttributeFromDTO(dto *openapi.TgvalidatordRequestAttribute) model.RequestAttribute {
	if dto == nil {
		return model.RequestAttribute{}
	}
	return model.RequestAttribute{
		ID:    safeString(dto.Id),
		Key:   safeString(dto.Key),
		Value: safeString(dto.Value),
	}
}
