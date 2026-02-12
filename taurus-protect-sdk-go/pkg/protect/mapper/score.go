package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// RefreshScoreResultFromDTO converts an OpenAPI TgvalidatordRefreshScoreReply to a domain RefreshScoreResult.
func RefreshScoreResultFromDTO(dto *openapi.TgvalidatordRefreshScoreReply) *model.RefreshScoreResult {
	if dto == nil {
		return nil
	}

	return &model.RefreshScoreResult{
		Scores: AddressScoresFromDTO(dto.Scores),
	}
}
