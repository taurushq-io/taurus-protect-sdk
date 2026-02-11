package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WebhookFromDTO converts an OpenAPI Webhook to a domain Webhook.
func WebhookFromDTO(dto *openapi.TgvalidatordWebhook) *model.Webhook {
	if dto == nil {
		return nil
	}

	webhook := &model.Webhook{
		ID:     safeString(dto.Id),
		Type:   safeString(dto.Type),
		URL:    safeString(dto.Url),
		Status: safeString(dto.Status),
	}

	// Convert timestamps
	if dto.TimeoutUntil != nil {
		webhook.TimeoutUntil = *dto.TimeoutUntil
	}
	if dto.CreatedAt != nil {
		webhook.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		webhook.UpdatedAt = *dto.UpdatedAt
	}

	return webhook
}

// WebhooksFromDTO converts a slice of OpenAPI Webhooks to domain Webhooks.
func WebhooksFromDTO(dtos []openapi.TgvalidatordWebhook) []*model.Webhook {
	if dtos == nil {
		return nil
	}
	webhooks := make([]*model.Webhook, len(dtos))
	for i := range dtos {
		webhooks[i] = WebhookFromDTO(&dtos[i])
	}
	return webhooks
}
