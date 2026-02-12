package mapper

import (
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WebhookCallFromDTO converts an OpenAPI WebhookCall to a domain WebhookCall.
func WebhookCallFromDTO(dto *openapi.TgvalidatordWebhookCall) *model.WebhookCall {
	if dto == nil {
		return nil
	}

	call := &model.WebhookCall{
		ID:            safeString(dto.Id),
		EventID:       safeString(dto.EventId),
		WebhookID:     safeString(dto.WebhookId),
		Payload:       safeString(dto.Payload),
		Status:        safeString(dto.Status),
		StatusMessage: safeString(dto.StatusMessage),
	}

	// Parse attempts count (comes as string from API)
	if dto.Attempts != nil {
		if attempts, err := strconv.ParseInt(*dto.Attempts, 10, 64); err == nil {
			call.Attempts = attempts
		}
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		call.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		call.UpdatedAt = *dto.UpdatedAt
	}

	return call
}

// WebhookCallsFromDTO converts a slice of OpenAPI WebhookCall to domain WebhookCalls.
func WebhookCallsFromDTO(dtos []openapi.TgvalidatordWebhookCall) []*model.WebhookCall {
	if dtos == nil {
		return nil
	}
	calls := make([]*model.WebhookCall, len(dtos))
	for i := range dtos {
		calls[i] = WebhookCallFromDTO(&dtos[i])
	}
	return calls
}
